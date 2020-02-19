///
// Copyright 2020 Insolar Technologies GmbH
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
///

// +build functest

package functest

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"runtime"
	"strings"
	"testing"

	"github.com/insolar/insolar/insolar/secrets"

	"github.com/insolar/insolar/api"
	"github.com/insolar/insolar/api/requester"
	"github.com/insolar/insolar/applicationbase/testutils/launchnet"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/platformpolicy"

	"github.com/insolar/rpc/v2/json2"

	"github.com/stretchr/testify/require"

	"github.com/pkg/errors"
)

type contractInfo struct {
	reference *insolar.Reference
	testName  string
}

var contracts = map[string]*contractInfo{}

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

func checkConvertRequesterError(t *testing.T, err error) *requester.Error {
	rv, ok := err.(*requester.Error)
	require.Truef(t, ok, "got wrong error %T (expected *requester.Error) with text '%s'", err, err.Error())
	return rv
}

func getRPSResponseBody(t testing.TB, URL string, postParams map[string]interface{}) []byte {
	jsonValue, _ := json.Marshal(postParams)

	postResp, err := http.Post(URL, "application/json", bytes.NewBuffer(jsonValue))
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, postResp.StatusCode)
	body, err := ioutil.ReadAll(postResp.Body)
	require.NoError(t, err)
	return body
}

func getSeed(t testing.TB) string {
	body := getRPSResponseBody(t, launchnet.TestRPCUrl, postParams{
		"jsonrpc": "2.0",
		"method":  "node.getSeed",
		"id":      "",
	})
	getSeedResponse := &getSeedResponse{}
	unmarshalRPCResponse(t, body, getSeedResponse)
	require.NotNil(t, getSeedResponse.Result)
	return getSeedResponse.Result.Seed
}

func getInfo(t testing.TB) infoResponse {
	pp := postParams{
		"jsonrpc": "2.0",
		"method":  "network.getInfo",
		"id":      1,
		"params":  map[string]string{},
	}
	body := getRPSResponseBody(t, launchnet.TestRPCUrl, pp)
	rpcInfoResponse := &rpcInfoResponse{}
	unmarshalRPCResponse(t, body, rpcInfoResponse)
	require.NotNil(t, rpcInfoResponse.Result)
	return rpcInfoResponse.Result
}

func getStatus(t testing.TB) statusResponse {
	body := getRPSResponseBody(t, launchnet.TestRPCUrl, postParams{
		"jsonrpc": "2.0",
		"method":  "node.getStatus",
		"id":      "1",
	})
	rpcStatusResponse := &rpcStatusResponse{}
	unmarshalRPCResponse(t, body, rpcStatusResponse)
	require.NotNil(t, rpcStatusResponse.Result)
	return rpcStatusResponse.Result
}

func unmarshalRPCResponse(t testing.TB, body []byte, response RPCResponseInterface) {
	err := json.Unmarshal(body, &response)
	require.NoError(t, err)
	require.Equal(t, "2.0", response.getRPCVersion())
	require.Nil(t, response.getError())
}

func unmarshalCallResponse(t testing.TB, body []byte, response *requester.ContractResponse) {
	err := json.Unmarshal(body, &response)
	require.NoError(t, err)
}

func signedRequest(t *testing.T, URL string, user *launchnet.User, method string, params interface{}) (interface{}, error) {
	res, refStr, err := makeSignedRequest(URL, user, method, params)

	if err != nil {
		var suffix string
		requesterError, ok := err.(*requester.Error)
		if ok {
			suffix = " [" + strings.Join(requesterError.Data.Trace, ": ") + "]"
		}
		t.Error("[" + method + "]" + err.Error() + suffix)
	}
	require.NotEmpty(t, refStr, "request ref is empty")
	require.NotEqual(t, insolar.NewEmptyReference().String(), refStr, "request ref is zero")

	_, err = insolar.NewReferenceFromString(refStr)
	require.Nil(t, err)

	return res, err
}

