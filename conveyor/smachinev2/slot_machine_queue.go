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
	"fmt"
	"github.com/insolar/insolar/conveyor/smachine/tools"
	"github.com/pkg/errors"
	"sync"
)

func NewSlotMachineSync(locker sync.Locker, signalCallback func()) SlotMachineSync {
	return SlotMachineSync{
		locker:        locker,
		updateQueue:   tools.NewSignalFuncQueue(locker, signalCallback),
		callbackQueue: tools.NewSignalFuncQueue(locker, signalCallback),
	}
}

type SlotMachineSync struct {
	locker sync.Locker

	updateQueue   tools.SyncQueue // func(w SlotWorker) // for detached/async ops, queued functions MUST BE panic-safe
	callbackQueue tools.SyncQueue // func(w DetachableSlotWorker) // for detached/async ops, queued functions MUST BE panic-safe

	detachQueues map[SlotID]*tools.SyncQueue
}

func (m *SlotMachineSync) IsZero() bool {
	return m.locker == nil
}

func (m *SlotMachineSync) AddAsyncUpdate(link SlotLink, fn func(link SlotLink, worker SlotWorker)) {
	if fn == nil {
		panic("illegal value")
	}

	m.updateQueue.Add(func(w interface{}) {
		fn(link, w.(SlotWorker))
	})
}

func (m *SlotMachineSync) ProcessUpdates(worker SlotWorker) bool {
	worker.EnsureMode(NonDetachableContext)

	tasks := m.updateQueue.Flush()
	if len(tasks) == 0 {
		return false
	}

	for _, fn := range tasks {
		fn(worker)
	}
	return true
}

func (m *SlotMachineSync) AddAsyncCallback(link SlotLink, fn func(link SlotLink, worker DetachableSlotWorker)) {
	if fn == nil {
		panic("illegal value")
	}

	m.callbackQueue.Add(func(w interface{}) {
		fn(link, w.(DetachableSlotWorker))
	})
}

func (m *SlotMachineSync) ProcessCallbacks(worker SlotWorker) (hasSignal, wasDetached bool) {

	if worker.HasSignal() {
		return true, false
	}

	tasks := m.callbackQueue.Flush()
	if len(tasks) == 0 {
		return false, false
	}

	for i, fn := range tasks {
		wasDetached = worker.DetachableCall(func(w DetachableSlotWorker) {
			fn(w)
		})

		if worker.HasSignal() {
			m.callbackQueue.AddAll(tasks[i+1:])
			return true, wasDetached
		}
	}
	return false, false
}

func (m *SlotMachineSync) flushDetachQueue(slotID SlotID) tools.SyncFuncList {
	dq := m.detachQueues[slotID]
	if dq == nil {
		return nil
	}
	delete(m.detachQueues, slotID)
	return dq.Flush()
}

func (m *SlotMachineSync) _slotCall(slotLink SlotLink, fn SlotDetachableFunc) {
	if !slotLink.IsValid() {
		detached := m.flushDetachQueue(slotLink.SlotID()) // cleanup
		if slotLink.s == nil || slotLink.s.machine == nil {
			return
		}

		slotLink.s.machine._handleMissedSlotCallback(slotLink, fn, detached)
		return
	}

	err := m._slotSafeCall(slotLink.s, fn)
	if err == nil {
		return
	}
	slotLink.s.machine._handleSlotAsyncPanic(slotLink, fn, err)
}

func recoverSlotPanic(msg string, recovered interface{}, prev error) error {
	if recovered == nil {
		return prev
	}
	if prev != nil {
		return errors.Wrap(prev, fmt.Sprintf("%s: %v", msg, recovered))
	}
	return errors.Errorf("%s: %v", msg, recovered)
}

func (m *SlotMachineSync) _slotSafeCall(slot *Slot, fn SlotDetachableFunc) (recovered error) {
	defer func() {
		recovered = recoverSlotPanic("async result has failed", recover(), recovered)
	}()
	return m.getWorker().TryExclusiveSlotCall(slot, fn), nil
}

func (m *SlotMachineSync) AppendSlotDetachQueue(id SlotID) {
	detached := m.flushDetachQueue(id)
	if len(detached) > 0 {
		m.callbackQueue.AddAll(detached)
	}
}
