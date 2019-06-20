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
	"sync"
	"time"

	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/network/consensus/common"
	"github.com/insolar/insolar/network/consensus/gcpv2/census"
	errors2 "github.com/insolar/insolar/network/consensus/gcpv2/errors"
	"github.com/insolar/insolar/network/consensus/gcpv2/packets"
)

type PhasedRoundController struct {
	rw sync.RWMutex

	/* Derived from the provided externally - set at init() or start(). Don't need mutex */
	fullCancel    context.CancelFunc /* cancels prepareCancel as well */
	prepareCancel context.CancelFunc

	/* Other fields - need mutex */
	isRunning bool
	prepR     *PrepRealm
	realm     FullRealm
}

func NewPhasedRoundController(strategy RoundStrategy, chronicle census.ConsensusChronicles, transport TransportFactory,
	config LocalNodeConfiguration) *PhasedRoundController {

	r := &PhasedRoundController{}

	r.realm.coreRealm = newCoreRealm(&r.rw)
	r.realm.coreRealm.errorFactory = errors2.NewMisbehaviorFactories(r.realm.coreRealm.captureMisbehavior)
	r.realm.coreRealm.strategy = strategy
	r.realm.coreRealm.config = config
	r.realm.coreRealm.chronicle = chronicle
	r.realm.coreRealm.packetSender = transport.GetPacketSender()
	r.realm.coreRealm.initialCensus = chronicle.GetLatestCensus()

	crypto := transport.GetCryptographyFactory()
	r.realm.coreRealm.verifierFactory = crypto
	r.realm.coreRealm.digest = crypto.GetDigestFactory()
	sks := config.GetSecretKeyStore()
	r.realm.coreRealm.signer = crypto.GetNodeSigner(sks)
	r.realm.coreRealm.packetBuilder = transport.GetPacketBuilder(r.realm.coreRealm.signer)

	population := r.realm.coreRealm.initialCensus.GetOnlinePopulation()
	r.realm.coreRealm.self = NewNodeAppearanceAsSelf(population.GetLocalProfile())

	return r
}

