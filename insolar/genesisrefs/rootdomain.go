package genesisrefs

import (
	"sync"

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
		Method:   insolar.GenesisNameRootDomain,
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
