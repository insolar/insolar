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

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	"github.com/pkg/errors"

	"github.com/insolar/component-manager"

	"github.com/insolar/insolar/application/api"
	"github.com/insolar/insolar/certificate"
	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/contractrequester"
	"github.com/insolar/insolar/cryptography"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/insolar/flow/dispatcher"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/jetcoordinator"
	"github.com/insolar/insolar/insolar/node"
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/keystore"
	"github.com/insolar/insolar/ledger/drop"
	"github.com/insolar/insolar/ledger/light/executor"
	"github.com/insolar/insolar/ledger/light/handle"
	"github.com/insolar/insolar/ledger/light/proc"
	"github.com/insolar/insolar/ledger/object"
	"github.com/insolar/insolar/log/logwatermill"
	"github.com/insolar/insolar/logicrunner/artifacts"
	"github.com/insolar/insolar/metrics"
	"github.com/insolar/insolar/network/servicenetwork"
	"github.com/insolar/insolar/platformpolicy"
	"github.com/insolar/insolar/server/internal"
)

type components struct {
	cmp               *component.Manager
	NodeRef, NodeRole string
	replicator        executor.LightReplicator
	cleaner           executor.Cleaner
}

func initTemporaryCertificateManager(ctx context.Context, cfg *configuration.Configuration) (*certificate.CertificateManager, error) {
	earlyComponents := component.NewManager(nil)

	keyStore, err := keystore.NewKeyStore(cfg.KeysPath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load KeyStore")
	}

	platformCryptographyScheme := platformpolicy.NewPlatformCryptographyScheme()
	keyProcessor := platformpolicy.NewKeyProcessor()

	cryptographyService := cryptography.NewCryptographyService()
	earlyComponents.Register(platformCryptographyScheme, keyStore)
	earlyComponents.Inject(cryptographyService, keyProcessor)

	publicKey, err := cryptographyService.GetPublicKey()
	if err != nil {
		return nil, errors.Wrap(err, "failed to retrieve node public key")
	}

	certManager, err := certificate.NewManagerReadCertificate(publicKey, keyProcessor, cfg.CertificatePath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create new CertificateManager")
	}

	return certManager, nil
}

