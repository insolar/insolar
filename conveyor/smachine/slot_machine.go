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
	"fmt"
	"github.com/insolar/insolar/conveyor/smachine/tools"
	"github.com/pkg/errors"
	"math"
	"sync"
	"sync/atomic"
	"time"
)

type SlotMachineConfig struct {
	SyncStrategy    WorkSynchronizationStrategy
	PollingPeriod   time.Duration
	PollingTruncate time.Duration
	SlotPageSize    uint16
}

type DependencyInjector interface {
	InjectDependencies(sm StateMachine, slotID SlotID, container *SlotMachine)
}

//WorkSynchronizationStrategy
func NewSlotMachine(config SlotMachineConfig, injector DependencyInjector, adapters *SharedRegistry) SlotMachine {
	ownsAdapters := false
	if adapters == nil {
		adapters = NewSharedRegistry()
		ownsAdapters = true
	}
	return SlotMachine{
		config:        config,
		injector:      injector,
		adapters:      adapters,
		ownsAdapters:  ownsAdapters,
		slotPool:      NewSlotPool(config.SyncStrategy.NewSlotPoolLocker(), config.SlotPageSize),
		activeSlots:   NewSlotQueue(ActiveSlots),
		prioritySlots: NewSlotQueue(ActiveSlots),
		workingSlots:  NewSlotQueue(WorkingSlots),
		syncQueue:     tools.NewSignalFuncQueue(&sync.Mutex{}, config.SyncStrategy.GetInternalSignalCallback()),
	}
}

type SlotMachine struct {
	config         SlotMachineConfig
	injector       DependencyInjector
	containerState SlotMachineState
	adapters       *SharedRegistry
	ownsAdapters   bool

	migrationCount uint16

	lastSlotID SlotID // atomic
	slotPool   SlotPool

	scanCount        uint32
	machineStartedAt time.Time
	scanStartedAt    time.Time
	scanWakeUpAt     time.Time

	workingSlots SlotQueue //slots are currently in processing

	hotWaitOnly   bool         // true when activeSlots has only slots added by "hot wait" / WaitAny
	activeSlots   SlotQueue    //they are are moved to workingSlots on every full Scan
	prioritySlots SlotQueue    //they are are moved to workingSlots on every partial or full Scan
	pollingSlots  PollingQueue //they are are moved to workingSlots on every full Scan when time has passed

	syncQueue    tools.SyncQueue // for detached/async ops, queued functions MUST BE panic-safe
	detachQueues map[SlotID]*tools.SyncQueue

	stepSync StepSyncCatalog
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

func (m *SlotMachine) GetAdapters() *SharedRegistry {
	return m.adapters
}

func (m *SlotMachine) OccupiedSlotCount() int {
	return m.slotPool.Count()
}

func (m *SlotMachine) AllocatedSlotCount() int {
	return m.slotPool.Capacity()
}

func (m *SlotMachine) IsZero() bool {
	return m.syncQueue.IsZero()
}

func (m *SlotMachine) IsEmpty() bool {
	return m.slotPool.IsEmpty()
}

func (m *SlotMachine) SetContainerState(s SlotMachineState) {
	m.containerState = s
}

func (m *SlotMachine) ScanOnceAsNested(context ExecutionContext) (repeatNow bool, nextPollTime time.Time) {
	worker := context.(*executionContext).worker.StartNested(m.containerState)
	defer worker.FinishNested(m.containerState)

	return m.ScanOnce(worker)
}

func (m *SlotMachine) ScanEventsOnly() (repeatNow bool, nextPollTime time.Time) {

	scanTime := time.Now()
	m.beforeScan(scanTime)

	repeatNow = m.scanEvents()

	return repeatNow, time.Time{}
}

func (m *SlotMachine) ProcessEventsOnly(scanTime time.Time) (waitForSignal bool, nextPollTime time.Time) {

	m.beforeScan(scanTime)

	m.pollingSlots.PrepareFor(scanTime.Add(m.config.PollingPeriod).Truncate(m.config.PollingTruncate))

	waitForSignal = !m.scanEvents()
	return waitForSignal, minTime(m.scanWakeUpAt, m.pollingSlots.GetNearestPollTime())
}

func (m *SlotMachine) ScanOnce(worker SlotWorker) (repeatNow bool, nextPollTime time.Time) {

	scanTime := time.Now()
	m.beforeScan(scanTime)

	switch {
	case m.machineStartedAt.IsZero():
		m.machineStartedAt = scanTime
		fallthrough
	case m.workingSlots.IsEmpty():
		m.scanCount++

		m.hotWaitOnly = true
		m.workingSlots.AppendAll(&m.prioritySlots)
		m.workingSlots.AppendAll(&m.activeSlots)
		m.pollingSlots.FilterOut(scanTime, m.workingSlots.AppendAll)
	default:
		// we were interrupted
		m.workingSlots.PrependAll(&m.prioritySlots)
	}

	m.pollingSlots.PrepareFor(scanTime.Add(m.config.PollingPeriod).Truncate(m.config.PollingTruncate))

	repeatNow = m.scanEvents()
	m.scanWorkingSlots(worker, scanTime)

	repeatNow = repeatNow || !m.hotWaitOnly
	return repeatNow, minTime(m.scanWakeUpAt, m.pollingSlots.GetNearestPollTime())
}

func minTime(t1, t2 time.Time) time.Time {
	if t1.IsZero() {
		return t2
	}
	if t2.IsZero() || t1.Before(t2) {
		return t1
	}
	return t2
}

func (m *SlotMachine) beforeScan(scanTime time.Time) {
	if m.machineStartedAt.IsZero() {
		m.machineStartedAt = scanTime
	}
	m.scanStartedAt = scanTime
	m.scanWakeUpAt = time.Time{}
}

func (m *SlotMachine) scanEvents() bool {

	syncQ := m.syncQueue.Flush()
	if len(syncQ) == 0 {
		return false
	}

	for _, fn := range syncQ {
		fn() // allows to resync detached
	}
	return true
}

func (m *SlotMachine) scanWorkingSlots(worker SlotWorker, scanStartTime time.Time) {
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

		wasDetached, err := worker.DetachableCall(func(workerCtx WorkerContext) {
			ec := executionContext{worker: workerCtx, slotContext: slotContext{s: currentSlot}}
			stopNow, stateUpdate, asyncCount = ec.executeNextStep()
		})

		if err != nil {
			stateUpdate = stateUpdatePanic(err)
		}

		if wasDetached {
			// MUST NOT apply any changes in the current routine, as it is no more considered as safe
			slotLink := currentSlot.NewLink()
			m.applyDetachedStateUpdate(slotLink, stateUpdate, asyncCount, prevStepNo)
			return
		}

		currentSlot.stopWorking(asyncCount, prevStepNo)

		if !stateUpdate.IsZero() {
			m.applyStateUpdate(currentSlot, false, stateUpdate)
		}
		if stopNow {
			return
		}
	}
}

