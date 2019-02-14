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
	"context"
	"fmt"
	"testing"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/testutils"
	"github.com/stretchr/testify/require"
)

const testRealPulse = core.PulseNumber(1000)

func mockQueue(t *testing.T) *NonBlockingQueueMock {
	qMock := NewNonBlockingQueueMock(t)
	qMock.SinkPushFunc = func(p interface{}) (r bool) {
		return true
	}
	return qMock
}

func mockQueueReturnFalse(t *testing.T) *NonBlockingQueueMock {
	qMock := NewNonBlockingQueueMock(t)
	qMock.SinkPushFunc = func(p interface{}) (r bool) {
		return false
	}
	return qMock
}

func mockSlot(t *testing.T, q *NonBlockingQueueMock) Slot {
	slot := Slot{
		inputQueue: q,
		pulseState: Present,
	}
	return slot
}

func testPulseConveyer(t *testing.T, isQueueOk bool) *PulseConveyer {
	var q *NonBlockingQueueMock
	if isQueueOk {
		q = mockQueue(t)
	} else {
		q = mockQueueReturnFalse(t)
	}
	slot := mockSlot(t, q)
	slotMap := make(map[core.PulseNumber]Slot)
	slotMap[testRealPulse] = slot

	c := &PulseConveyer{
		slotMap: slotMap,
	}
	return c
}

func TestNewPulseConveyer(t *testing.T) {
	c := NewPulseConveyer()
	require.NotNil(t, c)
}

func TestNewSlot(t *testing.T) {
	s := NewSlot(Future)
	require.NotNil(t, s)
	require.Equal(t, Future, s.pulseState)
	require.Empty(t, s.inputQueue.RemoveAll())
}

func TestConveyer_SinkPush(t *testing.T) {
	c := testPulseConveyer(t, true)
	data := "fancy_data"

	ok := c.sinkPush(testRealPulse, data)
	require.True(t, ok)
	c.slotMap[testRealPulse].inputQueue.(*NonBlockingQueueMock).SinkPushMock.Expect(data)
}

func TestConveyer_SinkPush_Err(t *testing.T) {
	c := testPulseConveyer(t, false)
	data := "fancy_data"

	ok := c.sinkPush(testRealPulse, data)
	require.False(t, ok)
	c.slotMap[testRealPulse].inputQueue.(*NonBlockingQueueMock).SinkPushMock.Expect(data)
}

func TestConveyer_SinkPushAll(t *testing.T) {
	c := testPulseConveyer(t, true)
	data1 := "fancy_data_1"
	data2 := "fancy_data_2"
	data := []interface{}{data1, data2}

	ok := c.sinkPushAll(testRealPulse, data)
	require.True(t, ok)
	c.slotMap[testRealPulse].inputQueue.(*NonBlockingQueueMock).SinkPushMock.Expect(data)
}

func TestConveyer_SinkPushAll_Err(t *testing.T) {
	c := testPulseConveyer(t, false)
	data1 := "fancy_data_1"
	data2 := "fancy_data_2"
	data := []interface{}{data1, data2}

	ok := c.sinkPushAll(testRealPulse, data)
	require.False(t, ok)
	c.slotMap[testRealPulse].inputQueue.(*NonBlockingQueueMock).SinkPushMock.Expect(data)
}

func TestConveyer_Start(t *testing.T) {
	ctx := context.Background()
	c := testPulseConveyer(t, false)
	ps := testutils.NewPulseStorageMock(t)
	ps.CurrentFunc = func(p context.Context) (r *core.Pulse, r1 error) {
		return &core.Pulse{}, nil
	}
	c.PulseStorage = ps

	err := c.Start(ctx)
	require.NoError(t, err)
	require.Len(t, c.slotMap, 2)
}

func TestConveyer_Start_Err(t *testing.T) {
	ctx := context.Background()
	c := NewPulseConveyer()
	ps := testutils.NewPulseStorageMock(t)
	errMessage := "some fancy error"
	ps.CurrentFunc = func(p context.Context) (r *core.Pulse, r1 error) {
		return &core.Pulse{}, fmt.Errorf(errMessage)
	}
	c.PulseStorage = ps

	err := c.Start(ctx)
	require.Error(t, err)
	require.EqualError(t, err, errMessage)
	require.Empty(t, c.slotMap)
}
