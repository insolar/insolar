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
)

type InitHandler func(element slot.SlotElementHelper) (interface{}, uint32, error)
type TransitHandler func(element slot.SlotElementHelper) (interface{}, uint32, error)
type MigrationHandler func(element slot.SlotElementHelper) (interface{}, uint32, error)
type ErrorHandler func(element slot.SlotElementHelper, err error) (interface{}, uint32)
type AdapterResponseHandler func(element slot.SlotElementHelper, err error) (interface{}, uint32)
type NestedHandler func(element slot.SlotElementHelper, err error) (interface{}, uint32)

type TransitionErrorHandler func(element slot.SlotElementHelper, err error) (interface{}, uint32)
type ResponseErrorHandler func(element slot.SlotElementHelper, err error) (interface{}, uint32)

type State struct {
	Transit TransitHandler
	Migrate MigrationHandler
	Error   ErrorHandler
}

type StateMachine struct {
	Id int
	InitHandler     InitHandler
	States          []State
	FinalizeHandler interface{}
}

func (sm *StateMachine) GetTypeID() int {
	return sm.Id
}

func (sm *StateMachine) GetTransitionHandler(state uint16) TransitHandler {
	return sm.States[state].Transit
}

func (sm *StateMachine) GetMigrationHandler(state uint16) MigrationHandler {
	return sm.States[state].Migrate
}

func (sm *StateMachine) GetTransitionErrorHandler(state uint16) ErrorHandler {
	return sm.States[state].Error
}

func (sm *StateMachine) GetResponseHandler(state uint16) AdapterResponseHandler {
	return func(element slot.SlotElementHelper, err error) (interface{}, uint32) {
		return nil, 0
	}
}

func (sm *StateMachine) GetNestedHandler() NestedHandler {
	return func(element slot.SlotElementHelper, err error) (interface{}, uint32) {
		return nil, 0
	}
}

func (sm *StateMachine) GetResponseErrorHandler(state uint16) ResponseErrorHandler {
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
