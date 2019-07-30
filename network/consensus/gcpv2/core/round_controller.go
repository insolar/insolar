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
	"fmt"
	"github.com/insolar/insolar/network/consensus/gcpv2/core/coreapi"
	"sync"
	"time"

	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/network/consensus/common/endpoints"
	"github.com/insolar/insolar/network/consensus/common/pulse"
	"github.com/insolar/insolar/network/consensus/gcpv2/api"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/census"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/profiles"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/transport"

	errors2 "github.com/insolar/insolar/network/consensus/gcpv2/core/errors"
)

type RoundStrategyFactory interface {
	CreateRoundStrategy(chronicle api.ConsensusChronicles, config api.LocalNodeConfiguration) (RoundStrategy, PhaseControllersBundle)
}

type RoundStrategy interface {
	GetBaselineWeightForNeighbours() uint32
	ShuffleNodeSequence(n int, swap func(i, j int))

	ConfigureRoundContext(ctx context.Context, expectedPulse pulse.Number, self profiles.LocalNode) context.Context
	AdjustConsensusTimings(timings *api.RoundTimings)
}

var _ api.RoundController = &PhasedRoundController{}

type PhasedRoundController struct {
	rw sync.RWMutex

	/* Derived from the provided externally - set at init() or start(). Don't need mutex */
	chronicle api.ConsensusChronicles
	bundle    PhaseControllersBundle

	// fullCancel     context.CancelFunc /* cancels prepareCancel as well */
	prepareCancel  context.CancelFunc
	prevPulseRound api.RoundController
	minDuration    time.Duration

	roundWorker RoundStateMachineWorker

	/* Other fields - need mutex */
	prepR *PrepRealm
	realm FullRealm
}

func NewPhasedRoundController(strategy RoundStrategy, chronicle api.ConsensusChronicles, bundle PhaseControllersBundle,
	transport transport.Factory, config api.LocalNodeConfiguration,
	controlFeeder api.ConsensusControlFeeder, candidateFeeder api.CandidateControlFeeder, ephemeralFeeder api.EphemeralControlFeeder,
	prevPulseRound api.RoundController) *PhasedRoundController {

	r := &PhasedRoundController{chronicle: chronicle, prevPulseRound: prevPulseRound, bundle: bundle}

	r.realm.coreRealm.initBefore(&r.rw, strategy, transport, config, chronicle.GetLatestCensus(),
		controlFeeder, candidateFeeder, ephemeralFeeder)

	nbhSizes := r.realm.initBefore(transport)
	r.realm.coreRealm.initBeforePopulation(nbhSizes)

	return r
}

func (r *PhasedRoundController) PrepareConsensusRound(upstream api.UpstreamController) {
	r.rw.Lock()
	defer r.rw.Unlock()

	r.realm.coreRealm.roundContext = r.roundWorker.preInit(
		r.realm.coreRealm.strategy.ConfigureRoundContext(
			r.realm.config.GetParentContext(),
			r.realm.initialCensus.GetExpectedPulseNumber(),
			r.realm.GetLocalProfile(),
		), upstream)

	r.realm.coreRealm.stateMachine = &r.roundWorker

	r.realm.coreRealm.postponedPacketFn = func(packet transport.PacketParser, from endpoints.Inbound, verifyFlags coreapi.PacketVerifyFlags) bool {
		// There is no real context for delayed reprocessing, so we use the round context
		ctx := r.realm.coreRealm.roundContext
		_, err := r.handlePacket(ctx, packet, from, verifyFlags)
		if err != nil {
			inslogger.FromContext(ctx).Error(err)
		}
		return true
	}

	inslogger.FromContext(r.realm.roundContext).Debugf(
		"Starting consensus round: self={%v}, bundle=%v, census=%+v", r.realm.GetLocalProfile(), r.bundle, r.realm.initialCensus)

	preps := r.bundle.CreatePrepPhaseControllers()
	if len(preps) == 0 {
		panic("illegal state - no prep realm")
	}

	prep := PrepRealm{coreRealm: &r.realm.coreRealm}
	prep.init(
		func(successful bool) {
			// RUNS under lock
			if r.prepR == nil {
				return
			}
			defer r.prepR.stop() // initiates handover from PrepRealm
			r.prepR = nil
			r.roundWorker.Start() // ensures that worker was started
			r._startFullRealm(successful)
		})

	var prepCtx context.Context
	// r.prepareCancel will be cancelled through r.fullCancel()
	prepCtx, r.prepareCancel = context.WithCancel(r.realm.roundContext)

	r.prepR = &prep
	r.prepR.beforeStart(prepCtx, preps)

	r.roundWorker.init(func() {
		// requires r.roundWorker.StartXXX to happen under lock
		r._setStartedAt()
		if r.prepR != nil { // PrepRealm can be finished before starting
			r.prepR._startWorkers(prepCtx, preps)
		}
	},
		// both further handlers MUST not use round's lock inside
		r.onConsensusStopper,
		r.onConsensusFinished,
	)

	r.realm.coreRealm.pollingWorker.Start(r.realm.roundContext, 100*time.Millisecond)
	r.prepR.prepareEphemeralPolling(prepCtx)
}

