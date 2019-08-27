//
// Copyright 2019 Insolar Technologies GmbH
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

// +build functest

package functest

import (
	"testing"

	"github.com/insolar/insolar/testutils/launchnet"

	"github.com/stretchr/testify/require"
)

const NOTEXISTINGPUBLICKEY = "not_existing_public_key"

func getNodeRefSignedCall(t *testing.T, params map[string]interface{}) (string, error) {
	res, err := signedRequest(t, launchnet.TestRPCUrl, &launchnet.Root, "contract.getNodeRef", params)
	if err != nil {
		return "", err
	}
	return res.(string), nil
}

func TestGetNodeRefByPublicKey(t *testing.T) {
	const testRole = "light_material"
	ref, err := registerNodeSignedCall(t, map[string]interface{}{"publicKey": TESTPUBLICKEY, "role": testRole})
	require.NoError(t, err)
	require.NotNil(t, ref)

	nodeRef, err := getNodeRefSignedCall(t, map[string]interface{}{"publicKey": TESTPUBLICKEY})
	require.NoError(t, err)
	require.Equal(t, ref, nodeRef)
}

func TestGetNodeRefByNotExistsPK(t *testing.T) {
	const testRole = "light_material"
	ref, err := registerNodeSignedCall(t, map[string]interface{}{"publicKey": TESTPUBLICKEY, "role": testRole})
	require.NoError(t, err)
	require.NotNil(t, ref)

	_, err = signedRequestWithEmptyRequestRef(t, launchnet.TestRPCUrl, &launchnet.Root,
		"contract.getNodeRef", map[string]interface{}{"publicKey": NOTEXISTINGPUBLICKEY})
	require.Error(t, err)
	require.Contains(t, err.Error(), "network node was not found by public key:")
}

func TestGetNodeRefInvalidParams(t *testing.T) {
	const testRole = "light_material"
	ref, err := registerNodeSignedCall(t, map[string]interface{}{"publicKey": TESTPUBLICKEY, "role": testRole})
	require.NoError(t, err)
	require.NotNil(t, ref)

	_, err = signedRequestWithEmptyRequestRef(t, launchnet.TestRPCUrl, &launchnet.Root,
		"contract.getNodeRef", map[string]interface{}{"publicKey": 123})
	require.Error(t, err)
	require.Contains(t, err.Error(), "incorrect input: failed to get 'publicKey' param")
}
