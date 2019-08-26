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

package purgatory

import (
	"context"
	"fmt"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/transport"
	"sync"

	"github.com/insolar/insolar/network/consensus/gcpv2/core/packetdispatch"
	"github.com/insolar/insolar/network/consensus/gcpv2/core/population"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/network/consensus/common/cryptkit"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/member"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/misbehavior"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/profiles"
	"github.com/insolar/insolar/network/consensus/gcpv2/censusimpl"
)

func NewRealmPurgatory(population population.RealmPopulation, _ profiles.Factory, svf cryptkit.SignatureVerifierFactory,
	hook *population.Hook, postponedPacketFn packetdispatch.PostponedPacketFunc) RealmPurgatory {
	return RealmPurgatory{
		population: population,
		// profileFactory:    pf,
		svFactory:         svf,
		hook:              hook,
		postponedPacketFn: postponedPacketFn,
	}
}

type AnnouncingMember interface {
	IsJoiner() bool
	GetNodeID() insolar.ShortNodeID
	Blames() misbehavior.BlameFactory
	Frauds() misbehavior.FraudFactory
	GetReportProfile() profiles.BaseNode

	DispatchAnnouncement(ctx context.Context, rank member.Rank, profile profiles.StaticProfile,
		announcement profiles.MemberAnnouncement) error

	ApplyNeighbourEvidence(n *population.NodeAppearance, an profiles.MemberAnnouncement,
		cappedTrust bool, applyAfterChecks population.MembershipApplyFunc) (bool, error)

	GetStatic() profiles.StaticProfile
}

type RealmPurgatory struct {
	population population.RealmPopulation
	svFactory  cryptkit.SignatureVerifierFactory
	// profileFactory    profiles.Factory
	postponedPacketFn packetdispatch.PostponedPacketFunc

	hook *population.Hook

	/* LOCK WARNING!
	This lock is engaged inside NodePhantom's lock.
	DO NOT call NodePhantom methods under this lock.
	*/
	rw sync.RWMutex

	phantomByID map[insolar.ShortNodeID]*NodePhantom

	// phantomByEP map[string]*NodePhantom
}

// type PurgatoryNodeState int
//
// const PurgatoryDuplicatePK PurgatoryNodeState = -1
// const PurgatoryExistingMember PurgatoryNodeState = -2

func (p *RealmPurgatory) GetPhantomNode(id insolar.ShortNodeID) *NodePhantom {
	p.rw.RLock()
	defer p.rw.RUnlock()

	return p.phantomByID[id]
}

func (p *RealmPurgatory) getPhantomNode(id insolar.ShortNodeID) (*NodePhantom, bool) {
	p.rw.RLock()
	defer p.rw.RUnlock()

	np, ok := p.phantomByID[id]
	return np, ok
}

func (p *RealmPurgatory) getOrCreatePhantom(id insolar.ShortNodeID) AnnouncingMember {

	p.rw.Lock()
	defer p.rw.Unlock()

	np, ok := p.phantomByID[id]
	if ok {
		if np == nil { // avoid interface-nil
			return nil
		}
		return np
	}

	na := p.population.GetNodeAppearance(id)
	if na != nil {
		return na
	}

	if p.phantomByID == nil {
		p.phantomByID = make(map[insolar.ShortNodeID]*NodePhantom)
	}
	limiter := p.population.CreatePacketLimiter(false /* doesnt matter here */)
	np = NewNodePhantom(p, id, limiter)
	p.phantomByID[id] = np
	return np
}

func (p *RealmPurgatory) getOrCreateMember(id insolar.ShortNodeID) AnnouncingMember {

	na := p.population.GetNodeAppearance(id)
	if na != nil { // main path
		return na
	}

	np, ok := p.getPhantomNode(id) // read lock
	if !ok {
		am := p.getOrCreatePhantom(id) // write lock
		if am != nil {
			return am
		}
	} else if np != nil {
		return np
	}

	// NB! np == NIL - it means that phantom was moved to a normal population
	na = p.population.GetNodeAppearance(id)
	if na == nil {
		// nil entry in the purgatory means that there MUST have be a relevant NodeAppearance
		panic("illegal state")
	}
	return na
}

func (p *RealmPurgatory) FindMember(id insolar.ShortNodeID, introducedBy insolar.ShortNodeID) AnnouncingMember {
	am, _ := p.getMember(id, introducedBy)
	return am
}

func (p *RealmPurgatory) getMember(id insolar.ShortNodeID, introducedBy insolar.ShortNodeID) (AnnouncingMember, *population.NodeAppearance) {

	na := p.population.GetNodeAppearance(id)
	if na != nil { // main path
		return na, na
	}

	np, ok := p.getPhantomNode(id) // read lock
	if !ok {
		return nil, nil
	}
	if np != nil {
		np.IntroducedBy(introducedBy)
		return np, nil
	}

	na = p.population.GetNodeAppearance(id)
	if na == nil {
		// nil entry in the purgatory means that there MUST have be a relevant NodeAppearance
		panic("illegal state")
	}
	return na, na
}

