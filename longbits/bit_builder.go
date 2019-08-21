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
	"math/bits"
)

func NewBitBuilder(expectedLen int) BitBuilder {
	if expectedLen == 0 {
		return BitBuilder{}
	}
	return BitBuilder{bytes: make([]byte, 0, expectedLen)}
}

type BitBuilder struct {
	bytes       []byte
	accumulator uint16
}

func (p BitBuilder) IsZero() bool {
	return p.accumulator == 0 && p.bytes == nil
}

func (p BitBuilder) Len() int {
	return (1+len(p.bytes))<<3 - bits.LeadingZeros8(uint8(p.accumulator>>8))
}

const accInit = 0x80

func (p BitBuilder) ensure() {
	if p.accumulator == 0 {
		if len(p.bytes) != 0 {
			panic("illegal state")
		}
		p.accumulator = accInit
	}
}

func (p BitBuilder) Append(bit bool) BitBuilder {
	p.ensure()

	p.accumulator <<= 1
	if bit {
		p.accumulator |= 1
	}
	if p.accumulator >= accInit<<8 {
		p.bytes = append(p.bytes, byte(p.accumulator))
		p.accumulator = accInit
	}

	return p
}

func (p BitBuilder) AppendN(bitCount int, bit bool) BitBuilder {
	p.ensure()

	if bitCount == 0 {
		return p
	}
	if bit {
		return p.appendN1(bitCount)
	}
	return p.appendN0(bitCount)
}

func (p BitBuilder) appendN0(bitCount int) BitBuilder {
	p.ensure()

	if p.accumulator != accInit {
		alignCount := bits.LeadingZeros8(uint8(p.accumulator >> 8))
		if alignCount < bitCount {
			p.accumulator <<= uint8(bitCount)
			return p
		}

		bitCount -= alignCount
		p.accumulator <<= uint8(alignCount)
		p.bytes = append(p.bytes, byte(p.accumulator))
		p.accumulator = accInit
	}

	if bitCount == 0 {
		return p
	}

	alignCount := uint8(bitCount) & 0x7
	bitCount >>= 3

	if bitCount > 0 {
		p.bytes = append(p.bytes, make([]byte, bitCount)...)
	}

	p.accumulator <<= alignCount
	return p
}

func (p BitBuilder) appendN1(bitCount int) BitBuilder {
	p.ensure()

	if p.accumulator != accInit {
		alignCount := bits.LeadingZeros8(uint8(p.accumulator >> 8))
		if alignCount < bitCount {
			p.accumulator <<= uint8(bitCount)
			p.accumulator |= 0xFF >> uint8(8-bitCount)
			return p
		}

		bitCount -= alignCount
		p.accumulator <<= uint8(alignCount)
		p.accumulator |= 0xFF >> uint8(8-alignCount)

		p.bytes = append(p.bytes, byte(p.accumulator))
		p.accumulator = accInit
	}

	if bitCount == 0 {
		return p
	}

	alignCount := uint8(bitCount) & 0x7
	bitCount >>= 3

	if bitCount > 0 {
		i := len(p.bytes)
		p.bytes = append(p.bytes, make([]byte, bitCount)...)
		for ; i < len(p.bytes); i++ {
			p.bytes[i] = 0xFF
		}
	}

	p.accumulator <<= alignCount
	p.accumulator |= 0xFF >> (8 - alignCount)
	return p
}

func (p BitBuilder) Done() ([]byte, int) {
	if p.accumulator <= accInit {
		return p.bytes, len(p.bytes) << 3
	}

	p.bytes = append(p.bytes, byte(p.accumulator))
	return p.bytes, len(p.bytes)<<3 - bits.LeadingZeros8(uint8(p.accumulator>>8))
}

func (p BitBuilder) DoneAndCopy() ([]byte, int) {
	b, l := p.Done()
	if len(b) == 0 {
		return nil, l
	}
	return append(make([]byte, 0, len(b)), b...), l
}

func (p BitBuilder) DoneToByteString() (ByteString, int) {
	b, l := p.Done()
	return NewByteString(b), l
}

func (p BitBuilder) Copy() BitBuilder {
	if p.bytes == nil {
		return BitBuilder{accumulator: p.accumulator}
	}
	return BitBuilder{accumulator: p.accumulator, bytes: append(make([]byte, 0, cap(p.bytes)), p.bytes...)}
}
