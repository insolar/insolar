// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package core

import (
	"context"
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/network/consensus/common/cryptkit"
	"github.com/insolar/insolar/network/consensus/common/endpoints"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/census"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/member"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/phases"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/proofs"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/transport"
	"github.com/insolar/insolar/network/consensus/gcpv2/core/coreapi"
	"github.com/insolar/insolar/network/consensus/gcpv2/core/errors"
	"github.com/insolar/insolar/network/consensus/gcpv2/core/packetdispatch"
	"github.com/insolar/insolar/network/consensus/gcpv2/core/population"
	"github.com/insolar/insolar/pulse"
)

/*
	PrepRealm is a functionally limited and temporary realm that is used when this node doesn't know pulse or last consensus.
	It can ONLY pre-processed packets, but is disallowed to send them.

	Pre-processed packets as postponed by default and processing will be repeated when FullRealm is activated.
*/
type PrepRealm struct {
	/* Provided externally. Don't need mutex */
	*coreRealm // points the core part realms, it is shared between of all Realms of a Round

	completeFn func(successful bool, startedAt time.Time) // MUST be called under lock, consequent calls are ignored

	/* Derived from the provided externally - set at init() or start(). Don't need mutex */
	packetDispatchers []population.PacketDispatcher
	packetRecorder    packetdispatch.PacketRecorder
	// queueToFull       chan packetrecorder.PostponedPacket
	// phase2ExtLimit    uint8

	/* Other fields - need mutex */

	limiters           sync.Map
	lastCloudStateHash cryptkit.DigestHolder
	disableEphemeral   bool                       // blocks polling
	prepSelf           *population.NodeAppearance /* local copy to avoid race */
}

func (p *PrepRealm) init(completeFn func(successful bool, startedAt time.Time)) {
	p.completeFn = completeFn
	if p.coreRealm.self == nil {
		panic("illegal state")
	}
	p.prepSelf = p.coreRealm.self
}

