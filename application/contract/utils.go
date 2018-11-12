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

package contract

import (
	"crypto"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/platformpolicy"
)

// TODO: this file should be removed

// Sign signs given seed.
func Sign(data []byte, key crypto.PrivateKey) ([]byte, error) {
	signature, err := platformpolicy.NewPlatformCryptographyScheme().Signer(key).Sign(data)
	if err != nil {
		return nil, err
	}
	return signature.Bytes(), nil
}

// Verify verifies signature.
func Verify(data []byte, signatureRaw []byte, pubKey string) bool {
	return platformpolicy.NewPlatformCryptographyScheme().Verifier(pubKey).Verify(core.SignatureFromBytes(signatureRaw), data)
}

func GeneratePrivateKey() (crypto.PrivateKey, error) {
	return platformpolicy.NewKeyProcessor().GeneratePrivateKey()
}

func ExportPublicKey(publicKey crypto.PublicKey) (string, error) {
	key, err := platformpolicy.NewKeyProcessor().ExportPublicKey(publicKey)
	return string(key), err
}

func ExtractPublicKey(privateKey crypto.PrivateKey) crypto.PublicKey {
	return platformpolicy.NewKeyProcessor().ExtractPublicKey(privateKey)
}
