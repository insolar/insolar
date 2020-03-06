// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

// +build functest

package functest

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/applicationbase/testutils/testresponse"
)

func TestGetSeed(t *testing.T) {
	seed1 := testresponse.GetSeed(t)
	seed2 := testresponse.GetSeed(t)

	require.NotEqual(t, seed1, seed2)

}
