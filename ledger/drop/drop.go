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

package drop

import (
	"bytes"
	"context"

	"github.com/ugorji/go/codec"

	"github.com/insolar/insolar/insolar"
)

//go:generate minimock -i github.com/insolar/insolar/ledger/drop.Modifier -o ./ -s _mock.go

// Modifier provides an interface for modifying jetdrops.
type Modifier interface {
	Set(ctx context.Context, drop Drop) error
}

//go:generate minimock -i github.com/insolar/insolar/ledger/drop.Accessor -o ./ -s _mock.go

// Accessor provides an interface for accessing jetdrops.
type Accessor interface {
	ForPulse(ctx context.Context, jetID insolar.JetID, pulse insolar.PulseNumber) (Drop, error)
}

//go:generate minimock -i github.com/insolar/insolar/ledger/drop.Cleaner -o ./ -s _mock.go

// Cleaner provides an interface for removing jetdrops from a storage.
type Cleaner interface {
	DeleteForPN(ctx context.Context, pulse insolar.PulseNumber)
}

// Drop is a blockchain block.
// It contains hashes of the current block and the previous one.
type Drop struct {
	// Pulse number (probably we should save it too).
	Pulse insolar.PulseNumber

	// PrevHash is a hash of all record hashes belongs to previous pulse.
	PrevHash []byte

	// Hash is a hash of all record hashes belongs to one pulse and previous drop hash.
	Hash []byte

	// Size represents data about physical size of the current jet.Drop.
	Size uint64

	// JetID represents data about JetID of the current jet.Drop.
	JetID insolar.JetID
}

// MustEncode serializes jet drop.
func MustEncode(drop *Drop) []byte {
	var buf bytes.Buffer
	enc := codec.NewEncoder(&buf, &codec.CborHandle{})
	err := enc.Encode(drop)
	if err != nil {
		panic(err)
	}
	return buf.Bytes()
}

// Decode deserializes jet drop.
func Decode(buf []byte) (*Drop, error) {
	dec := codec.NewDecoder(bytes.NewReader(buf), &codec.CborHandle{})
	var drop Drop
	err := dec.Decode(&drop)
	if err != nil {
		return nil, err
	}
	return &drop, nil
}

// MustDecode deserializes a jet.Drop
func MustDecode(buf []byte) (dr Drop) {
	dec := codec.NewDecoderBytes(buf, &codec.CborHandle{})
	dec.MustDecode(&dr)
	return dr
}
