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
	"io"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/longbits"
	"github.com/insolar/insolar/network/consensus/common/cryptkit"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/member"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/phases"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/profiles"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/proofs"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/statevector"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/transport"
	"github.com/insolar/insolar/pulse"
)

type EmuPacketWrapper struct {
	parser transport.PacketParser
}

func UnwrapPacketParser(payload interface{}) transport.PacketParser {
	if v, ok := payload.(EmuPacketWrapper); ok {
		return v.parser
	}
	return nil
}

func WrapPacketParser(payload transport.PacketParser) interface{} {
	return EmuPacketWrapper{parser: payload}
}

func (v EmuPacketWrapper) String() string {
	return fmt.Sprintf("Wrap{%v}", v.parser)
}

var _ transport.PulsePacketReader = &EmuPulsarNetPacket{}
var _ proofs.OriginalPulsarPacket = &EmuPulsarNetPacket{}
var _ transport.PacketParser = &EmuPulsarNetPacket{}

type EmuPulsarNetPacket struct {
	pulseData pulse.Data
}

func (r *EmuPulsarNetPacket) GetPulseDataDigest() cryptkit.DigestHolder {
	return nil
}

func (r *EmuPulsarNetPacket) ParsePacketBody() (transport.PacketParser, error) {
	return nil, nil
}

func (r *EmuPulsarNetPacket) IsRelayForbidden() bool {
	return false
}

func (r *EmuPulsarNetPacket) AsByteString() longbits.ByteString {
	panic("implement me")
}

func (r *EmuPulsarNetPacket) WriteTo(w io.Writer) (n int64, err error) {
	panic("implement me")
}

func (r *EmuPulsarNetPacket) Read(p []byte) (n int, err error) {
	panic("implement me")
}

func (r *EmuPulsarNetPacket) AsBytes() []byte {
	panic("implement me")
}

func (r *EmuPulsarNetPacket) FixedByteSize() int {
	panic("implement me")
}

func (r *EmuPulsarNetPacket) GetSourceID() insolar.ShortNodeID {
	return insolar.AbsentShortNodeID
}

func (r *EmuPulsarNetPacket) GetReceiverID() insolar.ShortNodeID {
	return insolar.AbsentShortNodeID
}

func (r *EmuPulsarNetPacket) GetTargetID() insolar.ShortNodeID {
	return insolar.AbsentShortNodeID
}

func (r *EmuPulsarNetPacket) OriginalPulsarPacket() {
}

func (r *EmuPulsarNetPacket) GetPacketSignature() cryptkit.SignedDigest {
	return cryptkit.SignedDigest{}
}

func (*EmuPulsarNetPacket) GetPacketType() phases.PacketType {
	return phases.PacketPulsarPulse
}

func (*EmuPulsarNetPacket) GetMemberPacket() transport.MemberPacketReader {
	return nil
}

func (r *EmuPulsarNetPacket) GetPulseData() pulse.Data {
	return r.pulseData
}

func (r *EmuPulsarNetPacket) GetPulseDataEvidence() proofs.OriginalPulsarPacket {
	return r
}

func (r *EmuPulsarNetPacket) GetPulseNumber() pulse.Number {
	return r.pulseData.PulseNumber
}

func (r *EmuPulsarNetPacket) GetPulsePacket() transport.PulsePacketReader {
	return r
}

func (r *EmuPulsarNetPacket) String() string {
	return fmt.Sprintf("pd:{%v}, pulsar:*", r.pulseData)
}

// var _ gcp_v2.PhasePacketReader = &basePacket{}
// var _ gcp_v2.MemberPacketReader = &basePacket{}
// var _ cryptkit.SignedEvidenceHolder = &basePacket{}

type basePacket struct {
	src           insolar.ShortNodeID
	tgt           insolar.ShortNodeID
	isAlternative bool
	nodeCount     uint16
	mp            profiles.MembershipProfile
	isLeaving     bool
	leaveReason   uint32
	joiner        transport.JoinerAnnouncementReader
	// joinerAnnouncer insolar.ShortNodeID
	cloudIntro *proofs.NodeWelcomePackage

	profiles.BriefCandidateProfile
	profiles.CandidateProfileExtension
}

