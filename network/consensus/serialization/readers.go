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
	"github.com/insolar/insolar/network"
	"io"
	"time"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/network/consensus/adapters"
	"github.com/insolar/insolar/network/consensus/common/cryptkit"
	"github.com/insolar/insolar/network/consensus/common/endpoints"
	"github.com/insolar/insolar/network/consensus/common/longbits"
	"github.com/insolar/insolar/network/consensus/common/pulse"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/member"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/phases"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/proofs"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/transport"
	"github.com/insolar/insolar/network/consensusv1/packets"
)

type originalPulsarPacket struct {
	longbits.FixedReader
}

func (p *originalPulsarPacket) OriginalPulsarPacket() {}

type packetData struct {
	data   []byte
	packet *Packet
}

func (p *packetData) GetPulseNumber() pulse.Number {
	return p.packet.getPulseNumber()
}

type PacketParser struct {
	packetData
	digester     cryptkit.DataDigester
	signMethod   cryptkit.SignMethod
	keyProcessor insolar.KeyProcessor
}

func (p *PacketParser) ParsePacketBody() (transport.PacketParser, error) {
	return nil, nil
}

func newPacketParser(
	ctx context.Context,
	reader io.Reader,
	digester cryptkit.DataDigester,
	signMethod cryptkit.SignMethod,
	keyProcessor insolar.KeyProcessor,
) (*PacketParser, error) {

	capture := network.NewCapturingReader(reader)
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
	digester     cryptkit.DataDigester
	signMethod   cryptkit.SignMethod
	keyProcessor insolar.KeyProcessor
}

func NewPacketParserFactory(
	digester cryptkit.DataDigester,
	signMethod cryptkit.SignMethod,
	keyProcessor insolar.KeyProcessor,
) *PacketParserFactory {

	return &PacketParserFactory{
		digester:     digester,
		signMethod:   signMethod,
		keyProcessor: keyProcessor,
	}
}

func (f *PacketParserFactory) ParsePacket(ctx context.Context, reader io.Reader) (transport.PacketParser, error) {
	return newPacketParser(ctx, reader, f.digester, f.signMethod, f.keyProcessor)
}

func (p *PacketParser) GetPulsePacket() transport.PulsePacketReader {
	return &PulsePacketReader{
		data: p.packetData.data,
		body: p.packet.EncryptableBody.(*PulsarPacketBody),
	}
}

