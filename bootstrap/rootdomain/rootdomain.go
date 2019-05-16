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
	"github.com/insolar/insolar/insolar/message"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/platformpolicy"
)

// Name is the constant name of root domain.
const Name = "rootdomain"

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
	parcel := &message.Parcel{
		Msg: &message.GenesisRequest{
			Name: Name,
		},
	}
	req := record.Request{
		Parcel:      message.ParcelToBytes(parcel),
		MessageHash: message.ParcelMessageHash(r.PCS, parcel),
		Object:      insolar.GenesisRecord.ID(),
	}
	virtRec := record.Wrap(req)
	hash := record.HashVirtual(r.PCS.ReferenceHasher(), virtRec)
	return *insolar.NewID(genesisPulse, hash)
}
