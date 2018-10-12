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
	"encoding/gob"
	"errors"
	"math/big"
	"sort"

	ecdsahelper "github.com/insolar/insolar/cryptohelpers/ecdsa"
	"golang.org/x/crypto/sha3"

	"github.com/insolar/insolar/core"
)

func checkPayloadSignature(request *Payload) (bool, error) {
	return checkSignature(request.Body, request.PublicKey, request.Signature)
}

func checkSignature(data interface{}, pub string, signature []byte) (bool, error) {
	b := bytes.Buffer{}
	e := gob.NewEncoder(&b)
	err := e.Encode(data)
	if err != nil {
		return false, err
	}

	r := big.Int{}
	s := big.Int{}
	sigLen := len(signature)
	r.SetBytes(signature[:(sigLen / 2)])
	s.SetBytes(signature[(sigLen / 2):])

	h := sha3.New256()
	_, err = h.Write(b.Bytes())
	if err != nil {
		return false, err
	}
	calculatedHash := h.Sum(nil)
	publicKey, err := ecdsahelper.ImportPublicKey(pub)
	if err != nil {
		return false, err
	}

	return ecdsa.Verify(publicKey, calculatedHash, &r, &s), nil
}

func singData(privateKey *ecdsa.PrivateKey, data interface{}) ([]byte, error) {
	b := bytes.Buffer{}
	e := gob.NewEncoder(&b)
	err := e.Encode(data)
	if err != nil {
		return nil, err
	}

	h := sha3.New256()
	_, err = h.Write(b.Bytes())
	if err != nil {
		return nil, err
	}
	calculatedHash := h.Sum(nil)
	r, s, err := ecdsa.Sign(rand.Reader, privateKey, calculatedHash)
	if err != nil {
		return nil, err
	}

	return append(r.Bytes(), s.Bytes()...), nil
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
		h := sha3.New256()
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
