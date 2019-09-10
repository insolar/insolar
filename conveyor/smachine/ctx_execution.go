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

import "context"

type slotContextMode uint8

const (
	inactiveContext slotContextMode = iota
	discardedContext
	constructContext
	initContext
	execContext
	migrateContext
	bargeInContext
)

type contextTemplate struct {
	marker struct{}
	mode   slotContextMode
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

func (p *slotContext) GetParent() SlotLink {
	return p.s.parent
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

func (p *slotContext) JumpExt(fn StateFunc, mf MigrateFunc, flags StepFlags) StateUpdate {
	p.ensureAtLeastState(initContext)
	if fn == nil {
		panic("illegal value")
	}
	return stateUpdateNext(&p.marker, fn, mf, true, flags)
}

func (p *slotContext) Jump(fn StateFunc) StateUpdate {
	p.ensureAtLeastState(initContext)
	return stateUpdateNext(&p.marker, fn, nil, true, 0)
}

func (p *slotContext) Stop() StateUpdate {
	p.ensureAtLeastState(initContext)
	return stateUpdateStop(&p.marker)
}

func (p *slotContext) Replace(fn CreateFunc) StateUpdate {
	p.ensureAtLeastState(migrateContext)
	if fn == nil {
		panic("illegal value")
	}
	return stateUpdateReplace(&p.marker, fn)
}

func (p *slotContext) Repeat(limit int) StateUpdate {
	p.ensureExactState(execContext)
	return stateUpdateRepeat(&p.marker, limit)
}

func (p *slotContext) Stay() StateUpdate {
	p.ensureAtLeastState(migrateContext)
	return stateUpdateNoChange(&p.marker)
}

func (p *slotContext) WakeUp() StateUpdate {
	p.ensureAtLeastState(migrateContext)
	return stateUpdateRepeat(&p.marker, 0)
}

func (p *slotContext) Share(data interface{}, wakeUpOnUse bool) SharedDataLink {
	p.ensureAtLeastState(initContext)
	return SharedDataLink{p.s.NewStepLink(), wakeUpOnUse, data}
}

var _ ConstructionContext = &constructionContext{}

type constructionContext struct {
	contextTemplate
	ctx    context.Context
	slotID SlotID
	parent SlotLink
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

	return fn(p).ensureMarker(&p.marker)
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

	return fn(p).ensureMarker(&p.marker)
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
	return &conditionalUpdate{marker: &p.marker, updMode: stateUpdNext}
}

func (p *executionContext) Poll() StateConditionalUpdate {
	p.ensureExactState(execContext)
	return &conditionalUpdate{marker: &p.marker, updMode: stateUpdPoll}
}

func (p *executionContext) WaitForActive(link SlotLink) StateConditionalUpdate {
	p.ensureExactState(execContext)
	p.s.ensureLocal(link)
	return &conditionalUpdate{marker: &p.marker, updMode: stateUpdWait, dependency: link}
}

func (p *executionContext) Wait() StateConditionalUpdate {
	p.ensureExactState(execContext)
	return &conditionalUpdate{marker: &p.marker}
}

func (p *executionContext) WaitForShared(link SharedDataLink) StateConditionalUpdate {
	p.ensureExactState(execContext)
	r, releaseFn := p.worker.AttachToShared(p.s, link.link, link.wakeup)
	if releaseFn != nil {
		defer releaseFn()
	}
	if r == SharedDataBusyRemote {
		panic("not implemented")
	}
	return &conditionalUpdate{marker: &p.marker, dependency: link.link.SlotLink}
}

func (p *executionContext) UseShared(a SharedDataAccessor) SharedAccessReport {
	p.ensureExactState(execContext)
	r, releaseFn := p.worker.AttachToShared(p.s, a.link.link, a.link.wakeup)
	if releaseFn != nil {
		defer releaseFn()
	}
	if r >= SharedDataAvailableLocal {
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
		stateUpdate := current.Transition(p).ensureMarker(&p.marker)

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
	marker     *struct{}
	kickOff    StepPrepareFunc
	dependency SlotLink
	updMode    stateUpdType
	//flags      stepFlags
}

func (c *conditionalUpdate) ThenJump(fn StateFunc) StateUpdate {
	if fn == nil {
		panic("illegal value")
	}
	return c.then(fn, nil, 0)
}

func (c *conditionalUpdate) ThenJumpExt(fn StateFunc, mf MigrateFunc, flags StepFlags) StateUpdate {
	if fn == nil {
		panic("illegal value")
	}
	return c.then(fn, mf, 0)
}

func (c *conditionalUpdate) ThenRepeat() StateUpdate {
	return c.then(nil, nil, 0)
}

func (c *conditionalUpdate) then(fn StateFunc, mf MigrateFunc, flags StepFlags) StateUpdate {
	slotStep := SlotStep{Transition: fn, Migration: mf, StepFlags: flags}
	switch c.updMode {
	case stateUpdNext: // Yield
		return stateUpdateYield(c.marker, slotStep, c.kickOff)
	case stateUpdPoll: // Poll
		return stateUpdatePoll(c.marker, slotStep, c.kickOff)
	case stateUpdWait: // Wait & WaitForActive
		if c.dependency.IsEmpty() {
			return stateUpdateWait(c.marker, slotStep, c.kickOff)
		}
		return stateUpdateWaitForSlot(c.marker, c.dependency, slotStep, c.kickOff)
	default: // WaitForShared
		if c.kickOff != nil {
			panic("illegal state")
		}
		return stateUpdateWaitForShared(c.marker, c.dependency, slotStep)
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

func (b bargeInRequest) WithJumpExt(fn StateFunc, mf MigrateFunc, sf StepFlags) BargeInFunc {
	b.u.ensureAtLeastState(initContext)
	if fn == nil {
		panic("illegal value")
	}
	bfn := b.m.createBargeIn(b.link, func(ctx BargeInContext) StateUpdate {
		return ctx.JumpExt(fn, mf, sf)
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
	r.StepFlags |= StepResetAllFlags
	return p.s.step
}

func (p *bargingInContext) executeBargeIn(bargeInFn BargeInApplyFunc) StateUpdate {
	p.setState(inactiveContext, bargeInContext)
	defer p.setState(bargeInContext, discardedContext)

	return bargeInFn(p).ensureMarker(&p.marker)
}
