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
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"regexp"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/insolar/insolar/application/genesisrefs"

	"github.com/insolar/insolar/application/api"
	"github.com/insolar/insolar/application/api/requester"
	"github.com/insolar/insolar/application/testutils/launchnet"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/platformpolicy"

	"github.com/insolar/rpc/v2/json2"

	"github.com/stretchr/testify/require"

	"github.com/pkg/errors"
)

const (
	countTwoActiveDaemon = iota + 2
	countThreeActiveDaemon
)

const TestDepositAmount string = "1000000000000000000"

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

func checkConvertRequesterError(t *testing.T, err error) *requester.Error {
	rv, ok := err.(*requester.Error)
	require.Truef(t, ok, "got wrong error %T (expected *requester.Error) with text '%s'", err, err.Error())
	return rv
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

func createMigrationMemberForMA(t *testing.T) *launchnet.User {
	member, err := newUserWithKeys()
	require.NoError(t, err)
	member.Ref = launchnet.Root.Ref

	result, err := signedRequest(t, launchnet.TestRPCUrlPublic, member, "member.migrationCreate", nil)
	require.NoError(t, err)
	ref, ok := result.(map[string]interface{})["reference"].(string)
	require.True(t, ok)
	member.Ref = ref

	ma, ok := result.(map[string]interface{})["migrationAddress"].(string)
	require.True(t, ok)
	member.MigrationAddress = ma
	return member

}

func generateMigrationAddress() (string, error) {
	bytes := make([]byte, 20)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func getBalanceNoErr(t *testing.T, caller *launchnet.User, reference string) *big.Int {
	balance, _ := getBalanceAndDepositsNoErr(t, caller, reference)
	return balance
}

func getAdminDepositBalance(t *testing.T, caller *launchnet.User, reference string) (*big.Int, error) {
	_, deposits := getBalanceAndDepositsNoErr(t, caller, reference)
	mapd, ok := deposits[genesisrefs.FundsDepositName].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("can't parse deposit")
	}
	amount, ok := new(big.Int).SetString(mapd["balance"].(string), 10)
	if !ok {
		return nil, fmt.Errorf("can't parse deposit balance")
	}
	return amount, nil
}

func getBalanceAndDepositsNoErr(t *testing.T, caller *launchnet.User, reference string) (*big.Int, map[string]interface{}) {
	balance, deposits, err := getBalanceAndDeposits(t, caller, reference)
	require.NoError(t, err)
	return balance, deposits
}

func getBalanceAndDeposits(t *testing.T, caller *launchnet.User, reference string) (*big.Int, map[string]interface{}, error) {
	res, err := signedRequest(t, launchnet.TestRPCUrl, caller, "member.getBalance", map[string]interface{}{"reference": reference})
	if err != nil {
		return nil, nil, err
	}
	balance, ok := new(big.Int).SetString(res.(map[string]interface{})["balance"].(string), 10)
	if !ok {
		return nil, nil, fmt.Errorf("can't parse balance")
	}
	depositsSliced, ok := res.(map[string]interface{})["deposits"].([]interface{})
	if !ok {
		return balance, nil, fmt.Errorf("can't parse deposits")
	}

	var depositsMap = map[string]interface{}{}
	for _, d := range depositsSliced {
		dMap := d.(map[string]interface{})
		ethTxHash, ok := dMap["ethTxHash"].(string)
		if !ok {
			return balance, nil, fmt.Errorf("can't parse ethTxHash")
		}

		confirmerReferencesSliced, ok := dMap["confirmerReferences"].([]interface{})
		if !ok {
			return balance, nil, fmt.Errorf("can't parse confirmerReferences")
		}

		var confirmerReferences = map[string]interface{}{}
		for _, cr := range confirmerReferencesSliced {
			crMap := cr.(map[string]interface{})
			reference, ok := crMap["reference"].(string)
			if !ok {
				return balance, nil, fmt.Errorf("can't parse reference")
			}
			amount, ok := crMap["amount"]
			if !ok {
				return balance, nil, fmt.Errorf("can't get amount")
			}
			confirmerReferences[reference] = amount
		}

		dMap["confirmerReferences"] = confirmerReferences
		depositsMap[ethTxHash] = dMap
	}

	return balance, depositsMap, nil
}

func migrate(t *testing.T, memberRef string, amount string, tx string, ma string, mdNum int) map[string]interface{} {
	anotherMember := createMember(t)

	_, err := signedRequest(t,
		launchnet.TestRPCUrl,
		launchnet.MigrationDaemons[mdNum],
		"deposit.migration",
		map[string]interface{}{"amount": amount, "ethTxHash": tx, "migrationAddress": ma})
	require.NoError(t, err)
	_, deposits := getBalanceAndDepositsNoErr(t, anotherMember, memberRef)
	deposit, ok := deposits[tx].(map[string]interface{})
	require.True(t, ok)
	confirmations := deposit["confirmerReferences"].(map[string]interface{})
	require.Equal(t, amount+"0", confirmations[launchnet.MigrationDaemons[mdNum].Ref])

	return deposit
}

const migrationAmount = "360000"

func fullMigration(t *testing.T, txHash string) *launchnet.User {
	activeDaemons := activateDaemons(t, countTwoActiveDaemon)

	member := createMigrationMemberForMA(t)
	for i := range activeDaemons {
		migrate(t, member.Ref, migrationAmount, txHash, member.MigrationAddress, i)
	}
	return member
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
		"id":      "",
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
		"id":      "",
	})
	rpcStatusResponse := &rpcStatusResponse{}
	unmarshalRPCResponse(t, body, rpcStatusResponse)
	require.NotNil(t, rpcStatusResponse.Result)
	return rpcStatusResponse.Result
}

