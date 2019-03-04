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
	"github.com/insolar/insolar/core"
	"github.com/stretchr/testify/require"
)

const testRealPulse = core.PulseNumber(1000)
const testPulseDelta = 10
const testUnknownPastPulse = core.PulseNumber(500)
const testUnknownFuturePulse = core.PulseNumber(2000)

func mockQueue(t *testing.T) *queue.IQueueMock {
	qMock := queue.NewIQueueMock(t)
	qMock.SinkPushFunc = func(p interface{}) (r error) {
		return nil
	}
	qMock.SinkPushAllFunc = func(p []interface{}) (r error) {
		return nil
	}
	qMock.PushSignalFunc = func(p uint32, p1 queue.SyncDone) (r error) {
		p1.Done()
		return nil
	}
	return qMock
}

func mockQueueReturnFalse(t *testing.T) *queue.IQueueMock {
	qMock := queue.NewIQueueMock(t)
	qMock.SinkPushFunc = func(p interface{}) (r error) {
		return errors.New("test error")
	}
	qMock.SinkPushAllFunc = func(p []interface{}) (r error) {
		return errors.New("test error")
	}
	qMock.PushSignalFunc = func(p uint32, p1 queue.SyncDone) (r error) {
		return errors.New("test error")
	}
	return qMock
}

func mockSlot(t *testing.T, q *queue.IQueueMock, pulseNumber core.PulseNumber) *Slot {
	slot := &Slot{
		inputQueue:  q,
		pulseNumber: pulseNumber,
	}
	return slot
}

type mockSyncDone struct {
	doneCount int
}

func (s *mockSyncDone) Done() {
	s.doneCount = s.doneCount + 1
}

func mockCallback() queue.SyncDone {
	return &mockSyncDone{doneCount: 0}
}

func testPulseConveyor(t *testing.T, isQueueOk bool) *PulseConveyor {
	var q *queue.IQueueMock
	if isQueueOk {
		q = mockQueue(t)
	} else {
		q = mockQueueReturnFalse(t)
	}
	presentSlot := mockSlot(t, q, testRealPulse)
	futureSlot := mockSlot(t, q, testRealPulse+testPulseDelta)
	slotMap := make(map[core.PulseNumber]*Slot)
	slotMap[testRealPulse] = presentSlot
	slotMap[testRealPulse+testPulseDelta] = futureSlot
	slotMap[core.AntiquePulseNumber] = mockSlot(t, q, core.AntiquePulseNumber)

	return &PulseConveyor{
		state:              Active,
		slotMap:            slotMap,
		futurePulseNumber:  &futureSlot.pulseNumber,
		presentPulseNumber: &presentSlot.pulseNumber,
	}
}

func TestNewPulseConveyor(t *testing.T) {
	c := NewPulseConveyor()
	require.NotNil(t, c)
}

func TestConveyor_GetState(t *testing.T) {
	c := testPulseConveyor(t, true)
	c.state = PreparingPulse

	state := c.GetState()
	require.Equal(t, PreparingPulse, state)
}

