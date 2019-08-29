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
	"sync"
)

var _ ExecutionAdapterSink = &ChannelAdapter{}

func NewChannelAdapter(ctx context.Context, chanLen int) ChannelAdapter {
	return ChannelAdapter{
		ctx: ctx,
		c:   make(chan ChannelRecord, chanLen),
	}
}

type ChannelAdapter struct {
	ctx context.Context
	c   chan ChannelRecord
	m   sync.Mutex
	q   []ChannelRecord
}

func (c *ChannelAdapter) Channel() <-chan ChannelRecord {
	return c.c
}

func (c *ChannelAdapter) Close() {
	defer func() {
		_ = recover()
	}()

	c.m.Lock()
	defer c.m.Unlock()
	c.q = nil
	close(c.c)
}

func (c *ChannelAdapter) CallAsync(stepLink StepLink, fn AdapterCallFunc, callback AdapterCallbackFunc) {
	r := ChannelRecord{stepLink, fn, callback, nil}

	if !c.append(r, false) || !c.send(r) {
		c.append(r, true)
	}
}

func (c *ChannelAdapter) CallAsyncWithCancel(stepLink StepLink, fn AdapterCallFunc, callback AdapterCallbackFunc) (cancelFn context.CancelFunc) {
	cancel := syncrun.NewChainedCancel()
	r := ChannelRecord{stepLink, fn, callback, cancel}

	if !c.append(r, false) || !c.send(r) {
		c.append(r, true)
	}
	return cancel.Cancel
}

func (c *ChannelAdapter) append(r ChannelRecord, force bool) bool {
	c.m.Lock()
	defer c.m.Unlock()

	switch {
	case len(c.q) > 0:
		break
	case !force:
		return false
	default:
		go c.sendWorker() // wont start because of lock
	}
	c.q = append(c.q, r)
	return true
}

func (c *ChannelAdapter) send(r ChannelRecord) bool {
	select {
	case c.c <- r:
		return true
	default:
		return false
	}
}

func (c *ChannelAdapter) sendWorker() {

	var done <-chan struct{}
	if c.ctx != nil {
		done = c.ctx.Done()
	}

	defer func() {
		_ = recover()
	}()
	for {
		var r ChannelRecord
		c.m.Lock()
		switch len(c.q) {
		case 0:
			return
		case 1:
			r = c.q[0]
			c.q = nil
		default:
			r, c.q[0] = c.q[0], r
			c.q = c.q[1:] // TODO potential memory leak on same speed of read & write
		}
		c.m.Unlock()

		select {
		case <-done:
			return
		case c.c <- r:
		}
	}
}

type ChannelRecord struct {
	stepLink StepLink
	callFunc AdapterCallFunc
	callback AdapterCallbackFunc
	cancel   *syncrun.ChainedCancel
}

func (c ChannelRecord) IsCancelled() bool {
	return !c.stepLink.IsAtStep() || c.cancel != nil && c.cancel.IsCancelled()
}

func (c ChannelRecord) RunCall() AsyncResultFunc {
	return c.callFunc()
}

func (c ChannelRecord) SendResult(result AsyncResultFunc) {
	if result == nil {
		c.callback(func(ctx AsyncResultContext) {
		})
	} else {
		c.callback(result)
	}
}

func (c ChannelRecord) SendCancel() {
	c.callback(nil)
}

func (c ChannelRecord) RunAndSendResult() bool {
	if c.IsCancelled() {
		c.callback(nil)
		return false
	}

	result, recovered := c.safeCall()

	if c.IsCancelled() && recovered == nil {
		c.callback(nil)
		return false
	}

	wrapCallback(result, recovered, c.stepLink, c.cancel, c.callback)
	return true
}

// just to make sure that ChannelRecord doesn't leak into a closure
func wrapCallback(result AsyncResultFunc, recovered interface{}, stepLink StepLink, cancel *syncrun.ChainedCancel, callback AdapterCallbackFunc) {

	if recovered != nil {
		result = nil
	}

	callback(func(ctx AsyncResultContext) {
		if recovered != nil {
			panic(recovered)
		}
		if result == nil || !stepLink.IsAtStep() {
			return
		}
		result(ctx)
	})
}

func (c ChannelRecord) safeCall() (result AsyncResultFunc, recovered interface{}) {
	defer func() {
		recovered = recover()
	}()
	return c.callFunc(), nil
}
