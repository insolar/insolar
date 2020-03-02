// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package bits

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestResetBits(t *testing.T) {
	orig := []byte{0xFF}
	got := ResetBits(orig, 5)
	require.NotEqual(t, &orig, &got, "without overflow returns a new slice")

	gotWithOverflow := ResetBits(orig, 9)
	require.Equal(t, []byte{0xFF}, gotWithOverflow, "returns equals slice on overflow")
	require.Equal(t, &orig, &gotWithOverflow, "on overflow returns the same slice")
	require.Equal(t, []byte{0xFF}, orig, "original unchanged after resetBits")
}
