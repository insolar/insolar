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
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/api/requester"
	"github.com/insolar/insolar/testutils/launchnet"
)

func TestGetRequest(t *testing.T) {
	postResp, err := http.Get(launchnet.TestRPCUrl)
	defer postResp.Body.Close()
	require.NoError(t, err)
	require.Equal(t, http.StatusMethodNotAllowed, postResp.StatusCode)
}

func TestWrongUrl(t *testing.T) {
	jsonValue, _ := json.Marshal(postParams{})
	testURL := launchnet.HOST + launchnet.AdminPort + "/not_api"
	postResp, err := http.Post(testURL, "application/json", bytes.NewBuffer(jsonValue))
	defer postResp.Body.Close()
	require.NoError(t, err)
	require.Equal(t, http.StatusNotFound, postResp.StatusCode)
}

func TestWrongJson(t *testing.T) {
	postResp, err := http.Post(launchnet.TestRPCUrl, "application/json", bytes.NewBuffer([]byte("some not json value")))
	defer postResp.Body.Close()
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, postResp.StatusCode)
	body, err := ioutil.ReadAll(postResp.Body)
	require.NoError(t, err)

	response := &requester.ContractResponse{}
	unmarshalCallResponse(t, body, response)
	require.NotNil(t, response.Error)

	require.Equal(t, "invalid character 's' looking for beginning of value", response.Error.Message)
	require.Nil(t, response.Result)
}
