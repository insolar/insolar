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

func (p2p *Phase2Packet) ContainsRequests() bool {
	for _, vote := range p2p.votesAndAnswers {
		if vote.Type() == TypeMissingNode {
			return true
		}
	}
	return false
}

func (p2p *Phase2Packet) AddVote(vote ReferendumVote) {
	// TODO: check size

	p2p.votesAndAnswers = append(p2p.votesAndAnswers, vote)
	p2p.packetHeader.f01 = true
}

func (p2p *Phase2Packet) GetVotes() []ReferendumVote {
	return p2p.votesAndAnswers
}
