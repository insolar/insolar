/*
 *    Copyright 2018 Insolar
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package functest

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWrongUrl(t *testing.T) {
	jsonValue, _ := json.Marshal(postParams{})
	testURL := HOST + "/not_api"
	postResp, err := http.Post(testURL, "application/json", bytes.NewBuffer(jsonValue))
	require.NoError(t, err)
	require.Equal(t, http.StatusNotFound, postResp.StatusCode)
}

func TestGetRequest(t *testing.T) {
	postResp, err := http.Get(TestCallUrl)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, postResp.StatusCode)
	body, err := ioutil.ReadAll(postResp.Body)
	require.NoError(t, err)

	getResponse := &response{}
	unmarshalCallResponse(t, body, getResponse)
	require.NotNil(t, getResponse.Error)

	require.Equal(t, "[ UnmarshalRequest ] Empty body", getResponse.Error)
	require.Nil(t, getResponse.Result)
}

func TestWrongJson(t *testing.T) {
	postResp, err := http.Post(TestCallUrl, "application/json", bytes.NewBuffer([]byte("some not json value")))
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, postResp.StatusCode)
	body, err := ioutil.ReadAll(postResp.Body)
	require.NoError(t, err)

	response := &response{}
	unmarshalCallResponse(t, body, response)
	require.NotNil(t, response.Error)

	require.Equal(t, "[ UnmarshalRequest ] Can't unmarshal input params: invalid character 's' looking for beginning of value", response.Error)
	require.Nil(t, response.Result)
}
