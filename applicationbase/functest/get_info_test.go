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
