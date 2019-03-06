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
	"strings"
)

type handler struct {
	machine *stateMachine
	state int
	Name string
	Params []string
	Results []string
}


func (h *handler) checkInputEventType(idx int, setEventType bool) {
	if setEventType && h.machine.InputEventType == nil {
		h.machine.InputEventType = &h.Params[idx]
	} else if h.machine.InputEventType == nil || h.Params[idx] != *h.machine.InputEventType {
		exitWithError("%s should have input event same type as Init payload", h.Name)
	}
}

func (h *handler) checkPayloadParameter(idx int) {
	if !strings.HasPrefix(h.Params[1], "*") {
		exitWithError("%s payload must be a pointer", h.Name)
	}
	if h.machine.PayloadType == nil || h.Params[idx] != *h.machine.PayloadType {
		exitWithError("%s returned payload should be same type as Init payload", h.Name)
	}
}

func (h *handler) checkInterfaceParameter(idx int) {
	if h.Params[idx] != "interface{}" {
		exitWithError("%d parameter for %s should be an interface{}", idx, h.Name)
	}
}

func (h *handler) checkAdapterResponseParameter(idx int) {
	if h.Params[idx] != "adapter.IAdapterResponse" {
		exitWithError("%d parameter for %s should be an AdapterResponse", idx, h.Name)
	}
}

func (h *handler) checkErrorParameter(idx int) {
	if h.Params[idx] != "error" {
		exitWithError("%d parameter for %s should be an error", idx, h.Name)
	}
}

func (h *handler) checkCommonHandlerReturns(setPayload bool) {
	if len(h.Results) != 3 {
		exitWithError("%s should return three values", h.Name)
	}
	if setPayload && h.machine.PayloadType == nil {
		if !strings.HasPrefix(h.Results[0], "*") {
			exitWithError("%s payload must be a pointer", h.Name)
		}
		h.machine.PayloadType = &h.Results[0]
	} else if h.machine.PayloadType == nil || h.Results[0] != *h.machine.PayloadType {
		exitWithError("%s returned payload should be same type as Init payload", h.Name)
	}
	if h.Results[1] != "common.ElUpdate" {
		exitWithError("%s returned state should be ElUpdate", h.Name)
	}
	if h.Results[2] != "error" {
		exitWithError("%s returned error must be of type error", h.Name)
	}
}

func (h *handler) checkErrorHandlerReturns() {
	if len(h.Results) != 2 {
		exitWithError("%s should return two values", h.Name)
	}
	if h.Results[0] != *h.machine.PayloadType {
		exitWithError("%s returned payload should be same type as Init payload", h.Name)
	}
	if h.Results[1] != "common.ElUpdate" {
		exitWithError("%s returned state should be ElUpdate", h.Name)
	}
}

func (h *handler) setAsState() {
	if len(h.Params) != 0 {
		exitWithError("%s state must don't have any parameters", h.Name)
	}
	if len(h.Results) != 1 || h.Results[0] != "common.ElState" {
		exitWithError("%s state should returns only common.ElState")
	}
	h.machine.States = append(h.machine.States, state{Name: h.Name})
}

func (h *handler) setAsInit() {
	if len(h.Params) != 2 {
		exitWithError("%s must have only two parameters", h.Name)
	}
	h.checkInputEventType(0, true)
	h.checkInterfaceParameter(1)
	h.checkCommonHandlerReturns(true)
	if h.machine.States[0].Transition != nil {
		exitWithError("%s init handler already exists", h.machine.Name)
	}
	h.machine.States[0].Transition = h
}

func (h *handler) setAsInitFuture() {
	if len(h.Params) != 2 {
		exitWithError("%s must have only two parameters", h.Name)
	}
	h.checkInputEventType(0, true)
	h.checkInterfaceParameter(1)
	h.checkCommonHandlerReturns(true)
	if h.machine.States[0].TransitionFuture != nil {
		exitWithError("%s init (future) handler already exists", h.machine.Name)
	}
	h.machine.States[0].TransitionFuture = h
}

func (h *handler) setAsInitPast() {
	if len(h.Params) != 2 {
		exitWithError("%s must have only two parameters", h.Name)
	}
	h.checkInputEventType(0, true)
	h.checkInterfaceParameter(1)
	h.checkCommonHandlerReturns(true)
	if h.machine.States[0].TransitionPast != nil {
		exitWithError("%s init (past) handler already exists", h.machine.Name)
	}
	h.machine.States[0].TransitionPast = h
}

func (h *handler) setAsErrorState() {
	if len(h.Params) != 3 {
		exitWithError("%s must have three parameters", h.Name)
	}
	h.checkInterfaceParameter(0)
	h.checkInterfaceParameter(1)
	h.checkErrorParameter(2)
	h.checkErrorHandlerReturns()
	h.machine.States[h.state].ErrorState = h
}

func (h *handler) setAsErrorStateFuture() {
	if len(h.Params) != 3 {
		exitWithError("%s must have three parameters", h.Name)
	}
	h.checkInterfaceParameter(0)
	h.checkInterfaceParameter(1)
	h.checkErrorParameter(2)
	h.checkErrorHandlerReturns()
	h.machine.States[h.state].ErrorStateFuture = h
}

func (h *handler) setAsErrorStatePast() {
	if len(h.Params) != 3 {
		exitWithError("%s must have three parameters", h.Name)
	}
	h.checkInterfaceParameter(0)
	h.checkInterfaceParameter(1)
	h.checkErrorParameter(2)
	h.checkErrorHandlerReturns()
	h.machine.States[h.state].ErrorStatePast = h
}

