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

type BitBuilderOrder byte

const (
	// Least significant bit is first - first AppendBit() appends the least significant bit
	LSB BitBuilderOrder = 0
	// Most significant bit is first - first AppendBit() appends the most significant bit
	MSB BitBuilderOrder = 1

	initLSB = 0x01
	initMSB = 0x80
)

func NewBitBuilder(direction BitBuilderOrder, expectedLen int) BitBuilder {
	if expectedLen > 0 {
		return AppendBitBuilder(make([]byte, 0, expectedLen), direction)
	}
	return AppendBitBuilder(nil, direction)
}

func AppendBitBuilder(appendTo []byte, direction BitBuilderOrder) BitBuilder {
	switch direction {
	case LSB:
		return BitBuilder{accInit: initLSB, accBit: initLSB, bytes: appendTo}
	case MSB:
		return BitBuilder{accInit: initMSB, accBit: initMSB, bytes: appendTo}
	default:
		panic("illegal value")
	}
}

// supports to be created as BitBuilder{} - it equals NewBitBuilder(LSB, 0)
type BitBuilder struct {
	bytes       []byte
	accumulator byte
	accInit     byte
	accBit      byte
}

func (p *BitBuilder) IsZero() bool {
	return p.accInit == 0
}

func (p *BitBuilder) Len() int {
	return len(p.bytes)<<3 + int(p.AlignOffset())
}

func (p *BitBuilder) ensure() {
	if p.accInit == 0 {
		if len(p.bytes) != 0 {
			panic("illegal state")
		}
		p.accInit = initLSB
		p.accBit = initLSB
	} else if p.accBit == 0 {
		panic("illegal state")
	}
}

func (p *BitBuilder) AppendAlignedByte(b byte) {
	p.ensure()
	if p.accBit != p.accInit {
		panic("illegal state")
	}
	if p._rightShift() {
		b = bits.Reverse8(b)
	}
	p.bytes = append(p.bytes, b)
}

func shiftLeft(b, n byte) byte {
	return b << n
}

func shiftRight(b, n byte) byte {
	return b >> n
}

func (p *BitBuilder) _align(rightShift bool) uint8 {
	switch {
	case p.accBit == p.accInit:
		return 0
	case rightShift:
		return uint8(bits.LeadingZeros8(p.accBit))
	default:
		return uint8(bits.TrailingZeros8(p.accBit))
	}
}

func (p *BitBuilder) align() (rightShift bool, ofs uint8) {
	switch rightShift := p._rightShift(); {
	case p.accBit == p.accInit:
		return rightShift, 0
	case rightShift:
		return true, uint8(bits.LeadingZeros8(p.accBit))
	default:
		return false, uint8(bits.TrailingZeros8(p.accBit))
	}
}

func (p *BitBuilder) _rightShift() bool {
	switch {
	case p.accInit == initLSB:
		return false
	case p.accInit == initMSB:
		return true
	default:
		panic("illegal state")
	}
}

func shifters(rightShift bool) (normFn, revFn func(byte, byte) byte) {
	if rightShift {
		return shiftRight, shiftLeft
	}
	return shiftLeft, shiftRight
}

func (p *BitBuilder) AlignOffset() uint8 {
	_, ofs := p.align()
	return ofs
}

func (p *BitBuilder) CompletedByteCount() int {
	return len(p.bytes)
}

func (p *BitBuilder) PadWithBit(bit int) {
	p.PadWith(bit != 0)
}

func (p *BitBuilder) PadWith(bit bool) {
	if bit {
		p.appendN1(-1)
	}
	p.appendN0(-1)
}

func (p *BitBuilder) AppendBit(bit int) {
	p.Append(bit != 0)
}

func (p *BitBuilder) Append(bit bool) {
	p.ensure()

	if bit {
		p.accumulator |= p.accBit
	}

	if p._rightShift() {
		p.accBit >>= 1
	} else {
		p.accBit <<= 1
	}

	if p.accBit == 0 {
		p.bytes = append(p.bytes, p.accumulator)
		p.accumulator = 0
		p.accBit = p.accInit
	}
}

func (p *BitBuilder) AppendSubByte(value byte, bitLen uint8) {
	if bitLen >= 8 {
		if bitLen != 8 {
			panic("illegal value")
		}
		p.AppendByte(value)
		return
	}
	switch bitLen {
	case 0:
		return
	case 1:
		p.Append(value&1 != 0)
		return
	}

	p.ensure()
	rightShift, usedCount := p.align()
	normFn, revFn := shifters(rightShift)
	if rightShift {
		value = bits.Reverse8(value)
	}

	value &= revFn(0xFF, 8-bitLen)

	remainCount := 8 - usedCount
	switch {
	case usedCount == 0:
		p.accBit = normFn(p.accBit, bitLen)
		p.accumulator = value
		return
	case remainCount > bitLen:
		p.accBit = normFn(p.accBit, bitLen)
		p.accumulator |= normFn(value, usedCount)
		return
	default:
		p.accumulator |= normFn(value, usedCount)
		bitLen -= remainCount
	}

	p.bytes = append(p.bytes, p.accumulator)
	p.accBit = p.accInit
	if bitLen == 0 {
		p.accumulator = 0
		return
	}
	p.accBit = normFn(p.accBit, bitLen)
	p.accumulator = revFn(value, remainCount)
}

