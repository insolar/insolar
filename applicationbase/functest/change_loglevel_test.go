// +build functest

package functest

import (
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/insolar/insolar/applicationbase/testutils/launchnet"

	"github.com/stretchr/testify/require"
)

func TestChangeLogLevelOk(t *testing.T) {
	launchnet.RunOnlyWithLaunchnet(t)
	url := launchnet.HostDebug + "/debug/loglevel?level=debug"
	resp, err := http.Get(url)
	require.NoError(t, err)
	defer resp.Body.Close() //nolint: errcheck
	require.Equal(t, http.StatusOK, resp.StatusCode)
	resBody, err := ioutil.ReadAll(resp.Body)
	require.Nil(t, err)
	require.Equal(t, "New log level: 'debug'\n", string(resBody))
}

func TestChangeLogLevelFail(t *testing.T) {
	launchnet.RunOnlyWithLaunchnet(t)
	url := launchnet.HostDebug + "/debug/loglevel?level=ololo"
	resp, err := http.Get(url)
	require.NoError(t, err)
	defer resp.Body.Close() //nolint: errcheck
	require.NotEqual(t, http.StatusOK, resp.StatusCode)
}
