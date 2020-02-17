// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package profiles

import (
	"fmt"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/network/consensus/common/cryptkit"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/member"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/proofs"
)

type MembershipProfile struct {
	Index          member.Index
	Mode           member.OpMode
	Power          member.Power
	RequestedPower member.Power
	proofs.NodeAnnouncedState
}

// TODO support joiner in MembershipProfile
// func (v MembershipProfile) IsJoiner() bool {
//
// }

func NewMembershipProfile(mode member.OpMode, power member.Power, index member.Index,
	nsh cryptkit.SignedDigestHolder, nas proofs.MemberAnnouncementSignature,
	ep member.Power) MembershipProfile {

	return MembershipProfile{
		Index:          index,
		Power:          power,
		Mode:           mode,
		RequestedPower: ep,
		NodeAnnouncedState: proofs.NodeAnnouncedState{
			StateEvidence:     nsh,
			AnnounceSignature: nas,
		},
	}
}

func NewMembershipProfileForJoiner(brief BriefCandidateProfile) MembershipProfile {

	return MembershipProfile{
		Index:          member.JoinerIndex,
		Power:          0,
		Mode:           0,
		RequestedPower: brief.GetStartPower(),
		NodeAnnouncedState: proofs.NodeAnnouncedState{
			StateEvidence:     brief.GetBriefIntroSignedDigest(),
			AnnounceSignature: brief.GetBriefIntroSignedDigest().GetSignatureHolder(),
		},
	}
}

func NewMembershipProfileByNode(np ActiveNode, nsh cryptkit.SignedDigestHolder, nas proofs.MemberAnnouncementSignature,
	ep member.Power) MembershipProfile {

	idx := member.JoinerIndex
	if !np.IsJoiner() {
		idx = np.GetIndex()
	}

	return NewMembershipProfile(np.GetOpMode(), np.GetDeclaredPower(), idx, nsh, nas, ep)
}

func (p MembershipProfile) IsEmpty() bool {
	return p.StateEvidence == nil || p.AnnounceSignature == nil
}

func (p MembershipProfile) IsJoiner() bool {
	return p.Index.IsJoiner()
}

func (p MembershipProfile) CanIntroduceJoiner() bool {
	return p.Mode.CanIntroduceJoiner(p.Index.IsJoiner())
}

func (p MembershipProfile) AsRank(nc int) member.Rank {
	if p.Index.IsJoiner() {
		return member.JoinerRank
	}
	return member.NewMembershipRank(p.Mode, p.Power, p.Index, member.AsIndex(nc))
}

func (p MembershipProfile) AsRankUint16(nc uint16) member.Rank {
	if p.Index.IsJoiner() {
		return member.JoinerRank
	}
	return member.NewMembershipRank(p.Mode, p.Power, p.Index, member.AsIndexUint16(nc))
}

func (p MembershipProfile) Equals(o MembershipProfile) bool {
	if p.Index != o.Index || p.Power != o.Power || p.IsEmpty() || o.IsEmpty() || p.RequestedPower != o.RequestedPower {
		return false
	}

	return p.NodeAnnouncedState.Equals(o.NodeAnnouncedState)
}

func (p MembershipProfile) StringParts() string {
	if p.Power == p.RequestedPower {
		return fmt.Sprintf("pw:%v se:%v cs:%v", p.Power, p.StateEvidence, p.AnnounceSignature)
	}

	return fmt.Sprintf("pw:%v->%v se:%v cs:%v", p.Power, p.RequestedPower, p.StateEvidence, p.AnnounceSignature)
}

func (p MembershipProfile) String() string {
	index := "joiner"
	if !p.Index.IsJoiner() {
		index = fmt.Sprintf("idx:%d", p.Index)
	}
	return fmt.Sprintf("%s %s", index, p.StringParts())
}

type JoinerAnnouncement struct {
	JoinerProfile  StaticProfile
	IntroducedByID insolar.ShortNodeID
	JoinerSecret   cryptkit.DigestHolder
}

