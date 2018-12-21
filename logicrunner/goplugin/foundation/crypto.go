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

package foundation

import (
	"crypto"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/platformpolicy"
)

// TODO: this file should be removed

var platformCryptographyScheme = platformpolicy.NewPlatformCryptographyScheme()
var keyProcessor = platformpolicy.NewKeyProcessor()

// Sign signs given seed.
func Sign(data []byte, key crypto.PrivateKey) ([]byte, error) {
	signature, err := platformCryptographyScheme.Signer(key).Sign(data)
	if err != nil {
		return nil, err
	}
	return signature.Bytes(), nil
}

// Verify verifies signature.
func Verify(data []byte, signatureRaw []byte, publicKey crypto.PublicKey) bool {
	return platformCryptographyScheme.Verifier(publicKey).Verify(core.SignatureFromBytes(signatureRaw), data)
}

func GeneratePrivateKey() (crypto.PrivateKey, error) {
	return keyProcessor.GeneratePrivateKey()
}

func ImportPublicKey(publicKey string) (crypto.PublicKey, error) {
	return keyProcessor.ImportPublicKeyPEM([]byte(publicKey))
}

func ExportPublicKey(publicKey crypto.PublicKey) (string, error) {
	key, err := keyProcessor.ExportPublicKeyPEM(publicKey)
	return string(key), err
}

func ExtractPublicKey(privateKey crypto.PrivateKey) crypto.PublicKey {
	return keyProcessor.ExtractPublicKey(privateKey)
}
