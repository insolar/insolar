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
	"context"
	"github.com/insolar/insolar/network/consensus/common/cryptography_containers"
	"github.com/insolar/insolar/network/consensus/common/long_bits"
	"github.com/insolar/insolar/network/consensus/gcpv2/api"
	"github.com/insolar/insolar/network/consensus/gcpv2/gcp_types"
	"io"
	"math/rand"

	"github.com/insolar/insolar/network/consensus/gcpv2/packets"
)

var EmuDefaultPacketBuilder api.PacketBuilder = &emuPacketBuilder{}
var EmuDefaultCryptography api.TransportCryptographyFactory = &emuTransportCryptography{}

func NewEmuTransport(sender api.PacketSender) api.TransportFactory {
	return &emuTransport{sender}
}

type emuTransport struct {
	sender api.PacketSender
}

func (r *emuTransport) GetPacketSender() api.PacketSender {
	return r.sender
}

func (r *emuTransport) GetPacketBuilder(signer cryptography_containers.DigestSigner) api.PacketBuilder {
	return EmuDefaultPacketBuilder
}

func (r *emuTransport) GetCryptographyFactory() api.TransportCryptographyFactory {
	return EmuDefaultCryptography
}

type emuPackerCloner interface {
	clonePacketFor(target gcp_types.NodeProfile, sendOptions api.PacketSendOptions) packets.PacketParser
}

type emuPacketSender struct {
	cloner emuPackerCloner
}

func (r *emuPacketSender) SendTo(ctx context.Context, t gcp_types.NodeProfile, sendOptions api.PacketSendOptions, s api.PacketSender) {
	c := r.cloner.clonePacketFor(t, sendOptions)
	s.SendPacketToTransport(ctx, t, sendOptions, c)
}

func (r *emuPacketSender) SendToMany(ctx context.Context, targetCount int, s api.PacketSender,
	filter func(ctx context.Context, targetIndex int) (gcp_types.NodeProfile, api.PacketSendOptions)) {

	for i := 0; i < targetCount; i++ {
		sendTo, sendOptions := filter(ctx, i)
		if sendTo != nil {
			c := r.cloner.clonePacketFor(sendTo, sendOptions)
			s.SendPacketToTransport(ctx, sendTo, sendOptions, c)
		}

		select {
		case <-ctx.Done():
			return
		default:
		}
	}
}

type emuPacketBuilder struct {
}

func (r *emuPacketBuilder) GetNeighbourhoodSize() gcp_types.NeighbourhoodSizes {
	return gcp_types.NeighbourhoodSizes{NeighbourhoodSize: 5, NeighbourhoodTrustThreshold: 2, JoinersPerNeighbourhood: 2, JoinersBoost: 1}
}

func (r *emuPacketBuilder) PreparePhase0Packet(sender *packets.NodeAnnouncementProfile, pulsarPacket packets.OriginalPulsarPacket,
	options api.PacketSendOptions) api.PreparedPacketSender {
	v := EmuPhase0NetPacket{
		basePacket: basePacket{
			src:       sender.GetNodeID(),
			nodeCount: sender.GetNodeCount(),
			mp:        sender.GetMembershipProfile(),
		},
		pulsePacket: pulsarPacket.(*EmuPulsarNetPacket)}
	return &emuPacketSender{&v}
}

func (r *EmuPhase0NetPacket) clonePacketFor(t gcp_types.NodeProfile, sendOptions api.PacketSendOptions) packets.PacketParser {
	c := *r
	c.tgt = t.GetShortNodeID()
	return &c
}

func (r *emuPacketBuilder) PreparePhase1Packet(sender *packets.NodeAnnouncementProfile, pulsarPacket packets.OriginalPulsarPacket,
	welcome *gcp_types.NodeWelcomePackage, options api.PacketSendOptions) api.PreparedPacketSender {

	pp := pulsarPacket.(*EmuPulsarNetPacket)
	if pp == nil || !pp.pulseData.IsValidPulseData() {
		panic("pulse data is missing or invalid")
	}

	v := EmuPhase1NetPacket{
		EmuPhase0NetPacket: EmuPhase0NetPacket{
			basePacket: basePacket{
				src:         sender.GetNodeID(),
				nodeCount:   sender.GetNodeCount(),
				mp:          sender.GetMembershipProfile(),
				isLeaving:   sender.IsLeaving(),
				leaveReason: sender.GetLeaveReason(),
			},
			pulsePacket: pp},
	}
	v.pn = pp.pulseData.PulseNumber
	v.isRequest = options&api.RequestForPhase1 != 0
	if v.isRequest || options&api.SendWithoutPulseData != 0 {
		v.pulsePacket = nil
	}

	return &emuPacketSender{&v}
}

func (r *EmuPhase1NetPacket) clonePacketFor(t gcp_types.NodeProfile, sendOptions api.PacketSendOptions) packets.PacketParser {
	c := *r
	c.tgt = t.GetShortNodeID()

	//if !t.IsJoiner() {
	//	c.selfIntro = nil
	//}
	if sendOptions&api.SendWithoutPulseData != 0 {
		c.pulsePacket = nil
	}

	return &c
}

