// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/network/LICENSE.md.

package misbehavior

import (
	"testing"

	"github.com/pkg/errors"

	"github.com/stretchr/testify/require"
)

func TestIs(t *testing.T) {
	// TODO
	/*r := NewReportMock(t)
	require.True(t, Is(r))*/

	err := errors.New("test")
	require.False(t, Is(err))
}

func TestReportOf(t *testing.T) {
	// TODO
	/*r := NewReportMock(t)
	require.NotNil(t, Of(r))*/

	err := errors.New("test")
	require.Nil(t, Of(err))
}

func TestCategory(t *testing.T) {
	require.Equal(t, Category(1), Type(1<<32).Category())
}

func TestType(t *testing.T) {
	require.Equal(t, 1, Type(1).Type())

	require.Zero(t, Type(1<<32).Type())
}

func TestCategoryOf(t *testing.T) {
	require.Zero(t, Category(0).Of(0))

	require.Equal(t, Type(1<<32), Category(1).Of(0))

	require.Equal(t, Type((1<<32)+1), Category(1).Of(1))
}