func (m *SlotMachine) Migrate(cleanupWeak bool) {
	m.migrationCount++

	if m.ownsAdapters {
		m.adapters.migrate(m.containerState, m.migrationCount)
	}

	m.slotPool.ScanAndCleanup(cleanupWeak, m.disposeSlot, m.migratePage)
}

func (m *SlotMachine) Cleanup() {
	m.slotPool.ScanAndCleanup(true, m.disposeSlot, m.verifyPage)
}

func (m *SlotMachine) verifyPage(slotPage []Slot) (isPageEmptyOrWeak, hasWeakSlots bool) {
	isPageEmptyOrWeak = true
	hasWeakSlots = false

	for i := range slotPage {
		slot := &slotPage[i]

		switch {
		case slot.isEmpty():
			continue
		case slot.isWorking():
			break
		case slot.step.Flags&StepWeak != 0:
			hasWeakSlots = true
			continue
		}
		return false, hasWeakSlots
	}
	return isPageEmptyOrWeak, hasWeakSlots
}

func (m *SlotMachine) migratePage(slotPage []Slot) (isPageEmptyOrWeak, hasWeakSlots bool) {
	isPageEmptyOrWeak = true
	hasWeakSlots = false
	for i := range slotPage {
		isSlotEmptyOrWeak, isSlotAvailable := m.migrateSlot(&slotPage[i])
		switch {
		case !isSlotEmptyOrWeak:
			isPageEmptyOrWeak = false
		case isSlotAvailable:
			hasWeakSlots = true
		}
	}
	return isPageEmptyOrWeak, hasWeakSlots
}

func (m *SlotMachine) AddNew(ctx context.Context, parent SlotLink, sm StateMachine) SlotLink {
	if sm == nil {
		panic("illegal value")
	}
	if ctx == nil {
		panic("illegal value")
	}

	_, link := m.addStateMachine(ctx, nil, m.allocateNextSlotID(), parent, sm)
	return link
}

