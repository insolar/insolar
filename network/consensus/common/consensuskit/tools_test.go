// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package consensuskit

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBftMajority(t *testing.T) {
	require.Equal(t, 4, BftMajority(5))

	require.Zero(t, BftMajority(0))

	require.Equal(t, -3, BftMajority(-5))
}

func TestBftMinority(t *testing.T) {
	require.Equal(t, 1, BftMinority(5))

	require.Zero(t, BftMinority(0))

	require.Equal(t, -2, BftMinority(-5))
}