func signedRequestWithEmptyRequestRef(t *testing.T, URL string, user *launchnet.User, method string, params interface{}) (interface{}, error) {
	res, refStr, err := makeSignedRequest(URL, user, method, params)

	require.Equal(t, "", refStr)

	return res, err
}

func makeSignedRequest(URL string, user *launchnet.User, method string, params interface{}) (interface{}, string, error) {
	ctx := context.TODO()
	rootCfg, err := requester.CreateUserConfig(user.Ref, user.PrivKey, user.PubKey)
	if err != nil {
		var suffix string
		if requesterError, ok := err.(*requester.Error); ok {
			suffix = " [" + strings.Join(requesterError.Data.Trace, ": ") + "]"
		}
		fmt.Println(err.Error() + suffix)
		return nil, "", err
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

	seed, err := requester.GetSeed(URL)
	if err != nil {
		return nil, "", err
	}

	res, err := requester.SendWithSeed(ctx, URL, rootCfg, &requester.Params{
		CallSite:   method,
		CallParams: params,
		PublicKey:  user.PubKey,
		Reference:  user.Ref,
		Test:       caller}, seed)

	if err != nil {
		return nil, "", err
	}

	resp := requester.ContractResponse{}
	err = json.Unmarshal(res, &resp)
	if err != nil {
		return nil, "", err
	}

	if resp.Error != nil {
		return nil, "", resp.Error
	}

	if resp.Result == nil {
		return nil, "", errors.New("Error and result are nil")
	}
	return resp.Result.CallResult, resp.Result.RequestReference, nil

}

func newUserWithKeys() (*launchnet.User, error) {
	privateKey, err := secrets.GeneratePrivateKeyEthereum()
	if err != nil {
		return nil, err
	}

	privKeyStr, err := secrets.ExportPrivateKeyPEM(privateKey)
	if err != nil {
		return nil, err
	}
	publicKey := secrets.ExtractPublicKey(privateKey)
	pubKeyStr, err := secrets.ExportPublicKeyPEM(publicKey)
	if err != nil {
		return nil, err
	}
	return &launchnet.User{
		PrivKey: string(privKeyStr),
		PubKey:  string(pubKeyStr),
	}, nil
}

func createMember(t *testing.T) *launchnet.User {
	member, err := newUserWithKeys()
	require.NoError(t, err)
	member.Ref = launchnet.Root.Ref

	result, err := signedRequest(t, launchnet.TestRPCUrlPublic, member, "member.create", nil)
	require.NoError(t, err)
	ref, ok := result.(map[string]interface{})["reference"].(string)
	require.True(t, ok)
	member.Ref = ref
	return member
}

// uploadContractOnce is needed for running tests with count
// use unique names when uploading contracts otherwise your contract won't be uploaded
func uploadContractOnce(t testing.TB, name string, code string) *insolar.Reference {
	return uploadContractOnceExt(t, name, code, false)
}

func uploadContractOnceExt(t testing.TB, name string, code string, panicIsLogicalError bool) *insolar.Reference {
	if _, ok := contracts[name]; !ok {
		ref := uploadContract(t, name, code, panicIsLogicalError)
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

func uploadContract(t testing.TB, contractName string, contractCode string, panicIsLogicalError bool) *insolar.Reference {
	uploadBody := getRPSResponseBody(t, launchnet.TestRPCUrl, postParams{
		"jsonrpc": "2.0",
		"method":  "funcTestContract.upload",
		"id":      "",
		"params": map[string]interface{}{
			"name":                contractName,
			"code":                contractCode,
			"panicIsLogicalError": panicIsLogicalError,
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
	require.NoError(t, err, "unmarshal error")
	require.Empty(t, uploadRes.Error, "upload result error %#v", uploadRes)

	prototypeRef, err := insolar.NewReferenceFromString(uploadRes.Result.PrototypeRef)
	require.NoError(t, err)
	require.False(t, prototypeRef.IsEmpty())

	return prototypeRef
}

func callConstructorNoChecks(t testing.TB, prototypeRef *insolar.Reference, method string, args ...interface{}) callResult {
	argsSerialized, err := insolar.Serialize(args)
	require.NoError(t, err)

	objectBody := getRPSResponseBody(t, launchnet.TestRPCUrl, postParams{
		"jsonrpc": "2.0",
		"method":  "funcTestContract.callConstructor",
		"id":      "",
		"params": map[string]interface{}{
			"PrototypeRefString": prototypeRef.String(),
			"Method":             method,
			"MethodArgs":         argsSerialized,
		},
	})
	require.NotEmpty(t, objectBody)

	callConstructorRes := callResult{}
	err = json.Unmarshal(objectBody, &callConstructorRes)
	require.NoError(t, err)

	return callConstructorRes
}

func callConstructor(t testing.TB, prototypeRef *insolar.Reference, method string, args ...interface{}) *insolar.Reference {
	callConstructorRes := callConstructorNoChecks(t, prototypeRef, method, args...)
	require.Empty(t, callConstructorRes.Error)

	require.NotEmpty(t, callConstructorRes.Result.Object)

	objectRef, err := insolar.NewReferenceFromString(callConstructorRes.Result.Object)
	require.NoError(t, err)

	require.NotEqual(t, insolar.NewReferenceFromBytes(make([]byte, insolar.RecordRefSize)), objectRef)

	return objectRef
}

func callConstructorExpectSystemError(t testing.TB, prototypeRef *insolar.Reference, method string, args ...interface{}) string {
	callConstructorRes := callConstructorNoChecks(t, prototypeRef, method, args...)

	require.NotEmpty(t, callConstructorRes.Error)

	return callConstructorRes.Error.Message
}

type callResult struct {
	Version string              `json:"jsonrpc"`
	ID      string              `json:"id"`
	Result  api.CallMethodReply `json:"result"`
	Error   json2.Error         `json:"error"`
}

func callMethod(t testing.TB, objectRef *insolar.Reference, method string, args ...interface{}) api.CallMethodReply {
	callRes := callMethodNoChecks(t, objectRef, method, args...)
	require.Empty(t, callRes.Error)

	return callRes.Result
}

func callMethodNoChecks(t testing.TB, objectRef *insolar.Reference, method string, args ...interface{}) callResult {
	argsSerialized, err := insolar.Serialize(args)
	require.NoError(t, err)

	respBody := getRPSResponseBody(t, launchnet.TestRPCUrl, postParams{
		"jsonrpc": "2.0",
		"method":  "funcTestContract.callMethod",
		"id":      "",
		"params": map[string]interface{}{
			"ObjectRefString": objectRef.String(),
			"Method":          method,
			"MethodArgs":      argsSerialized,
		},
	})
	require.NotEmpty(t, respBody)

	callRes := struct {
		Version string              `json:"jsonrpc"`
		ID      string              `json:"id"`
		Result  api.CallMethodReply `json:"result"`
		Error   json2.Error         `json:"error"`
	}{}

	err = json.Unmarshal(respBody, &callRes)
	require.NoError(t, err)

	return callRes
}

func expectedError(t *testing.T, trace []string, expected string) {
	found := hasSubstring(trace, expected)
	require.True(t, found, "Expected error (%s) not found in trace: %v", expected, trace)
}

func hasSubstring(trace []string, expected string) bool {
	found := false
	for _, trace := range trace {
		found = strings.Contains(trace, expected)
		if found {
			return found
		}
	}
	return found
}

func generateNodePublicKey(t *testing.T) string {
	ks := platformpolicy.NewKeyProcessor()

	privKey, err := ks.GeneratePrivateKey()
	require.NoError(t, err)

	pubKeyStr, err := ks.ExportPublicKeyPEM(ks.ExtractPublicKey(privKey))
	require.NoError(t, err)

	return string(pubKeyStr)
}
