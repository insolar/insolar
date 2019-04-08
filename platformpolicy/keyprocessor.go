package platformpolicy

import (
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/platformpolicy/commoncrypto"
	"github.com/insolar/insolar/platformpolicy/customcrypto/secp256k1"
)

const (
	COMMON    AlgorithmType = "COMMON"
	SECP256K1 AlgorithmType = "SECP256K1"
)

type AlgorithmType string

func NewKeyProcessor() insolar.KeyProcessor {
	return NewKeyProcessor1(SECP256K1)
}

func NewKeyProcessor1(algorithmType AlgorithmType) insolar.KeyProcessor {
	switch algorithmType {
	case SECP256K1:
		return secp256k1.NewKeyProcessor()
	default:
		return commoncrypto.NewKeyProcessor()
	}
}
