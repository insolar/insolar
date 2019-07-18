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
	"github.com/insolar/insolar/network/consensus/common/pulse"
	"github.com/insolar/insolar/network/consensus/gcpv2/api"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/census"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/misbehavior"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/phases"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/profiles"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/proofs"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/transport"

	"github.com/insolar/insolar/network/consensus/gcpv2/censusimpl"
)

type FullRealm struct {
	coreRealm
	nodeContext nodeContext

	/* Derived from the ones provided externally - set at init() or start(). Don't need mutex */
	packetBuilder   transport.PacketBuilder
	packetSender    transport.PacketSender
	controlFeeder   api.ConsensusControlFeeder
	candidateFeeder api.CandidateControlFeeder
	profileFactory  profiles.Factory

	timings api.RoundTimings

	census     census.Active
	population RealmPopulation
	purgatory  RealmPurgatory

	packetDispatchers []PacketDispatcher

	/* Other fields - need mutex */
	isFinished bool
}

func (r *FullRealm) dispatchPacket(ctx context.Context, packet transport.PacketParser, from endpoints.Inbound,
	verifyFlags PacketVerifyFlags) error {

	pt := packet.GetPacketType()

	var sourceNode MemberPacketReceiver
	var sourceID insolar.ShortNodeID

	switch {
	case pt.GetLimitPerSender() == 0 || int(pt) >= len(r.packetDispatchers) || r.packetDispatchers[pt] == nil:
		return fmt.Errorf("packet type (%v) is unknown", pt)
	case pt.IsMemberPacket():
		strict, err := VerifyPacketRoute(ctx, packet, r.GetSelfNodeID())
		if err != nil {
			return err
		}
		if strict {
			verifyFlags |= RequireStrictVerify
		}

		sourceID = packet.GetSourceID()
		sourceNode = r.getMemberReceiver(sourceID)
	}

	if sourceNode != nil && !sourceNode.CanReceivePacket(pt) {
		return fmt.Errorf("packet type (%v) limit exceeded: from=%v(%v)", pt, sourceNode.GetNodeID(), from)
	}

	pd := r.packetDispatchers[pt] // was checked above for != nil

	if verifyFlags&(SkipVerify|SuccesfullyVerified) == 0 {
		var err error
		strict := verifyFlags&RequireStrictVerify != 0
		switch {
		case sourceNode != nil:
			err = sourceNode.VerifyPacketAuthenticity(packet.GetPacketSignature(), from, strict)
			if err != nil {
				return err
			}
			verifyFlags |= SuccesfullyVerified
		case pd.HasCustomVerifyForHost(from, strict):
			// skip default
		default:
			err = r.coreRealm.VerifyPacketAuthenticity(packet.GetPacketSignature(), from, strict)
			if err != nil {
				return err
			}
			verifyFlags |= SuccesfullyVerified
		}
	}

	//this enables lazy parsing - packet is fully parsed AFTER validation, hence makes it less prone to exploits for non-members
	var err error
	packet, err = LazyPacketParse(packet)
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
		memberCreated := false
		memberCreated, err = pd.DispatchUnknownMemberPacket(ctx, sourceID, memberPacket, from)
		if err != nil {
			return err
		}
		if !memberCreated {
			return fmt.Errorf("packet type (%v) from unknown sourceID(%v): from=%v", pt, sourceID, from)
		}

		sourceNode = r.getMemberReceiver(sourceID)
		if sourceNode == nil {
			return fmt.Errorf("inconsistent behavior for packet type (%v) from unknown sourceID(%v): from=%v", pt, sourceID, from)
		}
	}

	if !sourceNode.SetPacketReceived(pt) {
		return fmt.Errorf("packet type (%v) limit exceeded: from=%v(%v)", pt, sourceNode.GetNodeID(), from)
	}

	return sourceNode.DispatchMemberPacket(ctx, packet, from, verifyFlags, pd)
}

