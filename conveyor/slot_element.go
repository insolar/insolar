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

import "github.com/insolar/insolar/conveyor/interfaces/slot"

// ActivationStatus represents status of work for slot element
type ActivationStatus int

//go:generate stringer -type=ActivationStatus
const (
	EmptyElement = ActivationStatus(iota)
	ActiveElement
	NotActiveElement
)

type StateMachineType interface {
	GetTypeID() int
	GetMigrationHandler(state int) MigrationHandler
	GetTransitionHandler(state int) TransitHandler
	GetResponseHandler(state int) AdapterResponseHandler
	GetNestedHandler(state int) NestedHandler

	GetTransitionErrorHandler(state int) TransitionErrorHandler
	GetResponseErrorHandler(state int) ResponseErrorHandler
}

type slotElement struct {
	id               uint32
	nodeID           uint32
	parentElementID  uint32
	inputEvent       interface{}
	payload          interface{} // nolint
	postponedError   error
	stateMachineType StateMachineType
	state            uint16

	nextElement      *slotElement
	activationStatus ActivationStatus
}

// newSlotElement creates new slot element with provided activation status
func newSlotElement(activationStatus ActivationStatus) *slotElement {
	return &slotElement{activationStatus: activationStatus}
}

type MigrationHandler func()
type TransitHandler func()
type AdapterResponseHandler func()
type NestedHandler func()
type TransitionErrorHandler func()
type ResponseErrorHandler func()

// GetMigrationHandler implements StateMachineType
func (se *slotElement) GetMigrationHandler(state int) MigrationHandler {
	panic("implement me")
	//return matrix.Matrix.GetHandlers(se.stateMachineType.GetTypeID(), state).Migrate
}

// GetTransitionHandler implements StateMachineType
func (se *slotElement) GetTransitionHandler(state int) TransitHandler {
	panic("implement me")
	//return matrix.Matrix.GetHandlers(se.stateMachineType.GetTypeID(), state).Transit
}

// GetResponseHandler implements StateMachineType
func (se *slotElement) GetResponseHandler(state int) AdapterResponseHandler {
	panic("implement me")
}

// GetNestedHandler implements StateMachineType
func (se *slotElement) GetNestedHandler(state int) NestedHandler {
	panic("implement me")
}

// GetTransitionErrorHandler implements StateMachineType
func (se *slotElement) GetTransitionErrorHandler(state int) TransitionErrorHandler {
	panic("implement me")
}

// GetResponseErrorHandler implements StateMachineType
func (se *slotElement) GetResponseErrorHandler(state int) ResponseErrorHandler {
	panic("implement me")
}

// ---- SlotElementRestrictedHelper

// GetParentElementID implements SlotElementRestrictedHelper
func (se *slotElement) GetParentElementID() uint32 {
	return se.parentElementID
}

// GetInputEvent implements SlotElementRestrictedHelper
func (se *slotElement) GetInputEvent() interface{} {
	return se.inputEvent
}

// GetPayload implements SlotElementRestrictedHelper
func (se *slotElement) GetPayload() interface{} {
	return se.payload
}

// Reactivate implements SlotElementRestrictedHelper
func (se *slotElement) Reactivate() {
	panic("implement me")
}

// LeaveSequence implements SlotElementRestrictedHelper
func (se *slotElement) LeaveSequence() {
	panic("implement me")
}

// ---- SlotElementReadOnly

// LeaveSequence implements SlotElementReadOnly
func (se *slotElement) GetElementID() uint32 {
	return se.id
}

// GetNodeID implements SlotElementReadOnly
func (se *slotElement) GetNodeID() uint32 {
	return se.nodeID
}

// GetType implements SlotElementReadOnly
func (se *slotElement) GetType() int {
	return se.stateMachineType.GetTypeID()
}

// GetState implements SlotElementReadOnly
func (se *slotElement) GetState() uint16 {
	return se.state
}

// ---- SlotElementHelper

// InformParent implements SlotElementHelper
func (se *slotElement) InformParent(payload interface{}) bool {
	panic("implement me")
}

// DeactivateTill implements SlotElementHelper
func (se *slotElement) DeactivateTill(reactivateOn slot.ReactivateMode) {
	panic("implement me")
}

// SendTask implements SlotElementHelper
func (se *slotElement) SendTask(adapterID uint32, taskPayload interface{}, respHandlerID uint32) error {
	panic("implement me")
}
