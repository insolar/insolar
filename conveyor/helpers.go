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

package conveyor

import (
	"fmt"
)

const (
	smShift       = 10
	maxStateValue = (1 << smShift) - 1
)

// first part of elUpdate is State Machine, second part is State.
func extractStates(elUpdate uint32) (uint32, uint32) {
	sm := elUpdate >> smShift
	state := elUpdate & ((1 << smShift) - 1)

	return sm, state
}

// 'state' MUST be less than maxStateValue ( 2^smShift )
func joinStates(sm uint32, state uint32) uint32 {
	if state > maxStateValue {
		panic(fmt.Sprint("Invalid state: ", state))
	}
	result := sm
	result = result << smShift
	return result + state
}
