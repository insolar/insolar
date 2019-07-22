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

type MembershipAnnouncement struct {
	Membership  MembershipProfile
	IsLeaving   bool
	LeaveReason uint32
	JoinerID    insolar.ShortNodeID
}

func NewMembershipAnnouncement(mp MembershipProfile) MembershipAnnouncement {
	return MembershipAnnouncement{
		Membership: mp,
	}
}

func NewMembershipAnnouncementWithJoinerID(mp MembershipProfile, joinerID insolar.ShortNodeID) MembershipAnnouncement {
	return MembershipAnnouncement{
		Membership: mp,
		JoinerID:   joinerID,
	}
}

func NewMembershipAnnouncementWithLeave(mp MembershipProfile, leaveReason uint32) MembershipAnnouncement {
	return MembershipAnnouncement{
		Membership:  mp,
		IsLeaving:   true,
		LeaveReason: leaveReason,
	}
}