/* LOCK - runs under RoundController lock */
func (r *FullRealm) start(census census.Active, population census.OnlinePopulation, bundle PhaseControllersBundle) {
	r.initBasics(census)
	allControllers, perNodeControllers := r.initHandlers(bundle, population.GetCount())
	r.initPopulation(bundle.IsDynamicPopulationRequired(), population, perNodeControllers)
	r.initSelf()
	r.startWorkers(allControllers)
}

func (r *FullRealm) initBefore(transport transport.Factory, controlFeeder api.ConsensusControlFeeder,
	candidateFeeder api.CandidateControlFeeder) transport.NeighbourhoodSizes {
	r.packetSender = transport.GetPacketSender()
	r.packetBuilder = transport.GetPacketBuilder(r.signer)
	r.controlFeeder = controlFeeder
	r.candidateFeeder = candidateFeeder
	return r.packetBuilder.GetNeighbourhoodSize()
}

func (r *FullRealm) initBasics(census census.Active) {

	r.census = census
	r.profileFactory = census.GetProfileFactory(r.verifierFactory)

	r.timings = r.config.GetConsensusTimings(r.pulseData.NextPulseDelta, r.IsJoiner())
	r.strategy.AdjustConsensusTimings(&r.timings)

	r.nodeContext.initFull(r.verifierFactory, uint8(r.nbhSizes.NeighbourhoodTrustThreshold),
		func(report misbehavior.Report) interface{} {
			r.census.GetMisbehaviorRegistry().AddReport(report)
			return nil
		})
}

func (r *FullRealm) initHandlers(bundle PhaseControllersBundle, nodeCount int) ([]PhaseController, []PerNodePacketDispatcherFactory) {

	r.packetDispatchers = make([]PacketDispatcher, phases.PacketTypeCount)
	controllers, nodeCallback := bundle.CreateFullPhaseControllers(nodeCount)

	if len(controllers) == 0 {
		panic("no phase controllers")
	}
	r.nodeContext.setNodeToPhaseCallback(nodeCallback)
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

	return controllers, individualHandlers
}

func (r *FullRealm) initPopulation(needsDynamic bool, population census.OnlinePopulation, individualHandlers []PerNodePacketDispatcherFactory) {

	initNodeFn := func(ctx context.Context, n *NodeAppearance) {
		n.callback = &r.nodeContext
		for k, ctl := range individualHandlers {
			var ph PhasePerNodePacketFunc
			ctx, ph = ctl.CreatePerNodePacketHandler(ctx, n)
			if ph == nil {
				continue
			}
			if n.handlers == nil {
				n.handlers = make([]PhasePerNodePacketFunc, len(individualHandlers))
			}
			n.handlers[k] = ph
		}
	}

	if needsDynamic {
		r.population = NewDynamicRealmPopulation(r.strategy, population, r.config.GetNodeCountHint(),
			r.nbhSizes.ExtendingNeighbourhoodLimit, initNodeFn)
	} else {
		r.population = NewMemberRealmPopulation(r.strategy, population,
			r.nbhSizes.ExtendingNeighbourhoodLimit, initNodeFn)
	}

	r.purgatory = NewRealmPurgatory(r.population, r.profileFactory, r.verifierFactory, &r.nodeContext, r.postponedPacketFn)
}

func (r *FullRealm) initSelf() {
	newSelf := r.population.GetSelf()
	prevSelf := r.self

	if newSelf.GetNodeID() != prevSelf.GetNodeID() {
		panic("inconsistent transition of self between realms")
	}

	prevSelf.copySelfTo(newSelf)
	r.self = newSelf

	if !newSelf.profile.IsJoiner() {
		// joiners are not allowed to request leave
		newSelf.requestedLeave, newSelf.requestedLeaveReason = r.controlFeeder.GetRequiredGracefulLeave()
	}

	if !newSelf.requestedLeave {
		// leaver is not allowed to add new nodes
		jc := r.registerNextJoinCandidate()
		if jc != nil {
			newSelf.requestedJoinerID = jc.GetNodeID()
		}
	}
	newSelf.callback.updatePopulationVersion()
}

