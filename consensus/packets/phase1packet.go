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

package packets

import (
	"github.com/pkg/errors"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/platformpolicy/keys"
)

type Phase1Packet struct {
	// -------------------- Header
	packetHeader PacketHeader

	// -------------------- Section 1 ( Pulse )
	pulseData      PulseDataExt // optional
	proofNodePulse NodePulseProof

	// -------------------- Section 2 ( Claims ) ( optional )
	claims []ReferendumClaim

	// --------------------
	// signature contains signature of Header + Section 1 + Section 2
	Signature [SignatureLength]byte
}

func (p1p *Phase1Packet) Clone() ConsensusPacket {
	clone := *p1p
	clone.claims = make([]ReferendumClaim, len(p1p.claims))
	for i := 0; i < len(p1p.claims); i++ {
		clone.claims[i] = p1p.claims[i].Clone()
	}
	return &clone
}

func (p1p *Phase1Packet) GetOrigin() insolar.ShortNodeID {
	return p1p.packetHeader.OriginNodeID
}

func (p1p *Phase1Packet) GetTarget() insolar.ShortNodeID {
	return p1p.packetHeader.TargetNodeID
}

func (p1p *Phase1Packet) SetRouting(origin, target insolar.ShortNodeID) {
	p1p.packetHeader.OriginNodeID = origin
	p1p.packetHeader.TargetNodeID = target
	p1p.packetHeader.HasRouting = true
}

func (p1p *Phase1Packet) GetType() PacketType {
	return p1p.packetHeader.PacketT
}

func (p1p *Phase1Packet) Verify(crypto insolar.CryptographyService, key keys.PublicKey) error {
	raw, err := p1p.rawBytes()
	if err != nil {
		return errors.Wrap(err, "Failed to get raw part of phase 1 packet")
	}
	valid := crypto.Verify(key, insolar.SignatureFromBytes(p1p.Signature[:]), raw)
	if !valid {
		return errors.New("bad signature")
	}
	return nil
}

func (p1p *Phase1Packet) Sign(cryptographyService insolar.CryptographyService) error {
	raw, err := p1p.rawBytes()
	if err != nil {
		return errors.Wrap(err, "Failed to get raw part of phase 1 packet")
	}
	signature, err := cryptographyService.Sign(raw)
	if err != nil {
		return errors.Wrap(err, "Failed to sign phase 1 packet")
	}
	copy(p1p.Signature[:], signature.Bytes()[:SignatureLength])
	return nil
}

func NewPhase1Packet(pulse insolar.Pulse) *Phase1Packet {
	result := &Phase1Packet{}
	result.packetHeader.PacketT = Phase1
	result.packetHeader.Pulse = uint32(pulse.PulseNumber)
	result.pulseData = pulseToDataExt(pulse)
	return result
}

func pulseToDataExt(pulse insolar.Pulse) PulseDataExt {
	result := PulseDataExt{}
	result.Entropy = pulse.Entropy
	result.EpochPulseNo = uint32(pulse.EpochPulseNumber)
	result.NextPulseDelta = uint16(pulse.NextPulseNumber - pulse.PulseNumber)
	result.PrevPulseDelta = uint16(pulse.PulseNumber - pulse.PrevPulseNumber)
	result.OriginID = pulse.OriginID
	result.PulseTimestamp = uint32(pulse.PulseTimestamp)
	return result
}

func dataToPulse(number insolar.PulseNumber, data PulseDataExt) insolar.Pulse {
	result := insolar.Pulse{}
	result.PulseNumber = number
	result.Entropy = data.Entropy
	result.EpochPulseNumber = int(data.EpochPulseNo)
	result.NextPulseNumber = number + insolar.PulseNumber(data.NextPulseDelta)
	result.PrevPulseNumber = number - insolar.PulseNumber(data.PrevPulseDelta)
	result.OriginID = data.OriginID
	result.PulseTimestamp = int64(data.PulseTimestamp)
	return result
}

func (p1p *Phase1Packet) hasPulseDataExt() bool { // nolint: megacheck
	return p1p.packetHeader.f00
}

func (p1p *Phase1Packet) hasSection2() bool {
	return p1p.packetHeader.f01
}

func (p1p *Phase1Packet) GetPulseNumber() insolar.PulseNumber {
	return insolar.PulseNumber(p1p.packetHeader.Pulse)
}

func (p1p *Phase1Packet) GetPulse() insolar.Pulse {
	return dataToPulse(insolar.PulseNumber(p1p.packetHeader.Pulse), p1p.pulseData)
}

func (p1p *Phase1Packet) GetPulseProof() *NodePulseProof {
	return &p1p.proofNodePulse
}

// SetPulseProof sets PulseProof and check struct fields len, returns error if invalid len
func (p1p *Phase1Packet) SetPulseProof(proofStateHash, proofSignature []byte) error {
	if len(proofStateHash) == HashLength && len(proofSignature) == SignatureLength {
		copy(p1p.proofNodePulse.NodeStateHash[:], proofStateHash[:HashLength])
		copy(p1p.proofNodePulse.NodeSignature[:], proofSignature[:SignatureLength])
		return nil
	}

	return errors.New("invalid proof fields len")
}

// AddClaim adds claim if phase1Packet has space for it and returns true, otherwise returns false
func (p1p *Phase1Packet) AddClaim(claim ReferendumClaim) bool {
	if claim == nil {
		log.Warn("claim is nil")
		return false
	}

	getClaimSize := func(claims ...ReferendumClaim) int {
		result := 0
		for _, cl := range claims {
			result += int(getClaimWithHeaderSize(cl))
			result += HeaderSize
		}
		return result
	}

	claimSize := getClaimSize(append(p1p.claims, claim)...)

	if claimSize > phase1PacketSizeForClaims {
		return false
	}

	p1p.claims = append(p1p.claims, claim)
	p1p.updateHeader()
	return true
}

func (p1p *Phase1Packet) updateHeader() {
	p1p.packetHeader.f01 = len(p1p.claims) > 0
}

// TODO: I AM AWFUL WORKAROUND, NEED TO REWORK
func (p1p *Phase1Packet) RemoveAnnounceClaim() {
	for i, claim := range p1p.claims {
		if claim.Type() == TypeNodeAnnounceClaim {
			p1p.claims = append(p1p.claims[:i], p1p.claims[i+1:]...)
		}
	}
	p1p.updateHeader()
}

func (p1p *Phase1Packet) GetAnnounceClaim() *NodeAnnounceClaim {
	for _, claim := range p1p.claims {
		c, ok := claim.(*NodeAnnounceClaim)
		if !ok {
			continue
		}
		return c
	}
	return nil
}

func (p1p *Phase1Packet) GetClaims() []ReferendumClaim {
	return p1p.claims
}
