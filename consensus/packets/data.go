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
type ClaimType uint8
type ReferendumType uint8

const (
	Phase1 = PacketType(iota + 1)
	Phase2
)

const (
	TypeNodeJoinClaim = ClaimType(iota + 1)
	TypeCapabilityPollingAndActivation
	TypeNodeViolationBlame
	TypeNodeBroadcast
	TypeNodeLeaveClaim
)

const HashLength = 64
const SignatureLength = 71
const ReferenceLength = 64

// ----------------------------------PHASE 1--------------------------------

type Phase1Packet struct {
	// -------------------- Header
	PacketHeader PacketHeader

	// -------------------- Section 1 ( Pulse )
	PulseData      PulseDataExt // optional
	ProofNodePulse NodePulseProof

	// -------------------- Section 2 ( Claims ) ( optional )
	claims []ReferendumClaim

	// --------------------
	// signature contains signature of Header + Section 1 + Section 2
	Signature []byte
}

func NewPhase1Packet() *Phase1Packet {
	return &Phase1Packet{
		Signature: make([]byte, SignatureLength),
	}
}

func (p1p *Phase1Packet) hasPulseDataExt() bool { // nolint: megacheck
	return p1p.PacketHeader.f00
}

func (p1p *Phase1Packet) hasSection2() bool {
	return p1p.PacketHeader.f01
}

func (p1p *Phase1Packet) SetPacketHeader(header *RoutingHeader) error {
	if header.PacketType != types.Phase1 {
		return errors.New("Phase1Packet.SetPacketHeader: wrong packet type")
	}
	p1p.PacketHeader.setRoutingFields(header, Phase1)

	return nil
}

func (p1p *Phase1Packet) GetPulse() core.Pulse {
	//TODO: need convert method with pulse signature check
	return core.Pulse{
		PulseNumber: core.PulseNumber(p1p.PacketHeader.Pulse),
		Entropy:     p1p.PulseData.Entropy,
	}
}

func (p1p *Phase1Packet) GetPulseProof() *NodePulseProof {
	return &p1p.ProofNodePulse
}

func (p1p *Phase1Packet) GetPacketHeader() (*RoutingHeader, error) {
	header := &RoutingHeader{}

	if p1p.PacketHeader.PacketT != Phase1 {
		return nil, errors.New("Phase1Packet.GetPacketHeader: wrong packet type")
	}

	header.PacketType = types.Phase1
	header.OriginID = p1p.PacketHeader.OriginNodeID
	header.TargetID = p1p.PacketHeader.TargetNodeID

	return header, nil
}

