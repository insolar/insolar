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
	return newStateUpdate(marker, uint16(upd), step, param)
}

func slotMachineUpdateUint(marker ContextMarker, upd stateUpdType, step SlotStep, param uint32) StateUpdate {
	return newStateUpdateUint(marker, uint16(upd), step, param)
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
	return newStateUpdateUint(marker, uint16(stateUpdNextLoop), SlotStep{}, ulimit)
}

func stateUpdateNext(marker ContextMarker, slotStep SlotStep, canLoop bool) StateUpdate {
	slotStep.ensureTransition()
	if canLoop {
		return slotMachineUpdateUint(marker, stateUpdNextLoop, slotStep, math.MaxUint32)
	}

	return slotMachineUpdateUint(marker, stateUpdNext, slotStep, 0)
}

type StepPrepareFunc func()

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

func stateUpdateWaitForSlot(marker ContextMarker, waitOn SlotLink, slotStep SlotStep) StateUpdate {
	return newStateUpdateLink(marker, uint16(stateUpdSleep), waitOn, slotStep, nil)
}

func stateUpdateWaitForShared(marker ContextMarker, waitOn SlotLink, slotStep SlotStep) StateUpdate {
	return newStateUpdateLink(marker, uint16(stateUpdWaitForShared), waitOn, slotStep, nil)
}

func stateUpdateWaitForEvent(marker ContextMarker, slotStep SlotStep, prepare StepPrepareFunc, until uint32) StateUpdate {
	return newStateUpdateUintAndParam(marker, uint16(stateUpdWaitForEvent), slotStep, until, prepareToParam(prepare))
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
	return slotMachineUpdate(marker, stateUpdError, SlotStep{}, err)
}

func stateUpdatePanic(recovered error) StateUpdate {
	return slotMachineUpdate(0, stateUpdPanic, SlotStep{}, recovered)
}

func stateUpdateExpired(slotStep SlotStep, info interface{}) StateUpdate {
	return slotMachineUpdate(0, stateUpdExpired, slotStep, info)
}

type stateUpdType uint32

func (u stateUpdType) HasStep() bool {
	return u >= stateUpdWaitForActive
}

func (u stateUpdType) HasPrepare() bool {
	return u >= stateUpdNext
}

const (
	_ stateUpdType = iota

	// no step
	stateUpdNoChange
	stateUpdStop
	stateUpdError // external handler, cant be detached
	stateUpdPanic // external handler, cant be detached
	stateUpdExpired
	stateUpdReplace // external handler, cant be detached
	stateUpdReplaceWith

	// step, no prepare
	stateUpdWaitForActive
	stateUpdWaitForShared
	stateUpdRepeat   // supports short-loop
	stateUpdNextLoop // supports short-loop

	// step and prepare // prepare should be executed in-sync
	stateUpdNext
	stateUpdPoll
	stateUpdSleep
	stateUpdWaitForEvent
)
