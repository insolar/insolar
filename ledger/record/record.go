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

package record

import (
	"encoding/binary"
	"io"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/cryptohelpers/hash"
)

// ProjectionType is a "view filter" for record.
// E.g. we can read whole object or just it's hash.
type ProjectionType uint32

// Memory is actual contracts' state, variables etc.
type Memory []byte

// RelativePulseNumber - special value of PulseNum, it means a Drop-relative Pulse Number.
// It is only allowed for Storage.
const RelativePulseNumber core.PulseNumber = 65536

// TypeID encodes a record object type.
type TypeID uint32

// ID is a composite identifier for records.
//
// Hash is a bytes slice here to avoid copy of Hash array.
type ID struct {
	Pulse core.PulseNumber
	Hash  []byte
}

// Record is base interface for all records.
type Record interface {
	Type() TypeID
}

// SHA3Hash224 hashes Record by it's CBOR representation and type identifier.
func SHA3Hash224(rec Record) []byte {
	cborBlob := MustEncode(rec)
	return hash.SHA3hash224(getTypeIDbyRecord(rec), hashableBytes(cborBlob))
}

// WriteHash implements hash.Writer interface.
func (id ID) WriteHash(w io.Writer) {
	b := ID2Bytes(id)
	err := binary.Write(w, binary.BigEndian, b)
	if err != nil {
		panic("binary.Write failed:" + err.Error())
	}
}

// IsEqual checks equality of IDs.
func (id ID) IsEqual(id2 ID) bool {
	if (id.Hash == nil) != (id2.Hash == nil) {
		return false
	}
	if len(id.Hash) != len(id2.Hash) {
		return false
	}
	for i := range id.Hash {
		if id.Hash[i] != id2.Hash[i] {
			return false
		}
	}
	return true
}

// WriteHash implements hash.Writer interface.
func (id TypeID) WriteHash(w io.Writer) {
	err := binary.Write(w, binary.BigEndian, id)
	if err != nil {
		panic("binary.Write failed:" + err.Error())
	}
}

// CoreID generates Reference byte representation (key without prefix).
func (id *ID) CoreID() *core.RecordID {
	var b core.RecordID
	_ = copy(b[:], ID2Bytes(*id))
	return &b
}

// Reference allows to address any record across the whole network.
type Reference struct {
	Domain ID
	Record ID
}

// CoreRef generates Reference byte representation (key without prefix).
func (ref *Reference) CoreRef() *core.RecordRef {
	var b core.RecordRef
	// Record part should go first so we can iterate keys of a certain slot
	_ = copy(b[:core.RecordIDSize], ID2Bytes(ref.Record))
	_ = copy(b[core.RecordIDSize:], ID2Bytes(ref.Domain))
	return &b
}

// IsEqual checks equality of References.
func (ref Reference) IsEqual(ref2 Reference) bool {
	if !ref.Domain.IsEqual(ref2.Domain) {
		return false
	}
	return ref.Record.IsEqual(ref2.Record)
}

// IsNotEqual checks non equality of References.
func (ref Reference) IsNotEqual(ref2 Reference) bool {
	return !ref.IsEqual(ref2)
}

// String returns Base58 string representation of Reference.
func (ref Reference) String() string {
	return ref.CoreRef().String()
}
