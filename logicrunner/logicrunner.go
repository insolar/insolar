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

// Package logicrunner - infrastructure for executing smartcontracts
package logicrunner

import (
	"context"
	"strconv"
	"sync"

	watermillMsg "github.com/ThreeDotsLabs/watermill/message"
	"github.com/pkg/errors"
	"go.opencensus.io/trace"

	"github.com/insolar/go-actors/actor/system"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/insolar/flow/dispatcher"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/instrumentation/instracer"
	"github.com/insolar/insolar/logicrunner/artifacts"
	"github.com/insolar/insolar/logicrunner/builtin"
	lrCommon "github.com/insolar/insolar/logicrunner/common"
	"github.com/insolar/insolar/logicrunner/goplugin"
	"github.com/insolar/insolar/logicrunner/machinesmanager"
	"github.com/insolar/insolar/logicrunner/shutdown"
	"github.com/insolar/insolar/logicrunner/writecontroller"
	"github.com/insolar/insolar/network"
)

// LogicRunner is a general interface of contract executor
type LogicRunner struct {
	ContractRequester          insolar.ContractRequester          `inject:""`
	NodeNetwork                network.NodeNetwork                `inject:""`
	PlatformCryptographyScheme insolar.PlatformCryptographyScheme `inject:""`
	PulseAccessor              pulse.Accessor                     `inject:""`
	ArtifactManager            artifacts.Client                   `inject:""`
	DescriptorsCache           artifacts.DescriptorsCache         `inject:""`
	JetCoordinator             jet.Coordinator                    `inject:""`
	RequestsExecutor           RequestsExecutor                   `inject:""`
	MachinesManager            machinesmanager.MachinesManager    `inject:""`
	JetStorage                 jet.Storage                        `inject:""`
	Publisher                  watermillMsg.Publisher
	Sender                     bus.Sender
	SenderWithRetry            *bus.WaitOKSender
	StateStorage               StateStorage
	ResultsMatcher             ResultMatcher
	OutgoingSender             OutgoingRequestSender
	WriteController            writecontroller.WriteController
	FlowDispatcher             dispatcher.Dispatcher
	ShutdownFlag               shutdown.Flag

	Cfg *configuration.LogicRunner

	rpc *lrCommon.RPC
}

// NewLogicRunner is constructor for LogicRunner
func NewLogicRunner(cfg *configuration.LogicRunner, publisher watermillMsg.Publisher, sender bus.Sender) (*LogicRunner, error) {
	if cfg == nil {
		return nil, errors.New("LogicRunner have nil configuration")
	}
	res := LogicRunner{
		Cfg:       cfg,
		Publisher: publisher,
		Sender:    sender,
	}

	res.ResultsMatcher = newResultsMatcher(&res)
	return &res, nil
}

func (lr *LogicRunner) LRI() {}

func (lr *LogicRunner) Init(ctx context.Context) error {
	lr.ShutdownFlag = shutdown.NewFlag()

	as := system.New()
	lr.OutgoingSender = NewOutgoingRequestSender(as, lr.ContractRequester, lr.ArtifactManager)

	lr.StateStorage = NewStateStorage(
		lr.Publisher,
		lr.RequestsExecutor,
		lr.Sender,
		lr.JetCoordinator,
		lr.PulseAccessor,
		lr.ArtifactManager,
		lr.OutgoingSender,
		lr.ShutdownFlag,
	)

	lr.rpc = lrCommon.NewRPC(
		NewRPCMethods(
			lr.ArtifactManager,
			lr.DescriptorsCache,
			lr.ContractRequester,
			lr.StateStorage,
			lr.OutgoingSender,
		),
		lr.Cfg,
	)

	lr.WriteController = writecontroller.NewWriteController()
	err := lr.WriteController.Open(ctx, insolar.FirstPulseNumber)
	if err != nil {
		return errors.Wrap(err, "failed to initialize write controller")
	}

	lr.initHandlers()

	return nil
}

