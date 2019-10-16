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
	"github.com/insolar/insolar/conveyor/injector"
	"math"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

type MigrationFunc func(migrationCount uint32)

type SlotMachineConfig struct {
	PollingPeriod        time.Duration
	PollingTruncate      time.Duration
	SlotPageSize         uint16
	ScanCountLimit       int
	CleanupWeakOnMigrate bool
}

const maxLoopCount = 10000

func NewSlotMachine(config SlotMachineConfig,
	eventCallback, signalCallback func(),
	parentRegistry injector.DependencyRegistry,
) SlotMachine {
	return SlotMachine{
		config:         config,
		parentRegistry: parentRegistry,
		slotPool:       NewSlotPool(config.SlotPageSize, false),
		activeSlots:    NewSlotQueue(ActiveSlots),
		prioritySlots:  NewSlotQueue(ActiveSlots),
		workingSlots:   NewSlotQueue(WorkingSlots),
		syncQueue:      NewSlotMachineSync(eventCallback, signalCallback),
	}
}

var _ injector.DependencyRegistry = &SlotMachine{}

type SlotMachine struct {
	config     SlotMachineConfig
	lastSlotID SlotID // atomic
	slotPool   SlotPool

	parentRegistry injector.DependencyRegistry
	localRegistry  sync.Map // is used for both dependencies and tracking of dependencies

	scanCount        uint32
	machineStartedAt time.Time
	scanStartedAt    time.Time

	migrators      []MigrationFunc
	migrationCount uint32

	hotWaitOnly  bool      // true when activeSlots & prioritySlots have only slots added by "hot wait"
	scanWakeUpAt time.Time // when all slots are waiting, this is the earliest time requested for wakeup

	activeSlots   SlotQueue    //they are are moved to workingSlots on every full Scan
	prioritySlots SlotQueue    //they are are moved to workingSlots on every full Scan (placed first)
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

func (m *SlotMachine) IsActive() bool {
	return m.syncQueue.IsActive()
}

func (m *SlotMachine) Stop() bool {
	return m.syncQueue.SetStopping()
}

/* ------- Methods for dependency injections - safe for concurrent use ------------- */

type dependencyKey string // is applied to avoid key interference with aliases

func (m *SlotMachine) FindDependency(id string) (interface{}, bool) {
	if v, ok := m.localRegistry.Load(dependencyKey(id)); ok {
		return v, true
	}
	if m.parentRegistry != nil {
		return m.parentRegistry.FindDependency(id)
	}
	return nil, false
}

func (m *SlotMachine) PutDependency(id string, v interface{}) {
	if id == "" {
		panic("illegal key")
	}
	m.localRegistry.Store(dependencyKey(id), v)
}

func (m *SlotMachine) TryPutDependency(id string, v interface{}) bool {
	if id == "" {
		panic("illegal key")
	}
	_, loaded := m.localRegistry.LoadOrStore(dependencyKey(id), v)
	return !loaded
}

/* -------------- Methods to run state machines --------------- */

func (m *SlotMachine) ScanOnce(worker AttachedSlotWorker) (repeatNow bool, nextPollTime time.Time) {
	status := m.syncQueue.GetStatus()
	if status == SlotMachineInactive {
		return false, time.Time{}
	}

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
	}
	m.pollingSlots.PrepareFor(scanTime.Add(m.config.PollingPeriod).Truncate(m.config.PollingTruncate))

	if status == SlotMachineStopping {
		return m.stopAll(worker), time.Time{}
	}

	repeatNow = m.syncQueue.ProcessUpdates(worker)
	hasUpdates, hasSignal, wasDetached := m.syncQueue.ProcessCallbacks(worker)
	if hasUpdates {
		repeatNow = true
	}

	if !hasSignal && !wasDetached {
		m.executeWorkingSlots(worker)
	}

	repeatNow = repeatNow || !m.hotWaitOnly
	return repeatNow, minTime(m.scanWakeUpAt, m.pollingSlots.GetNearestPollTime())
}

