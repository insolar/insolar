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

package syncrun

import (
	"context"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"sync"
	"sync/atomic"
	"time"
)

const (
	runStatusUninitialized = iota
	runStatusInitialized
	runStatusStarted
	runStatusStopping
	runStatusStopped
)

type Status uint8

func (v Status) IsRunning() bool {
	return v == runStatusStarted
}

func (v Status) IsStopping() bool {
	return v == runStatusStopping
}

func (v Status) IsStoppingOrStopped() bool {
	return v >= runStatusStopping
}

func (v Status) WasStarted() bool {
	return v >= runStatusStarted
}

func (v Status) WasInitialized() bool {
	return v >= runStatusInitialized
}

type WorkFunc func(context.Context) error

/* Ensures that all methods are called in the exact sequence */
type SyncingWorker struct {
	runStatus     int32 // atomic
	ctx           context.Context
	cancelFn      func()
	beforeStartFn func(context.Context)
	afterStopFn   func(context.Context, interface{})
	asyncCmd      chan WorkFunc
	timeout       *time.Timer
}

func (p *SyncingWorker) Init(queueLen int, beforeStartFn func(context.Context), afterStopFn func(context.Context, interface{})) {
	if queueLen <= 0 {
		panic("illegal value")
	}

	if !atomic.CompareAndSwapInt32(&p.runStatus, runStatusUninitialized, runStatusInitialized) {
		panic("illegal state")
	}

	p.asyncCmd = make(chan WorkFunc, queueLen)

	if beforeStartFn == nil {
		beforeStartFn = func(context.Context) {}
	}
	if afterStopFn == nil {
		afterStopFn = func(context.Context, interface{}) {}
	}
	p.beforeStartFn = beforeStartFn
	p.afterStopFn = afterStopFn
}

func (p *SyncingWorker) GetStatus() Status {
	return Status(atomic.LoadInt32(&p.runStatus))
}

func (p *SyncingWorker) Stop() {
	if p.cancelFn == nil {
		return
	}
	p.cancelFn()
}

func (p *SyncingWorker) ensureStart() (bool, Status) {
	for {
		if atomic.CompareAndSwapInt32(&p.runStatus, runStatusInitialized, runStatusStarted) {
			break
		}
		s := atomic.LoadInt32(&p.runStatus)
		if s != runStatusInitialized {
			return false, Status(s)
		}
	}
	return true, runStatusStarted
}

func (p *SyncingWorker) ensureStartUnprepared(ctx context.Context) (bool, Status) {
	if ctx == nil {
		panic("illegal value")
	}
	if p.ctx != nil {
		panic("illegal state - was prepared")
	}
	return p.ensureStart()
}

func (p *SyncingWorker) TryStartWithDeadline(ctx context.Context, d time.Time) Status {
	if canStart, status := p.ensureStartUnprepared(ctx); !canStart {
		return status
	}

	p.ctx, p.cancelFn = context.WithDeadline(ctx, d)
	p.run()
	return runStatusStarted
}

func (p *SyncingWorker) TryStart(ctx context.Context) Status {
	if canStart, status := p.ensureStartUnprepared(ctx); !canStart {
		return status
	}

	p.ctx, p.cancelFn = context.WithCancel(ctx)
	p.run()
	return runStatusStarted
}

func (p *SyncingWorker) AttachContext(ctx context.Context) context.Context {
	if ctx == nil {
		panic("illegal value")
	}

	if atomic.LoadInt32(&p.runStatus) == runStatusUninitialized && p.ctx == nil {
		p.ctx, p.cancelFn = context.WithCancel(ctx)
		return p.ctx
	}

	panic("illegal state")
}

func (p *SyncingWorker) TryStartAttached() Status {
	if p.ctx == nil {
		panic("illegal state - was not prepared")
	}

	if canStart, status := p.ensureStart(); !canStart {
		return status
	}

	p.run()
	return runStatusStarted
}

func (p *SyncingWorker) run() {
	p.beforeStartFn(p.ctx)
	p.beforeStartFn = nil //avoid unnecessary retention

	go p._run()
}

func (p *SyncingWorker) _run() {
	err := p.runCommandsAndCleanup()

	p.afterStopFn(p.ctx, err)
	p.afterStopFn = nil
	atomic.StoreInt32(&p.runStatus, runStatusStopped)
}

func (p *SyncingWorker) runCommandsAndCleanup() (result interface{}) {

	isClosed := false
	cleanup := func() {
		if !isClosed {
			isClosed = true
			p.closeQueue()
		}
		p.cancelFn()
		if p.timeout != nil {
			p.timeout.Stop()
		}
	}

	defer func() {
		recovered := recover()
		if recovered != nil {
			result = recovered
		}
		cleanup()
	}()

	result = p.runCommands()

	if !atomic.CompareAndSwapInt32(&p.runStatus, runStatusStarted, runStatusStopping) {
		inslogger.FromContext(p.ctx).Error("illegal state")
		return
	}
	cleanup()

	var log insolar.Logger
	for cmd := range p.asyncCmd {
		err2 := cmd(p.ctx)
		if err2 != nil {
			if log == nil {
				log = inslogger.FromContext(p.ctx)
			}
			log.Error("ignored error: ", err2)
		}
	}
	return result
}

func (p *SyncingWorker) runCommands() (result interface{}) {

	defer func() {
		recovered := recover()
		if recovered != nil {
			result = recovered
		}
	}()

	for {
		var expired <-chan time.Time
		if p.timeout != nil {
			expired = p.timeout.C
		}

		select {
		case <-p.ctx.Done():
			return p.ctx.Err()
		case <-expired:
			return context.DeadlineExceeded
		case cmd, ok := <-p.asyncCmd:
			if !ok {
				return result
			}
			result = cmd(p.ctx)
			if result != nil {
				return result
			}
		}
	}
}

func (p *SyncingWorker) AsyncCall(fn func(context.Context) error) (successful bool) {
	if p.asyncCmd == nil {
		return false
	}
	defer func() {
		_ = recover()
		successful = false
	}()
	p.asyncCmd <- fn
	return true
}

func (p *SyncingWorker) SetDynamicDeadline(d time.Time) bool {
	return p.SyncCall(func(context.Context) error {
		if p.timeout != nil {
			p.timeout.Stop()
		}
		p.timeout = time.NewTimer(time.Until(d))
		return nil
	})
}

func (p *SyncingWorker) SyncCall(fn func(context.Context) error) (successful bool) {
	if fn == nil {
		panic("illegal value")
	}

	wg := sync.WaitGroup{}
	wg.Add(1)

	successful = p.AsyncCall(func(ctx context.Context) error {
		defer wg.Done()
		return fn(ctx)
	})
	if !successful {
		return false
	}
	wg.Wait()
	return true
}

func (p *SyncingWorker) closeQueue() {
	if p.asyncCmd == nil {
		return
	}

	defer func() {
		_ = recover()
	}()

	close(p.asyncCmd)
}

func (p *SyncingWorker) AsyncStop() {
	if p.asyncCmd == nil {
		return
	}
	defer func() {
		_ = recover()
	}()
	p.asyncCmd <- func(context.Context) error {
		p.Stop()
		return nil
	}
	close(p.asyncCmd)
}

func (p *SyncingWorker) GetContext() context.Context {
	if p.ctx == nil {
		panic("illegal state")
	}
	return p.ctx
}
