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

package fsm

import (
	"fmt"
)

// Element State ID
type StateID uint32

// Element State Machine Type ID
type ID uint32

// ElementState is StateID + (ID << 10)
type ElementState uint32

const (
	stateShift    = 10
	maxStateValue = (1 << stateShift) - 1
)

// NewElementState constructs element from ID and StateID
// state MUST be less than maxStateValue ( 2^stateShift )
func NewElementState(stateMachine ID, state StateID) ElementState {
	if state > maxStateValue {
		panic(fmt.Sprint("Invalid state: ", state))
	}
	result := (uint32(stateMachine) << stateShift) + uint32(state)
	return ElementState(result)
}

// Parse method returns ID and StateID from ElementState
func (es ElementState) Parse() (ID, StateID) {
	sm := es >> stateShift
	state := es & ((1 << stateShift) - 1)
	return ID(sm), StateID(state)
}
