package integration_test

import (
	"context"
	"crypto"
	"sync"

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
	"github.com/insolar/insolar/insolar/reply"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/instrumentation/instracer"
	"github.com/insolar/insolar/keystore"
	"github.com/insolar/insolar/ledger/drop"
	"github.com/insolar/insolar/ledger/light/artifactmanager"
	"github.com/insolar/insolar/ledger/light/executor"
	"github.com/insolar/insolar/ledger/light/hot"
	"github.com/insolar/insolar/ledger/light/pulsemanager"
	"github.com/insolar/insolar/ledger/light/replication"
	"github.com/insolar/insolar/ledger/object"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/platformpolicy"
	"github.com/pkg/errors"
)

type Server struct {
	pm      insolar.PulseManager
	network *networkMock
	handler message.HandlerFunc
	pulse   insolar.Pulse
	lock    sync.RWMutex
}

func DefaultLightConfig() configuration.Configuration {
	cfg := configuration.Configuration{}
	cfg.KeysPath = "testdata/bootstrap_keys.json"
	cfg.Ledger.LightChainLimit = 5
	return cfg
}

func NewServer(ctx context.Context, cfg configuration.Configuration) (*Server, error) {
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
		Tokens = delegationtoken.NewDelegationTokenFactory()
		Bus = &stub{}
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
		records := object.NewRecordMemory()
		indexes := object.NewIndexStorageMemory()
		writeController := hot.NewWriteController()

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
		handler.IndexLocker = idLocker
		handler.Records = records
		handler.Nodes = Nodes
		handler.HotDataWaiter = waiter
		handler.JetReleaser = waiter
		handler.WriteAccessor = writeController
		handler.Sender = WmBus
		handler.IndexStorage = indexes
		err := handler.Init(ctx)
		if err != nil {
			return nil, err
		}

		jetCalculator := executor.NewJetCalculator(Coordinator, Jets)
		var lightCleaner = replication.NewCleaner(
			Jets.(jet.Cleaner),
			Nodes,
			drops,
			records,
			indexes,
			Pulses,
			Pulses,
			indexes,
			handler.FilamentCalculator,
			conf.LightChainLimit,
		)

		lthSyncer := replication.NewReplicatorDefault(
			jetCalculator,
			lightCleaner,
			Bus,
			Pulses,
			drops,
			records,
			indexes,
			Jets,
		)

		jetSplitter := executor.NewJetSplitter(cfg.Ledger.JetSplit, jetCalculator, Jets, Jets, drops, drops, Pulses, records)

		hotSender := executor.NewHotSender(
			drops,
			indexes,
			Pulses,
			Jets,
			conf.LightChainLimit,
			WmBus,
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
		pm.GIL = &stub{}
		pm.NodeNet = NodeNetwork

		PulseManager = pm
		Handler = handler
	}

	netMock := newNetworkMock(WmBus.IncomingMessageRouter(Handler.FlowDispatcher.Process), NodeNetwork.GetOrigin().ID())
	startWatermill(ctx, logger, pubSub, WmBus, netMock.Handle, Handler.FlowDispatcher.Process)

	s := &Server{
		pm:      PulseManager,
		handler: Handler.FlowDispatcher.Process,
		pulse:   *insolar.GenesisPulse,
		network: netMock,
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

func (s *Server) Pulse(ctx context.Context) {
	s.pulse = insolar.Pulse{
		PulseNumber: s.pulse.PulseNumber + 10,
	}
	err := s.pm.Set(ctx, s.pulse)
	if err != nil {
		panic(err)
	}
}

func (s *Server) Send(pl payload.Payload) {
	msg, err := payload.NewMessage(pl)
	if err != nil {
		panic(err)
	}
	meta := payload.Meta{
		Payload: msg.Payload,
		Pulse:   s.pulse.PulseNumber,
		ID:      []byte(watermill.NewUUID()),
	}
	buf, err := meta.Marshal()
	if err != nil {
		panic(err)
	}
	msg.Payload = buf
	msg.Metadata.Set(bus.MetaPulse, s.pulse.PulseNumber.String())
	msg.Metadata.Set(bus.MetaSpanData, string(instracer.MustSerialize(context.Background())))

	_, err = s.handler(msg)
	if err != nil {
		panic(err)
	}
}

func (s *Server) Receive(h func(meta payload.Meta, pl payload.Payload)) {
	s.network.SetReceiver(h)
}

type networkMock struct {
	outHandler func(meta payload.Meta, pl payload.Payload)
	inHandler  message.HandlerFunc
	me         insolar.Reference
	lock       sync.RWMutex
}

func newNetworkMock(in message.HandlerFunc, me insolar.Reference) *networkMock {
	return &networkMock{inHandler: in, me: me}
}

func (p *networkMock) Handle(msg *message.Message) ([]*message.Message, error) {
	p.lock.RLock()
	defer p.lock.RUnlock()

	meta := payload.Meta{}
	err := meta.Unmarshal(msg.Payload)
	if err != nil {
		panic(errors.Wrap(err, "failed to unmarshal meta"))
	}
	if p.outHandler != nil {
		pl, err := payload.Unmarshal(meta.Payload)
		if err != nil {
			panic(errors.Wrap(err, "failed to unmarshal payload"))
		}
		go p.outHandler(meta, pl)
	}

	if meta.Receiver != p.me {
		return nil, nil
	}

	msg.Metadata.Set(bus.MetaPulse, meta.Pulse.String())
	return p.inHandler(msg)
}

func (p *networkMock) SetReceiver(r func(meta payload.Meta, pl payload.Payload)) {
	p.lock.Lock()
	p.outHandler = r
	p.lock.Unlock()
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

func (*stub) Send(context.Context, insolar.Message, *insolar.MessageSendOptions) (insolar.Reply, error) {
	return &reply.OK{}, nil
}

func (*stub) Register(p insolar.MessageType, handler insolar.MessageHandler) error {
	return nil
}

func (*stub) MustRegister(p insolar.MessageType, handler insolar.MessageHandler) {
}

func (*stub) OnPulse(context.Context, insolar.Pulse) error {
	return nil
}

func (*stub) Acquire(ctx context.Context) {}

func (*stub) Release(ctx context.Context) {}

func (*stub) MoveSyncToActive(ctx context.Context, number insolar.PulseNumber) error {
	return nil
}
