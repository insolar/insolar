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
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/insolar/insolar/conveyor/interfaces/constant"
	"github.com/insolar/insolar/conveyor/interfaces/iadapter"
	"github.com/insolar/insolar/conveyor/interfaces/islot"
	"github.com/insolar/insolar/conveyor/interfaces/istatemachine"
	"github.com/insolar/insolar/conveyor/queue"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/instrumentation/inslogger"
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

	ctxLogger core.Logger
}

func newWorker(slot *Slot) worker {

	slot.slotState = Initializing

	w := worker{
		slot:               slot,
		nextWorkerState:    ReadInputQueue,
		postponedResponses: make([]queue.OutputElement, 0),
		stop:               false,
		ctxLogger:          inslogger.FromContext(context.Background()),
	}

	w.setLoggerFields()

	return w
}

type MachineType int

const (
	InputEvent MachineType = iota + 1
	NestedCall
)

func GetStateMachineByType(mtype MachineType) istatemachine.StateMachineType {
	//panic("implement me") // TODO:
	sm := istatemachine.NewStateMachineTypeMock(&testing.T{})
	sm.GetTransitionHandlerFunc = func(p constant.PulseState, p1 uint32) (r istatemachine.TransitHandler) {
		return func(element islot.SlotElementHelper) (interface{}, uint32, error) {
			return nil, 0, nil
		}
	}
	sm.GetMigrationHandlerFunc = func(p constant.PulseState, p1 uint32) (r istatemachine.MigrationHandler) {
		return func(element islot.SlotElementHelper) (interface{}, uint32, error) {
			return nil, 0, nil
		}
	}
	return sm
}

func (w *worker) setLoggerFields() {
	ctx, _ := inslogger.WithField(context.Background(), "pulseState", w.slot.pulseState.String())
	_, w.ctxLogger = inslogger.WithField(ctx, "slotState", w.slot.slotState.String())
}

func (w *worker) changePulseState() {
	w.ctxLogger.Debugf("[ changePulseState ] starts ... ( pulseState: %s )", w.slot.pulseState.String())
	switch w.slot.pulseState {
	case constant.Future:
		w.slot.pulseState = constant.Present
	case constant.Present:
		w.slot.pulseState = constant.Past
	case constant.Past:
		w.ctxLogger.Error("[ changePulseState ] Try to change pulse state for 'Past' slot. Skip it")
	case constant.Antique:
		w.ctxLogger.Error("[ changePulseState ] Try to change pulse state for 'Antique' slot. Skip it")
	default:
		panic("[ changePulseState ] Unknown state: " + w.slot.pulseState.String())
	}
	w.setLoggerFields()
}

func (w *worker) processPendingPulseSignalWorking(hasActivate bool, element *queue.OutputElement, activateSyncDone queue.SyncDone) bool {
	w.slot.slotState = Suspending
	w.ctxLogger.Info("[ processSignalsWorking ] Got PendingPulseSignal. Set slot state to 'Suspending'")
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
	w.ctxLogger.Debugf("[ processSignalsWorking ] starts ... ( len: %d. pulseState: %s", len(elements), w.slot.pulseState.String())
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
				w.ctxLogger.Info("[ processSignalsWorking ] Got ActivatePulseSignal")
				hasActivate = true
				activateSyncDone = el.GetData().(queue.SyncDone)
			case CancelSignal:
				w.stop = true // TODO: do it more correctly
				w.slot.slotState = Suspending
				w.ctxLogger.Info("[ processSignalsWorking ] Got CancelSignal. Set slot state to 'Suspending'")
			default:
				panic(fmt.Sprintf("[ processSignalsWorking ] Unknown signal: %+v", el.GetItemType()))
			}
		} else {
			break
		}
	}

	if hasActivate && !hasPending {
		w.ctxLogger.Error("[ processSignals ] Got ActivatePulseSignal and don't get PendingPulseSignal. Skip it. Continue working")
	}

	return numSignals
}

