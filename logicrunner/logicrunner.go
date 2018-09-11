/*
 *    Copyright 2018 INS Ecosystem
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
	"github.com/pkg/errors"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/logicrunner/builtin"
	"github.com/insolar/insolar/logicrunner/goplugin"
	"github.com/insolar/insolar/messagerouter/message"
)

// LogicRunner is a general interface of contract executor
type LogicRunner struct {
	Executors       [core.MachineTypesLastID]core.MachineLogicExecutor
	ArtifactManager core.ArtifactManager
	Cfg             configuration.LogicRunner
}

// NewLogicRunner is constructor for `LogicRunner`
func NewLogicRunner(cfg configuration.LogicRunner) (*LogicRunner, error) {
	res := LogicRunner{
		ArtifactManager: nil,
		Cfg:             cfg,
	}
	return &res, nil
}

// Start starts logic runner component
func (lr *LogicRunner) Start(c core.Components) error {
	am := c["core.Ledger"].(core.Ledger).GetManager()
	mr := c["core.MessageRouter"].(core.MessageRouter)
	lr.ArtifactManager = am

	if lr.Cfg.BuiltIn != nil {
		bi := builtin.NewBuiltIn(mr, am)
		if err := lr.RegisterExecutor(core.MachineTypeBuiltin, bi); err != nil {
			return err
		}
	}

	if lr.Cfg.GoPlugin != nil {
		gp, err := goplugin.NewGoPlugin(lr.Cfg.GoPlugin, mr, am)
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

type withCodeDescriptor interface {
	CodeDescriptor() (core.CodeDescriptor, error)
}

func (lr *LogicRunner) executorFromDescriptor(from withCodeDescriptor) (core.MachineLogicExecutor, error) {
	codeDesc, err := from.CodeDescriptor()
	if err != nil {
		return nil, errors.Wrap(err, "couldn't get code descriptor")
	}

	mt, err := codeDesc.MachineType()
	if err != nil {
		return nil, errors.Wrap(err, "couldn't get machine type")
	}

	executor, err := lr.GetExecutor(mt)
	if err != nil {
		return nil, errors.Wrap(err, "no executer registered")
	}

	return executor, nil
}

// Execute runs a method on an object, ATM just thin proxy to `GoPlugin.Exec`
func (lr *LogicRunner) Execute(msg core.Message) *core.Response {
	lr.ArtifactManager.SetArchPref(
		[]core.MachineType{
			core.MachineTypeBuiltin,
			core.MachineTypeGoPlugin,
		},
	)

	switch m := msg.(type) {
	case *message.CallMethodMessage:
		objDesc, err := lr.ArtifactManager.GetLatestObj(m.ObjectRef)
		if err != nil {
			return &core.Response{Error: errors.Wrap(err, "couldn't get object")}
		}

		data, err := objDesc.Memory()
		if err != nil {
			return &core.Response{Error: errors.Wrap(err, "couldn't get object's data")}
		}

		codeDesc, err := objDesc.CodeDescriptor()
		if err != nil {
			return &core.Response{Error: errors.Wrap(err, "couldn't get object's code descriptor")}
		}

		executor, err := lr.executorFromDescriptor(objDesc)
		if err != nil {
			return &core.Response{Error: errors.Wrap(err, "couldn't get executor")}
		}

		newData, result, err := executor.CallMethod(*codeDesc.Ref(), data, m.Method, m.Arguments)
		if err != nil {
			return &core.Response{Error: errors.Wrap(err, "executer error")}
		}

		_, err = lr.ArtifactManager.UpdateObj(
			core.RecordRef{}, core.RecordRef{}, m.ObjectRef, newData,
		)
		if err != nil {
			return &core.Response{Error: errors.Wrap(err, "couldn't update object")}
		}

		return &core.Response{Data: newData, Result: result}

	case *message.CallConstructorMessage:
		classDesc, err := lr.ArtifactManager.GetLatestClass(m.ClassRef)
		if err != nil {
			return &core.Response{Error: errors.Wrap(err, "couldn't get class")}
		}

		codeDesc, err := classDesc.CodeDescriptor()
		if err != nil {
			return &core.Response{Error: errors.Wrap(err, "couldn't get class's code descriptor")}
		}

		executor, err := lr.executorFromDescriptor(classDesc)
		if err != nil {
			return &core.Response{Error: errors.Wrap(err, "couldn't get executor")}
		}

		newData, err := executor.CallConstructor(*codeDesc.Ref(), m.Name, m.Arguments)
		if err != nil {
			return &core.Response{Error: errors.Wrap(err, "executer error")}
		}

		return &core.Response{Data: newData}

	case *message.DelegateMessage:
		// TODO: should be InjectDelegate
		ref, err := lr.ArtifactManager.ActivateObjDelegate(
			core.RecordRef{}, core.RecordRef{}, m.Class, m.Into, m.Body,
		)
		if err != nil {
			return &core.Response{Error: errors.Wrap(err, "couldn't save new object")}
		}
		return &core.Response{Data: []byte(ref.String())}

	default:
		panic("Unknown message type")
	}
}
