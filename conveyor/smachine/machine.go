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
	"github.com/insolar/insolar/network/consensus/common/syncrun"
	"sync"
	"time"
)

type SlotMachineConfig struct {
	BeforeDetach  time.Duration
	PollingPeriod time.Duration
	SlotPageSize  uint16
}

func NewSlotMachine(config SlotMachineConfig) SlotMachine {
	return SlotMachine{
		config:       config,
		unusedSlots:  NewSlotQueue(UnusedSlots),
		activeSlots:  NewSlotQueue(ActiveSlots),
		workingSlots: NewSlotQueue(WorkingSlots),
		syncQueue:    NewSyncQueue(&sync.Mutex{}),
	}
}

type SlotMachine struct {
	config SlotMachineConfig

	adapters map[AdapterID]AdapterExecutor

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

func (m *SlotMachine) IsEmpty() bool {
	return m.syncQueue.IsZero()
}

func (m *SlotMachine) GetAdapterQueue(adapter ExecutionAdapter) AdapterExecutor {
	return m.adapters[adapter.GetAdapterID()]
}

func (m *SlotMachine) ScanOnceAsNested(context ExecutionContext) bool {
	workCtl := context.(*executionContext).worker.workCtl
	return m.ScanOnce(workCtl)
}

func (m *SlotMachine) RegisterAdapter(adapterID AdapterID, adapterExecutor AdapterExecutor) ExecutionAdapter {
	if adapterID.IsEmpty() {
		panic("illegal value")
	}
	if adapterExecutor == nil {
		panic("illegal value")
	}

	if m.adapters == nil {
		m.adapters = make(map[AdapterID]AdapterExecutor)
	}
	if m.adapters[adapterID] != nil {
		panic("duplicate adapter id: " + adapterID)
	}
	m.adapters[adapterID] = adapterExecutor

	return &adapterExecHelper{adapterID, adapterExecutor}
}

func (m *SlotMachine) ScanOnce(workCtl WorkerController) (hasUpdates bool) {

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

	m.scanWorkingSlots(workCtl, scanTime)

	return true
}

func (m *SlotMachine) scanWorkingSlots(workCtl WorkerController, scanTime time.Time) {
	sw := SlotWorker{workCtl: workCtl, machine: m}

	for {
		currentSlot := m.workingSlots.First()
		if currentSlot == nil {
			return
		}
		currentSlot.removeFromQueue()
		slotLink := currentSlot.NewLink()

		currentSlot.startWorking(m.scanCount, time.Now().Sub(m.startedAt))

		stopNow := false
		var stateUpdate StateUpdate
		var asyncCount uint16

		wasDetached, err := sw.detachableCall(func() {
			ec := executionContext{worker: &sw, slotContext: slotContext{s: currentSlot}}
			stopNow, stateUpdate, asyncCount = ec.executeNextStep()
		})

		if err != nil {
			stateUpdate = stateUpdateFailed(err)
		}

		if wasDetached {
			// MUST NOT apply any changes in the current routine, as it is no more considered as safe
			m.applyDetachedStateUpdate(slotLink, stateUpdate, asyncCount)
			return
		}

		currentSlot.stopWorking(asyncCount)
		m.applyStateUpdate(currentSlot, stateUpdate)

		if stopNow {
			return
		}
	}
}

func (m *SlotMachine) Migrate() {
	m.migrationCount++

	if len(m.slots) == 0 {
		return
	}

	m.migratePage(m.slots[0][:m.slotPgPos])
	for _, slotPage := range m.slots[1:] {
		m.migratePage(slotPage)
	}

	// TODO inform adapters
}

func (m *SlotMachine) migratePage(slotPage []Slot) {
	for i := range slotPage {
		m.migrate(&slotPage[i])
	}
}

func (m *SlotMachine) AddNew(parent SlotLink, sm StateMachine) SlotLink {
	if sm == nil {
		panic("illegal state")
	}

	m.slotCount++
	_, link := m.addStateMachine(nil, m.slotCount, parent, sm)
	return link
}

func (m *SlotMachine) addStateMachine(slot *Slot, newSlotID SlotID, parent SlotLink, sm StateMachine) (bool, SlotLink) {
	smd := sm.GetStateMachineDeclaration()
	if smd == nil {
		panic("illegal state")
	}

	if slot == nil {
		slot = m.allocateSlot()
	}
	m.prepareSlot(slot, newSlotID, parent, smd)
	link := slot.NewLink() // should happen BEFORE start as slot can die immediately

	return m.startSlot(slot, sm), link
}

func (m *SlotMachine) prepareSlot(slot *Slot, slotID SlotID, parent SlotLink, smd StateMachineDeclaration) {
	if smd == nil {
		panic("illegal state")
	}

	slot.init(slotID, parent, smd)
	slot.migrationCount = m.migrationCount
	slot.lastWorkScan = uint8(m.scanCount)
}

func (m *SlotMachine) startSlot(slot *Slot, sm StateMachine) bool {
	initState := slot.machine.GetInitStateFor(sm)
	if initState == nil {
		panic("illegal state")
	}

	ic := initializationContext{slotContext{s: slot}}
	stateUpdate := ic.executeInitialization(initState)

	return m.applyStateUpdate(slot, stateUpdate)
}

func (m *SlotMachine) migrate(slot *Slot) bool {
	if slot.isEmpty() || slot.isWorking() {
		return false
	}

	if m.migrationCount < slot.migrationCount {
		panic("illegal state")
	}

	for m.migrationCount != slot.migrationCount {
		migrateFn := slot.step.Migration
		if migrateFn == nil {
			migrateFn = slot.machine.GetMigrateFn(slot.step.Transition)
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

func (m *SlotMachine) applyAsyncStateUpdate(stepLink StepLink, resultFn AsyncResultFunc, recovered interface{}) {

	m.syncSafe(stepLink.SlotLink, func(slot *Slot) {
		if !slot.isWorking() {
			m._applyAsyncStateUpdate(slot, resultFn)
			return
		}
		/* this is an async result for a handler that was detached - we have to postpone it until reattachment */
		if m.detachQueues == nil {
			m.detachQueues = make(map[SlotID]SyncFuncList)
		}

		dq := m.detachQueues[stepLink.SlotID()]
		dq = append(dq, func() {
			m.safeWrapAndHandleError(stepLink.SlotLink, func(slot *Slot) {
				m._applyAsyncStateUpdate(slot, resultFn)
			})
		})
		m.detachQueues[stepLink.SlotID()] = dq
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

		if !m.applyStateUpdate(slot, stateUpdate) || !m.migrate(slot) {
			// TODO log warning on missed update(s)
			//if len(detachQueue) > 0 {
			//}
			return
		}

		for _, fn := range detachQueue {
			fn()
		}
	})
}

func (m *SlotMachine) applyStateUpdate(slot *Slot, stateUpdate StateUpdate) bool {

	if slot.isWorking() { // there must be no calls here from inside of custom handlers
		panic("illegal state")
	}
	slot.incStep()

	updType, slotStep, param := ExtractStateUpdate(stateUpdate)
	su := stateUpdType(updType)

	if su.HasStep() {
		switch {
		case su.HasPrepare():
			recovered := runStepPrepareFn(param, slot)
			if recovered != nil {
				su = stateUpdDispose
				// TODO report recovered error
				break
			}
			fallthrough
		default:
			if slotStep.Transition != nil {
				slot.setNextStep(slotStep)
			}
		}
	}

	switch su {
	case stateUpdNoChange, stateUpdExpired:
		// can't be here
		panic("illegal value")

	case stateUpdNext, stateUpdNextLoop, stateUpdRepeat:
		if slot.step.Transition == nil {
			break
		}
		m.addSlotToActiveOrWorkingQueue(slot)
		return true

	case stateUpdPoll:
		if slot.step.Transition == nil {
			break
		}
		m.pollingSlots.AddLast(slot)
		return true

	case stateUpdWait:
		if slot.step.Transition == nil {
			break
		}

		waitOn := stateUpdate.link.s
		switch {
		case waitOn == slot:
			// don't wait for self
			fallthrough
		case !stateUpdate.link.IsValid():
			// don't wait for an expired slot
			m.addSlotToActiveOrWorkingQueue(slot)
			return true
		default:
			switch waitOn.QueueType() {
			case ActiveSlots, WorkingSlots:
				// don't wait
				m.addSlotToActiveOrWorkingQueue(slot)
				return true
			case NoQueue:
				waitOn.makeQueueHead()
				fallthrough
			case AnotherSlotQueue, PollingSlots:
				waitOn.queue.AddLast(slot)
				return true
			default:
				panic("illegal state")
			}
		}
		return true

	case stateUpdReplace:
		parent := slot.parent
		m.disposeSlot(slot)
		ok, _ := m.applySlotCreate(slot, parent, param.(CreateFunc)) // recursive call inside
		return ok

	case stateUpdStop:
		m.disposeSlot(slot)
		m.unusedSlots.AddLast(slot)
		return false

	case stateUpdDispose:
		break

	default:
		panic("illegal state")
	}

	fmt.Printf("slot fail: %v %+v %+v", slot.GetID(), stateUpdate, stateUpdate.param)
	m.disposeSlot(slot)
	m.unusedSlots.AddLast(slot)
	return false
}

func (m *SlotMachine) applySlotCreate(slot *Slot, parent SlotLink, fnCreate CreateFunc) (bool, SlotLink) {
	m.slotCount++
	cc := constructionContext{parent: parent, slotID: m.slotCount}
	sm := cc.executeCreate(fnCreate)

	return m.addStateMachine(slot, cc.slotID, cc.parent, sm)
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
	// TODO inform adapters

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

type AdapterCallback struct {
	stepLink StepLink
	callback AdapterCallbackFunc
	cancel   *syncrun.ChainedCancel
}

func (c AdapterCallback) IsZero() bool {
	return c.stepLink.IsEmpty() || c.callback == nil
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

	// NB! Do NOT ignore "result = nil" - it MUST decrement async call count
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
