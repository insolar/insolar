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
// +build slowtest

package integration

import (
	"context"
	"flag"
	"io"
	"math"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/insolar/insolar/log/logwatermill"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	"github.com/pkg/errors"

	"github.com/insolar/insolar/component"
	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/contractrequester"
	"github.com/insolar/insolar/cryptography"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/bus/meta"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/jetcoordinator"
	"github.com/insolar/insolar/insolar/node"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/insolar/utils"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/keystore"
	"github.com/insolar/insolar/logicrunner"
	"github.com/insolar/insolar/logicrunner/artifacts"
	"github.com/insolar/insolar/logicrunner/logicexecutor"
	"github.com/insolar/insolar/logicrunner/machinesmanager"
	"github.com/insolar/insolar/logicrunner/pulsemanager"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/platformpolicy"
	"github.com/insolar/insolar/virtual/integration/mimic"
)

func NodeLight() insolar.Reference {
	return light.ref
}

const PulseStep insolar.PulseNumber = 10

type Server struct {
	test              *testing.T
	ctx               context.Context
	componentManager  *component.Manager
	stopper           func()
	lock              sync.RWMutex
	clientSender      bus.Sender
	logicRunner       *logicrunner.LogicRunner
	contractRequester *contractrequester.ContractRequester

	pulse        insolar.Pulse
	pulseStorage *pulse.StorageMem
	pulseManager insolar.PulseManager

	mimic mimic.Ledger

	ExternalPubSub, IncomingPubSub *gochannel.GoChannel
}

func DefaultVMConfig() configuration.Configuration {
	cfg := configuration.Configuration{}
	cfg.KeysPath = "testdata/bootstrap_keys.json"
	cfg.Ledger.LightChainLimit = math.MaxInt32
	cfg.LogicRunner = configuration.NewLogicRunner()
	cfg.Bus.ReplyTimeout = 5 * time.Second
	cfg.Log = configuration.NewLog()
	cfg.Log.Level = insolar.DebugLevel.String()
	cfg.Log.Formatter = insolar.TextFormat.String()
	return cfg
}

func checkError(ctx context.Context, err error, message string) {
	if err == nil {
		return
	}
	inslogger.FromContext(ctx).Fatalf("%v: %v", message, err.Error())
}

var verboseWM bool

func init() {
	flag.BoolVar(&verboseWM, "verbose-wm", false, "flag to enable watermill logging")
}

