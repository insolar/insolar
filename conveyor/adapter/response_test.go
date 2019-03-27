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

package adapter

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/insolar/insolar/conveyor/adapter/adapterid"
	"github.com/insolar/insolar/conveyor/interfaces/slot"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/messagebus"
	"github.com/insolar/insolar/testutils"
	"github.com/stretchr/testify/require"
)

type replyMock int

func (replyMock) Type() insolar.ReplyType {
	return insolar.ReplyType(124)
}
func testResponseSenderTask(t *testing.T) SendResponseTask {
	res := replyMock(111)
	return SendResponseTask{Future: mockConveyorFuture(t, res), Result: res}
}

func mockConveyorFuture(t *testing.T, result insolar.Reply) *testutils.ConveyorFutureMock {
	cfMock := testutils.NewConveyorFutureMock(t)
	cfMock.SetResultFunc = func(p insolar.Reply) {
		if result != p {
			panic(fmt.Sprintf("Expected result %v, get %v", result, p))
		}
	}
	return cfMock
}

func startResponseSendAdapter() *CancellableQueueAdapter {
	queueAdapter := NewAdapterWithQueue(NewSendResponseProcessor(), adapterid.SendResponse)
	adapter := queueAdapter.(*CancellableQueueAdapter)
	started := make(chan bool, 1)
	adapter.StartProcessing(started)
	<-started
	return adapter
}

func TestResponseSendAdapter_PushTask_IncorrectPayload(t *testing.T) {
	adapter := startResponseSendAdapter()
	resp := &mockResponseSink{}

	err := adapter.PushTask(resp, 33, 22, 22)
	require.NoError(t, err)
	// TODO: wait until swa.taskHolder is empty (i.e. task was done)
	time.Sleep(200 * time.Millisecond)
	require.Error(t, resp.GetResponse().(error))
	require.Contains(t, resp.GetResponse().(error).Error(), "[ SendResponseProcessor.Process ] Incorrect payload type: int")
}

func TestResponseSendAdapter_PushTask(t *testing.T) {
	adapter := startResponseSendAdapter()
	resp := &mockResponseSink{}

	err := adapter.PushTask(resp, 33, 22, testResponseSenderTask(t))
	require.NoError(t, err)
	// TODO: wait until swa.taskHolder is empty (i.e. task was done)
	time.Sleep(200 * time.Millisecond)
	require.Contains(t, resp.GetResponse().(string), "Response was send successfully")
}

func TestResponseSendAdapter_PushTask_AfterStopProcessing(t *testing.T) {
	adapter := startResponseSendAdapter()
	resp := &mockResponseSink{}
	adapter.StopProcessing()

	err := adapter.PushTask(resp, 34, 22, testResponseSenderTask(t))
	// TODO: check swa.taskHolder is empty (i.e. task wasn't created)
	require.Contains(t, err.Error(), "Queue is blocked")
}

func TestResponseSendAdapter_CancelPulseTasks(t *testing.T) {
	adapter := startResponseSendAdapter()
	resp := &mockResponseSink{response: "startValue"}
	err := adapter.PushTask(resp, 33, 22, testResponseSenderTask(t))
	require.NoError(t, err)

	adapter.CancelPulseTasks(resp.GetPulseNumber())
	// TODO: wait until swa.taskHolder is empty (i.e. task was done)
	time.Sleep(200 * time.Millisecond)
	require.Nil(t, resp.GetResponse())
}

func TestResponseSendAdapter_FlushPulseTasks(t *testing.T) {
	adapter := startResponseSendAdapter()
	resp := &mockResponseSink{response: "startValue"}
	err := adapter.PushTask(resp, 34, 22, testResponseSenderTask(t))
	require.NoError(t, err)

	adapter.FlushPulseTasks(resp.GetPulseNumber())
	// TODO: wait until swa.taskHolder is empty (i.e. task was done)
	time.Sleep(200 * time.Millisecond)
	require.Equal(t, "startValue", resp.GetResponse())
}

func TestResponseSendAdapter_Parallel(t *testing.T) {
	adapter := startResponseSendAdapter()

	numIterations := 200
	parallelPushTasks := 27

	wg := sync.WaitGroup{}
	wg.Add(parallelPushTasks)

	// PushTask
	for i := 0; i < parallelPushTasks; i++ {
		go func(wg *sync.WaitGroup, adapter TaskSink) {
			for i := 0; i < numIterations; i++ {
				resp := &mockResponseSink{}
				adapter.PushTask(resp, 34, 22, testResponseSenderTask(t))
			}
			wg.Done()
		}(&wg, adapter)
	}

	wg.Wait()

	adapter.StopProcessing()

}

type mockReply struct {
	data string
}

func (mr *mockReply) Type() insolar.ReplyType {
	return 0
}

func TestSendResponseHelper(t *testing.T) {
	f := messagebus.NewFuture()
	event := insolar.ConveyorPendingMessage{Future: f}
	testReply := &mockReply{data: "Put-in"}

	slotElementHelperMock := slot.NewSlotElementHelperMock(t)
	slotElementHelperMock.GetInputEventFunc = func() (r interface{}) {
		return event
	}
	slotElementHelperMock.SendTaskFunc = func(p adapterid.ID, response interface{}, p2 uint32) (r error) {
		f := response.(SendResponseTask).Future
		f.SetResult(testReply)
		return nil
	}

	adapterCatalog := newHelperCatalog()
	err := adapterCatalog.sendResponseHelper.SendResponse(slotElementHelperMock, testReply, 42)
	require.NoError(t, err)

	gotReply, err := f.GetResult(time.Second)
	require.NoError(t, err)
	require.Equal(t, testReply, gotReply)
}

func TestSendResponseHelper_BadInput(t *testing.T) {
	slotElementHelperMock := slot.NewSlotElementHelperMock(t)
	slotElementHelperMock.GetInputEventFunc = func() (r interface{}) {
		return 33
	}
	adapterCatalog := newHelperCatalog()
	err := adapterCatalog.sendResponseHelper.SendResponse(slotElementHelperMock, &mockReply{}, 44)
	require.Error(t, err)
	require.Contains(t, err.Error(), "Input event is not insolar.ConveyorPendingMessage")
}
