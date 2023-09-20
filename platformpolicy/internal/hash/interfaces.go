package hash

import (
	"github.com/insolar/insolar/insolar"
)

type AlgorithmProvider interface {
	Hash224bits() insolar.Hasher
	Hash256bits() insolar.Hasher
	Hash512bits() insolar.Hasher
}
