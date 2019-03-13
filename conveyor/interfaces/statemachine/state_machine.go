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

package statemachine

import (
	"github.com/insolar/insolar/conveyor/interfaces/adapter"
	"github.com/insolar/insolar/conveyor/interfaces/slot"
)

// TODO: move here helpers for types below
// Element State ID
type StateID uint32

// Element State Machine Type ID
type ID uint32

// ElementState is StateID + (ID << 10)
type ElementState uint32

const (
	stateShift = 10
)

// NewElementState constructs element from ID and StateID
func NewElementState(stateMachine ID, state StateID) ElementState {
	result := (uint32(stateMachine) << stateShift) + uint32(state)
	return ElementState(result)
}

// Parse method returns ID and StateID from ElementState
func (es ElementState) Parse() (ID, StateID) {
	sm := es >> stateShift
	state := es & ((1 << stateShift) - 1)
	return ID(sm), StateID(state)
}

// Types below describes different types of raw handlers
type TransitHandler func(element slot.SlotElementHelper) (interface{}, ElementState, error)
type MigrationHandler func(element slot.SlotElementHelper) (interface{}, ElementState, error)
type AdapterResponseHandler func(element slot.SlotElementHelper, response adapter.IAdapterResponse) (interface{}, ElementState, error)
type NestedHandler func(element slot.SlotElementHelper, err error) (interface{}, ElementState)
type TransitionErrorHandler func(element slot.SlotElementHelper, err error) (interface{}, ElementState)
type ResponseErrorHandler func(element slot.SlotElementHelper, response adapter.IAdapterResponse, err error) (interface{}, ElementState)

// StateMachine describes access to element's state machine
//go:generate minimock -i github.com/insolar/insolar/conveyor/interfaces/statemachine.StateMachine -o ./ -s _mock.go
type StateMachine interface {
	GetTypeID() ID
	GetMigrationHandler(state StateID) MigrationHandler
	GetTransitionHandler(state StateID) TransitHandler
	GetResponseHandler(state StateID) AdapterResponseHandler
	GetNestedHandler(state StateID) NestedHandler
	GetTransitionErrorHandler(state StateID) TransitionErrorHandler
	GetResponseErrorHandler(state StateID) ResponseErrorHandler
}
