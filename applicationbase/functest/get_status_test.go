// +build functest

package functest

import (
	"testing"

	"github.com/insolar/insolar/applicationbase/testutils/launchnet"

	"github.com/stretchr/testify/require"
)

func TestGetStatus(t *testing.T) {
	status := getStatus(t)
	require.NotNil(t, status)

	numNodes, err := launchnet.GetNodesCount(AppPath)
	require.NoError(t, err)

	require.Equal(t, "CompleteNetworkState", status.NetworkState)
	require.Equal(t, numNodes, status.WorkingListSize)
}
