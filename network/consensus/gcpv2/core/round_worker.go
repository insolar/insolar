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

package core

import (
	"context"
	"sync/atomic"
	"time"

	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/network/consensus/common/pulse"
	"github.com/insolar/insolar/network/consensus/gcpv2/api"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/census"
)

const (
	runStatusUninitialized = iota
	runStatusInitialized
	runStatusStarted
	runStatusStopping
	runStatusStopped
)

type RoundState uint8

const (
	RoundInactive RoundState = iota
	RoundAwaitingPulse
	RoundPulseDetected
	RoundPulseAccepted
	RoundPulsePreparing
	RoundPulseCommitted
	RoundConsensusFinished
	RoundStopped
	RoundTimedOut
	RoundAborted
)

var _ api.RoundStateCallback = &RoundStateMachineWorker{}

type RoundStateMachineWorker struct {
	api.UpstreamController

	ctx      context.Context
	cancelFn context.CancelFunc

	runStatus  int32 // atomic
	roundState uint32

	timeout <-chan time.Time

	asyncCmd   chan func()
	syncCmd    chan func()
	starterFn  func()
	stopperFn  func()
	finishedFn func()
}

func (p *RoundStateMachineWorker) OnPulseDetected() {
	p.applyState(RoundPulseDetected)
}

func (p *RoundStateMachineWorker) OnFullRoundStarting() {
	p.applyState(RoundPulseAccepted)
}

func (p *RoundStateMachineWorker) PreparePulseChange(report api.UpstreamReport, ch chan<- api.UpstreamState) {
	p.forceState(RoundPulsePreparing)
	p.UpstreamController.PreparePulseChange(report, ch)
}

func (p *RoundStateMachineWorker) CommitPulseChange(report api.UpstreamReport, pd pulse.Data, activeCensus census.Operational) {
	p.forceState(RoundPulseCommitted)
	p.UpstreamController.CommitPulseChange(report, pd, activeCensus)
}

func (p *RoundStateMachineWorker) CommitPulseChangeByStateless(report api.UpstreamReport, pd pulse.Data, activeCensus census.Operational) {
	p.forceState(RoundPulsePreparing)
	p.applyState(RoundPulseCommitted)
	p.UpstreamController.CommitPulseChange(report, pd, activeCensus)
}

func (p *RoundStateMachineWorker) CancelPulseChange() {
	p.applyState(RoundPulseAccepted)
	p.UpstreamController.CancelPulseChange()
}

func (p *RoundStateMachineWorker) ConsensusFinished(report api.UpstreamReport, expectedCensus census.Operational) {
	p.forceState(RoundConsensusFinished)
	p.UpstreamController.ConsensusFinished(report, expectedCensus)
}

func (p *RoundStateMachineWorker) SetTimeout(deadline time.Time) {
	p.sync(func() {
		p.timeout = time.After(time.Until(deadline))
	})
}

func (p *RoundStateMachineWorker) onUnexpectedPulse(pulse.Number) {

}

func (p *RoundStateMachineWorker) onNextPulse(pulse.Number) {
	p.cancelFn()
}

func (p *RoundStateMachineWorker) Stop() {
	p.cancelFn()
}

func (p *RoundStateMachineWorker) preInit(ctx context.Context, upstream api.UpstreamController) context.Context {
	if p.cancelFn != nil {
		panic("illegal state - was initialized")
	}

	p.UpstreamController = upstream
	p.asyncCmd = make(chan func(), 10)
	p.syncCmd = make(chan func())
	p.ctx, p.cancelFn = context.WithCancel(ctx)
	return p.ctx
}

func (p *RoundStateMachineWorker) init(starterFn func(), stopperFn func(), finishedFn func()) {
	if !atomic.CompareAndSwapInt32(&p.runStatus, runStatusUninitialized, runStatusInitialized) {
		panic("illegal state")
	}
	p.starterFn = starterFn
	p.stopperFn = stopperFn
	p.finishedFn = finishedFn
}

