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
	"github.com/insolar/insolar/network/consensus/common/syncrun"
	"sync"
	"time"
)

type SlotMachineConfig struct {
	BeforeDetach  time.Duration
	PollingPeriod time.Duration
	SlotPageSize  uint16
}

type DependencyInjector interface {
	InjectDependencies(sm StateMachine, slotID SlotID, container *SlotMachine)
}

func NewSlotMachine(config SlotMachineConfig, injector DependencyInjector) SlotMachine {
	return SlotMachine{
		config:       config,
		injector:     injector,
		unusedSlots:  NewSlotQueue(UnusedSlots),
		activeSlots:  NewSlotQueue(ActiveSlots),
		workingSlots: NewSlotQueue(WorkingSlots),
		syncQueue:    NewSyncQueue(&sync.Mutex{}),
	}
}

type SlotMachine struct {
	config   SlotMachineConfig
	injector DependencyInjector

	machineState SlotMachineState

	adapters map[AdapterID]*adapterExecHelper

	slotCount SlotID
	slots     [][]Slot
	slotPgPos uint16

	migrationCount uint16

	scanCount uint32
	startedAt time.Time

	unusedSlots  SlotQueue
	workingSlots SlotQueue //slots are currently in processing

	activeSlots  SlotQueue //they are are moved to workingSlots on every Scan
	pollingSlots *SlotQueue

	pollingSeq     []PollingSlotQueue
	pollingSeqHead uint16
	pollingSeqTail uint16

	syncQueue    SyncQueue // for detached/async ops, queued functions MUST BE panic-safe
	detachQueues map[SlotID]SyncFuncList

	stepSync StepSyncCatalog
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

func (m *SlotMachine) IsEmpty() bool {
	return m.syncQueue.IsZero()
}

func (m *SlotMachine) ScanOnceAsNested(context ExecutionContext) bool {
	worker := context.(*executionContext).worker.StartNested(m.machineState)
	defer worker.FinishNested(m.machineState)

	return m.ScanOnce(worker)
}

func (m *SlotMachine) RegisterAdapter(adapterID AdapterID, adapterExecutor AdapterExecutor) ExecutionAdapter {
	if adapterID.IsEmpty() {
		panic("illegal value")
	}
	if adapterExecutor == nil {
		panic("illegal value")
	}

	if m.adapters == nil {
		m.adapters = make(map[AdapterID]*adapterExecHelper)
	}
	if m.adapters[adapterID] != nil {
		panic("duplicate adapter id: " + adapterID)
	}
	adapterExecutor.RegisterOn(m.machineState)
	r := &adapterExecHelper{adapterID, adapterExecutor}
	m.adapters[adapterID] = r

	return r
}

func (m *SlotMachine) GetAdapter(adapterID AdapterID) ExecutionAdapter {
	return m.adapters[adapterID]
}

func (m *SlotMachine) ScanOnce(worker SlotWorker) (hasUpdates bool) {

	scanTime := time.Now()

	if m.startedAt.IsZero() {
		m.startedAt = scanTime
	}
	m.scanCount++

	syncQ := m.syncQueue.Flush()
	hasUpdates = len(syncQ) > 0

	for _, fn := range syncQ {
		fn() // allow to resync detached
	}

	if m.workingSlots.IsEmpty() {
		m.workingSlots.AppendAll(&m.activeSlots)
		m.preparePollingSlots(scanTime)
	}
	m.allocatePollingSlots(scanTime)

	if m.workingSlots.IsEmpty() {
		return hasUpdates
	}

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
			stateUpdate = stateUpdateFailed(err)
		}

		if wasDetached {
			// MUST NOT apply any changes in the current routine, as it is no more considered as safe
			slotLink := currentSlot.NewLink()
			m.applyDetachedStateUpdate(slotLink, stateUpdate, asyncCount)
			return
		}

		currentSlot.stopWorking(asyncCount)

		if !stateUpdate.IsZero() {
			m.applyStateUpdate(currentSlot, stateUpdate)
		}
		if stopNow {
			return
		}
	}
}

