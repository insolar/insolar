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
	"github.com/insolar/insolar/network/consensus/gcpv2/api/profiles"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/transport"
	"github.com/insolar/insolar/network/consensus/gcpv2/core"
)

func ValidateIntrosOnMember(reader transport.ExtendedIntroReader, brief transport.BriefIntroductionReader,
	fullIntroRequired bool, n core.AnnouncingMember) error {

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

func ApplyUnknownAnnouncement(ctx context.Context, announcerID insolar.ShortNodeID,
	reader transport.AnnouncementPacketReader, brief transport.BriefIntroductionReader,
	fullIntroRequired bool, realm *core.FullRealm) (bool, error) {

	//var err error
	//err := ValidateIntrosOnMember(reader, brief, fullIntroRequired, nil)
	//if err != nil {
	//	return false, err
	//}

	na := reader.GetAnnouncementReader()
	nr := na.GetNodeRank()

	ma := AnnouncementFromReader(na)
	// TODO verify announcement content and signature

	purgatory := realm.GetPurgatory()
	if reader.HasFullIntro() {
		full := reader.GetFullIntroduction()
		intro := realm.GetProfileFactory().CreateFullIntroProfile(full)
		return purgatory.SelfFromMemberAnnouncement(ctx, announcerID, intro, nr, ma)
	} else if brief != nil {
		intro := realm.GetProfileFactory().CreateBriefIntroProfile(brief)
		return purgatory.SelfFromMemberAnnouncement(ctx, announcerID, intro, nr, ma)
	} else {
		return purgatory.SelfFromMemberAnnouncement(ctx, announcerID, nil, nr, ma)
	}
}

func ApplyMemberAnnouncement(ctx context.Context, reader transport.AnnouncementPacketReader, brief transport.BriefIntroductionReader,
	fullIntroRequired bool, n *core.NodeAppearance, realm *core.FullRealm) (bool, insolar.ShortNodeID, error) {

	//err := ValidateIntrosOnMember(reader, brief, fullIntroRequired, n)
	//if err != nil {
	//	return false, 0, err
	//}

	na := reader.GetAnnouncementReader()
	nr := na.GetNodeRank()

	if n.GetRank(realm.GetNodeCount()) != nr {
		return false, 0, n.Frauds().NewMismatchedNeighbourRank(n.GetReportProfile())
	}

	var err error
	var matches = true
	announcerID := n.GetNodeID()

	ma := AnnouncementFromReader(na)
	// TODO verify announcement content and signature

	if reader.HasFullIntro() {
		full := reader.GetFullIntroduction()
		matches = n.UpgradeDynamicNodeProfile(ctx, full)
		if !matches {
			// TODO should be fraud
			return false, 0, n.Blames().NewProtocolViolation(n.GetReportProfile(), "announcement is incorrect")
		}
	} else if brief != nil {
		matches = profiles.EqualStaticProfiles(n.GetReportProfile().GetStatic(), brief)
		if !matches {
			// TODO should be fraud
			return false, 0, n.Blames().NewProtocolViolation(n.GetReportProfile(), "announcement is incorrect")
		}
	}
	if !matches {
		// TODO should be fraud
		return false, 0, n.Blames().NewProtocolViolation(n.GetReportProfile(), "announcement is incorrect")
	}
	if !n.CanIntroduceJoiner() && !ma.JoinerID.IsAbsent() {
		return false, 0, n.Blames().NewProtocolViolation(n.GetReportProfile(), "node is not allowed to add a joiner")
	}

	modified, err := n.ApplyNodeMembership(ma)

	if err == nil && modified && !ma.JoinerID.IsAbsent() {
		purgatory := realm.GetPurgatory()

		ja := na.GetJoinerAnnouncement()
		// originID := ja.GetJoinerIntroducedByID() // applies to neighbourhood only

		var joinerIntroProfile profiles.StaticProfile
		if ja.HasFullIntro() {
			joinerIntroProfile = realm.GetProfileFactory().CreateFullIntroProfile(ja.GetFullIntroduction())
		} else {
			joinerIntroProfile = realm.GetProfileFactory().CreateUpgradableIntroProfile(ja.GetBriefIntroduction())
		}
		err = purgatory.JoinerFromMemberAnnouncement(ctx, ma.JoinerID, joinerIntroProfile, announcerID)
		if err == nil && ma.JoinerID == realm.GetSelfNodeID() {
			//we trust more to these who has introduced us
			// It is also REQUIRED as vector calculation requires at least one trusted node to work properly
			n.UpdateNodeTrustLevel(member.TrustBySome)
		}
	}
	return modified, ma.JoinerID, err
}

func ApplyNeighbourJoinerAnnouncement(ctx context.Context, sender *core.NodeAppearance,
	joinerAnnouncedBySender insolar.ShortNodeID, neighbour core.AnnouncingMember, joinerAnnouncedByNeighbour insolar.ShortNodeID,
	neighbourJoinerAnnouncement transport.JoinerAnnouncementReader, realm *core.FullRealm) error {

	if joinerAnnouncedByNeighbour.IsAbsent() {
		if neighbourJoinerAnnouncement != nil {
			return neighbour.Blames().NewProtocolViolation(sender.GetReportProfile(), "joiner profile is unexpected on neighbourhood")
		}
		return nil
	}

	purgatory := realm.GetPurgatory()

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

		brief := neighbourJoinerAnnouncement.GetBriefIntroduction()
		nbIntroProfile := realm.GetProfileFactory().CreateUpgradableIntroProfile(brief)

		introducedByID := neighbourJoinerAnnouncement.GetJoinerIntroducedByID()
		if introducedByID.IsAbsent() {
			introducedByID = sender.GetNodeID()
		}
		return purgatory.JoinerFromNeighbourhood(ctx, neighbourID, nbIntroProfile, introducedByID)
	}

	if neighbourJoinerAnnouncement == nil {
		return neighbour.Blames().NewProtocolViolation(sender.GetReportProfile(), "joiner profile was not expected in neighbourhood")
	}
	return purgatory.JoinerFromNeighbourhood(ctx, neighbourID, nil, sender.GetNodeID())
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

type ResolvedNeighbour struct {
	Neighbour    core.AnnouncingMember
	Announcement profiles.MembershipAnnouncement
}

func VerifyNeighbourhood(ctx context.Context, neighbourhood []transport.MembershipAnnouncementReader,
	n *core.NodeAppearance, realm *core.FullRealm) ([]ResolvedNeighbour, error) {

	hasThis := false
	hasSelf := false
	neighbours := make([]ResolvedNeighbour, len(neighbourhood))
	//nc := realm.GetNodeCount()
	purgatory := realm.GetPurgatory()
	localID := realm.GetSelfNodeID()
	senderID := n.GetNodeID()

	for idx, nb := range neighbourhood {
		nid := nb.GetNodeID()
		if nid == n.GetNodeID() {
			hasSelf = true
		}
		nr := nb.GetNodeRank()
		nba := AnnouncementFromReader(nb)
		neighbour, err := purgatory.MemberFromNeighbourhood(ctx, nid, nr, nba, senderID)
		if err != nil {
			return nil, err
		}
		if neighbour == nil {
			return nil, n.Frauds().NewUnknownNeighbour(n.GetReportProfile())
		}

		// TODO may vary for dynamic population
		//if neighbour.GetRank(nc) != nr {
		//	return nil, n.Frauds().NewMismatchedNeighbourRank(n.GetReportProfile())
		//}

		// TODO validate node proof - if fails, then fraud on sender
		// neighbourProfile.IsValidPacketSignature(nshEvidence.GetSignature())

		neighbours[idx].Announcement = nba
		neighbours[idx].Neighbour = neighbour

		if nid == localID {
			hasThis = true
		}
	}

	if !hasThis || hasSelf {
		return nil, n.Frauds().NewNeighbourMissingTarget(n.GetReportProfile())
	}
	if hasSelf {
		return nil, n.Frauds().NewNeighbourContainsSource(n.GetReportProfile())
	}
	return neighbours, nil
}
