/*
 *    Copyright 2018 Insolar
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package sign

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/asn1"
	"math/big"

	"github.com/insolar/insolar/core"
	"github.com/pkg/errors"
)

type signature struct {
	R, S *big.Int
}

func fromRS(r, s *big.Int) *signature {
	return &signature{R: r, S: s}
}

func (p *signature) Marshal() ([]byte, error) {
	signature, err := asn1.Marshal(p)
	if err != nil {
		return nil, errors.Wrap(err, "[ Marshall ] Could't marshal signature")
	}
	return signature, nil
}

func (p *signature) Unmarshal(signatureRaw []byte) error {
	rest, err := asn1.Unmarshal(signatureRaw, p)
	if len(rest) != 0 {
		return errors.New("[ Unmarshal ] len of rest must be 0")
	}
	if err != nil {
		return errors.Wrap(err, "[ Unmarshal ] Could't unmarshal signature")
	}
	return nil
}

type ecdsaSignerWrapper struct {
	privateKey crypto.PrivateKey
	hasher     core.Hasher
}

func (sw *ecdsaSignerWrapper) Sign(data []byte) ([]byte, error) {
	hash := sw.hasher.Hash(data)
	privateKey := sw.privateKey.(*ecdsa.PrivateKey)

	r, s, err := ecdsa.Sign(rand.Reader, privateKey, hash)
	if err != nil {
		return nil, errors.Wrap(err, "[ Sign ] could't sign data")
	}

	signature, err := fromRS(r, s).Marshal()
	if err != nil {
		return nil, errors.Wrap(err, "[ Sign ] could't sign data")
	}

	return signature, nil
}

type ecdsaVerifyWrapper struct {
	publicKey crypto.PublicKey
	hasher    core.Hasher
}

func (sw *ecdsaVerifyWrapper) Verify(data, signatureRaw []byte) bool {
	var signature signature
	err := signature.Unmarshal(signatureRaw)
	if err != nil {
		return false
	}

	hash := sw.hasher.Hash(data)
	publicKey := sw.publicKey.(*ecdsa.PublicKey)

	return ecdsa.Verify(publicKey, hash, signature.R, signature.S)
}
