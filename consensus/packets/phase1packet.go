/*
 *    Copyright 2018 Insolar
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
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
			result += claimHeaderSize
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
