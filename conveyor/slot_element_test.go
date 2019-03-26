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

	"github.com/insolar/insolar/conveyor/adapter"
	"github.com/insolar/insolar/conveyor/interfaces/fsm"
	"github.com/insolar/insolar/conveyor/interfaces/slot"
	"github.com/insolar/insolar/conveyor/interfaces/statemachine"
	"github.com/insolar/insolar/insolar"
	"github.com/stretchr/testify/require"
)

func TestSlotElement_setDeleteState(t *testing.T) {
	el := newSlotElement(ActiveElement, nil)
	el.setDeleteState()
	require.Equal(t, el.activationStatus, EmptyElement)
}

func TestSlotElement_Reactivate(t *testing.T) {
	el := newSlotElement(NotActiveElement, nil)
	el.Reactivate()
	require.Equal(t, el.activationStatus, ActiveElement)
}

func TestSlotElement_DeactivateTill_Empty(t *testing.T) {
	el := newSlotElement(ActiveElement, nil)
	require.Panics(t, func() {
		el.DeactivateTill(slot.Empty)
	})
}

func TestSlotElement_DeactivateTill_Tick(t *testing.T) {
	el := newSlotElement(ActiveElement, nil)
	require.Panics(t, func() {
		el.DeactivateTill(slot.Tick)
	})
}

func TestSlotElement_DeactivateTill_SeqHead(t *testing.T) {
	el := newSlotElement(ActiveElement, nil)
	require.Panics(t, func() {
		el.DeactivateTill(slot.SeqHead)
	})
}

func TestSlotElement_DeactivateTill_Response(t *testing.T) {
	el := newSlotElement(ActiveElement, nil)
	el.DeactivateTill(slot.Response)
	require.Equal(t, el.activationStatus, NotActiveElement)
}

func TestSlotElement_update(t *testing.T) {
	el := newSlotElement(ActiveElement, nil)
	testStateID := fsm.StateID(42)
	testPayLoad := 142
	testStateMachine := statemachine.NewStateMachineMock(t)
	require.NotEqual(t, testStateID, el.GetState())
	require.NotEqual(t, testPayLoad, el.GetPayload())
	require.NotEqual(t, testStateMachine, el.stateMachine)

	el.update(testStateID, testPayLoad, testStateMachine)

	require.Equal(t, testStateID, el.GetState())
	require.Equal(t, testPayLoad, el.GetPayload())
	require.Equal(t, testStateMachine, el.stateMachine)
}

func TestSlotElement_SendTask_NoSuchAdapterID(t *testing.T) {
	el := newSlotElement(ActiveElement, nil)
	// make it empty for test
	adapter.Storage = adapter.NewStorage()
	require.PanicsWithValue(t, "[ SendTask ] No such adapter: 142", func() {
		el.SendTask(142, 22, 44)
	})
}

func TestSlotElement_SendTask(t *testing.T) {
	slot := newSlot(44, 44, func(number insolar.PulseNumber) {

	})
	el := newSlotElement(ActiveElement, slot)

	el.SendTask(uint32(adapter.SendResponseAdapterID), 42, 42)
	adapter := adapter.Storage.GetAdapterByID(uint32(adapter.SendResponseAdapterID))
}
