//
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
//

package watchdog

import (
	"context"
	"time"
)

func DoneOf(ctx context.Context) <-chan struct{} {
	if h := from(ctx); h != nil {
		return h.done(ctx)
	}
	return ctx.Done()
}

func Beat(ctx context.Context) {
	ForcedBeat(ctx, false)
}

func ForcedBeat(ctx context.Context, forced bool) {
	if h := from(ctx); h != nil {
		h.beat(ctx, forced)
	}
}

func WithFrame(ctx context.Context, frameName string) context.Context {
	if h := from(ctx); h != nil {
		return h.root.createSubFrame(ctx, frameName, h)
	}
	return ctx
}

func Call(ctx context.Context, frameName string, fn func(context.Context)) {

	frame := from(ctx)
	if frame == nil {
		fn(ctx)
		return
	}

	frame.root.createSubFrame(ctx, frameName, frame).call(fn)
}

func WithFactory(ctx context.Context, name string, factory HeartbeatGeneratorFactory) context.Context {
	r := frameRoot{factory}
	return r.createSubFrame(ctx, name, nil)
}

func WithoutFactory(ctx context.Context) context.Context {
	return context.WithValue(ctx, watchdogKey, nil) // stop search
}

func FromContext(ctx context.Context) (bool, HeartbeatGeneratorFactory) {
	f := from(ctx)
	if f == nil {
		return false, nil
	}
	if f.root == nil {
		return true, nil
	}
	return true, f.root.factory
}

func from(ctx context.Context) *frame {
	h, ok := ctx.Value(watchdogKey).(*frame)
	if ok {
		return h
	}
	return nil
}

var watchdogKey = &struct{}{}

type frameRoot struct {
	factory HeartbeatGeneratorFactory
}

func (r *frameRoot) createSubFrame(ctx context.Context, name string, parent *frame) *frame {
	if parent != nil {
		name = parent.name + "/" + name
	}
	return &frame{r, ctx, r.factory.CreateGenerator(name), name}
}

type frame struct {
	root      *frameRoot
	context   context.Context
	generator *HeartbeatGenerator
	name      string
}

func (h *frame) Deadline() (deadline time.Time, ok bool) {
	h.beat(h.context, false)
	return h.context.Deadline()
}

func (h *frame) Value(key interface{}) interface{} {
	if watchdogKey == key {
		return h
	}
	h.beat(h.context, false)
	return h.context.Value(key)
}

func (h *frame) Err() error {
	err := h.context.Err()
	if err != nil {
		h.generator.Cancel()
	} else {
		h.generator.Heartbeat()
	}
	return err
}

func (h *frame) Done() <-chan struct{} {
	return h.done(h.context)
}

func (h *frame) beat(ctx context.Context, forced bool) {
	if ctx.Err() != nil {
		h.generator.Cancel()
	} else {
		h.generator.ForcedHeartbeat(forced)
	}
}

func (h *frame) start() {
	h.generator.ForcedHeartbeat(true)
}

func (h *frame) cancel() {
	h.generator.Cancel()
}

func (h *frame) done(ctx context.Context) <-chan struct{} {
	ch := ctx.Done()
	select {
	case <-ch:
		h.generator.Cancel()
	default:
		h.generator.ForcedHeartbeat(false)
	}
	return ch
}

func (h *frame) call(fn func(context.Context)) {
	h.start()
	defer h.cancel()
	fn(h)
}
