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

package sign

import (
	"crypto"

	"github.com/insolar/insolar/insolar"
)

type ecdsaProvider struct {
}

func NewECDSAProvider() AlgorithmProvider {
	return &ecdsaProvider{}
}

func (p *ecdsaProvider) DataSigner(privateKey crypto.PrivateKey, hasher insolar.Hasher) insolar.Signer {
	return &ecdsaDataSignerWrapper{
		ecdsaDigestSignerWrapper: ecdsaDigestSignerWrapper{
			privateKey: MustConvertPrivateKeyToEcdsa(privateKey),
		},
		hasher: hasher,
	}
}
func (p *ecdsaProvider) DigestSigner(privateKey crypto.PrivateKey) insolar.Signer {
	return &ecdsaDigestSignerWrapper{
		privateKey: MustConvertPrivateKeyToEcdsa(privateKey),
	}
}

func (p *ecdsaProvider) DataVerifier(publicKey crypto.PublicKey, hasher insolar.Hasher) insolar.Verifier {
	return &ecdsaDataVerifyWrapper{
		ecdsaDigestVerifyWrapper: ecdsaDigestVerifyWrapper{
			publicKey: MustConvertPublicKeyToEcdsa(publicKey),
		},
		hasher: hasher,
	}
}

func (p *ecdsaProvider) DigestVerifier(publicKey crypto.PublicKey) insolar.Verifier {
	return &ecdsaDigestVerifyWrapper{
		publicKey: MustConvertPublicKeyToEcdsa(publicKey),
	}
}
