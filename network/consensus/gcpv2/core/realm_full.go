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
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/network/consensus/common/cryptkit"
	"github.com/insolar/insolar/network/consensus/common/endpoints"
	"github.com/insolar/insolar/network/consensus/gcpv2/api"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/census"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/member"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/misbehavior"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/phases"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/profiles"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/proofs"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/transport"
	"github.com/insolar/insolar/network/consensus/gcpv2/censusimpl"
	"github.com/insolar/insolar/network/consensus/gcpv2/core/coreapi"
	"github.com/insolar/insolar/network/consensus/gcpv2/core/packetdispatch"
	pop "github.com/insolar/insolar/network/consensus/gcpv2/core/population"
	"github.com/insolar/insolar/network/consensus/gcpv2/core/purgatory"
	"github.com/insolar/insolar/pulse"
)

var _ pulse.DataHolder = &FullRealm{}

type FullRealm struct {
	coreRealm
	// nodeContext pop.Hook

	/* Derived from the ones provided externally - set at init() or start(). Don't need mutex */
	packetBuilder  transport.PacketBuilder
	packetSender   transport.PacketSender
	profileFactory profiles.Factory

	timings api.RoundTimings

	census         census.Active
	population     pop.RealmPopulation
	populationHook *pop.Hook
	purgatory      purgatory.RealmPurgatory

	packetDispatchers []pop.PacketDispatcher

	/* Other fields - need mutex */
	isFinished bool
}

func (r *FullRealm) dispatchPacket(ctx context.Context, packet transport.PacketParser, from endpoints.Inbound,
	verifyFlags coreapi.PacketVerifyFlags) error {

	pt := packet.GetPacketType()

	var sourceNode packetdispatch.MemberPacketReceiver
	var sourceID insolar.ShortNodeID

	switch {
	case pt.GetLimitPerSender() == 0 || int(pt) >= len(r.packetDispatchers) || r.packetDispatchers[pt] == nil:
		return fmt.Errorf("packet type (%v) is unknown", pt)
	case pt.IsMemberPacket():
		selfID := r.GetSelfNodeID()
		strict, err := coreapi.VerifyPacketRoute(ctx, packet, selfID, from)
		if err != nil {
			return err
		}
		if strict {
			verifyFlags |= coreapi.RequireStrictVerify
		}

		sourceID = packet.GetSourceID()
		sourceNode = r.getMemberReceiver(sourceID)
	}

	pd := r.packetDispatchers[pt] // was checked above for != nil

	var err error
	verifyFlags, err = r.coreRealm.VerifyPacketAuthenticity(ctx, packet, from, sourceNode, coreapi.SkipVerify,
		pd, verifyFlags)

	if err != nil {
		return err
	}

	var canHandle bool
	canHandle, err = r.coreRealm.VerifyPacketPulseNumber(ctx, packet, from, r.GetPulseNumber(), r.GetNextPulseNumber(),
		"full:dispatchPacket")

	if !canHandle || err != nil {
		return err
	}

	if sourceNode != nil && !sourceNode.CanReceivePacket(pt) {
		return fmt.Errorf("packet type (%v) limit exceeded: from=%v(%v)", pt, sourceNode.GetNodeID(), from)
	}

	// this enables lazy parsing - packet is fully parsed AFTER validation, hence makes it less prone to exploits for non-members
	packet, err = coreapi.LazyPacketParse(packet)
	if err != nil {
		return err
	}

	if !pt.IsMemberPacket() {
		return pd.DispatchHostPacket(ctx, packet, from, verifyFlags)
	}

	// now it is safe to parse the rest of the packet
	memberPacket := packet.GetMemberPacket()
	if memberPacket == nil {
		return fmt.Errorf("packet type (%v) can't be parsed: from=%v", pt, from)
	}

	if sourceNode == nil {
		memberTriggered := false
		memberTriggered, err = pd.TriggerUnknownMember(ctx, sourceID, memberPacket, from)
		if err != nil {
			return err
		}
		if !memberTriggered {
			return fmt.Errorf("packet type (%v) from unknown sourceID(%v): from=%v", pt, sourceID, from)
		}

		sourceNode = r.getMemberReceiver(sourceID)
		if sourceNode == nil {
			return fmt.Errorf("inconsistent behavior for packet type (%v) from unknown sourceID(%v): from=%v", pt, sourceID, from)
		}
	}

	if !sourceNode.SetPacketReceived(pt) {
		inslogger.FromContext(ctx).Infof("packet type (%v) limit exceeded: from=%v(%v)", pt, sourceNode.GetNodeID(), from)
		return nil
	}

	return sourceNode.DispatchMemberPacket(ctx, packet, from, verifyFlags, pd)
}

