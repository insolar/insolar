// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package common

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCBORSerializer_Serialize(t *testing.T) {
	ser := NewCBORSerializer()
	buf := []byte{}
	err := ser.Serialize([]interface{}{nil}, &buf)
	require.NoError(t, err)
	require.Equal(t, []byte{0x81, 0xf6}, buf)
}
