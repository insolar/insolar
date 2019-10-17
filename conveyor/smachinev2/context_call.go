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
	"github.com/insolar/insolar/network/consensus/common/syncrun"
)

func NewExecutionAdapter(adapterID AdapterID, executor AdapterExecutor) ExecutionAdapter {
	if adapterID.IsEmpty() {
		panic("illegal value")
	}
	if executor == nil {
		panic("illegal value")
	}
	return adapterExecHelper{adapterID, executor}
}

var _ ExecutionAdapter = &adapterExecHelper{}

type adapterExecHelper struct {
	adapterID AdapterID
	executor  AdapterExecutor
}

func (p adapterExecHelper) IsEmpty() bool {
	return p.adapterID.IsEmpty()
}

func (p adapterExecHelper) GetAdapterID() AdapterID {
	return p.adapterID
}

func (p adapterExecHelper) PrepareSync(ctx ExecutionContext, fn AdapterCallFunc) SyncCallRequester {
	return &adapterCallRequest{ctx: ctx.(*executionContext), fn: fn, executor: p.executor, mode: adapterSyncCallContext}
}

func (p adapterExecHelper) PrepareAsync(ctx ExecutionContext, fn AdapterCallFunc) AsyncCallRequester {
	return &adapterCallRequest{ctx: ctx.(*executionContext), fn: fn, executor: p.executor, mode: adapterAsyncCallContext}
}

func (p adapterExecHelper) PrepareNotify(ctx ExecutionContext, fn AdapterNotifyFunc) NotifyRequester {
	return &adapterNotifyRequest{ctx: ctx.(*executionContext), fn: fn, executor: p.executor, mode: adapterAsyncCallContext}
}

const (
	adapterSyncCallContext     = 1
	adapterAsyncCallContext    = 2
	adapterCallContextDisposed = 3
)

/* ============================================================== */

type adapterCallRequest struct {
	ctx      *executionContext
	fn       AdapterCallFunc
	executor AdapterExecutor
	mode     uint8

	flags    AsyncCallFlags
	nestedFn CreateFactoryFunc
	cancel   *syncrun.ChainedCancel
}

func (c *adapterCallRequest) discard() {
	c.mode = adapterCallContextDisposed
}

func (c *adapterCallRequest) ensureMode(mode uint8) {
	if c.mode != mode {
		panic("illegal state")
	}
}

func (c *adapterCallRequest) WithCancel(fn *context.CancelFunc) AsyncCallRequester {
	if c.cancel != nil {
		*fn = c.cancel.Cancel
		return c
	}

	r := *c
	r.cancel = syncrun.NewChainedCancel()
	*fn = r.cancel.Cancel
	return &r
}

func (c *adapterCallRequest) WithNested(nestedFn CreateFactoryFunc) AsyncCallRequester {
	r := *c
	r.nestedFn = nestedFn
	return &r
}

func (c *adapterCallRequest) WithFlags(flags AsyncCallFlags) AsyncCallRequester {
	r := *c
	r.flags = flags
	return &r
}

func (c *adapterCallRequest) Start() {
	c.ensureMode(adapterAsyncCallContext)
	defer c.discard()

	c._startAsync()
}

func (c *adapterCallRequest) DelayedStart() CallConditionalBuilder {
	c.ensureMode(adapterAsyncCallContext)
	defer c.discard()

	return callConditionalBuilder{c.ctx, c._startAsync}
}

func (c *adapterCallRequest) TryCall() bool {
	c.ensureMode(adapterSyncCallContext)
	defer c.discard()

	return c._startSync(true)
}

func (c *adapterCallRequest) Call() {
	c.ensureMode(adapterSyncCallContext)
	defer c.discard()

	if !c._startSync(false) {
		panic("call was cancelled")
	}
}

