//
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
//

package smachine

import "time"

var _ ExecutionContext = &executionContext{}

type executionContext struct {
	slotContext
	countAsyncCalls uint16
}

func (p *executionContext) GetPendingCallCount() int {
	p.ensure(updCtxExec)
	return int(p.s.asyncCallCount) + int(p.countAsyncCalls)
}

func (p *executionContext) Yield() StateConditionalBuilder {
	ncu := p.newConditionalUpdate(stateUpdNext)
	return &ncu
}

func (p *executionContext) Poll() StateConditionalBuilder {
	ncu := p.newConditionalUpdate(stateUpdPoll)
	return &ncu
}

func (p *executionContext) Sleep() StateConditionalBuilder {
	ncu := p.newConditionalUpdate(stateUpdSleep)
	return &ncu
}

func (p *executionContext) WaitAny() StateConditionalBuilder {
	ncu := p.newConditionalUpdate(stateUpdWaitForEvent)
	return &ncu
}

func (p *executionContext) WaitAnyUntil(until time.Time) StateConditionalBuilder {
	ncu := p.newConditionalUpdate(stateUpdWaitForEvent)

	ncu.until = p.s.machine.toRelativeTime(until)
	if ncu.until != 0 && !until.After(time.Now()) {
		ncu.decision = Passed
	}

	return &ncu
}

func (p *executionContext) newConditionalUpdate(updType stateUpdKind) conditionalUpdate {
	p.ensure(updCtxExec)
	return conditionalUpdate{template: newStateUpdateTemplate(p.mode, p.getMarker(), updType)}
}

func (p *executionContext) waitFor(link SlotLink, updMode stateUpdKind) StateConditionalBuilder {
	p.ensure(updCtxExec)
	if link.IsEmpty() {
		panic("illegal value")
		//		return &conditionalUpdate{marker: p.getMarker()}
	}

	switch isValid, isBusy := link.getIsValidAndBusy(); {
	case !isValid:
		ncu := p.newConditionalUpdate(stateUpdNext)
		ncu.decision = Impossible
		return &ncu
	case !isBusy:
		ncu := p.newConditionalUpdate(stateUpdNext)
		ncu.decision = Passed
		return &ncu
	default:
		ncu := p.newConditionalUpdate(updMode)
		ncu.decision = NotPassed
		ncu.dependency = link
		return &ncu
	}
}

func (p *executionContext) WaitActivation(link SlotLink) StateConditionalBuilder {
	return p.waitFor(link, stateUpdWaitForActive)
}

func (p *executionContext) WaitShared(link SharedDataLink) StateConditionalBuilder {
	return p.waitFor(link.link, stateUpdWaitForShared)
}

func (p *executionContext) UseShared(a SharedDataAccessor) SharedAccessReport {
	p.ensure(updCtxExec)

	switch a.accessLocal(p.s) {
	case Passed:
		return SharedSlotAvailableAlways
	case Impossible:
		return SharedSlotAbsent
	}
	return p.s.machine.useSlotAsShared(a.link, a.accessFn, p.w)
}

func (p *executionContext) executeNextStep() (stateUpdate StateUpdate, sut StateUpdateType, asyncCallCount uint16) {
	p.setMode(updCtxExec)
	defer func() {
		p.discardAndUpdate("execution", recover(), &stateUpdate)
	}()

	current := p.s.step

	stateUpdate = current.Transition(p).ensureMarker(p.getMarker())
	sut = typeOfStateUpdateForPrepare(p.mode, stateUpdate)
	sut.Prepare(p.s, &stateUpdate)

	return stateUpdate, sut, p.countAsyncCalls
}

/* ========================================================================= */

var _ ConditionalBuilder = &conditionalUpdate{}
var _ StateConditionalBuilder = &conditionalUpdate{}

type conditionalUpdate struct {
	template   StateUpdateTemplate
	kickOff    StepPrepareFunc
	dependency SlotLink
	until      uint32
	decision   Decision
}

func (c *conditionalUpdate) GetDecision() Decision {
	if c.decision.IsZero() {
		return NotPassed
	}
	return c.decision
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

func (c *conditionalUpdate) then(slotStep SlotStep) StateUpdate {
	if c.dependency.IsEmpty() {
		if c.until == 0 {
			return c.template.newStep(slotStep, c.kickOff)
		}
		return c.template.newStepUntil(slotStep, c.kickOff, c.until)
	} else {
		if c.until != 0 {
			panic("illegal value")
		}
		return c.template.newStepLink(slotStep, c.dependency)
	}
}
