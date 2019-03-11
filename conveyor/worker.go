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
	"math/rand"
	"testing"
	"time"

	"github.com/insolar/insolar/conveyor/interfaces/constant"
	"github.com/insolar/insolar/conveyor/interfaces/iadapter"
	"github.com/insolar/insolar/conveyor/interfaces/islot"
	"github.com/insolar/insolar/conveyor/interfaces/istatemachine"
	"github.com/insolar/insolar/conveyor/queue"
	"github.com/insolar/insolar/log"
	"github.com/pkg/errors"
)

// WorkerState shows islot working mode
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
	stop               bool

	nodeState int // TODO: remove it when right implementation of node state calculation appears
}

func newWorkerStateMachineImpl(slot *Slot) workerStateMachineImpl {
	slot.slotState = Initializing
	return workerStateMachineImpl{
		slot:               slot,
		nextWorkerState:    Unknown,
		postponedResponses: make([]queue.OutputElement, 0),
		stop:               false,
	}
}

type MachineType int

const (
	InputEvent MachineType = iota + 1
	NestedCall
)

func GetStateMachineByType(mtype MachineType) istatemachine.StateMachineType {
	//panic("implement me") // TODO:
	sm := istatemachine.NewStateMachineTypeMock(&testing.T{})
	sm.GetMigrationHandlerFunc = func(p constant.PulseState, p1 uint32) (r istatemachine.MigrationHandler) {
		return func(element islot.SlotElementHelper) (interface{}, uint32, error) {
			r := rand.Int() % 100
			state := uint32(0)
			if r < 80 {
				state = 777777777
			}
			return element.GetElementID(), state, nil
		}
	}

	sm.GetTransitionHandlerFunc = func(p constant.PulseState, p1 uint32) (r istatemachine.TransitHandler) {
		return func(element islot.SlotElementHelper) (interface{}, uint32, error) {
			r := rand.Int() % 100
			state := uint32(0)
			if r < 80 {
				state = 888888888
			}
			return element.GetElementID(), state, nil
		}
	}
	return sm
}

func (w *workerStateMachineImpl) changePulseState() {
	switch w.slot.pulseState {
	case constant.Future:
		w.slot.pulseState = constant.Present
	case constant.Present:
		w.slot.pulseState = constant.Past
	case constant.Past:
		log.Error("[ changePulseState ] Try to change pulse state for 'Past' islot. Skip it")
	default:
		panic("[ changePulseState ] Unknown state: " + w.slot.pulseState.String())
	}
}

// TODO: is it ok?
type emptySyncDone struct{}

func (m emptySyncDone) Done() {}

// If we have both signals ( PendingPulseSignal and ActivatePulseSignal ),
// then change islot state and push ActivatePulseSignal back to queue.
func (w *workerStateMachineImpl) processSignalsWorking(elements []queue.OutputElement) int {
	numSignals := 0
	hasPending := false
	hasActivate := false
	for i := 0; i < len(elements); i++ {
		el := elements[i]
		if el.IsSignal() {
			numSignals++
			switch el.GetItemType() {
			case PendingPulseSignal:
				w.slot.slotState = Suspending
				hasPending = true
			case ActivatePulseSignal:
				hasActivate = true
				if hasPending {
					w.slot.inputQueue.PushSignal(ActivatePulseSignal, emptySyncDone{})
					break
				}
			default:
				panic(fmt.Sprintf("[ processSignalsWorking ] Unknown signal: %+v", el.GetItemType()))
			}
		} else {
			break
		}
	}

	if hasActivate && !hasPending {
		log.Error("[ processSignals ] Got ActivatePulseSignal and don't get PendingPulseSignal. Skip it. Continue working")
	}

	return numSignals
}

func (w *workerStateMachineImpl) readInputQueueWorking() error {
	elements := w.slot.inputQueue.RemoveAll()

	numSignals := w.processSignalsWorking(elements)
	// remove signals
	elements = elements[numSignals:]
	for i := 0; i < len(elements); i++ {
		el := elements[i]

		_, err := w.slot.createElement(GetStateMachineByType(InputEvent), 0, el)
		if err != nil {
			return errors.Wrapf(err, "[ readInputQueueWorking ] Can't createElement: %+v", el)
		}
	}

	return nil
}

