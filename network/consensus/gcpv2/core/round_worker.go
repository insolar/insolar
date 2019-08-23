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
	"github.com/insolar/insolar/network/consensus/common/capacity"
	"github.com/insolar/insolar/network/consensus/gcpv2/api"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/census"
	"github.com/insolar/insolar/pulse"
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

	trafficControl          api.TrafficControlFeeder
	trafficThrottleDuration time.Duration

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
	p.trafficControl.SetTrafficLimit(capacity.LevelMinimal, p.trafficThrottleDuration)
}

func (p *RoundStateMachineWorker) OnFullRoundStarting() {
	p.applyState(RoundPulseAccepted)
}

func (p *RoundStateMachineWorker) PreparePulseChange(report api.UpstreamReport, ch chan<- api.UpstreamState) {
	p.applyState(RoundPulsePreparing)
	p.trafficControl.SetTrafficLimit(capacity.LevelZero, p.trafficThrottleDuration)
	p.UpstreamController.PreparePulseChange(report, ch)
}

func (p *RoundStateMachineWorker) CommitPulseChange(report api.UpstreamReport, pd pulse.Data, activeCensus census.Operational) {
	p.applyState(RoundPulseCommitted)
	p.trafficControl.SetTrafficLimit(capacity.LevelReduced, p.trafficThrottleDuration)
	p.UpstreamController.CommitPulseChange(report, pd, activeCensus)
}

func (p *RoundStateMachineWorker) CommitPulseChangeByStateless(report api.UpstreamReport, pd pulse.Data, activeCensus census.Operational) {
	p.applyState(RoundPulsePreparing)
	p.applyState(RoundPulseCommitted)
	p.trafficControl.SetTrafficLimit(capacity.LevelReduced, p.trafficThrottleDuration)
	p.UpstreamController.CommitPulseChange(report, pd, activeCensus)
}

func (p *RoundStateMachineWorker) CancelPulseChange() {
	p.applyState(RoundPulseAccepted)
	p.trafficControl.SetTrafficLimit(capacity.LevelMinimal, p.trafficThrottleDuration)
	p.UpstreamController.CancelPulseChange()
}

func (p *RoundStateMachineWorker) ConsensusFinished(report api.UpstreamReport, expectedCensus census.Operational) {
	p.applyState(RoundConsensusFinished)
	p.trafficControl.ResumeTraffic()
	p.UpstreamController.ConsensusFinished(report, expectedCensus)
}

func (p *RoundStateMachineWorker) SetTimeout(deadline time.Time) {
	p.sync(func() {
		p.timeout = time.After(time.Until(deadline))
	})
}

func (p *RoundStateMachineWorker) OnRoundStopped(ctx context.Context) {
	err := ctx.Err()
	switch {
	case err == nil:
		p.applyState(RoundStopped)
	case err == context.DeadlineExceeded:
		p.applyState(RoundTimedOut)
	default:
		p.applyState(RoundAborted)
	}
}

func (p *RoundStateMachineWorker) OnPrepRoundFailed() {
	p.applyState(RoundAborted)
}

func (p *RoundStateMachineWorker) onUnexpectedPulse(pulse.Number) {

}

func (p *RoundStateMachineWorker) onNextPulse(pulse.Number) {
	p.cancelFn()
}

func (p *RoundStateMachineWorker) Stop() {
	p.cancelFn()
}

func (p *RoundStateMachineWorker) preInit(ctx context.Context, upstream api.UpstreamController,
	trafficControl api.TrafficControlFeeder, trafficThrottleDuration time.Duration) context.Context {
	if p.cancelFn != nil {
		panic("illegal state - was initialized")
	}

	p.UpstreamController = upstream
	p.trafficControl = trafficControl
	p.trafficThrottleDuration = trafficThrottleDuration

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

func (p *RoundStateMachineWorker) SafeStartAndGetIsRunning() bool {
	return p.startAndGetIsRunning(true)
}

func (p *RoundStateMachineWorker) Start() {
	p.startAndGetIsRunning(false)
}

func (p *RoundStateMachineWorker) startAndGetIsRunning(safe bool) bool {
	if !atomic.CompareAndSwapInt32(&p.runStatus, runStatusInitialized, runStatusStarted) {
		if atomic.LoadInt32(&p.runStatus) >= runStatusStopping {
			if safe {
				return false
			}
			panic("illegal state")
		}
		return true // isRunning
	}
	if p.starterFn != nil {
		p.starterFn()
	}
	atomic.CompareAndSwapUint32(&p.roundState, uint32(RoundInactive), uint32(RoundAwaitingPulse))
	go p.stateWorker()
	return true
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
		p.trafficControl.ResumeTraffic()
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

func (p *RoundStateMachineWorker) IsStartedAndRunning() (bool, bool) {
	s := atomic.LoadInt32(&p.runStatus)
	return s >= runStatusStarted, s == runStatusStarted
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

func (p *RoundStateMachineWorker) applyState(newState RoundState) {
	for {
		attention := false

		curState := p.GetState()
		switch { // normal transitions
		case curState == newState:
			return // no transition
		case newState == curState+1 && curState <= RoundConsensusFinished:
			break // next non-stopped step
		case curState == RoundConsensusFinished && newState > RoundConsensusFinished:
			break // stop
		case curState == RoundPulsePreparing && newState == RoundPulseAccepted:
			break // prepare was cancelled by caller
		case curState > RoundConsensusFinished && newState > RoundConsensusFinished:
			// the first state is correct, don't change
			return
		case curState > RoundConsensusFinished:
			// attempt to restart from a final state
			inslogger.FromContext(p.ctx).Errorf("reset transition attempt: current=%v new=%v", curState, newState)
			return
		// case curState < RoundConsensusFinished && newState == RoundConsensusFinished:
		//	break // early finish
		default:
			attention = true
		}
		if atomic.CompareAndSwapUint32(&p.roundState, uint32(curState), uint32(newState)) {
			if attention {
				if newState < curState {
					inslogger.FromContext(p.ctx).Warnf("backward state transition: current=%v new=%v", curState, newState)
				} else {
					inslogger.FromContext(p.ctx).Infof("fast-forward state transition: current=%v new=%v", curState, newState)
				}

				switch { // transition from a state that require cancellation
				case curState == RoundPulsePreparing:
					p.UpstreamController.CancelPulseChange()
				case curState < RoundConsensusFinished && newState > RoundConsensusFinished:
					p.trafficControl.ResumeTraffic()
					p.UpstreamController.ConsensusAborted()
				}
			}
			return
		}
	}
}
