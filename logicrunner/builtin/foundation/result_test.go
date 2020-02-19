// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

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
