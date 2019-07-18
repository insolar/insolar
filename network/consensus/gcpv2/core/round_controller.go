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
	"github.com/insolar/insolar/instrumentation/inslogger"
	"sync"
	"time"

	"github.com/insolar/insolar/network/consensus/common/endpoints"
	"github.com/insolar/insolar/network/consensus/common/pulse"
	"github.com/insolar/insolar/network/consensus/gcpv2/api"
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

type PhasedRoundController struct {
	rw sync.RWMutex

	/* Derived from the provided externally - set at init() or start(). Don't need mutex */
	chronicle      api.ConsensusChronicles
	bundle         PhaseControllersBundle
	fullCancel     context.CancelFunc /* cancels prepareCancel as well */
	prepareCancel  context.CancelFunc
	prevPulseRound api.RoundController

	/* Other fields - need mutex */
	isRunning bool
	prepR     *PrepRealm
	realm     FullRealm
}

func NewPhasedRoundController(strategy RoundStrategy, chronicle api.ConsensusChronicles, bundle PhaseControllersBundle,
	transport transport.Factory, config api.LocalNodeConfiguration,
	controlFeeder api.ConsensusControlFeeder, candidateFeeder api.CandidateControlFeeder,
	prevPulseRound api.RoundController) *PhasedRoundController {

	r := &PhasedRoundController{chronicle: chronicle, prevPulseRound: prevPulseRound, bundle: bundle}

	r.realm.coreRealm.initBefore(&r.rw, strategy, transport, config, chronicle.GetLatestCensus())
	nbhSizes := r.realm.initBefore(transport, controlFeeder, candidateFeeder)
	r.realm.coreRealm.initBeforePopulation(controlFeeder.GetRequiredPowerLevel(), nbhSizes)

	return r
}

func (r *PhasedRoundController) StartConsensusRound(upstream api.UpstreamController) {
	r.rw.Lock()
	defer r.rw.Unlock()

	if r.fullCancel != nil {
		panic("was started once")
	}

	ctx := r.realm.config.GetParentContext()
	ctx, r.fullCancel = context.WithCancel(ctx)

	r.realm.coreRealm.roundContext = r.realm.coreRealm.strategy.ConfigureRoundContext(
		ctx,
		r.realm.initialCensus.GetExpectedPulseNumber(),
		r.realm.GetLocalProfile(),
	)

	r.isRunning = true

	r.realm.coreRealm.roundStartedAt = time.Now()
	r.realm.coreRealm.upstream = upstream
	r.realm.coreRealm.postponedPacketFn = func(packet transport.PacketParser, from endpoints.Inbound, verifyFlags PacketVerifyFlags) {
		// There is no real context for delayed reprocessing, so we use the round context
		ctx := r.realm.coreRealm.roundContext
		err := r.handlePacket(ctx, packet, from, verifyFlags)

		if err != nil {
			inslogger.FromContext(ctx).Error(err)
		}
	}

	preps := r.bundle.CreatePrepPhaseControllers()

	if len(preps) > 0 {
		prep := PrepRealm{coreRealm: &r.realm.coreRealm}
		prep.init(
			r.bundle.IsEphemeralPulseAllowed(),
			func(successful bool) {
				if r.prepR == nil {
					return
				}
				defer r.prepR.stop() // initiates handover from PrepRealm
				r.prepR = nil
				r.startFullRealm()
			})

		// r.prepareCancel will be cancelled through r.fullCancel()
		ctx, r.prepareCancel = context.WithCancel(r.realm.roundContext)

		r.prepR = &prep
		r.prepR.start(ctx, preps)
	} else {
		r.prepR = nil
		r.startFullRealm()
	}
}

/*
Returns true when this round was running.
*/
func (r *PhasedRoundController) StopConsensusRound() {
	r.rw.Lock()
	defer r.rw.Unlock()

	r.prevPulseRound = nil // prevents memory leak

	if r.fullCancel == nil || !r.isRunning {
		return
	}
	r.isRunning = false
	r.fullCancel()
}

/* LOCK: simple */
func (r *PhasedRoundController) IsRunning() bool {
	r.rw.RLock()
	defer r.rw.RUnlock()
	return r.isRunning
}

/* Checks if Controller can handle a new packet, and which Realm should do it. If result = nil, then FullRealm is used */
/* LOCK: simple */
func (r *PhasedRoundController) beforeHandlePacket() (prep *PrepRealm, current pulse.Number, possibleNext pulse.Number, err error) {
	r.rw.RLock()
	defer r.rw.RUnlock()

	if r.fullCancel == nil {
		return nil, 0, 0, errors2.NewRoundStateError("not started")
	}
	if !r.isRunning {
		return nil, 0, 0, errors2.NewRoundStateError("stopped")
	}

	if r.prepR != nil {
		return r.prepR, r.realm.coreRealm.initialCensus.GetExpectedPulseNumber(), 0, nil
	}
	return nil, r.realm.GetPulseNumber(), r.realm.GetNextPulseNumber(), nil
}

