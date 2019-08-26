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
	"math"
	"time"
)

type stateUpdateFn func(m *SlotMachine, slot *Slot, upd StateUpdate) bool

func slotMachineUpdate(marker *struct{}, upd stateUpdType, param interface{}, apply stateUpdateFn) StateUpdate {
	return StateUpdate{
		marker:  marker,
		updType: uint32(upd),
		param:   param,
		apply:   apply,
	}
}

func (u StateUpdate) getStateUpdateFn() stateUpdateFn {
	return u.apply.(stateUpdateFn)
}

func stateUpdateNoChange(marker *struct{}) StateUpdate {
	return slotMachineUpdate(marker, stateUpdNoChange, nil, nil)
}

func stateUpdateRepeat(marker *struct{}, limit int) StateUpdate {
	ulimit := uint32(0)

	switch {
	case limit > math.MaxUint32:
		ulimit = math.MaxUint32
	case limit > 0:
		ulimit = uint32(limit)
	}
	return slotMachineUpdate(marker, stateUpdRepeat, ulimit, nil)
}

func (u StateUpdate) getRepeatLimit() uint32 {
	return u.param.(uint32)
}

func toUnixNano(t time.Time) int64 {
	if t.IsZero() {
		return 0
	}
	r := t.UnixNano()
	if r <= 0 {
		return 1
	}
	return r
}

func stateUpdateNextOnly(marker *struct{}, sf StateFunc, mf MigrateFunc) StateUpdate {
	return _stateUpdateNext(marker, sf, mf, true, 0)
}

func stateUpdateActivate(marker *struct{}, sf StateFunc, mf MigrateFunc, flags stepFlags) StateUpdate {
	return _stateUpdateNext(marker, sf, mf, false, flags)
}

func _stateUpdateNext(marker *struct{}, sf StateFunc, mf MigrateFunc, canRepeat bool, flags stepFlags) StateUpdate {
	if sf == nil {
		panic("illegal value")
	}

	var param interface{}
	slotStep := SlotStep{transition: sf, migration: mf, stepFlags: flags}
	if canRepeat {
		// enables a shortcut on executionContext
		// NB! apply will NOT be called when a shortcut is possible
		param = &slotStep
	}
	return slotMachineUpdate(marker, stateUpdNext, param, func(m *SlotMachine, slot *Slot, upd StateUpdate) bool {
		slot.setNextStep(slotStep)
		m.addSlotToActiveOrWorkingQueue(slot)
		return true
	})
}

func (u StateUpdate) getShortLoopStep() *SlotStep {
	if s, ok := u.param.(*SlotStep); ok {
		return s
	}
	return nil
}

func stateUpdateDeactivate(marker *struct{}, slotStep SlotStep) StateUpdate {
	if slotStep.IsEmpty() {
		panic("illegal value")
	}
	return slotMachineUpdate(marker, stateUpdNext, nil, func(m *SlotMachine, slot *Slot, upd StateUpdate) bool {
		slot.setNextStep(slotStep)
		m.timeReqSlots.AddLast(slot)
		return true
	})
}

func stateUpdateWaitForSlot(marker *struct{}, waitOn SlotLink, slotStep SlotStep) StateUpdate {
	if slotStep.IsEmpty() {
		panic("illegal value")
	}
	if slotStep.HasTimeout() {
		panic("illegal value - slot wait can't be combined with time wait")
	}

	return slotMachineUpdate(marker, stateUpdNext, nil, func(m *SlotMachine, slot *Slot, upd StateUpdate) bool {
		switch {
		case waitOn.s == slot:
			// don't wait
		case !waitOn.IsValid():
			// don't wait
		default:
			switch waitOn.s.QueueType() {
			case ActiveSlots, WorkingSlots:
				// don't wait
			case NoQueue:
				waitOn.s.makeQueueHead()
				fallthrough
			case AnotherSlotQueue, PollingSlots:
				slot.setNextStep(slotStep)
				waitOn.s.queue.AddLast(slot)
				return true
			default:
				panic("illegal state")
			}
		}
		slot.setNextStep(slotStep)
		m.addSlotToActiveOrWorkingQueue(slot)
		return true
	})
}

func stateUpdateReplace(marker *struct{}, cf CreateFunc) StateUpdate {
	if cf == nil {
		panic("illegal state")
	}
	return slotMachineUpdate(marker, stateUpdNext, nil, func(m *SlotMachine, slot *Slot, upd StateUpdate) bool {
		parent := slot.parent
		m.disposeSlot(slot)
		ok, _ := m.applySlotCreate(slot, parent, cf) // recursive call inside
		return ok
	})
}

func stateUpdateStop(marker *struct{}) StateUpdate {
	return slotMachineUpdate(marker, stateUpdNext, nil, func(m *SlotMachine, slot *Slot, upd StateUpdate) bool {
		m.disposeSlot(slot)
		m.unusedSlots.AddLast(slot)
		return false
	})
}

func stateUpdateFailed(err error) StateUpdate {
	return slotMachineUpdate(nil, stateUpdDispose, err, nil)
}

func stateUpdateExpired(info interface{}) StateUpdate {
	return slotMachineUpdate(nil, stateUpdExpired, info, nil)
}

type stateUpdType uint32

const (
	_ stateUpdType = iota
	stateUpdNoChange
	stateUpdRepeat // supports short-loop
	stateUpdNext   // supports short-loop
	stateUpdDispose
	stateUpdExpired

	//stateUpdFlagNoWakeup = 1 << 5
	//stateUpdFlagHasAsync = 1 << 6
	//stateUpdFlagYield    = 1 << 7
	//stateUpdateMask     stateUpdateFlags = 0x0F
)
