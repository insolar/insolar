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
	"testing"
	"time"

	"github.com/insolar/insolar/conveyor/adapter"
	"github.com/insolar/insolar/conveyor/interfaces/constant"
	"github.com/insolar/insolar/conveyor/interfaces/iadapter"
	"github.com/insolar/insolar/conveyor/interfaces/islot"
	"github.com/insolar/insolar/conveyor/interfaces/istatemachine"
	"github.com/insolar/insolar/conveyor/queue"
	"github.com/insolar/insolar/core"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

var testPulseStates = []constant.PulseState{constant.Future, constant.Present, constant.Past}

func makeSlotAndWorker(pulseState constant.PulseState, pulseNumber core.PulseNumber) (*Slot, worker) {
	slot := NewSlot(pulseState, pulseNumber, nil)
	worker := newWorker(slot)
	slot.removeSlotCallback = func(number core.PulseNumber) {}

	return slot, worker
}

func Test_changePulseState(t *testing.T) {
	slot, worker := makeSlotAndWorker(constant.Future, 22)

	worker.changePulseState()
	require.Equal(t, constant.Present, slot.pulseState)

	worker.changePulseState()
	require.Equal(t, constant.Past, slot.pulseState)

	worker.changePulseState()
	require.Equal(t, constant.Past, slot.pulseState)

	slot.pulseState = 99999
	require.PanicsWithValue(t, "[ changePulseState ] Unknown state: PulseState(99999)", worker.changePulseState)
}

func areSlotStatesEqual(s1 *Slot, s2 *Slot, t *testing.T, excludePulseStateCheck bool) {
	if !excludePulseStateCheck {
		require.Equal(t, s1.pulseState, s2.pulseState)
	}
	require.Equal(t, s1.stateMachine, s2.stateMachine)
	require.Equal(t, s1.pulse, s2.pulse)
	require.Equal(t, s1.slotState, s2.slotState)
}

// ---- processSignalsWorking

func Test_processSignalsWorking_EmptyInput(t *testing.T) {

	for _, tt := range testPulseStates {
		t.Run(tt.String(), func(t *testing.T) {
			slot, worker := makeSlotAndWorker(tt, 22)

			oldSlot := *slot
			require.Equal(t, 0, worker.processSignalsWorking([]queue.OutputElement{}))
			areSlotStatesEqual(&oldSlot, slot, t, false)
		})
	}
}

func Test_processSignalsWorking_NonSignals(t *testing.T) {

	for _, tt := range testPulseStates {
		t.Run(tt.String(), func(t *testing.T) {
			slot, worker := makeSlotAndWorker(tt, 22)
			oldSlot := *slot

			nonSignals := []queue.OutputElement{
				*queue.NewOutputElement(emptySyncDone{}, 0),
				*queue.NewOutputElement(emptySyncDone{}, 0),
				*queue.NewOutputElement(emptySyncDone{}, 0),
			}
			require.Equal(t, 0, worker.processSignalsWorking(nonSignals))

			areSlotStatesEqual(&oldSlot, slot, t, false)
		})
	}
}

func Test_processSignalsWorking_BadSignal(t *testing.T) {

	for _, tt := range testPulseStates {
		t.Run(tt.String(), func(t *testing.T) {
			slot, worker := makeSlotAndWorker(tt, 22)
			oldSlot := *slot

			badSignal := []queue.OutputElement{*queue.NewOutputElement(emptySyncDone{}, 9999999)}
			require.PanicsWithValue(t, "[ processSignalsWorking ] Unknown signal: 9999999", func() {
				worker.processSignalsWorking(badSignal)
			})
			areSlotStatesEqual(&oldSlot, slot, t, false)
		})
	}
}

func Test_processSignalsWorking_PendingPulseSignal(t *testing.T) {

	for _, tt := range testPulseStates {
		t.Run(tt.String(), func(t *testing.T) {
			slot, worker := makeSlotAndWorker(tt, 22)
			pendingSignal := []queue.OutputElement{*queue.NewOutputElement(emptySyncDone{}, PendingPulseSignal)}
			require.Equal(t, 1, worker.processSignalsWorking(pendingSignal))
			require.Equal(t, Suspending, slot.slotState)
		})
	}
}

type emptySyncDone struct{}

func (m emptySyncDone) Done() {}

