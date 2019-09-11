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
	"github.com/insolar/insolar/network/consensus/common/syncrun"
	"github.com/pkg/errors"
	"sync"
	"time"
)

type SlotMachineConfig struct {
	PollingPeriod   time.Duration
	PollingTruncate time.Duration
	SlotPageSize    uint16
}

type DependencyInjector interface {
	InjectDependencies(sm StateMachine, slotID SlotID, container *SlotMachine)
}

func NewSlotMachine(config SlotMachineConfig, injector DependencyInjector, adapters *AdapterRegistry, internalSignal func()) SlotMachine {
	ownsAdapters := false
	if adapters == nil {
		adapters = NewAdapterRegistry()
		ownsAdapters = true
	}
	return SlotMachine{
		config:        config,
		injector:      injector,
		adapters:      adapters,
		ownsAdapters:  ownsAdapters,
		unusedSlots:   NewSlotQueue(UnusedSlots),
		activeSlots:   NewSlotQueue(ActiveSlots),
		prioritySlots: NewSlotQueue(ActiveSlots),
		workingSlots:  NewSlotQueue(WorkingSlots),
		syncQueue:     tools.NewSignalFuncQueue(&sync.Mutex{}, internalSignal),
	}
}

type SlotMachine struct {
	config       SlotMachineConfig
	injector     DependencyInjector
	adapters     *AdapterRegistry
	ownsAdapters bool

	containerState SlotMachineState

	lastSlotID SlotID
	slots      [][]Slot
	slotPgPos  uint16

	migrationCount uint16

	scanCount uint32
	startedAt time.Time

	unusedSlots  SlotQueue
	workingSlots SlotQueue //slots are currently in processing

	activeWaitOnly bool
	activeSlots    SlotQueue //they are are moved to workingSlots on every Scan
	prioritySlots  SlotQueue
	pollingSlots   PollingQueue

	syncQueue    tools.SyncQueue // for detached/async ops, queued functions MUST BE panic-safe
	detachQueues map[SlotID]tools.SyncFuncList

	stepSync StepSyncCatalog
}

func (m *SlotMachine) GetAdapters() *AdapterRegistry {
	return m.adapters
}

func (m *SlotMachine) OccupiedSlotCount() int {
	n := len(m.slots)
	if n == 0 {
		return 0
	}
	return (n-1)*int(m.config.SlotPageSize) + int(m.slotPgPos) - m.unusedSlots.Count()
}

func (m *SlotMachine) AllocatedSlotCount() int {
	return len(m.slots) * int(m.config.SlotPageSize)
}

func (m *SlotMachine) IsZero() bool {
	return m.syncQueue.IsZero()
}

func (m *SlotMachine) IsEmpty() bool {
	return m.OccupiedSlotCount() == 0
}

func (m *SlotMachine) SetContainerState(s SlotMachineState) {
	m.containerState = s
}

func (m *SlotMachine) ScanOnceAsNested(context ExecutionContext) bool {
	worker := context.(*executionContext).worker.StartNested(m.containerState)
	defer worker.FinishNested(m.containerState)

	return m.ScanOnce(worker)
}

func (m *SlotMachine) ScanEventsOnly() (hasUpdates bool) {

	scanTime := time.Now()
	if m.startedAt.IsZero() {
		m.startedAt = scanTime
	}

	return m.scanEvents()
}

func (m *SlotMachine) scanEvents() (hasUpdates bool) {

	syncQ := m.syncQueue.Flush()
	hasUpdates = len(syncQ) > 0

	for _, fn := range syncQ {
		fn() // allows to resync detached
	}

	return true
}

