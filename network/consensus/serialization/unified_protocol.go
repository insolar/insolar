//
// Modified BSD 3-Clause Clear License
//
// Copyright (c) 2019 Insolar Technologies GmbH
//
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without modification,
// are permitted (subject to the limitations in the disclaimer below) provided that
// the following conditions are met:
//  * Redistributions of source code must retain the above copyright notice, this list
//    of conditions and the following disclaimer.
//  * Redistributions in binary form must reproduce the above copyright notice, this list
//    of conditions and the following disclaimer in the documentation and/or other materials
//    provided with the distribution.
//  * Neither the name of Insolar Technologies GmbH nor the names of its contributors
//    may be used to endorse or promote products derived from this software without
//    specific prior written permission.
//
// NO EXPRESS OR IMPLIED LICENSES TO ANY PARTY'S PATENT RIGHTS ARE GRANTED
// BY THIS LICENSE. THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS
// AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES,
// INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY
// AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL
// THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT,
// INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING,
// BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS
// OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
// ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
//
// Notwithstanding any other provisions of this license, it is prohibited to:
//    (a) use this software,
//
//    (b) prepare modifications and derivative works of this software,
//
//    (c) distribute this software (including without limitation in source code, binary or
//        object code form), and
//
//    (d) reproduce copies of this software
//
//    for any commercial purposes, and/or
//
//    for the purposes of making available this software to third parties as a service,
//    including, without limitation, any software-as-a-service, platform-as-a-service,
//    infrastructure-as-a-service or other similar online service, irrespective of
//    whether it competes with the products or services of Insolar Technologies GmbH.
//

package serialization

import (
	"io"
	"math/bits"

	"github.com/insolar/insolar/network/consensus/common"
	"github.com/insolar/insolar/network/consensus/gcpv2/packets"
)

const (
	packetTypeMask      = 15 // 0b00001111
	packetTypeBitSize   = 4
	protocolTypeShift   = 4
	protocolTypeBitSize = 4

	payloadLengthMask    = 16383 // 0b0011111111111111
	payloadLengthBitSize = 14
	headerShift          = 14
	headerBitSize        = 2
)

type FlagType uint8

const (
	IsRelayRestricted = FlagType(0)
	IsBodyEncrypted   = FlagType(1)
)

type ProtocolType uint8

const (
	GlobulaConsensus = ProtocolType(1)
)

/*
	ByteSize=16
*/
type UnifiedProtocolPacketHeader struct {
	/*
		Functions of TargetID, SourceID and RelayId depends on ProtocolType
	*/
	ReceiverID uint32 // NB! MUST for Signature calculation must be considered as 0, actual value can be different

	ProtocolAndPacketType  uint8  `insolar-transport:"[0:3]=header:Packet;[4:7]=header:Protocol"` // [00-03]PacketType [04-07]ProtocolType
	PacketFlags            uint8  `insolar-transport:"[0]=IsRelayRestricted;[1]=IsBodyEncrypted;[2:]=flags:PacketFlags"`
	HeaderAndPayloadLength uint16 // [00-13] ByteLength of Payload, [14-15] reserved = 0
	SourceID               uint32 // may differ from actual sender when relay is in use, MUST NOT =0
	TargetID               uint32 // indicates final destination, if =0 then there is no relay allowed by sender and receiver MUST decline a packet if actual sender != source
}

func (p UnifiedProtocolPacketHeader) SerializeTo(writer io.Writer, signer common.DataSigner) error {
	return serializeTo(writer, signer, p)
}

func (p UnifiedProtocolPacketHeader) GetPacketType() packets.PacketType {
	return packets.PacketType(p.ProtocolAndPacketType) & packetTypeMask
}

func (p *UnifiedProtocolPacketHeader) SetPacketType(packetType packets.PacketType) {
	if bits.Len(uint(packetType)) > packetTypeBitSize {
		panic("invalid packet type")
	}

	p.ProtocolAndPacketType |= uint8(packetType)
}

func (p UnifiedProtocolPacketHeader) GetProtocolType() ProtocolType {
	return ProtocolType(p.ProtocolAndPacketType) >> protocolTypeShift
}

func (p *UnifiedProtocolPacketHeader) SetProtocolType(protocolType ProtocolType) {
	if bits.Len(uint(protocolType)) > protocolTypeBitSize {
		panic("invalid protocol type")
	}

	p.ProtocolAndPacketType |= uint8(protocolType << protocolTypeShift)
}

func (p UnifiedProtocolPacketHeader) GetPayloadLength() uint16 {
	return p.HeaderAndPayloadLength & payloadLengthMask
}

func (p UnifiedProtocolPacketHeader) SetPayloadLength(payloadLength uint16) {
	if bits.Len(uint(payloadLength)) > payloadLengthBitSize {
		panic("invalid payload length")
	}

	p.HeaderAndPayloadLength |= payloadLength
}

func (p *UnifiedProtocolPacketHeader) GetHeader() uint16 {
	return p.HeaderAndPayloadLength >> headerShift
}

func (p *UnifiedProtocolPacketHeader) GetFlag(f FlagType) bool {
	if f > 5 {
		panic("invalid flag index")
	}

	return p.getFlag(f + 2)
}

func (p *UnifiedProtocolPacketHeader) SetFlag(f FlagType) {
	if f > 5 {
		panic("invalid flag index")
	}

	p.setFlag(f + 2)
}

func (p *UnifiedProtocolPacketHeader) IsRelayRestricted() bool {
	return p.getFlag(IsRelayRestricted)
}

func (p *UnifiedProtocolPacketHeader) SetIsRelayRestricted() {
	p.setFlag(IsRelayRestricted)
}

func (p *UnifiedProtocolPacketHeader) IsBodyEncrypted() bool {
	return p.getFlag(IsBodyEncrypted)
}

func (p *UnifiedProtocolPacketHeader) SetIsBodyEncrypted() {
	p.setFlag(IsBodyEncrypted)
}

func (p *UnifiedProtocolPacketHeader) getFlag(f FlagType) bool {
	return hasBit(uint(p.PacketFlags), uint(f))
}

func (p *UnifiedProtocolPacketHeader) setFlag(f FlagType) {
	setBit(uint(p.PacketFlags), uint(f))
}
