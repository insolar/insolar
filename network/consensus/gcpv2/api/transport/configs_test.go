// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/network/LICENSE.md.

package transport

import (
	"math"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestVerifySizes(t *testing.T) {
	ns := &NeighbourhoodSizes{}
	ns.NeighbourhoodSize = 1
	require.Panics(t, ns.VerifySizes)

	ns.NeighbourhoodSize = 5
	ns.NeighbourhoodTrustThreshold = 0
	require.Panics(t, ns.VerifySizes)

	ns.NeighbourhoodTrustThreshold = math.MaxUint8 + 1
	require.Panics(t, ns.VerifySizes)

	ns.NeighbourhoodTrustThreshold = 1
	ns.JoinersPerNeighbourhood = 0
	require.Panics(t, ns.VerifySizes)

	ns.JoinersPerNeighbourhood = 1
	require.Panics(t, ns.VerifySizes)

	ns.JoinersPerNeighbourhood = 2
	ns.JoinersBoost = -1
	require.Panics(t, ns.VerifySizes)

	ns.JoinersBoost = 0
	ns.NeighbourhoodSize = 0
	require.Panics(t, ns.VerifySizes)

	ns.NeighbourhoodSize = 5
	require.NotPanics(t, ns.VerifySizes)
}