func Test_processSignalsWorking_ActivatePulseSignal(t *testing.T) {

	for _, tt := range testPulseStates {
		t.Run(tt.String(), func(t *testing.T) {
			slot, worker := makeSlotAndWorker(tt, 22)
			oldSlot := *slot
			activateSignal := []queue.OutputElement{*queue.NewOutputElement(emptySyncDone{}, ActivatePulseSignal)}
			require.Equal(t, 1, worker.processSignalsWorking(activateSignal))

			areSlotStatesEqual(&oldSlot, slot, t, false)
		})
	}
}

func Test_processSignalsWorking_ActivateAndPendingPulseSignals(t *testing.T) {

	for _, tt := range testPulseStates {
		t.Run(tt.String(), func(t *testing.T) {
			slot, worker := makeSlotAndWorker(tt, 22)
			signals := []queue.OutputElement{
				*queue.NewOutputElement(emptySyncDone{}, ActivatePulseSignal),
				*queue.NewOutputElement(emptySyncDone{}, PendingPulseSignal),
			}

			require.Equal(t, 2, worker.processSignalsWorking(signals))
			require.Equal(t, Suspending, slot.slotState)
			inputElements := slot.inputQueue.RemoveAll()
			require.Len(t, inputElements, 1)
			require.Equal(t, ActivatePulseSignal, int(inputElements[0].GetItemType()))
		})
	}
}

// ---- readInputQueueWorking

func Test_readInputQueueWorking_EmptyInputQueue(t *testing.T) {

	for _, tt := range testPulseStates {
		t.Run(tt.String(), func(t *testing.T) {
			slot, worker := makeSlotAndWorker(tt, 22)
			oldSlot := *slot
			require.NoError(t, worker.readInputQueueWorking())

			areSlotStatesEqual(&oldSlot, slot, t, false)
		})
	}
}

func Test_readInputQueueWorking_SignalOnly(t *testing.T) {

	for _, tt := range testPulseStates {
		t.Run(tt.String(), func(t *testing.T) {

			slot, worker := makeSlotAndWorker(tt, 22)
			oldSlot := *slot

			require.NoError(t, slot.inputQueue.PushSignal(ActivatePulseSignal, mockCallback()))
			require.NoError(t, worker.readInputQueueWorking())
			areSlotStatesEqual(&oldSlot, slot, t, false)
		})
	}
}

func Test_readInputQueueWorking_EventOnly(t *testing.T) {

	for _, tt := range testPulseStates {
		t.Run(tt.String(), func(t *testing.T) {
			slot, worker := makeSlotAndWorker(tt, 22)
			oldSlot := *slot
			var payLoad interface{}
			payLoad = 99
			require.NoError(t, slot.inputQueue.SinkPush(payLoad))
			require.NoError(t, worker.readInputQueueWorking())

			areSlotStatesEqual(&oldSlot, slot, t, false)
			el := slot.popElement(ActiveElement)
			require.Equal(t, payLoad, el.payload)
		})
	}
}

func Test_readInputQueueWorking_SignalsAndEvents(t *testing.T) {

	for _, tt := range testPulseStates {
		t.Run(tt.String(), func(t *testing.T) {
			slot, worker := makeSlotAndWorker(tt, 22)
			oldSlot := *slot

			slot.inputQueue.PushSignal(ActivatePulseSignal, mockCallback())

			numElements := 20
			for i := 0; i < numElements; i++ {
				require.NoError(t, slot.inputQueue.SinkPush(i))
			}

			require.NoError(t, worker.readInputQueueWorking())
			areSlotStatesEqual(&oldSlot, slot, t, false)

			for i := 0; i < numElements; i++ {
				el := slot.popElement(ActiveElement)
				require.Equal(t, i, el.payload)
			}
		})
	}
}

// ---- processSignalsSuspending

func Test_processSignalsSuspending_EmptyInput(t *testing.T) {

	for _, tt := range testPulseStates {
		t.Run(tt.String(), func(t *testing.T) {
			slot, worker := makeSlotAndWorker(tt, 22)

			oldSlot := *slot
			require.Equal(t, 0, worker.processSignalsSuspending([]queue.OutputElement{}))
			areSlotStatesEqual(&oldSlot, slot, t, false)
		})
	}
}

