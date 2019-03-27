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
	"github.com/insolar/insolar/conveyor/handler"
	"github.com/insolar/insolar/conveyor/interfaces/fsm"
	"github.com/insolar/insolar/conveyor/interfaces/slot"
)

// State struct contains predefined set of handlers
type State struct {
	Migration            handler.MigrationHandler
	Transition           handler.TransitHandler
	AdapterResponse      handler.AdapterResponseHandler
	ErrorState           handler.TransitionErrorHandler
	AdapterResponseError handler.ResponseErrorHandler
}

// StateMachine is a type for conveyor state machines
type StateMachine struct {
	ID     fsm.ID
	States []State
}

// GetTypeID method returns StateMachine ID
func (sm *StateMachine) GetTypeID() fsm.ID {
	return sm.ID
}

// GetMigrationHandler method returns migration handler
func (sm *StateMachine) GetMigrationHandler(state fsm.StateID) handler.MigrationHandler {
	return sm.States[state].Migration
}

// GetTransitionHandler method returns transition handler
func (sm *StateMachine) GetTransitionHandler(state fsm.StateID) handler.TransitHandler {
	return sm.States[state].Transition
}

// GetResponseHandler returns response handler
func (sm *StateMachine) GetResponseHandler(state fsm.StateID) handler.AdapterResponseHandler {
	return sm.States[state].AdapterResponse
}

// GetNestedHandler returns nested handler
func (sm *StateMachine) GetNestedHandler(state fsm.StateID) handler.NestedHandler {
	return func(element slot.SlotElementHelper, err error) (interface{}, fsm.ElementState) {
		// TODO: Implement me
		return nil, 0
	}
}

// GetTransitionErrorHandler returns transition error handler
func (sm *StateMachine) GetTransitionErrorHandler(state fsm.StateID) handler.TransitionErrorHandler {
	return sm.States[state].ErrorState
}

// GetResponseErrorHandler returns response error handler
func (sm *StateMachine) GetResponseErrorHandler(state fsm.StateID) handler.ResponseErrorHandler {
	return sm.States[state].AdapterResponseError
}
