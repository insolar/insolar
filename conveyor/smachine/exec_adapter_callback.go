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

type AdapterCallbackFunc func(AsyncResultFunc, error)

func NewAdapterCallback(caller StepLink, callback AdapterCallbackFunc, flags AsyncCallFlags, nestedFn CreateFactoryFunc) *AdapterCallback {
	return &AdapterCallback{caller, callback, nil, nestedFn, 0, flags}
}

type AdapterCallback struct {
	caller     StepLink
	callbackFn AdapterCallbackFunc
	cancel     *syncrun.ChainedCancel
	nestedFn   CreateFactoryFunc
	state      uint32 // atomic
	flags      AsyncCallFlags
}

func (c *AdapterCallback) Prepare(needCancel bool) (context.CancelFunc, func(context.CancelFunc)) {
	if !atomic.CompareAndSwapUint32(&c.state, 0, 1) {
		panic("illegal state - in use")
	}
	if c.cancel != nil {
		panic("illegal state")
	}
	if !needCancel {
		return nil, nil
	}

	c.cancel = syncrun.NewChainedCancel()
	return c.cancel.Cancel, c.cancel.SetChain
}

const stepBondTolerance = 1

func (c *AdapterCallback) canCall() bool {
	return c.flags&CallBoundToStep == 0 || c.caller.IsNearStep(stepBondTolerance)
}

func (c *AdapterCallback) IsCancelled() bool {
	return c.cancel.IsCancelled() || !c.canCall()
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
	case !c.canCall():
		return
	case c.callbackFn != nil:
		c.callbackFn(resultFn, err)
		return
	}

	c.caller.s.machine.queueAsyncCallback(c.caller.SlotLink, func(slot *Slot, worker DetachableSlotWorker) StateUpdate {
		slot.decAsyncCount()

		if err == nil && resultFn != nil /* not a cancellation callback */ && !c.cancel.IsCancelled() {
			rc := asyncResultContext{slot: slot}
			wakeup := rc.executeResult(resultFn)

			if (wakeup || c.flags&AutoWakeUp != 0) && (c.flags&WakeUpBoundToStep == 0 || c.caller.IsNearStep(stepBondTolerance)) {
				return newStateUpdateTemplate(updCtxAsyncCallback, 0, stateUpdRepeat).newUint(0)
			}
		}

		return newStateUpdateTemplate(updCtxAsyncCallback, 0, stateUpdNoChange).newNoArg()
	}, err)
}

func (c *AdapterCallback) SendNested(adapterId AdapterId, defaultFactoryFn CreateFactoryFunc, payload interface{}) bool {
	if c.caller.getActiveMachine() == nil {
		return false
	}
	return c._createNested(adapterId, c.nestedFn, payload) ||
		c._createNested(adapterId, defaultFactoryFn, payload)
}

func (c *AdapterCallback) _createNested(adapterId AdapterId, factoryFn CreateFactoryFunc, payload interface{}) bool {
	if factoryFn == nil {
		return false
	}
	createFn := factoryFn(payload)
	if createFn == nil {
		return false
	}
	_, ok := c.caller.s.machine.AddNested(adapterId, c.caller.SlotLink, createFn)
	return ok
}