func Test_processSignalsSuspending_NonSignals(t *testing.T) {

	for _, tt := range testPulseStates {
		t.Run(tt.String(), func(t *testing.T) {
			slot, worker := makeSlotAndWorker(tt, 22)
			oldSlot := *slot

			nonSignals := []queue.OutputElement{
				*queue.NewOutputElement(1, 0),
				*queue.NewOutputElement(2, 0),
				*queue.NewOutputElement(3, 0),
			}
			require.Equal(t, 0, worker.processSignalsSuspending(nonSignals))

			areSlotStatesEqual(&oldSlot, slot, t, false)
		})
	}
}

func Test_processSignalsSuspending_BadSignal(t *testing.T) {

	for _, tt := range testPulseStates {
		t.Run(tt.String(), func(t *testing.T) {
			slot, worker := makeSlotAndWorker(tt, 22)
			oldSlot := *slot

			badSignal := []queue.OutputElement{*queue.NewOutputElement(1, 9999999)}
			require.PanicsWithValue(t, "[ processSignalsSuspending ] Unknown signal: 9999999", func() {
				worker.processSignalsSuspending(badSignal)
			})
			areSlotStatesEqual(&oldSlot, slot, t, false)
		})
	}
}

func Test_processSignalsSuspending_PendingPulseSignal(t *testing.T) {

	for _, tt := range testPulseStates {
		t.Run(tt.String(), func(t *testing.T) {
			slot, worker := makeSlotAndWorker(tt, 22)
			oldSlot := *slot
			pendingSignal := []queue.OutputElement{*queue.NewOutputElement(1, PendingPulseSignal)}
			require.Equal(t, 1, worker.processSignalsSuspending(pendingSignal))
			areSlotStatesEqual(&oldSlot, slot, t, false)
		})
	}
}

func Test_processSignalsSuspending_ActivatePulseSignal(t *testing.T) {

	for _, tt := range testPulseStates {
		t.Run(tt.String(), func(t *testing.T) {
			slot, worker := makeSlotAndWorker(tt, 22)
			activateSignal := []queue.OutputElement{*queue.NewOutputElement(emptySyncDone{}, ActivatePulseSignal)}
			require.Equal(t, 1, worker.processSignalsSuspending(activateSignal))

			require.Equal(t, Initializing, slot.slotState)
		})
	}
}

// ---- readInputQueueWorking

func Test_readInputQueueSuspending_EmptyInputQueue(t *testing.T) {

	for _, tt := range testPulseStates {
		t.Run(tt.String(), func(t *testing.T) {
			slot, worker := makeSlotAndWorker(tt, 22)
			oldSlot := *slot
			require.NoError(t, worker.readInputQueueSuspending())

			areSlotStatesEqual(&oldSlot, slot, t, false)
		})
	}
}

func Test_readInputQueueSuspending_SignalOnly(t *testing.T) {

	for _, tt := range testPulseStates {
		t.Run(tt.String(), func(t *testing.T) {
			slot, worker := makeSlotAndWorker(tt, 22)
			oldSlot := *slot

			require.NoError(t, slot.inputQueue.PushSignal(PendingPulseSignal, mockCallback()))
			require.NoError(t, worker.readInputQueueSuspending())
			areSlotStatesEqual(&oldSlot, slot, t, false)
		})
	}
}

func Test_readInputQueueSuspending_EventOnly(t *testing.T) {

	tests := []constant.PulseState{constant.Future, constant.Present}

	for _, tt := range tests {
		t.Run(tt.String(), func(t *testing.T) {
			slot, worker := makeSlotAndWorker(tt, 22)
			oldSlot := *slot
			var payLoad interface{}
			payLoad = 99
			require.NoError(t, slot.inputQueue.SinkPush(payLoad))
			require.NoError(t, worker.readInputQueueSuspending())

			areSlotStatesEqual(&oldSlot, slot, t, false)
			el := slot.popElement(ActiveElement)
			require.Equal(t, payLoad, el.payload)
		})
	}
}

func Test_readInputQueueSuspending_EventOnly_Past(t *testing.T) {
	slot, worker := makeSlotAndWorker(constant.Past, 4444)
	var payLoad interface{}
	payLoad = 99
	require.NoError(t, slot.inputQueue.SinkPush(payLoad))
	require.NoError(t, worker.readInputQueueSuspending())

	require.Equal(t, Working, slot.slotState)

	el := slot.popElement(ActiveElement)
	require.Equal(t, payLoad, el.payload)
}

