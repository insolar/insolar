///
// Modified BSD 3-Clause Clear License
//
// Copyright (c) 2019 Insolar Technologies GmbH
//
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without modification,
// are permitted (subject to the limitations in the disclaimer below) provided that
// the following conditions are met:
//  * Redistributions of source code must retain the above copyright notice, this list
//    of conditions and the following disclaimer.
//  * Redistributions in binary form must reproduce the above copyright notice, this list
//    of conditions and the following disclaimer in the documentation and/or other materials
//    provided with the distribution.
//  * Neither the name of Insolar Technologies GmbH nor the names of its contributors
//    may be used to endorse or promote products derived from this software without
//    specific prior written permission.
//
// NO EXPRESS OR IMPLIED LICENSES TO ANY PARTY'S PATENT RIGHTS ARE GRANTED
// BY THIS LICENSE. THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS
// AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES,
// INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY
// AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL
// THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT,
// INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING,
// BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS
// OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
// ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
//
// Notwithstanding any other provisions of this license, it is prohibited to:
//    (a) use this software,
//
//    (b) prepare modifications and derivative works of this software,
//
//    (c) distribute this software (including without limitation in source code, binary or
//        object code form), and
//
//    (d) reproduce copies of this software
//
//    for any commercial purposes, and/or
//
//    for the purposes of making available this software to third parties as a service,
//    including, without limitation, any software-as-a-service, platform-as-a-service,
//    infrastructure-as-a-service or other similar online service, irrespective of
//    whether it competes with the products or services of Insolar Technologies GmbH.
///

package smachine

import (
	"fmt"
	"math"
	"sync"
	"time"
)

type SlotMachineConfig struct {
	BeforeDetach      time.Duration
	DefaultPollPeriod time.Duration
	SlotPageSize      uint16
}

func NewSlotMachine(config SlotMachineConfig) SlotMachine {
	return SlotMachine{
		config:       config,
		unusedSlots:  NewQueueHead(UnusedList),
		activeSlots:  NewQueueHead(ActiveList),
		workingSlots: NewQueueHead(ActiveList),
		pollingSlots: NewQueueHead(PollingList),
		syncQueue:    NewSyncQueue(&sync.Mutex{}),
	}
}

