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
	"sync/atomic"
	"time"

	"github.com/insolar/insolar/conveyor/queue"
	"github.com/insolar/insolar/log"
	"github.com/pkg/errors"
)

type processElement struct {
	elementId   idType
	handlerId   idType
	taskPayload interface{}
	respSink    PulseConveyorSlotResponseSink
	cancelInfo  *cancelInfoT
}

type cancelInfoT struct {
	id     uint64
	cancel chan bool
	flush  chan bool
}

func newCancelInfo(id uint64) *cancelInfoT {
	return &cancelInfoT{
		id:     id,
		cancel: make(chan bool, 1),
		flush:  make(chan bool, 1),
	}
}

type taskHolderT struct {
	taskHolderLock sync.Mutex
	tasks          map[uint32][]*cancelInfoT
}

func newTaskHolder() taskHolderT {
	return taskHolderT{
		tasks: make(map[uint32][]*cancelInfoT),
	}
}

func (th *taskHolderT) add(info *cancelInfoT, pulseNumber uint32) {
	log.Infof("[ taskHolderT.add ] Adding pulseNumber: %d. Id: %d", pulseNumber, info.id)
	th.taskHolderLock.Lock()
	defer th.taskHolderLock.Unlock()

	el, ok := th.tasks[pulseNumber]
	if !ok {
		th.tasks[pulseNumber] = []*cancelInfoT{info}
	} else {
		th.tasks[pulseNumber] = append(el, info)
	}
}

func processStop(cancelList []*cancelInfoT, flush bool) {
	for _, el := range cancelList {
		if flush {
			log.Info("[ processStop ] flush: ", el.id)
			el.flush <- true
		} else {
			log.Info("[ processStop ] cancel: ", el.id)
			el.cancel <- true
		}
	}
}

func (th *taskHolderT) stop(pulseNumber uint32, flush bool) {
	log.Infof("[ taskHolderT.stop ] Stopping pulseNumber: %d, flush: %d", pulseNumber, flush)
	th.taskHolderLock.Lock()
	defer th.taskHolderLock.Unlock()

	cancelList, ok := th.tasks[pulseNumber]
	if !ok {
		log.Info("[ taskHolderT.stop ] No such pulseNumber: ", pulseNumber)
		return
	}

	processStop(cancelList, flush)

	delete(th.tasks, pulseNumber)
}

func (th *taskHolderT) stopAll(flush bool) {
	th.taskHolderLock.Lock()
	defer th.taskHolderLock.Unlock()

	log.Infof("[ taskHolderT.stopAll ] flush: ", flush)

	for _, cancelList := range th.tasks {
		processStop(cancelList, flush)
	}

	th.tasks = make(map[uint32][]*cancelInfoT)

}

// SimpleWaitAdapter holds all adapter logic
type SimpleWaitAdapter struct {
	queue             queue.IQueue
	processingStarted uint32
	stopProcessing    uint32
	processingStopped chan bool
	// adapterID comes from configuration
	adapterID uint32

	taskHolder taskHolderT
}

type simpleWaitAdapterInputData struct {
	waitPeriodMilliseconds int
}

type simpleWaitAdapterOutputData struct {
	info string
}

// NewSimpleWaitAdapter creates new instance of SimpleWaitAdapter
func NewSimpleWaitAdapter() PulseConveyorAdapterTaskSink {
	adapter := &SimpleWaitAdapter{
		queue:             queue.NewMutexQueue(),
		processingStarted: 0,
		stopProcessing:    0,
		processingStopped: make(chan bool, 1),
		taskHolder:        newTaskHolder(),
	}
	started := make(chan bool, 1)
	go adapter.StartProcessing(started)
	<-started

	return adapter
}

// StopProcessing is blocking
func (swa *SimpleWaitAdapter) StopProcessing() {
	if atomic.LoadUint32(&swa.stopProcessing) != 0 {
		log.Infof("[ StopProcessing ]  Nothing done")
		return
	}
	atomic.StoreUint32(&swa.stopProcessing, 1)
	<-swa.processingStopped
}

