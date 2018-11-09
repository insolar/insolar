package message

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"testing"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/cryptohelpers/hash"
	"github.com/insolar/insolar/testutils"
	"github.com/stretchr/testify/assert"
)

func TestValidateToken(t *testing.T) {
	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	ref := testutils.RandomRef()

	tmp := core.Message(&GenesisRequest{})
	parcel, err := NewParcel(context.TODO(), tmp, ref, key, 1234, nil)
	assert.NoError(t, err)

	serialized := ToBytes(parcel.Message())
	msgHash := hash.IntegrityHasher().Hash(serialized)

	err = ValidateRoutingToken(&key.PublicKey, parcel.GetToken(), msgHash)
	assert.NoError(t, err)
}
