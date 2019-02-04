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

type Phase3Packet struct {
	// -------------------- Header
	packetHeader PacketHeader

	// -------------------- Section 1
	globuleHashSignature    GlobuleHashSignature
	bitset                  BitSet
	SignatureHeaderSection1 [SignatureLength]byte
}

func (p3p *Phase3Packet) Clone() ConsensusPacket {
	clone := *p3p
	return &clone
}

func (p3p *Phase3Packet) GetType() PacketType {
	return p3p.packetHeader.PacketT
}

func (p3p *Phase3Packet) GetOrigin() core.ShortNodeID {
	return p3p.packetHeader.OriginNodeID
}

func (p3p *Phase3Packet) GetTarget() core.ShortNodeID {
	return p3p.packetHeader.TargetNodeID
}

func (p3p *Phase3Packet) SetRouting(origin, target core.ShortNodeID) {
	p3p.packetHeader.OriginNodeID = origin
	p3p.packetHeader.TargetNodeID = target
	p3p.packetHeader.HasRouting = true
}

func (p3p *Phase3Packet) Verify(crypto core.CryptographyService, key crypto.PublicKey) error {
	raw, err := p3p.rawBytes()
	if err != nil {
		return errors.Wrap(err, "Failed to get raw part of phase 3 packet")
	}
	valid := crypto.Verify(key, core.SignatureFromBytes(p3p.SignatureHeaderSection1[:]), raw)
	if !valid {
		return errors.New("bad signature")
	}
	return nil
}

func (p3p *Phase3Packet) Sign(cryptographyService core.CryptographyService) error {
	raw, err := p3p.rawBytes()
	if err != nil {
		return errors.Wrap(err, "Failed to get raw part of phase 3 packet")
	}
	signature, err := cryptographyService.Sign(raw)
	if err != nil {
		return errors.Wrap(err, "Failed to sign phase 3 packet")
	}
	copy(p3p.SignatureHeaderSection1[:], signature.Bytes()[:SignatureLength])
	return nil
}

func NewPhase3Packet(globuleHashSignature GlobuleHashSignature, bitSet BitSet) *Phase3Packet {
	result := &Phase3Packet{
		globuleHashSignature: globuleHashSignature,
		bitset:               bitSet,
	}
	result.packetHeader.PacketT = Phase3
	return result
}

func (p3p *Phase3Packet) GetBitset() BitSet {
	return p3p.bitset
}

func (p3p *Phase3Packet) GetGlobuleHashSignature() GlobuleHashSignature {
	return p3p.globuleHashSignature
}

func (p3p *Phase3Packet) GetPulseNumber() core.PulseNumber {
	return core.PulseNumber(p3p.packetHeader.Pulse)
}