func (r *basePacket) GetLastCloudStateHash() cryptkit.DigestHolder {
	return r.cloudIntro.LastCloudStateHash
}

func (r *basePacket) GetCloudIdentity() cryptkit.DigestHolder {
	return r.cloudIntro.CloudIdentity
}

func (r *basePacket) HasFullIntro() bool {
	return r.CandidateProfileExtension != nil && r.BriefCandidateProfile != nil
}

func (r *basePacket) GetCloudIntroduction() transport.CloudIntroductionReader {
	if !r.HasCloudIntro() {
		return nil
	}
	return r
}

func (r *basePacket) GetFullIntroduction() transport.FullIntroductionReader {
	if r.HasFullIntro() {
		return r
	}
	return nil
}

func (r *basePacket) HasCloudIntro() bool {
	return r.cloudIntro != nil
}

func (r *basePacket) HasJoinerSecret() bool {
	return r.cloudIntro != nil && r.cloudIntro.JoinerSecret != nil
}

func (r *basePacket) GetJoinerSecret() cryptkit.DigestHolder {
	if !r.HasJoinerSecret() {
		return nil
	}
	return r.cloudIntro.JoinerSecret
}

// func (r *basePacket) GetJoinerIntroducedByID() insolar.ShortNodeID {
//	if r.joiner == nil {
//		return insolar.AbsentShortNodeID
//	}
//	return r.joinerAnnouncer
// }

func (r *basePacket) ParsePacketBody() (transport.PacketParser, error) {
	return nil, nil
}

func (r *basePacket) GetRequestedPower() member.Power {
	return r.mp.RequestedPower
}

func (r *basePacket) IsLeaving() bool {
	return r.isLeaving
}

func (r *basePacket) GetLeaveReason() uint32 {
	return r.leaveReason
}

func (r *basePacket) GetJoinerID() insolar.ShortNodeID {
	if r.joiner == nil {
		return insolar.AbsentShortNodeID
	}
	return r.joiner.GetBriefIntroduction().GetStaticNodeID()
}

func (r *basePacket) GetJoinerAnnouncement() transport.JoinerAnnouncementReader {
	return r.joiner
}

func (r *basePacket) GetAnnouncementSignature() proofs.MemberAnnouncementSignature {
	return r.mp.AnnounceSignature
}

func (r *basePacket) GetNodeID() insolar.ShortNodeID {
	return r.tgt
}

func (r *basePacket) GetNodeRank() member.Rank {
	return r.mp.AsRankUint16(r.nodeCount)
}

func (r *basePacket) GetAnnouncementReader() transport.MembershipAnnouncementReader {
	return r
}

