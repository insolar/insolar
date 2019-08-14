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
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/network/consensus/common/capacity"
	"github.com/insolar/insolar/network/consensus/common/syncrun"
	"github.com/insolar/insolar/network/consensus/common/timer"
	"sync/atomic"
	"time"

	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/network/consensus/common/pulse"
	"github.com/insolar/insolar/network/consensus/gcpv2/api"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/census"
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
	worker                  syncrun.SyncingWorker
	upstream                api.UpstreamController
	controlFeeder           api.PulseControlFeeder
	trafficControl          api.TrafficControlFeeder
	trafficThrottleDuration time.Duration

	finishedFn func()

	roundState uint32
}

func (p *RoundStateMachineWorker) OnPulseDetected() {
	p.applyState(RoundPulseDetected)
	p.trafficControl.SetTrafficLimit(capacity.LevelMinimal, p.trafficThrottleDuration)
	p.controlFeeder.OnPulseDetected()
}

func (p *RoundStateMachineWorker) OnFullRoundStarting() {
	p.applyState(RoundPulseAccepted)
	p.controlFeeder.OnFullRoundStarting()
}

func (p *RoundStateMachineWorker) PreparePulseChange(report api.UpstreamReport, ch chan<- api.UpstreamState) {
	p.applyState(RoundPulsePreparing)
	p.trafficControl.SetTrafficLimit(capacity.LevelZero, p.trafficThrottleDuration)
	p.controlFeeder.OnPreparePulseChange(report)
	p.upstream.PreparePulseChange(report, ch)
}

func (p *RoundStateMachineWorker) CommitPulseChangeByStateless(report api.UpstreamReport, pd pulse.Data, activeCensus census.Operational) {
	p.applyState(RoundPulsePreparing) // simulate prepare, but ignore upstream
	p.controlFeeder.OnPreparePulseChange(report)
	p.CommitPulseChange(report, pd, activeCensus)
}

func (p *RoundStateMachineWorker) CommitPulseChange(report api.UpstreamReport, pd pulse.Data, activeCensus census.Operational) {
	p.applyState(RoundPulseCommitted)
	p.trafficControl.SetTrafficLimit(capacity.LevelReduced, p.trafficThrottleDuration)
	p.controlFeeder.OnCommitPulseChange(report, pd, activeCensus)
	p.upstream.CommitPulseChange(report, pd, activeCensus)
}

func (p *RoundStateMachineWorker) CancelPulseChange() {
	p.applyState(RoundPulseAccepted)
	p.trafficControl.SetTrafficLimit(capacity.LevelMinimal, p.trafficThrottleDuration)
	p.controlFeeder.OnCancelPulseChange()
	p.upstream.CancelPulseChange()
}

func (p *RoundStateMachineWorker) ConsensusFinished(report api.UpstreamReport, expectedCensus census.Operational) {
	p.applyState(RoundConsensusFinished)
	p.trafficControl.ResumeTraffic()
	p.controlFeeder.OnConsensusFinished(report, expectedCensus)
	p.upstream.ConsensusFinished(report, expectedCensus)
}

func (p *RoundStateMachineWorker) ConsensusAborted() {
	p.applyState(RoundAborted)
	p.trafficControl.ResumeTraffic()
	p.controlFeeder.OnConsensusAborted()
	p.upstream.ConsensusAborted()
}

func (p *RoundStateMachineWorker) SetTimeoutTimer(deadline timer.Holder) {
	p.worker.SetDynamicDeadlineTimer(deadline)
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
	p.worker.Stop()
}

func (p *RoundStateMachineWorker) Stop() {
	p.worker.Stop()
}

func (p *RoundStateMachineWorker) preInit(ctx context.Context, upstream api.UpstreamController,
	controlFeeder api.PulseControlFeeder, trafficControl api.TrafficControlFeeder,
	trafficThrottleDuration time.Duration) context.Context {

	p.upstream = upstream
	p.controlFeeder = controlFeeder
	p.trafficControl = trafficControl
	p.trafficThrottleDuration = trafficThrottleDuration

	return p.worker.AttachContext(ctx)
}

func (p *RoundStateMachineWorker) init(starterFn func(), stopperFn func(), finishedFn func()) {
	p.worker.Init(syncrun.SyncingWorkerConfig{
		QueueLen: 10,
		BeforeStartFn: func(ctx context.Context) {
			if starterFn != nil {
				starterFn()
			}
			atomic.CompareAndSwapUint32(&p.roundState, uint32(RoundInactive), uint32(RoundAwaitingPulse))
		},
		AfterStopFn: func(ctx context.Context, err interface{}) {
			if stopperFn != nil {
				stopperFn()
			}
			p.applyFinishStatus(err)
		}})
	p.finishedFn = finishedFn
}

func (p *RoundStateMachineWorker) applyFinishStatus(err interface{}) {

	finishStatus := RoundAborted
	if err == nil {
		if atomic.CompareAndSwapUint32(&p.roundState, uint32(RoundConsensusFinished), uint32(RoundStopped)) || p.GetState() == RoundStopped {
			return
		}
	} else if err2, ok := err.(error); ok && err2 == context.DeadlineExceeded {
		finishStatus = RoundTimedOut
	}

	p.applyState(finishStatus)
	//	atomic.StoreUint32(&p.roundState, uint32(finishStatus))
}

