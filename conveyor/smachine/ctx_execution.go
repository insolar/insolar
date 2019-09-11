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
	"unsafe"
)

type slotContextMode uint8

const (
	inactiveContext slotContextMode = iota
	discardedContext
	constructContext
	initContext
	execContext
	migrateContext
	bargeInContext
	failContext
)

type contextTemplate struct {
	mode slotContextMode
}

func (p *contextTemplate) getMarker() ContextMarker {
	return ContextMarker(uintptr(unsafe.Pointer(p)))
}

func (p *contextTemplate) ensureExactState(expected slotContextMode) {
	if p.mode != expected {
		panic("illegal state")
	}
}

func (p *contextTemplate) ensureAtLeastState(s slotContextMode) {
	if p.mode >= s {
		return
	}
	panic("illegal state")
}

func (p *contextTemplate) setState(expected, updated slotContextMode) {
	p.ensureExactState(expected)
	p.mode = updated
}

type slotContext struct {
	contextTemplate
	s *Slot
}

func (p *slotContext) GetSlotID() SlotID {
	return p.s.GetID()
}

func (p *slotContext) SlotLink() SlotLink {
	return p.s.NewLink()
}

func (p *slotContext) GetContext() context.Context {
	return p.s.ctx
}

func (p *slotContext) GetContainer() SlotMachineState {
	return p.s.machine.containerState
}

func (p *slotContext) GetParent() SlotLink {
	return p.s.parent
}

func (p *slotContext) SetDefaultErrorHandler(fn ErrorHandlerFunc) {
	p.ensureAtLeastState(initContext)
	p.s.defErrorHandler = fn
}

func (p *slotContext) SetDefaultMigration(fn MigrateFunc) {
	p.ensureAtLeastState(initContext)
	p.s.defMigrate = fn
}

func (p *slotContext) SetDefaultFlags(flags StepFlags) {
	p.ensureAtLeastState(initContext)
	if flags&StepResetAllFlags != 0 {
		p.s.defFlags = flags &^ StepResetAllFlags
	} else {
		p.s.defFlags |= flags
	}
}

func (p *slotContext) JumpExt(step SlotStep) StateUpdate {
	p.ensureAtLeastState(initContext)
	return stateUpdateNext(p.getMarker(), step, true)
}

func (p *slotContext) Jump(fn StateFunc) StateUpdate {
	p.ensureAtLeastState(initContext)
	return stateUpdateNext(p.getMarker(), SlotStep{Transition: fn}, true)
}

func (p *slotContext) Stop() StateUpdate {
	p.ensureAtLeastState(initContext)
	return stateUpdateStop(p.getMarker())
}

func (p *slotContext) Error(err error) StateUpdate {
	p.ensureAtLeastState(initContext)
	return stateUpdateError(p.getMarker(), err)
}

func (p *slotContext) Replace(fn CreateFunc) StateUpdate {
	p.ensureAtLeastState(migrateContext)
	return stateUpdateReplace(p.getMarker(), fn)
}

func (p *slotContext) ReplaceWith(sm StateMachine) StateUpdate {
	p.ensureAtLeastState(migrateContext)
	return stateUpdateReplaceWith(p.getMarker(), sm)
}

func (p *slotContext) Repeat(limit int) StateUpdate {
	p.ensureExactState(execContext)
	return stateUpdateRepeat(p.getMarker(), limit)
}

func (p *slotContext) Stay() StateUpdate {
	p.ensureAtLeastState(migrateContext)
	return stateUpdateNoChange(p.getMarker())
}

func (p *slotContext) WakeUp() StateUpdate {
	p.ensureAtLeastState(migrateContext)
	return stateUpdateRepeat(p.getMarker(), 0)
}

func (p *slotContext) Share(data interface{}, wakeUpOnUse bool) SharedDataLink {
	p.ensureAtLeastState(initContext)
	return SharedDataLink{p.s.NewStepLink(), wakeUpOnUse, data}
}

func (p *slotContext) AffectedStep() SlotStep {
	p.ensureExactState(migrateContext)
	r := p.s.step
	r.Flags |= StepResetAllFlags
	return p.s.step
}

var _ ConstructionContext = &constructionContext{}

type constructionContext struct {
	contextTemplate
	ctx     context.Context
	slotID  SlotID
	parent  SlotLink
	machine *SlotMachine
}

func (p *constructionContext) SetContext(ctx context.Context) {
	if ctx == nil {
		panic("illegal value")
	}
	p.ctx = ctx
}

