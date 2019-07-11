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
	"context"
	"errors"
	"io"

	"github.com/insolar/insolar/network/consensus/common"
	"github.com/insolar/insolar/network/consensus/gcpv2/packets"
)

const (
	packetTypeBitSize = 4
	packetTypeMask    = 1<<packetTypeBitSize - 1 // 0b00001111
	packetTypeMax     = packetTypeMask

	protocolTypeBitSize = 4
	protocolTypeShift   = protocolTypeBitSize
	protocolTypeMax     = 1<<protocolTypeBitSize - 1

	payloadLengthBitSize = 14
	payloadLengthMask    = 1<<payloadLengthBitSize - 1 // 0b0011111111111111
	payloadLengthMax     = payloadLengthMask

	pulseNumberBitSize = 30
	pulseNumberMask    = 1<<pulseNumberBitSize - 1 // 0b00111111111111111111111111111111
	pulseNumberMax     = pulseNumberMask
)

type Flag uint8

const (
	flagIsRelayRestricted = Flag(0)
	flagIsBodyEncrypted   = Flag(1)

	FlagHasPulsePacket = Flag(0)
)

const (
	reservedFlagSize = 2
	maxFlagIndex     = 5
)

type ProtocolType uint8

const (
	ProtocolTypePulsar           = ProtocolType(0)
	ProtocolTypeGlobulaConsensus = ProtocolType(1)
)

func (pt ProtocolType) NewBody() PacketBody {
	switch pt {
	case ProtocolTypePulsar:
		return &PulsarPacketBody{}
	case ProtocolTypeGlobulaConsensus:
		return &GlobulaConsensusPacketBody{}
	}

	return nil
}

var ErrNilBody = errors.New("body is nil")

