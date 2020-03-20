// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

// +build slowtest

package integration_test

import (
	"context"
	"crypto"
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/pkg/errors"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"

	component "github.com/insolar/component-manager"
	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/cryptography"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/insolar/flow/dispatcher"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/jetcoordinator"
	"github.com/insolar/insolar/insolar/node"
	"github.com/insolar/insolar/insolar/payload"
	insolarPulse "github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/keystore"
	"github.com/insolar/insolar/ledger/drop"
	"github.com/insolar/insolar/ledger/light/executor"
	"github.com/insolar/insolar/ledger/light/handle"
	"github.com/insolar/insolar/ledger/light/proc"
	"github.com/insolar/insolar/ledger/object"
	"github.com/insolar/insolar/log/logwatermill"
	"github.com/insolar/insolar/metrics"
	"github.com/insolar/insolar/network"
	networknode "github.com/insolar/insolar/network/node"
	"github.com/insolar/insolar/platformpolicy"
	"github.com/insolar/insolar/pulse"
)

var (
	light = nodeMock{
		ref:     gen.Reference(),
		shortID: 1,
		role:    insolar.StaticRoleLightMaterial,
	}
	heavy = nodeMock{
		ref:     gen.Reference(),
		shortID: 2,
		role:    insolar.StaticRoleHeavyMaterial,
	}
	virtual = nodeMock{
		ref:     gen.Reference(),
		shortID: 3,
		role:    insolar.StaticRoleVirtual,
	}
)

func NodeHeavy() insolar.Reference {
	return heavy.ref
}

const PulseStep insolar.PulseNumber = 10

type Server struct {
	pm           insolar.PulseManager
	pulse        insolar.Pulse
	lock         sync.RWMutex
	clientSender bus.Sender
	replicator   executor.LightReplicator
	cleaner      executor.Cleaner

	serverPubSub *gochannel.GoChannel
	clientPubSub *gochannel.GoChannel

	metrics *metrics.Metrics
}

func DefaultLightConfig() configuration.LightConfig {
	cfg := configuration.LightConfig{}
	cfg.KeysPath = "testdata/bootstrap_keys.json"
	cfg.LightChainLimit = math.MaxInt32
	cfg.Ledger.JetSplit.DepthLimit = math.MaxUint8
	cfg.Ledger.JetSplit.ThresholdOverflowCount = math.MaxInt32
	cfg.Ledger.JetSplit.ThresholdRecordsCount = math.MaxInt32
	cfg.Bus.ReplyTimeout = time.Minute
	return cfg
}

func DefaultLightInitialState() *payload.LightInitialState {
	return &payload.LightInitialState{
		NetworkStart: true,
		JetIDs:       []insolar.JetID{insolar.ZeroJetID},
		Pulse: insolarPulse.PulseProto{
			PulseNumber: pulse.MinTimePulse,
		},
		Drops: []drop.Drop{
			{JetID: insolar.ZeroJetID, Pulse: pulse.MinTimePulse},
		},
		LightChainLimit: math.MaxInt32,
	}
}

func DefaultHeavyResponse(pl payload.Payload) []payload.Payload {
	switch pl.(type) {
	case *payload.Replication, *payload.GotHotConfirmation:
		return nil
	case *payload.GetLightInitialState:
		return []payload.Payload{DefaultLightInitialState()}
	case *payload.SearchIndex:
		return []payload.Payload{&payload.SearchIndexInfo{}}
	}

	panic(fmt.Sprintf("unexpected message to heavy %T", pl))
}

func defaultReceiveCallback(meta payload.Meta, pl payload.Payload) []payload.Payload {
	if meta.Receiver == NodeHeavy() {
		return DefaultHeavyResponse(pl)
	}
	return nil
}