func (m *SlotMachine) AddAsyncNew(ctx context.Context, parent SlotLink, sm StateMachine) SlotID {
	if sm == nil || sm.GetStateMachineDeclaration() == nil {
		panic("illegal value")
	}
	if ctx == nil {
		panic("illegal value")
	}

	id := m.allocateNextSlotID()
	m.syncQueue.Add(func() {
		m._addStateMachine(ctx, nil, id, parent, sm)
	})

	// TODO return SlotLink
	return id
}

func (m *SlotMachine) addStateMachine(ctx context.Context, slot *Slot, newSlotID SlotID, parent SlotLink,
	sm StateMachine) (bool, SlotLink) {
	b, l := m._addStateMachine(ctx, slot, newSlotID, parent, sm)
	if l.IsEmpty() {
		panic("illegal state")
	}
	return b, l
}

func (m *SlotMachine) _addStateMachine(ctx context.Context, slot *Slot, newSlotID SlotID, parent SlotLink,
	sm StateMachine) (bool, SlotLink) {

	if sm == nil {
		return false, NoLink()
	}
	smd := sm.GetStateMachineDeclaration()
	if smd == nil {
		return false, NoLink()
	}

	if slot == nil {
		slot = m.allocateSlot()
	}
	m.prepareSlot(ctx, slot, newSlotID, parent, smd)
	link := slot.NewLink() // should happen BEFORE start as slot can die immediately
	if m.injector != nil {
		m.injector.InjectDependencies(sm, link.id, m)
	}

	return m.startSlot(slot, sm), link
}

func (m *SlotMachine) prepareSlot(ctx context.Context, slot *Slot, slotID SlotID, parent SlotLink, smd StateMachineDeclaration) {
	if smd == nil {
		panic("illegal state")
	}

	slot.init(ctx, slotID, parent, smd, m)
	slot.migrationCount = m.migrationCount
	slot.lastWorkScan = uint8(m.scanCount)
}

func (m *SlotMachine) startSlot(slot *Slot, sm StateMachine) bool {

	initState := slot.declaration.GetInitStateFor(sm)
	if initState == nil {
		panic("illegal state")
	}

	ic := initializationContext{slotContext{s: slot}}
	stateUpdate := ic.executeInitialization(initState)

	return m.applyStateUpdate(slot, false, stateUpdate)
}

func (m *SlotMachine) migrateSlot(slot *Slot) (isEmptyOrWeak, isAvailable bool) {
	if slot.isEmpty() {
		return true, false
	}

	if slot.isWorking() {
		return false, false
	}

	if m.migrationCount < slot.migrationCount {
		panic("illegal state")
	}

	for m.migrationCount != slot.migrationCount {
		migrateFn := slot.getMigration()
		if migrateFn == nil {
			slot.migrationCount = m.migrationCount
			break
		}

		mc := migrationContext{slotContext{s: slot}}
		stateUpdate := mc.executeMigrate(migrateFn)

		if !m.applyStateUpdate(slot, true, stateUpdate) {
			if !slot.isEmpty() {
				panic("illegal state")
			}
			return true, false // slot was stopped
		}

		slot.migrationCount++
	}

	return slot.step.Flags&StepWeak != 0, true
}

func (m *SlotMachine) allocateSlot() *Slot {
	return m.slotPool.AllocateSlot(m)
}

func (m *SlotMachine) slotAsyncCallback(slotLink SlotLink, fn func(*Slot)) {
	m.syncQueue.Add(func() {
		m.slotDirectCallback(slotLink, fn)
	})
}

