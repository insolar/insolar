//
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
//

package smachine

import (
	"context"
	"fmt"
	"math"
	"sync"
	"sync/atomic"
	"time"

	"github.com/insolar/insolar/conveyor/injector"
	"github.com/insolar/insolar/conveyor/tools"
)

type MigrationFunc func(migrationCount uint32)

type SlotMachineConfig struct {
	PollingPeriod        time.Duration
	PollingTruncate      time.Duration
	BoostNewSlotDuration time.Duration
	SlotPageSize         uint16
	ScanCountLimit       int
	CleanupWeakOnMigrate bool

	SlotIdGenerateFn  func() SlotID
	SlotMachineLogger SlotMachineLogger
	SlotAliasRegistry SlotAliasRegistry
}

type SlotAliasRegistry interface {
	PublishAlias(key interface{}, slot SlotLink) bool
	UnpublishAlias(key interface{})
	GetPublishedAlias(key interface{}) SlotLink
}

const maxLoopCount = 10000

func NewSlotMachine(config SlotMachineConfig,
	eventCallback, signalCallback func(),
	parentRegistry injector.DependencyRegistry,
) *SlotMachine {
	if config.ScanCountLimit <= 0 || config.ScanCountLimit > maxLoopCount {
		config.ScanCountLimit = maxLoopCount
	}

	m := &SlotMachine{
		config:         config,
		parentRegistry: parentRegistry,
		slotPool:       newSlotPool(config.SlotPageSize, false),
		syncQueue:      newSlotMachineSync(eventCallback, signalCallback),
	}

	if m.config.SlotIdGenerateFn == nil {
		m.config.SlotIdGenerateFn = m._allocateNextSlotID
	}

	m.slotPool.initSlotPool()

	m.activeSlots.initSlotQueue(ActiveSlots)
	m.boostedSlots.initSlotQueue(ActiveSlots)
	m.prioritySlots.initSlotQueue(ActiveSlots)

	m.workingSlots.initSlotQueue(WorkingSlots)

	return m
}

var _ injector.DependencyRegistry = &SlotMachine{}

type SlotMachine struct {
	config     SlotMachineConfig
	lastSlotID SlotID // atomic
	slotPool   SlotPool

	parentRegistry injector.DependencyRegistry
	localRegistry  sync.Map // is used for both dependencies and tracking of dependencies

	machineStartedAt time.Time
	scanStartedAt    time.Time

	boostPermitLatest   *chainedBoostPermit
	boostPermitEarliest *chainedBoostPermit

	scanAndMigrateCounts uint64 // atomic

	migrators []MigrationFunc

	hotWaitOnly  bool      // true when activeSlots & prioritySlots have only slots added by "hot wait"
	scanWakeUpAt time.Time // when all slots are waiting, this is the earliest time requested for wakeup

	prioritySlots    SlotQueue //they are are moved to workingSlots every time when enough non-priority slots are processed
	boostedSlots     SlotQueue
	activeSlots      SlotQueue    //they are are moved to workingSlots on every full Scan
	pollingSlots     PollingQueue //they are are moved to workingSlots on every full Scan when time has passed
	workingSlots     SlotQueue    //slots are currently in processing
	nonPriorityCount uint32       // number of non-priority slots processed since the last replenishment of workingSlots

	syncQueue SlotMachineSync
}

type ScanMode uint8

const (
	ScanDefault ScanMode = iota
	ScanPriorityOnly
	ScanEventsOnly
)

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

func (m *SlotMachine) getScanAndMigrateCounts() (scanCount, migrateCount uint32) {
	v := atomic.LoadUint64(&m.scanAndMigrateCounts)
	return uint32(v), uint32(v >> 32)
}

func (m *SlotMachine) getScanCount() uint32 {
	v := atomic.LoadUint64(&m.scanAndMigrateCounts)
	return uint32(v)
}

func (m *SlotMachine) incScanCount() uint32 {
	for {
		v := atomic.LoadUint64(&m.scanAndMigrateCounts)
		vv := uint64(uint32(v)+1) | v&^math.MaxUint32
		if atomic.CompareAndSwapUint64(&m.scanAndMigrateCounts, v, vv) {
			return uint32(vv)
		}
	}
}

func (m *SlotMachine) incMigrateCount() uint32 {
	return uint32(atomic.AddUint64(&m.scanAndMigrateCounts, 1<<32) >> 32)
}

