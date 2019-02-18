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

package conveyer

import (
	"errors"
	"fmt"
	"testing"

	"github.com/insolar/insolar/conveyer/queue"
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
	return qMock
}

func mockSlot(t *testing.T, q *queue.IQueueMock, pulseNumber core.PulseNumber) *Slot {
	slot := &Slot{
		inputQueue:  q,
		pulseNumber: pulseNumber,
	}
	return slot
}

func testPulseConveyer(t *testing.T, isQueueOk bool) *PulseConveyer {
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
	slotMap[AntiqueSlotPulse] = mockSlot(t, q, AntiqueSlotPulse)

	return &PulseConveyer{
		state:              Active,
		slotMap:            slotMap,
		futurePulseNumber:  &futureSlot.pulseNumber,
		presentPulseNumber: &presentSlot.pulseNumber,
	}
}

func TestNewPulseConveyer(t *testing.T) {
	c := NewPulseConveyer()
	require.NotNil(t, c)
}

func TestNewSlot(t *testing.T) {
	s := NewSlot(Future, testRealPulse)
	require.NotNil(t, s)
	require.Equal(t, Future, s.pulseState)
	require.Empty(t, s.inputQueue.RemoveAll())
}

func TestConveyer_GetState(t *testing.T) {
	c := testPulseConveyer(t, true)
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

func TestConveyer_IsOperational(t *testing.T) {
	c := testPulseConveyer(t, true)
	for _, tt := range tests {
		t.Run(tt.state.String(), func(t *testing.T) {
			c.state = tt.state
			res := c.IsOperational()
			require.Equal(t, tt.expectedResult, res)
		})
	}
}

func TestConveyer_SinkPush(t *testing.T) {
	c := testPulseConveyer(t, true)
	data := "fancy_data"

	err := c.SinkPush(testRealPulse, data)
	require.NoError(t, err)
	c.slotMap[testRealPulse].inputQueue.(*queue.IQueueMock).SinkPushMock.Expect(data)
}

func TestConveyer_SinkPush_QueueErr(t *testing.T) {
	c := testPulseConveyer(t, false)
	data := "fancy_data"

	err := c.SinkPush(testRealPulse, data)
	require.EqualError(t, err, "[ SinkPush ] can't push to queue: test error")
	c.slotMap[testRealPulse].inputQueue.(*queue.IQueueMock).SinkPushMock.Expect(data)
}

func TestConveyer_SinkPush_AntiqueSlot(t *testing.T) {
	c := testPulseConveyer(t, true)
	data := "fancy_data"

	err := c.SinkPush(testUnknownPastPulse, data)
	require.NoError(t, err)
	c.slotMap[AntiqueSlotPulse].inputQueue.(*queue.IQueueMock).SinkPushMock.Expect(data)
}

func TestConveyer_SinkPush_UnknownSlot(t *testing.T) {
	c := testPulseConveyer(t, true)
	data := "fancy_data"

	err := c.SinkPush(testUnknownFuturePulse, data)
	require.EqualError(t, err, fmt.Sprintf("[ SinkPush ] can't get slot by pulse number %d", testUnknownFuturePulse))
}

func TestConveyer_SinkPush_NotOperational(t *testing.T) {
	c := testPulseConveyer(t, true)
	c.state = Inactive
	data := "fancy_data"

	err := c.SinkPush(testUnknownFuturePulse, data)
	fmt.Println(err.Error())
	require.EqualError(t, err, "[ SinkPush ] conveyer is not operational now")
}

func TestConveyer_SinkPushAll(t *testing.T) {
	c := testPulseConveyer(t, true)
	data1 := "fancy_data_1"
	data2 := "fancy_data_2"
	data := []interface{}{data1, data2}

	err := c.SinkPushAll(testRealPulse, data)
	require.NoError(t, err)
	c.slotMap[testRealPulse].inputQueue.(*queue.IQueueMock).SinkPushMock.Expect(data)
}

func TestConveyer_SinkPushAll_QueueErr(t *testing.T) {
	c := testPulseConveyer(t, false)
	data1 := "fancy_data_1"
	data2 := "fancy_data_2"
	data := []interface{}{data1, data2}

	err := c.SinkPushAll(testRealPulse, data)
	require.EqualError(t, err, "[ SinkPushAll ] can't push to queue: test error")
	c.slotMap[testRealPulse].inputQueue.(*queue.IQueueMock).SinkPushMock.Expect(data)
}

func TestConveyer_SinkPushAll_AntiqueSlot(t *testing.T) {
	c := testPulseConveyer(t, true)
	data1 := "fancy_data_1"
	data2 := "fancy_data_2"
	data := []interface{}{data1, data2}

	err := c.SinkPushAll(testUnknownPastPulse, data)
	require.NoError(t, err)
	c.slotMap[AntiqueSlotPulse].inputQueue.(*queue.IQueueMock).SinkPushMock.Expect(data)
}

func TestConveyer_SinkPushAll_UnknownSlot(t *testing.T) {
	c := testPulseConveyer(t, true)
	data1 := "fancy_data_1"
	data2 := "fancy_data_2"
	data := []interface{}{data1, data2}

	err := c.SinkPushAll(testUnknownFuturePulse, data)
	require.EqualError(t, err, fmt.Sprintf("[ SinkPushAll ] can't get slot by pulse number %d", testUnknownFuturePulse))
}

func TestConveyer_SinkPushAll_NotOperational(t *testing.T) {
	c := testPulseConveyer(t, true)
	c.state = Inactive
	data1 := "fancy_data_1"
	data2 := "fancy_data_2"
	data := []interface{}{data1, data2}

	err := c.SinkPushAll(testUnknownFuturePulse, data)
	require.EqualError(t, err, "[ SinkPushAll ] conveyer is not operational now")
}

func TestConveyer_PreparePulse(t *testing.T) {
	c := testPulseConveyer(t, true)
	c.futurePulseNumber = nil
	pulse := core.Pulse{PulseNumber: testRealPulse + testPulseDelta}

	err := c.PreparePulse(pulse)
	require.NoError(t, err)
	require.NotNil(t, c.futurePulseData)
	require.Equal(t, PreparingPulse, c.state)
}

func TestConveyer_PreparePulse_NotOperational(t *testing.T) {
	c := testPulseConveyer(t, true)
	pulse := core.Pulse{PulseNumber: testRealPulse + testPulseDelta}
	c.state = Inactive

	err := c.PreparePulse(pulse)

	require.EqualError(t, err, "[ PreparePulse ] conveyer is not operational now")
	require.Nil(t, c.futurePulseData)
	require.Equal(t, Inactive, c.state)
}

func TestConveyer_PreparePulse_AlreadyDone(t *testing.T) {
	c := testPulseConveyer(t, true)
	pulse := core.Pulse{PulseNumber: testRealPulse + testPulseDelta}
	c.futurePulseData = &pulse
	c.state = PreparingPulse

	err := c.PreparePulse(pulse)

	require.EqualError(t, err, "[ PreparePulse ] preparation was already done")
	require.Equal(t, PreparingPulse, c.state)
}

func TestConveyer_PreparePulse_NotFuture(t *testing.T) {
	c := testPulseConveyer(t, true)
	pulse := core.Pulse{PulseNumber: testRealPulse + testPulseDelta + 10}

	err := c.PreparePulse(pulse)
	require.EqualError(t, err, "[ PreparePulse ] received future pulse is different from expected")
	require.Nil(t, c.futurePulseData)
	require.Equal(t, Active, c.state)
}

func TestConveyer_ActivatePulse(t *testing.T) {
	c := testPulseConveyer(t, true)
	pulse := core.Pulse{PulseNumber: testRealPulse + testPulseDelta}
	c.futurePulseData = &pulse
	c.state = PreparingPulse

	err := c.ActivatePulse()

	require.NoError(t, err)
	require.Nil(t, c.futurePulseData)
	require.Equal(t, Active, c.state)
}

func TestConveyer_ActivatePulse_NotOperational(t *testing.T) {
	c := testPulseConveyer(t, true)
	pulse := core.Pulse{PulseNumber: testRealPulse + testPulseDelta}
	c.futurePulseData = &pulse
	c.state = Inactive

	err := c.ActivatePulse()

	require.EqualError(t, err, "[ ActivatePulse ] conveyer is not operational now")
	require.Equal(t, &pulse, c.futurePulseData)
	require.Equal(t, Inactive, c.state)
}

func TestConveyer_ActivatePulse_NoPrepare(t *testing.T) {
	c := testPulseConveyer(t, true)
	c.futurePulseData = nil

	err := c.ActivatePulse()

	require.EqualError(t, err, "[ ActivatePulse ] preparation missing")
	require.Equal(t, Active, c.state)
}
