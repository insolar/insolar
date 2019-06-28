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
	common2 "github.com/insolar/insolar/network/consensus/gcpv2/common"
)

type GlobulaConsensusProtocolV2Packet struct {
	Header      UnifiedProtocolPacketHeader `insolar-transport:"Protocol=0x01;Packet=0-4"` // ByteSize<=16
	PulseNumber common.PulseNumber          `insolar-transport:"[30-31]=0"`                // [30-31] MUST = 0, ByteSize=4

	// Phases 0-2
	NodeRank         SerializedGlobulaNodeRank `insolar-transport:"Packet=0,1,2"` // ByteSize=4, current rank
	GlobulaNodeState *CompactGlobulaNodeState  `insolar-transport:"Packet=1,2"`   // ByteSize=128

	PulsarPacket *EmbeddedPulsarData `insolar-transport:"optional=PacketFlags[1];Packet=0,1"` // ByteSize>=124

	//These 3 fields below are mutually exclusive by consensus logic
	SelfIntroToJoiner *IntroductionToJoiner   `insolar-transport:"optional=PacketFlags[2];Packet=1,2"` // ByteSize=234
	SelfIntroOfJoiner *IntroductionOfJoiner   `insolar-transport:"optional=PacketFlags[2];Packet=1"`   // ByteSize>82
	MemberAnnounce    *MembershipAnnouncement `insolar-transport:"optional=PacketFlags[3];Packet=1"`   // ByteSize=41, 202, 214

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
	GlobulaNodeBitset  *NodeAppearanceBitset `insolar-transport:"Packet=3"`                                                          // ByteSize=1..335
	TrustedStateVector *GlobulaStateVector   `insolar-transport:"Packet=3;TrustedStateVector.ExpectedRank[30-31]=flags:Phase3Flags"` // ByteSize=96
	DoubtedStateVector *GlobulaStateVector   `insolar-transport:"optional=Phase3Flags[0];Packet=3"`                                  // ByteSize=96
	// FraudStateVector *GlobulaStateVector `insolar-transport:"optional=Phase3Flags[1];Packet=3"` //ByteSize=96

	// Claim Section
	Neighbourhood []NodeNeighbourClaim `insolar-transport:"Packet=2"` // ByteSize=N*[133, 174, 335, 347]

	// 	ReferendumVotes []ReferendumVote `insolar-transport:"Packet=3"`

	// End Of Packet
	EndOfClaims     EmptyClaim     // ByteSize=1 - indicates end of claims
	PacketSignature common.Bits512 // ByteSize=64
}

/*  SIZES ARE OUTDATED

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

type CompactGlobulaNodeState struct {
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
	SignedGlobulaStateHash common.Bits512 // ByteSize=64
	/*
		Hash(all MembershipAnnouncement.Signature of this vector).FoldTo224()
	*/
	MembershipAnnouncementHash common.Bits224            // ByteSize=28
	ExpectedRank               SerializedGlobulaNodeRank // ByteSize=4
}

type NodeAppearanceBitset struct {
	// ByteSize=[1..252]
	FlagsAndLoLength uint8 // [00-05] LoByteLength, [06] Compressed, [07] HasHiLength (to be compatible with Protobuf VarInt)
	HiLength         uint8 // [00-06] HiByteLength, [07] MUST = 0 (to be compatible with Protobuf VarInt)
	Bytes            []byte
}

// type CloudIdentity common.Bits512 //ByteSize=64

type IntroductionToJoiner struct {
	// ByteSize=234
	LastCloudStateHash CloudStateHash // ByteSize=64
	SelfIntro          NodeBriefIntro // ByteSize=102 | 106 | 118
	SelfIntroSignature common.Bits512 // ByteSize=64
}

type IntroductionOfJoiner struct {
	// ByteSize>82
	ExtraIntro NodeExtraIntro
}

type NodeBriefIntro struct {
	// ByteSize=6 + (4 | 6 | 18) + 92 = 102 | 106 | 118
	PrimaryRoleAndFlags uint8 `insolar-transport:"[0:5]=header:NodePrimaryRole;[6:7]=header:AddrMode"` //AddrMode =0 reserved, =1 Relay, =2 IPv4 =3 IPv6
	SpecialRoles        common2.NodeSpecialRole
	ShortID             common.ShortNodeID

	// 4 | 6 | 18 bytes
	InboundRelayID common.ShortNodeID `insolar-transport:"AddrMode=2"`
	BasePort       uint16             `insolar-transport:"AddrMode=0,1"`
	PrimaryIPv4    uint32             `insolar-transport:"AddrMode=0"`
	PrimaryIPv6    [4]uint32          `insolar-transport:"AddrMode=1"`

	// 92 bytes
	ExtraIntroHash common.Bits224
	NodePK         common.Bits512 // works as a unique node identity
}

type NodeExtraIntro struct {
	// ByteSize>=82
	IssuedAt     common.PulseNumber
	IssuedAtTime uint64

	EndpointLen    uint8
	ExtraEndpoints []uint16

	ProofLen               uint8
	NodeReferenceWithProof []common.Bits512

	DiscoveryIssuerNodeId         common.ShortNodeID
	FullIntroSignatureByDiscovery common.Bits512
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

type MembershipAnnouncement struct {
	// ByteSize = 41 | 202 | 214

	LeaverAnnouncement *LeaverAnnouncement //first byte = 0 // ByteSize = 9
	JoinerAnnouncement *JoinerAnnouncement //first byte != 0 // ByteSize = 166 | 170 | 182

	NodeShortSignature common.Bits256 // ByteSize = 32
}

type LeaverAnnouncement struct {
	// ByteSize = 9
	Reserved0   uint8 // == 0
	LeaveReason uint32
	PulseNumber common.PulseNumber //MUST be equal to the current pulse number
}

type JoinerAnnouncement struct {
	// ByteSize = 166 | 170 | 182
	JoinerIntro          NodeBriefIntro // ByteSize=102 | 106 | 118
	JoinerIntroSignature common.Bits512 // ByteSize=64
}

type NodeNeighbourClaim struct {
	// ByteSize=133, 174, 335, 347
	ClaimHeader `insolar-transport:"ClaimType=1"`

	NodeRank         SerializedGlobulaNodeRank `insolar-transport:"[30-31]=flags:reserved"` // ByteSize=4
	GlobulaNodeState CompactGlobulaNodeState   // ByteSize=128

	MembershipAnnouncement *MembershipAnnouncement //presence is detected by claim length // ByteSize = 41 | 202 | 214
}
