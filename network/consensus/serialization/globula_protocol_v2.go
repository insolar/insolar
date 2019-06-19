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

package serialization

import (
	"github.com/insolar/insolar/network/consensus/common"
)

type GlobulaConsensusProtocolV2Packet struct {
	Header      UnifiedProtocolPacketHeader `insolar-transport:"Protocol=0x01;Packet=0-4"` // ByteSize<=16
	PulseNumber common.PulseNumber          `insolar-transport:"[30-31]=0"`                // [30-31] MUST = 0, ByteSize=4

	// Phases 0-2
	NodeRank          SerializedGlobulaNodeRank `insolar-transport:"Packet=0,1,2"`                       // ByteSize=4, current rank
	GlobulaNodeState  *EmbeddedGlobulaNodeState `insolar-transport:"Packet=1,2"`                         // ByteSize=128
	SelfIntroToJoiner *IntroductionToJoiner     `insolar-transport:"optional=PacketFlags[2];Packet=1,2"` // ByteSize=?
	PulsarPacket      *EmbeddedPulsarData       `insolar-transport:"optional=PacketFlags[1];Packet=0,1"` // ByteSize>=124

	// Phase 3
	/*
		GlobulaNodeBitset is a 5-state bitset, each node has a state at the same index as was given in rank.
		Node have following states:
		0 - z-value (same as missing value) Trusted node
		1 - Doubted node
		2 -
		3 - Fraud node
		4 - Missing node
	*/
	GlobulaNodeBitset  *NodeApperanceBitset `insolar-transport:"Packet=3"`                                                          // ByteSize=1..335
	TrustedStateVector *GlobulaStateVector  `insolar-transport:"Packet=3;TrustedStateVector.ExpectedRank[28-31]=flags:Phase3Flags"` // ByteSize=96
	DoubtedStateVector *GlobulaStateVector  `insolar-transport:"optional=Phase3Flags[0];Packet=3"`                                  // ByteSize=96
	// FraudStateVector *GlobulaStateVector `insolar-transport:"optional=Phase3Flags[1];Packet=3"` //ByteSize=96

	// Claim Section
	LeaveClaim *LeaveAnnouncementClaim `insolar-transport:"Packet=1"` // ByteSize=5, exclusive
	JoinClaim  *JoinRequestClaim       `insolar-transport:"Packet=1"` // ByteSize=177, exclusive

	Neighbourhood []NodeNeighbourClaim `insolar-transport:"Packet=2"` // ByteSize=N*[197,374]

	// 	ReferendumVotes []ReferendumVote `insolar-transport:"Packet=3"`

	// End Of Packet
	EndOfClaims     EmptyClaim     // ByteSize=1 - indicates end of claims
	PacketSignature common.Bits512 // ByteSize=64
}

/*
	Phase0 size: >=148
	Phase1 size: >=405 normal, >=645 to joiner, >=822 to joiner with JoinClaim :: w/o pulse data -124
	Phase2 size: 217 + N*[197,374] ... 1800 byte => (8+self) members .. 4+2 .. 4 joining neighbours
	Phase3 size: >=218 <=728

	Network traffic 1001 nodes:
				     IN          OUT
		Phase0: <    148 000 	148 000
		Phase1: <    645 000    645 000
		Phase2: <  1 600 000  1 600 000    //neighbourhood = 5-7
		Phase3: <    728 000 	728 000

*/

type EmbeddedPulsarData struct {
	// ByteSize>=124
	Header UnifiedProtocolPacketHeader // ByteSize=16

	// PulseNumber common.PulseNumber //available externally
	PulsarPulsePacketExt // ByteSize>=108
}

type EmbeddedGlobulaNodeState struct {
	// ByteSize=128
	// PulseDataHash            common.Bits256 //available externally
	// FoldedLastCloudStateHash common.Bits224 //available externally
	// NodeRank                 GlobulaNodeRank //available externally

	NodeStateHash             common.Bits512 // ByteSize=64
	GlobulaNodeStateSignature common.Bits512 // ByteSize=64, :=Sign(NodePK, Merkle512(NodeStateHash, (LastCloudStateHash.FoldTo224() << 32 | GlobulaNodeRank)))
}

type GlobulaStateVector struct {
	// ByteSize=96
	/*
		GlobulaStateHash = merkle(GlobulaNodeStateSignature of all nodes of this vector)
		SignedGlobulaStateHash = sign(GlobulaStateHash, SK(sending node))
	*/
	SignedGlobulaStateHash common.Bits512            // ByteSize=64
	Phase1SignaturesHash   common.Bits224            // ByteSize=28
	ExpectedRank           SerializedGlobulaNodeRank // ByteSize=4
}

