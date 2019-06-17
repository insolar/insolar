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

	"github.com/stretchr/testify/require"
)

const NOTEXISTINGPUBLICKEY = "not_existing_public_key"

func getNodeRefSignedCall(params map[string]interface{}) (string, error) {
	res, err := signedRequest(&root, "GetNodeRef", params)
	if err != nil {
		return "", err
	}
	return res.(string), nil
}

func TestGetNodeRefByPK(t *testing.T) {
	const testRole = "light_material"
	ref, err := registerNodeSignedCall(map[string]interface{}{"public": TESTPUBLICKEY, "role": testRole})
	require.NoError(t, err)
	require.NotNil(t, ref)

	nodeRef, err := getNodeRefSignedCall(map[string]interface{}{"public": TESTPUBLICKEY})
	require.NoError(t, err)
	require.Equal(t, ref, nodeRef)
}

func TestGetNodeRefByNotExistsPK(t *testing.T) {
	const testRole = "light_material"
	ref, err := registerNodeSignedCall(map[string]interface{}{"public": TESTPUBLICKEY, "role": testRole})
	require.NoError(t, err)
	require.NotNil(t, ref)

	nodeRef, err := getNodeRefSignedCall(map[string]interface{}{"public": NOTEXISTINGPUBLICKEY})
	require.Equal(t, "", nodeRef)
	require.Contains(t, err.Error(), "[ GetNodeRefByPK ] NetworkNode not found by PK: ")
}

func TestGetNodeRefInvalidParams(t *testing.T) {
	const testRole = "light_material"
	ref, err := registerNodeSignedCall(map[string]interface{}{"public": TESTPUBLICKEY, "role": testRole})
	require.NoError(t, err)
	require.NotNil(t, ref)

	nodeRef, err := getNodeRefSignedCall(map[string]interface{}{"public": 123})
	require.Equal(t, "", nodeRef)
	require.Contains(t, err.Error(), "[ getNodeRefCall ] Can't unmarshal params: ")
}
