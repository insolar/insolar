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
	"sync"
	"testing"
	"time"

	"github.com/insolar/insolar/conveyor/interfaces/islot"
	"github.com/insolar/insolar/log"
	"github.com/stretchr/testify/require"
)

type mockResponseSink struct {
	response string
	lock     sync.Mutex
}

func (m *mockResponseSink) PushResponse(adapterID idType, elementID idType, handlerID idType, respPayload interface{}) {
	m.lock.Lock()
	defer m.lock.Unlock()
	log.Infof("[ mockResponseSink.PushResponse] PushResponse: %+v", respPayload)
	m.response = respPayload.(string)
}

func (m *mockResponseSink) GetResponse() string {
	m.lock.Lock()
	defer m.lock.Unlock()

	return m.response
}

func (m *mockResponseSink) PushNestedEvent(adapterID idType, parentElementID idType, handlerID idType, eventPayload interface{}) {
	m.lock.Lock()
	defer m.lock.Unlock()
	log.Infof("[ mockResponseSink.PushNestedEvent] PushNestedEvent: %+v", eventPayload)
	m.response = eventPayload.(string)
}

func (m *mockResponseSink) GetPulseNumber() uint32 {
	return 142
}

func (m *mockResponseSink) GetNodeID() uint32 {
	return 42
}

func (m *mockResponseSink) GetSlotDetails() islot.SlotDetails {
	return islot.NewSlotDetailsMock(&testing.T{})
}

func TestFunctionality(t *testing.T) {
	adapter := NewSimpleWaitAdapter().(*SimpleWaitAdapter)
	started := make(chan bool, 1)
	adapter.StartProcessing(started)
	<-started

	adapter.CancelElementTasks(33, 22)
	adapter.CancelPulseTasks(44)
	adapter.FlushNodeTasks(55)
	adapter.FlushPulseTasks(66)

	err := adapter.PushTask(&mockResponseSink{}, 33, 22, 22)
	require.Contains(t, err.Error(), "Incorrect payload type")

	resp := &mockResponseSink{}
	err = adapter.PushTask(resp, 33, 22, simpleWaitAdapterInputData{waitPeriodMilliseconds: 20})
	require.NoError(t, err)
	time.Sleep(200 * time.Millisecond)
	require.Contains(t, resp.GetResponse(), "Work completed successfully")

	// CancelPulseTasks test
	resp = &mockResponseSink{}
	err = adapter.PushTask(resp, 33, 22, simpleWaitAdapterInputData{waitPeriodMilliseconds: 200000000})
	require.NoError(t, err)
	adapter.CancelPulseTasks(resp.GetPulseNumber())
	time.Sleep(200 * time.Millisecond)
	require.Contains(t, resp.GetResponse(), "Cancel. Return Response")

	// FlushPulseTasks
	resp = &mockResponseSink{}
	err = adapter.PushTask(resp, 34, 22, simpleWaitAdapterInputData{waitPeriodMilliseconds: 200000000})
	require.NoError(t, err)
	adapter.FlushPulseTasks(resp.GetPulseNumber())

	// FlushPulseTasks
	resp = &mockResponseSink{}
	err = adapter.PushTask(resp, 34, 22, simpleWaitAdapterInputData{waitPeriodMilliseconds: 200000000})
	require.NoError(t, err)
	adapter.FlushNodeTasks(resp.GetPulseNumber())

	adapter.StopProcessing()
	adapter.StopProcessing()

	err = adapter.PushTask(resp, 34, 22, simpleWaitAdapterInputData{waitPeriodMilliseconds: 200000000})
	require.Contains(t, err.Error(), "Queue is blocked")
}

func TestParallel(t *testing.T) {

	adapter := NewSimpleWaitAdapter().(*SimpleWaitAdapter)
	started := make(chan bool, 1)
	adapter.StartProcessing(started)
	<-started

	res := mockResponseSink{}

	pulseNumber := res.GetPulseNumber()

	numIterations := 200
	parallelPushTasks := 27
	parallelCancelElement := 19
	parallelCancelPulse := 13
	parallelFlushPulse := 9
	parallelFlushNode := 5

	wg := sync.WaitGroup{}
	wg.Add(parallelPushTasks + parallelCancelElement + parallelCancelPulse + parallelFlushPulse + parallelFlushNode)

	// PushTask
	for i := 0; i < parallelPushTasks; i++ {
		go func(wg *sync.WaitGroup, adapter PulseConveyorAdapterTaskSink) {
			for i := 0; i < numIterations; i++ {
				resp := &mockResponseSink{}
				adapter.PushTask(resp, 34, 22, simpleWaitAdapterInputData{waitPeriodMilliseconds: 20})
			}
			wg.Done()
		}(&wg, adapter)
	}

	// CancelElementTasks
	for i := 0; i < parallelCancelElement; i++ {
		go func(wg *sync.WaitGroup, adapter PulseConveyorAdapterTaskSink) {
			for i := 0; i < numIterations; i++ {
				adapter.CancelElementTasks(pulseNumber, 22)
			}
			wg.Done()
		}(&wg, adapter)
	}

	// CancelPulseTasks
	for i := 0; i < parallelCancelPulse; i++ {
		go func(wg *sync.WaitGroup, adapter PulseConveyorAdapterTaskSink) {
			for i := 0; i < numIterations; i++ {
				adapter.CancelPulseTasks(pulseNumber)
			}
			wg.Done()
		}(&wg, adapter)
	}

	// FlushPulseTasks
	for i := 0; i < parallelFlushPulse; i++ {
		go func(wg *sync.WaitGroup, adapter PulseConveyorAdapterTaskSink) {
			for i := 0; i < numIterations; i++ {
				adapter.FlushPulseTasks(pulseNumber)
			}
			wg.Done()
		}(&wg, adapter)
	}

	// FlushNodeTasks
	for i := 0; i < parallelFlushNode; i++ {
		go func(wg *sync.WaitGroup, adapter PulseConveyorAdapterTaskSink) {
			for i := 0; i < numIterations; i++ {
				adapter.FlushNodeTasks(pulseNumber)
			}
			wg.Done()
		}(&wg, adapter)
	}

	wg.Wait()

	adapter.StopProcessing()

}
