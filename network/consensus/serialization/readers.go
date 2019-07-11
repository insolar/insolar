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
	"bytes"
	"context"
	"fmt"
	"io"
	"time"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/network/consensus/adapters"
	"github.com/insolar/insolar/network/consensus/common/cryptography_containers"
	"github.com/insolar/insolar/network/consensus/common/endpoints"
	"github.com/insolar/insolar/network/consensus/common/long_bits"
	"github.com/insolar/insolar/network/consensus/common/pulse_data"
	"github.com/insolar/insolar/network/consensus/gcpv2/gcp_types"
	"github.com/insolar/insolar/network/consensus/gcpv2/packets"
	"github.com/insolar/insolar/network/utils"
)

type originalPulsarPacket struct {
	long_bits.FixedReader
}

func (p *originalPulsarPacket) OriginalPulsarPacket() {}

type packetData struct {
	data   []byte
	packet *Packet
}

func (p *packetData) GetPulseNumber() pulse_data.PulseNumber {
	return p.packet.getPulseNumber()
}

type PacketParser struct {
	packetData
	digester     cryptography_containers.DataDigester
	signMethod   cryptography_containers.SignMethod
	keyProcessor insolar.KeyProcessor
}

func newPacketParser(
	ctx context.Context,
	reader io.Reader,
	digester cryptography_containers.DataDigester,
	signMethod cryptography_containers.SignMethod,
	keyProcessor insolar.KeyProcessor,
) (*PacketParser, error) {

	capture := utils.NewCapturingReader(reader)
	parser := &PacketParser{
		packetData: packetData{
			packet: new(Packet),
		},
		digester:     digester,
		signMethod:   signMethod,
		keyProcessor: keyProcessor,
	}

	_, err := parser.packet.DeserializeFrom(ctx, capture)
	if err != nil {
		return nil, err
	}

	parser.data = capture.Captured()

	return parser, nil
}

type PacketParserFactory struct {
	digester     cryptography_containers.DataDigester
	signMethod   cryptography_containers.SignMethod
	keyProcessor insolar.KeyProcessor
}

func NewPacketParserFactory(
	digester cryptography_containers.DataDigester,
	signMethod cryptography_containers.SignMethod,
	keyProcessor insolar.KeyProcessor,
) *PacketParserFactory {

	return &PacketParserFactory{
		digester:     digester,
		signMethod:   signMethod,
		keyProcessor: keyProcessor,
	}
}

func (f *PacketParserFactory) ParsePacket(ctx context.Context, reader io.Reader) (packets.PacketParser, error) {
	return newPacketParser(ctx, reader, f.digester, f.signMethod, f.keyProcessor)
}

func (p *PacketParser) GetPulsePacket() packets.PulsePacketReader {
	return &PulsePacketReader{
		data:        p.packetData.data,
		pulseNumber: p.packet.getPulseNumber(),
		body:        p.packet.EncryptableBody.(*PulsarPacketBody),
	}
}

func (p *PacketParser) GetMemberPacket() packets.MemberPacketReader {
	return &MemberPacketReader{
		PacketParser: *p,
		body:         p.packet.EncryptableBody.(*GlobulaConsensusPacketBody),
	}
}

func (p *PacketParser) GetSourceID() insolar.ShortNodeID {
	return p.packet.Header.GetSourceID()
}

func (p *PacketParser) GetReceiverID() insolar.ShortNodeID {
	return insolar.ShortNodeID(p.packet.Header.ReceiverID)
}

func (p *PacketParser) GetTargetID() insolar.ShortNodeID {
	return insolar.ShortNodeID(p.packet.Header.TargetID)
}

func (p *PacketParser) GetPacketType() gcp_types.PacketType {
	return p.packet.Header.GetPacketType()
}

func (p *PacketParser) IsRelayForbidden() bool {
	return p.packet.Header.IsRelayRestricted()
}

func (p *PacketParser) GetPacketSignature() cryptography_containers.SignedDigest {
	signature := cryptography_containers.NewSignature(&p.packet.PacketSignature, p.digester.GetDigestMethod().SignedBy(p.signMethod))
	digest := p.digester.GetDigestOf(bytes.NewReader(p.data))
	return cryptography_containers.NewSignedDigest(digest, signature)
}

type PulsePacketReader struct {
	data        []byte
	body        *PulsarPacketBody
	pulseNumber pulse_data.PulseNumber
}

