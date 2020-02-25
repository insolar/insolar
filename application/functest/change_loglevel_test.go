// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

// +build functest

package functest

import (
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/insolar/insolar/application/testutils/launchnet"

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
