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
	"github.com/insolar/insolar/conveyor/interfaces/islot"
	"github.com/insolar/insolar/conveyor/interfaces/istatemachine"
)

// ActivationStatus represents status of work for slot element
type ActivationStatus int

//go:generate stringer -type=ActivationStatus
const (
	EmptyElement = ActivationStatus(iota)
	ActiveElement
	NotActiveElement
)

type slotElement struct {
	id               uint32
	nodeID           uint32
	parentElementID  uint32
	inputEvent       interface{}
	payload          interface{} // nolint
	postponedError   error
	stateMachineType istatemachine.StateMachineType
	state            uint32

	nextElement      *slotElement
	prevElement      *slotElement
	activationStatus ActivationStatus
}

// newSlotElement creates new slot element with provided activation status
func newSlotElement(activationStatus ActivationStatus) *slotElement {
	return &slotElement{activationStatus: activationStatus}
}

// ---- SlotElementRestrictedHelper

func (se *slotElement) setDeleteState() {
	se.activationStatus = EmptyElement
}

func (se *slotElement) update(state uint32, payload interface{}, sm istatemachine.StateMachineType) {
	se.state = state
	se.payload = payload
	se.stateMachineType = sm
}

func (se *slotElement) isDeactivated() bool {
	return se.activationStatus == NotActiveElement
}

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
func (se *slotElement) GetState() uint32 {
	return se.state
}

// ---- SlotElementHelper

// InformParent implements SlotElementHelper
func (se *slotElement) InformParent(payload interface{}) bool {
	panic("implement me")
}

// DeactivateTill implements SlotElementHelper
func (se *slotElement) DeactivateTill(reactivateOn islot.ReactivateMode) {
	panic("implement me")
}

// SendTask implements SlotElementHelper
func (se *slotElement) SendTask(adapterID uint32, taskPayload interface{}, respHandlerID uint32) error {
	panic("implement me")
}
