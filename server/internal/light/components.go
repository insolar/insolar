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

package light

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
	"github.com/insolar/insolar/insolar/message"
	"github.com/insolar/insolar/keystore"
	"github.com/insolar/insolar/ledger"
	"github.com/insolar/insolar/logicrunner/artifacts"
	"github.com/insolar/insolar/messagebus"
	"github.com/insolar/insolar/metrics"
	"github.com/insolar/insolar/network/nodenetwork"
	"github.com/insolar/insolar/network/servicenetwork"
	"github.com/insolar/insolar/network/state"
	"github.com/insolar/insolar/network/termination"
	"github.com/insolar/insolar/networkcoordinator"
	"github.com/insolar/insolar/platformpolicy"
)

type components struct {
	cmp               component.Manager
	NodeRef, NodeRole string
}

func newComponents(ctx context.Context, cfg configuration.Configuration) *components {
	// Cryptography.
	var (
		KeyProcessor  insolar.KeyProcessor
		CryptoScheme  insolar.PlatformCryptographyScheme
		CryptoService insolar.CryptographyService
		CertManager   insolar.CertificateManager
	)
	{
		var err error
		// Private key storage.
		ks, err := keystore.NewKeyStore(cfg.KeysPath)
		checkError(ctx, err, "failed to load KeyStore")
		// Public key manipulations.
		KeyProcessor = platformpolicy.NewKeyProcessor()
		// Platform cryptography.
		CryptoScheme = platformpolicy.NewPlatformCryptographyScheme()
		// Sign, verify, etc.
		CryptoService = cryptography.NewCryptographyService()

		c := component.Manager{}
		c.Inject(CryptoService, CryptoScheme, KeyProcessor, ks)

		publicKey, err := CryptoService.GetPublicKey()
		checkError(ctx, err, "failed to retrieve node public key")

		// Node certificate.
		CertManager, err = certificate.NewManagerReadCertificate(publicKey, KeyProcessor, cfg.CertificatePath)
		checkError(ctx, err, "failed to start Certificate")
	}

	c := &components{}
	c.cmp = component.Manager{}
	c.NodeRef = CertManager.GetCertificate().GetNodeRef().String()
	c.NodeRole = CertManager.GetCertificate().GetRole().String()

	// Network.
	var (
		NetworkService     insolar.Network
		NetworkCoordinator insolar.NetworkCoordinator
		NodeNetwork        insolar.NodeNetwork
		NetworkSwitcher    insolar.NetworkSwitcher
		Termination        insolar.TerminationHandler
	)
	{
		var err error
		// External communication.
		NetworkService, err = servicenetwork.NewServiceNetwork(cfg, &c.cmp, false)
		checkError(ctx, err, "failed to start Network")

		Termination = termination.NewHandler(NetworkService)

		// Node info.
		NodeNetwork, err = nodenetwork.NewNodeNetwork(cfg.Host, CertManager.GetCertificate())
		checkError(ctx, err, "failed to start NodeNetwork")

		NetworkSwitcher, err = state.NewNetworkSwitcher()
		checkError(ctx, err, "failed to start NetworkSwitcher")

		NetworkCoordinator, err = networkcoordinator.New()
		checkError(ctx, err, "failed to start NetworkCoordinator")
	}

	// API.
	var (
		Requester insolar.ContractRequester
		Genesis   insolar.GenesisDataProvider
		API       insolar.APIRunner
	)
	{
		var err error
		Requester, err = contractrequester.New()
		checkError(ctx, err, "failed to start ContractRequester")

		Genesis, err = genesisdataprovider.New()
		checkError(ctx, err, "failed to start GenesisDataProvider")

		API, err = api.NewRunner(&cfg.APIRunner)
		checkError(ctx, err, "failed to start ApiRunner")
	}

	// Communication.
	var (
		Tokens  insolar.DelegationTokenFactory
		Parcels message.ParcelFactory
		Bus     insolar.MessageBus
	)
	{
		var err error
		Tokens = delegationtoken.NewDelegationTokenFactory()
		Parcels = messagebus.NewParcelFactory()
		Bus, err = messagebus.NewMessageBus(cfg)
		checkError(ctx, err, "failed to start MessageBus")
	}

	metricsHandler, err := metrics.NewMetrics(
		ctx,
		cfg.Metrics,
		metrics.GetInsolarRegistry(c.NodeRole),
		c.NodeRole,
	)
	checkError(ctx, err, "failed to start Metrics")

	c.cmp.Register(
		Termination,
		CryptoScheme,
		CryptoService,
		KeyProcessor,
		CertManager,
		NodeNetwork,
		NetworkService,
	)

	components := ledger.GetLedgerComponents(cfg.Ledger, CertManager.GetCertificate())

	c.cmp.Inject(
		append(components, []interface{}{
			Bus,
			Requester,
			Tokens,
			Parcels,
			artifacts.NewClient(),
			Genesis,
			API,
			metricsHandler,
			NetworkSwitcher,
			NetworkCoordinator,
			CryptoService,
			KeyProcessor,
		}...)...,
	)

	err = c.cmp.Init(ctx)
	checkError(ctx, err, "failed to init components")

	return c
}

func (c *components) Start(ctx context.Context) error {
	return c.cmp.Start(ctx)
}

func (c *components) Stop(ctx context.Context) error {
	return c.cmp.Stop(ctx)
}
