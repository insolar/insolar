//
// Copyright 2019 Insolar Technologies GmbH
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package insolar

import (
	"encoding/binary"

	"github.com/insolar/insolar/longbits"
	"github.com/insolar/insolar/reference"
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

func NewGlobalReference(local ID, base ID) *Reference {
	global := reference.NewGlobal(base, local)
	return &global
}

// NewReferenceFromBase58 deserializes reference from base58 encoded string
func NewReferenceFromBase58(input string) (*Reference, error) {
	global, err := reference.DefaultDecoder().Decode(input)
	if err != nil {
		return nil, err
	}
	return &global, nil
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

// NewIDFromBase58 deserializes ID from base58 encoded string
func NewIDFromBase58(input string) (*ID, error) {
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