func (w *worker) readInputQueueWorking() error {
	w.ctxLogger.Debugf("[ readInputQueueWorking ] starts ... ( pulseState: %s )", w.slot.pulseState.String())
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

func updateElement(element *slotElement, payload interface{}, fullState uint32) {
	log.Debugf("[ updateElement ] starts ... ( element: %+v. fullstate: %d )", element, fullState)
	if fullState != 0 {
		sm, state := extractStates(fullState)
		machineType := element.stateMachineType
		if sm != 0 {
			machineType = GetStateMachineByType(MachineType(sm))
		}
		element.update(state, payload, machineType)
		return
	}
	element.setDeleteState()
}

func (w *worker) processResponse(resp queue.OutputElement) error {
	adapterResp, ok := resp.GetData().(iadapter.IAdapterResponse)
	if !ok {
		panic(fmt.Sprintf("[ processResponse ] Bad type in adapter response queue: %T", resp.GetData()))
	}
	element := w.slot.extractSlotElementByID(adapterResp.GetElementID())
	if element == nil {
		w.ctxLogger.Warnf("[ processResponse ] Unknown element id: %d. AdapterResp: %+v", adapterResp.GetElementID(), adapterResp)
		return nil
	}

	respHandler := element.stateMachineType.GetResponseHandler(w.slot.pulseState, element.state)

	payload, newState, err := respHandler(element, adapterResp)
	if err != nil {
		w.ctxLogger.Error("[ processResponse ] Response handler errors: ", err)
		respErrorHandler := element.stateMachineType.GetResponseErrorHandler(w.slot.pulseState, element.state)
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
	w.ctxLogger.Debugf("[ readResponseQueue ] starts ... ( pulseState: %s )", w.slot.pulseState.String())
	w.ctxLogger.Info("[ readResponseQueue ] Set next worker state to 'ProcessElements'")
	w.nextWorkerState = ProcessElements
	w.postponedResponses = append(w.postponedResponses, w.slot.responseQueue.RemoveAll()...)

	if w.slot.pulseState == constant.Future {
		w.postponedResponses = make([]queue.OutputElement, 0)
		w.ctxLogger.Warnf("[ readResponseQueue ] Shouldn't have responses in 'Future'. Drop them ( %d )", len(w.postponedResponses))
		return nil
	}

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
			w.ctxLogger.Info("[ readResponseQueue ] Set next worker state to 'ReadInputQueue'")
			w.nextWorkerState = ReadInputQueue
			break
		}
	}

	w.postponedResponses = w.postponedResponses[numProcessedElements:]

	return nil
}

func (w *worker) waitQueuesOrTick() {
	w.ctxLogger.Debugf("[ waitQueuesOrTick ] starts ... ( pulseState: %s )", w.slot.pulseState.String())
	time.Sleep(time.Millisecond * 300)
	//panic("[ waitQueuesOrTick ] implement me") // TODO :
}