func activateDaemons(t *testing.T, countDaemon int) []*launchnet.User {
	var activeDaemons []*launchnet.User
	for i := 0; i < countDaemon; i++ {
		if len(launchnet.MigrationDaemons[i].Ref) > 0 {
			res, err := signedRequest(t, launchnet.TestRPCUrl, &launchnet.MigrationAdmin, "migration.checkDaemon",
				map[string]interface{}{"reference": launchnet.MigrationDaemons[i].Ref})
			require.NoError(t, err)

			status := res.(map[string]interface{})["status"].(string)

			if status == "inactive" {
				_, err := signedRequest(t, launchnet.TestRPCUrl, &launchnet.MigrationAdmin,
					"migration.activateDaemon", map[string]interface{}{"reference": launchnet.MigrationDaemons[i].Ref})
				require.NoError(t, err)
			}
			activeDaemons = append(activeDaemons, launchnet.MigrationDaemons[i])
		}
	}
	return activeDaemons
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

	var errMsg string
	if err != nil {
		var suffix string
		requesterError, ok := err.(*requester.Error)
		if ok {
			suffix = " [" + strings.Join(requesterError.Data.Trace, ": ") + "]"
		}
		t.Error(err.Error() + suffix)
	}
	require.NotEqual(t, "", refStr, "request ref is empty: %s", errMsg)
	require.NotEqual(t, insolar.NewEmptyReference().String(), refStr, "request ref is zero: %s", errMsg)

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
	uploadBody := getRPSResponseBody(t, launchnet.TestRPCUrl, postParams{
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

	prototypeRef, err := insolar.NewReferenceFromString(uploadRes.Result.PrototypeRef)
	require.NoError(t, err)
	require.False(t, prototypeRef.IsEmpty())

	return prototypeRef
}

func callConstructor(t testing.TB, prototypeRef *insolar.Reference, method string, args ...interface{}) *insolar.Reference {
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

	objectRef, err := insolar.NewReferenceFromString(callConstructorRes.Result.Object)
	require.NoError(t, err)

	require.NotEqual(t, insolar.NewReferenceFromBytes(make([]byte, insolar.RecordRefSize)), objectRef)

	return objectRef
}

func callConstructorExpectSystemError(t testing.TB, prototypeRef *insolar.Reference, method string, args ...interface{}) string {
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

	callConstructorRes := struct {
		Version string              `json:"jsonrpc"`
		ID      string              `json:"id"`
		Result  api.CallMethodReply `json:"result"`
		Error   json2.Error         `json:"error"`
	}{}

	err = json.Unmarshal(objectBody, &callConstructorRes)
	require.NoError(t, err)
	require.NotEmpty(t, callConstructorRes.Error)

	return callConstructorRes.Error.Message
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
		res, _, err := makeSignedRequest(launchnet.TestRPCUrlPublic, daemon, "member.get", nil)
		if err != nil {
			return errors.Wrap(err, "[ setup ] get member by public key failed ,key ")
		}
		launchnet.MigrationDaemons[i].Ref = res.(map[string]interface{})["reference"].(string)
	}
	return nil
}

func getAddressCount(t *testing.T, startWithIndex int) map[int]int {
	result, err := signedRequest(t, launchnet.TestRPCUrl, &launchnet.MigrationAdmin, "migration.getAddressCount",
		map[string]interface{}{"startWithIndex": startWithIndex})
	require.NoError(t, err)
	resultsSliced, ok := result.([]interface{})
	require.True(t, ok)

	var migrationShardsMap = map[int]int{}
	for _, r := range resultsSliced {
		rMap := r.(map[string]interface{})
		shardIndex, ok := rMap["shardIndex"].(float64)
		require.True(t, ok)
		freeCount, ok := rMap["freeCount"].(float64)
		require.True(t, ok)
		migrationShardsMap[int(shardIndex)] = int(freeCount)
	}
	return migrationShardsMap
}

func verifyFundsMembersAndDeposits(t *testing.T, m *launchnet.User) error {
	res2, err := signedRequest(t, launchnet.TestRPCUrlPublic, m, "member.get", nil)
	if err != nil {
		return err
	}
	decodedRes2, ok := res2.(map[string]interface{})
	m.Ref = decodedRes2["reference"].(string)
	if !ok {
		return errors.New(fmt.Sprintf("failed to decode: expected map[string]interface{}, got %T", res2))
	}
	balance, deposits := getBalanceAndDepositsNoErr(t, m, decodedRes2["reference"].(string))
	if big.NewInt(0).Cmp(balance) != 0 {
		return errors.New("balance should be zero, current value: " + balance.String())
	}
	deposit, ok := deposits["genesis_deposit"].(map[string]interface{})
	if deposit["amount"] != TestDepositAmount {
		return errors.New(fmt.Sprintf("deposit amount should be %s, current value: %s", TestDepositAmount, deposit["amount"]))
	}
	if deposit["balance"] != TestDepositAmount {
		return errors.New(fmt.Sprintf("deposit balance should be %s, current value: %s", TestDepositAmount, deposit["balance"]))
	}
	return nil
}

func verifyFundsMembersExist(t *testing.T, m *launchnet.User) error {
	res2, err := signedRequest(t, launchnet.TestRPCUrlPublic, m, "member.get", nil)
	if err != nil {
		return err
	}
	decodedRes2, ok := res2.(map[string]interface{})
	m.Ref = decodedRes2["reference"].(string)
	if !ok {
		return errors.New(fmt.Sprintf("failed to decode: expected map[string]interface{}, got %T", res2))
	}
	_, deposits := getBalanceAndDepositsNoErr(t, m, decodedRes2["reference"].(string))
	deposit, ok := deposits["genesis_deposit"].(map[string]interface{})
	if deposit["amount"] != TestDepositAmount {
		return errors.New(fmt.Sprintf("deposit amount should be %s, current value: %s", TestDepositAmount, deposit["amount"]))
	}
	return nil
}
