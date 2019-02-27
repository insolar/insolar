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

package conveyor

import (
	"github.com/insolar/insolar/conveyor/generator/common"
)

// ActivationStatus represents status of work for slot element
type ActivationStatus int

//go:generate stringer -type=ActivationStatus
const (
	EmptyElement = ActivationStatus(iota)
	ActiveElement
	NotActiveElement
)

type stateMachineTypeI interface {
	getTypeID() int
	getMigrationHandler(state int) common.MigrationHandler
	getTransitionHandler(state int) common.TransitHandler
	getResponseHandler(state int) common.AdapterResponseHandler
	getNestedHandler(state int) common.NestedHandler

	getTransitionErrorHandler(state int) common.TransitionErrorHandler
	getResponseErrorHandler(state int) common.ResponseErrorHandler
}

type reactivateMode interface {
}

type slotElementHelper interface {
	slotElementRestrictedHelper
	informParent(payload interface{}) bool
	deactivateTill(reactivateOn reactivateMode)
	sendTask(adapterID uint32, taskPayload interface{}, respHandlerID uint32) error
	// joinSequence( sequenceKey map-key,sequenceOrder uint64 )
	// isSequenceHead() bool
}

type slotElementRestrictedHelper interface {
	slotElementReadOnly

	getParentElementID() uint32
	getInputEvent() interface{}
	getPayload() interface{}

	reactivate()
	leaveSequence()
}

type slotElementReadOnly interface {
	getElementID() uint32
	getNodeID() uint32
	getType() int
	getState() int
}

type slotElement struct {
	id               uint32
	nodeID           uint32
	parentElementID  uint32
	inputEvent       interface{}
	payload          interface{} // nolint
	stateMachineType stateMachineTypeI
	state            uint16

	nextElement      *slotElement
	activationStatus ActivationStatus
}

// newSlotElement creates new slot element with provided activation status
func newSlotElement(activationStatus ActivationStatus) *slotElement {
	return &slotElement{activationStatus: activationStatus}
}