func (r *PulsePacketReader) GetPulseData() pulse_data.PulseData {
	return pulse_data.PulseData{
		PulseNumber:  r.pulseNumber,
		PulseDataExt: r.body.PulseDataExt,
	}
}

func (r *PulsePacketReader) GetPulseDataEvidence() packets.OriginalPulsarPacket {
	return &originalPulsarPacket{
		FixedReader: long_bits.NewFixedReader(r.data),
	}
}

type MemberPacketReader struct {
	PacketParser
	body *GlobulaConsensusPacketBody
}

func (r *MemberPacketReader) AsPhase0Packet() packets.Phase0PacketReader {
	return &Phase0PacketReader{
		MemberPacketReader: *r,
	}
}

func (r *MemberPacketReader) AsPhase1Packet() packets.Phase1PacketReader {
	return &Phase1PacketReader{
		MemberPacketReader: *r,
		ExtendedIntroReader: ExtendedIntroReader{
			MemberPacketReader: *r,
		},
	}
}

func (r *MemberPacketReader) AsPhase2Packet() packets.Phase2PacketReader {
	return &Phase2PacketReader{
		MemberPacketReader: *r,
		ExtendedIntroReader: ExtendedIntroReader{
			MemberPacketReader: *r,
		},
	}
}

func (r *MemberPacketReader) AsPhase3Packet() packets.Phase3PacketReader {
	return &Phase3PacketReader{*r}
}

type Phase0PacketReader struct {
	MemberPacketReader
}

func (r *Phase0PacketReader) GetNodeRank() gcp_types.MembershipRank {
	return r.body.CurrentRank
}

func (r *Phase0PacketReader) GetEmbeddedPulsePacket() packets.PulsePacketReader {
	return &PulsePacketReader{
		data:        r.body.PulsarPacket.Data,
		pulseNumber: r.GetPulseNumber(),
		body:        &r.body.PulsarPacket.PulsarPacketBody,
	}
}

type ExtendedIntroReader struct {
	MemberPacketReader
}

func (r *ExtendedIntroReader) HasFullIntro() bool {
	flags := r.packet.Header.GetFlagRangeInt(1, 2)
	return flags == 2 || flags == 3
}

func (r *ExtendedIntroReader) HasCloudIntro() bool {
	flags := r.packet.Header.GetFlagRangeInt(1, 2)
	return flags == 2 || flags == 3
}

func (r *ExtendedIntroReader) HasJoinerSecret() bool {
	return r.packet.Header.GetFlagRangeInt(1, 2) == 3
}

func (r *ExtendedIntroReader) GetFullIntroduction() packets.FullIntroductionReader {
	if !r.HasFullIntro() {
		return nil
	}

	return &FullIntroductionReader{
		MemberPacketReader: r.MemberPacketReader,
		intro:              r.body.FullSelfIntro,
	}
}

func (r *ExtendedIntroReader) GetCloudIntroduction() packets.CloudIntroductionReader {
	if !r.HasCloudIntro() {
		return nil
	}

	return &CloudIntroductionReader{
		MemberPacketReader: r.MemberPacketReader,
	}
}

func (r *ExtendedIntroReader) GetJoinerSecret() cryptography_containers.DigestHolder {
	if !r.HasJoinerSecret() {
		return nil
	}

	return cryptography_containers.NewDigest(&r.body.JoinerSecret, r.digester.GetDigestMethod()).AsDigestHolder()
}

type Phase1PacketReader struct {
	MemberPacketReader
	ExtendedIntroReader
}

func (r *Phase1PacketReader) HasPulseData() bool {
	return r.packet.Header.hasFlag(0)
}

func (r *Phase1PacketReader) GetEmbeddedPulsePacket() packets.PulsePacketReader {
	return &PulsePacketReader{
		data:        r.body.PulsarPacket.Data,
		pulseNumber: r.GetPulseNumber(),
		body:        &r.body.PulsarPacket.PulsarPacketBody,
	}
}

func (r *Phase1PacketReader) GetAnnouncementReader() packets.MembershipAnnouncementReader {
	return &MembershipAnnouncementReader{
		MemberPacketReader: r.MemberPacketReader,
	}
}

type Phase2PacketReader struct {
	MemberPacketReader
	ExtendedIntroReader
}

func (r *Phase2PacketReader) GetBriefIntroduction() packets.BriefIntroductionReader {
	flags := r.packet.Header.GetFlagRangeInt(1, 2)
	if flags != 1 {
		return nil
	}

	return &FullIntroductionReader{
		MemberPacketReader: r.MemberPacketReader,
		intro: NodeFullIntro{
			NodeBriefIntro: r.body.BriefSelfIntro,
		},
	}
}