func (r *basePacket) GetNodeStateHashEvidence() proofs.NodeStateHashEvidence {
	return r.mp.StateEvidence
}
func (r *basePacket) GetEvidence() cryptkit.SignedData {
	v := longbits.NewBits64(0)
	d := cryptkit.NewDigest(&v, "stub")
	s := cryptkit.NewSignature(&v, "stub")
	return cryptkit.NewSignedData(&v, d, s)
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

func (r *basePacket) GetPacketSignature() cryptkit.SignedDigest {
	return cryptkit.SignedDigest{}
}

func (r *basePacket) GetPulsePacket() transport.PulsePacketReader {
	return nil
}

func (r *basePacket) AsPhase0Packet() transport.Phase0PacketReader {
	return nil
}

func (r *basePacket) AsPhase1Packet() transport.Phase1PacketReader {
	return nil
}

func (r *basePacket) AsPhase2Packet() transport.Phase2PacketReader {
	return nil
}

func (r *basePacket) AsPhase3Packet() transport.Phase3PacketReader {
	return nil
}

func (r *basePacket) String() string {
	intro := ""
	if r.HasFullIntro() {
		intro = " intro:full"
	} else if r.BriefCandidateProfile != nil {
		intro = " intro:brief"
	}
	cloud := ""
	if r.HasCloudIntro() {
		cloud = fmt.Sprintf(" cloud:%v", *r.cloudIntro)
	}
	announcement := ""
	if r.isLeaving {
		announcement = fmt.Sprintf(" leave:%d", r.leaveReason)
	} else if r.joiner != nil {
		joinerID := r.GetJoinerID()
		if r.joiner.HasFullIntro() {
			announcement = fmt.Sprintf(" join:%d+full", joinerID)
		} else {
			announcement = fmt.Sprintf(" join:%d", joinerID)
		}
	}
	return fmt.Sprintf("s:%v t:%v%s%s%s", r.src, r.tgt, announcement, intro, cloud)
}

func (r *basePacket) adjustBySender(profile *transport.NodeAnnouncementProfile) {
	if profile.GetNodeRank().IsJoiner() {
		r.mp.AnnounceSignature = nil
		r.mp.StateEvidence = nil
	}
}

var _ transport.Phase0PacketReader = &EmuPhase0NetPacket{}
var _ transport.MemberPacketReader = &EmuPhase0NetPacket{}
var _ transport.PacketParser = &EmuPhase0NetPacket{}
var _ emuPackerCloner = &EmuPhase0NetPacket{}

type EmuPhase0NetPacket struct {
	basePacket
	pulsePacket *EmuPulsarNetPacket
	pn          pulse.Number
}

func (r *EmuPhase0NetPacket) GetPacketType() phases.PacketType {
	return phases.PacketPhase0
}

func (r *EmuPhase0NetPacket) GetMemberPacket() transport.MemberPacketReader {
	return r
}

func (r *EmuPhase0NetPacket) AsPhase0Packet() transport.Phase0PacketReader {
	return r
}

func (r *EmuPhase0NetPacket) GetPulseNumber() pulse.Number {
	if r.pulsePacket == nil {
		return r.pn
	}
	return r.pulsePacket.pulseData.PulseNumber
}

func (r *EmuPhase0NetPacket) GetEmbeddedPulsePacket() transport.PulsePacketReader {
	return r.pulsePacket
}

func (r *EmuPhase0NetPacket) String() string {
	return fmt.Sprintf("ph:0 %v pulsePkt: {%v} mp:{%v} nc:%d ", r.basePacket.String(), r.pulsePacket, r.mp, r.nodeCount)
}

var _ transport.Phase1PacketReader = &EmuPhase1NetPacket{}
var _ transport.MemberPacketReader = &EmuPhase1NetPacket{}
var _ transport.PacketParser = &EmuPhase1NetPacket{}

type EmuPhase1NetPacket struct {
	EmuPhase0NetPacket
	// packetType uint8 // to reuse this type for Phase1 and Phase1Req
}

func (r *EmuPhase1NetPacket) String() string {
	suffix := ""
	if r.isAlternative {
		suffix = "rq"
	}
	return fmt.Sprintf("ph:1%s %s pulsePkt:{%v} mp:{%v} nc:%d", suffix, r.basePacket.String(), r.pulsePacket, r.mp, r.nodeCount)
}

func (r *EmuPhase1NetPacket) GetPacketType() phases.PacketType {
	if r.isAlternative {
		return phases.PacketReqPhase1
	} else {
		return phases.PacketPhase1
	}
}

func (r *EmuPhase1NetPacket) AsPhase0Packet() transport.Phase0PacketReader {
	return nil
}

func (r *EmuPhase1NetPacket) AsPhase1Packet() transport.Phase1PacketReader {
	return r
}

func (r *EmuPhase1NetPacket) GetNodeStateHashEvidence() proofs.NodeStateHashEvidence {
	return r.mp.StateEvidence
}

func (r *EmuPhase1NetPacket) HasPulseData() bool {
	return r.pulsePacket != nil
}

func (r *EmuPhase1NetPacket) GetMemberPacket() transport.MemberPacketReader {
	return r
}

var _ transport.Phase2PacketReader = &EmuPhase2NetPacket{}
var _ transport.MemberPacketReader = &EmuPhase2NetPacket{}
var _ transport.PacketParser = &EmuPhase2NetPacket{}

type EmuPhase2NetPacket struct {
	basePacket
	pulseNumber   pulse.Number
	neighbourhood []transport.MembershipAnnouncementReader
}

func (r *EmuPhase2NetPacket) GetBriefIntroduction() transport.BriefIntroductionReader {
	return r.BriefCandidateProfile
}

func (r *EmuPhase2NetPacket) String() string {
	suffix := ""
	if r.isAlternative {
		suffix = "xt"
	}
	return fmt.Sprintf("ph:2%s %s pn:%v mp:{%v} nc:%d ngbh:%v", suffix, r.basePacket.String(), r.pulseNumber, r.mp, r.nodeCount, r.neighbourhood)
}

func (r *EmuPhase2NetPacket) GetPacketType() phases.PacketType {
	if r.isAlternative {
		return phases.PacketExtPhase2
	} else {
		return phases.PacketPhase2
	}
}

func (r *EmuPhase2NetPacket) GetNeighbourhood() []transport.MembershipAnnouncementReader {
	return r.neighbourhood
}

func (r *EmuPhase2NetPacket) AsPhase2Packet() transport.Phase2PacketReader {
	return r
}

func (r *EmuPhase2NetPacket) GetPulseNumber() pulse.Number {
	return r.pulseNumber
}

func (r *EmuPhase2NetPacket) GetMemberPacket() transport.MemberPacketReader {
	return r
}

var _ transport.Phase3PacketReader = &EmuPhase3NetPacket{}
var _ transport.MemberPacketReader = &EmuPhase3NetPacket{}
var _ transport.PacketParser = &EmuPhase3NetPacket{}

type EmuPhase3NetPacket struct {
	basePacket
	pulseNumber pulse.Number
	vectors     statevector.Vector
}

func (r *EmuPhase3NetPacket) GetTrustedExpectedRank() member.Rank {
	return r.vectors.Trusted.ExpectedRank
}

func (r *EmuPhase3NetPacket) GetDoubtedExpectedRank() member.Rank {
	return r.vectors.Doubted.ExpectedRank
}

func (r *EmuPhase3NetPacket) GetTrustedGlobulaAnnouncementHash() proofs.GlobulaAnnouncementHash {
	return r.vectors.Trusted.AnnouncementHash
}

func (r *EmuPhase3NetPacket) GetTrustedGlobulaStateSignature() proofs.GlobulaStateSignature {
	return r.vectors.Trusted.StateSignature
}

func (r *EmuPhase3NetPacket) GetDoubtedGlobulaAnnouncementHash() proofs.GlobulaAnnouncementHash {
	return r.vectors.Doubted.AnnouncementHash
}

func (r *EmuPhase3NetPacket) GetDoubtedGlobulaStateSignature() proofs.GlobulaStateSignature {
	return r.vectors.Doubted.StateSignature
}

func (r *EmuPhase3NetPacket) GetPacketType() phases.PacketType {
	if r.isAlternative {
		return phases.PacketFastPhase3
	} else {
		return phases.PacketPhase3
	}
}

func (r *EmuPhase3NetPacket) String() string {
	suffix := ""
	if r.isAlternative {
		suffix = "ft"
	}
	return fmt.Sprintf("ph:3%s %s, pn:%v set:%v gahT:%v gahD:%v", suffix, r.basePacket.String(), r.pulseNumber,
		r.vectors.Bitset, r.GetTrustedGlobulaAnnouncementHash(), r.GetDoubtedGlobulaAnnouncementHash())
}

func (r *EmuPhase3NetPacket) GetBitset() member.StateBitset {
	return r.vectors.Bitset
}

func (r *EmuPhase3NetPacket) AsPhase3Packet() transport.Phase3PacketReader {
	return r
}

func (r *EmuPhase3NetPacket) GetPulseNumber() pulse.Number {
	return r.pulseNumber
}

func (r *EmuPhase3NetPacket) GetMemberPacket() transport.MemberPacketReader {
	return r
}
