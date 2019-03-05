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
	"github.com/insolar/insolar/conveyor/interfaces/statemachine"
	"github.com/insolar/insolar/conveyor/interfaces/slot"
	"github.com/insolar/insolar/conveyor/interfaces/constant"
)

type State struct {
	Migration statemachine.MigrationHandler
	MigrationFuturePresent statemachine.MigrationHandler
	Transition statemachine.TransitHandler
	TransitionFuture statemachine.TransitHandler
	TransitionPast statemachine.TransitHandler
	AdapterResponse statemachine.AdapterResponseHandler
	AdapterResponseFuture statemachine.AdapterResponseHandler
	AdapterResponsePast statemachine.AdapterResponseHandler
	ErrorState statemachine.TransitionErrorHandler
	ErrorStateFuture statemachine.TransitionErrorHandler
	ErrorStatePast statemachine.TransitionErrorHandler
	AdapterResponseError statemachine.ResponseErrorHandler
	AdapterResponseErrorFuture statemachine.ResponseErrorHandler
	AdapterResponseErrorPast statemachine.ResponseErrorHandler
	/*Finalization *handler
	FinalizationFuture *handler
	FinalizationPast *handler*/
}

type StateMachine struct {
	Id     int
	States []State
}

func (sm *StateMachine) GetTypeID() int {
	return sm.Id
}

func (sm *StateMachine) GetMigrationHandler(slotType constant.PulseState, state uint32) statemachine.MigrationHandler {
	switch slotType {
	case constant.Future:
		return sm.States[state].MigrationFuturePresent
	case constant.Present:
		return sm.States[state].Migration
	default:
		panic("migration handler can't be called for past tense")
	}
}

func (sm *StateMachine) GetTransitionHandler(slotType constant.PulseState, state uint32) statemachine.TransitHandler {
	switch slotType {
	case constant.Future:
		return sm.States[state].TransitionFuture
	case constant.Present:
		return sm.States[state].Transition
	case constant.Past:
		return sm.States[state].TransitionPast
	case constant.Antique:
		return sm.States[state].TransitionPast
	default:
		panic("handler can't be called for unallocated tense")
	}
}

func (sm *StateMachine) GetResponseHandler(slotType constant.PulseState, state uint32) statemachine.AdapterResponseHandler {
	switch slotType {
	case constant.Future:
		return sm.States[state].AdapterResponseFuture
	case constant.Present:
		return sm.States[state].AdapterResponse
	case constant.Past:
		return sm.States[state].AdapterResponsePast
	case constant.Antique:
		return sm.States[state].AdapterResponsePast
	default:
		panic("handler can't be called for unallocated tense")
	}
}

func (sm *StateMachine) GetNestedHandler(slotType constant.PulseState, state uint32) statemachine.NestedHandler {
	return func(element slot.SlotElementHelper, err error) (interface{}, uint32) {
		// todo needs implementation
		return nil, 0
	}
}

func (sm *StateMachine) GetTransitionErrorHandler(slotType constant.PulseState, state uint32) statemachine.TransitionErrorHandler {
	switch slotType {
	case constant.Future:
		return sm.States[state].ErrorStateFuture
	case constant.Present:
		return sm.States[state].ErrorState
	case constant.Past:
		return sm.States[state].ErrorStatePast
	case constant.Antique:
		return sm.States[state].ErrorStatePast
	default:
		panic("handler can't be called for unallocated tense")
	}
}

func (sm *StateMachine) GetResponseErrorHandler(slotType constant.PulseState, state uint32) statemachine.ResponseErrorHandler {
	switch slotType {
	case constant.Future:
		return sm.States[state].AdapterResponseErrorFuture
	case constant.Present:
		return sm.States[state].AdapterResponseError
	case constant.Past:
		return sm.States[state].AdapterResponseErrorPast
	case constant.Antique:
		return sm.States[state].AdapterResponseErrorPast
	default:
		panic("handler can't be called for unallocated tense")
	}
}

type ElState uint32 //Element State Machine Type ID
type ElType uint32  //Element State ID
type ElNewState uint32

func (s ElState) ToInt() uint32 {
	return uint32(s)
}

type RawHandlerT func(element slot.SlotElementHelper) (err error, new_state uint32, new_payload interface{})

type ElUpdate uint32 ///Element State ID + Element Machine Type ID << 10

func (s ElUpdate) ToInt() uint32 {
	return uint32(s)
}