func (p *BitBuilder) AppendNBit(bitCount int, bit int) {
	p.AppendN(bitCount, bit != 0)
}

func (p *BitBuilder) AppendN(bitCount int, bit bool) {
	p.ensure()
	switch {
	case bitCount == 0:
	case bitCount == 1:
		p.Append(bit)
	case bitCount < 0:
		panic("invalid bitCount value")
	case bit:
		p.appendN1(bitCount)
	default:
		p.appendN0(bitCount)
	}
}

func (p *BitBuilder) appendN0(bitCount int) {
	p.ensure()

	rightShift, usedCount := p.align()
	normFn, _ := shifters(rightShift)

	if usedCount == 0 {
		if bitCount < 0 {
			return
		}
	} else {
		switch {
		case bitCount < 0:
			bitCount = 0
		default:
			alignCount := 8 - int(usedCount)
			if alignCount > bitCount {
				p.accBit = normFn(p.accBit, uint8(bitCount))
				return
			}
			bitCount -= alignCount
		}
		p.bytes = append(p.bytes, byte(p.accumulator))
		p.accumulator = 0
		p.accBit = p.accInit
		if bitCount == 0 {
			return
		}
	}

	if alignCount := uint8(bitCount) & 0x7; alignCount > 0 {
		p.accBit = normFn(p.accBit, alignCount)
	}
	if byteCount := bitCount >> 3; byteCount > 0 {
		p.bytes = append(p.bytes, make([]byte, byteCount)...)
	}
}

func (p *BitBuilder) appendN1(bitCount int) {
	p.ensure()

	rightShift, usedCount := p.align()
	normFn, revFn := shifters(rightShift)

	if usedCount == 0 {
		if bitCount < 0 {
			return
		}
	} else {
		switch {
		case bitCount < 0:
			bitCount = 0
		default:
			alignCount := 8 - int(usedCount)
			if alignCount > bitCount {
				filler := revFn(0xFF, uint8(alignCount-bitCount)) // make some zeros
				p.accumulator |= normFn(filler, usedCount)
				p.accBit = normFn(p.accBit, uint8(bitCount))
				return
			}
			bitCount -= alignCount
		}
		p.accumulator |= normFn(0xFF, usedCount)
		p.bytes = append(p.bytes, byte(p.accumulator))
		p.accumulator = 0
		p.accBit = p.accInit
		if bitCount == 0 {
			return
		}
	}

	if alignCount := uint8(bitCount) & 0x7; alignCount > 0 {
		p.accBit = normFn(p.accBit, alignCount)
		p.accumulator = revFn(0xFF, 8-alignCount)
	}

	if byteCount := bitCount >> 3; byteCount > 0 {
		lowBound := len(p.bytes)
		p.bytes = append(p.bytes, make([]byte, byteCount)...)
		for i := len(p.bytes) - 1; i >= lowBound; i-- {
			p.bytes[i] = 0xFF
		}
	}
}

func (p *BitBuilder) AppendByte(b byte) {
	p.ensure()

	rightShift, usedCount := p.align()
	normFn, revFn := shifters(rightShift)

	if rightShift {
		b = bits.Reverse8(b)
	}

	if usedCount == 0 {
		p.bytes = append(p.bytes, b)
		return
	}
	nextByte := p.accumulator | normFn(b, usedCount)
	p.bytes = append(p.bytes, nextByte)

	p.accumulator = revFn(b, 8-usedCount)
}

func (p *BitBuilder) dump() []byte {
	_, usedCount := p.align()

	bytes := append(make([]byte, 0, cap(p.bytes)), p.bytes...)
	if usedCount > 0 {
		bytes = append(bytes, byte(p.accumulator))
	}
	return bytes
}

func (p *BitBuilder) Done() ([]byte, int) {
	_, usedCount := p.align()

	bytes := p.bytes
	p.bytes = nil
	if usedCount > 0 {
		bytes = append(bytes, byte(p.accumulator))
		p.accumulator = 0
		p.accBit = p.accInit
		return bytes, (len(p.bytes)-1)<<3 + int(usedCount)
	}
	return bytes, len(p.bytes) << 3
}

func (p *BitBuilder) DoneToBytes() []byte {
	b, _ := p.Done()
	return b
}

func (p *BitBuilder) DoneToByteString() (ByteString, int) {
	b, l := p.Done()
	return CopyBytes(b), l
}

func (p *BitBuilder) Copy() *BitBuilder {
	c := *p
	if p.bytes != nil {
		c.bytes = append(make([]byte, 0, cap(p.bytes)), p.bytes...)
	}
	return &c
}
