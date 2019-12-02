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
	"io"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"

	"github.com/insolar/component-manager"

	"github.com/insolar/insolar/application/api"
	"github.com/insolar/insolar/certificate"
	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/contractrequester"
	"github.com/insolar/insolar/cryptography"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/jetcoordinator"
	"github.com/insolar/insolar/insolar/node"
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/keystore"
	"github.com/insolar/insolar/log/logwatermill"
	"github.com/insolar/insolar/logicrunner"
	"github.com/insolar/insolar/logicrunner/artifacts"
	"github.com/insolar/insolar/logicrunner/logicexecutor"
	"github.com/insolar/insolar/logicrunner/machinesmanager"
	"github.com/insolar/insolar/logicrunner/pulsemanager"
	"github.com/insolar/insolar/metrics"
	"github.com/insolar/insolar/network/servicenetwork"
	"github.com/insolar/insolar/platformpolicy"
	"github.com/insolar/insolar/server/internal"
)

type bootstrapComponents struct {
	CryptographyService        insolar.CryptographyService
	PlatformCryptographyScheme insolar.PlatformCryptographyScheme
	KeyStore                   insolar.KeyStore
	KeyProcessor               insolar.KeyProcessor
}

func initBootstrapComponents(ctx context.Context, cfg configuration.Configuration) bootstrapComponents {
	earlyComponents := component.NewManager(nil)

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
	cryptographyService insolar.CryptographyService,
	keyProcessor insolar.KeyProcessor,
) *certificate.CertificateManager {
	var certManager *certificate.CertificateManager
	var err error

	publicKey, err := cryptographyService.GetPublicKey()
	checkError(ctx, err, "failed to retrieve node public key")

	certManager, err = certificate.NewManagerReadCertificate(publicKey, keyProcessor, cfg.CertificatePath)
	checkError(ctx, err, "failed to start Certificate")

	return certManager
}