func NewServer(
	t *testing.T,
	ctx context.Context,
	cfg configuration.Configuration,
	receiveCallback func(meta payload.Meta, pl payload.Payload) []payload.Payload,
	mmanager machinesmanager.MachinesManager) (*Server, error) {

	traceID := "main_" + utils.RandTraceID()
	ctx, logger := inslogger.InitNodeLogger(ctx, cfg.Log, traceID, "", "")

	if mmanager == nil {
		mmanager = machinesmanager.NewMachinesManager()
	}

	cm := component.Manager{}

	// Cryptography.
	var (
		KeyProcessor  insolar.KeyProcessor
		CryptoScheme  insolar.PlatformCryptographyScheme
		CryptoService insolar.CryptographyService
		KeyStore      insolar.KeyStore
	)
	{
		var err error
		// Private key storage.
		KeyStore, err = keystore.NewKeyStore(cfg.KeysPath)
		if err != nil {
			return nil, errors.Wrap(err, "failed to load KeyStore")
		}
		// Public key manipulations.
		KeyProcessor = platformpolicy.NewKeyProcessor()
		// Platform cryptography.
		CryptoScheme = platformpolicy.NewPlatformCryptographyScheme()
		// Sign, verify, etc.
		CryptoService = cryptography.NewCryptographyService()
	}

	// Network.
	var (
		NodeNetwork network.NodeNetwork
	)
	{
		NodeNetwork = newNodeNetMock(&virtual)
	}

	// Role calculations.
	var (
		Coordinator *jetcoordinator.Coordinator
		Pulses      *pulse.StorageMem
		Jets        *jet.Store
		Nodes       *node.Storage
	)
	{
		Nodes = node.NewStorage()
		Pulses = pulse.NewStorageMem()
		Jets = jet.NewStore()

		Coordinator = jetcoordinator.NewJetCoordinator(cfg.Ledger.LightChainLimit, virtual.ref)
		Coordinator.PulseCalculator = Pulses
		Coordinator.PulseAccessor = Pulses
		Coordinator.JetAccessor = Jets
		Coordinator.PlatformCryptographyScheme = CryptoScheme
		Coordinator.Nodes = Nodes
	}

	// PulseManager
	var (
		PulseManager *pulsemanager.PulseManager
	)
	{
		PulseManager = pulsemanager.NewPulseManager()
	}

	wmLogger := watermill.LoggerAdapter(watermill.NopLogger{})

	if verboseWM {
		wmLogger = logwatermill.NewWatermillLogAdapter(logger)
	}

	// Communication.
	var (
		ClientBus                      *bus.Bus
		ExternalPubSub, IncomingPubSub *gochannel.GoChannel
	)
	{
		pubsub := gochannel.NewGoChannel(gochannel.Config{}, wmLogger)
		ExternalPubSub = pubsub
		IncomingPubSub = pubsub

		c := jetcoordinator.NewJetCoordinator(cfg.Ledger.LightChainLimit, virtual.ref)
		c.PulseCalculator = Pulses
		c.PulseAccessor = Pulses
		c.JetAccessor = Jets
		c.PlatformCryptographyScheme = CryptoScheme
		c.Nodes = Nodes
		ClientBus = bus.NewBus(cfg.Bus, IncomingPubSub, Pulses, c, CryptoScheme)
	}

	logicRunner, err := logicrunner.NewLogicRunner(&cfg.LogicRunner, IncomingPubSub, ClientBus)
	checkError(ctx, err, "failed to start LogicRunner")

	contractRequester, err := contractrequester.New(
		ClientBus,
		Pulses,
		Coordinator,
		CryptoScheme,
	)
	checkError(ctx, err, "failed to start ContractRequester")

	// TODO: remove this hack in INS-3341
	contractRequester.LR = logicRunner

	cm.Inject(CryptoScheme,
		KeyStore,
		CryptoService,
		KeyProcessor,
		Coordinator,
		logicRunner,

		ClientBus,
		IncomingPubSub,
		contractRequester,
		artifacts.NewClient(ClientBus),
		artifacts.NewDescriptorsCache(),
		Pulses,
		Jets,
		Nodes,

		logicexecutor.NewLogicExecutor(),
		logicrunner.NewRequestsExecutor(),
		mmanager,
		NodeNetwork,
		PulseManager)

	err = cm.Init(ctx)
	checkError(ctx, err, "failed to init components")

	err = cm.Start(ctx)
	checkError(ctx, err, "failed to start components")

	var (
		LedgerMimic mimic.Ledger
	)
	{
		LedgerMimic = mimic.NewMimicLedger(ctx, CryptoScheme, Pulses)
	}

	// Start routers with handlers.
	outHandler := func(msg *message.Message) error {
		var err error

		if msg.Metadata.Get(meta.Type) == meta.TypeReply {
			err = ExternalPubSub.Publish(getIncomingTopic(msg), msg)
			if err != nil {
				panic(errors.Wrap(err, "failed to publish to self"))
			}
			return nil
		}

		msgMeta := payload.Meta{}
		err = msgMeta.Unmarshal(msg.Payload)
		if err != nil {
			panic(errors.Wrap(err, "failed to unmarshal meta"))
		}

		// Republish as incoming to self.
		if msgMeta.Receiver == virtual.ID() {
			err = ExternalPubSub.Publish(getIncomingTopic(msg), msg)
			if err != nil {
				panic(errors.Wrap(err, "failed to publish to self"))
			}
			return nil
		}

		pl, err := payload.Unmarshal(msgMeta.Payload)
		if err != nil {
			panic(errors.Wrap(err, "failed to unmarshal payload"))
		}
		if msgMeta.Receiver == NodeLight() {
			go func() {
				var replies []payload.Payload
				if receiveCallback != nil {
					replies = receiveCallback(msgMeta, pl)
				} else {
					replies = LedgerMimic.ProcessMessage(msgMeta, pl)
				}

				for _, rep := range replies {
					msg, err := payload.NewMessage(rep)
					if err != nil {
						panic(err)
					}
					ClientBus.Reply(context.Background(), msgMeta, msg)
				}
			}()
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

	stopper := startWatermill(
		ctx, wmLogger, IncomingPubSub, ClientBus,
		outHandler,
		logicRunner.FlowDispatcher.Process,
		contractRequester.FlowDispatcher.Process,
	)

	PulseManager.AddDispatcher(logicRunner.FlowDispatcher, contractRequester.FlowDispatcher)

	inslogger.FromContext(ctx).WithFields(map[string]interface{}{
		"light":   light.ID().String(),
		"virtual": virtual.ID().String(),
	}).Info("started test server")

	s := &Server{
		ctx:               ctx,
		contractRequester: contractRequester,
		test:              t,
		pulseStorage:      Pulses,
		pulseManager:      PulseManager,
		componentManager:  &cm,
		stopper:           stopper,
		pulse:             *insolar.GenesisPulse,
		clientSender:      ClientBus,
		mimic:             LedgerMimic,
	}
	return s, nil
}

func (s *Server) Stop(ctx context.Context) {
	panicIfErr(s.componentManager.Stop(ctx))
	s.stopper()
}

func (s *Server) IncrementPulse(ctx context.Context) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.pulse = insolar.Pulse{
		PulseNumber: s.pulse.PulseNumber + PulseStep,
	}
	err := s.pulseManager.Set(ctx, s.pulse)
	if err != nil {
		panic(err)
	}
}

func (s *Server) SendToSelf(ctx context.Context, pl payload.Payload) (<-chan *message.Message, func()) {
	msg, err := payload.NewMessage(pl)
	if err != nil {
		panic(err)
	}
	msg.Metadata.Set(meta.TraceID, s.test.Name())
	return s.clientSender.SendTarget(ctx, msg, virtual.ID())
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

func getIncomingTopic(msg *message.Message) string {
	topic := bus.TopicIncoming
	if msg.Metadata.Get(meta.Type) == meta.TypeReturnResults {
		topic = bus.TopicIncomingRequestResults
	}
	return topic
}

func TestMain(m *testing.M) {
	flag.Parse()
	os.Exit(m.Run())
}
