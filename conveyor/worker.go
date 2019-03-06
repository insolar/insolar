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

	adapter2 "github.com/insolar/insolar/conveyor/interfaces/adapter"
	"github.com/insolar/insolar/conveyor/interfaces/constant"
	"github.com/insolar/insolar/conveyor/interfaces/statemachine"
	"github.com/insolar/insolar/conveyor/queue"
	"github.com/insolar/insolar/log"
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
	slot               *Slot
	nextWorkerState    WorkerState
	postponedResponses []queue.OutputElement
}

func newWorkerStateMachineImpl(slot *Slot) workerStateMachineImpl {
	slot.slotState = Initializing
	return workerStateMachineImpl{
		slot:            slot,
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
	return w.slot.slotState == Working
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
			_, err := w.slot.createElement(GetStateMachineByType(InputEvent), 0, el)
			if err != nil {
				return errors.Wrap(err, "[ readInputQueue ] Can't createElement")
			}
		}
	}

	return nil
}

func setNewState(element *slotElement, payLoad interface{}, fullState uint32) {
	sm, state := extractStates(fullState)
	element.state = state
	element.payload = payLoad
	if sm != 0 {
		element.stateMachineType = GetStateMachineByType(MachineType(sm))
	}
}

func (w *workerStateMachineImpl) readResponseQueue() error {
	w.postponedResponses = append(w.postponedResponses, w.slot.responseQueue.RemoveAll()...)
	w.nextWorkerState = ProcessElements

	totalNumElements := len(w.postponedResponses)
	numProcessedElements := 0
	for i := 0; i < totalNumElements; i++ {
		resp := w.postponedResponses[i]
		if resp.GetItemType() > 9999 {
			// TODO: check isNestedEvent
		} else {
			adapterResp, ok := resp.GetData().(adapter2.IAdapterResponse)
			if !ok {
				panic(fmt.Sprintf("Bad type in adapter response queue: %T", resp.GetData()))
			}
			element := w.slot.getSlotElementByID(adapterResp.GetElementID())

			respHandler := element.stateMachineType.GetResponseHandler(w.slot.pulseState, element.state)
			if respHandler == nil {
				panic(fmt.Sprintf("No response handler. State: %d. \nAdapterResp: %+v", element.state, adapterResp))
			}

			payLoad, newState, err := respHandler(element, adapterResp)
			if err != nil {
				log.Error("[ readResponseQueue ] Response handler errors: ", err)
				respErrorHandler := element.stateMachineType.GetResponseErrorHandler(w.slot.pulseState, element.state)
				if respErrorHandler == nil {
					panic(fmt.Sprintf("No response error handler. State: %d. \nAdapterResp: %+v", element.state, adapterResp))
				}

				payLoad, newState = respErrorHandler(element, adapterResp, err)
			}

			if newState == 0 {
				// TODO: call finalization handler
			}

			setNewState(element, payLoad, newState)
		}

		numProcessedElements++

		if w.slot.inputQueue.HasSignal() {
			w.nextWorkerState = ReadInputQueue
			break
		}
	}

	w.postponedResponses = w.postponedResponses[totalNumElements:]

	return nil
}

func (w *workerStateMachineImpl) waitQueuesOrTick() {
	panic("implement me")
}

func (w *workerStateMachineImpl) processingElements() {
	if !w.slot.hasElements(ActiveElement) {
		if w.slot.pulseState == constant.Past {
			if w.slot.hasExpired() {
				w.slot.slotState = Suspending
				log.Info("[ processingElements ] Set slot state to 'Suspending'")
				return
			}
		}
		w.waitQueuesOrTick()
	}

	if w.slot.inputQueue.HasSignal() {
		w.nextWorkerState = ReadInputQueue
		return
	}

	element := w.slot.popElement(ActiveElement)
	lastState := element.state
	for ; element != nil; element = w.slot.popElement(ActiveElement) {
		for lastState < element.state {
			lastState = element.state
			transitionHandler := element.stateMachineType.GetTransitionHandler(w.slot.pulseState, element.state)
			payLoad, newState, err := transitionHandler(element)

			if err != nil {
				log.Error("[ processingElements ] Transition handler error: ", err)
				errorHandler := element.stateMachineType.GetTransitionErrorHandler(w.slot.pulseState, element.state)
				payLoad, newState = errorHandler(element, err)
			}

			if newState == 0 {
				element.setDeleteState()
			}
			setNewState(element, payLoad, newState)
			w.slot.pushElement(element)

			if w.slot.inputQueue.HasSignal() {
				w.nextWorkerState = ReadInputQueue
				return
			}
		}

	}
}

func (w *workerStateMachineImpl) working() {

	for w.isWorkingState() {
		err := w.readInputQueue()
		if err != nil {
			panic("implement me")
		}

		if !w.isWorkingState() {
			break
		}

		err = w.readResponseQueue()
		if err != nil {
			panic("implement me")
		}

		if !w.isWorkingState() {
			break
		}
		if w.nextWorkerState == ReadInputQueue {
			continue
		}

		w.processingElements()

		if !w.isWorkingState() {
			break
		}
		if w.nextWorkerState == ReadInputQueue {
			continue
		}

	}
}

func (w *workerStateMachineImpl) suspending() {

}
