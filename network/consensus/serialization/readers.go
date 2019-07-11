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
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/network/consensus/common/cryptography_containers"
	"github.com/insolar/insolar/network/consensus/common/long_bits"
	"github.com/insolar/insolar/network/consensus/common/pulse_data"
	"github.com/insolar/insolar/network/consensus/gcpv2/gcp_types"
	"io"

	"github.com/insolar/insolar/network/consensus/gcpv2/packets"
	"github.com/insolar/insolar/network/utils"
)

type PacketParser struct {
	packetData
	digester   cryptography_containers.DataDigester
	signMethod cryptography_containers.SignMethod
}

func newPacketParser(ctx context.Context, reader io.Reader, digester cryptography_containers.DataDigester, signMethod cryptography_containers.SignMethod) (*PacketParser, error) {
	capture := utils.NewCapturingReader(reader)
	parser := &PacketParser{
		packetData: packetData{
			packet: new(Packet),
		},
		digester:   digester,
		signMethod: signMethod,
	}

	_, err := parser.packet.DeserializeFrom(ctx, capture)
	if err != nil {
		return nil, err
	}

	parser.data = capture.Captured()

	return parser, nil
}

type PacketParserFactory struct {
	digester   cryptography_containers.DataDigester
	signMethod cryptography_containers.SignMethod
}

func NewPacketParserFactory(digester cryptography_containers.DataDigester, signMethod cryptography_containers.SignMethod) *PacketParserFactory {
	return &PacketParserFactory{
		digester:   digester,
		signMethod: signMethod,
	}
}

func (f *PacketParserFactory) ParsePacket(ctx context.Context, reader io.Reader) (packets.PacketParser, error) {
	return newPacketParser(ctx, reader, f.digester, f.signMethod)
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
	return &Phase0PacketReader{*r}
}

func (r *MemberPacketReader) AsPhase1Packet() packets.Phase1PacketReader {
	return &Phase1PacketReader{*r}
}

func (r *MemberPacketReader) AsPhase2Packet() packets.Phase2PacketReader {
	return &Phase2PacketReader{*r}
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

type Phase1PacketReader struct {
	MemberPacketReader
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

func (r *Phase1PacketReader) GetCloudIntroduction() packets.CloudIntroductionReader {
	return &CloudIntroductionReader{r.MemberPacketReader}
}

func (r *Phase1PacketReader) GetFullIntroduction() packets.FullIntroductionReader {
	panic("implement me")
}

func (r *Phase1PacketReader) GetAnnouncementReader() packets.MembershipAnnouncementReader {
	panic("implement me")
}

type Phase2PacketReader struct {
	MemberPacketReader
}

func (r *Phase2PacketReader) GetBriefIntroduction() packets.BriefIntroductionReader {
	panic("implement me")
}

func (r *Phase2PacketReader) GetAnnouncementReader() packets.MembershipAnnouncementReader {
	panic("implement me")
}

func (r *Phase2PacketReader) GetNeighbourhood() []packets.MembershipAnnouncementReader {
	panic("implement me")
}

type Phase3PacketReader struct {
	MemberPacketReader
}

func (r *Phase3PacketReader) GetTrustedGlobulaAnnouncementHash() gcp_types.GlobulaAnnouncementHash {
	panic("implement me")
}

func (r *Phase3PacketReader) GetTrustedGlobulaStateSignature() gcp_types.GlobulaStateSignature {
	panic("implement me")
}

func (r *Phase3PacketReader) GetDoubtedGlobulaAnnouncementHash() gcp_types.GlobulaAnnouncementHash {
	panic("implement me")
}

func (r *Phase3PacketReader) GetDoubtedGlobulaStateSignature() gcp_types.GlobulaStateSignature {
	panic("implement me")
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

func (r *CloudIntroductionReader) GetJoinerSecret() cryptography_containers.DigestHolder {
	if r.packet.Header.GetFlagRangeInt(1, 2) != 3 {
		return nil
	}

	digest := cryptography_containers.NewDigest(&r.body.JoinerSecret, r.digester.GetDigestMethod())
	return digest.AsDigestHolder()
}

func (r *CloudIntroductionReader) GetCloudIdentity() cryptography_containers.DigestHolder {
	digest := cryptography_containers.NewDigest(&r.body.CloudIntro.CloudIdentity, r.digester.GetDigestMethod())
	return digest.AsDigestHolder()
}

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
