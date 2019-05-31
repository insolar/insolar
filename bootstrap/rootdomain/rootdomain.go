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

// Record provides methods to calculate root domain's identifiers.
type Record struct {
	PCS insolar.PlatformCryptographyScheme
}

// RootDomain is the root domain instance.
var RootDomain = &Record{
	PCS: platformpolicy.NewPlatformCryptographyScheme(),
}

// Ref returns insolar.Reference to root domain object.
func (r Record) Ref() insolar.Reference {
	return *insolar.NewReference(r.ID())
}

// ID returns insolar.ID  to root domain object.
func (r Record) ID() insolar.ID {
	req := record.Request{
		CallType: record.CTGenesis,
		Method:   insolar.GenesisNameRootDomain,
	}
	virtRec := record.Wrap(req)
	hash := record.HashVirtual(r.PCS.ReferenceHasher(), virtRec)
	return *insolar.NewID(genesisPulse, hash)
}

// GenesisRef returns reference to any genesis records based on the root domain.
func GenesisRef(name string) insolar.Reference {
	pcs := platformpolicy.NewPlatformCryptographyScheme()
	req := record.Request{
		CallType: record.CTGenesis,
		Method:   name,
	}
	virtRec := record.Wrap(req)
	hash := record.HashVirtual(pcs.ReferenceHasher(), virtRec)
	id := insolar.NewID(insolar.FirstPulseNumber, hash)
	return *insolar.NewReference(*id)
}