func (r *emuPacketBuilder) PreparePhase2Packet(sender *packets.NodeAnnouncementProfile,
	welcome *gcp_types.NodeWelcomePackage, neighbourhood []packets.MembershipAnnouncementReader,
	options api.PacketSendOptions) api.PreparedPacketSender {

	v := EmuPhase2NetPacket{
		basePacket: basePacket{
			src:         sender.GetNodeID(),
			nodeCount:   sender.GetNodeCount(),
			mp:          sender.GetMembershipProfile(),
			isLeaving:   sender.IsLeaving(),
			leaveReason: sender.GetLeaveReason(),
		},
		pulseNumber:   sender.GetPulseNumber(),
		neighbourhood: neighbourhood,
	}
	return &emuPacketSender{&v}
}

func (r *EmuPhase2NetPacket) clonePacketFor(t gcp_types.NodeProfile, sendOptions api.PacketSendOptions) packets.PacketParser {
	c := *r
	c.tgt = t.GetShortNodeID()

	//if !t.IsJoiner() || len(c.intros) == 1 /* the only joiner */ {
	//	c.intros = nil
	//} else {
	//	c.intros = make([]common2.NodeIntroduction, 0, len(r.intros)-1)
	//	for _, ni := range r.intros {
	//		if ni.GetShortNodeID() == t.GetShortNodeID() {
	//			continue
	//		}
	//		c.intros = append(c.intros, ni)
	//	}
	//}
	return &c
}

func (r *emuPacketBuilder) PreparePhase3Packet(sender *packets.NodeAnnouncementProfile, vectors gcp_types.HashedNodeVector,
	options api.PacketSendOptions) api.PreparedPacketSender {

	v := EmuPhase3NetPacket{
		basePacket: basePacket{
			src:       sender.GetNodeID(),
			nodeCount: sender.GetNodeCount(),
			mp:        sender.GetMembershipProfile(),
		},
		pulseNumber: sender.GetPulseNumber(),
		vectors:     vectors,
	}
	return &emuPacketSender{&v}
}

func (r *EmuPhase3NetPacket) clonePacketFor(t gcp_types.NodeProfile, sendOptions api.PacketSendOptions) packets.PacketParser {
	c := *r
	c.tgt = t.GetShortNodeID()
	return &c
}

type emuTransportCryptography struct {
}

func (r *emuTransportCryptography) GetPublicKeyStore(skh cryptography_containers.SignatureKeyHolder) cryptography_containers.PublicKeyStore {
	return nil
}

func (r *emuTransportCryptography) GetPacketDigester() cryptography_containers.DataDigester {
	panic("not implemented")
}

func (r *emuTransportCryptography) GetGshDigester() cryptography_containers.SequenceDigester {
	return &gshDigester{}
}

func (r *emuTransportCryptography) IsDigestMethodSupported(m cryptography_containers.DigestMethod) bool {
	return true
}

func (r *emuTransportCryptography) IsValidDataSignature(data io.Reader, signature cryptography_containers.SignatureHolder) bool {
	return true
}

func (r *emuTransportCryptography) IsSignOfSignatureMethodSupported(m cryptography_containers.SignatureMethod) bool {
	return true
}

func (r *emuTransportCryptography) IsDigestOfSignatureMethodSupported(m cryptography_containers.SignatureMethod) bool {
	return true
}

func (r *emuTransportCryptography) IsSignMethodSupported(m cryptography_containers.SignMethod) bool {
	return true
}

func (r *emuTransportCryptography) IsValidDigestSignature(digest cryptography_containers.DigestHolder, signature cryptography_containers.SignatureHolder) bool {
	return true
}

func (r *emuTransportCryptography) SignDigest(digest cryptography_containers.Digest) cryptography_containers.Signature {
	return cryptography_containers.NewSignature(digest, digest.GetDigestMethod().SignedBy(r.GetSignMethod()))
}

func (r *emuTransportCryptography) GetSignMethod() cryptography_containers.SignMethod {
	return "emuSing"
}

func (r *emuTransportCryptography) GetSignatureVerifierWithPKS(pks cryptography_containers.PublicKeyStore) cryptography_containers.SignatureVerifier {
	return r
}

func (r *emuTransportCryptography) GetDigestFactory() cryptography_containers.DigestFactory {
	return r
}

func (r *emuTransportCryptography) GetNodeSigner(sks cryptography_containers.SecretKeyStore) cryptography_containers.DigestSigner {
	return r
}

type gshDigester struct {
	// TODO do test or a proper digest calc
	rnd      *rand.Rand
	lastSeed int64
}

func (s *gshDigester) AddNext(digest cryptography_containers.DigestHolder) {
	// it is a dirty emulation of digest
	if s.rnd == nil {
		s.rnd = rand.New(rand.NewSource(0))
	}
	s.lastSeed = int64(s.rnd.Uint64() ^ digest.FoldToUint64())
	s.rnd.Seed(s.lastSeed)
}

func (s *gshDigester) GetDigestMethod() cryptography_containers.DigestMethod {
	return "emuDigest64"
}

func (s *gshDigester) ForkSequence() cryptography_containers.SequenceDigester {
	cp := gshDigester{}
	if s.rnd != nil {
		cp.rnd = rand.New(rand.NewSource(s.lastSeed))
	}
	return &cp
}

func (s *gshDigester) FinishSequence() cryptography_containers.Digest {
	if s.rnd == nil {
		panic("nothing")
	}
	bits := long_bits.NewBits64(s.rnd.Uint64())
	s.rnd = nil
	return cryptography_containers.NewDigest(&bits, s.GetDigestMethod())
}