func (m *SlotMachine) Migrate() {
	m.migrationCount++

	for _, adapter := range m.adapters {
		adapter.executor.Migrate(m.machineState, m.migrationCount)
	}

	if len(m.slots) == 0 {
		return
	}

	m.migratePage(m.slots[0][:m.slotPgPos])
	for _, slotPage := range m.slots[1:] {
		m.migratePage(slotPage)
	}
}

func (m *SlotMachine) migratePage(slotPage []Slot) {
	for i := range slotPage {
		m.migrateSlot(&slotPage[i])
	}
}

func (m *SlotMachine) AddNew(ctx context.Context, parent SlotLink, sm StateMachine) SlotLink {
	if sm == nil {
		panic("illegal state")
	}
	if ctx == nil {
		panic("illegal value")
	}

	m.slotCount++
	_, link := m.addStateMachine(ctx, nil, m.slotCount, parent, sm)
	return link
}

func (m *SlotMachine) addStateMachine(ctx context.Context, slot *Slot, newSlotID SlotID, parent SlotLink, sm StateMachine) (bool, SlotLink) {
	smd := sm.GetStateMachineDeclaration()
	if smd == nil {
		panic("illegal state")
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

	slot.init(ctx, slotID, parent, smd, m.machineState)
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

	return m.applyStateUpdate(slot, stateUpdate)
}

func (m *SlotMachine) migrateSlot(slot *Slot) bool {
	if slot.isEmpty() || slot.isWorking() {
		return false
	}

	if m.migrationCount < slot.migrationCount {
		panic("illegal state")
	}

	for m.migrationCount != slot.migrationCount {
		migrateFn := slot.step.Migration
		if migrateFn == nil {
			migrateFn = slot.declaration.GetMigrateFn(slot.step.Transition)
			if migrateFn == nil {
				migrateFn = slot.defMigrate
			}
		}
		if migrateFn == nil {
			slot.migrationCount = m.migrationCount
			break
		}

		mc := migrationContext{slotContext{s: slot}}
		stateUpdate := mc.executeMigrate(migrateFn)

		if !m.applyStateUpdate(slot, stateUpdate) {
			return false // slot was stopped
		}

		slot.migrationCount++
	}
	return true
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

func (m *SlotMachine) syncSafe(slotLink SlotLink, fn func(*Slot)) {
	m.syncQueue.Add(func() {
		m.safeWrapAndHandleError(slotLink, fn)
	})
}

func (m *SlotMachine) safeWrapAndHandleError(slotLink SlotLink, fn func(*Slot)) {
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

	stateUpdate := stateUpdateFailed(err)
	m.slotAccessError("adapter callback panic", slotLink, stateUpdate)
	m.applyStateUpdate(slotLink.s, stateUpdate)
}

func recoverToErr(msg string, recovered interface{}, defErr error) error {
	if recovered == nil {
		return defErr
	}
	return fmt.Errorf("%s: %v", msg, recovered)
}

func safeSlotCall(slot *Slot, fn func(*Slot)) (err error) {
	defer func() {
		err = recoverToErr("async result has failed", recover(), err)
	}()
	fn(slot)
	return nil
}

func (m *SlotMachine) applyAsyncStateUpdate(link SlotLink, resultFn AsyncResultFunc, recovered interface{}) {

	m.syncSafe(link, func(slot *Slot) {
		if !slot.isWorking() {
			m._applyAsyncStateUpdate(slot, resultFn)
			return
		}
		/* this is an async result for a handler that was detached - we have to postpone it until reattachment */
		if m.detachQueues == nil {
			m.detachQueues = make(map[SlotID]SyncFuncList)
		}

		dq := m.detachQueues[link.SlotID()]
		dq = append(dq, func() {
			m.safeWrapAndHandleError(link, func(slot *Slot) {
				m._applyAsyncStateUpdate(slot, resultFn)
			})
		})
		m.detachQueues[link.SlotID()] = dq
	})
}

func (m *SlotMachine) getDetachQueue(slotID SlotID) SyncFuncList {
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

	switch slot.QueueType() {
	case ActiveSlots, WorkingSlots:
		// do nothing
	//case AnotherSlot:
	//	// can't wake up?
	case UnusedSlots:
		panic("illegal state")
	default:
		slot.removeFromQueue()
		m.activeSlots.AddLast(slot)
	}
}

func (m *SlotMachine) applyDetachedStateUpdate(slotLink SlotLink, stateUpdate StateUpdate, asyncCount uint16) {
	m.syncSafe(slotLink, func(slot *Slot) {
		slot.stopWorking(asyncCount)
		detachQueue := m.getDetachQueue(slotLink.SlotID())

		if !stateUpdate.IsZero() && (!m.applyStateUpdate(slot, stateUpdate) || !m.migrateSlot(slot)) {
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

func (m *SlotMachine) applyStateUpdate(slot *Slot, stateUpdate StateUpdate) bool {

	stillActive, recovered := m._applyStateUpdate(slot, stateUpdate)
	if recovered == nil {
		return stillActive
	}
	// TODO report recovered error

	fmt.Printf("slot fail: %v %+v %+v", slot.GetID(), stateUpdate, stateUpdate.param1)
	m.disposeSlot(slot)
	m.unusedSlots.AddLast(slot)
	return false
}

func (m *SlotMachine) _applyStateUpdate(slot *Slot, stateUpdate StateUpdate) (stillActive bool, recovered interface{}) {

	if slot.isWorking() { // there must be no calls here from inside of custom handlers
		return false, "illegal state / slot is working"
	}

	defer func() {
		r := recover()
		switch {
		case r == nil:
			break
		case recovered == nil:
			recovered = r
		default:
			recovered = fmt.Sprintf("recovered: before={%v} after={%v}", recovered, r)
		}
	}()

	su := stateUpdType(stateUpdate.updType)

	if su.HasStep() {
		if su.HasPrepare() && stateUpdate.param1 != nil {
			fn := stateUpdate.param1.(StepPrepareFunc) // panic is handled
			fn(slot)
		}
		slot.setNextStep(stateUpdate.step)
		if slot.step.Transition == nil {
			return false, "missing transition"
		}
	}

	switch su {
	case stateUpdRepeat, stateUpdNext, stateUpdNextLoop:
		m.addSlotToActiveOrWorkingQueue(slot)
		return true, nil

	case stateUpdPoll:
		m.pollingSlots.AddLast(slot)
		return true, nil

	case stateUpdWait:
		waitOn := stateUpdate.getLink()
		switch {
		case waitOn.s == slot || !waitOn.IsValid():
			// don't wait for self
			// don't wait for an expired slot
			m.addSlotToActiveOrWorkingQueue(slot)
		default:
			switch waitOn.s.QueueType() {
			case ActiveSlots, WorkingSlots:
				// don't wait
				m.addSlotToActiveOrWorkingQueue(slot)
			case NoQueue:
				waitOn.s.makeQueueHead()
				fallthrough
			case AnotherSlotQueue, PollingSlots:
				waitOn.s.queue.AddLast(slot)
			default:
				return false, "illegal slot queue"
			}
		}
		return true, nil

	case stateUpdReplace:
		cf := stateUpdate.param1.(CreateFunc) // panic is handled
		parent := slot.parent
		ctx := slot.ctx
		m.disposeSlot(slot)
		ok, _ := m.applySlotCreate(ctx, slot, parent, cf) // NB! recursive call inside
		return ok, nil

	case stateUpdStop:
		m.disposeSlot(slot)
		m.unusedSlots.AddLast(slot)
		return false, nil

	case 0, stateUpdNoChange, stateUpdExpired:
		// can't be here
		return false, "unexpected state update"

	case stateUpdDispose:
		return false, "stateUpdDispose"

	default:
		return false, "unknown update type"
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
			stateUpdate := bc.executeBargeIn(applyFn).ensureMarker(&bc.marker)

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
	m.slotCount++
	cc := constructionContext{ctx: ctx, parent: parent, slotID: m.slotCount}
	sm := cc.executeCreate(fnCreate)

	return m.addStateMachine(cc.ctx, slot, cc.slotID, cc.parent, sm)
}

func (m *SlotMachine) slotAccessError(msg string, link SlotLink, update StateUpdate) {
	// TODO logging
	fmt.Printf("%s: %+v %+v", msg, link, update)
}

func (m *SlotMachine) addSlotToActiveOrWorkingQueue(slot *Slot) {
	if slot.isLastScan(m.scanCount) {
		m.activeSlots.AddLast(slot)
	} else {
		m.workingSlots.AddLast(slot)
	}
}

func (m *SlotMachine) disposeSlot(slot *Slot) {
	dep := slot.dependency
	slot.dependency = nil

	if slot.dependency != nil {
		dep.OnSlotDisposed()
	}

	if slot.QueueType() == AnotherSlotQueue && slot.isQueueHead() {
		for {
			next := slot.QueueNext()
			if next == nil {
				break
			}
			next.removeFromQueue()
			m.addSlotToActiveOrWorkingQueue(next)
		}
		slot.vacateQueueHead()
	} else {
		slot.removeFromQueue()
	}

	slot.dispose()
}

func (m *SlotMachine) preparePollingSlots(scanTime time.Time) {

	for m.pollingSeqHead != m.pollingSeqTail { // FIXME it won't pick the only non-empty polling slot queue
		ps := &m.pollingSeq[m.pollingSeqTail]

		if !ps.IsEmpty() && ps.pollingTime.After(scanTime) {
			break
		}

		m.pollingSeqTail++
		if int(m.pollingSeqTail) >= len(m.pollingSeq) {
			m.pollingSeqTail = 0
		}

		m.workingSlots.AppendAll(&ps.SlotQueue)
	}
}

func (m *SlotMachine) allocatePollingSlots(scanTime time.Time) {
	var pollingQueue *PollingSlotQueue
	switch {
	case m.pollingSlots == nil:
		if m.pollingSeqHead != 0 {
			panic("illegal state")
		}
		if len(m.pollingSeq) == 0 {
			m.growPollingSlots()
		}
	case !m.pollingSlots.IsEmpty():
		m.growPollingSlots()
		m.pollingSeqHead++
		if int(m.pollingSeqHead) >= len(m.pollingSeq) {
			m.pollingSeqHead = 0
		}
	}
	pollingQueue = &m.pollingSeq[m.pollingSeqHead]

	if !pollingQueue.SlotQueue.IsEmpty() {
		panic("illegal state")
	}

	m.pollingSlots = &pollingQueue.SlotQueue
	pollingQueue.pollingTime = scanTime.Add(m.config.PollingPeriod)
}

func (m *SlotMachine) growPollingSlots() {
	switch {
	case m.pollingSeqHead+1 == m.pollingSeqTail:
		// full
		sLen := len(m.pollingSeq)

		cp := make([]PollingSlotQueue, sLen, 1+(sLen<<2)/3)
		copy(cp, m.pollingSeq[m.pollingSeqTail:])
		copy(cp[m.pollingSeqTail:], m.pollingSeq[:m.pollingSeqTail])
		m.pollingSeq = cp

		m.pollingSeqTail = 0
		m.pollingSeqHead = uint16(sLen - 1)
		fallthrough
	case m.pollingSeqTail == 0 && int(m.pollingSeqHead)+1 >= len(m.pollingSeq):
		// full
		for {
			m.pollingSeq = append(m.pollingSeq, PollingSlotQueue{SlotQueue: NewSlotQueue(PollingSlots)})
			if len(m.pollingSeq) == cap(m.pollingSeq) {
				break
			}
		}
	}
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
