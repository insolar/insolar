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

package announce

import (
	"context"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/member"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/misbehavior"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/profiles"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/transport"
	"github.com/insolar/insolar/network/consensus/gcpv2/core"
)

type AnnouncingMember interface {
	IsJoiner() bool
	GetNodeID() insolar.ShortNodeID
	Blames() misbehavior.BlameFactory
	Frauds() misbehavior.FraudFactory
	GetReportProfile() profiles.BaseNode
	ApplyNodeMembership(announcement profiles.MembershipAnnouncement) (bool, error)
	GetRank(nodeCount int) member.Rank
	CanIntroduceJoiner() bool
}

func ValidateIntrosOnMember(reader transport.ExtendedIntroReader, brief transport.BriefIntroductionReader,
	fullIntroRequired bool, n AnnouncingMember) error {

	if reader.HasJoinerSecret() {
		return n.Blames().NewProtocolViolation(n.GetReportProfile(), "joiner secret was not expected")
	}

	if reader.HasCloudIntro() {
		return n.Blames().NewProtocolViolation(n.GetReportProfile(), "cloud intro can NOT be sent by joiner")
	}

	if reader.HasFullIntro() || brief != nil {
		if !n.IsJoiner() {
			return n.Blames().NewProtocolViolation(n.GetReportProfile(), "intro(s) were not expected")
		}
		if reader.HasFullIntro() {
			return nil
		}
		if fullIntroRequired {
			return n.Blames().NewProtocolViolation(n.GetReportProfile(), "joiner MUST send a full intro")
		}
		if brief == nil {
			return n.Blames().NewProtocolViolation(n.GetReportProfile(), "joiner MUST send at least a brief intro")
		}
		return nil
	}
	if n.IsJoiner() {
		return n.Blames().NewProtocolViolation(n.GetReportProfile(), "joiner MUST send a brief or a full intro")
	}
	return nil
}

func ApplyUnknownJoinerAnnouncement(ctx context.Context, announcerID insolar.ShortNodeID,
	reader transport.AnnouncementPacketReader, brief transport.BriefIntroductionReader,
	fullIntroRequired bool, realm *core.FullRealm) (bool, error) {

	var err error
	//err := ValidateIntrosOnMember(reader, brief, fullIntroRequired, nil)
	//if err != nil {
	//	return false, err
	//}

	na := reader.GetAnnouncementReader()
	if !na.GetNodeRank().IsJoiner() {
		return false, nil
	}

	purgatory := realm.GetPurgatory()
	if reader.HasFullIntro() {
		// announcer is joiner
		err = purgatory.FromSelfIntroduction(ctx, announcerID, nil, reader.GetFullIntroduction())
	} else if brief != nil {
		// announcer is joiner
		err = purgatory.FromSelfIntroduction(ctx, announcerID, brief, nil)
	}
	return err == nil, err
}

func ApplyMemberAnnouncement(ctx context.Context, reader transport.AnnouncementPacketReader, brief transport.BriefIntroductionReader,
	fullIntroRequired bool, n AnnouncingMember, realm *core.FullRealm) (bool, insolar.ShortNodeID, error) {

	err := ValidateIntrosOnMember(reader, brief, fullIntroRequired, n)
	if err != nil {
		return false, 0, err
	}

	na := reader.GetAnnouncementReader()
	nr := na.GetNodeRank()

	if n.GetRank(realm.GetNodeCount()) != nr {
		return false, 0, n.Frauds().NewMismatchedNeighbourRank(n.GetReportProfile())
	}

	purgatory := realm.GetPurgatory()
	announcerID := n.GetNodeID()
	if reader.HasFullIntro() {
		// announcer is joiner
		err = purgatory.FromSelfIntroduction(ctx, announcerID, nil, reader.GetFullIntroduction())
	} else if brief != nil {
		// announcer is joiner
		err = purgatory.FromSelfIntroduction(ctx, announcerID, brief, nil)
	}
	if err != nil {
		return false, 0, err
	}

	ma := AnnouncementFromReader(na)

	if !n.CanIntroduceJoiner() && !ma.JoinerID.IsAbsent() {
		return false, 0, n.Blames().NewProtocolViolation(n.GetReportProfile(), "joiner is not allowed to add a joiner")
	}

	modified, err := n.ApplyNodeMembership(ma)

	if err == nil && modified && !ma.JoinerID.IsAbsent() {
		ja := na.GetJoinerAnnouncement()
		err = purgatory.FromMemberAnnouncement(ctx, ma.JoinerID, ja.GetBriefIntroduction(), nil, announcerID)
	}
	return modified, ma.JoinerID, err
}

func ApplyNeighbourJoinerAnnouncement(ctx context.Context, purgatory *core.RealmPurgatory, sender AnnouncingMember,
	joinerAnnouncedBySender insolar.ShortNodeID, neighbour AnnouncingMember, joinerAnnouncedByNeighbour insolar.ShortNodeID,
	neighbourJoinerAnnouncement transport.JoinerAnnouncementReader) error {

	if joinerAnnouncedByNeighbour.IsAbsent() {
		if neighbourJoinerAnnouncement != nil {
			return neighbour.Blames().NewProtocolViolation(sender.GetReportProfile(), "joiner profile is unexpected on neighbourhood")
		}
		return nil
	}

	neighbourID := neighbour.GetNodeID()
	if neighbour.IsJoiner() {
		if neighbourID == joinerAnnouncedBySender {
			if neighbourJoinerAnnouncement == nil {
				return nil //ok, we've got details from the sender's announcement
			}
			return neighbour.Blames().NewProtocolViolation(sender.GetReportProfile(), "joiner profile is duplicated in neighbourhood")
		}
		if neighbourJoinerAnnouncement == nil {
			return neighbour.Blames().NewProtocolViolation(sender.GetReportProfile(), "joiner profile is missing in neighbourhood")
		}
		return purgatory.BriefSelfFromNeighbourhood(ctx, neighbourID,
			neighbourJoinerAnnouncement.GetBriefIntroduction(), neighbourJoinerAnnouncement.GetJoinerIntroducedByID())
	}

	if neighbourJoinerAnnouncement == nil {
		return neighbour.Blames().NewProtocolViolation(sender.GetReportProfile(), "joiner profile is missing in neighbourhood")
	}
	return purgatory.NoticeFromNeighbourhood(ctx, joinerAnnouncedByNeighbour, neighbourID, sender.GetNodeID())
}

func AnnouncementFromReader(nb transport.MembershipAnnouncementReader) profiles.MembershipAnnouncement {

	nr := nb.GetNodeRank()
	mp := profiles.NewMembershipProfile(nr.GetMode(), nr.GetPower(), nr.GetIndex(), nb.GetNodeStateHashEvidence(),
		nb.GetAnnouncementSignature(), nb.GetRequestedPower())

	switch {
	case nb.IsLeaving():
		return profiles.NewMembershipAnnouncementWithLeave(mp, nb.GetLeaveReason())
	default:
		return profiles.NewMembershipAnnouncementWithJoinerID(mp, nb.GetJoinerID())
	}
}
