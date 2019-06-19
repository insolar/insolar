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
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"net"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/rpc/v2/json2"
	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/api"
	"github.com/insolar/insolar/api/requester"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/platformpolicy"

	"github.com/pkg/errors"
)

const sendRetryCount = 5

var contracts map[string]*insolar.Reference

type postParams map[string]interface{}

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

type bootstrapNode struct {
	PublicKey string `json:"public_key"`
}

type infoResponse struct {
	RootDomain string `json:"RootDomain"`
	RootMember string `json:"RootMember"`
	NodeDomain string `json:"NodeDomain"`
	TraceID    string `json:"TraceID"`
}

type rpcInfoResponse struct {
	RPCResponse
	Result infoResponse `json:"result"`
}

type statusResponse struct {
	NetworkState    string `json:"NetworkState"`
	WorkingListSize int    `json:"WorkingListSize"`
}

type rpcStatusResponse struct {
	RPCResponse
	Result statusResponse `json:"result"`
}

func createMember(t *testing.T, name string) *user {
	member, err := newUserWithKeys()
	require.NoError(t, err)
	member.ref = root.ref

	addBurnAddresses(t)

	result, err := signedRequest(member, "contract.createMember", map[string]interface{}{})
	require.NoError(t, err)
	ref, ok := result.(string)
	require.True(t, ok)
	member.ref = ref
	return member
}

func addBurnAddresses(t *testing.T) {
	_, err := signedRequest(&migrationAdmin, "wallet.addBurnAddresses", map[string]interface{}{"burnAddresses": []string{"fake_ba"}})
	require.NoError(t, err)
}

func getBalanceNoErr(t *testing.T, caller *user, reference string) *big.Int {
	balance, err := getBalance(caller, reference)
	require.NoError(t, err)
	return balance
}

func getBalance(caller *user, reference string) (*big.Int, error) {
	res, err := signedRequest(caller, "wallet.getMyBalance", map[string]interface{}{})
	if err != nil {
		return nil, err
	}
	amount, ok := new(big.Int).SetString(res.(string), 10)
	if !ok {
		return nil, fmt.Errorf("[ Transfer ] can't parse input amount")
	}
	return amount, nil
}

func getRPSResponseBody(t *testing.T, postParams map[string]interface{}) []byte {
	jsonValue, _ := json.Marshal(postParams)
	postResp, err := http.Post(TestRPCUrl, "application/json", bytes.NewBuffer(jsonValue))
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, postResp.StatusCode)
	body, err := ioutil.ReadAll(postResp.Body)
	require.NoError(t, err)
	return body
}

func getSeed(t *testing.T) string {
	body := getRPSResponseBody(t, postParams{
		"jsonrpc": "2.0",
		"method":  "node.GetSeed",
		"id":      "",
	})
	getSeedResponse := &getSeedResponse{}
	unmarshalRPCResponse(t, body, getSeedResponse)
	require.NotNil(t, getSeedResponse.Result)
	return getSeedResponse.Result.Seed
}

func getInfo(t *testing.T) infoResponse {
	pp := postParams{
		"jsonrpc": "2.0",
		"method":  "network.GetInfo",
		"id":      "",
	}
	body := getRPSResponseBody(t, pp)
	rpcInfoResponse := &rpcInfoResponse{}
	unmarshalRPCResponse(t, body, rpcInfoResponse)
	require.NotNil(t, rpcInfoResponse.Result)
	return rpcInfoResponse.Result
}

func getStatus(t *testing.T) statusResponse {
	body := getRPSResponseBody(t, postParams{
		"jsonrpc": "2.0",
		"method":  "node.GetStatus",
		"id":      "",
	})
	rpcStatusResponse := &rpcStatusResponse{}
	unmarshalRPCResponse(t, body, rpcStatusResponse)
	require.NotNil(t, rpcStatusResponse.Result)
	return rpcStatusResponse.Result
}

func unmarshalRPCResponse(t *testing.T, body []byte, response RPCResponseInterface) {
	err := json.Unmarshal(body, &response)
	require.NoError(t, err)
	require.Equal(t, "2.0", response.getRPCVersion())
	require.Nil(t, response.getError())
}

func unmarshalCallResponse(t *testing.T, body []byte, response *requester.ContractAnswer) {
	err := json.Unmarshal(body, &response)
	require.NoError(t, err)
}

