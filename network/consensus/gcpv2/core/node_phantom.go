///
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
///

package core

import (
	"context"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/network/consensus/common/cryptkit"
	"github.com/insolar/insolar/network/consensus/common/endpoints"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/phases"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/profiles"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/transport"
	"sync"
)

func NewNodePhantom(purgatory *RealmPurgatory, nodeID insolar.ShortNodeID, limiter phases.PacketLimiter) *NodePhantom {
	return &NodePhantom{
		purgatory:        purgatory,
		nodeID:           nodeID,
		limiter:          limiter,
		postponedPackets: make([]PostponedPacket, 0, 1+limiter.GetRemainingPacketCountDefault()<<1),
	}
}

var _ MemberPacketReceiver = &NodePhantom{}
var _ MemberPacketSender = &NodePhantom{}

type NodePhantom struct {
	purgatory *RealmPurgatory

	nodeID  insolar.ShortNodeID
	mutex   sync.Mutex
	limiter phases.PacketLimiter

	figment figment
	//figments map[string]*figment

	postponedPackets []PostponedPacket
}

func (p *NodePhantom) GetStatic() profiles.StaticProfile {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	sp := p.figment.profile
	if sp == nil {
		panic("illegal state")
	}
	return sp
}

func (p *NodePhantom) SetPacketSent(pt phases.PacketType) bool {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	var allowed bool
	allowed, p.limiter = p.limiter.SetPacketSent(pt)
	return allowed
}

func (p *NodePhantom) GetNodeID() insolar.ShortNodeID {
	return p.nodeID
}

func (p *NodePhantom) CanReceivePacket(pt phases.PacketType) bool {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	return p.limiter.CanReceivePacket(pt)
}

func (p *NodePhantom) VerifyPacketAuthenticity(ps cryptkit.SignedDigest, from endpoints.Inbound, strictFrom bool) error {
	return nil
}

func (p *NodePhantom) SetPacketReceived(pt phases.PacketType) bool {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	var allowed bool
	allowed, p.limiter = p.limiter.SetPacketReceived(pt)
	return allowed
}

func (p *NodePhantom) DispatchMemberPacket(ctx context.Context, packet transport.PacketParser, from endpoints.Inbound,
	flags PacketVerifyFlags, pd PacketDispatcher) error {

	if p.WasAscent() {
		// MUST be outside of locks
		p.purgatory.sendPostponedPacket(ctx, packet, from, flags)
		return nil
	}

	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.postponedPackets = append(p.postponedPackets, PostponedPacket{packet, from, flags})
	return nil
}

func (p *NodePhantom) ApplyNodeIntro(ctx context.Context,
	brief profiles.BriefCandidateProfile, full profiles.CandidateProfile,
	announcerID insolar.ShortNodeID, originID insolar.ShortNodeID) error {

	p.mutex.Lock()
	defer p.mutex.Unlock()

	if (brief == nil) == (full == nil) {
		panic("illegal value")
	}

	return p.figment.applyNodeIntro(ctx, p, brief, full, announcerID, originID)
}

func (p *NodePhantom) WasAscent() bool {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	return p.postponedPackets == nil
}

func (p *NodePhantom) ascend(ctx context.Context, nsp profiles.StaticProfile, sv cryptkit.SignatureVerifier) {

	if p.postponedPackets == nil {
		panic("illegal state")
	}

	packets := p.postponedPackets
	p.postponedPackets = nil //mark

	p.purgatory.ascendFromPurgatory(ctx, p.nodeID, nsp, sv, packets)
}

type figment struct {
	phantom     *NodePhantom
	originID    insolar.ShortNodeID
	announcerID insolar.ShortNodeID

	profile profiles.StaticProfile

	//announceSignature proofs.MemberAnnouncementSignature // one-time set
	//stateEvidence     proofs.NodeStateHashEvidence       // one-time set
	//firstFraudDetails *misbehavior.FraudError
	//neighborReports int
}

func (p *figment) applyNodeIntro(ctx context.Context, phantom *NodePhantom,
	brief profiles.BriefCandidateProfile, full profiles.CandidateProfile,
	announcerID insolar.ShortNodeID, originID insolar.ShortNodeID) error {

	if p.phantom == nil {
		p.phantom = phantom
	}
	if p.originID.IsAbsent() {
		p.originID = originID
	}

	ascentWithBrief := p.phantom.purgatory.IsBriefAscensionAllowed()

	hasUpdate, hasMismatch := p.updateProfile(brief, full, ascentWithBrief)
	if hasMismatch {
		panic("InconsistentNeighbourAnnouncement") // TODO
		//return p.RegisterFraud(p.Frauds().NewInconsistentNeighbourAnnouncement(p.GetReportProfile()))
	}

	if p.announcerID.IsAbsent() {
		if announcerID.IsAbsent() || announcerID == p.phantom.nodeID /* self-ascension is not allowed */ {
			return nil
		}
		// TODO do we need to double-check, that the announcerID is an active node?
		p.announcerID = announcerID
		if p.profile == nil {
			return nil
		}
	} else if !hasUpdate {
		return nil
	}

	if p.profile.GetExtension() != nil || ascentWithBrief {
		p.phantom.ascend(ctx, p.profile, nil)
	} else {
		p.phantom.purgatory.onBriefProfileCreated(p.phantom)
	}
	return nil
}

func (p *figment) updateProfile(brief profiles.BriefCandidateProfile, full profiles.CandidateProfile, ascentWithBrief bool) (bool, bool) {

	switch {
	case p.profile == nil:
		switch {
		case full != nil:
			p.profile = p.phantom.purgatory.GetProfileFactory().CreateFullIntroProfile(full)
		case brief != nil:
			if ascentWithBrief {
				p.profile = p.phantom.purgatory.GetProfileFactory().CreateUpgradableIntroProfile(brief)
			} else {
				p.profile = p.phantom.purgatory.GetProfileFactory().CreateBriefIntroProfile(brief)
			}
		default:
			return false, false
		}
		return true, false
	case full != nil:
		matches, created := profiles.ApplyNodeIntro(p.profile, brief, full)
		if !matches {
			if created != nil {
				panic("illegal state")
			}
			return false, true
		}
		return created != nil, false
	case brief != nil:
		matches, created := profiles.ApplyNodeIntro(p.profile, brief, full)
		if created != nil { //it is impossible to create a full profile on a brief intro
			panic("illegal state")
		}
		return false, !matches
	default:
		return false, false
	}
}
