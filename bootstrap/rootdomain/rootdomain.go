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

package rootdomain

import (
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/platformpolicy"
)

var genesisPulse = insolar.GenesisPulse.PulseNumber

type Record struct {
	PCS insolar.PlatformCryptographyScheme
}

var RootDomain = &Record{
	PCS: platformpolicy.NewPlatformCryptographyScheme(),
}

func (r Record) Ref() insolar.Reference {
	id := r.ID()
	return *insolar.NewReference(id, id)
}

func (r Record) ID() insolar.ID {
	req := record.Request{
		CallType: record.CTGenesis,
		Method: Name,
	}
	virtRec := record.Wrap(req)
	hash := record.HashVirtual(r.PCS.ReferenceHasher(), virtRec)
	return *insolar.NewID(genesisPulse, hash)
}

// GenesisRef returns reference for genesis records based on the root domain.
func GenesisRef(name string) insolar.Reference {
	pcs := platformpolicy.NewPlatformCryptographyScheme()
	parcel := &message.Parcel{
		Msg: &message.GenesisRequest{
			Name: name,
		},
	}
	vrec := record.Wrap(record.Request{
		Parcel:      message.ParcelToBytes(parcel),
		MessageHash: message.ParcelMessageHash(pcs, parcel),
		Object:      RootDomain.ID(),
	})
	id := insolar.NewID(insolar.FirstPulseNumber, record.HashVirtual(pcs.ReferenceHasher(), vrec))
	return *insolar.NewReference(RootDomain.ID(), *id)
}