func (p *constructionContext) GetContext() context.Context {
	return p.ctx
}

func (p *constructionContext) GetContainer() SlotMachineState {
	return p.machine.containerState
}

func (p *constructionContext) GetSlotID() SlotID {
	if p.slotID == 0 {
		panic("illegal state")
	}
	return p.slotID
}

func (p *constructionContext) GetParent() SlotLink {
	return p.parent
}

func (p *constructionContext) executeCreate(nextCreate CreateFunc) StateMachine {
	p.setState(inactiveContext, constructContext)
	defer p.setState(constructContext, discardedContext)

	return nextCreate(p)
}

var _ MigrationContext = &migrationContext{}

type migrationContext struct {
	slotContext
}

func (p *migrationContext) executeMigrate(fn MigrateFunc) StateUpdate {
	p.setState(inactiveContext, migrateContext)
	defer p.setState(migrateContext, discardedContext)

	return fn(p).ensureMarker(p.getMarker())
}

var _ InitializationContext = &initializationContext{}

type initializationContext struct {
	slotContext
	machine *SlotMachine
}

func (p *initializationContext) BargeInWithParam(applyFn BargeInApplyFunc) BargeInParamFunc {
	p.ensureAtLeastState(initContext)

	return p.machine.createBargeIn(p.s.NewStepLink().AnyStep(), applyFn)
}

func (p *initializationContext) BargeIn() BargeInRequester {
	p.ensureAtLeastState(initContext)
	return &bargeInRequest{&p.contextTemplate, p.machine, p.s.NewStepLink().AnyStep()}
}

func (p *initializationContext) executeInitialization(fn InitFunc) StateUpdate {
	p.setState(inactiveContext, initContext)
	defer p.setState(initContext, discardedContext)

	return fn(p).ensureMarker(p.getMarker())
}

var _ ExecutionContext = &executionContext{}

type executionContext struct {
	slotContext
	machine         *SlotMachine
	worker          WorkerContext
	countAsyncCalls uint16
}

func (p *executionContext) StepLink() StepLink {
	return p.s.NewStepLink()
}

func (p *executionContext) GetPendingCallCount() int {
	return int(p.s.asyncCallCount) + int(p.countAsyncCalls)
}

func (p *executionContext) Yield() StateConditionalUpdate {
	p.ensureExactState(execContext)
	return &conditionalUpdate{marker: p.getMarker(), updMode: stateUpdNext}
}

func (p *executionContext) Poll() StateConditionalUpdate {
	p.ensureExactState(execContext)
	return &conditionalUpdate{marker: p.getMarker(), updMode: stateUpdPoll}
}

func (p *executionContext) Sleep() StateConditionalUpdate {
	p.ensureExactState(execContext)
	return &conditionalUpdate{marker: p.getMarker(), updMode: stateUpdSleep}
}

func (p *executionContext) WaitForEvent() StateConditionalUpdate {
	p.ensureExactState(execContext)
	return &conditionalUpdate{marker: p.getMarker(), updMode: stateUpdWaitForEvent}
}

func (p *executionContext) waitFor(link SlotLink, updMode stateUpdType) StateConditionalUpdate {
	p.ensureExactState(execContext)
	if link.IsEmpty() {
		panic("illegal value")
		//		return &conditionalUpdate{marker: p.getMarker()}
	}

	r, releaseFn := p.worker.AttachTo(p.s, link, false)
	if releaseFn != nil {
		defer releaseFn()
	}
	switch r {
	case SharedSlotAbsent, SharedSlotLocalAvailable:
		// no wait
		return &conditionalUpdate{marker: p.getMarker()}

	case SharedSlotRemoteBusy, SharedSlotRemoteAvailable:
		// state has to be re-detected upon applying the update
		panic("not implemented")
	}
	return &conditionalUpdate{marker: p.getMarker(), updMode: updMode, dependency: link}
}

func (p *executionContext) WaitForActive(link SlotLink) StateConditionalUpdate {
	return p.waitFor(link, stateUpdWaitForActive)
}

func (p *executionContext) WaitForShared(link SharedDataLink) StateConditionalUpdate {
	return p.waitFor(link.link.SlotLink, stateUpdWaitForShared)
}

