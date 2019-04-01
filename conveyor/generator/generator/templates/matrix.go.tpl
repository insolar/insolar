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
	"github.com/insolar/insolar/conveyor/fsm"
	{{range .}}"{{fileToImport .File}}"
	{{end}}
)

type StateMachineSet struct{
	stateMachines []StateMachine
}

func newStateMachineSet() *StateMachineSet {
	return &StateMachineSet{
		stateMachines: make([]StateMachine, 1),
	}
}

func (s *StateMachineSet) addMachine(machine StateMachine) {
	s.stateMachines = append(s.stateMachines, machine)
}

func ( s *StateMachineSet ) GetStateMachineByID(id fsm.ID) StateMachine{
	return s.stateMachines[id]
}

type Matrix struct {
	future *StateMachineSet
	present *StateMachineSet
	past *StateMachineSet
}

const (
	{{range $i, $m := .}}{{if (isNull $i)}}{{$m.Name}}  fsm.ID = iota + 1
	{{else}}{{$m.Name}}
	{{end}}{{end}}
)

func NewMatrix() *Matrix {
	m := Matrix{
		future: newStateMachineSet(),
		present: newStateMachineSet(),
		past: newStateMachineSet(),
	}
	{{range .}}
	m.future.addMachine({{.Package}}.Raw{{.Name}}FutureFactory())
	m.present.addMachine({{.Package}}.Raw{{.Name}}PresentFactory())
	m.past.addMachine({{.Package}}.Raw{{.Name}}PastFactory())
	{{end}}
	return &m
}

func (m *Matrix) GetInitialStateMachine() StateMachine {
	return m.present.stateMachines[Initial]
}

func (m *Matrix) GetFutureConfig() SetAccessor{
	return m.future
}

func (m *Matrix) GetPresentConfig() SetAccessor{
	return m.present
}

func (m *Matrix) GetPastConfig() SetAccessor{
	return m.past
}
