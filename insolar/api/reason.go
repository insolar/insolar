// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

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
