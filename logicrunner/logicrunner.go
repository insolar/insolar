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
	"bytes"
	"context"
	"net"
	"strconv"
	"sync"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	watermillMsg "github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/infrastructure/gochannel"

	"github.com/insolar/insolar/log"

	wmBus "github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/insolar/flow/bus"
	"github.com/insolar/insolar/insolar/flow/dispatcher"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/record"

	"go.opencensus.io/trace"

	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/logicrunner/artifacts"

	"github.com/insolar/insolar/instrumentation/instracer"

	"github.com/pkg/errors"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/message"
	"github.com/insolar/insolar/insolar/reply"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/logicrunner/builtin"
	"github.com/insolar/insolar/logicrunner/goplugin"
)

const maxQueueLength = 10

type Ref = insolar.Reference

// Context of one contract execution
type ObjectState struct {
	sync.Mutex

	ExecutionState *ExecutionState
	Validation     *ExecutionState
	Consensus      *Consensus
}

type CurrentExecution struct {
	Context       context.Context
	LogicContext  *insolar.LogicCallContext
	RequestRef    *Ref
	Request       *record.Request
	RequesterNode *Ref
	SentResult    bool
}

type ExecutionQueueElement struct {
	ctx        context.Context
	parcel     insolar.Parcel
	request    *Ref
	fromLedger bool
}

type Error struct {
	Err      error
	Request  *Ref
	Contract *Ref
	Method   string
}

func (lre Error) Error() string {
	var buffer bytes.Buffer

	buffer.WriteString(lre.Err.Error())
	if lre.Contract != nil {
		buffer.WriteString(" Contract=" + lre.Contract.String())
	}
	if lre.Method != "" {
		buffer.WriteString(" Method=" + lre.Method)
	}
	if lre.Request != nil {
		buffer.WriteString(" Request=" + lre.Request.String())
	}

	return buffer.String()
}

func (st *ObjectState) MustModeState(mode string) (res *ExecutionState) {
	switch mode {
	case "execution":
		res = st.ExecutionState
	case "validation":
		res = st.Validation
	default:
		panic("'" + mode + "' is unknown object processing mode")
	}
	if res == nil {
		panic("object is not in " + mode + " mode")
	}
	if res.Current == nil {
		panic("object "+ res.Ref.String() +" has no Current")
	}
	return res
}

func (st *ObjectState) WrapError(err error, message string) error {
	if err == nil {
		err = errors.New(message)
	} else {
		err = errors.Wrap(err, message)
	}
	return Error{
		Err: err,
	}
}

func makeWMMessage(ctx context.Context, payLoad watermillMsg.Payload, msgType string) *watermillMsg.Message {
	wmMsg := watermillMsg.NewMessage(watermill.NewUUID(), payLoad)
	wmMsg.Metadata.Set(wmBus.MetaTraceID, inslogger.TraceID(ctx))
	wmMsg.Metadata.Set(MessageTypeField, msgType)

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
	JetCoordinator             jet.Coordinator                    `inject:""`

	Executors    [insolar.MachineTypesLastID]insolar.MachineLogicExecutor
	machinePrefs []insolar.MachineType
	Cfg          *configuration.LogicRunner

	state      map[Ref]*ObjectState // if object exists, we are validating or executing it right now
	stateMutex sync.RWMutex

	sock net.Listener

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

	err := initHandlers(&res)
	if err != nil {
		return nil, errors.Wrap(err, "Error while init handlers for logic runner:")
	}

	return &res, nil
}

