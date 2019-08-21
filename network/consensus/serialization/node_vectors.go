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
	"fmt"
	"io"
	"math"

	"github.com/pkg/errors"

	"github.com/insolar/insolar/longbits"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/member"
)

const (
	compressedBitIndex = 6
	hasHiLenBitIndex   = 7

	loByteLengthBitSize = 6
	loByteLengthMask    = 1<<loByteLengthBitSize - 1 // 0b00111111
	loByteLengthMax     = loByteLengthMask

	hiByteLengthBitSize = 7
	hiByteLengthMask    = 1<<hiByteLengthBitSize - 1 // 0b01111111
	hiByteLengthMax     = hiByteLengthMask
	hiByteLengthShift   = loByteLengthBitSize - 1

	byteLengthMax = (hiByteLengthMax << loByteLengthBitSize) | loByteLengthMax
)

type NodeVectors struct {
	// ByteSize=133..599
	/*
		GlobulaNodeBitset is a 5-state bitset, each node has a state at the same index as it was given in the current rank.
		Node have following states:
		0 - z-value (same as missing value) Trusted node
		1 - Doubted node
		2 -
		3 - Fraud node
		4 - Missing node
	*/
	StateVectorMask        NodeAppearanceBitset // ByteSize=1..335
	MainStateVector        GlobulaStateVector   // ByteSize=132
	AdditionalStateVectors []GlobulaStateVector `insolar-transport:"count=PacketFlags[1:2]"` // ByteSize=count * 132
}

func (nv NodeVectors) String() string {
	return fmt.Sprintf(
		"<bitset=%v trusted=%s doubted=%s>",
		nv.StateVectorMask.Bytes,
		nv.MainStateVector,
		nv.AdditionalStateVectors,
	)
}

func (nv *NodeVectors) SerializeTo(ctx SerializeContext, writer io.Writer) error {
	if err := nv.StateVectorMask.SerializeTo(ctx, writer); err != nil {
		return errors.Wrap(err, "failed to serialize StateVectorMask")
	}

	if err := nv.MainStateVector.SerializeTo(ctx, writer); err != nil {
		return errors.Wrap(err, "failed to serialize MainStateVector")
	}

	for i := 0; i < int(ctx.GetFlagRangeInt(1, 2)); i++ {
		if err := nv.AdditionalStateVectors[i].SerializeTo(ctx, writer); err != nil {
			return errors.Wrapf(err, "failed to serialize AdditionalStateVectors[%d]", i)
		}
	}

	return nil
}

func (nv *NodeVectors) DeserializeFrom(ctx DeserializeContext, reader io.Reader) error {
	if err := nv.StateVectorMask.DeserializeFrom(ctx, reader); err != nil {
		return errors.Wrap(err, "failed to deserialize StateVectorMask")
	}

	if err := nv.MainStateVector.DeserializeFrom(ctx, reader); err != nil {
		return errors.Wrap(err, "failed to deserialize MainStateVector")
	}

	length := ctx.GetFlagRangeInt(1, 2)
	if length > 0 {
		nv.AdditionalStateVectors = make([]GlobulaStateVector, length)
		for i := 0; i < int(length); i++ {
			if err := nv.AdditionalStateVectors[i].DeserializeFrom(ctx, reader); err != nil {
				return errors.Wrapf(err, "failed to deserialize AdditionalStateVectors[%d]", i)
			}
		}
	}

	return nil
}

type NodeAppearanceBitset struct {
	// ByteSize=1..335
	FlagsAndLoLength uint8 // [00-05] LoByteLength, [06] Compressed, [07] HasHiLength (to be compatible with Protobuf VarInt)
	HiLength         uint8 // [00-06] HiByteLength, [07] MUST = 0 (to be compatible with Protobuf VarInt)
	Bytes            []byte
}

func (nab *NodeAppearanceBitset) SetBitset(bitset member.StateBitset) {
	length := bitset.Len()
	if length > math.MaxUint16 {
		panic("invalid length")
	}

	nab.setLength(uint16(length))
	nab.setIsCompressed(false)
	nab.Bytes = make([]byte, length)

	// TODO: 1 entry == 1 byte. im too lazy
	for i, entry := range bitset {
		nab.Bytes[i] = byte(entry)
	}
}

func (nab *NodeAppearanceBitset) GetBitset() member.StateBitset {
	length := nab.getLength()
	if nab.isCompressed() {
		panic("not implemented")
	}

	bitset := make([]member.BitsetEntry, length)
	for i, b := range nab.Bytes {
		bitset[i] = member.BitsetEntry(b)
	}

	return bitset
}

func (nab *NodeAppearanceBitset) isCompressed() bool {
	return hasBit(uint(nab.FlagsAndLoLength), compressedBitIndex)
}

