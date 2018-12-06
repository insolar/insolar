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

type Phase2Packet struct {
	// -------------------- Header
	packetHeader PacketHeader

	// -------------------- Section 1
	globuleHashSignature    [HashLength]byte
	bitSet                  BitSet
	SignatureHeaderSection1 [SignatureLength]byte

	// -------------------- Section 2 (optional)
	votesAndAnswers         []ReferendumVote
	SignatureHeaderSection2 [SignatureLength]byte
}

func (p2p *Phase2Packet) GetPulseNumber() core.PulseNumber {
	return core.PulseNumber(p2p.packetHeader.Pulse)
}

func (p2p *Phase2Packet) isPhase3Needed() bool {
	return p2p.packetHeader.f00
}

func (p2p *Phase2Packet) hasSection2() bool {
	return p2p.packetHeader.f01
}

func (p2p *Phase2Packet) SetPacketHeader(header *RoutingHeader) error {
	if header.PacketType != types.Phase2 {
		return errors.New("Phase2Packet.SetPacketHeader: wrong packet type")
	}

	p2p.packetHeader.setRoutingFields(header, Phase2)

	return nil
}

func (p2p *Phase2Packet) GetPacketHeader() (*RoutingHeader, error) {
	header := &RoutingHeader{}

	if p2p.packetHeader.PacketT != Phase2 {
		return nil, errors.New("Phase2Packet.GetPacketHeader: wrong packet type")
	}

	header.PacketType = types.Phase2
	header.OriginID = p2p.packetHeader.OriginNodeID
	header.TargetID = p2p.packetHeader.TargetNodeID

	return header, nil
}

func (p2p *Phase2Packet) GetGlobuleHashSignature() []byte {
	return p2p.globuleHashSignature[:]
}

func (p2p *Phase2Packet) SetGlobuleHashSignature(globuleHashSignature []byte) error {
	if len(globuleHashSignature) == SignatureLength {
		copy(p2p.globuleHashSignature[:], globuleHashSignature[:SignatureLength])
		return nil
	}

	return errors.New("invalid proof fields len")
}

func (p2p *Phase2Packet) GetBitSet() BitSet {
	return p2p.bitSet
}

func (p2p *Phase2Packet) SetBitSet(bitset BitSet) {
	p2p.bitSet = bitset
}
