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
	"sync/atomic"
)

type Slot struct {
	idAndStep uint64       //atomic access
	machine   *SlotMachine // set only once

	/* -----------------------------------
	   Slot fields to support processing queues
	   -----------------------------------
	   SYNC: these portion of slot can be accessed:
		- if queue is assigned - by the goroutine owning the queue's head
	    - if queue is unassigned - by the goroutine of the machine
	*/
	prevInQueue *Slot
	nextInQueue *Slot
	queue       *QueueHead

	/* SYNC: this portion of slot can ONLY be accessed by
	- the same goroutine that either has set BUSY
	- or for non-BUSY - by the goroutine of the machine
	*/
	slotData
}

type slotDeclarationData struct {
	declaration StateMachineDeclaration

	shadowMigrate   ShadowMigrateFunc
	stepLogger      StepLoggerFunc
	defMigrate      MigrateFunc
	defErrorHandler ErrorHandlerFunc
	defTerminate    TerminationHandlerFunc
	defResult       interface{}

	slotReplaceData
}

// transferred with Replace()
type slotReplaceData struct {
	parent SlotLink
	ctx    context.Context

	// DO NOT modify
	inheritable map[string]interface{}

	defFlags StepFlags
}

func (v slotReplaceData) takeOutForReplace() slotReplaceData {
	return v
}

type slotData struct {
	slotDeclarationData

	slotFlags      slotFlags
	lastWorkScan   uint8  // to check if a slot was executed in this cycle
	asyncCallCount uint16 // pending calls, overflow panics
	migrationCount uint32 // can be wrapped by overflow

	step SlotStep

	dependency SlotDependency
}

type slotFlags uint8

const (
	slotWokenUp slotFlags = 1 << iota
	slotHasBargeIn
	slotHasAliases
	slotHadAsync
)

type SlotDependency interface {
	IsReleaseOnWorking() bool
	IsReleaseOnStepping() bool

	Release() (SlotDependency, []PostponedDependency, []StepLink)
	ReleaseAll() ([]PostponedDependency, []StepLink)
}

type PostponedDependency interface {
	PostponedActivate(appendTo []StepLink) []StepLink
}

const (
	slotFlagBusyShift = 32 + iota
	stepIncrementShift
)

const stepIncrement uint64 = 1 << stepIncrementShift
const slotFlagBusy uint64 = 1 << slotFlagBusyShift
const slotIdMask = slotFlagBusy - 1

//see also numberOfReservedSteps

/*
	Slot and step numbers are interpreted together with isBusy flag:
	[step] [busy] [slotId]
	[  0 ] [  0 ] [    0 ] slot is unused and is owned by a slot pool
	[  N ] [  0 ] [   >0 ] slot is in use and is not owned by a slot pool
	[  N ] [  1 ] [    N ] slot is in use and is not owned by a slot pool

	[  0 ] [  0 ] [   >0 ] illegal/reserved state
	[  0 ] [  1 ] [   >0 ] illegal/reserved state
	[ >0 ] [  0 ] [    0 ] illegal/reserved state


	[ >0 ] [  0 ] [   >0 ] slot is allocated and idle (only slot machine can access it); includes:
	 [ 1 ] [  0 ] [   >0 ] slot is constructed and awaits initialization
	 [>1 ] [  0 ] [   >0 ] slot is operational and idle

	[ >0 ] [  1 ] [   >0 ] slot is allocated and busy (only busy-holder can access it, can't be migrated); includes:
	 [ 1 ] [  1 ] [   >0 ] slot is allocated and is being constructed/initialized
	 [>1 ] [  1 ] [   >0 ] initialized slot is busy / working
	 [ 1 ] [  1 ] [    0 ] uninitialized slot is being finalized
	 [>1 ] [  1 ] [    0 ] initialized slot is being finalized


	Usual slot lifecycle is:
	[  0 ] [  0 ] [    0 ] slot is unused and is owned by a slot pool
	[  1 ] [  1 ] [   >0 ] slot is allocated and is being constructed/initialized (can't be migrated)
	[  1 ] [  0 ] [   >0 ] slot is constructed and awaits initialization (can't be migrated)
    [ >1 ] [  N ] [   >0 ] slot is operational (idle / working)
    [ >1 ] [  1 ] [    0 ] slot is being finalized (can't be migrated)
	[  0 ] [  0 ] [    0 ] slot is unused and is owned by a slot pool

	Step number is a cyclic incrementing counter with reserved values:
	= 0 - slot is is just allocated (+busy flag) or is not used by a state machine
	= 1 - slot is initializing, can only appear once for a state machine
	On overflow, step will change to =numberOfReservedSteps
*/

func (s *Slot) _getState() (id SlotID, step uint32, isBusy bool) {
	v := atomic.LoadUint64(&s.idAndStep)
	return SlotID(v), uint32(v >> stepIncrementShift), v&slotFlagBusy != 0
}