func (p *PacketParser) GetMemberPacket() transport.MemberPacketReader {
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

func (p *PacketParser) GetPacketType() phases.PacketType {
	return p.packet.Header.GetPacketType()
}

func (p *PacketParser) IsRelayForbidden() bool {
	return p.packet.Header.IsRelayRestricted()
}

func (p *PacketParser) GetPacketSignature() cryptkit.SignedDigest {
	payloadReader := bytes.NewReader(p.data[:len(p.data)-signatureSize])

	signature := cryptkit.NewSignature(&p.packet.PacketSignature, p.digester.GetDigestMethod().SignedBy(p.signMethod))
	digest := p.digester.GetDigestOf(payloadReader)
	return cryptkit.NewSignedDigest(digest, signature)
}

type PulsePacketReader struct {
	data []byte
	body *PulsarPacketBody
}

func (r *PulsePacketReader) GetPulseData() pulse.Data {
	return r.body.getPulseData()
}

func (r *PulsePacketReader) GetPulseDataEvidence() proofs.OriginalPulsarPacket {
	return &originalPulsarPacket{
		FixedReader: longbits.NewFixedReader(r.data),
	}
}

type MemberPacketReader struct {
	PacketParser
	body *GlobulaConsensusPacketBody
}

func (r *MemberPacketReader) AsPhase0Packet() transport.Phase0PacketReader {
	return &Phase0PacketReader{
		MemberPacketReader: *r,
		EmbeddedPulseReader: EmbeddedPulseReader{
			MemberPacketReader: *r,
		},
	}
}

func (r *MemberPacketReader) AsPhase1Packet() transport.Phase1PacketReader {
	return &Phase1PacketReader{
		MemberPacketReader: *r,
		ExtendedIntroReader: ExtendedIntroReader{
			MemberPacketReader: *r,
		},
		EmbeddedPulseReader: EmbeddedPulseReader{
			MemberPacketReader: *r,
		},
	}
}

func (r *MemberPacketReader) AsPhase2Packet() transport.Phase2PacketReader {
	return &Phase2PacketReader{
		MemberPacketReader: *r,
		ExtendedIntroReader: ExtendedIntroReader{
			MemberPacketReader: *r,
		},
	}
}

func (r *MemberPacketReader) AsPhase3Packet() transport.Phase3PacketReader {
	return &Phase3PacketReader{*r}
}

type EmbeddedPulseReader struct {
	MemberPacketReader
}

func (r *EmbeddedPulseReader) HasPulseData() bool {
	return r.packet.Header.HasFlag(FlagHasPulsePacket)
}

func (r *EmbeddedPulseReader) GetEmbeddedPulsePacket() transport.PulsePacketReader {
	if !r.HasPulseData() {
		return nil
	}

	return &PulsePacketReader{
		data: r.body.PulsarPacket.Data,
		body: &r.body.PulsarPacket.PulsarPacketBody,
	}
}

type Phase0PacketReader struct {
	MemberPacketReader
	EmbeddedPulseReader
}

func (r *Phase0PacketReader) GetNodeRank() member.Rank {
	return r.body.CurrentRank
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

func (r *ExtendedIntroReader) GetFullIntroduction() transport.FullIntroductionReader {
	if !r.HasFullIntro() {
		return nil
	}

	return &FullIntroductionReader{
		MemberPacketReader: r.MemberPacketReader,
		intro:              r.body.FullSelfIntro,
	}
}

func (r *ExtendedIntroReader) GetCloudIntroduction() transport.CloudIntroductionReader {
	if !r.HasCloudIntro() {
		return nil
	}

	return &CloudIntroductionReader{
		MemberPacketReader: r.MemberPacketReader,
	}
}

func (r *ExtendedIntroReader) GetJoinerSecret() cryptkit.SignatureHolder {
	if !r.HasJoinerSecret() {
		return nil
	}

	return cryptkit.NewSignature(
		&r.body.JoinerSecret,
		r.digester.GetDigestMethod().SignedBy(r.signMethod),
	).AsSignatureHolder()
}

type Phase1PacketReader struct {
	MemberPacketReader
	ExtendedIntroReader
	EmbeddedPulseReader
}

func (r *Phase1PacketReader) GetAnnouncementReader() transport.MembershipAnnouncementReader {
	return &MembershipAnnouncementReader{
		MemberPacketReader: r.MemberPacketReader,
	}
}

type Phase2PacketReader struct {
	MemberPacketReader
	ExtendedIntroReader
}

func (r *Phase2PacketReader) GetBriefIntroduction() transport.BriefIntroductionReader {
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

func (r *Phase2PacketReader) GetAnnouncementReader() transport.MembershipAnnouncementReader {
	return &MembershipAnnouncementReader{
		MemberPacketReader: r.MemberPacketReader,
	}
}

func (r *Phase2PacketReader) GetNeighbourhood() []transport.MembershipAnnouncementReader {
	readers := make([]transport.MembershipAnnouncementReader, r.body.Neighbourhood.NeighbourCount)
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

func (r *Phase3PacketReader) GetTrustedGlobulaAnnouncementHash() proofs.GlobulaAnnouncementHash {
	return cryptkit.NewDigest(&r.body.Vectors.MainStateVector.VectorHash, r.digester.GetDigestMethod()).AsDigestHolder()
}

func (r *Phase3PacketReader) GetTrustedExpectedRank() member.Rank {
	return r.body.Vectors.MainStateVector.ExpectedRank
}

func (r *Phase3PacketReader) GetTrustedGlobulaStateSignature() proofs.GlobulaStateSignature {
	return cryptkit.NewSignature(
		&r.body.Vectors.MainStateVector.SignedGlobulaStateHash,
		r.digester.GetDigestMethod().SignedBy(r.signMethod),
	).AsSignatureHolder()
}

func (r *Phase3PacketReader) GetDoubtedGlobulaAnnouncementHash() proofs.GlobulaAnnouncementHash {
	if !r.hasDoubtedVector() {
		return nil
	}

	return cryptkit.NewDigest(&r.body.Vectors.AdditionalStateVectors[0].VectorHash, r.digester.GetDigestMethod()).AsDigestHolder()
}

func (r *Phase3PacketReader) GetDoubtedExpectedRank() member.Rank {
	if !r.hasDoubtedVector() {
		return 0
	}

	return r.body.Vectors.AdditionalStateVectors[0].ExpectedRank
}

func (r *Phase3PacketReader) GetDoubtedGlobulaStateSignature() proofs.GlobulaStateSignature {
	if !r.hasDoubtedVector() {
		return nil
	}

	return cryptkit.NewSignature(
		&r.body.Vectors.AdditionalStateVectors[0].SignedGlobulaStateHash,
		r.digester.GetDigestMethod().SignedBy(r.signMethod),
	).AsSignatureHolder()
}

func (r *Phase3PacketReader) GetBitset() member.StateBitset {
	return r.body.Vectors.StateVectorMask.GetBitset()
}

type CloudIntroductionReader struct {
	MemberPacketReader
}

func (r *CloudIntroductionReader) GetLastCloudStateHash() cryptkit.DigestHolder {
	digest := cryptkit.NewDigest(&r.body.CloudIntro.LastCloudStateHash, r.digester.GetDigestMethod())
	return digest.AsDigestHolder()
}

func (r *CloudIntroductionReader) hasJoinerSecret() bool {
	return r.packet.Header.GetFlagRangeInt(1, 2) == 3
}

func (r *CloudIntroductionReader) GetJoinerSecret() cryptkit.DigestHolder {
	if !r.hasJoinerSecret() {
		return nil
	}

	return cryptkit.NewDigest(&r.body.JoinerSecret, r.digester.GetDigestMethod()).AsDigestHolder()
}

func (r *CloudIntroductionReader) GetCloudIdentity() cryptkit.DigestHolder {
	digest := cryptkit.NewDigest(&r.body.CloudIntro.CloudIdentity, r.digester.GetDigestMethod())
	return digest.AsDigestHolder()
}

type FullIntroductionReader struct {
	MemberPacketReader
	intro NodeFullIntro
}

func (r *FullIntroductionReader) GetStaticNodeID() insolar.ShortNodeID {
	return r.intro.ShortID
}

func (r *FullIntroductionReader) GetPrimaryRole() member.PrimaryRole {
	return r.intro.getPrimaryRole()
}

func (r *FullIntroductionReader) GetSpecialRoles() member.SpecialRole {
	return r.intro.SpecialRoles
}

func (r *FullIntroductionReader) GetStartPower() member.Power {
	return r.intro.StartPower
}

func (r *FullIntroductionReader) GetNodePublicKey() cryptkit.SignatureKeyHolder {
	return adapters.NewECDSASignatureKeyHolderFromBits(r.intro.NodePK, r.keyProcessor)
}

func (r *FullIntroductionReader) GetDefaultEndpoint() endpoints.Outbound {
	return adapters.NewOutbound(packets.NodeAddress(r.intro.Endpoint).String())
}

func (r *FullIntroductionReader) GetJoinerSignature() cryptkit.SignatureHolder {
	return cryptkit.NewSignature(
		&r.intro.JoinerSignature,
		r.digester.GetDigestMethod().SignedBy(r.signMethod),
	).AsSignatureHolder()
}

func (r *FullIntroductionReader) GetIssuedAtPulse() pulse.Number {
	return r.intro.IssuedAtPulse
}

func (r *FullIntroductionReader) GetIssuedAtTime() time.Time {
	return time.Unix(0, int64(r.intro.IssuedAtTime))
}

func (r *FullIntroductionReader) GetPowerLevels() member.PowerSet {
	return r.intro.PowerLevels
}

func (r *FullIntroductionReader) GetExtraEndpoints() []endpoints.Outbound {
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

func (r *FullIntroductionReader) GetIssuerSignature() cryptkit.SignatureHolder {
	return cryptkit.NewSignature(
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

func (r *MembershipAnnouncementReader) GetNodeRank() member.Rank {
	return r.body.Announcement.CurrentRank
}

func (r *MembershipAnnouncementReader) GetRequestedPower() member.Power {
	if !r.hasRank() {
		return 0
	}

	return r.body.Announcement.RequestedPower
}

func (r *MembershipAnnouncementReader) GetNodeStateHashEvidence() proofs.NodeStateHashEvidence {
	if !r.hasRank() {
		return nil
	}

	return &NodeStateHashReader{
		MemberPacketReader: r.MemberPacketReader,
		gns:                r.body.Announcement.Member.NodeState,
	}
}

func (r *MembershipAnnouncementReader) GetAnnouncementSignature() proofs.MemberAnnouncementSignature {
	if !r.hasRank() {
		return nil
	}

	return cryptkit.NewSignature(
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

func (r *MembershipAnnouncementReader) GetJoinerIntroducedByID() insolar.ShortNodeID {
	panic("implement me") // TODO
}

func (r *MembershipAnnouncementReader) GetJoinerAnnouncement() transport.JoinerAnnouncementReader {
	if !r.hasRank() {
		return nil
	}

	if r.body.Announcement.Member.AnnounceID == insolar.ShortNodeID(r.packet.Header.SourceID) ||
		r.body.Announcement.Member.AnnounceID.IsAbsent() {
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

func (r *JoinerAnnouncementReader) GetBriefIntroduction() transport.BriefIntroductionReader {
	return &FullIntroductionReader{
		MemberPacketReader: r.MemberPacketReader,
		intro: NodeFullIntro{
			NodeBriefIntro: r.joiner.NodeBriefIntro,
		},
	}
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

func (r *NeighbourAnnouncementReader) GetNodeRank() member.Rank {
	return r.neighbour.CurrentRank
}

func (r *NeighbourAnnouncementReader) GetRequestedPower() member.Power {
	return r.neighbour.RequestedPower
}

func (r *NeighbourAnnouncementReader) GetNodeStateHashEvidence() proofs.NodeStateHashEvidence {
	return &NodeStateHashReader{
		MemberPacketReader: r.MemberPacketReader,
		gns:                r.neighbour.Member.NodeState,
	}
}

func (r *NeighbourAnnouncementReader) GetAnnouncementSignature() proofs.MemberAnnouncementSignature {
	return cryptkit.NewSignature(
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

func (r *NeighbourAnnouncementReader) GetJoinerIntroducedByID() insolar.ShortNodeID {
	if r.hasRank() {
		return 0
	}

	return r.neighbour.JoinerIntroducedBy
}

func (r *NeighbourAnnouncementReader) GetJoinerID() insolar.ShortNodeID {
	if r.hasRank() {
		return 0
	}

	return r.neighbour.NeighbourNodeID
}

func (r *NeighbourAnnouncementReader) GetJoinerAnnouncement() transport.JoinerAnnouncementReader {
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

func (r *NodeStateHashReader) GetNodeStateHash() proofs.NodeStateHash {
	return cryptkit.NewDigest(
		&r.gns.NodeStateHash,
		r.digester.GetDigestMethod(),
	).AsDigestHolder()
}

func (r *NodeStateHashReader) GetGlobulaNodeStateSignature() cryptkit.SignatureHolder {
	return cryptkit.NewSignature(
		&r.gns.GlobulaNodeStateSignature,
		r.digester.GetDigestMethod().SignedBy(r.signMethod),
	).AsSignatureHolder()
}
