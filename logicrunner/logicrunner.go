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

	"github.com/pkg/errors"
	"go.opencensus.io/trace"

	"github.com/ThreeDotsLabs/watermill"
	watermillMsg "github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/infrastructure/gochannel"

	"github.com/insolar/insolar/log"

	"github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/insolar/flow/dispatcher"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/message"
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/insolar/reply"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/instrumentation/instracer"
	"github.com/insolar/insolar/logicrunner/artifacts"
	"github.com/insolar/insolar/logicrunner/builtin"
	lrCommon "github.com/insolar/insolar/logicrunner/common"
	"github.com/insolar/insolar/logicrunner/goplugin"
)

const maxQueueLength = 10

type Ref = insolar.Reference

// Context of one contract execution
type ObjectState struct {
	sync.Mutex

	ExecutionState *ExecutionState
	Validation     *ExecutionState
}

func (st *ObjectState) GetModeState(mode insolar.CallMode) (rv *ExecutionState, err error) {
	switch mode {
	case insolar.ExecuteCallMode:
		rv = st.ExecutionState
	case insolar.ValidateCallMode:
		rv = st.Validation
	default:
		err = errors.Errorf("'%d' is unknown object processing mode", mode)
	}

	if rv == nil && err != nil {
		err = errors.Errorf("object is not in '%s' mode", mode)
	}
	return rv, err
}

func (st *ObjectState) MustModeState(mode insolar.CallMode) *ExecutionState {
	res, err := st.GetModeState(mode)
	if err != nil {
		panic(err)
	}
	if res.CurrentList.Empty() {
		panic("object " + res.Ref.String() + " has no Current")
	}
	return res
}

func makeWMMessage(ctx context.Context, payLoad watermillMsg.Payload, msgType string) *watermillMsg.Message {
	wmMsg := watermillMsg.NewMessage(watermill.NewUUID(), payLoad)
	wmMsg.Metadata.Set(bus.MetaTraceID, inslogger.TraceID(ctx))

	sp, err := instracer.Serialize(ctx)
	if err == nil {
		wmMsg.Metadata.Set(bus.MetaSpanData, string(sp))
	} else {
		inslogger.FromContext(ctx).Error(err)
	}

	wmMsg.Metadata.Set(bus.MetaType, msgType)
	return wmMsg
}

// LogicRunner is a general interface of contract executor
type LogicRunner struct {
	MessageBus                 insolar.MessageBus                 `inject:""`
	ContractRequester          insolar.ContractRequester          `inject:""`
	NodeNetwork                insolar.NodeNetwork                `inject:""`
	PlatformCryptographyScheme insolar.PlatformCryptographyScheme `inject:""`
	ParcelFactory              message.ParcelFactory              `inject:""`
	PulseAccessor              pulse.Accessor                     `inject:""`
	ArtifactManager            artifacts.Client                   `inject:""`
	DescriptorsCache           artifacts.DescriptorsCache         `inject:""`
	JetCoordinator             jet.Coordinator                    `inject:""`
	LogicExecutor              LogicExecutor                      `inject:""`
	MachinesManager            MachinesManager                    `inject:""`

	Cfg *configuration.LogicRunner

	state      map[Ref]*ObjectState // if object exists, we are validating or executing it right now
	stateMutex sync.RWMutex

	rpc *lrCommon.RPC

	stopLock   sync.Mutex
	isStopping bool
	stopChan   chan struct{}

	// Inner dispatcher will be merged with FlowDispatcher after
	// complete migration to watermill.
	FlowDispatcher      *dispatcher.Dispatcher
	innerFlowDispatcher *dispatcher.Dispatcher
	publisher           watermillMsg.Publisher
	router              *watermillMsg.Router
}

// NewLogicRunner is constructor for LogicRunner
func NewLogicRunner(cfg *configuration.LogicRunner) (*LogicRunner, error) {
	if cfg == nil {
		return nil, errors.New("LogicRunner have nil configuration")
	}
	res := LogicRunner{
		Cfg:   cfg,
		state: make(map[Ref]*ObjectState),
	}
	res.rpc = lrCommon.NewRPC(NewRPCMethods(&res), cfg)

	return &res, nil
}

func (lr *LogicRunner) LRI() {}

