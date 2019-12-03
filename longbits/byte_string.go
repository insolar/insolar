//
// Copyright 2019 Insolar Technologies GmbH
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package longbits

import (
	"io"
	"math/bits"
	"reflect"
	"unsafe"

	"github.com/insolar/insolar/longbits/bytehash"
)

const EmptyByteString = ByteString("")

func Wrap(s string) ByteString {
	return ByteString(s)
}

func Copy(v []byte) ByteString {
	return ByteString(v)
}

// WARNING! The given array MUST be immutable
func WrapBytes(b []byte) ByteString {
	if len(b) == 0 {
		return EmptyByteString
	}
	return wrap(b)
}

func wrap(b []byte) ByteString {
	pSlice := (*reflect.SliceHeader)(unsafe.Pointer(&b))

	var res ByteString
	pString := (*reflect.StringHeader)(unsafe.Pointer(&res))

	pString.Data = pSlice.Data
	pString.Len = pSlice.Len

	return res
}

func Zero(len int) ByteString {
	return Fill(len, 0)
}

func Fill(len int, fill byte) ByteString {
	if len == 0 {
		return EmptyByteString
	}
	b := make([]byte, len)
	if fill != 0 {
		for i := len - 1; i >= 0; i-- {
			b[i] = fill
		}
	}
	return wrap(b)
}

var _ FoldableReader = EmptyByteString.AsReader()

type ByteString string

// TODO check behavior with nil/zero strings

func (v ByteString) IsEmpty() bool {
	return len(v) == 0
}

func (v ByteString) AsReader() FoldableReader {
	return v
}

func (v ByteString) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write([]byte(v))
	return int64(n), err
}

func (v ByteString) Hash() uint32 {
	return bytehash.HashStr(string(v))
}

func (v ByteString) HashWithSeed(seed uint32) uint32 {
	return bytehash.HashStrWithSeed(string(v), uint(seed))
}

func (v ByteString) Read(b []byte) (n int, err error) {
	return copy(b, v), nil
}

func (v ByteString) ReadAt(b []byte, off int64) (n int, err error) {
	if n, err = VerifyReadAt(b, off, len(v)); err != nil || n == 0 {
		return n, err
	} else {
		return copy(b, v[off:]), nil
	}
}

func (v ByteString) AsBytes() []byte {
	return ([]byte)(v)
}

func (v ByteString) AsByteString() ByteString {
	return v
}

func (v ByteString) FixedByteSize() int {
	return len(v)
}

func (v ByteString) FoldToUint64() uint64 {
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
