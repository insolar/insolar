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

	EncryptableBody PacketBody
	EncryptionData  []byte

	PacketSignature common.Bits512 // ByteSize=64
}

type PacketBody struct {
	/*
		PacketFlags - flags =1 outside of the prescribed phases should cause packet read error
		[0]   - valid for Phase 0, 1: HasPulsarData : full pulsar data data is present
		[1:2]
			for Phase 1, 2: HasIntro : introduction is present
				0 - no intro
				1 - brief intro (this option is only allowed Phase 2 only)
				2 - full intro + cloud intro
				3 - full intro + cloud intro + joiner secret (only for member-to-joiner packet)
			for Phase 3: ExtraVectorCount : number of additional vectors inside NodeVectors
	*/

	// Phases 0-2
	// - Phase0 is not sent to joiners and suspects, and PulsarPacket field must not be sent by joiners
	PulsarPacket *EmbeddedPulsarData     `insolar-transport:"Packet=0,1;optional=PacketFlags[0]"` // ByteSize>=124
	Announcement *MembershipAnnouncement `insolar-transport:"Packet=1,2"`                         // ByteSize= (JOINER) 5, (MEMBER) 201, 205 (MEMBER+JOINER) 196, 198, 208

	/*
		FullSelfIntro MUST be included when any of the following are true
			1. sender or receiver is a joiner
			2. sender or receiver is suspect and the other node was joined after this node became suspect
	*/
	BriefSelfIntro *NodeBriefIntro `insolar-transport:"Packet=  2;optional=PacketFlags[1:2]=1"`   // ByteSize= 135, 137, 147
	FullSelfIntro  *NodeFullIntro  `insolar-transport:"Packet=1,2;optional=PacketFlags[1:2]=2,3"` // ByteSize>= 221, 223, 233
	CloudIntro     *CloudIntro     `insolar-transport:"Packet=1,2;optional=PacketFlags[1:2]=2,3"` // ByteSize= 128
	JoinerSecret   common.Bits512  `insolar-transport:"Packet=1,2;optional=PacketFlags[1:2]=3"`   // ByteSize= 64

	Neighbourhood Neighbourhood `insolar-transport:"Packet=2"` // ByteSize= 1 + N * (205 .. 220)
	Vectors       NodeVectors   `insolar-transport:"Packet=3"` // ByteSize=133..599

	Claims ClaimList `insolar-transport:"Packet=1,3"` // ByteSize= 1 + ...
}

/*

Phase0 packet: >=208
Phase1 packet: >=717 																(claims ~700 bytes)
Phase2 packet: 293 + N * (205 .. 220) ... 1500 byte => (6+self) members/joiners
				w=5 -> 1173 byte
Phase3 packet: >=218 <=684															(claims ~700 bytes)

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
	// ByteSize=128

	CloudIdentity      common.Bits512 // ByteSize=64
	LastCloudStateHash common.Bits512
}

type Neighbourhood struct {
	// ByteSize= 1 + N * (205 .. 220)
	NeighbourCount uint8
	Neighbours     []NeighbourAnnouncement
}

type NeighbourAnnouncement struct {
	// ByteSize(JOINER) = 73 + (135, 137, 147) = 208, 210, 220
	// ByteSize(MEMBER) = 73 + (132, 136) = 205, 209
	NeighbourNodeID common.ShortNodeID // ByteSize=4 // !=0

	CurrentRank    common2.MembershipRank // ByteSize=4
	RequestedPower common2.MemberPower    // ByteSize=1

	/*
		As joiner has no state before joining, its announcement and relevant signature are considered equal to
		NodeBriefIntro and related signature, and CurrentRank of joiner will always be ZERO, as joiner has no index/nodeCount/power.

		The field "Joiner" MUST BE OMITTED when	this joiner is introduced by the sending node
	*/
	Joiner *JoinAnnouncement `insolar-transport:"optional=CurrentRank==0"` // ByteSize = 135, 137, 147

	/* For non-joiner */
	Member *NodeAnnouncement `insolar-transport:"optional=CurrentRank!=0"` // ByteSize = 132, 136

	/* AnnounceSignature is copied from the original Phase1 */
	AnnounceSignature common.Bits512 `insolar-transport:"optional"` // ByteSize = 64
}

type MembershipAnnouncement struct {
	// ByteSize(MEMBER) = 69 + (132, 136) = 201, 205
	// ByteSize(MEMBER + JOINER) = 69 + (167, 169, 181) = 196, 198, 208
	// ByteSize(JOINER) = 4

	/*
		This field MUST be excluded from the packet, but considered for signature calculation.
		Value of this field equals SourceID
	*/
	ShortID common.ShortNodeID `insolar-transport:"ignore=send"` // ByteSize = 0

	CurrentRank common2.MembershipRank // ByteSize=4

	/* For non-joiner ONLY */
	RequestedPower    common2.MemberPower `insolar-transport:"optional=CurrentRank!=0"` // ByteSize=1
	Member            *NodeAnnouncement   `insolar-transport:"optional=CurrentRank!=0"` // ByteSize = 132, 136, 267, 269, 279
	AnnounceSignature common.Bits512      `insolar-transport:"optional=CurrentRank!=0"` // ByteSize = 64
	// AnnounceSignature = sign(LastCloudHash + hash(NodeFullIntro) + CurrentRank + fields of MembershipAnnouncement, SK(sender))
}

type NodeAnnouncement struct {
	// ByteSize(MembershipAnnouncement) = 132, 136, 267, 269, 279
	// ByteSize(NeighbourAnnouncement) = 132, 136

	NodeState  common2.CompactGlobulaNodeState // ByteSize=128
	AnnounceID common.ShortNodeID              // ByteSize=4 // =0 - no announcement, =self - is leaver, else has joiner
	/*
		1. When is in MembershipAnnouncement
			"Leaver" is present when AnnounceID = Header.SourceID (sender is leaving)
		2. When is in NeighbourAnnouncement
			"Leaver" is present when AnnounceID = NeighbourNodeID (neighbour is leaving)
	*/
	Leaver *LeaveAnnouncement `insolar-transport:"optional"` // ByteSize = 4
	/*
		1. "Joiner" is NEVER present when "Leaver" is present
		2. when AnnounceID != 0 (sender/neighbour has introduced a joiner with AnnounceID)
			a. "Joiner" is present when is in MembershipAnnouncement
			b. "Joiner" is NEVER present when is in NeighbourAnnouncement
	*/
	Joiner *JoinAnnouncement `insolar-transport:"optional"` // ByteSize = 135, 137, 147
}

type JoinAnnouncement struct {
	// ByteSize= 135, 137, 147
	NodeBriefIntro
}

type LeaveAnnouncement struct {
	// ByteSize = 4
	LeaveReason uint32
}
