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

package matrix

import (
	"github.com/insolar/insolar/conveyor/fsm"
	"github.com/insolar/insolar/conveyor/handler"
)

// StateMachine describes access to element's state machine
//go:generate minimock -i github.com/insolar/insolar/conveyor/generator/matrix.StateMachine -o ./ -s _mock.go
type StateMachine interface {
	GetTypeID() fsm.ID
	GetMigrationHandler(state fsm.StateID) handler.MigrationHandler
	GetTransitionHandler(state fsm.StateID) handler.TransitHandler
	GetResponseHandler(state fsm.StateID) handler.AdapterResponseHandler
	GetNestedHandler(state fsm.StateID) handler.NestedHandler
	GetTransitionErrorHandler(state fsm.StateID) handler.TransitionErrorHandler
	GetResponseErrorHandler(state fsm.StateID) handler.ResponseErrorHandler
}

// SetAccessor gives access to set of state machines
type SetAccessor interface {
	GetStateMachineByID(id int) StateMachine
}

type StateMachineHolder interface {
	GetFutureConfig() SetAccessor
	GetPresentConfig() SetAccessor
	GetPastConfig() SetAccessor
	GetInitialStateMachine() StateMachine
}
