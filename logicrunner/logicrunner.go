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
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/messagerouter/types"
	"github.com/pkg/errors"
)

// Object is an inner representation of storage object for transfwering it over API
type Object struct {
	MachineType core.MachineType
	Reference   types.Reference
	Data        []byte
}

// ArtifactManager interface
type ArtifactManager interface {
	Get(ref string) (data []byte, codeRef types.Reference, err error)
}

// LogicRunner is a general interface of contract executor
type LogicRunner struct {
	Executors       [core.MachineTypesTotalCount]types.MachineLogicExecutor
	ArtifactManager ArtifactManager
}

// NewLogicRunner is constructor for `LogicRunner`
func NewLogicRunner(am ArtifactManager) (*LogicRunner, error) {
	res := LogicRunner{ArtifactManager: am}

	return &res, nil
}

// RegisterExecutor registers an executor for particular `MachineType`
func (r *LogicRunner) RegisterExecutor(t core.MachineType, e types.MachineLogicExecutor) error {
	r.Executors[int(t)] = e
	return nil
}

// GetExecutor returns an executor for the `MachineType` if it was registered (`RegisterExecutor`),
// returns error otherwise
func (r *LogicRunner) GetExecutor(t core.MachineType) (types.MachineLogicExecutor, error) {
	if res := r.Executors[int(t)]; res != nil {
		return res, nil
	}

	return nil, errors.New("No executor registered for machine")
}

// Execute runs a method on an object, ATM just thin proxy to `GoPlugin.Exec`
func (r *LogicRunner) Execute(msg types.Message) *types.Response {
	data, codeRef, err := r.ArtifactManager.Get(string(msg.Reference))
	if err != nil {
		return &types.Response{Error: errors.Wrap(err, "couldn't ")}
	}

	executor, err := r.GetExecutor(core.MachineTypeGoPlugin)
	if err != nil {
		return &types.Response{Error: errors.Wrap(err, "no executer registered")}
	}

	newData, result, err := executor.CallMethod(codeRef, data, msg.Method, msg.Arguments)
	if err != nil {
		return &types.Response{Error: errors.Wrap(err, "executer error")}
	}

	return &types.Response{Data: newData, Result: result}
}
