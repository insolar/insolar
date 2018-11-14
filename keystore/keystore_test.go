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

package keystore

import (
	"crypto/ecdsa"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testKeys    = "testdata/keys.json"
	testBadKeys = "testdata/bad_keys.json"
)

func TestNewKeyStore(t *testing.T) {
	ks, err := NewKeyStore(testKeys)
	assert.NoError(t, err)
	assert.NotNil(t, ks)
}

func TestNewKeyStore_Fails(t *testing.T) {
	ks, err := NewKeyStore(testBadKeys)
	assert.Error(t, err)
	assert.Nil(t, ks)
}

func TestKeyStore_GetPrivateKey(t *testing.T) {
	ks, err := NewKeyStore(testKeys)
	assert.NoError(t, err)

	pk, err := ks.GetPrivateKey("")
	assert.NotNil(t, pk)
	assert.NoError(t, err)
}

func TestKeyStore_GetPrivateKeyReturnsECDSA(t *testing.T) {
	ks, err := NewKeyStore(testKeys)
	assert.NoError(t, err)

	pk, err := ks.GetPrivateKey("")
	assert.NotNil(t, pk)
	assert.NoError(t, err)

	ecdsaPK, ok := pk.(*ecdsa.PrivateKey)
	assert.NotNil(t, ecdsaPK)
	assert.True(t, ok)
}