type SlotMachine struct {
	config SlotMachineConfig

	adapters map[AdapterID]ExecutionAdapterSink

	slotCount SlotID
	slots     [][]Slot
	slotPgPos uint16

	migrationCount uint16

	unusedSlots  QueueHead
	workingSlots QueueHead //slots are currently in processing

	activeSlots  QueueHead //they are are moved to workingSlots on every Scan
	pollingSlots QueueHead
	nextPollTime time.Time

	syncQueue    SyncQueue // for detached/async ops, queued functions MUST BE panic-safe
	detachQueues map[SlotID]SyncFuncList
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

func (m *SlotMachine) GetAdapterQueue(adapter ExecutionAdapter) ExecutionAdapterSink {
	return m.adapters[adapter.GetAdapterID()]
}

func (m *SlotMachine) ScanOnce(workCtl WorkerController) (hasUpdates bool) {

	syncQ := m.syncQueue.Flush()
	hasUpdates = len(syncQ) > 0

	for _, fn := range syncQ {
		fn() // allow to resync detached
	}

	if m.workingSlots.IsEmpty() {
		m.workingSlots.AppendAll(&m.activeSlots)

		if !m.nextPollTime.IsZero() && m.nextPollTime.Before(time.Now()) {
			m.flushPollingSlotsTo(&m.workingSlots)
		}
	}

	if m.workingSlots.IsEmpty() {
		return hasUpdates
	}

	m.scanWorkingSlots(workCtl)
	return true
}

func (m *SlotMachine) scanWorkingSlots(workCtl WorkerController) {
	sw := SlotWorker{workCtl: workCtl, machine: m}

	for {
		currentSlot := m.workingSlots.First()
		if currentSlot == nil {
			return
		}
		currentSlot.remove()
		slotLink := currentSlot.NewLink()

		stopNow := false
		var stateUpdate StateUpdate

		currentSlot.setWorking()
		wasDetached, err := sw.detachableCall(func() {
			ec := executionContext{worker: &sw, slotContext: slotContext{s: currentSlot}}
			var asyncCount uint32
			stopNow, stateUpdate, asyncCount = ec.executeNextStep()
			currentSlot.asyncCallCount += asyncCount
		})

		if err != nil {
			stateUpdate = StateUpdate{flags: stateUpdateFailed, param: err}
		}

		if wasDetached {
			// MUST NOT apply any changes in the current routine, as it is no more considered as safe
			m.applyDetachedStateUpdate(slotLink, stateUpdate)
			return
		}

		currentSlot.setNotWorking()
		m.applyStateUpdate(currentSlot, stateUpdate)

		if stopNow {
			return
		}
	}
}

func (m *SlotMachine) flushPollingSlotsTo(target *QueueHead) {
	m.nextPollTime = time.Time{}
	target.AppendAll(&m.pollingSlots)
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

	m.flushPollingSlotsTo(&m.activeSlots)
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
	if slot.isWorking() {
		return true
	}

	if m.migrationCount < slot.migrationCount {
		panic("illegal state")
	}

	for m.migrationCount != slot.migrationCount {
		migrateFn := slot.nextState.migration
		if migrateFn == nil {
			migrateFn = slot.machine.GetMigrateFn(slot.nextState.transition)
			if migrateFn == nil {
				migrateFn = slot.migrateSlot
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
		slot.remove()
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

func (m *SlotMachine) beforeSlotStop(slot *Slot) {
	// TODO inform adapters
}

func (m *SlotMachine) setNextPollingTime(nextPollingOfSlot time.Time) {
	if m.nextPollTime.IsZero() {
		if m.config.DefaultPollPeriod <= 0 || m.config.DefaultPollPeriod >= math.MaxInt64 {
			if !nextPollingOfSlot.IsZero() {
				m.nextPollTime = nextPollingOfSlot
			}
			return
		}
		m.nextPollTime = time.Now().Add(m.config.DefaultPollPeriod)
	}

	if !nextPollingOfSlot.IsZero() && m.nextPollTime.After(nextPollingOfSlot) {
		m.nextPollTime = nextPollingOfSlot
	}
}

func (m *SlotMachine) syncSafe(slotLink SlotLink, fn func(*Slot)) {
	m.syncQueue.Add(func() {
		m.safeWrapAndHandleError(slotLink, fn)
	})
}

func (m *SlotMachine) safeWrapAndHandleError(slotLink SlotLink, fn func(*Slot)) {
	slot := slotLink.Get()
	if slot == nil {
		m.getDetachQueue(slotLink.SlotID()) // cleanup

		stateUpdate := StateUpdate{flags: stateUpdateNoChange, param: fn}
		m.slotAccessError("slot has expired on adapter callback", slotLink, stateUpdate)
		return
	}

	err := safeSlotCall(slot, fn)
	if err == nil {
		return
	}

	stateUpdate := StateUpdate{flags: stateUpdateFailed, param: err}
	m.slotAccessError("adapter callback panic", slot.NewLink(), stateUpdate)
	m.applyStateUpdate(slot, stateUpdate)
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

func (m *SlotMachine) applyAsyncStateUpdate(slotLink SlotLink, resultFn AsyncResultFunc) {

	m.syncSafe(slotLink, func(slot *Slot) {
		if !slot.isWorking() {
			m._applyAsyncStateUpdate(slot, resultFn)
			return
		}
		/* this is an async result for a handler that was detached - we have to postpone it until reattachment */
		if m.detachQueues == nil {
			m.detachQueues = make(map[SlotID]SyncFuncList)
		}

		dq := m.detachQueues[slot.slotID]
		dq = append(dq, func() {
			m.safeWrapAndHandleError(slotLink, func(slot *Slot) {
				m._applyAsyncStateUpdate(slot, resultFn)
			})
		})
		m.detachQueues[slot.slotID] = dq
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

	switch slot.DependencyType() {
	case ActiveList:
		// do nothing
	//case AnotherSlot:
	//	// can't wake up?
	case UnusedList:
		panic("illegal state")
	default:
		slot.remove()
		m.activeSlots.AddLast(slot)
	}
}

func (m *SlotMachine) applyDetachedStateUpdate(slotLink SlotLink, stateUpdate StateUpdate) {
	m.syncSafe(slotLink, func(slot *Slot) {
		slot.setNotWorking()
		detachQueue := m.getDetachQueue(slotLink.SlotID())

		if !m.applyStateUpdate(slot, stateUpdate) {
			return
		}
		if !m.migrate(slot) {
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

	switch stateUpdate.getMode() {
	case stateUpdateNoChange:
		// not for here
		panic("illegal value")

	case stateUpdateNext:
		if stateUpdate.nextStep.transition == nil {
			panic("illegal state")
		}
		slot.nextState = stateUpdate.nextStep
		m.activeSlots.AddLast(slot)
		return true
	case stateUpdateRepeat:
		if slot.nextState.transition == nil {
			panic("illegal state")
		}
		m.activeSlots.AddLast(slot)
		return true
	case stateUpdateReplace:
		fn := stateUpdate.getCreateFn()
		if fn == nil {
			panic("illegal state")
		}

		parent := slot.parent
		m.beforeSlotStop(slot)
		slot.dispose()

		ok, _ := m.applySlotCreate(slot, parent, fn) // recursive call inside
		return ok
	case stateUpdateStop, stateUpdateFailed:
		if !slot.isEmpty() {
			m.beforeSlotStop(slot)
			slot.dispose()
		}
		slot.remove()
		m.unusedSlots.AddLast(slot)
		return false
	case stateUpdateHotWait:
		m.setNextPollingTime(stateUpdate.wakeupTime)
		m.pollingSlots.AddLast(slot)
		return true
	case stateUpdateColdWait:
		// don't add anywhere
		return true
	default:
		panic("illegal state")
	}
}

func (m *SlotMachine) applySlotCreate(slot *Slot, parent SlotLink, fnCreate CreateFunc) (bool, SlotLink) {
	m.slotCount++
	cc := constructionContext{parent: parent, slotID: m.slotCount}
	sm := cc.executeCreate(fnCreate)

	return m.addStateMachine(slot, cc.slotID, cc.parent, sm)
}

func (m *SlotMachine) slotAccessError(msg string, link SlotLink, update StateUpdate) {
	// TODO logging
}

func (m *SlotMachine) ScanOnceAsNested(context ExecutionContext) bool {
	workCtl := context.(*executionContext).worker.workCtl
	return m.ScanOnce(workCtl)
}
