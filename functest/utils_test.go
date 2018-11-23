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

	"github.com/insolar/insolar/api/requester"
	"github.com/insolar/insolar/platformpolicy"
	"github.com/stretchr/testify/require"
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

type RPCResponseInterface interface {
	getRPCVersion() string
	getError() map[string]interface{}
}

type RPCResponse struct {
	RPCVersion string                 `json:"jsonrpc"`
	Error      map[string]interface{} `json:"error"`
}

func (r *RPCResponse) getRPCVersion() string {
	return r.RPCVersion
}

func (r *RPCResponse) getError() map[string]interface{} {
	return r.Error
}

type getSeedResponse struct {
	RPCResponse
	Result struct {
		Seed    string `json:"Seed"`
		TraceID string `json:"TraceID"`
	} `json:"result"`
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
	Error      string `json:"error"`
	RootDomain string `json:"root_domain"`
	RootMember string `json:"root_member"`
	NodeDomain string `json:"node_domain"`
}

type newInfoResponse struct {
	RootDomain string `json:"RootDomain"`
	RootMember string `json:"RootMember"`
	NodeDomain string `json:"NodeDomain"`
	TraceID    string `json:"TraceID"`
}

type rpcInfoResponse struct {
	RPCResponse
	Result newInfoResponse `json:"result"`
}

func createMember(t *testing.T, name string) *user {
	member, err := newUserWithKeys()
	require.NoError(t, err)
	result, err := signedRequest(&root, "CreateMember", name, member.pubKey)
	require.NoError(t, err)
	ref, ok := result.(string)
	require.True(t, ok)
	member.ref = ref
	return member
}

func getBalanceNoErr(t *testing.T, caller *user, reference string) int {
	balance, err := getBalance(caller, reference)
	require.NoError(t, err)
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
	postResp, err := http.Post(TestAPIURL, "application/json", bytes.NewBuffer(jsonValue))
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, postResp.StatusCode)
	body, err := ioutil.ReadAll(postResp.Body)
	require.NoError(t, err)
	return body
}

func getRPSResponseBody(t *testing.T, postParams map[string]interface{}) []byte {
	jsonValue, _ := json.Marshal(postParams)
	postResp, err := http.Post(TestRPC, "application/json", bytes.NewBuffer(jsonValue))
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, postResp.StatusCode)
	body, err := ioutil.ReadAll(postResp.Body)
	require.NoError(t, err)
	return body
}

func getSeed(t *testing.T) string {
	body := getRPSResponseBody(t, postParams{
		"jsonrpc": "2.0",
		"method":  "seed.Get",
		"id":      "",
	})
	getSeedResponse := &getSeedResponse{}
	unmarshalRPCResponse(t, body, getSeedResponse)
	require.NotNil(t, getSeedResponse.Result)
	return getSeedResponse.Result.Seed
}

func getInfo(t *testing.T) newInfoResponse {
	body := getRPSResponseBody(t, postParams{
		"jsonrpc": "2.0",
		"method":  "info.Get",
		"id":      "",
	})
	rpcInfoResponse := &rpcInfoResponse{}
	unmarshalRPCResponse(t, body, rpcInfoResponse)
	require.NotNil(t, rpcInfoResponse.Result)
	// info = rpcInfoResponse.Result
	return rpcInfoResponse.Result
}

func unmarshalResponse(t *testing.T, body []byte, response responseInterface) {
	err := json.Unmarshal(body, &response)
	require.NoError(t, err)
	require.Nil(t, response.getError())
}

func unmarshalRPCResponse(t *testing.T, body []byte, response RPCResponseInterface) {
	err := json.Unmarshal(body, &response)
	require.NoError(t, err)
	require.Equal(t, "2.0", response.getRPCVersion())
	require.Nil(t, response.getError())
}

func unmarshalResponseWithError(t *testing.T, body []byte, response responseInterface) {
	err := json.Unmarshal(body, &response)
	require.NoError(t, err)
	require.NotNil(t, response.getError())
}

func unmarshalCallResponse(t *testing.T, body []byte, response *response) {
	err := json.Unmarshal(body, &response)
	require.NoError(t, err)
}

type response struct {
	Result interface{}
	Error  string
}

func signedRequest(user *user, method string, params ...interface{}) (interface{}, error) {
	ctx := context.TODO()
	rootCfg, err := requester.CreateUserConfig(user.ref, user.privKey)
	if err != nil {
		return nil, err
	}
	res, err := requester.Send(ctx, TestAPIURL, rootCfg, &requester.RequestConfigJSON{
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
	ks := platformpolicy.NewKeyProcessor()

	privateKey, err := ks.GeneratePrivateKey()
	if err != nil {
		return nil, err
	}

	privKeyStr, err := ks.ExportPrivateKey(privateKey)
	if err != nil {
		return nil, err
	}
	publicKey := ks.ExtractPublicKey(privateKey)
	pubKeyStr, err := ks.ExportPublicKey(publicKey)
	if err != nil {
		return nil, err
	}
	return &user{
		privKey: string(privKeyStr),
		pubKey:  string(pubKeyStr),
	}, nil
}
