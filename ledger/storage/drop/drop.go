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

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/ledger/storage/jet"
	"github.com/ugorji/go/codec"
)

//go:generate minimock -i github.com/insolar/insolar/ledger/storage/drop.Modifier -o ./ -s _mock.go

// Modifier provides an interface for modifying jetdrops.
type Modifier interface {
	Set(ctx context.Context, drop jet.Drop) error
}

//go:generate minimock -i github.com/insolar/insolar/ledger/storage/drop.Accessor -o ./ -s _mock.go

// Accessor provides an interface for accessing jetdrops.
type Accessor interface {
	ForPulse(ctx context.Context, jetID core.JetID, pulse core.PulseNumber) (jet.Drop, error)
}

//go:generate minimock -i github.com/insolar/insolar/ledger/storage/drop.Cleaner -o ./ -s _mock.go

// Cleaner provides an interface for removing jetdrops from a storage.
type Cleaner interface {
	Delete(pulse core.PulseNumber)
}

// Serialize serializes a drop
func Serialize(dr jet.Drop) []byte {
	buff := bytes.NewBuffer(nil)
	enc := codec.NewEncoder(buff, &codec.CborHandle{})
	enc.MustEncode(dr)
	return buff.Bytes()
}

// Deserialize deserializes a jet.Drop
func Deserialize(buf []byte) (dr jet.Drop) {
	dec := codec.NewDecoderBytes(buf, &codec.CborHandle{})
	dec.MustDecode(&dr)
	return dr
}
