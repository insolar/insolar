///
// Modified BSD 3-Clause Clear License
//
// Copyright (c) 2019 Insolar Technologies GmbH
//
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without modification,
// are permitted (subject to the limitations in the disclaimer below) provided that
// the following conditions are met:
//  * Redistributions of source code must retain the above copyright notice, this list
//    of conditions and the following disclaimer.
//  * Redistributions in binary form must reproduce the above copyright notice, this list
//    of conditions and the following disclaimer in the documentation and/or other materials
//    provided with the distribution.
//  * Neither the name of Insolar Technologies GmbH nor the names of its contributors
//    may be used to endorse or promote products derived from this software without
//    specific prior written permission.
//
// NO EXPRESS OR IMPLIED LICENSES TO ANY PARTY'S PATENT RIGHTS ARE GRANTED
// BY THIS LICENSE. THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS
// AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES,
// INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY
// AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL
// THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT,
// INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING,
// BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS
// OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
// ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
//
// Notwithstanding any other provisions of this license, it is prohibited to:
//    (a) use this software,
//
//    (b) prepare modifications and derivative works of this software,
//
//    (c) distribute this software (including without limitation in source code, binary or
//        object code form), and
//
//    (d) reproduce copies of this software
//
//    for any commercial purposes, and/or
//
//    for the purposes of making available this software to third parties as a service,
//    including, without limitation, any software-as-a-service, platform-as-a-service,
//    infrastructure-as-a-service or other similar online service, irrespective of
//    whether it competes with the products or services of Insolar Technologies GmbH.
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
