///
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
///

package ref

import (
	"encoding/binary"
	"fmt"
	"github.com/insolar/insolar/network/consensus/common/longbits"
	"github.com/insolar/insolar/network/consensus/common/pulse"
	"io"
)

var byteOrder = binary.BigEndian

func NewRecordID(pn pulse.Number, hash longbits.Bits224) Local {
	return NewLocal(pn, 0, hash) // scope is not allowed for RecordID
}

func NewLocal(pn pulse.Number, scope uint8, hash longbits.Bits224) Local {
	if !pn.IsSpecialOrTimePulse() {
		panic("illegal value")
	}
	return Local{pulseAndScope: pn.WithFlags(scope), hash: hash}
}

const pulseAndScopeSize = 4

type Local struct {
	pulseAndScope uint32 // pulse + scope
	hash          longbits.Bits224
}

func (v Local) IsEmpty() bool {
	return v.pulseAndScope == 0
}

func (v Local) GetPulseNumber() pulse.Number {
	return pulse.OfUint32(v.pulseAndScope)
}

func (v Local) GetHash() longbits.Bits224 {
	return v.hash
}

func (v Local) getScope() uint8 {
	return uint8(pulse.FlagsOf(v.pulseAndScope))
}

func (v *Local) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write(v.pulseAndScopeAsBytes())
	if err != nil {
		return int64(n), err
	}

	var n2 int64
	n2, err = v.hash.WriteTo(w)
	return int64(n) + n2, err
}

func (v *Local) Read(p []byte) (n int, err error) {
	if len(p) < pulseAndScopeSize {
		return copy(p, v.pulseAndScopeAsBytes()), nil
	}

	byteOrder.PutUint32(p, v.pulseAndScope)
	return pulseAndScopeSize + copy(p[pulseAndScopeSize:], v.hash.AsBytes()), nil
}

func (v Local) len() int {
	return pulseAndScopeSize + len(v.hash)
}

func (v Local) String() string {
	sc := v.getScope()
	if sc != 0 {
		return fmt.Sprintf("%d(%d)_0x%08x", v.GetPulseNumber(), sc, v.hash.FoldToUint64())
	}
	return fmt.Sprintf("%d_0x%08x", v.GetPulseNumber(), v.hash.FoldToUint64())
}

func (v *Local) AsByteString() longbits.ByteString {
	return longbits.NewByteString(v.AsBytes())
}

func (v *Local) AsBytes() []byte {
	val := make([]byte, v.len())
	byteOrder.PutUint32(val, v.pulseAndScope)
	_, _ = v.hash.Read(val[pulseAndScopeSize:])
	return val
}

func (v *Local) pulseAndScopeAsBytes() []byte {
	val := make([]byte, pulseAndScopeSize)
	byteOrder.PutUint32(val, v.pulseAndScope)
	return val
}

func (v *Local) AsReader() io.ByteReader {
	return v.asReader(uint8(v.len()))
}

func (v *Local) asReader(limit uint8) *byteReader {
	return &byteReader{v: v, s: limit}
}

func (v *Local) asWriter() *byteWriter {
	return &byteWriter{v: v}
}

type byteReader struct {
	v *Local
	o uint8
	s uint8
}

func (p *byteReader) ReadByte() (byte, error) {
	switch p.o {
	case 0, 1, 2, 3:
		v := byte(p.v.pulseAndScope >> (p.o << 3))
		p.o++
		return v, nil
	default:
		if p.o >= p.s {
			return 0, io.EOF
		}
		p.o++
		return p.v.hash[p.o-1-pulseAndScopeSize], nil
	}
}

var _ io.ByteWriter = &byteWriter{}

type byteWriter struct {
	v *Local
	o uint8
}

func (p *byteWriter) WriteByte(c byte) error {
	switch p.o {
	case 0, 1, 2, 3:
		shift := p.o << 3
		p.v.pulseAndScope = uint32(c)<<shift | p.v.pulseAndScope&^0xFF<<shift
	default:
		if p.isFull() {
			return io.ErrUnexpectedEOF
		}
		p.v.hash[p.o-pulseAndScopeSize] = c
	}
	p.o++
	return nil
}

func (p *byteWriter) isFull() bool {
	return int(p.o) >= p.v.len()
}
