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
	"time"

	"github.com/pkg/errors"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/eventbus/event"
	"github.com/insolar/insolar/eventbus/reaction"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/logicrunner/builtin"
	"github.com/insolar/insolar/logicrunner/goplugin"
)

// LogicRunner is a general interface of contract executor
type LogicRunner struct {
	Executors       [core.MachineTypesLastID]core.MachineLogicExecutor
	ArtifactManager core.ArtifactManager
	EventBus        core.EventBus
	Cfg             configuration.LogicRunner
}

// NewLogicRunner is constructor for LogicRunner
func NewLogicRunner(cfg configuration.LogicRunner) (*LogicRunner, error) {
	res := LogicRunner{
		ArtifactManager: nil,
		Cfg:             cfg,
	}
	return &res, nil
}

// Start starts logic runner component
func (lr *LogicRunner) Start(c core.Components) error {
	am := c.Ledger.GetArtifactManager()
	lr.ArtifactManager = am
	eventBus := c.EventBus
	lr.EventBus = eventBus

	if lr.Cfg.BuiltIn != nil {
		bi := builtin.NewBuiltIn(eventBus, am)
		if err := lr.RegisterExecutor(core.MachineTypeBuiltin, bi); err != nil {
			return err
		}
	}

	if lr.Cfg.GoPlugin != nil {
		gp, err := goplugin.NewGoPlugin(lr.Cfg.GoPlugin, eventBus, am)
		if err != nil {
			return err
		}
		if err := lr.RegisterExecutor(core.MachineTypeGoPlugin, gp); err != nil {
			return err
		}
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
func (lr *LogicRunner) Execute(e core.Event) (core.Reaction, error) {
	ctx := core.LogicCallContext{
		Time: time.Now(), // TODO: probably we should take it from e
	}

	machinePref := []core.MachineType{
		core.MachineTypeBuiltin,
		core.MachineTypeGoPlugin,
	}

	switch m := e.(type) {
	case *event.CallMethodEvent:
		return lr.executeMethodCall(ctx, m, machinePref)

	case *event.CallConstructorEvent:
		classDesc, err := lr.ArtifactManager.GetClass(m.ClassRef, nil)
		if err != nil {
			return nil, errors.Wrap(err, "couldn't get class")
		}
		ctx.Class = classDesc.HeadRef()

		codeDesc, err := classDesc.CodeDescriptor(machinePref)
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

		return &reaction.CommonReaction{Data: newData}, nil

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

func (lr *LogicRunner) getObjectEvent(objref core.RecordRef, machinePref []core.MachineType) (*objectBody, error) {
	objDesc, err := lr.ArtifactManager.GetObject(objref, nil)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't get object")
	}

	classDesc, err := objDesc.ClassDescriptor(nil)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't get object's class")
	}

	codeDesc, err := classDesc.CodeDescriptor(machinePref)
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

func (lr *LogicRunner) executeMethodCall(ctx core.LogicCallContext, e *event.CallMethodEvent, machinePref []core.MachineType) (core.Reaction, error) {
	objbody, err := lr.getObjectEvent(e.ObjectRef, machinePref)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't get object")
	}

	ctx.Callee = &e.ObjectRef
	ctx.Class = &objbody.Class

	executor, err := lr.GetExecutor(objbody.MachineType)
	if err != nil {
		return nil, errors.Wrap(err, "no executer registered")
	}

	executer := func() (*reaction.CommonReaction, error) {
		newData, result, err := executor.CallMethod(
			&ctx, objbody.Code, objbody.Body, e.Method, e.Arguments,
		)
		if err != nil {
			return nil, errors.Wrap(err, "executer error")
		}

		_, err = lr.ArtifactManager.UpdateObject(
			core.RecordRef{}, core.RecordRef{}, e.ObjectRef, newData,
		)
		if err != nil {
			return nil, errors.Wrap(err, "couldn't update object")
		}

		return &reaction.CommonReaction{Data: newData, Result: result}, nil
	}

	switch e.ReturnMode {
	case event.ReturnResult:
		return executer()
	case event.ReturnNoWait:
		go func() {
			_, err := executer()
			if err != nil {
				log.Error(err)
			}
		}()
		return &reaction.CommonReaction{}, nil
	}
	return nil, errors.Errorf("Invalid ReturnMode #%d", e.ReturnMode)
}