func (c *adapterCallRequest) _startAsync() {
	if c.cancel != nil && c.cancel.IsCancelled() {
		return
	}

	c.ctx.countAsyncCalls++
	if c.ctx.countAsyncCalls == 0 {
		panic("overflow")
	}

	stepLink := c.ctx.s.NewStepLink()
	callback := NewAdapterCallback(stepLink, nil, c.flags, c.nestedFn)
	cancelFn := c.executor.StartCall(c.fn, callback, c.cancel != nil)

	if c.cancel != nil {
		c.cancel.SetChain(cancelFn)
	}
}

func (c *adapterCallRequest) _startSync(isTry bool) bool {
	resultFn := c._startSyncWithResult(isTry)

	if resultFn == nil {
		return false
	}

	rc := asyncResultContext{slot: c.ctx.s}
	rc.executeResult(resultFn)
	return true
}

func (c *adapterCallRequest) _startSyncWithResult(isTry bool) AsyncResultFunc {

	if ok, result := c.executor.TrySyncCall(c.fn); ok {
		return result
	}

	workerMark := c.ctx.w.GetSignalMark()
	if workerMark == nil {
		return nil
	}

	type resultType struct {
		fn  AsyncResultFunc
		err error
	}
	resultCh := make(chan resultType, 1)

	callback := NewAdapterCallback(c.ctx.s.NewStepLink(), func(fn AsyncResultFunc, err error) {
		resultCh <- resultType{fn, err}
		close(resultCh) // prevent repeated callbacks
	}, 0, c.nestedFn)

	cancelFn := c.executor.StartCall(c.fn, callback, false)

	select {
	case result := <-resultCh:
		if result.err != nil {
			panic(result.err)
		}
		return result.fn

	case <-workerMark.ChannelIf(c.ctx.flags&IgnoreSignal == 0, nil):
		if cancelFn != nil {
			cancelFn()
		}
		if isTry {
			return nil
		}
		panic("signal")

	case <-c.ctx.s.machine.GetStoppingSignal():
		if cancelFn != nil {
			cancelFn()
		}
		return nil
	}
}

/* ============================================================== */

type adapterNotifyRequest struct {
	ctx      *executionContext
	fn       AdapterNotifyFunc
	executor AdapterExecutor
	mode     uint8
}

func (c *adapterNotifyRequest) discard() {
	c.mode = adapterCallContextDisposed
}

func (c *adapterNotifyRequest) ensure() {
	if c.mode != adapterAsyncCallContext {
		panic("illegal state")
	}
}

func (c *adapterNotifyRequest) Send() {
	c.ensure()
	defer c.discard()

	c._startAsync()
}

func (c *adapterNotifyRequest) DelayedSend() CallConditionalBuilder {
	c.ensure()
	return callConditionalBuilder{c.ctx, c._startAsync}
}

func (c *adapterNotifyRequest) _startAsync() {
	c.executor.SendNotify(c.fn)
}

/* ============================================================== */

var _ CallConditionalBuilder = callConditionalBuilder{}

type callConditionalBuilder struct {
	ctx     *executionContext
	kickOff StepPrepareFunc
}

func (v callConditionalBuilder) newConditionalUpdate(updType stateUpdKind) ConditionalBuilder {
	cu := v.ctx.newConditionalUpdate(updType)
	cu.kickOff = v.kickOff
	return &cu
}

func (v callConditionalBuilder) Sleep() ConditionalBuilder {
	return v.newConditionalUpdate(stateUpdSleep)
}

func (v callConditionalBuilder) Poll() ConditionalBuilder {
	return v.newConditionalUpdate(stateUpdPoll)
}

func (v callConditionalBuilder) WaitAny() ConditionalBuilder {
	return v.newConditionalUpdate(stateUpdWaitForEvent)
}

func (v callConditionalBuilder) ThenJump(fn StateFunc) StateUpdate {
	return v.WaitAny().ThenJump(fn)
}

func (v callConditionalBuilder) ThenJumpExt(step SlotStep) StateUpdate {
	return v.WaitAny().ThenJumpExt(step)
}

func (v callConditionalBuilder) ThenRepeat() StateUpdate {
	return v.WaitAny().ThenRepeat()
}
