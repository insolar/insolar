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
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"regexp"
	"runtime"
	"testing"
	"time"

	"github.com/insolar/insolar/api"
	"github.com/insolar/insolar/api/requester"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/logicrunner/builtin/foundation"
	"github.com/insolar/insolar/platformpolicy"
	"github.com/insolar/insolar/testutils"
	"github.com/insolar/insolar/testutils/launchnet"

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

func createMember(t *testing.T) *launchnet.User {
	member, err := newUserWithKeys()
	require.NoError(t, err)
	member.Ref = launchnet.Root.Ref

	result, err := signedRequest(t, member, "member.create", nil)
	require.NoError(t, err)
	ref, ok := result.(map[string]interface{})["reference"].(string)
	require.True(t, ok)
	member.Ref = ref
	return member
}

func createMigrationMemberForMA(t *testing.T, ma string) *launchnet.User {
	member, err := newUserWithKeys()
	require.NoError(t, err)
	member.Ref = launchnet.Root.Ref

	_, err = signedRequest(t, &launchnet.MigrationAdmin, "migration.addAddresses", map[string]interface{}{"migrationAddresses": []string{ma}})
	require.NoError(t, err)

	result, err := signedRequest(t, member, "member.migrationCreate", nil)
	require.NoError(t, err)
	ref, ok := result.(map[string]interface{})["reference"].(string)
	require.True(t, ok)
	member.Ref = ref
	return member

}

func addMigrationAddress(t *testing.T) {
	ba := testutils.RandomString()
	_, err := signedRequest(t, &launchnet.MigrationAdmin, "migration.addAddresses", map[string]interface{}{"migrationAddresses": []string{ba}})
	require.NoError(t, err)
}

func getBalanceNoErr(t *testing.T, caller *launchnet.User, reference string) *big.Int {
	balance, err := getBalance(t, caller, reference)
	require.NoError(t, err)
	return balance
}

func getBalance(t *testing.T, caller *launchnet.User, reference string) (*big.Int, error) {
	res, err := signedRequest(t, caller, "member.getBalance", map[string]interface{}{"reference": reference})
	if err != nil {
		return nil, err
	}
	amount, ok := new(big.Int).SetString(res.(map[string]interface{})["balance"].(string), 10)
	if !ok {
		return nil, fmt.Errorf("can't parse input amount")
	}
	return amount, nil
}

func migrate(t *testing.T, memberRef string, amount string, tx string, ma string, mdNum int) map[string]interface{} {
	anotherMember := createMember(t)

	_, err := signedRequest(t,
		launchnet.MigrationDaemons[mdNum],
		"deposit.migration",
		map[string]interface{}{"amount": amount, "ethTxHash": tx, "migrationAddress": ma})
	require.NoError(t, err)
	res, err := signedRequest(t, anotherMember, "member.getBalance", map[string]interface{}{"reference": memberRef})
	require.NoError(t, err)
	deposits, ok := res.(map[string]interface{})["deposits"].(map[string]interface{})
	require.True(t, ok)
	deposit, ok := deposits[tx].(map[string]interface{})
	sm := make(foundation.StableMap)
	require.NoError(t, err)
	confirmerReferencesMap := deposit["confirmerReferences"].(string)
	decoded, err := base64.StdEncoding.DecodeString(confirmerReferencesMap)
	require.NoError(t, err)
	err = sm.UnmarshalBinary(decoded)
	require.True(t, ok)
	require.Equal(t, sm[launchnet.MigrationDaemons[mdNum].Ref], amount)

	return deposit
}

func generateMigrationAddress() string {
	return testutils.RandomString()
}

const migrationAmount = "360000"

func fullMigration(t *testing.T, txHash string) *launchnet.User {
	activateDaemons(t)

	migrationAddress := testutils.RandomString()
	member := createMigrationMemberForMA(t, migrationAddress)

	migrate(t, member.Ref, migrationAmount, txHash, migrationAddress, 0)
	migrate(t, member.Ref, migrationAmount, txHash, migrationAddress, 2)
	migrate(t, member.Ref, migrationAmount, txHash, migrationAddress, 1)

	return member
}

func getRPSResponseBody(t testing.TB, postParams map[string]interface{}) []byte {
	jsonValue, _ := json.Marshal(postParams)

	postResp, err := http.Post(launchnet.TestRPCUrl, "application/json", bytes.NewBuffer(jsonValue))
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, postResp.StatusCode)
	body, err := ioutil.ReadAll(postResp.Body)
	require.NoError(t, err)
	return body
}