func (v JoinerAnnouncement) IsEmpty() bool {
	return v.JoinerProfile == nil
}

type MembershipAnnouncement struct {
	Membership   MembershipProfile
	IsLeaving    bool
	LeaveReason  uint32
	JoinerID     insolar.ShortNodeID
	JoinerSecret cryptkit.DigestHolder
}

type MemberAnnouncement struct {
	MemberID insolar.ShortNodeID
	MembershipAnnouncement
	Joiner        JoinerAnnouncement
	AnnouncedByID insolar.ShortNodeID
}

func NewMemberAnnouncement(memberID insolar.ShortNodeID, mp MembershipProfile,
	announcerID insolar.ShortNodeID) MemberAnnouncement {

	return MemberAnnouncement{
		MemberID:               memberID,
		MembershipAnnouncement: NewMembershipAnnouncement(mp),
		AnnouncedByID:          announcerID,
	}
}

func NewJoinerAnnouncement(brief StaticProfile,
	announcerID insolar.ShortNodeID) MemberAnnouncement {

	// TODO joiner secret
	return MemberAnnouncement{
		MemberID:               brief.GetStaticNodeID(),
		MembershipAnnouncement: NewMembershipAnnouncement(NewMembershipProfileForJoiner(brief)),
		AnnouncedByID:          announcerID,
		Joiner: JoinerAnnouncement{
			JoinerProfile:  brief,
			IntroducedByID: announcerID,
		},
	}
}

func NewJoinerIDAnnouncement(joinerID, announcerID insolar.ShortNodeID) MemberAnnouncement {

	return MemberAnnouncement{
		MemberID:      joinerID,
		AnnouncedByID: announcerID,
	}
}

func NewMemberAnnouncementWithLeave(memberID insolar.ShortNodeID, mp MembershipProfile, leaveReason uint32,
	announcerID insolar.ShortNodeID) MemberAnnouncement {

	return MemberAnnouncement{
		MemberID:               memberID,
		MembershipAnnouncement: NewMembershipAnnouncementWithLeave(mp, leaveReason),
		AnnouncedByID:          announcerID,
	}
}

func NewMemberAnnouncementWithJoinerID(memberID insolar.ShortNodeID, mp MembershipProfile,
	joinerID insolar.ShortNodeID, joinerSecret cryptkit.DigestHolder,
	announcerID insolar.ShortNodeID) MemberAnnouncement {

	return MemberAnnouncement{
		MemberID:               memberID,
		MembershipAnnouncement: NewMembershipAnnouncementWithJoinerID(mp, joinerID, joinerSecret),
		AnnouncedByID:          announcerID,
	}
}

func NewMemberAnnouncementWithJoiner(memberID insolar.ShortNodeID, mp MembershipProfile, joiner JoinerAnnouncement,
	announcerID insolar.ShortNodeID) MemberAnnouncement {

	return MemberAnnouncement{
		MemberID: memberID,
		MembershipAnnouncement: NewMembershipAnnouncementWithJoinerID(mp,
			joiner.JoinerProfile.GetStaticNodeID(), joiner.JoinerSecret),
		Joiner:        joiner,
		AnnouncedByID: announcerID,
	}
}

func NewMembershipAnnouncement(mp MembershipProfile) MembershipAnnouncement {
	return MembershipAnnouncement{
		Membership: mp,
	}
}

func NewMembershipAnnouncementWithJoinerID(mp MembershipProfile,
	joinerID insolar.ShortNodeID, joinerSecret cryptkit.DigestHolder) MembershipAnnouncement {

	return MembershipAnnouncement{
		Membership:   mp,
		JoinerID:     joinerID,
		JoinerSecret: joinerSecret,
	}
}

func NewMembershipAnnouncementWithLeave(mp MembershipProfile, leaveReason uint32) MembershipAnnouncement {
	return MembershipAnnouncement{
		Membership:  mp,
		IsLeaving:   true,
		LeaveReason: leaveReason,
	}
}