func initHandlers(lr *LogicRunner) error {
	wmLogger := log.NewWatermillLogAdapter(inslogger.FromContext(context.Background()))
	pubSub := gochannel.NewGoChannel(gochannel.Config{}, wmLogger)

	dep := &Dependencies{
		Publisher: pubSub,
		lr:        lr,
	}

	initHandle := func(msg bus.Message) *Init {
		return &Init{
			dep:     dep,
			Message: msg,
		}
	}
	lr.FlowDispatcher = dispatcher.NewDispatcher(func(msg bus.Message) flow.Handle {
		return initHandle(msg).Present
	}, func(msg bus.Message) flow.Handle {
		return initHandle(msg).Future
	})

	innerInitHandle := func(msg bus.Message) *InnerInit {
		innerMsg := msg.WatermillMsg
		return &InnerInit{
			dep:     dep,
			Message: innerMsg,
		}
	}

	lr.innerFlowDispatcher = dispatcher.NewDispatcher(func(msg bus.Message) flow.Handle {
		return innerInitHandle(msg).Present
	}, func(msg bus.Message) flow.Handle {
		return innerInitHandle(msg).Present
	})

	router, err := watermillMsg.NewRouter(watermillMsg.RouterConfig{}, wmLogger)
	if err != nil {
		return errors.Wrap(err, "Error while creating new watermill router")
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

	return nil
}

// Start starts logic runner component
func (lr *LogicRunner) Start(ctx context.Context) error {
	if lr.Cfg.BuiltIn != nil {
		bi := builtin.NewBuiltIn(lr.MessageBus, lr.ArtifactManager)
		if err := lr.RegisterExecutor(insolar.MachineTypeBuiltin, bi); err != nil {
			return err
		}
		lr.machinePrefs = append(lr.machinePrefs, insolar.MachineTypeBuiltin)
	}

	if lr.Cfg.GoPlugin != nil {
		if lr.Cfg.RPCListen != "" {
			StartRPC(ctx, lr)
		}

		gp, err := goplugin.NewGoPlugin(lr.Cfg, lr.MessageBus, lr.ArtifactManager)
		if err != nil {
			return err
		}
		if err := lr.RegisterExecutor(insolar.MachineTypeGoPlugin, gp); err != nil {
			return err
		}
		lr.machinePrefs = append(lr.machinePrefs, insolar.MachineTypeGoPlugin)
	}

	lr.RegisterHandlers()

	return nil
}

func (lr *LogicRunner) RegisterHandlers() {
	lr.MessageBus.MustRegister(insolar.TypeCallMethod, lr.FlowDispatcher.WrapBusHandle)
	lr.MessageBus.MustRegister(insolar.TypeExecutorResults, lr.FlowDispatcher.WrapBusHandle)
	lr.MessageBus.MustRegister(insolar.TypeValidateCaseBind, lr.HandleValidateCaseBindMessage)
	lr.MessageBus.MustRegister(insolar.TypeValidationResults, lr.HandleValidationResultsMessage)
	lr.MessageBus.MustRegister(insolar.TypePendingFinished, lr.FlowDispatcher.WrapBusHandle)
	lr.MessageBus.MustRegister(insolar.TypeStillExecuting, lr.FlowDispatcher.WrapBusHandle)
	lr.MessageBus.MustRegister(insolar.TypeAbandonedRequestsNotification, lr.FlowDispatcher.WrapBusHandle)
}

// Stop stops logic runner component and its executors
func (lr *LogicRunner) Stop(ctx context.Context) error {
	reterr := error(nil)
	for _, e := range lr.Executors {
		if e == nil {
			continue
		}
		err := e.Stop()
		if err != nil {
			reterr = errors.Wrap(reterr, err.Error())
		}
	}

	if lr.sock != nil {
		if err := lr.sock.Close(); err != nil {
			return err
		}
	}

	lr.router.Close()

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
		return errors.New("can't execute this object")
	}
	return nil
}

func loggerWithTargetID(ctx context.Context, msg insolar.Parcel) context.Context {
	context, _ := inslogger.WithField(ctx, "targetid", msg.DefaultTarget().String())
	return context
}

