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
	"math/big"
	"regexp"
	"runtime"
	"strings"

	"io/ioutil"
	"net/http"
	"testing"

	"github.com/insolar/insolar/api"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/rpc/v2/json2"

	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/api/requester"
	"github.com/insolar/insolar/platformpolicy"

	"github.com/pkg/errors"
)

const sendRetryCount = 5

type contractInfo struct {
	reference *insolar.Reference
	testName  string
}

var contracts map[string]*contractInfo

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
		Seed    string `json:"seed"`
		TraceID string `json:"traceID"`
	} `json:"result"`
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
	NetworkState    string `json:"networkState"`
	WorkingListSize int    `json:"workingListSize"`
}

type rpcStatusResponse struct {
	RPCResponse
	Result statusResponse `json:"result"`
}

func createMember(t *testing.T) *user {
	member, err := newUserWithKeys()
	require.NoError(t, err)
	member.ref = root.ref

	addBurnAddress(t)

	result, err := retryableCreateMember(member, "contract.createMember", map[string]interface{}{}, true)
	require.NoError(t, err)
	ref, ok := result.(string)
	require.True(t, ok)
	member.ref = ref
	return member
}

func addBurnAddress(t *testing.T) {
	_, err := signedRequest(&migrationAdmin, "migration.addBurnAddresses", map[string]interface{}{"burnAddresses": []string{"fake_ba"}})
	require.NoError(t, err)
}

func getBalanceNoErr(t *testing.T, caller *user, reference string) *big.Int {
	balance, err := getBalance(caller, reference)
	require.NoError(t, err)
	return balance
}

func getBalance(caller *user, reference string) (*big.Int, error) {
	res, err := signedRequest(caller, "wallet.getBalance", map[string]interface{}{"reference": reference})
	if err != nil {
		return nil, err
	}
	amount, ok := new(big.Int).SetString(res.(string), 10)
	if !ok {
		return nil, fmt.Errorf("can't parse input amount")
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
		"method":  "node.getSeed",
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
		"method":  "network.getInfo",
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
		"method":  "node.getStatus",
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

func retryableCreateMember(user *user, method string, params map[string]interface{}, updatePublicKey bool) (interface{}, error) {
	// TODO: delete this after deduplication (INS-2778)
	var result interface{}
	var err error
	currentIterNum := 1
	for ; currentIterNum <= sendRetryCount; currentIterNum++ {
		result, err = signedRequest(user, method, params)
		if err == nil || !strings.Contains(err.Error(), "member for this publicKey already exist") {
			return result, err
		}
		fmt.Printf("CreateMember request was duplicated, retry. Attempt for duplicated: %d/%d\n", currentIterNum, sendRetryCount)
		newUser, nErr := newUserWithKeys()
		if nErr != nil {
			return nil, nErr
		}
		user.privKey = newUser.privKey
		if updatePublicKey {
			user.pubKey = newUser.pubKey
		}
	}
	return result, err
}

func signedRequest(user *user, method string, params map[string]interface{}) (interface{}, error) {
	ctx := context.TODO()
	rootCfg, err := requester.CreateUserConfig(user.ref, user.privKey, user.pubKey)
	if err != nil {
		return nil, err
	}

	var caller string
	fpcs := make([]uintptr, 1)
	for i := 2; i < 10; i++ {
		if n := runtime.Callers(i, fpcs); n == 0 {
			break
		}
		caller = runtime.FuncForPC(fpcs[0] - 1).Name()
		if ok, _ := regexp.MatchString(`\.Test`, caller); ok {
			break
		}
		caller = ""
	}

	var resp requester.ContractAnswer
	currentIterNum := 1
	for ; currentIterNum <= sendRetryCount; currentIterNum++ {
		res, err := requester.Send(ctx, TestAPIURL, rootCfg, &requester.Request{
			JSONRPC: "2.0",
			ID:      1,
			Method:  "api.call",
			Params:  requester.Params{CallSite: method, CallParams: params, PublicKey: user.pubKey},
			Test:    caller,
		})

		if err != nil {
			return nil, err
		}

		resp = requester.ContractAnswer{}
		err = json.Unmarshal(res, &resp)
		if err != nil {
			return nil, err
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

// uploadContractOnce is needed for running tests with count
// use unique names when uploading contracts otherwise your contract won't be uploaded
func uploadContractOnce(t *testing.T, name string, code string) *insolar.Reference {
	if _, ok := contracts[name]; !ok {
		ref := uploadContract(t, name, code)
		contracts[name] = &contractInfo{
			reference: ref,
			testName:  t.Name(),
		}
	}
	require.Equal(
		t, contracts[name].testName, t.Name(),
		"[ uploadContractOnce ] You cant use name of contract multiple times: "+contracts[name].testName,
	)
	return contracts[name].reference
}

func uploadContract(t *testing.T, contractName string, contractCode string) *insolar.Reference {
	uploadBody := getRPSResponseBody(t, postParams{
		"jsonrpc": "2.0",
		"method":  "contract.upload",
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

func callConstructor(t *testing.T, prototypeRef *insolar.Reference, method string, args ...interface{}) *insolar.Reference {
	argsSerialized, err := insolar.Serialize(args)
	require.NoError(t, err)

	objectBody := getRPSResponseBody(t, postParams{
		"jsonrpc": "2.0",
		"method":  "contract.callConstructor",
		"id":      "",
		"params": map[string]interface{}{
			"PrototypeRefString": prototypeRef.String(),
			"Method":             method,
			"MethodArgs":         argsSerialized,
		},
	})
	require.NotEmpty(t, objectBody)

	callConstructorRes := struct {
		Version string                   `json:"jsonrpc"`
		ID      string                   `json:"id"`
		Result  api.CallConstructorReply `json:"result"`
		Error   json2.Error              `json:"error"`
	}{}

	err = json.Unmarshal(objectBody, &callConstructorRes)
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
		"method":  "contract.callMethod",
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
