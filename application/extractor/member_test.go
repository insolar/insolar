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

func TestPublicKeyResponse(t *testing.T) {
	testValue := "test_public_key"

	data, err := core.Serialize([]interface{}{testValue, nil})
	require.NoError(t, err)

	result, err := PublicKeyResponse(data)

	require.NoError(t, err)
	require.Equal(t, testValue, result)
}

func TestPublicKeyResponse_ErrorResponse(t *testing.T) {
	testValue := "test_public_key"
	contractErr := &foundation.Error{S: "Custom test error"}

	data, err := core.Serialize([]interface{}{testValue, contractErr})
	require.NoError(t, err)

	result, err := PublicKeyResponse(data)

	require.Contains(t, err.Error(), "Has error in response")
	require.Contains(t, err.Error(), "Custom test error")
	require.Equal(t, "", result)
}

func TestPublicKeyResponse_UnmarshalError(t *testing.T) {
	testValue := "some_no_valid_data"

	data, err := core.Serialize(testValue)
	require.NoError(t, err)

	result, err := PublicKeyResponse(data)

	require.Contains(t, err.Error(), "Can't unmarshal")
	require.Equal(t, "", result)
}

func TestCallResponse(t *testing.T) {
	testValue := map[interface{}]interface{}{
		"string_value": "test_string",
		"int_value":    uint64(1),
	}

	data, err := core.Serialize([]interface{}{testValue, nil})
	require.NoError(t, err)

	result, contractErr, err := CallResponse(data)

	require.NoError(t, err)
	require.Nil(t, contractErr)
	require.Equal(t, testValue, result)
}

func TestCallResponse_UnmarshalError(t *testing.T) {
	testValue := "some_no_valid_data"

	data, err := core.Serialize(testValue)
	require.NoError(t, err)

	result, contractErr, err := CallResponse(data)

	require.Contains(t, err.Error(), "Can't unmarshal response")
	require.Nil(t, contractErr)
	require.Nil(t, result)
}