func (lr *LogicRunner) CheckExecutionLoop(
	ctx context.Context, es *ExecutionState, parcel insolar.Parcel,
) bool {
	if es.Current == nil {
		return false
	}

	if es.Current.SentResult {
		return false
	}

	if es.Current.Request.ReturnMode == record.ReturnNoWait {
		return false
	}

	msg, ok := parcel.Message().(*message.CallMethod)
	if ok && msg.ReturnMode == record.ReturnNoWait {
		return false
	}

	if inslogger.TraceID(es.Current.Context) != inslogger.TraceID(ctx) {
		return false
	}

	inslogger.FromContext(ctx).Debug("loop detected")

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

	pulse := lr.pulse(ctx)
	meCurrent, _ := lr.JetCoordinator.IsAuthorized(
		ctx, insolar.DynamicRoleVirtualExecutor, *es.Ref.Record(), pulse.PulseNumber, lr.JetCoordinator.Me(),
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

func (lr *LogicRunner) executeOrValidate(
	ctx context.Context, es *ExecutionState, parcel insolar.Parcel,
) {
	ctx, span := instracer.StartSpan(ctx, "LogicRunner.ExecuteOrValidate")
	defer span.End()

	msg := parcel.Message().(*message.CallMethod)

	es.Current.LogicContext = &insolar.LogicCallContext{
		Mode:            "execution",
		Caller:          msg.GetCaller(),
		Callee:          &es.Ref,
		Request:         es.Current.RequestRef,
		Time:            time.Now(), // TODO: probably we should take it earlier
		Pulse:           *lr.pulse(ctx),
		TraceID:         inslogger.TraceID(ctx),
		CallerPrototype: &msg.CallerPrototype,
		Immutable:       msg.Immutable,
	}

	var re insolar.Reply
	var err error
	switch msg.CallType {
	case record.CTMethod:
		re, err = lr.executeMethodCall(ctx, es, msg)

	case record.CTSaveAsChild, record.CTSaveAsDelegate:
		re, err = lr.executeConstructorCall(ctx, es, msg)

	default:
		panic("Unknown e type")
	}
	errstr := ""
	if err != nil {
		inslogger.FromContext(ctx).Warn("contract execution error: ", err)
		errstr = err.Error()
	}

	es.Lock()
	defer es.Unlock()

	es.Current.SentResult = true
	if es.Current.Request.ReturnMode != record.ReturnResult {
		return
	}

	target := *es.Current.RequesterNode
	request := *es.Current.RequestRef
	seq := es.Current.Request.Sequence

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

func (lr *LogicRunner) getExecStateFromRef(ctx context.Context, rawRef []byte) *ExecutionState {
	ref := Ref{}.FromSlice(rawRef)

	// TODO: we should stop processing here if 'ref' doesn't exist. Made UpsertObjectState return error if so?
	os := lr.UpsertObjectState(ref)

	os.Lock()
	defer os.Unlock()
	if os.ExecutionState == nil {
		inslogger.FromContext(ctx).Info("[ ProcessExecutionQueue ] got not existing reference. It's strange")
		return nil
	}
	es := os.ExecutionState

	return es
}

func (lr *LogicRunner) unsafeGetLedgerPendingRequest(ctx context.Context, es *ExecutionState) *insolar.Reference {
	es.Lock()
	if es.LedgerQueueElement != nil || !es.LedgerHasMoreRequests {
		es.Unlock()
		return nil
	}
	es.Unlock()

	ledgerHasMore := true

	id := *es.Ref.Record()

	parcel, err := lr.ArtifactManager.GetPendingRequest(ctx, id)
	if err != nil {
		if err != insolar.ErrNoPendingRequest {
			inslogger.FromContext(ctx).Debug("GetPendingRequest failed with error")
			return nil
		}

		ledgerHasMore = false
	}
	es.Lock()
	defer es.Unlock()

	if !ledgerHasMore {
		es.LedgerHasMoreRequests = ledgerHasMore
		return nil
	}

	msg := parcel.Message().(*message.CallMethod)

	pulse := lr.pulse(ctx).PulseNumber
	authorized, err := lr.JetCoordinator.IsAuthorized(
		ctx, insolar.DynamicRoleVirtualExecutor, id, pulse, lr.JetCoordinator.Me(),
	)
	if err != nil {
		inslogger.FromContext(ctx).Debug("Authorization failed with error in getLedgerPendingRequest")
		return nil
	}

	if !authorized {
		inslogger.FromContext(ctx).Debug("pulse changed, can't process abandoned messages for this object")
		return nil
	}

	request := msg.GetReference()
	request.SetRecord(id)

	es.LedgerHasMoreRequests = ledgerHasMore
	es.LedgerQueueElement = &ExecutionQueueElement{
		ctx:        ctx,
		parcel:     parcel,
		request:    &request,
		fromLedger: true,
	}

	return msg.DefaultTarget()
}

func (lr *LogicRunner) executeMethodCall(ctx context.Context, es *ExecutionState, m *message.CallMethod) (insolar.Reply, error) {

	objDesc, err := lr.ArtifactManager.GetObject(ctx, *m.Object)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't get object")
	}
	es.ObjectDescriptor = objDesc

	if es.PrototypeDescriptor == nil {
		protoRef, err := objDesc.Prototype()
		if err != nil {
			return nil, errors.Wrap(err, "couldn't get prototype")
		}

		protoDesc, codeDesc, err := lr.getDescriptorsByPrototypeRef(ctx, *protoRef)
		if err != nil {
			return nil, errors.Wrap(err, "couldn't get descriptors by prototype reference")
		}

		es.PrototypeDescriptor = protoDesc
		es.CodeDescriptor = codeDesc
	}

	current := *es.Current
	current.LogicContext.Prototype = es.PrototypeDescriptor.HeadRef()
	current.LogicContext.Code = es.CodeDescriptor.Ref()
	current.LogicContext.Parent = es.ObjectDescriptor.Parent()
	// it's needed to assure that we call method on ref, that has same prototype as proxy, that we import in contract code
	if m.Prototype != nil && !m.Prototype.Equal(*es.PrototypeDescriptor.HeadRef()) {
		return nil, errors.New("proxy call error: try to call method of prototype as method of another prototype")
	}

	executor, err := lr.GetExecutor(es.CodeDescriptor.MachineType())
	if err != nil {
		return nil, es.WrapError(err, "no executor registered")
	}

	newData, result, err := executor.CallMethod(
		ctx, current.LogicContext, *es.CodeDescriptor.Ref(), es.ObjectDescriptor.Memory(), m.Method, m.Arguments,
	)
	if err != nil {
		return nil, es.WrapError(err, "executor error")
	}

	am := lr.ArtifactManager
	if es.deactivate {
		_, err := am.DeactivateObject(
			ctx, Ref{}, *current.RequestRef, es.ObjectDescriptor,
		)
		if err != nil {
			return nil, es.WrapError(err, "couldn't deactivate object")
		}
	} else if !bytes.Equal(es.ObjectDescriptor.Memory(), newData) {
		_, err := am.UpdateObject(ctx, Ref{}, *current.RequestRef, es.ObjectDescriptor, newData)
		if err != nil {
			return nil, es.WrapError(err, "couldn't update object")
		}
	}
	_, err = am.RegisterResult(ctx, *m.Object, *current.RequestRef, result)
	if err != nil {
		return nil, es.WrapError(err, "couldn't save results")
	}

	return &reply.CallMethod{Result: result}, nil
}

func (lr *LogicRunner) getDescriptorsByPrototypeRef(
	ctx context.Context, protoRef Ref,
) (
	artifacts.ObjectDescriptor, artifacts.CodeDescriptor, error,
) {
	protoDesc, err := lr.ArtifactManager.GetObject(ctx, protoRef)
	if err != nil {
		return nil, nil, errors.Wrap(err, "couldn't get prototype descriptor")
	}
	codeRef, err := protoDesc.Code()
	if err != nil {
		return nil, nil, errors.Wrap(err, "couldn't get code reference")
	}
	codeDesc, err := lr.ArtifactManager.GetCode(ctx, *codeRef)
	if err != nil {
		return nil, nil, errors.Wrap(err, "couldn't get code descriptor")
	}

	return protoDesc, codeDesc, nil
}

func (lr *LogicRunner) getDescriptorsByObjectRef(
	ctx context.Context, objRef Ref,
) (
	artifacts.ObjectDescriptor, artifacts.ObjectDescriptor, artifacts.CodeDescriptor, error,
) {
	ctx, span := instracer.StartSpan(ctx, "LogicRunner.getDescriptorsByObjectRef")
	defer span.End()

	objDesc, err := lr.ArtifactManager.GetObject(ctx, objRef)
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "couldn't get object")
	}

	protoRef, err := objDesc.Prototype()
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "couldn't get prototype reference")
	}

	protoDesc, codeDesc, err := lr.getDescriptorsByPrototypeRef(ctx, *protoRef)
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "couldn't resolve prototype reference to descriptors")
	}

	return objDesc, protoDesc, codeDesc, nil
}