/* LOCK - runs under RoundController lock */
func (r *FullRealm) start(census census.Active, population census.OnlinePopulation, bundle PhaseControllersBundle) {
	r.initBasics(census)

	isDynamic := bundle.IsDynamicPopulationRequired()
	perNodeControllers, nodeCallback, startFn := r.initHandlers(isDynamic, population.GetIndexedCapacity(), bundle)

	r.initPopulation(isDynamic, population, perNodeControllers, nodeCallback)
	startFn()
}

func (r *FullRealm) initBefore(transport transport.Factory) transport.NeighbourhoodSizes {

	if r.ephemeralFeeder != nil {
		r.timings = r.ephemeralFeeder.GetEphemeralTimings(r.config)
		r.strategy.AdjustConsensusTimings(&r.timings)
	}

	r.packetSender = transport.GetPacketSender()
	r.packetBuilder = transport.GetPacketBuilder(r.signer)
	return r.packetBuilder.GetNeighbourhoodSize()
}

func (r *FullRealm) initBasics(census census.Active) {

	r.census = census
	r.profileFactory = census.GetProfileFactory(r.assistant)

	if r.ephemeralFeeder == nil { // ephemeral timings are initialized earlier, in initBefore()
		r.timings = r.config.GetConsensusTimings(r.pulseData.NextPulseDelta)
		r.strategy.AdjustConsensusTimings(&r.timings)
	}

	if r.expectedPopulationSize == 0 {
		r.expectedPopulationSize = member.AsIndex(r.config.GetNodeCountHint())
	}
}

func (r *FullRealm) initHandlers(needsDynamic bool, populationCount int,
	bundle PhaseControllersBundle) ([]PerNodePacketDispatcherFactory, NodeUpdateCallback, func()) {

	r.packetDispatchers = make([]pop.PacketDispatcher, phases.PacketTypeCount)

	nodeCount := populationCount
	if needsDynamic && int(r.expectedPopulationSize) > nodeCount {
		nodeCount = r.expectedPopulationSize.AsInt()
	}

	controllers, nodeCallback := bundle.CreateFullPhaseControllers(nodeCount)

	if len(controllers) == 0 {
		panic("no phase controllers")
	}
	individualHandlers := make([]PerNodePacketDispatcherFactory, 0, len(controllers))

	for _, ctl := range controllers {
		for _, pt := range ctl.GetPacketType() {
			if r.packetDispatchers[pt] != nil {
				panic("multiple controllers for packet type")
			}
			pd, nf := ctl.CreatePacketDispatcher(pt, len(individualHandlers), r)
			r.packetDispatchers[pt] = pd
			if nf != nil {
				individualHandlers = append(individualHandlers, nf)
			}
		}
	}

	return individualHandlers, nodeCallback,
		func() {
			for _, ctl := range controllers {
				ctl.BeforeStart(r.roundContext, r)
			}
			for _, ctl := range controllers {
				ctl.StartWorker(r.roundContext, r)
			}
		}
}

