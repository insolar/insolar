// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package platformpolicy

import (
	"crypto"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/platformpolicy/internal/hash"
	"github.com/insolar/insolar/platformpolicy/internal/sign"
)

type platformCryptographyScheme struct {
	hashProvider hash.AlgorithmProvider
	signProvider sign.AlgorithmProvider
}

func (pcs *platformCryptographyScheme) PublicKeySize() int {
	return sign.TwoBigIntBytesLength
}

func (pcs *platformCryptographyScheme) SignatureSize() int {
	return sign.TwoBigIntBytesLength
}

func (pcs *platformCryptographyScheme) ReferenceHashSize() int {
	return pcs.hashProvider.Hash224bits().Size()
}

func (pcs *platformCryptographyScheme) IntegrityHashSize() int {
	return pcs.hashProvider.Hash512bits().Size()
}

func (pcs *platformCryptographyScheme) ReferenceHasher() insolar.Hasher {
	return pcs.hashProvider.Hash224bits()
}

func (pcs *platformCryptographyScheme) IntegrityHasher() insolar.Hasher {
	return pcs.hashProvider.Hash512bits()
}

func (pcs *platformCryptographyScheme) DataSigner(privateKey crypto.PrivateKey, hasher insolar.Hasher) insolar.Signer {
	return pcs.signProvider.DataSigner(privateKey, hasher)
}

func (pcs *platformCryptographyScheme) DigestSigner(privateKey crypto.PrivateKey) insolar.Signer {
	return pcs.signProvider.DigestSigner(privateKey)
}

func (pcs *platformCryptographyScheme) DataVerifier(publicKey crypto.PublicKey, hasher insolar.Hasher) insolar.Verifier {
	return pcs.signProvider.DataVerifier(publicKey, hasher)
}

func (pcs *platformCryptographyScheme) DigestVerifier(publicKey crypto.PublicKey) insolar.Verifier {
	return pcs.signProvider.DigestVerifier(publicKey)
}

func NewPlatformCryptographyScheme() insolar.PlatformCryptographyScheme {
	return &platformCryptographyScheme{
		hashProvider: hash.NewSHA3Provider(),
		signProvider: sign.NewECDSAProvider(),
	}
}