func (s *Slot) GetState() (id SlotID, step uint32, isBusy bool) {
	v := atomic.LoadUint64(&s.idAndStep)
	if SlotID(v).IsUnknown() || v <= slotFlagBusy {
		panic("illegal state")
	}
	return SlotID(v), uint32(v >> stepIncrementShift), v&slotFlagBusy != 0
}

func (s *Slot) GetSlotID() SlotID {
	v := atomic.LoadUint64(&s.idAndStep)
	if SlotID(v).IsUnknown() || v <= slotFlagBusy {
		panic("illegal state")
	}
	return SlotID(v)
}

func (s *Slot) invalidateSlotId() {
	for {
		v := atomic.LoadUint64(&s.idAndStep)
		if SlotID(v).IsUnknown() || v <= slotFlagBusy {
			panic("illegal state")
		}
		if atomic.CompareAndSwapUint64(&s.idAndStep, v, slotFlagBusy|v&^slotIdMask) {
			return
		}
	}
}

func (s *Slot) isEmpty() bool {
	return atomic.LoadUint64(&s.idAndStep) == 0
}

func (s *Slot) isBusy() bool {
	return atomic.LoadUint64(&s.idAndStep)&slotFlagBusy != 0
}

func (s *Slot) isInitializing() bool {
	v := atomic.LoadUint64(&s.idAndStep)
	return v&^(slotFlagBusy-1) == slotFlagBusy|stepIncrement
}

func (s *Slot) ensureInitializing() {
	if !s.isInitializing() {
		panic("illegal state")
	}
}

func (s *Slot) _slotAllocated(id SlotID) {
	if id == 0 {
		atomic.StoreUint64(&s.idAndStep, slotFlagBusy)
	} else {
		atomic.StoreUint64(&s.idAndStep, uint64(id)|stepIncrement|slotFlagBusy)
	}
}

func (s *Slot) _trySetFlag(f uint64) (bool, uint64) {
	for {
		v := atomic.LoadUint64(&s.idAndStep)
		if v&f != 0 {
			return false, 0
		}

		if atomic.CompareAndSwapUint64(&s.idAndStep, v, v|f) {
			return true, v
		}
	}
}

func (s *Slot) _setFlag(f uint64) uint64 {
	ok, v := s._trySetFlag(f)
	if !ok {
		panic("illegal state")
	}
	return v
}

func (s *Slot) _unsetFlag(f uint64) uint64 {
	for {
		v := atomic.LoadUint64(&s.idAndStep)
		if v&f == 0 {
			panic("illegal state")
		}

		if atomic.CompareAndSwapUint64(&s.idAndStep, v, v&^f) {
			return v
		}
	}
}

func (s *Slot) incStep() {
	for {
		v := atomic.LoadUint64(&s.idAndStep)
		if SlotID(v) == 0 {
			panic("illegal state")
		}
		update := v + stepIncrement
		if update < stepIncrement {
			// overflow, skip steps 0 and 1
			update += stepIncrement * numberOfReservedSteps
		}
		if atomic.CompareAndSwapUint64(&s.idAndStep, v, update) {
			return
		}
	}
}

func (s *Slot) _tryStartSlot(minStepNo uint32) (isEmpty, isStarted bool, prevStepNo uint32) {
	for {
		v := atomic.LoadUint64(&s.idAndStep)
		if v == 0 /* isEmpty() */ {
			return true, false, 0
		}

		prevStepNo = uint32(v >> stepIncrementShift)
		if v&slotFlagBusy != 0 /* isBusy() */ || v < slotFlagBusy /* reserved */ || prevStepNo < minStepNo {
			return false, false, prevStepNo
		}

		if atomic.CompareAndSwapUint64(&s.idAndStep, v, v|slotFlagBusy) {
			return false, true, prevStepNo
		}
	}
}

func (s *Slot) _tryStartWithId(slotId SlotID, minStepNo uint32) (isValid, isStarted bool, prevStepNo uint32) {
	for {
		v := atomic.LoadUint64(&s.idAndStep)
		if v == 0 /* isEmpty() */ || SlotID(v) != slotId {
			return false, false, 0
		}

		prevStepNo = uint32(v >> stepIncrementShift)
		if v&slotFlagBusy != 0 /* isBusy() */ || v < slotFlagBusy /* reserved */ || prevStepNo < minStepNo {
			return false, false, prevStepNo
		}

		if atomic.CompareAndSwapUint64(&s.idAndStep, v, v|slotFlagBusy) {
			return false, true, prevStepNo
		}
	}
}

func (s *Slot) stopWorking() (prevStepNo uint32) {
	for {
		v := atomic.LoadUint64(&s.idAndStep)
		if v&slotFlagBusy == 0 {
			panic("illegal state")
		}

		if atomic.CompareAndSwapUint64(&s.idAndStep, v, v&^slotFlagBusy) {
			return uint32(v >> stepIncrementShift)
		}
	}
}

