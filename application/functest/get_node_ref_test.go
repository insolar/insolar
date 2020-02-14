// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

// +build functest

package functest

import (
	"testing"

	"github.com/insolar/insolar/application/testutils/launchnet"

	"github.com/stretchr/testify/require"
)

func getNodeRefSignedCall(t *testing.T, params map[string]interface{}) (string, error) {
	res, err := signedRequest(t, launchnet.TestRPCUrl, &launchnet.Root, "contract.getNodeRef", params)
	if err != nil {
		return "", err
	}
	return res.(string), nil
}

func TestGetNodeRefByPublicKey(t *testing.T) {
	const testRole = "light_material"
	publicKey := generateNodePublicKey(t)
	ref, err := registerNodeSignedCall(t, map[string]interface{}{"publicKey": publicKey, "role": testRole})
	require.NoError(t, err)
	require.NotNil(t, ref)

	nodeRef, err := getNodeRefSignedCall(t, map[string]interface{}{"publicKey": publicKey})
	require.NoError(t, err)
	require.Equal(t, ref, nodeRef)
}

func TestGetNodeRefByNotExistsPK(t *testing.T) {
	const testRole = "light_material"
	ref, err := registerNodeSignedCall(t, map[string]interface{}{"publicKey": generateNodePublicKey(t), "role": testRole})
	require.NoError(t, err)
	require.NotNil(t, ref)

	notExistingPublicKey := generateNodePublicKey(t)
	_, err = signedRequestWithEmptyRequestRef(t, launchnet.TestRPCUrl, &launchnet.Root,
		"contract.getNodeRef", map[string]interface{}{"publicKey": notExistingPublicKey})
	data := checkConvertRequesterError(t, err).Data
	require.Contains(t, data.Trace, "network node was not found by public key")
}

func TestGetNodeRefInvalidParams(t *testing.T) {
	const testRole = "light_material"
	ref, err := registerNodeSignedCall(t, map[string]interface{}{"publicKey": generateNodePublicKey(t), "role": testRole})
	require.NoError(t, err)
	require.NotNil(t, ref)

	_, err = signedRequestWithEmptyRequestRef(t, launchnet.TestRPCUrl, &launchnet.Root,
		"contract.getNodeRef", map[string]interface{}{"publicKey": 123})
	data := checkConvertRequesterError(t, err).Data
	require.Contains(t, data.Trace, "failed to get 'publicKey' param")
}
