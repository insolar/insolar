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

package tests

import (
	"fmt"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/network/consensus/common/cryptography_containers"
	"github.com/insolar/insolar/network/consensus/common/long_bits"
	"github.com/insolar/insolar/network/consensus/common/pulse_data"
	"github.com/insolar/insolar/network/consensus/gcpv2/gcp_types"

	"github.com/insolar/insolar/network/consensus/adapters"
	"github.com/insolar/insolar/network/consensus/gcpv2/packets"
)

type EmuPacketWrapper struct {
	parser packets.PacketParser
}

func UnwrapPacketParser(payload interface{}) packets.PacketParser {
	if v, ok := payload.(EmuPacketWrapper); ok {
		return v.parser
	}
	return nil
}

func WrapPacketParser(payload packets.PacketParser) interface{} {
	return EmuPacketWrapper{parser: payload}
}

func (v EmuPacketWrapper) String() string {
	return fmt.Sprintf("Wrap{%v}", v.parser)
}

var _ cryptography_containers.SignedEvidenceHolder = &basePacket{}

type basePacket struct {
	src         insolar.ShortNodeID
	tgt         insolar.ShortNodeID
	nodeCount   uint16
	mp          gcp_types.MembershipProfile
	isLeaving   bool
	leaveReason uint32

	sd cryptography_containers.SignedDigest
}

func (r *basePacket) GetRequestedPower() gcp_types.MemberPower {
	return r.mp.RequestedPower
}

func (r *basePacket) IsLeaving() bool {
	return r.isLeaving
}

func (r *basePacket) GetLeaveReason() uint32 {
	return r.leaveReason
}

func (r *basePacket) GetJoinerID() insolar.ShortNodeID {
	return 0
}

func (r *basePacket) GetJoinerAnnouncement() packets.JoinerAnnouncementReader {
	return nil
}

func (r *basePacket) GetNodeStateHashEvidence() gcp_types.NodeStateHashEvidence {
	return r.mp.StateEvidence
}

func (r *basePacket) GetAnnouncementSignature() gcp_types.MemberAnnouncementSignature {
	return r.mp.AnnounceSignature
}

func (r *basePacket) GetNodeID() insolar.ShortNodeID {
	return r.tgt
}

func (r *basePacket) GetNodeRank() gcp_types.MembershipRank {
	return gcp_types.NewMembershipRank(r.mp.Mode, r.mp.Power, r.mp.Index, r.nodeCount)
}

func (r *basePacket) GetAnnouncementReader() packets.MembershipAnnouncementReader {
	return r
}

func (r *basePacket) GetEvidence() cryptography_containers.SignedData {
	v := long_bits.NewBits64(0)
	d := cryptography_containers.NewDigest(&v, "stub")
	s := cryptography_containers.NewSignature(&v, "stub")
	return cryptography_containers.NewSignedData(&v, d, s)
}

func (r *basePacket) GetSourceID() insolar.ShortNodeID {
	return r.src
}

func (r *basePacket) GetReceiverID() insolar.ShortNodeID {
	return r.tgt
}

func (r *basePacket) GetTargetID() insolar.ShortNodeID {
	return r.tgt
}

func (r *basePacket) IsRelayForbidden() bool {
	return true
}

func (r *basePacket) GetPacketSignature() cryptography_containers.SignedDigest {
	return r.sd
}

func (r *basePacket) GetPulsePacket() packets.PulsePacketReader {
	return nil
}

func (r *basePacket) AsPhase0Packet() packets.Phase0PacketReader {
	return nil
}

func (r *basePacket) AsPhase1Packet() packets.Phase1PacketReader {
	return nil
}

func (r *basePacket) AsPhase2Packet() packets.Phase2PacketReader {
	return nil
}

func (r *basePacket) AsPhase3Packet() packets.Phase3PacketReader {
	return nil
}

func (r *basePacket) String() string {
	return fmt.Sprintf("s:%v, t:%v", r.src, r.tgt)
}

var _ packets.Phase0PacketReader = &EmuPhase0NetPacket{}
var _ packets.MemberPacketReader = &EmuPhase0NetPacket{}
var _ packets.PacketParser = &EmuPhase0NetPacket{}
var _ emuPackerCloner = &EmuPhase0NetPacket{}

type EmuPhase0NetPacket struct {
	basePacket
	pulsePacket packets.OriginalPulsarPacket
	pn          pulse_data.PulseNumber
}

func (r *EmuPhase0NetPacket) GetPacketType() gcp_types.PacketType {
	return gcp_types.PacketPhase0
}

func (r *EmuPhase0NetPacket) GetMemberPacket() packets.MemberPacketReader {
	return r
}

func (r *EmuPhase0NetPacket) AsPhase0Packet() packets.Phase0PacketReader {
	return r
}

func (r *EmuPhase0NetPacket) GetPulseNumber() pulse_data.PulseNumber {
	if r.pulsePacket == nil {
		return r.pn
	}
	return r.pulsePacket.(*adapters.PulsePacketReader).GetPulseData().PulseNumber
}

func (r *EmuPhase0NetPacket) GetEmbeddedPulsePacket() packets.PulsePacketReader {
	return r.pulsePacket.(*adapters.PulsePacketReader)
}

func (r *EmuPhase0NetPacket) String() string {
	return fmt.Sprintf("ph:0 %v, pulsePkt: {%v}", r.basePacket.String(), r.pulsePacket)
}

