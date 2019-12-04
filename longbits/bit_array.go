//
//    Copyright 2019 Insolar Technologies
//
//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.
//

package longbits

import "math/bits"

type IndexedBits interface {
	BitMasked(index int) (value byte, mask uint8)
	BitBool(index int) bool
	BitValue(index int) byte
	Byte(index int) byte
	BitLen() int
	SearchBit(startAt int, bit bool) int
}

type BitSlice []byte

func BitPos(index int) (bytePos int, bitPos uint8) {
	if index < 0 {
		panic("illegal value")
	}
	return index >> 3, uint8(index & 0x07)
}

func (v BitSlice) BitMasked(index int) (value byte, mask uint8) {
	bytePos, bitPos := v.BitPos(index)
	mask = 1 << bitPos
	return v[bytePos] & mask, mask
}

func (v BitSlice) BitBool(index int) bool {
	if b, _ := v.BitMasked(index); b != 0 {
		return true
	}
	return false
}

func (v BitSlice) BitValue(index int) byte {
	if b, _ := v.BitMasked(index); b != 0 {
		return 1
	}
	return 0
}

func (v BitSlice) Byte(index int) byte {
	return v[index]
}

func (v BitSlice) BitPos(index int) (bytePos int, bitPos uint8) {
	bytePos, bitPos = BitPos(index)
	if bytePos >= len(v) {
		panic("out of bounds")
	}
	return bytePos, bitPos
}

func (v BitSlice) BitLen() int {
	return len(v) << 3
}

func (v BitSlice) SearchBit(startAt int, bit bool) int {
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

func (v BitSlice) searchBit1(startAt int) int {
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

func (v BitSlice) searchBit0(startAt int) int {
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
