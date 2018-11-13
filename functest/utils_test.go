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
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/insolar/insolar/api/requesters"
	"github.com/insolar/insolar/cryptohelpers/ecdsa"
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
	Role          int    `json:"role"`
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

type infoResponse struct {
	Error      string            `json:"error"`
	RootDomain string            `json:"root_domain"`
	RootMember string            `json:"root_member"`
	Prototypes map[string]string `json:"prototypes"`
}

func createMember(t *testing.T, name string) *user {
	member, err := newUserWithKeys()
	assert.NoError(t, err)
	result, err := signedRequest(&root, "CreateMember", name, member.pubKey)
	assert.NoError(t, err)
	ref, ok := result.(string)
	assert.True(t, ok)
	member.ref = ref
	return member
}

func getBalanceNoErr(t *testing.T, caller *user, reference string) int {
	balance, err := getBalance(caller, reference)
	assert.NoError(t, err)
	return balance
}

func getBalance(caller *user, reference string) (int, error) {
	res, err := signedRequest(caller, "GetBalance", reference)
	if err != nil {
		return 0, err
	}
	amount, ok := res.(float64)
	if !ok {
		return 0, errors.New("result is not int")
	}
	return int(amount), nil
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

func getInfo(t *testing.T) infoResponse {
	resp, err := http.Get(TestURL + "/info")
	assert.NoError(t, err)
	body, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	err = json.Unmarshal(body, &info)
	assert.NoError(t, err)
	return info
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

type response struct {
	Result interface{}
	Error  string
}

func signedRequest(user *user, method string, params ...interface{}) (interface{}, error) {
	ctx := context.TODO()
	rootCfg, err := requesters.CreateUserConfig(user.ref, user.privKey)
	if err != nil {
		return nil, err
	}
	res, err := requesters.Send(ctx, TestURL, rootCfg, &requesters.RequestConfigJSON{
		Method: method,
		Params: params,
	})
	if err != nil {
		return nil, err
	}
	var resp = response{}
	err = json.Unmarshal(res, &resp)

	if err != nil {
		return nil, err
	}
	if resp.Error != "" {
		return resp.Result, errors.New(resp.Error)
	}
	return resp.Result, nil
}

func newUserWithKeys() (*user, error) {
	key, err := ecdsa.GeneratePrivateKey()
	if err != nil {
		return nil, err
	}
	privKeyStr, err := ecdsa.ExportPrivateKey(key)
	if err != nil {
		return nil, err
	}
	pubKeyStr, err := ecdsa.ExportPublicKey(&key.PublicKey)
	if err != nil {
		return nil, err
	}
	return &user{
		privKey: privKeyStr,
		pubKey:  pubKeyStr,
	}, nil
}
