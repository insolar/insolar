//
// Copyright 2019 Insolar Technologies GmbH
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package virtual

import (
	"context"

	"github.com/insolar/insolar/api"
	"github.com/insolar/insolar/certificate"
	"github.com/insolar/insolar/component"
	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/contractrequester"
	"github.com/insolar/insolar/cryptography"
	"github.com/insolar/insolar/genesisdataprovider"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/delegationtoken"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/jetcoordinator"
	"github.com/insolar/insolar/insolar/node"
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/keystore"
	"github.com/insolar/insolar/logicrunner"
	"github.com/insolar/insolar/logicrunner/artifacts"
	"github.com/insolar/insolar/logicrunner/pulsemanager"
	"github.com/insolar/insolar/messagebus"
	"github.com/insolar/insolar/metrics"
	"github.com/insolar/insolar/network/nodenetwork"
	"github.com/insolar/insolar/network/servicenetwork"
	"github.com/insolar/insolar/network/termination"
	"github.com/insolar/insolar/networkcoordinator"
	"github.com/insolar/insolar/platformpolicy"
	"github.com/insolar/insolar/pulsar"
	"github.com/insolar/insolar/pulsar/entropygenerator"
	"github.com/insolar/insolar/version/manager"
)

type bootstrapComponents struct {
	CryptographyService        insolar.CryptographyService
	PlatformCryptographyScheme insolar.PlatformCryptographyScheme
	KeyStore                   insolar.KeyStore
	KeyProcessor               insolar.KeyProcessor
}

func initBootstrapComponents(ctx context.Context, cfg configuration.Configuration) bootstrapComponents {
	earlyComponents := component.Manager{}

	keyStore, err := keystore.NewKeyStore(cfg.KeysPath)
	checkError(ctx, err, "failed to load KeyStore: ")

	platformCryptographyScheme := platformpolicy.NewPlatformCryptographyScheme()
	keyProcessor := platformpolicy.NewKeyProcessor()

	cryptographyService := cryptography.NewCryptographyService()
	earlyComponents.Register(platformCryptographyScheme, keyStore)
	earlyComponents.Inject(cryptographyService, keyProcessor)

	return bootstrapComponents{
		CryptographyService:        cryptographyService,
		PlatformCryptographyScheme: platformCryptographyScheme,
		KeyStore:                   keyStore,
		KeyProcessor:               keyProcessor,
	}
}

func initCertificateManager(
	ctx context.Context,
	cfg configuration.Configuration,
	isBootstrap bool,
	cryptographyService insolar.CryptographyService,
	keyProcessor insolar.KeyProcessor,
) *certificate.CertificateManager {
	var certManager *certificate.CertificateManager
	var err error

	publicKey, err := cryptographyService.GetPublicKey()
	checkError(ctx, err, "failed to retrieve node public key")

	if isBootstrap {
		certManager, err = certificate.NewManagerCertificateWithKeys(publicKey, keyProcessor)
		checkError(ctx, err, "failed to start Certificate (bootstrap mode)")
	} else {
		certManager, err = certificate.NewManagerReadCertificate(publicKey, keyProcessor, cfg.CertificatePath)
		checkError(ctx, err, "failed to start Certificate")
	}

	return certManager
}

// initComponents creates and links all insolard components
func initComponents(
	ctx context.Context,
	cfg configuration.Configuration,
	cryptographyService insolar.CryptographyService,
	platformCryptographyScheme insolar.PlatformCryptographyScheme,
	keyStore insolar.KeyStore,
	keyProcessor insolar.KeyProcessor,
	certManager insolar.CertificateManager,
	isGenesis bool,

) (*component.Manager, insolar.TerminationHandler, error) {
	cm := component.Manager{}

	nodeNetwork, err := nodenetwork.NewNodeNetwork(cfg.Host.Transport, certManager.GetCertificate())
	checkError(ctx, err, "failed to start NodeNetwork")

	logicRunner, err := logicrunner.NewLogicRunner(&cfg.LogicRunner)
	checkError(ctx, err, "failed to start LogicRunner")

	nw, err := servicenetwork.NewServiceNetwork(cfg, &cm, isGenesis)
	checkError(ctx, err, "failed to start Network")

	terminationHandler := termination.NewHandler(nw)

	delegationTokenFactory := delegationtoken.NewDelegationTokenFactory()
	parcelFactory := messagebus.NewParcelFactory()

	messageBus, err := messagebus.NewMessageBus(cfg)
	checkError(ctx, err, "failed to start MessageBus")

	contractRequester, err := contractrequester.New()
	checkError(ctx, err, "failed to start ContractRequester")

	genesisDataProvider, err := genesisdataprovider.New()
	checkError(ctx, err, "failed to start GenesisDataProvider")

	apiRunner, err := api.NewRunner(&cfg.APIRunner)
	checkError(ctx, err, "failed to start ApiRunner")

	metricsHandler, err := metrics.NewMetrics(ctx, cfg.Metrics, metrics.GetInsolarRegistry("virtual"), "virtual")
	checkError(ctx, err, "failed to start Metrics")

	networkCoordinator, err := networkcoordinator.New()
	checkError(ctx, err, "failed to start NetworkCoordinator")

	_, err = manager.NewVersionManager(cfg.VersionManager)
	checkError(ctx, err, "failed to load VersionManager: ")

	// move to logic runner ??
	err = logicRunner.OnPulse(ctx, *pulsar.NewPulse(cfg.Pulsar.NumberDelta, 0, &entropygenerator.StandardEntropyGenerator{}))
	checkError(ctx, err, "failed init pulse for LogicRunner")

	cm.Register(
		terminationHandler,
		platformCryptographyScheme,
		keyStore,
		cryptographyService,
		keyProcessor,
		certManager,
		nodeNetwork,
		nw,
		pulsemanager.NewPulseManager(),
	)

	components := []interface{}{
		messageBus,
		contractRequester,
		logicRunner,
		artifacts.NewClient(),
		pulse.NewStorageMem(),
		jet.NewStore(),
		jetcoordinator.NewJetCoordinator(cfg.Ledger.LightChainLimit),
		node.NewStorage(),
		delegationTokenFactory,
		parcelFactory,
	}
	components = append(components, []interface{}{
		genesisDataProvider,
		apiRunner,
		metricsHandler,
		networkCoordinator,
		cryptographyService,
		keyProcessor,
	}...)

	cm.Inject(components...)

	return &cm, terminationHandler, nil
}
