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

package packets

import (
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/network/consensus/common/cryptography_containers"
	"github.com/insolar/insolar/network/consensus/common/long_bits"
	"github.com/insolar/insolar/network/consensus/common/pulse_data"
	"github.com/insolar/insolar/network/consensus/gcpv2/gcp_types"
)

//go:generate minimock -i github.com/insolar/insolar/network/consensus/gcpv2/packets.PacketParser -o ../packets -s _mock.go

type PacketParser interface {
	GetPacketType() gcp_types.PacketType

	/* Should return UnknownPulseNumber when it is not possible to identify it for a packet */
	GetPulseNumber() pulse_data.PulseNumber

	GetPulsePacket() PulsePacketReader
	GetMemberPacket() MemberPacketReader

	GetSourceID() insolar.ShortNodeID
	GetReceiverID() insolar.ShortNodeID

	IsRelayForbidden() bool
	GetTargetID() insolar.ShortNodeID

	GetPacketSignature() cryptography_containers.SignedDigest
}

type PulsePacketReader interface {
	// GetPulsarId() PulsarId
	GetPulseData() pulse_data.PulseData
	GetPulseDataEvidence() OriginalPulsarPacket
}

type MemberPacketReader interface {
	GetPacketType() gcp_types.PacketType

	AsPhase0Packet() Phase0PacketReader
	AsPhase1Packet() Phase1PacketReader
	AsPhase2Packet() Phase2PacketReader
	AsPhase3Packet() Phase3PacketReader

	GetPacketSignature() cryptography_containers.SignedDigest
}

type PhasePacketReader interface {
	GetPulseNumber() pulse_data.PulseNumber
}

type Phase0PacketReader interface {
	PhasePacketReader

	GetNodeRank() gcp_types.MembershipRank
	GetEmbeddedPulsePacket() PulsePacketReader
}

type ExtendedIntroReader interface {
	HasFullIntro() bool
	GetFullIntroduction() FullIntroductionReader

	HasCloudIntro() bool
	GetCloudIntroduction() CloudIntroductionReader

	HasJoinerSecret() bool
	GetJoinerSecret() cryptography_containers.SignatureHolder
}

type Phase1PacketReader interface {
	PhasePacketReader
	ExtendedIntroReader

	HasPulseData() bool /* PulseData/PulsarData is optional for Phase1 */
	GetEmbeddedPulsePacket() PulsePacketReader

	GetAnnouncementReader() MembershipAnnouncementReader
}

type Phase2PacketReader interface {
	PhasePacketReader
	ExtendedIntroReader

	GetBriefIntroduction() BriefIntroductionReader
	GetAnnouncementReader() MembershipAnnouncementReader
	GetNeighbourhood() []MembershipAnnouncementReader
}

type Phase3PacketReader interface {
	PhasePacketReader

	GetBitset() gcp_types.NodeBitset

	//GetTrustedExpectedRank() common2.MembershipRank
	GetTrustedGlobulaAnnouncementHash() gcp_types.GlobulaAnnouncementHash
	GetTrustedGlobulaStateSignature() gcp_types.GlobulaStateSignature
	GetTrustedExpectedRank() gcp_types.MembershipRank

	//GetDoubtedExpectedRank() common2.MembershipRank
	GetDoubtedGlobulaAnnouncementHash() gcp_types.GlobulaAnnouncementHash
	GetDoubtedGlobulaStateSignature() gcp_types.GlobulaStateSignature
	GetDoubtedExpectedRank() gcp_types.MembershipRank
}

type MembershipAnnouncementReader interface {
	GetNodeID() insolar.ShortNodeID
	GetNodeRank() gcp_types.MembershipRank
	GetRequestedPower() gcp_types.MemberPower
	GetNodeStateHashEvidence() gcp_types.NodeStateHashEvidence
	GetAnnouncementSignature() gcp_types.MemberAnnouncementSignature

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
	GetBriefIntro() BriefIntroductionReader
	GetBriefIntroSignature() cryptography_containers.SignatureHolder
}

type CloudIntroductionReader interface {
	GetLastCloudStateHash() cryptography_containers.DigestHolder
	GetJoinerSecret() cryptography_containers.DigestHolder
	GetCloudIdentity() cryptography_containers.DigestHolder
}

type BriefIntroductionReader interface {
	gcp_types.BriefCandidateProfile
}

//go:generate minimock -i github.com/insolar/insolar/network/consensus/gcpv2/packets.FullIntroductionReader -o ../packets -s _mock.go
type FullIntroductionReader interface {
	gcp_types.CandidateProfile
}

type OriginalPulsarPacket interface {
	long_bits.FixedReader
	OriginalPulsarPacket()
}