func Test_readInputQueueSuspending_SignalsAndEvents(t *testing.T) {
	tests := []constant.PulseState{constant.Future, constant.Present}

	for _, tt := range tests {
		t.Run(tt.String(), func(t *testing.T) {
			slot, worker := makeSlotAndWorker(tt, 4444)
			oldSlot := *slot

			slot.inputQueue.PushSignal(PendingPulseSignal, mockCallback())

			numElements := 20
			for i := 0; i < numElements; i++ {
				require.NoError(t, slot.inputQueue.SinkPush(i))
			}

			require.NoError(t, worker.readInputQueueSuspending())
			areSlotStatesEqual(&oldSlot, slot, t, false)

			for i := 0; i < numElements; i++ {
				el := slot.popElement(ActiveElement)
				require.Equal(t, i, el.payload)
			}
		})
	}
}

func Test_readInputQueueSuspending_SignalsAndEvents_Past(t *testing.T) {
	slot, worker := makeSlotAndWorker(constant.Past, 44444)
	slot.inputQueue.PushSignal(ActivatePulseSignal, mockCallback())

	numElements := 20
	for i := 0; i < numElements; i++ {
		require.NoError(t, slot.inputQueue.SinkPush(i))
	}

	require.NoError(t, worker.readInputQueueSuspending())

	for i := 0; i < numElements; i++ {
		el := slot.popElement(ActiveElement)
		require.Equal(t, i, el.payload)
	}

	require.Equal(t, Working, slot.slotState)

}

// ---- migrate

var testActivationStatus = []ActivationStatus{ActiveElement, NotActiveElement}

func Test_migrate_EmptyList(t *testing.T) {
	for _, tps := range testPulseStates {
		t.Run(tps.String(), func(t *testing.T) {

			for _, tas := range testActivationStatus {
				t.Run(tas.String(), func(t *testing.T) {

					_, worker := makeSlotAndWorker(tps, 44444)
					require.NoError(t, worker.migrate(tas))
				})
			}
		})
	}
}

func Test_migrate_NoMigrationHandler(t *testing.T) {
	sm := istatemachine.NewStateMachineTypeMock(t)
	sm.GetMigrationHandlerFunc = func(p constant.PulseState, p1 uint32) (r istatemachine.MigrationHandler) {
		return nil
	}

	for _, tps := range testPulseStates {
		t.Run(tps.String(), func(t *testing.T) {
			for _, tas := range testActivationStatus {
				t.Run(tas.String(), func(t *testing.T) {

					slot, worker := makeSlotAndWorker(tps, 44444)
					oldSlot := *slot

					_, err := slot.createElement(sm, 22, queue.OutputElement{})
					require.NoError(t, err)
					numActiveElements := slot.len(tas)
					require.NoError(t, worker.migrate(tas))
					areSlotStatesEqual(&oldSlot, slot, t, false)
					require.Equal(t, numActiveElements, slot.len(tas))
				})
			}
		})
	}
}

// pop element and move it to targetStatus list
func moveLastElementToState(slot *Slot, targetStatus ActivationStatus, t *testing.T) {
	element := slot.popElement(ActiveElement)
	require.NotNil(t, element)
	element.activationStatus = targetStatus
	err := slot.pushElement(element)
	require.NoError(t, err)
}

func Test_migrate_MigrationHandlerOk(t *testing.T) {
	initState := uint32(44)
	migrationState := initState + 1
	initPayLoad := 99
	migrationPayLoad := initPayLoad + 1
	sm := istatemachine.NewStateMachineTypeMock(t)
	sm.GetMigrationHandlerFunc = func(p constant.PulseState, p1 uint32) (r istatemachine.MigrationHandler) {
		return func(element islot.SlotElementHelper) (interface{}, uint32, error) {
			return migrationPayLoad, joinStates(0, migrationState), nil
		}
	}

	for _, tps := range testPulseStates {
		t.Run(tps.String(), func(t *testing.T) {
			for _, tas := range testActivationStatus {
				t.Run(tas.String(), func(t *testing.T) {

					slot, worker := makeSlotAndWorker(tps, 4444)
					event := queue.NewOutputElement(initPayLoad, 0)

					_, err := slot.createElement(sm, initState, *event)
					require.NoError(t, err)

					moveLastElementToState(slot, tas, t)

					require.NoError(t, worker.migrate(tas))
					element := slot.popElement(tas)
					require.NotNil(t, element)
					require.Equal(t, migrationState, element.state)
					require.Equal(t, migrationPayLoad, element.payload)
				})
			}
		})
	}
}

