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

func slotMachineUpdate(marker ContextMarker, upd stateUpdType, step SlotStep, param interface{}) StateUpdate {
	return NewStateUpdate(marker, uint16(upd), step, param)
}

func slotMachineUpdateUint(marker ContextMarker, upd stateUpdType, step SlotStep, param uint32) StateUpdate {
	return NewStateUpdateUint(marker, uint16(upd), step, param)
}

func stateUpdateNoChange(marker ContextMarker) StateUpdate {
	return slotMachineUpdate(marker, stateUpdNoChange, SlotStep{}, nil)
}

func stateUpdateRepeat(marker ContextMarker, limit int) StateUpdate {
	ulimit := uint32(0)

	switch {
	case limit > math.MaxUint32:
		ulimit = math.MaxUint32
	case limit > 0:
		ulimit = uint32(limit)
	}
	return NewStateUpdateUint(marker, uint16(stateUpdRepeat), SlotStep{}, ulimit)
}

func stateUpdateNext(marker ContextMarker, slotStep SlotStep, canLoop bool) StateUpdate {
	slotStep.ensureTransition()
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

func stateUpdateYield(marker ContextMarker, slotStep SlotStep, prepare StepPrepareFunc) StateUpdate {
	return slotMachineUpdate(marker, stateUpdNext, slotStep, prepareToParam(prepare))
}

func stateUpdatePoll(marker ContextMarker, slotStep SlotStep, prepare StepPrepareFunc) StateUpdate {
	return slotMachineUpdate(marker, stateUpdPoll, slotStep, prepareToParam(prepare))
}

func stateUpdateSleep(marker ContextMarker, slotStep SlotStep, prepare StepPrepareFunc) StateUpdate {
	return slotMachineUpdate(marker, stateUpdSleep, slotStep, prepareToParam(prepare))
}

func stateUpdateWaitForSlot(marker ContextMarker, waitOn SlotLink, slotStep SlotStep, prepare StepPrepareFunc) StateUpdate {
	return NewStateUpdateLink(marker, uint16(stateUpdSleep), waitOn, slotStep, prepareToParam(prepare))
}

func stateUpdateWaitForShared(marker ContextMarker, waitOn SlotLink, slotStep SlotStep) StateUpdate {
	return NewStateUpdateLink(marker, uint16(stateUpdWaitForShared), waitOn, slotStep, nil)
}

func stateUpdateWaitForEvent(marker ContextMarker, slotStep SlotStep) StateUpdate {
	return NewStateUpdate(marker, uint16(stateUpdWaitForEvent), slotStep, nil)
}

func stateUpdateReplace(marker ContextMarker, cf CreateFunc) StateUpdate {
	if cf == nil {
		panic("illegal state")
	}
	return slotMachineUpdate(marker, stateUpdReplace, SlotStep{}, cf)
}

func stateUpdateReplaceWith(marker ContextMarker, sm StateMachine) StateUpdate {
	if sm == nil {
		panic("illegal state")
	}
	return slotMachineUpdate(marker, stateUpdReplaceWith, SlotStep{}, sm)
}

func stateUpdateStop(marker ContextMarker) StateUpdate {
	return slotMachineUpdate(marker, stateUpdStop, SlotStep{}, nil)
}

func stateUpdateError(marker ContextMarker, err error) StateUpdate {
	return slotMachineUpdate(marker, stateUpdFailed, SlotStep{}, err)
}

func stateUpdateInternalError(err error) StateUpdate {
	return slotMachineUpdate(0, stateUpdFailed, SlotStep{}, err)
}

func stateUpdateExpired(slotStep SlotStep, info interface{}) StateUpdate {
	return slotMachineUpdate(0, stateUpdExpired, slotStep, info)
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
	stateUpdFailed
	stateUpdExpired
	stateUpdReplace
	stateUpdReplaceWith

	stateUpdRepeat   // supports short-loop, no prepare
	stateUpdNextLoop // supports short-loop, no prepare
	stateUpdNext
	stateUpdPoll
	stateUpdSleep
	stateUpdWaitForShared
	stateUpdWaitForEvent

	//stateUpdFlagNoWakeup = 1 << 5
	//stateUpdFlagHasAsync = 1 << 6
	//stateUpdFlagYield    = 1 << 7
	//stateUpdateMask     stateUpdateFlags = 0x0F
)
