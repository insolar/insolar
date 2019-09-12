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

var _ ExecutionAdapter = &adapterExecHelper{}

type adapterExecHelper struct {
	adapterID AdapterID
	executor  AdapterExecutor
}

func (p *adapterExecHelper) IsEmpty() bool {
	return p.adapterID.IsEmpty()
}

func (p *adapterExecHelper) GetAdapterID() AdapterID {
	return p.adapterID
}

func (p *adapterExecHelper) PrepareSync(ctx ExecutionContext, fn AdapterCallFunc) SyncCallRequester {
	return &adapterCallRequest{ctx: ctx.(*executionContext), fn: fn, executor: p.executor, mode: adapterSyncCallContext}
}

func (p *adapterExecHelper) PrepareAsync(ctx ExecutionContext, fn AdapterCallFunc) AsyncCallRequester {
	return &adapterCallRequest{ctx: ctx.(*executionContext), fn: fn, executor: p.executor, mode: adapterAsyncCallContext}
}

const (
	adapterSyncCallContext     = 1
	adapterAsyncCallContext    = 2
	adapterCallContextDisposed = 3
)

type adapterCallRequest struct {
	ctx      *executionContext
	fn       AdapterCallFunc
	executor AdapterExecutor
	mode     uint8

	stepBound bool
	cancel    *syncrun.ChainedCancel
}

func (c *adapterCallRequest) discard() {
	c.mode = adapterCallContextDisposed
}

func (c *adapterCallRequest) ensureMode(mode uint8) {
	if c.mode != mode {
		panic("illegal state")
	}
}

func (c *adapterCallRequest) GetCancel(fn *context.CancelFunc) AsyncCallRequester {
	if c.cancel != nil {
		*fn = c.cancel.Cancel
		return c
	}

	r := *c
	r.cancel = syncrun.NewChainedCancel()
	*fn = r.cancel.Cancel
	return &r
}

func (c *adapterCallRequest) CancelOnStep(attach bool) AsyncCallRequester {
	r := *c
	r.stepBound = attach
	return &r
}

func (c *adapterCallRequest) Start() {
	c.ensureMode(adapterAsyncCallContext)
	defer c.discard()

	c._startAsync()
}

func (c *adapterCallRequest) DelayedStart() CallConditionalUpdate {
	c.ensureMode(adapterAsyncCallContext)
	defer c.discard()

	return &conditionalUpdate{marker: c.ctx.getMarker(), updMode: stateUpdNext, kickOff: c._startAsync}
}

func (c *adapterCallRequest) TryCall() bool {
	c.ensureMode(adapterSyncCallContext)
	defer c.discard()

	return c._startSync()
}

func (c *adapterCallRequest) Call() {
	c.ensureMode(adapterSyncCallContext)
	defer c.discard()

	if !c._startSync() {
		panic("call was cancelled")
	}
}

func (c *adapterCallRequest) _startAsync() {
	var stepLink StepLink
	stepLink = c.ctx.s.NewStepLink()
	if !c.stepBound {
		stepLink = stepLink.AnyStep()
	}

	if c.cancel != nil && c.cancel.IsCancelled() {
		return
	}

	c.ctx.countAsyncCalls++

	cancelFn := c.executor.StartCall(stepLink, c.fn, func(fn AsyncResultFunc, recovered interface{}) {
		c.ctx.s.machine.applyAsyncStateUpdate(stepLink.SlotLink, fn, recovered)
	}, c.cancel != nil)

	if c.cancel != nil {
		c.cancel.SetChain(cancelFn)
	}
}

func (c *adapterCallRequest) _startSync() bool {
	resultFn := c._startSyncWithResult()

	if resultFn == nil {
		return false
	}

	rc := asyncResultContext{slot: c.ctx.s}
	rc.executeResult(resultFn)
	return true
}

func (c *adapterCallRequest) _startSyncWithResult() AsyncResultFunc {

	if ok, result := c.executor.TrySyncCall(c.fn); ok {
		return result
	}

	ok, wc := c.ctx.worker.GetCond()
	if !ok {
		return nil
	}

	var resultFn AsyncResultFunc
	var resultRecovered interface{}
	var callState int

	stepLink := c.ctx.s.NewStepLink()
	cancelFn := c.executor.StartCall(stepLink, c.fn, func(fn AsyncResultFunc, recovered interface{}) {
		wc.L.Lock()
		if callState == 0 {
			resultFn = fn
			resultRecovered = recovered
			callState = 1
			wc.Broadcast()
		}
		wc.L.Unlock()
	}, false)

	wc.L.Lock()
	if callState == 0 {
		wc.Wait()

		if callState == 0 {
			/* Cond can be triggered by Worker for emergency stop */
			callState = 2

			wc.L.Unlock()

			if cancelFn != nil {
				cancelFn()
			}
			return nil
		}
	}
	wc.L.Unlock()

	if resultRecovered != nil {
		panic(resultRecovered)
	}
	return resultFn
}

func NewAdapterCallback(stepLink StepLink, callback AdapterCallbackFunc, cancel *syncrun.ChainedCancel) AdapterCallback {
	return AdapterCallback{stepLink, callback, cancel}
}

type AdapterCallback struct {
	stepLink StepLink
	callback AdapterCallbackFunc
	cancel   *syncrun.ChainedCancel
}

func (c AdapterCallback) IsZero() bool {
	return c.stepLink.IsEmpty()
}

func (c AdapterCallback) IsCancelled() bool {
	return !c.stepLink.IsAtStep() || c.cancel != nil && c.cancel.IsCancelled()
}

func (c AdapterCallback) SendResult(result AsyncResultFunc) {
	if c.IsZero() {
		panic("illegal state")
	}
	_sendResult(c.stepLink, result, c.callback, c.cancel)
}

// just to make sure that outer struct doesn't leak into a closure
func _sendResult(stepLink StepLink, result AsyncResultFunc, callback AdapterCallbackFunc, cancel *syncrun.ChainedCancel) {

	if result == nil {
		// NB! Do NOT ignore "result = nil" - it MUST decrement async call count
		callback(func(ctx AsyncResultContext) {}, nil)
		return
	}

	callback(func(ctx AsyncResultContext) {
		if result == nil || !stepLink.IsAtStep() || cancel != nil && cancel.IsCancelled() {
			return
		}
		result(ctx)
	}, nil)
}

func (c AdapterCallback) SendPanic(recovered interface{}) {
	if c.IsZero() {
		panic("illegal state")
	}
	c.callback(nil, recovered)
}

func (c AdapterCallback) SendCancel() {
	if c.IsZero() {
		panic("illegal state")
	}
	c.callback(nil, nil)
}
