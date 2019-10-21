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
	"errors"
	"runtime"
)

const (
	_ stateUpdKind = iota

	stateUpdNoChange
	stateUpdStop
	stateUpdError
	stateUpdPanic
	stateUpdReplace
	stateUpdReplaceWith

	stateUpdInternalRepeatNow // this is a special op. MUST NOT be used anywhere else.

	stateUpdRepeat   // supports short-loop
	stateUpdNextLoop // supports short-loop

	stateUpdNext
	stateUpdPoll
	stateUpdSleep
	stateUpdWaitForEvent
	stateUpdWaitForActive
	stateUpdWaitForShared
)

const stateUpdWakeup = stateUpdRepeat

var stateUpdateTypes []StateUpdateType

// init() is used instead of variable initializer to avoid "initialization loop" error
func init() {
	stateUpdateTypes = []StateUpdateType{
		stateUpdNoChange: {
			filter: updCtxMigrate | updCtxBargeIn | updCtxAsyncCallback,

			apply: func(slot *Slot, stateUpdate StateUpdate, worker FixedSlotWorker) (isAvailable bool, err error) {
				//if !slot.isInQueue() {
				//	return false, errors.New("unexpected state update")
				//}
				return true, nil
			},
		},

		stateUpdInternalRepeatNow: {
			filter: updCtxInternal, // can't be created by a template
			apply: func(slot *Slot, stateUpdate StateUpdate, worker FixedSlotWorker) (isAvailable bool, err error) {
				if slot.isInQueue() {
					return false, errors.New("unexpected internal repeat")
				}
				m := slot.machine
				m.workingSlots.AddFirst(slot)
				return true, nil
			},
		},

		stateUpdStop: {
			filter: updCtxExec | updCtxInit | updCtxMigrate | updCtxBargeIn,
			apply:  stateUpdateDefaultStop,
		},

		stateUpdError: {
			filter:    updCtxExec | updCtxInit | updCtxMigrate | updCtxBargeIn,
			params:    updParamVar,
			varVerify: stateUpdateDefaultVerifyError,
			apply:     stateUpdateDefaultError,
		},

		stateUpdPanic: {
			filter:    updCtxInternal, // can't be created by a template
			params:    updParamVar,
			varVerify: stateUpdateDefaultVerifyError,
			apply:     stateUpdateDefaultError,
		},

		stateUpdReplaceWith: {
			filter: updCtxExec | updCtxMigrate,
			params: updParamVar,
			varVerify: func(v interface{}) {
				sm := v.(StateMachine)
				if sm == nil {
					panic("illegal value")
				}
			},

			prepare: func(slot *Slot, stateUpdate *StateUpdate) {
				m := slot.machine

				sm, ok := stateUpdate.param1.(StateMachine)
				if !ok {
					panic("illegal value")
				}

				newSlot := m.allocateSlot()
				newSlot.slotReplaceData = slot.slotReplaceData.takeOutForReplace()
				if m.prepareNewSlot(newSlot, slot, nil, sm, true) {
					// prevent this slot from firing the termination handler
					slot.defResultHandler = nil
					slot.defResult = nil
				}

				stateUpdate.param1 = nil
				stateUpdate.link = newSlot
			},

			apply: stateUpdateDefaultReplace,
		},

		stateUpdReplace: {
			filter: updCtxExec | updCtxMigrate,
			params: updParamVar,

			prepare: func(slot *Slot, stateUpdate *StateUpdate) {
				m := slot.machine

				fn, ok := stateUpdate.param1.(CreateFunc)
				if !ok {
					panic("illegal value")
				}
				newSlot := m.allocateSlot()
				newSlot.slotReplaceData = slot.slotReplaceData.takeOutForReplace()
				if m.prepareNewSlot(newSlot, slot, fn, nil, true) {
					// prevent this slot from firing the termination handler
					slot.defResultHandler = nil
					slot.defResult = nil
				}

				stateUpdate.param1 = nil
				stateUpdate.link = newSlot
			},

			apply: stateUpdateDefaultReplace,
		},

		stateUpdRepeat: {
			filter: updCtxExec | updCtxBargeIn | updCtxAsyncCallback,
			params: updParamUint,

			shortLoop: func(slot *Slot, stateUpdate StateUpdate, loopCount uint32) bool {
				return loopCount < stateUpdate.param0
			},

			apply: func(slot *Slot, stateUpdate StateUpdate, worker FixedSlotWorker) (isAvailable bool, err error) {
				slot.activateSlot(worker)
				return true, nil
			},
		},

		stateUpdNextLoop: {
			filter: updCtxExec | updCtxInit | updCtxBargeIn,
			params: updParamStep | updParamUint,

			shortLoop: func(slot *Slot, stateUpdate StateUpdate, loopCount uint32) bool {
				if loopCount >= stateUpdate.param0 {
					return false
				}
				ns := stateUpdate.step.Transition
				if ns != nil && !slot.declaration.IsConsecutive(slot.step.Transition, ns) {
					return false
				}
				slot.setNextStep(stateUpdate.step)
				return true
			},

			apply: stateUpdateDefaultJump,
		},

		stateUpdNext: {
			filter:    updCtxExec | updCtxInit | updCtxBargeIn | updCtxMigrate,
			params:    updParamStep | updParamVar,
			prepare:   stateUpdateDefaultNoArgPrepare,
			varVerify: stateUpdateDefaultVerifyNoArgFn,
			apply:     stateUpdateDefaultJump,
		},

		stateUpdPoll: {
			filter:  updCtxExec,
			params:  updParamStep | updParamVar,
			prepare: stateUpdateDefaultNoArgPrepare,
			apply: func(slot *Slot, stateUpdate StateUpdate, worker FixedSlotWorker) (isAvailable bool, err error) {
				m := slot.machine
				slot.setNextStep(stateUpdate.step)
				m.updateSlotQueue(slot, worker, deactivateSlot)
				m.pollingSlots.Add(slot)
				return true, nil
			},
		},

		stateUpdSleep: {
			filter:  updCtxExec,
			params:  updParamStep | updParamVar,
			prepare: stateUpdateDefaultNoArgPrepare,
			apply: func(slot *Slot, stateUpdate StateUpdate, worker FixedSlotWorker) (isAvailable bool, err error) {
				m := slot.machine
				slot.setNextStep(stateUpdate.step)
				m.updateSlotQueue(slot, worker, deactivateSlot)
				return true, nil
			},
		},

		stateUpdWaitForEvent: {
			filter:  updCtxExec,
			params:  updParamStep | updParamUint | updParamVar,
			prepare: stateUpdateDefaultNoArgPrepare,
			apply: func(slot *Slot, stateUpdate StateUpdate, worker FixedSlotWorker) (isAvailable bool, err error) {
				m := slot.machine
				slot.setNextStep(stateUpdate.step)
				m.updateSlotQueue(slot, worker, activateHotWaitSlot)

				if stateUpdate.param0 > 0 {
					m.scanWakeUpAt = minTime(m.scanWakeUpAt, m.fromRelativeTime(stateUpdate.param0))
				}
				return true, nil
			},
		},

		stateUpdWaitForActive: {
			filter: updCtxExec,
			params: updParamStep | updParamLink,
			//		prepare: stateUpdateDefaultNoArgPrepare,
			apply: func(slot *Slot, stateUpdate StateUpdate, worker FixedSlotWorker) (isAvailable bool, err error) {
				m := slot.machine
				slot.setNextStep(stateUpdate.step)
				waitOn := stateUpdate.getLink()

				if waitOn.s == slot || !waitOn.IsValid() {
					// don't wait for self
					// don't wait for an expired slot
					m.updateSlotQueue(slot, worker, activateSlot)
					return
				}
				panic("not implemented") // TODO requires sync
				//switch waitOn.s.QueueType() {
				//case ActiveSlots, WorkingSlots:
				//	// don't wait
				//	m.updateSlotQueue(slot, worker, activateSlot)
				//case NoQueue:
				//	waitOn.s.makeQueueHead()
				//	fallthrough
				//case ActivationOfSlot, PollingSlots:
				//	m.updateSlotQueue(slot, worker, deactivateSlot)
				//	waitOn.s.queue.AddLast(slot)
				//default:
				//	return false, errors.New("illegal slot queue")
				//}
				return true, nil
			},
		},

		stateUpdWaitForShared: {
			filter:  updCtxExec,
			params:  updParamStep | updParamLink,
			prepare: stateUpdateDefaultNoArgPrepare,
			apply: func(slot *Slot, stateUpdate StateUpdate, worker FixedSlotWorker) (isAvailable bool, err error) {
				m := slot.machine
				slot.setNextStep(stateUpdate.step)

				waitOn := stateUpdate.getLink()
				if waitOn.s == slot || !waitOn.IsValid() {
					// don't wait for self
					// don't wait for an expired slot
					m.updateSlotQueue(slot, worker, activateSlot)
					return
				}

				wakeupLink := slot.NewLink()
				m.syncQueue.AddAsyncCallback(waitOn, func(waitOn SlotLink, worker DetachableSlotWorker) bool {
					switch {
					case !wakeupLink.IsValid():
						return true
					case waitOn.isValidAndBusy():
						// add this back
						return false
					case !worker.NonDetachableCall(wakeupLink.s.activateSlot):
						m.syncQueue.AddAsyncUpdate(wakeupLink, SlotLink.activateSlot)
					}
					return true
				})

				return true, nil
			},
		},
	}

	for i := range stateUpdateTypes {
		if stateUpdateTypes[i].filter != 0 {
			stateUpdateTypes[i].updKind = stateUpdKind(i)
		}
	}
}

