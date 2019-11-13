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
)

type bitReader struct {
	accumulator byte
	accInit     byte
	accBit      byte
}

func (p *bitReader) align() (rightShift bool, ofs uint8) {
	if p.accInit == initFirstLow {
		if p.accBit == 0 {
			return false, 0
		}
		return false, uint8(bits.TrailingZeros8(p.accBit))
	}

	if p.accBit == 0 {
		return true, 0
	}
	return true, uint8(bits.LeadingZeros8(p.accBit))
}

func (p *bitReader) AlignOffset() uint8 {
	_, ofs := p.align()
	return ofs
}

func (p *bitReader) ensure(readFn func() (byte, error)) error {
	if p.accBit != 0 {
		return nil
	}
	if p.accInit == 0 {
		p.accInit = initFirstLow
	}
	if b, err := readFn(); err != nil {
		return err
	} else {
		p.accumulator = b
	}
	p.accBit = p.accInit
	return nil
}

func (p *bitReader) readNext(readFn func() (byte, error)) (int, uint8, error) {
	if e := p.ensure(readFn); e != nil {
		return 0, 0, e
	}

	m := p.accBit
	if rightShift, _ := p.align(); rightShift {
		p.accBit >>= 1
	} else {
		p.accBit <<= 1
	}
	return int(p.accumulator & m), m, nil
}

func (p *bitReader) readByte(readFn func() (byte, error)) (byte, error) {
	switch p.accBit {
	case 0:
		if p.accInit == 0 {
			p.accInit = initFirstLow
		}
		return readFn()
	case p.accInit:
		p.accBit = 0
		return p.accumulator, nil
	}

	v := p.accumulator
	if b, err := readFn(); err != nil {
		return 0, err
	} else {
		p.accumulator = b
	}

	w := p.accumulator
	if rightShift, usedBits := p.align(); rightShift {
		v <<= usedBits
		w >>= 8 - usedBits
	} else {
		v >>= usedBits
		w <<= 8 - usedBits
	}
	return v | w, nil
}

func (p *bitReader) readSubByte(bitLen uint8, readFn func() (byte, error)) (uint8, error) {
	switch {
	case bitLen == 0:
		return 0, nil
	case bitLen == 1:
		switch v, _, e := p.readNext(readFn); {
		case e != nil:
			return 0, e
		case v != 0:
			return 1, nil
		default:
			return 0, nil
		}
	case bitLen == 8:
		return p.readByte(readFn)
	case bitLen > 8:
		panic("illegal value")
	}

	if e := p.ensure(readFn); e != nil {
		return 0, e
	}

	rightShift, usedBits := p.align()
	remainBits := 8 - usedBits

	if bitLen <= remainBits {
		if rightShift {
			p.accBit >>= bitLen
			v := p.accumulator &^ ((p.accBit - 1) | p.accBit)
			return v << usedBits, nil
		} else {
			p.accBit <<= bitLen
			v := p.accumulator & (p.accBit - 1)
			return v >> usedBits, nil
		}
	}
	bitLen -= remainBits

	v := p.accumulator
	p.accBit = 0
	if e := p.ensure(readFn); e != nil {
		return 0, e
	}

	if rightShift {
		p.accBit >>= bitLen
		v &= 0xFF >> usedBits
		w := uint16(v)<<8 | uint16(p.accumulator)
		w <<= bitLen
		v = uint8(w >> 8)
	} else {
		p.accBit <<= bitLen
		v &= 0xFF << usedBits
		w := uint16(v) | uint16(p.accumulator)<<8
		w >>= bitLen
		v = uint8(w)
		v >>= usedBits - bitLen
	}
	return v, nil
}

func newBitReader(direction BitBuilderDirection) bitReader {
	switch direction {
	case FirstLow:
		return bitReader{accInit: initFirstLow}
	case FirstHigh:
		return bitReader{accInit: initFirstHigh}
	default:
		panic("illegal value")
	}
}

func NewBitIoReader(direction BitBuilderDirection, byteReader io.ByteReader) *BitIoReader {
	if byteReader == nil {
		panic("illegal value")
	}
	return &BitIoReader{byteReader: byteReader, bitReader: newBitReader(direction)}
}

type BitIoReader struct {
	byteReader io.ByteReader
	bitReader
}

func (p *BitIoReader) ReadBool() (bool, error) {
	if v, _, e := p.ReadNext(); v != 0 {
		return true, e
	} else {
		return false, e
	}
}

func (p *BitIoReader) ReadBit() (int, error) {
	if v, _, e := p.ReadNext(); v != 0 {
		return 1, e
	} else {
		return 0, e
	}
}

func (p *BitIoReader) ReadNext() (int, uint8, error) {
	return p.readNext(p.byteReader.ReadByte)
}

func (p *BitIoReader) ReadByte() (byte, error) {
	return p.readByte(p.byteReader.ReadByte)
}

func (p *BitIoReader) ReadSubByte(bitLen uint8) (byte, error) {
	return p.readSubByte(bitLen, p.byteReader.ReadByte)
}

func NewBitArrayReader(direction BitBuilderDirection, bytes []byte) *BitArrayReader {
	return &BitArrayReader{bytes: bytes, bitReader: newBitReader(direction)}
}

type BitArrayReader struct {
	bytes []byte
	bitReader
}

func (p *BitArrayReader) IsArrayDepleted() bool {
	return len(p.bytes) == 0
}

func (p *BitArrayReader) ReadBool() bool {
	if v, _ := p.ReadNext(); v != 0 {
		return true
	} else {
		return false
	}
}

func (p *BitArrayReader) ReadBit() int {
	if v, _ := p.ReadNext(); v != 0 {
		return 1
	} else {
		return 0
	}
}

func (p *BitArrayReader) _read() (uint8, error) {
	v := p.bytes[0]
	p.bytes = p.bytes[1:]
	return v, nil
}

func (p *BitArrayReader) ReadNext() (int, uint8) {
	i, b, _ := p.readNext(p._read)
	return i, b
}

func (p *BitArrayReader) ReadByte() byte {
	b, _ := p.readByte(p._read)
	return b
}

func (p *BitArrayReader) ReadSubByte(bitLen uint8) byte {
	b, _ := p.readSubByte(bitLen, p._read)
	return b
}

func NewBitStrReader(direction BitBuilderDirection, bytes ByteString) *BitStrReader {
	return &BitStrReader{bytes: string(bytes), bitReader: newBitReader(direction)}
}

type BitStrReader struct {
	bytes string
	bitReader
}

func (p *BitStrReader) IsArrayDepleted() bool {
	return len(p.bytes) == 0
}

func (p *BitStrReader) ReadBool() bool {
	if v, _ := p.ReadNext(); v != 0 {
		return true
	} else {
		return false
	}
}

func (p *BitStrReader) ReadBit() int {
	if v, _ := p.ReadNext(); v != 0 {
		return 1
	} else {
		return 0
	}
}

func (p *BitStrReader) _read() (uint8, error) {
	v := p.bytes[0]
	p.bytes = p.bytes[1:]
	return v, nil
}

func (p *BitStrReader) ReadNext() (int, uint8) {
	i, b, _ := p.readNext(p._read)
	return i, b
}

func (p *BitStrReader) ReadByte() byte {
	b, _ := p.readByte(p._read)
	return b
}

func (p *BitStrReader) ReadSubByte(bitLen uint8) byte {
	b, _ := p.readSubByte(bitLen, p._read)
	return b
}
