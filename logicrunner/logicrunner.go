// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

// Package logicrunner - infrastructure for executing smartcontracts
package logicrunner

import (
	"context"
	"strconv"
	"sync"
	"time"

	watermillMsg "github.com/ThreeDotsLabs/watermill/message"
	"github.com/pkg/errors"
	"go.opencensus.io/stats"

	"github.com/insolar/go-actors/actor/system"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/insolar/flow/dispatcher"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/payload"
	insolarPulse "github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/instrumentation/instracer"
	"github.com/insolar/insolar/logicrunner/artifacts"
	"github.com/insolar/insolar/logicrunner/builtin"
	lrCommon "github.com/insolar/insolar/logicrunner/common"
	"github.com/insolar/insolar/logicrunner/goplugin"
	"github.com/insolar/insolar/logicrunner/machinesmanager"
	"github.com/insolar/insolar/logicrunner/metrics"
	"github.com/insolar/insolar/logicrunner/shutdown"
	"github.com/insolar/insolar/logicrunner/writecontroller"
	"github.com/insolar/insolar/pulse"
)

// LogicRunner is a general interface of contract executor
type LogicRunner struct {
	ContractRequester          insolar.ContractRequester          `inject:""`
	PlatformCryptographyScheme insolar.PlatformCryptographyScheme `inject:""`
	PulseAccessor              insolarPulse.Accessor              `inject:""`
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

	builtinContracts builtin.BuiltinContracts
}

// NewLogicRunner is constructor for LogicRunner
func NewLogicRunner(
	cfg *configuration.LogicRunner, publisher watermillMsg.Publisher, sender bus.Sender, builtinContracts builtin.BuiltinContracts,
) (*LogicRunner, error) {
	if cfg == nil {
		return nil, errors.New("LogicRunner have nil configuration")
	}
	res := LogicRunner{
		Cfg:              cfg,
		Publisher:        publisher,
		Sender:           sender,
		builtinContracts: builtinContracts,
	}

	return &res, nil
}

func (lr *LogicRunner) LRI() {}

func (lr *LogicRunner) Init(ctx context.Context) error {
	lr.ShutdownFlag = shutdown.NewFlag()

	as := system.New()
	lr.OutgoingSender = NewOutgoingRequestSender(as, lr.ContractRequester, lr.ArtifactManager, lr.PulseAccessor)

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
	lr.ResultsMatcher = newResultsMatcher(lr.Sender, lr.PulseAccessor)

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
	err := lr.WriteController.Open(ctx, pulse.MinTimePulse)
	if err != nil {
		return errors.Wrap(err, "failed to initialize write controller")
	}

	lr.initHandlers()

	return nil
}

func (lr *LogicRunner) initHandlers() {
	dep := &Dependencies{
		ArtifactManager:  lr.ArtifactManager,
		Publisher:        lr.Publisher,
		StateStorage:     lr.StateStorage,
		ResultsMatcher:   lr.ResultsMatcher,
		Sender:           lr.Sender,
		JetStorage:       lr.JetStorage,
		JetCoordinator:   lr.JetCoordinator,
		WriteAccessor:    lr.WriteController,
		OutgoingSender:   lr.OutgoingSender,
		RequestsExecutor: lr.RequestsExecutor,
		PulseAccessor:    lr.PulseAccessor,
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
		lr.builtinContracts,
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

func (lr *LogicRunner) OnPulse(ctx context.Context, oldPulse insolar.Pulse, newPulse insolar.Pulse) error {
	onPulseStart := time.Now()
	ctx, span := instracer.StartSpan(ctx, "pulse.logicrunner")
	defer func(ctx context.Context) {
		stats.Record(ctx,
			metrics.LogicRunnerOnPulseTiming.M(float64(time.Since(onPulseStart).Nanoseconds())/1e6))
		span.Finish()
	}(ctx)

	err := lr.WriteController.CloseAndWait(ctx, oldPulse.PulseNumber)
	if err != nil {
		return errors.Wrap(err, "failed to close pulse on write controller")
	}

	lr.ResultsMatcher.Clear(ctx)

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
	spanMessages.SetTag("numMessages", strconv.Itoa(len(messages)))

	var sendWg sync.WaitGroup

	for ref, msg := range messages {
		sendWg.Add(len(msg))
		for _, msg := range msg {
			go lr.sendOnPulseMessage(ctx, ref, msg, &sendWg)
		}
	}

	sendWg.Wait()
	spanMessages.Finish()
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
	done, err := lr.WriteController.Begin(ctx, flow.Pulse(ctx))
	if err != nil {
		if err == writecontroller.ErrWriteClosed {
			return flow.ErrCancelled
		}
		return errors.Wrap(err, "couldn't obtain writecontroller lock")
	}
	defer done()

	lr.ResultsMatcher.AddUnwantedResponse(ctx, *m)
	return nil
}