func InitHandlers(lr *LogicRunner, s bus.Sender) (*watermillMsg.Router, error) {
	wmLogger := log.NewWatermillLogAdapter(inslogger.FromContext(context.Background()))
	pubSub := gochannel.NewGoChannel(gochannel.Config{}, wmLogger)

	dep := &Dependencies{
		Publisher: pubSub,
		lr:        lr,
		Sender:    s,
	}

	initHandle := func(msg *watermillMsg.Message) *Init {
		return &Init{
			dep:     dep,
			Message: msg,
		}
	}
	lr.FlowDispatcher = dispatcher.NewDispatcher(func(msg *watermillMsg.Message) flow.Handle {
		return initHandle(msg).Present
	}, func(msg *watermillMsg.Message) flow.Handle {
		return initHandle(msg).Future
	})

	innerInitHandle := func(msg *watermillMsg.Message) *InnerInit {
		return &InnerInit{
			dep:     dep,
			Message: msg,
		}
	}

	lr.innerFlowDispatcher = dispatcher.NewDispatcher(func(msg *watermillMsg.Message) flow.Handle {
		return innerInitHandle(msg).Present
	}, func(msg *watermillMsg.Message) flow.Handle {
		return innerInitHandle(msg).Present
	})

	router, err := watermillMsg.NewRouter(watermillMsg.RouterConfig{}, wmLogger)
	if err != nil {
		return nil, errors.Wrap(err, "Error while creating new watermill router")
	}

	router.AddNoPublisherHandler(
		"InnerMsgHandler",
		InnerMsgTopic,
		pubSub,
		lr.innerFlowDispatcher.InnerSubscriber,
	)
	go func() {
		if err := router.Run(); err != nil {
			ctx := context.Background()
			inslogger.FromContext(ctx).Error("Error while running router", err)
		}
	}()
	<-router.Running()

	lr.router = router
	lr.publisher = pubSub

	return router, nil
}

