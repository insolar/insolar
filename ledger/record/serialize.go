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
	"time"

	"github.com/ugorji/go/codec"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/cryptohelpers/hash"
	"github.com/insolar/insolar/log"
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
	start := time.Now()
	cborH := &codec.CborHandle{}
	rec := getRecordByTypeID(raw.Type)
	dec := codec.NewDecoder(bytes.NewReader(raw.Data), cborH)
	err := dec.Decode(rec)
	since := time.Since(start)
	if err != nil {
		panic(err)
	}
	if raw.Type == codeRecordID {
		log.Debugf("ToRecord func in record/serialize: for TypeID %s, time inside - %s", raw.Type, since)
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
	rec := core.GenRecordID(id.Pulse, id.Hash)
	return rec[:]
}

// record type ids for record types
// in use mostly for hashing and deserialization
// (we don't use iota for clarity and predictable ids,
// not depended on definition order)
const (
	// meta
	childRecordID   TypeID = 10
	genesisRecordID TypeID = 11

	// request
	callRequestRecordID TypeID = 20

	// result
	classActivateRecordID  TypeID = 30
	objectActivateRecordID TypeID = 31
	codeRecordID           TypeID = 32
	classAmendRecordID     TypeID = 33
	deactivationRecordID   TypeID = 34
	objectAmendRecordID    TypeID = 35
	typeRecordID           TypeID = 36
)

// getRecordByTypeID returns Record interface with concrete record type under the hood.
// This is useful with deserialization cases.
func getRecordByTypeID(id TypeID) Record { // nolint: gocyclo
	switch id {
	// request records
	case callRequestRecordID:
		return &CallRequest{}
	case classActivateRecordID:
		return &ClassActivateRecord{}
	case objectActivateRecordID:
		return &ObjectActivateRecord{}
	case codeRecordID:
		return &CodeRecord{}
	case classAmendRecordID:
		return &ClassAmendRecord{}
	case deactivationRecordID:
		return &DeactivationRecord{}
	case objectAmendRecordID:
		return &ObjectAmendRecord{}
	case typeRecordID:
		return &TypeRecord{}
	case childRecordID:
		return &ChildRecord{}
	case genesisRecordID:
		return &GenesisRecord{}
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
	case *ClassActivateRecord:
		return classActivateRecordID
	case *ObjectActivateRecord:
		return objectActivateRecordID
	case *CodeRecord:
		return codeRecordID
	case *ClassAmendRecord:
		return classAmendRecordID
	case *DeactivationRecord:
		return deactivationRecordID
	case *ObjectAmendRecord:
		return objectAmendRecordID
	case *TypeRecord:
		return typeRecordID
	case *ChildRecord:
		return childRecordID
	case *GenesisRecord:
		return genesisRecordID
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

// MustEncodeToRaw wraps EncodeToRaw, panics on encoding errors.
func MustEncodeToRaw(rec Record) *Raw {
	raw, err := EncodeToRaw(rec)
	if err != nil {
		panic(err)
	}
	return raw
}
