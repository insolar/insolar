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
	"reflect"
	"runtime"
	"strings"

	"github.com/insolar/insolar/conveyor/fsm"
)

type handlerType uint

const (
	adapterHelperParamIndex   = 4
	responseAdapterParamIndex = 4
)

const (
	Transition = handlerType(iota)
	Migration
	AdapterResponse
)

type handler struct {
	machine    *StateMachine
	importPath string
	Name       string
	params     []string
	results    []string
	states     []fsm.ElementState
}

func newHandler(machine *StateMachine, f interface{}, states []fsm.ElementState) *handler {
	tp := reflect.TypeOf(f)
	if tp.Kind().String() != "func" {
		exitWithError("[%s %s] handler must be function", machine.Name, tp.Name())
	}

	fullName := runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
	lastDotIndex := strings.LastIndex(fullName, ".")

	handler := &handler{
		machine:    machine,
		importPath: fullName[:lastDotIndex],
		Name:       fullName[lastDotIndex+1:],
		params:     make([]string, tp.NumIn()),
		results:    make([]string, tp.NumOut()),
		states:     states,
	}

	// check common input types
	if tp.NumIn() < 1 || tp.In(0).String() != "context.Context" {
		exitWithError("[%s %s] first parameter should be context.Context\n", machine.Name, handler.Name)
	}
	if tp.NumIn() < 2 || tp.In(1).String() != "fsm.SlotElementHelper" {
		exitWithError("[%s %s] second parameter should be fsm.SlotElementHelper\n", machine.Name, handler.Name)
	}
	// check common return types
	if tp.NumOut() < 1 || tp.Out(0).String() != "fsm.ElementState" {
		exitWithError("[%s %s] first returned value should be fsm.ElementState\n", machine.Name, handler.Name)
	}

	for i := 0; i < tp.NumIn(); i++ {
		handler.params[i] = tp.In(i).String()
	}

	for i := 0; i < tp.NumOut(); i++ {
		handler.results[i] = tp.Out(i).String()
	}
	return handler
}

func newInitHandler(machine *StateMachine, f interface{}, states []fsm.ElementState) *handler {
	h := newHandler(machine, f, states)
	const inputEventTypeParamIndex = 2
	const payloadTypeResultIndex = 1
	// todo check input and results len
	if h.machine.InputEventType == nil {
		h.machine.setImportFor(h.params[inputEventTypeParamIndex])
		h.machine.InputEventType = &h.params[inputEventTypeParamIndex]
	}
	if h.machine.PayloadType == nil {
		h.machine.setImportFor(h.results[payloadTypeResultIndex])
		h.machine.PayloadType = &h.results[payloadTypeResultIndex]
	}
	return h
}

func (h *handler) GetResponseAdapterType() string {
	responseAdapterParam := h.params[responseAdapterParamIndex]
	return responseAdapterParam
}

func (h *handler) GetAdapterHelperType() *string {
	if len(h.params) <= adapterHelperParamIndex {
		return nil
	}

	const adapterPrefix = "adapter."
	if strings.HasPrefix(h.params[adapterHelperParamIndex], adapterPrefix) {
		result := h.params[adapterHelperParamIndex][len(adapterPrefix):]
		return &result
	}
	return &h.params[adapterHelperParamIndex]
}
