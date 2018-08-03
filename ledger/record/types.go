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
)

const (
	// HashSize is a record hash size. We use 224-bit SHA-3 hash (28 bytes).
	HashSize = 28
	// PulseNumSize - 4 bytes is a PulseNum size (uint32)
	PulseNumSize = 4
	// IDSize is an record identifier size.
	IDSize = PulseNumSize + HashSize
)

// PulseNum is a sequential number of Pulse.
// Upper 2 bits are reserved for use in references (scope), must be zero otherwise.
// Valid Absolute PulseNum must be >65536.
// If PulseNum <65536 it is a relative PulseNum
type PulseNum uint32

// ID evaluates record ID on PulseNum for Record
func (pn PulseNum) ID(rec Record) ID {
	return Key2ID(pn.Key(rec))
}

// Key evaluates record Key on PulseNum for Record
func (pn PulseNum) Key(rec Record) Key {
	raw, err := EncodeToRaw(rec)
	if err != nil {
		panic(err)
	}
	return Key{
		Pulse: pn,
		Hash:  raw.Hash(),
	}
}

// SpecialPulseNumber - special value of PulseNum, it means a Drop-relative Pulse Number.
// It is only allowed for Storage.
const SpecialPulseNumber PulseNum = 65536

// Hash is hash sum of record, 24-byte array.
type Hash [HashSize]byte

// ID is a record ID. Compounds PulseNum and Type
type ID [IDSize]byte

// WriteHash implements hash.Writer interface.
func (id ID) WriteHash(w io.Writer) {
	err := binary.Write(w, binary.BigEndian, id)
	if err != nil {
		panic("binary.Write failed:" + err.Error())
	}
}

// Key is a composite key for storage methods.
//
// Key and ID converts one to another in both directions.
// Hash is a bytes slice here to avoid copy to Hash array.
type Key struct {
	Pulse PulseNum
	Hash  []byte
}

// TypeID encodes a record object type.
type TypeID uint32

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