func (nab *NodeAppearanceBitset) setIsCompressed(compressed bool) {
	nab.FlagsAndLoLength = uint8(toggleBit(uint(nab.FlagsAndLoLength), compressedBitIndex, compressed))
}

func (nab *NodeAppearanceBitset) hasHiLength() bool {
	return hasBit(uint(nab.FlagsAndLoLength), hasHiLenBitIndex)
}

func (nab *NodeAppearanceBitset) setHasHiLength(has bool) {
	nab.FlagsAndLoLength = uint8(toggleBit(uint(nab.FlagsAndLoLength), hasHiLenBitIndex, has))
}

func (nab *NodeAppearanceBitset) getLoLength() uint8 {
	return nab.FlagsAndLoLength & loByteLengthMask
}

func (nab *NodeAppearanceBitset) clearLoLength() {
	nab.FlagsAndLoLength ^= nab.FlagsAndLoLength & loByteLengthMask
}

func (nab *NodeAppearanceBitset) clearHiLength() {
	nab.HiLength ^= nab.HiLength & hiByteLengthMask
}

func (nab *NodeAppearanceBitset) getHiLength() uint8 {
	return nab.HiLength & hiByteLengthMask
}

func (nab *NodeAppearanceBitset) setLoLength(length uint8) {
	if length > loByteLengthMax {
		panic("invalid length")
	}

	nab.FlagsAndLoLength |= length
}

func (nab *NodeAppearanceBitset) setHiLength(length uint8) {
	if length > hiByteLengthMax {
		panic("invalid length")
	}

	nab.HiLength |= length
}

func (nab *NodeAppearanceBitset) getLength() uint16 {
	length := uint16(nab.getLoLength())
	if nab.hasHiLength() {
		return (uint16(nab.getHiLength()) << hiByteLengthShift) | length
	}

	return length
}

func (nab *NodeAppearanceBitset) setLength(length uint16) {
	if length > byteLengthMax {
		panic("invalid length")
	}

	nab.setHasHiLength(length > loByteLengthMax)
	nab.clearHiLength()
	nab.clearLoLength()

	if length > loByteLengthMax {
		nab.setLoLength(uint8(length & loByteLengthMask))
		nab.setHiLength(uint8(length >> hiByteLengthShift))
	} else {
		nab.setLoLength(uint8(length))
	}
}

func (nab *NodeAppearanceBitset) SerializeTo(ctx SerializeContext, writer io.Writer) error {
	if err := write(writer, nab.FlagsAndLoLength); err != nil {
		return errors.Wrap(err, "failed to serialize FlagsAndLoLength")
	}

	if nab.hasHiLength() {
		if err := write(writer, nab.HiLength); err != nil {
			return errors.Wrap(err, "failed to serialize HiLength")
		}
	}

	if nab.getLength() > 0 {
		if err := write(writer, nab.Bytes); err != nil {
			return errors.Wrap(err, "failed to serialize Bytes")
		}
	}

	return nil
}

func (nab *NodeAppearanceBitset) DeserializeFrom(ctx DeserializeContext, reader io.Reader) error {
	if err := read(reader, &nab.FlagsAndLoLength); err != nil {
		return errors.Wrap(err, "failed to deserialize FlagsAndLoLength")
	}

	if nab.hasHiLength() {
		if err := read(reader, &nab.HiLength); err != nil {
			return errors.Wrap(err, "failed to serialize HiLength")
		}
	}

	length := nab.getLength()
	if length > 0 {
		nab.Bytes = make([]byte, length)
		if err := read(reader, &nab.Bytes); err != nil {
			return errors.Wrapf(err, "failed to serialize Bytes")
		}
	}
	return nil
}

type GlobulaStateVector struct {
	// ByteSize=132
	ExpectedRank           member.Rank      // ByteSize=4
	VectorHash             longbits.Bits512 // ByteSize=64
	SignedGlobulaStateHash longbits.Bits512 // ByteSize=64
}

func (gsv GlobulaStateVector) String() string {
	return fmt.Sprintf("<rank=%s gsh=%s §gsh=%s>", gsv.ExpectedRank, gsv.VectorHash, gsv.SignedGlobulaStateHash)
}

func (gsv *GlobulaStateVector) SerializeTo(_ SerializeContext, writer io.Writer) error {
	return write(writer, gsv)
}

func (gsv *GlobulaStateVector) DeserializeFrom(_ DeserializeContext, reader io.Reader) error {
	return read(reader, gsv)
}

// const (
// 	stateBitSize = 3
// 	bitsInByte   = 8
// )
//
// func bitsetByteSize(entryLen int) uint16 {
// 	return uint16(math.Ceil(float64(entryLen*stateBitSize) / bitsInByte))
// }
