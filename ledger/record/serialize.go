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
	"bytes"
	"fmt"
	"io"

	"github.com/ugorji/go/codec"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/core/hash"
)

// Raw struct contains raw serialized record.
// We need raw blob to not have dependency on record structure changes in future,
// and have ability of consistent hash checking on old records.
type Raw struct {
	Type TypeID
	Data []byte
}

// DecodeToRaw decodes bytes to Raw struct from CBOR.
func DecodeToRaw(b []byte) (*Raw, error) {
	cborH := &codec.CborHandle{}
	var rec Raw
	dec := codec.NewDecoderBytes(b, cborH)
	err := dec.Decode(&rec)
	if err != nil {
		return nil, err
	}
	return &rec, nil
}

// MustEncodeRaw wraps EncodeRaw, panics on encode errors.
func MustEncodeRaw(raw *Raw) []byte {
	b, err := EncodeRaw(raw)
	if err != nil {
		panic(err)
	}
	return b
}

// EncodeRaw encodes Raw to CBOR.
func EncodeRaw(raw *Raw) ([]byte, error) {
	cborH := &codec.CborHandle{}
	var b bytes.Buffer
	enc := codec.NewEncoder(&b, cborH)
	err := enc.Encode(raw)
	return b.Bytes(), err
}

// hashableBytes exists just to allow []byte implements hash.Writer
type hashableBytes []byte

func (b hashableBytes) WriteHash(w io.Writer) {
	_, err := w.Write(b)
	if err != nil {
		panic(err)
	}
}

// Hash generates hash for Raw record.
func (raw *Raw) Hash() []byte {
	return hash.SHA3hash224(raw.Type, hashableBytes(raw.Data))
}

// ToRecord decodes Raw to Record.
func (raw *Raw) ToRecord() Record {
	cborH := &codec.CborHandle{}
	rec := getRecordByTypeID(raw.Type)
	dec := codec.NewDecoder(bytes.NewReader(raw.Data), cborH)
	err := dec.Decode(rec)
	if err != nil {
		panic(err)
	}
	return rec
}

// Bytes2ID converts ID from byte representation to struct.
func Bytes2ID(b []byte) ID {
	return ID{
		Pulse: core.Bytes2PulseNumber(b[:core.PulseNumberSize]),
		Hash:  b[core.PulseNumberSize:],
	}
}

// Core2Reference converts commonly used reference to Ledger-specific.
func Core2Reference(cRef core.RecordRef) Reference {
	return Reference{
		Record: Bytes2ID(cRef[:core.RecordIDSize]),
		Domain: Bytes2ID(cRef[core.RecordIDSize:]),
	}
}

// ID2Bytes converts ID struct to it's byte representation.
func ID2Bytes(id ID) []byte {
	var buf = make([]byte, core.RecordIDSize)
	copy(buf[:core.PulseNumberSize], id.Pulse.Bytes())
	copy(buf[core.PulseNumberSize:], id.Hash)
	return buf
}

// record type ids for record types
// in use mostly for hashing and deserialization
// (we don't use iota for clarity and predictable ids,
// not depended on definition order)
const (
	// request record ids
	callRequestRecordID TypeID = 1
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
	deactivationRecordID        TypeID = 22
	objectAmendRecordID         TypeID = 23
	statefulCallResultID        TypeID = 24
	statefulExceptionResultID   TypeID = 25
	enforcedObjectAmendRecordID TypeID = 26
	objectAppendRecordID        TypeID = 27
	typeRecordID                TypeID = 28
	childRecordID               TypeID = 29
)

