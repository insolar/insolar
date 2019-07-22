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
	"fmt"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/network/consensus/common/cryptkit"
	"github.com/insolar/insolar/network/consensus/common/endpoints"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/member"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/misbehavior"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/phases"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/profiles"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/transport"
	"github.com/insolar/insolar/network/consensus/gcpv2/core/packetrecorder"
	"sync"
)

func NewNodePhantom(purgatory *RealmPurgatory, nodeID insolar.ShortNodeID, limiter phases.PacketLimiter) *NodePhantom {
	return &NodePhantom{
		purgatory: purgatory,
		nodeID:    nodeID,
		limiter:   limiter,
		recorder:  packetrecorder.NewUnsafePacketRecorder(int(limiter.GetRemainingPacketCountDefault())),
	}
}

var _ MemberPacketReceiver = &NodePhantom{}
var _ MemberPacketSender = &NodePhantom{}

type NodePhantom struct {
	purgatory *RealmPurgatory

	nodeID  insolar.ShortNodeID
	mutex   sync.Mutex
	limiter phases.PacketLimiter

	figment  figment
	recorder packetrecorder.UnsafePacketRecorder
	//figments map[string]*figment
}

func (p *NodePhantom) ApplyNeighbourEvidence(n *NodeAppearance, an profiles.MembershipAnnouncement, cappedTrust bool) (bool, error) {
	return false, nil
}

func (p *NodePhantom) Blames() misbehavior.BlameFactory {
	return p.purgatory.callback.GetBlameFactory()
}

func (p *NodePhantom) Frauds() misbehavior.FraudFactory {
	return p.purgatory.callback.GetFraudFactory()
}

func (p *NodePhantom) GetReportProfile() profiles.BaseNode {
	panic("implement me")
}

func (p *NodePhantom) IsJoiner() bool {
	return p.figment.rank.IsJoiner()
}

func (p *NodePhantom) EncryptJoinerSecret(joinerSecret cryptkit.DigestHolder) cryptkit.DigestHolder {
	// TODO encryption of joinerSecret
	return joinerSecret
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
	flags packetrecorder.PacketVerifyFlags, pd PacketDispatcher) error {

	_, err := pd.DispatchUnknownMemberPacket(ctx, p.nodeID, packet.GetMemberPacket(), from)
	if err != nil {
		return err
	}

	p.postponePacket(ctx, packet, from, flags)
	return nil
}

func (p *NodePhantom) postponePacket(ctx context.Context, packet transport.PacketParser, from endpoints.Inbound, flags packetrecorder.PacketVerifyFlags) {

	inslogger.FromContext(ctx).Debugf("packet added to purgatory: s=%d t=%d pt=%v",
		packet.GetSourceID(), packet.GetTargetID(), packet.GetPacketType())

	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.recorder.Record(packet, from, flags)
}

func (p *NodePhantom) DispatchAnnouncement(ctx context.Context, rank member.Rank, profile profiles.StaticProfile,
	announcement *profiles.MembershipAnnouncement, introducedByID insolar.ShortNodeID) error {

	p.mutex.Lock()
	defer p.mutex.Unlock()

	return p.figment.dispatchAnnouncement(ctx, p, rank, profile, announcement, introducedByID)
}

func (p *NodePhantom) isAscent() bool {
	return !p.recorder.IsRecording()
}

func (p *NodePhantom) ascend(ctx context.Context, nsp profiles.StaticProfile, rank member.Rank, sv cryptkit.SignatureVerifier) {

	if p.isAscent() {
		panic("illegal state")
	}

	p.purgatory.ascendFromPurgatory(ctx, p.nodeID, nsp, rank, sv)
	p.recorder.Playback(p.purgatory.postponedPacketFn)
}

type figment struct {
	phantom     *NodePhantom
	announcerID insolar.ShortNodeID
	rank        member.Rank

	profile profiles.StaticProfile

	//announceSignature proofs.MemberAnnouncementSignature // one-time set
	//stateEvidence     proofs.NodeStateHashEvidence       // one-time set
	//firstFraudDetails *misbehavior.FraudError
	//neighborReports int
}

func (p *figment) dispatchAnnouncement(ctx context.Context, phantom *NodePhantom, rank member.Rank, profile profiles.StaticProfile,
	announcement *profiles.MembershipAnnouncement, announcedBy insolar.ShortNodeID) error {

	hasUpdate := false
	if p.phantom == nil {
		p.phantom = phantom
		p.rank = rank
	}
	if p.announcerID.IsAbsent() && announcedBy != phantom.nodeID && !announcedBy.IsAbsent() {
		p.announcerID = announcedBy
		hasUpdate = true
	}

	ascentWithBrief := p.phantom.purgatory.IsBriefAscensionAllowed()

	hasProfileUpdate, hasMismatch := p.updateProfile(rank, profile, ascentWithBrief)
	if hasMismatch {
		panic(fmt.Sprintf("inconsistent neighbour announcement: local=%d, phantom=%d, announcer=%d, rank=%v, profile=%+v, firmentRank=%v, figmentProfile=%+v, ann=%+v",
			p.phantom.purgatory.callback.localNodeID, p.phantom.nodeID, announcedBy, rank, profile, p.rank, p.profile, announcement))
		// TODO return p.RegisterFraud(p.Frauds().NewInconsistentNeighbourAnnouncement(p.GetReportProfile()))
	}
	if hasProfileUpdate {
		hasUpdate = true
	}

	if p.profile == nil || !hasUpdate {
		return nil
	}
	if p.rank.IsJoiner() && p.announcerID.IsAbsent() {
		/* self-ascension is not allowed for joiners */
		return nil
	}

	if p.profile.GetExtension() != nil || ascentWithBrief {
		inslogger.FromContext(ctx).Debugf("Phantom node ascension: s=%d, t=%d, full=%v",
			p.phantom.purgatory.callback.localNodeID, p.phantom.nodeID, p.profile.GetExtension() != nil)
		p.phantom.ascend(ctx, p.profile, p.rank, nil)
	} else {
		inslogger.FromContext(ctx).Debugf("Phantom node added: s=%d, t=%d, full=%v",
			p.phantom.purgatory.callback.localNodeID, p.phantom.nodeID, p.profile.GetExtension() != nil)
		p.phantom.purgatory.onBriefProfileCreated(p.phantom)
	}
	return nil
}

func (p *figment) updateProfile(rank member.Rank, profile profiles.StaticProfile, ascentWithBrief bool) (updated bool, mismatched bool) {

	switch {
	case rank != p.rank:
		return false, true
	case profile == nil:
		return false, false
	case p.profile == nil:
		p.profile = profile
		return true, false
	case profile.GetExtension() == nil:
		return false, !profiles.EqualStaticProfiles(p.profile, profile)
	case p.profile.GetExtension() == nil:
		if !profiles.EqualStaticProfiles(p.profile, profile) {
			return false, true
		}
		p.profile = profile
		return true, false
	default:
		return false, !profiles.EqualStaticProfiles(p.profile, profile) ||
			!profiles.EqualStaticExtensions(p.profile.GetExtension(), profile.GetExtension())
	}
}
