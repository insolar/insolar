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
	"time"
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

func (p *contextTemplate) template(updType stateUpdType) StateUpdateTemplate {
	return newStateUpdateTemplate(p.mode, p.getMarker(), updType)
}

func (p *contextTemplate) setDiscarded() {
	p.mode = updCtxDiscarded
}

/* ========================================================================= */

type slotContext struct {
	contextTemplate
	s *Slot
	w SlotWorker
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
	return p.template(stateUpdRepeat).newUint(0)
}

func (p *slotContext) Share(data interface{}, wakeUpAfterUse bool) SharedDataLink {
	p.ensureAtLeast(updCtxInit)
	return SharedDataLink{p.s.NewStepLink(), wakeUpAfterUse, data}
}

func (p *slotContext) AffectedStep() SlotStep {
	p.ensureAny3(updCtxMigrate, updCtxBargeIn, updCtxFail)
	r := p.s.step
	r.Flags |= StepResetAllFlags
	return p.s.step
}

func (p *slotContext) NewChild(ctx context.Context, fn CreateFunc) SlotLink {
	p.ensureAny2(updCtxExec, updCtxFail)
	if fn == nil {
		panic("illegal value")
	}
	if ctx == nil {
		panic("illegal value")
	}

	m := p.s.machine
	newSlot := m.allocateSlot()
	newSlot.ctx = p.s.ctx
	newSlot.parent = p.s.NewLink()
	link := newSlot.NewLink()

	m.prepareNewSlot(newSlot, p.s, fn, nil)

	p.w.NonDetachableOrAsyncCall(newSlot, func(s *Slot, w WorkerContext) {
		s.machine.startNewSlot(s, w)
	})

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

/* ========================================================================= */

var _ ExecutionContext = &executionContext{}

type executionContext struct {
	slotContext
	worker          WorkerContext
	countAsyncCalls uint16
}

func (p *executionContext) GetPendingCallCount() int {
	p.ensure(updCtxExec)
	return int(p.s.asyncCallCount) + int(p.countAsyncCalls)
}

//func (p *executionContext) Yield() StateConditionalUpdate {
//	p.ensureExactState(execContext)
//	return &conditionalUpdate{marker: p.getMarker(), updMode: stateUpdNext}
//}
//
//func (p *executionContext) Poll() StateConditionalUpdate {
//	p.ensureExactState(execContext)
//	return &conditionalUpdate{marker: p.getMarker(), updMode: stateUpdPoll}
//}
//
//func (p *executionContext) Sleep() StateConditionalUpdate {
//	p.ensureExactState(execContext)
//	return &conditionalUpdate{marker: p.getMarker(), updMode: stateUpdSleep}
//}
//
//func (p *executionContext) WaitAny() StateConditionalUpdate {
//	p.ensureExactState(execContext)
//	return &conditionalUpdate{marker: p.getMarker(), updMode: stateUpdWaitForEvent}
//}
//
//func (p *executionContext) WaitAnyUntil(until time.Time) StateConditionalUpdate {
//	p.ensureExactState(execContext)
//	u := p.s.machine.toRelativeTime(until)
//	if u != 0 && !until.After(time.Now()) {
//		return &conditionalUpdate{marker: p.getMarker()}
//	}
//	return &conditionalUpdate{marker: p.getMarker(), updMode: stateUpdWaitForEvent, until: u}
//}
//
//func (p *executionContext) waitFor(link SlotLink, updMode stateUpdType) StateConditionalUpdate {
//	p.ensureExactState(execContext)
//	if link.IsEmpty() {
//		panic("illegal value")
//		//		return &conditionalUpdate{marker: p.getMarker()}
//	}
//
//	r, releaseFn := p.worker.AttachTo(p.s, link, false)
//	if releaseFn != nil {
//		defer releaseFn()
//	}
//	switch r {
//	case SharedSlotAbsent, SharedSlotAvailableAlways:
//		// no wait
//		return &conditionalUpdate{marker: p.getMarker()}
//
//	case SharedSlotRemoteBusy, SharedSlotRemoteAvailable:
//		// state has to be re-detected upon applying the update
//		panic("not implemented")
//	}
//	return &conditionalUpdate{marker: p.getMarker(), updMode: updMode, dependency: link}
//}
//
//func (p *executionContext) WaitActivation(link SlotLink) StateConditionalUpdate {
//	return p.waitFor(link, stateUpdWaitForActive)
//}
//
//func (p *executionContext) WaitShared(link SharedDataLink) StateConditionalUpdate {
//	return p.waitFor(link.link.SlotLink, stateUpdWaitForShared)
//}
//
//func (p *executionContext) UseShared(a SharedDataAccessor) SharedAccessReport {
//	p.ensure(updCtxExec)
//	r, releaseFn := p.worker.AttachTo(p.s, a.link.link.SlotLink, a.link.wakeup)
//	if releaseFn != nil {
//		defer releaseFn()
//	}
//	if r.IsAvailable() {
//		a.accessFn(a.link.data)
//	}
//	return r
//}

//func (p *executionContext) SyncOneStep(key string, weight int32, broadcastFn BroadcastReceiveFunc) Syncronizer {
//	p.ensureExactState(execContext)
//	return p.machine.stepSync.Join(p.s, key, weight, broadcastFn)
//}

func (p *executionContext) executeNextStep() (bool, StateUpdate, uint16) {
	p.setMode(updCtxExec)
	defer p.setDiscarded()

	for loopCount := uint32(0); ; loopCount++ {

		canLoop, hasSignal := p.worker.CanLoopOrHasSignal(loopCount)
		if hasSignal || !canLoop {
			return hasSignal, StateUpdate{}, p.countAsyncCalls
		}

		current := p.s.step
		stateUpdate := current.Transition(p).ensureMarker(p.getMarker())
		sut := typeOfStateUpdateForMode(p.mode, stateUpdate)

		if sut.ShortLoop(p.s, stateUpdate, loopCount) {
			continue
		}

		sut.Prepare(p.s, &stateUpdate)
		return false, stateUpdate, p.countAsyncCalls
	}
}

/* ========================================================================= */

var _ ConstructionContext = &constructionContext{}

type constructionContext struct {
	contextTemplate
	s *Slot
}

func (p *constructionContext) SlotLink() SlotLink {
	return p.s.NewLink()
}

func (p *constructionContext) GetContext() context.Context {
	return p.s.ctx
}

func (p *constructionContext) SetContext(ctx context.Context) {
	if ctx == nil {
		panic("illegal value")
	}
	p.s.ctx = ctx
}

//func (p *constructionContext) GetContainer() SlotMachineState {
//	return p.machine.containerState
//}

func (p *constructionContext) ParentLink() SlotLink {
	return p.s.parent
}

func (p *constructionContext) SetParent(parent SlotLink) {
	p.s.parent = parent
}

func (p *constructionContext) executeCreate(nextCreate CreateFunc) StateMachine {
	p.setMode(updCtxConstruction)
	defer p.setDiscarded()

	return nextCreate(p)
}

/* ========================================================================= */

var _ InitializationContext = &initializationContext{}

type initializationContext struct {
	slotContext
}

func (p *initializationContext) executeInitialization(fn InitFunc) StateUpdate {
	p.setMode(updCtxInit)
	defer p.setDiscarded()

	return p.ensureAndPrepare(p.s, fn(p))
}
