// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package extractor

import (
	"testing"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/logicrunner/builtin/foundation"
	"github.com/stretchr/testify/require"
)

func TestPublicKeyResponse(t *testing.T) {
	testValue := "test_public_key"

	data, err := foundation.MarshalMethodResult(testValue, nil)
	require.NoError(t, err)

	result, err := PublicKeyResponse(data)

	require.NoError(t, err)
	require.Equal(t, testValue, result)
}

func TestPublicKeyResponse_ErrorResponse(t *testing.T) {
	testValue := "test_public_key"
	contractErr := &foundation.Error{S: "Custom test error"}

	data, err := foundation.MarshalMethodResult(testValue, contractErr)
	require.NoError(t, err)

	result, err := PublicKeyResponse(data)

	require.Error(t, err)
	require.Contains(t, err.Error(), "Has error in response")
	require.Contains(t, err.Error(), "Custom test error")
	require.Equal(t, "", result)
}

func TestPublicKeyResponse_UnmarshalError(t *testing.T) {
	testValue := "some_no_valid_data"

	data, err := insolar.Serialize(testValue)
	require.NoError(t, err)

	result, err := PublicKeyResponse(data)
	require.Error(t, err)
	require.Contains(t, err.Error(), "Can't unmarshal")
	require.Equal(t, "", result)
}

func TestCallResponse(t *testing.T) {
	testValue := map[string]interface{}{
		"string_value": "test_string",
		"int_value":    uint64(1),
	}

	data, err := foundation.MarshalMethodResult(testValue, nil)
	require.NoError(t, err)

	result, contractErr, err := CallResponse(data)

	require.NoError(t, err)
	require.Nil(t, contractErr)
	require.Equal(t, testValue, result)
}

func TestCallResponse_UnmarshalError(t *testing.T) {
	testValue := "some_no_valid_data"

	data, err := insolar.Serialize(testValue)
	require.NoError(t, err)

	result, contractErr, err := CallResponse(data)

	require.Error(t, err)
	require.Contains(t, err.Error(), "Can't unmarshal response")
	require.Nil(t, contractErr)
	require.Nil(t, result)
}
