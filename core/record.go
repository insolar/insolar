/*
 *    Copyright 2018 Insolar
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package core

import (
	"encoding/binary"

	"github.com/jbenet/go-base58"

	"github.com/insolar/insolar/cryptoproviders/hash"
)

const (
	// RecordHashSize is a record hash size. We use 224-bit SHA-3 hash (28 bytes).
	RecordHashSize = 28
	// RecordIDSize is relative record address.
	RecordIDSize = PulseNumberSize + RecordHashSize
	// RecordRefSize is absolute records address (including domain ID).
	RecordRefSize = RecordIDSize * 2
)

// RecordID is a unified record ID.
type RecordID [RecordIDSize]byte

// NewRecordID generates RecordID byte representation.
func NewRecordID(pulse PulseNumber, hash []byte) *RecordID {
	var id RecordID
	copy(id[:PulseNumberSize], pulse.Bytes())
	copy(id[PulseNumberSize:], hash)
	return &id
}

// Bytes returns byte slice of RecordID.
func (id *RecordID) Bytes() []byte {
	return id[:]
}

// Pulse returns byte slice of RecordID.
func (id *RecordID) Pulse() PulseNumber {
	pulse := binary.BigEndian.Uint32(id[:PulseNumberSize])
	return PulseNumber(pulse)
}

// Equal checks if reference points to the same record.
func (id *RecordID) Equal(other *RecordID) bool {
	if id == nil || other == nil {
		return false
	}
	return *id == *other
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
	var id RecordID
	copy(id[:], ref[:RecordIDSize])
	return &id
}

// String outputs base58 RecordRef representation.
func (ref RecordRef) String() string {
	return base58.Encode(ref[:])
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

// GenRequest calculates RecordRef for request message from pulse number and request's payload.
func GenRequest(pn PulseNumber, payload []byte) *RecordRef {
	ref := NewRecordRef(
		RecordID{},
		*NewRecordID(pn, hash.IDHashBytes(payload)),
	)
	return ref
}

// NewRefFromBase58 deserializes reference from base58 encoded string.
func NewRefFromBase58(str string) RecordRef {
	// TODO: if str < 20 bytes, always returns 0. need to check this.
	decoded := base58.Decode(str)
	var ref RecordRef
	copy(ref[:], decoded)
	return ref
}