func (r *Phase2PacketReader) GetAnnouncementReader() packets.MembershipAnnouncementReader {
	return &MembershipAnnouncementReader{
		MemberPacketReader: r.MemberPacketReader,
	}
}

func (r *Phase2PacketReader) GetNeighbourhood() []packets.MembershipAnnouncementReader {
	readers := make([]packets.MembershipAnnouncementReader, r.body.Neighbourhood.NeighbourCount)
	for i := 0; i < int(r.body.Neighbourhood.NeighbourCount); i++ {
		readers[i] = &NeighbourAnnouncementReader{
			MemberPacketReader: r.MemberPacketReader,
			neighbour:          r.body.Neighbourhood.Neighbours[i],
		}
	}

	return readers
}

type Phase3PacketReader struct {
	MemberPacketReader
}

func (r *Phase3PacketReader) hasDoubtedVector() bool {
	return r.packet.Header.GetFlagRangeInt(1, 2) > 0
}

func (r *Phase3PacketReader) GetTrustedGlobulaAnnouncementHash() gcp_types.GlobulaAnnouncementHash {
	return cryptography_containers.NewDigest(&r.body.Vectors.MainStateVector.VectorHash, r.digester.GetDigestMethod()).AsDigestHolder()
}
func (r *Phase3PacketReader) GetTrustedExpectedRank() gcp_types.MembershipRank {
	return r.body.Vectors.MainStateVector.ExpectedRank
}

func (r *Phase3PacketReader) GetTrustedGlobulaStateSignature() gcp_types.GlobulaStateSignature {
	return cryptography_containers.NewSignature(
		&r.body.Vectors.MainStateVector.SignedGlobulaStateHash,
		r.digester.GetDigestMethod().SignedBy(r.signMethod),
	).AsSignatureHolder()
}

func (r *Phase3PacketReader) GetDoubtedGlobulaAnnouncementHash() gcp_types.GlobulaAnnouncementHash {
	if !r.hasDoubtedVector() {
		return nil
	}

	return cryptography_containers.NewDigest(&r.body.Vectors.AdditionalStateVectors[0].VectorHash, r.digester.GetDigestMethod()).AsDigestHolder()
}

func (r *Phase3PacketReader) GetDoubtedExpectedRank() gcp_types.MembershipRank {
	if !r.hasDoubtedVector() {
		return 0
	}

	return r.body.Vectors.AdditionalStateVectors[0].ExpectedRank
}

func (r *Phase3PacketReader) GetDoubtedGlobulaStateSignature() gcp_types.GlobulaStateSignature {
	if !r.hasDoubtedVector() {
		return nil
	}

	return cryptography_containers.NewSignature(
		&r.body.Vectors.AdditionalStateVectors[0].SignedGlobulaStateHash,
		r.digester.GetDigestMethod().SignedBy(r.signMethod),
	).AsSignatureHolder()
}

func (r *Phase3PacketReader) GetBitset() gcp_types.NodeBitset {
	return r.body.Vectors.StateVectorMask.GetBitset()
}

type CloudIntroductionReader struct {
	MemberPacketReader
}

func (r *CloudIntroductionReader) GetLastCloudStateHash() cryptography_containers.DigestHolder {
	digest := cryptography_containers.NewDigest(&r.body.CloudIntro.LastCloudStateHash, r.digester.GetDigestMethod())
	return digest.AsDigestHolder()
}

func (r *CloudIntroductionReader) hasJoinerSecret() bool {
	return r.packet.Header.GetFlagRangeInt(1, 2) == 3
}

func (r *CloudIntroductionReader) GetJoinerSecret() cryptography_containers.DigestHolder {
	if !r.hasJoinerSecret() {
		return nil
	}

	return cryptography_containers.NewDigest(&r.body.JoinerSecret, r.digester.GetDigestMethod()).AsDigestHolder()
}

func (r *CloudIntroductionReader) GetCloudIdentity() cryptography_containers.DigestHolder {
	digest := cryptography_containers.NewDigest(&r.body.CloudIntro.CloudIdentity, r.digester.GetDigestMethod())
	return digest.AsDigestHolder()
}

type FullIntroductionReader struct {
	MemberPacketReader
	intro NodeFullIntro
}

func (r *FullIntroductionReader) GetNodeID() insolar.ShortNodeID {
	return r.intro.ShortID
}

