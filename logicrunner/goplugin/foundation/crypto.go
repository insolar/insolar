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

package foundation

import (
	"crypto/ecdsa"
	"github.com/go-jose"
	"github.com/insolar/x-crypto/x509"
)

// TODO: this file should be removed

// Verify verifies signature.
// TODO jose verify
func Verify(jwsraw string, publicKey []byte) (bool, error) {

	var jwk jose.JSONWebKey
	err := jwk.UnmarshalJSON([]byte(publicKey))
	if err != nil {
		return false, err
	}

	obj, err := jose.ParseSigned(jwsraw)
	if err != nil {
		return false, err
	}

	_, err = obj.Verify(jwk)

	if err != nil {
		return false, err
	}

	return true, nil
}

// JWK format to der
func ExportPublicKey(publicKey []byte) (string, error) {
	var jwk jose.JSONWebKey
	err := jwk.UnmarshalJSON([]byte(publicKey))
	if err != nil || !jwk.Valid() {
		return "", err
	}
	//pk, _ := jwk.Key.(*jose.RawJSONWebKey).EcPublicKey()
	pk, _ := jwk.Key.(*ecdsa.PublicKey)

	pkder, _ := x509.MarshalPKIXPublicKey(pk)

	return string(pkder), nil
}
