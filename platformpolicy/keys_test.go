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

package platformpolicy

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExportImportPrivateKey(t *testing.T) {
	ks := NewKeyProcessor()

	privateKey, _ := ks.GeneratePrivateKey()

	encoded, err := ks.ExportPrivateKeyPEM(privateKey)
	require.NoError(t, err)
	decoded, err := ks.ImportPrivateKeyPEM(encoded)
	require.NoError(t, err)

	assert.ObjectsAreEqual(decoded, privateKey)
}

func TestExportImportPublicKey(t *testing.T) {
	ks := NewKeyProcessor()

	privateKey, _ := ks.GeneratePrivateKey()
	publicKey := ks.ExtractPublicKey(privateKey)

	encoded, err := ks.ExportPublicKeyPEM(publicKey)
	require.NoError(t, err)
	decoded, err := ks.ImportPublicKeyPEM(encoded)
	require.NoError(t, err)

	assert.ObjectsAreEqual(decoded, privateKey)
}

func TestExportImportPublicKeyBinary(t *testing.T) {
	ks := NewKeyProcessor()

	privateKey, _ := ks.GeneratePrivateKey()
	publicKey := ks.ExtractPublicKey(privateKey)

	encoded, err := ks.ExportPublicKeyPEM(publicKey)
	require.NoError(t, err)

	bin, err := ks.ExportPublicKeyBinary(publicKey)
	require.NoError(t, err)
	assert.Len(t, bin, 66)

	binPK, err := ks.ImportPublicKeyBinary(bin)
	require.NoError(t, err)

	encodedBinPK, err := ks.ExportPublicKeyPEM(binPK)
	require.NoError(t, err)

	assert.Equal(t, encoded, encodedBinPK)
}