func (m *SlotMachine) slotDirectCallback(slotLink SlotLink, fn func(*Slot)) {
	if !slotLink.IsValid() {
		detached := m.pullDetachQueue(slotLink.SlotID()) // cleanup

		m._handleMissedSlotCallback(slotLink, fn, detached)
		return
	}

	err := slotDirectCallbackSafe(slotLink.s, fn)
	if err == nil {
		return
	}

	stateUpdate := stateUpdatePanic(err)
	m.applyStateUpdate(slotLink.s, false, stateUpdate)
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

func slotDirectCallbackSafe(slot *Slot, fn func(*Slot)) (recovered error) {
	defer func() {
		recovered = recoverSlotPanic("async result has failed", recover(), recovered)
	}()
	fn(slot)
	return nil
}

func (m *SlotMachine) applyAsyncStateUpdate(link SlotLink, resultFn AsyncResultFunc, recovered interface{}) {

	m.slotAsyncCallback(link, func(slot *Slot) {
		if !slot.isWorking() {
			m._applyAsyncStateUpdate(slot, resultFn)
			return
		}
		/* this is an async result for a handler that was detached - we have to postpone it until reattachment */
		if m.detachQueues == nil {
			m.detachQueues = make(map[SlotID]*tools.SyncQueue)
		}

		dq := m.detachQueues[link.SlotID()]
		if dq == nil {
			dqs := tools.NewNoSyncQueue()
			dq = &dqs
		}

		dq.Add(func() {
			m.slotDirectCallback(link, func(slot *Slot) {
				m._applyAsyncStateUpdate(slot, resultFn)
			})
		})
		m.detachQueues[link.SlotID()] = dq
	})
}

func (m *SlotMachine) pullDetachQueue(slotID SlotID) tools.SyncFuncList {
	dq := m.detachQueues[slotID]
	if dq == nil {
		return nil
	}
	delete(m.detachQueues, slotID)
	return dq.Flush()
}

func (m *SlotMachine) _applyAsyncStateUpdate(slot *Slot, resultFn AsyncResultFunc) {
	slot.asyncCallCount--

	if resultFn == nil {
		return
	}
	rc := asyncResultContext{slot: slot}
	if !rc.executeResult(resultFn) {
		return
	}
	m._applyInplaceUpdate(slot, true, activateSlot)
}

func (m *SlotMachine) applyDetachedStateUpdate(slotLink SlotLink, stateUpdate StateUpdate,
	asyncCount uint16, prevStepNo uint32) {

	m.slotAsyncCallback(slotLink, func(slot *Slot) {
		slot.stopWorking(asyncCount, prevStepNo)
		detachQueue := m.pullDetachQueue(slotLink.SlotID())

		if !stateUpdate.IsZero() && !m.applyStateUpdate(slot, false, stateUpdate) {
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

func (m *SlotMachine) applyStateUpdate(slot *Slot, inplaceUpdate bool, stateUpdate StateUpdate) bool {

	if slot.machine != m {
		panic("illegal state")
	}

	isAvailable, recovered := m._applyStateUpdate(slot, inplaceUpdate, stateUpdate)
	if recovered == nil {
		return isAvailable
	}

	m._handleStateUpdateError(slot, inplaceUpdate, stateUpdate, recovered)
	m.disposeSlot(slot)
	m.slotPool.RecycleSlot(slot)
	return false
}

func (m *SlotMachine) _applyStateUpdate(slot *Slot, inplaceUpdate bool, stateUpdate StateUpdate) (isAvailable bool, errMsg interface{}) {

	if slot.isWorking() { // there must be no calls here from inside of custom handlers
		return false, "illegal state / slot is working"
	}

	defer func() {
		r := recover()
		if r == nil {
			return
		}
		var err error
		if errMsg != nil {
			err = fmt.Errorf("%s", errMsg)
		}
		errMsg = recoverSlotPanic("apply has failed", r, err)
	}()

	switch stateUpdType(stateUpdate.updType) {
	case stateUpdRepeat, stateUpdNextLoop:
		if stateUpdate.param1 != nil {
			return false, "unexpected param1"
		}
		fallthrough

	case stateUpdNext:
		applyNextStep(slot, stateUpdate)
		m._applyInplaceUpdate(slot, inplaceUpdate, activateSlot)
		return true, nil

	case stateUpdPoll:
		applyNextStep(slot, stateUpdate)
		m._applyInplaceUpdate(slot, inplaceUpdate, deactivateSlot)
		m.pollingSlots.Add(slot)
		return true, nil

	case stateUpdSleep:
		applyNextStep(slot, stateUpdate)
		m._applyInplaceUpdate(slot, inplaceUpdate, deactivateSlot)
		return true, nil

	case stateUpdWaitForEvent:
		applyNextStep(slot, stateUpdate)
		m._applyInplaceUpdate(slot, inplaceUpdate, activateHotWaitSlot)
		if stateUpdate.param0 > 0 {
			m.scanWakeUpAt = minTime(m.scanWakeUpAt, m.fromRelativeTime(stateUpdate.param0))
		}
		return true, nil

	case stateUpdWaitForActive:
		applyNextStep(slot, stateUpdate)
		waitOn := stateUpdate.getLink()
		switch {
		case waitOn.s == slot || !waitOn.IsValid():
			// don't wait for self
			// don't wait for an expired slot
			m._applyInplaceUpdate(slot, inplaceUpdate, activateSlot)
		default:
			switch waitOn.s.QueueType() {
			case ActiveSlots, WorkingSlots:
				// don't wait
				m._applyInplaceUpdate(slot, inplaceUpdate, activateSlot)
			case NoQueue:
				waitOn.s.makeQueueHead()
				fallthrough
			case ActivationOfSlot, PollingSlots:
				m._applyInplaceUpdate(slot, inplaceUpdate, deactivateSlot)
				waitOn.s.queue.AddLast(slot)
			default:
				return false, "illegal slot queue"
			}
		}
		return true, nil

	case stateUpdWaitForShared:
		//applyNextStep(slot, stateUpdate)
		panic("not implemented") // TODO

	case stateUpdReplace:
		// disposeSlot below can handle both in-place and off-place updates
		cf := stateUpdate.param1.(CreateFunc) // panic is handled
		parent := slot.parent
		ctx := slot.ctx
		m.disposeSlot(slot)
		ok, link := m.applySlotCreate(ctx, slot, parent, cf) // NB! recursive call inside
		if link.IsEmpty() {
			return false, "replacement was not created"
		}
		return ok, nil

	case stateUpdReplaceWith:
		sm := stateUpdate.param1.(StateMachine)
		if sm == nil {
			return false, "state machine was not provided"
		}
		smd := sm.GetStateMachineDeclaration()
		if smd == nil {
			return false, "state machine declaration is missing"
		}
		parent := slot.parent
		ctx := slot.ctx
		m.disposeSlot(slot)
		ok, link := m._addStateMachine(ctx, slot, m.allocateNextSlotID(), parent, sm)
		if link.IsEmpty() {
			return false, "replacement was not assigned"
		}
		return ok, nil

	case stateUpdStop:
		// disposeSlot can handle both in-place and off-place updates
		m.disposeSlot(slot)
		m.slotPool.RecycleSlot(slot)
		return false, nil

	case stateUpdNoChange:
		// only applicable for in-place updates
		if inplaceUpdate {
			return false, nil
		}
		return false, "unexpected state update"

	case stateUpdExpired:
		// can't be here
		return false, "unexpected state update"

	case stateUpdError, stateUpdPanic:
		return false, "error was reported"

	default:
		return false, "unknown update type"
	}
}

type slotActivationMode uint8

const (
	deactivateSlot slotActivationMode = iota
	activateSlot
	activateHotWaitSlot
)

func (m *SlotMachine) _applyInplaceUpdate(slot *Slot, inplaceUpdate bool, activation slotActivationMode) {

	if !slot.isQueueHead() {
		if inplaceUpdate {
			switch activation {
			case activateSlot:
				switch slot.QueueType() {
				case ActiveSlots, WorkingSlots:
					return
				}
			case activateHotWaitSlot:
				if slot.QueueType() == ActiveSlots {
					return
				}
			}
			slot.removeFromQueue()
		} else {
			slot.ensureNotInQueue()
		}

		if activation == deactivateSlot {
			return
		}
		m.activateSlot(slot, activation == activateHotWaitSlot)
		return
	}

	if slot.QueueType() != ActivationOfSlot {
		panic("illegal state")
	}

	if activation == deactivateSlot {
		if !inplaceUpdate {
			slot.ensureNotInQueue()
		}
		return
	}

	nextDep := slot.removeHeadedQueue()
	m.activateSlot(slot, activation == activateHotWaitSlot)
	m._reactivateDependencies(nextDep, activation == activateHotWaitSlot)
}

func applyNextStep(slot *Slot, stateUpdate StateUpdate) {
	slot.setNextStep(stateUpdate.step)
	if slot.step.Transition == nil {
		panic("missing transition")
	}
}

func (m *SlotMachine) createBargeIn(link StepLink, applyFn BargeInApplyFunc) BargeInParamFunc {
	return func(param interface{}) bool {
		if !link.IsValid() {
			return false
		}
		m.applyAsyncStateUpdate(link.SlotLink, func(ctx AsyncResultContext) {
			valid, atExactStep := link.isValidAndAtExactStep()
			if !valid {
				return
			}
			slot := link.s
			// it was not initiated as async call, so the counter needs adjustment
			slot.asyncCallCount++

			bc := bargingInContext{slotContext{s: slot}, param, atExactStep}
			stateUpdate := bc.executeBargeIn(applyFn)

			switch stateUpdType(stateUpdate.updType) {
			case stateUpdNoChange:
				return
			case stateUpdRepeat:
				// wakeup
				break
			case stateUpdNextLoop, stateUpdNext:
				slot.setNextStep(stateUpdate.step)
			default:
				panic("illegal value")
			}
			ctx.WakeUp()
		}, nil)
		return true
	}
}

func (m *SlotMachine) applySlotCreate(ctx context.Context, slot *Slot, parent SlotLink, fnCreate CreateFunc) (bool, SlotLink) {
	cc := constructionContext{ctx: ctx, parent: parent, slotID: m.allocateNextSlotID(), machine: m}
	sm := cc.executeCreate(fnCreate)
	return m._addStateMachine(cc.ctx, slot, cc.slotID, cc.parent, sm)
}

func (m *SlotMachine) activateSlot(slot *Slot, hotWait bool) {
	switch {
	case hotWait:
		m._addSlotToActiveQueue(slot)
	case slot.isLastScan(m.scanCount):
		m.hotWaitOnly = false
		m._addSlotToActiveQueue(slot)
	default:
		m._addSlotToWorkingQueue(slot)
	}
}

func (m *SlotMachine) activateDependedSlot(slot *Slot, addToActiveOnly bool) {
	if addToActiveOnly {
		m.hotWaitOnly = false
		m._addSlotToActiveQueue(slot)
	} else {
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

func (m *SlotMachine) disposeSlot(slot *Slot) {
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

	slot.dispose()
}

func (m *SlotMachine) _reactivateDependencies(current *Slot, activeWait bool) {

	addToActive := activeWait || current.isLastScan(m.scanCount)

	for current != nil {
		next := current.nextInQueue
		current.nextInQueue = nil
		current.prevInQueue = nil

		nm := next.machine
		if m == nm {
			m.activateDependedSlot(next, addToActive)
		} else {
			// TODO inefficient for multiple dependencies
			nm.syncQueue.Add(func() {
				nm.activateDependedSlot(next, addToActive)
			})
		}
	}
}

func (m *SlotMachine) _handleMissedSlotCallback(link SlotLink, missed func(*Slot), detached tools.SyncFuncList) {
	// TODO logging

	if missed == nil {
		fmt.Printf("callback(s) on expired slot: count=%d, %v", len(detached), detached)
	} else {
		fmt.Printf("callback(s) on expired slot: count=%d, %p, %v", 1+len(detached), missed, detached)
	}
}

func (m *SlotMachine) _handleStateUpdateError(slot *Slot, inplaceUpdate bool, stateUpdate StateUpdate, recovered interface{}) {
	var slotError error

	isAsync := inplaceUpdate
	isPanic := true
	switch msg := recovered.(type) {
	case string:
		switch stateUpdType(stateUpdate.updType) {
		case stateUpdError:
			isPanic = false
			fallthrough
		case stateUpdPanic:
			switch p := stateUpdate.param1.(type) {
			case error:
				slotError = p
			default:
				slotError = fmt.Errorf("unknown error/panic payload: %v", p)
			}
		default:
			slotError = fmt.Errorf("%s: update=%v", msg, stateUpdate)
		}
	case error:
		slotError = errors.WithMessage(msg, fmt.Sprintf("internal error: update=%v", stateUpdate))
	default:
		slotError = fmt.Errorf("internal error: msg=%v update=%v", msg, stateUpdate)
	}

	// TODO log error

	handler := slot.getErrorHandler()
	if handler != nil {
		ec := failureContext{slotContext{s: slot}, isPanic, isAsync, slotError}
		err := ec.executeErrorHandlerSafe(handler)
		if err == nil {
			return
		}
		slotError = errors.WithMessage(slotError, err.Error())
	}

	// TODO log unhandled error
	fmt.Printf("slot fail: %v %+v %v", slot.GetID(), stateUpdate, slotError)
}

func (m *SlotMachine) toRelativeTime(t time.Time) uint32 {

	if m.scanStartedAt.IsZero() {
		panic("illegal state")
	}

	if t.IsZero() {
		return 0
	}

	d := t.Sub(m.scanStartedAt)
	if d <= 0 {
		return 1
	}

	d /= time.Microsecond
	if d == 0 {
		return 1
	}
	d++
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
