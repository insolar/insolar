package sign

import (
	"crypto"

	"github.com/insolar/insolar/insolar"
)

type AlgorithmProvider interface {
	DataSigner(crypto.PrivateKey, insolar.Hasher) insolar.Signer
	DigestSigner(crypto.PrivateKey) insolar.Signer
	DataVerifier(crypto.PublicKey, insolar.Hasher) insolar.Verifier
	DigestVerifier(crypto.PublicKey) insolar.Verifier
}