func (p *RoundStateMachineWorker) SafeStartAndGetIsRunning() bool {
	return p.startAndGetIsRunning(true)
}

func (p *RoundStateMachineWorker) Start() {
	p.startAndGetIsRunning(false)
}

func (p *RoundStateMachineWorker) startAndGetIsRunning(safe bool) bool {
	if p.worker.TryStartAttached().IsStoppingOrStopped() {
		if safe {
			return false
		}
		panic("illegal state")
	}
	return true
}

func (p *RoundStateMachineWorker) IsStartedAndRunning() (bool, bool) {
	s := p.worker.GetStatus()
	return s.WasStarted(), s.IsRunning()
}

func (p *RoundStateMachineWorker) IsRunning() bool {
	return p.worker.GetStatus().IsRunning()
}

func (p *RoundStateMachineWorker) EnsureRunning() {
	s := p.worker.GetStatus()
	switch {
	case s.IsRunning():
		return
	case !s.WasStarted():
		panic("illegal state - not started")
	default:
		panic("illegal state - stopped")
	}
}

func (p *RoundStateMachineWorker) GetState() RoundState {
	return RoundState(atomic.LoadUint32(&p.roundState))
}

func (p *RoundStateMachineWorker) applyState(newState RoundState) {

	var transitionCompletionAction func()
	var curState RoundState
	doFinish := false

	logLevel := insolar.NoLevel
	logMsg := ""

loop:
	for {
		logLevel, logMsg = insolar.NoLevel, "state transition"
		curState = p.GetState()

		switch { // normal transitions
		case curState == newState:
			switch {
			case curState == RoundPulseDetected:
				return // OK
			case curState > RoundConsensusFinished:
				return // OK
			default:
				logLevel, logMsg = insolar.WarnLevel, "state self-loop transition"
				break loop // INCORRECT, do not apply
			}
		case curState+1 == newState && curState < RoundConsensusFinished:
			// next step from a non-stopped state
			if newState == RoundConsensusFinished {
				transitionCompletionAction = p.finishedFn
			}
		case curState == RoundConsensusFinished && newState > RoundConsensusFinished:
			// classification of stop
		case curState == RoundPulsePreparing && newState == RoundPulseAccepted:
			// prepare was cancelled by caller
		case curState > RoundConsensusFinished && newState > RoundConsensusFinished:
			// the first state is correct, don't change
			return // OK
		case curState > RoundConsensusFinished:
			// attempt to restart from a final state
			logLevel, logMsg = insolar.ErrorLevel, "state reset transition"
			break loop // TRANSITION IS NOT ALLOWED
		default:
			switch {
			case curState == RoundInactive:
				logLevel, logMsg = insolar.WarnLevel, "transition from inactive state"
			case curState > newState:
				logLevel, logMsg = insolar.ErrorLevel, "backward state transition"
			default:
				// DebugLevel
				logLevel, logMsg = insolar.WarnLevel, "fast-forward state transition"
			}

			doFinish = curState < RoundConsensusFinished && newState > RoundConsensusFinished

			switch { // transition from a state that require cancellation
			case curState == RoundPulsePreparing && newState > RoundConsensusFinished:
				transitionCompletionAction = func() {
					p.controlFeeder.OnFailedPreparePulseChange()
					p.upstream.CancelPulseChange()
					p.upstream.ConsensusAborted()
				}
			case curState == RoundPulsePreparing && newState == RoundConsensusFinished:
				transitionCompletionAction = func() {
					p.controlFeeder.OnFailedPreparePulseChange()
					p.upstream.CancelPulseChange()
				}
			case curState == RoundPulsePreparing:
				transitionCompletionAction = p.upstream.CancelPulseChange
			case curState < RoundConsensusFinished && newState > RoundConsensusFinished:
				transitionCompletionAction = p.upstream.ConsensusAborted
			}
		}

		if atomic.CompareAndSwapUint32(&p.roundState, uint32(curState), uint32(newState)) {
			if doFinish && p.finishedFn != nil {
				p.finishedFn()
			}
			if transitionCompletionAction != nil {
				transitionCompletionAction()
			}

			break loop
		}
	}

	if logLevel == insolar.NoLevel {
		return
	}

	log := inslogger.FromContext(p.worker.GetContext())
	switch logLevel {
	case insolar.ErrorLevel:
		log.Errorf("forbidden %s: current=%v new=%v", logMsg, curState, newState)
	case insolar.WarnLevel:
		log.Warnf("unexpected %s: current=%v new=%v", logMsg, curState, newState)
	case insolar.DebugLevel:
		log.Debugf("unexpected %s: current=%v new=%v", logMsg, curState, newState)
	default:
		log.Infof("unexpected %s: current=%v new=%v", logMsg, curState, newState)
	}
}