func (m *SlotMachine) CopyConfig() SlotMachineConfig {
	return m.config
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

func (m *SlotMachine) AddDependency(v interface{}) {
	if !m.TryPutDependency(injector.GetDefaultInjectionId(v), v) {
		panic(fmt.Errorf("duplicate dependency: %T %[1]v", v))
	}
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

func (m *SlotMachine) GetPublishedGlobalAlias(key interface{}) SlotLink {
	return m.getGlobalPublished(key)
}

/* -------------- Methods to run state machines --------------- */

func (m *SlotMachine) RunToStop(worker AttachableSlotWorker, signal *tools.SignalVersion) {
	m.Stop()
	worker.AttachTo(m, signal, uint32(m.config.ScanCountLimit), func(worker AttachedSlotWorker) {
		for !m.syncQueue.IsInactive() && !worker.HasSignal() {
			m.ScanOnce(ScanDefault, worker)
		}
	})
}

func (m *SlotMachine) ScanNested(outerCtx ExecutionContext, scanMode ScanMode,
	loopLimit uint32, worker AttachableSlotWorker,
) (repeatNow bool, nextPollTime time.Time) {
	if loopLimit == 0 {
		loopLimit = uint32(m.config.ScanCountLimit)
	}
	ec := outerCtx.(*executionContext)

	worker.AttachAsNested(m, ec.w, loopLimit, func(worker AttachedSlotWorker) {
		repeatNow, nextPollTime = m.ScanOnce(scanMode, worker)
	})
	return repeatNow, nextPollTime
}

func (m *SlotMachine) ScanOnce(scanMode ScanMode, worker AttachedSlotWorker) (repeatNow bool, nextPollTime time.Time) {
	status := m.syncQueue.GetStatus()
	if status == SlotMachineInactive {
		return false, time.Time{}
	}

	scanTime := time.Now()
	m.beforeScan(scanTime)
	currentScanNo := uint32(0)

	switch {
	case m.config.BoostNewSlotDuration > 0:
		m.boostPermitLatest = m.boostPermitLatest.reuseOrNew(scanTime)
		switch m.boostPermitEarliest {
		case nil:
			m.boostPermitEarliest = m.boostPermitLatest
		case m.boostPermitLatest:
			// reuse, no need to check
		default:
			m.boostPermitEarliest = m.boostPermitEarliest.discardOlderThan(scanTime.Add(-m.config.BoostNewSlotDuration))
			if m.boostPermitEarliest == nil {
				panic("unexpected")
			}
		}
	case m.boostPermitLatest != nil:
		panic("unexpected")
	}

	switch {
	case m.machineStartedAt.IsZero():
		m.machineStartedAt = scanTime
		fallthrough
	case scanMode == ScanEventsOnly:
		// no scans
		currentScanNo = m.getScanCount()

	case !m.workingSlots.IsEmpty():
		// we were interrupted
		currentScanNo = m.getScanCount()
		if m.nonPriorityCount >= uint32(m.config.ScanCountLimit) {
			m.nonPriorityCount = 0
			m.workingSlots.AppendAll(&m.prioritySlots)
			if scanMode != ScanPriorityOnly {
				m.workingSlots.AppendAll(&m.boostedSlots)
			}
		}

	default:
		currentScanNo = m.incScanCount()

		m.hotWaitOnly = true
		m.nonPriorityCount = 0
		m.workingSlots.AppendAll(&m.prioritySlots)

		if scanMode != ScanPriorityOnly {
			m.workingSlots.AppendAll(&m.boostedSlots)
			m.workingSlots.AppendAll(&m.activeSlots)
		}
		m.pollingSlots.FilterOut(scanTime, m.workingSlots.AppendAll)
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

	if scanMode != ScanEventsOnly && !hasSignal && !wasDetached {
		m.executeWorkingSlots(currentScanNo, scanMode == ScanPriorityOnly, worker)
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

func (m *SlotMachine) executeWorkingSlots(currentScanNo uint32, priorityOnly bool, worker AttachedSlotWorker) {
	limit := m.config.ScanCountLimit
	for i := 0; i < limit; i++ {
		currentSlot := m.workingSlots.First()
		if currentSlot == nil {
			return
		}
		loopLimit := 1 + ((limit - i) / m.workingSlots.Count())

		prevStepNo := currentSlot.startWorking(currentScanNo) // its counterpart is in slotPostExecution()
		currentSlot.removeFromQueue()

		switch {
		case currentSlot.step.Flags&StepPriority != 0:
			// break
		case priorityOnly:
			if currentSlot.isBoosted() {
				m.boostedSlots.AddLast(currentSlot)
			} else {
				m.activeSlots.AddLast(currentSlot)
			}
			currentSlot.stopWorking()
			continue
		case currentSlot.isBoosted():
			// break
		default:
			m.nonPriorityCount++
		}

		if stopNow, loopExtraIncrement := m._executeSlot(currentSlot, prevStepNo, worker, loopLimit); stopNow {
			return
		} else {
			i += loopExtraIncrement
		}
	}
}

func (m *SlotMachine) _executeSlot(slot *Slot, prevStepNo uint32, worker AttachedSlotWorker, loopLimit int) (hasSignal bool, loopCount int) {

	inactivityNano := slot.touch(time.Now().UnixNano())

	if dep := slot.dependency; dep != nil && dep.IsReleaseOnWorking() {
		released := slot._releaseDependency()
		m.activateDependants(released, slot.NewLink(), worker)
	}
	slot.slotFlags &^= slotWokenUp

	var stateUpdate StateUpdate
	wasDetached := worker.DetachableCall(func(worker DetachableSlotWorker) {

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
		m.asyncPostSlotExecution(slot, stateUpdate, prevStepNo, inactivityNano)
		return true, loopCount
	}

	hasAsync := m.slotPostExecution(slot, stateUpdate, worker, prevStepNo, false, inactivityNano)
	if hasAsync && !hasSignal {
		_, hasSignal, wasDetached = m.syncQueue.ProcessCallbacks(worker)
		return hasSignal || wasDetached, loopCount
	}
	return hasSignal, loopCount
}

const durationUnknownNano = time.Duration(1)
const durationNotApplicableNano = time.Duration(0)

func (m *SlotMachine) _executeSlotInitByCreator(slot *Slot, worker DetachableSlotWorker) {

	slot.ensureInitializing()
	m._boostNewSlot(slot)

	slot.touch(time.Now().UnixNano())

	ec := executionContext{slotContext: slotContext{s: slot, w: worker}}
	stateUpdate, _, asyncCnt := ec.executeNextStep()

	slot.addAsyncCount(asyncCnt)
	if !worker.NonDetachableCall(func(worker FixedSlotWorker) {
		m.slotPostExecution(slot, stateUpdate, worker, 0, false, durationUnknownNano)
	}) {
		m.asyncPostSlotExecution(slot, stateUpdate, 0, durationUnknownNano)
	}
}

func (m *SlotMachine) slotPostExecution(slot *Slot, stateUpdate StateUpdate, worker FixedSlotWorker,
	prevStepNo uint32, wasAsync bool, inactivityNano time.Duration) (hasAsync bool) {

	activityNano := durationNotApplicableNano
	if !wasAsync && inactivityNano > durationNotApplicableNano {
		activityNano = slot.touch(time.Now().UnixNano())
	}

	slot.logStepUpdate(prevStepNo, stateUpdate, wasAsync, inactivityNano, activityNano)

	slot.updateBoostFlag()

	if !stateUpdate.IsZero() && !m.applyStateUpdate(slot, stateUpdate, worker) {
		return false
	}

	if slot.canMigrateWorking(prevStepNo, wasAsync) {
		_, migrateCount := m.getScanAndMigrateCounts()
		if _, isAvailable := m._migrateSlot(migrateCount, slot, prevStepNo, worker); !isAvailable {
			return false
		}
	}

	hasAsync = wasAsync || slot.hasAsyncOrBargeIn()
	m.stopSlotWorking(slot, prevStepNo, worker)
	return hasAsync
}

func (m *SlotMachine) queueAsyncCallback(link SlotLink,
	callbackFn func(*Slot, DetachableSlotWorker, error) StateUpdate, prevErr error) {

	if callbackFn == nil && prevErr == nil || !m._canCallback(link) {
		return
	}

	m.syncQueue.AddAsyncCallback(link, func(link SlotLink, worker DetachableSlotWorker) (isDone bool) {
		if !m._canCallback(link) {
			return true
		}
		if worker == nil {
			if m.IsActive() {
				step, _ := link.GetStepLink()
				m.logInternal(step, "async detachment retry limit exceeded", nil)
			}
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
				stateUpdate = callbackFn(slot, worker, prevErr)
			}
		}()

		if worker.NonDetachableCall(func(worker FixedSlotWorker) {
			m.slotPostExecution(slot, stateUpdate, worker, prevStepNo, true, durationNotApplicableNano)
		}) {
			m.syncQueue.ProcessDetachQueue(link, worker)
		} else {
			m.asyncPostSlotExecution(slot, stateUpdate, prevStepNo, durationNotApplicableNano)
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

func (m *SlotMachine) asyncPostSlotExecution(s *Slot, stateUpdate StateUpdate, prevStepNo uint32, inactivityNano time.Duration) {
	m.syncQueue.AddAsyncUpdate(s.NewLink(), func(link SlotLink, worker FixedSlotWorker) {
		if !link.IsValid() {
			return
		}
		slot := link.s
		if m.slotPostExecution(slot, stateUpdate, worker, prevStepNo, true, inactivityNano) {
			m.syncQueue.FlushSlotDetachQueue(link)
		}
	})
}

/* -- Methods to migrate slots ------------------------------ */

func (m *SlotMachine) TryMigrateNested(outerCtx ExecutionContext) bool {
	ec := outerCtx.(*executionContext)
	return ec.w.NonDetachableCall(m.migrate)
}

func (m *SlotMachine) MigrateNested(outerCtx MigrationContext) {
	mc := outerCtx.(*migrationContext)
	m.migrate(mc.fixedWorker)
}

func (m *SlotMachine) migrate(worker FixedSlotWorker) {
	m.migrateWithBefore(worker, nil)
}

func (m *SlotMachine) migrateWithBefore(worker FixedSlotWorker, beforeFn func()) {
	migrateCount := m.incMigrateCount()

	if beforeFn != nil {
		beforeFn()
	}

	for _, fn := range m.migrators {
		fn(migrateCount)
	}

	m.slotPool.ScanAndCleanup(m.config.CleanupWeakOnMigrate, worker, m.recycleSlot,
		func(slotPage []Slot, worker FixedSlotWorker) (isPageEmptyOrWeak, hasWeakSlots bool) {
			return m.migratePage(migrateCount, slotPage, worker)
		})

	m.syncQueue.CleanupDetachQueues()
}

func (m *SlotMachine) AddMigrationCallback(fn MigrationFunc) {
	if fn == nil {
		panic("illegal value")
	}
	m.migrators = append(m.migrators, fn)
}

func (m *SlotMachine) migratePage(migrateCount uint32, slotPage []Slot, worker FixedSlotWorker) (isPageEmptyOrWeak, hasWeakSlots bool) {
	isPageEmptyOrWeak = true
	hasWeakSlots = false
	for i := range slotPage {
		isSlotEmptyOrWeak, isSlotAvailable := m.migrateSlot(migrateCount, &slotPage[i], worker)
		switch {
		case !isSlotEmptyOrWeak:
			isPageEmptyOrWeak = false
		case isSlotAvailable:
			hasWeakSlots = true
		}
	}
	return isPageEmptyOrWeak, hasWeakSlots
}

func (m *SlotMachine) migrateSlot(migrateCount uint32, slot *Slot, w FixedSlotWorker) (isEmptyOrWeak, isAvailable bool) {
	isEmpty, isStarted, prevStepNo := slot.tryStartMigrate()
	if !isStarted {
		return isEmpty, false
	}
	isEmptyOrWeak, isAvailable = m._migrateSlot(migrateCount, slot, prevStepNo, w)
	if isAvailable {
		m.stopSlotWorking(slot, prevStepNo, w)
	}
	return isEmptyOrWeak, isAvailable
}

func (m *SlotMachine) _migrateSlot(migrateCount uint32, slot *Slot, prevStepNo uint32, worker FixedSlotWorker) (isEmptyOrWeak, isAvailable bool) {

	inactivityNano := slot.touch(time.Now().UnixNano())

	for delta := migrateCount - slot.migrationCount; delta > 0; {
		migrateFn := slot.getMigration()

		if migrateFn != nil {
			if slot.shadowMigrate != nil {
				slot.shadowMigrate(slot.migrationCount, 1)
			}

			mc := migrationContext{
				slotContext: slotContext{s: slot, w: migrateWorkerWrapper{worker}},
				fixedWorker: worker,
			}

			stateUpdate, skipMultiple := mc.executeMigration(migrateFn)

			activityNano := slot.touch(time.Now().UnixNano())
			slot.logStepMigrate(prevStepNo, stateUpdate, inactivityNano, activityNano)
			inactivityNano = durationUnknownNano

			if !m.applyStateUpdate(slot, stateUpdate, worker) {
				return true, false
			}
			slot.migrationCount++
			delta--
			if delta == 0 {
				break
			}
			if !skipMultiple {
				continue
			}
		}

		if slot.shadowMigrate != nil {
			slot.shadowMigrate(slot.migrationCount, delta)
		}
		slot.migrationCount = migrateCount
		break
	}

	return slot.step.Flags&StepWeak != 0, true
}

/* -- Methods to allocate slots ------------------------------ */

// SAFE for concurrent use
func (m *SlotMachine) allocateNextSlotID() SlotID {
	return m.config.SlotIdGenerateFn()
}

func (m *SlotMachine) _allocateNextSlotID() SlotID {
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

// SAFE for concurrent use
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
		case slot.isBusy():
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

		switch isEmpty, isStarted, _ := slot._tryStartSlot(1); {
		case isEmpty:
			//continue
		case isStarted:
			m.recycleSlot(slot, w)
		default:
			hasWorking = true
		}
	}
	return !hasWorking, false
}

func (m *SlotMachine) recycleSlot(slot *Slot, worker FixedSlotWorker) {
	m.recycleSlotWithError(slot, worker, nil)
}

func (m *SlotMachine) recycleSlotWithError(slot *Slot, worker FixedSlotWorker, err error) {

	var link StepLink
	hasPanic := false
	func() {
		defer func() {
			recovered := recover()
			hasPanic = recovered != nil
			err = RecoverSlotPanicWithStack("internal panic - recycleSlot", recovered, err)
		}()

		link = slot.NewStepLink()
		slot.invalidateSlotId() // slotId is reset here and all links are invalid since this moment

		th := slot.defTerminate
		if th != nil {
			slot.defTerminate = nil // avoid self-loops
			m.runTerminationHandler(slot.ctx, th, TerminationData{
				Slot:   link,
				Parent: slot.parent,
				Result: slot.defResult,
				Error:  err,
				worker: worker,
			})
		}

		if slot.slotFlags&(slotHadAsync|slotHasBargeIn|slotHasAliases) != 0 {
			defer m.syncQueue.FlushSlotDetachQueue(link.SlotLink)
		}

		if slot.slotFlags&slotHasAliases != 0 {
			// cleanup aliases associated with the slot
			// MUST happen before releasing of dependencies
			m.unregisterBoundAliases(link.SlotID())
		}

		{
			// cleanup synchronization dependency
			if slot.dependency != nil {
				released := slot._releaseDependency()
				m.activateDependants(released, link.SlotLink, worker)
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
	}()

	if hasPanic {
		m.logCritical(link, "recycle", err)
	} else {
		slot.logInternal(link, "recycle", err)
	}
	m._recycleSlot(slot)
}

// SAFE for concurrent use
// This method can be called concurrently but ONLY to release new (empty) slots - slot MUST NOT have any kind of dependencies
func (m *SlotMachine) recycleEmptySlot(slot *Slot, err error) {
	if slot.slotFlags != 0 {
		panic("illegal state")
	}

	th := slot.defTerminate
	if th != nil { // it can be already set by construction - we must invoke it
		slot.defTerminate = nil // avoid self-loops
		m.runTerminationHandler(slot.ctx, th, TerminationData{
			Slot:   slot.NewStepLink(),
			Parent: slot.parent,
			Result: slot.defResult,
			Error:  err,
			worker: nil,
		})
	}

	// slot.invalidateSlotId() // empty slot doesn't need early invalidation

	m._recycleSlot(slot) // SAFE for concurrent use
}

func (m *SlotMachine) _recycleSlot(slot *Slot) {
	slot.dispose()               // check state and cleanup fields
	m.slotPool.RecycleSlot(slot) // SAFE for concurrent use
}

func (m *SlotMachine) OccupiedSlotCount() int {
	return m.slotPool.Count()
}

func (m *SlotMachine) AllocatedSlotCount() int {
	return m.slotPool.Capacity()
}

/* -- General purpose synchronization ------------------------------ */

func (m *SlotMachine) ScheduleCall(fn MachineCallFunc, isSignal bool) bool {
	if fn == nil {
		panic("illegal value")
	}
	callFn := func(_ SlotLink, worker FixedSlotWorker) {
		mc := machineCallContext{m: m, w: worker}
		err := mc.executeCall(fn)
		if err != nil {
			m.logInternal(StepLink{}, "schedule call execution", err)
		}
	}
	if isSignal {
		return m.syncQueue.AddAsyncSignal(SlotLink{}, callFn)
	} else {
		return m.syncQueue.AddAsyncUpdate(SlotLink{}, callFn)
	}
}

// SAFE for concurrent use
func (m *SlotMachine) runTerminationHandler(ctx context.Context, th TerminationHandlerFunc, td TerminationData) {
	if ctx == nil {
		ctx = context.Background() // TODO provide a default context for SlotMachine?
	}

	m.syncQueue.AddAsyncCallback(td.Slot.SlotLink, func(_ SlotLink, _ DetachableSlotWorker) bool {
		err := func() (err error) {
			defer func() {
				err = RecoverSlotPanicWithStack("termination handler", recover(), nil)
			}()
			th(ctx, td)
			return nil
		}()
		if err != nil {
			m.logInternal(StepLink{SlotLink: td.Slot.SlotLink}, "failed termination handler", err)
		}
		return true
	})
}

/* -- Methods to create and start new machines ------------------------------ */

func (m *SlotMachine) AddNew(ctx context.Context, sm StateMachine, defValues CreateDefaultValues) SlotLink {
	switch {
	case ctx != nil:
		defValues.Context = ctx
	case defValues.Context == nil:
		panic("illegal value")
	}

	link, ok := m.prepareNewSlotWithDefaults(nil, nil, sm, defValues)
	if ok {
		m.syncQueue.AddAsyncUpdate(link, m._startAddedSlot)
	}
	return link
}

func (m *SlotMachine) AddNewByFunc(ctx context.Context, cf CreateFunc, defValues CreateDefaultValues) (SlotLink, bool) {
	switch {
	case ctx != nil:
		defValues.Context = ctx
	case defValues.Context == nil:
		panic("illegal value")
	}

	link, ok := m.prepareNewSlotWithDefaults(nil, cf, nil, defValues)
	if ok {
		m.syncQueue.AddAsyncUpdate(link, m._startAddedSlot)
	}
	return link, ok
}

func (m *SlotMachine) AddNested(_ AdapterId, parent SlotLink, cf CreateFunc) (SlotLink, bool) {
	if parent.IsEmpty() {
		panic("illegal value")
	}
	// TODO pass adapterId into injections?

	link, ok := m.prepareNewSlot(nil, cf, nil,
		prepareSlotValue{slotReplaceData: slotReplaceData{parent: parent}})

	if ok {
		m.syncQueue.AddAsyncUpdate(link, m._startAddedSlot)
	}
	return link, ok
}

type prepareSlotValue struct {
	slotReplaceData
	overrides     map[string]interface{}
	terminate     TerminationHandlerFunc
	tracerId      TracerId
	isReplacement bool
}

func (m *SlotMachine) prepareNewSlotWithDefaults(creator *Slot, fn CreateFunc, sm StateMachine, defValues CreateDefaultValues) (SlotLink, bool) {
	return m.prepareNewSlot(creator, fn, sm, prepareSlotValue{
		slotReplaceData: slotReplaceData{
			parent: defValues.Parent,
			ctx:    defValues.Context,
		},
		overrides: defValues.OverriddenDependencies,
		tracerId:  defValues.TracerId,
		terminate: defValues.TerminationHandler,
	})
}

func mergeDefaultValues(target *prepareSlotValue, source CreateDefaultValues) {
	if source.Context != nil {
		target.ctx = source.Context
	}
	if !source.Parent.IsEmpty() {
		target.parent = source.Parent
	}
	if source.TerminationHandler != nil {
		target.terminate = source.TerminationHandler
	}
	if len(source.TracerId) > 0 {
		target.tracerId = source.TracerId
	}

	switch {
	case source.OverriddenDependencies == nil:
	case target.overrides == nil:
		target.overrides = source.OverriddenDependencies
	default:
		for k, v := range source.OverriddenDependencies {
			target.overrides[k] = v
		}
	}
}

// caller MUST be busy-holder of both creator and slot, then this method is SAFE for concurrent use
func (m *SlotMachine) prepareNewSlot(creator *Slot, fn CreateFunc, sm StateMachine, defValues prepareSlotValue) (SlotLink, bool) {
	switch {
	case (fn == nil) == (sm == nil):
		panic("illegal value")
	case !m.IsActive():
		return SlotLink{}, false
	}

	slot := m.allocateSlot()
	defer func() {
		if slot != nil {
			m.recycleEmptySlot(slot, nil) // all construction errors are reported to caller
		}
	}()

	slot.slotReplaceData = defValues.slotReplaceData
	slot.defTerminate = defValues.terminate // terminate handler must be executed even if construction has failed

	switch {
	case slot.ctx != nil:
	case creator != nil:
		slot.ctx = creator.ctx
	case slot.parent.IsValid():
		// TODO this can be racy when the parent is from another SlotMachine running under a different worker ...
		slot.ctx = slot.parent.s.ctx
	}
	if slot.ctx == nil {
		slot.ctx = context.Background() // TODO provide SlotMachine context?
	}

	cc := constructionContext{s: slot, injects: defValues.overrides, tracerId: defValues.tracerId}
	if defValues.isReplacement {
		cc.inherit = InheritResolvedDependencies
	}

	if fn != nil {
		sm = cc.executeCreate(fn)
		if sm == nil {
			return slot.NewLink(), false // slot will be released by defer
		}
	}

	decl := sm.GetStateMachineDeclaration()
	if decl == nil {
		panic(fmt.Errorf("illegal state - declaration is missing: %v", sm))
	}
	slot.declaration = decl

	link := slot.NewLink()

	// get injects sorted out
	var localInjects []interface{}
	slot.inheritable, localInjects = m.prepareInjects(creator, link, sm, cc.inherit, defValues.isReplacement,
		cc.injects, defValues.inheritable)

	// Step Logger
	var stepLoggerFactory StepLoggerFactoryFunc
	if m.config.SlotMachineLogger != nil {
		stepLoggerFactory = m.config.SlotMachineLogger.CreateStepLogger
	}
	if stepLogger, ok := decl.GetStepLogger(slot.ctx, sm, cc.tracerId, stepLoggerFactory); ok {
		slot.stepLogger = stepLogger
	} else if stepLoggerFactory != nil {
		slot.stepLogger = stepLoggerFactory(slot.ctx, sm, cc.tracerId)
	}
	if slot.stepLogger == nil && len(cc.tracerId) > 0 {
		slot.stepLogger = StepLoggerStub{cc.tracerId}
	}

	slot.setTracing(cc.isTracing)

	// get Init step
	initFn := slot.declaration.GetInitStateFor(sm)
	if initFn == nil {
		panic(fmt.Errorf("illegal state - initialization is missing: %v", sm))
	}

	// Setup Slot counters
	if creator != nil {
		slot.migrationCount = creator.migrationCount
		slot.lastWorkScan = creator.lastWorkScan
	} else {
		scanCount, migrateCount := m.getScanAndMigrateCounts()
		slot.migrationCount = migrateCount
		slot.lastWorkScan = uint8(scanCount)
	}

	// shadow migrate for injected dependencies
	slot.shadowMigrate = buildShadowMigrator(localInjects, slot.declaration.GetShadowMigrateFor(sm))

	// final touch
	slot.step = SlotStep{Transition: initFn.defaultInit}
	slot.stepDecl = &defaultInitDecl

	slot = nil //protect from defer
	return link, true
}

var defaultInitDecl = StepDeclaration{stepDeclExt: stepDeclExt{Name: "<init>"}}

func (v InitFunc) defaultInit(ctx ExecutionContext) StateUpdate {
	ec := ctx.(*executionContext)
	if ec.s.shadowMigrate != nil {
		ec.s.shadowMigrate(ec.s.migrationCount, 0)
	}
	ic := initializationContext{ec.clone(updCtxInactive)}
	su := ic.executeInitialization(v)
	su.marker = ec.getMarker()
	return su
}

func (m *SlotMachine) prepareInjects(creator *Slot, link SlotLink, sm StateMachine, mode DependencyInheritanceMode, isReplacement bool,
	constructorOverrides, defValuesInjects map[string]interface{},
) (map[string]interface{}, []interface{}) {

	var overrides []map[string]interface{}
	if len(constructorOverrides) > 0 {
		overrides = append(overrides, constructorOverrides)
	}
	if len(defValuesInjects) > 0 {
		overrides = append(overrides, defValuesInjects)
	}
	if mode&InheritResolvedDependencies != 0 && creator != nil && len(creator.inheritable) > 0 {
		overrides = append(overrides, creator.inheritable)
	}

	var localDeps injector.DependencyRegistryFunc
	if len(overrides) > 0 {
		localDeps = func(id string) (interface{}, bool) {
			for _, om := range overrides {
				if v, ok := om[id]; ok {
					return v, true
				}
			}
			return nil, false
		}
	}

	var addedInjects []interface{}

	dResolver := injector.NewDependencyResolver(sm, m, localDeps, func(_ string, v interface{}, from injector.DependencyOrigin) {
		if from&(injector.DependencyFromLocal|injector.DependencyFromProvider) != 0 {
			addedInjects = append(addedInjects, v)
		}
	})
	dInjector := injector.NewDependencyInjectorFor(&dResolver)

	link.s.declaration.InjectDependencies(sm, link, &dInjector)

	switch {
	case isReplacement:
		// replacing SM should take all inherited dependencies, even if unused
		if mode&copyAllDependencies != 0 {
			for _, o := range overrides {
				dResolver.ResolveAndMerge(o)
			}
		}
		if mode&DiscardResolvedDependencies != 0 {
			return nil, addedInjects
		}

	case mode&DiscardResolvedDependencies != 0:
		return nil, addedInjects

	case mode&copyAllDependencies != 0:
		localInjects := addedInjects // keep only injects that were explicitly used
		for _, o := range overrides {
			dResolver.ResolveAndMerge(o)
		}
		return dResolver.Flush(), localInjects
	}

	return dResolver.Flush(), addedInjects
}

func (m *SlotMachine) _boostNewSlot(slot *Slot) {
	if bp := m.boostPermitLatest; bp != nil && slot.boost == nil {
		bp.use()
		slot.boost = &bp.boostPermit
		if slot.boost.isActive() {
			slot.slotFlags |= slotIsBoosted
		}
	}
}

func (m *SlotMachine) startNewSlot(slot *Slot, worker FixedSlotWorker) {
	slot.ensureInitializing()
	m._boostNewSlot(slot)
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
		m._boostNewSlot(slot)
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
	m._boostNewSlot(slot)
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
		// hot wait ignores boosted to reduce interference
		switch {
		case slot.isPriority():
			m.prioritySlots.AddLast(slot)
		//case slot.isBoosted():
		//	m.boostedSlots.AddLast(slot)
		default:
			m.activeSlots.AddLast(slot)
		}
	case slot.isLastScan(m.getScanCount()):
		m.hotWaitOnly = false
		switch {
		case slot.isPriority():
			m.prioritySlots.AddLast(slot)
		case slot.isBoosted():
			m.boostedSlots.AddLast(slot)
		default:
			m.activeSlots.AddLast(slot)
		}
	default:
		// addSlotToWorkingQueue
		if slot.isPriority() {
			m.workingSlots.AddFirst(slot)
		} else {
			m.workingSlots.AddLast(slot)
		}
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
			recovered := recover()
			isPanic = recovered != nil
			err = RecoverSlotPanicWithStack("apply state update panic", recovered, err)
		}()
		isAvailable, err = typeOfStateUpdate(stateUpdate).Apply(slot, stateUpdate, w)
	}()

	if err == nil {
		return isAvailable
	}

	return m.handleSlotUpdateError(slot, w, stateUpdate, isPanic, err)
}

func (m *SlotMachine) handleSlotUpdateError(slot *Slot, worker FixedSlotWorker, stateUpdate StateUpdate, isPanic bool, err error) bool {

	canRecover := false
	isAsync := false
	if se, ok := err.(SlotPanicError); ok {
		isAsync = se.IsAsync
	}
	canRecover = isAsync // || !isPanic

	action := ErrorHandlerDefault

	eh := slot.getErrorHandler()
	if eh != nil {
		fc := failureContext{isPanic: isPanic, isAsync: isAsync, canRecover: canRecover, err: err, result: slot.defResult}

		ok := false
		if ok, action, err = fc.executeFailure(eh); ok {
			// do not change result on failure of the error handler
			slot.defResult = fc.result
		} else {
			action = ErrorHandlerDefault
		}
	}

	switch action {
	case ErrorHandlerRecover, ErrorHandlerRecoverAndWakeUp:
		switch {
		case !canRecover:
			slot.logStepError(errorHandlerRecoverDenied, stateUpdate, isAsync, err)
		case action == ErrorHandlerRecoverAndWakeUp:
			slot.activateSlot(worker)
			fallthrough
		default:
			slot.logStepError(ErrorHandlerRecover, stateUpdate, isAsync, err)
			return true
		}
	case ErrorHandlerMute:
		slot.logStepError(ErrorHandlerMute, stateUpdate, isAsync, err)
	default:
		slot.logStepError(ErrorHandlerDefault, stateUpdate, isAsync, err)
	}

	m.recycleSlotWithError(slot, worker, err)
	return false
}

func (m *SlotMachine) logCritical(link StepLink, msg string, err error) {
	if sml := m.config.SlotMachineLogger; sml != nil {
		sml.LogMachineCritical(SlotMachineData{m.getScanCount(), link, err}, msg)
	}
}

func (m *SlotMachine) logInternal(link StepLink, msg string, err error) {
	if sml := m.config.SlotMachineLogger; sml != nil {
		sml.LogMachineInternal(SlotMachineData{m.getScanCount(), link, err}, msg)
	}
}

/* ------ BargeIn support -------------------------- */

func (m *SlotMachine) createBargeIn(link StepLink, applyFn BargeInApplyFunc) BargeInParamFunc {

	link.s.slotFlags |= slotHasBargeIn
	return func(param interface{}) bool {
		if !link.IsValid() {
			return false
		}
		m.queueAsyncCallback(link.SlotLink, func(slot *Slot, worker DetachableSlotWorker, _ error) StateUpdate {
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
	m.slotPostExecution(slot, stateUpdate, worker, 0, false, durationNotApplicableNano)
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
				m.slotPostExecution(slot, stateUpdate, worker, prevStepNo, true, durationNotApplicableNano)
				return
			}
			// Plan B
			m.queueAsyncCallback(link.SlotLink, func(slot *Slot, worker DetachableSlotWorker, _ error) StateUpdate {
				if link.IsAtStep() {
					return stateUpdate
				}
				return StateUpdate{} // no change
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
		switch {
		case !wakeUp.IsValid():
			// requester is dead - no need to to anything
			return true
		case worker != nil && waitOn.isValidAndBusy():
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

	released := slot._releaseDependency()
	m.activateDependants(released, slot.NewLink(), worker)
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

func (m *SlotMachine) activateDependants(links []StepLink, ignore SlotLink, worker FixedSlotWorker) {
	for _, link := range links {
		switch {
		case link.SlotLink == ignore:
			continue
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

func (m *SlotMachine) activateDependantByDetachable(links []StepLink, ignore SlotLink, worker DetachableSlotWorker) bool {
	if len(links) == 0 {
		return false
	}

	if worker.NonDetachableCall(func(worker FixedSlotWorker) {
		m.activateDependants(links, ignore, worker)
	}) {
		return true
	}

	m.syncQueue.AddAsyncUpdate(SlotLink{}, func(_ SlotLink, worker FixedSlotWorker) {
		m.activateDependants(links, ignore, worker)
	})
	return true
}

func (m *SlotMachine) GetStoppingSignal() <-chan struct{} {
	return m.syncQueue.GetStoppingSignal()
}

// Must support nil receiver
func (m *SlotMachine) GetMachineId() string {
	return fmt.Sprintf("%p", m)
}

// UNSAFE!
func (m *SlotMachine) HasPriorityWork() bool {
	return !m.prioritySlots.IsEmpty() || !m.workingSlots.IsEmpty()
}
