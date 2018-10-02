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

// RecordID is a unified record ID.
type RecordID [RecordIDSize]byte

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

// NewRefFromBase58 deserializes reference from base58 encoded string.
func NewRefFromBase58(str string) RecordRef {
	// TODO: if str < 20 bytes, always returns 0. need to check this.
	decoded := base58.Decode(str)
	var ref RecordRef
	copy(ref[:], decoded)
	return ref
}

// RandomRef generates random RecordRef
func RandomRef() RecordRef {
	ref := [64]byte{}
	rand.Read(ref[:]) // nolint
	return ref
}