func (m *SlotMachine) ScanOnce(worker SlotWorker) (hasUpdates bool) {

	scanTime := time.Now()

	switch {
	case m.startedAt.IsZero():
		m.startedAt = scanTime
		fallthrough
	case m.workingSlots.IsEmpty():
		m.scanCount++

		m.activeWaitOnly = true
		m.workingSlots.AppendAll(&m.prioritySlots)
		m.workingSlots.AppendAll(&m.activeSlots)
		m.pollingSlots.FilterOut(scanTime, &m.workingSlots)
	default:
		// we were interrupted
		m.workingSlots.PrependAll(&m.prioritySlots)
	}

	hasUpdates = m.scanEvents()

	if m.workingSlots.IsEmpty() {
		return hasUpdates
	}

	m.pollingSlots.PrepareFor(scanTime.Add(m.config.PollingPeriod).Truncate(m.config.PollingTruncate))
	m.scanWorkingSlots(worker, scanTime)

	return true
}

func (m *SlotMachine) scanWorkingSlots(worker SlotWorker, scanStartTime time.Time) {
	for {
		currentSlot := m.workingSlots.First()
		if currentSlot == nil {
			return
		}
		currentSlot.removeFromQueue()
		currentSlot.startWorking(m.scanCount) //, time.Now().Sub(m.startedAt))

		stopNow := false
		var stateUpdate StateUpdate
		var asyncCount uint16

		// TODO consider use of sync.Pool for executionContext if they are allocated on heap

		wasDetached, err := worker.DetachableCall(func(workerCtx WorkerContext) {
			ec := executionContext{machine: m, worker: workerCtx, slotContext: slotContext{s: currentSlot}}
			stopNow, stateUpdate, asyncCount = ec.executeNextStep()
		})

		if err != nil {
			stateUpdate = stateUpdatePanic(err)
		}

		if wasDetached {
			// MUST NOT apply any changes in the current routine, as it is no more considered as safe
			slotLink := currentSlot.NewLink()
			m.applyDetachedStateUpdate(slotLink, stateUpdate, asyncCount)
			return
		}

		currentSlot.stopWorking(asyncCount)

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

	m.scanAndCleanup(cleanupWeak, m.migratePage)
}

func (m *SlotMachine) Cleanup() {
	m.scanAndCleanup(true, m.verifyPage)
}

func (m *SlotMachine) scanAndCleanup(cleanupWeak bool, scanPage func(slotPage []Slot) (isPageEmptyOrWeak, hasWeakSlots bool)) {
	if len(m.slots) == 0 || len(m.slots) == 1 && m.slotPgPos == 0 {
		return
	}

	isAllEmptyOrWeak, hasSomeWeakSlots := scanPage(m.slots[0][:m.slotPgPos])

	j := 1
	for i, slotPage := range m.slots[1:] {
		isPageEmptyOrWeak, hasWeakSlots := scanPage(slotPage)
		switch {
		case !isPageEmptyOrWeak:
			isAllEmptyOrWeak = false
		case !hasWeakSlots:
			cleanupEmptyPage(slotPage)
			m.slots[i+1] = nil
			continue
		default:
			hasSomeWeakSlots = true
		}

		if j != i+1 {
			m.slots[j] = slotPage
			m.slots[i+1] = nil
		}
		j++
	}

	if isAllEmptyOrWeak && (cleanupWeak || !hasSomeWeakSlots) {
		for _, slotPage := range m.slots {
			m.disposePageSlots(slotPage)
		}
		m.slots = m.slots[:1]
		m.slotPgPos = 0
		return
	}

	if len(m.slots) > j {
		m.slots = m.slots[:j]
	}
}

func (m *SlotMachine) disposePageSlots(slotPage []Slot) {
	for i := range slotPage {
		slot := &slotPage[i]
		m.disposeSlot(slot)
	}
}

func cleanupEmptyPage(slotPage []Slot) {
	for i := range slotPage {
		slot := &slotPage[i]
		if slot.QueueType() != UnusedSlots {
			panic("illegal state")
		}
		slot.removeFromQueue()
	}
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

	m.lastSlotID++
	_, link := m.addStateMachine(ctx, nil, m.lastSlotID, parent, sm)
	return link
}

func (m *SlotMachine) AddAsyncNew(ctx context.Context, parent SlotLink, sm StateMachine) {
	if sm == nil {
		panic("illegal value")
	}
	if ctx == nil {
		panic("illegal value")
	}

	m.syncQueue.Add(func() {
		m.lastSlotID++
		m.addStateMachine(ctx, nil, m.lastSlotID, parent, sm)
	})
}

//func (m *SlotMachine) AddSharedState(ctx context.Context, parent SlotLink, ss SharedState) SharedStateAdapter {
//	if ss == nil {
//		panic("illegal value")
//	}
//	if ctx == nil {
//		panic("illegal value")
//	}
//
//	state := &sharedStateMachine{ state: ss }
//
//	m.lastSlotID++
//	ok, link := m.addStateMachine(ctx, nil, m.lastSlotID, parent, state)
//	if !ok {
//		panic("illegal state")
//	}
//	state.accessFn
//	return link
//}

func (m *SlotMachine) addStateMachine(ctx context.Context, slot *Slot, newSlotID SlotID, parent SlotLink,
	sm StateMachine) (bool, SlotLink) {

	smd := sm.GetStateMachineDeclaration()
	if smd == nil {
		panic("illegal value")
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

	ic := initializationContext{slotContext{s: slot}, m}
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

func (m *SlotMachine) allocateSlot() (slot *Slot) {
	switch {
	case !m.unusedSlots.IsEmpty():
		slot = m.unusedSlots.First()
		slot.removeFromQueue()
	case m.slots == nil:
		m.slots = make([][]Slot, 1)
		m.slots[0] = make([]Slot, m.config.SlotPageSize)
		m.slotPgPos = 1
		slot = &m.slots[0][0]
	default:
		lenSlots := len(m.slots[0])
		if int(m.slotPgPos) == lenSlots {
			m.slots = append(m.slots, m.slots[0])
			m.slots[0] = make([]Slot, lenSlots)
			m.slotPgPos = 0
		}
		slot = &m.slots[0][m.slotPgPos]
		m.slotPgPos++
	}
	return slot
}

func (m *SlotMachine) slotAsyncCall(slotLink SlotLink, fn func(*Slot)) {
	m.syncQueue.Add(func() {
		m.slotDirectCall(slotLink, fn)
	})
}

func (m *SlotMachine) slotDirectCall(slotLink SlotLink, fn func(*Slot)) {
	if !slotLink.IsValid() {
		m.getDetachQueue(slotLink.SlotID()) // cleanup

		stateUpdate := stateUpdateExpired(SlotStep{}, fn)
		m.slotAccessError("slot has expired on adapter callback", slotLink, stateUpdate)
		return
	}

	err := safeSlotCall(slotLink.s, fn)
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

func safeSlotCall(slot *Slot, fn func(*Slot)) (recovered error) {
	defer func() {
		recovered = recoverSlotPanic("async result has failed", recover(), recovered)
	}()
	fn(slot)
	return nil
}

func (m *SlotMachine) applyAsyncStateUpdate(link SlotLink, resultFn AsyncResultFunc, recovered interface{}) {

	m.slotAsyncCall(link, func(slot *Slot) {
		if !slot.isWorking() {
			m._applyAsyncStateUpdate(slot, resultFn)
			return
		}
		/* this is an async result for a handler that was detached - we have to postpone it until reattachment */
		if m.detachQueues == nil {
			m.detachQueues = make(map[SlotID]tools.SyncFuncList)
		}

		dq := m.detachQueues[link.SlotID()]
		dq = append(dq, func() {
			m.slotDirectCall(link, func(slot *Slot) {
				m._applyAsyncStateUpdate(slot, resultFn)
			})
		})
		m.detachQueues[link.SlotID()] = dq
	})
}

func (m *SlotMachine) getDetachQueue(slotID SlotID) tools.SyncFuncList {
	dq := m.detachQueues[slotID]
	if dq == nil {
		return nil
	}
	delete(m.detachQueues, slotID)
	return dq
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

func (m *SlotMachine) applyDetachedStateUpdate(slotLink SlotLink, stateUpdate StateUpdate, asyncCount uint16) {
	m.slotAsyncCall(slotLink, func(slot *Slot) {
		slot.stopWorking(asyncCount)
		detachQueue := m.getDetachQueue(slotLink.SlotID())

		if !stateUpdate.IsZero() && !m.applyStateUpdate(slot, false, stateUpdate) {
			//if len(detachQueue) > 0 {
			// TODO log warning on post-mortem update(s)
			//}
			return
		}

		if _, isAvailable := m.migrateSlot(slot); !isAvailable {
			//if len(detachQueue) > 0 {
			// TODO log warning on post-mortem update(s)
			//}
			return
		}

		for _, fn := range detachQueue {
			fn()
		}
	})
}

func (m *SlotMachine) applyStateUpdate(slot *Slot, inplaceUpdate bool, stateUpdate StateUpdate) bool {

	isAvailable, recovered := m._applyStateUpdate(slot, inplaceUpdate, stateUpdate)
	if recovered == nil {
		return isAvailable
	}

	m._handleStateUpdateError(slot, inplaceUpdate, stateUpdate, recovered)
	m.disposeSlot(slot)
	m.unusedSlots.AddLast(slot)
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
		applySlotPrepareAndNextStep(slot, stateUpdate)
		m._applyInplaceUpdate(slot, inplaceUpdate, activateSlot)
		return true, nil

	case stateUpdPoll:
		applySlotPrepareAndNextStep(slot, stateUpdate)
		m._applyInplaceUpdate(slot, inplaceUpdate, deactivateSlot)
		m.pollingSlots.Add(slot)
		return true, nil

	case stateUpdSleep:
		applySlotPrepareAndNextStep(slot, stateUpdate)
		m._applyInplaceUpdate(slot, inplaceUpdate, deactivateSlot)
		return true, nil

	case stateUpdWaitForEvent:
		applySlotPrepareAndNextStep(slot, stateUpdate)
		m._applyInplaceUpdate(slot, inplaceUpdate, activeWaitSlot)
		return true, nil

	case stateUpdWaitForActive:
		applySlotPrepareAndNextStep(slot, stateUpdate)
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
		//applySlotPrepareAndNextStep(slot, stateUpdate)
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
		m.lastSlotID++
		return m.addStateMachine(ctx, slot, m.lastSlotID, parent, sm)

	case stateUpdStop:
		// disposeSlot can handle both in-place and off-place updates
		m.disposeSlot(slot)
		m.unusedSlots.AddLast(slot)
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
	activeWaitSlot
)

func (m *SlotMachine) _applyInplaceUpdate(slot *Slot, inplaceUpdate bool, activation slotActivationMode) {

	if inplaceUpdate && !m.removeSlotQueue(slot, activation == activateSlot) {
		return
	}
	if activation == deactivateSlot {
		return
	}
	m.addSlotToActiveOrWorkingQueue(slot, activation == activeWaitSlot)
}

func applySlotPrepareAndNextStep(slot *Slot, stateUpdate StateUpdate) {
	if stateUpdate.param1 != nil {
		fn := stateUpdate.param1.(StepPrepareFunc) // panic is handled
		fn(slot)
	}
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
	m.lastSlotID++
	cc := constructionContext{ctx: ctx, parent: parent, slotID: m.lastSlotID, machine: m}
	sm := cc.executeCreate(fnCreate)

	if sm != nil {
		smd := sm.GetStateMachineDeclaration()
		if smd != nil {
			return m.addStateMachine(cc.ctx, slot, cc.slotID, cc.parent, sm)
		}
	}

	return false, NoLink()
}

func (m *SlotMachine) slotAccessError(msg string, link SlotLink, update StateUpdate) {
	// TODO logging
	fmt.Printf("%s: %+v %+v", msg, link, update)
}

func (m *SlotMachine) addSlotToActiveOrWorkingQueueNoWait(slot *Slot) {
	m.addSlotToActiveOrWorkingQueue(slot, false)
}

func (m *SlotMachine) addSlotToActiveOrWorkingQueue(slot *Slot, activeWait bool) {
	switch {
	case activeWait:
		break
	case slot.isLastScan(m.scanCount):
		m.activeWaitOnly = false
		break
	default:
		if slot.isPriority() {
			m.workingSlots.AddFirst(slot)
		} else {
			m.workingSlots.AddLast(slot)
		}
		return
	}
	if slot.isPriority() {
		m.prioritySlots.AddLast(slot)
	} else {
		m.activeSlots.AddLast(slot)
	}
}

func (m *SlotMachine) removeSlotQueue(slot *Slot, keepActive bool) bool {
	if !slot.isQueueHead() {
		if keepActive {
			switch slot.QueueType() {
			case ActiveSlots, WorkingSlots:
				return false
			}
		}
		//switch slot.QueueType() {
		slot.removeFromQueue()
		return true
	}

	switch slot.QueueType() {
	case ActivationOfSlot:
		slot.removeHeadedQueue(m.addSlotToActiveOrWorkingQueueNoWait)
	default:
		panic("illegal state")
	}
	return true
}

func (m *SlotMachine) disposeSlot(slot *Slot) {
	dep := slot.dependency
	slot.dependency = nil
	if dep != nil {
		dep.OnSlotDisposed()
	}

	if slot.isQueueHead() {
		slot.removeHeadedQueue(m.addSlotToActiveOrWorkingQueueNoWait)
	} else {
		slot.removeFromQueue()
	}

	slot.dispose()
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
		ec := failureContext{slotContext{s: slot}, m, isPanic, isAsync, slotError}
		err := ec.executeErrorHandlerSafe(handler)
		if err == nil {
			return
		}
		slotError = errors.WithMessage(slotError, err.Error())
	}

	// TODO log unhandled error
	fmt.Printf("slot fail: %v %+v %v", slot.GetID(), stateUpdate, slotError)
}

func NewAdapterCallback(stepLink StepLink, callback AdapterCallbackFunc, cancel *syncrun.ChainedCancel) AdapterCallback {
	return AdapterCallback{stepLink, callback, cancel}
}

type AdapterCallback struct {
	stepLink StepLink
	callback AdapterCallbackFunc
	cancel   *syncrun.ChainedCancel
}

func (c AdapterCallback) IsZero() bool {
	return c.stepLink.IsEmpty()
}

func (c AdapterCallback) IsCancelled() bool {
	return !c.stepLink.IsAtStep() || c.cancel != nil && c.cancel.IsCancelled()
}

func (c AdapterCallback) SendResult(result AsyncResultFunc) {
	if c.IsZero() {
		panic("illegal state")
	}
	_sendResult(c.stepLink, result, c.callback, c.cancel)
}

// just to make sure that outer struct doesn't leak into a closure
func _sendResult(stepLink StepLink, result AsyncResultFunc, callback AdapterCallbackFunc, cancel *syncrun.ChainedCancel) {

	if result == nil {
		// NB! Do NOT ignore "result = nil" - it MUST decrement async call count
		callback(func(ctx AsyncResultContext) {}, nil)
		return
	}

	callback(func(ctx AsyncResultContext) {
		if result == nil || !stepLink.IsAtStep() || cancel != nil && cancel.IsCancelled() {
			return
		}
		result(ctx)
	}, nil)
}

func (c AdapterCallback) SendPanic(recovered interface{}) {
	if c.IsZero() {
		panic("illegal state")
	}
	c.callback(nil, recovered)
}

func (c AdapterCallback) SendCancel() {
	if c.IsZero() {
		panic("illegal state")
	}
	c.callback(nil, nil)
}
