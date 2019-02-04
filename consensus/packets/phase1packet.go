/*
 * The Clear BSD License
 *
 * Copyright (c) 2019 Insolar Technologies
 *
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without modification, are permitted (subject to the limitations in the disclaimer below) provided that the following conditions are met:
 *
 *  Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.
 *  Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.
 *  Neither the name of Insolar Technologies nor the names of its contributors may be used to endorse or promote products derived from this software without specific prior written permission.
 *
 * NO EXPRESS OR IMPLIED LICENSES TO ANY PARTY'S PATENT RIGHTS ARE GRANTED BY THIS LICENSE. THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 *
 */

package packets

import (
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/network/transport/packet/types"
	"github.com/pkg/errors"
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

func NewPhase1Packet() *Phase1Packet {
	return &Phase1Packet{}
}

func (p1p *Phase1Packet) hasPulseDataExt() bool { // nolint: megacheck
	return p1p.packetHeader.f00
}

func (p1p *Phase1Packet) hasSection2() bool {
	return p1p.packetHeader.f01
}

func (p1p *Phase1Packet) SetPacketHeader(header *RoutingHeader) error {
	if header.PacketType != types.Phase1 {
		return errors.New("Phase1Packet.SetPacketHeader: wrong packet type")
	}
	p1p.packetHeader.setRoutingFields(header, Phase1)

	return nil
}

func (p1p *Phase1Packet) GetPulseNumber() core.PulseNumber {
	return core.PulseNumber(p1p.packetHeader.Pulse)
}

func (p1p *Phase1Packet) GetPulse() core.Pulse {
	//TODO: need convert method with pulse signature check
	return core.Pulse{
		PulseNumber: core.PulseNumber(p1p.packetHeader.Pulse),
		Entropy:     p1p.pulseData.Entropy,
	}
}

func (p1p *Phase1Packet) GetPulseProof() *NodePulseProof {
	return &p1p.proofNodePulse
}

func (p1p *Phase1Packet) GetPacketHeader() (*RoutingHeader, error) {
	header := &RoutingHeader{}

	if p1p.packetHeader.PacketT != Phase1 {
		return nil, errors.New("Phase1Packet.GetPacketHeader: wrong packet type")
	}

	header.PacketType = types.Phase1
	header.OriginID = p1p.packetHeader.OriginNodeID
	header.TargetID = p1p.packetHeader.TargetNodeID

	return header, nil
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
	p1p.packetHeader.f01 = true
	return true
}

func (p1p *Phase1Packet) GetClaims() []ReferendumClaim {
	return p1p.claims
}

func (ph *PacketHeader) setRoutingFields(header *RoutingHeader, packetType PacketType) {
	ph.TargetNodeID = header.TargetID
	ph.OriginNodeID = header.OriginID
	ph.HasRouting = true
	ph.PacketT = packetType
}
