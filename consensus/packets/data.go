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

type PacketType uint8
type ReferendumType uint8

const (
	Phase1 = PacketType(iota + 1)
	Phase2
)

const HashLength = 64
const SignatureLength = 71
const ReferenceLength = 64

// ----------------------------------PHASE 1--------------------------------

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
	signature uint64
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

// TODO this
func (p1p *Phase1Packet) AddClaim(claim ReferendumClaim) error {

	if phase1PacketSizeForClaims-int(getClaimWithHeaderSize(claim)) < 0 {
		return errors.New("No space for claim")
	}

	p1p.claims = append(p1p.claims, claim)
	p1p.packetHeader.f01 = true
	return nil
}

type PacketHeader struct {
	PacketT    PacketType
	HasRouting bool
	//-----------------
	f01   bool
	f00   bool
	Pulse uint32
	//-----------------
	OriginNodeID core.ShortNodeID
	TargetNodeID core.ShortNodeID
}

func (ph *PacketHeader) setRoutingFields(header *RoutingHeader, packetType PacketType) {
	ph.TargetNodeID = header.TargetID
	ph.OriginNodeID = header.OriginID
	ph.HasRouting = true
	ph.PacketT = packetType
}

// PulseDataExt is a pulse data extension.
type PulseDataExt struct {
	NextPulseDelta uint16
	PrevPulseDelta uint16
	OriginID       [16]byte
	EpochPulseNo   uint32
	PulseTimestamp uint32
	Entropy        core.Entropy
}

// PulseData is a pulse data.
type PulseData struct {
	PulseNumber uint32
	Data        *PulseDataExt
}

type NodePulseProof struct {
	NodeStateHash [HashLength]byte
	NodeSignature [SignatureLength]byte
}

func (npp *NodePulseProof) StateHash() []byte {
	return npp.NodeStateHash[:]
}

func (npp *NodePulseProof) Signature() []byte {
	return npp.NodeSignature[:]
}

// ----------------------------------PHASE 2--------------------------------

type ReferendumVote struct {
	Type   ReferendumType
	Length uint16
}

type NodeListVote struct {
	NodeListCount uint16
	NodeListHash  [32]byte
}

type DeviantBitSet struct {
	CompressedSet     bool
	HighBitLengthFlag bool
	LowBitLength      uint8
	//------------------
	HighBitLength uint8
	Payload       []byte
}

type Phase2Packet struct {
	// -------------------- Header
	packetHeader PacketHeader

	// -------------------- Section 1
	globuleHashSignature    [HashLength]byte
	deviantBitSet           DeviantBitSet
	signatureHeaderSection1 [SignatureLength]byte

	// -------------------- Section 2 (optional)
	votesAndAnswers         []ReferendumVote
	signatureHeaderSection2 [SignatureLength]byte
}

func (phase2Packet *Phase2Packet) isPhase3Needed() bool {
	return phase2Packet.packetHeader.f00
}

func (phase2Packet *Phase2Packet) hasSection2() bool {
	return phase2Packet.packetHeader.f01
}

func (phase2Packet *Phase2Packet) SetPacketHeader(header *RoutingHeader) error {
	if header.PacketType != types.Phase2 {
		return errors.New("Phase2Packet.SetPacketHeader: wrong packet type")
	}

	phase2Packet.packetHeader.setRoutingFields(header, Phase2)

	return nil
}

func (phase2Packet *Phase2Packet) GetPacketHeader() (*RoutingHeader, error) {
	header := &RoutingHeader{}

	if phase2Packet.packetHeader.PacketT != Phase2 {
		return nil, errors.New("Phase2Packet.GetPacketHeader: wrong packet type")
	}

	header.PacketType = types.Phase2
	header.OriginID = phase2Packet.packetHeader.OriginNodeID
	header.TargetID = phase2Packet.packetHeader.TargetNodeID

	return header, nil
}

func (phase2Packet *Phase2Packet) GetGlobuleHashSignature() []byte {
	return phase2Packet.globuleHashSignature[:]
}

func (phase2Packet *Phase2Packet) SetGlobuleHashSignature(globuleHashSignature []byte) error {
	if len(globuleHashSignature) == SignatureLength {
		copy(phase2Packet.globuleHashSignature[:], globuleHashSignature[:SignatureLength])
		return nil
	}

	return errors.New("invalid proof fields len")
}
