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
	"sync"
	"sync/atomic"

	"github.com/insolar/insolar/conveyor/tools"
)

func NewSlotMachineSync(eventCallback, signalCallback func()) SlotMachineSync {
	return SlotMachineSync{
		signalQueue:   tools.NewSignalFuncQueue(&sync.Mutex{}, signalCallback),
		updateQueue:   tools.NewSignalFuncQueue(&sync.Mutex{}, eventCallback),
		callbackQueue: tools.NewSignalFuncQueue(&sync.Mutex{}, eventCallback),
		machineStatus: uint32(SlotMachineActive),
	}
}

const maxDetachRetries = 100

type SlotMachineStatus uint8

const (
	SlotMachineInactive SlotMachineStatus = iota
	SlotMachineStopping
	SlotMachineActive
)

type SlotMachineSync struct {
	machineStatus uint32 // atomic

	signalQueue   tools.SyncQueue // func(w FixedSlotWorker) // for detached/async ops, queued functions MUST BE panic-safe
	updateQueue   tools.SyncQueue // func(w FixedSlotWorker) // for detached/async ops, queued functions MUST BE panic-safe
	callbackQueue tools.SyncQueue // func(w DetachableSlotWorker) // for detached/async ops, queued functions MUST BE panic-safe

	detachLock   sync.RWMutex
	detachQueues map[SlotLink]*tools.SyncQueue
}

func (m *SlotMachineSync) IsZero() bool {
	return m.updateQueue.Locker() == nil
}

func (m *SlotMachineSync) GetStatus() SlotMachineStatus {
	return SlotMachineStatus(atomic.LoadUint32(&m.machineStatus))
}

func (m *SlotMachineSync) IsActive() bool {
	return m.GetStatus() == SlotMachineActive
}

func (m *SlotMachineSync) IsInactive() bool {
	return m.GetStatus() < SlotMachineStopping
}

func (m *SlotMachineSync) SetStopping() bool {
	return atomic.CompareAndSwapUint32(&m.machineStatus, uint32(SlotMachineActive), uint32(SlotMachineStopping))
}

func (m *SlotMachineSync) SetInactive() bool {
	return atomic.SwapUint32(&m.machineStatus, uint32(SlotMachineInactive)) != uint32(SlotMachineInactive)
}

func (m *SlotMachineSync) FlushAll() {
	m.signalQueue.Flush()
	m.updateQueue.Flush()
	m.callbackQueue.Flush()

	m.detachLock.Lock()
	m.detachQueues = nil
	m.detachLock.Unlock()
}

/* This method MUST ONLY be used for own operations of SlotMachine, no StateMachine handlers are allowed  */
func (m *SlotMachineSync) AddAsyncSignal(link SlotLink, fn func(link SlotLink, worker FixedSlotWorker)) bool {
	switch {
	case fn == nil:
		panic("illegal value")
	case m.IsInactive():
		return false
	}

	m.signalQueue.Add(func(w interface{}) {
		fn(link, w.(FixedSlotWorker))
	})
	return true
}

/* This method MUST ONLY be used for own operations of SlotMachine, no StateMachine handlers are allowed  */
func (m *SlotMachineSync) AddAsyncUpdate(link SlotLink, fn func(link SlotLink, worker FixedSlotWorker)) bool {
	switch {
	case fn == nil:
		panic("illegal value")
	case m.IsInactive():
		return false
	}

	m.updateQueue.Add(func(w interface{}) {
		fn(link, w.(FixedSlotWorker))
	})
	return true
}

func (m *SlotMachineSync) ProcessUpdates(worker FixedSlotWorker) (hasUpdates bool) {
	switch {
	case worker == nil:
		panic("illegal value")
	case m.IsInactive():
		return
	}

	tasks := m.signalQueue.Flush()
	if len(tasks) > 0 {
		hasUpdates = true
		for _, fn := range tasks {
			fn(worker)
		}
	}

	tasks = m.updateQueue.Flush()
	if len(tasks) > 0 {
		hasUpdates = true
		for _, fn := range tasks {
			fn(worker)
		}
	}

	return hasUpdates
}

func (m *SlotMachineSync) CanProcessCallbacks() bool {
	return m.IsActive() // callbacks are cancelled on stopping
}

type AsyncCallbackFunc func(link SlotLink, worker DetachableSlotWorker) bool

func (m *SlotMachineSync) AddAsyncCallback(link SlotLink, fn AsyncCallbackFunc) bool {
	switch {
	case fn == nil:
		panic("illegal value")
	case !m.CanProcessCallbacks(): // callbacks are cancelled on stopping
		fn(link, nil)
		return false
	}

	m._addAsyncCallback(&m.callbackQueue, link, fn, 0)
	return true
}

