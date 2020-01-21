// Copyright 2020 Insolar Network Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