func (h *handler) setAsMigration() {
	if len(h.Params) != 2 {
		exitWithError("%s must have only two parameters", h.Name)
	}
	h.checkInputEventType(0, false)
	h.checkPayloadParameter(1)
	h.checkCommonHandlerReturns(false)
	h.machine.States[h.state].Migration = h
}

func (h *handler) setAsMigrationFuturePresent() {
	if len(h.Params) != 2 {
		exitWithError("%s must have only two parameters", h.Name)
	}
	h.checkInputEventType(0, false)
	h.checkPayloadParameter(1)
	h.checkCommonHandlerReturns(false)
	h.machine.States[h.state].MigrationFuturePresent = h
}

func (h *handler) setAsTransition() {
	if len(h.Params) < 2 {
		exitWithError("%s must have two or more parameters", h.Name)
	}
	h.checkInputEventType(0, false)
	h.checkPayloadParameter(1)
	// todo check adapter handler
	h.checkCommonHandlerReturns(false)
	h.machine.States[h.state].Transition = h
}

func (h *handler) setAsTransitionFuture() {
	if len(h.Params) < 2 {
		exitWithError("%s must have two or more parameters", h.Name)
	}
	h.checkInputEventType(0, false)
	h.checkPayloadParameter(1)
	// todo check adapter handler
	h.checkCommonHandlerReturns(false)
	h.machine.States[h.state].TransitionFuture = h
}

func (h *handler) setAsTransitionPast() {
	if len(h.Params) < 2 {
		exitWithError("%s must have two or more parameters", h.Name)
	}
	h.checkInputEventType(0, false)
	h.checkPayloadParameter(1)
	// todo check adapter handler
	h.checkCommonHandlerReturns(false)
	h.machine.States[h.state].TransitionPast = h
}

/*func (h *handler) setAsFinalization() {
	if len(h.Params) != 2 {
		exitWithError("%s must have only two parameters", h.Name)
	}
	h.checkInputEventType(0, false)
	h.checkPayloadParameter(1)
	if len(h.Results) != 0 {
		exitWithError("%s must don't have any returns", h.Name)
	}
	h.machine.States[h.state].Finalization = h
}

func (h *handler) setAsFinalizationFuture() {
	if len(h.Params) != 2 {
		exitWithError("%s must have only two parameters", h.Name)
	}
	h.checkInputEventType(0, false)
	h.checkPayloadParameter(1)
	if len(h.Results) != 0 {
		exitWithError("%s must don't have any returns", h.Name)
	}
	h.machine.States[h.state].FinalizationFuture = h
}

func (h *handler) setAsFinalizationPast() {
	if len(h.Params) != 2 {
		exitWithError("%s must have only two parameters", h.Name)
	}
	h.checkInputEventType(0, false)
	h.checkPayloadParameter(1)
	if len(h.Results) != 0 {
		exitWithError("%s must don't have any returns", h.Name)
	}
	h.machine.States[h.state].FinalizationPast = h
}*/

func (h *handler) setAsAdapterResponse() {
	if len(h.Params) != 3 {
		exitWithError("%s must have three parameters", h.Name)
	}
	h.checkInputEventType(0, false)
	h.checkPayloadParameter(1)
	// todo check respPayload here
	h.checkCommonHandlerReturns(false)
	h.machine.States[h.state].AdapterResponse = h
}

func (h *handler) setAsAdapterResponseFuture() {
	if len(h.Params) != 3 {
		exitWithError("%s must have three parameters", h.Name)
	}
	h.checkInputEventType(0, false)
	h.checkPayloadParameter(1)
	// todo check respPayload here
	h.checkCommonHandlerReturns(false)
	h.machine.States[h.state].AdapterResponseFuture = h
}

func (h *handler) setAsAdapterResponsePast() {
	if len(h.Params) != 3 {
		exitWithError("%s must have three parameters", h.Name)
	}
	h.checkInputEventType(0, false)
	h.checkPayloadParameter(1)
	// todo check respPayload here
	h.checkCommonHandlerReturns(false)
	h.machine.States[h.state].AdapterResponsePast = h
}

func (h *handler) setAsAdapterResponseError() {
	if len(h.Params) != 4 {
		exitWithError("%s must have four parameters", h.Name)
	}
	h.checkInterfaceParameter(0)
	h.checkInterfaceParameter(1)
	h.checkAdapterResponseParameter(2)
	h.checkErrorParameter(3)
	h.checkErrorHandlerReturns()
	h.machine.States[h.state].AdapterResponseError = h
}

func (h *handler) setAsAdapterResponseErrorFuture() {
	if len(h.Params) != 4 {
		exitWithError("%s must have four parameters", h.Name)
	}
	h.checkInterfaceParameter(0)
	h.checkInterfaceParameter(1)
	h.checkAdapterResponseParameter(2)
	h.checkErrorParameter(3)
	h.checkErrorHandlerReturns()
	h.machine.States[h.state].AdapterResponseErrorFuture = h
}

func (h *handler) setAsAdapterResponseErrorPast() {
	if len(h.Params) != 4 {
		exitWithError("%s must have four parameters", h.Name)
	}
	h.checkInterfaceParameter(0)
	h.checkInterfaceParameter(1)
	h.checkAdapterResponseParameter(2)
	h.checkErrorParameter(3)
	h.checkErrorHandlerReturns()
	h.machine.States[h.state].AdapterResponseErrorPast = h
}
