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

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/core/message"
	"github.com/insolar/insolar/core/reply"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/logicrunner/builtin"
	"github.com/insolar/insolar/logicrunner/goplugin"
)

// LogicRunner is a general interface of contract executor
type LogicRunner struct {
	Executors            [core.MachineTypesLastID]core.MachineLogicExecutor
	ArtifactManager      core.ArtifactManager
	MessageBus           core.MessageBus
	machinePrefs         []core.MachineType
	Cfg                  *configuration.LogicRunner
	caseBind             core.CaseBind
	caseBindMutex        sync.Mutex
	caseBindReplays      map[core.RecordRef]core.CaseBindReplay
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
		caseBindReplays: make(map[core.RecordRef]core.CaseBindReplay),
	}
	return &res, nil
}

// Start starts logic runner component
func (lr *LogicRunner) Start(c core.Components) error {
	am := c.Ledger.GetArtifactManager()
	lr.ArtifactManager = am
	messageBus := c.MessageBus
	lr.MessageBus = messageBus

	if lr.Cfg.BuiltIn != nil {
		bi := builtin.NewBuiltIn(messageBus, am)
		if err := lr.RegisterExecutor(core.MachineTypeBuiltin, bi); err != nil {
			return err
		}
		lr.machinePrefs = append(lr.machinePrefs, core.MachineTypeBuiltin)
	}

	if lr.Cfg.GoPlugin != nil {
		if lr.Cfg.RPCListen != "" {
			StartRPC(lr)
		}

		gp, err := goplugin.NewGoPlugin(lr.Cfg, messageBus, am)
		if err != nil {
			return err
		}
		if err := lr.RegisterExecutor(core.MachineTypeGoPlugin, gp); err != nil {
			return err
		}
		lr.machinePrefs = append(lr.machinePrefs, core.MachineTypeGoPlugin)
	}

	// TODO: use separate handlers
	if err := messageBus.Register(core.TypeCallMethod, lr.Execute); err != nil {
		return err
	}
	if err := messageBus.Register(core.TypeCallConstructor, lr.Execute); err != nil {
		return err
	}

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

	ctx := core.LogicCallContext{
		Caller: msg.GetCaller(),
		Time:   time.Now(), // TODO: probably we should take it from e
		Pulse:  lr.caseBind.Pulse,
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

type objectBody struct {
	Body        []byte
	Code        core.RecordRef
	Class       core.RecordRef
	MachineType core.MachineType
}

func (lr *LogicRunner) getObjectMessage(objref core.RecordRef) (*objectBody, error) {
	objDesc, err := lr.ArtifactManager.GetObject(objref, nil)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't get object")
	}

	classDesc, err := objDesc.ClassDescriptor(nil)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't get object's class")
	}

	codeDesc, err := classDesc.CodeDescriptor(lr.machinePrefs)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't get object's code descriptor")
	}

	return &objectBody{
		Body:        objDesc.Memory(),
		Code:        *codeDesc.Ref(),
		Class:       *classDesc.HeadRef(),
		MachineType: codeDesc.MachineType(),
	}, nil
}

func (lr *LogicRunner) executeMethodCall(ctx core.LogicCallContext, e *message.CallMethod, vb ValidationBehaviour) (core.Reply, error) {
	vb.Begin(e.ObjectRef, core.CaseRecord{
		Type: core.CaseRecordTypeStart,
		Resp: e,
	})

	objbody, err := lr.getObjectMessage(e.ObjectRef)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't get object message")
	}

	ctx.Callee = &e.ObjectRef
	ctx.Class = &objbody.Class
	vb.ModifyContext(&ctx)

	executor, err := lr.GetExecutor(objbody.MachineType)
	if err != nil {
		return nil, errors.Wrap(err, "no executor registered")
	}

	executer := func() (*reply.CallMethod, error) {
		newData, result, err := executor.CallMethod(
			&ctx, objbody.Code, objbody.Body, e.Method, e.Arguments,
		)
		if err != nil {
			return nil, errors.Wrap(err, "executor error")
		}

		if vb.NeedSave() {
			_, err = lr.ArtifactManager.UpdateObject(
				core.RecordRef{}, core.RecordRef{}, e.ObjectRef, newData,
			)
			if err != nil {
				return nil, errors.Wrap(err, "couldn't update object")
			}
		}

		re := &reply.CallMethod{Data: newData, Result: result}

		vb.End(e.ObjectRef, core.CaseRecord{
			Type: core.CaseRecordTypeResult,
			Resp: re,
		})

		return re, nil
	}

	switch e.ReturnMode {
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
	return nil, errors.Errorf("Invalid ReturnMode #%d", e.ReturnMode)
}

func (lr *LogicRunner) executeConstructorCall(ctx core.LogicCallContext, m *message.CallConstructor, vb ValidationBehaviour) (core.Reply, error) {
	vb.Begin(m.ClassRef, core.CaseRecord{
		Type: core.CaseRecordTypeStart,
		Resp: m,
	})

	classDesc, err := lr.ArtifactManager.GetClass(m.ClassRef, nil)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't get class")
	}
	ctx.Class = classDesc.HeadRef()

	codeDesc, err := classDesc.CodeDescriptor(lr.machinePrefs)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't get class's code descriptor")
	}

	executor, err := lr.GetExecutor(codeDesc.MachineType())
	if err != nil {
		return nil, errors.Wrap(err, "no executer registered")
	}

	newData, err := executor.CallConstructor(&ctx, *codeDesc.Ref(), m.Name, m.Arguments)
	if err != nil {
		return nil, errors.Wrap(err, "executer error")
	}

	switch m.SaveAs {
	case message.Child:
		ref, err := lr.ArtifactManager.ActivateObject(
			core.RecordRef{}, core.RandomRef(), m.ClassRef, m.ParentRef, newData,
		)

		vb.End(m.ClassRef, core.CaseRecord{
			Type: core.CaseRecordTypeResult,
			Resp: &reply.CallConstructor{Object: ref},
		})

		return &reply.CallConstructor{Object: ref}, err
	case message.Delegate:
		ref, err := lr.ArtifactManager.ActivateObjectDelegate(
			core.RecordRef{}, core.RandomRef(), m.ClassRef, m.ParentRef, newData,
		)

		vb.End(m.ClassRef, core.CaseRecord{
			Type: core.CaseRecordTypeResult,
			Resp: &reply.CallConstructor{Object: ref},
		})

		return &reply.CallConstructor{Object: ref}, err
	default:
		return nil, errors.New("unsupported type of save object")
	}
}

func (lr *LogicRunner) OnPulse(pulse core.Pulse) error {
	lr.caseBindMutex.Lock()
	lr.caseBind = core.CaseBind{
		Pulse:   pulse,
		Records: make(map[core.RecordRef][]core.CaseRecord),
	}
	lr.caseBindMutex.Unlock()
	return nil
}
