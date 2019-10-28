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
	"errors"

	"github.com/insolar/insolar/network/consensus/common/syncrun"
)

type AdapterCall struct {
	CallFn   AdapterCallFunc
	Callback *AdapterCallback
}

var ErrCancelledCall = errors.New("cancelled")

type AdapterCallDelegateFunc func(
	// Not nil
	callFn AdapterCallFunc,
	// Nil for notify calls and when there is no nested call factories are available.
	// Returns false when nested call is impossible (outer call is cancelled or finished)
	nestedCallFn NestedCallFunc,
	// Nil when cancellation is not traced / not configured on SlotMachine adapter / and for notifications
	chainCancel *syncrun.ChainedCancel) (AsyncResultFunc, error)

func (c AdapterCall) DelegateAndSendResult(defaultNestedFn CreateFactoryFunc, delegate AdapterCallDelegateFunc) error {
	switch {
	case delegate == nil:
		panic("illegal value")
	case c.Callback == nil:
		return c.delegateNotify(delegate)
	case c.Callback.IsCancelled():
		c.Callback.SendCancel()
		return nil
	}

	result, err := func() (result AsyncResultFunc, err error) {
		defer func() {
			err = RecoverAsyncSlotPanicWithStack("async call", recover(), err)
		}()
		nestedCallFn := c.Callback.getNestedCallHandler(defaultNestedFn)
		return delegate(c.CallFn, nestedCallFn, c.Callback.ChainedCancel())
	}()

	switch {
	case err == nil:
		if !c.Callback.IsCancelled() {
			c.Callback.SendResult(result)
		}
		fallthrough
	case err == ErrCancelledCall:
		c.Callback.SendCancel()
	default:
		c.Callback.SendPanic(err)
	}
	return nil
}

func (c AdapterCall) delegateNotify(delegate AdapterCallDelegateFunc) error {
	result, err := func() (result AsyncResultFunc, err error) {
		defer func() {
			err = RecoverAsyncSlotPanicWithStack("async notify", recover(), err)
		}()
		return delegate(c.CallFn, nil, nil)
	}()
	switch {
	case err == nil:
		if result == nil {
			return nil
		}
		return errors.New("result is unexpected")
	case err == ErrCancelledCall:
		// can't send cancel
		return nil
	default:
		return err
	}
}

func (c AdapterCall) RunAndSendResult(arg interface{}) error {
	return c.DelegateAndSendResult(nil,
		func(callFn AdapterCallFunc, _ NestedCallFunc, _ *syncrun.ChainedCancel) (AsyncResultFunc, error) {
			return callFn(arg), nil
		})
}
