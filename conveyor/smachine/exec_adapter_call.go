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

type AdapterCallDelegateFunc func(callFn AdapterCallFunc, notify bool, chainCancel *syncrun.ChainedCancel) (AsyncResultFunc, error)

func (c AdapterCall) DelegateAndSendResult(delegate AdapterCallDelegateFunc) error {
	if delegate == nil {
		panic("illegal value")
	}

	if c.Callback == nil {
		result, err := func() (result AsyncResultFunc, err error) {
			defer func() {
				err = RecoverAsyncSlotPanicWithStack("async call", recover(), err)
			}()
			return delegate(c.CallFn, true, nil)
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

	if c.Callback.IsCancelled() {
		c.Callback.SendCancel()
		return nil
	}

	result, err := func() (result AsyncResultFunc, err error) {
		defer func() {
			err = RecoverAsyncSlotPanicWithStack("async call", recover(), err)
		}()

		return delegate(c.CallFn, false, c.Callback.ChainedCancel())
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

func (c AdapterCall) RunAndSendResult(arg interface{}) error {
	return c.DelegateAndSendResult(func(callFn AdapterCallFunc, _ bool, _ *syncrun.ChainedCancel) (AsyncResultFunc, error) {
		return callFn(arg), nil
	})
}
