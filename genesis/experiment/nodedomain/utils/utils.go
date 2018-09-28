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
	"crypto/sha256"
	"encoding/asn1"
	"math/big"

	"github.com/jbenet/go-base58"
	"github.com/pkg/errors"
)

func SerializePublicKey(key ecdsa.PublicKey) (string, error) {

	escdaPair := EcdsaPair{First: key.X, Second: key.Y}
	serKey, err := asn1.Marshal(escdaPair)
	if err != nil {
		return "", errors.Wrap(err, "[ SerializePublicKey ]")
	}

	return base58.Encode(serKey), nil
}

func MakePublicKey(pair EcdsaPair) *ecdsa.PublicKey {
	return &ecdsa.PublicKey{
		Curve: GetCurve(),
		X:     pair.First,
		Y:     pair.Second,
	}
}

func DeserializePublicKey(serPubKey string) (*ecdsa.PublicKey, error) {
	rawPubKey := base58.Decode(serPubKey)
	pk := &EcdsaPair{}
	rest, err := asn1.Unmarshal(rawPubKey, pk)

	if err != nil {
		return nil, errors.Wrap(err, "[ DeserializePublicKey ]")
	}
	if len(rest) != 0 {
		return nil, errors.New("[ DeserializePublicKey ] Rest exists. Wtf?")
	}

	return MakePublicKey(*pk), nil
}

func MakeHash(seed []byte) [sha256.Size]byte {
	return sha256.Sum256(seed)
}

func GetCurve() elliptic.Curve {
	return elliptic.P256()
}

type EcdsaPair struct {
	First  *big.Int
	Second *big.Int
}