func (p *PrepRealm) dispatchPacket(ctx context.Context, packet transport.PacketParser, from endpoints.Inbound,
	verifyFlags coreapi.PacketVerifyFlags) error {

	pt := packet.GetPacketType()
	selfID := p.prepSelf.GetNodeID()

	var limiterKey string
	switch {
	case pt.GetLimitPerSender() == 0:
		return errors.UnknownPacketType(pt)
	case pt.IsMemberPacket():
		strict, err := coreapi.VerifyPacketRoute(ctx, packet, selfID, from)
		if err != nil {
			return err
		}
		if strict {
			verifyFlags |= coreapi.RequireStrictVerify
		}
		limiterKey = endpoints.ShortNodeIDAsByteString(packet.GetSourceID())
	default:
		limiterKey = string(from.AsByteString())
	}

	/*
		We use limiter here explicitly to ensure that the node's postpone queue can't be overflown during PrepPhase
	*/
	limiter := phases.NewAtomicPacketLimiter(phases.NewPacketLimiter(p.nbhSizes.ExtendingNeighbourhoodLimit))
	{
		limiterI, _ := p.limiters.LoadOrStore(limiterKey, limiter)
		limiter = limiterI.(*phases.AtomicPacketLimiter)
	}

	// if !limiter.GetPacketLimiter().CanReceivePacket(pt) {
	//	return fmt.Errorf("packet type (%v) limit exceeded: from=%v", pt, from)
	// }

	var pd population.PacketDispatcher

	if int(pt) < len(p.packetDispatchers) {
		pd = p.packetDispatchers[pt]
	}

	var err error
	verifyFlags, err = p.coreRealm.VerifyPacketAuthenticity(ctx, packet, from, nil, coreapi.DefaultVerify,
		pd, verifyFlags)

	if err != nil {
		return err
	}

	var canHandle bool
	canHandle, err = p.coreRealm.VerifyPacketPulseNumber(ctx, packet, from, p.initialCensus.GetExpectedPulseNumber(), 0,
		"prep:dispatchPacket")

	if !canHandle || err != nil {
		return err
	}

	if !limiter.SetPacketReceived(pt) {
		return errors.LimitExceeded(pt, packet.GetSourceID(), from)
	}

	if pd != nil {
		// this enables lazy parsing - packet is fully parsed AFTER validation, hence makes it less prone to exploits for non-members
		packet, err = coreapi.LazyPacketParse(packet)
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
	p.packetRecorder = packetdispatch.NewPacketRecorder(int(packetsPerSender) * 100)

	p.packetDispatchers = make([]population.PacketDispatcher, phases.PacketTypeCount)
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

// runs under lock
func (p *PrepRealm) _startWorkers(ctx context.Context, controllers []PrepPhaseController) {

	if p.originalPulse != nil {
		// we were set for FullRealm, so skip prep workers
		return
	}

	for _, ctl := range controllers {
		ctl.StartWorker(ctx, p)
	}
}

func (p *PrepRealm) prepareEphemeralPolling(ctxPrep context.Context) {
	if p.ephemeralFeeder == nil || !p.ephemeralFeeder.IsActive() {
		return
	}

	minDuration := p.ephemeralFeeder.GetMinDuration()
	beforeNextRound := p.ephemeralFeeder.GetMaxDuration()

	var startTimer *time.Timer
	var startCh <-chan time.Time

	pop := p.initialCensus.GetOnlinePopulation()
	local := pop.GetLocalProfile()
	if pop.GetIndexedCount() < 2 || local.IsJoiner() || !local.GetStatic().GetSpecialRoles().IsDiscovery() {
		beforeNextRound = 0
	}

	if beforeNextRound > 0 && beforeNextRound < math.MaxInt64 {
		if beforeNextRound < minDuration {
			beforeNextRound = minDuration
		}

		if beforeNextRound < time.Second {
			beforeNextRound = time.Second
		}

		startTimer = time.NewTimer(beforeNextRound)
		startCh = startTimer.C
	}

	p.AddPoll(func(ctxOfPolling context.Context) bool {
		select {
		case <-ctxOfPolling.Done():
		case <-ctxPrep.Done():
		default:
			select {
			case <-startCh:
				go p.pushEphemeralPulse(ctxPrep)
			default:
				if !p.checkEphemeralStartByCandidate(ctxPrep) {
					return true // stay in polling
				}
				go p.pushEphemeralPulse(ctxPrep)
				// stop polling anyway - repeating of unsuccessful is bad
			}
		}
		if startTimer != nil {
			startTimer.Stop()
		}
		return false
	})
}

func (p *PrepRealm) pushEphemeralPulse(ctx context.Context) {

	p.Lock()
	defer p.Unlock()

	if p.disableEphemeral || p.ephemeralFeeder == nil {
		return // ephemeral mode was deactivated
	}

	pde := p.ephemeralFeeder.CreateEphemeralPulsePacket(p.initialCensus)
	ok, pn := p._applyPulseData(ctx, time.Now(), pde, false)
	if !ok && pn != pde.GetPulseNumber() {
		inslogger.FromContext(ctx).Error("active ephemeral start has failed, going to passive")
	}
}

func (p *PrepRealm) checkEphemeralStartByCandidate(ctx context.Context) bool {
	jc, _ := p.candidateFeeder.PickNextJoinCandidate()
	if jc != nil {
		inslogger.FromContext(ctx).Debug("ephemeral polling has found a candidate: ", jc)
		return true
	}
	return false
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

func (p *PrepRealm) ApplyPulseData(ctx context.Context, startedAt time.Time, pp transport.PulsePacketReader, fromPulsar bool, from endpoints.Inbound) error {

	pde := pp.GetPulseDataEvidence()
	pd := pp.GetPulseData()
	pn := pd.PulseNumber
	if pde.GetPulseData() != pd {
		return fmt.Errorf("pulse data and pulse data evidence are mismatched: %v, %v", pd, pde)
	}
	if pd.IsEmpty() {
		return fmt.Errorf("pulse data is empty: %v", pd)
	}

	p.Lock()
	defer p.Unlock()

	ok, epn := p._applyPulseData(ctx, startedAt, pde, fromPulsar)
	if ok || !epn.IsUnknown() && epn == pn {
		return nil
	}

	// TODO blame pulsar and/or node
	localID := p.self.GetNodeID()

	return errors.NewPulseRoundMismatchErrorDef(pn, epn, localID, from, "prep:ApplyPulseData")
}

func (p *PrepRealm) _applyPulseData(_ context.Context, startedAt time.Time, pdp proofs.OriginalPulsarPacket, fromPulsar bool) (bool, pulse.Number) {

	pd := pdp.GetPulseData()

	valid := false
	switch {
	case p.originalPulse != nil:
		return false, p.pulseData.PulseNumber // got something already
	case fromPulsar || p.ephemeralFeeder == nil:
		// Pulsars are NEVER ALLOWED to send ephemeral pulses
		valid = pd.IsValidPulsarData()
	default:
		valid = pd.IsValidPulseData()
	}

	if !valid {
		return false, pulse.Unknown // TODO improve logging on mismatch cases
	}

	switch {
	case p.ephemeralFeeder != nil && pd.IsFromPulsar():
		if fromPulsar { // we cant receive pulsar packets directly from pulsars when ephemeral
			panic("illegal state")
		}
		if !p.ephemeralFeeder.CanStopEphemeralByPulse(pd, p.prepSelf.GetProfile()) {
			return false, pulse.Unknown
		}
		p.disableEphemeral = true
		// real pulse can't be validated vs ephemeral pulse
	default:
		epn := pulse.Unknown
		if p.initialCensus.IsActive() {
			epn = p.initialCensus.GetExpectedPulseNumber()
		} else {
			epn = p.initialCensus.GetPulseNumber()
		}

		if !epn.IsUnknownOrEqualTo(pd.PulseNumber) {
			return false, epn
		}
	}

	if p.originalPulse != nil || !p.pulseData.IsEmpty() {
		return false, pd.PulseNumber
	}

	p.originalPulse = pdp
	p.pulseData = pd

	p.completeFn(true, startedAt)

	return true, pd.PulseNumber
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
