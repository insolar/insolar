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

	"github.com/pkg/errors"
	"go.opencensus.io/stats"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/conveyor"
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
	"github.com/insolar/insolar/logicrunner/s_artifact"
	"github.com/insolar/insolar/logicrunner/s_contract_requester"
	"github.com/insolar/insolar/logicrunner/s_contract_runner"
	"github.com/insolar/insolar/logicrunner/s_jet_storage"
	"github.com/insolar/insolar/logicrunner/s_sender"
	"github.com/insolar/insolar/logicrunner/shutdown"
	"github.com/insolar/insolar/logicrunner/sm_object"
	"github.com/insolar/insolar/logicrunner/statemachine"
	"github.com/insolar/insolar/logicrunner_old/requestexecutor"
)

// LogicRunner is a general interface of contract executor
type LogicRunner struct {
	ContractRequester insolar.ContractRequester `inject:""`
	PulseAccessor     insolarPulse.Accessor     `inject:""`
	ArtifactManager   artifacts.Client          `inject:""`
	JetStorage        jet.Storage               `inject:""`

	LogicExecutor    logicexecutor.LogicExecutor
	DescriptorsCache artifacts.DescriptorsCache
	RequestsExecutor requestexecutor.RequestExecutor
	MachinesManager  machinesmanager.MachinesManager
	Sender           bus.Sender
	FlowDispatcher   dispatcher.Dispatcher
	ShutdownFlag     shutdown.Flag
	ContractRunner   s_contract_runner.ContractRunnerService

	Conveyor       *conveyor.PulseConveyor
	ConveyorWorker lrCommon.ConveyorWorker

	ObjectCatalog            *sm_object.LocalObjectCatalog
	ArtifactClientService    *s_artifact.ArtifactClientServiceAdapter
	ContractRequesterService *s_contract_requester.ContractRequesterServiceAdapter
	ContractRunnerService    *s_contract_runner.ContractRunnerServiceAdapter
	SenderService            *s_sender.SenderServiceAdapter
	JetStorageService        *s_jet_storage.JetStorageServiceAdapter

	Cfg *configuration.LogicRunner

	rpc *lrCommon.RPC
}

// NewLogicRunner is constructor for LogicRunner
func NewLogicRunner(cfg *configuration.LogicRunner, sender bus.Sender) (*LogicRunner, error) {
	if cfg == nil {
		return nil, errors.New("LogicRunner have nil configuration")
	}
	res := LogicRunner{
		Cfg:    cfg,
		Sender: sender,
	}

	return &res, nil
}

func (lr *LogicRunner) LRI() {}

func (lr *LogicRunner) Init(ctx context.Context) error {
	lr.ShutdownFlag = shutdown.NewFlag()
	lr.MachinesManager = machinesmanager.NewMachinesManager()
	lr.DescriptorsCache = artifacts.NewDescriptorsCache(lr.ArtifactManager)
	lr.LogicExecutor = logicexecutor.NewLogicExecutor(lr.MachinesManager, lr.DescriptorsCache)
	lr.RequestsExecutor = requestexecutor.NewRequestsExecutor(lr.Sender, lr.LogicExecutor, lr.ArtifactManager, lr.PulseAccessor)
	lr.ContractRunner = s_contract_runner.CreateContractRunner(lr.LogicExecutor, lr.MachinesManager, lr.ArtifactManager)
	lr.rpc = lrCommon.NewRPC(lr.ContractRunner, lr.Cfg)

	// configuration steps for slot machine
	machineConfig := smachine.SlotMachineConfig{
		PollingPeriod:     500 * time.Millisecond,
		PollingTruncate:   1 * time.Millisecond,
		SlotPageSize:      1000,
		ScanCountLimit:    100000,
		SlotMachineLogger: statemachine.ConveyorLoggerFactory{},
	}

	lr.ObjectCatalog = &sm_object.LocalObjectCatalog{}
	lr.ArtifactClientService = s_artifact.CreateArtifactClientService(lr.ArtifactManager)
	lr.ContractRequesterService = s_contract_requester.CreateContractRequesterService(lr.ContractRequester)
	lr.ContractRunnerService = s_contract_runner.CreateContractRunnerService(lr.ContractRunner)
	lr.SenderService = s_sender.CreateSenderService(lr.Sender, lr.PulseAccessor)
	lr.JetStorageService = s_jet_storage.CreateJetStorageService(lr.JetStorage)

	defaultHandlers := statemachine.DefaultHandlersFactory

	lr.Conveyor = conveyor.NewPulseConveyor(context.Background(), conveyor.PulseConveyorConfig{
		ConveyorMachineConfig: machineConfig,
		SlotMachineConfig:     machineConfig,
		EventlessSleep:        100 * time.Millisecond,
		MinCachePulseAge:      100,
		MaxPastPulseAge:       1000,
	}, defaultHandlers, nil)

	lr.Conveyor.AddDependency(lr.ObjectCatalog)
	lr.Conveyor.AddDependency(lr.ArtifactClientService)
	lr.Conveyor.AddDependency(lr.ContractRequesterService)
	lr.Conveyor.AddDependency(lr.ContractRunnerService)
	lr.Conveyor.AddDependency(lr.SenderService)
	lr.Conveyor.AddDependency(lr.JetStorageService)

	lr.ConveyorWorker = lrCommon.NewConveyorWorker()
	lr.ConveyorWorker.AttachTo(lr.Conveyor)

	lr.FlowDispatcher = lrCommon.NewConveyorDispatcher(lr.Conveyor)

	return nil
}

// Start starts logic runner component
func (lr *LogicRunner) Start(ctx context.Context) error {
	if lr.Cfg.RPCListen != "" {
		lr.rpc.Start(ctx)
	}

	if lr.Cfg.BuiltIn != nil {
		bi := builtin.NewBuiltIn(lr.ArtifactManager, lr.ContractRunner)

		err := lr.MachinesManager.RegisterExecutor(insolar.MachineTypeBuiltin, bi)
		if err != nil {
			return err
		}
	}

	if lr.Cfg.GoPlugin != nil {
		gp := goplugin.NewGoPlugin(lr.Cfg)

		err := lr.MachinesManager.RegisterExecutor(insolar.MachineTypeGoPlugin, gp)
		if err != nil {
			return err
		}
	}

	lr.ArtifactManager.InjectFinish()

	return nil
}

// Stop stops logic runner component and its executors
func (lr *LogicRunner) Stop(ctx context.Context) error {
	lr.ConveyorWorker.Stop()

	return lr.rpc.Stop(ctx)
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

	lr.stopIfNeeded(ctx)

	return nil
}

func (lr *LogicRunner) stopIfNeeded(ctx context.Context) {
	lr.ShutdownFlag.Done(ctx, func() bool { return true })
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
	ctx = inslogger.ContextWithTrace(ctx, data.LogTraceID)
	ctx = inslogger.WithLoggerLevel(ctx, data.LogLevel)
	if data.TraceSpanData != nil {
		parentSpan := instracer.MustDeserialize(data.TraceSpanData)
		return instracer.WithParentSpan(ctx, parentSpan)
	}
	return ctx
}

func (lr *LogicRunner) AddUnwantedResponse(ctx context.Context, msg insolar.Payload) error {
	return nil
}
