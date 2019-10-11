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
	"math"
	"unsafe"
)

type contextTemplate struct {
	mode updCtxMode
}

func (p *contextTemplate) getMarker() ContextMarker {
	return ContextMarker(uintptr(unsafe.Pointer(p))) // TODO replace with getting an atomic counter?
}

func (p *contextTemplate) ensureAndPrepare(s *Slot, stateUpdate StateUpdate) StateUpdate {
	stateUpdate.ensureMarker(p.getMarker())

	sut := typeOfStateUpdateForMode(p.mode, stateUpdate)
	sut.Prepare(s, &stateUpdate)

	return stateUpdate
}

func (p *contextTemplate) setMode(mode updCtxMode) {
	if mode == updCtxInactive {
		panic("illegal value")
	}
	if p.mode != updCtxInactive {
		panic("illegal state")
	}
	p.mode = mode
}

func (p *contextTemplate) ensureAtLeast(mode updCtxMode) {
	if p.mode < mode {
		panic("illegal state")
	}
}

func (p *contextTemplate) ensure(mode0 updCtxMode) {
	if p.mode != mode0 {
		panic("illegal state")
	}
}

func (p *contextTemplate) ensureAny2(mode0, mode1 updCtxMode) {
	if p.mode != mode0 && p.mode != mode1 {
		panic("illegal state")
	}
}

func (p *contextTemplate) ensureAny3(mode0, mode1, mode2 updCtxMode) {
	if p.mode != mode0 && p.mode != mode1 && p.mode != mode2 {
		panic("illegal state")
	}
}

func (p *contextTemplate) ensureValid() {
	if p.mode <= updCtxDiscarded {
		panic("illegal state")
	}
}

func (p *contextTemplate) template(updType stateUpdKind) StateUpdateTemplate {
	return newStateUpdateTemplate(p.mode, p.getMarker(), updType)
}

func (p *contextTemplate) setDiscarded() {
	p.mode = updCtxDiscarded
}

func (p *contextTemplate) discardAndCapture(msg string, recovered interface{}, err *error) {
	p.mode = updCtxDiscarded
	if recovered == nil {
		return
	}
	*err = recoverSlotPanic(msg, recovered, *err)
}

func (p *contextTemplate) discardAndUpdate(msg string, recovered interface{}, update *StateUpdate) {
	p.mode = updCtxDiscarded
	recoverSlotPanicAsUpdate(update, msg, recovered, nil)
}

/* ========================================================================= */

type slotContext struct {
	contextTemplate
	s *Slot
	w DetachableSlotWorker
}

func (p *slotContext) clone() slotContext {
	p.ensureValid()
	return slotContext{s: p.s, w: p.w}
}

func (p *slotContext) SlotLink() SlotLink {
	p.ensureValid()
	return p.s.NewLink()
}

func (p *slotContext) StepLink() StepLink {
	p.ensure(updCtxExec)
	return p.s.NewStepLink()
}

func (p *slotContext) GetContext() context.Context {
	p.ensureValid()
	return p.s.ctx
}

func (p *slotContext) ParentLink() SlotLink {
	p.ensureValid()
	return p.s.parent
}

func (p *slotContext) SetDefaultErrorHandler(fn ErrorHandlerFunc) {
	p.ensureAtLeast(updCtxInit)
	p.s.defErrorHandler = fn
}

func (p *slotContext) SetDefaultMigration(fn MigrateFunc) {
	p.ensureAtLeast(updCtxInit)
	p.s.defMigrate = fn
}

func (p *slotContext) SetDefaultFlags(flags StepFlags) {
	p.ensureAtLeast(updCtxInit)
	if flags&StepResetAllFlags != 0 {
		p.s.defFlags = flags &^ StepResetAllFlags
	} else {
		p.s.defFlags |= flags
	}
}

func (p *slotContext) JumpExt(step SlotStep) StateUpdate {
	return p.template(stateUpdNext).newStep(step, nil)
}

func (p *slotContext) Jump(fn StateFunc) StateUpdate {
	return p.template(stateUpdNextLoop).newStepUint(SlotStep{Transition: fn}, math.MaxUint32)
}

func (p *slotContext) Stop() StateUpdate {
	return p.template(stateUpdStop).newNoArg()
}

func (p *slotContext) Error(err error) StateUpdate {
	return p.template(stateUpdError).newError(err)
}

func (p *slotContext) Replace(fn CreateFunc) StateUpdate {
	return p.template(stateUpdReplace).newVar(fn)
}

