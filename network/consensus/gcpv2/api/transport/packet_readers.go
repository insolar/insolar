// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package transport

import (
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/network/consensus/common/cryptkit"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/member"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/phases"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/profiles"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/proofs"
	"github.com/insolar/insolar/pulse"
)

//go:generate minimock -i github.com/insolar/insolar/network/consensus/gcpv2/api/transport.PacketParser -o . -s _mock.go -g

type PacketParser interface {
	GetPacketType() phases.PacketType
	/* Should return Unknown when it is not possible to identify it for a packet */
	GetPulseNumber() pulse.Number

	GetSourceID() insolar.ShortNodeID
	GetReceiverID() insolar.ShortNodeID
	GetTargetID() insolar.ShortNodeID
	IsRelayForbidden() bool

	GetPacketSignature() cryptkit.SignedDigest
	ParsePacketBody() (PacketParser, error) // enables lazy parsing / parsing only after packet validation

	GetPulsePacket() PulsePacketReader
	GetMemberPacket() MemberPacketReader
}

type PulsePacketReader interface {
	// GetPulsarId() PulsarId
	GetPulseData() pulse.Data
	GetPulseDataEvidence() proofs.OriginalPulsarPacket
}

type MemberPacketReader interface {
	GetPacketType() phases.PacketType

	AsPhase0Packet() Phase0PacketReader
	AsPhase1Packet() Phase1PacketReader
	AsPhase2Packet() Phase2PacketReader
	AsPhase3Packet() Phase3PacketReader

	GetPacketSignature() cryptkit.SignedDigest
}

type PhasePacketReader interface {
	GetPulseNumber() pulse.Number
}

type Phase0PacketReader interface {
	PhasePacketReader

	GetNodeRank() member.Rank
	GetEmbeddedPulsePacket() PulsePacketReader
}

type ExtendedIntroReader interface {
	HasFullIntro() bool
	GetFullIntroduction() FullIntroductionReader

	HasCloudIntro() bool
	GetCloudIntroduction() CloudIntroductionReader

	HasJoinerSecret() bool
	GetJoinerSecret() cryptkit.DigestHolder
}

type AnnouncementPacketReader interface {
	ExtendedIntroReader
	GetAnnouncementReader() MembershipAnnouncementReader
}

type Phase1PacketReader interface {
	PhasePacketReader
	AnnouncementPacketReader

	HasPulseData() bool /* PulseData/PulsarData is optional for Phase1 */
	GetEmbeddedPulsePacket() PulsePacketReader
}

type Phase2PacketReader interface {
	PhasePacketReader
	AnnouncementPacketReader

	GetBriefIntroduction() BriefIntroductionReader
	GetNeighbourhood() []MembershipAnnouncementReader
}

type HashedVectorReader interface {
	GetBitset() member.StateBitset

	GetTrustedGlobulaAnnouncementHash() proofs.GlobulaAnnouncementHash
	GetTrustedGlobulaStateSignature() proofs.GlobulaStateSignature
	GetTrustedExpectedRank() member.Rank

	GetDoubtedGlobulaAnnouncementHash() proofs.GlobulaAnnouncementHash
	GetDoubtedGlobulaStateSignature() proofs.GlobulaStateSignature
	GetDoubtedExpectedRank() member.Rank
}

type Phase3PacketReader interface {
	PhasePacketReader
	HashedVectorReader
}

type MembershipAnnouncementReader interface {
	GetNodeID() insolar.ShortNodeID
	GetNodeRank() member.Rank
	GetRequestedPower() member.Power
	GetNodeStateHashEvidence() proofs.NodeStateHashEvidence
	GetAnnouncementSignature() proofs.MemberAnnouncementSignature

	// Methods below are not applicable when GetNodeRank().IsJoiner()
	IsLeaving() bool
	GetLeaveReason() uint32

	/*
		If GetJoinerID() == 0 then there is no joiner announced by the member
		If this reader is part of Neighbourhood then nonzero GetJoinerID() will be equal to GetNodeID()
	*/
	GetJoinerID() insolar.ShortNodeID

	/* Can be nil when this reader is part of Neighbourhood - then joiner data is in the sender's announcement */
	GetJoinerAnnouncement() JoinerAnnouncementReader
}

type JoinerAnnouncementReader interface {
	GetJoinerIntroducedByID() insolar.ShortNodeID
	GetBriefIntroduction() BriefIntroductionReader

	HasFullIntro() bool
	GetFullIntroduction() FullIntroductionReader
}

type CloudIntroductionReader interface {
	GetLastCloudStateHash() cryptkit.DigestHolder
	// GetAnnouncedJoinerSecret() cryptkit.DigestHolder
	GetCloudIdentity() cryptkit.DigestHolder
}

type BriefIntroductionReader interface {
	profiles.BriefCandidateProfile
}

//go:generate minimock -i github.com/insolar/insolar/network/consensus/gcpv2/api/transport.FullIntroductionReader -o . -s _mock.go -g
type FullIntroductionReader interface {
	profiles.CandidateProfile
}
