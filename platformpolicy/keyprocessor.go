package platformpolicy

import (
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/platformpolicy/commoncrypto"
)

const (
	COMMON    AlgorithmType = "COMMON"
	SECP256K1 AlgorithmType = "SECP256K1"
)

type AlgorithmType string

func NewKeyProcessor() insolar.KeyProcessor {
	return newKeyProcessor(SECP256K1)
}

func newKeyProcessor(algorithmType AlgorithmType) insolar.KeyProcessor {
	return commoncrypto.NewKeyProcessor()
}
