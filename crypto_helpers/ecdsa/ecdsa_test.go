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
	"crypto/ecdsa"
	"crypto/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExportImportPrivateKey(t *testing.T) {
	privateKey, _ := ecdsa.GenerateKey(GetCurve(), rand.Reader)

	encoded, err := ExportPrivateKey(privateKey)
	decoded, err := ImportPrivateKey(encoded)

	assert.NoError(t, err)
	assert.ObjectsAreEqual(decoded, privateKey)
}

func TestExportImportPublicKey(t *testing.T) {
	privateKey, _ := ecdsa.GenerateKey(GetCurve(), rand.Reader)
	publicKey := &privateKey.PublicKey

	encoded, err := ExportPublicKey(publicKey)
	decoded, err := ImportPublicKey(encoded)

	assert.NoError(t, err)
	assert.ObjectsAreEqual(decoded, privateKey)
}
