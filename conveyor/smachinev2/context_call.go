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
	"sync/atomic"
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

func (c *adapterCallRequest) WithAutoCancelOnStep(attach bool) AsyncCallRequester {
	r := *c
	r.stepBound = attach
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

	cu := c.ctx.newConditionalUpdate(stateUpdWaitForEvent)
	cu.kickOff = c._startAsync
	return &cu
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
	if c.ctx.countAsyncCalls == 0 {
		panic("overflow")
	}
	cancelFn := c.executor.StartCall(stepLink, c.fn, _asyncCallback(stepLink), c.cancel != nil)

	if c.cancel != nil {
		c.cancel.SetChain(cancelFn)
	}
}

func _asyncCallback(stepLink StepLink) AdapterCallbackFunc {
	callbackGuard := uint32(0)

	return func(resultFn AsyncResultFunc, err error) {
		if atomic.SwapUint32(&callbackGuard, 1) != 0 {
			panic("repeated callback")
		}

		if !stepLink.IsValid() {
			return
		}

		stepLink.s.machine.queueAsyncCallback(stepLink.SlotLink, func(slot *Slot, worker DetachableSlotWorker) StateUpdate {
			slot.decAsyncCount()

			if err == nil && resultFn != nil {
				rc := asyncResultContext{slot: slot}
				if wakeup := rc.executeResult(resultFn); wakeup {
					return newStateUpdateTemplate(updCtxAsyncCallback, 0, stateUpdRepeat).newUint(0)
				}
			}

			return newStateUpdateTemplate(updCtxAsyncCallback, 0, stateUpdNoChange).newNoArg()
		}, err)
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

	ok, wc := c.ctx.w.GetCond()
	if !ok {
		return nil
	}

	var resultFn AsyncResultFunc
	var resultErr error
	var callState int

	stepLink := c.ctx.s.NewStepLink()
	cancelFn := c.executor.StartCall(stepLink, c.fn, func(fn AsyncResultFunc, err error) {
		wc.L.Lock()
		switch callState {
		case 0:
			resultFn = fn
			resultErr = err
			callState = 1
			wc.Broadcast()
		case 1:
			wc.L.Unlock()
			panic("repeated callback")
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

	if resultErr != nil {
		panic(resultErr)
	}
	return resultFn
}
