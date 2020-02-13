// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

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
