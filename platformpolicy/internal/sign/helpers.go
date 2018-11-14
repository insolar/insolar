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
)

func MustConvertPublicKeyToEcdsa(publicKey crypto.PublicKey) *ecdsa.PublicKey {
	ecdsaPublicKey, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		panic("[ Sign ] Failed to convert public key to ecdsa public key")
	}
	return ecdsaPublicKey
}

func MustConvertPrivateKeyToEcdsa(privateKey crypto.PrivateKey) *ecdsa.PrivateKey {
	ecdsaPrivateKey, ok := privateKey.(*ecdsa.PrivateKey)
	if !ok {
		panic("[ Sign ] Failed to convert private key to ecdsa private key")
	}
	return ecdsaPrivateKey
}