func (r *PhasedRoundController) StartConsensusRound(upstream UpstreamPulseController) {
	r.rw.Lock()
	defer r.rw.Unlock()

	if r.fullCancel != nil {
		panic("was started once")
	}
	r.isRunning = true

	r.realm.roundStartedAt = time.Now()
	r.realm.coreRealm.upstream = upstream

	ctx := r.realm.config.GetParentContext()
	ctx, r.fullCancel = context.WithCancel(ctx)

	r.realm.roundContext = r.realm.strategy.CreateRoundContext(ctx)
	r.realm.logger = inslogger.FromContext(r.realm.roundContext)
	ctx, r.prepareCancel = context.WithCancel(r.realm.roundContext)

	preps := r.realm.strategy.GetPrepPhaseControllers()

	if len(preps) > 0 {
		r.prepR = &PrepRealm{
			coreRealm:         &r.realm.coreRealm,
			completeFn:        r.finishPreparation,
			postponedPacketFn: r.handlePostponedPacket,
		}
		r.prepR.start(ctx, preps,
			10000 /* Should be excessively enough to avoid lockups */)
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

	if r.fullCancel == nil || !r.isRunning {
		return // false
	}
	r.isRunning = false
	r.fullCancel()
	//return true
}

/* LOCK: simple */
func (r *PhasedRoundController) IsRunning() bool {
	r.rw.RLock()
	defer r.rw.RUnlock()
	return r.isRunning
}

/* Checks if Controller can handle a new packet, and which Realm should do it. If result = nil, then FullRealm is used */
/* LOCK: simple */
func (r *PhasedRoundController) beforeHandlePacket() (*PrepRealm, common.PulseNumber, error) {
	r.rw.RLock()
	defer r.rw.RUnlock()

	if r.fullCancel == nil {
		return nil, 0, errors2.NewRoundStateError("not started")
	}
	if !r.isRunning {
		return nil, 0, errors2.NewRoundStateError("stopped")
	}

	if r.prepR == nil {
		return nil, r.realm.GetPulseNumber(), nil
	} else {
		return r.prepR, r.realm.initialCensus.GetExpectedPulseNumber(), nil
	}
}

/*
LOCK - must be called under LOCK
Removes PrepRealm and starts FullRealm. Should only be used from PrepRealm.
Nodes and PulseData(?) must be available.
*/
func (r *PhasedRoundController) finishPreparation(successful bool) {
	if r.prepR == nil {
		return
	}
	prep := r.prepR
	r.prepR = nil

	r.startFullRealm()

	prep.stop() // initiates handover from PrepRealm
}

func (r *PhasedRoundController) startFullRealm() {

	chronicle := r.realm.chronicle
	lastCensus := chronicle.GetLatestCensus()
	pd := &r.realm.pulseData

	if lastCensus.IsActive() && lastCensus.GetPulseNumber().IsUnknown() {
		/* This is the priming lastCensus */
		b := chronicle.GetActiveCensus().CreateBuilder(pd.PulseNumber)
		r.prepareNewMembers(b.GetOnlinePopulationBuilder())
		b.BuildAndMakeExpected(lastCensus.GetMandateRegistry().GetPrimingCloudHash())
		chronicle.GetExpectedCensus().MakeActive(*pd)
	} else {
		if lastCensus.GetPulseNumber() != pd.PulseNumber {
			panic("illegal state - pulse number of expected census and of the realm are mismatched")
		}
		if !lastCensus.IsActive() {
			/* Auto-activation of the prepared lastCensus */
			expCensus := chronicle.GetExpectedCensus()
			lastCensus = expCensus.MakeActive(*pd)
		}
	}

	r.realm.start()
}

func (r *PhasedRoundController) handlePostponedPacket(packet packets.PacketParser, from common.HostIdentityHolder) {
	// NB! we may need to handle errors from delayed packets
	_ = r.handlePacket(packet, from, true)
}

func (r *PhasedRoundController) HandlePacket(packet packets.PacketParser, from common.HostIdentityHolder) error {
	return r.handlePacket(packet, from, false)
}

func (r *PhasedRoundController) handlePacket(packet packets.PacketParser, from common.HostIdentityHolder, preVerified bool) error {

	/* a separate method with lock is to ensure that further packet processing is not connected to a lock */
	prep, filterPN, err := r.beforeHandlePacket()
	if err != nil {
		return err
	}

	if !filterPN.IsUnknown() {
		pn := packet.GetPulseNumber()
		if !pn.IsUnknown() && filterPN != pn {
			return errors2.NewPulseRoundMismatchError(pn,
				fmt.Sprintf("packet pulse number mismatched: expected=%v, actual=%v", filterPN, pn))
		}
	}

	pt := packet.GetPacketType()
	if pt.IsMemberPacket() {
		memberPacket := packet.GetMemberPacket()
		if memberPacket == nil {
			panic("missing parser for phased packet")
		}
		selfId := r.realm.coreRealm.GetSelfNodeId()
		sid := memberPacket.GetSourceShortNodeId()
		if sid == selfId {
			return fmt.Errorf("loopback, source ShortNodeID(%v) == this ShortNodeID(%v)", sid, selfId)
		}
		if memberPacket.HasTargetShortNodeId() {
			tid := memberPacket.GetTargetShortNodeId()
			if tid != selfId {
				return fmt.Errorf("target ShortNodeID(%v) != this ShortNodeID(%v)", tid, selfId)
			}
		}

		if prep == nil { // Full realm is active - we can use node projections
			src, err := r.realm.GetNodeApperance(sid)
			if err != nil {
				return err
			}
			err = src.VerifyPacketAuthenticity(packet, from, preVerified)
			if err != nil {
				return err
			}
			route := &r.realm.handlers[pt]
			if route.HasMemberHandler() {
				return route.handleMemberPacket(memberPacket, src)
			}
			return route.handleHostPacket(packet, from)
		}
	}

	if !preVerified {
		err = r.realm.coreRealm.VerifyPacketAuthenticity(packet, from)
		if err != nil {
			return err
		}
	}

	if prep != nil { // Prep realm is active
		h := prep.handlers[pt]
		var explicitPostpone = false
		if h != nil {
			explicitPostpone, err = h(packet, from)
			if !explicitPostpone || err != nil {
				return err
			}
		}
		// if packet is not handled, then we may need to leave it for FullRealm
		if prep.PostponePacket(packet, from) {
			return nil
		}
		if explicitPostpone {
			return fmt.Errorf("unable to postpone packet explicitly: type=%v", pt)
		} else {
			return errPacketIsNotAllowed
		}
	} else {
		return r.realm.handlers[pt].handleHostPacket(packet, from)
	}
}

/* Initiates cancellation of this round */
func (r *PhasedRoundController) cancelRound() {
	panic("not implemented")
}

/* Initiates cancellation of this round */
func (r *PhasedRoundController) finishRound() {
	panic("not implemented")
}

func (r *PhasedRoundController) prepareNewMembers(pop census.OnlinePopulationBuilder) {

	for _, p := range pop.GetUnorderedProfiles() {
		if p.GetSignatureVerifier() != nil {
			continue
		}
		v := r.realm.GetSignatureVerifier(p.GetNodePublicKeyStore())
		p.SetSignatureVerifier(v)
	}
}