type NodeApperanceBitset struct {
	// ByteSize=[1..252]
	FlagsAndLoLength uint8 // [00-05] LoByteLength, [06] Compressed, [07] HasHiLength (to be compatible with Protobuf VarInt)
	HiLength         uint8 // [00-06] HiByteLength, [07] MUST = 0 (to be compatible with Protobuf VarInt)
	bytes            []byte
}

// type CloudIdentity common.Bits512 //ByteSize=64

type IntroductionToJoiner struct {
	// ByteSize=240
	// CloudIdentity CloudIdentity //ByteSize=64
	LastCloudStateHash CloudStateHash   // ByteSize=64
	SelfIntro          NodeIntroduction // ByteSize=176
}

type NodeBriefIntro struct {
	// ByteSize=142
	ProtocolVersionAndFlags uint16
	ShortNodeId             common.ShortNodeID
	InboundRelayId          common.ShortNodeID // =0 - no relay is needed to send packets to the node
	IssuedAt                common.PulseNumber
	NodePK                  common.Bits512 // works as a unique node identity
	NodeSignature           common.Bits512
}

type NodeExtIntro struct {
	// ByteSize>=128
	NodeIntroHash          common.Bits256
	IssuedByPKFolded       common.Bits256 // DiscoveryNode
	NodeReferenceWithProof []common.Bits512
	IssuerSignature        common.Bits512
}

type NodeIntroduction struct {
	// ByteSize=176
	ProtocolVersionAndFlags uint16
	Reserved0               uint8
	ValidAsRequestFor       uint8 // for how long this intro can be used for joining, but it is not validity of Intro packet itself
	ShortNodeId             common.ShortNodeID
	InboundRelayId          common.ShortNodeID // =0 - no relay is needed to send packets to the node
	IssuedAt                common.PulseNumber
	NodePK                  common.Bits512 // works as a unique node identity
	NodeSignature           common.Bits512
}

type JoinerMandate struct {
	NodeIntroHash          common.Bits256
	IssuedByPK             common.Bits512 // DiscoveryNode
	NodeReferenceWithProof []common.Bits512
	IssuerSignature        common.Bits512
}

type ClaimHeader struct {
	TypeAndLength uint16 `insolar-transport:"header;[0-9]=length;[10-15]=header:ClaimType;group=Claims"` // [00-09] ByteLength [10-15] ClaimClass
	// actual payload
}

type GenericClaim struct {
	// ByteSize>=1
	ClaimHeader
	Payload []byte
}

type EmptyClaim struct {
	// ByteSize=1
	ClaimHeader `insolar-transport:"delimiter;ClaimType=0;length=header"`
}

type LeaveAnnouncementClaim struct {
	// ByteSize=5
	ClaimHeader `insolar-transport:"exclusive;ClaimType=1"` // Must be the only claim per packet, identified by len<128
	LeaveReason uint32
}

type JoinRequestClaim struct {
	// ByteSize=177
	ClaimHeader `insolar-transport:"exclusive;ClaimType=2"` // Must be the only per packet, identified by len>128
	Intro       NodeIntroduction
}

type JoinerMandateClaim struct {
	ClaimHeader   `insolar-transport:"exclusive;ClaimType=3"` // Must be the only per packet, identified by len>128
	JoinerMandate JoinerMandate
}

type NodeNeighbourClaim struct {
	// ByteSize=[197,374]
	ClaimHeader `insolar-transport:"ClaimType=4"`

	NodeRank         SerializedGlobulaNodeRank `insolar-transport:"[30-31]=flags:nodeRank"` // ByteSize=4
	GlobulaNodeState EmbeddedGlobulaNodeState  // ByteSize=128

	// As claimClass=1 is exclusive, this packet signature is of signature claims
	Phase1PacketSignature *common.Bits512 `insolar-transport:"optional=nodeRank[31]"` // ByteSize=64

	// Only claimClass=1 is allowed here
	LeaveClaim    *LeaveAnnouncementClaim // ByteSize=5
	JoinClaim     *JoinRequestClaim       // ByteSize=177
	JoinerMandate *JoinerMandateClaim

	// EndOfClaims EmptyClaim - not included, end of claims identified by len of NodeNeighbourClaim
}

type PowerActivationClaim struct {
	// ByteSize>64
	ClaimHeader   `insolar-transport:"claimClass=2"`
	NodeReference common.Bits512
	// merkle proofs?
	// or signed by HME
}
