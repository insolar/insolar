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

package replica

import (
	"context"
	"crypto"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/heavy/sequence"
)

//go:generate minimock -i github.com/insolar/insolar/ledger/heavy/replica.Integrity -o ./ -s _mock.go

// Integrity is an interface to add and validate transferred record integrity.
type Integrity interface {
	// Wrap returns the serialized sequence of db records signed by node signature.
	Wrap([]sequence.Item) []byte
	// UnwrapAndValidate returns the deserialized sequence of db records if signature is valid or returns empty slice.
	UnwrapAndValidate([]byte) []sequence.Item
}

func NewIntegrity(crypto insolar.CryptographyService, parentPubKey crypto.PublicKey) Integrity {
	return &wrapper{crypto, parentPubKey}
}

type wrapper struct {
	crypto       insolar.CryptographyService
	parentPubKey crypto.PublicKey
}

func (w *wrapper) Wrap(items []sequence.Item) []byte {
	data, err := insolar.Serialize(items)
	if err != nil {
		inslogger.FromContext(context.Background()).Errorf("failed to serialize sequence items")
		return []byte{}
	}
	signature, err := w.crypto.Sign(data)
	if err != nil {
		inslogger.FromContext(context.Background()).Errorf("failed to sign sequence items")
		return []byte{}
	}
	pack := Packet{data, signature.Bytes()}
	packet, err := insolar.Serialize(&pack)
	if err != nil {
		inslogger.FromContext(context.Background()).Errorf("failed to serialize wrapped packet")
		return []byte{}
	}
	return packet
}

func (w *wrapper) UnwrapAndValidate(rawPacket []byte) []sequence.Item {
	var pack Packet
	err := insolar.Deserialize(rawPacket, &pack)
	if err != nil {
		inslogger.FromContext(context.Background()).Errorf("failed to deserialize wrapped packet")
		return []sequence.Item{}
	}
	if !w.crypto.Verify(w.parentPubKey, insolar.SignatureFromBytes(pack.Signature), pack.Data) {
		inslogger.FromContext(context.Background()).Errorf("invalid packet signature")
		return []sequence.Item{}
	}
	var seq []sequence.Item
	err = insolar.Deserialize(pack.Data, &seq)
	if err != nil {
		inslogger.FromContext(context.Background()).Errorf("failed to deserialize sequence items")
		return []sequence.Item{}
	}
	return seq
}

type Packet struct {
	Data      []byte
	Signature []byte
}
