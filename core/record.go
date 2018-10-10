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
	"crypto/rand"

	"github.com/jbenet/go-base58"

	"github.com/insolar/insolar/cryptohelpers/hash"
)

const (
	// RecordHashSize is a record hash size. We use 224-bit SHA-3 hash (28 bytes).
	RecordHashSize = 28
	// RecordIDSize is relative record address.
	RecordIDSize = PulseNumberSize + RecordHashSize
	// RecordRefSize is absolute records address (including domain ID).
	RecordRefSize = RecordIDSize * 2
)

// RecordRef is a unified record reference.
type RecordRef [RecordRefSize]byte

// ComposeRecordRef returns RecordRef composed from domain and record
func ComposeRecordRef(domain RecordID, record RecordID) (ref RecordRef) {
	(&ref).SetDomain(domain)
	(&ref).SetRecord(record)
	return
}

// SetRecord set record's RecordID.
func (ref *RecordRef) SetRecord(recID RecordID) {
	copy(ref[:RecordIDSize], recID[:])
}

// SetDomain set domain's RecordID.
func (ref *RecordRef) SetDomain(recID RecordID) {
	copy(ref[RecordIDSize:], recID[:])
}

// GetRecordID returns record's RecordID.
func (ref *RecordRef) GetRecordID() (id RecordID) {
	copy(id[:], ref[:RecordIDSize])
	return id
}

// GetDomainID returns domain's RecordID.
func (ref *RecordRef) GetDomainID() (id RecordID) {
	copy(id[:], ref[RecordIDSize:])
	return id
}

// RecordID is a unified record ID.
type RecordID [RecordIDSize]byte

// GenRecordID generates RecordID byte representation.
func GenRecordID(pn PulseNumber, h []byte) (recid RecordID) {
	copy(recid[:PulseNumberSize], pn.Bytes())
	copy(recid[PulseNumberSize:], h)
	return
}

// Bytes returns byte slice of RecordID.
func (id *RecordID) Bytes() []byte {
	return id[:]
}

// String outputs base58 RecordRef representation.
func (ref RecordRef) String() string {
	return base58.Encode(ref[:])
}

// Equal checks if reference points to the same record.
func (ref RecordRef) Equal(other RecordRef) bool {
	return ref == other
}

// Domain returns domain ID part of reference.
func (ref RecordRef) Domain() RecordID {
	var domain RecordID
	copy(domain[:], ref[RecordIDSize:])
	return domain
}

// GenRequest calculates RecordRef for request message from pulse number and request's payload.
func GenRequest(pn PulseNumber, payload []byte) *RecordRef {
	ref := ComposeRecordRef(
		RecordID{},
		GenRecordID(pn, hash.IDHashBytes(payload)),
	)
	return &ref
}

// NewRefFromBase58 deserializes reference from base58 encoded string.
func NewRefFromBase58(str string) RecordRef {
	// TODO: if str < 20 bytes, always returns 0. need to check this.
	decoded := base58.Decode(str)
	var ref RecordRef
	copy(ref[:], decoded)
	return ref
}

// RandomRef generates random RecordRef
// TODO: DO NOT USE THIS IN PRODUCTION.
// For tests copy this code or move it to test utils.
func RandomRef() RecordRef {
	ref := [64]byte{}
	rand.Read(ref[:]) // nolint
	return ref
}