func NewServer(
	ctx context.Context,
	cfg configuration.LightConfig,
	receiveCallback func(meta payload.Meta, pl payload.Payload) []payload.Payload,
) (*Server, error) {
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

		c := component.NewManager(nil)
		c.Inject(CryptoService, CryptoScheme, KeyProcessor, ks)
	}

	// Network.
	var (
		NodeNetwork network.NodeNetwork
	)
	{
		NodeNetwork = newNodeNetMock(&light)
	}

	// Role calculations.
	var (
		Coordinator jet.Coordinator
		Pulses      *insolarPulse.StorageMem
		Jets        jet.Storage
		Nodes       *node.Storage
	)
	{
		Nodes = node.NewStorage()
		Pulses = insolarPulse.NewStorageMem()
		Jets = jet.NewStore()

		c := jetcoordinator.NewJetCoordinator(cfg.LightChainLimit, light.ref)
		c.PulseCalculator = Pulses
		c.PulseAccessor = Pulses
		c.JetAccessor = Jets
		c.PlatformCryptographyScheme = CryptoScheme
		c.Nodes = Nodes

		Coordinator = c
	}

	logger := logwatermill.NewWatermillLogAdapter(inslogger.FromContext(ctx))

	// Communication.
	var (
		ServerBus, ClientBus       *bus.Bus
		ServerPubSub, ClientPubSub *gochannel.GoChannel
	)
	{
		ServerPubSub = gochannel.NewGoChannel(gochannel.Config{}, logger)
		ClientPubSub = gochannel.NewGoChannel(gochannel.Config{}, logger)
		ServerBus = bus.NewBus(cfg.Bus, ServerPubSub, Pulses, Coordinator, CryptoScheme)

		c := jetcoordinator.NewJetCoordinator(cfg.LightChainLimit, virtual.ref)
		c.PulseCalculator = Pulses
		c.PulseAccessor = Pulses
		c.JetAccessor = Jets
		c.PlatformCryptographyScheme = CryptoScheme
		c.Nodes = Nodes
		ClientBus = bus.NewBus(cfg.Bus, ClientPubSub, Pulses, c, CryptoScheme)
	}

	// Light components.
	var (
		PulseManager   insolar.PulseManager
		Replicator     executor.LightReplicator
		Cleaner        executor.Cleaner
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

		jetFetcher := executor.NewFetcher(Nodes, Jets, ServerBus, Coordinator)
		filamentCalculator := executor.NewFilamentCalculator(
			indexes,
			records,
			Coordinator,
			jetFetcher,
			ServerBus,
			Pulses,
		)
		requestChecker := executor.NewRequestChecker(
			filamentCalculator,
			Coordinator,
			jetFetcher,
			CryptoScheme,
			ServerBus,
		)
		detachedNotifier := executor.NewDetachedNotifierDefault(ServerBus)

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
			cfg.LightChainLimit,
			conf.CleanerDelay,
			conf.FilamentCacheLimit,
		)
		Cleaner = lightCleaner

		lthSyncer := executor.NewReplicatorDefault(
			jetCalculator,
			lightCleaner,
			ServerBus,
			Pulses,
			drops,
			records,
			indexes,
			Jets,
		)
		Replicator = lthSyncer

		jetSplitter := executor.NewJetSplitter(cfg.Ledger.JetSplit, jetCalculator, Jets, Jets, drops, drops, Pulses, records)

		hotSender := executor.NewHotSender(
			drops,
			indexes,
			Pulses,
			Jets,
			cfg.LightChainLimit,
			ServerBus,
		)

		stateIniter := executor.NewStateIniter(
			cfg.LightChainLimit, Jets, hotWaitReleaser, drops, Nodes, ServerBus, Pulses, Pulses, jetCalculator, indexes,
		)

		metricsRegistry := executor.NewMetricsRegistry()

		dep := proc.NewDependencies(
			CryptoScheme,
			Coordinator,
			Jets,
			Pulses,
			ServerBus,
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
			configuration.NewLedgerLight(),
			metricsRegistry,
		)

		initHandle := func(msg *message.Message) *handle.Init {
			return handle.NewInit(dep, ServerBus, msg)
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
			[]dispatcher.Dispatcher{FlowDispatcher},
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

	// Start routers with handlers.
	{
		outHandler := func(msg *message.Message) error {
			meta := payload.Meta{}
			err := meta.Unmarshal(msg.Payload)
			if err != nil {
				panic(errors.Wrap(err, "failed to unmarshal meta"))
			}

			pl, err := payload.Unmarshal(meta.Payload)
			if err != nil {
				panic(nil)
			}
			go func() {
				var replies []payload.Payload
				if receiveCallback != nil {
					replies = receiveCallback(meta, pl)
				} else {
					replies = defaultReceiveCallback(meta, pl)
				}

				for _, rep := range replies {
					msg, err := payload.NewMessage(rep)
					if err != nil {
						panic(err)
					}
					ClientBus.Reply(context.Background(), meta, msg)
				}
			}()

			// Republish as incoming to self.
			if meta.Receiver == light.ID() {
				err = ServerPubSub.Publish(bus.TopicIncoming, msg)
				if err != nil {
					panic(err)
				}
				return nil
			}

			clientHandler := func(msg *message.Message) (messages []*message.Message, e error) {
				return nil, nil
			}
			// Republish as incoming to client.
			_, err = ClientBus.IncomingMessageRouter(clientHandler)(msg)

			if err != nil {
				panic(err)
			}
			return nil
		}

		inRouter, err := message.NewRouter(message.RouterConfig{}, logger)
		if err != nil {
			panic(err)
		}
		outRouter, err := message.NewRouter(message.RouterConfig{}, logger)
		if err != nil {
			panic(err)
		}

		outRouter.AddNoPublisherHandler(
			"Outgoing",
			bus.TopicOutgoing,
			ServerPubSub,
			outHandler,
		)

		inRouter.AddMiddleware(
			middleware.InstantAck,
			ServerBus.IncomingMessageRouter,
		)
		inRouter.AddNoPublisherHandler(
			"Incoming",
			bus.TopicIncoming,
			ServerPubSub,
			FlowDispatcher.Process,
		)
		inRouter.AddNoPublisherHandler(
			"OutgoingFromClient",
			bus.TopicOutgoing,
			ClientPubSub,
			FlowDispatcher.Process,
		)

		startRouter(ctx, inRouter)
		startRouter(ctx, outRouter)
	}

	inslogger.FromContext(ctx).WithFields(map[string]interface{}{
		"light":   light.ID().String(),
		"virtual": virtual.ID().String(),
		"heavy":   heavy.ID().String(),
	}).Info("started test server")

	m := metrics.NewMetrics(configuration.NewMetrics(), metrics.GetInsolarRegistry("test-server"), "test-server")
	if err := m.Init(ctx); err != nil {
		panic(err)
	}

	s := &Server{
		pm:           PulseManager,
		pulse:        *insolar.GenesisPulse,
		clientSender: ClientBus,
		replicator:   Replicator,
		cleaner:      Cleaner,
		metrics:      m,
		serverPubSub: ServerPubSub,
		clientPubSub: ClientPubSub,
	}
	return s, nil
}