func setNewElementState(element *slotElement, payLoad interface{}, fullState uint32) {
	if fullState == 0 {
		element.setDeleteState()
	} else {
		sm, state := extractStates(fullState)
		element.state = state
		element.payload = payLoad
		if sm != 0 {
			element.stateMachineType = GetStateMachineByType(MachineType(sm))
		}
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
			panic("Nested request is Not implemented")
		} else {
			adapterResp, ok := resp.GetData().(iadapter.IAdapterResponse)
			if !ok {
				panic(fmt.Sprintf("[ readResponseQueue ] Bad type in iadapter response queue: %T", resp.GetData()))
			}
			element := w.slot.extractSlotElementByID(adapterResp.GetElementID())
			if element == nil {
				log.Warnf("[ readResponseQueue ] Unknown element id: %d. AdapterResp: %+v", adapterResp.GetElementID(), adapterResp)
				continue
			}

			respHandler := element.stateMachineType.GetResponseHandler(w.slot.pulseState, element.state)

			payLoad, newState, err := respHandler(element, adapterResp)
			if err != nil {
				log.Error("[ readResponseQueue ] Response handler errors: ", err)
				respErrorHandler := element.stateMachineType.GetResponseErrorHandler(w.slot.pulseState, element.state)
				if respErrorHandler == nil {
					panic(fmt.Sprintf("[ readResponseQueue ] No response error handler. State: %d. AdapterResp: %+v", element.state, adapterResp))
				}

				payLoad, newState = respErrorHandler(element, adapterResp, err)
			}

			setNewElementState(element, payLoad, newState)
			w.slot.pushElement(element)
			err = w.slot.pushElement(element)
			if err != nil {
				return errors.Wrapf(err, "[ readResponseQueue ] Can't pushElement: %+v", element)
			}
		}

		numProcessedElements++

		if w.slot.inputQueue.HasSignal() {
			w.nextWorkerState = ReadInputQueue
			break
		}
	}

	w.postponedResponses = w.postponedResponses[numProcessedElements:]

	return nil
}

func (w *workerStateMachineImpl) waitQueuesOrTick() {
	log.Info("[ waitQueuesOrTick ] sleep ...")
	time.Sleep(time.Millisecond * 300)
	//panic("[ waitQueuesOrTick ] implement me") // TODO :
}

func (w *workerStateMachineImpl) processingElements() {
	if !w.slot.hasElements(ActiveElement) {
		if w.slot.pulseState == constant.Past {
			if w.slot.hasExpired() {
				w.slot.slotState = Suspending
				log.Info("[ processingElements ] Set islot state to 'Suspending'")
				return
			}
		}
		w.waitQueuesOrTick()
	}

	if w.slot.inputQueue.HasSignal() {
		log.Info("[ processingElements ] Set next worker state to 'ReadInputQueue'")
		w.nextWorkerState = ReadInputQueue
		return
	}

	lastState := uint32(0)
	numActiveElements := w.slot.len(ActiveElement)
	for ; numActiveElements > 0; numActiveElements-- {
		element := w.slot.popElement(ActiveElement)
		for lastState < element.state {
			lastState = element.state
			transitionHandler := element.stateMachineType.GetTransitionHandler(w.slot.pulseState, element.state)
			payLoad, newState, err := transitionHandler(element)

			if err != nil {
				log.Error("[ processingElements ] Transition handler error: ", err)
				errorHandler := element.stateMachineType.GetTransitionErrorHandler(w.slot.pulseState, element.state)
				payLoad, newState = errorHandler(element, err)
			}

			setNewElementState(element, payLoad, newState)
			w.slot.pushElement(element)

			if w.slot.inputQueue.HasSignal() {
				w.nextWorkerState = ReadInputQueue
				log.Info("[ processingElements ] Set next worker state to 'ReadInputQueue'")
				return
			}
		}
	}
}

func (w *workerStateMachineImpl) working() {

	for w.slot.isWorking() {
		err := w.readInputQueueWorking()
		if err != nil {
			panic(fmt.Sprintf("[ working ] readInputQueueWorking. Error: %s", err))
		}

		if !w.slot.isWorking() {
			log.Info("[ working ] Break after readInputQueueWorking")
			break
		}

		err = w.readResponseQueue()
		if err != nil {
			panic(fmt.Sprintf("[ working ] readResponseQueue. implement me: %s", err))
		}

		if !w.slot.isWorking() {
			log.Info("[ working ] Break after readResponseQueue")
			break
		}
		if w.nextWorkerState == ReadInputQueue {
			continue
		}

		w.processingElements()

		if !w.slot.isWorking() {
			log.Info("[ working ] Break after processingElements")
			break
		}
	}
}

func (w *workerStateMachineImpl) calculateNodeState() {
	// TODO: приходит PreparePulse, в нём есть callback, вызываем какой-то адаптер, куда передаем этот callback
	w.nodeState = 555
}