// getRecordByTypeID returns Record interface with concrete record type under the hood.
// This is useful with deserialization cases.
func getRecordByTypeID(id TypeID) Record { // nolint: gocyclo
	switch id {
	// request records
	case callRequestRecordID:
		return &CallRequest{}
	case wipeOutRecordID:
		return &WipeOutRecord{}
	case readRecordResultID:
		return &ReadRecordResult{}
	case statelessCallResultID:
		return &StatelessCallResult{}
	case statelessExceptionResultID:
		return &StatelessExceptionResult{}
	case readObjectResultID:
		return &ReadObjectResult{}
	case specialResultID:
		return &SpecialResult{}
	case lockUnlockResultID:
		return &LockUnlockResult{}
	case rejectionResultID:
		return &RejectionResult{}
	case activationRecordID:
		return &ActivationRecord{}
	case classActivateRecordID:
		return &ClassActivateRecord{}
	case objectActivateRecordID:
		return &ObjectActivateRecord{}
	case codeRecordID:
		return &CodeRecord{}
	case amendRecordID:
		return &AmendRecord{}
	case classAmendRecordID:
		return &ClassAmendRecord{}
	case deactivationRecordID:
		return &DeactivationRecord{}
	case objectAmendRecordID:
		return &ObjectAmendRecord{}
	case statefulCallResultID:
		return &StatefulCallResult{}
	case statefulExceptionResultID:
		return &StatefulExceptionResult{}
	case enforcedObjectAmendRecordID:
		return &EnforcedObjectAmendRecord{}
	case objectAppendRecordID:
		return &ObjectAppendRecord{}
	case typeRecordID:
		return &TypeRecord{}
	case childRecordID:
		return &ChildRecord{}
	default:
		panic(fmt.Errorf("unknown record type id %v", id))
	}
}

// getRecordByTypeID returns record's TypeID based on concrete record type of Record interface.
func getTypeIDbyRecord(rec Record) TypeID { // nolint: gocyclo, megacheck
	switch v := rec.(type) {
	// request records
	case *CallRequest:
		return callRequestRecordID
	// result records
	case *ResultRecord:
		return resultRecordID
	case *WipeOutRecord:
		return wipeOutRecordID
	case *ReadRecordResult:
		return readRecordResultID
	case *StatelessCallResult:
		return statelessCallResultID
	case *StatelessExceptionResult:
		return statelessExceptionResultID
	case *ReadObjectResult:
		return readObjectResultID
	case *SpecialResult:
		return specialResultID
	case *LockUnlockResult:
		return lockUnlockResultID
	case *RejectionResult:
		return rejectionResultID
	case *ActivationRecord:
		return activationRecordID
	case *ClassActivateRecord:
		return classActivateRecordID
	case *ObjectActivateRecord:
		return objectActivateRecordID
	case *CodeRecord:
		return codeRecordID
	case *AmendRecord:
		return amendRecordID
	case *ClassAmendRecord:
		return classAmendRecordID
	case *DeactivationRecord:
		return deactivationRecordID
	case *ObjectAmendRecord:
		return objectAmendRecordID
	case *StatefulCallResult:
		return statefulCallResultID
	case *StatefulExceptionResult:
		return statefulExceptionResultID
	case *EnforcedObjectAmendRecord:
		return enforcedObjectAmendRecordID
	case *ObjectAppendRecord:
		return objectAppendRecordID
	case *TypeRecord:
		return typeRecordID
	case *ChildRecord:
		return childRecordID
	default:
		panic(fmt.Errorf("can't find record id by type %T", v))
	}
}

// Encode serializes record to CBOR.
func Encode(rec Record) ([]byte, error) {
	cborH := &codec.CborHandle{}
	var b bytes.Buffer
	enc := codec.NewEncoder(&b, cborH)
	err := enc.Encode(rec)
	return b.Bytes(), err
}

// MustEncode wraps Encode, panics on encoding errors.
func MustEncode(rec Record) []byte {
	b, err := Encode(rec)
	if err != nil {
		panic(err)
	}
	return b
}

// EncodeToRaw converts record to Raw record.
func EncodeToRaw(rec Record) (*Raw, error) {
	b, err := Encode(rec)
	if err != nil {
		panic(err)
	}
	return &Raw{
		Type: getTypeIDbyRecord(rec),
		Data: b,
	}, nil
}
