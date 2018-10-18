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

	"github.com/insolar/insolar/testutils"
	"github.com/stretchr/testify/assert"
)

type postParams map[string]interface{}

type errorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type responseInterface interface {
	getError() *errorResponse
}

type baseResponse struct {
	Qid string         `json:"qid"`
	Err *errorResponse `json:"error"`
}

func (r *baseResponse) getError() *errorResponse {
	return r.Err
}

type createMemberResponse struct {
	baseResponse
	Reference string `json:"reference"`
}

type sendMoneyResponse struct {
	baseResponse
	Success bool `json:"success"`
}

type getBalanceResponse struct {
	baseResponse
	Amount   uint   `json:"amount"`
	Currency string `json:"currency"`
}

type getSeedResponse struct {
	baseResponse
	Seed string `json:"seed"`
}

type isAuthorized struct {
	baseResponse
	PublicKey     string `json:"public_key"`
	Role          []int  `json:"role"`
	NetCoordCheck bool   `json:"netcoord_auth_success"`
}

type userInfo struct {
	Member string `json:"member"`
	Wallet uint   `json:"wallet"`
}

type dumpUserInfoResponse struct {
	baseResponse
	DumpInfo userInfo `json:"dump_info"`
}

type dumpAllUsersResponse struct {
	baseResponse
	DumpInfo []userInfo `json:"dump_info"`
}

type bootstrapNode struct {
	PublicKey string `json:"public_key"`
	Host      string `json:"host"`
}

type certificate struct {
	MajorityRule   int             `json:"majority_rule"`
	PublicKey      string          `json:"public_key"`
	Reference      string          `json:"reference"`
	Roles          []string        `json:"roles"`
	BootstrapNodes []bootstrapNode `json:"bootstrap_nodes"`
}

type registerNodeResponse struct {
	baseResponse
	Certificate certificate `json:"certificate"`
}

func createMember(t *testing.T) string {
	body := getResponseBody(t, postParams{
		"query_type": "create_member",
		"name":       testutils.RandomString(),
		"public_key": "000",
	})

	firstMemberResponse := &createMemberResponse{}
	unmarshalResponse(t, body, firstMemberResponse)

	return firstMemberResponse.Reference
}

func getBalance(t *testing.T, reference string) int {
	body := getResponseBody(t, postParams{
		"query_type": "get_balance",
		"reference":  reference,
	})

	firstBalanceResponse := &getBalanceResponse{}
	unmarshalResponse(t, body, firstBalanceResponse)

	return int(firstBalanceResponse.Amount)
}

func getResponseBody(t *testing.T, postParams map[string]interface{}) []byte {
	jsonValue, _ := json.Marshal(postParams)
	postResp, err := http.Post(TestURL, "application/json", bytes.NewBuffer(jsonValue))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, postResp.StatusCode)
	body, err := ioutil.ReadAll(postResp.Body)
	assert.NoError(t, err)
	return body
}

func getSeed(t *testing.T) string {
	body := getResponseBody(t, postParams{
		"query_type": "get_seed",
	})

	getSeedResponse := &getSeedResponse{}
	unmarshalResponse(t, body, getSeedResponse)

	return getSeedResponse.Seed
}

func unmarshalResponse(t *testing.T, body []byte, response responseInterface) {
	err := json.Unmarshal(body, &response)
	assert.NoError(t, err)
	assert.Nil(t, response.getError())
}

func unmarshalResponseWithError(t *testing.T, body []byte, response responseInterface) {
	err := json.Unmarshal(body, &response)
	assert.NoError(t, err)
	assert.NotNil(t, response.getError())
}
