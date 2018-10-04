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

package ecdsa

import (
	"crypto/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExportImportPrivateKey(t *testing.T) {
	privateKey, _ := GeneratePrivateKey()

	encoded, err := ExportPrivateKey(privateKey)
	assert.NoError(t, err)
	decoded, err := ImportPrivateKey(encoded)
	assert.NoError(t, err)

	assert.ObjectsAreEqual(decoded, privateKey)
}

func TestExportImportPublicKey(t *testing.T) {
	privateKey, _ := GeneratePrivateKey()
	publicKey := &privateKey.PublicKey

	encoded, err := ExportPublicKey(publicKey)
	assert.NoError(t, err)
	decoded, err := ImportPublicKey(encoded)
	assert.NoError(t, err)

	assert.ObjectsAreEqual(decoded, privateKey)
}

func makeSeed() []byte {
	seed := make([]byte, 32)
	_, err := rand.Read(seed)
	if err != nil {
		panic(err)
	}

	return seed
}

func TestSignVerify(t *testing.T) {
	privateKey, _ := GeneratePrivateKey()
	seed := makeSeed()
	sign, err := Sign(seed, privateKey)
	assert.NoError(t, err)

	pubKeyStr, err := ExportPublicKey(&privateKey.PublicKey)
	assert.NoError(t, err)

	ok, err := Verify(seed, sign, pubKeyStr)
	assert.NoError(t, err)
	assert.True(t, ok)
}

func TestImportExportSignature(t *testing.T) {
	privateKey, _ := GeneratePrivateKey()
	seed := makeSeed()
	sign, err := Sign(seed, privateKey)
	assert.NoError(t, err)

	signStr := ExportSignature(sign)
	decodedSign, err := ImportSignature(signStr)
	assert.NoError(t, err)
	assert.Equal(t, sign, decodedSign)
}
