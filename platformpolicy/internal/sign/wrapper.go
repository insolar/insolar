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
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/asn1"
	"fmt"
	"math/big"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/log"
	"github.com/pkg/errors"
)

const bigIntLength = 32
const lenBytes = 2

type ecdsaSignature struct {
	R, S *big.Int
}

func (p ecdsaSignature) Marshal() ([]byte, error) {
	signature, err := asn1.Marshal(p)
	if err != nil {
		return nil, errors.Wrap(err, "[ Marshall ] Could't marshal ecdsaSignature")
	}
	return signature, nil
}

func (p *ecdsaSignature) Unmarshal(signatureRaw []byte) error {
	rest, err := asn1.Unmarshal(signatureRaw, p)
	if len(rest) != 0 {
		return errors.New("[ Unmarshal ] len of rest must be 0")
	}
	if err != nil {
		return errors.Wrap(err, "[ Unmarshal ] Could't unmarshal ecdsaSignature")
	}
	return nil
}

type ecdsaSignerWrapper struct {
	privateKey *ecdsa.PrivateKey
	hasher     core.Hasher
}

func (sw *ecdsaSignerWrapper) Sign(data []byte) (*core.Signature, error) {
	hash := sw.hasher.Hash(data)

	r, s, err := ecdsa.Sign(rand.Reader, sw.privateKey, hash)
	if err != nil {
		return nil, errors.Wrap(err, "[ Sign ] could't sign data")
	}

	ecdsaSignature := makeSignature(r, s)
	if err != nil {
		return nil, errors.Wrap(err, "[ Sign ] could't sign data")
	}

	signature := core.SignatureFromBytes(ecdsaSignature)
	return &signature, nil
}

type ecdsaVerifyWrapper struct {
	publicKey *ecdsa.PublicKey
	hasher    core.Hasher
}

func (sw *ecdsaVerifyWrapper) Verify(signature core.Signature, data []byte) bool {
	if signature.Bytes() == nil {
		return false
	}
	r, s, err := getRSFromBytes(signature.Bytes())
	if err != nil {
		log.Error(err)
		return false
	}
	ecdsaSignature := ecdsaSignature{r, s}
	hash := sw.hasher.Hash(data)

	return ecdsa.Verify(sw.publicKey, hash, ecdsaSignature.R, ecdsaSignature.S)
}

func makeSignature(r, s *big.Int) []byte {
	if (len(r.Bytes()) > bigIntLength) ||
		(len(s.Bytes()) > bigIntLength) {
		err := fmt.Sprintf("[ makeSignature ] wrong r, s length. r: %d; s: %d; needed: %d", len(r.Bytes()), len(s.Bytes()), bigIntLength)
		panic(err)
	}
	rLen := uint8(len(r.Bytes()))
	sLen := uint8(len(s.Bytes()))
	res := make([]byte, rLen+sLen+lenBytes)
	res[0] = rLen
	copy(res[1:rLen+lenBytes], r.Bytes())
	res[rLen+1] = sLen
	copy(res[rLen+lenBytes:], s.Bytes())
	return res[:]
}

func getRSFromBytes(data []byte) (*big.Int, *big.Int, error) {
	if len(data) > (bigIntLength*lenBytes + lenBytes) {
		err := fmt.Sprintf("[ getRSFromBytes ] wrong data length to get a r, s. recv len: %d", len(data))
		return nil, nil, errors.New(err)
	}
	r := new(big.Int)
	s := new(big.Int)
	rLen := data[0]
	if int(rLen+1) > len(data) {
		err := fmt.Sprintf("[ getRSFromBytes ] wrong data to parse r, s")
		return nil, nil, errors.New(err)
	}
	sLen := data[rLen+1]
	if int(rLen+sLen+lenBytes) != len(data) {
		err := fmt.Sprintf("[ getRSFromBytes ] wrong data to parse r, s")
		return nil, nil, errors.New(err)
	}
	rBytes := make([]byte, rLen)
	sBytes := make([]byte, sLen)
	copy(rBytes, data[1:rLen+lenBytes])
	copy(sBytes, data[rLen+lenBytes:])
	r.SetBytes(rBytes)
	s.SetBytes(sBytes)
	return r, s, nil
}
