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

package genesis

import (
	"github.com/insolar/insolar/bootstrap/rootdomain"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/message"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/platformpolicy"
)

func refByName(name string) insolar.Reference {
	pcs := platformpolicy.NewPlatformCryptographyScheme()
	parcel := &message.Parcel{
		Msg: &message.GenesisRequest{
			Name: name,
		},
	}
	vrec := record.Wrap(record.Request{
		Parcel:      message.ParcelToBytes(parcel),
		MessageHash: message.ParcelHash(pcs, parcel),
		Object:      rootdomain.RootDomain.ID(),
	})
	id := insolar.NewID(insolar.FirstPulseNumber, record.HashVirtual(pcs.ReferenceHasher(), vrec))
	return *insolar.NewReference(rootdomain.RootDomain.ID(), *id)
}
