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

package common

import (
	"github.com/insolar/insolar/conveyor/interfaces/slot"
	"github.com/insolar/insolar/conveyor/interfaces/statemachine"
)

// TODO: move this code from common package
// State struct contains predefined set of handlers
type State struct {
	Migration            statemachine.MigrationHandler
	Transition           statemachine.TransitHandler
	AdapterResponse      statemachine.AdapterResponseHandler
	ErrorState           statemachine.TransitionErrorHandler
	AdapterResponseError statemachine.ResponseErrorHandler
	// TODO: Finalization handlers
}

// StateMachine is a type for conveyor state machines
type StateMachine struct {
	ID     int
	States []State
}

// GetTypeID method returns StateMachine ID
func (sm *StateMachine) GetTypeID() int {
	return sm.ID
}

// GetMigrationHandler method returns migration handler
func (sm *StateMachine) GetMigrationHandler(state uint32) statemachine.MigrationHandler {
	return sm.States[state].Migration
}

// GetTransitionHandler method returns transition handler
func (sm *StateMachine) GetTransitionHandler(state uint32) statemachine.TransitHandler {
	return sm.States[state].Transition
}

// GetResponseHandler returns response handler
func (sm *StateMachine) GetResponseHandler(state uint32) statemachine.AdapterResponseHandler {
	return sm.States[state].AdapterResponse
}

// GetNestedHandler returns nested handler
func (sm *StateMachine) GetNestedHandler(state uint32) statemachine.NestedHandler {
	return func(element slot.SlotElementHelper, err error) (interface{}, uint32) {
		// TODO: Implement me
		return nil, 0
	}
}

// GetTransitionErrorHandler returns transition error handler
func (sm *StateMachine) GetTransitionErrorHandler(state uint32) statemachine.TransitionErrorHandler {
	return sm.States[state].ErrorState
}

// GetResponseErrorHandler returns response error handler
func (sm *StateMachine) GetResponseErrorHandler(state uint32) statemachine.ResponseErrorHandler {
	return sm.States[state].AdapterResponseError
}