func (r *FullRealm) initPopulation(needsDynamic bool, population census.OnlinePopulation,
	individualHandlers []PerNodePacketDispatcherFactory, nodeCallback NodeUpdateCallback) {

	initNodeFn := func(ctx context.Context, n *pop.NodeAppearance) []pop.DispatchMemberPacketFunc {
		if len(individualHandlers) == 0 {
			return nil
		}
		result := make([]pop.DispatchMemberPacketFunc, len(individualHandlers))
		for k, ctl := range individualHandlers {
			ctx, result[k] = ctl.CreatePerNodePacketHandler(ctx, n)
		}
		return result
	}

	var notifyAll func()

	log := inslogger.FromContext(r.roundContext)

	hookCfg := pop.NewSharedNodeContext(r.assistant, r, uint8(r.nbhSizes.NeighbourhoodTrustThreshold), r.getEphemeralMode(),
		func(report misbehavior.Report) interface{} {
			log.Warnf("Got Report: %+v", report)
			r.census.GetMisbehaviorRegistry().AddReport(report)
			return nil
		})

	if needsDynamic {
		expectedSize := r.expectedPopulationSize.AsInt()
		if population.GetIndexedCapacity() > expectedSize {
			expectedSize = population.GetIndexedCapacity()
		}

		popStruct := pop.NewDynamicRealmPopulation(population, expectedSize, r.nbhSizes.ExtendingNeighbourhoodLimit,
			r.strategy.ShuffleNodeSequence, r.strategy.GetBaselineWeightForNeighbours(), hookCfg, initNodeFn)

		r.population = popStruct
		popStruct.InitCallback(nodeCallback)
		r.populationHook = popStruct.GetHook()
		notifyAll = popStruct.NotifyAllOnAdded

		// TODO probably should happen at later stages, closer to Phase3 analysis
		r.population.SealIndexed(expectedSize)
	} else {
		if population.GetIndexedCount() == 0 || !population.IsValid() ||
			population.GetIndexedCount() != population.GetIndexedCapacity() {
			panic("dynamic population is required for joiner or suspect")
		}

		popStruct := pop.NewFixedRealmPopulation(population, r.nbhSizes.ExtendingNeighbourhoodLimit,
			r.strategy.ShuffleNodeSequence, r.strategy.GetBaselineWeightForNeighbours(), hookCfg, initNodeFn)

		r.population = popStruct
		popStruct.InitCallback(nodeCallback)
		r.populationHook = popStruct.GetHook()
		notifyAll = popStruct.NotifyAllOnAdded
	}

	r.initSelf() // should happen before notifications - just in case someone will access GetSelf
	r.purgatory = purgatory.NewRealmPurgatory(r.population, r.profileFactory, r.assistant,
		r.populationHook, r.postponedPacketFn)

	notifyAll()
}

func (r *FullRealm) initSelf() {
	newSelf := r.population.GetSelf()
	prevSelf := r.self

	if newSelf.GetNodeID() != prevSelf.GetNodeID() {
		panic("inconsistent transition of self between realms")
	}

	prevSelf.CopySelfTo(newSelf)
	r.self = newSelf
}

func (r *FullRealm) registerNextJoinCandidate() (*pop.NodeAppearance, cryptkit.DigestHolder) {

	if !r.GetSelf().CanIntroduceJoiner() {
		return nil, nil
	}

	for {
		cp, secret := r.candidateFeeder.PickNextJoinCandidate()
		if cp == nil {
			return nil, nil
		}
		if r.GetPopulation().GetNodeAppearance(cp.GetStaticNodeID()) == nil {
			nip := r.profileFactory.CreateFullIntroProfile(cp)
			sv := r.GetSignatureVerifier(nip.GetPublicKeyStore())
			np := censusimpl.NewJoinerProfile(nip, sv)
			na := pop.NewLocalJoinerNodeAppearance(&np, r.GetSelfNodeID(), secret)

			nna, err := r.population.AddToDynamics(r.roundContext, &na)
			if err != nil {
				inslogger.FromContext(r.roundContext).Error(err)
			} else if nna != nil {
				inslogger.FromContext(r.roundContext).Debugf("Candidate/joiner added as a dynamic node: s=%d, t=%d, full=%v",
					r.GetSelfNodeID(), np.GetNodeID(), np.GetExtension() != nil)

				return nna, secret
			}
		}

		inslogger.FromContext(r.roundContext).Debugf("Candidate/joiner was rejected due to duplicate id: s=%d, t=%d",
			r.GetSelfNodeID(), cp.GetStaticNodeID())

		r.candidateFeeder.RemoveJoinCandidate(false, cp.GetStaticNodeID())
	}
}

func (r *FullRealm) Frauds() misbehavior.FraudFactory {
	return r.populationHook.GetFraudFactory()
}

func (r *FullRealm) Blames() misbehavior.BlameFactory {
	return r.populationHook.GetBlameFactory()
}

func (r *FullRealm) GetSelf() *pop.NodeAppearance {
	return r.self
}

func (r *FullRealm) GetSelfNodeID() insolar.ShortNodeID {
	return r.self.GetNodeID()
}

func (r *FullRealm) IsJoiner() bool {
	return r.self.IsJoiner()
}

func (r *FullRealm) GetPacketSender() transport.PacketSender {
	return r.packetSender
}

