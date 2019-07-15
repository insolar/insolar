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

//go:generate minimock -i github.com/insolar/insolar/network/consensus/gcpv2/api/profiles.NodeIntroProfile -o . -s _mock.go

type nodeIntroProfile interface {
	GetShortNodeID() insolar.ShortNodeID
	GetPrimaryRole() member.PrimaryRole
	GetSpecialRoles() member.SpecialRole
	GetNodePublicKey() cryptkit.SignatureKeyHolder
	GetStartPower() member.Power
}

type NodeIntroProfile interface { //brief intro
	Host
	nodeIntroProfile
	GetAnnouncementSignature() cryptkit.SignatureHolder

	HasIntroduction() bool             // must be always true for LocalNode
	GetIntroduction() NodeIntroduction // not null, full intro, will panic when HasIntroduction() == false
}

type BaseNode interface {
	// TODO Rename
	NodeIntroProfile
	GetSignatureVerifier() cryptkit.SignatureVerifier
	GetOpMode() member.OpMode
}

const NodeIndexBits = 10 // DO NOT change it, otherwise nasty consequences will come
const NodeIndexMask = 1<<NodeIndexBits - 1
const MaxNodeIndex = NodeIndexMask

//go:generate minimock -i github.com/insolar/insolar/network/consensus/gcpv2/api/profiles.ActiveNode -o . -s _mock.go

type ActiveNode interface {
	// TODO Rename
	BaseNode
	GetIndex() int // 0 for joiners
	IsJoiner() bool
	GetDeclaredPower() member.Power
}

type EvictedNode interface {
	// TODO Rename
	BaseNode
	GetLeaveReason() uint32
}

type BriefCandidateProfile interface {
	nodeIntroProfile

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
	CreateBriefIntroProfile(candidate BriefCandidateProfile) NodeIntroProfile
	/* This method MUST: (1) ensure same values of both params; (2) create a new copy of NodeIntroProfile */
	CreateFullIntroProfile(candidate CandidateProfile) NodeIntroProfile
}

//go:generate minimock -i github.com/insolar/insolar/network/consensus/gcpv2/api/profiles.LocalNode -o . -s _mock.go

type LocalNode interface {
	ActiveNode
	LocalNodeProfile()
}

type Updatable interface {
	ActiveNode
	SetOpMode(m member.OpMode)
	SetPower(declaredPower member.Power)
	SetRank(index int, m member.OpMode, declaredPower member.Power)
	SetSignatureVerifier(verifier cryptkit.SignatureVerifier)
	// Update certificate / mandate

	SetOpModeAndLeaveReason(exitCode uint32)
	GetLeaveReason() uint32
	SetIndex(index int)
}

type MembershipProfile struct {
	Index          uint16
	Mode           member.OpMode
	Power          member.Power
	RequestedPower member.Power
	proofs.NodeAnnouncedState
}

func NewMembershipProfile(mode member.OpMode, power member.Power, index uint16, nsh proofs.NodeStateHashEvidence, nas proofs.MemberAnnouncementSignature,
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

	return NewMembershipProfile(np.GetOpMode(), np.GetDeclaredPower(),
		uint16(np.GetIndex()), nsh, nas, ep)
}

func (p MembershipProfile) IsEmpty() bool {
	return p.StateEvidence == nil || p.AnnounceSignature == nil
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
	Joiner      NodeIntroProfile
}

func NewMembershipAnnouncement(mp MembershipProfile) MembershipAnnouncement {
	return MembershipAnnouncement{
		Membership: mp,
	}
}

func NewMembershipAnnouncementWithJoiner(mp MembershipProfile, joiner NodeIntroProfile) MembershipAnnouncement {
	return MembershipAnnouncement{
		Membership: mp,
		Joiner:     joiner,
	}
}

func NewMembershipAnnouncementWithLeave(mp MembershipProfile, leaveReason uint32) MembershipAnnouncement {
	return MembershipAnnouncement{
		Membership:  mp,
		IsLeaving:   true,
		LeaveReason: leaveReason,
	}
}

func EqualIntroProfiles(p NodeIntroProfile, o NodeIntroProfile) bool {
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
