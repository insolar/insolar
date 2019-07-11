package integration

import (
	"context"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/infrastructure/gochannel"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"github.com/insolar/insolar/component"
	"github.com/insolar/insolar/configuration"
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
	"github.com/insolar/insolar/ledger/blob"
	"github.com/insolar/insolar/ledger/drop"
	"github.com/insolar/insolar/ledger/light/artifactmanager"
	"github.com/insolar/insolar/ledger/light/executor"
	"github.com/insolar/insolar/ledger/light/hot"
	"github.com/insolar/insolar/ledger/light/pulsemanager"
	"github.com/insolar/insolar/ledger/light/replication"
	"github.com/insolar/insolar/ledger/object"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/messagebus"
	"github.com/insolar/insolar/network/servicenetwork"
	"github.com/insolar/insolar/platformpolicy"
	"github.com/pkg/errors"
)

type LightServer struct {
	pm insolar.PulseManager
}

func DefaultLightConfig() configuration.Configuration {
	cfg := configuration.Configuration{}
	cfg.KeysPath = "testdata/bootstrap_keys.json"
	return cfg
}

func NewLightServer(ctx context.Context, cfg configuration.Configuration) (*LightServer, error) {
	// Cryptography.
	var (
		KeyProcessor  insolar.KeyProcessor
		CryptoScheme  insolar.PlatformCryptographyScheme
		CryptoService insolar.CryptographyService
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
	}

	logger := log.NewWatermillLogAdapter(inslogger.FromContext(ctx))
	pubSub := gochannel.NewGoChannel(gochannel.Config{}, logger)

	// Network.
	var (
		NetworkService *servicenetwork.ServiceNetwork
		NodeNetwork    insolar.NodeNetwork
	)
	{

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
		c.NodeNet = NodeNetwork
		c.PlatformCryptographyScheme = CryptoScheme
		c.Nodes = Nodes

		Coordinator = c
	}

	// Communication.
	var (
		Tokens insolar.DelegationTokenFactory
		Bus    insolar.MessageBus
		WmBus  *bus.Bus
	)
	{
		var err error
		Tokens = delegationtoken.NewDelegationTokenFactory()
		Bus, err = messagebus.NewMessageBus(cfg)
		if err != nil {
			return nil, errors.Wrap(err, "failed to start MessageBus")
		}
		WmBus = bus.NewBus(pubSub, Pulses, Coordinator, CryptoScheme)
	}

	// Light components.
	var (
		PulseManager insolar.PulseManager
		Handler      *artifactmanager.MessageHandler
	)
	{
		conf := cfg.Ledger
		idLocker := object.NewIndexLocker()
		drops := drop.NewStorageMemory()
		blobs := blob.NewStorageMemory()
		records := object.NewRecordMemory()
		indexes := object.NewIndexStorageMemory()
		writeController := hot.NewWriteController()

		c := component.Manager{}
		c.Inject(CryptoScheme)

		waiter := hot.NewChannelWaiter()

		handler := artifactmanager.NewMessageHandler(&conf)
		handler.PulseCalculator = Pulses

		handler.Bus = Bus
		handler.PCS = CryptoScheme
		handler.JetCoordinator = Coordinator
		handler.CryptographyService = CryptoService
		handler.DelegationTokenFactory = Tokens
		handler.JetStorage = Jets
		handler.DropModifier = drops
		handler.BlobModifier = blobs
		handler.BlobAccessor = blobs
		handler.Blobs = blobs
		handler.IndexLocker = idLocker
		handler.Records = records
		handler.Nodes = Nodes
		handler.HotDataWaiter = waiter
		handler.JetReleaser = waiter
		handler.WriteAccessor = writeController
		handler.Sender = WmBus
		handler.IndexStorage = indexes

		jetCalculator := executor.NewJetCalculator(Coordinator, Jets)
		var lightCleaner = replication.NewCleaner(
			Jets.(jet.Cleaner),
			Nodes,
			drops,
			blobs,
			records,
			indexes,
			Pulses,
			Pulses,
			conf.LightChainLimit,
		)

		lthSyncer := replication.NewReplicatorDefault(
			jetCalculator,
			lightCleaner,
			Bus,
			Pulses,
			drops,
			blobs,
			records,
			indexes,
			Jets,
		)

		jetSplitter := executor.NewJetSplitter(jetCalculator, Jets, Jets, drops, drops, Pulses)

		hotSender := executor.NewHotSender(
			Bus,
			drops,
			indexes,
			Pulses,
			Jets,
			conf.LightChainLimit,
		)

		pm := pulsemanager.NewPulseManager(
			jetSplitter,
			lthSyncer,
			writeController,
			hotSender,
		)
		pm.MessageHandler = handler
		pm.Bus = Bus
		pm.NodeNet = NodeNetwork
		pm.JetReleaser = waiter
		pm.JetModifier = Jets
		pm.NodeSetter = Nodes
		pm.Nodes = Nodes
		pm.PulseAccessor = Pulses
		pm.PulseCalculator = Pulses
		pm.PulseAppender = Pulses
		pm.ActiveListSwapper = &stub{}
		pm.GIL = &stub{}
		pm.NodeNet = &stub{}

		PulseManager = pm
		Handler = handler
	}

	startWatermill(ctx, logger, pubSub, WmBus, NetworkService.SendMessageHandler, Handler.FlowDispatcher.Process)

	return &LightServer{pm: PulseManager}, nil
}

func startWatermill(
	ctx context.Context,
	logger watermill.LoggerAdapter,
	pubSub message.PubSub,
	b *bus.Bus,
	outHandler, inHandler message.HandlerFunc,
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
		pubSub,
		outHandler,
	)

	inRouter.AddMiddleware(
		middleware.InstantAck,
		b.IncomingMessageRouter,
	)

	inRouter.AddNoPublisherHandler(
		"IncomingHandler",
		bus.TopicIncoming,
		pubSub,
		inHandler,
	)

	startRouter(ctx, inRouter)
	startRouter(ctx, outRouter)
}

func startRouter(ctx context.Context, router *message.Router) {
	go func() {
		if err := router.Run(); err != nil {
			inslogger.FromContext(ctx).Error("Error while running router", err)
		}
	}()
	<-router.Running()
}

type stub struct{}

func (*stub) GetOrigin() insolar.NetworkNode {
	panic("implement me")
}

func (*stub) GetWorkingNode(ref insolar.Reference) insolar.NetworkNode {
	panic("implement me")
}

func (*stub) GetWorkingNodes() []insolar.NetworkNode {
	return []insolar.NetworkNode{}
}

func (*stub) GetWorkingNodesByRole(role insolar.DynamicRole) []insolar.Reference {
	panic("implement me")
}

func (*stub) Acquire(ctx context.Context) {
}

func (*stub) Release(ctx context.Context) {
}

func (*stub) MoveSyncToActive(ctx context.Context, number insolar.PulseNumber) error {
	return nil
}

func (s *LightServer) Pulse(ctx context.Context, pulse insolar.Pulse) error {
	return s.pm.Set(ctx, pulse)
}
