package serialization

import (
	"bytes"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewClaimList(t *testing.T) {
	list := NewClaimList()
	assert.Equal(t, claimTypeEmpty, list.EndOfClaims.ClaimType())
	assert.Equal(t, 0, list.EndOfClaims.Length())
	assert.Len(t, list.Claims, 0)

	payload := []byte{1, 2, 3, 4, 5}
	claim := NewGenericClaim(payload)
	assert.Equal(t, claimTypeGeneric, claim.ClaimType())
	assert.Equal(t, len(claim.Payload), claim.Length())
	assert.Equal(t, payload, claim.Payload)
	list.Push(claim)
	assert.Len(t, list.Claims, 1)
	assert.Equal(t, claim, list.Claims[0])
}

func TestClaimList_SerializeDeserialize(t *testing.T) {
	list := NewClaimList()
	list.Push(NewGenericClaim([]byte{1, 2, 3, 4, 5}))
	list2 := ClaimList{}

	buf := make([]byte, 0)
	rw := bytes.NewBuffer(buf)
	w := newTrackableWriter(rw)
	packetCtx := newPacketContext(context.Background(), &Header{})
	serializeCtx := newSerializeContext(packetCtx, w, digester, signer, nil)

	err := list.SerializeTo(serializeCtx, rw)
	assert.NoError(t, err)

	r := newTrackableReader(rw)
	deserializeCtx := newDeserializeContext(packetCtx, r, nil)
	err = list2.DeserializeFrom(deserializeCtx, rw)
	assert.NoError(t, err)

	assert.Equal(t, list, list2)
}
