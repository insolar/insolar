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
	"github.com/insolar/insolar/network"
	"strconv"
	"sync"

	"github.com/insolar/go-actors/actor/system"

	watermillMsg "github.com/ThreeDotsLabs/watermill/message"
	"github.com/pkg/errors"
	"go.opencensus.io/trace"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/insolar/flow/dispatcher"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/message"
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/instrumentation/instracer"
	"github.com/insolar/insolar/logicrunner/artifacts"
	"github.com/insolar/insolar/logicrunner/builtin"
	lrCommon "github.com/insolar/insolar/logicrunner/common"
	"github.com/insolar/insolar/logicrunner/goplugin"
	"github.com/insolar/insolar/logicrunner/handles"
	"github.com/insolar/insolar/logicrunner/machinesmanager"
	"github.com/insolar/insolar/logicrunner/outgoingsender"
	"github.com/insolar/insolar/logicrunner/requestsexecutor"
	"github.com/insolar/insolar/logicrunner/resultmatcher"
	"github.com/insolar/insolar/logicrunner/statestorage"
	"github.com/insolar/insolar/logicrunner/writecontroller"
)

// LogicRunner is a general interface of contract executor
type LogicRunner struct {
	MessageBus                 insolar.MessageBus                 `inject:""`
	ContractRequester          insolar.ContractRequester          `inject:""`
	NodeNetwork                network.NodeNetwork                `inject:""`
	PlatformCryptographyScheme insolar.PlatformCryptographyScheme `inject:""`
	ParcelFactory              message.ParcelFactory              `inject:""`
	PulseAccessor              pulse.Accessor                     `inject:""`
	ArtifactManager            artifacts.Client                   `inject:""`
	DescriptorsCache           artifacts.DescriptorsCache         `inject:""`
	JetCoordinator             jet.Coordinator                    `inject:""`
	RequestsExecutor           requestsexecutor.RequestsExecutor  `inject:""`
	MachinesManager            machinesmanager.MachinesManager    `inject:""`
	JetStorage                 jet.Storage                        `inject:""`
	Publisher                  watermillMsg.Publisher
	Sender                     bus.Sender
	SenderWithRetry            *bus.WaitOKSender
	StateStorage               statestorage.StateStorage
	ResultsMatcher             resultmatcher.ResultMatcher
	OutgoingSender             outgoingsender.OutgoingRequestSender
	WriteController            writecontroller.WriteController
	FlowDispatcher             dispatcher.Dispatcher

	Cfg *configuration.LogicRunner

	rpc *lrCommon.RPC

	stopLock   sync.Mutex
	isStopping bool
	stopChan   chan struct{}
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
	lr.ResultsMatcher = resultmatcher.NewResultsMatcher(lr.MessageBus, lr.PulseAccessor, lr.JetCoordinator)

	as := system.New()
	lr.OutgoingSender = outgoingsender.NewOutgoingRequestSender(as, lr.ContractRequester, lr.ArtifactManager)

	lr.StateStorage = statestorage.NewStateStorage(
		lr.Publisher,
		lr.RequestsExecutor,
		lr.MessageBus,
		lr.JetCoordinator,
		lr.PulseAccessor,
		lr.ArtifactManager,
		lr.OutgoingSender,
	)

	lr.SenderWithRetry = bus.NewWaitOKWithRetrySender(lr.Sender, lr.PulseAccessor, 3)

	lr.rpc = lrCommon.NewRPC(
		NewRPCMethods(lr.ArtifactManager, lr.DescriptorsCache, lr.ContractRequester, lr.StateStorage, lr.OutgoingSender),
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
	dep := &handles.Dependencies{
		Publisher:         lr.Publisher,
		StateStorage:      lr.StateStorage,
		ResultsMatcher:    lr.ResultsMatcher,
		JetCoordinator:    lr.JetCoordinator,
		ArtifactManager:   lr.ArtifactManager,
		ContractRequester: lr.ContractRequester,
		MessageBus:        lr.MessageBus,
		Sender:            lr.Sender,
		JetStorage:        lr.JetStorage,
		WriteAccessor:     lr.WriteController,
		OutgoingSender:    lr.OutgoingSender,
		RequestsExecutor:  lr.RequestsExecutor,
	}

	initHandle := func(msg *watermillMsg.Message) *handles.Init {
		return &handles.Init{
			Dep:     dep,
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

	gp, err := goplugin.NewGoPlugin(lr.Cfg, lr.MessageBus, lr.ArtifactManager)
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
	if err := lr.rpc.Stop(ctx); err != nil {
		return err
	}

	return reterr
}

func (lr *LogicRunner) GracefulStop(ctx context.Context) error {
	inslogger.FromContext(ctx).Debug("LogicRunner.GracefulStop starts ...")

	lr.stopLock.Lock()
	if !lr.isStopping {
		lr.isStopping = true
		lr.stopChan = make(chan struct{}, 1)
	}
	lr.stopLock.Unlock()

	inslogger.FromContext(ctx).Debug("LogicRunner.GracefulStop wait ...")
	<-lr.stopChan
	inslogger.FromContext(ctx).Debug("LogicRunner.GracefulStop ends ...")
	return nil
}

func (lr *LogicRunner) OnPulse(ctx context.Context, oldPulse insolar.Pulse, newPulse insolar.Pulse) error {
	ctx, span := instracer.StartSpan(ctx, "pulse.logicrunner")
	defer span.End()

	err := lr.WriteController.CloseAndWait(ctx, oldPulse.PulseNumber)
	if err != nil {
		return errors.Wrap(err, "failed to close pulse on write controller")
	}

	lr.ResultsMatcher.Clear()

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
	// lock is required to access LogicRunner.state
	lr.StateStorage.Lock()
	defer lr.StateStorage.Unlock()

	if lr.StateStorage.IsEmpty() {
		lr.stopLock.Lock()
		if lr.isStopping {
			inslogger.FromContext(ctx).Debug("LogicRunner ready to stop")
			lr.stopChan <- struct{}{}
		}
		lr.stopLock.Unlock()
	}
}

func (lr *LogicRunner) sendOnPulseMessagesAsync(ctx context.Context, messages []insolar.Message) {
	ctx, spanMessages := instracer.StartSpan(ctx, "pulse.logicrunner sending messages")
	spanMessages.AddAttributes(trace.StringAttribute("numMessages", strconv.Itoa(len(messages))))

	var sendWg sync.WaitGroup
	sendWg.Add(len(messages))

	for _, msg := range messages {
		go lr.sendOnPulseMessage(ctx, msg, &sendWg)
	}

	sendWg.Wait()
	spanMessages.End()
}

func (lr *LogicRunner) sendOnPulseMessage(ctx context.Context, msg insolar.Message, sendWg *sync.WaitGroup) {
	defer sendWg.Done()
	_, err := lr.MessageBus.Send(ctx, msg, nil)
	if err != nil {
		inslogger.FromContext(ctx).Error(errors.Wrap(err, "error while sending validation data on pulse"))
	}
}

func (lr *LogicRunner) AddUnwantedResponse(ctx context.Context, msg insolar.Message) error {
	m := msg.(*message.ReturnResults)
	return lr.ResultsMatcher.AddUnwantedResponse(ctx, m)
}
