// +build functest

package functest

import (
	"testing"

	"github.com/insolar/insolar/applicationbase/testutils/launchnet"
	"github.com/insolar/insolar/applicationbase/testutils/testrequest"
	"github.com/insolar/insolar/applicationbase/testutils/testresponse"

	"github.com/stretchr/testify/require"
)

func getNodeRefSignedCall(t *testing.T, params map[string]interface{}) (string, error) {
	res, err := testrequest.SignedRequest(t, launchnet.TestRPCUrl, &Root, "contract.getNodeRef", params)
	if err != nil {
		return "", err
	}
	return res.(string), nil
}

func TestGetNodeRefByPublicKey(t *testing.T) {
	const testRole = "light_material"
	publicKey := testrequest.GenerateNodePublicKey(t)
	ref, err := registerNodeSignedCall(t, map[string]interface{}{"publicKey": publicKey, "role": testRole})
	require.NoError(t, err)
	require.NotNil(t, ref)

	nodeRef, err := getNodeRefSignedCall(t, map[string]interface{}{"publicKey": publicKey})
	require.NoError(t, err)
	require.Equal(t, ref, nodeRef)
}

func TestGetNodeRefByNotExistsPK(t *testing.T) {
	const testRole = "light_material"
	ref, err := registerNodeSignedCall(t, map[string]interface{}{"publicKey": testrequest.GenerateNodePublicKey(t), "role": testRole})
	require.NoError(t, err)
	require.NotNil(t, ref)

	notExistingPublicKey := testrequest.GenerateNodePublicKey(t)
	_, err = testrequest.SignedRequestWithEmptyRequestRef(t, launchnet.TestRPCUrl, &Root,
		"contract.getNodeRef", map[string]interface{}{"publicKey": notExistingPublicKey})
	data := testresponse.CheckConvertRequesterError(t, err).Data
	require.Contains(t, data.Trace, "network node was not found by public key")
}

func TestGetNodeRefInvalidParams(t *testing.T) {
	const testRole = "light_material"
	ref, err := registerNodeSignedCall(t, map[string]interface{}{"publicKey": testrequest.GenerateNodePublicKey(t), "role": testRole})
	require.NoError(t, err)
	require.NotNil(t, ref)

	_, err = testrequest.SignedRequestWithEmptyRequestRef(t, launchnet.TestRPCUrl, &Root,
		"contract.getNodeRef", map[string]interface{}{"publicKey": 123})
	data := testresponse.CheckConvertRequesterError(t, err).Data
	require.Contains(t, data.Trace, "failed to get 'publicKey' param")
}
