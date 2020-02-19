// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package longbits

import (
	"io"
	"math/bits"
	"strings"
)

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
