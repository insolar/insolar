package api

import (
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/platformpolicy"
)

func MakeReason(pulse insolar.PulseNumber, data []byte) insolar.Reference {
	hasher := platformpolicy.NewPlatformCryptographyScheme().ReferenceHasher()
	return *insolar.NewReference(*insolar.NewID(pulse, hasher.Hash(data)))
}
