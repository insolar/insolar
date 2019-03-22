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

	"github.com/insolar/insolar/core"
)

//go:generate minimock -i github.com/insolar/insolar/ledger/storage/drop.Modifier -o ./ -s _mock.go

// Modifier provides an interface for modifying jetdrops.
type Modifier interface {
	Set(ctx context.Context, drop Drop) error
}

//go:generate minimock -i github.com/insolar/insolar/ledger/storage/drop.Accessor -o ./ -s _mock.go

// Accessor provides an interface for accessing jetdrops.
type Accessor interface {
	ForPulse(ctx context.Context, jetID core.JetID, pulse core.PulseNumber) (Drop, error)
}

// Drop is a blockchain block.
// It contains hashes of the current block and the previous one.
type Drop struct {
	// Pulse number (probably we should save it too).
	Pulse core.PulseNumber

	// PrevHash is a hash of all record hashes belongs to previous pulse.
	PrevHash []byte

	// Hash is a hash of all record hashes belongs to one pulse and previous drop hash.
	Hash []byte

	// Size represents data about physical size of the current jet.Drop.
	Size uint64

	// JetID represents data about JetID of the current jet.Drop.
	JetID core.JetID
}

// Encode serializes jet drop.
func Encode(drop *Drop) ([]byte, error) {
	var buf bytes.Buffer
	enc := codec.NewEncoder(&buf, &codec.CborHandle{})
	err := enc.Encode(drop)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
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

//go:generate minimock -i github.com/insolar/insolar/ledger/storage/drop.Cleaner -o ./ -s _mock.go

// Cleaner provides an interface for removing jetdrops from a storage.
type Cleaner interface {
	Delete(pulse core.PulseNumber)
}

// Serialize serializes a drop
func Serialize(dr Drop) []byte {
	buff := bytes.NewBuffer(nil)
	enc := codec.NewEncoder(buff, &codec.CborHandle{})
	enc.MustEncode(dr)
	return buff.Bytes()
}

// Deserialize deserializes a jet.Drop
func Deserialize(buf []byte) (dr Drop) {
	dec := codec.NewDecoderBytes(buf, &codec.CborHandle{})
	dec.MustDecode(&dr)
	return dr
}