func (r *FullRealm) GetPacketBuilder() transport.PacketBuilder {
	return r.packetBuilder
}

func (r *FullRealm) GetSigner() cryptkit.DigestSigner {
	return r.signer
}

func (r *FullRealm) GetPopulation() pop.RealmPopulation {
	return r.population
}

func (r *FullRealm) GetNodeCount() int {
	return r.population.GetIndexedCount()
}

func (r *FullRealm) GetPulseNumber() pulse.Number {
	return r.pulseData.PulseNumber
}

func (r *FullRealm) GetNextPulseNumber() pulse.Number {
	return r.pulseData.GetNextPulseNumber()
}

func (r *FullRealm) GetOriginalPulse() proofs.OriginalPulsarPacket {
	// NB! locks for this field are only needed for PrepRealm
	return r.coreRealm.originalPulse
}

func (r *FullRealm) GetPulseDataDigest() cryptkit.DigestHolder {
	return r.originalPulse.GetPulseDataDigest()
}

func (r *FullRealm) GetPulseData() pulse.Data {
	return r.pulseData
}

func (r *FullRealm) GetLastCloudStateHash() proofs.CloudStateHash {
	return r.census.GetCloudStateHash()
}

func (r *FullRealm) getUpstreamReport() api.UpstreamReport {
	if !r.pulseData.PulseNumber.IsTimePulse() {
		panic("pulse number was not set")
	}

	sp := r.GetSelf().GetProfile()
	return api.UpstreamReport{
		PulseNumber: r.pulseData.PulseNumber,
		MemberPower: sp.GetDeclaredPower(),
		MemberMode:  sp.GetOpMode(),
		IsJoiner:    sp.IsJoiner(),
		// IsEphemeral: false,
	}
}

func (r *FullRealm) PreparePulseChange() (bool, <-chan api.UpstreamState) {
	report := r.getUpstreamReport()

	if r.IsLocalStateful() {
		inslogger.FromContext(r.roundContext).Warnf("PreparePulseChange: self=%s, eph=%v", r.self, r.populationHook.GetEphemeralMode())
		ch := make(chan api.UpstreamState, 1)
		r.stateMachine.PreparePulseChange(report, ch)
		return true, ch
	}

	inslogger.FromContext(r.roundContext).Warnf("PrepareAndCommitStatelessPulseChange: self=%s, eph=%v", r.self, r.populationHook.GetEphemeralMode())
	r.stateMachine.CommitPulseChangeByStateless(report, r.pulseData, r.census)
	return false, nil
}

func (r *FullRealm) CommitPulseChange() {
	report := r.getUpstreamReport()
	inslogger.FromContext(r.roundContext).Warnf("CommitPulseChange: self=%s", r.self)

	r.stateMachine.CommitPulseChange(report, r.pulseData, r.census)
}

func (r *FullRealm) GetTimings() api.RoundTimings {
	return r.timings
}

func (r *FullRealm) GetNeighbourhoodSizes() transport.NeighbourhoodSizes {
	return r.nbhSizes
}

func (r *FullRealm) GetLocalProfile() profiles.LocalNode {
	return r.self.GetProfile().(profiles.LocalNode)
}

func (r *FullRealm) IsLocalStateful() bool {
	return r.self.IsStateful()
}

func (r *FullRealm) ApplyLocalState(nsh proofs.NodeStateHash) bool {

	if (nsh == nil) == r.IsLocalStateful() {
		panic("illegal value")
	}

	mp := r.self.GetNodeMembershipProfileOrEmpty()
	ma := r.buildLocalMemberAnnouncementDraft(mp)

	if nsh != nil {
		v := nsh.SignWith(r.signer)
		ma.Membership.StateEvidence = v
		ma.Membership.AnnounceSignature = v.GetSignatureHolder()
	} else {
		v := r.self.GetStatelessAnnouncementEvidence()
		// v := nsh.SignWith(r.signer)
		ma.Membership.StateEvidence = v
		ma.Membership.AnnounceSignature = v.GetSignatureHolder()
	}

	// TODO use r.GetLastCloudStateHash() + digest(PulseData) + r.digest.GetGshDigester() to build digest for signing

	// TODO Hack! MUST provide announcement hash

	return r.self.SetLocalNodeState(ma)
}

