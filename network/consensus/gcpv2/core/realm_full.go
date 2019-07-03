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
	"github.com/insolar/insolar/network/consensus/gcpv2/errors"

	"github.com/insolar/insolar/network/consensus/gcpv2/packets"

	common2 "github.com/insolar/insolar/network/consensus/gcpv2/common"

	"github.com/insolar/insolar/network/consensus/gcpv2/census"

	"github.com/insolar/insolar/network/consensus/common"
)

type FullRealm struct {
	coreRealm
	nodeContext

	/* Derived from the ones provided externally - set at init() or start(). Don't need mutex */
	packetBuilder   PacketBuilder
	packetSender    PacketSender
	controlFeeder   ConsensusControlFeeder
	candidateFeeder CandidateControlFeeder
	profileFactory  common2.NodeProfileFactory

	handlers []packetDispatcher

	timings  common2.RoundTimings
	nbhSizes common2.NeighbourhoodSizes

	census     census.ActiveCensus
	population RealmPopulation

	/* Other fields - need mutex */
	isFinished bool
}

/* LOCK - runs under RoundController lock */
func (r *FullRealm) start(census census.ActiveCensus, population census.OnlinePopulation) {
	r.initBasics(census)
	allCtls, perNodeCtls := r.initHandlers(population.GetCount())
	r.initPopulation(population, perNodeCtls)
	r.initSelf()
	r.startWorkers(allCtls)
}

func (r *FullRealm) init(transport TransportFactory, controlFeeder ConsensusControlFeeder, candidateFeeder CandidateControlFeeder) {
	r.packetSender = transport.GetPacketSender()
	r.packetBuilder = transport.GetPacketBuilder(r.signer)
	r.controlFeeder = controlFeeder
	r.candidateFeeder = candidateFeeder
}

func (r *FullRealm) initBasics(census census.ActiveCensus) {

	r.census = census
	r.profileFactory = census.GetProfileFactory(r.verifierFactory)

	r.timings = r.config.GetConsensusTimings(r.pulseData.NextPulseDelta, r.IsJoiner())
	r.strategy.AdjustConsensusTimings(&r.timings)

	r.nbhSizes = r.packetBuilder.GetNeighbourhoodSize()

	r.nodeContext.initFull(uint8(r.nbhSizes.NeighbourhoodTrustThreshold),
		func(report errors.MisbehaviorReport) interface{} {
			r.census.GetMisbehaviorRegistry().AddReport(report)
			return nil
		})
}

func (r *FullRealm) initHandlers(nodeCount int) (allControllers []PhaseController, perNodeControllers []PhaseController) {
	r.handlers = make([]packetDispatcher, packets.MaxPacketType)

	controllers, nodeCallback := r.strategy.GetFullPhaseControllers(nodeCount)

	if len(controllers) == 0 {
		panic("no phase controllers")
	}
	r.nodeContext.setNodeToPhaseCallback(nodeCallback)
	individualHandlers := make([]PhaseController, 0, len(controllers))

	for _, ctl := range controllers {
		pt := ctl.GetPacketType()
		dh := &r.handlers[pt]
		if dh.init(r, ctl) {
			dh.setRedirectHandler(len(individualHandlers))
			individualHandlers = append(individualHandlers, ctl)
		}
		if dh.HasMemberHandler() && !pt.IsMemberPacket() {
			panic("only member packet types can be handled as member/per-node")
		}
	}

	return controllers, individualHandlers
}

func (r *FullRealm) initPopulation(population census.OnlinePopulation, individualHandlers []PhaseController) {

	r.population = NewMemberRealmPopulation(r.strategy, population,
		func(ctx context.Context, n *NodeAppearance) {
			n.callback = &r.nodeContext
			for k, ctl := range individualHandlers {
				var ph PhasePerNodePacketFunc
				ph, ctx = ctl.CreatePerNodePacketHandler(k, n, r, ctx)
				if ph == nil {
					continue
				}
				if n.handlers == nil {
					n.handlers = make([]PhasePerNodePacketFunc, len(individualHandlers))
				}
				n.handlers[k] = ph
			}
		})
}

func (r *FullRealm) initSelf() {
	newSelf := r.population.GetSelf()
	prevSelf := r.self

	if newSelf.GetShortNodeID() != prevSelf.GetShortNodeID() {
		panic("inconsistent transition of self between realms")
	}

	prevSelf.copySelfTo(newSelf)
	r.self = newSelf

	newSelf.requestedLeave, newSelf.requestedLeaveExitCode = r.controlFeeder.GetRequiredGracefulLeave()
	if !newSelf.requestedLeave {
		newSelf.requestedJoiner = r.pickNextJoinCandidate()
	}
	newSelf.callback.updatePopulationVersion()
}

