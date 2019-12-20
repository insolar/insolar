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
	"fmt"
	"sync/atomic"

	"github.com/insolar/insolar/network/consensus/common/syncrun"
)

type AdapterCallbackFunc func(AsyncResultFunc, error) bool

func NewAdapterCallback(adapterId AdapterId, caller StepLink, callbackOverride AdapterCallbackFunc, flags AsyncCallFlags,
	nestedFn CreateFactoryFunc) *AdapterCallback {
	return &AdapterCallback{adapterId, caller, callbackOverride, nil, nestedFn, 0, flags}
}

type AdapterCallback struct {
	adapterId  AdapterId
	caller     StepLink
	callbackFn AdapterCallbackFunc
	cancel     *syncrun.ChainedCancel
	nestedFn   CreateFactoryFunc
	state      uint32 // atomic
	flags      AsyncCallFlags
}

func (c *AdapterCallback) Prepare(needCancel bool) context.CancelFunc {
	if !atomic.CompareAndSwapUint32(&c.state, 0, 1) {
		panic("illegal state - in use")
	}
	if c.cancel != nil {
		panic("illegal state")
	}
	if !needCancel {
		return nil
	}

	c.cancel = syncrun.NewChainedCancel()
	return c.cancel.Cancel
}

const stepBondTolerance uint32 = 1

func (c *AdapterCallback) canCall() bool {
	return c.flags&CallBoundToStep == 0 || c.caller.IsNearStep(stepBondTolerance)
}

func (c *AdapterCallback) IsCancelled() bool {
	return c.cancel.IsCancelled() || !c.canCall()
}

func (c *AdapterCallback) ChainedCancel() *syncrun.ChainedCancel {
	return c.cancel
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
		if c.callbackFn(resultFn, err) {
			return
		}
	}

	c.caller.s.machine.queueAsyncCallback(c.caller.SlotLink, func(slot *Slot, worker DetachableSlotWorker, err error) StateUpdate {
		slot.decAsyncCount()

		wakeupAllowed := c.flags&WakeUpBoundToStep == 0 || c.caller.IsNearStep(stepBondTolerance)
		wakeup := false
		switch {
		case err != nil:
			return StateUpdate{} // result will be replaced by queueAsyncCallback()

		case resultFn == nil /* cancelled by adapter */ || c.cancel.IsCancelled():
			wakeup = c.flags&WakeUpOnCancel != 0

		default:
			rc := asyncResultContext{s: slot}
			wakeupResult := rc.executeResult(resultFn)
			wakeup = wakeupResult || c.flags&WakeUpOnResult != 0
		}

		if wakeup && wakeupAllowed {
			return newStateUpdateTemplate(updCtxAsyncCallback, 0, stateUpdWakeup).newNoArg()
		}
		return newStateUpdateTemplate(updCtxAsyncCallback, 0, stateUpdNoChange).newNoArg()
	}, err)
}

func (c *AdapterCallback) SendNested(defaultFactoryFn CreateFactoryFunc, payload interface{}) error {

	m := c.caller.getActiveMachine()
	if m == nil {
		return fmt.Errorf("target SlotMachine is stopping/stopped")
	}

	createFn := func(factoryFn CreateFactoryFunc) (bool, error) {
		if factoryFn == nil {
			return false, nil
		}
		if cf := factoryFn(payload); cf != nil {
			switch link, ok := m.AddNested(c.adapterId, c.caller.SlotLink, cf); {
			case ok:
				return true, nil
			case link.IsEmpty():
				return true, fmt.Errorf("target SlotMachine is stopping/stopped")
			default:
				return true, fmt.Errorf("cancelled by constructor")
			}
		}
		return false, nil
	}

	if ok, err := createFn(c.nestedFn); ok {
		return err
	}
	if ok, err := createFn(defaultFactoryFn); ok {
		return err
	}
	return fmt.Errorf("unknown payload for nested call")
}

type NestedCallFunc func(interface{}) error

func (c *AdapterCallback) getNestedCallHandler(defFactoryFn CreateFactoryFunc) NestedCallFunc {
	if defFactoryFn == nil && c.nestedFn == nil {
		return nil
	}
	return func(v interface{}) error {
		return c.SendNested(defFactoryFn, v)
	}
}
