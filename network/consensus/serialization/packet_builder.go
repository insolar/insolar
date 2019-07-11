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
	"github.com/insolar/insolar/network/consensus/common/cryptography_containers"
	"github.com/insolar/insolar/network/consensus/gcpv2/api"
	"github.com/insolar/insolar/network/consensus/gcpv2/api_2"

	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/network/consensus/gcpv2/packets"
)

type PacketBuilder struct {
	crypto      api_2.TransportCryptographyFactory
	localConfig api_2.LocalNodeConfiguration
}

func NewPacketBuilder(crypto api_2.TransportCryptographyFactory, localConfig api_2.LocalNodeConfiguration) *PacketBuilder {
	return &PacketBuilder{
		crypto:      crypto,
		localConfig: localConfig,
	}
}

func (p *PacketBuilder) GetNeighbourhoodSize() api.NeighbourhoodSizes {
	return api.NeighbourhoodSizes{
		NeighbourhoodSize:           5,
		NeighbourhoodTrustThreshold: 2,
		JoinersPerNeighbourhood:     2,
		JoinersBoost:                1,
	}
}

func (p *PacketBuilder) preparePacket(sender *packets.NodeAnnouncementProfile, packetType api.PacketType) *Packet {
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

func (p *PacketBuilder) prepareWrapper(packet *Packet) *preparedPacketWrapper {
	return &preparedPacketWrapper{
		packet:   packet,
		digester: p.crypto.GetDigestFactory().GetPacketDigester(),
		signer:   p.crypto.GetNodeSigner(p.localConfig.GetSecretKeyStore()),
	}
}

func (p *PacketBuilder) PreparePhase0Packet(sender *packets.NodeAnnouncementProfile, pulsarPacket packets.OriginalPulsarPacket,
	options api_2.PacketSendOptions) api_2.PreparedPacketSender {

	packet := p.preparePacket(sender, api.PacketPhase0)
	if (options & api_2.SendWithoutPulseData) == 0 {
		packet.Header.SetFlag(FlagHasPulsePacket)
	}

	body := packet.EncryptableBody.(*GlobulaConsensusPacketBody)
	body.CurrentRank = sender.GetNodeRank()
	body.PulsarPacket.Data = pulsarPacket.AsBytes()

	return p.prepareWrapper(packet)
}

func (p *PacketBuilder) PreparePhase1Packet(sender *packets.NodeAnnouncementProfile, pulsarPacket packets.OriginalPulsarPacket,
	options api_2.PacketSendOptions) api_2.PreparedPacketSender {

	packet := p.preparePacket(sender, api.PacketPhase1)
	if (options & api_2.SendWithoutPulseData) == 0 {
		packet.Header.SetFlag(FlagHasPulsePacket)
	}

	body := packet.EncryptableBody.(*GlobulaConsensusPacketBody)
	body.PulsarPacket.Data = pulsarPacket.AsBytes()

	// TODO: fixed linter :)
	body.FullSelfIntro.setAddrMode(body.FullSelfIntro.getAddrMode())
	body.FullSelfIntro.setPrimaryRole(body.FullSelfIntro.getPrimaryRole())

	return p.prepareWrapper(packet)
}

func (p *PacketBuilder) PreparePhase2Packet(sender *packets.NodeAnnouncementProfile,
	neighbourhood []packets.MembershipAnnouncementReader, options api_2.PacketSendOptions) api_2.PreparedPacketSender {

	packet := p.preparePacket(sender, api.PacketPhase2)

	return p.prepareWrapper(packet)
}

func (p *PacketBuilder) PreparePhase3Packet(sender *packets.NodeAnnouncementProfile,
	vectors api.HashedNodeVector, options api_2.PacketSendOptions) api_2.PreparedPacketSender {

	packet := p.preparePacket(sender, api.PacketPhase3)

	body := packet.EncryptableBody.(*GlobulaConsensusPacketBody)
	body.Vectors.StateVectorMask.SetBitset(vectors.Bitset)

	return p.prepareWrapper(packet)
}

type preparedPacketWrapper struct {
	packet   *Packet
	buf      [packetMaxSize]byte
	digester cryptography_containers.DataDigester
	signer   cryptography_containers.DigestSigner
}

func (p *preparedPacketWrapper) SendTo(ctx context.Context, target api.NodeProfile, sendOptions api_2.PacketSendOptions, sender api_2.PacketSender) {
	p.packet.Header.TargetID = uint32(target.GetShortNodeID())

	if (sendOptions & api_2.SendWithoutPulseData) != 0 {
		p.packet.Header.ClearFlag(FlagHasPulsePacket)
	}

	buf := bytes.NewBuffer(p.buf[0:0:packetMaxSize])
	_, err := p.packet.SerializeTo(ctx, buf, p.digester, p.signer)
	if err != nil {
		inslogger.FromContext(ctx).Error(err)
	}

	sender.SendPacketToTransport(ctx, target, sendOptions, p.buf[:buf.Len()])
}

func (p *preparedPacketWrapper) SendToMany(ctx context.Context, targetCount int, sender api_2.PacketSender,
	filter func(ctx context.Context, targetIndex int) (api.NodeProfile, api_2.PacketSendOptions)) {

	for i := 0; i <= targetCount; i++ {
		if np, options := filter(ctx, i); np != nil {
			p.SendTo(ctx, np, options, sender)
		}
	}
}
