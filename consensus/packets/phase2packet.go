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
	"crypto"

	"github.com/insolar/insolar/core"
	"github.com/pkg/errors"
)

type GlobuleHashSignature [SignatureLength]byte

type Phase2Packet struct {
	// -------------------- Header
	packetHeader PacketHeader

	// -------------------- Section 1
	globuleHashSignature    GlobuleHashSignature
	bitSet                  BitSet
	SignatureHeaderSection1 [SignatureLength]byte

	// -------------------- Section 2 (optional)
	votesAndAnswers         []ReferendumVote
	SignatureHeaderSection2 [SignatureLength]byte
}

func (p2p *Phase2Packet) Clone() ConsensusPacket {
	clone := *p2p
	return &clone
}

func NewPhase2Packet() *Phase2Packet {
	result := &Phase2Packet{}
	result.packetHeader.PacketT = Phase2
	return result
}

func (p2p *Phase2Packet) GetType() PacketType {
	return p2p.packetHeader.PacketT
}

func (p2p *Phase2Packet) GetOrigin() core.ShortNodeID {
	return p2p.packetHeader.OriginNodeID
}

func (p2p *Phase2Packet) GetTarget() core.ShortNodeID {
	return p2p.packetHeader.TargetNodeID
}

func (p2p *Phase2Packet) SetRouting(origin, target core.ShortNodeID) {
	p2p.packetHeader.OriginNodeID = origin
	p2p.packetHeader.TargetNodeID = target
	p2p.packetHeader.HasRouting = true
}

func (p2p *Phase2Packet) Verify(crypto core.CryptographyService, key crypto.PublicKey) error {
	raw, err := p2p.rawFirstPart()
	if err != nil {
		return errors.Wrap(err, "Failed to get raw first part of phase 2 packet")
	}
	valid := crypto.Verify(key, core.SignatureFromBytes(p2p.SignatureHeaderSection1[:]), raw)
	if !valid {
		return errors.New("first part bad signature")
	}

	if !p2p.hasSection2() {
		return nil
	}

	raw, err = p2p.rawSecondPart()
	if err != nil {
		return errors.Wrap(err, "Failed to get raw second part of phase 2 packet")
	}
	valid = crypto.Verify(key, core.SignatureFromBytes(p2p.SignatureHeaderSection2[:]), raw)
	if !valid {
		return errors.New("second part bad signature")
	}
	return nil
}

func (p2p *Phase2Packet) Sign(cryptographyService core.CryptographyService) error {
	raw, err := p2p.rawFirstPart()
	if err != nil {
		return errors.Wrap(err, "Failed to get raw first part of phase 2 packet")
	}
	signature, err := cryptographyService.Sign(raw)
	if err != nil {
		return errors.Wrap(err, "Failed to sign first part of phase 2 packet")
	}
	copy(p2p.SignatureHeaderSection1[:], signature.Bytes()[:SignatureLength])

	if !p2p.hasSection2() {
		return nil
	}

	raw, err = p2p.rawSecondPart()
	if err != nil {
		return errors.Wrap(err, "Failed to get raw second part of phase 2 packet")
	}
	signature, err = cryptographyService.Sign(raw)
	if err != nil {
		return errors.Wrap(err, "Failed to sign second part of phase 2 packet")
	}
	copy(p2p.SignatureHeaderSection2[:], signature.Bytes()[:SignatureLength])

	return nil
}

func (p2p *Phase2Packet) GetPulseNumber() core.PulseNumber {
	return core.PulseNumber(p2p.packetHeader.Pulse)
}

func (p2p *Phase2Packet) IsPhase3Needed() bool {
	return p2p.packetHeader.f00
}

func (p2p *Phase2Packet) hasSection2() bool {
	return p2p.packetHeader.f01
}

func (p2p *Phase2Packet) GetGlobuleHashSignature() GlobuleHashSignature {
	return p2p.globuleHashSignature
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

func (p2p *Phase2Packet) ContainsResponses() bool {
	for _, vote := range p2p.votesAndAnswers {
		if vote.Type() == TypeMissingNodeSupplementaryVote || vote.Type() == TypeMissingNodeClaim {
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
