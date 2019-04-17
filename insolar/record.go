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
	// RecordHashOffset is a offset where hash bytes starts in ID.
	RecordHashOffset = PulseNumberSize
	// RecordRefSize is absolute records address (including domain ID).
	RecordRefSize = RecordIDSize * 2
	// RecordRefIDSeparator is character that separates ID from DomainID in serialized Reference.
	RecordRefIDSeparator = "."
)

// ID is a unified record ID.
type ID [RecordIDSize]byte

// String implements stringer on ID and returns base58 encoded value
func (id *ID) String() string {
	return base58.Encode(id[:])
}

// NewID generates ID byte representation.
func NewID(pulse PulseNumber, hash []byte) *ID {
	var id ID
	copy(id[:PulseNumberSize], pulse.Bytes())
	copy(id[RecordHashOffset:], hash)
	return &id
}

// Bytes returns byte slice of ID.
func (id ID) Bytes() []byte {
	return id[:]
}

// Pulse returns a copy of Pulse part of ID.
func (id *ID) Pulse() PulseNumber {
	pulse := binary.BigEndian.Uint32(id[:PulseNumberSize])
	return PulseNumber(pulse)
}

// Hash returns a copy of Hash part of ID.
func (id *ID) Hash() []byte {
	recHash := make([]byte, RecordHashSize)
	copy(recHash, id[RecordHashOffset:])
	return recHash
}

// Equal checks if reference points to the same record.
func (id *ID) Equal(other ID) bool {
	if id == nil {
		return false
	}
	return *id == other
}

// NewIDFromBase58 deserializes ID from base58 encoded string.
func NewIDFromBase58(str string) (*ID, error) {
	decoded := base58.Decode(str)
	if len(decoded) != RecordIDSize {
		return nil, errors.New("bad ID size")
	}
	var id ID
	copy(id[:], decoded)
	return &id, nil
}

// MarshalJSON serializes ID into JSON.
func (id *ID) MarshalJSON() ([]byte, error) {
	if id == nil {
		return json.Marshal(nil)
	}
	return json.Marshal(id.String())
}

// Reference is a unified record reference.
type Reference [RecordRefSize]byte

// NewReference returns Reference composed from domain and record
func NewReference(domain ID, record ID) *Reference {
	var ref Reference
	ref.SetDomain(domain)
	ref.SetRecord(record)
	return &ref
}

// SetDomain set domain's ID.
func (ref *Reference) SetDomain(recID ID) {
	copy(ref[RecordIDSize:], recID[:])
}

// SetRecord set record's ID.
func (ref *Reference) SetRecord(recID ID) {
	copy(ref[:RecordIDSize], recID[:])
}

// Domain returns domain ID part of reference.
func (ref Reference) Domain() *ID {
	var id ID
	copy(id[:], ref[RecordIDSize:])
	return &id
}

// Record returns record's ID.
func (ref *Reference) Record() *ID {
	if ref == nil {
		return nil
	}
	var id ID
	copy(id[:], ref[:RecordIDSize])
	return &id
}

// String outputs base58 Reference representation.
func (ref Reference) String() string {
	return ref.Record().String() + RecordRefIDSeparator + ref.Domain().String()
}

// FromSlice : After CBOR Marshal/Unmarshal Ref can be converted to byte slice, this converts it back
func (ref Reference) FromSlice(from []byte) Reference {
	for i := 0; i < RecordRefSize; i++ {
		ref[i] = from[i]
	}
	return ref
}

// Bytes returns byte slice of Reference.
func (ref Reference) Bytes() []byte {
	return ref[:]
}

// Equal checks if reference points to the same record.
func (ref Reference) Equal(other Reference) bool {
	return ref == other
}

// IsEmpty - check for void
func (ref Reference) IsEmpty() bool {
	return ref.Equal(Reference{})
}

// Compare compares two record references
func (ref Reference) Compare(other Reference) int {
	return bytes.Compare(ref.Bytes(), other.Bytes())
}

// NewReferenceFromBase58 deserializes reference from base58 encoded string.
func NewReferenceFromBase58(str string) (*Reference, error) {
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
	return NewReference(*domainID, *recordID), nil
}

// MarshalJSON serializes reference into JSON.
func (ref *Reference) MarshalJSON() ([]byte, error) {
	if ref == nil {
		return json.Marshal(nil)
	}
	return json.Marshal(ref.String())
}

func (ref Reference) Marshal() ([]byte, error) {
	return ref[:], nil
}

func (ref *Reference) MarshalTo(data []byte) (int, error) {
	copy(data, ref[:])
	return RecordRefSize, nil
}

func (ref *Reference) Unmarshal(data []byte) error {
	if len(data) != RecordRefSize {
		return errors.New("Not enough bytes to unpack Reference")
	}
	copy(ref[:], data)
	return nil
}
func (ref *Reference) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, ref)
}
func (ref *Reference) Size() int { return RecordRefSize }

func (id ID) Marshal() ([]byte, error) { return id[:], nil }
func (id *ID) MarshalTo(data []byte) (int, error) {
	copy(data, id[:])
	return RecordIDSize, nil
}
func (id *ID) Unmarshal(data []byte) error {
	if len(data) != RecordIDSize {
		return errors.New("Not enough bytes to unpack ID")
	}
	copy(id[:], data)
	return nil
}
func (id *ID) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, id)
}
func (id *ID) Size() int { return RecordIDSize }
func (id ID) Compare(other ID) int {
	return bytes.Compare(id.Bytes(), other.Bytes())
}

// DebugString prints ID in human readable form.
func (id *ID) DebugString() string {
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
