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
	"github.com/insolar/insolar/conveyor/generator/matrix"
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
	GetTypeID() int
	GetMigrationHandler(state int) common.MigrationHandler
	GetTransitionHandler(state int) common.TransitHandler
	GetResponseHandler(state int) common.AdapterResponseHandler
	GetNestedHandler(state int) common.NestedHandler

	GetTransitionErrorHandler(state int) common.TransitionErrorHandler
	GetResponseErrorHandler(state int) common.ResponseErrorHandler
}

type reactivateMode interface {
}

type SlotElementHelperI interface {
	SlotElementRestrictedHelper
	informParent(payload interface{}) bool
	deactivateTill(reactivateOn reactivateMode)
	sendTask(adapterID uint32, taskPayload interface{}, respHandlerID uint32) error
	// joinSequence( sequenceKey map-key,sequenceOrder uint64 )
	// isSequenceHead() bool
}

type SlotElementRestrictedHelper interface {
	SlotElementReadOnly

	GetParentElementID() uint32
	GetInputEvent() interface{}
	GetPayload() interface{}

	Reactivate()
	LeaveSequence()
}

type SlotElementReadOnly interface {
	GetElementID() uint32
	GetNodeID() uint32
	GetType() int
	GetState() int
}

type slotElement struct {
	id               uint32
	nodeID           uint32
	parentElementID  uint32
	inputEvent       interface{}
	payload          interface{} // nolint
	postponedError   error
	stateMachineType stateMachineTypeI
	state            uint16

	nextElement      *slotElement
	activationStatus ActivationStatus
}

// newSlotElement creates new slot element with provided activation status
func newSlotElement(activationStatus ActivationStatus) *slotElement {
	return &slotElement{activationStatus: activationStatus}
}

func (se *slotElement) GetTypeID() int {
	return se.stateMachineType.GetTypeID()
}

func (se *slotElement) GetMigrationHandler(state int) common.MigrationHandler {
	return matrix.Matrix.GetHandlers(se.stateMachineType.GetTypeID(), state).Migrate
}

func (se *slotElement) GetTransitionHandler(state int) common.TransitHandler {
	return matrix.Matrix.GetHandlers(se.stateMachineType.GetTypeID(), state).Transit
}

// TODO: implement me
func (se *slotElement) GetResponseHandler(state int) common.AdapterResponseHandler {
	panic("implement me")
}

// TODO: implement me
func (se *slotElement) GetNestedHandler(state int) common.NestedHandler {
	panic("implement me")
}

// TODO: implement me
func (se *slotElement) GetTransitionErrorHandler(state int) common.TransitionErrorHandler {
	panic("implement me")
}

// TODO: implement me
func (se *slotElement) GetResponseErrorHandler(state int) common.ResponseErrorHandler {
	panic("implement me")
}

// GetParentElementID return parentElementID
func (se *slotElement) GetParentElementID() uint32 {
	return se.parentElementID
}

// GetInputEvent return inputEvent
func (se *slotElement) GetInputEvent() interface{} {
	return se.inputEvent
}

// GetPayload return payload
func (se *slotElement) GetPayload() interface{} {
	return se.payload
}

func (se *slotElement) Reactivate() {
	panic("implement me")
}

func (se *slotElement) LeaveSequence() {
	panic("implement me")
}
