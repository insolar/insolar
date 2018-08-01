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
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/ugorji/go/codec"
	"golang.org/x/crypto/sha3"
)

// Raw struct contains raw serialized record.
// We need raw blob to not have dependency on record structure changes in future,
// and have ability of consistent hash checking on old records.
type Raw struct {
	Type TypeID
	Data []byte
}

// Hash returns 28 bytes of SHA3 hash on Data field.
func (raw *Raw) Hash() Hash {
	return sha3.Sum224(raw.Data)
}

// Decode decodes Data field of Raw struct as record from CBOR format.
func (raw *Raw) Decode() Record {
	cborH := &codec.CborHandle{}
	rec := getRecordByTypeID(raw.Type)
	dec := codec.NewDecoder(bytes.NewReader(raw.Data), cborH)
	err := dec.Decode(rec)
	if err != nil {
		panic(err)
	}
	return rec
}

// Key2ID converts Key with PulseNum and Hash pair to binary representation (record.ID).
func Key2ID(k Key) ID {
	var id ID
	var err error
	buf := bytes.NewBuffer(id[:0])

	err = binary.Write(buf, binary.BigEndian, k.Pulse)
	if err != nil {
		panic("binary.Write failed to write PulseNum:" + err.Error())
	}
	err = binary.Write(buf, binary.BigEndian, k.Hash)
	if err != nil {
		panic("binary.Write failed to write Hash:" + err.Error())
	}
	return id
}

// ID2Key converts ID to Key with PulseNum and Hash pair.
func ID2Key(id ID) Key {
	return Key{
		Pulse: PulseNum(binary.BigEndian.Uint32(id[:PulseNumSize])),
		Hash:  id[PulseNumSize:],
	}
}

// record type ids for record types
// in use mostly for hashing and deserialization
// (we don't use iota for clarity and predictable ids,
// not depended on defenition order)
const (
	// request record ids
	requestRecordID       TypeID = 1
	callRequestID         TypeID = 2
	lockUnlockRequestID   TypeID = 3
	readRecordRequestID   TypeID = 4
	readObjectID          TypeID = 5
	readObjectCompositeID TypeID = 6
	// result record ids
	resultRecordID              TypeID = 7
	wipeOutRecordID             TypeID = 8
	readRecordResultID          TypeID = 9
	statelessCallResultID       TypeID = 10
	statelessExceptionResultID  TypeID = 11
	readObjectResultID          TypeID = 12
	specialResultID             TypeID = 13
	lockUnlockResultID          TypeID = 14
	rejectionResultID           TypeID = 15
	activationRecordID          TypeID = 16
	classActivateRecordID       TypeID = 17
	objectActivateRecordID      TypeID = 18
	codeRecordID                TypeID = 19
	amendRecordID               TypeID = 20
	classAmendRecordID          TypeID = 21
	memoryMigrationCodeID       TypeID = 22
	deactivationRecordID        TypeID = 23
	objectAmendRecordID         TypeID = 24
	statefulCallResultID        TypeID = 25
	statefulExceptionResultID   TypeID = 26
	enforcedObjectAmendRecordID TypeID = 27
	objectAppendRecordID        TypeID = 28
)

// getRecordByTypeID returns Record interface with concrete record type under the hood.
// This is useful with deserialization cases.
func getRecordByTypeID(id TypeID) Record {
	switch id {
	case requestRecordID:
		return &RequestRecord{}
	case callRequestID:
		return &CallRequest{}
	case lockUnlockRequestID:
		return &LockUnlockRequest{}
	case readRecordRequestID:
		return &ReadRecordRequest{}
	case readObjectID:
		return &ReadObject{}
	case readObjectCompositeID:
		return &ReadObjectComposite{}
	default:
		panic(fmt.Errorf("unknown record type id %v", id))
	}
}

// getRecordByTypeID returns record's TypeID based on concrete record type of Record interface.
func getTypeIDbyRecord(rec Record) TypeID {
	switch v := rec.(type) {
	case RequestRecord, *RequestRecord:
		return requestRecordID
	case CallRequest, *CallRequest:
		return callRequestID
	case LockUnlockRequest, *LockUnlockRequest:
		return lockUnlockRequestID
	case ReadRecordRequest, *ReadRecordRequest:
		return readRecordRequestID
	case ReadObject, *ReadObject:
		return readObjectID
	case ReadObjectComposite, *ReadObjectComposite:
		return readObjectCompositeID
	default:
		panic(fmt.Errorf("can't find record id by type %T", v))
	}
}

// Encode serializes record to CBOR blob.
func Encode(rec Record) ([]byte, error) {
	cborH := &codec.CborHandle{}
	var b bytes.Buffer
	enc := codec.NewEncoder(&b, cborH)
	err := enc.Encode(rec)
	return b.Bytes(), err
}

// MustEncode is helper that wraps a call to a function Encode and panics if the error is non-nil.
func MustEncode(rec Record) []byte {
	b, err := Encode(rec)
	if err != nil {
		panic(err)
	}
	return b
}

// encodeToRaw converts concrete record to Raw record.
func encodeToRaw(rec Record) (Raw, error) {
	b, err := Encode(rec)
	if err != nil {
		panic(err)
	}
	return Raw{
		Type: getTypeIDbyRecord(rec),
		Data: b,
	}, nil
}