func (r *FullIntroductionReader) GetNodePrimaryRole() gcp_types.NodePrimaryRole {
	return r.intro.getPrimaryRole()
}

func (r *FullIntroductionReader) GetNodeSpecialRoles() gcp_types.NodeSpecialRole {
	return r.intro.SpecialRoles
}

func (r *FullIntroductionReader) GetStartPower() gcp_types.MemberPower {
	return r.intro.StartPower
}

func (r *FullIntroductionReader) GetNodePK() cryptography_containers.SignatureKeyHolder {
	return adapters.NewECDSASignatureKeyHolderFromBits(r.intro.NodePK, r.keyProcessor)
}

func (r *FullIntroductionReader) GetNodeEndpoint() endpoints.NodeEndpoint {
	ip := int2ip(r.intro.PrimaryIPv4)

	return adapters.NewNodeEndpoint(fmt.Sprintf("%s:%d", ip.String(), r.intro.BasePort))
}

func (r *FullIntroductionReader) GetJoinerSignature() cryptography_containers.SignatureHolder {
	return cryptography_containers.NewSignature(
		&r.intro.JoinerSignature,
		r.digester.GetDigestMethod().SignedBy(r.signMethod),
	).AsSignatureHolder()
}

func (r *FullIntroductionReader) GetIssuedAtPulse() pulse_data.PulseNumber {
	return r.intro.IssuedAtPulse
}

func (r *FullIntroductionReader) GetIssuedAtTime() time.Time {
	return time.Unix(0, int64(r.intro.IssuedAtTime))
}

func (r *FullIntroductionReader) GetPowerLevels() gcp_types.MemberPowerSet {
	return r.intro.PowerLevels
}

func (r *FullIntroductionReader) GetExtraEndpoints() []endpoints.NodeEndpoint {
	// TODO:
	return nil
}

func (r *FullIntroductionReader) GetReference() insolar.Reference {
	if r.body.FullSelfIntro.ProofLen > 0 {
		ref := insolar.Reference{}
		copy(ref[:], r.intro.NodeRefProof[0].AsBytes())
		return ref
	}

	return insolar.Reference{}
}

func (r *FullIntroductionReader) GetIssuerID() insolar.ShortNodeID {
	return r.intro.DiscoveryIssuerNodeID
}

func (r *FullIntroductionReader) GetIssuerSignature() cryptography_containers.SignatureHolder {
	return cryptography_containers.NewSignature(
		&r.intro.IssuerSignature,
		r.digester.GetDigestMethod().SignedBy(r.signMethod),
	).AsSignatureHolder()
}

type MembershipAnnouncementReader struct {
	MemberPacketReader
}

func (r *MembershipAnnouncementReader) hasRank() bool {
	return r.body.Announcement.CurrentRank != 0
}

func (r *MembershipAnnouncementReader) GetNodeID() insolar.ShortNodeID {
	return r.body.Announcement.ShortID
}

func (r *MembershipAnnouncementReader) GetNodeRank() gcp_types.MembershipRank {
	return r.body.Announcement.CurrentRank
}

func (r *MembershipAnnouncementReader) GetRequestedPower() gcp_types.MemberPower {
	if !r.hasRank() {
		return 0
	}

	return r.body.Announcement.RequestedPower
}

func (r *MembershipAnnouncementReader) GetNodeStateHashEvidence() gcp_types.NodeStateHashEvidence {
	if !r.hasRank() {
		return nil
	}

	return &NodeStateHashReader{
		MemberPacketReader: r.MemberPacketReader,
		gns:                r.body.Announcement.Member.NodeState,
	}
}

func (r *MembershipAnnouncementReader) GetAnnouncementSignature() gcp_types.MemberAnnouncementSignature {
	if !r.hasRank() {
		return nil
	}

	return cryptography_containers.NewSignature(
		&r.body.Announcement.AnnounceSignature,
		r.digester.GetDigestMethod().SignedBy(r.signMethod),
	).AsSignatureHolder()
}

func (r *MembershipAnnouncementReader) IsLeaving() bool {
	if !r.hasRank() {
		return false
	}

	return r.body.Announcement.Member.AnnounceID == insolar.ShortNodeID(r.packet.Header.SourceID)
}

func (r *MembershipAnnouncementReader) GetLeaveReason() uint32 {
	if !r.hasRank() {
		return 0
	}

	if r.body.Announcement.Member.AnnounceID != insolar.ShortNodeID(r.packet.Header.SourceID) {
		return 0
	}

	return r.body.Announcement.Member.Leaver.LeaveReason
}