func (lr *LogicRunner) initializeBuiltin(_ context.Context) error {
	bi := builtin.NewBuiltIn(lr.ArtifactManager, NewRPCMethods(lr))
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
	if err := lr.router.Close(); err != nil {
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

func (lr *LogicRunner) CheckOurRole(ctx context.Context, msg insolar.Message, role insolar.DynamicRole) error {
	// TODO do map of supported objects for pulse, go to jetCoordinator only if map is empty for ref
	target := msg.DefaultTarget()
	isAuthorized, err := lr.JetCoordinator.IsAuthorized(
		ctx, role, *target.Record(), lr.pulse(ctx).PulseNumber, lr.JetCoordinator.Me(),
	)
	if err != nil {
		return errors.Wrap(err, "authorization failed with error")
	}
	if !isAuthorized {
		return errors.New("can't executeAndReply this object")
	}
	return nil
}

func loggerWithTargetID(ctx context.Context, msg insolar.Parcel) context.Context {
	ctx, _ = inslogger.WithField(ctx, "targetid", msg.DefaultTarget().String())
	return ctx
}

// values here (boolean flags) are inverted here, since it's common "predicate" checking function
func noLoopCheckerPredicate(current *Transcript, args interface{}) bool {
	apiReqID := args.(string)
	if current.Request.ReturnMode == record.ReturnNoWait ||
		current.Request.APIRequestID != apiReqID {
		return true
	}
	return false
}

func (lr *LogicRunner) CheckExecutionLoop(
	ctx context.Context, es *ExecutionState, parcel insolar.Parcel) bool {
	if es.CurrentList.Empty() {
		return false
	}

	msg, ok := parcel.Message().(*message.CallMethod)
	if ok && msg.ReturnMode == record.ReturnNoWait {
		return false
	}

	if es.CurrentList.Check(noLoopCheckerPredicate, msg.APIRequestID) {
		return false
	}

	inslogger.FromContext(ctx).Error("loop detected")
	return true
}

// finishPendingIfNeeded checks whether last execution was a pending one.
// If this is true as a side effect the function sends a PendingFinished
// message to the current executor
func (lr *LogicRunner) finishPendingIfNeeded(ctx context.Context, es *ExecutionState) {
	es.Lock()
	defer es.Unlock()

	if es.pending != message.InPending {
		return
	}

	es.pending = message.NotPending
	es.PendingConfirmed = false

	pulseObj := lr.pulse(ctx)
	meCurrent, _ := lr.JetCoordinator.IsAuthorized(
		ctx, insolar.DynamicRoleVirtualExecutor, *es.Ref.Record(), pulseObj.PulseNumber, lr.JetCoordinator.Me(),
	)
	if !meCurrent {
		go func() {
			msg := message.PendingFinished{Reference: es.Ref}
			_, err := lr.MessageBus.Send(ctx, &msg, nil)
			if err != nil {
				inslogger.FromContext(ctx).Error("Unable to send PendingFinished message:", err)
			}
		}()
	}
}

func (lr *LogicRunner) sendRequestReply(
	ctx context.Context, current *Transcript, re insolar.Reply, errstr string,
) {
	if current.Request.ReturnMode != record.ReturnResult {
		return
	}

	target := *current.RequesterNode
	request := *current.RequestRef
	seq := current.Request.Sequence

	go func() {
		inslogger.FromContext(ctx).Debugf("Sending Method Results for %#v", request)

		_, err := lr.MessageBus.Send(
			ctx,
			&message.ReturnResults{
				Caller:   lr.NodeNetwork.GetOrigin().ID(),
				Target:   target,
				Sequence: seq,
				Reply:    re,
				Error:    errstr,
			},
			&insolar.MessageSendOptions{
				Receiver: &target,
			},
		)
		if err != nil {
			inslogger.FromContext(ctx).Error("couldn't deliver results: ", err)
		}
	}()
}

func (lr *LogicRunner) executeLogic(ctx context.Context, current *Transcript) (insolar.Reply, error) {
	ctx, span := instracer.StartSpan(ctx, "LogicRunner.executeLogic")
	defer span.End()

	if current.Request.CallType == record.CTMethod {
		objDesc, err := lr.ArtifactManager.GetObject(ctx, *current.Request.Object)
		if err != nil {
			return nil, errors.Wrap(err, "couldn't get object")
		}
		current.ObjectDescriptor = objDesc
	}

	res, err := lr.LogicExecutor.Execute(ctx, current)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't executeAndReply request")
	}

	am := lr.ArtifactManager
	request := current.Request

	switch {
	case res.Activation:
		err = lr.ArtifactManager.ActivateObject(
			ctx, *current.RequestRef, *request.Base, *request.Prototype,
			request.CallType == record.CTSaveAsDelegate,
			res.NewMemory,
		)
		if err != nil {
			return nil, errors.Wrap(err, "couldn't activate object")
		}
		return &reply.CallConstructor{Object: current.RequestRef}, err
	case res.Deactivation:
		err := am.DeactivateObject(
			ctx, *current.RequestRef, current.ObjectDescriptor, res.Result,
		)
		if err != nil {
			return nil, errors.Wrap(err, "couldn't deactivate object")
		}
		return &reply.CallMethod{Result: res.Result}, nil
	case res.NewMemory != nil:
		err := am.UpdateObject(
			ctx, *current.RequestRef, current.ObjectDescriptor, res.NewMemory, res.Result,
		)
		if err != nil {
			return nil, errors.Wrap(err, "couldn't update object")
		}
		return &reply.CallMethod{Result: res.Result}, nil
	default:
		_, err = am.RegisterResult(ctx, *request.Object, *current.RequestRef, res.Result)
		if err != nil {
			return nil, errors.Wrap(err, "couldn't save results")
		}
		return &reply.CallMethod{Result: res.Result}, nil
	}
}

func (lr *LogicRunner) startGetLedgerPendingRequest(ctx context.Context, es *ExecutionState) {
	err := lr.publisher.Publish(InnerMsgTopic, makeWMMessage(ctx, es.Ref.Bytes(), getLedgerPendingRequestMsg))
	if err != nil {
		inslogger.FromContext(ctx).Warnf("can't send getLedgerPendingRequestMsg: ", err)
	}
}

func (lr *LogicRunner) OnPulse(ctx context.Context, pulse insolar.Pulse) error {
	lr.stateMutex.Lock()

	lr.FlowDispatcher.ChangePulse(ctx, pulse)
	lr.innerFlowDispatcher.ChangePulse(ctx, pulse)

	ctx, span := instracer.StartSpan(ctx, "pulse.logicrunner")
	defer span.End()

	messages := make([]insolar.Message, 0)

	for ref, state := range lr.state {
		meNext, _ := lr.JetCoordinator.IsAuthorized(
			ctx, insolar.DynamicRoleVirtualExecutor, *ref.Record(), pulse.PulseNumber, lr.JetCoordinator.Me(),
		)
		state.Lock()

		if es := state.ExecutionState; es != nil {
			es.Lock()

			toSend := es.OnPulse(ctx, meNext)
			messages = append(messages, toSend...)

			if !meNext {
				if es.CurrentList.Empty() {
					state.ExecutionState = nil
				}
			} else {
				if es.pending == message.NotPending && es.LedgerHasMoreRequests {
					lr.startGetLedgerPendingRequest(ctx, es)
				}
				if es.pending == message.NotPending {
					es.Broker.StartProcessorIfNeeded(ctx)
				}
			}

			es.Unlock()
		}

		if state.ExecutionState == nil && state.Validation == nil {
			delete(lr.state, ref)
		}

		state.Unlock()
	}

	lr.stateMutex.Unlock()

	if len(messages) > 0 {
		go lr.sendOnPulseMessagesAsync(ctx, messages)
	}

	lr.stopIfNeeded(ctx)

	return nil
}

func (lr *LogicRunner) stopIfNeeded(ctx context.Context) {
	// lock is required to access LogicRunner.state
	lr.stateMutex.Lock()
	defer lr.stateMutex.Unlock()

	if len(lr.state) == 0 {
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

func convertQueueToMessageQueue(ctx context.Context, queue []*Transcript) []message.ExecutionQueueElement {
	mq := make([]message.ExecutionQueueElement, 0)
	var traces string
	for _, elem := range queue {
		mq = append(mq, message.ExecutionQueueElement{
			Parcel:  elem.Parcel,
			Request: elem.RequestRef,
		})

		traces += inslogger.TraceID(elem.Context) + ", "
	}

	inslogger.FromContext(ctx).Debug("convertQueueToMessageQueue: ", traces)

	return mq
}
