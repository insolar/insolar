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

package reference

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/insolar/insolar/longbits"
	"github.com/insolar/insolar/pulse"
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
		v := byte(p.v.pulseAndScope >> ((3 - p.o) << 3))
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
		shift := (3 - p.o) << 3
		p.v.pulseAndScope = uint32(c)<<shift | p.v.pulseAndScope&^(0xFF<<shift)
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
