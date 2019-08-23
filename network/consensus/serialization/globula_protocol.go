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
	"fmt"
	"io"

	"github.com/pkg/errors"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/longbits"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/member"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/phases"
	"github.com/insolar/insolar/network/consensus/serialization/pulseserialization"
)

type GlobulaConsensusPacketBody struct {
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
	CurrentRank  member.Rank            `insolar-transport:"Packet=0"`                           // ByteSize=4
	PulsarPacket EmbeddedPulsarData     `insolar-transport:"Packet=0,1;optional=PacketFlags[0]"` // ByteSize>=124
	Announcement MembershipAnnouncement `insolar-transport:"Packet=1,2"`                         // ByteSize= (JOINER) 5, (MEMBER) 201, 205 (MEMBER+JOINER) 196, 198, 208

	// This field can be included by sender who has introduced a joiner to facilitate joining process, and contains full intro data of the joiner
	// This field  is not mandatory and can be omitted, e.g. when network is stable or some space is required for claims
	JoinerExt NodeExtendedIntro `insolar-transport:"Packet=1;optional=PacketFlags[3]"`

	/*
		FullSelfIntro MUST be included when any of the following are true
			1. sender or receiver is a joiner
			2. sender or receiver is suspect and the other node was joined after this node became suspect
	*/
	BriefSelfIntro NodeBriefIntro   `insolar-transport:"Packet=  2;optional=PacketFlags[1:2]=1"`   // ByteSize= 135, 137, 147
	FullSelfIntro  NodeFullIntro    `insolar-transport:"Packet=1,2;optional=PacketFlags[1:2]=2,3"` // ByteSize>= 221, 223, 233
	CloudIntro     CloudIntro       `insolar-transport:"Packet=1,2;optional=PacketFlags[1:2]=2,3"` // ByteSize= 128
	JoinerSecret   longbits.Bits512 `insolar-transport:"Packet=1,2;optional=PacketFlags[1:2]=3"`   // ByteSize= 64

	Neighbourhood Neighbourhood `insolar-transport:"Packet=2"` // ByteSize= 1 + N * (205 .. 220)
	Vectors       NodeVectors   `insolar-transport:"Packet=3"` // ByteSize=133..599

	Claims ClaimList `insolar-transport:"Packet=1,3"` // ByteSize= 1 + ...
}

func (b *GlobulaConsensusPacketBody) String(ctx PacketContext) string {
	flags := ctx.GetFlagRangeInt(1, 2)
	hasBrief := flags == 1
	hasFull := flags == 2 || flags == 3
	intro := "no"
	if hasBrief {
		intro = "brief"
	}
	if hasFull {
		intro = "full"
	}

	switch ctx.GetPacketType().GetPayloadEquivalent() {
	case phases.PacketPhase0:
		return fmt.Sprintf("<current_rank=%s>", b.CurrentRank)
	case phases.PacketPhase1:
		return fmt.Sprintf("<membership=%s intro=%s>", b.Announcement, intro)
	case phases.PacketPhase2:
		return fmt.Sprintf("<membership=%s intro=%s neighbourhood=%s>", b.Announcement, intro, b.Neighbourhood)
	case phases.PacketPhase3:
		return fmt.Sprintf("<vectors=%s>", b.Vectors)
	default:
		return "unknown packet"
	}
}