func (m *SlotMachine) beforeScan(scanTime time.Time) {
	if m.machineStartedAt.IsZero() {
		m.machineStartedAt = scanTime
	}
	m.scanStartedAt = scanTime
	m.scanWakeUpAt = time.Time{}
}

func (m *SlotMachine) stopAll(worker AttachedSlotWorker) (repeatNow bool) {
	clean := m.slotPool.ScanAndCleanup(true, worker, m.recycleSlot, m.stopPage)
	hasUpdates := m.syncQueue.ProcessUpdates(worker)
	hasCallbacks, _, _ := m.syncQueue.ProcessCallbacks(worker)

	if hasUpdates || hasCallbacks || !clean || !m.syncQueue.CleanupDetachQueues() || !m.slotPool.IsEmpty() {
		return true
	}

	m.syncQueue.SetInactive()
	return false
}

func (m *SlotMachine) executeWorkingSlots(worker AttachedSlotWorker) {
	limit := m.config.ScanCountLimit
	if limit <= 0 || limit > maxLoopCount {
		limit = maxLoopCount
	}

	for i := 0; i < limit; i++ {
		currentSlot := m.workingSlots.First()
		if currentSlot == nil {
			return
		}
		loopLimit := 1 + ((limit - i) / m.workingSlots.Count())

		prevStepNo := currentSlot.startWorking(m.scanCount) // its counterpart is in slotPostExecution()
		currentSlot.removeFromQueue()

		if stopNow, loopExtraIncrement := m._executeSlot(currentSlot, prevStepNo, worker, loopLimit); stopNow {
			return
		} else {
			i += loopExtraIncrement
		}
	}
}

func (m *SlotMachine) _executeSlot(slot *Slot, prevStepNo uint32, worker AttachedSlotWorker, loopLimit int) (hasSignal bool, loopCount int) {

	if dep := slot.dependency; dep != nil && dep.IsReleaseOnWorking() {
		slot.dependency = nil
		m.activateDependants(dep.Release(), worker)
	}
	slot.slotFlags &^= slotWokenUp

	// TODO consider use of sync.Pool for executionContext if they are allocated on heap
	var stateUpdate StateUpdate
	wasDetached := worker.DetachableCall(func(worker DetachableSlotWorker) {
		//defer func() {
		//	// kill slot on fail
		//}()

		for ; loopCount < loopLimit; loopCount++ {
			canLoop := false
			canLoop, hasSignal = worker.CanLoopOrHasSignal(loopCount)
			if !canLoop || hasSignal {
				if loopCount == 0 {
					// a very special update type, not to be used anywhere else
					stateUpdate = StateUpdate{updKind: uint8(stateUpdInternalRepeatNow)}
				} else {
					stateUpdate = newStateUpdateTemplate(updCtxExec, 0, stateUpdRepeat).newUint(0)
				}
				return
			}

			var asyncCnt uint16
			var sut StateUpdateType

			ec := executionContext{slotContext: slotContext{s: slot, w: worker}}
			stateUpdate, sut, asyncCnt = ec.executeNextStep()

			slot.addAsyncCount(asyncCnt)
			if !sut.ShortLoop(slot, stateUpdate, uint32(loopCount)) {
				return
			}
		}
	})

	if wasDetached {
		// MUST NOT apply any changes in the current routine, as it is no more safe to update queues
		m.asyncPostSlotExecution(slot, stateUpdate, prevStepNo)
		return true, loopCount
	}

	hasAsync := m.slotPostExecution(slot, stateUpdate, worker, prevStepNo, false)
	if hasAsync && !hasSignal {
		_, hasSignal, wasDetached = m.syncQueue.ProcessCallbacks(worker)
		return hasSignal || wasDetached, loopCount
	}
	return hasSignal, loopCount
}

