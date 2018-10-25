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
	"net"
	"sync"
	"time"

	"github.com/pkg/errors"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/core/message"
	"github.com/insolar/insolar/core/reply"
	"github.com/insolar/insolar/inscontext"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/logicrunner/builtin"
	"github.com/insolar/insolar/logicrunner/goplugin"
)

type Ref = core.RecordRef

// Context of one contract execution
type ExecutionState struct {
	sync.Mutex
	mainContext core.Context
	callContext *core.LogicCallContext
	deactivate  bool
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
func (lr *LogicRunner) Start(ctx core.Context, c core.Components) error {
	am := c.Ledger.GetArtifactManager()
	lr.ArtifactManager = am
	messageBus := c.MessageBus
	lr.MessageBus = messageBus
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
			StartRPC(lr)
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
func (lr *LogicRunner) Stop(ctx core.Context) error {
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
func (lr *LogicRunner) Execute(ctx core.Context, inmsg core.Message) (core.Reply, error) {
	// TODO do not pass here message.ValidateCaseBind and message.ExecutorResults
	msg, ok := inmsg.(message.IBaseLogicMessage)
	if !ok {
		return nil, errors.New("Execute( ! message.IBaseLogicMessage )")
	}

	ref := msg.GetReference()

	es := lr.UpsertExecution(ref)
	ctx.Log().Warnf("LOCKING, %s", ref.String())
	es.Lock()
	ctx.Log().Warnf("LOCKED")
	es.mainContext = ctx

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
	isAuthorized, err := lr.Ledger.GetJetCoordinator().IsAuthorized(vb.GetRole(), *msg.Target(), lr.pulse().PulseNumber, lr.Network.GetNodeID())

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

	es.mainContext.Log().Warnf("VALIDATE IS %b", validate)
	reqref, err := vb.RegisterRequest(msg)

	if err != nil {
		return nil, errors.Wrap(err, "Can't create request")
	}

	es.callContext = &core.LogicCallContext{
		Caller:  msg.GetCaller(),
		Callee:  &ref,
		Request: reqref,
		Time:    time.Now(), // TODO: probably we should take it from e
		Pulse:   *lr.pulse(),
	}

	switch m := msg.(type) {
	case *message.CallMethod:
		re, err := lr.executeMethodCall(es, m, vb)
		return re, err

	case *message.CallConstructor:
		re, err := lr.executeConstructorCall(es, m, vb)
		return re, err

	default:
		panic("Unknown e type")
	}
}

func (lr *LogicRunner) pulse() *core.Pulse {
	pulse, err := lr.Ledger.GetPulseManager().Current()
	if err != nil {
		panic(err)
	}
	return pulse
}

type objectBody struct {
	Object core.ObjectDescriptor
	Class  core.ClassDescriptor
	Code   core.CodeDescriptor
}

func (lr *LogicRunner) getObjectMessage(objref Ref) (*objectBody, error) {
	ctx := inscontext.TODO()
	cr, step := lr.getNextValidationStep(objref)
	if step >= 0 { // validate
		if core.CaseRecordTypeGetObject != cr.Type {
			return nil, errors.New("Wrong validation type on RouteCall")
		}
		sig := HashInterface(objref)
		if !bytes.Equal(cr.ReqSig, sig) {
			return nil, errors.New("Wrong validation sig on RouteCall")
		}
		return cr.Resp.(*objectBody), nil
	}

	objDesc, err := lr.ArtifactManager.GetObject(ctx, objref, nil, false)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't get object")
	}

	classDesc, err := lr.ArtifactManager.GetClass(ctx, *objDesc.Class(), nil)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't get object's class")
	}

	codeDesc := classDesc.CodeDescriptor()
	ob := &objectBody{
		Object: objDesc,
		Class:  classDesc,
		Code:   codeDesc,
	}
	lr.addObjectCaseRecord(objref, core.CaseRecord{
		Type:   core.CaseRecordTypeGetObject,
		ReqSig: HashInterface(objref),
		Resp:   ob,
	})
	return ob, nil
}

func (lr *LogicRunner) executeMethodCall(es *ExecutionState, m *message.CallMethod, vb ValidationBehaviour) (core.Reply, error) {
	insctx := es.mainContext

	objbody, err := lr.getObjectMessage(m.ObjectRef)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't get object message")
	}

	es.callContext.Class = objbody.Class.HeadRef()
	vb.ModifyContext(es.callContext)

	executor, err := lr.GetExecutor(objbody.Code.MachineType())
	if err != nil {
		return nil, errors.Wrap(err, "no executor registered")
	}

	executer := func() (*reply.CallMethod, error) {
		newData, result, err := executor.CallMethod(
			es.callContext, *objbody.Code.Ref(), objbody.Object.Memory(), m.Method, m.Arguments,
		)
		if err != nil {
			return nil, errors.Wrap(err, "executor error")
		}

		if vb.NeedSave() {
			am := lr.ArtifactManager
			if es.deactivate {
				_, err = am.DeactivateObject(
					insctx, Ref{}, *es.callContext.Request, objbody.Object,
				)
			} else {
				_, err = am.UpdateObject(
					insctx, Ref{}, *es.callContext.Request, objbody.Object, newData,
				)
			}
			if err != nil {
				return nil, errors.Wrap(err, "couldn't update object")
			}
		}

		re := &reply.CallMethod{Data: newData, Result: result}

		vb.End(m.ObjectRef, core.CaseRecord{
			Type: core.CaseRecordTypeResult,
			Resp: re,
		})
		insctx.Log().Warnf("UNLOCK METHOD")
		es.Unlock()
		return re, nil
	}

	switch m.ReturnMode {
	case message.ReturnResult:
		return executer()
	case message.ReturnNoWait:
		go func() {
			_, err := executer()
			if err != nil {
				log.Error(err)
			}
		}()
		return &reply.CallMethod{}, nil
	}
	return nil, errors.Errorf("Invalid ReturnMode #%d", m.ReturnMode)
}

