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
)

type Phase3Packet struct {
	// -------------------- Header
	packetHeader PacketHeader

	// -------------------- Section 1
	globuleHashSignature    GlobuleHashSignature
	bitset                  BitSet
	SignatureHeaderSection1 [SignatureLength]byte
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
	panic("implement me")
}

func (p3p *Phase3Packet) Sign(crypto core.CryptographyService) error {
	panic("implement me")
}

func NewPhase3Packet(globuleHashSignature GlobuleHashSignature, bitSet BitSet) Phase3Packet {
	return Phase3Packet{
		globuleHashSignature: globuleHashSignature,
		bitset:               bitSet,
	}
}

func (p3p *Phase3Packet) GetBitset() BitSet {
	return p3p.bitset
}

func (p3p *Phase3Packet) GetGlobuleHashSignature() GlobuleHashSignature {
	return p3p.globuleHashSignature
}
