package cryptography

import (
	"testing"

	"github.com/stretchr/testify/require"
)

const TestBadCert = "testdata/bad_keys.json"

func TestReadPrivateKey_BadPrivateKey(t *testing.T) {
	_, err := NewStorageBoundCryptographyService(TestBadCert)
	require.Contains(t, err.Error(), "Failed to create KeyStore")
}
