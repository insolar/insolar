//
// Copyright 2019 Insolar Technologies GmbH
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

// +build functest

package functest

import (
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/insolar/insolar/testutils/launchnet"

	"github.com/stretchr/testify/require"
)

func TestChangeLogLevelOk(t *testing.T) {
	url := launchnet.HostDebug + "/debug/loglevel?level=debug"
	resp, err := http.Get(url)
	defer resp.Body.Close()
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	resBody, err := ioutil.ReadAll(resp.Body)
	require.Nil(t, err)
	require.Equal(t, "New log level: 'debug'\n", string(resBody))
}

func TestChangeLogLevelFail(t *testing.T) {
	url := launchnet.HostDebug + "/debug/loglevel?level=ololo"
	resp, err := http.Get(url)
	defer resp.Body.Close()
	require.NoError(t, err)
	require.NotEqual(t, http.StatusOK, resp.StatusCode)
}