func (p *slotContext) ReplaceWith(sm StateMachine) StateUpdate {
	return p.template(stateUpdReplaceWith).newVar(sm)
}

func (p *slotContext) Repeat(limit int) StateUpdate {
	ulimit := uint32(0)
	switch {
	case limit <= 0:
	case limit > math.MaxUint32:
		ulimit = math.MaxUint32
	default:
		ulimit = uint32(limit)
	}

	return p.template(stateUpdRepeat).newUint(ulimit)
}

func (p *slotContext) Stay() StateUpdate {
	return p.template(stateUpdNoChange).newNoArg()
}

func (p *slotContext) WakeUp() StateUpdate {
	return p.template(stateUpdWakeup).newUint(0)
}

func (p *slotContext) Share(data interface{}, wakeUpAfterUse bool) SharedDataLink {
	p.ensureAtLeast(updCtxInit)
	if data == nil {
		panic("illegal value")
	}
	return SharedDataLink{p.s.NewLink(), wakeUpAfterUse, data}
}

func (p *slotContext) AffectedStep() SlotStep {
	p.ensureAny3(updCtxMigrate, updCtxBargeIn, updCtxFail)
	r := p.s.step
	r.Flags |= StepResetAllFlags
	return p.s.step
}

func (p *slotContext) NewChild(ctx context.Context, fn CreateFunc) SlotLink {
	return p._newChild(ctx, fn, false)
}

func (p *slotContext) InitChild(ctx context.Context, fn CreateFunc) SlotLink {
	return p._newChild(ctx, fn, true)
}

func (p *slotContext) _newChild(ctx context.Context, fn CreateFunc, runInit bool) SlotLink {
	p.ensureAny2(updCtxExec, updCtxFail)
	if fn == nil {
		panic("illegal value")
	}
	if ctx == nil {
		panic("illegal value")
	}

	m := p.s.machine
	newSlot := m.allocateSlot()
	newSlot.ctx = ctx
	newSlot.parent = p.s.NewLink()
	link := newSlot.NewLink()

	m.prepareNewSlot(newSlot, p.s, fn, nil)
	m.startNewSlotByDetachable(newSlot, runInit, p.w)
	return link
}

func (p *slotContext) BargeInWithParam(applyFn BargeInApplyFunc) BargeInParamFunc {
	p.ensureAny2(updCtxExec, updCtxInit)
	return p.s.machine.createBargeIn(p.s.NewStepLink().AnyStep(), applyFn)
}

func (p *slotContext) BargeIn() BargeInRequester {
	p.ensureAny2(updCtxExec, updCtxInit)
	return &bargeInRequest{&p.contextTemplate, p.s.machine, p.s.NewStepLink().AnyStep()}
}

func (p *slotContext) BargeInThisStepOnly() BargeInRequester {
	p.ensure(updCtxExec)
	return &bargeInRequest{&p.contextTemplate, p.s.machine, p.s.NewStepLink()}
}

func (p *slotContext) Check(link SyncLink) Decision {
	p.ensureAtLeast(updCtxInit)
	dep := p.s.dependency
	if dep != nil {
		d := link.controller.CheckDependency(dep)
		if d.IsValid() {
			return d
		}
	}
	return link.controller.CheckState()
}

func (p *slotContext) AcquireForThisStep(link SyncLink) Decision {
	return p.acquire(link, true)
}

func (p *slotContext) Acquire(link SyncLink) Decision {
	return p.acquire(link, false)
}

func (p *slotContext) acquire(link SyncLink, oneStep bool) Decision {
	p.ensureAtLeast(updCtxInit)
	dep := p.s.dependency
	if dep != nil {
		d := link.controller.UseDependency(dep, oneStep)
		if d.IsValid() {
			return d
		}
	}

	d := Impossible
	p.w.NonDetachableCall(func(worker FixedSlotWorker) {
		p.s.releaseDependency(worker)
		d, p.s.dependency = link.controller.CreateDependency(p.s, oneStep)
	})
	return d
}

func (p *slotContext) Release(link SyncLink) bool {
	p.ensureAtLeast(updCtxInit)
	dep := p.s.dependency
	if dep == nil {
		return false
	}

	if !p.w.NonDetachableCall(p.s.releaseDependency) {
		m := p.s.machine
		m.syncQueue.AddAsyncUpdate(p.s.NewLink(), func(link SlotLink, worker FixedSlotWorker) {
			if !link.IsValid() || link.s.dependency != dep {
				return
			}
			link.s.releaseDependency(worker)
		})
	}
	return true
}