func (m *SlotMachine) _executeSlotInitByCreator(currentSlot *Slot, worker DetachableSlotWorker) {

	currentSlot.ensureInitializing()

	ec := executionContext{slotContext: slotContext{s: currentSlot, w: worker}}
	stateUpdate, _, asyncCnt := ec.executeNextStep()

	currentSlot.addAsyncCount(asyncCnt)
	if !worker.NonDetachableCall(func(worker FixedSlotWorker) {
		m.slotPostExecution(currentSlot, stateUpdate, worker, 0, false)
	}) {
		m.asyncPostSlotExecution(currentSlot, stateUpdate, 0)
	}
}

func (m *SlotMachine) slotPostExecution(slot *Slot, stateUpdate StateUpdate, worker FixedSlotWorker,
	prevStepNo uint32, wasAsync bool) (hasAsync bool) {

	if !stateUpdate.IsZero() && !m.applyStateUpdate(slot, stateUpdate, worker) {
		return false
	}

	if slot.canMigrateWorking(prevStepNo, wasAsync) {
		if _, isAvailable := m._migrateSlot(slot, worker); !isAvailable {
			return false
		}
	}

	hasAsync = wasAsync || slot.hasAsyncOrBargeIn()
	m.stopSlotWorking(slot, prevStepNo, worker)
	return hasAsync
}

func (m *SlotMachine) queueAsyncCallback(link SlotLink,
	callbackFn func(slot *Slot, worker DetachableSlotWorker) StateUpdate, prevErr error) {

	if callbackFn == nil && prevErr == nil || !m._canCallback(link) {
		return
	}

	m.syncQueue.AddAsyncCallback(link, func(link SlotLink, worker DetachableSlotWorker) (isDone bool) {
		if !m._canCallback(link) {
			return true
		}
		if worker == nil {
			// TODO _handleAsyncDetachmentLimitExceeded
			return true
		}

		slot, isStarted, prevStepNo := link.tryStartWorking()
		if !isStarted {
			return false
		}
		var stateUpdate StateUpdate
		func() {
			defer func() {
				recoverSlotPanicAsUpdate(&stateUpdate, "async callback panic", recover(), prevErr)
			}()
			if callbackFn != nil {
				stateUpdate = callbackFn(slot, worker)
			}
		}()

		if worker.NonDetachableCall(func(worker FixedSlotWorker) {
			m.slotPostExecution(slot, stateUpdate, worker, prevStepNo, true)
		}) {
			m.syncQueue.ProcessDetachQueue(link, worker)
		} else {
			m.asyncPostSlotExecution(slot, stateUpdate, prevStepNo)
		}

		return true
	})
}

func (m *SlotMachine) _canCallback(link SlotLink) bool {
	if link.s.machine != m {
		panic("illegal state")
	}
	if link.IsValid() {
		return true
	}
	return false
}

func (m *SlotMachine) asyncPostSlotExecution(s *Slot, stateUpdate StateUpdate, prevStepNo uint32) {
	m.syncQueue.AddAsyncUpdate(s.NewLink(), func(link SlotLink, worker FixedSlotWorker) {
		if !link.IsValid() {
			return
		}
		slot := link.s
		if m.slotPostExecution(slot, stateUpdate, worker, prevStepNo, true) {
			m.syncQueue.FlushSlotDetachQueue(link)
		}
	})
}

/* -- Methods to migrate slots ------------------------------ */

func (m *SlotMachine) migrate(worker FixedSlotWorker) {
	m.migrationCount++

	for _, fn := range m.migrators {
		fn(m.migrationCount)
	}

	m.slotPool.ScanAndCleanup(m.config.CleanupWeakOnMigrate, worker, m.recycleSlot, m.migratePage)
	m.syncQueue.CleanupDetachQueues()
}

func (m *SlotMachine) AddMigrationCallback(fn MigrationFunc) {
	if fn == nil {
		panic("illegal value")
	}
	m.migrators = append(m.migrators, fn)
}

