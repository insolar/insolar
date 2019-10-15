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
	"github.com/insolar/insolar/network/consensus/common/syncrun"
)

func NewAdapterCallback(stepLink StepLink, callback AdapterCallbackFunc, nested AdapterNestedEventFunc,
	cancel *syncrun.ChainedCancel) AdapterCallback {
	return AdapterCallback{stepLink, callback, nested, cancel}
}

type AdapterCallback struct {
	stepLink StepLink
	callback AdapterCallbackFunc
	nested   AdapterNestedEventFunc
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
	// c.stepLink must be checked by callback
	_sendResult(result, c.callback, c.cancel)
}

// just to make sure that outer struct doesn't leak into a closure
func _sendResult(result AsyncResultFunc, callback AdapterCallbackFunc, cancel *syncrun.ChainedCancel) {

	if callback == nil {
		if result == nil {
			return
		}
		panic("illegal state")
	}

	if result == nil {
		// NB! Do NOT ignore "result = nil" - it callback may need to decrement async call count
		callback(func(ctx AsyncResultContext) {}, nil)
		return
	}

	callback(func(ctx AsyncResultContext) {
		if result == nil || cancel != nil && cancel.IsCancelled() {
			return
		}
		result(ctx)
	}, nil)
}

func (c AdapterCallback) SendPanic(err error) {
	if c.IsZero() {
		panic("illegal state")
	}
	c.callback(nil, err)
}

func (c AdapterCallback) SendCancel() {
	if c.IsZero() {
		panic("illegal state")
	}
	c.callback(nil, nil)
}

func (c AdapterCallback) SendNested(payload interface{}) {
	if c.IsZero() {
		panic("illegal state")
	}
	c.nested(c.stepLink, payload, false)
}

func (c AdapterCallback) SendNestedWithCancel(payload interface{}) context.CancelFunc {
	if c.IsZero() {
		panic("illegal state")
	}
	return c.nested(c.stepLink, payload, true)
}
