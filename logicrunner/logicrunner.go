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
)

// MachineType is a type of virtual machine
type MachineType int

// Real constants of MachineType
const (
	MachineTypeBuiltin MachineType = iota
	MachineTypeGoPlugin
)

// Executor is an interface implementers for one particular machine type
type Executor interface {
	Exec(object Object, method string, args Arguments) (newObjectState []byte, methodResults Arguments, err error)
}

// LogicRunner is a general interface of contract executor
type LogicRunner struct {
	Executors [MachineTypeGoPlugin + 1]Executor
}

// NewLogicRunner is constructor for `LogicRunner`
func NewLogicRunner() (*LogicRunner, error) {
	res := LogicRunner{}

	return &res, nil
}

// RegisterExecutor registers an executor for particular `MachineType`
func (r *LogicRunner) RegisterExecutor(t MachineType, e Executor) error {
	r.Executors[int(t)] = e
	return nil
}

// GetExecutor returns an executor for the `MachineType` if it was registered (`RegisterExecutor`),
// returns error otherwise
func (r *LogicRunner) GetExecutor(t MachineType) (Executor, error) {
	if res := r.Executors[int(t)]; res != nil {
		return res, nil
	}

	return nil, errors.New("No executor registered for machine")
}

// Execute runs a method on an object, ATM just thin proxy to `GoPlugin.Exec`
func (r *LogicRunner) Execute(object Object, method string, args Arguments) ([]byte, Arguments, error) {
	e, err := r.GetExecutor(object.MachineType)
	if err != nil {
		return nil, nil, errors.Wrap(err, "no executer registered")
	}

	return e.Exec(object, method, args)
}