func (r *FullRealm) buildLocalMemberAnnouncementDraft(mp profiles.MembershipProfile) profiles.MemberAnnouncement {

	lp := r.self.GetProfile()

	if lp.IsJoiner() {
		return profiles.NewJoinerAnnouncement(lp.GetStatic(), insolar.AbsentShortNodeID)
	}

	localID := lp.GetNodeID()
	if isLeave, leaveReason := r.controlFeeder.GetRequiredGracefulLeave(); isLeave {
		return profiles.NewMemberAnnouncementWithLeave(localID, mp, leaveReason, insolar.AbsentShortNodeID)
	}

	r.self.CanIntroduceJoiner()
	if !r.unsafeRound && lp.CanIntroduceJoiner() {
		jc, secret := r.registerNextJoinCandidate()
		if jc != nil {
			return profiles.NewMemberAnnouncementWithJoinerID(localID, mp, jc.GetNodeID(), secret, localID)
		}
	}

	return profiles.NewMemberAnnouncement(localID, mp, insolar.AbsentShortNodeID)
}

func (r *FullRealm) CreateAnnouncement(n *pop.NodeAppearance, isJoinerProfileRequired bool) *transport.NodeAnnouncementProfile {
	ma := n.GetRequestedAnnouncement()
	if ma.Membership.IsEmpty() {
		panic("illegal state")
	}

	var joiner *transport.JoinerAnnouncement
	if !ma.JoinerID.IsAbsent() && isJoinerProfileRequired {
		joiner = r.GetPurgatory().GetJoinerAnnouncement(ma.JoinerID, n.GetNodeID())
		switch {
		case joiner != nil:
			break
		case n == r.self:
			panic(fmt.Sprintf("illegal state - local joiner is missing: %d", ma.JoinerID))
		default:
			panic(fmt.Sprintf("illegal state - joiner is missing: s=%d n=%d j=%d",
				r.self.GetNodeID(), n.GetNodeID(), ma.JoinerID))
		}
	} else if ma.Membership.IsJoiner() {
		joiner = n.GetAnnouncementAsJoiner()
		if joiner == nil {
			panic("illegal state")
		}
	}

	return transport.NewNodeAnnouncement(n.GetProfile(), ma, r.GetNodeCount(), r.pulseData.PulseNumber, joiner)
}

func (r *FullRealm) CreateLocalAnnouncement() *transport.NodeAnnouncementProfile {
	return r.CreateAnnouncement(r.self, true)
}

func (r *FullRealm) CreateLocalPhase0Announcement() *transport.NodeAnnouncementProfile {
	ma := r.self.GetRequestedAnnouncement()
	return transport.NewNodeAnnouncement(r.self.GetProfile(), ma, r.GetNodeCount(), r.pulseData.PulseNumber, nil)
}

func (r *FullRealm) finishRound(ctx context.Context, builder census.Builder, csh proofs.CloudStateHash) {
	r.Lock()
	defer r.Unlock()

	if r.isFinished {
		panic("illegal state")
	}
	r.isFinished = true

	pb := builder.GetPopulationBuilder()
	local := pb.GetLocalProfile()

	var expected census.Expected
	if csh != nil {
		expected = builder.Build(csh).MakeExpected()
	} else {
		expected = builder.BuildAsBroken(csh).MakeExpected()
	}

	if expected.GetOnlinePopulation().GetLocalProfile().IsJoiner() {
		panic("DEBUG FAIL-FAST: local remains as joiner")
	}

	if expected.GetOnlinePopulation().GetLocalProfile().GetOpMode().IsMistrustful() {
		panic("DEBUG FAIL-FAST: local was marked as fraud suspect")
	}

	isNextEphemeral := false
	if r.ephemeralFeeder != nil {
		if r.ephemeralFeeder.CanStopEphemeralByCensus(expected) {
			expected = expected.Rebuild(pulse.Unknown).MakeExpected()
		} else {
			isNextEphemeral = true
		}
	}

	newSelf := expected.GetOnlinePopulation().GetLocalProfile()
	report := api.UpstreamReport{
		PulseNumber: r.pulseData.PulseNumber,
		MemberPower: newSelf.GetDeclaredPower(),
		MemberMode:  newSelf.GetOpMode(),
	}

	// MUST happen before ephemeralFeeder to trigger proper worker stare
	r.stateMachine.ConsensusFinished(report, expected)

	if r.ephemeralFeeder != nil {
		r.ephemeralFeeder.EphemeralConsensusFinished(isNextEphemeral, r.roundStartedAt, expected)
		if !isNextEphemeral {
			r.ephemeralFeeder.OnEphemeralCancelled()
			r.ephemeralFeeder = nil
		}
	}

	nextNP := expected.GetPulseNumber()
	rs := r.self.GetRequestedState()
	if expected.GetOnlinePopulation().IsValid() {
		switch {
		case rs.IsLeaving:
			r.controlFeeder.OnAppliedGracefulLeave(rs.LeaveReason, nextNP)
		case !rs.JoinerID.IsAbsent():
			r.candidateFeeder.RemoveJoinCandidate(true, rs.JoinerID)
		}
	} else {
		inslogger.FromContext(ctx).Debugf("got a broken population: s=%d %v", local.GetNodeID(), expected.GetOnlinePopulation())
	}
	pw := rs.RequestedPower
	if !local.IsPowered() {
		pw = 0
	}
	r.controlFeeder.OnAppliedMembershipProfile(local.GetOpMode(), pw, nextNP)
}

