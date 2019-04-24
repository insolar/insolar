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

package generator

import (
	"runtime"
	"strings"

	"github.com/insolar/insolar/conveyor/fsm"
)

type StateMachine struct {
	Package              string
	File                 string
	Name                 string
	InputEventType       *string
	PayloadType          *string
	States               []*state
	AdapterHelperCatalog map[string]string
	parsedImports        map[string]string
	usedImports          map[string]string
}

func (g *Generator) AddMachine(name string) *StateMachine {
	_, file, _, _ := runtime.Caller(1)
	pkg, imports := parseSource(file)
	machine := &StateMachine{
		File:                 file,
		Package:              pkg,
		Name:                 name,
		AdapterHelperCatalog: g.adapterHelperCatalog,
		parsedImports:        imports,
		usedImports:          make(map[string]string),
	}
	g.stateMachines = append(g.stateMachines, machine)
	return machine
}

func (m *StateMachine) setImportFor(param string) {
	if param[0] == '*' {
		param = param[1:]
	}
	if dotIndex := strings.Index(param, "."); dotIndex > 0 {
		pkgName := param[:dotIndex]
		if pkgName == m.Package {
			return
		}
		if _, exists := m.parsedImports[pkgName]; !exists {
			exitWithError("Couldn't find import for package %s\n", pkgName)
		}
		m.usedImports[pkgName] = m.parsedImports[pkgName]
	}
}

func (m *StateMachine) GetImports() []string {
	var result []string
	for _, v := range m.usedImports {
		result = append(result, v)
	}
	return result
}

func (m *StateMachine) GetInputEventType() string {
	if strings.HasPrefix(*m.InputEventType, m.Package) {
		return (*m.InputEventType)[len(m.Package)+1:]
	}
	return *m.InputEventType
}

func (m *StateMachine) GetPayloadType() string {
	if strings.HasPrefix(*m.PayloadType, "*"+m.Package) {
		return "*" + (*m.PayloadType)[len(m.Package)+2:]
	}
	return *m.PayloadType
}

func (m *StateMachine) createStateUnlessExists(current fsm.ElementState, returned []fsm.ElementState) {
	// todo check is state not contains stateMachine
	for _, s := range append(returned, current) {
		if len(m.States) <= int(s) {
			for i := len(m.States); i <= int(s); i++ {
				m.States = append(m.States, &state{})
			}
		}
		if m.States[int(s)] == nil {
			m.States[int(s)] = &state{}
		}
	}
}

func (m *StateMachine) addHandler(s fsm.ElementState, ht handlerType, ps PulseState, h *handler) {
	if m.States[s].handlers[ps] == nil {
		m.States[s].handlers[ps] = make(map[handlerType]*handler)
	}
	if m.States[s].handlers[ps][ht] != nil {
		exitWithError("handler already exists %s", h.Name)
	}
	m.States[s].handlers[ps][ht] = h
}

func (m *StateMachine) Init(handlerFunc interface{}, states ...fsm.ElementState) *StateMachine {
	m.createStateUnlessExists(0, states)
	m.addHandler(0, Transition, Present, newInitHandler(m, handlerFunc, states))
	return m
}

func (m *StateMachine) InitFuture(handlerFunc interface{}, states ...fsm.ElementState) *StateMachine {
	m.createStateUnlessExists(0, states)
	m.addHandler(0, Transition, Future, newInitHandler(m, handlerFunc, states))
	return m
}

func (m *StateMachine) InitPast(handlerFunc interface{}, states ...fsm.ElementState) *StateMachine {
	m.createStateUnlessExists(0, states)
	m.addHandler(0, Transition, Past, newInitHandler(m, handlerFunc, states))
	return m
}

func (m *StateMachine) Transition(state fsm.ElementState, handlerFunc interface{}, states ...fsm.ElementState) *StateMachine {
	m.createStateUnlessExists(state, states)
	m.addHandler(state, Transition, Present, newHandler(m, handlerFunc, states))
	return m
}

func (m *StateMachine) TransitionFuture(state fsm.ElementState, handlerFunc interface{}, states ...fsm.ElementState) *StateMachine {
	m.createStateUnlessExists(state, states)
	m.addHandler(state, Transition, Future, newHandler(m, handlerFunc, states))
	return m
}

func (m *StateMachine) TransitionPast(state fsm.ElementState, handlerFunc interface{}, states ...fsm.ElementState) *StateMachine {
	m.createStateUnlessExists(state, states)
	m.addHandler(state, Transition, Past, newHandler(m, handlerFunc, states))
	return m
}

func (m *StateMachine) Migration(state fsm.ElementState, handlerFunc interface{}, states ...fsm.ElementState) *StateMachine {
	m.createStateUnlessExists(state, states)
	m.addHandler(state, Migration, Past, newHandler(m, handlerFunc, states))
	return m
}

func (m *StateMachine) MigrationFuturePresent(state fsm.ElementState, handlerFunc interface{}, states ...fsm.ElementState) *StateMachine {
	m.createStateUnlessExists(state, states)
	m.addHandler(state, Migration, Present, newHandler(m, handlerFunc, states))
	return m
}

func (m *StateMachine) AdapterResponse(state fsm.ElementState, handlerFunc interface{}, states ...fsm.ElementState) *StateMachine {
	m.createStateUnlessExists(state, states)
	handler := newHandler(m, handlerFunc, states)
	m.addHandler(state, AdapterResponse, Present, handler)
	m.setImportFor(handler.params[responseAdapterParamIndex])
	return m
}

func (m *StateMachine) AdapterResponseFuture(state fsm.ElementState, handlerFunc interface{}, states ...fsm.ElementState) *StateMachine {
	m.createStateUnlessExists(state, states)
	handler := newHandler(m, handlerFunc, states)
	m.addHandler(state, AdapterResponse, Future, handler)
	m.setImportFor(handler.params[responseAdapterParamIndex])
	return m
}

func (m *StateMachine) AdapterResponsePast(state fsm.ElementState, handlerFunc interface{}, states ...fsm.ElementState) *StateMachine {
	m.createStateUnlessExists(state, states)
	handler := newHandler(m, handlerFunc, states)
	m.addHandler(state, AdapterResponse, Past, handler)
	m.setImportFor(handler.params[responseAdapterParamIndex])
	return m
}
