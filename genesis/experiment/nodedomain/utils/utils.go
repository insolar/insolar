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

package utils

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/asn1"
	"math/big"

	"github.com/jbenet/go-base58"
	"github.com/pkg/errors"
)

// SerializePublicKey serializes pk string
func SerializePublicKey(key ecdsa.PublicKey) (string, error) {

	escdaPair := EcdsaPair{First: key.X, Second: key.Y}
	serKey, err := asn1.Marshal(escdaPair)
	if err != nil {
		return "", errors.Wrap(err, "[ SerializePublicKey ]")
	}

	return base58.Encode(serKey), nil
}

// MakePublicKey makes PK from EcdsaPair
func MakePublicKey(pair EcdsaPair) *ecdsa.PublicKey {
	return &ecdsa.PublicKey{
		Curve: GetCurve(),
		X:     pair.First,
		Y:     pair.Second,
	}
}

// DeserializePublicKey deserializes pk from string
func DeserializePublicKey(serPubKey string) (*ecdsa.PublicKey, error) {
	rawPubKey := base58.Decode(serPubKey)
	pk := &EcdsaPair{}
	rest, err := asn1.Unmarshal(rawPubKey, pk)

	if err != nil {
		return nil, errors.Wrap(err, "[ DeserializePublicKey ]")
	}
	if len(rest) != 0 {
		return nil, errors.New("[ DeserializePublicKey ] len of rest must be 0")
	}

	return MakePublicKey(*pk), nil
}

// MakeHash makes hash from seed
func MakeHash(seed []byte) [sha256.Size]byte {
	return sha256.Sum256(seed)
}

// GetCurve gets default curve
func GetCurve() elliptic.Curve {
	return elliptic.P256()
}

// EcdsaPair represents two ints for ecdsa
type EcdsaPair struct {
	First  *big.Int
	Second *big.Int
}

// Sign signs given seed
func Sign(seed []byte, key *ecdsa.PrivateKey) ([]byte, error) {

	hash := MakeHash(seed)

	r, s, err := ecdsa.Sign(rand.Reader, key, hash[:])

	if err != nil {
		return nil, errors.Wrap(err, "[ Sign ]")
	}

	data, err := asn1.Marshal(EcdsaPair{First: r, Second: s})
	if err != nil {
		return nil, errors.Wrap(err, "[ Sign ]")
	}

	return data, nil
}

// Verifies signature
func Verify(seed []byte, signatureRaw []byte, pubKey string) (bool, error) {
	var ecdsaPair EcdsaPair
	rest, err := asn1.Unmarshal(signatureRaw, &ecdsaPair)
	if err != nil {
		return false, errors.Wrap(err, "[ Verify ]")
	}
	if len(rest) != 0 {
		return false, errors.New("[ Verify ] len of  rest must be 0")
	}

	savedKey, err := DeserializePublicKey(pubKey)
	if err != nil {
		return false, errors.Wrap(err, "[ Verify ]")
	}

	hash := MakeHash(seed)

	return ecdsa.Verify(savedKey, hash[:], ecdsaPair.First, ecdsaPair.Second), nil
}

// MakeSeed makes random seed
func MakeSeed() []byte {
	seed := make([]byte, 32)
	_, err := rand.Read(seed)
	if err != nil {
		panic(err)
	}

	return seed
}