func (r *FullRealm) GetProfileFactory() profiles.Factory {
	return r.profileFactory
}

func (r *FullRealm) GetPurgatory() *purgatory.RealmPurgatory {
	return &r.purgatory
}

func (r *FullRealm) getMemberReceiver(id insolar.ShortNodeID) packetdispatch.MemberPacketReceiver {
	// Purgatory MUST be checked first to avoid "missing" a node during its transition from the purgatory to normal population
	pn := r.GetPurgatory().GetPhantomNode(id)
	if pn != nil {
		return pn
	}
	na := r.GetPopulation().GetNodeAppearance(id)
	if na != nil {
		return na
	}
	return nil
}

func (r *FullRealm) GetWelcomePackage() *proofs.NodeWelcomePackage {
	return &proofs.NodeWelcomePackage{
		CloudIdentity:      r.census.GetMandateRegistry().GetCloudIdentity(),
		LastCloudStateHash: r.census.GetCloudStateHash(),
	}
}

func (r *FullRealm) BuildNextPopulation(ctx context.Context, ranks []profiles.PopulationRank,
	gsh proofs.GlobulaStateHash, csh proofs.CloudStateHash) bool {

	b := r.census.CreateBuilder(ctx, r.GetNextPulseNumber())
	pb := b.GetPopulationBuilder()
	selfID := r.GetSelfNodeID()
	localUpdProfile := pb.GetLocalProfile()
	selfMode := member.ModeEvictedGracefully

	idx := member.AsIndex(0)
	for _, pr := range ranks {
		prevAP := pr.Profile
		nodeID := prevAP.GetNodeID()

		nextAP := localUpdProfile

		if nodeID == selfID {
			selfMode = pr.OpMode
		} else {
			static := prevAP.GetStatic()
			static, _ = r.profileFactory.TryConvertUpgradableIntroProfile(static)
			nextAP = pb.AddProfile(static)
		}
		if pr.OpMode.IsPowerless() && pr.Power != 0 {
			panic("illegal state")
		}

		nextAP.SetSignatureVerifier(prevAP.GetSignatureVerifier())
		if pr.OpMode.IsEvictedGracefully() {
			na := r.self
			if nodeID != selfID {
				na = r.population.GetNodeAppearance(nodeID)
			}
			leave, leaveReason := na.GetRequestedLeave()
			if !leave {
				panic("illegal state")
			}
			nextAP.SetOpModeAndLeaveReason(idx, leaveReason)
		} else {
			nextAP.SetRank(idx, pr.OpMode, pr.Power)
		}
		idx++
	}

	b.SetGlobulaStateHash(gsh)
	b.SealCensus()

	r.finishRound(ctx, b, csh)

	if selfMode.IsEvicted() {
		inslogger.FromContext(ctx).Info("Node has left")
		return false
	}
	return true
}

func (r *FullRealm) MonitorOtherPulses(packet transport.PulsePacketReader, from endpoints.Inbound) error {

	return r.Blames().NewMismatchedPulsarPacket(from, r.GetOriginalPulse(), packet.GetPulseDataEvidence())
}

func (r *FullRealm) NotifyRoundStopped(ctx context.Context) {
	r.stateMachine.OnRoundStopped(ctx)
}