func (b *GlobulaConsensusPacketBody) SerializeTo(ctx SerializeContext, writer io.Writer) error {
	packetType := ctx.GetPacketType().GetPayloadEquivalent()

	if packetType == phases.PacketPhase0 {
		if err := write(writer, b.CurrentRank); err != nil {
			return errors.Wrap(err, "failed to serialize CurrentRank")
		}
	}

	if packetType == phases.PacketPhase0 || packetType == phases.PacketPhase1 {
		if ctx.HasFlag(FlagHasPulsePacket) { // valid for Phase 0, 1: HasPulsarData : full pulsar data data is present
			if err := b.PulsarPacket.SerializeTo(ctx, writer); err != nil {
				return errors.Wrap(err, "failed to serialize PulsarPacket")
			}
		}
	}

	if packetType == phases.PacketPhase1 || packetType == phases.PacketPhase2 {
		if err := b.Announcement.SerializeTo(ctx, writer); err != nil {
			return errors.Wrap(err, "failed to serialize Announcement")
		}
	}

	if packetType == phases.PacketPhase1 {
		if ctx.HasFlag(FlagHasJoinerExt) {
			if err := b.JoinerExt.SerializeTo(ctx, writer); err != nil {
				return errors.Wrap(err, "failed to serialize JoinerExt")
			}
		}
	}

	if packetType == phases.PacketPhase2 {
		if ctx.GetFlagRangeInt(1, 2) == 1 { // [1:2] == 1 - has brief intro (this option is only allowed Phase 2 only)
			if err := b.BriefSelfIntro.SerializeTo(ctx, writer); err != nil {
				return errors.Wrap(err, "failed to serialize BriefSelfIntro")
			}
		}
	}

	if packetType == phases.PacketPhase1 || packetType == phases.PacketPhase2 {
		if ctx.GetFlagRangeInt(1, 2) == 2 || ctx.GetFlagRangeInt(1, 2) == 3 { // [1:2] in (2, 3) - has full intro + cloud intro
			if err := b.FullSelfIntro.SerializeTo(ctx, writer); err != nil {
				return errors.Wrap(err, "failed to serialize FullSelfIntro")
			}

			if err := b.CloudIntro.SerializeTo(ctx, writer); err != nil {
				return errors.Wrap(err, "failed to serialize CloudIntro")
			}
		}

		if ctx.GetFlagRangeInt(1, 2) == 3 { // [1:2] == 3 - has joiner secret (only for member-to-joiner packet)
			if err := write(writer, b.JoinerSecret); err != nil {
				return errors.Wrap(err, "failed to serialize JoinerSecret")
			}
		}
	}

	if packetType == phases.PacketPhase2 {
		if err := b.Neighbourhood.SerializeTo(ctx, writer); err != nil {
			return errors.Wrap(err, "failed to serialize Neighbourhood")
		}
	}

	if packetType == phases.PacketPhase3 {
		if err := b.Vectors.SerializeTo(ctx, writer); err != nil {
			return errors.Wrap(err, "failed to serialize Vectors")
		}
	}

	if packetType == phases.PacketPhase1 || packetType == phases.PacketPhase3 {
		if err := b.Claims.SerializeTo(ctx, writer); err != nil {
			return errors.Wrap(err, "failed to serialize Claims")
		}
	}

	return nil
}

