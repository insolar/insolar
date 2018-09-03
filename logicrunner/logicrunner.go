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
	"github.com/pkg/errors"
	"github.com/derekparker/delve/pkg/config"
	"github.com/insolar/insolar/logicrunner/builtin"
)

// Object is an inner representation of storage object for transfwering it over API
type Object struct {
	MachineType core.MachineType
	Reference   core.RecordRef
	Data        []byte
}

// LogicRunner is a general interface of contract executor
type LogicRunner struct {
	Executors       [core.MachineTypesTotalCount]core.MachineLogicExecutor
	ArtifactManager core.ArtifactManager
}

// NewLogicRunner is constructor for `LogicRunner`
func NewLogicRunner(config config.Config) (*LogicRunner, error) {
	res := LogicRunner{ArtifactManager: nil}

	return &res, nil
}

func (lr *LogicRunner) Start(c core.Components)  error {
	lr.ArtifactManager = c["core.ArtifactManager"].(core.ArtifactManager)
	mr := c["core.MessageRouter"].(core.MessageRouter)

	bi, err := builtin.NewBuiltIn(lr.ArtifactManager)
	if err != nil {
		return err
	}

	err = lr.RegisterExecutor(core.MachineTypeGoPlugin, gp)
	if err != nil {
		return err
	}
	return nil
}

func (lr *LogicRunner) Stop()  error {
	ch := make(chan error)

	for _,e := range lr.Executors {
		e.
	}

	ch <- nil
	return ch

}

// RegisterExecutor registers an executor for particular `MachineType`
func (r *LogicRunner) RegisterExecutor(t core.MachineType, e core.MachineLogicExecutor) error {
	r.Executors[int(t)] = e
	return nil
}

// GetExecutor returns an executor for the `MachineType` if it was registered (`RegisterExecutor`),
// returns error otherwise
func (r *LogicRunner) GetExecutor(t core.MachineType) (core.MachineLogicExecutor, error) {
	if res := r.Executors[int(t)]; res != nil {
		return res, nil
	}

	return nil, errors.New("No executor registered for machine")
}

// Execute runs a method on an object, ATM just thin proxy to `GoPlugin.Exec`
func (r *LogicRunner) Execute(msg core.Message) *core.Response {
	data, codeRef, err := r.ArtifactManager.Get(msg.Reference)
	if err != nil {
		return &core.Response{Error: errors.Wrap(err, "couldn't ")}
	}

	executor, err := r.GetExecutor(core.MachineTypeGoPlugin)
	if err != nil {
		return &core.Response{Error: errors.Wrap(err, "no executer registered")}
	}

	if msg.Constructor {
		newData, err := executor.CallConstructor(codeRef, msg.Method, msg.Arguments)
		if err != nil {
			return &core.Response{Error: errors.Wrap(err, "executer error")}
		}
		return &core.Response{Data: newData}
	}

	newData, result, err := executor.CallMethod(codeRef, data, msg.Method, msg.Arguments)
	if err != nil {
		return &core.Response{Error: errors.Wrap(err, "executer error")}
	}

	return &core.Response{Data: newData, Result: result}
}
