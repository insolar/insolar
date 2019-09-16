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
	"context"
	"github.com/insolar/insolar/conveyor/smachine/tools"
	"math"
	"sync"
	"sync/atomic"
	"time"
)

type SlotMachineConfig struct {
	//SyncStrategy    WorkSynchronizationStrategy
	PollingPeriod   time.Duration
	PollingTruncate time.Duration
	SlotPageSize    uint16
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
		injector: injector,
		//adapters:      adapters,
		//ownsAdapters:  ownsAdapters,
		slotPool:      NewSlotPool(config.SyncStrategy.NewSlotPoolLocker(), config.SlotPageSize),
		activeSlots:   NewSlotQueue(ActiveSlots),
		prioritySlots: NewSlotQueue(ActiveSlots),
		workingSlots:  NewSlotQueue(WorkingSlots),
		syncQueue:     NewSlotMachineSync(&sync.Mutex{}, config.SyncStrategy.GetInternalSignalCallback()),
	}
}

type SlotMachine struct {
	lastSlotID SlotID // atomic
	slotPool   SlotPool

	injector DependencyInjector

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

	syncQueue SlotMachineSync
}

func (m *SlotMachine) IsZero() bool {
	return m.syncQueue.IsZero()
}

func (m *SlotMachine) IsEmpty() bool {
	return m.slotPool.IsEmpty()
}

/* -------------------------------- */

func (m *SlotMachine) scanWorkingSlots(worker SlotWorker) {
	for {
		currentSlot := m.workingSlots.First()
		if currentSlot == nil {
			return
		}
		currentSlot.removeFromQueue()
		prevStepNo := currentSlot.startWorking(m.scanCount)

		stopNow := false
		var stateUpdate StateUpdate
		var asyncCount uint16

		// TODO consider use of sync.Pool for executionContext if they are allocated on heap

		wasDetached, err := worker.DetachableCall(func(workerCtx DetachableSlotWorker) {
			ec := executionContext{worker: workerCtx, slotContext: slotContext{s: currentSlot}}
			stopNow, stateUpdate, asyncCount = ec.executeNextStep()
		})

		if err != nil {
			stateUpdate = newStateUpdateTemplate(updCtxExec, 0, stateUpdPanic).newError(err)
		}

		if wasDetached {
			// MUST NOT apply any changes in the current routine, as it is no more safe to update queues
			slotLink := currentSlot.NewLink()
			m.applyDetachedStateUpdate(slotLink, stateUpdate, asyncCount, prevStepNo)
			return
		}

		if !stateUpdate.IsZero() {
			m.applyStateUpdate(currentSlot, stateUpdate, worker)
		}
		if currentSlot.stopWorking(asyncCount, prevStepNo) {
			// TODO apply migrate after initialization
		}
		if stopNow {
			return
		}

		// TODO if has async or barge-in then check for a detached queue
		if currentSlot.hasAsyncOrBargeIn() {
			// TODO process queues
		}

	}
}

/* -------------------------------- */

func (m *SlotMachine) recycleSlot(slot *Slot, w DetachableSlotWorker) {
	dependants := m._recycleSlot(slot)
	if dependants != nil {
		w.ActivateLinkedList(dependants, false)
	}
}

func (m *SlotMachine) recycleEmptySlot(slot *Slot) {
	if m._recycleSlot(slot) != nil {
		panic("illegal state")
	}
}

func (m *SlotMachine) _recycleSlot(slot *Slot) *Slot {
	dependants := m._cleanupSlot(slot)
	slot.dispose()

	m.slotPool.RecycleSlot(slot)
	return dependants
}

func (m *SlotMachine) _cleanupSlot(slot *Slot) *Slot {
	dep := slot.dependency
	slot.dependency = nil
	if dep != nil {
		dep.OnSlotDisposed()
	}

	if slot.isQueueHead() {
		return slot.removeHeadedQueue()
	} else {
		slot.removeFromQueue()
	}
	return nil
}

/* -------------------------------- */

func (m *SlotMachine) AddNew(ctx context.Context, parent SlotLink, sm StateMachine) SlotLink {
	link := m._addNewPrepare(ctx, parent, sm)
	m._startAddedSlot(link.s)
	return link
}

func (m *SlotMachine) AddNewAsync(ctx context.Context, parent SlotLink, sm StateMachine) SlotLink {
	link := m._addNewPrepare(ctx, parent, sm)
	// TODO async
	m._startAddedSlot(link.s)
	return link
}

func (m *SlotMachine) _addNewPrepare(ctx context.Context, parent SlotLink, sm StateMachine) SlotLink {
	if sm == nil {
		panic("illegal value")
	}
	if ctx == nil {
		panic("illegal value")
	}

	newSlot := m.allocateSlot()
	newSlot.parent = parent
	newSlot.ctx = ctx
	link := newSlot.NewLink()

	m.prepareNewSlot(newSlot, nil, nil, sm)
	return link
}

// TODO migrate MUST not take slots at step=0
// TODO slot must apply migrate after apply when step=0, but only outside of detachable
func (m *SlotMachine) prepareNewSlot(slot, creator *Slot, fn CreateFunc, sm StateMachine) bool {

	defer func() {
		recovered := recover()
		if recovered != nil {
			m.recycleEmptySlot(slot)
			panic(recovered)
		}
	}()

	if fn != nil {
		if sm != nil {
			panic("illegal value")
		}
		cc := constructionContext{s: slot}
		sm = cc.executeCreate(fn)
	}

	if sm == nil {
		return false
	}
	decl := sm.GetStateMachineDeclaration()
	if decl == nil {
		panic("illegal state")
	}
	slot.declaration = decl

	link := slot.NewLink()
	if !decl.InjectDependencies(sm, link, m, m.injector) && m.injector != nil {
		m.injector.InjectDependencies(sm, link, m)
	}

	initFn := slot.declaration.GetInitStateFor(sm)
	if initFn == nil {
		panic("illegal value")
	}

	if creator != nil {
		slot.migrationCount = creator.migrationCount
		slot.lastWorkScan = creator.lastWorkScan - 1
	} else {
		slot.migrationCount = m.migrationCount
		slot.lastWorkScan = uint8(m.scanCount - 1)
	}

	slot.step = SlotStep{Transition: func(ctx ExecutionContext) StateUpdate {
		ic := initializationContext{ctx.(*executionContext).clone()}
		slot.incStep()
		return ic.executeInitialization(initFn)
	}}

	return true
}

