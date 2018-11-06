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
	noWait      bool
	validate    bool
	insContext  context.Context
	callContext *core.LogicCallContext
	deactivate  bool
	request     *Ref
	traceID     string

	objectbody *ObjectBody
}

func (es *ExecutionState) Unlock() {
	es.Mutex.Unlock()
	es.traceID = "Done"
}

// LogicRunner is a general interface of contract executor
type LogicRunner struct {
	Executors       [core.MachineTypesLastID]core.MachineLogicExecutor
	ArtifactManager core.ArtifactManager
	MessageBus      core.MessageBus
	Ledger          core.Ledger
	Network         core.Network
	machinePrefs    []core.MachineType
	Cfg             *configuration.LogicRunner
	execution       map[Ref]*ExecutionState // if object exists, we are validating or executing it right now
	executionMutex  sync.Mutex

	// TODO move caseBind to context
	caseBind      core.CaseBind
	caseBindMutex sync.Mutex

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
		ArtifactManager: nil,
		Ledger:          nil,
		Cfg:             cfg,
		execution:       make(map[Ref]*ExecutionState),
		caseBind:        core.CaseBind{Records: make(map[Ref][]core.CaseRecord)},
		caseBindReplays: make(map[Ref]core.CaseBindReplay),
	}
	return &res, nil
}

// Start starts logic runner component
func (lr *LogicRunner) Start(ctx context.Context, c core.Components) error {
	am := c.Ledger.GetArtifactManager()
	lr.ArtifactManager = am
	lr.MessageBus = c.MessageBus
	lr.Ledger = c.Ledger
	lr.Network = c.Network

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

	// TODO: use separate handlers
	if err := lr.MessageBus.Register(core.TypeCallMethod, lr.Execute); err != nil {
		return err
	}
	if err := lr.MessageBus.Register(core.TypeCallConstructor, lr.Execute); err != nil {
		return err
	}

	if err := lr.MessageBus.Register(core.TypeExecutorResults, lr.ExecutorResults); err != nil {
		return err
	}
	if err := lr.MessageBus.Register(core.TypeValidateCaseBind, lr.ValidateCaseBind); err != nil {
		return err
	}
	if err := lr.MessageBus.Register(core.TypeValidationResults, lr.ProcessValidationResults); err != nil {
		return err
	}

	return nil
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

// Execute runs a method on an object, ATM just thin proxy to `GoPlugin.Exec`
func (lr *LogicRunner) Execute(ctx context.Context, inmsg core.SignedMessage) (core.Reply, error) {
	// TODO do not pass here message.ValidateCaseBind and message.ExecutorResults
	msg, ok := inmsg.Message().(message.IBaseLogicMessage)
	if !ok {
		return nil, errors.New("Execute( ! message.IBaseLogicMessage )")
	}
	ref := msg.GetReference()

	es := lr.UpsertExecution(ref)
	if lr.execution[ref].traceID == inslogger.TraceID(ctx) {
		return nil, errors.Errorf("loop detected")
	}
	fuse := true
	es.Lock()
	defer func() {
		if fuse {
			es.Unlock()
		}
	}()
	lr.execution[ref].traceID = inslogger.TraceID(ctx)
	es.insContext = ctx

	lr.caseBindReplaysMutex.Lock()
	cb, validate := lr.caseBindReplays[ref]
	lr.caseBindReplaysMutex.Unlock()

	var vb ValidationBehaviour
	if validate {
		vb = ValidationChecker{lr: lr, cb: cb}
	} else {
		vb = ValidationSaver{lr: lr}
	}

	// TODO do map of supported objects for pulse, go to jetCoordinator only if map is empty for ref
	isAuthorized, err := lr.Ledger.GetJetCoordinator().IsAuthorized(
		ctx,
		vb.GetRole(),
		*msg.Target(),
		lr.pulse().PulseNumber,
		lr.Network.GetNodeID(),
	)

	if err != nil {
		return nil, errors.New("Authorization failed with error: " + err.Error())
	}
	if !isAuthorized {
		return nil, errors.New("Can't execute this object")
	}

	vb.Begin(ref, core.CaseRecord{
		Type: core.CaseRecordTypeStart,
		Resp: msg,
	})

	vb.Begin(ref, core.CaseRecord{
		Type: core.CaseRecordTypeTraceID,
		Resp: inslogger.TraceID(ctx),
	})

	es.request, err = vb.RegisterRequest(msg)

	if err != nil {
		return nil, errors.Wrap(err, "Can't create request")
	}

	es.callContext = &core.LogicCallContext{
		Caller:  msg.GetCaller(),
		Callee:  &ref,
		Request: es.request,
		Time:    time.Now(), // TODO: probably we should take it from e
		Pulse:   *lr.pulse(),
		TraceID: inslogger.TraceID(ctx),
	}

	switch m := msg.(type) {
	case *message.CallMethod:
		fuse = false
		re, err := lr.executeMethodCall(es, m, vb)
		return re, err

	case *message.CallConstructor:
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
			if !bytes.Equal(cr.ReqSig, HashInterface(objref)) {
				return errors.New("Wrong validation sig on CaseRecordTypeSignObject")
			}
			if !bytes.Equal(cr.Resp.([]byte), HashInterface(es.objectbody)) {
				return errors.New("Wrong validation comparision on CaseRecordTypeSignObject")
			}

		} else {
			lr.addObjectCaseRecord(objref, core.CaseRecord{
				Type:   core.CaseRecordTypeSignObject,
				ReqSig: HashInterface(objref),
				Resp:   HashInterface(es.objectbody),
			})
		}
		return nil
	}

	if step >= 0 { // validate
		if core.CaseRecordTypeGetObject != cr.Type {
			return errors.Errorf("Wrong validation type on CaseRecordTypeGetObject %d", cr.Type)
		}
		sig := HashInterface(objref)
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
		ReqSig: HashInterface(objref),
		Resp:   &bcopy,
	})
	return nil
}

