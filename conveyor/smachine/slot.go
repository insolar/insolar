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
	idAndStep    uint64 //atomic access
	parent       SlotLink
	machineState SlotMachineState
	ctx          context.Context

	declaration StateMachineDeclaration
	step        SlotStep

	defMigrate   MigrateFunc
	defFlags     StepFlags
	workState    slotWorkFlags
	lastWorkScan uint8

	asyncCallCount uint16 // pending calls
	migrationCount uint16 // can be wrapped by overflow

	dependency SlotDependency

	/* -----------------------------------
	   Slot fields to support processing queue
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

type slotWorkFlags uint8

const (
	Working slotWorkFlags = 1 << iota
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
	return uint32(s.idAndStep >> 32)
}

func (s *Slot) GetAtomicIDAndStep() (SlotID, uint32) {
	v := atomic.LoadUint64(&s.idAndStep)
	return SlotID(v), uint32(v >> 32)
}

const stepIncrement = 1 << 32

func (s *Slot) init(ctx context.Context, id SlotID, parent SlotLink, decl StateMachineDeclaration,
	machineState SlotMachineState) {

	if decl == nil {
		panic("illegal state")
	}
	if id.IsUnknown() {
		panic("illegal value")
	}
	s.ensureNotInQueue()
	atomic.StoreUint64(&s.idAndStep, uint64(id)+stepIncrement)
	s.parent = parent
	s.declaration = decl
	s.machineState = machineState
	s.ctx = ctx
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
	s.forcedDispose()
}

func (s *Slot) forcedDispose() {
	atomic.StoreUint64(&s.idAndStep, 0)
	*s = Slot{}
}

func (s *Slot) NewLink() SlotLink {
	return SlotLink{s.GetID(), s}
}

func (s *Slot) NewExactStepLink() StepLink {
	return StepLink{s.NewLink(), s.GetStep()}
}

func (s *Slot) NewAnyStepLink() StepLink {
	return StepLink{s.NewLink(), 0}
}

func (s *Slot) isEmpty() bool {
	return s.declaration == nil
}

func (s *Slot) isWorking() bool {
	return s.workState&Working != 0
}

func (s *Slot) isLastScan(scanNo uint32) bool {
	return s.lastWorkScan == uint8(scanNo)
}

func (s *Slot) startWorking(scanNo uint32) /* , timeMark time.Duration) */ {
	if s.workState&Working != 0 {
		panic("illegal state")
	}
	s.lastWorkScan = uint8(scanNo)
	s.workState |= Working
}

func (s *Slot) stopWorking(asyncCount uint16) {
	if s.workState&Working == 0 {
		panic("illegal state")
	}
	s.asyncCallCount += asyncCount
	s.workState &^= Working
}

func (s *Slot) setNextStep(step SlotStep) {
	if step.Transition == nil {
		panic("illegal value")
	}
	if step.StepFlags&StepResetAllFlags == 0 {
		step.StepFlags |= s.defFlags
	}
	s.step = step
}
