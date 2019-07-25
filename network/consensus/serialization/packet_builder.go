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

	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/network/consensus/common/cryptkit"
	"github.com/insolar/insolar/network/consensus/gcpv2/api"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/phases"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/profiles"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/proofs"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/statevector"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/transport"
)

var defaultNeighbourhoodSize = transport.NeighbourhoodSizes{
	NeighbourhoodSize:           5,
	NeighbourhoodTrustThreshold: 2,
	JoinersPerNeighbourhood:     2,
	JoinersBoost:                1,
}

type fullReader struct {
	profiles.StaticProfile
	profiles.StaticProfileExtension
}

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

func (pb *PacketBuilder) GetNeighbourhoodSize() transport.NeighbourhoodSizes {
	return defaultNeighbourhoodSize
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

func (pb *PacketBuilder) preparePacketSender(packet *Packet) *PreparedPacketSender {
	return &PreparedPacketSender{
		packet:   packet,
		digester: pb.crypto.GetDigestFactory().GetPacketDigester(),
		signer:   pb.crypto.GetNodeSigner(pb.localConfig.GetSecretKeyStore()),
	}
}

func (pb *PacketBuilder) PreparePhase0Packet(
	sender *transport.NodeAnnouncementProfile,
	pulsarPacket proofs.OriginalPulsarPacket,
	options transport.PacketPrepareOptions,
) transport.PreparedPacketSender {

	packet := pb.preparePacket(sender, phases.PacketPhase0)
	fillPhase0(packet.EncryptableBody.(*GlobulaConsensusPacketBody), sender, pulsarPacket)

	if !options.HasAny(transport.PrepareWithoutPulseData) {
		packet.Header.SetFlag(FlagHasPulsePacket)
	}

	return pb.preparePacketSender(packet)
}

func (pb *PacketBuilder) PreparePhase1Packet(
	sender *transport.NodeAnnouncementProfile,
	pulsarPacket proofs.OriginalPulsarPacket,
	welcome *proofs.NodeWelcomePackage,
	options transport.PacketPrepareOptions,
) transport.PreparedPacketSender {

	packet := pb.preparePacket(sender, phases.PacketPhase1)
	fillPhase1(packet.EncryptableBody.(*GlobulaConsensusPacketBody), sender, pulsarPacket, welcome)

	if !options.HasAny(transport.PrepareWithoutPulseData) {
		packet.Header.SetFlag(FlagHasPulsePacket)
	}

	if options.HasAny(transport.AlternativePhasePacket) {
		packet.Header.setPacketType(phases.PacketReqPhase1)
	}

	if !options.HasAny(transport.OnlyBriefIntroAboutJoiner) {
		packet.Header.SetFlag(FlagHasJoinerExt)
	}

	packet.Header.SetFlag(FlagSelfIntro1)
	if welcome != nil {
		packet.Header.ClearFlag(FlagSelfIntro1)
		packet.Header.SetFlag(FlagSelfIntro2)

		if welcome.JoinerSecret != nil {
			packet.Header.SetFlag(FlagSelfIntro1)
		}
	}

	return pb.preparePacketSender(packet)
}

func (pb *PacketBuilder) PreparePhase2Packet(
	sender *transport.NodeAnnouncementProfile,
	welcome *proofs.NodeWelcomePackage,
	neighbourhood []transport.MembershipAnnouncementReader,
	options transport.PacketPrepareOptions,
) transport.PreparedPacketSender {

	packet := pb.preparePacket(sender, phases.PacketPhase2)
	fullPhase2(packet.EncryptableBody.(*GlobulaConsensusPacketBody), sender, welcome, neighbourhood)

	if options.HasAny(transport.AlternativePhasePacket) {
		packet.Header.setPacketType(phases.PacketExtPhase2)
	}

	if !options.HasAny(transport.OnlyBriefIntroAboutJoiner) {
		packet.Header.SetFlag(FlagHasJoinerExt)
	}

	packet.Header.SetFlag(FlagSelfIntro1)
	if welcome != nil {
		packet.Header.ClearFlag(FlagSelfIntro1)
		packet.Header.SetFlag(FlagSelfIntro2)

		if welcome.JoinerSecret != nil {
			packet.Header.SetFlag(FlagSelfIntro1)
		}
	}

	return pb.preparePacketSender(packet)
}

func (pb *PacketBuilder) PreparePhase3Packet(
	sender *transport.NodeAnnouncementProfile,
	vectors statevector.Vector,
	options transport.PacketPrepareOptions,
) transport.PreparedPacketSender {

	packet := pb.preparePacket(sender, phases.PacketPhase3)
	fillPhase3(packet.EncryptableBody.(*GlobulaConsensusPacketBody), vectors)

	if options.HasAny(transport.AlternativePhasePacket) {
		packet.Header.setPacketType(phases.PacketFastPhase3)
	}

	if vectors.Doubted.AnnouncementHash != nil {
		packet.Header.SetFlag(1)
	}

	return pb.preparePacketSender(packet)
}

type PreparedPacketSender struct {
	packet   *Packet
	digester cryptkit.DataDigester
	signer   cryptkit.DigestSigner
}

func (p *PreparedPacketSender) Copy() *PreparedPacketSender {
	ppsCopy := *p
	pCopy := *p.packet
	ppsCopy.packet = &pCopy
	return &ppsCopy
}

func (p *PreparedPacketSender) SendTo(
	ctx context.Context,
	target transport.TargetProfile,
	sendOptions transport.PacketSendOptions,
	sender transport.PacketSender,
) {

	p.packet.Header.TargetID = uint32(target.GetNodeID())
	p.packet.Header.ReceiverID = uint32(target.GetNodeID())

	ctx, _ = inslogger.WithFields(ctx, map[string]interface{}{
		"receiver_node_id": p.packet.Header.ReceiverID,
		"target_node_id":   p.packet.Header.TargetID,
		"packet_type":      p.packet.Header.GetPacketType().String(),
		"packet_pulse":     p.packet.getPulseNumber(),
	})

	p.beforeSend(sender, sendOptions, target)

	var buf [packetMaxSize]byte
	buffer := bytes.NewBuffer(buf[0:0:packetMaxSize])

	_, err := p.packet.SerializeTo(ctx, buffer, p.digester, p.signer)
	if err != nil {
		inslogger.FromContext(ctx).Error(err)
	}

	sender.SendPacketToTransport(ctx, target, sendOptions, buffer.Bytes())
}

func (p *PreparedPacketSender) SendToMany(
	ctx context.Context,
	targetCount int,
	sender transport.PacketSender,
	filter transport.ProfileFilter,
) {

	for i := 0; i < targetCount; i++ {
		if np, options := filter(ctx, i); np != nil {
			p.Copy().SendTo(ctx, np, options, sender)
		}
	}
}

func (p *PreparedPacketSender) beforeSend(
	sender transport.PacketSender,
	sendOptions transport.PacketSendOptions,
	target transport.TargetProfile,
) {

	packetType := p.packet.Header.GetPacketType().GetPayloadEquivalent()

	if packetType == phases.PacketPhase0 || packetType == phases.PacketPhase1 {
		if sendOptions.HasAny(transport.SendWithoutPulseData) {
			p.packet.Header.ClearFlag(FlagHasPulsePacket)
		}
	}

	if packetType == phases.PacketPhase1 || packetType == phases.PacketPhase2 {
		if !target.IsJoiner() || !sendOptions.HasAny(transport.TargetNeedsIntro) {
			p.packet.Header.ClearFlag(FlagSelfIntro1)
			p.packet.Header.ClearFlag(FlagSelfIntro2)
		}
	}
}
