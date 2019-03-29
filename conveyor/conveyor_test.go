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
	"errors"
	"fmt"
	"testing"

	"github.com/insolar/insolar/conveyor/queue"
	"github.com/insolar/insolar/insolar"
	"github.com/stretchr/testify/require"
)

const testRealPulse = insolar.PulseNumber(1000)
const testPulseDelta = 10
const testUnknownPastPulse = insolar.PulseNumber(500)
const testUnknownFuturePulse = insolar.PulseNumber(2000)

func mockSlot(t *testing.T, isQueueOk bool) TaskPusher {
	slot := NewTaskPusherMock(t)
	if isQueueOk {
		slot.SinkPushFunc = func(p interface{}) error {
			return nil
		}

		slot.SinkPushAllFunc = func(p []interface{}) error {
			return nil
		}

		slot.PushSignalFunc = func(signalType uint32, callback queue.SyncDone) error {
			callback.SetResult(333)
			return nil
		}
	} else {
		slot.SinkPushFunc = func(p interface{}) error {
			return errors.New("test error")
		}

		slot.SinkPushAllFunc = func(p []interface{}) error {
			return errors.New("test error")
		}

		slot.PushSignalFunc = func(signalType uint32, callback queue.SyncDone) error {
			return errors.New("test error")
		}
	}
	return slot
}

type mockSyncDone struct {
	waiter    chan interface{}
	doneCount int
}

func (s *mockSyncDone) GetResult() []byte {
	result := <-s.waiter
	hash, ok := result.([]byte)
	if !ok {
		return []byte{}
	}
	return hash
}

func (s *mockSyncDone) SetResult(result interface{}) {
	s.waiter <- result
	s.doneCount += 1
}

func mockCallback() queue.SyncDone {
	return &mockSyncDone{
		doneCount: 0,
		waiter:    make(chan interface{}, 3),
	}
}

func testPulseConveyor(t *testing.T, isQueueOk bool) *PulseConveyor {
	presentSlot := mockSlot(t, isQueueOk)
	futureSlot := mockSlot(t, isQueueOk)

	presentPulse := testRealPulse
	futurePulse := testRealPulse + testPulseDelta

	slotMap := make(map[insolar.PulseNumber]TaskPusher)
	slotMap[presentPulse] = presentSlot
	slotMap[futurePulse] = futureSlot
	slotMap[insolar.AntiquePulseNumber] = mockSlot(t, isQueueOk)

	return &PulseConveyor{
		state:              insolar.ConveyorActive,
		slotMap:            slotMap,
		futurePulseNumber:  &futurePulse,
		presentPulseNumber: &presentPulse,
	}
}

func TestNewPulseConveyor(t *testing.T) {
	c, err := NewPulseConveyor()
	require.NotNil(t, c)
	require.NoError(t, err)
}

func TestConveyor_GetState(t *testing.T) {
	c := testPulseConveyor(t, true)
	c.state = insolar.ConveyorPreparingPulse

	state := c.GetState()
	require.Equal(t, insolar.ConveyorPreparingPulse, state)
}

var tests = []struct {
	state          insolar.ConveyorState
	expectedResult bool
}{
	{insolar.ConveyorActive, true},
	{insolar.ConveyorPreparingPulse, true},
	{insolar.ConveyorShuttingDown, false},
	{insolar.ConveyorInactive, false},
}

func TestConveyor_IsOperational(t *testing.T) {
	c := testPulseConveyor(t, true)
	for _, tt := range tests {
		t.Run(tt.state.String(), func(t *testing.T) {
			c.state = tt.state
			res := c.IsOperational()
			require.Equal(t, tt.expectedResult, res)
		})
	}
}

func TestConveyor_InitiateShutdown(t *testing.T) {
	c := testPulseConveyor(t, true)

	c.InitiateShutdown(true)
	require.Equal(t, insolar.ConveyorShuttingDown, c.state)
}

func TestConveyor_SinkPush(t *testing.T) {
	c := testPulseConveyor(t, true)
	data := "fancy_data"

	err := c.SinkPush(testRealPulse, data)
	require.NoError(t, err)
}

func TestConveyor_SinkPush_QueueErr(t *testing.T) {
	c := testPulseConveyor(t, false)
	data := "fancy_data"

	err := c.SinkPush(testRealPulse, data)
	require.EqualError(t, err, "[ SinkPush ] can't push to queue: test error")
}

func TestConveyor_SinkPush_AntiqueSlot(t *testing.T) {
	c := testPulseConveyor(t, true)
	data := "fancy_data"

	err := c.SinkPush(testUnknownPastPulse, data)
	require.NoError(t, err)
}

