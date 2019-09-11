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
	idAndStep uint64 //atomic access
	parent    SlotLink
	machine   *SlotMachine
	ctx       context.Context

	declaration StateMachineDeclaration
	step        SlotStep

	defMigrate      MigrateFunc
	defErrorHandler ErrorHandlerFunc
	defFlags        StepFlags

	workState    slotWorkFlags
	lastWorkScan uint8

	asyncCallCount uint16 // pending calls
	migrationCount uint16 // can be wrapped by overflow

	dependency SlotDependency

	/* -----------------------------------
	   Slot fields to support processing queues
	   ----------------------------------- */
	prevInQueue *Slot
	nextInQueue *Slot
	queue       *QueueHead
}

type SlotDependency interface {
	GetKey() string
	GetWeight() int32

	OnStepChanged()
	OnSlotDisposed()
	OnBroadcast(payload interface{}) (accepted, wakeup bool)

	Remove()
}

const stepIncrementShift = 32
const stepIncrement = 1 << stepIncrementShift

type slotWorkFlags uint8

const (
	slotWorking slotWorkFlags = 1 << iota
)

func (s *Slot) ensureNotInQueue() {
	if s.queue != nil || s.nextInQueue != nil || s.prevInQueue != nil {
		panic("illegal state")
	}
}

func (s *Slot) ensureInQueue() {
	if s.queue == nil || s.nextInQueue == nil || s.prevInQueue == nil {
		panic("illegal state")
	}
}

func (s *Slot) GetID() SlotID {
	return SlotID(s.idAndStep)
}

func (s *Slot) GetStep() uint32 {
	return uint32(s.idAndStep >> stepIncrementShift)
}

func (s *Slot) GetAtomicIDAndStep() (SlotID, uint32) {
	v := atomic.LoadUint64(&s.idAndStep)
	return SlotID(v), uint32(v >> stepIncrementShift)
}

func (s *Slot) init(ctx context.Context, id SlotID, parent SlotLink, decl StateMachineDeclaration,
	machine *SlotMachine) {

	if decl == nil {
		panic("illegal value")
	}
	if machine == nil {
		panic("illegal value")
	}
	if id.IsUnknown() {
		panic("illegal value")
	}
	switch {
	case s.machine == machine:
		break
	case s.machine == nil:
		s.machine = machine
	default:
		panic("illegal value")
	}

	s.ensureNotInQueue()
	s.parent = parent
	s.declaration = decl
	s.ctx = ctx
	atomic.StoreUint64(&s.idAndStep, uint64(id)+stepIncrement)
}

func (s *Slot) incStep() bool {
	for {
		v := atomic.LoadUint64(&s.idAndStep)
		if v == 0 {
			return false
		}
		update := v + stepIncrement
		if update < stepIncrement {
			// overflow, skip 0 step value
			update += stepIncrement
		}
		if atomic.CompareAndSwapUint64(&s.idAndStep, v, update) {
			return true
		}
	}
}

func (s *Slot) dispose() {
	s.ensureNotInQueue()
	if s.dependency != nil {
		panic("illegal state")
	}

	atomic.StoreUint64(&s.idAndStep, 0)
	// this may cause data racing error ... but there is none
	*s = Slot{machine: s.machine}
}

func (s *Slot) NewLink() SlotLink {
	return SlotLink{s.GetID(), s}
}

func (s *Slot) NewStepLink() StepLink {
	return StepLink{s.NewLink(), s.GetStep()}
}

func (s *Slot) isEmpty() bool {
	return s.declaration == nil
}

func (s *Slot) isWorking() bool {
	return s.workState&slotWorking != 0
}

func (s *Slot) isLastScan(scanNo uint32) bool {
	return s.lastWorkScan == uint8(scanNo)
}

func (s *Slot) startWorking(scanNo uint32) /* , timeMark time.Duration) */ {
	if s.workState&slotWorking != 0 {
		panic("illegal state")
	}
	s.lastWorkScan = uint8(scanNo)
	s.workState |= slotWorking
}

func (s *Slot) stopWorking(asyncCount uint16) {
	if s.workState&slotWorking == 0 {
		panic("illegal state")
	}
	s.asyncCallCount += asyncCount
	s.workState &^= slotWorking
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

func (s *Slot) removeHeadedQueue(moveTo func(slot *Slot)) {
	for {
		next := s.QueueNext()
		if next == nil {
			break
		}
		next.removeFromQueue()
		moveTo(next)
	}
	s.vacateQueueHead()
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
