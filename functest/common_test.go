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

	"github.com/insolar/insolar/api"
	"github.com/insolar/insolar/testutils"
	"github.com/stretchr/testify/assert"
)

func TestWrongUrl(t *testing.T) {
	jsonValue, _ := json.Marshal(postParams{
		"query_type": "dump_all_users",
	})
	testURL := HOST + "/not_api/v1"
	postResp, err := http.Post(testURL, "application/json", bytes.NewBuffer(jsonValue))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, postResp.StatusCode)
}

func TestGetRequest(t *testing.T) {
	postResp, err := http.Get(TestURL)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, postResp.StatusCode)
	body, err := ioutil.ReadAll(postResp.Body)
	assert.NoError(t, err)

	getResponse := &baseResponse{}
	unmarshalResponseWithError(t, body, getResponse)

	assert.Equal(t, api.BadRequest, getResponse.Err.Code)
	assert.Equal(t, "Bad request", getResponse.Err.Message)
}

func TestWrongJson(t *testing.T) {
	postResp, err := http.Post(TestURL, "application/json", bytes.NewBuffer([]byte("some not json value")))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, postResp.StatusCode)
	body, err := ioutil.ReadAll(postResp.Body)
	assert.NoError(t, err)

	response := &baseResponse{}
	unmarshalResponseWithError(t, body, response)

	assert.Equal(t, api.BadRequest, response.Err.Code)
	assert.Equal(t, "Bad request", response.Err.Message)
}

func TestWrongQueryType(t *testing.T) {
	body := getResponseBody(t, postParams{
		"query_type": "wrong_query_type",
	})

	response := &baseResponse{}
	unmarshalResponseWithError(t, body, response)

	assert.Equal(t, api.BadRequest, response.Err.Code)
	assert.Equal(t, "Wrong query parameter 'query_type' = 'wrong_query_type'", response.Err.Message)
}

func TestWithoutQueryType(t *testing.T) {
	body := getResponseBody(t, postParams{})

	response := &baseResponse{}
	unmarshalResponseWithError(t, body, response)

	assert.Equal(t, api.BadRequest, response.Err.Code)
	assert.Equal(t, "Wrong query parameter 'query_type' = ''", response.Err.Message)
}

func TestTooMuchParams(t *testing.T) {
	body := getResponseBody(t, postParams{
		"query_type": "create_member",
		"some_param": "irrelevant info",
		"name":       testutils.RandomString(),
	})

	firstMemberResponse := &createMemberResponse{}
	unmarshalResponse(t, body, firstMemberResponse)

	firstMemberRef := firstMemberResponse.Reference
	assert.NotEqual(t, "", firstMemberRef)
}

func TestQueryTypeAsIntParams(t *testing.T) {
	body := getResponseBody(t, postParams{
		"query_type": 100,
	})

	response := &baseResponse{}
	unmarshalResponseWithError(t, body, response)

	assert.Equal(t, api.BadRequest, response.Err.Code)
	assert.Equal(t, "Bad request", response.Err.Message)
}

func TestWrongTypeInParams(t *testing.T) {
	body := getResponseBody(t, postParams{
		"query_type": "create_member",
		"name":       128182187,
	})

	response := &baseResponse{}
	unmarshalResponseWithError(t, body, response)

	assert.Equal(t, api.BadRequest, response.Err.Code)
	assert.Equal(t, "Bad request", response.Err.Message)
}
