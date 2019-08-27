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

	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/insolar/flow/dispatcher"
	"github.com/insolar/insolar/ledger/light/handle"
	"github.com/insolar/insolar/ledger/light/proc"
	"github.com/insolar/insolar/network"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	"github.com/pkg/errors"

	"github.com/insolar/insolar/api"
	"github.com/insolar/insolar/certificate"
	"github.com/insolar/insolar/component"
	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/contractrequester"
	"github.com/insolar/insolar/cryptography"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/delegationtoken"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/jetcoordinator"
	"github.com/insolar/insolar/insolar/node"
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/keystore"
	"github.com/insolar/insolar/ledger/drop"
	"github.com/insolar/insolar/ledger/light/executor"
	"github.com/insolar/insolar/ledger/object"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/logicrunner/artifacts"
	"github.com/insolar/insolar/metrics"
	"github.com/insolar/insolar/network/nodenetwork"
	"github.com/insolar/insolar/network/servicenetwork"
	"github.com/insolar/insolar/network/termination"
	"github.com/insolar/insolar/platformpolicy"
	"github.com/insolar/insolar/server/internal"
)

type components struct {
	cmp               component.Manager
	NodeRef, NodeRole string
	replicator        executor.LightReplicator
	cleaner           executor.Cleaner
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

		c := component.Manager{}
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
	comps.cmp = component.Manager{}
	comps.NodeRef = CertManager.GetCertificate().GetNodeRef().String()
	comps.NodeRole = CertManager.GetCertificate().GetRole().String()

	// Watermill stuff.
	var (
		wmLogger   *log.WatermillLogAdapter
		publisher  message.Publisher
		subscriber message.Subscriber
	)
	{
		wmLogger = log.NewWatermillLogAdapter(inslogger.FromContext(ctx))
		pubsub := gochannel.NewGoChannel(gochannel.Config{}, wmLogger)
		subscriber = pubsub
		publisher = pubsub
		// Wrapped watermill publisher for introspection.
		publisher = internal.PublisherWrapper(ctx, &comps.cmp, cfg.Introspection, publisher)
	}

	// Network.
	var (
		NetworkService *servicenetwork.ServiceNetwork
		NodeNetwork    network.NodeNetwork
		Termination    insolar.TerminationHandler
	)
	{
		var err error
		// External communication.
		NetworkService, err = servicenetwork.NewServiceNetwork(cfg, &comps.cmp)
		if err != nil {
			return nil, errors.Wrap(err, "failed to start Network")
		}

		Termination = termination.NewHandler(NetworkService)

		// Node info.
		NodeNetwork, err = nodenetwork.NewNodeNetwork(cfg.Host.Transport, CertManager.GetCertificate())
		if err != nil {
			return nil, errors.Wrap(err, "failed to start NodeNetwork")
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

		c := jetcoordinator.NewJetCoordinator(cfg.Ledger.LightChainLimit)
		c.PulseCalculator = Pulses
		c.PulseAccessor = Pulses
		c.JetAccessor = Jets
		c.OriginProvider = NodeNetwork
		c.PlatformCryptographyScheme = CryptoScheme
		c.Nodes = Nodes

		Coordinator = c
	}

	// Communication.
	var (
		Tokens insolar.DelegationTokenFactory
		Sender *bus.Bus
	)
	{
		Tokens = delegationtoken.NewDelegationTokenFactory()
		Sender = bus.NewBus(cfg.Bus, publisher, Pulses, Coordinator, CryptoScheme)
	}

	// API.
	var (
		Requester       *contractrequester.ContractRequester
		ArtifactsClient = artifacts.NewClient(Sender)
		APIWrapper      *api.RunnerWrapper
	)
	{
		var err error
		Requester, err = contractrequester.New()
		if err != nil {
			return nil, errors.Wrap(err, "failed to start ContractRequester")
		}

		API, err := api.NewRunner(
			&cfg.APIRunner,
			CertManager,
			Requester,
			NodeNetwork,
			NetworkService,
			Pulses,
			ArtifactsClient,
			Coordinator,
			NetworkService,
		)
		if err != nil {
			return nil, errors.Wrap(err, "failed to start ApiRunner")
		}

		AdminAPIRunner, err := api.NewRunner(
			&cfg.AdminAPIRunner,
			CertManager,
			Requester,
			NodeNetwork,
			NetworkService,
			Pulses,
			ArtifactsClient,
			Coordinator,
			NetworkService,
		)
		if err != nil {
			return nil, errors.Wrap(err, "failed to start AdminAPIRunner")
		}

		APIWrapper = api.NewWrapper(API, AdminAPIRunner)
	}

	metricsHandler, err := metrics.NewMetrics(
		ctx,
		cfg.Metrics,
		metrics.GetInsolarRegistry(comps.NodeRole),
		comps.NodeRole,
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to start Metrics")
	}

	// Light components.
	var (
		PulseManager   insolar.PulseManager
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

		c := component.Manager{}
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
			NodeNetwork,
			FlowDispatcher,
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
		)
	}

	comps.cmp.Inject(
		Sender,
		Jets,
		Pulses,
		Coordinator,
		PulseManager,
		metricsHandler,
		Requester,
		Tokens,
		ArtifactsClient,
		APIWrapper,
		KeyProcessor,
		Termination,
		CryptoScheme,
		CryptoService,
		CertManager,
		NodeNetwork,
		NetworkService,
		publisher,
	)

	err = comps.cmp.Init(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to init components")
	}

	comps.startWatermill(ctx, wmLogger, subscriber, Sender, NetworkService.SendMessageHandler, FlowDispatcher.Process, Requester.ReceiveResult)

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
