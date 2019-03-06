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

type State struct {
	Migration              statemachine.MigrationHandler
	// MigrationFuturePresent statemachine.MigrationHandler
	Transition statemachine.TransitHandler
	// TransitionFuture statemachine.TransitHandler
	// TransitionPast statemachine.TransitHandler
	AdapterResponse statemachine.AdapterResponseHandler
	// AdapterResponseFuture statemachine.AdapterResponseHandler
	// AdapterResponsePast statemachine.AdapterResponseHandler
	ErrorState statemachine.TransitionErrorHandler
	// ErrorStateFuture statemachine.TransitionErrorHandler
	// ErrorStatePast statemachine.TransitionErrorHandler
	AdapterResponseError statemachine.ResponseErrorHandler
	// AdapterResponseErrorFuture statemachine.ResponseErrorHandler
	// AdapterResponseErrorPast statemachine.ResponseErrorHandler
	/*Finalization *handler
	FinalizationFuture *handler
	FinalizationPast *handler*/
}

type StateMachine struct {
	ID     int
	States []State
}

func (sm *StateMachine) GetTypeID() int {
	return sm.ID
}

func (sm *StateMachine) GetMigrationHandler(state uint32) statemachine.MigrationHandler {
	return sm.States[state].Migration
}

func (sm *StateMachine) GetTransitionHandler(state uint32) statemachine.TransitHandler {
	return sm.States[state].Transition
}

func (sm *StateMachine) GetResponseHandler(state uint32) statemachine.AdapterResponseHandler {
	return sm.States[state].AdapterResponse
}

func (sm *StateMachine) GetNestedHandler(state uint32) statemachine.NestedHandler {
	return func(element slot.SlotElementHelper, err error) (interface{}, uint32) {
		// todo needs implementation
		return nil, 0
	}
}

func (sm *StateMachine) GetTransitionErrorHandler(state uint32) statemachine.TransitionErrorHandler {
	return sm.States[state].ErrorState
}

func (sm *StateMachine) GetResponseErrorHandler(state uint32) statemachine.ResponseErrorHandler {
	return sm.States[state].AdapterResponseError
}

type ElState uint32 //Element State Machine Type ID
type ElType uint32  //Element State ID
func (s ElState) ToInt() uint32 {
	return uint32(s)
}

type ElUpdate uint32 ///Element State ID + Element Machine Type ID << 10
func (s ElUpdate) ToInt() uint32 {
	return uint32(s)
}
