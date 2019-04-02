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
	"testing"
	"time"

	"github.com/insolar/insolar/conveyor/adapter"
	"github.com/insolar/insolar/conveyor/adapter/adapterid"
	"github.com/insolar/insolar/conveyor/fsm"
	"github.com/insolar/insolar/conveyor/generator/matrix"
	"github.com/insolar/insolar/conveyor/generator/state_machines/sample/custom"
	"github.com/insolar/insolar/insolar"
)

func Test_Generated_State_Machine(t *testing.T) {
	cnt := 0
	element := fsm.NewSlotElementHelperMock(t)
	element.GetInputEventFunc = func() interface{} {
		cnt++
		if cnt == 3 {
			return insolar.ConveyorPendingMessage{}
		}
		return custom.Event{}
	}
	element.GetPayloadFunc = func() interface{} {
		return &custom.Payload{}
	}

	active := true

	element.DeactivateTillFunc = func(p fsm.ReactivateMode) {
		active = false
	}

	machines := matrix.NewMatrix().GetPresentConfig()
	var stateID fsm.StateID = 0
	var elementState fsm.ElementState = 0
	var err error

	element.SendTaskFunc = func(adapterID adapterid.ID, taskPayload interface{}, respHandlerID uint32) error {
		go func() {
			time.Sleep(time.Second)
			p := taskPayload.(adapter.SendResponseTask)
			_, elementState, err = machines.GetStateMachineByID(matrix.SampleStateMachine).GetResponseHandler(fsm.StateID(respHandlerID))(element, p.Result)
			if err != nil {
				panic(err)
			}
			_, stateID = elementState.Parse()
			active = true
		}()
		return nil
	}

	for {
		if active {
			handler := machines.GetStateMachineByID(matrix.SampleStateMachine).GetTransitionHandler(stateID)
			if handler != nil {
				_, elementState, err = handler(element)
				if err != nil {
					panic(err)
				}
				_, stateID = elementState.Parse()
				if stateID == 0 {
					break
				}
			}
		}
	}
}
