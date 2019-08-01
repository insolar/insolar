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

package longbits

import (
	"io"
	"math/bits"
	"strings"
)

//type BinaryString interface {
//	FoldableReader
//
//	BitLen() int
//	BitMasked(index int) (value byte, mask uint8)
//	BitBool(index int) bool
//	BitByte(index int) byte
//
//	SearchBit(startAt int, bit bool) int
//
//	Len() int
//	Byte(index int) byte
//}

const EmptyByteString = ByteString("")

func NewByteStringOf(s string) ByteString {
	return ByteString(s)
}

func NewByteString(v []byte) ByteString {
	return ByteString(v)
}

func NewZeroByteString(len int) ByteString {
	return ByteString(make([]byte, len))
}

func NewFillByteString(len int, fill byte) ByteString {
	if fill == 0 {
		return NewZeroByteString(len)
	}
	return ByteString(strings.Repeat(string([]byte{fill}), len))
}

var _ FoldableReader = EmptyByteString.AsReader()

type ByteString string

// TODO check behavior with nil/zero strings

func (v ByteString) AsReader() FoldableReader {
	return &v
}

func (v *ByteString) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write([]byte(*v))
	return int64(n), err
}

func (v *ByteString) Read(p []byte) (n int, err error) {
	return copy(p, *v), nil
}

func (v *ByteString) AsBytes() []byte {
	return ([]byte)(*v)
}

func (v *ByteString) AsByteString() ByteString {
	return *v
}

func (v *ByteString) FixedByteSize() int {
	return len(*v)
}

func (v *ByteString) FoldToUint64() uint64 {
	folded := v.FoldToBits64()
	return folded.FoldToUint64()
}

func (v ByteString) BitMasked(index int) (value byte, mask uint8) {
	bytePos, bitPos := v.BitPos(index)
	mask = 1 << bitPos
	return v[bytePos] & mask, mask
}

func (v ByteString) BitBool(index int) bool {
	if b, _ := v.BitMasked(index); b != 0 {
		return true
	}
	return false
}

func (v ByteString) BitByte(index int) byte {
	if b, _ := v.BitMasked(index); b != 0 {
		return 1
	}
	return 0
}

func (v ByteString) Byte(index int) byte {
	return v[index]
}

func (v ByteString) BitPos(index int) (bytePos int, bitPos uint8) {
	bytePos, bitPos = BitPos(index)
	if bytePos >= len(v) {
		panic("out of bounds")
	}
	return bytePos, bitPos
}

func (v ByteString) BitLen() int {
	return len(v) << 3
}

func (v ByteString) SearchBit(startAt int, bit bool) int {
	if startAt < 0 {
		panic("illegal value")
	}
	if startAt>>3 >= len(v) {
		return -1
	}

	if bit {
		return v.searchBit1(startAt)
	}
	return v.searchBit0(startAt)
}

func (v ByteString) searchBit1(startAt int) int {
	bytePos := startAt >> 3
	bitPos := byte(startAt & 0x7)

	b := v[bytePos] >> bitPos
	if b != 0 {
		return bytePos<<3 + int(bitPos) + bits.TrailingZeros8(b)
	}

	for bytePos++; bytePos < len(v); bytePos++ {
		b := v[bytePos]
		if b != 0 {
			return bytePos<<3 + bits.TrailingZeros8(b)
		}
	}
	return -1

}

func (v ByteString) searchBit0(startAt int) int {
	bytePos := startAt >> 3
	bitPos := byte(startAt & 0x7)

	b := (^v[bytePos]) >> bitPos
	if b != 0 {
		return bytePos<<3 + int(bitPos) + bits.TrailingZeros8(b)
	}

	for bytePos++; bytePos < len(v); bytePos++ {
		b := v[bytePos]
		if b != 0xFF {
			return bytePos<<3 + bits.TrailingZeros8(^b)
		}
	}
	return -1
}

func BitPos(index int) (bytePos int, bitPos uint8) {
	if index < 0 {
		panic("illegal value")
	}
	if index == 0 {
		return 0, 0
	}
	return index >> 3, uint8(index & 0x07)
}

func (v ByteString) FoldToBits64() Bits64 {
	var folded Bits64
	if len(v) == 0 {
		return folded
	}

	alignedLen := len(v) & (len(folded) - 1)
	copy(folded[alignedLen:], v)

	for i := 0; i < alignedLen; i += len(folded) {
		folded[0] ^= v[i+0]
		folded[1] ^= v[i+1]
		folded[2] ^= v[i+2]
		folded[3] ^= v[i+3]
		folded[4] ^= v[i+4]
		folded[5] ^= v[i+5]
		folded[6] ^= v[i+6]
		folded[7] ^= v[i+7]
	}
	return folded
}

func (v ByteString) String() string {
	return bitsToStringDefault(&v)
}
