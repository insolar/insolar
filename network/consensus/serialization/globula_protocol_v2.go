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
	Header      UnifiedProtocolPacketHeader `insolar-transport:"Protocol=0x01;Packet=0-4"` // ByteSize=16
	PulseNumber common.PulseNumber          `insolar-transport:"[30-31]=0"`                // [30-31] MUST ==0, ByteSize=4

	// Phases 0-2
	// - Phase0 is not sent to joiners and suspects, and PulsarPacket field must not be sent by joiners
	PulsarPacket *EmbeddedPulsarData     `insolar-transport:"optional=PacketFlags[0];Packet=0,1"` // ByteSize>=124
	Announcement *MembershipAnnouncement `insolar-transport:"Packet=1,2"`                         // ByteSize= (JOINER) 166, 168, 180, (MEMBER) 197, 201, 359, 361, 373

	// ONLY from member to joiner
	CloudToJoiner      *CloudIntro     `insolar-transport:"optional=PacketFlags[1];Packet=1"` // ByteSize= 192
	BriefIntroToJoiner *BriefSelfIntro `insolar-transport:"optional=PacketFlags[1];Packet=2"` // ByteSize= 70, 72, 84

	/*
		FullSelfIntro MUST be included when any of the following are true
			1. sender or receiver is a joiner
			2. sender or receiver is suspect and the other node was joined after this node became suspect
	*/
	FullSelfIntro *FullSelfIntro `insolar-transport:"optional=PacketFlags[1];Packet=1"` // ByteSize> 152

	Neighbourhood Neighbourhood `insolar-transport:"Packet=2"` // ByteSize= 1 + N * (107 - 209)
	Vectors       NodeVectors   `insolar-transport:"Packet=3"` // ByteSize=133..599

	Claims          ClaimList      `insolar-transport:"Packet=1,3"` // ByteSize= 1 + ...
	PacketSignature common.Bits512 // ByteSize=64
}

/*

Phase0 packet: >=208
Phase1 packet: >=754 normal, >=926
Phase2 packet: 561 (w/j) + N * (107 - 209) ... 1900 byte => (6+self) members/joiners
				w=5 -> 1397 byte
Phase3 packet: >=218 <=684

Network traffic ~1000 nodes:
			     IN          OUT
	Phase0: <    208 000 	208 000
	Phase1: <    800 000    800 000
	Phase2: <  1 400 000  1 400 000    //neighbourhood = 5
	Phase3: <    600 000 	600 000

	Total: ~	 3MB		3MB
*/

type EmbeddedPulsarData struct {
	// ByteSize>=124
	Header UnifiedProtocolPacketHeader // ByteSize=16

	// PulseNumber common.PulseNumber //available externally
	PulsarPulsePacketExt // ByteSize>=108
}

type CloudIntro struct {
	// ByteSize=192

	CloudIdentity      common.Bits512 // ByteSize=64
	JoinerSecret       common.Bits512
	LastCloudStateHash common.Bits512
}

type BriefSelfIntro struct {
	// ByteSize= 70, 72, 84
	NodeBriefIntroExt // ByteSize= 70, 72, 84
}

type FullSelfIntro struct {
	// ByteSize= >=82 + (70, 72, 84) = >152
	NodeBriefIntroExt
	NodeFullIntroExt
}

type Neighbourhood struct {
	// ByteSize=1 + N * (170 - 209)
	NeighbourCount uint8
	Neighbours     []NeighbourAnnouncement
}

type NeighbourAnnouncement struct {
	// ByteSize= 4 + x
	// ByteSize(JOINER) = 171, 173, 185
	// ByteSize(MEMBER) = 205, 209
	NeighbourNodeID common.ShortNodeID // ByteSize=4 // !=0
	MembershipAnnouncement
}
type MembershipAnnouncement struct {
	// ByteSize= 5 + (162, 164, 176, 196, 200)
	// ByteSize(JOINER MEMBER) = 167, 169, 181
	// ByteSize(MEMBER) = 201, 205
	// ByteSize(SELF) = 201, 205, 363, 365, 377
	// ByteSize(JOINER SELF) = 167, 169, 181

	CurrentRank    common2.MembershipRank // ByteSize=4
	RequestedPower common2.MemberPower    // ByteSize=1

	/*
		As joiner has no state before joining, its announcement and relevant signature are considered equal to
		NodeBriefIntro and its signature.
		CurrentRank of joiner will always be ZERO, as joiner has no index/nodeCount/power.
		The field "Joiner" MUST BE OMITTED when	this joiner is introduced by the sending node (NeighbourNodeID == StateUpdate.AnnounceID)
	*/
	Joiner *JoinAnnouncement `insolar-transport:"optional=CurrentRank==0"` // ByteSize = 162, 164, 176

	/* For non-joiner */
	Member *MemberAnnouncement `insolar-transport:"optional=CurrentRank!=0"` // ByteSize = 197, 201
}

type MemberAnnouncement struct {
	// ByteSize = 196 + (0, 4, 162, 164, 176) = 196, 200, 358, 360, 372
	// ByteSize(SELF) = 196, 200, 358, 360, 372
	// ByteSize(MEMBER) = 196, 200

	NodeState  common2.CompactGlobulaNodeState // ByteSize=128
	AnnounceID common.ShortNodeID              // ByteSize=4 // =0 - no announcement, =self - is leaver, else has joiner
	/*

		Presence of the following fields depends on current parsing context:
		1. When MembershipAnnouncement is used (for node own StateUpdate), then
			a. "Leaver" is present when AnnounceID = thisNodeID, and it means that this node is leaving
			b. otherwise "Joiner" is present when AnnounceID != 0, and it means that this node has introduced a joiner with ID = AnnounceID

			AnnounceSignature = sign(LastCloudHash + hash(NodeFullIntro) + CurrentRank + fields of MembershipAnnouncement, SK(sender))

		2. When NeighbourAnnouncement is used for a neighbour node, then
			a. "Leaver" is present when AnnounceID = NeighbourNodeID, and it means that the neighbour node is leaving
			b. while AnnounceID != 0 means that this node has introduced a joiner with ID = AnnounceID, but field "Joiner" will NOT be present.

			AnnounceSignature is copied from the original Phase1
	*/
	Leaver            *LeaveAnnouncement `insolar-transport:"optional"` // ByteSize = 4
	Joiner            *JoinAnnouncement  `insolar-transport:"optional"` // ByteSize = 162, 164, 176
	AnnounceSignature common.Bits512     `insolar-transport:"optional"` // ByteSize = 64
}

type JoinAnnouncement struct {
	// ByteSize = 162, 164, 176
	JoinerIntro     NodeBriefIntroExt // ByteSize= 98, 100, 112 // NodeId is available outside
	JoinerSignature common.Bits512    // ByteSize=64 // = sign(NodeBriefIntro, SK(joiner))
}

type LeaveAnnouncement struct {
	// ByteSize = 4
	LeaveReason uint32
}