func (r *PhasedRoundController) onConsensusStopper() {
	latest := r.chronicle.GetLatestCensus()

	inslogger.FromContext(r.realm.roundContext).Debugf(
		"Stopping consensus round: self={%v}, bundle=%v, census=%+v", r.realm.GetLocalProfile(), r.bundle, latest)

	if latest.GetOnlinePopulation().GetLocalProfile().IsJoiner() {
		panic("DEBUG FAIL-FAST: local remains as joiner")
	}

	if r.chronicle.GetExpectedCensus() == nil {
		panic("DEBUG FAIL-FAST: consensus didn't finish")
	}

	// TODO print purgatory
}

func (r *PhasedRoundController) onConsensusFinished() {
	r.rw.Lock()
	defer r.rw.Unlock()
	r._onConsensusFinished()
}

func (r *PhasedRoundController) _onConsensusFinished() {
	// prevents memory leak and disallows older controller to handle messages after a consensus is done
	if r.prevPulseRound != nil {
		r.prevPulseRound.StopConsensusRound()
	}
	r.prevPulseRound = nil
}

func (r *PhasedRoundController) _setStartedAt() {
	if r.realm.roundStartedAt.IsZero() { // can be called a few times
		r.realm.roundStartedAt = time.Now()
	}
}

func (r *PhasedRoundController) StopConsensusRound() {
	r.rw.Lock()
	defer r.rw.Unlock()
	r.roundWorker.Stop()
	r._onConsensusFinished() // double-check, just to be on a safe side
}

func (r *PhasedRoundController) IsRunning() bool {
	return r.roundWorker.IsRunning()
}

func (r *PhasedRoundController) beforeHandlePacket() (prep *PrepRealm, current pulse.Number,
	possibleNext pulse.Number, prev api.RoundController) {

	r.rw.RLock()
	defer r.rw.RUnlock()
	if r.prepR != nil {
		return r.prepR, r.realm.coreRealm.initialCensus.GetExpectedPulseNumber(), 0, r.prevPulseRound
	}
	return nil, r.realm.GetPulseNumber(), r.realm.GetNextPulseNumber(), r.prevPulseRound
}

/*
RUNS under lock.
Can be called from a polling function (for ephemeral), and happen BEFORE PrepRealm start
*/
func (r *PhasedRoundController) _startFullRealm(prepWasSuccessful bool) {

	r.roundWorker.OnFullRoundStarting()

	chronicle := r.chronicle
	lastCensus := chronicle.GetLatestCensus()
	pd := &r.realm.pulseData

	if lastCensus.GetCensusState() == census.PrimingCensus {
		/* This is the priming census */
		priming := lastCensus.GetMandateRegistry().GetPrimingCloudHash()
		lastCensus.(census.Prime).MakeExpected(pd.PulseNumber, priming, priming).MakeActive(*pd)
	} else {
		if lastCensus.GetPulseNumber() != pd.PulseNumber {
			// TODO inform control feeder when our pulse is less
			panic(fmt.Sprintf("illegal state - pulse number of expected census (%v) and of the realm (%v) are mismatched for %v", lastCensus.GetPulseNumber(), pd.PulseNumber, r.realm.GetSelfNodeID()))
		}
		if !lastCensus.IsActive() {
			/* Auto-activation of the prepared lastCensus */
			expCensus := chronicle.GetExpectedCensus()
			lastCensus = expCensus.MakeActive(*pd)
		}
	}

	active := chronicle.GetActiveCensus()
	r.realm.start(active, active.GetOnlinePopulation(), r.bundle)
	r.roundWorker.SetTimeout(r.realm.roundStartedAt.Add(r.realm.timings.EndOfConsensus))
}

