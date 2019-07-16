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
	"github.com/insolar/insolar/network/consensus/common/endpoints"
	"github.com/insolar/insolar/network/consensus/common/pulse"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/member"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/power"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/proofs"
	"time"
)

//go:generate minimock -i github.com/insolar/insolar/network/consensus/gcpv2/api/profiles.Host -o . -s _mock.go

type Host interface {
	GetDefaultEndpoint() endpoints.Outbound
	GetPublicKeyStore() cryptkit.PublicKeyStore
	IsAcceptableHost(from endpoints.Inbound) bool
	// GetHostType()
}

type NodeIntroduction interface {
	// full intro
	GetShortNodeID() insolar.ShortNodeID
	GetReference() insolar.Reference
	IsAllowedPower(p member.Power) bool
	ConvertPowerRequest(request power.Request) member.Power
}

//go:generate minimock -i github.com/insolar/insolar/network/consensus/gcpv2/api/profiles.StaticProfile -o . -s _mock.go

type staticProfile interface {
	GetShortNodeID() insolar.ShortNodeID
	GetPrimaryRole() member.PrimaryRole
	GetSpecialRoles() member.SpecialRole
	GetNodePublicKey() cryptkit.SignatureKeyHolder
	GetStartPower() member.Power
}

type StaticProfile interface { //brief intro
	Host
	staticProfile
	GetAnnouncementSignature() cryptkit.SignatureHolder

	GetIntroduction() NodeIntroduction // must be always be not null for LocalNode, full intro, == nil when has no full
}

type BaseNode interface {
	StaticProfile
	//GetShortNodeID() insolar.ShortNodeID

	/*
		As dynamic nodes may update static part info, code inside consenus logic MUST access static profile
		by getting it GetStatic() to ensure consistency among attributes
	*/
	GetStatic() StaticProfile
	GetSignatureVerifier() cryptkit.SignatureVerifier
	GetOpMode() member.OpMode
}

//go:generate minimock -i github.com/insolar/insolar/network/consensus/gcpv2/api/profiles.ActiveNode -o . -s _mock.go

type ActiveNode interface {
	BaseNode
	GetIndex() member.Index
	IsJoiner() bool
	GetDeclaredPower() member.Power
}

type EvictedNode interface {
	BaseNode
	GetLeaveReason() uint32
}

type BriefCandidateProfile interface {
	staticProfile

	GetDefaultEndpoint() endpoints.Outbound
	GetJoinerSignature() cryptkit.SignatureHolder
}

//go:generate minimock -i github.com/insolar/insolar/network/consensus/gcpv2/api/profiles.CandidateProfile -o . -s _mock.go

type CandidateProfile interface {
	BriefCandidateProfile

	GetPowerLevels() member.PowerSet
	GetExtraEndpoints() []endpoints.Outbound

	GetReference() insolar.Reference
	// NodeRefProof	[]common.Bits512

	GetIssuedAtPulse() pulse.Number // =0 when a node was connected during zeronet
	GetIssuedAtTime() time.Time
	GetIssuerID() insolar.ShortNodeID
	GetIssuerSignature() cryptkit.SignatureHolder
}

type Factory interface {
	CreateBriefIntroProfile(candidate BriefCandidateProfile) StaticProfile
	/* This method MUST: (1) ensure same values of both params; (2) create a new copy of StaticProfile */
	CreateFullIntroProfile(candidate CandidateProfile) StaticProfile
}

//go:generate minimock -i github.com/insolar/insolar/network/consensus/gcpv2/api/profiles.LocalNode -o . -s _mock.go

type LocalNode interface {
	ActiveNode
	LocalNodeProfile()
}

type Updatable interface {
	ActiveNode

	AsActiveNode() ActiveNode

	SetOpMode(m member.OpMode)
	SetPower(declaredPower member.Power)
	SetRank(index member.Index, m member.OpMode, declaredPower member.Power)
	SetSignatureVerifier(verifier cryptkit.SignatureVerifier)
	// Update certificate / mandate

	SetOpModeAndLeaveReason(exitCode uint32)
	GetLeaveReason() uint32
	SetIndex(index member.Index)
}

type MembershipProfile struct {
	Index          member.Index
	Mode           member.OpMode
	Power          member.Power
	RequestedPower member.Power
	proofs.NodeAnnouncedState
}

// TODO support joiner in MembershipProfile
//func (v MembershipProfile) IsJoiner() bool {
//
//}

func NewMembershipProfile(mode member.OpMode, power member.Power, index member.Index,
	nsh proofs.NodeStateHashEvidence, nas proofs.MemberAnnouncementSignature,
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

func NewMembershipProfileByNode(np ActiveNode, nsh proofs.NodeStateHashEvidence, nas proofs.MemberAnnouncementSignature,
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

func (p MembershipProfile) AsRank(nc int) member.Rank {
	if p.Index.IsJoiner() {
		return member.JoinerRank
	}
	return member.NewMembershipRank(p.Mode, p.Power, p.Index, member.AsIndex(nc).AsUint16())
}

func (p MembershipProfile) AsRankUint16(nc uint16) member.Rank {
	if p.Index.IsJoiner() {
		return member.JoinerRank
	}
	return member.NewMembershipRank(p.Mode, p.Power, p.Index, member.Index(nc).AsUint16())
}

func (p MembershipProfile) Equals(o MembershipProfile) bool {
	if p.Index != o.Index || p.Power != o.Power || p.IsEmpty() || o.IsEmpty() || p.RequestedPower != o.RequestedPower {
		return false
	}

	if p.StateEvidence != o.StateEvidence {
		if !p.StateEvidence.GetNodeStateHash().Equals(o.StateEvidence.GetNodeStateHash()) {
			return false
		}
		if !p.StateEvidence.GetGlobulaNodeStateSignature().Equals(o.StateEvidence.GetGlobulaNodeStateSignature()) {
			return false
		}
	}

	return p.AnnounceSignature == o.AnnounceSignature || p.AnnounceSignature.Equals(o.AnnounceSignature)
}

func (p MembershipProfile) StringParts() string {
	if p.Power == p.RequestedPower {
		return fmt.Sprintf("pw:%v se:%v cs:%v", p.Power, p.StateEvidence, p.AnnounceSignature)
	}

	return fmt.Sprintf("pw:%v->%v se:%v cs:%v", p.Power, p.RequestedPower, p.StateEvidence, p.AnnounceSignature)
}

func (p MembershipProfile) String() string {
	return fmt.Sprintf("idx:%03d %s", p.Index, p.StringParts())
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

func EqualIntroProfiles(p StaticProfile, o StaticProfile) bool {
	if p == nil || o == nil {
		return false
	}
	if p == o {
		return true
	}

	if p.GetShortNodeID() != o.GetShortNodeID() || p.GetPrimaryRole() != o.GetPrimaryRole() ||
		p.GetSpecialRoles() != o.GetSpecialRoles() || p.GetStartPower() != o.GetStartPower() ||
		!p.GetNodePublicKey().Equals(o.GetNodePublicKey()) {
		return false
	}

	return endpoints.EqualEndpoints(p.GetDefaultEndpoint(), o.GetDefaultEndpoint())
}

func MatchIntroAndRank(np ActiveNode, nc int, nr member.Rank) bool {
	if nr.IsJoiner() {
		return np.IsJoiner()
	}

	return int(nr.GetTotalCount()) == nc && nr.GetIndex() == np.GetIndex() && nr.GetMode() == np.GetOpMode() &&
		nr.GetPower() == np.GetDeclaredPower()
}