func Test_migrate_MigrationHandler_LastStateOfStateMachine(t *testing.T) {
	sm := istatemachine.NewStateMachineTypeMock(t)
	sm.GetMigrationHandlerFunc = func(p constant.PulseState, p1 uint32) (r istatemachine.MigrationHandler) {
		return func(element islot.SlotElementHelper) (interface{}, uint32, error) {
			return element.GetPayload(), 0, nil
		}
	}

	for _, tps := range testPulseStates {
		t.Run(tps.String(), func(t *testing.T) {
			for _, tas := range testActivationStatus {
				t.Run(tas.String(), func(t *testing.T) {

					slot, worker := makeSlotAndWorker(tps, 444)
					_, err := slot.createElement(sm, 22, queue.OutputElement{})
					require.NoError(t, err)
					oldSlot := *slot

					moveLastElementToState(slot, tas, t)

					numEmptyElements := slot.len(EmptyElement)
					require.NoError(t, worker.migrate(tas))

					areSlotStatesEqual(&oldSlot, slot, t, false)

					element := slot.popElement(tas)
					require.Nil(t, element)
					require.Equal(t, numEmptyElements+1, slot.len(EmptyElement))
				})
			}
		})
	}
}

func Test_migrate_MigrationHandler_Error(t *testing.T) {
	sm := istatemachine.NewStateMachineTypeMock(t)
	sm.GetMigrationHandlerFunc = func(p constant.PulseState, p1 uint32) (r istatemachine.MigrationHandler) {
		return func(element islot.SlotElementHelper) (interface{}, uint32, error) {
			return element.GetPayload(), 0, errors.New("Test Error")
		}
	}

	transitionErrorState := uint32(999)
	transitionErrorPayLoad := 777
	sm.GetTransitionErrorHandlerFunc = func(p constant.PulseState, p1 uint32) (r istatemachine.TransitionErrorHandler) {
		return func(element islot.SlotElementHelper, err error) (interface{}, uint32) {
			return transitionErrorPayLoad, joinStates(0, transitionErrorState)
		}
	}

	for _, tps := range testPulseStates {
		t.Run(tps.String(), func(t *testing.T) {
			for _, tas := range testActivationStatus {
				t.Run(tas.String(), func(t *testing.T) {

					slot, worker := makeSlotAndWorker(tps, 22)
					_, err := slot.createElement(sm, 22, queue.OutputElement{})
					require.NoError(t, err)
					moveLastElementToState(slot, tas, t)

					require.NoError(t, worker.migrate(tas))
					element := slot.popElement(tas)
					require.NotNil(t, element)
					require.Equal(t, transitionErrorState, element.state)
					require.Equal(t, transitionErrorPayLoad, element.payload)
				})
			}
		})
	}
}

// ---- suspending

func Test_suspending_Past(t *testing.T) {
	slot, worker := makeSlotAndWorker(constant.Past, 22)
	removeSlot := false
	slot.removeSlotCallback = func(number core.PulseNumber) {
		removeSlot = true
	}
	oldSlot := *slot

	// to predict infinite loop
	require.NoError(t, slot.inputQueue.PushSignal(ActivatePulseSignal, mockCallback()))

	slot.slotState = Suspending
	worker.suspending()
	areSlotStatesEqual(&oldSlot, slot, t, false)
	require.True(t, removeSlot)
}

func Test_suspending_Present(t *testing.T) {
	slot, worker := makeSlotAndWorker(constant.Present, 22)
	oldSlot := *slot

	// to predict infinite loop
	require.NoError(t, slot.inputQueue.PushSignal(ActivatePulseSignal, mockCallback()))

	slot.slotState = Suspending
	require.Equal(t, 0, worker.nodeState)
	worker.suspending()
	areSlotStatesEqual(&oldSlot, slot, t, true)
	require.Equal(t, 555, worker.nodeState)
}

func Test_suspending_ReadInputQueue(t *testing.T) {
	slot, worker := makeSlotAndWorker(constant.Present, 22)

	// to predict infinite loop
	require.NoError(t, slot.inputQueue.PushSignal(ActivatePulseSignal, mockCallback()))

	require.Equal(t, 0, worker.nodeState)
	slot.slotState = Suspending
	worker.suspending()
	require.Equal(t, constant.Past, slot.pulseState)
}

