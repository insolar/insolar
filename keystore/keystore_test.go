package keystore

import (
	"crypto/ecdsa"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	testKeys    = "testdata/keys.json"
	testBadKeys = "testdata/bad_keys.json"
)

func TestNewKeyStore(t *testing.T) {
	ks, err := NewKeyStore(testKeys)
	require.NoError(t, err)
	require.NotNil(t, ks)
}

func TestNewKeyStore_Fails(t *testing.T) {
	ks, err := NewKeyStore(testBadKeys)
	require.Error(t, err)
	require.Nil(t, ks)
}

func TestKeyStore_GetPrivateKey(t *testing.T) {
	ks, err := NewKeyStore(testKeys)
	require.NoError(t, err)

	pk, err := ks.GetPrivateKey("")
	require.NotNil(t, pk)
	require.NoError(t, err)
}

func TestKeyStore_GetPrivateKeyReturnsECDSA(t *testing.T) {
	ks, err := NewKeyStore(testKeys)
	require.NoError(t, err)

	pk, err := ks.GetPrivateKey("")
	require.NotNil(t, pk)
	require.NoError(t, err)

	ecdsaPK, ok := pk.(*ecdsa.PrivateKey)
	require.NotNil(t, ecdsaPK)
	require.True(t, ok)
}
