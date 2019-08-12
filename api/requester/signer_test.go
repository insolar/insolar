package requester

import (
	"testing"

	"github.com/insolar/insolar/logicrunner/builtin/foundation"
	"github.com/stretchr/testify/require"
)

func TestSignP256(t *testing.T) {
	testPublicKey := "-----BEGIN PUBLIC KEY-----\nMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAE3k6XSS/9dA8Jp/8rvRR72Cs9QBW2\nVSq5qhEnWWx7YBRK36C2O5J6+xGhlF+FsLkX11uuI0Ui0eOkmKFwa7pgVw==\n-----END PUBLIC KEY-----\n"
	testPrivateKey := "-----BEGIN PRIVATE KEY-----\nMHcCAQEEIAdn1/HospNWbdFwRasn4GRl48P1u6UG8PIagoa22nwhoAoGCCqGSM49\nAwEHoUQDQgAE3k6XSS/9dA8Jp/8rvRR72Cs9QBW2VSq5qhEnWWx7YBRK36C2O5J6\n+xGhlF+FsLkX11uuI0Ui0eOkmKFwa7pgVw==\n-----END PRIVATE KEY-----\n"
	testMessage := []byte("test")

	signature, err := Sign(testPrivateKey, testMessage)
	require.NoError(t, err)
	err = foundation.VerifySignature(testMessage, signature, testPublicKey, testPublicKey, false)
	require.NoError(t, err)
}

func TestSignP256K(t *testing.T) {
	testPublicKey := "-----BEGIN PUBLIC KEY-----\nMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAE3k6XSS/9dA8Jp/8rvRR72Cs9QBW2\nVSq5qhEnWWx7YBRK36C2O5J6+xGhlF+FsLkX11uuI0Ui0eOkmKFwa7pgVw==\n-----END PUBLIC KEY-----\n"
	testPrivateKey := "-----BEGIN PRIVATE KEY-----\nMHcCAQEEIAdn1/HospNWbdFwRasn4GRl48P1u6UG8PIagoa22nwhoAoGCCqGSM49\nAwEHoUQDQgAE3k6XSS/9dA8Jp/8rvRR72Cs9QBW2VSq5qhEnWWx7YBRK36C2O5J6\n+xGhlF+FsLkX11uuI0Ui0eOkmKFwa7pgVw==\n-----END PRIVATE KEY-----\n"
	testMessage := []byte("test")

	signature, err := Sign(testPrivateKey, testMessage)
	require.NoError(t, err)
	err = foundation.VerifySignature(testMessage, signature, testPublicKey, testPublicKey, false)
	require.NoError(t, err)
}
