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

	"github.com/insolar/insolar/conveyor/adapter"
	"github.com/insolar/insolar/conveyor/interfaces/statemachine"
	"github.com/pkg/errors"
)

// WorkerState shows slot working mode
type WorkerState int

//go:generate stringer -type=WorkerState
const (
	Unknown = WorkerState(iota)
	ReadInputQueue
	ProcessElements
)

type workerStateMachineImpl struct {
	slot            *Slot
	slotState       SlotState
	nextWorkerState WorkerState
}

func newWorkerStateMachineImpl(slot *Slot) workerStateMachineImpl {
	return workerStateMachineImpl{
		slot:            slot,
		slotState:       Initializing,
		nextWorkerState: Unknown,
	}
}

type MachineType int

const (
	InputEvent MachineType = iota + 1
	NestedCall
)

func GetStateMachineByType(mtype MachineType) statemachine.StateMachineType {
	return nil
}

func (w *workerStateMachineImpl) isWorkingState() bool {
	return w.slotState == Working
}

func (w *workerStateMachineImpl) readInputQueue() error {
	elements := w.slot.inputQueue.RemoveAll()
	for i := 0; i < len(elements); i++ {
		el := elements[i]
		// check is it signal
		if el.GetItemType() > 0 {
			switch el.GetItemType() {
			case PendingPulseSignal:
				panic("implement me")
			case ActivatePulseSignal:
				panic("implement me")
			default:
				panic("implement me")
			}
		} else {
			// TODO: do it in one step
			el, err := w.slot.createElement(GetStateMachineByType(InputEvent), 0, el)
			if err != nil {
				return errors.Wrap(err, "[ readInputQueue ] Can't createElement")
			}
			err = w.slot.pushElement(ActiveElement, el)
			if err != nil {
				return errors.Wrap(err, "[ readInputQueue ] Can't pushElement")
			}
		}
	}

	return nil
}

func (w *workerStateMachineImpl) readResponseQueue() error {
	postponedResponses := w.slot.responseQueue.RemoveAll()
	w.nextWorkerState = ProcessElements

	for i := 0; i < len(postponedResponses); i++ {
		resp := postponedResponses[i]
		if resp.GetItemType() > 9999 { // TODO: check isNestedEvent

		} else {
			adapterResp, ok := resp.GetData().(adapter.AdapterResponse)
			if !ok {
				panic(fmt.Sprintf("Bad type in adapter response queue: %T", resp.GetData()))
			}
			element := w.slot.elements[adapterResp.ElementID]

			element.stateMachineType.GetResponseHandler(element.state)
			// Call ReponseHandler

		}
	}

	return nil
}

func (w *workerStateMachineImpl) working() {

	for w.isWorkingState() {
		err := w.readInputQueue()
		if err != nil {
			panic("implement me")
		}

	}

}
