/*
 *    Copyright 2019 Insolar Technologies
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

package drop

import (
	"bytes"
	"context"

	"github.com/ugorji/go/codec"

	"github.com/insolar/insolar/core"
)

// Modifier provides an interface for modifying jetdrops.
//go:generate minimock -i github.com/insolar/insolar/ledger/storage/drop.Modifier -o ./ -s _mock.go
type Modifier interface {
	Set(ctx context.Context, jetID core.JetID, drop Drop) error
}

// Accessor provides an interface for accessing jetdrops.
//go:generate minimock -i github.com/insolar/insolar/ledger/storage/drop.Accessor -o ./ -s _mock.go
type Accessor interface {
	ForPulse(ctx context.Context, jetID core.JetID, pulse core.PulseNumber) (Drop, error)
}

// Drop is a blockchain block.
// It contains hashes of the current block and the previous one.
type Drop struct {
	// nolint: golint
	// Pulse number (probably we should save it too).
	Pulse core.PulseNumber

	// PrevHash is a hash of all record hashes belongs to previous pulse.
	PrevHash []byte

	// Hash is a hash of all record hashes belongs to one pulse and previous drop hash.
	Hash []byte

	Size uint64
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