func (r *MembershipAnnouncementReader) GetJoinerID() insolar.ShortNodeID {
	if !r.hasRank() {
		return 0
	}

	if r.body.Announcement.Member.AnnounceID == insolar.ShortNodeID(r.packet.Header.SourceID) {
		return 0
	}

	return r.body.Announcement.Member.AnnounceID
}

func (r *MembershipAnnouncementReader) GetJoinerAnnouncement() packets.JoinerAnnouncementReader {
	if !r.hasRank() {
		return nil
	}

	if r.body.Announcement.Member.AnnounceID == insolar.ShortNodeID(r.packet.Header.SourceID) || r.body.Announcement.Member.AnnounceID == 0 {
		return nil
	}

	return &JoinerAnnouncementReader{
		MemberPacketReader: r.MemberPacketReader,
		joiner:             r.body.Announcement.Member.Joiner,
	}
}

type JoinerAnnouncementReader struct {
	MemberPacketReader
	joiner JoinAnnouncement
}

func (r *JoinerAnnouncementReader) GetBriefIntro() packets.BriefIntroductionReader {
	return &FullIntroductionReader{
		MemberPacketReader: r.MemberPacketReader,
		intro: NodeFullIntro{
			NodeBriefIntro: r.joiner.NodeBriefIntro,
		},
	}
}

func (r *JoinerAnnouncementReader) GetBriefIntroSignature() cryptography_containers.SignatureHolder {
	// TODO:
	return nil
}

type NeighbourAnnouncementReader struct {
	MemberPacketReader
	neighbour NeighbourAnnouncement
}

func (r *NeighbourAnnouncementReader) hasRank() bool {
	return r.neighbour.CurrentRank != 0
}

func (r *NeighbourAnnouncementReader) GetNodeID() insolar.ShortNodeID {
	return r.neighbour.NeighbourNodeID
}

func (r *NeighbourAnnouncementReader) GetNodeRank() gcp_types.MembershipRank {
	return r.neighbour.CurrentRank
}

func (r *NeighbourAnnouncementReader) GetRequestedPower() gcp_types.MemberPower {
	return r.neighbour.RequestedPower
}

func (r *NeighbourAnnouncementReader) GetNodeStateHashEvidence() gcp_types.NodeStateHashEvidence {
	return &NodeStateHashReader{
		MemberPacketReader: r.MemberPacketReader,
		gns:                r.neighbour.Member.NodeState,
	}
}

func (r *NeighbourAnnouncementReader) GetAnnouncementSignature() gcp_types.MemberAnnouncementSignature {
	return cryptography_containers.NewSignature(
		&r.neighbour.AnnounceSignature,
		r.digester.GetDigestMethod().SignedBy(r.signMethod),
	).AsSignatureHolder()
}

func (r *NeighbourAnnouncementReader) IsLeaving() bool {
	if !r.hasRank() {
		return false
	}

	return r.neighbour.Member.AnnounceID == r.neighbour.NeighbourNodeID
}

func (r *NeighbourAnnouncementReader) GetLeaveReason() uint32 {
	if !r.hasRank() {
		return 0
	}

	if r.neighbour.Member.AnnounceID != r.neighbour.NeighbourNodeID {
		return 0
	}

	return r.neighbour.Member.Leaver.LeaveReason
}

func (r *NeighbourAnnouncementReader) GetJoinerID() insolar.ShortNodeID {
	if r.hasRank() {
		return 0
	}

	return r.neighbour.NeighbourNodeID
}

func (r *NeighbourAnnouncementReader) GetJoinerAnnouncement() packets.JoinerAnnouncementReader {
	if !r.hasRank() {
		return nil
	}

	return &JoinerAnnouncementReader{
		MemberPacketReader: r.MemberPacketReader,
		joiner:             r.body.Announcement.Member.Joiner,
	}
}

type NodeStateHashReader struct {
	MemberPacketReader
	gns CompactGlobulaNodeState
}

func (r *NodeStateHashReader) GetNodeStateHash() gcp_types.NodeStateHash {
	return cryptography_containers.NewDigest(
		&r.gns.NodeStateHash,
		r.digester.GetDigestMethod(),
	).AsDigestHolder()
}

func (r *NodeStateHashReader) GetGlobulaNodeStateSignature() cryptography_containers.SignatureHolder {
	return cryptography_containers.NewSignature(
		&r.gns.GlobulaNodeStateSignature,
		r.digester.GetDigestMethod().SignedBy(r.signMethod),
	).AsSignatureHolder()
}
