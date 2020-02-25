// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package secrets

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestKeys_GetKeysFromFile(t *testing.T) {
	pair, err := ReadKeysFile("testdata/keypair.json", false)
	require.NoError(t, err, "read keys from json")
	assert.Equal(t, fmt.Sprintf("%T", pair.Private), "*ecdsa.PrivateKey", "private key has proper type")
	assert.Equal(t, fmt.Sprintf("%T", pair.Public), "*ecdsa.PublicKey", "public key has proper type")
}

func TestKeys_GetOnlyPublicKey(t *testing.T) {
	pair, err := ReadKeysFile("testdata/keypair.json", true)
	require.NoError(t, err, "read keys from json")
	assert.Equal(t, fmt.Sprintf("%T", pair.Private), "<nil>", "private key has proper type")
	assert.Equal(t, fmt.Sprintf("%T", pair.Public), "*ecdsa.PublicKey", "public key has proper type")
}

func TestKeys_GetOnlyPublic_WhenHasOnlyPublic(t *testing.T) {
	pair, err := ReadKeysFile("testdata/public_only.json", true)
	require.NoError(t, err, "read keys from json")
	assert.Equal(t, fmt.Sprintf("%T", pair.Private), "<nil>", "private key has proper type")
	assert.Equal(t, fmt.Sprintf("%T", pair.Public), "*ecdsa.PublicKey", "public key has proper type")
}