func (m *SlotMachine) migratePage(slotPage []Slot, worker FixedSlotWorker) (isPageEmptyOrWeak, hasWeakSlots bool) {
	isPageEmptyOrWeak = true
	hasWeakSlots = false
	for i := range slotPage {
		isSlotEmptyOrWeak, isSlotAvailable := m.migrateSlot(&slotPage[i], worker)
		switch {
		case !isSlotEmptyOrWeak:
			isPageEmptyOrWeak = false
		case isSlotAvailable:
			hasWeakSlots = true
		}
	}
	return isPageEmptyOrWeak, hasWeakSlots
}

func (m *SlotMachine) migrateSlot(slot *Slot, w FixedSlotWorker) (isEmptyOrWeak, isAvailable bool) {
	isEmpty, isStarted, prevStepNo := slot.tryStartMigrate()
	if !isStarted {
		return isEmpty, false
	}
	isEmptyOrWeak, isAvailable = m._migrateSlot(slot, w)
	if isAvailable {
		m.stopSlotWorking(slot, prevStepNo, w)
	}
	return isEmptyOrWeak, isAvailable
}

func (m *SlotMachine) _migrateSlot(slot *Slot, worker FixedSlotWorker) (isEmptyOrWeak, isAvailable bool) {

	for m.migrationCount != slot.migrationCount {
		migrateFn := slot.getMigration()
		if migrateFn == nil {
			slot.migrationCount = m.migrationCount
			break
		}

		mc := migrationContext{slotContext: slotContext{s: slot}}
		stateUpdate, skipMultiple := mc.executeMigration(migrateFn)

		if !m.applyStateUpdate(slot, stateUpdate, worker) {
			return true, false
		}

		if skipMultiple {
			slot.migrationCount = m.migrationCount
			break
		}
		slot.migrationCount++
	}

	return slot.step.Flags&StepWeak != 0, true
}

/* -- Methods to allocate slots ------------------------------ */

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

/* -- Methods to dispose/reuse slots ------------------------------ */

func (m *SlotMachine) Cleanup(worker FixedSlotWorker) {
	m.slotPool.ScanAndCleanup(true, worker, m.recycleSlot, m.verifyPage)
	m.syncQueue.CleanupDetachQueues()
}

