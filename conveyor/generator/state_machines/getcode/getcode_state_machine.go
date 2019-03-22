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

package getcode

import (
	"github.com/insolar/insolar/conveyor/interfaces/fsm"
	"github.com/insolar/insolar/conveyor/interfaces/iadapter"
)

// custom types
type Event struct{}
type Payload struct{}
type TA1 string
type TAR string

// conveyor: state_machine
type GetCodeStateMachine interface {
	GetTypeID() fsm.ID

	initPresentHandler(input Event, payload interface{}) (*Payload, fsm.ElementState, error)
	initFutureHandler(input Event, payload interface{}) (*Payload, fsm.ElementState, error)
	initPastHandler(input Event, payload interface{}) (*Payload, fsm.ElementState, error)

	errorPresentInit(input interface{}, payload interface{}, err error) (*Payload, fsm.ElementState)
	errorFutureInit(input interface{}, payload interface{}, err error) (*Payload, fsm.ElementState)
	errorPastInit(input interface{}, payload interface{}, err error) (*Payload, fsm.ElementState)

	// State Declaration
	stateFirst() fsm.StateID

	// Migration
	migrateFromPresentFirst(input Event, payload *Payload) (*Payload, fsm.ElementState, error)
	migrateFromFutureFirst(input Event, payload *Payload) (*Payload, fsm.ElementState, error)

	// Transition
	transitPresentFirst(input Event, payload *Payload /* todo: , adapterHelper TA1*/) (*Payload, fsm.ElementState, error)
	transitFutureFirst(input Event, payload *Payload /* todo: , adapterHelper TA1*/) (*Payload, fsm.ElementState, error)
	transitPastFirst(input Event, payload *Payload) (*Payload, fsm.ElementState, error)

	// TODO: Finalization
	// finalizePresentFirst(input Event, payload *Payload)
	// ...

	// Adapter Response
	responsePresentFirst(input Event, payload *Payload, respPayload TAR) (*Payload, fsm.ElementState, error)
	responseFutureFirst(input Event, payload *Payload, respPayload TAR) (*Payload, fsm.ElementState, error)
	responsePastFirst(input Event, payload *Payload, respPayload TAR) (*Payload, fsm.ElementState, error)

	// State Error
	errorPresentFirst(input interface{}, payload interface{}, err error) (*Payload, fsm.ElementState)
	errorFutureFirst(input interface{}, payload interface{}, err error) (*Payload, fsm.ElementState)
	errorPastFirst(input interface{}, payload interface{}, err error) (*Payload, fsm.ElementState)

	// Adapter Response Error
	errorResponsePresentFirst(input interface{}, payload interface{}, ar iadapter.Response, err error) (*Payload, fsm.ElementState)
	errorResponseFutureFirst(input interface{}, payload interface{}, ar iadapter.Response, err error) (*Payload, fsm.ElementState)
	errorResponsePastFirst(input interface{}, payload interface{}, ar iadapter.Response, err error) (*Payload, fsm.ElementState)

	// State Declaration
	stateSecond() fsm.StateID

	// Migration
	migrateFromPresentSecond(input Event, payload *Payload) (*Payload, fsm.ElementState, error)
	migrateFromFutureSecond(input Event, payload *Payload) (*Payload, fsm.ElementState, error)

	// Transition
	transitPresentSecond(input Event, payload *Payload /* todo: , adapterHelper TA1*/) (*Payload, fsm.ElementState, error)
	transitFutureSecond(input Event, payload *Payload /* todo: , adapterHelper TA1*/) (*Payload, fsm.ElementState, error)
	transitPastSecond(input Event, payload *Payload) (*Payload, fsm.ElementState, error)

	// TODO: Finalization
	// finalizePresentSecond(input Event, payload *Payload)
	// ...

	// Adapter Response
	responsePresentSecond(input Event, payload *Payload, respPayload TAR) (*Payload, fsm.ElementState, error)
	responseFutureSecond(input Event, payload *Payload, respPayload TAR) (*Payload, fsm.ElementState, error)
	responsePastSecond(input Event, payload *Payload, respPayload TAR) (*Payload, fsm.ElementState, error)

	// State Error
	errorPresentSecond(input interface{}, payload interface{}, err error) (*Payload, fsm.ElementState)
	errorFutureSecond(input interface{}, payload interface{}, err error) (*Payload, fsm.ElementState)
	errorPastSecond(input interface{}, payload interface{}, err error) (*Payload, fsm.ElementState)

	// Adapter Response Error
	errorResponsePresentSecond(input interface{}, payload interface{}, ar iadapter.Response, err error) (*Payload, fsm.ElementState)
	errorResponseFutureSecond(input interface{}, payload interface{}, ar iadapter.Response, err error) (*Payload, fsm.ElementState)
	errorResponsePastSecond(input interface{}, payload interface{}, ar iadapter.Response, err error) (*Payload, fsm.ElementState)
}

// type CleanTestStateMachine struct {
// 	BaseTestStateMachine
// }
//
// func (sm *CleanTestStateMachine) initPresentHandler(input Event, payload interface{}) (*Payload, fsm.ElementState, error) {
// 	return nil, fsm.NewElementState(sm.GetTypeID(), sm.stateFirst()), nil
// }
//
// func (sm *CleanTestStateMachine) initFutureHandler(input Event, payload interface{}) (*Payload, fsm.ElementState, error) {
// 	panic("implement me")
// }
func transitPresentFirst(input Event, payload *Payload /* todo: , adapterHelper TA1*/) (*Payload, fsm.ElementState, error) {
	return
}
