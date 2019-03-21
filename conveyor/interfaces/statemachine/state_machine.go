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
	"github.com/insolar/insolar/conveyor/interfaces/fsm"
	"github.com/insolar/insolar/conveyor/interfaces/iadapter"
	"github.com/insolar/insolar/conveyor/interfaces/slot"
)

// Types below describes different types of raw handlers
type TransitHandler func(element slot.SlotElementHelper) (interface{}, fsm.ElementState, error)
type MigrationHandler func(element slot.SlotElementHelper) (interface{}, fsm.ElementState, error)
type AdapterResponseHandler func(element slot.SlotElementHelper, response iadapter.Response) (interface{}, fsm.ElementState, error)
type NestedHandler func(element slot.SlotElementHelper, err error) (interface{}, fsm.ElementState)
type TransitionErrorHandler func(element slot.SlotElementHelper, err error) (interface{}, fsm.ElementState)
type ResponseErrorHandler func(element slot.SlotElementHelper, response iadapter.Response, err error) (interface{}, fsm.ElementState)

// StateMachine describes access to element's state machine
//go:generate minimock -i github.com/insolar/insolar/conveyor/interfaces/statemachine.StateMachine -o ./ -s _mock.go
type StateMachine interface {
	GetTypeID() fsm.ID
	GetMigrationHandler(state fsm.StateID) MigrationHandler
	GetTransitionHandler(state fsm.StateID) TransitHandler
	GetResponseHandler(state fsm.StateID) AdapterResponseHandler
	GetNestedHandler(state fsm.StateID) NestedHandler
	GetTransitionErrorHandler(state fsm.StateID) TransitionErrorHandler
	GetResponseErrorHandler(state fsm.StateID) ResponseErrorHandler
}

// SetAccessor gives access to set of state machines
type SetAccessor interface {
	GetStateMachineByID(id int) StateMachine
}
