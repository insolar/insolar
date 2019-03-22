/*
*    Copyright 2019 Insolar Technologies
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

package matrix

import (
	"github.com/insolar/insolar/conveyor/generator/state_machines/sample"
	"github.com/insolar/insolar/conveyor/interfaces/statemachine"
)

const numPulseStates = 3

type StateMachineSet struct {
	stateMachines []statemachine.StateMachine
}

func (s *StateMachineSet) GetStateMachineByID(id int) statemachine.StateMachine {
	return s.stateMachines[id]
}

type Matrix struct {
	matrix [numPulseStates]StateMachineSet
}

type MachineType int

const (
	TestStateMachine MachineType = iota + 1
)

func NewMatrix() *Matrix {
	m := Matrix{}

	// Fill m.matrix[i][0] with empty state machine, since 0 - is state of completion of state machine
	var emptyObject statemachine.StateMachine
	for i := 0; i < numPulseStates; i++ {
		m.matrix[i].stateMachines = append(m.matrix[i].stateMachines, emptyObject)
	}

	smsTestStateMachine := sample.RawTestStateMachineFactory()
	for i := 0; i < numPulseStates; i++ {
		m.matrix[i].stateMachines = append(m.matrix[i].stateMachines, smsTestStateMachine[i])
	}

	return &m
}

func (m *Matrix) GetInitialStateMachine() statemachine.StateMachine {
	return m.matrix[1].stateMachines[1]
}

func (m *Matrix) GetFutureConfig() statemachine.SetAccessor {
	return &m.matrix[0]
}

func (m *Matrix) GetPresentConfig() statemachine.SetAccessor {
	return &m.matrix[1]
}

func (m *Matrix) GetPastConfig() statemachine.SetAccessor {
	return &m.matrix[2]
}
