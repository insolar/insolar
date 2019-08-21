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
