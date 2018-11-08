package message

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"testing"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/testutils"
	"github.com/stretchr/testify/assert"
)

func TestValidateToken(t *testing.T) {
	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	ref := testutils.RandomRef()

	tmp := core.Message(&BootstrapRequest{})
	msg, err := NewSignedMessage(context.TODO(), tmp, ref, key, 1234, nil)
	assert.NoError(t, err)

	err = ValidateToken(&key.PublicKey,  msg)
	assert.NoError(t, err)
}
