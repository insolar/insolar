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

	"github.com/insolar/insolar/network/consensus/common/cryptkit"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/census"
	"github.com/insolar/insolar/network/consensus/gcpv2/core/errors"

	"github.com/insolar/insolar/network/consensus/gcpv2/api/member"
	"github.com/insolar/insolar/network/consensus/gcpv2/core/packetrecorder"

	"github.com/insolar/insolar/network/consensus/common/endpoints"
	"github.com/insolar/insolar/network/consensus/common/pulse"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/phases"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/proofs"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/transport"
)

/*
	PrepRealm is a functionally limited and temporary realm that is used when this node doesn't know pulse or last consensus.
	It can ONLY pre-processed packets, but is disallowed to send them.

	Pre-processed packets as postponed by default and processing will be repeated when FullRealm is activated.
*/
type PrepRealm struct {
	/* Provided externally. Don't need mutex */
	*coreRealm                                    // points the core part realms, it is shared between of all Realms of a Round
	completeFn              func(successful bool) //
	isEphemeralPulseAllowed bool

	/* Derived from the provided externally - set at init() or start(). Don't need mutex */
	packetDispatchers []PacketDispatcher
	packetRecorder    packetrecorder.PacketRecorder
	// queueToFull       chan packetrecorder.PostponedPacket
	// phase2ExtLimit    uint8

	limiters           sync.Map
	lastCloudStateHash cryptkit.DigestHolder

	/* Other fields - need mutex */
	// 	censusBuilder census.Builder
}

func (p *PrepRealm) init(isEphemeralPulseAllowed bool, completeFn func(successful bool)) {
	p.isEphemeralPulseAllowed = isEphemeralPulseAllowed
	p.completeFn = completeFn
}

func (p *PrepRealm) dispatchPacket(ctx context.Context, packet transport.PacketParser, from endpoints.Inbound,
	verifyFlags packetrecorder.PacketVerifyFlags) error {

	pt := packet.GetPacketType()
	selfID := p.GetSelfNodeID()

	var limiterKey string
	switch {
	case pt.GetLimitPerSender() == 0:
		return fmt.Errorf("packet type (%v) is unknown", pt)
	case pt.IsMemberPacket():
		strict, err := VerifyPacketRoute(ctx, packet, selfID)
		if err != nil {
			return err
		}
		if strict {
			verifyFlags = packetrecorder.RequireStrictVerify
		}
		limiterKey = endpoints.ShortNodeIDAsByteString(packet.GetSourceID())
	default:
		limiterKey = from.AsByteString()

		// TODO HACK - network doesnt have information about pulsars to validate packets, hackIgnoreVerification must be removed when fixed
		verifyFlags = packetrecorder.SkipVerify
	}

	/*
		We use limiter here explicitly to ensure that the node's postpone queue can't be overflown during PrepPhase
	*/
	limiter := phases.NewAtomicPacketLimiter(phases.NewPacketLimiter(p.nbhSizes.ExtendingNeighbourhoodLimit))
	{
		limiterI, _ := p.limiters.LoadOrStore(limiterKey, limiter)
		limiter = limiterI.(*phases.AtomicPacketLimiter)
	}

	if !limiter.GetPacketLimiter().CanReceivePacket(pt) {
		return fmt.Errorf("packet type (%v) limit exceeded: from=%v", pt, from)
	}

	var pd PacketDispatcher

	if int(pt) < len(p.packetDispatchers) {
		pd = p.packetDispatchers[pt]
	}

	if verifyFlags&(packetrecorder.SkipVerify|packetrecorder.SuccesfullyVerified) == 0 {
		strict := verifyFlags&packetrecorder.RequireStrictVerify != 0

		if pd == nil || !pd.HasCustomVerifyForHost(from, strict) {
			sourceID := packet.GetSourceID()

			err := p.coreRealm.VerifyPacketAuthenticity(packet.GetPacketSignature(), sourceID, from, strict)

			if err != nil {
				return err
			}
			verifyFlags |= packetrecorder.SuccesfullyVerified
		}
	}

	if !limiter.SetPacketReceived(pt) {
		return fmt.Errorf("packet type (%v) limit exceeded: from=%v", pt, from)
	}

	if pd != nil {
		// this enables lazy parsing - packet is fully parsed AFTER validation, hence makes it less prone to exploits for non-members
		var err error
		packet, err = LazyPacketParse(packet)
		if err != nil {
			return err
		}

		err = pd.DispatchHostPacket(ctx, packet, from, verifyFlags)
		if err != nil {
			// TODO an error to ignore postpone?
			return err
		}
	}

	p.packetRecorder.Record(packet, from, verifyFlags)
	return nil
}

