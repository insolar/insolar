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
	"testing"
	"time"

	"github.com/insolar/insolar/conveyor/interfaces/constant"
	"github.com/insolar/insolar/conveyor/interfaces/fsm"
	"github.com/insolar/insolar/conveyor/interfaces/iadapter"
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
	ReadResponseQueue
	ProcessElements
)

type worker struct { // nolint: unused
	slot               *Slot
	nextWorkerState    WorkerState
	postponedResponses []queue.OutputElement
	stop               bool

	activatePulseSync queue.SyncDone
	preparePulseSync  queue.SyncDone

	nodeState int // TODO: remove it when right implementation of node state calculation appears
}

func newWorker(slot *Slot) worker {
	slot.slotState = Initializing
	return worker{
		slot:               slot,
		nextWorkerState:    ReadInputQueue,
		postponedResponses: make([]queue.OutputElement, 0),
		stop:               false,
	}
}

type MachineType int

const (
	InputEvent MachineType = iota + 1
	NestedCall
)

func GetStateMachineByType(mtype MachineType) statemachine.StateMachine {
	//panic("implement me") // TODO:
	sm := statemachine.NewStateMachineMock(&testing.T{})
	return sm
}

func (w *worker) changePulseState() {
	log.Debugf("[ changePulseState ] starts ... ( pulseState: %s )", w.slot.pulseState.String())
	switch w.slot.pulseState {
	case constant.Future:
		w.slot.pulseState = constant.Present
	case constant.Present:
		w.slot.pulseState = constant.Past
	case constant.Past:
		log.Error("[ changePulseState ] Try to change pulse state for 'Past' slot. Skip it")
	case constant.Antique:
		log.Error("[ changePulseState ] Try to change pulse state for 'Antique' slot. Skip it")
	default:
		panic("[ changePulseState ] Unknown state: " + w.slot.pulseState.String())
	}
}

func (w *worker) processPendingPulseSignalWorking(hasActivate bool, element *queue.OutputElement, activateSyncDone queue.SyncDone) bool {
	w.slot.slotState = Suspending
	log.Info("[ processSignalsWorking ] Got PendingPulseSignal. Set slot state to 'Suspending'")
	w.preparePulseSync = element.GetData().(queue.SyncDone)
	if hasActivate {
		err := w.slot.inputQueue.PushSignal(ActivatePulseSignal, activateSyncDone)
		if err != nil {
			panic("[ processSignalsWorking ] Can't PushSignal: " + err.Error())
		}
		return false
	}

	return true
}

