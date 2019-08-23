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
	"time"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/network/consensus/common/cryptkit"
	"github.com/insolar/insolar/network/consensus/common/endpoints"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/member"
	"github.com/insolar/insolar/pulse"
)

//go:generate minimock -i github.com/insolar/insolar/network/consensus/gcpv2/api/profiles.Host -o . -s _mock.go -g

type Host interface {
	GetDefaultEndpoint() endpoints.Outbound
	GetPublicKeyStore() cryptkit.PublicKeyStore
	IsAcceptableHost(from endpoints.Inbound) bool
	// GetHostType()
}

//go:generate minimock -i github.com/insolar/insolar/network/consensus/gcpv2/api/profiles.StaticProfileExtension -o . -s _mock.go -g

type StaticProfileExtension interface {
	GetIntroducedNodeID() insolar.ShortNodeID
	CandidateProfileExtension
}

type staticProfile interface {
	GetStaticNodeID() insolar.ShortNodeID
	GetPrimaryRole() member.PrimaryRole
	GetSpecialRoles() member.SpecialRole
	GetNodePublicKey() cryptkit.SignatureKeyHolder
	GetStartPower() member.Power
	GetBriefIntroSignedDigest() cryptkit.SignedDigestHolder
}

//go:generate minimock -i github.com/insolar/insolar/network/consensus/gcpv2/api/profiles.StaticProfile -o . -s _mock.go -g

type StaticProfile interface { // brief intro
	Host
	staticProfile
	GetExtension() StaticProfileExtension // must be always be not null for LocalNode, full intro, == nil when has no full
}

//go:generate minimock -i github.com/insolar/insolar/network/consensus/gcpv2/api/profiles.BaseNode -o . -s _mock.go -g

type BaseNode interface {
	// StaticProfile
	GetNodeID() insolar.ShortNodeID

	/*
		As dynamic nodes may update static part info, code inside consenus logic MUST access static profile
		by getting it GetStatic() to ensure consistency among attributes
	*/
	GetStatic() StaticProfile
	GetSignatureVerifier() cryptkit.SignatureVerifier
	GetOpMode() member.OpMode
}

//go:generate minimock -i github.com/insolar/insolar/network/consensus/gcpv2/api/profiles.ActiveNode -o . -s _mock.go -g

type ActiveNode interface {
	BaseNode
	GetIndex() member.Index
	IsJoiner() bool
	IsPowered() bool
	IsVoter() bool
	IsStateful() bool
	CanIntroduceJoiner() bool
	HasFullProfile() bool

	GetDeclaredPower() member.Power
}

type EvictedNode interface {
	BaseNode
	GetLeaveReason() uint32
}

type BriefCandidateProfile interface {
	staticProfile

	GetDefaultEndpoint() endpoints.Outbound
}

//go:generate minimock -i github.com/insolar/insolar/network/consensus/gcpv2/api/profiles.CandidateProfile -o . -s _mock.go -g

type CandidateProfileExtension interface {
	GetPowerLevels() member.PowerSet
	GetExtraEndpoints() []endpoints.Outbound
	GetReference() insolar.Reference
	// NodeRefProof	[]common.Bits512

	GetIssuedAtPulse() pulse.Number // =0 when a node was connected during zeronet
	GetIssuedAtTime() time.Time
	GetIssuerID() insolar.ShortNodeID
	GetIssuerSignature() cryptkit.SignatureHolder
}

type CandidateProfile interface {
	BriefCandidateProfile
	CandidateProfileExtension
}

//go:generate minimock -i github.com/insolar/insolar/network/consensus/gcpv2/api/profiles.Factory -o . -s _mock.go -g

type Factory interface {
	CreateFullIntroProfile(candidate CandidateProfile) StaticProfile
	CreateBriefIntroProfile(candidate BriefCandidateProfile) StaticProfile
	CreateUpgradableIntroProfile(candidate BriefCandidateProfile) StaticProfile
	TryConvertUpgradableIntroProfile(profile StaticProfile) (StaticProfile, bool)
}

//go:generate minimock -i github.com/insolar/insolar/network/consensus/gcpv2/api/profiles.LocalNode -o . -s _mock.go -g

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

	SetOpModeAndLeaveReason(index member.Index, exitCode uint32)
	GetLeaveReason() uint32
	SetIndex(index member.Index)
}

type Upgradable interface {
	UpgradeProfile(upgradeData CandidateProfileExtension) bool
}
