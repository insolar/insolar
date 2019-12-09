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

	stateUpdWakeup
	stateUpdNext
	stateUpdPoll
	stateUpdSleep
	stateUpdWaitForEvent
	stateUpdWaitForActive
	stateUpdWaitForIdle
)

//const stateUpdWakeup = stateUpdRepeat

var stateUpdateTypes []StateUpdateType

// init() is used instead of variable initializer to avoid "initialization loop" error
func init() {
	stateUpdateTypes = []StateUpdateType{
		stateUpdNoChange: {
			name:   "noChange",
			filter: updCtxMigrate | updCtxBargeIn | updCtxAsyncCallback,

			apply: func(slot *Slot, stateUpdate StateUpdate, worker FixedSlotWorker) (isAvailable bool, err error) {
				//if !slot.isInQueue() {
				//	return false, errors.New("unexpected state update")
				//}
				return true, nil
			},
		},

		stateUpdInternalRepeatNow: {
			name:   "repeatNow",
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
			name:   "stop",
			filter: updCtxExec | updCtxInit | updCtxMigrate | updCtxBargeIn,
			apply:  stateUpdateDefaultStop,
		},

		stateUpdError: {
			name:      "error",
			filter:    updCtxExec | updCtxInit | updCtxMigrate | updCtxBargeIn,
			params:    updParamVar,
			varVerify: stateUpdateDefaultVerifyError,
			apply: func(slot *Slot, stateUpdate StateUpdate, worker FixedSlotWorker) (isAvailable bool, err error) {
				err = stateUpdate.param1.(error)
				if err == nil {
					err = errors.New("error argument is missing")
				}
				return slot.machine.handleSlotUpdateError(slot, worker, stateUpdate, false, err), nil
			},
		},

		stateUpdPanic: {
			name:      "panic",
			filter:    updCtxInternal, // can't be created by a template
			params:    updParamVar,
			varVerify: stateUpdateDefaultVerifyError,
			apply: func(slot *Slot, stateUpdate StateUpdate, worker FixedSlotWorker) (isAvailable bool, err error) {
				err = stateUpdate.param1.(error)
				if err == nil {
					err = errors.New("error argument is missing")
				}
				return slot.machine.handleSlotUpdateError(slot, worker, stateUpdate, true, err), nil
			},
		},

		stateUpdReplaceWith: {
			name:   "replaceWith",
			filter: updCtxExec,
			params: updParamVar,
			varVerify: func(v interface{}) {
				if v.(prepareReplaceData).sm == nil {
					panic("illegal value")
				}
			},

			prepare: statePrepareDefaultReplace,
			apply:   stateUpdateDefaultReplace,
		},

		stateUpdReplace: {
			name:   "replace",
			filter: updCtxExec,
			params: updParamVar,
			varVerify: func(v interface{}) {
				if v.(prepareReplaceData).fn == nil {
					panic("illegal value")
				}
			},

			prepare: statePrepareDefaultReplace,
			apply:   stateUpdateDefaultReplace,
		},

		stateUpdRepeat: {
			name:   "repeat",
			filter: updCtxExec,
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
			name:   "jumpLoop",
			filter: updCtxExec,
			params: updParamStep | updParamUint,

			shortLoop: func(slot *Slot, stateUpdate StateUpdate, loopCount uint32) bool {
				if loopCount >= stateUpdate.param0 {
					return false
				}
				switch nextStep := stateUpdate.step.Transition; {
				case nextStep == nil:
					slot.setNextStep(stateUpdate.step, nil)
					return false // the same step won't be short-looped

				case slot.stepDecl != nil:
					prevSeqId := slot.stepDecl.SeqId
					nextDecl := slot.declaration.GetStepDeclaration(nextStep)
					slot.setNextStep(stateUpdate.step, nextDecl)
					return nextDecl != nil && prevSeqId < nextDecl.SeqId // only proper further steps can be short-looped

				default:
					isConsec, nextDecl := slot.declaration.IsConsecutive(slot.step.Transition, nextStep)
					slot.setNextStep(stateUpdate.step, nextDecl)
					return isConsec
				}
			},

			apply: stateUpdateDefaultJump,
		},

		stateUpdWakeup: {
			name:   "wakeUp",
			filter: updCtxExec | updCtxBargeIn | updCtxAsyncCallback | updCtxMigrate,

			apply: func(slot *Slot, stateUpdate StateUpdate, worker FixedSlotWorker) (isAvailable bool, err error) {
				slot.activateSlot(worker)
				return true, nil
			},
		},

		stateUpdNext: {
			name:      "jump",
			filter:    updCtxExec | updCtxInit | updCtxBargeIn | updCtxMigrate,
			params:    updParamStep | updParamVar,
			prepare:   stateUpdateDefaultNoArgPrepare,
			varVerify: stateUpdateDefaultVerifyNoArgFn,
			apply:     stateUpdateDefaultJump,
		},

		stateUpdPoll: {
			name:    "poll",
			filter:  updCtxExec,
			params:  updParamStep | updParamVar,
			prepare: stateUpdateDefaultNoArgPrepare,
			apply: func(slot *Slot, stateUpdate StateUpdate, worker FixedSlotWorker) (isAvailable bool, err error) {
				m := slot.machine
				slot.setNextStep(stateUpdate.step, nil)
				m.updateSlotQueue(slot, worker, deactivateSlot)
				m.pollingSlots.AddToLatest(slot)
				return true, nil
			},
		},

		stateUpdSleep: {
			name:    "sleep",
			filter:  updCtxExec,
			params:  updParamStep | updParamVar,
			prepare: stateUpdateDefaultNoArgPrepare,
			apply: func(slot *Slot, stateUpdate StateUpdate, worker FixedSlotWorker) (isAvailable bool, err error) {
				m := slot.machine
				slot.setNextStep(stateUpdate.step, nil)
				m.updateSlotQueue(slot, worker, deactivateSlot)
				return true, nil
			},
		},

		stateUpdWaitForEvent: {
			name:    "waitEvent",
			filter:  updCtxExec,
			params:  updParamStep | updParamUint | updParamVar,
			prepare: stateUpdateDefaultNoArgPrepare,
			apply: func(slot *Slot, stateUpdate StateUpdate, worker FixedSlotWorker) (isAvailable bool, err error) {
				m := slot.machine
				slot.setNextStep(stateUpdate.step, nil)

				if stateUpdate.param0 == 0 {
					m.updateSlotQueue(slot, worker, activateHotWaitSlot)
					return true, nil
				}

				waitUntil := m.fromRelativeTime(stateUpdate.param0)

				if !m.scanStartedAt.Before(waitUntil) { // not(started<Until) ==> started>=Until
					slot.activateSlot(worker)
					return true, nil
				}

				m.updateSlotQueue(slot, worker, deactivateSlot)
				if !m.pollingSlots.AddToLatestBefore(waitUntil, slot) {
					m.scanWakeUpAt = minTime(m.scanWakeUpAt, waitUntil)
					m.updateSlotQueue(slot, worker, activateHotWaitSlot)
				}

				return true, nil
			},
		},

		stateUpdWaitForActive: {
			name:   "waitActive",
			filter: updCtxExec,
			params: updParamStep | updParamLink,
			//		prepare: stateUpdateDefaultNoArgPrepare,
			apply: func(slot *Slot, stateUpdate StateUpdate, worker FixedSlotWorker) (isAvailable bool, err error) {
				m := slot.machine
				slot.setNextStep(stateUpdate.step, nil)
				waitOn := stateUpdate.getLink()

				if waitOn.s == slot {
					// don't wait for self
					slot.activateSlot(worker)
					return true, nil
				}

				m.updateSlotQueue(slot, worker, deactivateSlot)

				// TODO work in progress
				panic("work in progress")

				//switch isValid, isBusy := waitOn.getIsValidAndBusy(); {
				//case !isValid:
				//	// don't wait for an expired or busy slot
				//	m.updateSlotQueue(slot, worker, activateSlot)
				//	return true, nil
				//case waitOn.isMachine(m):
				//	if isBusy {
				//		m.updateSlotQueue(slot, worker, activateHotWaitSlot)
				//		return true, nil
				//	}
				//	//m.queueOnSlot(waitOn.s, slot)
				//	panic("not implemented")
				////case worker.OuterCall(waitOn.s.machine, func(worker FixedSlotWorker) {
				////
				////}):
				//default:
				//	panic("not implemented") // TODO decide on action
				//}

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
				//return true, nil
			},
		},

		stateUpdWaitForIdle: {
			name:    "waitIdle",
			filter:  updCtxExec,
			params:  updParamStep | updParamLink,
			prepare: stateUpdateDefaultNoArgPrepare,
			apply: func(slot *Slot, stateUpdate StateUpdate, worker FixedSlotWorker) (isAvailable bool, err error) {
				m := slot.machine
				slot.setNextStep(stateUpdate.step, nil)

				waitOn := stateUpdate.getLink()
				mWaitOn := waitOn.getActiveMachine()
				if waitOn.s == slot || !waitOn.IsValid() || mWaitOn == nil {
					// don't wait for self
					// don't wait for an expired slot
					slot.activateSlot(worker)
					return true, nil
				}

				m.updateSlotQueue(slot, worker, deactivateSlot)

				wakeupLink := slot.NewLink()
				// here is a trick - we put a callback on the AWAITED object
				// because a callback is executed on non-busy object only
				// hence our call back will only be triggered when the object became available
				if !mWaitOn.syncQueue.AddAsyncCallback(waitOn, wakeupLink.activateOnNonBusy) {
					// callback was declined - wake up the requested back
					slot.activateSlot(worker)
				}
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
	if fn != nil {
		fn()
	}
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

func stateUpdateDefaultJump(slot *Slot, stateUpdate StateUpdate, worker FixedSlotWorker) (isAvailable bool, err error) {
	m := slot.machine
	slot.setNextStep(stateUpdate.step, nil)
	m.updateSlotQueue(slot, worker, activateSlot)
	return true, nil
}

func stateUpdateDefaultStop(slot *Slot, _ StateUpdate, worker FixedSlotWorker) (isAvailable bool, err error) {
	// recycleSlot can handle both in-place and off-place updates
	m := slot.machine
	m.recycleSlot(slot, worker)
	return false, nil
}

type prepareReplaceData struct {
	fn  CreateFunc
	sm  StateMachine
	def prepareSlotValue
}

func statePrepareDefaultReplace(creator *Slot, stateUpdate *StateUpdate) {
	m := creator.machine
	rep := stateUpdate.param1.(prepareReplaceData)

	if link, ok := m.prepareNewSlot(creator, rep.fn, rep.sm, rep.def); ok {
		// handover the termination handler
		link.s.defResult = creator.defResult
		link.s.defTerminate = creator.defTerminate
		creator.defTerminate = nil
		creator.defResult = nil

		stateUpdate.param1 = nil
		stateUpdate.link = link.s
	} else {
		panic(errors.New("replacement SM was not created"))
	}
}

func stateUpdateDefaultReplace(slot *Slot, stateUpdate StateUpdate, worker FixedSlotWorker) (isAvailable bool, err error) {
	replacementSlot := stateUpdate.link
	if replacementSlot == nil {
		return false, errors.New("replacement SM is missing")
	}
	m := replacementSlot.machine
	if slot.machine != m {
		return false, errors.New("replacement SM belongs to a different SlotMachine")
	}

	defer m.startNewSlot(replacementSlot, worker)
	return stateUpdateDefaultStop(slot, stateUpdate, worker)
}
