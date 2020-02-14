// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

// +build slowtest

package integration_test

import (
	"context"
	"crypto"
	"io/ioutil"
	"math"
	"sync"
	"time"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	"github.com/dgraph-io/badger"
	"github.com/pkg/errors"

	component "github.com/insolar/component-manager"
	"github.com/insolar/insolar/application"
	"github.com/insolar/insolar/application/genesis"
	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/cryptography"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/jetcoordinator"
	"github.com/insolar/insolar/insolar/node"
	"github.com/insolar/insolar/insolar/payload"
	insolarPulse "github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/insolar/store"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/keystore"
	"github.com/insolar/insolar/ledger/artifact"
	"github.com/insolar/insolar/ledger/drop"
	"github.com/insolar/insolar/ledger/heavy/executor"
	"github.com/insolar/insolar/ledger/heavy/handler"
	"github.com/insolar/insolar/ledger/heavy/pulsemanager"
	"github.com/insolar/insolar/ledger/object"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/log/logwatermill"
	"github.com/insolar/insolar/network"
	networknode "github.com/insolar/insolar/network/node"
	"github.com/insolar/insolar/platformpolicy"
	"github.com/insolar/insolar/pulse"
)

type badgerLogger struct {
	insolar.Logger
}

func (b badgerLogger) Warningf(fmt string, args ...interface{}) {
	b.Warnf(fmt, args...)
}

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
	JetKeeper    executor.JetKeeper
	replicator   executor.HeavyReplicator
	dbRollback   *executor.DBRollback

	serverPubSub *gochannel.GoChannel
	clientPubSub *gochannel.GoChannel
}

// After using it you have to remove directory configuration.Storage.DataDirectory by yourself
func DefaultHeavyConfig() configuration.Configuration {
	tmpDir, err := ioutil.TempDir("", "heavy-integr-test-")
	if err != nil {
		panic(err)
	}
	cfg := configuration.Configuration{}
	cfg.KeysPath = "testdata/bootstrap_keys.json"
	cfg.Ledger.LightChainLimit = math.MaxInt32
	cfg.Bus.ReplyTimeout = time.Minute
	cfg.Ledger.Storage = configuration.Storage{
		DataDirectory: tmpDir,
	}
	return cfg
}

func defaultReceiveCallback(meta payload.Meta, pl payload.Payload) []payload.Payload {
	return nil
}

