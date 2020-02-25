// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

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
