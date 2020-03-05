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

func TestGetInfo(t *testing.T) {
	info := testresponse.GetInfo(t)
	require.NotNil(t, info)
	require.NotEqual(t, "", info.NodeDomain)
}
