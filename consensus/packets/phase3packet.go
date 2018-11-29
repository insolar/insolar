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
	"github.com/insolar/insolar/network/transport/packet/types"
	"github.com/pkg/errors"
)

type Phase3Packet struct {
	// -------------------- Header
	packetHeader PacketHeader

	// -------------------- Section 1
	globuleHashSignature    [SignatureLength]byte
	deviantBitSet           BitSet
	SignatureHeaderSection1 [SignatureLength]byte
}

func NewPhase3Packet(globuleHash [SignatureLength]byte, bitSet BitSet) Phase3Packet {
	return Phase3Packet{
		globuleHashSignature: globuleHash,
		deviantBitSet:        bitSet,
	}
}

// SetPacketHeader set routing information for transport level.
func (p3p *Phase3Packet) SetPacketHeader(header *RoutingHeader) error {
	if header.PacketType != types.Phase3 {
		return errors.New("[ Phase3Packet.SetPacketHeader ] wrong packet type")
	}

	p3p.packetHeader.setRoutingFields(header, Phase3)
	return nil
}

// GetPacketHeader get routing information from transport level.
func (p3p *Phase3Packet) GetPacketHeader() (*RoutingHeader, error) {
	header := &RoutingHeader{}

	header.PacketType = types.Phase2
	header.OriginID = p3p.packetHeader.OriginNodeID
	header.TargetID = p3p.packetHeader.TargetNodeID

	return header, nil
}

func (p3p *Phase3Packet) GetBitset() BitSet {
	return p3p.deviantBitSet
}
