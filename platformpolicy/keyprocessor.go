package platformpolicy

import (
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/platformpolicy/commoncrypto"
)

const (
	COMMON AlgorithmType = "COMMON"
)

type AlgorithmType string

func NewKeyProcessor() insolar.KeyProcessor {
	return newKeyProcessor(COMMON)
}

func newKeyProcessor(algorithmType AlgorithmType) insolar.KeyProcessor {
	return commoncrypto.NewKeyProcessor()
}
