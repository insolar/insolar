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

type bargeInRequest struct {
	p    *contextTemplate
	m    *SlotMachine
	link StepLink
}

func (b bargeInRequest) WithWakeUp() BargeInFunc {
	b.p.ensureAny2(updCtxExec, updCtxInit)
	bfn := b.m.createBargeIn(b.link, func(ctx BargeInContext) StateUpdate {
		return ctx.WakeUp()
	})
	return func() bool {
		return bfn(nil)
	}
}

func (b bargeInRequest) WithJumpExt(step SlotStep) BargeInFunc {
	b.p.ensureAny2(updCtxExec, updCtxInit)
	step.ensureTransition()

	bfn := b.m.createBargeIn(b.link, func(ctx BargeInContext) StateUpdate {
		return ctx.JumpExt(step)
	})
	return func() bool {
		return bfn(nil)
	}
}

func (b bargeInRequest) WithJump(fn StateFunc) BargeInFunc {
	b.p.ensureAny2(updCtxExec, updCtxInit)
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

/* ========================================================================= */

var _ BargeInContext = &bargingInContext{}

type bargingInContext struct {
	slotContext
	param      interface{}
	atOriginal bool
}

func (p *bargingInContext) GetBargeInParam() interface{} {
	p.ensure(updCtxBargeIn)
	return p.param
}

func (p *bargingInContext) IsAtOriginalStep() bool {
	p.ensure(updCtxBargeIn)
	return p.atOriginal
}

func (p *bargingInContext) AffectedStep() SlotStep {
	p.ensure(updCtxBargeIn)
	r := p.s.step
	r.Flags |= StepResetAllFlags
	return p.s.step
}

func (p *bargingInContext) executeBargeIn(fn BargeInApplyFunc) StateUpdate {
	p.setMode(updCtxBargeIn)
	defer p.setDiscarded()

	return p.ensureAndPrepare(p.s, fn(p))
}

/* ========================================================================= */

var _ AsyncResultContext = &asyncResultContext{}

type asyncResultContext struct {
	contextTemplate
	slot   *Slot
	wakeup bool
}

func (p *asyncResultContext) SlotLink() SlotLink {
	p.setMode(updCtxAsyncCallback)
	return p.slot.NewLink()
}

func (p *asyncResultContext) ParentLink() SlotLink {
	return p.slot.parent
}

func (p *asyncResultContext) GetContext() context.Context {
	p.setMode(updCtxAsyncCallback)
	return p.slot.ctx
}

func (p *asyncResultContext) WakeUp() {
	p.setMode(updCtxAsyncCallback)
	p.wakeup = true
}

func (p *asyncResultContext) executeResult(fn AsyncResultFunc) bool {
	p.setMode(updCtxAsyncCallback)
	defer p.setDiscarded()

	fn(p)
	return p.wakeup
}
