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
	"math"
	"sync/atomic"
	"time"
)

type SlotMachineConfig struct {
	//SyncStrategy    WorkSynchronizationStrategy
	PollingPeriod   time.Duration
	PollingTruncate time.Duration
	SlotPageSize    uint16
}

type DependencyInjector interface {
	InjectDependencies(sm StateMachine, slotID SlotID, container *SlotMachine)
}

func NewSlotMachine(config SlotMachineConfig, injector DependencyInjector, // adapters *SharedRegistry
) SlotMachine {
	//ownsAdapters := false
	//if adapters == nil {
	//	adapters = NewSharedRegistry()
	//	ownsAdapters = true
	//}
	return SlotMachine{
		//config:        config,
		//injector:      injector,
		//adapters:      adapters,
		//ownsAdapters:  ownsAdapters,
		//slotPool:      NewSlotPool(config.SyncStrategy.NewSlotPoolLocker(), config.SlotPageSize),
		//activeSlots:   NewSlotQueue(ActiveSlots),
		//prioritySlots: NewSlotQueue(ActiveSlots),
		//workingSlots:  NewSlotQueue(WorkingSlots),
		//syncQueue:     tools.NewSignalFuncQueue(&sync.Mutex{}, config.SyncStrategy.GetInternalSignalCallback()),
	}
}

type SlotMachine struct {
	lastSlotID SlotID // atomic
	slotPool   SlotPool

	scanCount        uint32
	machineStartedAt time.Time
	scanStartedAt    time.Time

	scanWakeUpAt time.Time

	migrationCount uint16

	hotWaitOnly   bool         // true when activeSlots has only slots added by "hot wait" / WaitAny
	activeSlots   SlotQueue    //they are are moved to workingSlots on every full Scan
	prioritySlots SlotQueue    //they are are moved to workingSlots on every partial or full Scan
	pollingSlots  PollingQueue //they are are moved to workingSlots on every full Scan when time has passed
	workingSlots  SlotQueue    //slots are currently in processing

	syncQueue    tools.SyncQueue // for detached/async ops, queued functions MUST BE panic-safe
	detachQueues map[SlotID]*tools.SyncQueue
}

func (m *SlotMachine) IsZero() bool {
	return m.syncQueue.IsZero()
}

func (m *SlotMachine) IsEmpty() bool {
	return m.slotPool.IsEmpty()
}

func (m *SlotMachine) allocateNextSlotID() SlotID {
	for {
		r := atomic.LoadUint32((*uint32)(&m.lastSlotID))
		if r == math.MaxUint32 {
			panic("overflow")
		}
		if atomic.CompareAndSwapUint32((*uint32)(&m.lastSlotID), r, r+1) {
			return SlotID(r + 1)
		}
	}
}

func (m *SlotMachine) allocateSlot() *Slot {
	return m.slotPool.AllocateSlot(m, m.allocateNextSlotID())
}

/* -------------------------------- */

type slotActivationMode uint8

const (
	deactivateSlot slotActivationMode = iota
	activateSlot
	activateHotWaitSlot
)

func (m *SlotMachine) disposeSlot(slot *Slot, worker WorkerContext) {
	m._cleanupSlot(slot, worker)
	slot.dispose(UnknownSlotID)
}

func (m *SlotMachine) reuseSlot(slot *Slot, reuseFor SlotID, worker WorkerContext) {
	if reuseFor.IsUnknown() {
		panic("illegal state")
	}
	m._cleanupSlot(slot, worker)
	slot.dispose(reuseFor)
}

func (m *SlotMachine) _cleanupSlot(slot *Slot, worker WorkerContext) {
	dep := slot.dependency
	slot.dependency = nil
	if dep != nil {
		dep.OnSlotDisposed()
	}

	if slot.isQueueHead() {
		dependencies := slot.removeHeadedQueue()
		m._reactivateDependencies(dependencies, false)
	} else {
		slot.removeFromQueue()
	}
}

func (m *SlotMachine) updateSlotQueue(slot *Slot, context WorkerContext, mode slotActivationMode) {
	zzz
}

// TODO migrate MUST not take slots at step=0
// TODO slot must apply migrate after apply when step=0, but only outside of detachable
func (m *SlotMachine) prepareNewSlot(slot, creator *Slot, fn CreateFunc, sm StateMachine) {

	defer func() {
		recovered := recover()
		if recovered == nil {
			return
		}
		m._cleanupSlot(slot, nil)
		m.slotPool.RecycleSlot(slot)
		panic(recovered)
	}()

	if fn != nil {
		if sm != nil {
			panic("illegal value")
		}
		cc := constructionContext{s: slot}
		sm = cc.executeCreate(fn)
	}

	initFn := slot.declaration.GetInitStateFor(sm)
	if initFn == nil {
		panic("illegal value")
	}
	slot.migrationCount = creator.migrationCount
	slot.lastWorkScan = creator.lastWorkScan - 1

	slot.step = SlotStep{Transition: func(ctx ExecutionContext) StateUpdate {
		slot.incStep()
		exec := ctx.(*executionContext)
		ic := initializationContext{exec.clone()}
		return ic.executeInitialization(initFn)
	}}
}

func (m *SlotMachine) startNewSlot(slot *Slot, worker WorkerContext) {
	slot.ensureInitializing()
	m.updateSlotQueue(slot, worker, activateSlot)
	slot.stopWorking(0, 0)
}

/* -------------------------------- */

func (m *SlotMachine) createBargeIn(link StepLink, applyFn BargeInApplyFunc) BargeInParamFunc {
	//return func(param interface{}) bool {
	//	if !link.IsValid() {
	//		return false
	//	}
	//	m.applyAsyncStateUpdate(link.SlotLink, func(ctx AsyncResultContext) {
	//		valid, atExactStep := link.isValidAndAtExactStep()
	//		if !valid {
	//			return
	//		}
	//		slot := link.s
	//		// it was not initiated as async call, so the counter needs adjustment
	//		slot.asyncCallCount++
	//
	//		bc := bargingInContext{slotContext{s: slot}, param, atExactStep}
	//		stateUpdate := bc.executeBargeIn(applyFn)
	//
	//		switch stateUpdType(stateUpdate.updType) {
	//		case stateUpdNoChange:
	//			return
	//		case stateUpdRepeat:
	//			// wakeup
	//			break
	//		case stateUpdNextLoop, stateUpdNext:
	//			slot.setNextStep(stateUpdate.step)
	//		default:
	//			panic("illegal value")
	//		}
	//		ctx.WakeUp()
	//	}, nil)
	//	return true
	//}
}

/* -------------------------------- */

func minTime(t1, t2 time.Time) time.Time {
	if t1.IsZero() {
		return t2
	}
	if t2.IsZero() || t1.Before(t2) {
		return t1
	}
	return t2
}

func (m *SlotMachine) toRelativeTime(t time.Time) uint32 {
	if m.scanStartedAt.IsZero() {
		panic("illegal state")
	}
	if t.IsZero() {
		return 0
	}

	d := t.Sub(m.scanStartedAt)
	if d <= time.Microsecond {
		return 1
	}
	d = 1 + d/time.Microsecond
	if d > math.MaxUint32 {
		panic("illegal value")
	}
	return uint32(d)
}

func (m *SlotMachine) fromRelativeTime(rel uint32) time.Time {
	switch rel {
	case 0:
		return time.Time{}
	case 1:
		return m.scanStartedAt
	}
	return m.scanStartedAt.Add(time.Duration(rel-1) * time.Microsecond)
}
