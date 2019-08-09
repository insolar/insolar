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

package foundation

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMarshalMethodResult(t *testing.T) {
	data, err := MarshalMethodResult(10, nil)
	require.NoError(t, err)
	require.NotEmpty(t, data)

	var i int
	var contractErr *Error
	err = UnmarshalMethodResultSimplified(data, &i, &contractErr)
	require.NoError(t, err)
	require.Equal(t, 10, i)
	require.Nil(t, contractErr)
}

func TestMarshalMethodErrorResult(t *testing.T) {
	data, err := MarshalMethodErrorResult(errors.New("some"))
	require.NoError(t, err)
	require.NotEmpty(t, data)

	var i int
	var contractErr *Error
	err = UnmarshalMethodResultSimplified(data, &i, &contractErr)
	require.NoError(t, err)
	require.Equal(t, 0, i)
	require.Error(t, contractErr)
	require.Contains(t, contractErr.Error(), "some")
}