func newComponents(ctx context.Context, cfg configuration.Configuration) (*components, error) {
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
		if err != nil {
			return nil, errors.Wrap(err, "failed to load KeyStore")
		}
		// Public key manipulations.
		KeyProcessor = platformpolicy.NewKeyProcessor()
		// Platform cryptography.
		CryptoScheme = platformpolicy.NewPlatformCryptographyScheme()
		// Sign, verify, etc.
		CryptoService = cryptography.NewCryptographyService()

		c := component.NewManager(nil)
		c.Inject(CryptoService, CryptoScheme, KeyProcessor, ks)

		publicKey, err := CryptoService.GetPublicKey()
		if err != nil {
			return nil, errors.Wrap(err, "failed to retrieve node public key")
		}

		// Node certificate.
		CertManager, err = certificate.NewManagerReadCertificate(publicKey, KeyProcessor, cfg.CertificatePath)
		if err != nil {
			return nil, errors.Wrap(err, "failed to start Certificate")
		}
	}

	comps := &components{}
	comps.cmp = component.NewManager(nil)
	comps.NodeRef = CertManager.GetCertificate().GetNodeRef().String()
	comps.NodeRole = CertManager.GetCertificate().GetRole().String()

	// Watermill stuff.
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
		// Wrapped watermill publisher for introspection.
		publisher = internal.PublisherWrapper(ctx, comps.cmp, cfg.Introspection, publisher)
	}

	// Network.
	var (
		NetworkService *servicenetwork.ServiceNetwork
	)
	{
		var err error
		// External communication.
		NetworkService, err = servicenetwork.NewServiceNetwork(cfg, comps.cmp)
		if err != nil {
			return nil, errors.Wrap(err, "failed to start Network")
		}
	}

	// Role calculations.
	var (
		Coordinator jet.Coordinator
		Pulses      *pulse.StorageMem
		Jets        jet.Storage
		Nodes       *node.Storage
	)
	{
		Nodes = node.NewStorage()
		Pulses = pulse.NewStorageMem()
		Jets = jet.NewStore()

		c := jetcoordinator.NewJetCoordinator(cfg.Ledger.LightChainLimit, *CertManager.GetCertificate().GetNodeRef())
		c.PulseCalculator = Pulses
		c.PulseAccessor = Pulses
		c.JetAccessor = Jets
		c.PlatformCryptographyScheme = CryptoScheme
		c.Nodes = Nodes

		Coordinator = c
	}

	// Communication.
	var (
		Sender *bus.Bus
	)
	{
		Sender = bus.NewBus(cfg.Bus, publisher, Pulses, Coordinator, CryptoScheme)
	}

	// API.
	var (
		Requester           *contractrequester.ContractRequester
		ArtifactsClient     = artifacts.NewClient(Sender)
		AvailabilityChecker = api.NewNetworkChecker(cfg.AvailabilityChecker)
		APIWrapper          *api.RunnerWrapper
	)
	{
		var err error
		Requester, err = contractrequester.New(
			Sender,
			Pulses,
			Coordinator,
			CryptoScheme,
		)
		if err != nil {
			return nil, errors.Wrap(err, "failed to start ContractRequester")
		}

		API, err := api.NewRunner(
			&cfg.APIRunner,
			CertManager,
			Requester,
			NetworkService,
			NetworkService,
			Pulses,
			ArtifactsClient,
			Coordinator,
			NetworkService,
			AvailabilityChecker,
		)
		if err != nil {
			return nil, errors.Wrap(err, "failed to start ApiRunner")
		}

		AdminAPIRunner, err := api.NewRunner(
			&cfg.AdminAPIRunner,
			CertManager,
			Requester,
			NetworkService,
			NetworkService,
			Pulses,
			ArtifactsClient,
			Coordinator,
			NetworkService,
			AvailabilityChecker,
		)
		if err != nil {
			return nil, errors.Wrap(err, "failed to start AdminAPIRunner")
		}

		APIWrapper = api.NewWrapper(API, AdminAPIRunner)
	}

	metricsComp := metrics.NewMetrics(
		cfg.Metrics,
		metrics.GetInsolarRegistry(comps.NodeRole),
		comps.NodeRole,
	)

	metricsRegistry := executor.NewMetricsRegistry()

	// Light components.
	var (
		PulseManager   *executor.PulseManager
		FlowDispatcher dispatcher.Dispatcher
	)
	{
		conf := cfg.Ledger
		idLocker := object.NewIndexLocker()
		drops := drop.NewStorageMemory()
		records := object.NewRecordMemory()
		indexes := object.NewIndexStorageMemory()
		writeController := executor.NewWriteController()
		hotWaitReleaser := executor.NewChannelWaiter()

		c := component.NewManager(nil)
		c.Inject(CryptoScheme)

		jetFetcher := executor.NewFetcher(Nodes, Jets, Sender, Coordinator)
		filamentCalculator := executor.NewFilamentCalculator(
			indexes,
			records,
			Coordinator,
			jetFetcher,
			Sender,
			Pulses,
		)
		requestChecker := executor.NewRequestChecker(
			filamentCalculator,
			Coordinator,
			jetFetcher,
			CryptoScheme,
			Sender,
		)
		detachedNotifier := executor.NewDetachedNotifierDefault(Sender)

		jetCalculator := executor.NewJetCalculator(Coordinator, Jets)
		lightCleaner := executor.NewCleaner(
			Jets.(jet.Cleaner),
			Nodes,
			drops,
			records,
			indexes,
			Pulses,
			Pulses,
			indexes,
			filamentCalculator,
			conf.LightChainLimit,
			conf.CleanerDelay,
			conf.FilamentCacheLimit,
		)
		comps.cleaner = lightCleaner

		lthSyncer := executor.NewReplicatorDefault(
			jetCalculator,
			lightCleaner,
			Sender,
			Pulses,
			drops,
			records,
			indexes,
			Jets,
		)
		comps.replicator = lthSyncer

		jetSplitter := executor.NewJetSplitter(
			conf.JetSplit, jetCalculator, Jets, Jets, drops, drops, Pulses, records,
		)

		hotSender := executor.NewHotSender(
			drops,
			indexes,
			Pulses,
			Jets,
			conf.LightChainLimit,
			Sender,
		)

		stateIniter := executor.NewStateIniter(
			Jets, hotWaitReleaser, drops, Nodes, Sender, Pulses, Pulses, jetCalculator, indexes,
		)

		dep := proc.NewDependencies(
			CryptoScheme,
			Coordinator,
			Jets,
			Pulses,
			Sender,
			drops,
			idLocker,
			records,
			indexes,
			hotWaitReleaser,
			hotWaitReleaser,
			writeController,
			jetFetcher,
			filamentCalculator,
			requestChecker,
			detachedNotifier,
			conf,
			metricsRegistry,
		)

		initHandle := func(msg *message.Message) *handle.Init {
			return handle.NewInit(dep, Sender, msg)
		}

		FlowDispatcher = dispatcher.NewDispatcher(
			Pulses,
			func(msg *message.Message) flow.Handle {
				return initHandle(msg).Present
			}, func(msg *message.Message) flow.Handle {
				return initHandle(msg).Future
			}, func(msg *message.Message) flow.Handle {
				return initHandle(msg).Past
			},
		)

		PulseManager = executor.NewPulseManager(
			NetworkService,
			[]dispatcher.Dispatcher{FlowDispatcher, Requester.FlowDispatcher},
			Nodes,
			Pulses,
			Pulses,
			hotWaitReleaser,
			jetSplitter,
			lthSyncer,
			hotSender,
			writeController,
			stateIniter,
			hotWaitReleaser,
			metricsRegistry,
		)
	}

	comps.cmp.Inject(
		Sender,
		Jets,
		Pulses,
		Coordinator,
		PulseManager,
		metricsComp,
		Requester,
		ArtifactsClient,
		APIWrapper,
		AvailabilityChecker,
		KeyProcessor,
		CryptoScheme,
		CryptoService,
		CertManager,
		NetworkService,
		publisher,
	)

	err := comps.cmp.Init(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to init components")
	}

	comps.startWatermill(ctx, wmLogger, subscriber, Sender, NetworkService.SendMessageHandler, FlowDispatcher.Process, Requester.FlowDispatcher.Process)

	return comps, nil
}

func (c *components) Start(ctx context.Context) error {
	return c.cmp.Start(ctx)
}

func (c *components) Stop(ctx context.Context) error {
	c.replicator.Stop()
	c.cleaner.Stop()
	return c.cmp.Stop(ctx)
}

func (c *components) startWatermill(
	ctx context.Context,
	logger watermill.LoggerAdapter,
	sub message.Subscriber,
	b *bus.Bus,
	outHandler, inHandler, resultsHandler message.NoPublishHandlerFunc,
) {
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
}

func startRouter(ctx context.Context, router *message.Router) {
	go func() {
		if err := router.Run(ctx); err != nil {
			inslogger.FromContext(ctx).Error("Error while running router", err)
		}
	}()
	<-router.Running()
}