func (p *executionContext) UseShared(a SharedDataAccessor) SharedAccessReport {
	p.ensureExactState(execContext)
	r, releaseFn := p.worker.AttachTo(p.s, a.link.link.SlotLink, a.link.wakeup)
	if releaseFn != nil {
		defer releaseFn()
	}
	if r.IsAvailable() {
		a.accessFn(a.link.data)
	}
	return r
}

//func (p *executionContext) SyncOneStep(key string, weight int32, broadcastFn BroadcastReceiveFunc) Syncronizer {
//	p.ensureExactState(execContext)
//	return p.machine.stepSync.Join(p.s, key, weight, broadcastFn)
//}

func (p *executionContext) NewChild(ctx context.Context, fn CreateFunc) SlotLink {
	p.ensureExactState(execContext)
	if fn == nil {
		panic("illegal value")
	}
	if ctx == nil {
		panic("illegal value")
	}
	_, link := p.machine.applySlotCreate(ctx, nil, p.s.NewLink(), fn)
	if link.IsEmpty() {
		panic("machine was not created")
	}
	return link
}

func (p *executionContext) BargeInWithParam(applyFn BargeInApplyFunc) BargeInParamFunc {
	p.ensureAtLeastState(execContext)

	return p.machine.createBargeIn(p.s.NewStepLink().AnyStep(), applyFn)
}

func (p *executionContext) BargeIn() BargeInRequester {
	p.ensureAtLeastState(execContext)
	return &bargeInRequest{&p.contextTemplate, p.machine, p.s.NewStepLink().AnyStep()}
}

func (p *executionContext) BargeInThisStepOnly() BargeInRequester {
	p.ensureExactState(execContext)
	return &bargeInRequest{&p.contextTemplate, p.machine, p.s.NewStepLink()}
}

func (p *executionContext) executeNextStep() (bool, StateUpdate, uint16) {
	p.setState(inactiveContext, execContext)
	defer p.setState(execContext, discardedContext)

	for loopCount := uint32(0); ; loopCount++ {

		canLoop, hasSignal := p.worker.CanLoopOrHasSignal(loopCount)
		if hasSignal || !canLoop {
			return hasSignal, StateUpdate{}, p.countAsyncCalls
		}

		current := p.s.step
		stateUpdate := current.Transition(p).ensureMarker(p.getMarker())

		switch stateUpdType(stateUpdate.updType) { // fast path(s)
		case stateUpdRepeat:
			if loopCount < stateUpdate.param0 {
				continue
			}
		case stateUpdNextLoop:
			ns := stateUpdate.step.Transition
			if ns != nil && !p.s.declaration.IsConsecutive(current.Transition, ns) {
				break
			}
			p.s.setNextStep(stateUpdate.step)
			if loopCount < stateUpdate.param0 {
				continue
			}
			continue
		}
		return false, stateUpdate, p.countAsyncCalls
	}
}

var _ AsyncResultContext = &asyncResultContext{}

type asyncResultContext struct {
	slot   *Slot
	wakeup bool
}

func (p *asyncResultContext) GetContext() context.Context {
	return p.slot.ctx
}

func (p *asyncResultContext) GetContainer() SlotMachineState {
	return p.slot.machine.containerState
}

func (p *asyncResultContext) WakeUp() {
	p.wakeup = true
}

func (p *asyncResultContext) GetSlotID() SlotID {
	return p.slot.GetID()
}

func (p *asyncResultContext) GetParent() SlotLink {
	return p.slot.parent
}

func (p *asyncResultContext) executeResult(fn AsyncResultFunc) bool {
	fn(p)
	return p.wakeup
}

func (p *asyncResultContext) executeBroadcast(payload interface{}, fn BroadcastReceiveFunc) (accepted, wakeup bool) {
	accepted = fn(p, payload)
	wakeup = p.wakeup
	return
}

var _ ConditionalUpdate = &conditionalUpdate{}

type conditionalUpdate struct {
	marker     ContextMarker
	kickOff    StepPrepareFunc
	dependency SlotLink
	updMode    stateUpdType
	//flags      stepFlags
}

func (c *conditionalUpdate) ThenJump(fn StateFunc) StateUpdate {
	return c.ThenJumpExt(SlotStep{Transition: fn})
}

func (c *conditionalUpdate) ThenJumpExt(step SlotStep) StateUpdate {
	step.ensureTransition()
	return c.then(step)
}

func (c *conditionalUpdate) ThenRepeat() StateUpdate {
	return c.then(SlotStep{})
}

func (c *conditionalUpdate) IsAvailable() bool {
	return c.updMode == 0
}

