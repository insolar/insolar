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
	"net"
	"sync"
	"time"

	"github.com/pkg/errors"

	"bytes"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/core/message"
	"github.com/insolar/insolar/core/reply"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/logicrunner/builtin"
	"github.com/insolar/insolar/logicrunner/goplugin"
)

type Ref = core.RecordRef

// Context of one contract execution
type ExecutionContext struct {
	Pending bool           // execution moved from previous pulse
	TraceID []byte         // TraceID
	Queue   []core.Message // queued requests
}

// LogicRunner is a general interface of contract executor
type LogicRunner struct {
	Executors       [core.MachineTypesLastID]core.MachineLogicExecutor
	ArtifactManager core.ArtifactManager
	MessageBus      core.MessageBus
	machinePrefs    []core.MachineType
	Cfg             *configuration.LogicRunner
	context         map[Ref]ExecutionContext // if object exists, we are validating or executing it right now
	contextMutex    sync.Mutex

	JetCoordinator core.JetCoordinator
	NodeId         core.RecordRef

	// TODO refactor caseBind and caseBindReplays to one clear structure
	caseBind             core.CaseBind
	caseBindMutex        sync.Mutex
	caseBindReplays      map[Ref]core.CaseBindReplay
	caseBindReplaysMutex sync.Mutex
	sock                 net.Listener
}

// NewLogicRunner is constructor for LogicRunner
func NewLogicRunner(cfg *configuration.LogicRunner) (*LogicRunner, error) {
	if cfg == nil {
		return nil, errors.New("LogicRunner have nil configuration")
	}
	res := LogicRunner{
		ArtifactManager: nil,
		Cfg:             cfg,
		context:         make(map[Ref]ExecutionContext),
		caseBind:        core.CaseBind{Pulse: core.Pulse{}, Records: make(map[Ref][]core.CaseRecord)},
		caseBindReplays: make(map[Ref]core.CaseBindReplay),
	}
	return &res, nil
}

// Start starts logic runner component
func (lr *LogicRunner) Start(c core.Components) error {
	lr.ArtifactManager = c.Ledger.GetArtifactManager()
	lr.MessageBus = c.MessageBus

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

	// TODO - network rewors this
	lr.JetCoordinator = c.Ledger.GetJetCoordinator()
	lr.NodeId = c.Network.GetNodeID()

	return nil
}

// Stop stops logic runner component and its executors
func (lr *LogicRunner) Stop() error {
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

// RegisterExecutor registers an executor for particular `MachineType`
func (lr *LogicRunner) RegisterExecutor(t core.MachineType, e core.MachineLogicExecutor) error {
	lr.Executors[int(t)] = e
	return nil
}

// GetExecutor returns an executor for the `MachineType` if it was registered (`RegisterExecutor`),
// returns error otherwise
func (lr *LogicRunner) GetExecutor(t core.MachineType) (core.MachineLogicExecutor, error) {
	if res := lr.Executors[int(t)]; res != nil {
		return res, nil
	}

	return nil, errors.Errorf("No executor registered for machine %d", int(t))
}

func (lr *LogicRunner) GetContext(ref Ref) (ExecutionContext, bool) {
	lr.contextMutex.Lock()
	defer lr.contextMutex.Unlock()
	ret, ok := lr.context[ref]
	return ret, ok
}

func (lr *LogicRunner) SetContext(ref Ref, ec ExecutionContext) bool {
	lr.contextMutex.Lock()
	defer lr.contextMutex.Unlock()
	if _, ok := lr.context[ref]; ok {
		return false
	}
	lr.context[ref] = ec
	return true
}

// Execute runs a method on an object, ATM just thin proxy to `GoPlugin.Exec`
func (lr *LogicRunner) Execute(inmsg core.Message) (core.Reply, error) {
	msg, ok := inmsg.(message.IBaseLogicMessage)
	if !ok {
		return nil, errors.New("Execute( ! message.IBaseLogicMessage )")
	}

	ref := msg.GetReference()
	lr.caseBindReplaysMutex.Lock()
	cb, validate := lr.caseBindReplays[ref]
	lr.caseBindReplaysMutex.Unlock()

	var vb ValidationBehaviour
	if validate {
		vb = ValidationChecker{lr: lr, cb: cb}
	} else {
		vb = ValidationSaver{lr: lr}
	}
	isAuthorized, err := lr.JetCoordinator.IsAuthorized(vb.GetRole(), ref, lr.caseBind.Pulse.PulseNumber, *msg.GetCaller())
	reqref, err := lr.ArtifactManager.RegisterRequest(msg)
	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, errors.New("Authorization failed with error: " + err.Error())
	}

	if !isAuthorized {
		return nil, errors.New("Can't execute this object")
	}

	ctx := core.LogicCallContext{
		Caller:  msg.GetCaller(),
		Request: reqref,
		Time:    time.Now(), // TODO: probably we should take it from e
		Pulse:   lr.caseBind.Pulse,
	}

	switch m := msg.(type) {
	case *message.CallMethod:
		re, err := lr.executeMethodCall(ctx, m, vb)
		return re, err

	case *message.CallConstructor:
		re, err := lr.executeConstructorCall(ctx, m, vb)
		return re, err

	default:
		panic("Unknown e type")
	}
}