func (r *PhasedRoundController) ensureStarted() bool {

	isStarted, isRunning := r.roundWorker.IsStartedAndRunning()
	if isStarted {
		return isRunning
	}

	r.rw.Lock() // ensure that starting closure will run under lock
	defer r.rw.Unlock()
	return r.roundWorker.SafeStartAndGetIsRunning()
}

func (r *PhasedRoundController) HandlePacket(ctx context.Context, packet transport.PacketParser,
	from endpoints.Inbound) (api.RoundControlCode, error) {

	return r.handlePacket(ctx, packet, from, coreapi.DefaultVerify)
}

func (r *PhasedRoundController) handlePacket(ctx context.Context, packet transport.PacketParser, from endpoints.Inbound,
	verifyFlags coreapi.PacketVerifyFlags) (api.RoundControlCode, error) {

	pn := packet.GetPulseNumber()
	/* a separate method with lock is to ensure that further packet processing is not connected to a lock */
	prep, filterPN, _, prev := r.beforeHandlePacket()

	// TODO HACK - network doesnt have information about pulsars to validate packets, hackIgnoreVerification must be removed when fixed
	defaultOptions := coreapi.SkipVerify // coreapi.DefaultVerify

	if prev != nil && filterPN > pn {
		// something from a previous round?
		_, err := prev.HandlePacket(ctx, packet, from)
		return api.KeepRound, err
		//defaultOptions = coreapi.SkipVerify // validation was done by the prev controller
	}

	if r.realm.ephemeralFeeder != nil && !packet.GetPacketType().IsEphemeralPacket() {
		_, err := r.realm.VerifyPacketAuthenticity(ctx, packet, from, nil, coreapi.DefaultVerify, nil, defaultOptions)
		if err == nil {
			err = r.realm.ephemeralFeeder.OnNonEphemeralPacket(ctx, packet, from)
		}
		return api.KeepRound, err
	}

	var err error
	if prep != nil {
		// NB! Round may NOT be running yet here - ensure it is working before calling the state machine
		r.ensureStarted()

		if !pn.IsUnknown() && (filterPN.IsUnknown() || filterPN == pn) {
			r.roundWorker.OnPulseDetected()
		}

		err = prep.dispatchPacket(ctx, packet, from, defaultOptions) // prep realm can't inherit flags
	} else {
		err = r.realm.dispatchPacket(ctx, packet, from, verifyFlags|defaultOptions)
	}

	isNextPulse, nextPN := errors2.IsNextPulseArrivedError(err)
	if !isNextPulse {
		if errors2.IsNextPulseError(err) {
			r.roundWorker.onUnexpectedPulse(pn)
		}
		return api.KeepRound, err
	}

	if r.roundWorker.IsRunning() {
		canStop := false
		endOfConsensus := r.realm.GetStartedAt().Add(r.realm.timings.EndOfConsensus)
		if r.realm.ephemeralFeeder != nil {
			canStop = r.realm.ephemeralFeeder.CanStopOnHastyPulse(nextPN, endOfConsensus)
		} else {
			canStop = r.realm.controlFeeder.CanStopOnHastyPulse(nextPN, endOfConsensus)
		}

		if !canStop {
			r.roundWorker.onUnexpectedPulse(pn)
			return api.KeepRound, err
		}
		r.StopConsensusRound()
		// TODO should allow some time to handover a broken population?
		// TODO build a one-node population for Suspected mode
	} else {
		expected := r.chronicle.GetExpectedCensus()
		if expected == nil || !expected.GetOnlinePopulation().IsValid() {
			return api.NextRoundTerminate, fmt.Errorf("next population is invalid or not ready on: %v", err)
		}
	}
	r.roundWorker.onNextPulse(pn)
	return api.StartNextRound, err
}
