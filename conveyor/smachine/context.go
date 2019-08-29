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

type slotContextMode uint8

const (
	inactiveContext slotContextMode = iota
	discardedContext
	constructContext
	initContext
	execContext
	migrateContext
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

func (p *slotContext) GetSelf() SlotLink {
	return p.s.NewLink()
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

func (p *slotContext) JumpOverride(fn StateFunc, mf MigrateFunc, flags StepFlags) StateUpdate {
	p.ensureAtLeastState(initContext)
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

var _ ConstructionContext = &constructionContext{}

type constructionContext struct {
	contextTemplate
	slotID SlotID
	parent SlotLink
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

func (p *migrationContext) Stay() StateUpdate {
	return stateUpdateNoChange(&p.marker)
}

func (p *migrationContext) executeMigrate(fn MigrateFunc) StateUpdate {
	p.setState(inactiveContext, migrateContext)
	defer p.setState(migrateContext, discardedContext)

	return EnsureUpdateContext(&p.marker, fn(p))
}

var _ InitializationContext = &initializationContext{}

type initializationContext struct {
	slotContext
}

func (p *initializationContext) executeInitialization(fn InitFunc) StateUpdate {
	p.setState(inactiveContext, initContext)
	defer p.setState(initContext, discardedContext)

	return EnsureUpdateContext(&p.marker, fn(p))
}

var _ ExecutionContext = &executionContext{}

type executionContext struct {
	slotContext
	worker          *SlotWorker
	countAsyncCalls uint16
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
	if link.IsEmpty() {
		panic("illegal value")
	}
	return &conditionalUpdate{marker: &p.marker, dependency: link}
}

func (p *executionContext) Wait() StateConditionalUpdate {
	p.ensureExactState(execContext)
	return &conditionalUpdate{marker: &p.marker}
}

func (p *executionContext) SyncOneStep(key string, weight int32, broadcastFn BroadcastReceiveFunc) Syncronizer {
	p.ensureExactState(execContext)
	return p.worker.machine.stepSync.Join(p.s, key, weight, broadcastFn)
}

func (p *executionContext) NewChild(fn CreateFunc) SlotLink {
	p.ensureExactState(execContext)
	if fn == nil {
		panic("illegal value")
	}
	_, link := p.worker.machine.applySlotCreate(nil, p.s.NewLink(), fn)
	return link
}

func (p *executionContext) executeNextStep() (stopNow bool, stateUpdate StateUpdate, asyncCount uint16) {
	p.setState(inactiveContext, execContext)
	defer p.setState(execContext, discardedContext)

	loopLimit := p.worker.GetLoopLimit()

	for loopCount := uint32(0); loopCount < loopLimit; loopCount++ {
		if p.worker.HasSignal() {
			return true, stateUpdate, p.countAsyncCalls
		}

		current := p.s.step
		stateUpdate = EnsureUpdateContext(&p.marker, current.Transition(p))

		if p.countAsyncCalls != 0 {
			break
		}
		updType, updParam := ExtractStateUpdateParam(stateUpdate)

		switch stateUpdType(updType) { // fast path(s)
		case stateUpdRepeat:
			limit := getRepeatLimit(updParam)
			if loopCount < limit {
				continue
			}
		case stateUpdNextLoop:
			ns := getShortLoopTransition(updParam)
			if ns == nil || !p.s.machine.IsConsecutive(current.Transition, ns) {
				break
			}
			p.s.incStep()
			_, ss, _ := ExtractStateUpdate(stateUpdate)
			p.s.setNextStep(ss)
			continue
		}
		break
	}

	return false, stateUpdate, p.countAsyncCalls
}

var _ AsyncResultContext = &asyncResultContext{}

type asyncResultContext struct {
	slot   *Slot
	wakeup bool
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

func (c *conditionalUpdate) ThenJumpOverride(fn StateFunc, mf MigrateFunc, flags StepFlags) StateUpdate {
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
	case stateUpdNext: // yield
		return stateUpdateYield(c.marker, slotStep, c.kickOff)
	case stateUpdPoll: // poll
		return stateUpdatePoll(c.marker, slotStep, c.kickOff)
	default: // wait
		if c.dependency.IsEmpty() {
			return stateUpdateWait(c.marker, slotStep, c.kickOff)
		}
		return stateUpdateWaitForSlot(c.marker, c.dependency, slotStep, c.kickOff)
	}
}