func getSeed(t testing.TB) string {
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

func getInfo(t testing.TB) infoResponse {
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

func getStatus(t testing.TB) statusResponse {
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

func activateDaemons(t *testing.T) {

	if len(launchnet.MigrationDaemons[0].Ref) > 0 {
		res, err := signedRequest(t, &launchnet.MigrationAdmin, "migration.checkDaemon", map[string]interface{}{"reference": launchnet.MigrationDaemons[0].Ref})
		require.NoError(t, err)
		status := res.(map[string]interface{})["status"].(string)
		if status == "inactive" {
			for _, user := range launchnet.MigrationDaemons {
				_, err := signedRequest(t, &launchnet.MigrationAdmin, "migration.activateDaemon", map[string]interface{}{"reference": user.Ref})
				require.NoError(t, err)
			}
		}
	}

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

func signedRequest(t *testing.T, user *launchnet.User, method string, params interface{}) (interface{}, error) {
	res, refStr, err := makeSignedRequest(user, method, params)

	var errMsg string
	if err != nil {
		errMsg = err.Error()
	}
	emptyRef := insolar.NewEmptyReference()

	require.NotEqual(t, "", refStr, "request ref is empty: %s", errMsg)
	require.NotEqual(t, emptyRef.String(), refStr, "request ref is zero: %s", errMsg)

	_, err = insolar.NewReferenceFromBase58(refStr)
	require.Nil(t, err)

	return res, err
}

func signedRequestWithEmptyRequestRef(t *testing.T, user *launchnet.User, method string, params interface{}) (interface{}, error) {
	res, refStr, err := makeSignedRequest(user, method, params)

	require.Equal(t, "", refStr)

	return res, err
}

func makeSignedRequest(user *launchnet.User, method string, params interface{}) (interface{}, string, error) {
	ctx := context.TODO()
	rootCfg, err := requester.CreateUserConfig(user.Ref, user.PrivKey, user.PubKey)
	if err != nil {
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

	res, err := requester.Send(ctx, launchnet.TestAPIURL, rootCfg, &requester.Params{
		CallSite:   method,
		CallParams: params,
		PublicKey:  user.PubKey,
		Test:       caller})

	if err != nil {
		return nil, "", err
	}

	resp := requester.ContractResponse{}
	err = json.Unmarshal(res, &resp)
	if err != nil {
		return nil, "", err
	}

	if resp.Error != nil {
		return nil, "", errors.New(resp.Error.Message)
	}

	if resp.Result == nil {
		return nil, "", errors.New("Error and result are nil")
	}

	return resp.Result.CallResult, resp.Result.RequestReference, nil

}

func newUserWithKeys() (*launchnet.User, error) {
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
	return &launchnet.User{
		PrivKey: string(privKeyStr),
		PubKey:  string(pubKeyStr),
	}, nil
}

// uploadContractOnce is needed for running tests with count
// use unique names when uploading contracts otherwise your contract won't be uploaded
func uploadContractOnce(t testing.TB, name string, code string) *insolar.Reference {
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

func uploadContract(t testing.TB, contractName string, contractCode string) *insolar.Reference {
	uploadBody := getRPSResponseBody(t, postParams{
		"jsonrpc": "2.0",
		"method":  "funcTestContract.upload",
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
	require.NotEqual(t, insolar.NewReferenceFromBytes(emptyRef), prototypeRef)

	return prototypeRef
}

func callConstructor(t testing.TB, prototypeRef *insolar.Reference, method string, args ...interface{}) *insolar.Reference {
	argsSerialized, err := insolar.Serialize(args)
	require.NoError(t, err)

	objectBody := getRPSResponseBody(t, postParams{
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

	callConstructorRes := struct {
		Version string              `json:"jsonrpc"`
		ID      string              `json:"id"`
		Result  api.CallMethodReply `json:"result"`
		Error   json2.Error         `json:"error"`
	}{}

	err = json.Unmarshal(objectBody, &callConstructorRes)
	require.NoError(t, err)
	require.Empty(t, callConstructorRes.Error)

	require.NotEmpty(t, callConstructorRes.Result.Object)

	objectRef, err := insolar.NewReferenceFromBase58(callConstructorRes.Result.Object)
	require.NoError(t, err)

	require.NotEqual(t, insolar.NewReferenceFromBytes(make([]byte, insolar.RecordRefSize)), objectRef)

	return objectRef
}

type callRes struct {
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

func callMethodExpectError(t testing.TB, objectRef *insolar.Reference, method string, args ...interface{}) api.CallMethodReply {
	callRes := callMethodNoChecks(t, objectRef, method, args...)
	require.NotEmpty(t, callRes.Error)

	return callRes.Result
}

func callMethodNoChecks(t testing.TB, objectRef *insolar.Reference, method string, args ...interface{}) callRes {
	argsSerialized, err := insolar.Serialize(args)
	require.NoError(t, err)

	respBody := getRPSResponseBody(t, postParams{
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

func waitUntilRequestProcessed(
	customFunction func() api.CallMethodReply,
	functionTimeout time.Duration,
	timeoutBetweenAttempts time.Duration,
	attempts int) (*api.CallMethodReply, error) {

	var lastErr error
	for i := 0; i < attempts; i++ {
		reply, err := waitForFunction(customFunction, functionTimeout)
		if err == nil {
			return reply, nil
		}
		lastErr = err
		time.Sleep(timeoutBetweenAttempts)
	}
	return nil, errors.New("Timeout was exceeded. " + lastErr.Error())
}

func waitForFunction(customFunction func() api.CallMethodReply, functionTimeout time.Duration) (*api.CallMethodReply, error) {
	ch := make(chan api.CallMethodReply, 1)
	go func() {
		ch <- customFunction()
	}()

	select {
	case result := <-ch:
		if result.Error != nil {
			return nil, errors.New(result.Error.Error())
		}
		return &result, nil
	case <-time.After(functionTimeout):
		return nil, errors.New("timeout was exceeded")
	}
}

func setMigrationDaemonsRef() error {
	for i, mDaemon := range launchnet.MigrationDaemons {
		daemon := mDaemon
		daemon.Ref = launchnet.Root.Ref
		res, _, err := makeSignedRequest(daemon, "member.get", nil)
		if err != nil {
			return errors.Wrap(err, "[ setup ] get member by public key failed ,key ")
		}
		launchnet.MigrationDaemons[i].Ref = res.(map[string]interface{})["reference"].(string)
	}
	return nil
}