func (p *RealmPurgatory) ascendFromPurgatory(ctx context.Context, phantom *NodePhantom, nsp profiles.StaticProfile,
	rank member.Rank, sv cryptkit.SignatureVerifier, announcerID insolar.ShortNodeID, joinerSecret cryptkit.DigestHolder) {

	if sv == nil {
		sv = p.svFactory.CreateSignatureVerifierWithPKS(nsp.GetPublicKeyStore())
	}

	var np censusimpl.NodeProfileSlot
	if rank.IsJoiner() {
		np = censusimpl.NewJoinerProfile(nsp, sv)
	} else {
		np = censusimpl.NewNodeProfileExt(rank.GetIndex(), nsp, sv, rank.GetPower(), rank.GetMode())
	}

	nav := population.NewAscendedNodeAppearance(&np, phantom.limiter, announcerID, joinerSecret)
	na := &nav

	p.rw.Lock()
	defer p.rw.Unlock()
	p.phantomByID[phantom.nodeID] = nil // leave marker
	// delete(p.phantomByEP, ...)

	na, _ = p.population.AddToDynamics(ctx, na)

	inslogger.FromContext(ctx).Debugf("Candidate/joiner has ascended as dynamic node: s=%d, t=%d, full=%v",
		p.hook.GetLocalNodeID(), np.GetNodeID(), np.GetStatic().GetExtension() != nil)

	p.hook.OnPurgatoryNodeUpdate(p.hook.GetPopulationVersion(), na, population.FlagAscent)
}

func (p *RealmPurgatory) IsBriefAscensionAllowed() bool {
	// using false will delay processing of packets and may result in slower consensus
	// using true may produce NodeAppearance objects with Brief profiles
	return false
}

func (p *RealmPurgatory) UnknownAsSelfFromMemberAnnouncement(ctx context.Context, id insolar.ShortNodeID,
	profile profiles.StaticProfile, rank member.Rank, announcement profiles.MemberAnnouncement) (bool, error) {

	err := p.getOrCreateMember(id).DispatchAnnouncement(ctx, rank, profile, announcement)
	return err == nil, err
}

func (p *RealmPurgatory) FindJoinerProfile(nodeID insolar.ShortNodeID, introducedBy insolar.ShortNodeID) profiles.StaticProfile {
	am, _ := p.getMember(nodeID, introducedBy)
	if am != nil && am.IsJoiner() {
		return am.GetStatic()
	}
	return nil
}

func (p *RealmPurgatory) GetJoinerAnnouncement(nodeID insolar.ShortNodeID, introducedBy insolar.ShortNodeID) *transport.JoinerAnnouncement {
	am, na := p.getMember(nodeID, introducedBy)
	if am != nil && !am.IsJoiner() {
		return nil
	}

	if na != nil {
		return na.GetAnnouncementAsJoiner()
	}

	return am.(*NodePhantom).GetAnnouncementAsJoiner()
}

func (p *RealmPurgatory) onNodeUpdated(n *NodePhantom, flags population.UpdateFlags) {
	p.hook.OnPurgatoryNodeUpdate(n.purgatory.hook.UpdatePopulationVersion(), n, flags)
}

// WARNING! Is called under NodeAppearance lock
func (p *RealmPurgatory) AddJoinerAndEnsureAscendancy(
	announcement profiles.JoinerAnnouncement, announcedByID insolar.ShortNodeID) error {

	jp := announcement.JoinerProfile
	joinerID := jp.GetStaticNodeID()

	if announcedByID == joinerID {
		panic("illegal value - cant add itself")
	}

	err := p.getOrCreateMember(joinerID).DispatchAnnouncement(context.TODO(), // TODO context
		member.JoinerRank, jp,
		profiles.NewJoinerAnnouncement(jp, announcedByID))

	sp := p.FindJoinerProfile(joinerID, announcedByID)
	if sp == nil {
		panic(fmt.Sprintf("failed addition of a joiner: id=%d announcedByID=%d", joinerID, announcedByID))
	}
	return err
}

func (p *RealmPurgatory) VerifyNeighbour(announcement profiles.MemberAnnouncement, n *population.NodeAppearance) (AnnouncingMember, bool) {

	am, na := p.getMember(announcement.MemberID, announcement.AnnouncedByID)
	if na == nil {
		return am, false
	}

	return am, profiles.EqualStaticProfiles(na.GetStatic(), announcement.Joiner.JoinerProfile, false)
}

func (p *RealmPurgatory) UnknownFromNeighbourhood(ctx context.Context, rank member.Rank, announcement profiles.MemberAnnouncement,
	cappedTrust bool) error {

	m := p.getOrCreateMember(announcement.MemberID)
	if announcement.Membership.IsJoiner() {
		if announcement.Joiner.JoinerProfile == nil {
			panic("announcement.Joiner.JoinerProfile == nil") // it must be checked by caller
		}
		return m.DispatchAnnouncement(ctx, rank, announcement.Joiner.JoinerProfile, announcement)
	}
	return m.DispatchAnnouncement(ctx, rank, nil, announcement)
}

func (p *RealmPurgatory) IsJoinerSecretRequired() bool {
	return false
}
