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
	Transit statemachine.TransitHandler
	Migrate statemachine.MigrationHandler
	Error   statemachine.TransitionErrorHandler
}

type StateMachine struct {
	Id int
	InitHandler     statemachine.InitHandler
	States          []State
	FinalizeHandler interface{}
}

func (sm *StateMachine) GetTypeID() int {
	return sm.Id
}

func (sm *StateMachine) GetTransitionHandler(state int) statemachine.TransitHandler {
	return sm.States[state].Transit
}

func (sm *StateMachine) GetMigrationHandler(state int) statemachine.MigrationHandler {
	return sm.States[state].Migrate
}

func (sm *StateMachine) GetTransitionErrorHandler(state int) statemachine.TransitionErrorHandler {
	return sm.States[state].Error
}

func (sm *StateMachine) GetResponseHandler(state int) statemachine.AdapterResponseHandler {
	return func(element slot.SlotElementHelper, err error) (interface{}, uint32) {
		return nil, 0
	}
}

func (sm *StateMachine) GetNestedHandler() statemachine.NestedHandler {
	return func(element slot.SlotElementHelper, err error) (interface{}, uint32) {
		return nil, 0
	}
}

func (sm *StateMachine) GetResponseErrorHandler(state int) statemachine.ResponseErrorHandler {
	return func(element slot.SlotElementHelper, err error) (interface{}, uint32) {
		return nil, 0
	}
}

type ElState uint32 //Element State Machine Type ID
type ElType uint32  //Element State ID

func (s ElState) ToInt() uint32 {
	return uint32(s)
}

type RawHandlerT func(element slot.SlotElementHelper) (err error, new_state uint32, new_payload interface{})

type ElUpdate uint32 ///Element State ID + Element Machine Type ID << 10
