// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

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