func startRouter(ctx context.Context, router *message.Router) {
	go func() {
		if err := router.Run(ctx); err != nil {
			inslogger.FromContext(ctx).Error("Error while running router", err)
		}
	}()
	<-router.Running()
}

func (s *Server) SetPulse(ctx context.Context) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.pulse = insolar.Pulse{
		PulseNumber: s.pulse.PulseNumber + PulseStep,
	}
	err := s.pm.Set(ctx, s.pulse)
	if err != nil {
		panic(err)
	}
}

func (s *Server) Pulse() insolar.PulseNumber {
	s.lock.Lock()
	defer s.lock.Unlock()

	return s.pulse.PulseNumber
}

func (s *Server) Send(ctx context.Context, pl payload.Payload) (<-chan *message.Message, func()) {
	msg, err := payload.NewMessage(pl)
	if err != nil {
		panic(err)
	}
	return s.clientSender.SendTarget(ctx, msg, gen.Reference())
}

func (s *Server) Stop() {
	s.replicator.Stop()
	s.cleaner.Stop()
	err := s.clientPubSub.Close()
	if err != nil {
		panic(err)
	}
	err = s.serverPubSub.Close()
	if err != nil {
		panic(err)
	}
}

type nodeMock struct {
	ref     insolar.Reference
	shortID insolar.ShortNodeID
	role    insolar.StaticRole
}

func (n *nodeMock) ID() insolar.Reference {
	return n.ref
}

func (n *nodeMock) ShortID() insolar.ShortNodeID {
	return n.shortID
}

func (n *nodeMock) Role() insolar.StaticRole {
	return n.role
}

func (n *nodeMock) PublicKey() crypto.PublicKey {
	panic("implement me")
}

func (n *nodeMock) Address() string {
	return ""
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
	return insolar.NodeReady
}

func (n *nodeMock) GetPower() insolar.Power {
	return 1
}

type nodeNetMock struct {
	me insolar.NetworkNode
}

func (n *nodeNetMock) GetAccessor(insolar.PulseNumber) network.Accessor {
	return networknode.NewAccessor(networknode.NewSnapshot(insolar.GenesisPulse.PulseNumber, []insolar.NetworkNode{&virtual, &heavy, &light}))
}

func newNodeNetMock(me insolar.NetworkNode) *nodeNetMock {
	return &nodeNetMock{me: me}
}

func (n *nodeNetMock) GetOrigin() insolar.NetworkNode {
	return n.me
}
