package api

import (
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/platformpolicy"
)

func MakeReason(pulse insolar.PulseNumber, data []byte) insolar.Reference {
	hasher := platformpolicy.NewPlatformCryptographyScheme().ReferenceHasher()
	reasonID := *insolar.NewID(pulse, hasher.Hash(data))
	return *insolar.NewRecordReference(reasonID)
}
