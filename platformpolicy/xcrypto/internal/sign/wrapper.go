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

package sign

import (
	"crypto/rand"

	"github.com/pkg/errors"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/log"
	"github.com/insolar/x-crypto/ecdsa"
)

type ecdsaSignerWrapper struct {
	privateKey *ecdsa.PrivateKey
	hasher     insolar.Hasher
}

func (sw *ecdsaSignerWrapper) Sign(data []byte) (*insolar.Signature, error) {
	hash := sw.hasher.Hash(data)

	r, s, err := ecdsa.Sign(rand.Reader, sw.privateKey, hash)
	if err != nil {
		return nil, errors.Wrap(err, "[ Sign ] could't sign data")
	}

	ecdsaSignature := SerializeTwoBigInt(r, s)
	if err != nil {
		return nil, errors.Wrap(err, "[ Sign ] could't sign data")
	}

	signature := insolar.SignatureFromBytes(ecdsaSignature)
	return &signature, nil
}

type ecdsaVerifyWrapper struct {
	publicKey *ecdsa.PublicKey
	hasher    insolar.Hasher
}

func (sw *ecdsaVerifyWrapper) Verify(signature insolar.Signature, data []byte) bool {
	if signature.Bytes() == nil {
		return false
	}
	r, s, err := DeserializeTwoBigInt(signature.Bytes())
	if err != nil {
		log.Error(err)
		return false
	}

	hash := sw.hasher.Hash(data)
	return ecdsa.Verify(sw.publicKey, hash, r, s)
}
