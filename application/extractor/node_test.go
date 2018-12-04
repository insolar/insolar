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

package extractor

import (
	"testing"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/logicrunner/goplugin/foundation"
	"github.com/stretchr/testify/require"
)

func TestNodeInfoResponse(t *testing.T) {
	testPK := "test_public_key"
	testRole := core.StaticRoleVirtual

	testValue := struct {
		PublicKey string
		Role      core.StaticRole
	}{
		PublicKey: testPK,
		Role:      testRole,
	}

	data, err := core.Serialize([]interface{}{testValue, nil})
	require.NoError(t, err)

	pk, role, err := NodeInfoResponse(data)

	require.NoError(t, err)
	require.Equal(t, testPK, pk)
	require.Equal(t, testRole.String(), role)
}

func TestNodeInfoResponse_ErrorResponse(t *testing.T) {
	testPK := "test_public_key"
	testRole := core.StaticRoleVirtual

	testValue := struct {
		PublicKey string
		Role      core.StaticRole
	}{
		PublicKey: testPK,
		Role:      testRole,
	}
	contractErr := &foundation.Error{S: "Custom test error"}

	data, err := core.Serialize([]interface{}{testValue, contractErr})
	require.NoError(t, err)

	pk, role, err := NodeInfoResponse(data)

	require.Contains(t, err.Error(), "Has error in response")
	require.Contains(t, err.Error(), "Custom test error")
	require.Equal(t, "", pk)
	require.Equal(t, "", role)
}

func TestNodeInfoResponse_UnmarshalError(t *testing.T) {
	testValue := "some_no_valid_data"

	data, err := core.Serialize(testValue)
	require.NoError(t, err)

	pk, role, err := NodeInfoResponse(data)

	require.Contains(t, err.Error(), "Can't unmarshal response")
	require.Equal(t, "", pk)
	require.Equal(t, "", role)
}