func (w *workerStateMachineImpl) sendRemovalSignalToConveyor() {
	w.slot.removeSlotCallback(w.slot.pulseNumber)
	// TODO: how to do it?
	// catch conveyor lock, check input queue, if It's empty - remove islot from map, if it's not - got to Working state
}

func (w *workerStateMachineImpl) processSignalsSuspending(elements []queue.OutputElement) int {
	numSignals := 0
	// TODO: add check if many signals come
	for i := 0; i < len(elements); i++ {
		el := elements[i]
		if el.IsSignal() {
			numSignals++
			switch el.GetItemType() {
			case PendingPulseSignal:
				log.Warn("[ processSignalsSuspending ] Must not be PendingPulseSignal here. Skip it")
			case ActivatePulseSignal:
				w.changePulseState()
				w.slot.slotState = Initializing
			default:
				panic(fmt.Sprintf("[ processSignalsSuspending ] Unknown signal: %+v", el.GetItemType()))
			}
		} else {
			break
		}
	}

	return numSignals
}

func (w *workerStateMachineImpl) readInputQueueSuspending() error {
	elements := w.slot.inputQueue.RemoveAll()
	numSignals := w.processSignalsSuspending(elements)

	// remove signals
	elements = elements[numSignals:]

	for i := 0; i < len(elements); i++ {
		el := elements[i]

		_, err := w.slot.createElement(GetStateMachineByType(InputEvent), 0, el)
		if err != nil {
			return errors.Wrap(err, "[ readInputQueue ] Can't createElement")
		}
	}

	if len(elements) != 0 && w.slot.pulseState == constant.Past {
		w.slot.slotState = Working
	}

	return nil
}

func (w *workerStateMachineImpl) suspending() {
	log.Info("[ suspending ] workerStateMachineImpl.suspending starts ...")
	switch w.slot.pulseState {
	case constant.Past:
		w.sendRemovalSignalToConveyor()
	case constant.Present:
		w.calculateNodeState()
	}
	for w.slot.isSuspending() {
		err := w.readInputQueueSuspending()
		if err != nil {
			panic(fmt.Sprintf("[ suspending ] readInputQueueSuspending. Can't readInputQueueSuspending: %s", err))
		}
	}

	log.Infof("[ suspending ] Leaving suspending. pulseState: %s. slotState: %s",
		w.slot.pulseState.String(),
		w.slot.slotState.String(),
	)
}

func (w *workerStateMachineImpl) migrate(status ActivationStatus) error {
	log.Infof("[ migrate ] Starts ... ( %s )", status.String())
	numElements := w.slot.len(status)
	for ; numElements > 0; numElements-- {
		element := w.slot.popElement(status)
		migHandler := element.stateMachineType.GetMigrationHandler(w.slot.pulseState, element.state)
		payLoad, newState := element.payload, element.state
		var err error
		if migHandler == nil {
			log.Infof("[ migrate ] No migration handler for pulseState: %d, element.state: %d. Nothing done", w.slot.pulseState, element.state)
			err = w.slot.pushElement(element)
			if err != nil {
				return errors.Wrapf(err, "[ migrate ] Can't pushElement: %+v", element)
			}
			continue
		} else {
			payLoad, newState, err = migHandler(element)
			if err != nil {
				log.Error("[ migrate ] Response handler errors: ", err)
				respErrorHandler := element.stateMachineType.GetTransitionErrorHandler(w.slot.pulseState, element.state)

				payLoad, newState = respErrorHandler(element, err)
			}
		}

		if newState == 0 {
			element.setDeleteState()
		}
		setNewElementState(element, payLoad, newState)

		err = w.slot.pushElement(element)
		if err != nil {
			return errors.Wrapf(err, "[ migrate ] Can't pushElement: %+v", element)
		}
	}

	log.Info("[ migrate ] END")
	return nil

}

func (w *workerStateMachineImpl) initializing() {
	if w.slot.pulseState == constant.Future {
		log.Info("[ initializing ] pulseState is Future. Skip initializing")
		return
	} else {
		// TODO: Get init handler from config
	}

	err := w.migrate(ActiveElement)
	if err != nil {
		panic("[ initializing ] migrate ActiveElement: " + err.Error())
	}
	w.migrate(NotActiveElement)
	if err != nil {
		panic("[ initializing ] migrate NotActiveElement: " + err.Error())
	}
}

func (w *workerStateMachineImpl) run() {
	for !w.stop {
		switch w.slot.slotState {
		case Initializing:
			w.initializing()
			w.slot.slotState = Working
		case Working:
			w.working()
		case Suspending:
			w.suspending()
		default:
			panic("[ run ] Unknown islot state: " + w.slot.slotState.String())
		}
	}
}
