// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package reference

import (
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/pkg/errors"

	"github.com/insolar/insolar/longbits"
	"github.com/insolar/insolar/pulse"
)

var byteOrder = binary.BigEndian

func NewRecordID(pn pulse.Number, hash longbits.Bits224) Local {
	return NewLocal(pn, 0, hash) // scope is not allowed for RecordID
}

func NewLocal(pn pulse.Number, scope uint8, hash longbits.Bits224) Local {
	if !pn.IsSpecialOrTimePulse() {
		panic(fmt.Sprintf("illegal value: %d", pn))
	}
	return Local{pulseAndScope: pn.WithFlags(scope), hash: hash}
}

const JetDepthPosition = 0

type Local struct {
	pulseAndScope uint32 // pulse + scope
	hash          longbits.Bits224
}

// IsEmpty - check for void
func (v Local) IsEmpty() bool {
	return v.pulseAndScope == 0
}

// NotEmpty - check for non void
func (v Local) NotEmpty() bool {
	return !v.IsEmpty()
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
	if len(p) < LocalBinaryPulseAndScopeSize {
		return copy(p, v.pulseAndScopeAsBytes()), nil
	}

	byteOrder.PutUint32(p, v.pulseAndScope)
	return LocalBinaryPulseAndScopeSize + copy(p[LocalBinaryPulseAndScopeSize:], v.hash.AsBytes()), nil
}

func (v Local) len() int {
	return LocalBinaryPulseAndScopeSize + len(v.hash)
}

func (v *Local) AsByteString() longbits.ByteString {
	return longbits.NewByteString(v.AsBytes())
}

func (v *Local) AsBytes() []byte {
	val := make([]byte, v.len())
	byteOrder.PutUint32(val, v.pulseAndScope)
	_, _ = v.hash.Read(val[LocalBinaryPulseAndScopeSize:])
	return val
}

func (v *Local) pulseAndScopeAsBytes() []byte {
	val := make([]byte, LocalBinaryPulseAndScopeSize)
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
		return p.v.hash[p.o-1-LocalBinaryPulseAndScopeSize], nil
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
		p.v.hash[p.o-LocalBinaryPulseAndScopeSize] = c
	}
	p.o++
	return nil
}

func (p *byteWriter) isFull() bool {
	return int(p.o) >= p.v.len()
}

// Encode encodes Local to string with chosen encoder.
func (v Local) Encode(enc Encoder) string {
	repr, err := enc.EncodeRecord(&v)
	if err != nil {
		return ""
	}
	return repr
}

// String implements stringer on ID and returns base64 encoded value
func (v Local) String() string {
	return v.Encode(DefaultEncoder())
}

// Bytes returns byte slice of ID.
func (v Local) Bytes() []byte {
	return v.AsBytes()
}

// Equal checks if reference points to the same record
func (v *Local) Equal(other Local) bool {
	if v == nil {
		return false
	}
	return v.pulseAndScope == other.pulseAndScope && v.hash == other.hash
}

func (v Local) Compare(other Local) int {
	if v.pulseAndScope < other.pulseAndScope {
		return -1
	} else if v.pulseAndScope > other.pulseAndScope {
		return 1
	}

	return v.hash.Compare(other.hash)
}

// Pulse returns a copy of Pulse part of ID.
func (v Local) Pulse() pulse.Number {
	return v.GetPulseNumber()
}

// Hash returns a copy of Hash part of ID
func (v Local) Hash() []byte {
	rv := make([]byte, len(v.hash))
	copy(rv, v.hash[:])
	return rv
}

// MarshalJSON serializes ID into JSONFormat
func (v *Local) MarshalJSON() ([]byte, error) {
	if v == nil {
		return json.Marshal(nil)
	}
	return json.Marshal(v.String())
}

func (v Local) MarshalBinary() ([]byte, error) {
	return v.Marshal()
}

func (v Local) Marshal() ([]byte, error) {
	return v.AsBytes(), nil
}

func (v Local) MarshalTo(data []byte) (int, error) {
	if len(data) < LocalBinarySize {
		return 0, errors.New("not enough bytes to marshal reference.Local")
	}
	return copy(data, v.AsBytes()), nil
}

func (v *Local) UnmarshalJSON(data []byte) error {
	var repr interface{}

	err := json.Unmarshal(data, &repr)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal reference.Local")
	}

	switch realRepr := repr.(type) {
	case string:
		intermidiate, err := DefaultDecoder().Decode(realRepr)
		if err != nil {
			return errors.Wrap(err, "failed to unmarshal reference.Local")
		}
		*v = intermidiate.addressLocal
	case nil:
	default:
		return errors.Wrapf(err, "unexpected type %T when unmarshal reference.Local", repr)
	}

	return nil
}

func (v *Local) UnmarshalBinary(data []byte) error {
	return v.Unmarshal(data)
}

func (v *Local) Unmarshal(data []byte) error {
	if len(data) < LocalBinarySize {
		return errors.New("not enough bytes to unmarshal reference.Local")
	}

	writer := v.asWriter()
	for i := 0; i < LocalBinarySize; i++ {
		_ = writer.WriteByte(data[i])
	}

	return nil
}

func (v Local) Size() int {
	return LocalBinarySize
}

func (v Local) debugStringJet() string {
	depth, prefix := int(v.hash[JetDepthPosition]), v.hash[1:]

	if depth == 0 {
		return "[JET 0 -]"
	} else if depth > 8*(len(v.hash)-1) {
		return fmt.Sprintf("[JET: <wrong format> %d %b]", depth, prefix)
	}

	res := strings.Builder{}
	res.WriteString(fmt.Sprintf("[JET %d ", depth))

	for i := 0; i < depth; i++ {
		bytePos, bitPos := i/8, 7-i%8

		byteValue := prefix[bytePos]
		bitValue := byteValue >> uint(bitPos) & 0x01
		bitString := strconv.Itoa(int(bitValue))
		res.WriteString(bitString)
	}

	res.WriteString("]")
	return res.String()
}

// DebugString prints ID in human readable form.
func (v *Local) DebugString() string {
	if v == nil {
		return NilRef
	} else if v.Pulse().IsJet() {
		// TODO: remove this branch after finish transition to JetID
		return v.debugStringJet()
	}

	return fmt.Sprintf("%s [%d | %d | %s]", v.String(), v.Pulse(), v.getScope(), base64.RawURLEncoding.EncodeToString(v.Hash()))
}