/*
	ByteSize=16
*/
type Header struct {
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

func (h *Header) SerializeTo(_ SerializeContext, writer io.Writer) error {
	return write(writer, h)
}

func (h *Header) DeserializeFrom(_ DeserializeContext, reader io.Reader) error {
	return read(reader, h)
}

func (h *Header) GetSourceID() common.ShortNodeID {
	return common.ShortNodeID(h.SourceID)
}

func (h *Header) GetPacketType() packets.PacketType {
	return packets.PacketType(h.ProtocolAndPacketType) & packetTypeMask
}

func (h *Header) setPacketType(packetType packets.PacketType) {
	if packetType > packetTypeMax {
		panic("invalid packet type")
	}

	h.ProtocolAndPacketType |= uint8(packetType)
}

func (h *Header) GetProtocolType() ProtocolType {
	return ProtocolType(h.ProtocolAndPacketType) >> protocolTypeShift
}

func (h *Header) setProtocolType(protocolType ProtocolType) {
	if protocolType > protocolTypeMax {
		panic("invalid protocol type")
	}

	h.ProtocolAndPacketType |= uint8(protocolType << protocolTypeShift)
}

func (h *Header) getPayloadLength() uint16 {
	return h.HeaderAndPayloadLength & payloadLengthMask
}

func (h *Header) setPayloadLength(payloadLength uint16) {
	if payloadLength > payloadLengthMax {
		panic("invalid payload length")
	}

	h.HeaderAndPayloadLength |= payloadLength
}

func (h *Header) HasFlag(f Flag) bool {
	if f > maxFlagIndex {
		panic("invalid flag index")
	}

	return h.hasFlag(f + reservedFlagSize)
}

func (h *Header) ClearFlag(f Flag) {
	if f > maxFlagIndex {
		panic("invalid flag index")
	}

	h.clearFlag(f + reservedFlagSize)
}

func (h *Header) GetFlagRangeInt(from, to uint8) uint8 {
	if from > to {
		panic("invalid from range")
	}

	if to > maxFlagIndex {
		panic("invalid to range")
	}

	return h.getFlagRangeInt(from+reservedFlagSize, to+reservedFlagSize)
}

func (h *Header) SetFlag(f Flag) {
	if f > maxFlagIndex {
		panic("invalid flag index")
	}

	h.setFlag(f + reservedFlagSize)
}

func (h *Header) IsRelayRestricted() bool {
	return h.hasFlag(flagIsRelayRestricted)
}

func (h *Header) setIsRelayRestricted(restricted bool) {
	h.toggleFlag(flagIsRelayRestricted, restricted)
}

func (h *Header) IsBodyEncrypted() bool {
	return h.hasFlag(flagIsBodyEncrypted)
}

func (h *Header) setIsBodyEncrypted(encrypted bool) {
	h.toggleFlag(flagIsBodyEncrypted, encrypted)
}

func (h *Header) hasFlag(f Flag) bool {
	return hasBit(uint(h.PacketFlags), uint(f))
}

func (h *Header) toggleFlag(f Flag, val bool) {
	h.PacketFlags = uint8(toggleBit(uint(h.PacketFlags), uint(f), val))
}

func (h *Header) clearFlag(f Flag) {
	h.PacketFlags = uint8(clearBit(uint(h.PacketFlags), uint(f)))
}

func (h *Header) setFlag(f Flag) {
	h.PacketFlags = uint8(setBit(uint(h.PacketFlags), uint(f)))
}

func (h *Header) getFlagRangeInt(from, to uint8) uint8 {
	return uint8(uintFromBits(uint(h.PacketFlags), uint(from), uint(to)))
}

type Packet struct {
	Header      Header             `insolar-transport:"Protocol=0x01;Packet=0-4"` // ByteSize=16
	PulseNumber common.PulseNumber `insolar-transport:"[30-31]=0"`                // [30-31] MUST ==0, ByteSize=4

	EncryptableBody PacketBody
	EncryptionData  []byte

	PacketSignature common.Bits512 `insolar-transport:"generate=signature"` // ByteSize=64
}

func (p *Packet) setSignature(signature common.SignatureHolder) {
	copy(p.PacketSignature[:], signature.AsBytes())
}

func (p *Packet) setPayloadLength(payloadLength uint16) {
	p.Header.setPayloadLength(payloadLength)
}

func (p *Packet) getPulseNumber() common.PulseNumber {
	return p.PulseNumber & pulseNumberMask
}

func (p *Packet) setPulseNumber(pulseNumber common.PulseNumber) {
	if pulseNumber > pulseNumberMax {
		panic("invalid pulse number")
	}

	p.PulseNumber |= pulseNumber
}

func (p *Packet) SerializeTo(ctx context.Context, writer io.Writer, digester common.DataDigester, signer common.DigestSigner) (int64, error) {
	if p.EncryptableBody == nil {
		return 0, ErrMalformedPacketBody(ErrNilBody)
	}

	w := newTrackableWriter(writer)
	pctx := newPacketContext(ctx, &p.Header)
	sctx := newSerializeContext(pctx, w, digester, signer, p)

	if err := write(sctx, &p.PulseNumber); err != nil {
		return 0, ErrMalformedPulseNumber(err)
	}

	if err := p.EncryptableBody.SerializeTo(sctx, sctx); err != nil {
		return 0, ErrMalformedPacketBody(err)
	}

	return sctx.Finalize()
}

func (p *Packet) DeserializeFrom(ctx context.Context, reader io.Reader) (int64, error) {
	r := newTrackableReader(reader)

	if err := p.Header.DeserializeFrom(nil, r); err != nil {
		return r.totalRead, ErrMalformedHeader(err)
	}

	pctx := newPacketContext(ctx, &p.Header)
	dctx := newDeserializeContext(pctx, r, &p.Header)

	if err := read(dctx, &p.PulseNumber); err != nil {
		return r.totalRead, ErrMalformedPulseNumber(err)
	}

	p.EncryptableBody = pctx.GetProtocolType().NewBody()
	if p.EncryptableBody == nil {
		return 0, ErrMalformedPacketBody(ErrNilBody)
	}

	if err := p.EncryptableBody.DeserializeFrom(dctx, dctx); err != nil {
		return r.totalRead, ErrMalformedPacketBody(err)
	}

	if err := read(dctx, &p.PacketSignature); err != nil {
		return r.totalRead, ErrMalformedPacketSignature(err)
	}

	return dctx.Finalize()
}