var _ packets.Phase1PacketReader = &EmuPhase1NetPacket{}
var _ packets.MemberPacketReader = &EmuPhase1NetPacket{}
var _ packets.PacketParser = &EmuPhase1NetPacket{}

type EmuPhase1NetPacket struct {
	EmuPhase0NetPacket
	isRequest bool
	// packetType uint8 // to reuse this type for Phase1 and Phase1Req
}

func (r *EmuPhase1NetPacket) GetCloudIntroduction() packets.CloudIntroductionReader {
	panic("implement me")
}

func (r *EmuPhase1NetPacket) GetFullIntroduction() packets.FullIntroductionReader {
	panic("implement me")
}

func (r *EmuPhase1NetPacket) String() string {
	prefix := ""
	if r.isRequest {
		prefix = "rq"
	}
	return fmt.Sprintf("ph:1%s %s pulsePkt:{%v} mp:{%v} nc:%d", prefix, r.basePacket.String(), r.pulsePacket, r.mp, r.nodeCount)
}

func (r *EmuPhase1NetPacket) GetPacketType() gcp_types.PacketType {
	if r.isRequest {
		return gcp_types.PacketReqPhase1
	} else {
		return gcp_types.PacketPhase1
	}
}

func (r *EmuPhase1NetPacket) AsPhase0Packet() packets.Phase0PacketReader {
	return nil
}

func (r *EmuPhase1NetPacket) AsPhase1Packet() packets.Phase1PacketReader {
	return r
}

func (r *EmuPhase1NetPacket) GetNodeStateHashEvidence() gcp_types.NodeStateHashEvidence {
	return r.mp.StateEvidence
}

func (r *EmuPhase1NetPacket) HasPulseData() bool {
	return r.pulsePacket != nil
}

func (r *EmuPhase1NetPacket) GetMemberPacket() packets.MemberPacketReader {
	return r
}

var _ packets.Phase2PacketReader = &EmuPhase2NetPacket{}
var _ packets.MemberPacketReader = &EmuPhase2NetPacket{}
var _ packets.PacketParser = &EmuPhase2NetPacket{}

type EmuPhase2NetPacket struct {
	basePacket
	pulseNumber   pulse_data.PulseNumber
	neighbourhood []packets.MembershipAnnouncementReader
}

func (r *EmuPhase2NetPacket) GetBriefIntroduction() packets.BriefIntroductionReader {
	panic("implement me")
}

func (r *EmuPhase2NetPacket) String() string {
	return fmt.Sprintf("ph:2 %s pn:%v mp:{%v} nc:%d ngbh:%v", r.basePacket.String(), r.pulseNumber, r.mp, r.nodeCount, r.neighbourhood)
}

func (r *EmuPhase2NetPacket) GetNeighbourhood() []packets.MembershipAnnouncementReader {
	return r.neighbourhood
}

func (r *EmuPhase2NetPacket) GetPacketType() gcp_types.PacketType {
	return gcp_types.PacketPhase2
}

func (r *EmuPhase2NetPacket) AsPhase2Packet() packets.Phase2PacketReader {
	return r
}

func (r *EmuPhase2NetPacket) GetPulseNumber() pulse_data.PulseNumber {
	return r.pulseNumber
}

func (r *EmuPhase2NetPacket) GetMemberPacket() packets.MemberPacketReader {
	return r
}

var _ packets.Phase3PacketReader = &EmuPhase3NetPacket{}
var _ packets.MemberPacketReader = &EmuPhase3NetPacket{}
var _ packets.PacketParser = &EmuPhase3NetPacket{}

type EmuPhase3NetPacket struct {
	basePacket
	pulseNumber pulse_data.PulseNumber
	vectors     gcp_types.HashedNodeVector
}

func (r *EmuPhase3NetPacket) GetTrustedGlobulaAnnouncementHash() gcp_types.GlobulaAnnouncementHash {
	return r.vectors.TrustedAnnouncementVector
}

func (r *EmuPhase3NetPacket) GetTrustedGlobulaStateSignature() gcp_types.GlobulaStateSignature {
	return r.vectors.TrustedGlobulaStateVectorSignature
}

func (r *EmuPhase3NetPacket) GetDoubtedGlobulaAnnouncementHash() gcp_types.GlobulaAnnouncementHash {
	return r.vectors.DoubtedAnnouncementVector
}

func (r *EmuPhase3NetPacket) GetDoubtedGlobulaStateSignature() gcp_types.GlobulaStateSignature {
	return r.vectors.DoubtedGlobulaStateVectorSignature
}

func (r *EmuPhase3NetPacket) String() string {
	return fmt.Sprintf("ph:3 %s, pn:%v set:%v gahT:%v gahD:%v", r.basePacket.String(), r.pulseNumber,
		r.vectors.Bitset, r.GetTrustedGlobulaAnnouncementHash(), r.GetDoubtedGlobulaAnnouncementHash())
}

func (r *EmuPhase3NetPacket) GetBitset() gcp_types.NodeBitset {
	return r.vectors.Bitset
}

func (r *EmuPhase3NetPacket) GetPacketType() gcp_types.PacketType {
	return gcp_types.PacketPhase3
}

func (r *EmuPhase3NetPacket) AsPhase3Packet() packets.Phase3PacketReader {
	return r
}

func (r *EmuPhase3NetPacket) GetPulseNumber() pulse_data.PulseNumber {
	return r.pulseNumber
}

func (r *EmuPhase3NetPacket) GetMemberPacket() packets.MemberPacketReader {
	return r
}
