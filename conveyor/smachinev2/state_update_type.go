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
	"github.com/pkg/errors"
)

func newStateUpdateTemplate(contextType updCtxMode, marker ContextMarker, updKind stateUpdKind) StateUpdateTemplate {
	return stateUpdateTypes[updKind].template(marker, contextType)
}

func getStateUpdateType(updType stateUpdKind) StateUpdateType {
	return stateUpdateTypes[updType].get()
}

func typeOfStateUpdate(stateUpdate StateUpdate) StateUpdateType {
	return stateUpdateTypes[stateUpdate.updKind].get()
}

func typeOfStateUpdateForMode(contextMode updCtxMode, stateUpdate StateUpdate) StateUpdateType {
	return stateUpdateTypes[stateUpdate.updKind].getForMode(contextMode)
}

func newPanicStateUpdate(err error) StateUpdate {
	return StateUpdateTemplate{t: &stateUpdateTypes[stateUpdPanic]}.newError(err)
}

func getStateUpdateKind(stateUpdate StateUpdate) stateUpdKind {
	return stateUpdKind(stateUpdate.updKind)
}

type SlotUpdateFunc func(slot *Slot, stateUpdate StateUpdate, worker FixedSlotWorker) (isAvailable bool, err error)
type SlotUpdatePrepareFunc func(slot *Slot, stateUpdate *StateUpdate)
type SlotUpdateShortLoopFunc func(slot *Slot, stateUpdate StateUpdate, loopCount uint32) bool

type StateUpdateType struct {
	updKind stateUpdKind

	/* Runs within a valid ExecutionContext / detachable */
	shortLoop SlotUpdateShortLoopFunc
	/* Runs within a valid ExecutionContext / detachable */
	prepare SlotUpdatePrepareFunc

	/* Runs inside the Machine / non-detachable */
	apply SlotUpdateFunc

	filter    updCtxMode
	params    stateUpdParam
	varVerify func(interface{})
	//bargeIn func()
	//migrate bool
}

type StateUpdateTemplate struct {
	t *StateUpdateType

	marker  ContextMarker
	ctxType updCtxMode
}

type stateUpdBaseType = uint8
type stateUpdKind stateUpdBaseType
type stateUpdParam uint8
type updCtxMode uint32

const updCtxInactive updCtxMode = 0

const (
	updCtxDiscarded updCtxMode = 1 << iota
	updCtxConstruction
	updCtxBargeIn
	updCtxAsyncCallback
	updCtxFail
	updCtxInit
	updCtxExec
	updCtxMigrate
)

const (
	updParamStep stateUpdParam = 1 << iota
	updParamUint
	updParamLink
	updParamVar
)

func (v *StateUpdateType) verify(ctxType updCtxMode) {
	if ctxType <= updCtxDiscarded {
		panic("illegal value")
	}
	if v.updKind == 0 {
		panic("unknown type")
	}
	if v.apply == nil {
		panic("not implemented")
	}
	if ctxType&v.filter != ctxType {
		panic("not allowed")
	}
}

func (v *StateUpdateType) template(marker ContextMarker, ctxType updCtxMode) StateUpdateTemplate {
	v.verify(ctxType)
	return StateUpdateTemplate{v, marker, ctxType}
}

func (v *StateUpdateType) getForMode(ctxType updCtxMode) StateUpdateType {
	v.verify(ctxType)
	return *v
}

func (v *StateUpdateType) get() StateUpdateType {
	if v.updKind == 0 {
		panic("unknown type")
	}
	if v.apply == nil {
		panic("not implemented")
	}
	return *v
}

func (v *StateUpdateType) verifyVar(u interface{}) interface{} {
	if v.varVerify != nil {
		v.varVerify(u)
	}
	return u
}

func (v *StateUpdateType) ShortLoop(slot *Slot, stateUpdate StateUpdate, loopCount uint32) bool {
	if v.shortLoop == nil {
		return false
	}
	return v.shortLoop(slot, stateUpdate, loopCount)
}

func (v *StateUpdateType) Prepare(slot *Slot, stateUpdate *StateUpdate) {
	if v.prepare != nil {
		v.prepare(slot, stateUpdate)
	}
}

func (v *StateUpdateType) Apply(slot *Slot, stateUpdate StateUpdate, worker FixedSlotWorker) (isAvailable bool, err error) {
	if v.apply == nil {
		return false, errors.New("not implemented")
	}
	return v.apply(slot, stateUpdate, worker)
}

func (v StateUpdateTemplate) ensureTemplate(params stateUpdParam) {
	if v.t == nil {
		panic("illegal state")
	}
	if v.t.params&params != params {
		panic("illegal value")
	}
}

func (v StateUpdateTemplate) newNoArg() StateUpdate {
	v.ensureTemplate(0)
	return StateUpdate{
		marker:  v.marker,
		updKind: stateUpdBaseType(v.t.updKind),
	}
}

type StepPrepareFunc func()

func (v StateUpdateTemplate) newStep(slotStep SlotStep, prepare StepPrepareFunc) StateUpdate {
	v.ensureTemplate(updParamStep | updParamVar)
	return StateUpdate{
		marker:  v.marker,
		updKind: stateUpdBaseType(v.t.updKind),
		step:    slotStep,
		param1:  prepare,
	}
}

func (v StateUpdateTemplate) newStepUntil(slotStep SlotStep, prepare func(), until uint32) StateUpdate {
	v.ensureTemplate(updParamStep | updParamUint | updParamVar)
	return StateUpdate{
		marker:  v.marker,
		updKind: stateUpdBaseType(v.t.updKind),
		step:    slotStep,
		param1:  prepare,
		param0:  until,
	}
}

func (v StateUpdateTemplate) newStepUint(slotStep SlotStep, param uint32) StateUpdate {
	v.ensureTemplate(updParamStep | updParamUint)
	return StateUpdate{
		marker:  v.marker,
		updKind: stateUpdBaseType(v.t.updKind),
		param0:  param,
		step:    slotStep,
	}
}

func (v StateUpdateTemplate) newStepLink(slotStep SlotStep, link SlotLink) StateUpdate {
	v.ensureTemplate(updParamStep | updParamLink)
	return StateUpdate{
		marker:  v.marker,
		updKind: stateUpdBaseType(v.t.updKind),
		link:    link.s,
		param0:  uint32(link.id),
		step:    slotStep,
	}
}

func (v StateUpdateTemplate) newVar(u interface{}) StateUpdate {
	v.ensureTemplate(updParamVar)
	return StateUpdate{
		marker:  v.marker,
		updKind: stateUpdBaseType(v.t.updKind),
		param1:  v.t.verifyVar(u),
	}
}

func (v StateUpdateTemplate) newError(e error) StateUpdate {
	v.ensureTemplate(updParamVar)
	return StateUpdate{
		marker:  v.marker,
		updKind: stateUpdBaseType(v.t.updKind),
		param1:  v.t.verifyVar(e),
	}
}

func (v StateUpdateTemplate) newUint(param uint32) StateUpdate {
	v.ensureTemplate(updParamUint)
	return StateUpdate{
		marker:  v.marker,
		updKind: stateUpdBaseType(v.t.updKind),
		param0:  param,
	}
}
