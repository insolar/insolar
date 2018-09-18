package pulsar

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExportImportPrivateKey(t *testing.T) {
	privateKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)

	encoded, err := exportPrivateKey(privateKey)
	decoded, err := importPrivateKey(encoded)

	assert.NoError(t, err)
	assert.ObjectsAreEqual(decoded, privateKey)
}

func TestExportImportPublicKey(t *testing.T) {
	privateKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	publicKey := &privateKey.PublicKey

	encoded, err := exportPublicKey(publicKey)
	decoded, err := importPublicKey(encoded)

	assert.NoError(t, err)
	assert.ObjectsAreEqual(decoded, privateKey)
}