func TestConveyor_SinkPush_UnknownSlot(t *testing.T) {
	c := testPulseConveyor(t, true)
	data := "fancy_data"

	err := c.SinkPush(testUnknownFuturePulse, data)
	require.EqualError(t, err, fmt.Sprintf("[ SinkPush ] can't get slot by pulse number %d", testUnknownFuturePulse))
}

func TestConveyor_SinkPush_NotOperational(t *testing.T) {
	c := testPulseConveyor(t, true)
	c.state = insolar.ConveyorInactive
	data := "fancy_data"

	err := c.SinkPush(testUnknownFuturePulse, data)
	fmt.Println(err.Error())
	require.EqualError(t, err, "[ SinkPush ] conveyor is not operational now")
}

func TestConveyor_SinkPushAll(t *testing.T) {
	c := testPulseConveyor(t, true)
	data1 := "fancy_data_1"
	data2 := "fancy_data_2"
	data := []interface{}{data1, data2}

	err := c.SinkPushAll(testRealPulse, data)
	require.NoError(t, err)
}

func TestConveyor_SinkPushAll_QueueErr(t *testing.T) {
	c := testPulseConveyor(t, false)
	data1 := "fancy_data_1"
	data2 := "fancy_data_2"
	data := []interface{}{data1, data2}

	err := c.SinkPushAll(testRealPulse, data)
	require.EqualError(t, err, "[ SinkPushAll ] can't push to queue: test error")
}

func TestConveyor_SinkPushAll_AntiqueSlot(t *testing.T) {
	c := testPulseConveyor(t, true)
	data1 := "fancy_data_1"
	data2 := "fancy_data_2"
	data := []interface{}{data1, data2}

	err := c.SinkPushAll(testUnknownPastPulse, data)
	require.NoError(t, err)
}

func TestConveyor_SinkPushAll_UnknownSlot(t *testing.T) {
	c := testPulseConveyor(t, true)
	data1 := "fancy_data_1"
	data2 := "fancy_data_2"
	data := []interface{}{data1, data2}

	err := c.SinkPushAll(testUnknownFuturePulse, data)
	require.EqualError(t, err, fmt.Sprintf("[ SinkPushAll ] can't get slot by pulse number %d", testUnknownFuturePulse))
}

func TestConveyor_SinkPushAll_NotOperational(t *testing.T) {
	c := testPulseConveyor(t, true)
	c.state = insolar.ConveyorInactive
	data1 := "fancy_data_1"
	data2 := "fancy_data_2"
	data := []interface{}{data1, data2}

	err := c.SinkPushAll(testUnknownFuturePulse, data)
	require.EqualError(t, err, "[ SinkPushAll ] conveyor is not operational now")
}

func TestConveyor_PreparePulse(t *testing.T) {
	c := testPulseConveyor(t, true)
	c.futurePulseData = nil
	pulse := insolar.Pulse{PulseNumber: testRealPulse + testPulseDelta}
	callback := mockCallback()

	err := c.PreparePulse(pulse, callback)
	require.NoError(t, err)
	require.NotNil(t, c.futurePulseData)
	require.Equal(t, insolar.ConveyorPreparingPulse, c.state)
}

func TestConveyor_PreparePulse_ForFirstTime(t *testing.T) {
	c := testPulseConveyor(t, true)
	c.futurePulseNumber = nil
	pulse := insolar.Pulse{PulseNumber: testRealPulse + testPulseDelta}
	callback := mockCallback()

	err := c.PreparePulse(pulse, callback)
	require.NoError(t, err)
	require.NotNil(t, c.futurePulseData)
	require.NotNil(t, c.futurePulseNumber)
	require.Equal(t, insolar.ConveyorPreparingPulse, c.state)
}

func TestConveyor_PreparePulse_ShutDown(t *testing.T) {
	c := testPulseConveyor(t, true)
	pulse := insolar.Pulse{PulseNumber: testRealPulse + testPulseDelta}
	c.state = insolar.ConveyorShuttingDown
	callback := mockCallback()

	err := c.PreparePulse(pulse, callback)

	require.EqualError(t, err, "[ PreparePulse ] conveyor is shut down")
	require.Nil(t, c.futurePulseData)
	require.Equal(t, insolar.ConveyorShuttingDown, c.state)
}

func TestConveyor_PreparePulse_AlreadyDone(t *testing.T) {
	c := testPulseConveyor(t, true)
	pulse := insolar.Pulse{PulseNumber: testRealPulse + testPulseDelta}
	c.futurePulseData = &pulse
	c.state = insolar.ConveyorPreparingPulse
	callback := mockCallback()

	err := c.PreparePulse(pulse, callback)

	require.EqualError(t, err, "[ PreparePulse ] preparation was already done")
	require.Equal(t, insolar.ConveyorPreparingPulse, c.state)
}

