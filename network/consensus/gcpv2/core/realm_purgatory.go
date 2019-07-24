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
	"sync"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/network/consensus/common/cryptkit"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/member"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/misbehavior"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/profiles"
	"github.com/insolar/insolar/network/consensus/gcpv2/censusimpl"
	"github.com/insolar/insolar/network/consensus/gcpv2/core/packetrecorder"
)

func NewRealmPurgatory(population RealmPopulation, _ profiles.Factory, svf cryptkit.SignatureVerifierFactory,
	callback *nodeContext, postponedPacketFn packetrecorder.PostponedPacketFunc) RealmPurgatory {
	return RealmPurgatory{
		population: population,
		// profileFactory:    pf,
		svFactory:         svf,
		callback:          callback,
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
		announcement *profiles.MembershipAnnouncement, introducedByID insolar.ShortNodeID) error

	ApplyNeighbourEvidence(n *NodeAppearance, an profiles.MembershipAnnouncement, cappedTrust bool) (bool, error)
	GetStatic() profiles.StaticProfile
}

type RealmPurgatory struct {
	population RealmPopulation
	svFactory  cryptkit.SignatureVerifierFactory
	// profileFactory    profiles.Factory
	postponedPacketFn packetrecorder.PostponedPacketFunc

	callback *nodeContext

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
	limiter := p.population.CreatePacketLimiter()
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

func (p *RealmPurgatory) getMember(id insolar.ShortNodeID, introducedBy insolar.ShortNodeID) AnnouncingMember {

	na := p.population.GetNodeAppearance(id)
	if na != nil { // main path
		return na
	}

	np, ok := p.getPhantomNode(id) // read lock
	if !ok {
		return nil
	}
	if np != nil {
		// np.IntroducedBy(introducedBy) TODO do we need it?
		return np
	}

	na = p.population.GetNodeAppearance(id)
	if na == nil {
		// nil entry in the purgatory means that there MUST have be a relevant NodeAppearance
		panic("illegal state")
	}
	return na
}

func (p *RealmPurgatory) ascendFromPurgatory(ctx context.Context, id insolar.ShortNodeID, nsp profiles.StaticProfile,
	rank member.Rank, sv cryptkit.SignatureVerifier) {

	if sv == nil {
		sv = p.svFactory.GetSignatureVerifierWithPKS(nsp.GetPublicKeyStore())
	}
	var np censusimpl.NodeProfileSlot
	if rank.IsJoiner() {
		np = censusimpl.NewJoinerProfile(nsp, sv)
	} else {
		np = censusimpl.NewNodeProfileExt(rank.GetIndex(), nsp, sv, rank.GetPower(), rank.GetMode())
	}
	na := p.population.CreateNodeAppearance(ctx, &np)

	p.rw.Lock()
	defer p.rw.Unlock()
	p.phantomByID[id] = nil // leave marker
	// delete(p.phantomByEP, ...)
	na, _ = p.population.AddToDynamics(na)
	if na.IsJoiner() {
		_, err := na.ApplyNodeStateHashEvidenceForJoiner()
		if err != nil {
			inslogger.FromContext(ctx).Error(err)
		}
	}

	inslogger.FromContext(ctx).Debugf("Candidate/joiner has ascended as dynamic node: s=%d, t=%d, full=%v",
		p.callback.localNodeID, np.GetNodeID(), np.GetStatic().GetExtension() != nil)
}

func (p *RealmPurgatory) IsBriefAscensionAllowed() bool {
	// using false will delay processing of packets and may result in slower consensus
	// using true may produce NodeAppearance objects with Brief profiles
	return true
}

func (p *RealmPurgatory) SelfFromMemberAnnouncement(ctx context.Context, id insolar.ShortNodeID, profile profiles.StaticProfile,
	rank member.Rank, announcement profiles.MembershipAnnouncement) (bool, error) {

	err := p.getOrCreateMember(id).DispatchAnnouncement(ctx, rank, profile, &announcement, id)
	return err == nil, err
}

func (p *RealmPurgatory) JoinerFromMemberAnnouncement(ctx context.Context, id insolar.ShortNodeID, profile profiles.StaticProfile,
	introducedByID insolar.ShortNodeID) error {

	return p.getOrCreateMember(id).DispatchAnnouncement(ctx, member.JoinerRank, profile, nil, introducedByID)
}

func (p *RealmPurgatory) JoinerFromNeighbourhood(ctx context.Context, id insolar.ShortNodeID, profile profiles.StaticProfile,
	introducedByID insolar.ShortNodeID) error {

	return p.getOrCreateMember(id).DispatchAnnouncement(ctx, member.JoinerRank, profile, nil, introducedByID)
}

func (p *RealmPurgatory) MemberFromNeighbourhood(ctx context.Context, id insolar.ShortNodeID, rank member.Rank,
	announcement profiles.MembershipAnnouncement, introducedByID insolar.ShortNodeID) (AnnouncingMember, error) {

	am := p.getOrCreateMember(id)
	return am, am.DispatchAnnouncement(ctx, rank, nil, &announcement, introducedByID)
}

func (p *RealmPurgatory) FindJoinerProfile(nodeID insolar.ShortNodeID, introducedBy insolar.ShortNodeID) profiles.StaticProfile {
	am := p.getMember(nodeID, introducedBy)
	if am != nil && am.IsJoiner() {
		return am.GetStatic()
	}
	return nil
}

func (p *RealmPurgatory) onNodeUpdated(n *NodePhantom, flags UpdateFlags) {
	p.callback.onPurgatoryNodeUpdate(p.callback.updatePopulationVersion(), n, flags)
}
