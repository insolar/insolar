// Copyright 2020 Insolar Network Ltd.
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

package genesisrefs

import (
	"sync"

	"github.com/insolar/insolar/application"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/platformpolicy"
)

var genesisPulse = insolar.GenesisPulse.PulseNumber

// Record provides methods to calculate root domain's identifiers.
type Record struct {
	once                sync.Once
	rootDomainID        insolar.ID
	rootDomainReference insolar.Reference
	PCS                 insolar.PlatformCryptographyScheme
}

// RootDomain is the root domain instance.
var RootDomain = &Record{
	PCS: platformpolicy.NewPlatformCryptographyScheme(),
}

func (r *Record) initialize() {
	rootRecord := record.IncomingRequest{
		CallType: record.CTGenesis,
		Method:   application.GenesisNameRootDomain,
	}
	virtualRec := record.Wrap(&rootRecord)
	hash := record.HashVirtual(r.PCS.ReferenceHasher(), virtualRec)

	r.rootDomainID = *insolar.NewID(genesisPulse, hash)
	r.rootDomainReference = *insolar.NewReference(r.rootDomainID)
}

// ID returns insolar.ID  to root domain object.
func (r *Record) ID() insolar.ID {
	r.once.Do(r.initialize)

	return r.rootDomainID
}

// Reference returns insolar.Reference to root domain object
func (r *Record) Reference() insolar.Reference {
	r.once.Do(r.initialize)

	return r.rootDomainReference
}