// SetPulseProof sets PulseProof and check struct fields len, returns error if invalid len
func (p1p *Phase1Packet) SetPulseProof(proofStateHash, proofSignature []byte) error {
	if len(proofStateHash) == HashLength && len(proofSignature) == SignatureLength {
		copy(p1p.ProofNodePulse.NodeStateHash[:], proofStateHash[:HashLength])
		copy(p1p.ProofNodePulse.NodeSignature[:], proofSignature[:SignatureLength])
		return nil
	}

	return errors.New("invalid proof fields len")
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

// --------------REFERENDUM--------------

type ReferendumClaim interface {
	Serializer
	Type() ClaimType
	Length() uint16
}

// NodeBroadcast is a broadcast of info. Must be brief and only one entry per node.
// Type 4.
type NodeBroadcast struct {
	EmergencyLevel uint8
	length         uint16
}

func (nb *NodeBroadcast) Type() ClaimType {
	return TypeNodeBroadcast
}

func (nb *NodeBroadcast) Length() uint16 {
	return nb.length
}

// CapabilityPoolingAndActivation is a type 3.
type CapabilityPoolingAndActivation struct {
	PollingFlags   uint16
	CapabilityType uint16
	CapabilityRef  [ReferenceLength]byte
	length         uint16
}

func (cpa *CapabilityPoolingAndActivation) Type() ClaimType {
	return TypeCapabilityPollingAndActivation
}

func (cpa *CapabilityPoolingAndActivation) Length() uint16 {
	return cpa.length
}

// NodeViolationBlame is a type 2.
type NodeViolationBlame struct {
	BlameNodeID   uint32
	TypeViolation uint8
	claimType     ClaimType
	length        uint16
}

func (nvb *NodeViolationBlame) Type() ClaimType {
	return TypeNodeViolationBlame
}

func (nvb *NodeViolationBlame) Length() uint16 {
	return nvb.length
}

// NodeJoinClaim is a type 1, len == 272.
type NodeJoinClaim struct {
	NodeID                  uint32
	RelayNodeID             uint32
	ProtocolVersionAndFlags uint32
	JoinsAfter              uint32
	NodeRoleRecID           uint32
	NodeRef                 core.RecordRef
	NodePK                  [64]byte
	//length uint16
}

func (njc *NodeJoinClaim) Type() ClaimType {
	return TypeNodeJoinClaim
}

func (njc *NodeJoinClaim) Length() uint16 {
	return 0
}

// NodeLeaveClaim can be the only be issued by the node itself and must be the only claim record.
// Should be executed with the next pulse. Type 1, len == 0.
type NodeLeaveClaim struct {
	length uint16
}

func (nlc *NodeLeaveClaim) Type() ClaimType {
	return TypeNodeLeaveClaim
}

func (nlc *NodeLeaveClaim) Length() uint16 {
	return nlc.length
}

func NewNodeJoinClaim() *NodeJoinClaim {
	return &NodeJoinClaim{
		//length: 272,
	}
}

func NewNodViolationBlame() *NodeViolationBlame {
	return &NodeViolationBlame{
		claimType: TypeNodeViolationBlame,
	}
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
	PacketHeader PacketHeader

	// -------------------- Section 1
	GlobuleHashSignature    []byte
	DeviantBitSet           DeviantBitSet
	SignatureHeaderSection1 []byte

	// -------------------- Section 2 (optional)
	votesAndAnswers         []ReferendumVote
	SignatureHeaderSection2 []byte
}

func NewPhase2Packet() *Phase2Packet {
	return &Phase2Packet{
		SignatureHeaderSection1: make([]byte, SignatureLength),
		SignatureHeaderSection2: make([]byte, SignatureLength),
	}
}

func (phase2Packet *Phase2Packet) isPhase3Needed() bool {
	return phase2Packet.PacketHeader.f00
}

func (phase2Packet *Phase2Packet) hasSection2() bool {
	return phase2Packet.PacketHeader.f01
}

func (phase2Packet *Phase2Packet) SetPacketHeader(header *RoutingHeader) error {
	if header.PacketType != types.Phase2 {
		return errors.New("Phase2Packet.SetPacketHeader: wrong packet type")
	}

	phase2Packet.PacketHeader.setRoutingFields(header, Phase2)

	return nil
}

func (phase2Packet *Phase2Packet) GetPacketHeader() (*RoutingHeader, error) {
	header := &RoutingHeader{}

	if phase2Packet.PacketHeader.PacketT != Phase2 {
		return nil, errors.New("Phase2Packet.GetPacketHeader: wrong packet type")
	}

	header.PacketType = types.Phase2
	header.OriginID = phase2Packet.PacketHeader.OriginNodeID
	header.TargetID = phase2Packet.PacketHeader.TargetNodeID

	return header, nil
}

func (phase2Packet *Phase2Packet) GetGlobuleHashSignature() []byte {
	return phase2Packet.GlobuleHashSignature[:]
}

func (phase2Packet *Phase2Packet) SetGlobuleHashSignature(globuleHashSignature []byte) error {
	if len(globuleHashSignature) == SignatureLength {
		copy(phase2Packet.GlobuleHashSignature[:], globuleHashSignature[:SignatureLength])
		return nil
	}

	return errors.New("invalid proof fields len")
}