func (lr *LogicRunner) ValidateCaseBind(inmsg core.Message) (core.Reply, error) {
	msg, ok := inmsg.(*message.ValidateCaseBind)
	if !ok {
		return nil, errors.New("Execute( ! message.ValidateCaseBindInterface )")
	}

	passedStepsCount, validationError := lr.Validate(msg.GetReference(), msg.GetPulse(), msg.GetCaseRecords())
	_, err := lr.MessageBus.Send(&message.ValidationResults{
		RecordRef:        msg.GetReference(),
		PassedStepsCount: passedStepsCount,
		Error:            validationError,
	})

	return nil, err
}

func (lr *LogicRunner) ProcessValidationResults(inmsg core.Message) (core.Reply, error) {
	// Handle all validators Request
	// Do some staff if request don't come for a long time
	// Compare results of different validators and previous Executor
	return nil, nil
}

func (lr *LogicRunner) ExecutorResults(inmsg core.Message) (core.Reply, error) {
	// Coordinate this with ProcessValidationResults
	return nil, nil
}

type objectBody struct {
	Body        []byte
	Code        core.RecordRef
	Class       core.RecordRef
	MachineType core.MachineType
}

func (lr *LogicRunner) getObjectMessage(objref core.RecordRef) (*objectBody, error) {
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

	objDesc, err := lr.ArtifactManager.GetObject(objref, nil)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't get object")
	}

	classDesc, err := objDesc.ClassDescriptor(nil)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't get object's class")
	}

	codeDesc := classDesc.CodeDescriptor()
	ob := &objectBody{
		Body:        objDesc.Memory(),
		Code:        *codeDesc.Ref(),
		Class:       *classDesc.HeadRef(),
		MachineType: codeDesc.MachineType(),
	}
	lr.addObjectCaseRecord(objref, core.CaseRecord{
		Type:   core.CaseRecordTypeGetObject,
		ReqSig: HashInterface(objref),
		Resp:   ob,
	})
	return ob, nil
}

