/*
 *    Copyright 2018 INS Ecosystem
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

	"github.com/insolar/insolar/ledger/hash"
)

const (
	// HashSize is a record hash size. We use 224-bit SHA-3 hash (28 bytes).
	HashSize = 28
	// PulseNumSize - 4 bytes is a PulseNum size (uint32)
	PulseNumSize = 4
	// IDSize is the size in bytes of ID binary representation.
	IDSize = PulseNumSize + HashSize
	// RefIDSize is the size in bytes of Reference binary representation.
	RefIDSize = IDSize * 2
)

// ProjectionType is a "view filter" for record.
// E.g. we can read whole object or just it's hash.
type ProjectionType uint32

// Memory is actual contracts' state, variables etc.
type Memory []byte

// PulseNum is a sequential number of Pulse.
// Upper 2 bits are reserved for use in references (scope), must be zero otherwise.
// Valid Absolute PulseNum must be >65536.
// If PulseNum <65536 it is a relative PulseNum
type PulseNum uint32

// SpecialPulseNumber - special value of PulseNum, it means a Drop-relative Pulse Number.
// It is only allowed for Storage.
const SpecialPulseNumber PulseNum = 65536

// TypeID encodes a record object type.
type TypeID uint32

// ID is a composite identifier for records.
//
// Hash is a bytes slice here to avoid copy of Hash array.
type ID struct {
	Pulse PulseNum
	Hash  []byte
}

// Record is base interface for all records.
type Record interface {
	Domain() *Reference
}

// SHA3Hash224 hashes Record by it's CBOR representation and type identifier.
func SHA3Hash224(rec Record) []byte {
	cborBlob := MustEncode(rec)
	return hash.SHA3hash224(getTypeIDbyRecord(rec), hashableBytes(cborBlob))
}

// ID evaluates record ID on PulseNum for Record.
func (pn PulseNum) ID(rec Record) ID {
	raw, err := EncodeToRaw(rec)
	if err != nil {
		panic(err)
	}
	return ID{
		Pulse: pn,
		Hash:  raw.Hash(),
	}
}

// Bytes evaluates bytes representation of PulseNum and Record pair.
func (pn PulseNum) Bytes(rec Record) []byte {
	return ID2Bytes(pn.ID(rec))
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

// Reference allows to address any record across the whole network.
type Reference struct {
	Domain ID
	Record ID
}

// Key generates Reference byte representation (key without prefix).
func (ref *Reference) Key() []byte {
	b := make([]byte, RefIDSize)
	// Record part should go first so we can iterate keys of a certain slot
	_ = copy(b[:IDSize], ID2Bytes(ref.Record))
	_ = copy(b[IDSize:], ID2Bytes(ref.Domain))
	return b
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
