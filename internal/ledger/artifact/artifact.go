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
package artifact

import (
	"context"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/message"
	"github.com/insolar/insolar/ledger/storage"
	"github.com/insolar/insolar/ledger/storage/object"
)

type Contract struct {
	Name        string // ???
	Domain      insolar.Reference
	MachineType insolar.MachineType
	Binary      []byte
}

type Manager interface {
	RegisterPrototype(ctx context.Context, name string, domain insolar.Reference) (*insolar.Reference, error)
}

type Scope struct {
	PulseNumber                insolar.PulseNumber
	ObjectStorage              storage.ObjectStorage              `inject:""`
	PlatformCryptographyScheme insolar.PlatformCryptographyScheme `inject:""`
}

// NewScope creates new scope instance.
func NewScope(pn insolar.PulseNumber) *Scope { // nolint
	return &Scope{
		PulseNumber: pn,
	}
}

func (m *Scope) RegisterPrototype(ctx context.Context, name string, domain insolar.Reference) (*insolar.Reference, error) {
	parcel := &message.Parcel{
		Msg: &message.GenesisRequest{Name: name + "_proto"},
	}
	// RegisterRequest
	rec := &object.RequestRecord{
		Parcel:      message.ParcelToBytes(parcel),
		MessageHash: m.hashParcel(parcel),
		// TODO: figure out is it required or not?
		// Object:      *obj.Record(),
	}
	jetID := insolar.ZeroJetID
	protoID, err := m.ObjectStorage.SetRecord(
		ctx, insolar.ID(jetID), m.PulseNumber, rec)
	if err != nil {
		return nil, err
	}

	proto := insolar.NewReference(*domain.Domain(), *protoID)
	return proto, nil
}

func (m *Scope) hashParcel(parcel insolar.Parcel) []byte {
	return m.PlatformCryptographyScheme.IntegrityHasher().Hash(message.MustSerializeBytes(parcel.Message()))
}