// initComponents creates and links all insolard components
func initComponents(
	ctx context.Context,
	cfg configuration.Configuration,
	cryptographyService insolar.CryptographyService,
	pcs insolar.PlatformCryptographyScheme,
	keyStore insolar.KeyStore,
	keyProcessor insolar.KeyProcessor,
	certManager insolar.CertificateManager,

) (*component.Manager, func()) {
	cm := component.NewManager(nil)

	// Watermill.
	var (
		wmLogger   *logwatermill.WatermillLogAdapter
		publisher  message.Publisher
		subscriber message.Subscriber
	)
	{
		wmLogger = logwatermill.NewWatermillLogAdapter(inslogger.FromContext(ctx))
		pubsub := gochannel.NewGoChannel(gochannel.Config{}, wmLogger)
		subscriber = pubsub
		publisher = pubsub
		// Wrapped watermill Publisher for introspection.
		publisher = internal.PublisherWrapper(ctx, cm, cfg.Introspection, publisher)
	}

	nw, err := servicenetwork.NewServiceNetwork(cfg, cm)
	checkError(ctx, err, "failed to start Network")

	metricsComp := metrics.NewMetrics(cfg.Metrics, metrics.GetInsolarRegistry("virtual"), "virtual")

	jc := jetcoordinator.NewJetCoordinator(cfg.Ledger.LightChainLimit, *certManager.GetCertificate().GetNodeRef())
	pulses := pulse.NewStorageMem()

	b := bus.NewBus(cfg.Bus, publisher, pulses, jc, pcs)
	artifactsClient := artifacts.NewClient(b)
	cachedPulses := artifacts.NewPulseAccessorLRU(pulses, artifactsClient, cfg.LogicRunner.PulseLRUSize)

	logicRunner, err := logicrunner.NewLogicRunner(&cfg.LogicRunner, publisher, b)
	checkError(ctx, err, "failed to start LogicRunner")

	contractRequester, err := contractrequester.New(
		b,
		pulses,
		jc,
		pcs,
	)
	checkError(ctx, err, "failed to start ContractRequester")

	availabilityChecker := api.NewNetworkChecker(cfg.AvailabilityChecker)

	API, err := api.NewRunner(
		&cfg.APIRunner,
		certManager,
		contractRequester,
		nw,
		nw,
		pulses,
		artifactsClient,
		jc,
		nw,
		availabilityChecker,
	)
	checkError(ctx, err, "failed to start ApiRunner")

	AdminAPIRunner, err := api.NewRunner(
		&cfg.AdminAPIRunner,
		certManager,
		contractRequester,
		nw,
		nw,
		pulses,
		artifactsClient,
		jc,
		nw,
		availabilityChecker,
	)
	checkError(ctx, err, "failed to start AdminAPIRunner")

	APIWrapper := api.NewWrapper(API, AdminAPIRunner)

	// TODO: remove this hack in INS-3341
	contractRequester.LR = logicRunner

	pm := pulsemanager.NewPulseManager()

	cm.Register(
		pcs,
		keyStore,
		cryptographyService,
		keyProcessor,
		certManager,
		logicRunner,
		logicexecutor.NewLogicExecutor(cachedPulses),
		logicrunner.NewRequestsExecutor(),
		machinesmanager.NewMachinesManager(),
		APIWrapper,
		availabilityChecker,
		nw,
		pm,
	)

	components := []interface{}{
		b,
		publisher,
		contractRequester,
		artifactsClient,
		artifacts.NewDescriptorsCache(),
		jc,
		pulses,
		jet.NewStore(),
		node.NewStorage(),
	}
	components = append(components, []interface{}{
		metricsComp,
		cryptographyService,
		keyProcessor,
	}...)

	cm.Inject(components...)

	err = cm.Init(ctx)
	checkError(ctx, err, "failed to init components")

	// this should be done after Init due to inject
	pm.AddDispatcher(logicRunner.FlowDispatcher, contractRequester.FlowDispatcher)

	return cm, startWatermill(
		ctx, wmLogger, subscriber, b,
		nw.SendMessageHandler,
		logicRunner.FlowDispatcher.Process,
		contractRequester.FlowDispatcher.Process,
	)
}

func startWatermill(
	ctx context.Context,
	logger watermill.LoggerAdapter,
	sub message.Subscriber,
	b *bus.Bus,
	outHandler, inHandler, resultsHandler message.NoPublishHandlerFunc,
) func() {
	inRouter, err := message.NewRouter(message.RouterConfig{}, logger)
	if err != nil {
		panic(err)
	}
	outRouter, err := message.NewRouter(message.RouterConfig{}, logger)
	if err != nil {
		panic(err)
	}

	outRouter.AddNoPublisherHandler(
		"OutgoingHandler",
		bus.TopicOutgoing,
		sub,
		outHandler,
	)

	inRouter.AddMiddleware(
		b.IncomingMessageRouter,
	)

	inRouter.AddNoPublisherHandler(
		"IncomingHandler",
		bus.TopicIncoming,
		sub,
		inHandler,
	)

	inRouter.AddNoPublisherHandler(
		"IncomingRequestResultHandler",
		bus.TopicIncomingRequestResults,
		sub,
		resultsHandler)

	startRouter(ctx, inRouter)
	startRouter(ctx, outRouter)

	return stopWatermill(ctx, inRouter, outRouter)
}

func stopWatermill(ctx context.Context, routers ...io.Closer) func() {
	return func() {
		for _, r := range routers {
			err := r.Close()
			if err != nil {
				inslogger.FromContext(ctx).Error("Error while closing router", err)
			}
		}
	}
}

func startRouter(ctx context.Context, router *message.Router) {
	go func() {
		if err := router.Run(ctx); err != nil {
			inslogger.FromContext(ctx).Error("Error while running router", err)
		}
	}()
	<-router.Running()
}
