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
	"github.com/jbenet/go-base58"
)

const (
	// RecordIDSize is relative record address.
	RecordIDSize = 32
	// RecordRefSize is absolute records address (including domain ID).
	RecordRefSize = RecordIDSize * 2
)

// RecordRef is unified record reference.
type RecordRef [RecordRefSize]byte

func (ref RecordRef) String() string {
	return base58.Encode(ref[:])
}

// Equal checks if reference points to the same record.
func (ref RecordRef) Equal(other RecordRef) bool {
	return ref == other
}

// Domain returns domain ID part of reference.
func (ref RecordRef) Domain() [RecordIDSize]byte {
	var domain [RecordIDSize]byte
	copy(domain[:], ref[RecordIDSize:])
	return domain
}

// String2Ref deserializes reference from base58 encoded string.
func String2Ref(str string) RecordRef {
	decoded := base58.Decode(str)
	var ref RecordRef
	copy(ref[:], decoded)
	return ref
}