func signedRequest(user *user, method string, params map[string]interface{}) (interface{}, error) {
	ctx := context.TODO()
	rootCfg, err := requester.CreateUserConfig(user.ref, user.privKey)
	if err != nil {
		return nil, err
	}
	var resp requester.ContractAnswer
	currentIterNum := 1
	for ; currentIterNum <= sendRetryCount; currentIterNum++ {
		res, err := requester.Send(ctx, TestAPIURL, rootCfg, &requester.Request{
			JSONRPC: "2.0",
			ID:      1,
			Method:  "call.api",
			Params:  requester.Params{CallSite: method, CallParams: params},
		})

		if netErr, ok := errors.Cause(err).(net.Error); ok && netErr.Timeout() {
			fmt.Printf("Network timeout, retry. Attempt: %d/%d\n", currentIterNum, sendRetryCount)
			fmt.Printf("Method: %s\n", method)
			time.Sleep(time.Second)
			resp.Error = &requester.Error{Message: netErr.Error()}
			continue
		} else if err != nil {
			return nil, err
		}

		resp = requester.ContractAnswer{}
		err = json.Unmarshal(res, &resp)
		if err != nil {
			return nil, err
		}

		if resp.Error != nil && strings.Contains(resp.Error.Message, "Messagebus timeout exceeded") {
			fmt.Printf("Messagebus timeout exceeded, retry. Attempt: %d/%d\n", currentIterNum, sendRetryCount)
			fmt.Printf("Method: %s\n", method)
			time.Sleep(time.Second)
			continue
		}

		break
	}

	if resp.Error != nil {
		if currentIterNum > sendRetryCount {
			return nil, errors.New("Number of retries exceeded. " + resp.Error.Message)
		}

		return nil, errors.New(resp.Error.Message)
	} else {
		if resp.Result == nil {
			return nil, errors.New("Error and result are nil")
		} else {
			return resp.Result.ContractResult, nil
		}
	}
}

func newUserWithKeys() (*user, error) {
	ks := platformpolicy.NewKeyProcessor()

	privateKey, err := ks.GeneratePrivateKey()
	if err != nil {
		return nil, err
	}

	privKeyStr, err := ks.ExportPrivateKeyPEM(privateKey)
	if err != nil {
		return nil, err
	}
	publicKey := ks.ExtractPublicKey(privateKey)
	pubKeyStr, err := ks.ExportPublicKeyPEM(publicKey)
	if err != nil {
		return nil, err
	}
	return &user{
		privKey: string(privKeyStr),
		pubKey:  string(pubKeyStr),
	}, nil
}

// this is needed for running tests with count
func uploadContractOnce(t *testing.T, name string, code string) *insolar.Reference {
	if _, ok := contracts[name]; !ok {
		contracts[name] = uploadContract(t, name, code)
	}
	return contracts[name]
}

func uploadContract(t *testing.T, contractName string, contractCode string) *insolar.Reference {
	uploadBody := getRPSResponseBody(t, postParams{
		"jsonrpc": "2.0",
		"method":  "contract.Upload",
		"id":      "",
		"params": map[string]string{
			"name": contractName,
			"code": contractCode,
		},
	})
	require.NotEmpty(t, uploadBody)

	uploadRes := struct {
		Version string          `json:"jsonrpc"`
		ID      string          `json:"id"`
		Result  api.UploadReply `json:"result"`
		Error   json2.Error     `json:"error"`
	}{}

	err := json.Unmarshal(uploadBody, &uploadRes)
	require.NoError(t, err)
	require.Empty(t, uploadRes.Error)

	prototypeRef, err := insolar.NewReferenceFromBase58(uploadRes.Result.PrototypeRef)
	require.NoError(t, err)

	emptyRef := make([]byte, insolar.RecordRefSize)
	require.NotEqual(t, insolar.Reference{}.FromSlice(emptyRef), prototypeRef)

	return prototypeRef
}

func callConstructor(t *testing.T, prototypeRef *insolar.Reference) *insolar.Reference {
	objectBody := getRPSResponseBody(t, postParams{
		"jsonrpc": "2.0",
		"method":  "contract.CallConstructor",
		"id":      "",
		"params": map[string]string{
			"PrototypeRefString": prototypeRef.String(),
		},
	})
	require.NotEmpty(t, objectBody)

	callConstructorRes := struct {
		Version string                   `json:"jsonrpc"`
		ID      string                   `json:"id"`
		Result  api.CallConstructorReply `json:"result"`
		Error   json2.Error              `json:"error"`
	}{}

	err := json.Unmarshal(objectBody, &callConstructorRes)
	require.NoError(t, err)
	require.Empty(t, callConstructorRes.Error)

	objectRef, err := insolar.NewReferenceFromBase58(callConstructorRes.Result.ObjectRef)
	require.NoError(t, err)

	require.NotEqual(t, insolar.Reference{}.FromSlice(make([]byte, insolar.RecordRefSize)), objectRef)

	return objectRef
}

func callMethod(t *testing.T, objectRef *insolar.Reference, method string, args ...interface{}) api.CallMethodReply {
	argsSerialized, err := insolar.Serialize(args)
	require.NoError(t, err)

	callMethodBody := getRPSResponseBody(t, postParams{
		"jsonrpc": "2.0",
		"method":  "contract.CallMethod",
		"id":      "",
		"params": map[string]interface{}{
			"ObjectRefString": objectRef.String(),
			"Method":          method,
			"MethodArgs":      argsSerialized,
		},
	})
	require.NotEmpty(t, callMethodBody)

	callRes := struct {
		Version string              `json:"jsonrpc"`
		ID      string              `json:"id"`
		Result  api.CallMethodReply `json:"result"`
		Error   json2.Error         `json:"error"`
	}{}

	err = json.Unmarshal(callMethodBody, &callRes)
	require.NoError(t, err)
	require.Empty(t, callRes.Error)

	return callRes.Result
}
