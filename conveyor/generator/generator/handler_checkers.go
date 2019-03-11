/*
 *    Copyright 2019 INS Ecosystem
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

// Some default error checks for different types of handlers
// Handler types are: Init, ErrorState, Migration, Transition, etc

func (h *handler) checkInitHandler() {
	if len(h.params) != 2 {
		exitWithError("%s must have only two parameters", h.name)
	}
	h.checkInputEventType(0, true)
	h.checkInterfaceParameter(1)
	h.checkCommonHandlerReturns(true)
}

func (h *handler) checkErrorStateHandler() {
	if len(h.params) != 3 {
		exitWithError("%s must have three parameters", h.name)
	}
	h.checkInterfaceParameter(0)
	h.checkInterfaceParameter(1)
	h.checkErrorParameter(2)
	h.checkErrorHandlerReturns()
}

func (h *handler) checkMigrationHandler() {
	if len(h.params) != 2 {
		exitWithError("%s must have only two parameters", h.name)
	}
	h.checkInputEventType(0, false)
	h.checkPayloadParameter(1)
	h.checkCommonHandlerReturns(false)
}

func (h *handler) checkTransitionHandler() {
	if len(h.params) < 2 {
		exitWithError("%s must have two or more parameters", h.name)
	}
	h.checkInputEventType(0, false)
	h.checkPayloadParameter(1)
	// TODO: check adapter helpers
	h.checkCommonHandlerReturns(false)
}

func (h *handler) checkAdapterResponseHandler() {
	if len(h.params) != 3 {
		exitWithError("%s must have three parameters", h.name)
	}
	h.checkInputEventType(0, false)
	h.checkPayloadParameter(1)
	// TODO: check respPayload here
	h.checkCommonHandlerReturns(false)
}

func (h *handler) checkAdapterResponseErrorHandler() {

}

// Helper methods for checkers above

func (h *handler) checkInputEventType(idx int, setEventType bool) {
	if setEventType && h.machine.InputEventType == nil {
		h.machine.InputEventType = &h.params[idx]
	} else if h.machine.InputEventType == nil || h.params[idx] != *h.machine.InputEventType {
		exitWithError("%s should have input event same type as Init payload", h.name)
	}
}

func (h *handler) checkPayloadParameter(idx int) {
	if !strings.HasPrefix(h.params[1], "*") {
		exitWithError("%s payload must be a pointer", h.name)
	}
	if h.machine.PayloadType == nil || h.params[idx] != *h.machine.PayloadType {
		exitWithError("%s returned payload should be same type as Init payload", h.name)
	}
}

func (h *handler) checkInterfaceParameter(idx int) {
	if h.params[idx] != "interface{}" {
		exitWithError("%d parameter for %s should be an interface{}", idx, h.name)
	}
}

func (h *handler) checkAdapterResponseParameter(idx int) {
	if h.params[idx] != "adapter.IAdapterResponse" {
		exitWithError("%d parameter for %s should be an AdapterResponse", idx, h.name)
	}
}

func (h *handler) checkErrorParameter(idx int) {
	if h.params[idx] != "error" {
		exitWithError("%d parameter for %s should be an error", idx, h.name)
	}
}

func (h *handler) checkCommonHandlerReturns(setPayload bool) {
	if len(h.results) != 3 {
		exitWithError("%s should return three values", h.name)
	}
	if setPayload && h.machine.PayloadType == nil {
		if !strings.HasPrefix(h.results[0], "*") {
			exitWithError("%s payload must be a pointer", h.name)
		}
		h.machine.PayloadType = &h.results[0]
	} else if h.machine.PayloadType == nil || h.results[0] != *h.machine.PayloadType {
		exitWithError("%s returned payload should be same type as Init payload", h.name)
	}
	if h.results[1] != "common.ElUpdate" {
		exitWithError("%s returned state should be ElUpdate", h.name)
	}
	if h.results[2] != "error" {
		exitWithError("%s returned error must be of type error", h.name)
	}
}

func (h *handler) checkErrorHandlerReturns() {
	if len(h.results) != 2 {
		exitWithError("%s should return two values", h.name)
	}
	if h.results[0] != *h.machine.PayloadType {
		exitWithError("%s returned payload should be same type as Init payload", h.name)
	}
	if h.results[1] != "common.ElUpdate" {
		exitWithError("%s returned state should be ElUpdate", h.name)
	}
}
