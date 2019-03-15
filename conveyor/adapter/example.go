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
)

// queueTask is task for adapter with queue
type queueTask struct {
	cancelInfo *cancelInfo
	task       AdapterTask
}

type cancelInfo struct {
	id     uint64
	cancel chan bool
	flush  chan bool
}

func newCancelInfo(id uint64) *cancelInfo {
	return &cancelInfo{
		id:     id,
		cancel: make(chan bool, 1),
		flush:  make(chan bool, 1),
	}
}

func (ci *cancelInfo) ID() uint64 {
	return ci.id
}

func (ci *cancelInfo) Cancel() chan bool {
	return ci.cancel
}

func (ci *cancelInfo) Flush() chan bool {
	return ci.flush
}

type taskHolder struct {
	taskHolderLock sync.Mutex
	tasks          map[uint32][]*cancelInfo
}

func newTaskHolder() taskHolder {
	return taskHolder{
		tasks: make(map[uint32][]*cancelInfo),
	}
}

func (th *taskHolder) add(info *cancelInfo, pulseNumber uint32) {
	log.Infof("[ taskHolder.add ] Adding pulseNumber: %d. Id: %d", pulseNumber, info.id)
	th.taskHolderLock.Lock()
	defer th.taskHolderLock.Unlock()

	el, ok := th.tasks[pulseNumber]
	if !ok {
		th.tasks[pulseNumber] = []*cancelInfo{info}
	} else {
		th.tasks[pulseNumber] = append(el, info)
	}
}

func processStop(cancelList []*cancelInfo, flush bool) {
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

func (th *taskHolder) stop(pulseNumber uint32, flush bool) {
	log.Infof("[ taskHolder.stop ] Stopping pulseNumber: %d, flush: %d", pulseNumber, flush)
	th.taskHolderLock.Lock()
	defer th.taskHolderLock.Unlock()

	cancelList, ok := th.tasks[pulseNumber]
	if !ok {
		log.Info("[ taskHolder.stop ] No such pulseNumber: ", pulseNumber)
		return
	}

	processStop(cancelList, flush)

	delete(th.tasks, pulseNumber)
}

func (th *taskHolder) stopAll(flush bool) {
	th.taskHolderLock.Lock()
	defer th.taskHolderLock.Unlock()

	log.Infof("[ taskHolder.stopAll ] flush: ", flush)

	for _, cancelList := range th.tasks {
		processStop(cancelList, flush)
	}

	th.tasks = make(map[uint32][]*cancelInfo)

}

// AdapterWithQueue holds all adapter logic
type AdapterWithQueue struct {
	queue             queue.IQueue
	processingStarted uint32
	stopProcessing    uint32
	processingStopped chan bool
	// adapterID comes from configuration
	adapterID uint32

	taskHolder taskHolder
	worker     Worker
}

// StopProcessing is blocking
func (swa *AdapterWithQueue) StopProcessing() {
	if atomic.LoadUint32(&swa.stopProcessing) != 0 {
		log.Infof("[ StopProcessing ]  Nothing done")
		return
	}
	atomic.StoreUint32(&swa.stopProcessing, 1)
	<-swa.processingStopped
}

// StartProcessing start processing of input queue
func (swa *AdapterWithQueue) StartProcessing(started chan bool) {
	if atomic.LoadUint32(&swa.processingStarted) != 0 {
		log.Infof("[ StartProcessing ] processing already started. Nothing done")
		close(started)
		return
	}
	atomic.StoreUint32(&swa.processingStarted, 1)

	started <- true

	lastLoop := false
	for {

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
			task, ok := itask.GetData().(queueTask)
			if !ok {
				panic(fmt.Sprintf("[ StartProcessing ] How does it happen? Wrong Type: %T", itask.GetData()))
			}

			if swa.worker == nil {
				panic(fmt.Sprintf("[ StartProcessing ] Worker function wasn't provided"))
			}

			go swa.worker.Process(swa.adapterID, task.task, task.cancelInfo)
		}
	}

	swa.processingStopped <- true
}

var reqID uint64

func atomicLoadAndIncrementUint64(addr *uint64) uint64 {
	for {
		val := atomic.LoadUint64(addr)
		if atomic.CompareAndSwapUint64(addr, val, val+1) {
			return val
		}
	}
}

// PushTask implements PulseConveyorAdapterTaskSink
func (swa *AdapterWithQueue) PushTask(respSink AdaptorToSlotResponseSink,
	elementID idType,
	handlerID idType,
	taskPayload interface{}) error {

	cancelInfo := newCancelInfo(atomicLoadAndIncrementUint64(&reqID))
	swa.taskHolder.add(cancelInfo, respSink.GetPulseNumber())

	return swa.queue.SinkPush(
		queueTask{
			cancelInfo: cancelInfo,
			task: AdapterTask{
				respSink:    respSink,
				elementID:   elementID,
				handlerID:   handlerID,
				taskPayload: taskPayload,
			},
		},
	)
}

// CancelElementTasks: now cancels all pulseNumber's tasks
func (swa *AdapterWithQueue) CancelElementTasks(pulseNumber idType, elementID idType) {
	swa.taskHolder.stop(pulseNumber, false)
}

// CancelPulseTasks: now cancels all pulseNumber's tasks
func (swa *AdapterWithQueue) CancelPulseTasks(pulseNumber idType) {
	swa.taskHolder.stop(pulseNumber, false)
}

// FlushPulseTasks: now flush all pulseNumber's tasks
func (swa *AdapterWithQueue) FlushPulseTasks(pulseNumber uint32) {
	swa.taskHolder.stop(pulseNumber, true)
}

// FlushNodeTasks: now flush all tasks
func (swa *AdapterWithQueue) FlushNodeTasks(nodeID idType) {
	swa.taskHolder.stopAll(true)
}