func NewBadgerServer(
	ctx context.Context,
	cfg configuration.Configuration,
	genesisCfg application.GenesisHeavyConfig,
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
		Pulses      *insolarPulse.BadgerDB
		Jets        *jet.BadgerDBStore
		Nodes       *node.BadgerStorageDB
		DB          *store.BadgerDB
		DBRollback  *executor.DBRollback
	)
	{
		var err error
		options := badger.DefaultOptions(cfg.Ledger.Storage.DataDirectory)
		options.Logger = badgerLogger{Logger: inslogger.FromContext(ctx).WithField("component", "badger")}
		DB, err = store.NewBadgerDB(options)
		if err != nil {
			panic(errors.Wrap(err, "failed to initialize DB"))
		}
		Nodes = node.NewBadgerStorageDB(DB)
		Pulses = insolarPulse.NewBadgerDB(DB)
		Jets = jet.NewBadgerDBStore(DB)

		c := jetcoordinator.NewJetCoordinator(cfg.Ledger.LightChainLimit, light.ref)
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

		c := jetcoordinator.NewJetCoordinator(cfg.Ledger.LightChainLimit, virtual.ref)
		c.PulseCalculator = Pulses
		c.PulseAccessor = Pulses
		c.JetAccessor = Jets
		c.PlatformCryptographyScheme = CryptoScheme
		c.Nodes = Nodes
		ClientBus = bus.NewBus(cfg.Bus, ClientPubSub, Pulses, c, CryptoScheme)
	}

	var replicator executor.HeavyReplicator

	// Heavy components.
	var (
		PulseManager insolar.PulseManager
		Handler      *handler.Handler
		Genesis      *genesis.Genesis
		Records      *object.BadgerRecordDB
		JetKeeper    *executor.BadgerDBJetKeeper
	)
	{
		Records = object.NewBadgerRecordDB(DB)
		indexes := object.NewBadgerIndexDB(DB, Records)
		drops := drop.NewBadgerDB(DB)
		JetKeeper = executor.NewBadgerJetKeeper(Jets, DB, Pulses)
		DBRollback = executor.NewDBRollback(JetKeeper, drops, Records, indexes, Jets, Pulses, JetKeeper, Nodes)

		sp := insolarPulse.NewStartPulse()

		backupMaker, err := executor.NewBackupMaker(ctx, DB, cfg.Ledger, JetKeeper.TopSyncPulse(), DB)
		if err != nil {
			return nil, errors.Wrap(err, "failed create backuper")
		}

		gcRunInfo := executor.NewBadgerGCRunInfo(DB, cfg.Ledger.Storage.GCRunFrequency)
		replicator = executor.NewHeavyReplicatorDefault(Records, indexes, CryptoScheme, Pulses, drops, JetKeeper, backupMaker, Jets, gcRunInfo)

		pm := pulsemanager.NewPulseManager(nil)
		pm.NodeNet = NodeNetwork
		pm.NodeSetter = Nodes
		pm.Nodes = Nodes
		pm.PulseAppender = Pulses
		pm.PulseAccessor = Pulses
		pm.JetModifier = Jets
		pm.StartPulse = sp
		pm.FinalizationKeeper = executor.NewFinalizationKeeperDefault(JetKeeper, Pulses, cfg.Ledger.LightChainLimit)

		h := handler.New(cfg.Ledger, gcRunInfo)
		h.RecordAccessor = Records
		h.RecordModifier = Records
		h.JetCoordinator = Coordinator
		h.IndexAccessor = indexes
		h.IndexModifier = indexes
		h.DropModifier = drops
		h.PCS = CryptoScheme
		h.PulseAccessor = Pulses
		h.PulseCalculator = Pulses
		h.StartPulse = sp
		h.JetModifier = Jets
		h.JetAccessor = Jets
		h.JetTree = Jets
		h.JetKeeper = JetKeeper
		h.BackupMaker = backupMaker
		h.Sender = ClientBus
		h.Replicator = replicator

		PulseManager = pm
		Handler = h

		artifactManager := &artifact.Scope{
			PulseNumber:    pulse.MinTimePulse,
			PCS:            CryptoScheme,
			RecordAccessor: Records,
			RecordModifier: Records,
			IndexModifier:  indexes,
			IndexAccessor:  indexes,
		}
		Genesis = &genesis.Genesis{
			ArtifactManager: artifactManager,
			IndexModifier:   indexes,
			BaseRecord: &genesis.BadgerBaseRecord{
				DB:             DB,
				DropModifier:   drops,
				PulseAppender:  Pulses,
				PulseAccessor:  Pulses,
				RecordModifier: Records,
				IndexModifier:  indexes,
			},

			DiscoveryNodes:  genesisCfg.DiscoveryNodes,
			ContractsConfig: genesisCfg.ContractsConfig,
		}
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
			"OutgoingFromClient",
			bus.TopicOutgoing,
			ClientPubSub,
			Handler.Process,
		)

		startRouter(ctx, inRouter)
		startRouter(ctx, outRouter)
	}

	inslogger.FromContext(ctx).WithFields(map[string]interface{}{
		"light":   light.ID().String(),
		"virtual": virtual.ID().String(),
		"heavy":   heavy.ID().String(),
	}).Info("started test server")

	if err := Genesis.Start(ctx); err != nil {
		log.Fatalf("genesis failed on heavy with error: %v", err)
	}

	err := DBRollback.Start(ctx)
	if err != nil {
		log.Fatalf("rollback.Start return error: %v", err)
	}

	s := &Server{
		pm:           PulseManager,
		pulse:        *insolar.GenesisPulse,
		clientSender: ClientBus,
		JetKeeper:    JetKeeper,
		replicator:   replicator,
		dbRollback:   DBRollback,
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
