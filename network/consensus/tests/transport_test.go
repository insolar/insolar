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
	"io"
	"math/rand"

	"github.com/insolar/insolar/network/consensus/common"
	common2 "github.com/insolar/insolar/network/consensus/gcpv2/common"
	"github.com/insolar/insolar/network/consensus/gcpv2/core"
	"github.com/insolar/insolar/network/consensus/gcpv2/nodeset"
	"github.com/insolar/insolar/network/consensus/gcpv2/packets"
)

var EmuDefaultPacketBuilder core.PacketBuilder = &emuPacketBuilder{}
var EmuDefaultCryptography core.TransportCryptographyFactory = &emuTransportCryptography{}

func NewEmuTransport(sender core.PacketSender) core.TransportFactory {
	return &emuTransport{sender}
}

type emuTransport struct {
	sender core.PacketSender
}

func (r *emuTransport) GetPacketSender() core.PacketSender {
	return r.sender
}

func (r *emuTransport) GetPacketBuilder(signer common.DigestSigner) core.PacketBuilder {
	return EmuDefaultPacketBuilder
}

func (r *emuTransport) GetCryptographyFactory() core.TransportCryptographyFactory {
	return EmuDefaultCryptography
}

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
}

func (r *emuPacketBuilder) GetNeighbourhoodSize(populationCount int) common2.NeighbourhoodSizes {
	return common2.NeighbourhoodSizes{NeighbourhoodSize: 5, NeighbourhoodTrustThreshold: 2, JoinersPerNeighbourhood: 2, JoinersBoost: 1}
}

func (r *emuPacketBuilder) PreparePhase0Packet(sender common2.NodeProfile, pulsarPacket common2.OriginalPulsarPacket,
	mp common2.MembershipProfile, nodeCount int,
	options core.PacketSendOptions) core.PreparedPacketSender {
	v := EmuPhase0NetPacket{
		basePacket: basePacket{
			src:       sender.GetShortNodeID(),
			nodeCount: uint16(nodeCount),
			mp:        mp,
		},
		pulsePacket: pulsarPacket.(*EmuPulsarNetPacket)}
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
			basePacket: basePacket{
				src:       sender.GetShortNodeID(),
				nodeCount: uint16(nodeCount),
				mp:        mp,
			},
			pulsePacket: pp},
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
		basePacket: basePacket{
			src:       sender.GetShortNodeID(),
			nodeCount: uint16(nodeCount),
			mp:        mp,
		},
		pulseNumber:   pd.PulseNumber,
		neighbourhood: neighbourhood,
		intros:        intros}
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
		basePacket: basePacket{
			src: sender.GetShortNodeID()},
		bitset:      bitset,
		pulseNumber: pd.PulseNumber,
		gshTrusted:  gshTrusted,
		gshDoubted:  gshDoubted}
	return &emuPacketSender{&v}
}

func (r *EmuPhase3NetPacket) clonePacketFor(t common2.NodeProfile, sendOptions core.PacketSendOptions) packets.PacketParser {
	c := *r
	c.tgt = t.GetShortNodeID()
	return &c
}

type emuTransportCryptography struct {
}

func (r *emuTransportCryptography) GetPacketDigester() common.DataDigester {
	panic("not implemented")
}

func (r *emuTransportCryptography) GetGshDigester() common.SequenceDigester {
	return &gshDigester{}
}

func (r *emuTransportCryptography) IsDigestMethodSupported(m common.DigestMethod) bool {
	return true
}

func (r *emuTransportCryptography) IsValidDataSignature(data io.Reader, signature common.SignatureHolder) bool {
	return true
}

func (r *emuTransportCryptography) IsSignOfSignatureMethodSupported(m common.SignatureMethod) bool {
	return true
}

func (r *emuTransportCryptography) IsSignMethodSupported(m common.SignMethod) bool {
	return true
}

func (r *emuTransportCryptography) IsValidDigestSignature(digest common.DigestHolder, signature common.SignatureHolder) bool {
	return true
}

func (r *emuTransportCryptography) SignDigest(digest common.Digest) common.Signature {
	return common.NewSignature(digest, digest.GetDigestMethod().SignedBy(r.GetSignMethod()))
}

func (r *emuTransportCryptography) GetSignMethod() common.SignMethod {
	return "emuSing"
}

func (r *emuTransportCryptography) GetSignatureVerifierWithPKS(pks common.PublicKeyStore) common.SignatureVerifier {
	return r
}

func (r *emuTransportCryptography) GetDigestFactory() common.DigestFactory {
	return r
}

func (r *emuTransportCryptography) GetNodeSigner(sks common.SecretKeyStore) common.DigestSigner {
	return r
}

type gshDigester struct {
	// TODO do test or a proper digest calc
	rnd      *rand.Rand
	lastSeed int64
}

func (s *gshDigester) AddNext(digest common.DigestHolder) {
	// it is a dirty emulation of digest
	if s.rnd == nil {
		s.rnd = rand.New(rand.NewSource(0))
	}
	s.lastSeed = int64(s.rnd.Uint64() ^ digest.FoldToUint64())
	s.rnd.Seed(s.lastSeed)
}

func (s *gshDigester) GetDigestMethod() common.DigestMethod {
	return "emuDigest64"
}

func (s *gshDigester) ForkSequence() common.SequenceDigester {
	cp := gshDigester{}
	if s.rnd != nil {
		cp.rnd = rand.New(rand.NewSource(s.lastSeed))
	}
	return &cp
}

func (s *gshDigester) FinishSequence() common.Digest {
	if s.rnd == nil {
		panic("nothing")
	}
	bits := common.NewBits64(s.rnd.Uint64())
	s.rnd = nil
	return common.NewDigest(&bits, s.GetDigestMethod())
}