func (lr *LogicRunner) executeConstructorCall(es *ExecutionState, m *message.CallConstructor, vb ValidationBehaviour) (core.Reply, error) {
	insctx := es.mainContext
	classDesc, err := lr.ArtifactManager.GetClass(insctx, m.ClassRef, nil)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't get class")
	}
	es.callContext.Class = classDesc.HeadRef()

	codeDesc := classDesc.CodeDescriptor()
	executor, err := lr.GetExecutor(codeDesc.MachineType())
	if err != nil {
		return nil, errors.Wrap(err, "no executer registered")
	}

	newData, err := executor.CallConstructor(es.callContext, *codeDesc.Ref(), m.Name, m.Arguments)
	if err != nil {
		return nil, errors.Wrap(err, "executer error")
	}

	switch m.SaveAs {
	case message.Child, message.Delegate:
		if vb.NeedSave() {
			_, err = lr.ArtifactManager.ActivateObject(
				insctx,
				Ref{}, *es.callContext.Request, m.ClassRef, m.ParentRef, m.SaveAs == message.Delegate, newData,
			)
		}
		vb.End(m.ClassRef, core.CaseRecord{
			Type: core.CaseRecordTypeResult,
			Resp: &reply.CallConstructor{Object: es.callContext.Request},
		})
		insctx.Log().Warnf("CONSTRUCTOR")
		es.Unlock()
		return &reply.CallConstructor{Object: es.callContext.Request}, err
	default:
		insctx.Log().Warnf("CONSTRUCTOR")
		es.Unlock()
		return nil, errors.New("unsupported type of save object")
	}
}

func (lr *LogicRunner) OnPulse(pulse core.Pulse) error {
	lr.consensus = make(map[Ref]*Consensus)
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
			inscontext.TODO(),
			&message.ValidateCaseBind{RecordRef: ref, CaseRecords: records, Pulse: pulse},
		)
		if err != nil {
			panic("Error while sending caseBind data to validators: " + err.Error())
		}

		temp := message.ExecutorResults{RecordRef: ref, CaseRecords: records}
		_, err = lr.MessageBus.Send(inscontext.TODO(), &temp)
		if err != nil {
			return errors.New("error while sending caseBind data to new executor")
		}
	}

	return nil
}
