// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

// +build functest

package functest

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetInfo(t *testing.T) {
	info := getInfo(t)
	require.NotNil(t, info)
	require.NotEqual(t, "", info["rootDomain"])
	require.NotEqual(t, "", info["rootMember"])
	require.NotEqual(t, "", info["nodeDomain"])
}