func (m *SlotMachine) verifyPage(slotPage []Slot, _ FixedSlotWorker) (isPageEmptyOrWeak, hasWeakSlots bool) {
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

func (m *SlotMachine) stopPage(slotPage []Slot, w FixedSlotWorker) (isPageEmptyOrWeak, hasWeakSlots bool) {
	hasWorking := false

	for i := range slotPage {
		slot := &slotPage[i]

		switch {
		case slot.isEmpty():
			//continue
		case slot.isWorking():
			hasWorking = true
		default:
			m.recycleSlot(slot, w)
		}
	}
	return !hasWorking, false
}

func (m *SlotMachine) recycleSlot(slot *Slot, w FixedSlotWorker) {
	m._recycleSlot(slot, w)
}

func (m *SlotMachine) recycleEmptySlot(slot *Slot) {
	m._recycleSlot(slot, nil)
}

func (m *SlotMachine) _recycleSlot(slot *Slot, worker FixedSlotWorker) {

	if slot.slotFlags&slotHasAliases != 0 {
		// cleanup aliases associated with the slot
		slot.unregisterBoundAliases()
	}

	{ // cleanup synchronization dependency
		dep := slot.dependency
		if dep != nil {
			slot.dependency = nil
			m.activateDependants(dep.Release(), worker)
		}
	}

	{ // cleanup queues
		if slot.isQueueHead() {
			s := slot.removeHeadedQueue()
			m._activateDependantChain(s, worker)
		} else {
			slot.removeFromQueue()
		}
	}

	if slot.slotFlags&(slotHadAsync|slotHasBargeIn|slotHasAliases) != 0 {
		link := slot.NewLink()
		slot.dispose() // check state and cleanup fields
		m.syncQueue.FlushSlotDetachQueue(link)
	} else {
		slot.dispose() // check state and cleanup fields
	}

	m.slotPool.RecycleSlot(slot)
}

func (m *SlotMachine) OccupiedSlotCount() int {
	return m.slotPool.Count()
}

func (m *SlotMachine) AllocatedSlotCount() int {
	return m.slotPool.Capacity()
}

/* -- General purpose synchronization ------------------------------ */

func (m *SlotMachine) ScheduleCall(fn MachineCallFunc, isSignal bool) {
	if fn == nil {
		panic("illegal value")
	}
	callFn := func(_ SlotLink, worker FixedSlotWorker) {
		mc := machineCallContext{m: m, w: worker}
		err := mc.executeCall(fn)
		if err != nil {
			// TODO log call error
			runtime.KeepAlive(err)
		}
	}
	if isSignal {
		m.syncQueue.AddAsyncSignal(SlotLink{}, callFn)
	} else {
		m.syncQueue.AddAsyncUpdate(SlotLink{}, callFn)
	}
}

/* -- Methods to create and start new machines ------------------------------ */

func (m *SlotMachine) AddNew(ctx context.Context, parent SlotLink, sm StateMachine) SlotLink {
	link, ok := m._addNew(ctx, parent, sm)
	if ok {
		m.syncQueue.AddAsyncUpdate(link, m._startAddedSlot)
	}
	return link
}

func (m *SlotMachine) _addNew(ctx context.Context, parent SlotLink, sm StateMachine) (SlotLink, bool) {
	if sm == nil {
		panic("illegal value")
	}
	if ctx == nil {
		panic("illegal value")
	}

	if !m.IsActive() {
		return SlotLink{}, false
	}

	newSlot := m.allocateSlot()
	newSlot.parent = parent
	newSlot.ctx = ctx
	link := newSlot.NewLink()

	return link, m.prepareNewSlot(newSlot, nil, nil, sm, false)
}

func (m *SlotMachine) createNestedForAdapter(parent SlotLink, fn CreateFunc) (SlotLink, bool) {
	if fn == nil {
		panic("illegal value")
	}
	if !m.IsActive() {
		return SlotLink{}, false
	}

	newSlot := m.allocateSlot()
	newSlot.parent = parent
	if parent.IsValid() {
		newSlot.ctx = parent.s.ctx
	}
	link := newSlot.NewLink()

	if !m.prepareNewSlot(newSlot, nil, fn, nil, false) {
		return link, false
	}

	m.syncQueue.AddAsyncUpdate(link, m._startAddedSlot)
	return link, true
}

func (m *SlotMachine) prepareNewSlot(slot, creator *Slot, fn CreateFunc, sm StateMachine, inherit bool) bool {
	defer func() {
		recovered := recover()
		if recovered != nil {
			m.recycleEmptySlot(slot)
			panic(recovered)
		}
	}()

	var dInjector injector.DependencyInjector
	if fn != nil {
		if sm != nil {
			panic("illegal value")
		}
		cc := constructionContext{s: slot, inherit: inherit}
		sm = cc.executeCreate(fn)
		if sm == nil {
			m.recycleEmptySlot(slot)
			return false
		}

		if cc.inherit && creator.injected != nil { // TODO copy all custom injects from creator
			// use of FindLocalDependency for parentCopy of DependencyInjector
			// allows to get a copy of creator's injects without keeping a reference
			dInjector = injector.NewDependencyInjector(sm, m, creator.injected.FindLocalDependency)
		} else {
			dInjector = injector.NewDependencyInjector(sm, m, nil)
		}
		dInjector.ResolveAndPut(cc.injects)
	} else {
		dInjector = injector.NewDependencyInjector(sm, m, nil)
	}

	decl := sm.GetStateMachineDeclaration()
	if decl == nil {
		panic("illegal state")
	}
	slot.declaration = decl

	link := slot.NewLink()
	decl.InjectDependencies(sm, link, &dInjector)

	initFn := slot.declaration.GetInitStateFor(sm)
	if initFn == nil {
		panic("illegal value")
	}

	if creator != nil {
		slot.migrationCount = creator.migrationCount
		slot.lastWorkScan = creator.lastWorkScan
	} else {
		slot.migrationCount = m.migrationCount
		slot.lastWorkScan = uint8(m.scanCount)
	}

	shadowMigrateFn := slot.declaration.GetShadowMigrateFor(sm)
	if !dInjector.IsEmpty() {
		localInjects := dInjector.CopyAsRegistryNoParent()
		shadowMigrateFn = buildShadowMigrator(localInjects, shadowMigrateFn)
		slot.injected = localInjects
	}
	slot.shadowMigrate = shadowMigrateFn

	slot.step = SlotStep{Transition: func(ctx ExecutionContext) StateUpdate {
		ec := ctx.(*executionContext)
		ic := initializationContext{ec.clone(updCtxInactive)}
		slot.incStep()
		su := ic.executeInitialization(initFn)
		su.marker = ec.getMarker()
		return su
	}}

	return true
}

func (m *SlotMachine) startNewSlot(slot *Slot, worker FixedSlotWorker) {
	slot.ensureInitializing()
	m.stopSlotWorking(slot, 0, worker)
	m.updateSlotQueue(slot, worker, activateSlot)
}

func (m *SlotMachine) startNewSlotByDetachable(slot *Slot, runInit bool, w DetachableSlotWorker) {
	if runInit {
		m._executeSlotInitByCreator(slot, w)
		return
	}

	slot.ensureInitializing()
	if !w.NonDetachableCall(func(worker FixedSlotWorker) {
		m.stopSlotWorking(slot, 0, worker)
		m.updateSlotQueue(slot, worker, activateSlot)
	}) {
		m.syncQueue.AddAsyncUpdate(slot.NewLink(), m._startAddedSlot)
	}
}

func (m *SlotMachine) _startAddedSlot(link SlotLink, worker FixedSlotWorker) {
	if !link.IsValid() {
		panic("unexpected")
	}
	slot := link.s
	slot.ensureInitializing()
	m.stopSlotWorking(slot, 0, worker)
	list := m._updateSlotQueue(slot, false, activateSlot)
	if list != nil {
		panic("unexpected")
	}
}

/* -- Methods to manage processing queues, activate/deactivate slots ------------------------------ */

type slotActivationMode uint8

const (
	deactivateSlot slotActivationMode = iota
	activateSlot
	activateHotWaitSlot
)

func (m *SlotMachine) updateSlotQueue(slot *Slot, w FixedSlotWorker, activation slotActivationMode) {
	if slot.machine != m {
		panic("illegal state")
	}

	s := m._updateSlotQueue(slot, slot.isInQueue(), activation)
	m._activateDependantChain(s, w)
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

/* ---- slot state updates and error handling ---------------------------- */

func (m *SlotMachine) applyStateUpdate(slot *Slot, stateUpdate StateUpdate, w FixedSlotWorker) bool {
	if slot.machine != m {
		panic("illegal state")
	}

	var err error
	isAvailable := false
	isPanic := false

	func() {
		defer func() {
			isPanic = true
			err = RecoverSlotPanic("apply state update panic", recover(), err)
		}()
		isAvailable, err = typeOfStateUpdate(stateUpdate).Apply(slot, stateUpdate, w)
	}()

	if err == nil {
		return isAvailable
	}

	return m.handleSlotUpdateError(slot, w, isPanic, err)
}

func (m *SlotMachine) handleSlotUpdateError(slot *Slot, worker FixedSlotWorker, isPanic bool, err error) bool {

	canRecover := false
	res := ErrorHandlerDefault

	eh := slot.getErrorHandler()
	if eh != nil {
		fc := failureContext{isPanic: isPanic, err: err}
		if se, ok := err.(SlotPanicError); ok {
			fc.isAsync = se.IsAsync
		}
		canRecover = !fc.isPanic || fc.isAsync
		fc.canRecover = canRecover
		res, err = fc.executeFailure(eh)
	}

	recoverState := ""
	switch res {
	case ErrorHandlerMute:
		//recoverState = "recover=muted "
		break
	case ErrorHandlerRecover:
		if canRecover {
			return true
		}
		recoverState = "recover=failed "
		fallthrough
	default:
		// TODO log error m._handleStateUpdateError(slot, stateUpdate, w, err)
		fmt.Printf("SLOT ERROR: slot=%v %serr=%v\n", slot.GetSlotID(), recoverState, err)
		//runtime.KeepAlive(err)
	}

	m.recycleSlot(slot, worker)
	return false
}

/* ------ BargeIn support -------------------------- */

func (m *SlotMachine) createBargeIn(link StepLink, applyFn BargeInApplyFunc) BargeInParamFunc {

	link.s.slotFlags |= slotHasBargeIn
	return func(param interface{}) bool {
		if !link.IsValid() {
			return false
		}
		m.queueAsyncCallback(link.SlotLink, func(slot *Slot, worker DetachableSlotWorker) StateUpdate {
			_, atExactStep := link.isValidAndAtExactStep()
			bc := bargingInContext{slotContext{s: slot}, param, atExactStep}
			return bc.executeBargeIn(applyFn)
		}, nil)
		return true
	}
}

func (m *SlotMachine) bargeInNow(link SlotLink, param interface{}, applyFn BargeInApplyFunc, worker FixedSlotWorker) bool {
	if !link.isMachine(m) {
		return false
	}

	slot, isStarted, _ := link.tryStartWorking()
	if !isStarted {
		return false
	}

	releaseOnPanic := true
	defer func() {
		if releaseOnPanic {
			slot.stopWorking()
		}
	}()
	bc := bargingInContext{slotContext{s: link.s}, param, false}
	stateUpdate := bc.executeBargeInNow(applyFn)

	releaseOnPanic = false
	m.slotPostExecution(slot, stateUpdate, worker, 0, false)
	return true
}

func (m *SlotMachine) createLightBargeIn(link StepLink, stateUpdate StateUpdate) BargeInFunc {

	link.s.slotFlags |= slotHasBargeIn
	return func() bool {
		if !link.IsValid() {
			return false
		}
		m.syncQueue.AddAsyncUpdate(link.SlotLink, func(_ SlotLink, worker FixedSlotWorker) {
			if m._canCallback(link.SlotLink) || !link.IsAtStep() {
				return
			}
			// Plan A - faster one
			if slot, isStarted, prevStepNo := link.tryStartWorking(); isStarted {
				m.slotPostExecution(slot, stateUpdate, worker, prevStepNo, true)
				return
			}
			// Plan B
			m.queueAsyncCallback(link.SlotLink, func(slot *Slot, worker DetachableSlotWorker) StateUpdate {
				if link.IsAtStep() {
					return stateUpdate
				}
				return StateUpdate{}
			}, nil)
		})
		return true
	}
}

/* ----- Time operations --------------------------- */

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

/* ---- Unsorted ---------------------------- */

func (m *SlotMachine) wakeupOnDeactivationOf(slot *Slot, waitOn SlotLink, worker FixedSlotWorker) {
	if waitOn.s == slot || !waitOn.IsValid() {
		// don't wait for self
		// don't wait for an expired slot
		m.updateSlotQueue(slot, worker, activateSlot)
		return
	}

	wakeupLink := slot.NewLink()
	waitOnMachine := waitOn.s.machine
	waitOnMachine._wakeupOnDeactivateAsync(wakeupLink, waitOn)
}

//waitOn MUST belong to this machine!
func (m *SlotMachine) _wakeupOnDeactivateAsync(wakeUp, waitOn SlotLink) {
	m.syncQueue.AddAsyncCallback(waitOn, func(waitOn SlotLink, worker DetachableSlotWorker) bool {
		if !wakeUp.IsValid() {
			// requester is dead - no need to to anything
			return true
		}
		if worker != nil && waitOn.isValidAndBusy() {
			// have to wait further, add this back
			return false
		}

		switch {
		case worker == nil:
			break
		case wakeUp.isMachine(waitOn.s.machine):
			if worker.NonDetachableCall(wakeUp.activateSlot) {
				return true
			}
		default:
			if worker.NonDetachableOuterCall(wakeUp.s.machine, wakeUp.activateSlot) {
				return true
			}
		}
		wakeUp.s.machine.syncQueue.AddAsyncUpdate(wakeUp, SlotLink.activateSlot)
		return true
	})
}

func (m *SlotMachine) useSlotAsShared(link *SharedDataLink, accessFn SharedDataFunc, worker DetachableSlotWorker) SharedAccessReport {
	isValid, isBusy := link.link.getIsValidAndBusy()

	if !isValid {
		return SharedSlotAbsent
	}

	if !link.link.isMachine(m) { // isRemote
		panic("unimplemented") // TODO access to non-local slot machine

		//if isBusy {
		//	return SharedSlotRemoteBusy
		//}
	}

	if isBusy {
		return SharedSlotLocalBusy
	}
	data := link.getData()
	if data == nil {
		return SharedSlotAbsent
	}

	slot, isStarted, _ := link.link.tryStartWorking()
	if !isStarted {
		return SharedSlotLocalBusy
	}

	defer slot.stopWorking()
	wakeUp := accessFn(data)

	m.syncQueue.ProcessSlotCallbacksByDetachable(link.link, worker)
	if !wakeUp && link.flags&ShareDataWakesUpAfterUse == 0 || slot.slotFlags&slotWokenUp != 0 {
		return SharedSlotLocalAvailable
	}
	slot.slotFlags |= slotWokenUp

	if !worker.NonDetachableCall(slot.activateSlot) {
		stepLink := slot.NewStepLink() // remember the current step to avoid "back-fire" activation
		m.syncQueue.AddAsyncUpdate(stepLink.SlotLink, stepLink.activateSlotStepWithSlotLink)
	}
	return SharedSlotLocalAvailable
}

func (m *SlotMachine) stopSlotWorking(slot *Slot, prevStepNo uint32, worker FixedSlotWorker) {
	dep := slot.dependency
	newStepNo := slot.stopWorking()

	switch {
	case dep == nil:
		return
	case prevStepNo == 0:
		// there can be NO dependencies when a slot was just created
		panic("illegal state")
	case newStepNo == prevStepNo || newStepNo <= 1:
		// step didn't change or it is initialization (which is considered as a part of creation)
		return
	case !dep.IsReleaseOnStepping():
		return
	}

	slot.dependency = nil
	m.activateDependants(dep.Release(), worker)
}

func (m *SlotMachine) _activateDependantChain(chain *Slot, worker FixedSlotWorker) {
	for chain != nil {
		s := chain
		// we MUST cut the slot out of chain before any actions on the slot
		chain = chain._cutNext()

		switch {
		case m == s.machine:
			s.activateSlot(worker)
		case worker.OuterCall(s.machine, s.activateSlot):
			//
		default:
			link := s.NewStepLink() // remember the current step to avoid "back-fire" activation
			s.machine.syncQueue.AddAsyncUpdate(link.SlotLink, link.activateSlotStepWithSlotLink)
		}
	}
}

func (m *SlotMachine) activateDependants(links []StepLink, worker FixedSlotWorker) {
	for _, link := range links {
		switch {
		case link.isMachine(m):
			// slot will be activated if it is at the same step as it was when we've decided to activate it
			link.activateSlotStep(worker)
		case worker.OuterCall(link.s.machine, link.activateSlotStep):
			//
		default:
			link.s.machine.syncQueue.AddAsyncUpdate(link.SlotLink, link.activateSlotStepWithSlotLink)
		}
	}
}

func (m *SlotMachine) activateDependantByDetachable(links []StepLink, worker DetachableSlotWorker) bool {
	if len(links) == 0 {
		return false
	}

	if worker.NonDetachableCall(func(worker FixedSlotWorker) {
		m.activateDependants(links, worker)
	}) {
		return true
	}

	m.syncQueue.AddAsyncUpdate(SlotLink{}, func(_ SlotLink, worker FixedSlotWorker) {
		m.activateDependants(links, worker)
	})
	return true
}
