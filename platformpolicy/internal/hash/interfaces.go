// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package hash

import (
	"github.com/insolar/insolar/insolar"
)

type AlgorithmProvider interface {
	Hash224bits() insolar.Hasher
	Hash256bits() insolar.Hasher
	Hash512bits() insolar.Hasher
}