func (c *conditionalUpdate) then(slotStep SlotStep) StateUpdate {
	switch c.updMode {
	case stateUpdNext: // Yield & Call
		return stateUpdateYield(c.marker, slotStep, c.kickOff)
	case stateUpdPoll:
		return stateUpdatePoll(c.marker, slotStep, c.kickOff)
	case stateUpdSleep:
		return stateUpdateSleep(c.marker, slotStep, c.kickOff)

	case stateUpdWaitForEvent: // WaitForEvent
		return stateUpdateWaitForEvent(c.marker, slotStep, c.kickOff)

	case stateUpdWaitForActive: // WaitForActive
		if c.kickOff != nil {
			panic("illegal value")
		}
		return stateUpdateWaitForSlot(c.marker, c.dependency, slotStep)

	case stateUpdWaitForShared: // WaitForShared
		if c.kickOff != nil {
			panic("illegal value")
		}
		return stateUpdateWaitForShared(c.marker, c.dependency, slotStep)

	case 0:
		// WaitForShared or WaitForActive with true condition
		if c.kickOff != nil {
			panic("illegal value")
		}
		return stateUpdateNext(c.marker, slotStep, true)

	default:
		panic("illegal value")
	}
}

type bargeInRequest struct {
	u    *contextTemplate
	m    *SlotMachine
	link StepLink
}

func (b bargeInRequest) WithWakeUp() BargeInFunc {
	b.u.ensureAtLeastState(initContext)
	bfn := b.m.createBargeIn(b.link, func(ctx BargeInContext) StateUpdate {
		return ctx.WakeUp()
	})
	return func() bool {
		return bfn(nil)
	}
}

func (b bargeInRequest) WithJumpExt(step SlotStep) BargeInFunc {
	b.u.ensureAtLeastState(initContext)
	step.ensureTransition()

	bfn := b.m.createBargeIn(b.link, func(ctx BargeInContext) StateUpdate {
		return ctx.JumpExt(step)
	})
	return func() bool {
		return bfn(nil)
	}
}

func (b bargeInRequest) WithJump(fn StateFunc) BargeInFunc {
	b.u.ensureAtLeastState(initContext)
	if fn == nil {
		panic("illegal value")
	}
	bfn := b.m.createBargeIn(b.link, func(ctx BargeInContext) StateUpdate {
		return ctx.Jump(fn)
	})
	return func() bool {
		return bfn(nil)
	}
}

var _ BargeInContext = &bargingInContext{}

type bargingInContext struct {
	slotContext
	param      interface{}
	atOriginal bool
}

func (p *bargingInContext) GetBargeInParam() interface{} {
	p.ensureExactState(bargeInContext)
	return p.param
}

func (p *bargingInContext) IsAtOriginalStep() bool {
	p.ensureExactState(bargeInContext)
	return p.atOriginal
}

func (p *bargingInContext) AffectedStep() SlotStep {
	p.ensureExactState(bargeInContext)
	r := p.s.step
	r.Flags |= StepResetAllFlags
	return p.s.step
}

func (p *bargingInContext) executeBargeIn(bargeInFn BargeInApplyFunc) StateUpdate {
	p.setState(inactiveContext, bargeInContext)
	defer p.setState(bargeInContext, discardedContext)

	return bargeInFn(p).ensureMarker(p.getMarker())
}

var _ FailureContext = &failureContext{}

type failureContext struct {
	slotContext
	machine *SlotMachine
	isPanic bool
	isAsync bool
	err     error
}

func (p *failureContext) GetError() (isPanic, isAsync bool, err error) {
	p.ensureExactState(failContext)
	return p.isPanic, false, p.err
}

func (p *failureContext) NewChild(ctx context.Context, fn CreateFunc) SlotLink {
	p.ensureExactState(failContext)
	if fn == nil {
		panic("illegal value")
	}
	if ctx == nil {
		panic("illegal value")
	}
	_, link := p.machine.applySlotCreate(ctx, nil, p.s.NewLink(), fn)
	if link.IsEmpty() {
		panic("machine was not created")
	}
	return link
}

func (p *failureContext) executeErrorHandlerSafe(handlerFunc ErrorHandlerFunc) (err error) {
	p.setState(inactiveContext, failContext)
	defer func() {
		p.setState(failContext, discardedContext)
		err = recoverSlotPanic("error handling has failed", recover(), err)
	}()

	handlerFunc(p)
	return nil
}
