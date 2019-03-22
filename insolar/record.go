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

package insolar

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	base58 "github.com/jbenet/go-base58"
	"github.com/pkg/errors"
)

const (
	// RecordHashSize is a record hash size. We use 224-bit SHA-3 hash (28 bytes).
	RecordHashSize = 28
	// RecordIDSize is relative record address.
	RecordIDSize = PulseNumberSize + RecordHashSize
	// RecordHashOffset is a offset where hash bytes starts in RecordID.
	RecordHashOffset = PulseNumberSize
	// RecordRefSize is absolute records address (including domain ID).
	RecordRefSize = RecordIDSize * 2
	// RecordRefIDSeparator is character that separates RecordID from DomainID in serialized RecordRef.
	RecordRefIDSeparator = "."
)

// RecordID is a unified record ID.
type RecordID [RecordIDSize]byte

// String implements stringer on RecordID and returns base58 encoded value
func (id *RecordID) String() string {
	return base58.Encode(id[:])
}

// NewRecordID generates RecordID byte representation.
func NewRecordID(pulse PulseNumber, hash []byte) *RecordID {
	var id RecordID
	copy(id[:PulseNumberSize], pulse.Bytes())
	copy(id[RecordHashOffset:], hash)
	return &id
}

// Bytes returns byte slice of RecordID.
func (id *RecordID) Bytes() []byte {
	return id[:]
}

// Pulse returns a copy of Pulse part of RecordID.
func (id *RecordID) Pulse() PulseNumber {
	pulse := binary.BigEndian.Uint32(id[:PulseNumberSize])
	return PulseNumber(pulse)
}

// Hash returns a copy of Hash part of RecordID.
func (id *RecordID) Hash() []byte {
	recHash := make([]byte, RecordHashSize)
	copy(recHash, id[RecordHashOffset:])
	return recHash
}

// Equal checks if reference points to the same record.
func (id *RecordID) Equal(other *RecordID) bool {
	if id == nil || other == nil {
		return false
	}
	return *id == *other
}

// NewIDFromBase58 deserializes RecordID from base58 encoded string.
func NewIDFromBase58(str string) (*RecordID, error) {
	decoded := base58.Decode(str)
	if len(decoded) != RecordIDSize {
		return nil, errors.New("bad RecordID size")
	}
	var id RecordID
	copy(id[:], decoded)
	return &id, nil
}

// MarshalJSON serializes ID into JSON.
func (id *RecordID) MarshalJSON() ([]byte, error) {
	if id == nil {
		return json.Marshal(nil)
	}
	return json.Marshal(id.String())
}

// RecordRef is a unified record reference.
type RecordRef [RecordRefSize]byte

// NewRecordRef returns RecordRef composed from domain and record
func NewRecordRef(domain RecordID, record RecordID) *RecordRef {
	var ref RecordRef
	ref.SetDomain(domain)
	ref.SetRecord(record)
	return &ref
}

// SetDomain set domain's RecordID.
func (ref *RecordRef) SetDomain(recID RecordID) {
	copy(ref[RecordIDSize:], recID[:])
}

// SetRecord set record's RecordID.
func (ref *RecordRef) SetRecord(recID RecordID) {
	copy(ref[:RecordIDSize], recID[:])
}

// Domain returns domain ID part of reference.
func (ref RecordRef) Domain() *RecordID {
	var id RecordID
	copy(id[:], ref[RecordIDSize:])
	return &id
}

// Record returns record's RecordID.
func (ref *RecordRef) Record() *RecordID {
	if ref == nil {
		return nil
	}
	var id RecordID
	copy(id[:], ref[:RecordIDSize])
	return &id
}

// String outputs base58 RecordRef representation.
func (ref RecordRef) String() string {
	return ref.Record().String() + RecordRefIDSeparator + ref.Domain().String()
}

// FromSlice : After CBOR Marshal/Unmarshal Ref can be converted to byte slice, this converts it back
func (ref RecordRef) FromSlice(from []byte) RecordRef {
	for i := 0; i < RecordRefSize; i++ {
		ref[i] = from[i]
	}
	return ref
}

// Bytes returns byte slice of RecordRef.
func (ref RecordRef) Bytes() []byte {
	return ref[:]
}

// Equal checks if reference points to the same record.
func (ref RecordRef) Equal(other RecordRef) bool {
	return ref == other
}

// IsEmpty - check for void
func (ref RecordRef) IsEmpty() bool {
	return ref.Equal(RecordRef{})
}

// Compare compares two record references
func (ref RecordRef) Compare(other RecordRef) int {
	return bytes.Compare(ref.Bytes(), other.Bytes())
}

// NewRefFromBase58 deserializes reference from base58 encoded string.
func NewRefFromBase58(str string) (*RecordRef, error) {
	parts := strings.SplitN(str, RecordRefIDSeparator, 2)
	if len(parts) < 2 {
		return nil, errors.New("bad reference format")
	}
	recordID, err := NewIDFromBase58(parts[0])
	if err != nil {
		return nil, errors.Wrap(err, "bad record part")
	}
	domainID, err := NewIDFromBase58(parts[1])
	if err != nil {
		return nil, errors.Wrap(err, "bad domain part")
	}
	return NewRecordRef(*domainID, *recordID), nil
}

// MarshalJSON serializes reference into JSON.
func (ref *RecordRef) MarshalJSON() ([]byte, error) {
	if ref == nil {
		return json.Marshal(nil)
	}
	return json.Marshal(ref.String())
}

// DebugString prints ID in human readable form.
func (id *RecordID) DebugString() string {
	if id == nil {
		return "<nil>"
	}

	// TODO: remove this branch after finish transition to JetID
	pulse := NewPulseNumber(id[:PulseNumberSize])
	if pulse == PulseNumberJet {
		depth := int(id[PulseNumberSize])
		if depth == 0 {
			return "[JET 0 -]"
		}

		prefix := id[PulseNumberSize+1:]
		var res strings.Builder
		res.WriteString("[JET ")
		res.WriteString(strconv.Itoa(depth))
		res.WriteString(" ")

		for _, b := range prefix {
			for j := 7; j >= 0; j-- {
				if 0 == (b >> uint(j) & 0x01) {
					res.WriteString("0")
				} else {
					res.WriteString("1")
				}

				depth--
				if depth == 0 {
					res.WriteString("]")
					return res.String()
				}
			}
		}

		return fmt.Sprintf("[JET: <wrong format> %d %b]", depth, prefix)
	}

	return fmt.Sprintf("[%d | %s]", id.Pulse(), id.String())
}