func (r *PhasedRoundController) startFullRealm() {

	chronicle := r.chronicle
	lastCensus := chronicle.GetLatestCensus()
	pd := &r.realm.pulseData

	if lastCensus.IsActive() && lastCensus.GetPulseNumber().IsUnknown() {
		/* This is the priming lastCensus */
		b := chronicle.GetActiveCensus().CreateBuilder(pd.PulseNumber, true)
		r.realm.preparePrimingMembers(b.GetPopulationBuilder())
		priming := lastCensus.GetMandateRegistry().GetPrimingCloudHash()
		b.SetGlobulaStateHash(priming)
		b.SealCensus()
		b.BuildAndMakeExpected(priming)
		chronicle.GetExpectedCensus().MakeActive(*pd)
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
}

func (r *PhasedRoundController) HandlePacket(ctx context.Context, packet transport.PacketParser, from endpoints.Inbound) error {
	return r.handlePacket(ctx, packet, from, DefaultVerify)
}

func (r *PhasedRoundController) handlePacket(ctx context.Context, packet transport.PacketParser, from endpoints.Inbound,
	verifyFlags PacketVerifyFlags) error {

	/* a separate method with lock is to ensure that further packet processing is not connected to a lock */
	prep, filterPN, nextPN, err := r.beforeHandlePacket()
	if err != nil {
		return err
	}

	if !filterPN.IsUnknown() {
		pn := packet.GetPulseNumber()
		if !pn.IsUnknown() && filterPN != pn {
			if nextPN.IsUnknown() || nextPN != pn {
				return errors2.NewPulseRoundMismatchError(pn,
					fmt.Sprintf("packet pulse number mismatched: expected=%v, actual=%v", filterPN, pn))
			}
			return errors2.NewNextPulseArrivedError(pn)
		}
	}

	if prep != nil {
		return prep.dispatchPacket(ctx, packet, from, DefaultVerify) //prep realm can't inherit any flags
	}
	return r.realm.dispatchPacket(ctx, packet, from, verifyFlags)

	// if pt.IsMemberPacket() {
	//	memberPacket := packet.GetMemberPacket()
	//	if memberPacket == nil {
	//		panic("missing parser for phased packet")
	//	}
	//
	//	strictSenderCheck, err = r.verifyRoute(ctx, packet)
	//	if err != nil {
	//		return err
	//	}
	//
	//	if prep == nil { // Full realm is active - we can use node projections
	//		route, err := r.realm.getPacketDispatcher(pt)
	//		if err != nil {
	//			return err
	//		}
	//
	//		pop := r.realm.GetPopulation()
	//		sid := packet.GetSourceID()
	//		src := pop.GetNodeAppearance(sid)
	//		if src == nil {
	//			if route.HasUnknownMemberHandler() {
	//				src, err = route.dispatchUnknownMemberPacket(ctx, memberPacket, from)
	//				if err != nil {
	//					return err
	//				}
	//			}
	//			if src == nil {
	//				return fmt.Errorf("unknown source id (%v)", sid)
	//			}
	//		}
	//
	//		if !preVerified {
	//			err = src.VerifyPacketAuthenticity(packet, from, strictSenderCheck)
	//			if err != nil {
	//				return err
	//			}
	//		}
	//
	//		if route.HasMemberHandler() {
	//			return route.dispatchMemberPacket(ctx, memberPacket, src)
	//		}
	//		return route.dispatchHostPacket(ctx, packet, from)
	//	}
	// }
	//
	// //TODO HACK - network doesnt have information about pulsars to validate packets, hackIgnoreVerification must be removed when fixed
	// hackIgnoreVerification := !packet.GetPacketType().IsMemberPacket()
	//
	// if !preVerified && !hackIgnoreVerification {
	//	err = r.realm.coreRealm.VerifyPacketAuthenticity(packet, from, strictSenderCheck)
	//	if err != nil {
	//		return err
	//	}
	// }
	//
	// if prep != nil { // Prep realm is active
	//	return prep.handleHostPacket(ctx, packet, from)
	// }
	// route, err := r.realm.getPacketDispatcher(pt)
	// if err != nil {
	//	return err
	// }
	// return route.dispatchHostPacket(ctx, packet, from)
}

// /* Initiates cancellation of this round */
// func (r *PhasedRoundController) cancelRound() {
//	panic("not implemented")
// }
//
// /* Initiates cancellation of this round */
// func (r *PhasedRoundController) finishRound() {
//	panic("not implemented")
// }
