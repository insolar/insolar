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

import "context"

type bargeInBuilder struct {
	c      slotContext
	parent *slotContext
	link   StepLink
}

func (b bargeInBuilder) with(stateUpdate StateUpdate) BargeInFunc {
	b.parent.ensureValid()
	defer b.c.setDiscarded()
	return b.c.s.machine.createLightBargeIn(b.link,
		b.c.ensureAndPrepare(b.c.s, stateUpdate))
}

func (b bargeInBuilder) WithError(err error) BargeInFunc {
	return b.with(b.c.Error(err))
}

func (b bargeInBuilder) WithStop() BargeInFunc {
	return b.with(b.c.Stop())
}

func (b bargeInBuilder) WithWakeUp() BargeInFunc {
	return b.with(b.c.WakeUp())
}

func (b bargeInBuilder) WithJumpExt(step SlotStep) BargeInFunc {
	return b.with(b.c.JumpExt(step))
}

func (b bargeInBuilder) WithJump(fn StateFunc) BargeInFunc {
	return b.with(b.c.Jump(fn))
}

/* ========================================================================= */

var _ BargeInContext = &bargingInContext{}

type bargingInContext struct {
	slotContext
	param      interface{}
	atOriginal bool
}

func (p *bargingInContext) BargeInParam() interface{} {
	p.ensure(updCtxBargeIn)
	return p.param
}

func (p *bargingInContext) IsAtOriginalStep() bool {
	p.ensure(updCtxBargeIn)
	return p.atOriginal
}

func (p *bargingInContext) executeBargeIn(fn BargeInApplyFunc) (stateUpdate StateUpdate) {
	p.setMode(updCtxBargeIn)
	defer func() {
		p.discardAndUpdate("barge in", recover(), &stateUpdate)
	}()

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
	p.ensure(updCtxAsyncCallback)
	return p.slot.NewLink()
}

func (p *asyncResultContext) ParentLink() SlotLink {
	p.ensure(updCtxAsyncCallback)
	return p.slot.parent
}

func (p *asyncResultContext) GetContext() context.Context {
	p.ensure(updCtxAsyncCallback)
	return p.slot.ctx
}

func (p *asyncResultContext) WakeUp() {
	p.ensure(updCtxAsyncCallback)
	p.wakeup = true
}

func (p *asyncResultContext) executeResult(fn AsyncResultFunc) bool {
	p.setMode(updCtxAsyncCallback)
	defer p.setDiscarded()

	fn(p)
	return p.wakeup
}
