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
	"github.com/insolar/insolar/conveyor/interfaces/constant"
	"github.com/insolar/insolar/conveyor/interfaces/slot"
)

type TransitHandler func(element slot.SlotElementHelper) (interface{}, uint32, error)
type MigrationHandler func(element slot.SlotElementHelper) (interface{}, uint32, error)
type AdapterResponseHandler func(element slot.SlotElementHelper, response adapter.IAdapterResponse) (interface{}, uint32, error)
type NestedHandler func(element slot.SlotElementHelper, err error) (interface{}, uint32)

type TransitionErrorHandler func(element slot.SlotElementHelper, err error) (interface{}, uint32)
type ResponseErrorHandler func(element slot.SlotElementHelper, response adapter.IAdapterResponse, err error) (interface{}, uint32)

// StateMachineType describes access to element's state machine
//go:generate minimock -i github.com/insolar/insolar/conveyor/interfaces/statemachine.StateMachineType -o ./ -s _mock.go
type StateMachineType interface {
	GetTypeID() int
	GetMigrationHandler(pulseState constant.PulseState, state uint32) MigrationHandler
	GetTransitionHandler(pulseState constant.PulseState, state uint32) TransitHandler
	GetResponseHandler(pulseState constant.PulseState, state uint32) AdapterResponseHandler
	GetNestedHandler(pulseState constant.PulseState, state uint32) NestedHandler

	GetTransitionErrorHandler(pulseState constant.PulseState, state uint32) TransitionErrorHandler
	GetResponseErrorHandler(pulseState constant.PulseState, state uint32) ResponseErrorHandler
}
