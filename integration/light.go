package integration

import (
	"context"
	"crypto"

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
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/jetcoordinator"
	"github.com/insolar/insolar/insolar/node"
	"github.com/insolar/insolar/insolar/payload"
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
	"github.com/insolar/insolar/platformpolicy"
	"github.com/pkg/errors"
)

type LightServer struct {
	pm      insolar.PulseManager
	network *networkMock
	handler message.HandlerFunc
	pulse   insolar.Pulse
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
		NodeNetwork insolar.NodeNetwork
	)
	{
		NodeNetwork = newNodeNetMock()
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
		c.NodeNet = NodeNetwork

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
		pm.NodeNet = NodeNetwork

		PulseManager = pm
		Handler = handler
	}

	netMock := newNetworkMock(Handler.FlowDispatcher.Process, NodeNetwork.GetOrigin().ID())
	startWatermill(ctx, logger, pubSub, WmBus, netMock.Handle, Handler.FlowDispatcher.Process)

	s := &LightServer{
		pm:      PulseManager,
		handler: Handler.FlowDispatcher.Process,
		pulse:   *insolar.GenesisPulse,
	}
	return s, nil
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

func (s *LightServer) Pulse(ctx context.Context) error {
	s.pulse = insolar.Pulse{
		PulseNumber: s.pulse.PulseNumber + 10,
	}
	return s.pm.Set(ctx, s.pulse)
}

func (s *LightServer) Send(msg *message.Message) error {
	_, err := s.handler(msg)
	return err
}

func (s *LightServer) Receive(h func(*message.Message)) {
	s.network.outHandler = h
}

type networkMock struct {
	outHandler func(*message.Message)
	inHandler  message.HandlerFunc
	me         insolar.Reference
}

func newNetworkMock(in message.HandlerFunc, me insolar.Reference) *networkMock {
	return &networkMock{inHandler: in, me: me}
}

func (p *networkMock) Handle(msg *message.Message) ([]*message.Message, error) {
	if p.outHandler != nil {
		p.outHandler(msg)
	}

	meta := payload.Meta{}
	err := meta.Unmarshal(msg.Payload)
	if err != nil {
		panic(errors.Wrap(err, "failed to unmarshal message"))
	}
	if meta.Receiver != p.me {
		return nil, nil
	}

	return p.inHandler(msg)
}

type nodeMock struct {
	ref  insolar.Reference
	role insolar.StaticRole
}

func (n *nodeMock) ID() insolar.Reference {
	return n.ref
}

func (n *nodeMock) ShortID() insolar.ShortNodeID {
	panic("implement me")
}

func (n *nodeMock) Role() insolar.StaticRole {
	return n.role
}

func (n *nodeMock) PublicKey() crypto.PublicKey {
	panic("implement me")
}

func (n *nodeMock) Address() string {
	panic("implement me")
}

func (n *nodeMock) GetGlobuleID() insolar.GlobuleID {
	panic("implement me")
}

func (n *nodeMock) Version() string {
	panic("implement me")
}

func (n *nodeMock) LeavingETA() insolar.PulseNumber {
	panic("implement me")
}

func (n *nodeMock) GetState() insolar.NodeState {
	panic("implement me")
}

type nodeNetMock struct {
	nodes []insolar.NetworkNode
}

func newNodeNetMock() *nodeNetMock {
	return &nodeNetMock{nodes: []insolar.NetworkNode{
		&nodeMock{
			ref:  gen.Reference(),
			role: insolar.StaticRoleHeavyMaterial,
		},
		&nodeMock{
			ref:  gen.Reference(),
			role: insolar.StaticRoleVirtual,
		},
		&nodeMock{
			ref:  gen.Reference(),
			role: insolar.StaticRoleLightMaterial,
		},
		&nodeMock{
			ref:  gen.Reference(),
			role: insolar.StaticRoleVirtual,
		},
		&nodeMock{
			ref:  gen.Reference(),
			role: insolar.StaticRoleLightMaterial,
		},
	}}
}

func (n *nodeNetMock) GetOrigin() insolar.NetworkNode {
	return n.nodes[2]
}

func (n *nodeNetMock) GetWorkingNode(ref insolar.Reference) insolar.NetworkNode {
	panic("implement me")
}

func (n *nodeNetMock) GetWorkingNodes() []insolar.NetworkNode {
	return n.nodes
}

func (n *nodeNetMock) GetWorkingNodesByRole(role insolar.DynamicRole) []insolar.Reference {
	panic("implement me")
}

type stub struct{}

func (*stub) Acquire(ctx context.Context) {}

func (*stub) Release(ctx context.Context) {}

func (*stub) MoveSyncToActive(ctx context.Context, number insolar.PulseNumber) error {
	return nil
}
