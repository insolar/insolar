//
// Modified BSD 3-Clause Clear License
//
// Copyright (config) 2019 Insolar Technologies GmbH
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
//  * Neither the addr of Insolar Technologies GmbH nor the names of its contributors
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
//    (config) distribute this software (including without limitation in source code, binary or
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

	"github.com/insolar/insolar/network/consensus/common"
	common2 "github.com/insolar/insolar/network/consensus/gcpv2/common"
	"github.com/insolar/insolar/network/consensus/gcpv2/core"
	"github.com/insolar/insolar/network/consensus/gcpv2/nodeset"
	"github.com/insolar/insolar/network/consensus/gcpv2/packets"
)

type emuPackerCloner interface {
	clonePacketFor(target common2.NodeProfile, sendOptions core.PacketSendOptions) packets.PacketParser
}

type emuPacketSender struct {
	cloner emuPackerCloner
}

func (r *emuPacketSender) SendTo(t common2.NodeProfile, sendOptions core.PacketSendOptions, s core.PacketSender) {
	c := r.cloner.clonePacketFor(t, sendOptions)
	s.SendPacketToTransport(t, sendOptions, c)
}

type emuPacketBuilder struct {
	crypto      core.TransportCryptographyFactory
	localConfig core.LocalNodeConfiguration
}

func NewEmuPacketBuilder(crypto core.TransportCryptographyFactory, localConfig core.LocalNodeConfiguration) core.PacketBuilder {
	return &emuPacketBuilder{
		crypto:      crypto,
		localConfig: localConfig,
	}
}

func (r *emuPacketBuilder) GetNeighbourhoodSize(populationCount int) common2.NeighbourhoodSizes {
	return common2.NeighbourhoodSizes{NeighbourhoodSize: 5, NeighbourhoodTrustThreshold: 2, JoinersPerNeighbourhood: 2, JoinersBoost: 1}
}

func (r *emuPacketBuilder) defaultSign() common.SignedDigest {
	digester := r.crypto.GetDigestFactory().GetPacketDigester()
	digest := digester.GetDigestOf(bytes.NewReader([]byte{1, 3, 3, 7}))
	signer := r.crypto.GetNodeSigner(r.localConfig.GetSecretKeyStore())
	signature := signer.SignDigest(digest)
	sd := common.NewSignedDigest(digest, signature)
	return sd
}

func (r *emuPacketBuilder) defaultBasePacket(sender common2.NodeProfile, mp common2.MembershipProfile, nodeCount int) basePacket {
	sd := r.defaultSign()

	return basePacket{
		src:       sender.GetShortNodeID(),
		mp:        mp,
		nodeCount: uint16(nodeCount),
		sd:        sd,
	}
}

func (r *emuPacketBuilder) PreparePhase0Packet(sender common2.NodeProfile, pulsarPacket common2.OriginalPulsarPacket,
	mp common2.MembershipProfile, nodeCount int,
	options core.PacketSendOptions) core.PreparedPacketSender {
	v := EmuPhase0NetPacket{
		basePacket:  r.defaultBasePacket(sender, mp, nodeCount),
		pulsePacket: pulsarPacket.(*EmuPulsarNetPacket),
	}
	return &emuPacketSender{&v}
}

func (r *EmuPhase0NetPacket) clonePacketFor(t common2.NodeProfile, sendOptions core.PacketSendOptions) packets.PacketParser {
	c := *r
	c.tgt = t.GetShortNodeID()
	return &c
}

func (r *emuPacketBuilder) PreparePhase1Packet(sender common2.NodeProfile, pulsarPacket common2.OriginalPulsarPacket,
	mp common2.MembershipProfile, nodeCount int,
	options core.PacketSendOptions) core.PreparedPacketSender {

	pp := pulsarPacket.(*EmuPulsarNetPacket)
	if pp == nil || !pp.pulseData.IsValidPulseData() {
		panic("pulse data is missing or invalid")
	}

	v := EmuPhase1NetPacket{
		EmuPhase0NetPacket: EmuPhase0NetPacket{
			basePacket:  r.defaultBasePacket(sender, mp, nodeCount),
			pulsePacket: pp,
		},
		selfIntro: sender.GetIntroduction(),
	}
	v.pn = pp.pulseData.PulseNumber
	v.isRequest = options&core.RequestForPhase1 != 0
	if v.isRequest || options&core.SendWithoutPulseData != 0 {
		v.pulsePacket = nil
	}

	return &emuPacketSender{&v}
}

func (r *EmuPhase1NetPacket) clonePacketFor(t common2.NodeProfile, sendOptions core.PacketSendOptions) packets.PacketParser {
	c := *r
	c.tgt = t.GetShortNodeID()

	if !t.IsJoiner() {
		c.selfIntro = nil
	}
	if sendOptions&core.SendWithoutPulseData != 0 {
		c.pulsePacket = nil
	}

	return &c
}

func (r *emuPacketBuilder) PreparePhase2Packet(sender common2.NodeProfile, pd common.PulseData,
	mp common2.MembershipProfile, nodeCount int,
	neighbourhood []packets.NodeStateHashReportReader,
	intros []common2.NodeIntroduction, options core.PacketSendOptions) core.PreparedPacketSender {

	v := EmuPhase2NetPacket{
		basePacket:    r.defaultBasePacket(sender, mp, nodeCount),
		pulseNumber:   pd.PulseNumber,
		neighbourhood: neighbourhood,
		intros:        intros,
	}
	return &emuPacketSender{&v}
}

func (r *EmuPhase2NetPacket) clonePacketFor(t common2.NodeProfile, sendOptions core.PacketSendOptions) packets.PacketParser {
	c := *r
	c.tgt = t.GetShortNodeID()

	if !t.IsJoiner() || len(c.intros) == 1 /* the only joiner */ {
		c.intros = nil
	} else {
		c.intros = make([]common2.NodeIntroduction, 0, len(r.intros)-1)
		for _, ni := range r.intros {
			if ni.GetShortNodeID() == t.GetShortNodeID() {
				continue
			}
			c.intros = append(c.intros, ni)
		}
	}
	return &c
}

func (r *emuPacketBuilder) PreparePhase3Packet(sender common2.NodeProfile, pd common.PulseData,
	bitset nodeset.NodeBitset, gshTrusted common2.GlobulaStateHash, gshDoubted common2.GlobulaStateHash,
	options core.PacketSendOptions) core.PreparedPacketSender {

	v := EmuPhase3NetPacket{
		basePacket:  r.defaultBasePacket(sender, common2.MembershipProfile{}, 0),
		bitset:      bitset,
		pulseNumber: pd.PulseNumber,
		gshTrusted:  gshTrusted,
		gshDoubted:  gshDoubted,
	}
	return &emuPacketSender{&v}
}

func (r *EmuPhase3NetPacket) clonePacketFor(t common2.NodeProfile, sendOptions core.PacketSendOptions) packets.PacketParser {
	c := *r
	c.tgt = t.GetShortNodeID()
	return &c
}