// ---- working

func Test_working_ChangeStateToSuspending(t *testing.T) {
	for _, tt := range testPulseStates {
		t.Run(tt.String(), func(t *testing.T) {
			slot, worker := makeSlotAndWorker(tt, 22)
			slot.slotState = Working

			require.NoError(t, slot.inputQueue.PushSignal(PendingPulseSignal, mockCallback()))
			require.NoError(t, slot.inputQueue.PushSignal(ActivatePulseSignal, mockCallback()))
			worker.working()

			require.Equal(t, Suspending, slot.slotState)
			require.Equal(t, tt, slot.pulseState)
		})
	}
}

// ---- processingElements

func Test_processingElements_NoElementsInPast(t *testing.T) {
	slot, worker := makeSlotAndWorker(constant.Past, 22)

	worker.processingElements()
	require.Equal(t, Suspending, slot.slotState)
}

func Test_processingElements_AlreadyHasSignal(t *testing.T) {
	slot, worker := makeSlotAndWorker(constant.Present, 22)
	oldSlot := *slot

	require.NoError(t, slot.inputQueue.PushSignal(ActivatePulseSignal, mockCallback()))
	worker.processingElements()

	areSlotStatesEqual(&oldSlot, slot, t, false)
}

func Test_processingElements_OneEvent(t *testing.T) {

	transitionState := uint32(433)
	transitionPayload := 556

	sm := istatemachine.NewStateMachineTypeMock(t)
	sm.GetTransitionHandlerFunc = func(p constant.PulseState, p1 uint32) (r istatemachine.TransitHandler) {
		return func(element islot.SlotElementHelper) (interface{}, uint32, error) {
			return transitionPayload, joinStates(0, transitionState), nil
		}
	}

	for _, tps := range testPulseStates {
		t.Run(tps.String(), func(t *testing.T) {

			slot, worker := makeSlotAndWorker(tps, 22)
			oldSlot := *slot
			_, err := slot.createElement(sm, 22, queue.OutputElement{})
			require.NoError(t, err)

			worker.processingElements()

			areSlotStatesEqual(&oldSlot, slot, t, false)

			element := slot.popElement(ActiveElement)
			require.NotNil(t, element)

			require.Equal(t, transitionState, element.state)
			require.Equal(t, transitionPayload, element.payload)
		})
	}
}

func Test_processingElements_LastStateOfStateMachine(t *testing.T) {
	sm := istatemachine.NewStateMachineTypeMock(t)
	sm.GetTransitionHandlerFunc = func(p constant.PulseState, p1 uint32) (r istatemachine.TransitHandler) {
		return func(element islot.SlotElementHelper) (interface{}, uint32, error) {
			return element.GetPayload(), 0, nil
		}
	}

	for _, tps := range testPulseStates {
		t.Run(tps.String(), func(t *testing.T) {

			slot, worker := makeSlotAndWorker(tps, 22)
			oldSlot := *slot

			_, err := slot.createElement(sm, 22, queue.OutputElement{})
			require.NoError(t, err)

			numEmptyElements := slot.len(EmptyElement)
			worker.processingElements()
			element := slot.popElement(ActiveElement)
			require.Nil(t, element)

			require.Equal(t, numEmptyElements+1, slot.len(EmptyElement))
			areSlotStatesEqual(&oldSlot, slot, t, false)
		})
	}
}

func Test_processingElements_TransitionHandlerError(t *testing.T) {
	sm := istatemachine.NewStateMachineTypeMock(t)
	sm.GetTransitionHandlerFunc = func(p constant.PulseState, p1 uint32) (r istatemachine.TransitHandler) {
		return func(element islot.SlotElementHelper) (interface{}, uint32, error) {
			return nil, 0, errors.New("Test Error")
		}
	}

	transitionErrorState := uint32(999)
	transitionErrorPayLoad := 777

	sm.GetTransitionErrorHandlerFunc = func(p constant.PulseState, p1 uint32) (r istatemachine.TransitionErrorHandler) {
		return func(element islot.SlotElementHelper, err error) (interface{}, uint32) {
			return transitionErrorPayLoad, joinStates(0, transitionErrorState)
		}
	}

	for _, tps := range testPulseStates {
		t.Run(tps.String(), func(t *testing.T) {

			slot, worker := makeSlotAndWorker(tps, 22)

			oldSlot := *slot

			_, err := slot.createElement(sm, 22, queue.OutputElement{})
			require.NoError(t, err)

			worker.processingElements()
			element := slot.popElement(ActiveElement)
			require.NotNil(t, element)

			require.Equal(t, transitionErrorState, element.state)
			require.Equal(t, transitionErrorPayLoad, element.payload)

			areSlotStatesEqual(&oldSlot, slot, t, false)
		})
	}
}

