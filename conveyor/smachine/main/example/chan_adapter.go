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

package example

import (
	"context"
	"sync"

	"github.com/insolar/insolar/conveyor/smachine"
)

var _ smachine.AdapterExecutor = &ChannelAdapter{}

// This is a basic implementation of an adapter to run all calls in a separate goroutine.
// There an important property - it doesn't lock up.
func NewChannelAdapter(ctx context.Context, chanLen int, overflowLimit int) ChannelAdapter {
	return ChannelAdapter{
		ctx: ctx,
		c:   make(chan ChannelRecord, chanLen),
		o:   overflowLimit,
	}
}

type ChannelRecord = smachine.AdapterCall

type ChannelAdapter struct {
	ctx context.Context
	c   chan ChannelRecord
	m   sync.Mutex
	q   []ChannelRecord
	o   int
}

func (c *ChannelAdapter) TrySyncCall(fn smachine.AdapterCallFunc) (bool, smachine.AsyncResultFunc) {
	return false, nil
}

func (c *ChannelAdapter) SendNotify(fn smachine.AdapterNotifyFunc) {
	if fn == nil {
		panic("illegal value")
	}

	r := ChannelRecord{CallFn: func() smachine.AsyncResultFunc {
		fn()
		return nil
	}}
	c.queueCall(r)
}

func (c *ChannelAdapter) StartCall(fn smachine.AdapterCallFunc, callback *smachine.AdapterCallback, needCancel bool) context.CancelFunc {
	if fn == nil {
		panic("illegal value")
	}
	if callback == nil {
		panic("illegal value")
	}

	r := ChannelRecord{CallFn: fn, Callback: callback}
	cancelFn, setChainCancelFn := callback.Prepare(needCancel)

	nativeCancelFn := c.queueCall(r)
	if setChainCancelFn != nil {
		setChainCancelFn(nativeCancelFn)
	}

	return cancelFn
}

func (c *ChannelAdapter) queueCall(r ChannelRecord) context.CancelFunc {
	if !c.append(r, false) && !c.send(r) {
		c.append(r, true)
	}
	return nil
}

func (c *ChannelAdapter) Channel() <-chan ChannelRecord {
	return c.c
}

func (c *ChannelAdapter) Context() context.Context {
	return c.ctx
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

	if c.o >= 0 && len(c.q) > c.o {
		panic("overflow")
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
			c.m.Unlock()
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
