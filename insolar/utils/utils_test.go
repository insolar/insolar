// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRandTraceID(t *testing.T) {
	traceID := RandTraceID()
	require.NotEmpty(t, traceID)
}
