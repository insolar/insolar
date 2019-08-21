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
	"sync/atomic"
)

type Slot struct {
	idAndStep uint64 //atomic access
	parent    SlotLink

	machine StateMachineDeclaration

	nextState   SlotStep
	migrateSlot MigrateFunc

	lastActivation int64  // unix time
	asyncCallCount uint16 // pending calls
	migrationCount uint16
	workState      slotWorkState

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

type slotWorkState uint8

const (
	NotWorking slotWorkState = iota
	Working
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

func (s *Slot) init(id SlotID, parent SlotLink, machine StateMachineDeclaration) {
	if machine == nil {
		panic("illegal state")
	}
	if id.IsUnknown() {
		panic("illegal value")
	}
	s.ensureNotInQueue()
	atomic.StoreUint64(&s.idAndStep, uint64(id)+stepIncrement)
	s.parent = parent
	s.machine = machine
}

func (s *Slot) incStep(slotID SlotID) bool {
	for {
		v := atomic.LoadUint64(&s.idAndStep)
		if v == 0 || SlotID(v) != slotID {
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

func (s *Slot) NewStepLink() StepLink {
	return StepLink{s.NewLink(), s.GetStep(), nil}
}

func (s *Slot) isEmpty() bool {
	return s.machine == nil
}

func (s *Slot) isWorking() bool {
	return s.workState == Working
}

func (s *Slot) setWorking() {
	if s.workState != NotWorking {
		panic("illegal state")
	}
	s.workState = Working
}

func (s *Slot) setNotWorking() {
	s.workState = NotWorking
}
