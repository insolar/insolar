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
	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/logicrunner/builtin"
	"github.com/pkg/errors"
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
	lr.ArtifactManager = c["core.Ledger"].(core.Ledger).GetManager()
	mr := c["core.MessageRouter"].(core.MessageRouter)

	bi := builtin.NewBuiltIn(lr.ArtifactManager, mr)
	err := lr.RegisterExecutor(core.MachineTypeBuiltin, bi)
	if err != nil {
		return err
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

	return nil, errors.New("No executor registered for machine")
}

// Execute runs a method on an object, ATM just thin proxy to `GoPlugin.Exec`
func (lr *LogicRunner) Execute(msg core.Message) *core.Response {
	lr.ArtifactManager.SetArchPref([]core.MachineType{core.MachineTypeGoPlugin})
	objDesc, err := lr.ArtifactManager.GetLatestObj(msg.Reference)
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

	executor, err := lr.GetExecutor(core.MachineTypeGoPlugin)
	if err != nil {
		return &core.Response{Error: errors.Wrap(err, "no executer registered")}
	}

	if msg.Constructor {
		newData, err := executor.CallConstructor(*codeDesc.Ref(), msg.Method, msg.Arguments)
		if err != nil {
			return &core.Response{Error: errors.Wrap(err, "executer error")}
		}
		return &core.Response{Data: newData}
	}

	newData, result, err := executor.CallMethod(*codeDesc.Ref(), data, msg.Method, msg.Arguments)
	if err != nil {
		return &core.Response{Error: errors.Wrap(err, "executer error")}
	}

	return &core.Response{Data: newData, Result: result}
}
