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

package integration_test

import (
	"context"
	"crypto"
	"sync"

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

var (
	light = nodeMock{
		ref:  gen.Reference(),
		role: insolar.StaticRoleLightMaterial,
	}
	heavy = nodeMock{
		ref:  gen.Reference(),
		role: insolar.StaticRoleHeavyMaterial,
	}
	virtual = nodeMock{
		ref:  gen.Reference(),
		role: insolar.StaticRoleVirtual,
	}
)

type Server struct {
	pm           insolar.PulseManager
	pulse        insolar.Pulse
	lock         sync.RWMutex
	clientSender bus.Sender
}

func DefaultLightConfig() configuration.Configuration {
	cfg := configuration.Configuration{}
	cfg.KeysPath = "testdata/bootstrap_keys.json"
	cfg.Ledger.LightChainLimit = 5
	return cfg
}

func NewServer(ctx context.Context, cfg configuration.Configuration, receive func(meta payload.Meta, pl payload.Payload)) (*Server, error) {
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

	// Network.
	var (
		NodeNetwork insolar.NodeNetwork
	)
	{
		NodeNetwork = newNodeNetMock(&light)
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

	logger := log.NewWatermillLogAdapter(inslogger.FromContext(ctx))

	// Communication.
	var (
		Tokens                     insolar.DelegationTokenFactory
		Bus                        insolar.MessageBus
		ServerBus, ClientBus       *bus.Bus
		ServerPubSub, ClientPubSub message.PubSub
	)
	{
		Tokens = delegationtoken.NewDelegationTokenFactory()
		Bus = &stub{}
		ServerPubSub = gochannel.NewGoChannel(gochannel.Config{}, logger)
		ClientPubSub = gochannel.NewGoChannel(gochannel.Config{}, logger)
		ServerBus = bus.NewBus(ServerPubSub, Pulses, Coordinator, CryptoScheme)

		c := jetcoordinator.NewJetCoordinator(cfg.Ledger.LightChainLimit)
		c.PulseCalculator = Pulses
		c.PulseAccessor = Pulses
		c.JetAccessor = Jets
		c.NodeNet = newNodeNetMock(&virtual)
		c.PlatformCryptographyScheme = CryptoScheme
		c.Nodes = Nodes
		ClientBus = bus.NewBus(ClientPubSub, Pulses, c, CryptoScheme)
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
		handler.FlowDispatcher.PulseAccessor = Pulses

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
		handler.Sender = ServerBus
		handler.IndexStorage = indexes

		jetTreeUpdater := jet.NewFetcher(Nodes, Jets, Bus, Coordinator)
		filamentCalculator := executor.NewFilamentCalculator(
			indexes,
			records,
			Coordinator,
			jetTreeUpdater,
			ServerBus,
		)

		handler.JetTreeUpdater = jetTreeUpdater
		handler.FilamentCalculator = filamentCalculator

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
			ServerBus,
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
			ServerBus,
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

	// Start routers with handlers.
	{
		outHandler := func(msg *message.Message) ([]*message.Message, error) {
			meta := payload.Meta{}
			err := meta.Unmarshal(msg.Payload)
			if err != nil {
				panic(errors.Wrap(err, "failed to unmarshal meta"))
			}

			// Republish as incoming to self.
			if meta.Receiver == light.ID() {
				err = ServerPubSub.Publish(bus.TopicIncoming, msg)
				if err != nil {
					panic(err)
				}
				return nil, nil
			}

			if receive != nil {
				pl, err := payload.Unmarshal(meta.Payload)
				if err != nil {
					panic(nil)
				}
				receive(meta, pl)
			}

			clientHandler := func(msg *message.Message) (messages []*message.Message, e error) {
				return nil, nil
			}
			// Republish as incoming to client.
			_, err = ClientBus.IncomingMessageRouter(clientHandler)(msg)

			if err != nil {
				panic(err)
			}
			return nil, nil
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
			Handler.FlowDispatcher.Process,
		)
		inRouter.AddNoPublisherHandler(
			"OutgoingFromClient",
			bus.TopicOutgoing,
			ClientPubSub,
			Handler.FlowDispatcher.Process,
		)

		startRouter(ctx, inRouter)
		startRouter(ctx, outRouter)
	}

	inslogger.FromContext(ctx).WithFields(map[string]interface{}{
		"light":   light.ID().String(),
		"virtual": virtual.ID().String(),
		"heavy":   heavy.ID().String(),
	}).Info("started test server")

	s := &Server{
		pm:           PulseManager,
		pulse:        *insolar.GenesisPulse,
		clientSender: ClientBus,
	}
	return s, nil
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
	s.lock.Lock()
	defer s.lock.Unlock()

	s.pulse = insolar.Pulse{
		PulseNumber: s.pulse.PulseNumber + 10,
	}
	err := s.pm.Set(ctx, s.pulse)
	if err != nil {
		panic(err)
	}
}

func (s *Server) Send(ctx context.Context, pl payload.Payload) (<-chan *message.Message, func()) {
	msg, err := payload.NewMessage(pl)
	if err != nil {
		panic(err)
	}
	return s.clientSender.SendTarget(ctx, msg, insolar.Reference{})
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
	me insolar.NetworkNode
}

func newNodeNetMock(me insolar.NetworkNode) *nodeNetMock {
	return &nodeNetMock{me: me}
}

func (n *nodeNetMock) GetOrigin() insolar.NetworkNode {
	return n.me
}

func (n *nodeNetMock) GetWorkingNode(ref insolar.Reference) insolar.NetworkNode {
	panic("implement me")
}

func (n *nodeNetMock) GetWorkingNodes() []insolar.NetworkNode {
	return []insolar.NetworkNode{
		&virtual,
		&heavy,
		&light,
	}
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