// StartProcessing start processing of input queue
func (swa *SimpleWaitAdapter) StartProcessing(started chan bool) {
	if atomic.LoadUint32(&swa.processingStarted) != 0 {
		log.Infof("[ StartProcessing ] processing already started. Nothing done")
		close(started)
		return
	}
	atomic.StoreUint32(&swa.processingStarted, 1)

	started <- true

	lastLoop := false
	for true {

		var itasks []queue.OutputElement

		if atomic.LoadUint32(&swa.stopProcessing) != 0 {
			if lastLoop {
				log.Infof("[ StartProcessing ] Stop processing. EXIT")
				break
			}
			itasks = swa.queue.BlockAndRemoveAll()
			log.Info("[ StartProcessing ] Stop processing: one more loop")
			lastLoop = true
		} else {
			itasks = swa.queue.RemoveAll()
		}

		log.Infof("[ StartProcessing ] Got %d new tasks", len(itasks))

		if len(itasks) == 0 {
			log.Info("[ StartProcessing ] No tasks. Sleep a little bit")
			// TODO: do pretty wait
			time.Sleep(50 * time.Millisecond)
			continue
		}

		for _, itask := range itasks {
			task, ok := itask.GetData().(processElement)
			if !ok {
				panic(fmt.Sprintf("[ StartProcessing ] How does it happen? Wrong Type: %T", itask.GetData()))
			}

			go swa.doWork(task, task.cancelInfo)
		}
	}

	swa.processingStopped <- true
}

// it's function which make useful adapter's work
func (swa *SimpleWaitAdapter) doWork(task processElement, cancelInfo *cancelInfoT) {

	log.Info("[ doWork ] Start. cancelInfo.id: ", cancelInfo.id)

	payload := task.taskPayload.(simpleWaitAdapterInputData)

	var msg string
	select {
	case <-cancelInfo.cancel:
		msg = "Cancel. Return Response"
	case <-cancelInfo.flush:
		log.Info("[ SimpleWaitAdapter.doWork ] Flush. DON'T Return Response")
		return
	case <-time.After(time.Duration(payload.waitPeriodMilliseconds) * time.Millisecond):
		msg = fmt.Sprintf("Work completed successfully. Waited %d millisecond", payload.waitPeriodMilliseconds)
	}

	log.Info("[ SimpleWaitAdapter.doWork ] ", msg)

	task.respSink.PushResponse(swa.adapterID,
		task.elementId,
		task.handlerId,
		fmt.Sprintf(msg))

	// TODO: remove cancelInfo from swa.taskHolder
}

var reqId uint64 = 0

func atomicLoadAndIncrementUint64(addr *uint64) uint64 {
	for {
		val := atomic.LoadUint64(addr)
		if atomic.CompareAndSwapUint64(addr, val, val+1) {
			return val
		}
	}
}

// PushTask implements PulseConveyorAdapterTaskSink
func (swa *SimpleWaitAdapter) PushTask(respSink PulseConveyorSlotResponseSink,
	elementId idType,
	handlerId idType,
	taskPayload interface{}) error {

	payload, ok := taskPayload.(simpleWaitAdapterInputData)
	if !ok {
		return errors.Errorf("[ PushTask ] Incorrect payload type: %T", taskPayload)
	}

	cancelInfo := newCancelInfo(atomicLoadAndIncrementUint64(&reqId))
	swa.taskHolder.add(cancelInfo, respSink.GetPulseNumber())

	return swa.queue.SinkPush(processElement{
		respSink:    respSink,
		elementId:   elementId,
		handlerId:   handlerId,
		taskPayload: payload,
		cancelInfo:  cancelInfo,
	})
}

// CancelElementTasks: now cancels all pulseNumber's tasks
func (swa *SimpleWaitAdapter) CancelElementTasks(pulseNumber idType, elementId idType) {
	swa.taskHolder.stop(pulseNumber, false)
}

// CancelPulseTasks: now cancels all pulseNumber's tasks
func (swa *SimpleWaitAdapter) CancelPulseTasks(pulseNumber idType) {
	swa.taskHolder.stop(pulseNumber, false)
}

// FlushPulseTasks: now flush all pulseNumber's tasks
func (swa *SimpleWaitAdapter) FlushPulseTasks(pulseNumber uint32) {
	swa.taskHolder.stop(pulseNumber, true)
}

// FlushNodeTasks: now flush all tasks
func (swa *SimpleWaitAdapter) FlushNodeTasks(nodeId idType) {
	swa.taskHolder.stopAll(true)
}