func (lr *LogicRunner) executeMethodCall(es *ExecutionState, m *message.CallMethod, vb ValidationBehaviour) (core.Reply, error) {
	ctx := es.insContext

	es.noWait = false
	defer func() {
		if !es.noWait {
			es.Unlock()
		}
	}()

	err := lr.getObjectMessage(es, m.ObjectRef)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't get object message")
	}

	es.callContext.Prototype = es.objectbody.ClassHeadRef
	vb.ModifyContext(es.callContext)

	executor, err := lr.GetExecutor(es.objectbody.CodeMachineType)
	if err != nil {
		return nil, errors.Wrap(err, "no executor registered")
	}

	executer := func() (*reply.CallMethod, error) {
		defer func() {
			if es.noWait {
				es.Unlock()
			}
		}()
		newData, result, err := executor.CallMethod(
			ctx, es.callContext, *es.objectbody.CodeRef, es.objectbody.Object, m.Method, m.Arguments,
		)
		if err != nil {
			return nil, errors.Wrap(err, "executor error")
		}

		if vb.NeedSave() {
			am := lr.ArtifactManager
			if es.deactivate {
				_, err = am.DeactivateObject(
					ctx, Ref{}, *es.request, es.objectbody.objDescriptor,
				)
			} else {
				es.objectbody.objDescriptor, err = am.UpdateObject(
					ctx, Ref{}, *es.request, es.objectbody.objDescriptor, newData,
				)
			}
			if err != nil {
				return nil, errors.Wrap(err, "couldn't update object")
			}
			_, err = am.RegisterResult(ctx, *es.request, result)
			if err != nil {
				return nil, errors.Wrap(err, "couldn't save results")
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
		return executer()
	case message.ReturnNoWait:
		es.noWait = true
		go func() {
			_, err := executer()
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
	defer func() {
		es.Unlock()
	}()
	protoDesc, err := lr.ArtifactManager.GetObject(ctx, m.PrototypeRef, nil, false)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't get prototype")
	}
	es.callContext.Prototype = protoDesc.HeadRef()

	codeRef, err := protoDesc.Code()
	if err != nil {
		return nil, errors.Wrap(err, "couldn't code reference")
	}
	codeDesc, err := lr.ArtifactManager.GetCode(ctx, *codeRef)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't code")
	}
	executor, err := lr.GetExecutor(codeDesc.MachineType())
	if err != nil {
		return nil, errors.Wrap(err, "no executer registered")
	}

	newData, err := executor.CallConstructor(ctx, es.callContext, *codeDesc.Ref(), m.Name, m.Arguments)
	if err != nil {
		return nil, errors.Wrap(err, "executer error")
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
		return nil, errors.New("unsupported type of save object")
	}
}

func (lr *LogicRunner) OnPulse(pulse core.Pulse) error {
	insctx := context.TODO()

	lr.RefreshConsensus()
	// start of new Pulse, lock CaseBind data, copy it, clean original, unlock original
	objectsRecords := lr.refreshCaseBind()

	// TODO INS-666
	// TODO make refresh lr.Execution - Unlock mutexes n-1 time for each object, send some info for callers, do empty object

	if len(objectsRecords) == 0 {
		return nil
	}

	// send copy for validation
	for ref, records := range objectsRecords {
		_, err := lr.MessageBus.Send(
			insctx,
			&message.ValidateCaseBind{RecordRef: ref, CaseRecords: records, Pulse: pulse},
		)
		if err != nil {
			panic("Error while sending caseBind data to validators: " + err.Error())
		}

		results := message.ExecutorResults{RecordRef: ref, CaseRecords: records}
		_, err = lr.MessageBus.Send(insctx, &results)
		if err != nil {
			return errors.New("error while sending caseBind data to new executor")
		}
	}

	return nil
}
