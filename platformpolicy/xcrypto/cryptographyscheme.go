//
// Copyright 2019 Insolar Technologies GmbH
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
//

package xcrypto

import (
	"github.com/insolar/insolar/component"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/platformpolicy/algorithmprovider"
	"github.com/insolar/insolar/platformpolicy/keys"
	"github.com/insolar/insolar/platformpolicy/xcrypto/internal/hash"
	"github.com/insolar/insolar/platformpolicy/xcrypto/internal/sign"
)

type platformCryptographyScheme struct {
	HashProvider algorithmprovider.HashAlgorithmProvider `inject:""`
	SignProvider algorithmprovider.SignAlgorithmProvider `inject:""`
}

func (pcs *platformCryptographyScheme) PublicKeySize() int {
	return sign.TwoBigIntBytesLength
}

func (pcs *platformCryptographyScheme) SignatureSIze() int {
	return sign.TwoBigIntBytesLength
}

func (pcs *platformCryptographyScheme) ReferenceHasher() insolar.Hasher {
	return pcs.HashProvider.Hash224bits()
}

func (pcs *platformCryptographyScheme) IntegrityHasher() insolar.Hasher {
	return pcs.HashProvider.Hash512bits()
}

func (pcs *platformCryptographyScheme) Signer(privateKey keys.PrivateKey) insolar.Signer {
	return pcs.SignProvider.Sign(privateKey)
}

func (pcs *platformCryptographyScheme) Verifier(publicKey keys.PublicKey) insolar.Verifier {
	return pcs.SignProvider.Verify(publicKey)
}

func NewPlatformCryptographyScheme() insolar.PlatformCryptographyScheme {
	platformCryptographyScheme := &platformCryptographyScheme{}

	manager := component.Manager{}
	manager.Inject(
		platformCryptographyScheme,

		hash.NewSHA3Provider(),
		sign.NewECDSAProvider(),
	)

	return platformCryptographyScheme
}
