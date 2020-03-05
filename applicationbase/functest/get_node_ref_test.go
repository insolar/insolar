///
// Copyright 2020 Insolar Technologies GmbH
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
///

// +build functest

package functest

import (
	"testing"

	"github.com/insolar/insolar/applicationbase/testutils"
	"github.com/insolar/insolar/applicationbase/testutils/launchnet"
	"github.com/insolar/insolar/applicationbase/testutils/testresponse"

	"github.com/stretchr/testify/require"
)

func getNodeRefSignedCall(t *testing.T, params map[string]interface{}) (string, error) {
	res, err := testutils.SignedRequest(t, launchnet.TestRPCUrl, &Root, "contract.getNodeRef", params)
	if err != nil {
		return "", err
	}
	return res.(string), nil
}

func TestGetNodeRefByPublicKey(t *testing.T) {
	const testRole = "light_material"
	publicKey := testutils.GenerateNodePublicKey(t)
	ref, err := registerNodeSignedCall(t, map[string]interface{}{"publicKey": publicKey, "role": testRole})
	require.NoError(t, err)
	require.NotNil(t, ref)

	nodeRef, err := getNodeRefSignedCall(t, map[string]interface{}{"publicKey": publicKey})
	require.NoError(t, err)
	require.Equal(t, ref, nodeRef)
}

func TestGetNodeRefByNotExistsPK(t *testing.T) {
	const testRole = "light_material"
	ref, err := registerNodeSignedCall(t, map[string]interface{}{"publicKey": testutils.GenerateNodePublicKey(t), "role": testRole})
	require.NoError(t, err)
	require.NotNil(t, ref)

	notExistingPublicKey := testutils.GenerateNodePublicKey(t)
	_, err = testutils.SignedRequestWithEmptyRequestRef(t, launchnet.TestRPCUrl, &Root,
		"contract.getNodeRef", map[string]interface{}{"publicKey": notExistingPublicKey})
	data := testresponse.CheckConvertRequesterError(t, err).Data
	require.Contains(t, data.Trace, "network node was not found by public key")
}

func TestGetNodeRefInvalidParams(t *testing.T) {
	const testRole = "light_material"
	ref, err := registerNodeSignedCall(t, map[string]interface{}{"publicKey": testutils.GenerateNodePublicKey(t), "role": testRole})
	require.NoError(t, err)
	require.NotNil(t, ref)

	_, err = testutils.SignedRequestWithEmptyRequestRef(t, launchnet.TestRPCUrl, &Root,
		"contract.getNodeRef", map[string]interface{}{"publicKey": 123})
	data := testresponse.CheckConvertRequesterError(t, err).Data
	require.Contains(t, data.Trace, "failed to get 'publicKey' param")
}
