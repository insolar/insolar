// +build with_generated

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

package main

import (
	"github.com/insolar/insolar/conveyor/generator/state_machines/sample"
	"github.com/insolar/insolar/conveyor/interfaces/slot"
	"testing"
)

func Test_Generated_State_Machine(t *testing.T) {
	element := slot.NewSlotElementHelperMock(t)
	element.GetInputEventFunc = func() interface{} {
		return sample.Event{}
	}
	element.GetPayloadFunc = func() interface{} {
		return &sample.Payload{}
	}

	machines := sample.SMRHTestStateMachineFactory()

	var state uint32 = 0
	for {
		_, state, _ = machines[1].GetTransitionHandler(state)(element)
		if state == 0 {
			break
		}
	}
}