func (b *GlobulaConsensusPacketBody) DeserializeFrom(ctx DeserializeContext, reader io.Reader) error {
	packetType := ctx.GetPacketType().GetPayloadEquivalent()

	if packetType == phases.PacketPhase0 {
		if err := read(reader, &b.CurrentRank); err != nil {
			return errors.Wrap(err, "failed to deserialize CurrentRank")
		}
	}

	if packetType == phases.PacketPhase0 || packetType == phases.PacketPhase1 {
		if ctx.HasFlag(FlagHasPulsePacket) { // valid for Phase 0, 1: HasPulsarData : full pulsar data data is present
			if err := b.PulsarPacket.DeserializeFrom(ctx, reader); err != nil {
				return errors.Wrap(err, "failed to deserialize PulsarPacket")
			}
		}
	}

	if packetType == phases.PacketPhase1 || packetType == phases.PacketPhase2 {
		if err := b.Announcement.DeserializeFrom(ctx, reader); err != nil {
			return errors.Wrap(err, "failed to deserialize Announcement")
		}
	}

	if packetType == phases.PacketPhase1 {
		if ctx.HasFlag(FlagHasJoinerExt) {
			if err := b.JoinerExt.DeserializeFrom(ctx, reader); err != nil {
				return errors.Wrap(err, "failed to deserialize JoinerExt")
			}
		}
	}

	if packetType == phases.PacketPhase2 {
		if ctx.GetFlagRangeInt(1, 2) == 1 { // [1:2] == 1 - has brief intro (this option is only allowed Phase 2 only)
			if err := b.BriefSelfIntro.DeserializeFrom(ctx, reader); err != nil {
				return errors.Wrap(err, "failed to deserialize BriefSelfIntro")
			}
		}
	}

	if packetType == phases.PacketPhase1 || packetType == phases.PacketPhase2 {
		if ctx.GetFlagRangeInt(1, 2) == 2 || ctx.GetFlagRangeInt(1, 2) == 3 { // [1:2] in (2, 3) - has full intro + cloud intro
			if err := b.FullSelfIntro.DeserializeFrom(ctx, reader); err != nil {
				return errors.Wrap(err, "failed to deserialize FullSelfIntro")
			}

			if err := b.CloudIntro.DeserializeFrom(ctx, reader); err != nil {
				return errors.Wrap(err, "failed to deserialize CloudIntro")
			}

		}

		if ctx.GetFlagRangeInt(1, 2) == 3 { // [1:2] == 3 - has joiner secret (only for member-to-joiner packet)
			if err := read(reader, &b.JoinerSecret); err != nil {
				return errors.Wrap(err, "failed to deserialize JoinerSecret")
			}
		}
	}

	if packetType == phases.PacketPhase2 {
		if err := b.Neighbourhood.DeserializeFrom(ctx, reader); err != nil {
			return errors.Wrap(err, "failed to deserialize Neighbourhood")
		}
	}

	if packetType == phases.PacketPhase3 {
		if err := b.Vectors.DeserializeFrom(ctx, reader); err != nil {
			return errors.Wrap(err, "failed to deserialize Vectors")
		}
	}

	if packetType == phases.PacketPhase1 || packetType == phases.PacketPhase3 {
		if err := b.Claims.DeserializeFrom(ctx, reader); err != nil {
			return errors.Wrap(err, "failed to deserialize Claims")
		}
	}

	return nil
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
	Phase2: <  1 200 000  1 200 000    //neighbourhood = 5
	Phase3: <    600 000 	600 000

	Total:		~3MB	   ~3MB
*/

// TODO: HACK!
type EmbeddedPulsarData struct {
	Size uint16
	Data []byte

	// ByteSize>=124
	// Header      Header       `insolar-transport:"ignore=send"`           // ByteSize=16
	// PulseNumber pulse.Number `insolar-transport:"[30-31]=0;ignore=send"` // [30-31] MUST ==0, ByteSize=4

	PulsarPacketBody `insolar-transport:"ignore=send"` // ByteSize>=108
	// PulsarSignature  longbits.Bits512                  `insolar-transport:"ignore=send"` // ByteSize=64
}

func (pd *EmbeddedPulsarData) setData(data []byte) {
	pd.Size = uint16(len(data))
	pd.Data = data
}

func (pd *EmbeddedPulsarData) SerializeTo(ctx SerializeContext, writer io.Writer) error {
	if err := write(writer, pd.Size); err != nil {
		return errors.Wrap(err, "failed to serialize Data")
	}

	if err := write(writer, pd.Data); err != nil {
		return errors.Wrap(err, "failed to serialize Data")
	}

	return nil
}

func (pd *EmbeddedPulsarData) DeserializeFrom(ctx DeserializeContext, reader io.Reader) error {
	if err := read(reader, &pd.Size); err != nil {
		return errors.Wrap(err, "failed to deserialize Size")
	}

	if pd.Size == 0 {
		return errors.New("failed to deserialize PulseDataExt")
	}

	pd.Data = make([]byte, pd.Size)
	if err := read(reader, &pd.Data); err != nil {
		return errors.Wrap(err, "failed to deserialize Data")
	}

	p, err := pulseserialization.Deserialize(pd.Data)
	if err != nil {
		return errors.Wrap(err, "failed to deserialize PulsarPacket")
	}

	pd.PulsarPacketBody.PulseNumber = p.PulseNumber
	pd.PulsarPacketBody.PulseDataExt = p.DataExt

	return nil
}

type CloudIntro struct {
	// ByteSize=128

	CloudIdentity      longbits.Bits512 // ByteSize=64
	LastCloudStateHash longbits.Bits512
}

func (ci *CloudIntro) SerializeTo(ctx SerializeContext, writer io.Writer) error {
	return write(writer, ci)
}

func (ci *CloudIntro) DeserializeFrom(ctx DeserializeContext, reader io.Reader) error {
	return read(reader, ci)
}

type Neighbourhood struct {
	// ByteSize= 1 + N * (205 .. 220)
	NeighbourCount uint8
	FraudFlags     []uint8
	Neighbours     []NeighbourAnnouncement
}

func (n Neighbourhood) String() string {
	return fmt.Sprint(n.Neighbours)
}

func (n *Neighbourhood) SerializeTo(ctx SerializeContext, writer io.Writer) error {
	if err := write(writer, n.NeighbourCount); err != nil {
		return errors.Wrap(err, "failed to serialize NeighbourCount")
	}

	for i := 0; i < int(n.NeighbourCount); i++ {
		if err := n.Neighbours[i].SerializeTo(ctx, writer); err != nil {
			return errors.Wrapf(err, "failed to serialize Neighbours[%d]", i)
		}
	}

	return nil
}

func (n *Neighbourhood) DeserializeFrom(ctx DeserializeContext, reader io.Reader) error {
	if err := read(reader, &n.NeighbourCount); err != nil {
		return errors.Wrap(err, "failed to deserialize NeighbourCount")
	}

	if n.NeighbourCount > 0 {
		n.Neighbours = make([]NeighbourAnnouncement, n.NeighbourCount)
		for i := 0; i < int(n.NeighbourCount); i++ {
			if err := n.Neighbours[i].DeserializeFrom(ctx, reader); err != nil {
				return errors.Wrapf(err, "failed to serialize Neighbours[%d]", i)
			}
		}
	}

	return nil
}

type NeighbourAnnouncement struct {
	// ByteSize(JOINER) = 73 + (135, 137, 147) = 208, 210, 220
	// ByteSize(MEMBER) = 73 + (132, 136) = 205, 209
	NeighbourNodeID insolar.ShortNodeID // ByteSize=4 // !=0

	CurrentRank    member.Rank  // ByteSize=4
	RequestedPower member.Power // ByteSize=1

	/*
		As joiner has no state before joining, its announcement and relevant signature are considered equal to
		NodeBriefIntro and related signature, and CurrentRank of joiner will always be ZERO, as joiner has no index/nodeCount/power.

		Fields "Joiner" and "JoinerIntroducedBy" MUST BE OMITTED when this joiner is introduced by the sending node
	*/
	// TODO merge "Joiner" and "JoinerIntroducedBy" fields into NeighbourJoinerAnnouncement
	Joiner             JoinAnnouncement    `insolar-transport:"optional=CurrentRank==0"` // ByteSize = 135, 137, 147
	JoinerIntroducedBy insolar.ShortNodeID `insolar-transport:"optional=CurrentRank==0"`

	/* For non-joiner */
	Member NodeAnnouncement `insolar-transport:"optional=CurrentRank!=0"` // ByteSize = 132, 136

	/* AnnounceSignature is copied from the original Phase1 */
	AnnounceSignature longbits.Bits512 // ByteSize = 64
}

func (na NeighbourAnnouncement) String() string {
	if !na.Member.AnnounceID.IsAbsent() {
		return fmt.Sprintf(
			"<node_id=%d current_rank=%s power=%d announce=%s §announce=%s>",
			na.NeighbourNodeID,
			na.CurrentRank,
			na.RequestedPower,
			na.Member,
			na.AnnounceSignature,
		)
	}

	return fmt.Sprintf("<node_id=%d current_rank=%s power=%d>", na.NeighbourNodeID, na.CurrentRank, na.RequestedPower)
}

func (na *NeighbourAnnouncement) SerializeTo(ctx SerializeContext, writer io.Writer) error {
	if err := write(writer, na.NeighbourNodeID); err != nil {
		return errors.Wrap(err, "failed to serialize NeighbourNodeID")
	}

	if err := write(writer, na.CurrentRank); err != nil {
		return errors.Wrap(err, "failed to serialize CurrentRank")
	}

	if err := write(writer, na.RequestedPower); err != nil {
		return errors.Wrap(err, "failed to serialize RequestedPower")
	}

	if na.CurrentRank == 0 {
		announcedJoinerNodeID := ctx.GetAnnouncedJoinerNodeID()
		if announcedJoinerNodeID.IsAbsent() || na.NeighbourNodeID != announcedJoinerNodeID {
			if err := na.Joiner.SerializeTo(ctx, writer); err != nil {
				return errors.Wrap(err, "failed to serialize Joiner")
			}
			if err := write(writer, na.JoinerIntroducedBy); err != nil {
				return errors.Wrap(err, "failed to serialize JoinerIntroducedBy")
			}
		}
	} else {
		ctx.SetInContext(ContextNeighbourAnnouncement)
		ctx.SetNeighbourNodeID(na.NeighbourNodeID)
		defer ctx.SetInContext(NoContext)
		defer ctx.SetNeighbourNodeID(insolar.AbsentShortNodeID)

		if err := na.Member.SerializeTo(ctx, writer); err != nil {
			return errors.Wrap(err, "failed to serialize Member")
		}
	}

	if err := write(writer, na.AnnounceSignature); err != nil {
		return errors.Wrap(err, "failed to serialize AnnounceSignature")
	}

	return nil
}

func (na *NeighbourAnnouncement) DeserializeFrom(ctx DeserializeContext, reader io.Reader) error {
	if err := read(reader, &na.NeighbourNodeID); err != nil {
		return errors.Wrap(err, "failed to deserialize NeighbourNodeID")
	}

	if err := read(reader, &na.CurrentRank); err != nil {
		return errors.Wrap(err, "failed to deserialize CurrentRank")
	}

	if err := read(reader, &na.RequestedPower); err != nil {
		return errors.Wrap(err, "failed to deserialize RequestedPower")
	}

	if na.CurrentRank == 0 {
		announcedJoinerNodeID := ctx.GetAnnouncedJoinerNodeID()
		if announcedJoinerNodeID.IsAbsent() || na.NeighbourNodeID != announcedJoinerNodeID {
			if err := na.Joiner.DeserializeFrom(ctx, reader); err != nil {
				return errors.Wrap(err, "failed to deserialize Joiner")
			}
			if err := read(reader, &na.JoinerIntroducedBy); err != nil {
				return errors.Wrap(err, "failed to deserialize JoinerIntroducedBy")
			}
		}
	} else {
		ctx.SetInContext(ContextNeighbourAnnouncement)
		ctx.SetNeighbourNodeID(na.NeighbourNodeID)
		defer ctx.SetInContext(NoContext)
		defer ctx.SetNeighbourNodeID(insolar.AbsentShortNodeID)

		if err := na.Member.DeserializeFrom(ctx, reader); err != nil {
			return errors.Wrap(err, "failed to deserialize Member")
		}
	}

	if err := read(reader, &na.AnnounceSignature); err != nil {
		return errors.Wrap(err, "failed to deserialize AnnounceSignature")
	}

	return nil
}

type MembershipAnnouncement struct {
	// ByteSize(MEMBER) = 69 + (132, 136) = 201, 205
	// ByteSize(MEMBER + JOINER) = 69 + (167, 169, 181) = 196, 198, 208
	// ByteSize(JOINER) = 4

	/*
		This field MUST be excluded from the packet, but considered for signature calculation.
		Value of this field equals SourceID
	*/
	ShortID insolar.ShortNodeID `insolar-transport:"ignore=send"` // ByteSize = 0

	CurrentRank    member.Rank  // ByteSize=4
	RequestedPower member.Power // ByteSize=1

	// NodeState CompactGlobulaNodeState `insolar-transport:"optional=CurrentRank==0" ` // ByteSize=128 TODO: serialize, fill

	/* For non-joiner ONLY */
	Member            NodeAnnouncement `insolar-transport:"optional=CurrentRank!=0"` // ByteSize = 132, 136, 267, 269, 279
	AnnounceSignature longbits.Bits512 `insolar-transport:"optional=CurrentRank!=0"` // ByteSize = 64
	// AnnounceSignature = sign(LastCloudHash + hash(NodeFullIntro) + CurrentRank + fields of MembershipAnnouncement, SK(sender))
}

func (ma MembershipAnnouncement) String() string {
	if !ma.Member.AnnounceID.IsAbsent() {
		return fmt.Sprintf(
			"<current_rank=%s power=%d announce=%s §announce=%s>",
			ma.CurrentRank,
			ma.RequestedPower,
			ma.Member,
			ma.AnnounceSignature,
		)
	}

	return fmt.Sprintf("<current_rank=%s power=%d>", ma.CurrentRank, ma.RequestedPower)
}

func (ma *MembershipAnnouncement) SerializeTo(ctx SerializeContext, writer io.Writer) error {
	if err := write(writer, ma.CurrentRank); err != nil {
		return errors.Wrap(err, "failed to serialize CurrentRank")
	}

	if err := write(writer, ma.RequestedPower); err != nil {
		return errors.Wrap(err, "failed to serialize RequestedPower")
	}

	if ma.CurrentRank != 0 {
		ctx.SetInContext(ContextMembershipAnnouncement)
		defer ctx.SetInContext(NoContext)

		if err := ma.Member.SerializeTo(ctx, writer); err != nil {
			return errors.Wrap(err, "failed to serialize Member")
		}

		if err := write(writer, ma.AnnounceSignature); err != nil {
			return errors.Wrap(err, "failed to serialize AnnounceSignature")
		}
	}

	return nil
}

func (ma *MembershipAnnouncement) DeserializeFrom(ctx DeserializeContext, reader io.Reader) error {
	if err := read(reader, &ma.CurrentRank); err != nil {
		return errors.Wrap(err, "failed to deserialize CurrentRank")
	}

	if err := read(reader, &ma.RequestedPower); err != nil {
		return errors.Wrap(err, "failed to deserialize RequestedPower")
	}

	if ma.CurrentRank != 0 {
		ctx.SetInContext(ContextMembershipAnnouncement)
		defer ctx.SetInContext(NoContext)

		if err := ma.Member.DeserializeFrom(ctx, reader); err != nil {
			return errors.Wrap(err, "failed to deserialize Member")
		}

		if err := read(reader, &ma.AnnounceSignature); err != nil {
			return errors.Wrap(err, "failed to deserialize AnnounceSignature")
		}
	}

	return nil
}

type CompactGlobulaNodeState struct {
	// ByteSize=128
	// PulseDataHash            common.Bits256 //available externally
	// FoldedLastCloudStateHash common.Bits224 //available externally
	// NodeRank                 Rank //available externally

	NodeStateHash          longbits.Bits512 // ByteSize=64
	NodeStateHashSignature longbits.Bits512 // ByteSize=64, :=Sign(NodePK, Merkle512(NodeStateHash, (LastCloudStateHash.FoldTo224() << 32 | Rank)))
}

func (gns *CompactGlobulaNodeState) SerializeTo(_ SerializeContext, writer io.Writer) error {
	return write(writer, gns)
}

func (gns *CompactGlobulaNodeState) DeserializeFrom(_ DeserializeContext, reader io.Reader) error {
	return read(reader, gns)
}

type NodeAnnouncement struct {
	// ByteSize(MembershipAnnouncement) = 132, 136, 267, 269, 279
	// ByteSize(NeighbourAnnouncement) = 132, 136

	NodeState  CompactGlobulaNodeState // ByteSize=128
	AnnounceID insolar.ShortNodeID     // ByteSize=4 // =0 - no announcement, =self - is leaver, else has joiner
	/*
		1. When is in MembershipAnnouncement
			"Leaver" is present when AnnounceID = Header.SourceID (sender is leaving)
		2. When is in NeighbourAnnouncement
			"Leaver" is present when AnnounceID = NeighbourNodeID (neighbour is leaving)
	*/
	Leaver LeaveAnnouncement `insolar-transport:"optional"` // ByteSize = 4
	/*
		1. "Joiner" is NEVER present when "Leaver" is present
		2. when AnnounceID != 0 (sender/neighbour has introduced a joiner with AnnounceID)
			a. "Joiner" is present when is in MembershipAnnouncement
			b. "Joiner" is NEVER present when is in NeighbourAnnouncement
	*/
	Joiner JoinAnnouncement `insolar-transport:"optional"` // ByteSize = 135, 137, 147
}

func (na NodeAnnouncement) String() string {
	return fmt.Sprintf(
		"<announce_id=%d nsh=%s §nsh=%s>",
		na.AnnounceID,
		na.NodeState.NodeStateHash,
		na.NodeState.NodeStateHashSignature,
	)
}

func (na *NodeAnnouncement) SerializeTo(ctx SerializeContext, writer io.Writer) error {
	if err := na.NodeState.SerializeTo(ctx, writer); err != nil {
		return errors.Wrap(err, "failed to serialize NodeState")
	}

	if err := write(writer, na.AnnounceID); err != nil {
		return errors.Wrap(err, "failed to serialize AnnounceID")
	}

	if ctx.InContext(ContextMembershipAnnouncement) {
		if na.AnnounceID == ctx.GetSourceID() {
			if err := na.Leaver.SerializeTo(ctx, writer); err != nil {
				return errors.Wrap(err, "failed to serialize Leaver")
			}
		} else if !na.AnnounceID.IsAbsent() {
			if err := na.Joiner.SerializeTo(ctx, writer); err != nil {
				return errors.Wrap(err, "failed to serialize Joiner")
			}
			ctx.SetAnnouncedJoinerNodeID(na.AnnounceID)
		}
	}

	if ctx.InContext(ContextNeighbourAnnouncement) {
		if na.AnnounceID == ctx.GetNeighbourNodeID() {
			if err := na.Leaver.SerializeTo(ctx, writer); err != nil {
				return errors.Wrap(err, "failed to serialize Leaver")
			}
		}
	}

	return nil
}

func (na *NodeAnnouncement) DeserializeFrom(ctx DeserializeContext, reader io.Reader) error {
	if err := na.NodeState.DeserializeFrom(ctx, reader); err != nil {
		return errors.Wrap(err, "failed to deserialize NodeState")
	}

	if err := read(reader, &na.AnnounceID); err != nil {
		return errors.Wrap(err, "failed to deserialize AnnounceID")
	}

	if ctx.InContext(ContextMembershipAnnouncement) {
		if na.AnnounceID == ctx.GetSourceID() {
			if err := na.Leaver.DeserializeFrom(ctx, reader); err != nil {
				return errors.Wrap(err, "failed to deserialize Leaver")
			}
		} else if !na.AnnounceID.IsAbsent() {
			if err := na.Joiner.DeserializeFrom(ctx, reader); err != nil {
				return errors.Wrap(err, "failed to deserialize Joiner")
			}
			ctx.SetAnnouncedJoinerNodeID(na.AnnounceID)
		}
	}

	if ctx.InContext(ContextNeighbourAnnouncement) {
		if na.AnnounceID == ctx.GetNeighbourNodeID() {
			if err := na.Leaver.DeserializeFrom(ctx, reader); err != nil {
				return errors.Wrap(err, "failed to deserialize Leaver")
			}
		}
	}

	return nil
}

type JoinAnnouncement struct {
	// ByteSize= 135, 137, 147
	NodeBriefIntro
}

type LeaveAnnouncement struct {
	// ByteSize = 4
	LeaveReason uint32
}

func (la *LeaveAnnouncement) SerializeTo(_ SerializeContext, writer io.Writer) error {
	return write(writer, la)
}

func (la *LeaveAnnouncement) DeserializeFrom(_ DeserializeContext, reader io.Reader) error {
	return read(reader, la)
}