/* LOCK - runs under RoundController lock */
func (p *PrepRealm) beforeStart(ctx context.Context, controllers []PrepPhaseController) {

	if p.postponedPacketFn == nil {
		panic("illegal state")
	}
	limiter := phases.NewPacketLimiter(p.nbhSizes.ExtendingNeighbourhoodLimit)
	packetsPerSender := limiter.GetRemainingPacketCountDefault()
	p.packetRecorder = packetrecorder.NewPacketRecorder(int(packetsPerSender) * 100)

	p.packetDispatchers = make([]PacketDispatcher, phases.PacketTypeCount)
	for _, ctl := range controllers {
		for _, pt := range ctl.GetPacketType() {
			if p.packetDispatchers[pt] != nil {
				panic("multiple controllers for packet type")
			}
			p.packetDispatchers[pt] = ctl.CreatePacketDispatcher(pt, p)
		}
	}

	for _, ctl := range controllers {
		ctl.BeforeStart(ctx, p)
	}
}

func (p *PrepRealm) startWorkers(ctx context.Context, controllers []PrepPhaseController) {
	for _, ctl := range controllers {
		ctl.StartWorker(ctx, p)
	}
}

func (p *PrepRealm) stop() {
	p.packetRecorder.Playback(p.postponedPacketFn)
}

func (p *PrepRealm) GetOriginalPulse() proofs.OriginalPulsarPacket {
	p.RLock()
	defer p.RUnlock()

	// locks are only needed for PrepRealm
	return p.coreRealm.originalPulse
}

func (p *PrepRealm) GetMandateRegistry() census.MandateRegistry {
	return p.initialCensus.GetMandateRegistry()
}

func (p *PrepRealm) ApplyPulseData(pp transport.PulsePacketReader, fromPulsar bool, from endpoints.Inbound) error {
	pd := pp.GetPulseData()

	p.Lock()
	defer p.Unlock()

	valid := false
	switch {
	case p.originalPulse != nil:
		if pd == p.pulseData {
			return nil // got it already
		}
	case fromPulsar || !p.isEphemeralPulseAllowed:
		// Pulsars are NEVER ALLOWED to send ephemeral pulses
		valid = pd.IsValidPulsarData()
	default:
		valid = pd.IsValidPulseData()
	}
	if !valid {
		// if fromPulsar
		// TODO blame pulsar and/or node
		return fmt.Errorf("invalid pulse data")
	}

	epn := pulse.Unknown
	if p.initialCensus.GetCensusState() == census.PrimingCensus || p.initialCensus.IsActive() {
		epn = p.initialCensus.GetExpectedPulseNumber()
	} else {
		epn = p.initialCensus.GetPulseNumber()
	}

	//	sourceID := packet.GetSourceID()
	localID := p.self.GetNodeID()

	pn := pd.PulseNumber
	if !epn.IsUnknown() && epn != pn {
		return errors.NewPulseRoundMismatchError(pn,
			fmt.Sprintf("packet pulse number mismatched: expected=%v, actual=%v, local=%d, from=%v",
				epn, pn, localID, from))
	}

	//if p.IsJoiner() && p.lastCloudStateHash {
	//
	//}

	p.originalPulse = pp.GetPulseDataEvidence()
	p.pulseData = pd

	p.completeFn(true)

	return nil
}

func (p *PrepRealm) ApplyCloudIntro(lastCloudStateHash cryptkit.DigestHolder, populationCount int, from endpoints.Inbound) {

	p.Lock()
	defer p.Unlock()

	popCount := member.AsIndex(populationCount)
	if p.expectedPopulationSize < popCount {
		p.expectedPopulationSize = popCount
	}

	p.lastCloudStateHash = lastCloudStateHash
}