func TestConveyor_PreparePulse_NotFuture(t *testing.T) {
	c := testPulseConveyor(t, true)
	pulse := insolar.Pulse{PulseNumber: testRealPulse + testPulseDelta + 10}
	callback := mockCallback()

	err := c.PreparePulse(pulse, callback)
	require.EqualError(t, err, "[ PreparePulse ] received future pulse is different from expected")
	require.Nil(t, c.futurePulseData)
	require.Equal(t, insolar.ConveyorActive, c.state)
}

func TestConveyor_PreparePulse_PushSignalPresentPanic(t *testing.T) {
	c := testPulseConveyor(t, false)
	c.futurePulseNumber = nil
	pulse := insolar.Pulse{PulseNumber: testRealPulse + testPulseDelta}
	callback := mockCallback()
	oldState := c.state

	panicValue := fmt.Sprintf("[ PreparePulse ] can't send signal to present slot (for pulse %d), error - test error", c.presentPulseNumber)
	require.PanicsWithValue(t, panicValue, func() { c.PreparePulse(pulse, callback) })
	require.Nil(t, c.futurePulseData)
	require.Equal(t, oldState, c.state)
}

func TestConveyor_PreparePulse_PushSignalFuturePanic(t *testing.T) {
	c := testPulseConveyor(t, false)
	pulse := insolar.Pulse{PulseNumber: testRealPulse + testPulseDelta}
	callback := mockCallback()
	oldState := c.state

	panicValue := fmt.Sprintf("[ PreparePulse ] can't send signal to future slot (for pulse %d), error - test error", c.futurePulseNumber)
	require.PanicsWithValue(t, panicValue, func() { c.PreparePulse(pulse, callback) })
	require.Nil(t, c.futurePulseData)
	require.Equal(t, oldState, c.state)
}

func TestConveyor_ActivatePulse(t *testing.T) {
	c := testPulseConveyor(t, true)
	pulse := insolar.Pulse{PulseNumber: testRealPulse + testPulseDelta}
	c.futurePulseData = &pulse
	newFutureSlot := mockSlot(t, true)
	c.slotMap[pulse.NextPulseNumber] = newFutureSlot
	c.state = insolar.ConveyorPreparingPulse

	err := c.ActivatePulse()

	require.NoError(t, err)
	require.Nil(t, c.futurePulseData)
	require.Equal(t, insolar.ConveyorActive, c.state)
}

func TestConveyor_ActivatePulse_ShutDown(t *testing.T) {
	c := testPulseConveyor(t, true)
	pulse := insolar.Pulse{PulseNumber: testRealPulse + testPulseDelta}
	c.futurePulseData = &pulse
	c.state = insolar.ConveyorShuttingDown

	err := c.ActivatePulse()

	require.EqualError(t, err, "[ ActivatePulse ] conveyor is shut down")
	require.Equal(t, &pulse, c.futurePulseData)
	require.Equal(t, insolar.ConveyorShuttingDown, c.state)
}

func TestConveyor_ActivatePulse_NoPrepare(t *testing.T) {
	c := testPulseConveyor(t, true)
	c.futurePulseData = nil

	err := c.ActivatePulse()

	require.EqualError(t, err, "[ ActivatePulse ] preparation missing")
	require.Equal(t, insolar.ConveyorActive, c.state)
}

func TestConveyor_ActivatePulse_PushSignalErr(t *testing.T) {
	c := testPulseConveyor(t, false)
	pulse := insolar.Pulse{PulseNumber: testRealPulse + testPulseDelta}
	c.futurePulseData = &pulse
	newFutureSlot := mockSlot(t, false)
	c.slotMap[pulse.NextPulseNumber] = newFutureSlot
	c.state = insolar.ConveyorPreparingPulse

	panicValue := fmt.Sprintf("[ ActivatePulse ] can't send signal to future slot (for pulse %d), error - test error", c.futurePulseNumber)
	require.PanicsWithValue(t, panicValue, func() { c.ActivatePulse() })
	require.NotNil(t, c.futurePulseData)
	require.Equal(t, insolar.ConveyorPreparingPulse, c.state)
}

func TestConveyor_ActivatePreparePulse(t *testing.T) {
	c := testPulseConveyor(t, true)
	c.futurePulseData = nil
	pulse := insolar.Pulse{PulseNumber: testRealPulse + testPulseDelta}
	callback := mockCallback()

	err := c.PreparePulse(pulse, callback)
	require.NoError(t, err)

	err = c.ActivatePulse()
	require.NoError(t, err)
}
