///
//    Copyright 2019 Insolar Technologies
//
//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.
///

package smachine

import (
	"github.com/insolar/insolar/conveyor/smachine/tools"
	"sync"
)

func NewSlotMachineSync(signalCallback func()) SlotMachineSync {
	return SlotMachineSync{
		updateQueue:   tools.NewSignalFuncQueue(&sync.Mutex{}, signalCallback),
		callbackQueue: tools.NewSignalFuncQueue(&sync.Mutex{}, signalCallback),
	}
}

type SlotMachineSync struct {
	updateQueue   tools.SyncQueue // func(w SlotWorker) // for detached/async ops, queued functions MUST BE panic-safe
	callbackQueue tools.SyncQueue // func(w DetachableSlotWorker) // for detached/async ops, queued functions MUST BE panic-safe

	detachLock   sync.RWMutex
	detachQueues map[SlotLink]*tools.SyncQueue
}

func (m *SlotMachineSync) IsZero() bool {
	return m.updateQueue.Locker() == nil
}

/* This method MUST ONLY be used for own operations of SlotMachine, no StateMachine handlers are allowed  */
func (m *SlotMachineSync) AddAsyncUpdate(link SlotLink, fn func(link SlotLink, worker FixedSlotWorker)) {
	if fn == nil {
		panic("illegal value")
	}

	m.updateQueue.Add(func(w interface{}) {
		fn(link, w.(FixedSlotWorker))
	})
}

func (m *SlotMachineSync) ProcessUpdates(worker FixedSlotWorker) bool {
	tasks := m.updateQueue.Flush()
	if len(tasks) == 0 {
		return false
	}

	for _, fn := range tasks {
		fn(worker)
	}
	return true
}

type AsyncCallbackFunc func(link SlotLink, worker DetachableSlotWorker) bool

func (m *SlotMachineSync) AddAsyncCallback(link SlotLink, fn AsyncCallbackFunc) {
	if fn == nil {
		panic("illegal value")
	}

	m._addAsyncCallback(&m.callbackQueue, link, fn, 0)
}

func (m *SlotMachineSync) _addAsyncCallback(q *tools.SyncQueue, link SlotLink, fn AsyncCallbackFunc, repeatCount int) {
	q.Add(func(w interface{}) {
		if !fn(link, w.(DetachableSlotWorker)) {
			m._addDetachedCallback(link, fn, repeatCount+1)

			if !link.isValidAndBusy() { // double-check to avoid misses of detached elements
				m.AppendSlotDetachQueue(link)
			}
		}
	})
}

func (m *SlotMachineSync) ProcessCallbacks(worker FixedSlotWorker) (hasSignal, wasDetached bool) {

	if worker.HasSignal() {
		return true, false
	}

	tasks := m.callbackQueue.Flush()
	if len(tasks) == 0 {
		return false, false
	}

	hasSignal = false
	wasCalled := false

	wasDetached = worker.DetachableCall(func(w DetachableSlotWorker) {
		hasSignal = m.processCallbacks(tasks, w)
		wasCalled = true
	})

	if !wasCalled {
		m.callbackQueue.AddAll(tasks)
	}

	return hasSignal, wasDetached
}

func (m *SlotMachineSync) processCallbacks(tasks tools.SyncFuncList, worker DetachableSlotWorker) (hasSignal bool) {
	for i, fn := range tasks {
		fn(worker)
		if worker.HasSignal() {
			m.callbackQueue.AddAll(tasks[i+1:])
			return true
		}
	}
	return false
}

func (m *SlotMachineSync) _addDetachedCallback(link SlotLink, fn AsyncCallbackFunc, repeatCount int) {
	if repeatCount > 100 {
		fn(link, nil)
		return
	}

	m.detachLock.RLock()
	dq := m.detachQueues[link]
	m.detachLock.RUnlock()

	if dq == nil {
		dqv := tools.NewSignalFuncQueue(&sync.Mutex{}, nil)

		m.detachLock.Lock()
		dq = m.detachQueues[link]
		if dq == nil {
			dq = &dqv
			m.detachQueues[link] = dq
		}
		m.detachLock.Unlock()
	}

	m._addAsyncCallback(dq, link, fn, repeatCount+1)
}

func (m *SlotMachineSync) FlushDetachQueue(link SlotLink) tools.SyncFuncList {
	m.detachLock.RLock()
	dq := m.detachQueues[link]
	if dq != nil {
		delete(m.detachQueues, link)
	}
	m.detachLock.RUnlock()
	if dq == nil {
		return nil
	}
	return dq.Flush()
}

func (m *SlotMachineSync) ProcessDetachQueue(link SlotLink, worker DetachableSlotWorker) (hasSignal bool) {
	if worker.HasSignal() {
		return true
	}

	tasks := m.FlushDetachQueue(link)
	if len(tasks) == 0 {
		return false
	}

	return m.processCallbacks(tasks, worker)
}

func (m *SlotMachineSync) AppendSlotDetachQueue(link SlotLink) {
	detached := m.FlushDetachQueue(link)
	if len(detached) > 0 {
		m.callbackQueue.AddAll(detached)
	}
}

func (m *SlotMachineSync) CleanupDetachQueues(cleanupFn func(link SlotLink, list tools.SyncFuncList)) {
	m.detachLock.Lock()
	defer m.detachLock.Unlock()
	for link, dq := range m.detachQueues {
		if link.IsValid() {
			continue
		}
		delete(m.detachQueues, link)
		list := dq.Flush()
		if len(list) != 0 && cleanupFn != nil {
			cleanupFn(link, list)
		}
	}
}
