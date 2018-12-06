// +build functest

/*
 *    Copyright 2018 Insolar
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package functest

import (
	"testing"

	"github.com/stretchr/testify/require"
)

const TESTPUBLICKEY = "some_fancy_public_key"

func registerNodeSignedCall(params ...interface{}) (string, error) {
	res, err := signedRequest(&root, "RegisterNode", params...)
	if err != nil {
		return "", err
	}
	return res.(string), nil
}

func TestRegisterNodeVirtual(t *testing.T) {
	const testRole = "virtual"
	cert, err := registerNodeSignedCall(TESTPUBLICKEY, testRole)
	require.NoError(t, err)

	require.NotNil(t, cert)
}

func TestRegisterNodeHeavyMaterial(t *testing.T) {
	const testRole = "heavy_material"
	cert, err := registerNodeSignedCall(TESTPUBLICKEY, testRole)
	require.NoError(t, err)

	require.NotNil(t, cert)
}

func TestRegisterNodeLightMaterial(t *testing.T) {
	const testRole = "light_material"
	cert, err := registerNodeSignedCall(TESTPUBLICKEY, testRole)
	require.NoError(t, err)

	require.NotNil(t, cert)
}

func TestRegisterNodeNotExistRole(t *testing.T) {
	_, err := registerNodeSignedCall(TESTPUBLICKEY, "some_not_fancy_role")
	require.Contains(t, err.Error(),
		"[ RegisterNode ] Can't save as child: [ SaveAsChild ] on calling main API: couldn't save new object as child: executer error: problem with API call: Can't call constructor NewNodeRecord: Role is not supported: some_not_fancy_role")
}

func TestRegisterNodeByNoRoot(t *testing.T) {
	member := createMember(t, "Member1")
	const testRole = "virtual"
	_, err := signedRequest(member, "RegisterNode", TESTPUBLICKEY, testRole)
	require.Contains(t, err.Error(), "[ RegisterNode ] Only Root member can register node")
}