// ---- readResponseQueue

func Test_readResponseQueue_EmptyResponseQueue(t *testing.T) {
	for _, tt := range testPulseStates {
		t.Run(tt.String(), func(t *testing.T) {
			slot, worker := makeSlotAndWorker(tt, 22)
			oldSlot := *slot

			require.Empty(t, worker.postponedResponses)
			require.NoError(t, worker.readResponseQueue())
			require.Empty(t, worker.postponedResponses)
			areSlotStatesEqual(&oldSlot, slot, t, false)
		})
	}
}

func Test_readResponseQueue_OneEvent(t *testing.T) {
	responseState := uint32(446)
	sm := istatemachine.NewStateMachineTypeMock(t)
	sm.GetResponseHandlerFunc = func(p constant.PulseState, p1 uint32) (r istatemachine.AdapterResponseHandler) {
		return func(element islot.SlotElementHelper, response iadapter.IAdapterResponse) (interface{}, uint32, error) {
			return element.GetPayload(), joinStates(0, responseState), nil
		}
	}

	for _, tt := range testPulseStates {
		t.Run(tt.String(), func(t *testing.T) {
			slot, worker := makeSlotAndWorker(tt, 22)
			oldSlot := *slot

			resp := &adapter.AdapterResponse{}
			slot.responseQueue.SinkPush(resp)

			_, err := slot.createElement(sm, 22, queue.OutputElement{})
			require.NoError(t, err)

			require.Empty(t, worker.postponedResponses)
			require.NoError(t, worker.readResponseQueue())
			require.Empty(t, worker.postponedResponses)
			areSlotStatesEqual(&oldSlot, slot, t, false)

			element := slot.popElement(ActiveElement)
			require.NotEmpty(t, element)
			require.Equal(t, responseState, element.state)

		})
	}
}

func Test_readResponseQueue_BadTypeInResponseQueue(t *testing.T) {
	for _, tt := range testPulseStates {
		t.Run(tt.String(), func(t *testing.T) {
			slot, worker := makeSlotAndWorker(tt, 22)

			slot.responseQueue.SinkPush(76576)
			require.PanicsWithValue(t, "[ readResponseQueue ] Bad type in adapter response queue: int", func() {
				worker.readResponseQueue()
			})
		})
	}
}

func Test_readResponseQueue_BadElementIdInResponse(t *testing.T) {
	for _, tt := range testPulseStates {
		t.Run(tt.String(), func(t *testing.T) {
			slot, worker := makeSlotAndWorker(tt, 22)
			oldSlot := *slot
			resp := &adapter.AdapterResponse{}
			slot.responseQueue.SinkPush(resp)

			_, err := slot.createElement(nil, 33, queue.OutputElement{})
			require.NoError(t, err)

			// it changes element id
			{
				element := slot.popElement(ActiveElement)
				require.NotNil(t, element)
				element.setDeleteState()
				require.NoError(t, slot.pushElement(element))
			}

			worker.readResponseQueue()

			areSlotStatesEqual(&oldSlot, slot, t, false)

		})
	}
}

func Test_processingElements_ResponseHandlerError(t *testing.T) {
	sm := istatemachine.NewStateMachineTypeMock(t)
	sm.GetResponseHandlerFunc = func(p constant.PulseState, p1 uint32) (r istatemachine.AdapterResponseHandler) {
		return func(element islot.SlotElementHelper, response iadapter.IAdapterResponse) (interface{}, uint32, error) {
			return element.GetPayload(), 0, errors.New("Test Error")
		}
	}

	responseState := uint32(564)
	responsePayload := uint32(345)
	sm.GetResponseErrorHandlerFunc = func(p constant.PulseState, p1 uint32) (r istatemachine.ResponseErrorHandler) {
		return func(element islot.SlotElementHelper, response iadapter.IAdapterResponse, err error) (interface{}, uint32) {
			return responsePayload, joinStates(0, responseState)
		}
	}

	for _, tt := range testPulseStates {
		t.Run(tt.String(), func(t *testing.T) {
			slot, worker := makeSlotAndWorker(tt, 22)
			oldSlot := *slot
			resp := &adapter.AdapterResponse{}
			slot.responseQueue.SinkPush(resp)

			_, err := slot.createElement(sm, 33, queue.OutputElement{})
			require.NoError(t, err)

			worker.readResponseQueue()
			areSlotStatesEqual(&oldSlot, slot, t, false)

			element := slot.popElement(ActiveElement)
			require.NotNil(t, element)

			require.Equal(t, responseState, element.state)
			require.Equal(t, responsePayload, element.payload)
		})
	}

}

