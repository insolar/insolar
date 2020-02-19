// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package insolar

import (
	"encoding/binary"

	"github.com/insolar/insolar/longbits"
	"github.com/insolar/insolar/reference"
	"github.com/pkg/errors"
)

const (
	// RecordHashSize is a record hash size. We use 224-bit SHA-3 hash (28 bytes).
	RecordHashSize = 28
	// RecordIDSize is relative record address.
	RecordIDSize = PulseNumberSize + RecordHashSize
	// RecordHashOffset is a offset where hash bytes starts in ID.
	RecordHashOffset = PulseNumberSize
	// RecordRefSize is absolute records address (including domain ID).
	RecordRefSize = RecordIDSize * 2
	// RecordRefIDSeparator is character that separates ID from DomainID in serialized Reference.
	RecordRefIDSeparator = "."
)

type (
	// ID is a unified record ID
	ID = reference.Local
	// Reference is a unified record reference
	Reference = reference.Global
)

// NewReference returns Reference composed from domain and record.
func NewReference(id ID) *Reference {
	global := reference.NewSelfRef(id)
	return &global
}

func NewRecordReference(local ID) *Reference {
	global := reference.NewRecordRef(local)
	return &global
}

func NewGlobalReference(local ID, base ID) *Reference {
	global := reference.NewGlobal(base, local)
	return &global
}

// NewObjectReferenceFromString deserializes reference from base64 encoded string and checks if it object reference
func NewObjectReferenceFromString(input string) (*Reference, error) {
	global, err := NewReferenceFromString(input)
	if err != nil {
		return nil, err
	}
	if !global.IsObjectReference() {
		return nil, errors.New("provided reference is not object")
	}
	if !global.IsSelfScope() {
		return nil, errors.New("provided reference is not self-scoped")
	}
	return global, nil
}

// NewRecordReferenceFromString deserializes reference from base64 encoded string and checks if it record reference
func NewRecordReferenceFromString(input string) (*Reference, error) {
	global, err := NewReferenceFromString(input)
	if err != nil {
		return nil, err
	}
	if !global.IsRecordScope() {
		return nil, errors.New("provided reference is not record")
	}
	return global, nil
}

// NewReferenceFromString deserializes reference from base64 encoded string
func NewReferenceFromString(input string) (*Reference, error) {
	global, err := reference.DefaultDecoder().Decode(input)
	if err != nil {
		return nil, err
	}
	return &global, nil
}

// IsObjectReferenceString checks the validity of the reference
func IsObjectReferenceString(input string) bool {
	_, err := NewObjectReferenceFromString(input)
	return err == nil
}

// NewReferenceFromBytes : After CBOR Marshal/Unmarshal Ref can be converted to byte slice, this converts it back
func NewReferenceFromBytes(byteReference []byte) *Reference {
	g := reference.Global{}
	if err := g.Unmarshal(byteReference); err != nil {
		return nil
	}
	return &g
}

// NewEmptyReference returns empty Reference.
func NewEmptyReference() *Reference {
	return &Reference{}
}

// NewID generates ID byte representation
func NewID(p PulseNumber, hash []byte) *ID {
	hashB := longbits.Bits224{}
	copy(hashB[:], hash)

	local := reference.NewLocal(p, 0, hashB)
	return &local
}

// NewIDFromString deserializes ID from base64 encoded string
func NewIDFromString(input string) (*ID, error) {
	global, err := reference.DefaultDecoder().Decode(input)
	if err != nil {
		return nil, err
	}
	return global.GetLocal(), nil
}

// NewIDFromBytes converts byte slice to ID
func NewIDFromBytes(hash []byte) *ID {
	if hash == nil {
		return NewEmptyID()
	}
	pn := PulseNumber(binary.BigEndian.Uint32(hash[:reference.LocalBinaryPulseAndScopeSize]))
	return NewID(pn, hash[reference.LocalBinaryPulseAndScopeSize:])
}

func NewEmptyID() *ID {
	return &ID{}
}
