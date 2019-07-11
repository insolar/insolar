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
	"bytes"
	"context"

	"github.com/insolar/insolar/network/consensus/common/cryptography_containers"
	"github.com/insolar/insolar/network/consensus/gcpv2/api"
	"github.com/insolar/insolar/network/consensus/gcpv2/gcp_types"

	"github.com/insolar/insolar/network/consensus/adapters"
	"github.com/insolar/insolar/network/consensus/gcpv2/packets"
)

type emuPackerCloner interface {
	clonePacketFor(target gcp_types.NodeProfile, sendOptions api.PacketSendOptions) packets.PacketParser
}

type emuPacketSender struct {
	cloner emuPackerCloner
}

func (r *emuPacketSender) SendToMany(ctx context.Context, targetCount int, sender api.PacketSender,
	filter func(ctx context.Context, targetIndex int) (gcp_types.NodeProfile, api.PacketSendOptions)) {
	for i := 0; i < targetCount; i++ {
		sendTo, sendOptions := filter(ctx, i)
		if sendTo != nil {
			c := r.cloner.clonePacketFor(sendTo, sendOptions)
			sender.SendPacketToTransport(ctx, sendTo, sendOptions, c)
		}

		select {
		case <-ctx.Done():
			return
		default:
		}
	}
}

func (r *emuPacketSender) SendTo(ctx context.Context, t gcp_types.NodeProfile, sendOptions api.PacketSendOptions, s api.PacketSender) {
	c := r.cloner.clonePacketFor(t, sendOptions)
	s.SendPacketToTransport(ctx, t, sendOptions, c)
}

type emuPacketBuilder struct {
	crypto      api.TransportCryptographyFactory
	localConfig api.LocalNodeConfiguration
}

func NewEmuPacketBuilder(crypto api.TransportCryptographyFactory, localConfig api.LocalNodeConfiguration) api.PacketBuilder {
	return &emuPacketBuilder{
		crypto:      crypto,
		localConfig: localConfig,
	}
}

func (r *emuPacketBuilder) GetNeighbourhoodSize() gcp_types.NeighbourhoodSizes {
	return gcp_types.NeighbourhoodSizes{NeighbourhoodSize: 5, NeighbourhoodTrustThreshold: 2, JoinersPerNeighbourhood: 2, JoinersBoost: 1}
}

func (r *emuPacketBuilder) defaultSign() cryptography_containers.SignedDigest {
	digester := r.crypto.GetDigestFactory().GetPacketDigester()
	digest := digester.GetDigestOf(bytes.NewReader([]byte{1, 3, 3, 7}))
	signer := r.crypto.GetNodeSigner(r.localConfig.GetSecretKeyStore())
	signature := signer.SignDigest(digest)
	sd := cryptography_containers.NewSignedDigest(digest, signature)
	return sd
}

func (r *emuPacketBuilder) defaultBasePacket(sender *packets.NodeAnnouncementProfile) basePacket {
	sd := r.defaultSign()

	return basePacket{
		src:       sender.GetNodeID(),
		mp:        sender.GetMembershipProfile(),
		nodeCount: sender.GetNodeCount(),
		sd:        sd,
	}
}

func (r *emuPacketBuilder) PreparePhase0Packet(sender *packets.NodeAnnouncementProfile, pulsarPacket packets.OriginalPulsarPacket,
	options api.PacketSendOptions) api.PreparedPacketSender {
	v := EmuPhase0NetPacket{
		basePacket:  r.defaultBasePacket(sender),
		pulsePacket: pulsarPacket,
	}
	return &emuPacketSender{&v}
}

func (r *EmuPhase0NetPacket) clonePacketFor(t gcp_types.NodeProfile, sendOptions api.PacketSendOptions) packets.PacketParser {
	c := *r
	c.tgt = t.GetShortNodeID()
	return &c
}

func (r *emuPacketBuilder) PreparePhase1Packet(sender *packets.NodeAnnouncementProfile, pulsarPacket packets.OriginalPulsarPacket,
	options api.PacketSendOptions) api.PreparedPacketSender {

	pp := pulsarPacket.(*adapters.PulsePacketReader)
	pulseData := pp.GetPulseData()
	if pp == nil || !pulseData.IsValidPulseData() {
		panic("pulse data is missing or invalid")
	}

	v := EmuPhase1NetPacket{
		EmuPhase0NetPacket: EmuPhase0NetPacket{
			basePacket:  r.defaultBasePacket(sender),
			pulsePacket: pp,
		},
	}
	v.pn = pulseData.PulseNumber
	v.isRequest = options&api.RequestForPhase1 != 0
	if v.isRequest || options&api.SendWithoutPulseData != 0 {
		v.pulsePacket = nil
	}

	return &emuPacketSender{&v}
}

func (r *EmuPhase1NetPacket) clonePacketFor(t gcp_types.NodeProfile, sendOptions api.PacketSendOptions) packets.PacketParser {
	c := *r
	c.tgt = t.GetShortNodeID()

	if sendOptions&api.SendWithoutPulseData != 0 {
		c.pulsePacket = nil
	}

	return &c
}

func (r *emuPacketBuilder) PreparePhase2Packet(sender *packets.NodeAnnouncementProfile,
	neighbourhood []packets.MembershipAnnouncementReader,
	options api.PacketSendOptions) api.PreparedPacketSender {

	v := EmuPhase2NetPacket{
		basePacket:    r.defaultBasePacket(sender),
		pulseNumber:   sender.GetPulseNumber(),
		neighbourhood: neighbourhood,
	}
	return &emuPacketSender{&v}
}

func (r *EmuPhase2NetPacket) clonePacketFor(t gcp_types.NodeProfile, sendOptions api.PacketSendOptions) packets.PacketParser {
	c := *r
	c.tgt = t.GetShortNodeID()
	return &c
}

func (r *emuPacketBuilder) PreparePhase3Packet(sender *packets.NodeAnnouncementProfile, vectors gcp_types.HashedNodeVector,
	// bitset nodeset.NodeBitset, gshTrusted common2.GlobulaStateHash, gshDoubted common2.GlobulaStateHash,
	options api.PacketSendOptions) api.PreparedPacketSender {

	v := EmuPhase3NetPacket{
		basePacket: r.defaultBasePacket(sender),
		vectors:    vectors,
	}
	return &emuPacketSender{&v}
}

func (r *EmuPhase3NetPacket) clonePacketFor(t gcp_types.NodeProfile, sendOptions api.PacketSendOptions) packets.PacketParser {
	c := *r
	c.tgt = t.GetShortNodeID()
	return &c
}
