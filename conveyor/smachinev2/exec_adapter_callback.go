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

import (
	"context"
	"sync/atomic"

	"github.com/insolar/insolar/network/consensus/common/syncrun"
)

func NewAdapterCallback(caller StepLink, callback AdapterCallbackFunc, nestedFn CreateFactoryFunc) *AdapterCallback {
	return &AdapterCallback{caller, callback, nil, nestedFn, 0}
}

type AdapterCallback struct {
	caller     StepLink
	callbackFn AdapterCallbackFunc
	cancel     *syncrun.ChainedCancel
	nestedFn   CreateFactoryFunc
	state      uint32 // atomic
}

func (c *AdapterCallback) Prepare(requireCancel bool, chainedCancel context.CancelFunc) context.CancelFunc {
	if !atomic.CompareAndSwapUint32(&c.state, 0, 1) {
		panic("illegal state - in use")
	}
	if c.cancel != nil {
		panic("illegal state")
	}
	if chainedCancel == nil && !requireCancel {
		return nil
	}

	c.cancel = syncrun.NewChainedCancel()
	c.cancel.SetChain(chainedCancel)
	return c.cancel.Cancel
}

func (c *AdapterCallback) IsCancelled() bool {
	return !c.caller.IsAtStep() || c.cancel.IsCancelled()
}

func (c *AdapterCallback) SendResult(result AsyncResultFunc) {
	if result == nil {
		// NB! Do NOT ignore "result = nil" - callback need to decrement async call count
		result = func(ctx AsyncResultContext) {}
	}
	c.callback(false, result, nil)
}

func (c *AdapterCallback) SendPanic(err error) {
	if err == nil {
		panic("illegal value")
	}
	c.callback(false, nil, err)
}

// can be called repeatedly
func (c *AdapterCallback) SendCancel() {
	c.cancel.Cancel()
	c.callback(true, nil, nil)
}

func (c *AdapterCallback) callback(isCancel bool, resultFn AsyncResultFunc, err error) {
	switch atomic.SwapUint32(&c.state, 2) {
	case 2:
		if isCancel {
			// repeated cancel are safe
			return
		}
		panic("illegal state - repeated callback")
	case 0:
		panic("illegal state - unprepared")
	}

	switch {
	case c.callbackFn != nil:
		c.callbackFn(resultFn, err)
		return
	case !c.caller.IsAtStep():
		return
	}
	c.caller.s.machine.queueAsyncCallback(c.caller.SlotLink, func(slot *Slot, worker DetachableSlotWorker) StateUpdate {
		slot.decAsyncCount()

		if err == nil && resultFn != nil && !c.cancel.IsCancelled() {
			rc := asyncResultContext{slot: slot}
			if wakeup := rc.executeResult(resultFn); wakeup {
				return newStateUpdateTemplate(updCtxAsyncCallback, 0, stateUpdRepeat).newUint(0)
			}
		}

		return newStateUpdateTemplate(updCtxAsyncCallback, 0, stateUpdNoChange).newNoArg()
	}, err)
}

func (c *AdapterCallback) SendNested(defaultFactoryFn CreateFactoryFunc, payload interface{}) bool {
	if c.caller.getActiveMachine() == nil {
		return false
	}
	return c._createNested(c.nestedFn, payload) || c._createNested(defaultFactoryFn, payload)
}

func (c *AdapterCallback) _createNested(factoryFn CreateFactoryFunc, payload interface{}) bool {
	if factoryFn == nil {
		return false
	}
	createFn := factoryFn(payload)
	if createFn == nil {
		return false
	}
	_, ok := c.caller.s.machine.createNestedForAdapter(c.caller.SlotLink, createFn)
	return ok
}