func (r *FullRealm) pickNextJoinCandidate() *NodeAppearance {
	for {
		cp := r.candidateFeeder.PickNextJoinCandidate()
		if cp == nil {
			return nil
		}

		nip := r.profileFactory.CreateFullIntroProfile(cp)
		sv := r.GetSignatureVerifier(nip.GetNodePublicKeyStore())
		np := census.NewNodeProfile(0, common2.Working, nip, sv, 0)
		na := r.population.CreateNodeAppearance(r.roundContext, &np)
		nna, nodes := r.population.AddToDynamics(na)

		if !common2.EqualIntroProfiles(nna.profile, na.profile) {
			nodes = append(nodes, na)
			nna = nil
		}
		if nodes != nil {
			inslogger.FromContext(r.roundContext).Errorf("multiple joiners on same id(%v): %v", cp.GetNodeID(), nodes)
		}
		if nna != nil {
			return nna
		}
		r.candidateFeeder.RemoveJoinCandidate(false, cp.GetNodeID())
	}
}

func (r *FullRealm) startWorkers(controllers []PhaseController) {
	for _, ctl := range controllers {
		ctl.BeforeStart(r)
	}
	for _, ctl := range controllers {
		ctl.StartWorker(r.roundContext)
	}
}

func (r *FullRealm) GetPacketSender() PacketSender {
	return r.packetSender
}

func (r *FullRealm) GetPacketBuilder() PacketBuilder {
	return r.packetBuilder
}

func (r *FullRealm) GetSigner() common.DigestSigner {
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
	return r.population.GetNodeCount()
}

func (r *FullRealm) GetPulseNumber() common.PulseNumber {
	return r.pulseData.PulseNumber
}

func (r *FullRealm) GetNextPulseNumber() common.PulseNumber {
	return r.pulseData.GetNextPulseNumber()
}

func (r *FullRealm) GetOriginalPulse() common2.OriginalPulsarPacket {
	// NB! locks for this field are only needed for PrepRealm
	return r.coreRealm.originalPulse
}

func (r *FullRealm) GetPulseData() common.PulseData {
	return r.pulseData
}

func (r *FullRealm) GetLastCloudStateHash() common2.CloudStateHash {
	return r.census.GetCloudStateHash()
}

func (r *coreRealm) UpstreamPreparePulseChange() <-chan common2.NodeStateHash {
	if !r.pulseData.PulseNumber.IsTimePulse() {
		panic("pulse number was not set")
	}

	sp := r.GetSelf().GetProfile()
	report := MembershipUpstreamReport{
		PulseNumber:     r.pulseData.PulseNumber,
		MemberPower:     sp.GetDeclaredPower(),
		MembershipState: sp.GetState(),
	}
	return r.upstream.PreparePulseChange(report)
}

func (r *FullRealm) GetTimings() common2.RoundTimings {
	return r.timings
}

func (r *FullRealm) GetNeighbourhoodSizes() common2.NeighbourhoodSizes {
	return r.nbhSizes
}

func (r *FullRealm) GetLocalProfile() common2.LocalNodeProfile {
	return r.self.profile.(common2.LocalNodeProfile)
}

func (r *FullRealm) PrepareAndSetLocalNodeStateHashEvidence(nsh common2.NodeStateHash) {
	// TODO use r.GetLastCloudStateHash() + digest(PulseData) + r.digest.GetGshDigester() to build digest for signing

	//TODO Hack! MUST provide announcement hash
	nas := common.NewSignature(nsh, "stubSign")

	v := nsh.SignWith(r.signer)
	r.self.SetLocalNodeStateHashEvidence(common2.NewNodeStateHashEvidence(v), &nas)
}

func (r *FullRealm) CreateAnnouncement(n *NodeAppearance) *packets.NodeAnnouncementProfile {
	return packets.NewNodeAnnouncement(n.profile, n.GetNodeMembershipProfile(), r.GetNodeCount(), r.pulseData.PulseNumber)
}

func (r *FullRealm) CreateLocalAnnouncement() *packets.NodeAnnouncementProfile {
	return r.CreateAnnouncement(r.self)
}

func (r *FullRealm) CreateNextPopulationBuilder() census.Builder {
	return r.census.CreateBuilder(r.GetNextPulseNumber())
}

func (r *FullRealm) preparePrimingMembers(pop census.OnlinePopulationBuilder) {
	for _, p := range pop.GetUnorderedProfiles() {
		if p.GetSignatureVerifier() != nil {
			continue
		}
		v := r.GetSignatureVerifier(p.GetNodePublicKeyStore())
		p.SetSignatureVerifier(v)
	}
}

