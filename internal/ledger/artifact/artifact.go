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
	RegisterRequest(ctx context.Context, objectRef insolar.Reference, parcel insolar.Parcel) (*insolar.ID, error)
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

func (m *Scope) RegisterRequest(ctx context.Context, objectRef insolar.Reference, parcel insolar.Parcel) (*insolar.ID, error) {
	rec := &object.RequestRecord{
		Parcel:      message.ParcelToBytes(parcel),
		MessageHash: m.hashParcel(parcel),
		Object:      *objectRef.Record(),
	}
	return m.ObjectStorage.SetRecord(
		ctx,
		insolar.ID(insolar.ZeroJetID),
		m.PulseNumber,
		rec,
	)
}

func (m *Scope) hashParcel(parcel insolar.Parcel) []byte {
	return m.PlatformCryptographyScheme.IntegrityHasher().Hash(message.MustSerializeBytes(parcel.Message()))
}