func stateUpdateDefaultNoArgPrepare(_ *Slot, stateUpdate *StateUpdate) {
	fn := stateUpdate.param1.(StepPrepareFunc)
	if fn == nil {
		return
	}
	fn()
}

func stateUpdateDefaultVerifyNoArgFn(u interface{}) {
	runtime.KeepAlive(u.(StepPrepareFunc))
}

func stateUpdateDefaultVerifyError(u interface{}) {
	err := u.(error)
	if err == nil {
		panic("illegal value")
	}
}

func stateUpdateDefaultError(slot *Slot, stateUpdate StateUpdate, w FixedSlotWorker) (isAvailable bool, err error) {
	err = stateUpdate.param1.(error)
	if err == nil {
		err = errors.New("error argument is missing")
	}

	return slot.machine.handleSlotUpdateError(slot, w,
		getStateUpdateKind(stateUpdate) == stateUpdPanic, err), nil
}

func stateUpdateDefaultJump(slot *Slot, stateUpdate StateUpdate, worker FixedSlotWorker) (isAvailable bool, err error) {
	m := slot.machine
	slot.setNextStep(stateUpdate.step)
	m.updateSlotQueue(slot, worker, activateSlot)
	return true, nil
}

func stateUpdateDefaultStop(slot *Slot, _ StateUpdate, worker FixedSlotWorker) (isAvailable bool, err error) {
	// recycleSlot can handle both in-place and off-place updates
	m := slot.machine
	m.recycleSlot(slot, worker)
	return false, nil
}

func stateUpdateDefaultReplace(slot *Slot, stateUpdate StateUpdate, worker FixedSlotWorker) (isAvailable bool, err error) {
	if replacementSlot, ok := stateUpdate.param1.(*Slot); ok {
		m := replacementSlot.machine
		defer m.startNewSlot(replacementSlot, worker)
		return stateUpdateDefaultStop(slot, stateUpdate, worker)
	}

	return false, errors.New("replacement SM is missing")
}
