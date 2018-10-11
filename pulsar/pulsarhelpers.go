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

package pulsar

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/asn1"
	"math/big"
	"sort"

	ecdsa_helper "github.com/insolar/insolar/cryptohelpers/ecdsa"
	"github.com/pkg/errors"
	"github.com/ugorji/go/codec"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/ledger/hash"
)

type ecdsaSignature struct {
	R, S *big.Int
}

func checkPayloadSignature(request *Payload) (bool, error) {
	return checkSignature(request.Body, request.PublicKey, request.Signature)
}

func checkSignature(data interface{}, pub string, signature []byte) (bool, error) {
	cborH := &codec.CborHandle{}
	var b bytes.Buffer
	enc := codec.NewEncoder(&b, cborH)
	err := enc.Encode(data)
	if err != nil {
		return false, err
	}

	var ecdsaP ecdsaSignature
	rest, err := asn1.Unmarshal(signature, &ecdsaP)
	if err != nil {
		return false, errors.Wrap(err, "[ checkSignature ]")
	}
	if len(rest) != 0 {
		return false, errors.New("[ checkSignature ] len of  rest must be 0")
	}

	publicKey, err := ecdsa_helper.ImportPublicKey(pub)
	if err != nil {
		return false, err
	}

	return ecdsa.Verify(publicKey, b.Bytes(), ecdsaP.R, ecdsaP.S), nil
}

func signData(privateKey *ecdsa.PrivateKey, data interface{}) ([]byte, error) {
	cborH := &codec.CborHandle{}
	var b bytes.Buffer
	enc := codec.NewEncoder(&b, cborH)
	err := enc.Encode(data)
	if err != nil {
		return nil, err
	}
	return privateKey.Sign(rand.Reader, b.Bytes(), nil)
}

func selectByEntropy(entropy core.Entropy, values []string, count int) ([]string, error) { // nolint: megacheck
	type idxHash struct {
		idx  int
		hash []byte
	}

	if len(values) < count {
		return nil, errors.New("count value should be less than values size")
	}

	hashes := make([]*idxHash, 0, len(values))
	for i, value := range values {
		h := hash.NewSHA3()
		_, err := h.Write(entropy[:])
		if err != nil {
			return nil, err
		}
		_, err = h.Write([]byte(value))
		if err != nil {
			return nil, err
		}
		hashes = append(hashes, &idxHash{
			idx:  i,
			hash: h.Sum(nil),
		})
	}

	sort.SliceStable(hashes, func(i, j int) bool { return bytes.Compare(hashes[i].hash, hashes[j].hash) < 0 })

	selected := make([]string, 0, count)
	for i := 0; i < count; i++ {
		selected = append(selected, values[hashes[i].idx])
	}
	return selected, nil
}