func (lr *LogicRunner) executeMethodCall(ctx core.LogicCallContext, m *message.CallMethod, vb ValidationBehaviour) (core.Reply, error) {
	ec := ExecutionContext{}
	if !lr.SetContext(m.ObjectRef, ec) {
		return nil, errors.New("Method already executing")
	}
	vb.Begin(m.ObjectRef, core.CaseRecord{
		Type: core.CaseRecordTypeStart,
		Resp: m,
	})

	objbody, err := lr.getObjectMessage(m.ObjectRef)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't get object message")
	}

	ctx.Callee = &m.ObjectRef
	ctx.Class = &objbody.Class
	vb.ModifyContext(&ctx)

	executor, err := lr.GetExecutor(objbody.MachineType)
	if err != nil {
		return nil, errors.Wrap(err, "no executor registered")
	}

	executer := func() (*reply.CallMethod, error) {
		defer func() {
			lr.contextMutex.Lock()
			defer lr.contextMutex.Unlock()
			delete(lr.context, m.ObjectRef)
		}()
		newData, result, err := executor.CallMethod(
			&ctx, objbody.Code, objbody.Body, m.Method, m.Arguments,
		)
		if err != nil {
			return nil, errors.Wrap(err, "executor error")
		}

		// TODO: deactivation should be handled way better here
		if vb.NeedSave() && lr.lastObjectCaseRecord(m.ObjectRef).Type != core.CaseRecordTypeDeactivateObject {
			_, err = lr.ArtifactManager.UpdateObject(
				core.RecordRef{}, *ctx.Request, m.ObjectRef, newData,
			)
			if err != nil {
				return nil, errors.Wrap(err, "couldn't update object")
			}
		}

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

func (lr *LogicRunner) executeConstructorCall(ctx core.LogicCallContext, m *message.CallConstructor, vb ValidationBehaviour) (core.Reply, error) {
	ec := ExecutionContext{}
	if !lr.SetContext(m.GetRequest(), ec) {
		return nil, errors.New("Constructor already executing by you")
	}
	vb.Begin(m.ClassRef, core.CaseRecord{
		Type: core.CaseRecordTypeStart,
		Resp: m,
	})

	classDesc, err := lr.ArtifactManager.GetClass(m.ClassRef, nil)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't get class")
	}
	ctx.Class = classDesc.HeadRef()

	codeDesc := classDesc.CodeDescriptor()
	executor, err := lr.GetExecutor(codeDesc.MachineType())
	if err != nil {
		return nil, errors.Wrap(err, "no executer registered")
	}

	newData, err := executor.CallConstructor(&ctx, *codeDesc.Ref(), m.Name, m.Arguments)
	if err != nil {
		return nil, errors.Wrap(err, "executer error")
	}

	defer func() {
		lr.contextMutex.Lock()
		defer lr.contextMutex.Unlock()
		delete(lr.context, m.GetRequest())
	}()

	switch m.SaveAs {
	case message.Child:
		log.Warn()
		log.Warnf("M = %+v", m)
		if vb.NeedSave() {
			_, err = lr.ArtifactManager.ActivateObject(
				core.RecordRef{}, *ctx.Request, m.ClassRef, m.ParentRef, newData,
			)
		}
		vb.End(m.ClassRef, core.CaseRecord{
			Type: core.CaseRecordTypeResult,
			Resp: &reply.CallConstructor{Object: ctx.Request},
		})

		return &reply.CallConstructor{Object: ctx.Request}, err
	case message.Delegate:
		if vb.NeedSave() {
			_, err = lr.ArtifactManager.ActivateObjectDelegate(
				core.RecordRef{}, *ctx.Request, m.ClassRef, m.ParentRef, newData,
			)
		}
		vb.End(m.ClassRef, core.CaseRecord{
			Type: core.CaseRecordTypeResult,
			Resp: &reply.CallConstructor{Object: ctx.Request},
		})

		return &reply.CallConstructor{Object: ctx.Request}, err
	default:
		return nil, errors.New("unsupported type of save object")
	}
}

func (lr *LogicRunner) OnPulse(pulse core.Pulse) error {
	// start of new Pulse, lock CaseBind data, copy it, clean original, unlock original
	objectsRecords := lr.refreshCaseBind(pulse)

	if len(objectsRecords) == 0 {
		return nil
	}

	// send copy for validation
	for ref, records := range objectsRecords {
		_, err := lr.MessageBus.Send(&message.ValidateCaseBind{RecordRef: ref, CaseRecords: records, Pulse: pulse})
		if err != nil {
			panic("Error while sending caseBind data to validators: " + err.Error())
		}

		temp := message.ExecutorResults{RecordRef: ref, CaseRecords: records}
		_, err = lr.MessageBus.Send(&temp)
		if err != nil {
			return errors.New("error while sending caseBind data to new executor")
		}
	}

	return nil
}

// refreshCaseBind lock CaseBind data, copy it, clean original, unlock original, return copy
func (lr *LogicRunner) refreshCaseBind(pulse core.Pulse) map[core.RecordRef][]core.CaseRecord {
	lr.caseBindMutex.Lock()
	defer lr.caseBindMutex.Unlock()

	oldObjectsRecords := lr.caseBind.Records

	lr.caseBind = core.CaseBind{
		Pulse:   pulse,
		Records: make(map[core.RecordRef][]core.CaseRecord),
	}

	return oldObjectsRecords
}
