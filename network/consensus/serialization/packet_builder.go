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

	"github.com/insolar/insolar/network/consensus/common/cryptkit"
	"github.com/insolar/insolar/network/consensus/gcpv2/api"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/phases"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/profiles"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/proofs"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/statevector"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/transport"

	"github.com/insolar/insolar/instrumentation/inslogger"
)

type PacketBuilder struct {
	crypto      transport.CryptographyFactory
	localConfig api.LocalNodeConfiguration
}

func NewPacketBuilder(crypto transport.CryptographyFactory, localConfig api.LocalNodeConfiguration) *PacketBuilder {
	return &PacketBuilder{
		crypto:      crypto,
		localConfig: localConfig,
	}
}

func (p *PacketBuilder) GetNeighbourhoodSize() transport.NeighbourhoodSizes {
	return transport.NeighbourhoodSizes{
		NeighbourhoodSize:           5,
		NeighbourhoodTrustThreshold: 2,
		JoinersPerNeighbourhood:     2,
		JoinersBoost:                1,
	}
}

func (pb *PacketBuilder) preparePacket(sender *transport.NodeAnnouncementProfile, packetType phases.PacketType) *Packet {
	packet := &Packet{
		Header: Header{
			SourceID: uint32(sender.GetNodeID()),
		},
	}

	packet.Header.setProtocolType(ProtocolTypeGlobulaConsensus)
	packet.Header.setPacketType(packetType)
	packet.Header.setIsRelayRestricted(true)
	packet.Header.setIsBodyEncrypted(false)

	packet.setPulseNumber(sender.GetPulseNumber())
	packet.EncryptableBody = ProtocolTypeGlobulaConsensus.NewBody()

	return packet
}

func (pb *PacketBuilder) prepareWrapper(packet *Packet) *PreparedPacketSender {
	return &PreparedPacketSender{
		packet:   packet,
		digester: pb.crypto.GetDigestFactory().GetPacketDigester(),
		signer:   pb.crypto.GetNodeSigner(pb.localConfig.GetSecretKeyStore()),
	}
}

func (pb *PacketBuilder) PreparePhase0Packet(sender *transport.NodeAnnouncementProfile, pulsarPacket proofs.OriginalPulsarPacket,
	options transport.PacketSendOptions) transport.PreparedPacketSender {

	packet := pb.preparePacket(sender, phases.PacketPhase0)
	if (options & transport.SendWithoutPulseData) == 0 {
		packet.Header.SetFlag(FlagHasPulsePacket)
	}

	body := packet.EncryptableBody.(*GlobulaConsensusPacketBody)
	body.CurrentRank = sender.GetNodeRank()
	body.PulsarPacket.Data = pulsarPacket.AsBytes()

	return pb.prepareWrapper(packet)
}

func (pb *PacketBuilder) PreparePhase1Packet(sender *transport.NodeAnnouncementProfile, pulsarPacket proofs.OriginalPulsarPacket,
	options transport.PacketSendOptions) transport.PreparedPacketSender {

	packet := pb.preparePacket(sender, phases.PacketPhase1)
	if (options & transport.SendWithoutPulseData) == 0 {
		packet.Header.SetFlag(FlagHasPulsePacket)
	}

	body := packet.EncryptableBody.(*GlobulaConsensusPacketBody)
	body.PulsarPacket.Data = pulsarPacket.AsBytes()

	// TODO: fixed linter :)
	body.FullSelfIntro.setAddrMode(body.FullSelfIntro.getAddrMode())
	body.FullSelfIntro.setPrimaryRole(body.FullSelfIntro.getPrimaryRole())

	return pb.prepareWrapper(packet)
}

func (pb *PacketBuilder) PreparePhase2Packet(sender *transport.NodeAnnouncementProfile,
	neighbourhood []transport.MembershipAnnouncementReader, options transport.PacketSendOptions) transport.PreparedPacketSender {

	packet := pb.preparePacket(sender, phases.PacketPhase2)

	return pb.prepareWrapper(packet)
}

func (pb *PacketBuilder) PreparePhase3Packet(sender *transport.NodeAnnouncementProfile,
	vectors statevector.Vector, options transport.PacketSendOptions) transport.PreparedPacketSender {

	packet := pb.preparePacket(sender, phases.PacketPhase3)

	body := packet.EncryptableBody.(*GlobulaConsensusPacketBody)
	body.Vectors.StateVectorMask.SetBitset(vectors.Bitset)

	copy(body.Vectors.MainStateVector.VectorHash[:], vectors.Trusted.AnnouncementHash.AsBytes())
	copy(body.Vectors.MainStateVector.SignedGlobulaStateHash[:], vectors.Trusted.StateSignature.AsBytes())
	body.Vectors.MainStateVector.ExpectedRank = vectors.Trusted.ExpectedRank

	if vectors.Doubted.AnnouncementHash != nil {
		packet.Header.SetFlag(1)
		body.Vectors.AdditionalStateVectors = make([]GlobulaStateVector, 1)
		copy(body.Vectors.AdditionalStateVectors[0].VectorHash[:], vectors.Doubted.AnnouncementHash.AsBytes())
		copy(body.Vectors.AdditionalStateVectors[0].SignedGlobulaStateHash[:], vectors.Doubted.StateSignature.AsBytes())
		body.Vectors.AdditionalStateVectors[0].ExpectedRank = vectors.Doubted.ExpectedRank
	}

	return pb.prepareWrapper(packet)
}

type PreparedPacketSender struct {
	packet   *Packet
	buf      [packetMaxSize]byte
	digester cryptkit.DataDigester
	signer   cryptkit.DigestSigner
}

func (p *PreparedPacketSender) SendTo(ctx context.Context, target profiles.ActiveNode, sendOptions transport.PacketSendOptions, sender transport.PacketSender) {
	p.packet.Header.TargetID = uint32(target.GetShortNodeID())

	if (sendOptions & transport.SendWithoutPulseData) != 0 {
		p.packet.Header.ClearFlag(FlagHasPulsePacket)
	}

	buf := bytes.NewBuffer(p.buf[0:0:packetMaxSize])
	_, err := p.packet.SerializeTo(ctx, buf, p.digester, p.signer)
	if err != nil {
		inslogger.FromContext(ctx).Error(err)
	}

	sender.SendPacketToTransport(ctx, target, sendOptions, p.buf[:buf.Len()])
}

func (p *PreparedPacketSender) SendToMany(ctx context.Context, targetCount int, sender transport.PacketSender,
	filter func(ctx context.Context, targetIndex int) (profiles.ActiveNode, transport.PacketSendOptions)) {

	for i := 0; i <= targetCount; i++ {
		if np, options := filter(ctx, i); np != nil {
			p.SendTo(ctx, np, options, sender)
		}
	}
}