func (r *FullRealm) registerNextJoinCandidate() *NodeAppearance {

	if !r.GetSelf().CanIntroduceJoiner() {
		return nil
	}

	for {
		cp := r.candidateFeeder.PickNextJoinCandidate()
		if cp == nil {
			return nil
		}

		nip := r.profileFactory.CreateFullIntroProfile(cp)
		sv := r.GetSignatureVerifier(nip.GetPublicKeyStore())
		np := censusimpl.NewJoinerProfile(nip, sv, nip.GetStartPower())
		na := r.population.CreateNodeAppearance(r.roundContext, &np)
		nna, err := r.population.AddToDynamics(na)
		if err != nil {
			inslogger.FromContext(r.roundContext).Error(err)
		} else if nna != nil {
			return nna
		}
		r.candidateFeeder.RemoveJoinCandidate(false, cp.GetStaticNodeID())
	}
}

func (r *FullRealm) startWorkers(controllers []PhaseController) {
	for _, ctl := range controllers {
		ctl.BeforeStart(r)
	}
	for _, ctl := range controllers {
		ctl.StartWorker(r.roundContext, r)
	}
}

func (r *FullRealm) GetFraudFactory() misbehavior.FraudFactory {
	return r.nodeContext.fraudFactory
}

func (r *FullRealm) GetBlameFactory() misbehavior.BlameFactory {
	return r.nodeContext.blameFactory
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

func ShuffleNodeProjections(strategy RoundStrategy, nodeRefs []*NodeAppearance) {
	strategy.ShuffleNodeSequence(len(nodeRefs),
		func(i, j int) { nodeRefs[i], nodeRefs[j] = nodeRefs[j], nodeRefs[i] })
}

func (r *FullRealm) GetPopulation() RealmPopulation {
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

func (r *FullRealm) GetPulseData() pulse.Data {
	return r.pulseData
}

func (r *FullRealm) GetLastCloudStateHash() proofs.CloudStateHash {
	return r.census.GetCloudStateHash()
}

func (r *coreRealm) UpstreamPreparePulseChange() <-chan proofs.NodeStateHash {
	if !r.pulseData.PulseNumber.IsTimePulse() {
		panic("pulse number was not set")
	}

	sp := r.GetSelf().GetProfile()
	report := api.UpstreamReport{
		PulseNumber: r.pulseData.PulseNumber,
		MemberPower: sp.GetDeclaredPower(),
		MemberMode:  sp.GetOpMode(),
	}
	return r.upstream.PreparePulseChange(report)
}

func (r *FullRealm) CommitPulseChange() {
	if !r.pulseData.PulseNumber.IsTimePulse() {
		panic("pulse number was not set")
	}

	sp := r.GetSelf().GetProfile()
	report := api.UpstreamReport{
		PulseNumber: r.pulseData.PulseNumber,
		MemberPower: sp.GetDeclaredPower(),
		MemberMode:  sp.GetOpMode(),
	}
	go r.upstream.CommitPulseChange(report, r.pulseData, r.census)
}

func (r *FullRealm) GetTimings() api.RoundTimings {
	return r.timings
}

func (r *FullRealm) GetNeighbourhoodSizes() transport.NeighbourhoodSizes {
	return r.nbhSizes
}

func (r *FullRealm) GetLocalProfile() profiles.LocalNode {
	return r.self.profile.(profiles.LocalNode)
}

func (r *FullRealm) PrepareAndSetLocalNodeStateHashEvidence(nsh proofs.NodeStateHash) {
	// TODO use r.GetLastCloudStateHash() + digest(PulseData) + r.digest.GetGshDigester() to build digest for signing

	// TODO Hack! MUST provide announcement hash
	nas := cryptkit.NewSignature(nsh, "stubSign")

	v := nsh.SignWith(r.signer)
	r.self.SetLocalNodeStateHashEvidence(proofs.NewNodeStateHashEvidence(v), &nas)
}

func (r *FullRealm) CreateAnnouncement(n *NodeAppearance) *transport.NodeAnnouncementProfile {
	ma := n.GetRequestedAnnouncement()
	if ma.Membership.IsEmpty() {
		panic("illegal state")
	}

	var joiner *transport.JoinerAnnouncement
	if !ma.JoinerID.IsAbsent() {
		jna := r.GetPopulation().GetNodeAppearance(ma.JoinerID)
		switch {
		case jna != nil:
			jp := jna.GetProfile().GetStatic()
			joiner = transport.NewJoinerAnnouncement(jp, n.GetNodeID(), jp.GetJoinerSignature())
		case n == r.self:
			panic("illegal state - local joiner is missing")
		default:
			panic("illegal state - joiner is missing")
		}
	}

	return transport.NewNodeAnnouncement(n.profile, ma, r.GetNodeCount(), r.pulseData.PulseNumber, joiner)
}

func (r *FullRealm) CreateLocalAnnouncement() *transport.NodeAnnouncementProfile {
	return r.CreateAnnouncement(r.self)
}

func (r *FullRealm) CreateNextCensusBuilder() census.Builder {
	return r.census.CreateBuilder(r.GetNextPulseNumber(), true)
}

func (r *FullRealm) preparePrimingMembers(pop census.PopulationBuilder) {
	for _, p := range pop.GetUnorderedProfiles() {
		if p.GetSignatureVerifier() != nil {
			continue
		}
		v := r.GetSignatureVerifier(p.GetStatic().GetPublicKeyStore())
		p.SetSignatureVerifier(v)
	}
}

/* deprecated */
func (r *FullRealm) prepareRegularMembers(pop census.PopulationBuilder) {
	// cc := r.census.GetMandateRegistry().GetConsensusConfiguration()

	for _, p := range pop.GetUnorderedProfiles() {
		if p.GetSignatureVerifier() == nil {
			v := r.GetSignatureVerifier(p.GetStatic().GetPublicKeyStore())
			p.SetSignatureVerifier(v)
		}

		if p.GetOpMode().IsEvicted() {
			continue
		}

		var na *NodeAppearance
		if p.IsJoiner() {
			na = r.population.GetJoinerNodeAppearance(p.GetNodeID())
		} else {
			na = r.population.GetNodeAppearanceByIndex(p.GetIndex().AsInt())
		}
		rs := na.GetRequestedState()
		p.SetPower(rs.RequestedPower)
		p.SetOpMode(rs.RequestedMode)
	}
}

func (r *FullRealm) FinishRound(builder census.Builder, csh proofs.CloudStateHash) {
	r.Lock()
	defer r.Unlock()

	if r.isFinished {
		panic("illegal state")
	}
	r.isFinished = true

	local := builder.GetPopulationBuilder().GetLocalProfile()
	r.prepareRegularMembers(builder.GetPopulationBuilder())

	if local.GetOpMode().IsEvicted() {
		r.notifyConsensusFinished(local, nil)
		return
	}

	expected := builder.BuildAndMakeExpected(csh)
	r.notifyConsensusFinished(expected.GetOnlinePopulation().GetLocalProfile(), expected)
}

func (r *FullRealm) notifyConsensusFinished(newSelf profiles.ActiveNode, expectedCensus census.Operational) {
	report := api.UpstreamReport{
		PulseNumber: r.pulseData.PulseNumber,
		MemberPower: newSelf.GetDeclaredPower(),
		MemberMode:  newSelf.GetOpMode(),
	}

	r.controlFeeder.ConsensusFinished(report, expectedCensus)
	go r.upstream.ConsensusFinished(report, expectedCensus)
}

func (r *FullRealm) GetProfileFactory() profiles.Factory {
	return r.profileFactory
}

func (r *FullRealm) GetPurgatory() *RealmPurgatory {
	return &r.purgatory
}

func (r *FullRealm) getMemberReceiver(id insolar.ShortNodeID) MemberPacketReceiver {
	//Purgatory MUST be checked first to avoid "missing" a node during its transition from the purgatory to normal population
	node := r.GetPurgatory().GetPhantomNode(id)
	if node != nil {
		return node
	}
	return r.GetPopulation().GetNodeAppearance(id)
}