func (r *FullRealm) prepareRegularMembers(pop census.OnlinePopulationBuilder) {
	cc := r.census.GetMandateRegistry().GetConsensusConfiguration()

	pulsesInJustJoined := cc.GetPulsesForJustJoinedState()
	pulsesInSuspected := cc.GetPulsesForSuspectedState()

	for _, p := range pop.GetUnorderedProfiles() {
		if p.GetSignatureVerifier() == nil {
			v := r.GetSignatureVerifier(p.GetNodePublicKeyStore())
			p.SetSignatureVerifier(v)
		}
		ns := p.GetState()
		if ns.InSuspectedExceeded(int(pulsesInSuspected)) {
			panic("node must be removed as suspected")
		}
		ns = ns.UpdateOnNextPulse(int(pulsesInJustJoined))
		p.SetState(ns)

		idx := p.GetIndex()
		if idx >= 0 {
			na := r.population.GetNodeAppearanceByIndex(idx)
			p.SetPower(na.requestedPower)
		} else {
			na := r.population.GetJoinerNodeAppearance(p.GetShortNodeID())
			p.SetPower(na.requestedPower)
		}
	}
}

func (r *FullRealm) FinishRound(builder census.Builder, csh common2.CloudStateHash) {
	r.Lock()
	defer r.Unlock()

	if r.isFinished {
		panic("illegal state")
	}
	r.isFinished = true

	r.prepareRegularMembers(builder.GetOnlinePopulationBuilder())
	expected := builder.BuildAndMakeExpected(csh)

	r.upstreamMembershipConfirmed(expected)
}

func (r *FullRealm) upstreamMembershipConfirmed(expectedCensus census.OperationalCensus) {
	sp := r.GetSelf().GetProfile()
	report := MembershipUpstreamReport{
		PulseNumber:     r.pulseData.PulseNumber,
		MemberPower:     sp.GetDeclaredPower(),
		MembershipState: sp.GetState(),
	}

	r.upstream.MembershipConfirmed(report, expectedCensus)
}

func (r *FullRealm) getPacketDispatcher(pt packets.PacketType) (*packetDispatcher, error) {
	if int(pt) >= len(r.handlers) || !r.handlers[pt].IsEnabled() {
		return nil, fmt.Errorf("packet type (%v) has no handler", pt)
	}
	return &r.handlers[pt], nil
}

func (r *FullRealm) GetProfileFactory() common2.NodeProfileFactory {
	return r.profileFactory
}

func (r *FullRealm) CreatePurgatoryNode(ctx context.Context, intro packets.BriefIntroductionReader, from common.HostIdentityHolder) (*NodeAppearance, error) {

	panic("not implemented")
	//nip := r.profileFactory.CreateBriefIntroProfile(intro, intro.GetJoinerSignature())
	//if fIntro, ok := intro.(packets.FullIntroductionReader); ok && !fIntro.GetIssuerID().IsAbsent() {
	//	nip = r.profileFactory.CreateFullIntroProfile(nip, fIntro)
	//}
	//na := r.population.CreateNodeAppearance(r.roundContext, nip)
	//
	//nna, ps := r.population.AddToPurgatory(na)
	//
	//if !common2.EqualIntroProfiles(nna.profile, na.profile) {
	//	nodes = append(nodes, na)
	//	nna = nil
	//}
	//if nodes != nil {
	//	inslogger.FromContext(r.roundContext).Errorf("multiple joiners on same id(%v): %v", cp.GetNodeID(), nodes)
	//}
	//if nna != nil {
	//	newSelf.requestedJoiner = nna
	//	break
	//}

}

//func (r *FullRealm) getPurgatoryNode(profile common2.BriefCandidateProfile) *NodeAppearance {
//
//}
//
//func (r *FullRealm) createPurgatoryNode(profile common2.BriefCandidateProfile, nodeSignature common.SignatureHolder) *NodeAppearance {
//	pr := r.profileFactory.CreateBriefIntroProfile(profile, nodeSignature)
//
//}
//
//func (r *FullRealm) _registerPurgatoryNode(profile common2.BriefCandidateProfile) *NodeAppearance {
//
//}
//
//func (r *FullRealm) CreatePurgatoryNode(profile common2.BriefCandidateProfile) *NodeAppearance {
//	r.
//}
//
//func (r *FullRealm) UpgradeToDynamicNode(n *NodeAppearance, profile common2.CandidateProfile) {
//
//}
