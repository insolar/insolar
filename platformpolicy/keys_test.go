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
	assert.Len(t, bin, 64)

	binPK, err := ks.ImportPublicKeyBinary(bin)
	require.NoError(t, err)

	encodedBinPK, err := ks.ExportPublicKeyPEM(binPK)
	require.NoError(t, err)

	assert.Equal(t, encoded, encodedBinPK)
}
