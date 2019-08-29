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
	p.s.migrateSlot = fn
}

func (p *slotContext) JumpWithMigrate(fn StateFunc, mf MigrateFunc) StateUpdate {
	p.ensureAtLeastState(initContext)
	return stateUpdateNext(&p.marker, fn, mf, true, 0)
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

func (p *slotContext) Yield() StateUpdate {
	p.ensureExactState(execContext)
	return stateUpdateRepeat(&p.marker, 0)
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

func (p *migrationContext) Same() StateUpdate {
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

func (p *executionContext) WaitAny() StateConditionalUpdate {
	p.ensureExactState(execContext)
	return &conditionalUpdate{marker: &p.marker}
}

func (p *executionContext) SyncOneStep(key string, weight int32, broadcastFn BroadcastReceiveFunc) Syncronizer {
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

//func (p *executionContext) NextAdapterCall(a ExecutionAdapter, fn AdapterCallFunc, resultState StateFunc) (StateUpdate, context.CancelFunc) {
//	p.ensureExactState(execContext)
//	if resultState == nil {
//		panic("illegal value")
//	}
//	aq := p.worker.machine.GetAdapterQueue(a)
//
//	cf := &indirectCancel{}
//
//	stepLink := p.s.NewStepLink()
//	return StateUpdate{marker: &p.marker,
//		flags:    stateUpdateColdWait | stateUpdateHasAsync,
//		nextStep: SlotStep{transition: resultState},
//
//		param: func() {
//			cf.set(aq.CallAsyncWithCancel(stepLink, fn, func(fn AsyncResultFunc) {
//				p.worker.machine.applyAsyncStateUpdate(stepLink, fn)
//			}))
//		}}, cf.cancel
//}
//
//func (p *executionContext) AdapterSyncCall(a ExecutionAdapter, fn AdapterCallFunc) bool {
//	p.ensureExactState(execContext)
//	aq := p.worker.machine.GetAdapterQueue(a)
//
//	wc := p.worker.getCond()
//
//	var resultFn AsyncResultFunc
//	var stateFlag uint32
//
//	stepLink := p.s.NewStepLink()
//	aq.CallAsync(stepLink, fn, func(fn AsyncResultFunc) {
//		resultFn = fn
//		if !atomic.CompareAndSwapUint32(&stateFlag, 0, 1) {
//			return
//		}
//		wc.L.Lock()
//		wc.Broadcast()
//		wc.L.Unlock()
//	})
//
//	wc.L.Lock()
//	wc.Wait()
//	wc.L.Unlock()
//
//	if atomic.CompareAndSwapUint32(&stateFlag, 0, 2) {
//		stepLink.setCancelled()
//		return false
//	}
//	if resultFn == nil {
//		return false
//	}
//
//	rc := asyncResultContext{slot: p.s}
//	rc.executeResult(resultFn)
//	return true
//}
//
//func (p *executionContext) AdapterAsyncCall(a ExecutionAdapter, fn AdapterCallFunc) {
//	p.ensureExactState(execContext)
//	aq := p.worker.machine.GetAdapterQueue(a)
//
//	stepLink := p.s.NewStepLink()
//	p.countAsyncCalls++
//
//	aq.CallAsync(stepLink, fn, func(fn AsyncResultFunc) {
//		p.worker.machine.applyAsyncStateUpdate(stepLink, fn)
//	})
//}
//
//func (p *executionContext) AdapterAsyncCallWithCancel(a ExecutionAdapter, fn AdapterCallFunc) context.CancelFunc {
//	p.ensureExactState(execContext)
//	aq := p.worker.machine.GetAdapterQueue(a)
//
//	stepLink := p.s.NewStepLink()
//	p.countAsyncCalls++
//
//	return aq.CallAsyncWithCancel(stepLink, fn, func(fn AsyncResultFunc) {
//		p.worker.machine.applyAsyncStateUpdate(stepLink, fn)
//	})
//}

func (p *executionContext) executeNextStep() (stopNow bool, stateUpdate StateUpdate, asyncCount uint16) {
	p.setState(inactiveContext, execContext)
	defer p.setState(execContext, discardedContext)

	loopLimit := p.worker.GetLoopLimit()

	for loopCount := uint32(0); loopCount < loopLimit; loopCount++ {
		if p.worker.HasSignal() {
			return true, stateUpdate, p.countAsyncCalls
		}

		current := p.s.nextStep
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
			ns := getShortLoopStep(updParam)
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

func (p *asyncResultContext) GetSlotID() SlotID {
	return p.slot.GetID()
}

func (p *asyncResultContext) GetParent() SlotLink {
	return p.slot.parent
}

func (p *asyncResultContext) WakeUp() {
	p.wakeup = true
}

func (p *asyncResultContext) WakeUpAndJump() {
	panic("implement me")
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
	kickOff    func()
	dependency SlotLink
	poll       bool
	flags      stepFlags
}

func (c *conditionalUpdate) Poll() ConditionalUpdate {
	r := *c
	r.poll = true
	return &r
}

func (c *conditionalUpdate) Active(dependency SlotLink) ConditionalUpdate {
	r := *c
	r.dependency = dependency
	return &r
}

func (c *conditionalUpdate) AsyncJump(enable bool) CallConditionalUpdate {
	r := *c
	if enable {
		r.flags |= stepFlagAllowPreempt
	} else {
		r.flags &^= stepFlagAllowPreempt
	}
	return &r
}

func (c *conditionalUpdate) WakeUp(enable bool) ConditionalUpdate {
	r := *c
	if enable {
		r.flags = (r.flags &^ stepFlagAwakeMask) | stepFlagAwakeDefault
	} else {
		r.flags = (r.flags &^ stepFlagAwakeMask) | stepFlagAwakeDisable
	}
	return &r
}

func (c *conditionalUpdate) WakeUpAlways() ConditionalUpdate {
	r := *c
	r.flags = (r.flags &^ stepFlagAwakeMask) | stepFlagAwakeAlways
	return &r
}

func (c *conditionalUpdate) ThenJump(fn StateFunc) StateUpdate {
	if fn == nil {
		panic("illegal value")
	}
	return c.then(fn, nil)
}

func (c *conditionalUpdate) ThenJumpWithMigrate(fn StateFunc, mf MigrateFunc) StateUpdate {
	if fn == nil {
		panic("illegal value")
	}
	return c.then(fn, mf)
}

func (c *conditionalUpdate) ThenRepeat() StateUpdate {
	return c.then(nil, nil)
}

func (c *conditionalUpdate) then(fn StateFunc, mf MigrateFunc) StateUpdate {

	// TODO apply kickOff
	if c.kickOff != nil {
		panic("not implemented")
	}

	slotStep := SlotStep{Transition: fn, Migration: mf}
	switch {
	case fn == nil: // repeat
		panic("not implemented") // TODO repeat
	case c.dependency.IsEmpty():
		if c.poll {
			return stateUpdatePoll(c.marker, slotStep)
		}
		return stateUpdateWait(c.marker, slotStep)
	case c.poll:
		panic("not supported")
	default:
		return stateUpdateWaitForSlot(c.marker, c.dependency, slotStep)
	}
}