func (lr *LogicRunner) initHandlers() {
	dep := &Dependencies{
		Publisher:        lr.Publisher,
		StateStorage:     lr.StateStorage,
		ResultsMatcher:   lr.ResultsMatcher,
		lr:               lr,
		Sender:           lr.Sender,
		JetStorage:       lr.JetStorage,
		WriteAccessor:    lr.WriteController,
		OutgoingSender:   lr.OutgoingSender,
		RequestsExecutor: lr.RequestsExecutor,
	}

	initHandle := func(msg *watermillMsg.Message) *Init {
		return &Init{
			dep:     dep,
			Message: msg,
		}
	}
	lr.FlowDispatcher = dispatcher.NewDispatcher(lr.PulseAccessor,
		func(msg *watermillMsg.Message) flow.Handle {
			return initHandle(msg).Present
		}, func(msg *watermillMsg.Message) flow.Handle {
			return initHandle(msg).Future
		}, func(msg *watermillMsg.Message) flow.Handle {
			return initHandle(msg).Past
		})
}

func (lr *LogicRunner) initializeBuiltin(_ context.Context) error {
	bi := builtin.NewBuiltIn(
		lr.ArtifactManager,
		NewRPCMethods(lr.ArtifactManager, lr.DescriptorsCache, lr.ContractRequester, lr.StateStorage, lr.OutgoingSender),
	)
	if err := lr.MachinesManager.RegisterExecutor(insolar.MachineTypeBuiltin, bi); err != nil {
		return err
	}

	return nil
}

func (lr *LogicRunner) initializeGoPlugin(ctx context.Context) error {
	logger := inslogger.FromContext(ctx)
	if lr.Cfg.RPCListen == "" {
		logger.Error("Starting goplugin VM with RPC turned off")
	}

	gp, err := goplugin.NewGoPlugin(lr.Cfg, lr.ArtifactManager)
	if err != nil {
		return err
	}

	if err := lr.MachinesManager.RegisterExecutor(insolar.MachineTypeGoPlugin, gp); err != nil {
		return err
	}

	return nil
}

// Start starts logic runner component
func (lr *LogicRunner) Start(ctx context.Context) error {
	if lr.Cfg.BuiltIn != nil {
		if err := lr.initializeBuiltin(ctx); err != nil {
			return errors.Wrap(err, "Failed to initialize builtin VM")
		}
	}

	if lr.Cfg.GoPlugin != nil {
		if err := lr.initializeGoPlugin(ctx); err != nil {
			return errors.Wrap(err, "Failed to initialize goplugin VM")
		}
	}

	if lr.Cfg.RPCListen != "" {
		lr.rpc.Start(ctx)
	}

	lr.ArtifactManager.InjectFinish()

	return nil
}

// Stop stops logic runner component and its executors
func (lr *LogicRunner) Stop(ctx context.Context) error {
	reterr := error(nil)
	if lr.OutgoingSender != nil {
		lr.OutgoingSender.Stop(ctx)
	}
	if err := lr.rpc.Stop(ctx); err != nil {
		return err
	}

	return reterr
}

func (lr *LogicRunner) GracefulStop(ctx context.Context) error {
	waitFunction := lr.ShutdownFlag.Stop(ctx)
	waitFunction()

	return nil
}

func loggerWithTargetID(ctx context.Context, msg insolar.Parcel) context.Context {
	ctx, _ = inslogger.WithField(ctx, "targetid", msg.DefaultTarget().String())
	return ctx
}

func (lr *LogicRunner) OnPulse(ctx context.Context, oldPulse insolar.Pulse, newPulse insolar.Pulse) error {
	ctx, span := instracer.StartSpan(ctx, "pulse.logicrunner")
	defer span.End()

	err := lr.WriteController.CloseAndWait(ctx, oldPulse.PulseNumber)
	if err != nil {
		return errors.Wrap(err, "failed to close pulse on write controller")
	}

	messages := lr.StateStorage.OnPulse(ctx, newPulse)
	err = lr.WriteController.Open(ctx, newPulse.PulseNumber)
	if err != nil {
		return errors.Wrap(err, "failed to start new pulse on write controller")
	}

	if len(messages) > 0 {
		go lr.sendOnPulseMessagesAsync(ctx, messages)
	}

	lr.stopIfNeeded(ctx)

	return nil
}

