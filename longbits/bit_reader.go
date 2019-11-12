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
	byteReader  io.ByteReader
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

func (p *bitReader) readNext() (int, uint8) {
	if p.accBit == 0 {
		panic("illegal state")
	}
	m := p.accBit
	if rightShift, _ := p.align(); rightShift {
		p.accBit >>= 1
	} else {
		p.accBit <<= 1
	}
	return int(p.accumulator & m), m
}

func (p *bitReader) readByte(v byte) byte {
	if p.accBit == 0 {
		panic("illegal state")
	}
	w := p.accumulator
	if rightShift, usedBits := p.align(); rightShift {
		v <<= usedBits
		w >>= 8 - usedBits
	} else {
		v >>= usedBits
		w <<= 8 - usedBits
	}
	return v | w
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

func (p *BitIoReader) fillIn() error {
	var err error
	p.accumulator, err = p.byteReader.ReadByte()
	return err
}

func (p *BitIoReader) ReadNext() (int, uint8, error) {
	if p.accBit == 0 {
		if p.accInit == 0 {
			p.accInit = initFirstLow
		}
		if err := p.fillIn(); err != nil {
			return 0, 0, err
		}
		p.accBit = p.accInit
	}
	v, m := p.readNext()
	return v, m, nil
}

func (p *BitIoReader) ReadByte() (byte, error) {
	switch p.accBit {
	case 0:
		if p.accInit == 0 {
			p.accInit = initFirstLow
		}
		return p.byteReader.ReadByte()
	case p.accInit:
		p.accBit = 0
		return p.accumulator, nil
	}

	v := p.accumulator
	if err := p.fillIn(); err != nil {
		return 0, err
	}
	return p.readByte(v), nil
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

func (p *BitArrayReader) fillIn() {
	p.accumulator = p.bytes[0]
	p.bytes = p.bytes[1:]
}

func (p *BitArrayReader) ReadNext() (int, uint8) {
	if p.accBit == 0 {
		if p.accInit == 0 {
			p.accInit = initFirstLow
		}
		p.fillIn()
		p.accBit = p.accInit
	}
	return p.readNext()
}

func (p *BitArrayReader) ReadByte() byte {
	switch p.accBit {
	case 0:
		if p.accInit == 0 {
			p.accInit = initFirstLow
		}
		p.fillIn()
		return p.accumulator
	case p.accInit:
		p.accBit = 0
		return p.accumulator
	}

	v := p.accumulator
	p.fillIn()
	return p.readByte(v)
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

func (p *BitStrReader) fillIn() {
	p.accumulator = p.bytes[0]
	p.bytes = p.bytes[1:]
}

func (p *BitStrReader) ReadNext() (int, uint8) {
	if p.accBit == 0 {
		if p.accInit == 0 {
			p.accInit = initFirstLow
		}
		p.fillIn()
		p.accBit = p.accInit
	}
	return p.readNext()
}

func (p *BitStrReader) ReadByte() byte {
	switch p.accBit {
	case 0:
		if p.accInit == 0 {
			p.accInit = initFirstLow
		}
		p.fillIn()
		return p.accumulator
	case p.accInit:
		p.accBit = 0
		return p.accumulator
	}

	v := p.accumulator
	p.fillIn()
	return p.readByte(v)
}