// ---- initializing

func Test_initializing_EmptySlot(t *testing.T) {
	for _, tt := range testPulseStates {
		t.Run(tt.String(), func(t *testing.T) {
			slot, worker := makeSlotAndWorker(tt, 22)
			oldSlot := *slot
			worker.initializing()
			areSlotStatesEqual(&oldSlot, slot, t, false)
		})
	}
}

func Test_initializing_NotEmptySlot(t *testing.T) {
	migrationPayLoad := 555
	migrationState := uint32(99)
	sm := istatemachine.NewStateMachineTypeMock(t)
	sm.GetMigrationHandlerFunc = func(p constant.PulseState, p1 uint32) (r istatemachine.MigrationHandler) {
		return func(element islot.SlotElementHelper) (interface{}, uint32, error) {
			return migrationPayLoad, joinStates(0, migrationState), nil
		}
	}

	for _, tt := range testPulseStates {
		t.Run(tt.String(), func(t *testing.T) {
			slot, worker := makeSlotAndWorker(tt, 22)
			oldSlot := *slot

			element, err := slot.createElement(sm, 777, queue.OutputElement{})
			require.NoError(t, err)
			require.NotNil(t, element)

			worker.initializing()
			areSlotStatesEqual(&oldSlot, slot, t, false)
		})
	}
}

// ---- run

func Test_run(t *testing.T) {
	maxState := uint32(1000)
	sm := istatemachine.NewStateMachineTypeMock(t)
	sm.GetMigrationHandlerFunc = func(p constant.PulseState, state uint32) (r istatemachine.MigrationHandler) {
		return func(element islot.SlotElementHelper) (interface{}, uint32, error) {
			if state > maxState {
				state /= 2
			}
			return element.GetElementID(), joinStates(0, state+1), nil
		}
	}

	sm.GetTransitionHandlerFunc = func(p constant.PulseState, state uint32) (r istatemachine.TransitHandler) {
		return func(element islot.SlotElementHelper) (interface{}, uint32, error) {
			if state > maxState {
				state /= 2
			}
			return element.GetElementID(), joinStates(0, state+1), nil
		}
	}

	sm.GetResponseHandlerFunc = func(p constant.PulseState, state uint32) (r istatemachine.AdapterResponseHandler) {
		return func(element islot.SlotElementHelper, response iadapter.IAdapterResponse) (interface{}, uint32, error) {
			if state > maxState {
				state /= 2
			}
			return element.GetPayload(), joinStates(0, state+1), nil
		}
	}

	for _, tt := range testPulseStates {
		t.Run(tt.String(), func(t *testing.T) {
			slot, worker := makeSlotAndWorker(tt, 22)
			for i := 1; i < 8000; i++ {
				state := uint32(i)
				if state > maxState {
					state /= maxState
				}
				element, err := slot.createElement(sm, uint32(state), queue.OutputElement{})
				require.NoError(t, err)
				require.NotNil(t, element)
			}

			go func() {
				for i := 1; i < 10; i++ {
					resp := adapter.NewAdapterResponse(0, uint32(i), 0, 0)
					slot.responseQueue.SinkPush(resp)
					time.Sleep(time.Millisecond * 50)
				}
			}()

			go func() {
				time.Sleep(time.Millisecond * 400)
				slot.inputQueue.PushSignal(PendingPulseSignal, mockCallback())
			}()

			go func() {
				time.Sleep(time.Millisecond * 600)
				slot.inputQueue.PushSignal(ActivatePulseSignal, mockCallback())
			}()

			go func() {
				time.Sleep(time.Millisecond * 800)
				slot.inputQueue.PushSignal(CancelSignal, mockCallback())
			}()

			worker.run()

		})
	}

}