func (p *RoundStateMachineWorker) Start() {
	if !atomic.CompareAndSwapInt32(&p.runStatus, runStatusInitialized, runStatusStarted) {
		panic("illegal state")
	}
	if p.starterFn != nil {
		p.starterFn()
	}
	atomic.CompareAndSwapUint32(&p.roundState, uint32(RoundInactive), uint32(RoundAwaitingPulse))
	go p.stateWorker()
}

func (p *RoundStateMachineWorker) stateWorker() {

	exitState := p.runToLastState()

	switch {
	case exitState == RoundStopped:
		if atomic.CompareAndSwapUint32(&p.roundState, uint32(RoundConsensusFinished), uint32(RoundStopped)) ||
			p.GetState() == RoundStopped {
			break
		}
		exitState = RoundAborted
		fallthrough
	default:
		atomic.StoreUint32(&p.roundState, uint32(exitState))
	}
	atomic.StoreInt32(&p.runStatus, runStatusStopped)
}

func (p *RoundStateMachineWorker) runToLastState() (exitState RoundState) {
	defer func() {
		p.cancelFn()
		recovered := recover()
		if recovered != nil {
			exitState = RoundAborted
			// TODO log
		}
		if p.stopperFn != nil {
			p.stopperFn()
		}
	}()

	exitState = RoundAborted
	for {
		select {
		case <-p.ctx.Done():
			close(p.asyncCmd)
			close(p.syncCmd)
		case <-p.timeout:
			close(p.asyncCmd)
			close(p.syncCmd)
			exitState = RoundTimedOut
		case cmd, ok := <-p.asyncCmd:
			if ok {
				cmd()
				continue
			}
		case cmd, ok := <-p.syncCmd:
			if ok {
				p.flushAsync()
				cmd()
				continue
			}
		}

		if !atomic.CompareAndSwapInt32(&p.runStatus, runStatusStarted, runStatusStopping) {
			panic("illegal state")
		}
		p.cancelFn()

		for cmd := range p.asyncCmd { // ensure that a queued command is read
			cmd()
		}
		for cmd := range p.syncCmd { // ensure that a queued command is read
			cmd()
		}
		return RoundStopped
	}
}

func (p *RoundStateMachineWorker) flushAsync() {
	p.asyncCmd <- nil // there is a chance of deadlock when asyncCmd is full ...
	for {
		select {
		case cmd, ok := <-p.asyncCmd:
			if !ok || cmd == nil {
				return
			}
			cmd()
		default:
			return
		}
	}
}

func (p *RoundStateMachineWorker) IsRunning() bool {
	return atomic.LoadInt32(&p.runStatus) == runStatusStarted
}

func (p *RoundStateMachineWorker) EnsureRunning() {
	switch atomic.LoadInt32(&p.runStatus) {
	case runStatusStarted:
		return
	case runStatusUninitialized, runStatusInitialized:
		panic("illegal state - not started")
	default:
		panic("illegal state - stopped")
	}
}

func (p *RoundStateMachineWorker) sync(fn func()) {
	defer func() {
		_ = recover()
	}()
	p.syncCmd <- fn
}

func (p *RoundStateMachineWorker) async(fn func()) {
	defer func() {
		_ = recover()
	}()
	p.asyncCmd <- fn
}

func (p *RoundStateMachineWorker) GetState() RoundState {
	return RoundState(atomic.LoadUint32(&p.roundState))
}

func (p *RoundStateMachineWorker) forceState(newState RoundState) {
	p.applyState(newState)
}

func (p *RoundStateMachineWorker) applyState(newState RoundState) {
	for {
		curState := p.GetState()
		switch {
		case curState == newState:
			return
		case newState == curState+1:
		case curState == RoundAwaitingPulse && newState <= RoundPulsePreparing:
		case curState == RoundPulsePreparing && newState == RoundPulseAccepted:
		default:
			// invalid transition attempt
			inslogger.FromContext(p.ctx).Warnf("invalid state transition: current=%v new=%v", curState, newState)
			//return
		}
		if atomic.CompareAndSwapUint32(&p.roundState, uint32(curState), uint32(newState)) {
			return
		}
	}
}
