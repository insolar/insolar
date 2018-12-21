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
	"encoding/binary"
	"fmt"

	"github.com/insolar/insolar/core"
	"github.com/ugorji/go/codec"
)

// record type ids for record types
// in use mostly for hashing and deserialization
// (we don't use iota for clarity and predictable ids,
// not depended on definition order)
//go:generate stringer -type=TypeID
const (
	// meta
	typeGenesis TypeID = 10
	typeChild   TypeID = 11
	typeJet     TypeID = 12

	// request
	typeCallRequest TypeID = 20

	// result
	typeResult     TypeID = 30
	typeType       TypeID = 31
	typeCode       TypeID = 32
	typeActivate   TypeID = 33
	typeAmend      TypeID = 34
	typeDeactivate TypeID = 35
)

// getRecordByTypeID returns Record interface with concrete record type under the hood.
// This is useful with deserialization cases.
func getRecordByTypeID(id TypeID) Record { // nolint: gocyclo
	switch id {
	// request records
	case typeCallRequest:
		return &RequestRecord{}
	case typeActivate:
		return &ObjectActivateRecord{}
	case typeCode:
		return &CodeRecord{}
	case typeDeactivate:
		return &DeactivationRecord{}
	case typeAmend:
		return &ObjectAmendRecord{}
	case typeType:
		return &TypeRecord{}
	case typeChild:
		return &ChildRecord{}
	case typeGenesis:
		return &GenesisRecord{}
	case typeResult:
		return &ResultRecord{}
	case typeJet:
		return &JetRecord{}
	default:
		panic(fmt.Errorf("unknown record type id %v", id))
	}
}

// SerializeType returns binary representation of provided type.
func SerializeType(id TypeID) []byte {
	buf := make([]byte, TypeIDSize)
	binary.BigEndian.PutUint32(buf, uint32(id))
	return buf
}

// DeserializeType returns type from provided binary representation.
func DeserializeType(buf []byte) TypeID {
	return TypeID(binary.BigEndian.Uint32(buf))
}

// SerializeRecord returns binary representation of provided record.
func SerializeRecord(rec Record) []byte {
	typeBytes := SerializeType(rec.Type())
	buff := bytes.NewBuffer(typeBytes)
	enc := codec.NewEncoder(buff, &codec.CborHandle{})
	enc.MustEncode(rec)
	return buff.Bytes()
}

// DeserializeRecord returns record decoded from bytes.
func DeserializeRecord(buf []byte) Record {
	t := DeserializeType(buf[:TypeIDSize])
	dec := codec.NewDecoderBytes(buf[TypeIDSize:], &codec.CborHandle{})
	rec := getRecordByTypeID(t)
	dec.MustDecode(&rec)
	return rec
}

// CalculateIDForBlob calculate id for blob with using current pulse number
func CalculateIDForBlob(scheme core.PlatformCryptographyScheme, pulseNumber core.PulseNumber, blob []byte) *core.RecordID {
	hasher := scheme.IntegrityHasher()
	_, err := hasher.Write(blob)
	if err != nil {
		panic(err)
	}
	return core.NewRecordID(pulseNumber, hasher.Sum(nil))
}
