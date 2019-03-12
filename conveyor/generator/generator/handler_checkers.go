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
	if len(h.Params) != 2 {
		exitWithError("%s must have only two parameters", h.Name)
	}
	h.checkInputEventType(0, true)
	h.checkInterfaceParameter(1)
	h.checkCommonHandlerReturns(true)
}

func (h *handler) checkErrorStateHandler() {
	if len(h.Params) != 3 {
		exitWithError("%s must have three parameters", h.Name)
	}
	h.checkInterfaceParameter(0)
	h.checkInterfaceParameter(1)
	h.checkErrorParameter(2)
	h.checkErrorHandlerReturns()
}

func (h *handler) checkMigrationHandler() {
	if len(h.Params) != 2 {
		exitWithError("%s must have only two parameters", h.Name)
	}
	h.checkInputEventType(0, false)
	h.checkPayloadParameter(1)
	h.checkCommonHandlerReturns(false)
}

func (h *handler) checkTransitionHandler() {
	if len(h.Params) < 2 {
		exitWithError("%s must have two or more parameters", h.Name)
	}
	h.checkInputEventType(0, false)
	h.checkPayloadParameter(1)
	// TODO: check adapter helpers
	h.checkCommonHandlerReturns(false)
}

func (h *handler) checkAdapterResponseHandler() {
	if len(h.Params) != 3 {
		exitWithError("%s must have three parameters", h.Name)
	}
	h.checkInputEventType(0, false)
	h.checkPayloadParameter(1)
	// TODO: check respPayload here
	h.checkCommonHandlerReturns(false)
}

func (h *handler) checkAdapterResponseErrorHandler() {
	if len(h.Params) != 4 {
		exitWithError("%s must have four parameters", h.Name)
	}
	h.checkInterfaceParameter(0)
	h.checkInterfaceParameter(1)
	h.checkAdapterResponseParameter(2)
	h.checkErrorParameter(3)
	h.checkErrorHandlerReturns()
}

// Helper methods for checkers above

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
	if h.Results[1] != "common.ElementState" {
		exitWithError("%s returned state should be ElementState", h.Name)
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
	if h.Results[1] != "common.ElementState" {
		exitWithError("%s returned state should be ElementState", h.Name)
	}
}
