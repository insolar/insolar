/*
 *    Copyright 2018 Insolar
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

// Package logicrunner - infrastructure for executing smartcontracts
package logicrunner

import (
	"bytes"
	"context"
	"encoding/gob"
	"net"
	"sync"
	"time"

	"github.com/pkg/errors"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/core/message"
	"github.com/insolar/insolar/core/reply"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/logicrunner/builtin"
	"github.com/insolar/insolar/logicrunner/goplugin"
)

type Ref = core.RecordRef

// Context of one contract execution
type ExecutionState struct {
	sync.Mutex
	Ref    *Ref
	Method string

	validate    bool
	insContext  context.Context
	callContext *core.LogicCallContext
	deactivate  bool
	request     *Ref
	traceID     string
	queueLength int

	caseBind      core.CaseBind
	caseBindMutex sync.Mutex

	objectbody *ObjectBody
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

func (es *ExecutionState) ErrorWrap(err error, message string) error {
	if err == nil {
		err = errors.New(message)
	} else {
		err = errors.Wrap(err, message)
	}
	return Error{
		Err:      err,
		Request:  es.request,
		Contract: es.Ref,
		Method:   es.Method,
	}
}

func (es *ExecutionState) Lock() {
	es.queueLength++
	es.Mutex.Lock()
}

func (es *ExecutionState) Unlock() {
	es.queueLength--
	es.Mutex.Unlock()
	es.traceID = "Done"
}

func (es *ExecutionState) ReleaseQueue() {
	for es.queueLength > 1 {
		es.Unlock()
	}
}

func (es *ExecutionState) AddCaseRequest(record core.CaseRecord) {
	es.caseBindMutex.Lock()
	defer es.caseBindMutex.Unlock()

	es.caseBind.Requests = append(es.caseBind.Requests, core.CaseRequest{
		Request: record,
		Records: make([]core.CaseRecord, 0),
	})
}

func (es *ExecutionState) AddCaseRecord(record core.CaseRecord) {
	es.caseBindMutex.Lock()
	defer es.caseBindMutex.Unlock()

	requests := es.caseBind.Requests
	if len(requests) == 0 {
		panic("attempt to add record into case bind before any requests were added")
	}

	lastRequest := requests[len(requests)-1]
	lastRequest.Records = append(lastRequest.Records, record)
	requests[len(requests)-1] = lastRequest
	es.caseBind.Requests = requests
}

// LogicRunner is a general interface of contract executor
type LogicRunner struct {
	// FIXME: Ledger component is deprecated. Inject required sub-components.
	MessageBus                 core.MessageBus                 `inject:""`
	Ledger                     core.Ledger                     `inject:""`
	Network                    core.Network                    `inject:""`
	PlatformCryptographyScheme core.PlatformCryptographyScheme `inject:""`
	ParcelFactory              message.ParcelFactory           `inject:""`
	PulseManager               core.PulseManager               `inject:""`
	ArtifactManager            core.ArtifactManager            `inject:""`
	JetCoordinator             core.JetCoordinator             `inject:""`

	Executors      [core.MachineTypesLastID]core.MachineLogicExecutor
	machinePrefs   []core.MachineType
	Cfg            *configuration.LogicRunner
	execution      map[Ref]*ExecutionState // if object exists, we are validating or executing it right now
	executionMutex sync.Mutex

	caseBindReplays      map[Ref]core.CaseBindReplay
	caseBindReplaysMutex sync.Mutex
	consensus            map[Ref]*Consensus
	consensusMutex       sync.Mutex
	sock                 net.Listener
}

// NewLogicRunner is constructor for LogicRunner
func NewLogicRunner(cfg *configuration.LogicRunner) (*LogicRunner, error) {
	if cfg == nil {
		return nil, errors.New("LogicRunner have nil configuration")
	}
	res := LogicRunner{
		Cfg:             cfg,
		execution:       make(map[Ref]*ExecutionState),
		caseBindReplays: make(map[Ref]core.CaseBindReplay),
	}
	return &res, nil
}

// Start starts logic runner component
func (lr *LogicRunner) Start(ctx context.Context) error {
	if lr.Cfg.BuiltIn != nil {
		bi := builtin.NewBuiltIn(lr.MessageBus, lr.ArtifactManager)
		if err := lr.RegisterExecutor(core.MachineTypeBuiltin, bi); err != nil {
			return err
		}
		lr.machinePrefs = append(lr.machinePrefs, core.MachineTypeBuiltin)
	}

	if lr.Cfg.GoPlugin != nil {
		if lr.Cfg.RPCListen != "" {
			StartRPC(ctx, lr)
		}

		gp, err := goplugin.NewGoPlugin(lr.Cfg, lr.MessageBus, lr.ArtifactManager)
		if err != nil {
			return err
		}
		if err := lr.RegisterExecutor(core.MachineTypeGoPlugin, gp); err != nil {
			return err
		}
		lr.machinePrefs = append(lr.machinePrefs, core.MachineTypeGoPlugin)
	}

	lr.RegisterHandlers()

	return nil
}

func (lr *LogicRunner) RegisterHandlers() {
	lr.MessageBus.MustRegister(core.TypeCallMethod, lr.Execute)
	lr.MessageBus.MustRegister(core.TypeCallConstructor, lr.Execute)
	lr.MessageBus.MustRegister(core.TypeExecutorResults, lr.ExecutorResults)
	lr.MessageBus.MustRegister(core.TypeValidateCaseBind, lr.ValidateCaseBind)
	lr.MessageBus.MustRegister(core.TypeValidationResults, lr.ProcessValidationResults)
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

	return reterr
}

func (lr *LogicRunner) CheckOurRole(ctx context.Context, msg core.Message, role core.DynamicRole) error {
	// TODO do map of supported objects for pulse, go to jetCoordinator only if map is empty for ref
	target := message.ExtractTarget(msg)
	isAuthorized, err := lr.JetCoordinator.IsAuthorized(
		ctx, role, &target, lr.pulse(ctx).PulseNumber, lr.Network.GetNodeID(),
	)
	if err != nil {
		return errors.Wrap(err, "authorization failed with error")
	}
	if !isAuthorized {
		return errors.New("can't execute this object")
	}
	return nil
}

// Execute runs a method on an object, ATM just thin proxy to `GoPlugin.Exec`
func (lr *LogicRunner) Execute(ctx context.Context, parcel core.Parcel) (core.Reply, error) {
	msg, ok := parcel.Message().(message.IBaseLogicMessage)
	if !ok {
		return nil, errors.New("Execute( ! message.IBaseLogicMessage )")
	}
	ref := msg.GetReference()

	es := lr.UpsertExecution(ref)

	err := lr.CheckOurRole(ctx, msg, core.DynamicRoleVirtualExecutor)
	if err != nil {
		return nil, es.ErrorWrap(err, "can't play role")
	}

	if es.traceID == inslogger.TraceID(ctx) {
		return nil, es.ErrorWrap(nil, "loop detected")
	}
	entryPulse := lr.pulse(ctx).PulseNumber

	es.Lock()

	// unlock comes from OnPulse()
	// pulse changed while we was locked and we don't process anything
	if entryPulse != lr.pulse(ctx).PulseNumber {
		return nil, es.ErrorWrap(nil, "abort execution: new Pulse coming")
	}

	return lr.executeOrValidate(ctx, es, ValidationSaver{lr: lr}, parcel)
}

func (lr *LogicRunner) executeOrValidate(
	ctx context.Context, es *ExecutionState, vb ValidationBehaviour, parcel core.Parcel,
) (
	core.Reply, error,
) {
	fuse := true
	defer func() {
		if fuse {
			es.Unlock()
		}
	}()

	msg := parcel.Message().(message.IBaseLogicMessage)

	ref := *es.Ref

	es.traceID = inslogger.TraceID(ctx)
	es.insContext = ctx

	es.AddCaseRequest(core.CaseRecord{
		Type: core.CaseRecordTypeStart,
		Resp: msg,
	})

	vb.Begin(ref, core.CaseRecord{
		Type: core.CaseRecordTypeTraceID,
		Resp: inslogger.TraceID(ctx),
	})

	var err error
	es.request, err = vb.RegisterRequest(parcel)
	if err != nil {
		return nil, es.ErrorWrap(err, "can't create request")
	}

	es.callContext = &core.LogicCallContext{
		Caller:          msg.GetCaller(),
		Callee:          &ref,
		Request:         es.request,
		Time:            time.Now(), // TODO: probably we should take it from e
		Pulse:           *lr.pulse(ctx),
		TraceID:         inslogger.TraceID(ctx),
		CallerPrototype: msg.GetCallerPrototype(),
	}

	switch m := msg.(type) {
	case *message.CallMethod:
		es.Method = m.Method
		fuse = false
		re, err := lr.executeMethodCall(es, m, vb)
		return re, err

	case *message.CallConstructor:
		es.Method = m.Name
		fuse = false
		re, err := lr.executeConstructorCall(es, m, vb)
		return re, err

	default:
		panic("Unknown e type")
	}
}

// ObjectBody is an inner representation of object and all it accessory
// make it private again when we start it serialize before sending
type ObjectBody struct {
	objDescriptor   core.ObjectDescriptor
	Object          []byte
	ClassHeadRef    *Ref
	CodeMachineType core.MachineType
	CodeRef         *Ref
	Parent          *Ref
}

func init() {
	gob.Register(&ObjectBody{})
}

func (lr *LogicRunner) getObjectMessage(es *ExecutionState, objref Ref) error {
	ctx := es.insContext
	cr, step := lr.nextValidationStep(objref)
	// TODO: move this to vb, when vb become a part of es
	if es.objectbody != nil { // already have something
		if step > 0 { // check signature
			if core.CaseRecordTypeSignObject != cr.Type {
				return errors.Errorf("Wrong validation type on CaseRecordTypeSignObject %d, ", cr.Type)
			}
			if !bytes.Equal(cr.ReqSig, HashInterface(lr.PlatformCryptographyScheme, objref)) {
				return errors.New("Wrong validation sig on CaseRecordTypeSignObject")
			}
			if !bytes.Equal(cr.Resp.([]byte), HashInterface(lr.PlatformCryptographyScheme, es.objectbody)) {
				return errors.New("Wrong validation comparision on CaseRecordTypeSignObject")
			}

		} else {
			lr.addObjectCaseRecord(objref, core.CaseRecord{
				Type:   core.CaseRecordTypeSignObject,
				ReqSig: HashInterface(lr.PlatformCryptographyScheme, objref),
				Resp:   HashInterface(lr.PlatformCryptographyScheme, es.objectbody),
			})
		}
		return nil
	}

	if step >= 0 { // validate
		if core.CaseRecordTypeGetObject != cr.Type {
			return errors.Errorf("Wrong validation type on CaseRecordTypeGetObject %d", cr.Type)
		}
		sig := HashInterface(lr.PlatformCryptographyScheme, objref)
		if !bytes.Equal(cr.ReqSig, sig) {
			return errors.New("Wrong validation sig on CaseRecordTypeGetObject")
		}
		es.objectbody = cr.Resp.(*ObjectBody)
		return nil
	}

	objDesc, err := lr.ArtifactManager.GetObject(ctx, objref, nil, false)
	if err != nil {
		return errors.Wrap(err, "couldn't get object")
	}
	protoRef, err := objDesc.Prototype()
	if err != nil {
		return errors.Wrap(err, "couldn't get prototype reference")
	}
	protoDesc, err := lr.ArtifactManager.GetObject(ctx, *protoRef, nil, false)
	if err != nil {
		return errors.Wrap(err, "couldn't get object's class")
	}
	codeRef, err := protoDesc.Code()
	if err != nil {
		return errors.Wrap(err, "couldn't get code reference")
	}
	codeDesc, err := lr.ArtifactManager.GetCode(ctx, *codeRef)
	if err != nil {
		return errors.Wrap(err, "couldn't get code")
	}
	es.objectbody = &ObjectBody{
		objDescriptor:   objDesc,
		Object:          objDesc.Memory(),
		ClassHeadRef:    protoDesc.HeadRef(),
		CodeMachineType: codeDesc.MachineType(),
		CodeRef:         codeDesc.Ref(),
		Parent:          objDesc.Parent(),
	}
	bcopy := *es.objectbody
	copy(bcopy.Object, es.objectbody.Object)
	lr.addObjectCaseRecord(objref, core.CaseRecord{
		Type:   core.CaseRecordTypeGetObject,
		ReqSig: HashInterface(lr.PlatformCryptographyScheme, objref),
		Resp:   &bcopy,
	})
	return nil
}

func (lr *LogicRunner) executeMethodCall(es *ExecutionState, m *message.CallMethod, vb ValidationBehaviour) (core.Reply, error) {
	ctx := es.insContext

	delayedUnlock := false
	defer func() {
		if !delayedUnlock {
			es.Unlock()
		}
	}()

	err := lr.getObjectMessage(es, m.ObjectRef)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't get object message")
	}

	es.callContext.Prototype = es.objectbody.ClassHeadRef
	es.callContext.Code = es.objectbody.CodeRef
	es.callContext.Parent = es.objectbody.Parent

	vb.ModifyContext(es.callContext)

	executor, err := lr.GetExecutor(es.objectbody.CodeMachineType)
	if err != nil {
		return nil, es.ErrorWrap(err, "no executor registered")
	}

	executeFunction := func() (*reply.CallMethod, error) {
		newData, result, err := executor.CallMethod(
			ctx, es.callContext, *es.objectbody.CodeRef, es.objectbody.Object, m.Method, m.Arguments,
		)
		if err != nil {
			return nil, es.ErrorWrap(err, "executor error")
		}

		if vb.NeedSave() {
			am := lr.ArtifactManager
			if es.deactivate {
				_, err = am.DeactivateObject(
					ctx, Ref{}, *es.request, es.objectbody.objDescriptor,
				)
			} else {
				od, e := am.UpdateObject(ctx, Ref{}, *es.request, es.objectbody.objDescriptor, newData)
				err = e
				if od != nil && e == nil {
					es.objectbody.objDescriptor = od
				}
			}
			if err != nil {
				return nil, es.ErrorWrap(err, "couldn't update object")
			}
			_, err = am.RegisterResult(ctx, *es.request, result)
			if err != nil {
				return nil, es.ErrorWrap(err, "couldn't save results")
			}
		}

		es.objectbody.Object = newData
		re := &reply.CallMethod{Data: newData, Result: result}

		vb.End(m.ObjectRef, core.CaseRecord{
			Type: core.CaseRecordTypeResult,
			Resp: re,
		})
		return re, nil
	}

	switch m.ReturnMode {
	case message.ReturnResult:
		return executeFunction()
	case message.ReturnNoWait:
		delayedUnlock = true
		go func() {
			defer es.Unlock()
			_, err := executeFunction()
			if err != nil {
				inslogger.FromContext(ctx).Error(err)
			}
		}()
		return &reply.CallMethod{}, nil
	}
	return nil, errors.Errorf("Invalid ReturnMode #%d", m.ReturnMode)
}

func (lr *LogicRunner) executeConstructorCall(es *ExecutionState, m *message.CallConstructor, vb ValidationBehaviour) (core.Reply, error) {
	ctx := es.insContext
	defer es.Unlock()

	if es.callContext.Caller.IsEmpty() {
		return nil, es.ErrorWrap(nil, "Call constructor from nowhere")
	}

	protoDesc, err := lr.ArtifactManager.GetObject(ctx, m.PrototypeRef, nil, false)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't get prototype")
	}
	es.callContext.Prototype = protoDesc.HeadRef()

	codeRef, err := protoDesc.Code()
	if err != nil {
		return nil, es.ErrorWrap(err, "couldn't code reference")
	}
	codeDesc, err := lr.ArtifactManager.GetCode(ctx, *codeRef)
	if err != nil {
		return nil, es.ErrorWrap(err, "couldn't code")
	}
	es.callContext.Code = codeDesc.Ref()

	executor, err := lr.GetExecutor(codeDesc.MachineType())
	if err != nil {
		return nil, es.ErrorWrap(err, "no executer registered")
	}

	newData, err := executor.CallConstructor(ctx, es.callContext, *codeDesc.Ref(), m.Name, m.Arguments)
	if err != nil {
		return nil, es.ErrorWrap(err, "executer error")
	}

	switch m.SaveAs {
	case message.Child, message.Delegate:
		if vb.NeedSave() {
			_, err = lr.ArtifactManager.ActivateObject(
				ctx,
				Ref{}, *es.request, m.ParentRef, m.PrototypeRef, m.SaveAs == message.Delegate, newData,
			)
		}
		vb.End(m.GetReference(), core.CaseRecord{
			Type: core.CaseRecordTypeResult,
			Resp: &reply.CallConstructor{Object: es.request},
		})
		return &reply.CallConstructor{Object: es.request}, err
	default:
		return nil, es.ErrorWrap(nil, "unsupported type of save object")
	}
}

func (lr *LogicRunner) OnPulse(ctx context.Context, pulse core.Pulse) error {
	lr.RefreshConsensus()

	// start of new Pulse, lock CaseBind data, copy it, clean original, unlock original

	lr.executionMutex.Lock()
	defer lr.executionMutex.Unlock()

	messages := make([]core.Message, 0)

	// send copy for validation
	for ref, state := range lr.execution {
		messages = append(
			messages,
			&message.ValidateCaseBind{RecordRef: ref, CaseBind: state.caseBind, Pulse: pulse},
			&message.ExecutorResults{RecordRef: ref, CaseBind: state.caseBind},
		)

		// release unprocessed request
		state.ReleaseQueue()
	}

	// TODO: this not exactly correct
	lr.execution = make(map[Ref]*ExecutionState)

	for _, msg := range messages {
		_, err := lr.MessageBus.Send(ctx, msg, nil)
		if err != nil {
			return errors.New("error while sending caseBind data to new executor")
		}
	}

	return nil
}
