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
	"math"
)

func slotMachineUpdate(marker *struct{}, upd stateUpdType, step SlotStep, param interface{}) StateUpdate {
	return NewStateUpdate(marker, uint16(upd), step, param)
}

func slotMachineUpdateUint(marker *struct{}, upd stateUpdType, step SlotStep, param uint32) StateUpdate {
	return NewStateUpdateUint(marker, uint16(upd), step, param)
}

func stateUpdateNoChange(marker *struct{}) StateUpdate {
	return slotMachineUpdate(marker, stateUpdNoChange, SlotStep{}, nil)
}

func stateUpdateRepeat(marker *struct{}, limit int) StateUpdate {
	ulimit := uint32(0)

	switch {
	case limit > math.MaxUint32:
		ulimit = math.MaxUint32
	case limit > 0:
		ulimit = uint32(limit)
	}
	return NewStateUpdateUint(marker, uint16(stateUpdRepeat), SlotStep{}, ulimit)
}

func stateUpdateNext(marker *struct{}, sf StateFunc, mf MigrateFunc, canLoop bool, flags StepFlags) StateUpdate {
	if sf == nil {
		panic("illegal value")
	}

	slotStep := SlotStep{Transition: sf, Migration: mf, StepFlags: flags}
	if canLoop {
		return slotMachineUpdateUint(marker, stateUpdNextLoop, slotStep, math.MaxUint32)
	}

	return slotMachineUpdateUint(marker, stateUpdNext, slotStep, 0)
}

type StepPrepareFunc func(slot *Slot)

func prepareToParam(prepare StepPrepareFunc) interface{} {
	if prepare == nil {
		return nil
	}
	return prepare
}

func stateUpdateYield(marker *struct{}, slotStep SlotStep, prepare StepPrepareFunc) StateUpdate {
	return slotMachineUpdate(marker, stateUpdNext, slotStep, prepareToParam(prepare))
}

func stateUpdatePoll(marker *struct{}, slotStep SlotStep, prepare StepPrepareFunc) StateUpdate {
	return slotMachineUpdate(marker, stateUpdPoll, slotStep, prepareToParam(prepare))
}

func stateUpdateWait(marker *struct{}, slotStep SlotStep, prepare StepPrepareFunc) StateUpdate {
	return slotMachineUpdate(marker, stateUpdWait, slotStep, prepareToParam(prepare))
}

func stateUpdateWaitForSlot(marker *struct{}, waitOn SlotLink, slotStep SlotStep, prepare StepPrepareFunc) StateUpdate {
	return NewStateUpdateLink(marker, uint16(stateUpdWait), waitOn, slotStep, prepareToParam(prepare))
}

func stateUpdateReplace(marker *struct{}, cf CreateFunc) StateUpdate {
	if cf == nil {
		panic("illegal state")
	}
	return slotMachineUpdate(marker, stateUpdReplace, SlotStep{}, cf)
}

func stateUpdateStop(marker *struct{}) StateUpdate {
	return slotMachineUpdate(marker, stateUpdStop, SlotStep{}, nil)
}

func stateUpdateFailed(err error) StateUpdate {
	return slotMachineUpdate(nil, stateUpdDispose, SlotStep{}, err)
}

func stateUpdateExpired(slotStep SlotStep, info interface{}) StateUpdate {
	return slotMachineUpdate(nil, stateUpdExpired, slotStep, info)
}

type stateUpdType uint32

func (u stateUpdType) HasStep() bool {
	return u >= stateUpdRepeat
}

func (u stateUpdType) HasPrepare() bool {
	return u > stateUpdNextLoop
}

const (
	_ stateUpdType = iota
	stateUpdNoChange
	stateUpdStop
	stateUpdDispose
	stateUpdExpired
	stateUpdReplace

	stateUpdRepeat   // supports short-loop, no prepare
	stateUpdNextLoop // supports short-loop, no prepare
	stateUpdNext
	stateUpdPoll
	stateUpdWait

	//stateUpdFlagNoWakeup = 1 << 5
	//stateUpdFlagHasAsync = 1 << 6
	//stateUpdFlagYield    = 1 << 7
	//stateUpdateMask     stateUpdateFlags = 0x0F
)
