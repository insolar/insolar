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
	id         uint64
	cancel     chan bool
	flush      chan bool
	isCanceled uint32
	isFlushed  uint32
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

func (ci *cancelInfo) Cancel() <-chan bool {
	return ci.cancel
}

func (ci *cancelInfo) Flush() <-chan bool {
	return ci.flush
}

func (ci *cancelInfo) IsCanceled() bool {
	return atomic.LoadUint32(&ci.isCanceled) != 0
}

func (ci *cancelInfo) IsFlushed() bool {
	return atomic.LoadUint32(&ci.isFlushed) != 0
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
			atomic.StoreUint32(&el.isFlushed, 1)
			el.flush <- true
		} else {
			log.Info("[ processStop ] cancel: ", el.id)
			atomic.StoreUint32(&el.isCanceled, 1)
			el.cancel <- true
		}
	}
}

func (th *taskHolder) stop(pulseNumber uint32, flush bool) {
	log.Infof("[ taskHolder.stop ] Stopping pulseNumber: %d, flush: %s", pulseNumber, flush)
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

	log.Info("[ taskHolder.stopAll ] flush: ", flush)

	for _, cancelList := range th.tasks {
		processStop(cancelList, flush)
	}

	th.tasks = make(map[uint32][]*cancelInfo)

}

// CancellableQueueAdapter holds all adapter logic
type CancellableQueueAdapter struct {
	queue             queue.IQueue
	processingStarted uint32
	stopProcessing    uint32
	processingStopped chan bool
	// adapterID comes from configuration
	adapterID uint32

	taskHolder taskHolder
	processor  Processor
}

func (c *CancellableQueueAdapter) GetAdapterID() uint32 {
	return c.adapterID
}

// StopProcessing is blocking
func (a *CancellableQueueAdapter) StopProcessing() {
	if atomic.LoadUint32(&a.stopProcessing) != 0 {
		log.Infof("[ StopProcessing ]  Nothing done")
		return
	}
	atomic.StoreUint32(&a.stopProcessing, 1)
	<-a.processingStopped
}

// StartProcessing start processing of input queue
func (a *CancellableQueueAdapter) StartProcessing(started chan bool) {
	if atomic.LoadUint32(&a.processingStarted) != 0 {
		log.Infof("[ StartProcessing ] processing already started. Nothing done")
		close(started)
		return
	}
	atomic.StoreUint32(&a.processingStarted, 1)

	started <- true

	lastLoop := false
	for {

		var itasks []queue.OutputElement

		if atomic.LoadUint32(&a.stopProcessing) != 0 {
			if lastLoop {
				log.Infof("[ StartProcessing ] Stop processing. EXIT")
				break
			}
			itasks = a.queue.BlockAndRemoveAll()
			log.Info("[ StartProcessing ] Stop processing: one more loop")
			lastLoop = true
		} else {
			itasks = a.queue.RemoveAll()
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

			if a.processor == nil {
				panic(fmt.Sprintf("[ StartProcessing ] Processor function wasn't provided"))
			}

			go a.process(task)
		}
	}

	a.processingStopped <- true
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
func (a *CancellableQueueAdapter) PushTask(respSink AdapterToSlotResponseSink,
	elementID idType,
	handlerID idType,
	taskPayload interface{}) error {

	cancelInfo := newCancelInfo(atomicLoadAndIncrementUint64(&reqID))
	a.taskHolder.add(cancelInfo, respSink.GetPulseNumber())

	return a.queue.SinkPush(
		queueTask{
			cancelInfo: cancelInfo,
			task: AdapterTask{
				respSink:    respSink,
				elementID:   elementID,
				handlerID:   handlerID,
				TaskPayload: taskPayload,
			},
		},
	)
}

// CancelElementTasks: now cancels all pulseNumber's tasks
func (a *CancellableQueueAdapter) CancelElementTasks(pulseNumber idType, elementID idType) {
	a.taskHolder.stop(pulseNumber, false)
}

// CancelPulseTasks: now cancels all pulseNumber's tasks
func (a *CancellableQueueAdapter) CancelPulseTasks(pulseNumber idType) {
	a.taskHolder.stop(pulseNumber, false)
}

// FlushPulseTasks: now flush all pulseNumber's tasks
func (a *CancellableQueueAdapter) FlushPulseTasks(pulseNumber uint32) {
	a.taskHolder.stop(pulseNumber, true)
}

// FlushNodeTasks: now flush all tasks
func (a *CancellableQueueAdapter) FlushNodeTasks(nodeID idType) {
	a.taskHolder.stopAll(true)
}

func (a *CancellableQueueAdapter) process(cancellableTask queueTask) {
	adapterTask := cancellableTask.task
	respSink := adapterTask.respSink
	cancelInfo := cancellableTask.cancelInfo

	select {
	case <-cancelInfo.Cancel():
		log.Info("[ CancellableQueueAdapter.process ] Task was canceled")
		respSink.PushResponse(a.adapterID, adapterTask.elementID, adapterTask.handlerID, nil)
	case <-cancelInfo.Flush():
		log.Info("[ CancellableQueueAdapter.process ] Task was flushed. Don't push Response")
	default:
		helper := newNestedEventHelper(adapterTask, a.adapterID)
		respPayload := a.processor.Process(adapterTask, helper, cancelInfo)
		if !cancelInfo.IsFlushed() {
			respSink.PushResponse(a.adapterID, adapterTask.elementID, adapterTask.handlerID, respPayload)
			// TODO: remove cancelInfo from a.taskHolder
		}
	}
}

type nestedEventHelper struct {
	adapterTask AdapterTask
	adapterID   uint32
}

func newNestedEventHelper(adapterTask AdapterTask, adapterID uint32) NestedEventHelper {
	return &nestedEventHelper{
		adapterTask: adapterTask,
		adapterID:   adapterID,
	}
}

func (h *nestedEventHelper) Send(eventPayload interface{}) {
	h.adapterTask.respSink.PushNestedEvent(h.adapterID, h.adapterTask.elementID, h.adapterTask.handlerID, eventPayload)
}
