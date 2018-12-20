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

func NewPhase2Packet(globuleHashSignature GlobuleHashSignature, bitSet BitSet) *Phase2Packet {
	result := &Phase2Packet{
		globuleHashSignature: globuleHashSignature,
		bitSet:               bitSet,
	}
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