func (lr *LogicRunner) stopIfNeeded(ctx context.Context) {
	lr.ShutdownFlag.Done(ctx, func() bool {
		return lr.StateStorage.IsEmpty()
	})
}

func (lr *LogicRunner) sendOnPulseMessagesAsync(ctx context.Context, messages map[insolar.Reference][]payload.Payload) {
	ctx, spanMessages := instracer.StartSpan(ctx, "pulse.logicrunner sending messages")
	spanMessages.AddAttributes(trace.StringAttribute("numMessages", strconv.Itoa(len(messages))))

	var sendWg sync.WaitGroup

	for ref, msg := range messages {
		sendWg.Add(len(msg))
		for _, msg := range msg {
			go lr.sendOnPulseMessage(ctx, ref, msg, &sendWg)
		}
	}

	sendWg.Wait()
	spanMessages.End()
}

func (lr *LogicRunner) sendOnPulseMessage(ctx context.Context, objectRef insolar.Reference, payloadObj payload.Payload, sendWg *sync.WaitGroup) {
	defer sendWg.Done()

	msg, err := payload.NewMessage(payloadObj)
	if err != nil {
		inslogger.FromContext(ctx).Error("failed to serialize message: " + err.Error())
		return
	}

	// we dont really care about response, because we are sending this in the beginning of the pulse
	// so flow canceled should not happened, if it does, somebody already restarted
	_, done := lr.Sender.SendRole(ctx, msg, insolar.DynamicRoleVirtualExecutor, objectRef)
	done()
}

func contextWithServiceData(ctx context.Context, data *payload.ServiceData) context.Context {
	// ctx := inslogger.ContextWithTrace(context.Background(), data.LogTraceID)
	ctx = inslogger.ContextWithTrace(ctx, data.LogTraceID)
	ctx = inslogger.WithLoggerLevel(ctx, data.LogLevel)
	if data.TraceSpanData != nil {
		parentSpan := instracer.MustDeserialize(data.TraceSpanData)
		return instracer.WithParentSpan(ctx, parentSpan)
	}
	return ctx
}

func (lr *LogicRunner) AddUnwantedResponse(ctx context.Context, msg insolar.Payload) error {
	m := msg.(*payload.ReturnResults)
	currentPulse, err := lr.PulseAccessor.Latest(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to get current pulse")
	}
	done, err := lr.WriteController.Begin(ctx, currentPulse.PulseNumber)
	defer done()
	if err != nil {
		return flow.ErrCancelled
	}

	// TODO: move towards flow.Dispatcher in INS-3341
	err = lr.isStillExecutor(ctx, *m.Target.Record())
	if err != nil {
		return err
	}

	return lr.ResultsMatcher.AddUnwantedResponse(ctx, m)
}

func (lr *LogicRunner) isStillExecutor(ctx context.Context, object insolar.ID) error {
	currentPulse, err := lr.PulseAccessor.Latest(ctx)
	if err != nil {
		inslogger.FromContext(ctx).Error(errors.Wrap(err, "failed to get current pulse"))
		return flow.ErrCancelled
	}

	node, err := lr.JetCoordinator.VirtualExecutorForObject(ctx, object, currentPulse.PulseNumber)
	if err != nil {
		inslogger.FromContext(ctx).Error(errors.Wrap(err, "failed to calculate current executor"))
		return flow.ErrCancelled
	}

	if *node != lr.JetCoordinator.Me() {
		inslogger.FromContext(ctx).Debug("I'm not executor")
		return flow.ErrCancelled
	}

	return nil
}