// If we have both signals ( PendingPulseSignal and ActivatePulseSignal ),
// then change slot state and push ActivatePulseSignal back to queue.
func (w *worker) processSignalsWorking(elements []queue.OutputElement) int {
	log.Debugf("[ processSignalsWorking ] starts ... ( len: %d. pulseState: %s", w.slot.pulseState.String(), len(elements))
	numSignals := 0
	hasPending := false
	hasActivate := false
	var activateSyncDone queue.SyncDone
	for i := 0; i < len(elements); i++ {
		el := elements[i]
		if el.IsSignal() {
			numSignals++
			switch el.GetItemType() {
			case PendingPulseSignal:
				hasPending = true
				if !w.processPendingPulseSignalWorking(hasActivate, &el, activateSyncDone) {
					break
				}
			case ActivatePulseSignal:
				log.Info("[ processSignalsWorking ] Got ActivatePulseSignal")
				hasActivate = true
				activateSyncDone = el.GetData().(queue.SyncDone)
			case CancelSignal:
				w.stop = true // TODO: do it more correctly
				w.slot.slotState = Suspending
				log.Info("[ processSignalsWorking ] Got Cancel. Set slot state to 'Suspending'")
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

func (w *worker) readInputQueueWorking() error {
	log.Debugf("[ readInputQueueWorking ] starts ... ( pulseState: %s )", w.slot.pulseState.String())
	w.nextWorkerState = ReadResponseQueue
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

func updateElement(element *slotElement, payload interface{}, fullState fsm.ElementState) {
	log.Debugf("[ updateElement ] starts ... ( element: %+v. fullstate: %d )", element, fullState)
	if fullState != 0 {
		sm, state := fullState.Parse()
		machineType := element.stateMachine
		if sm != 0 {
			machineType = GetStateMachineByType(MachineType(sm))
		}
		element.update(state, payload, machineType)
		return
	}
	element.setDeleteState()
}

func (w *worker) processResponse(resp queue.OutputElement) error {
	adapterResp, ok := resp.GetData().(iadapter.Response)
	if !ok {
		panic(fmt.Sprintf("[ processResponse ] Bad type in adapter response queue: %T", resp.GetData()))
	}
	element := w.slot.extractSlotElementByID(adapterResp.GetElementID())
	if element == nil {
		log.Warnf("[ processResponse ] Unknown element id: %d. AdapterResp: %+v", adapterResp.GetElementID(), adapterResp)
		return nil
	}

	respHandler := element.stateMachine.GetResponseHandler(element.state)

	payload, newState, err := respHandler(element, adapterResp)
	if err != nil {
		log.Error("[ processResponse ] Response handler errors: ", err)
		respErrorHandler := element.stateMachine.GetResponseErrorHandler(element.state)
		if respErrorHandler == nil {
			panic(fmt.Sprintf("[ processResponse ] No response error handler. State: %d. AdapterResp: %+v", element.state, adapterResp))
		}

		payload, newState = respErrorHandler(element, adapterResp, err)
	}

	updateElement(element, payload, newState)
	err = w.slot.pushElement(element)
	if err != nil {
		return errors.Wrapf(err, "[ processResponse ] Can't pushElement: %+v", element)
	}

	return nil
}

func (w *worker) processNestedEvent(resp queue.OutputElement) {
	panic("Nested request is Not implemented")
}

func (w *worker) readResponseQueue() error {
	log.Debugf("[ readResponseQueue ] starts ... ( pulseState: %s )", w.slot.pulseState.String())
	w.nextWorkerState = ProcessElements
	w.postponedResponses = append(w.postponedResponses, w.slot.responseQueue.RemoveAll()...)

	totalNumElements := len(w.postponedResponses)
	numProcessedElements := 0
	for i := 0; i < totalNumElements; i++ {
		resp := w.postponedResponses[i]
		if resp.GetItemType() > 9999 { // TODO: check isNestedEvent
			w.processNestedEvent(resp)
		} else {
			err := w.processResponse(resp)
			if err != nil {
				return errors.Wrap(err, "[ readResponseQueue ] Can't processResponse")
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

func (w *worker) waitQueuesOrTick() {
	log.Debugf("[ waitQueuesOrTick ] starts ... ( pulseState: %s )", w.slot.pulseState.String())
	time.Sleep(time.Millisecond * 300)
	//panic("[ waitQueuesOrTick ] implement me") // TODO :
}

func (w *worker) processingElements() {
	log.Debugf("[ processingElements ] starts ... ( pulseState: %s )", w.slot.pulseState.String())
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
		log.Info("[ processingElements ] Set next worker state to 'ReadInputQueue'")
		w.nextWorkerState = ReadInputQueue
		return
	}

	numActiveElements := w.slot.len(ActiveElement)
	breakProcessing := false
	for ; numActiveElements > 0 && !breakProcessing; numActiveElements-- {
		element := w.slot.popElement(ActiveElement)
		lastState := -1
		for lastState < int(element.state) {
			lastState = int(element.state)

			if w.processOneElement(element) {
				break
			}

			if w.slot.inputQueue.HasSignal() {
				w.nextWorkerState = ReadInputQueue
				log.Info("[ processingElements ] Set next worker state to 'ReadInputQueue'")
				breakProcessing = true
				break
			}
		}

		err := w.slot.pushElement(element)
		if err != nil {
			panic(fmt.Sprintf("[ processingElements ] Can't push element: %+v", element))
		}
	}
}

func (w *worker) processOneElement(element *slotElement) bool {
	transitionHandler := element.stateMachine.GetTransitionHandler(element.state)
	payload, newState, err := transitionHandler(element)
	if err != nil {
		log.Error("[ processingElements ] Transition handler error: ", err)
		errorHandler := element.stateMachine.GetTransitionErrorHandler(element.state)
		payload, newState = errorHandler(element, err)
	}
	updateElement(element, payload, newState)

	stopProcessingElement := (newState == 0) || element.isDeactivated()

	return stopProcessingElement
}

func (w *worker) working() {
	log.Debugf("[ working ] starts ... ( pulseState: %s )", w.slot.pulseState.String())
	for w.slot.isWorking() {
		switch w.nextWorkerState {
		case ReadInputQueue:
			err := w.readInputQueueWorking()
			if err != nil {
				panic(fmt.Sprintf("[ working ] readInputQueueWorking. Error: %s", err))
			}
		case ReadResponseQueue:
			err := w.readResponseQueue()
			if err != nil {
				panic(fmt.Sprintf("[ working ] readResponseQueue. Error: %s", err))
			}
		case ProcessElements:
			w.processingElements()
		default:
			panic("[ working ] unknown nextWorkerState: " + w.nextWorkerState.String())
		}
	}
}

func (w *worker) calculateNodeState() {
	log.Debugf("[ calculateNodeState ] starts ... ( pulseState: %s )", w.slot.pulseState.String())
	// TODO: приходит PreparePulse, в нём есть callback, вызываем какой-то адаптер, куда передаем этот callback
	w.nodeState = 555
	if w.preparePulseSync != nil {
		w.preparePulseSync.Done()
	} else {
		log.Warn("[ calculateNodeState ] preparePulseSync is empty ")
	}
}

func (w *worker) sendRemovalSignalToConveyor() {
	log.Debugf("[ sendRemovalSignalToConveyor ] starts ... ( pulseState: %s )", w.slot.pulseState.String())
	w.slot.removeSlotCallback(w.slot.pulseNumber)
	// TODO: how to do it?
	// catch conveyor lock, check input queue, if It's empty - remove slot from map, if it's not - got to Working state
}

func (w *worker) processSignalsSuspending(elements []queue.OutputElement) int {
	log.Debugf("[ processSignalsSuspending ] starts ... ( pulseState: %s )", w.slot.pulseState.String())
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
				w.activatePulseSync = el.GetData().(queue.SyncDone)
				log.Info("[ processSignalsSuspending ] Set slot state to 'Initializing'")
			case CancelSignal:
				w.stop = true // TODO: do it more correctly
			default:
				panic(fmt.Sprintf("[ processSignalsSuspending ] Unknown signal: %+v", el.GetItemType()))
			}
		} else {
			break
		}
	}

	return numSignals
}

func (w *worker) readInputQueueSuspending() error {
	log.Debugf("[ readInputQueueSuspending ] starts ... ( pulseState: %s )", w.slot.pulseState.String())
	elements := w.slot.inputQueue.RemoveAll()
	numSignals := w.processSignalsSuspending(elements)

	// remove signals
	elements = elements[numSignals:]

	for i := 0; i < len(elements); i++ {
		el := elements[i]

		_, err := w.slot.createElement(GetStateMachineByType(InputEvent), 0, el)
		if err != nil {
			return errors.Wrap(err, "[ readInputQueueSuspending ] Can't createElement")
		}
	}

	if len(elements) != 0 && w.slot.pulseState == constant.Past {
		w.slot.slotState = Working
		log.Info("[ readInputQueueSuspending ] Set slot state to 'Working'")
	}

	return nil
}

func (w *worker) suspending() {
	log.Debugf("[ suspending ] starts ... ( pulseState: %s )", w.slot.pulseState.String())
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

func (w *worker) migrate(status ActivationStatus) error {
	log.Infof("[ migrate ] Starts ... ( status: %s. pulseState: %s )", status.String(), w.slot.pulseState.String())
	numElements := w.slot.len(status)
	for ; numElements > 0; numElements-- {
		element := w.slot.popElement(status)
		migHandler := element.stateMachine.GetMigrationHandler(element.state)
		var err error
		if migHandler == nil {
			log.Infof("[ migrate ] No migration handler for pulseState: %d, element.state: %d. Nothing done", w.slot.pulseState, element.state)
			err = w.slot.pushElement(element)
			if err != nil {
				return errors.Wrapf(err, "[ migrate ] Can't pushElement: %+v", element)
			}
			continue
		}

		payload, newState, err := migHandler(element)
		if err != nil {
			log.Error("[ migrate ] Response handler errors: ", err)
			respErrorHandler := element.stateMachine.GetTransitionErrorHandler(element.state)

			payload, newState = respErrorHandler(element, err)
		}

		updateElement(element, payload, newState)

		err = w.slot.pushElement(element)
		if err != nil {
			return errors.Wrapf(err, "[ migrate ] Can't pushElement: %+v", element)
		}
	}

	return nil

}

func (w *worker) getInitHandlersFromConfig() {
	// TODO: impolement me
}

func (w *worker) initializing() {
	log.Debugf("[ initializing ] starts ... ( pulseState: %s )", w.slot.pulseState.String())
	if w.slot.pulseState == constant.Future {
		log.Info("[ initializing ] pulseState is Future. Skip initializing")
		return
	}

	w.getInitHandlersFromConfig()

	err := w.migrate(ActiveElement)
	if err != nil {
		panic("[ initializing ] migrate ActiveElement: " + err.Error())
	}
	err = w.migrate(NotActiveElement)
	if err != nil {
		panic("[ initializing ] migrate NotActiveElement: " + err.Error())
	}
}

func (w *worker) run() {
	for !w.stop {
		switch w.slot.slotState {
		case Initializing:
			w.initializing()
			w.slot.slotState = Working
			log.Info("[ run ] Set slot state to 'Working'")
		case Working:
			if w.activatePulseSync != nil {
				w.activatePulseSync.Done()
			} else {
				log.Warn("[ run ] activatePulseSync is empty")
			}
			w.working()
		case Suspending:
			w.suspending()
		default:
			panic("[ run ] Unknown slot state: " + w.slot.slotState.String())
		}
	}
}