var tests = []struct {
	state          State
	expectedResult bool
}{
	{Active, true},
	{PreparingPulse, true},
	{ShuttingDown, false},
	{Inactive, false},
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

func TestConveyor_SinkPush(t *testing.T) {
	c := testPulseConveyor(t, true)
	data := "fancy_data"

	err := c.SinkPush(testRealPulse, data)
	require.NoError(t, err)
	c.slotMap[testRealPulse].inputQueue.(*queue.IQueueMock).SinkPushMock.Expect(data)
}

func TestConveyor_SinkPush_QueueErr(t *testing.T) {
	c := testPulseConveyor(t, false)
	data := "fancy_data"

	err := c.SinkPush(testRealPulse, data)
	require.EqualError(t, err, "[ SinkPush ] can't push to queue: test error")
	c.slotMap[testRealPulse].inputQueue.(*queue.IQueueMock).SinkPushMock.Expect(data)
}

func TestConveyor_SinkPush_AntiqueSlot(t *testing.T) {
	c := testPulseConveyor(t, true)
	data := "fancy_data"

	err := c.SinkPush(testUnknownPastPulse, data)
	require.NoError(t, err)
	c.slotMap[core.AntiquePulseNumber].inputQueue.(*queue.IQueueMock).SinkPushMock.Expect(data)
}

func TestConveyor_SinkPush_UnknownSlot(t *testing.T) {
	c := testPulseConveyor(t, true)
	data := "fancy_data"

	err := c.SinkPush(testUnknownFuturePulse, data)
	require.EqualError(t, err, fmt.Sprintf("[ SinkPush ] can't get slot by pulse number %d", testUnknownFuturePulse))
}

func TestConveyor_SinkPush_NotOperational(t *testing.T) {
	c := testPulseConveyor(t, true)
	c.state = Inactive
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
	c.slotMap[testRealPulse].inputQueue.(*queue.IQueueMock).SinkPushMock.Expect(data)
}

func TestConveyor_SinkPushAll_QueueErr(t *testing.T) {
	c := testPulseConveyor(t, false)
	data1 := "fancy_data_1"
	data2 := "fancy_data_2"
	data := []interface{}{data1, data2}

	err := c.SinkPushAll(testRealPulse, data)
	require.EqualError(t, err, "[ SinkPushAll ] can't push to queue: test error")
	c.slotMap[testRealPulse].inputQueue.(*queue.IQueueMock).SinkPushMock.Expect(data)
}

func TestConveyor_SinkPushAll_AntiqueSlot(t *testing.T) {
	c := testPulseConveyor(t, true)
	data1 := "fancy_data_1"
	data2 := "fancy_data_2"
	data := []interface{}{data1, data2}

	err := c.SinkPushAll(testUnknownPastPulse, data)
	require.NoError(t, err)
	c.slotMap[core.AntiquePulseNumber].inputQueue.(*queue.IQueueMock).SinkPushMock.Expect(data)
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
	c.state = Inactive
	data1 := "fancy_data_1"
	data2 := "fancy_data_2"
	data := []interface{}{data1, data2}

	err := c.SinkPushAll(testUnknownFuturePulse, data)
	require.EqualError(t, err, "[ SinkPushAll ] conveyor is not operational now")
}

func TestConveyor_PreparePulse(t *testing.T) {
	c := testPulseConveyor(t, true)
	c.futurePulseData = nil
	pulse := core.Pulse{PulseNumber: testRealPulse + testPulseDelta}
	callback := mockCallback()

	err := c.PreparePulse(pulse, callback)
	require.NoError(t, err)
	require.NotNil(t, c.futurePulseData)
	require.Equal(t, PreparingPulse, c.state)
}

func TestConveyor_PreparePulse_ForFirstTime(t *testing.T) {
	c := testPulseConveyor(t, true)
	c.futurePulseNumber = nil
	pulse := core.Pulse{PulseNumber: testRealPulse + testPulseDelta}
	callback := mockCallback()

	err := c.PreparePulse(pulse, callback)
	require.NoError(t, err)
	require.NotNil(t, c.futurePulseData)
	require.NotNil(t, c.futurePulseNumber)
	require.Equal(t, PreparingPulse, c.state)
}

func TestConveyor_PreparePulse_NotOperational(t *testing.T) {
	c := testPulseConveyor(t, true)
	pulse := core.Pulse{PulseNumber: testRealPulse + testPulseDelta}
	c.state = Inactive
	callback := mockCallback()

	err := c.PreparePulse(pulse, callback)

	require.EqualError(t, err, "[ PreparePulse ] conveyor is not operational now")
	require.Nil(t, c.futurePulseData)
	require.Equal(t, Inactive, c.state)
}

func TestConveyor_PreparePulse_AlreadyDone(t *testing.T) {
	c := testPulseConveyor(t, true)
	pulse := core.Pulse{PulseNumber: testRealPulse + testPulseDelta}
	c.futurePulseData = &pulse
	c.state = PreparingPulse
	callback := mockCallback()

	err := c.PreparePulse(pulse, callback)

	require.EqualError(t, err, "[ PreparePulse ] preparation was already done")
	require.Equal(t, PreparingPulse, c.state)
}

func TestConveyor_PreparePulse_NotFuture(t *testing.T) {
	c := testPulseConveyor(t, true)
	pulse := core.Pulse{PulseNumber: testRealPulse + testPulseDelta + 10}
	callback := mockCallback()

	err := c.PreparePulse(pulse, callback)
	require.EqualError(t, err, "[ PreparePulse ] received future pulse is different from expected")
	require.Nil(t, c.futurePulseData)
	require.Equal(t, Active, c.state)
}

func TestConveyor_PreparePulse_PushSignalPresentPanic(t *testing.T) {
	c := testPulseConveyor(t, false)
	c.futurePulseNumber = nil
	pulse := core.Pulse{PulseNumber: testRealPulse + testPulseDelta}
	callback := mockCallback()
	oldState := c.state

	panicValue := fmt.Sprintf("[ PreparePulse ] can't send signal to present slot (for pulse %d), error - test error", c.presentPulseNumber)
	require.PanicsWithValue(t, panicValue, func() { c.PreparePulse(pulse, callback) })
	require.Nil(t, c.futurePulseData)
	require.Equal(t, oldState, c.state)
}

func TestConveyor_PreparePulse_PushSignalFuturePanic(t *testing.T) {
	c := testPulseConveyor(t, false)
	pulse := core.Pulse{PulseNumber: testRealPulse + testPulseDelta}
	callback := mockCallback()
	oldState := c.state

	panicValue := fmt.Sprintf("[ PreparePulse ] can't send signal to future slot (for pulse %d), error - test error", c.futurePulseNumber)
	require.PanicsWithValue(t, panicValue, func() { c.PreparePulse(pulse, callback) })
	require.Nil(t, c.futurePulseData)
	require.Equal(t, oldState, c.state)
}

func TestConveyor_ActivatePulse(t *testing.T) {
	c := testPulseConveyor(t, true)
	pulse := core.Pulse{PulseNumber: testRealPulse + testPulseDelta}
	c.futurePulseData = &pulse
	newFutureSlot := mockSlot(t, mockQueue(t), pulse.NextPulseNumber)
	c.slotMap[pulse.NextPulseNumber] = newFutureSlot
	c.newFuturePulseNumber = &pulse.NextPulseNumber
	c.state = PreparingPulse

	err := c.ActivatePulse()

	require.NoError(t, err)
	require.Nil(t, c.futurePulseData)
	require.Equal(t, Active, c.state)
}

func TestConveyor_ActivatePulse_NotOperational(t *testing.T) {
	c := testPulseConveyor(t, true)
	pulse := core.Pulse{PulseNumber: testRealPulse + testPulseDelta}
	c.futurePulseData = &pulse
	c.state = Inactive

	err := c.ActivatePulse()

	require.EqualError(t, err, "[ ActivatePulse ] conveyor is not operational now")
	require.Equal(t, &pulse, c.futurePulseData)
	require.Equal(t, Inactive, c.state)
}

func TestConveyor_ActivatePulse_NoPrepare(t *testing.T) {
	c := testPulseConveyor(t, true)
	c.futurePulseData = nil

	err := c.ActivatePulse()

	require.EqualError(t, err, "[ ActivatePulse ] preparation missing")
	require.Equal(t, Active, c.state)
}

func TestConveyor_ActivatePulse_PushSignalErr(t *testing.T) {
	c := testPulseConveyor(t, false)
	pulse := core.Pulse{PulseNumber: testRealPulse + testPulseDelta}
	c.futurePulseData = &pulse
	newFutureSlot := NewSlot(Unallocated, pulse.NextPulseNumber)
	c.slotMap[pulse.NextPulseNumber] = newFutureSlot
	c.newFuturePulseNumber = &pulse.NextPulseNumber
	c.state = PreparingPulse

	panicValue := fmt.Sprintf("[ ActivatePulse ] can't send signal to present slot (for pulse %d), error - test error", c.futurePulseNumber)
	require.PanicsWithValue(t, panicValue, func() { c.ActivatePulse() })
	require.NotNil(t, c.futurePulseData)
	require.Equal(t, PreparingPulse, c.state)
}

func TestConveyor_ActivatePreparePulse(t *testing.T) {
	c := testPulseConveyor(t, true)
	c.futurePulseData = nil
	pulse := core.Pulse{PulseNumber: testRealPulse + testPulseDelta}
	callback := mockCallback()

	err := c.PreparePulse(pulse, callback)
	require.NoError(t, err)

	// mock queue in new slot, created in prepare func
	c.slotMap[pulse.NextPulseNumber].inputQueue = mockQueue(t)

	err = c.ActivatePulse()
	require.NoError(t, err)
}