func (m *SlotMachine) startNewSlot(slot *Slot, worker DetachableSlotWorker) {
	slot.ensureInitializing()
	slot.stopWorking(0, 0)
	m.updateSlotQueue(slot, worker, activateSlot)
}

func (m *SlotMachine) _startAddedSlot(slot *Slot) {
	slot.ensureInitializing()
	slot.stopWorking(0, 0)
	list := m._updateSlotQueue(slot, false, activateSlot)
	if list != nil {
		panic("unexpected")
	}
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

func (m *SlotMachine) updateSlotQueue(slot *Slot, w DetachableSlotWorker, activation slotActivationMode) {
	if slot.machine != m {
		panic("illegal state")
	}
	w.EnsureMode(NonDetachableContext)
	linkedList := m._updateSlotQueue(slot, w.IsInplaceUpdate(), activation)
	if linkedList == nil {
		return
	}
	w.ActivateLinkedList(linkedList, activation == activateHotWaitSlot)
}

func (m *SlotMachine) _updateSlotQueue(slot *Slot, inplaceUpdate bool, activation slotActivationMode) *Slot {
	if !slot.isQueueHead() {
		if inplaceUpdate {
			switch activation {
			case activateSlot:
				switch slot.QueueType() {
				case ActiveSlots, WorkingSlots:
					return nil
				}
			case activateHotWaitSlot:
				if slot.QueueType() == ActiveSlots {
					return nil
				}
			}
			slot.removeFromQueue()
		} else {
			slot.ensureNotInQueue()
		}

		if activation == deactivateSlot {
			return nil
		}
		m._activateSlot(slot, activation)
		return nil
	}

	if slot.QueueType() != ActivationOfSlot {
		panic("illegal state")
	}

	if activation == deactivateSlot {
		if !inplaceUpdate {
			slot.ensureNotInQueue()
		}
		return nil
	}

	nextDep := slot.removeHeadedQueue()
	m._activateSlot(slot, activation)
	return nextDep
}

func (m *SlotMachine) _activateSlot(slot *Slot, mode slotActivationMode) {
	switch {
	case mode == activateHotWaitSlot:
		m._addSlotToActiveQueue(slot)
	case slot.isLastScan(m.scanCount):
		m.hotWaitOnly = false
		m._addSlotToActiveQueue(slot)
	default:
		m._addSlotToWorkingQueue(slot)
	}
}

func (m *SlotMachine) _addSlotToActiveQueue(slot *Slot) {
	if slot.isPriority() {
		m.prioritySlots.AddLast(slot)
	} else {
		m.activeSlots.AddLast(slot)
	}
}

func (m *SlotMachine) _addSlotToWorkingQueue(slot *Slot) {
	if slot.isPriority() {
		m.workingSlots.AddFirst(slot)
	} else {
		m.workingSlots.AddLast(slot)
	}
}

/* -------------------------------- */

func (m *SlotMachine) applyStateUpdate(slot *Slot, stateUpdate StateUpdate, w DetachableSlotWorker) bool {

	if slot.machine != m {
		panic("illegal state")
	}

	isAvailable, err := typeOfStateUpdate(stateUpdate).Apply(slot, stateUpdate, w)
	if err == nil {
		// TODO migrate zero step slots
		return isAvailable
	}

	m._handleStateUpdateError(slot, stateUpdate, w, err)
	m.recycleSlot(slot, w)
	return false
}

func (m *SlotMachine) _applyStateUpdate(slot *Slot, stateUpdate StateUpdate, w DetachableSlotWorker) (isAvailable bool, err error) {
	defer func() {
		err = recoverSlotPanic("failed to apply update", recover(), err)
	}()

	return typeOfStateUpdate(stateUpdate).Apply(slot, stateUpdate, w)
}

func (m *SlotMachine) applyDetachedStateUpdate(slotLink SlotLink, stateUpdate StateUpdate, asyncCount uint16, prevStepNo uint32,
) {
	m.syncQueue.AddSlotCall(slotLink, func(slot *Slot, w DetachableSlotWorker) {
		slot.stopWorking(asyncCount, prevStepNo)
		detachQueue := m.pullDetachQueue(slotLink.SlotID())

		if !stateUpdate.IsZero() && !m.applyStateUpdate(slot, stateUpdate, w) {
			if len(detachQueue) > 0 {
				m._handleMissedSlotCallback(slotLink, nil, detachQueue)
			}
			return
		}

		if _, isAvailable := m.migrateSlot(slot); !isAvailable {
			if len(detachQueue) > 0 {
				m._handleMissedSlotCallback(slotLink, nil, detachQueue)
			}
			return
		}

		for _, fn := range detachQueue {
			fn()
		}
	})

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

func (m *SlotMachine) _handleSlotAsyncPanic(link SlotLink, callFunc SlotDetachableFunc, e error) {
	// TODO
}

func (m *SlotMachine) _handleMissedSlotCallback(link SlotLink, callFunc SlotDetachableFunc, lists tools.SyncFuncList) {
	// TODO
}