func (s *Slot) canMigrateWorking(prevStepNo uint32, migrateIsNeeded bool) bool {
	if prevStepNo > 1 {
		return migrateIsNeeded
	}
	return prevStepNo == 1 && atomic.LoadUint64(&s.idAndStep) >= stepIncrement*numberOfReservedSteps
}

func (s *Slot) tryStartMigrate() (isEmpty, isStarted bool, prevStepNo uint32) {
	isEmpty, isStarted, prevStepNo = s._tryStartSlot(2)
	return
}

func (s *Slot) startWorking(scanNo uint32) uint32 {
	if _, isStarted, prevStepNo := s._tryStartSlot(1); isStarted {
		s.lastWorkScan = uint8(scanNo)
		return prevStepNo
	}
	panic("illegal state")
}

func (s *Slot) isInQueue() bool {
	return s.queue != nil || s.nextInQueue != nil || s.prevInQueue != nil
}

func (s *Slot) ensureNotInQueue() {
	if s.isInQueue() {
		panic("illegal state")
	}
}

func (s *Slot) ensureInQueue() {
	if s.queue == nil || s.nextInQueue == nil || s.prevInQueue == nil {
		panic("illegal state")
	}
}

func (s *Slot) dispose() {
	s.ensureNotInQueue()
	if s.slotData.dependency != nil {
		panic("illegal state")
	}
	atomic.StoreUint64(&s.idAndStep, 0)
	s.slotData = slotData{}
}

func (s *Slot) NewLink() SlotLink {
	id, _, _ := s.GetState()
	return SlotLink{id, s}
}

func (s *Slot) NewStepLink() StepLink {
	id, step, _ := s.GetState()
	return StepLink{SlotLink{id, s}, step}
}

func (s *slotData) isLastScan(scanNo uint32) bool {
	return s.lastWorkScan == uint8(scanNo)
}

func (s *Slot) setNextStep(step SlotStep) {
	switch {
	case step.Transition == nil:
		if step.Flags != 0 || step.Migration != nil {
			panic("illegal value")
		}
		// leave as-is
		return

	case step.Flags&StepResetAllFlags == 0:
		step.Flags |= s.defFlags
	default:
		step.Flags &^= StepResetAllFlags
	}
	s.step = step
	s.incStep()
}

func (s *Slot) removeHeadedQueue() *Slot {
	nextDep, _, _ := s.queue.extractAll(nil)
	s.vacateQueueHead()
	return nextDep
}

func (s *Slot) ensureLocal(link SlotLink) {
	if s.machine == nil {
		panic("illegal state")
	}
	if s.machine != link.s.machine {
		panic("illegal state")
	}
}

func (s *Slot) isPriority() bool {
	return s.step.Flags&StepPriority != 0
}

func (s *Slot) getMigration() MigrateFunc {
	if s.step.Migration != nil {
		return s.step.Migration
	}
	return s.defMigrate
}

func (s *Slot) getErrorHandler() ErrorHandlerFunc {
	if s.step.Handler != nil {
		return s.step.Handler
	}
	return s.defErrorHandler
}

func (s *Slot) hasAsyncOrBargeIn() bool {
	return s.asyncCallCount > 0 || s.slotFlags&slotHasBargeIn != 0
}

func (s *Slot) addAsyncCount(asyncCnt uint16) {
	if asyncCnt == 0 {
		return
	}
	s.slotFlags |= slotHadAsync
	asyncCnt += s.asyncCallCount
	if asyncCnt <= s.asyncCallCount {
		panic("overflow")
	}
	s.asyncCallCount = asyncCnt
}

func (s *Slot) decAsyncCount() {
	if s.asyncCallCount == 0 {
		panic("underflow")
	}
	s.asyncCallCount--
}

func (s *Slot) _logStepUpdate(prevStepNo uint32, stateUpdate StateUpdate, flags StepLoggerFlags) {
	stepData := StepLoggerData{StepNo: s.NewStepLink(), Flags: flags, CurrentStep: s.step, NextStep: stateUpdate.step}

	if prevStepNo <= 1 {
		// nil all handlers as initialization transition can't be logged properly
		stepData.CurrentStep = SlotStep{Flags: s.step.Flags | StepResetAllFlags}
	}
	stepData.UpdateType, _ = getStateUpdateTypeName(stateUpdate)
	s.stepLogger(&stepData)
}

func (s *Slot) logStepUpdate(prevStepNo uint32, stateUpdate StateUpdate, wasAsync bool) {
	if s.stepLogger == nil {
		return
	}

	if wasAsync {
		s._logStepUpdate(prevStepNo, stateUpdate, StepLoggerDetached)
	} else {
		s._logStepUpdate(prevStepNo, stateUpdate, 0)
	}
}

func (s *Slot) logStepMigrate(prevStepNo uint32, stateUpdate StateUpdate) {
	if s.stepLogger == nil {
		return
	}
	s._logStepUpdate(prevStepNo, stateUpdate, StepLoggerMigrate)
}
