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
	"log"
	"reflect"
	"runtime"
	"strings"

	"github.com/insolar/insolar/conveyor/interfaces/fsm"
)

type handlerType uint
const (
	Transition = handlerType(iota)
	Migration
	AdapterResponse
)

type handler struct {
	machine 	*StateMachine
	importPath	string
	funcName	string
	params		[]string
	results		[]string
	states		[]fsm.ElementState
}

func NewHandler(machine *StateMachine, f interface{}, states []fsm.ElementState) *handler {
	tp := reflect.TypeOf(f)
	if tp.Kind().String() != "func" {
		log.Fatal("handler must be function")
	}

	fullName := runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
	lastDotIndex := strings.LastIndex(fullName, ".")

	handler := &handler{
		machine: machine,
		importPath: fullName[:lastDotIndex],
		funcName: fullName[lastDotIndex + 1:],
		params: make([]string, tp.NumIn()),
		results: make([]string, tp.NumOut()),
		states: states,
	}

	// check common input types
	if tp.NumIn() < 1 || tp.In(0).String() != "context.Context" {
		log.Fatal("["+handler.funcName+"] first parameter should be context.Context")
	}
	if tp.NumIn() < 2 || tp.In(1).String() != "slot.SlotElementHelper" {
		log.Fatal("["+handler.funcName+"] second parameter should be slot.SlotElementHelper")
	}
	// check common return types
	if tp.NumOut() < 1 || tp.Out(0).String() != "fsm.ElementState" {
		log.Fatal("["+handler.funcName+"] first returned value should be fsm.ElementState")
	}

	for i := 0; i < tp.NumIn(); i++ {
		handler.params[i] = tp.In(i).String()
	}

	for i := 0; i < tp.NumOut(); i++ {
		handler.results[i] = tp.Out(i).String()
	}
	return handler
}

func NewInitHandler(machine *StateMachine, f interface{}, states []fsm.ElementState) *handler {
	h := NewHandler(machine, f, states)
	// todo check input and results len
	if h.machine.InputEventType == nil {
		h.machine.InputEventType = &h.params[2]
	}
	if h.machine.PayloadType == nil {
		h.machine.PayloadType = &h.results[1]
	}
	return h
}

func (h *handler) GetName() string {
	return h.funcName
}

func (h *handler) GetResponseAdapterType() string {
	return h.params[4]
}

