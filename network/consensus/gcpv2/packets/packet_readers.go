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
	"github.com/insolar/insolar/network/consensus/common"
	common2 "github.com/insolar/insolar/network/consensus/gcpv2/common"
	"github.com/insolar/insolar/network/consensus/gcpv2/nodeset"
)

//go:generate minimock -i github.com/insolar/insolar/network/consensus/gcpv2/packets.PacketParser -o ../../testutils -s _mock.go

type PacketParser interface {
	GetPacketType() PacketType

	/* Should return UnknownPulseNumber when it is not possible to identify it for a packet */
	GetPulseNumber() common.PulseNumber

	GetPulsePacket() PulsePacketReader
	GetMemberPacket() MemberPacketReader

	GetSourceId() common.ShortNodeID
	GetReceiverId() common.ShortNodeID

	/* Returns zero when no relay */
	GetRelayTargetID() common.ShortNodeID

	GetPacketSignature() common.SignedDigest
}

type PulsePacketReader interface {
	// GetPulsarId() PulsarId
	GetPulseData() common.PulseData
	GetPulseDataEvidence() common2.OriginalPulsarPacket
}

type MemberPacketReader interface {
	GetPacketType() PacketType

	AsPhase0Packet() Phase0PacketReader
	AsPhase1Packet() Phase1PacketReader
	AsPhase2Packet() Phase2PacketReader
	AsPhase3Packet() Phase3PacketReader

	GetPacketSignature() common.SignedDigest
}

type PhasePacketReader interface {
	GetPulseNumber() common.PulseNumber
}

type Phase0PacketReader interface {
	PhasePacketReader

	GetNodeRank() common2.MembershipRank
	GetEmbeddedPulsePacket() PulsePacketReader
}

type Phase1PacketReader interface {
	PhasePacketReader

	HasPulseData() bool /* PulseData/PulsarData is optional for Phase1 */
	GetEmbeddedPulsePacket() PulsePacketReader

	//HasSelfIntro() bool
	GetCloudIntroduction() CloudIntroductionReader
	GetFullIntroduction() FullIntroductionReader

	GetAnnouncementReader() MembershipAnnouncementReader
	//GetNodeBriefIntro()
	//GetNodeFullIntro()
	//GetCloudIntro()
}

type Phase2PacketReader interface {
	PhasePacketReader

	GetBriefIntroduction() BriefIntroductionReader
	GetAnnouncementReader() MembershipAnnouncementReader
	GetNeighbourhood() []MembershipAnnouncementReader

	//GetNodeBriefIntro()
}

type Phase3PacketReader interface {
	PhasePacketReader

	GetBitset() nodeset.NodeBitset
	GetTrustedGsh() common2.GlobulaStateHash
	GetDoubtedGsh() common2.GlobulaStateHash

	GetTrustedCshEvidence() common.SignedEvidenceHolder
	GetDoubtedCshEvidence() common.SignedEvidenceHolder
}

type MembershipAnnouncementReader interface {
	GetNodeID() common.ShortNodeID
	GetNodeRank() common2.MembershipRank
	GetRequestedPower() common2.MemberPower
	GetNodeStateHashEvidence() common2.NodeStateHashEvidence
	GetAnnouncementSignature() common2.MemberAnnouncementSignature

	// Methods below are not applicable when GetNodeRank().IsJoiner()
	IsLeaving() bool
	GetLeaveReason() uint32

	/*
		If GetJoinerID() == 0 then there is no joiner announced by the member
		If this reader is part of Neighbourhood then nonzero GetJoinerID() will be equal to GetNodeID()
	*/
	GetJoinerID() common.ShortNodeID
	/* Can be nil when this reader is part of Neighbourhood - then joiner data is in the sender's announcmenet */
	GetJoinerAnnouncement() JoinerAnnouncementReader
}

type JoinerAnnouncementReader interface {
	GetBriefIntro() BriefIntroductionReader
	GetBriefIntroSignature() common.SignatureHolder
}

type CloudIntroductionReader interface {
	GetLastCloudStateHash() common.DigestHolder
	GetJoinerSecret() common.DigestHolder
	GetCloudIdentity() common.DigestHolder
}

type BriefIntroductionReader interface {
	common2.BriefCandidateProfile
}

type FullIntroductionReader interface {
	common2.CandidateProfile
}