func (w *worker) processingElements() {
	w.ctxLogger.Debugf("[ processingElements ] starts ... ( pulseState: %s )", w.slot.pulseState.String())
	if !w.slot.hasElements(ActiveElement) {
		if w.slot.pulseState == constant.Past {
			if w.slot.hasExpired() {
				w.slot.slotState = Suspending
				w.ctxLogger.Info("[ processingElements ] Set slot state to 'Suspending'")
				return
			}
		}
		w.waitQueuesOrTick()
	}

	if w.slot.inputQueue.HasSignal() {
		w.ctxLogger.Info("[ processingElements ] Set next worker state to 'ReadInputQueue'")
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
				w.ctxLogger.Info("[ processingElements ] Set next worker state to 'ReadInputQueue'")
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
	transitionHandler := element.stateMachineType.GetTransitionHandler(w.slot.pulseState, element.state)
	payload, newState, err := transitionHandler(element)
	if err != nil {
		w.ctxLogger.Error("[ processingElements ] Transition handler error: ", err)
		errorHandler := element.stateMachineType.GetTransitionErrorHandler(w.slot.pulseState, element.state)
		payload, newState = errorHandler(element, err)
	}
	updateElement(element, payload, newState)

	stopProcessingElement := (newState == 0) || element.isDeactivated()

	return stopProcessingElement
}

func (w *worker) working() {
	w.ctxLogger.Debugf("[ working ] starts ... ( pulseState: %s )", w.slot.pulseState.String())

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
	w.ctxLogger.Debugf("[ calculateNodeState ] starts ... ( pulseState: %s )", w.slot.pulseState.String())
	// TODO: приходит PreparePulse, в нём есть callback, вызываем какой-то адаптер, куда передаем этот callback
	w.preparePulseSync.SetResult(555)
}

func (w *worker) sendRemovalSignalToConveyor() {
	w.ctxLogger.Debugf("[ sendRemovalSignalToConveyor ] starts ... ( pulseState: %s )", w.slot.pulseState.String())
	w.slot.removeSlotCallback(w.slot.pulseNumber)
	// TODO: how to do it?
	// catch conveyor lock, check input queue, if It's empty - remove slot from map, if it's not - got to Working state
}

func (w *worker) processSignalsSuspending(elements []queue.OutputElement) int {
	w.ctxLogger.Debugf("[ processSignalsSuspending ] starts ... ( pulseState: %s )", w.slot.pulseState.String())
	numSignals := 0
	// TODO: add check if many signals come
	for i := 0; i < len(elements); i++ {
		el := elements[i]
		if el.IsSignal() {
			numSignals++
			switch el.GetItemType() {
			case PendingPulseSignal:
				w.ctxLogger.Warn("[ processSignalsSuspending ] Must not be PendingPulseSignal here. Skip it")
			case ActivatePulseSignal:
				w.changePulseState()
				w.slot.slotState = Initializing
				w.activatePulseSync = el.GetData().(queue.SyncDone)
				w.ctxLogger.Info("[ processSignalsSuspending ] Set slot state to 'Initializing'")
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
	w.ctxLogger.Debugf("[ readInputQueueSuspending ] starts ... ( pulseState: %s )", w.slot.pulseState.String())
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
		w.ctxLogger.Info("[ readInputQueueSuspending ] Set slot state to 'Working'")
	}

	return nil
}

func (w *worker) suspending() {
	w.ctxLogger.Debugf("[ suspending ] starts ... ( pulseState: %s )", w.slot.pulseState.String())
	switch w.slot.pulseState {
	case constant.Past:
		w.sendRemovalSignalToConveyor()
	case constant.Present:
		w.calculateNodeState()
	case constant.Future:
		if w.preparePulseSync != nil {
			w.preparePulseSync.SetResult(nil)
		} else {
			w.ctxLogger.Warn("[ suspending ] preparePulseSync is empty")
		}
	}

	for w.slot.isSuspending() {
		err := w.readInputQueueSuspending()
		if err != nil {
			panic(fmt.Sprintf("[ suspending ] readInputQueueSuspending. Can't readInputQueueSuspending: %s", err))
		}
	}

	w.ctxLogger.Infof("[ suspending ] Leaving suspending. pulseState: %s. slotState: %s",
		w.slot.pulseState.String(),
		w.slot.slotState.String(),
	)
}

func (w *worker) migrate(status ActivationStatus) error {
	w.ctxLogger.Infof("[ migrate ] Starts ... ( status: %s. pulseState: %s )", status.String(), w.slot.pulseState.String())
	numElements := w.slot.len(status)
	for ; numElements > 0; numElements-- {
		element := w.slot.popElement(status)
		migHandler := element.stateMachineType.GetMigrationHandler(w.slot.pulseState, element.state)
		var err error
		if migHandler == nil {
			w.ctxLogger.Infof("[ migrate ] No migration handler for pulseState: %d, element.state: %d. Nothing done", w.slot.pulseState, element.state)
			err = w.slot.pushElement(element)
			if err != nil {
				return errors.Wrapf(err, "[ migrate ] Can't pushElement: %+v", element)
			}
			continue
		}

		payload, newState, err := migHandler(element)
		if err != nil {
			w.ctxLogger.Error("[ migrate ] Response handler errors: ", err)
			respErrorHandler := element.stateMachineType.GetTransitionErrorHandler(w.slot.pulseState, element.state)

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
	w.ctxLogger.Debugf("[ initializing ] starts ... ( pulseState: %s )", w.slot.pulseState.String())
	if w.slot.pulseState == constant.Future {
		w.ctxLogger.Info("[ initializing ] pulseState is Future. Skip initializing")
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
			w.ctxLogger.Info("[ run ] Set slot state to 'Working'")
		case Working:
			if w.activatePulseSync != nil {
				w.activatePulseSync.SetResult(nil)
			} else {
				w.ctxLogger.Warn("[ run ] activatePulseSync is empty")
			}
			w.working()
		case Suspending:
			w.suspending()
		default:
			panic("[ run ] Unknown slot state: " + w.slot.slotState.String())
		}
	}
}