func (lr *LogicRunner) executeConstructorCall(
	ctx context.Context, es *ExecutionState, m *message.CallMethod,
) (
	insolar.Reply, error,
) {
	current := *es.Current
	if current.LogicContext.Caller.IsEmpty() {
		return nil, es.WrapError(nil, "Call constructor from nowhere")
	}

	if m.Prototype == nil {
		return nil, es.WrapError(nil, "prototype reference is required")
	}

	protoDesc, codeDesc, err := lr.getDescriptorsByPrototypeRef(ctx, *m.Prototype)
	if err != nil {
		return nil, es.WrapError(err, "couldn't descriptors")
	}

	current.LogicContext.Prototype = protoDesc.HeadRef()
	current.LogicContext.Code = codeDesc.Ref()

	executor, err := lr.GetExecutor(codeDesc.MachineType())
	if err != nil {
		return nil, es.WrapError(err, "no executer registered")
	}

	newData, err := executor.CallConstructor(ctx, current.LogicContext, *codeDesc.Ref(), m.Method, m.Arguments)
	if err != nil {
		return nil, es.WrapError(err, "executer error")
	}

	switch m.CallType {
	case record.CTSaveAsChild, record.CTSaveAsDelegate:
		_, err = lr.ArtifactManager.ActivateObject(
			ctx,
			Ref{}, *current.RequestRef, *m.Base, *m.Prototype, m.CallType == record.CTSaveAsDelegate, newData,
		)
		if err != nil {
			return nil, es.WrapError(err, "couldn't activate object")
		}
		_, err = lr.ArtifactManager.RegisterResult(ctx, *current.RequestRef, *current.RequestRef, nil)
		if err != nil {
			return nil, es.WrapError(err, "couldn't save results")
		}
		return &reply.CallConstructor{Object: current.RequestRef}, err

	default:
		return nil, es.WrapError(nil, "unsupported type of save object")
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

	ctx, spanStates := instracer.StartSpan(ctx, "pulse.logicrunner processing of states")
	for ref, state := range lr.state {
		meNext, _ := lr.JetCoordinator.IsAuthorized(
			ctx, insolar.DynamicRoleVirtualExecutor, *ref.Record(), pulse.PulseNumber, lr.JetCoordinator.Me(),
		)
		state.Lock()

		// some old stuff
		state.RefreshConsensus()

		if es := state.ExecutionState; es != nil {
			es.Lock()

			// if we are executor again we just continue working
			// without sending data on next executor (because we are next executor)
			if !meNext {
				sendExecResults := false

				if es.Current != nil {
					es.pending = message.InPending
					sendExecResults = true

					// TODO: this should return delegation token to continue execution of the pending
					messages = append(
						messages,
						&message.StillExecuting{
							Reference: ref,
						},
					)
				} else {
					if es.pending == message.InPending && !es.PendingConfirmed {
						inslogger.FromContext(ctx).Warn(
							"looks like pending executor died, continuing execution",
						)
						es.pending = message.NotPending
						sendExecResults = true
						es.LedgerHasMoreRequests = true
					}

					state.ExecutionState = nil
				}

				queue, ledgerHasMoreRequest := es.releaseQueue()
				if len(queue) > 0 || sendExecResults {
					// TODO: we also should send when executed something for validation
					// TODO: now validation is disabled
					messagesQueue := convertQueueToMessageQueue(queue)

					messages = append(
						messages,
						//&message.ValidateCaseBind{
						//	Reference: ref,
						//	Requests:  requests,
						//	Pulse:     pulse,
						//},
						&message.ExecutorResults{
							RecordRef:             ref,
							Pending:               es.pending,
							Queue:                 messagesQueue,
							LedgerHasMoreRequests: es.LedgerHasMoreRequests || ledgerHasMoreRequest,
						},
					)
				}
			} else {
				if es.Current != nil {
					// no pending should be as we are executing
					if es.pending == message.InPending {
						inslogger.FromContext(ctx).Warn(
							"we are executing ATM, but ES marked as pending, shouldn't be",
						)
						es.pending = message.NotPending
					}
				} else if es.pending == message.InPending && !es.PendingConfirmed {
					inslogger.FromContext(ctx).Warn(
						"looks like pending executor died, continuing execution",
					)
					es.pending = message.NotPending
					es.LedgerHasMoreRequests = true
					lr.startGetLedgerPendingRequest(ctx, es)
				}
				es.PendingConfirmed = false
			}

			es.Unlock()
		}

		if state.ExecutionState == nil && state.Validation == nil && state.Consensus == nil {
			delete(lr.state, ref)
		}

		state.Unlock()
	}
	spanStates.End()

	lr.stateMutex.Unlock()

	if len(messages) > 0 {
		go lr.sendOnPulseMessagesAsync(ctx, messages)
	}

	lr.stopIfNeeded(ctx)

	return nil
}

func (lr *LogicRunner) stopIfNeeded(ctx context.Context) {
	// lock is required to access lr.state
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

func (lr *LogicRunner) HandleAbandonedRequestsNotificationMessage(
	ctx context.Context, parcel insolar.Parcel,
) (
	insolar.Reply, error,
) {
	return lr.FlowDispatcher.WrapBusHandle(ctx, parcel)
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

func convertQueueToMessageQueue(queue []ExecutionQueueElement) []message.ExecutionQueueElement {
	mq := make([]message.ExecutionQueueElement, 0)
	for _, elem := range queue {
		mq = append(mq, message.ExecutionQueueElement{
			Parcel:  elem.parcel,
			Request: elem.request,
		})
	}

	return mq
}
