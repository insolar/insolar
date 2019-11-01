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
	"time"

	watermillMsg "github.com/ThreeDotsLabs/watermill/message"
	"github.com/pkg/errors"
	"go.opencensus.io/stats"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/conveyor"
	"github.com/insolar/insolar/conveyor/injector"
	"github.com/insolar/insolar/conveyor/smachine"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/bus"
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
	"github.com/insolar/insolar/logicrunner/logicexecutor"
	"github.com/insolar/insolar/logicrunner/machinesmanager"
	"github.com/insolar/insolar/logicrunner/metrics"
	"github.com/insolar/insolar/logicrunner/requestexecutor"
	"github.com/insolar/insolar/logicrunner/s_artifact"
	"github.com/insolar/insolar/logicrunner/s_contract_requester"
	"github.com/insolar/insolar/logicrunner/s_contract_runner"
	"github.com/insolar/insolar/logicrunner/s_sender"
	"github.com/insolar/insolar/logicrunner/shutdown"
	"github.com/insolar/insolar/logicrunner/sm_object"
	statemachine_go "github.com/insolar/insolar/logicrunner/statemachine"
)

// LogicRunner is a general interface of contract executor
type LogicRunner struct {
	PlatformCryptographyScheme insolar.PlatformCryptographyScheme `inject:""`
	ContractRequester          insolar.ContractRequester          `inject:""`
	PulseAccessor              insolarPulse.Accessor              `inject:""`
	ArtifactManager            artifacts.Client                   `inject:""`
	JetCoordinator             jet.Coordinator                    `inject:""`
	JetStorage                 jet.Storage                        `inject:""`

	LogicExecutor    logicexecutor.LogicExecutor
	DescriptorsCache artifacts.DescriptorsCache
	RequestsExecutor requestexecutor.RequestExecutor
	MachinesManager  machinesmanager.MachinesManager
	Publisher        watermillMsg.Publisher
	Sender           bus.Sender
	SenderWithRetry  *bus.WaitOKSender
	ResultsMatcher   ResultMatcher
	FlowDispatcher   dispatcher.Dispatcher
	ShutdownFlag     shutdown.Flag

	Conveyor                 *conveyor.PulseConveyor
	ObjectCatalog            *sm_object.LocalObjectCatalog
	ArtifactClientService    *s_artifact.ArtifactClientServiceAdapter
	ContractRequesterService *s_contract_requester.ContractRequesterServiceAdapter
	ContractRunnerService    *s_contract_runner.ContractRunnerServiceAdapter
	SenderService            *s_sender.SenderServiceAdapter

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

	return &res, nil
}

func (lr *LogicRunner) LRI() {}

func (lr *LogicRunner) Init(ctx context.Context) error {
	lr.ShutdownFlag = shutdown.NewFlag()
	lr.ResultsMatcher = newResultsMatcher(lr.Sender, lr.PulseAccessor)
	lr.MachinesManager = machinesmanager.NewMachinesManager()
	lr.DescriptorsCache = artifacts.NewDescriptorsCache(lr.ArtifactManager)
	lr.LogicExecutor = logicexecutor.NewLogicExecutor(lr.MachinesManager, lr.DescriptorsCache)
	lr.RequestsExecutor = requestexecutor.NewRequestsExecutor(lr.Sender, lr.LogicExecutor, lr.ArtifactManager, lr.PulseAccessor)

	// configuration steps for slot machine
	machineConfig := smachine.SlotMachineConfig{
		SlotPageSize:    1000,
		PollingPeriod:   500 * time.Millisecond,
		PollingTruncate: 1 * time.Millisecond,
		ScanCountLimit:  100000,
	}

	lr.ObjectCatalog = &sm_object.LocalObjectCatalog{}
	lr.ArtifactClientService = s_artifact.CreateArtifactClientService(lr.Sender)
	lr.ContractRequesterService = s_contract_requester.CreateContractRequesterService(lr.ContractRequester)
	lr.ContractRunnerService = s_contract_runner.CreateContractRunnerService(lr.LogicExecutor, lr.MachinesManager)
	lr.SenderService = s_sender.CreateSenderService(lr.Sender, lr.PulseAccessor)

	di := injector.SyncReadWriteContainer{}
	di.PutDependency("sm_object.LocalObjectCatalog", lr.ObjectCatalog)
	di.PutDependency("s_artifact.ArtifactClientServiceAdapter", lr.ArtifactClientService)
	di.PutDependency("s_contract_requester.ContractRequesterServiceAdapter", lr.ContractRequesterService)
	di.PutDependency("s_contract_runner.ContractRunnerServiceAdapter", lr.ContractRunnerService)
	di.PutDependency("s_sender.SenderServiceAdapter", lr.SenderService)

	defaultHandlers := statemachine_go.DefaultHandlersFactory
	lr.Conveyor = conveyor.NewPulseConveyor(ctx, machineConfig, defaultHandlers, machineConfig, di)

	rpcMethods := NewRPCMethods(lr.DescriptorsCache, lr.Conveyor, lr.PulseAccessor)
	lr.rpc = lrCommon.NewRPC(rpcMethods, lr.Cfg)
	return nil
}

func (lr *LogicRunner) initializeBuiltin(_ context.Context) error {
	rpcMethods := NewRPCMethods(lr.DescriptorsCache, lr.Conveyor, lr.PulseAccessor)
	bi := builtin.NewBuiltIn(lr.ArtifactManager, rpcMethods)
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
	// if lr.OutgoingSender != nil {
	// 	lr.OutgoingSender.Stop(ctx)
	// }
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

	// err := lr.WriteController.CloseAndWait(ctx, oldPulse.PulseNumber)
	// if err != nil {
	// 	return errors.Wrap(err, "failed to close pulse on write controller")
	// }

	lr.ResultsMatcher.Clear(ctx)

	// messages := lr.StateStorage.OnPulse(ctx, newPulse)
	// err = lr.WriteController.Open(ctx, newPulse.PulseNumber)
	// if err != nil {
	// 	return errors.Wrap(err, "failed to start new pulse on write controller")
	// }

	// if len(messages) > 0 {
	// 	go lr.sendOnPulseMessagesAsync(ctx, messages)
	// }

	lr.stopIfNeeded(ctx)

	return nil
}

func (lr *LogicRunner) stopIfNeeded(ctx context.Context) {
	// lr.ShutdownFlag.Done(ctx, func() bool {
	// 	return lr.StateStorage.IsEmpty()
	// })
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
	// done, err := lr.WriteController.Begin(ctx, flow.Pulse(ctx))
	// if err != nil {
	// 	if err == writecontroller.ErrWriteClosed {
	// 		return flow.ErrCancelled
	// 	}
	// 	return errors.Wrap(err, "couldn't obtain writecontroller lock")
	// }
	// defer done()

	lr.ResultsMatcher.AddUnwantedResponse(ctx, *m)
	return nil
}
