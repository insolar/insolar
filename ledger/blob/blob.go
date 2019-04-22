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

package blob

import (
	"bytes"
	"context"

	"github.com/insolar/insolar/insolar"
	"github.com/ugorji/go/codec"
)

//go:generate minimock -i github.com/insolar/insolar/ledger/blob.Accessor -o ./ -s _mock.go

// Accessor provides info about Blob-values from storage.
type Accessor interface {
	// ForID returns Blob for a provided id.
	ForID(ctx context.Context, id insolar.ID) (Blob, error)
}

//go:generate minimock -i github.com/insolar/insolar/ledger/blob.CollectionAccessor -o ./ -s _mock.go

// CollectionAccessor provides methods for querying blobs with specific search conditions.
type CollectionAccessor interface {
	// ForPulse returns []Blob for a provided jetID and a pulse number.
	ForPulse(ctx context.Context, jetID insolar.JetID, pn insolar.PulseNumber) []Blob
}

//go:generate minimock -i github.com/insolar/insolar/ledger/blob.Modifier -o ./ -s _mock.go

// Modifier provides methods for setting Blob-values to storage.
type Modifier interface {
	// Set saves new Blob-value in storage.
	Set(ctx context.Context, id insolar.ID, blob Blob) error
}

//go:generate minimock -i github.com/insolar/insolar/ledger/blob.Storage -o ./ -s _mock.go

// Storage is a union of Accessor and Modifier
type Storage interface {
	Accessor
	Modifier
}

//go:generate minimock -i github.com/insolar/insolar/ledger/blob.Cleaner -o ./ -s _mock.go

// Cleaner provides an interface for removing blobs from a storage.
type Cleaner interface {
	DeleteForPN(ctx context.Context, pulse insolar.PulseNumber)
}

// Blob represents blob-value with jetID.
type Blob struct {
	Value []byte
	JetID insolar.JetID
}

// Clone returns copy of argument blob.
func Clone(in Blob) (out Blob) {
	out.JetID = in.JetID
	if in.Value != nil {
		v := make([]byte, len(in.Value))
		copy(v, in.Value)
		out.Value = v
	}
	return
}

// MustEncode serializes a blob.
func MustEncode(blob *Blob) []byte {
	var buf bytes.Buffer
	enc := codec.NewEncoder(&buf, &codec.CborHandle{})
	err := enc.Encode(blob)
	if err != nil {
		panic(err)
	}
	return buf.Bytes()
}

// Decode deserializes a blob.
func Decode(buf []byte) (*Blob, error) {
	dec := codec.NewDecoder(bytes.NewReader(buf), &codec.CborHandle{})
	var blob Blob
	err := dec.Decode(&blob)
	if err != nil {
		return nil, err
	}
	return &blob, nil
}