func (m *SlotMachineSync) _addAsyncCallback(q *tools.SyncQueue, link SlotLink, fn AsyncCallbackFunc, repeatCount int) {
	q.Add(func(v interface{}) {
		w := v.(DetachableSlotWorker)
		if fn(link, w) {
			return
		}
		if repeatCount >= maxDetachRetries || w == nil {
			fn(link, nil)
			return
		}
		m._addDetachedCallback(link, fn, repeatCount+1)
	})
}

func (m *SlotMachineSync) ProcessCallbacks(worker AttachedSlotWorker) (hasUpdates, hasSignal, wasDetached bool) {
	switch {
	case worker == nil:
		panic("illegal value")
	case worker.HasSignal():
		return true, true, false
	}

	tasks := m.callbackQueue.Flush()
	if len(tasks) == 0 {
		return false, false, false
	}

	if !m.CanProcessCallbacks() {
		// cancel all callbacks
		return true, m.processCallbacks(tasks, nil), false
	}

	wasCalled := false
	hasSignal = false
	wasDetached = worker.DetachableCall(func(w DetachableSlotWorker) {
		hasSignal = m.processCallbacks(tasks, w)
		wasCalled = true
	})
	if !wasCalled {
		m.callbackQueue.AddAll(tasks)
	}
	return true, hasSignal, wasDetached
}

func (m *SlotMachineSync) ProcessSlotCallbacksByDetachable(link SlotLink, worker DetachableSlotWorker) (hasUpdates, hasSignal bool) {
	switch {
	case worker == nil:
		panic("illegal value")
	case m.IsInactive():
		return false, false
	case worker.HasSignal():
		return true, true
	}

	tasks := m._flushDetachQueue(link)
	if len(tasks) == 0 {
		return false, false
	}

	hasSignal = m.processCallbacks(tasks, worker)
	return true, hasSignal
}

func (m *SlotMachineSync) cancelCallbacks(tasks tools.SyncFuncList, worker SlotWorker) (hasSignal bool) {
	for i, fn := range tasks {
		fn(nil)
		if worker.HasSignal() {
			m.callbackQueue.AddAll(tasks[i+1:])
			return true
		}
	}
	return false
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
	m.detachLock.RLock()
	dq := m.detachQueues[link]
	m.detachLock.RUnlock()

	if dq == nil {
		dqv := tools.NewSignalFuncQueue(&sync.Mutex{}, nil)

		m.detachLock.Lock()
		dq = m.detachQueues[link]
		if dq == nil {
			dq = &dqv
			if m.detachQueues == nil {
				m.detachQueues = make(map[SlotLink]*tools.SyncQueue)
			}
			m.detachQueues[link] = dq
		}
		m.detachLock.Unlock()
	}

	m._addAsyncCallback(dq, link, fn, repeatCount)
}

func (m *SlotMachineSync) _flushDetachQueue(link SlotLink) tools.SyncFuncList {
	m.detachLock.RLock()
	dq := m.detachQueues[link]
	m.detachLock.RUnlock()
	if dq == nil {
		return nil
	}

	m.detachLock.Lock()
	dq = m.detachQueues[link]
	if dq != nil {
		delete(m.detachQueues, link)
	}
	m.detachLock.Unlock()

	if dq == nil {
		return nil
	}
	return dq.Flush()
}

func (m *SlotMachineSync) ProcessDetachQueue(link SlotLink, worker DetachableSlotWorker) (hasSignal bool) {
	switch {
	case worker == nil:
		panic("illegal value")
	case worker.HasSignal():
		return true
	}

	tasks := m._flushDetachQueue(link)
	if len(tasks) == 0 {
		return false
	}

	return m.processCallbacks(tasks, worker)
}

func (m *SlotMachineSync) FlushSlotDetachQueue(link SlotLink) {
	detached := m._flushDetachQueue(link)
	m.callbackQueue.AddAll(detached)
}

func (m *SlotMachineSync) CleanupDetachQueues() bool {
	m.detachLock.Lock()
	defer m.detachLock.Unlock()

	isClean := true
	for link, dq := range m.detachQueues {
		if link.IsValid() {
			continue
		}
		delete(m.detachQueues, link)
		list := dq.Flush()
		m.callbackQueue.AddAll(list)
		if len(list) != 0 {
			isClean = false
		}
	}

	return isClean
}
