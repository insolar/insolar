// Copyright 2020 Insolar Network Ltd.
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

package extractor

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/logicrunner/builtin/foundation"
)

func TestStringResponse(t *testing.T) {
	testValue := "test_string"

	data, err := foundation.MarshalMethodResult(testValue, nil)
	require.NoError(t, err)

	result, err := stringResponse(data)

	require.NoError(t, err)
	require.Equal(t, testValue, result)
}

func TestStringResponse_ErrorResponse(t *testing.T) {
	testValue := "test_string"
	contractErr := &foundation.Error{S: "Custom test error"}

	data, err := foundation.MarshalMethodResult(testValue, contractErr)
	require.NoError(t, err)

	result, err := stringResponse(data)

	require.Error(t, err)
	require.Contains(t, err.Error(), "Has error in response")
	require.Contains(t, err.Error(), "Custom test error")
	require.Equal(t, "", result)
}

func TestStringResponse_UnmarshalError(t *testing.T) {
	testValue := "some_no_valid_data"

	data, err := insolar.Serialize(testValue)
	require.NoError(t, err)

	result, err := stringResponse(data)

	require.Error(t, err)
	require.Contains(t, err.Error(), "Can't unmarshal")
	require.Equal(t, "", result)
}
