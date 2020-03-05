// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

// +build functest

package functest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"testing"
	"time"

	"github.com/insolar/insolar/application/genesisrefs"
	"github.com/insolar/insolar/applicationbase/testutils/testrequest"
	"github.com/insolar/insolar/applicationbase/testutils/testresponse"
	"github.com/insolar/insolar/insolar/secrets"

	"github.com/insolar/insolar/api"
	"github.com/insolar/insolar/api/requester"
	"github.com/insolar/insolar/applicationbase/testutils/launchnet"

	"github.com/stretchr/testify/require"

	"github.com/pkg/errors"
)

const (
	countOneActiveDaemon = iota + 1
	countTwoActiveDaemon
	countThreeActiveDaemon
)

const TestDepositAmount string = "1000000000000000000"

type infoResponse struct {
	RootDomain string `json:"RootDomain"`
	RootMember string `json:"RootMember"`
	NodeDomain string `json:"NodeDomain"`
	TraceID    string `json:"TraceID"`
}

type rpcInfoResponse struct {
	testresponse.RPCResponse
	Result infoResponse `json:"result"`
}

func checkConvertRequesterError(t *testing.T, err error) *requester.Error {
	rv, ok := err.(*requester.Error)
	require.Truef(t, ok, "got wrong error %T (expected *requester.Error) with text '%s'", err, err.Error())
	return rv
}

func createMember(t *testing.T) *AppUser {
	member, err := newUserWithKeys()
	require.NoError(t, err)
	member.Ref = Root.Ref

	result, err := testrequest.SignedRequest(t, launchnet.TestRPCUrlPublic, member, "member.create", nil)
	require.NoError(t, err)
	ref, ok := result.(map[string]interface{})["reference"].(string)
	require.True(t, ok)
	member.Ref = ref
	return member
}

func createMigrationMemberForMA(t *testing.T) *AppUser {
	member, err := newUserWithKeys()
	require.NoError(t, err)
	member.Ref = Root.Ref

	result, err := testrequest.SignedRequest(t, launchnet.TestRPCUrlPublic, member, "member.migrationCreate", nil)
	require.NoError(t, err)
	ref, ok := result.(map[string]interface{})["reference"].(string)
	require.True(t, ok)
	member.Ref = ref

	ma, ok := result.(map[string]interface{})["migrationAddress"].(string)
	require.True(t, ok)
	member.MigrationAddress = ma
	return member

}

func getBalanceNoErr(t *testing.T, caller *AppUser, reference string) *big.Int {
	balance, _ := getBalanceAndDepositsNoErr(t, caller, reference)
	return balance
}

func getAdminDepositBalance(t *testing.T, caller *AppUser, reference string) (*big.Int, error) {
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

func getBalanceAndDepositsNoErr(t *testing.T, caller *AppUser, reference string) (*big.Int, map[string]interface{}) {
	balance, deposits, err := getBalanceAndDeposits(t, caller, reference)
	require.NoError(t, err)
	return balance, deposits
}

func getBalanceAndDeposits(t *testing.T, caller *AppUser, reference string) (*big.Int, map[string]interface{}, error) {
	res, err := testrequest.SignedRequest(t, launchnet.TestRPCUrl, caller, "member.getBalance", map[string]interface{}{"reference": reference})
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

	_, err := testrequest.SignedRequest(t,
		launchnet.TestRPCUrl,
		MigrationDaemons[mdNum],
		"deposit.migration",
		map[string]interface{}{"amount": amount, "ethTxHash": tx, "migrationAddress": ma})
	require.NoError(t, err)
	_, deposits := getBalanceAndDepositsNoErr(t, anotherMember, memberRef)
	deposit, ok := deposits[tx].(map[string]interface{})
	require.True(t, ok)
	confirmations := deposit["confirmerReferences"].(map[string]interface{})
	require.Equal(t, amount+"0", confirmations[MigrationDaemons[mdNum].Ref])

	return deposit
}

const migrationAmount = "360000"

func fullMigration(t *testing.T, txHash string) *AppUser {
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

func getInfo(t testing.TB) infoResponse {
	pp := testresponse.PostParams{
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

func activateDaemons(t *testing.T, countDaemon int) []*AppUser {
	var activeDaemons []*AppUser
	for i := 0; i < countDaemon; i++ {
		if len(MigrationDaemons[i].Ref) > 0 {
			res, err := testrequest.SignedRequest(t, launchnet.TestRPCUrl, &MigrationAdmin, "migration.checkDaemon",
				map[string]interface{}{"reference": MigrationDaemons[i].Ref})
			require.NoError(t, err)

			status := res.(map[string]interface{})["status"].(string)

			if status == "inactive" {
				_, err := testrequest.SignedRequest(t, launchnet.TestRPCUrl, &MigrationAdmin,
					"migration.activateDaemon", map[string]interface{}{"reference": MigrationDaemons[i].Ref})
				require.NoError(t, err)
			}
			activeDaemons = append(activeDaemons, MigrationDaemons[i])
		}
	}
	return activeDaemons
}

func unmarshalRPCResponse(t testing.TB, body []byte, response testresponse.RPCResponseInterface) {
	err := json.Unmarshal(body, &response)
	require.NoError(t, err)
	require.Equal(t, "2.0", response.GetRPCVersion())
	require.Nil(t, response.GetError())
}

func unmarshalCallResponse(t testing.TB, body []byte, response *requester.ContractResponse) {
	err := json.Unmarshal(body, &response)
	require.NoError(t, err)
}

func newUserWithKeys() (*AppUser, error) {
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
	return &AppUser{
		PrivKey: string(privKeyStr),
		PubKey:  string(pubKeyStr),
	}, nil
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
	for i, mDaemon := range MigrationDaemons {
		daemon := mDaemon
		daemon.Ref = Root.Ref
		res, _, err := testrequest.MakeSignedRequest(launchnet.TestRPCUrlPublic, daemon, "member.get", nil)
		if err != nil {
			return errors.Wrap(err, "[ setup ] get member by public key failed ,key ")
		}
		MigrationDaemons[i].Ref = res.(map[string]interface{})["reference"].(string)
	}
	return nil
}

func getAddressCount(t *testing.T, startWithIndex int) map[int]int {
	result, err := testrequest.SignedRequest(t, launchnet.TestRPCUrl, &MigrationAdmin, "migration.getAddressCount",
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

func verifyFundsMembersAndDeposits(t *testing.T, m *AppUser, expectedBalance string) error {
	res2, err := testrequest.SignedRequest(t, launchnet.TestRPCUrlPublic, m, "member.get", nil)
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
	if deposit["amount"] != expectedBalance {
		return errors.New(fmt.Sprintf("deposit amount should be %s, current value: %s", expectedBalance, deposit["amount"]))
	}
	if deposit["balance"] != expectedBalance {
		return errors.New(fmt.Sprintf("deposit balance should be %s, current value: %s", expectedBalance, deposit["balance"]))
	}
	return nil
}

func verifyFundsMembersExist(t *testing.T, m *AppUser, expectedBalance string) error {
	res2, err := testrequest.SignedRequest(t, launchnet.TestRPCUrlPublic, m, "member.get", nil)
	if err != nil {
		return err
	}
	decodedRes2, ok := res2.(map[string]interface{})
	m.Ref = decodedRes2["reference"].(string)
	if !ok {
		return errors.New(fmt.Sprintf("failed to decode: expected map[string]interface{}, got %T", res2))
	}
	balance, deposits := getBalanceAndDepositsNoErr(t, m, decodedRes2["reference"].(string))
	require.Equal(t, expectedBalance, balance.String())
	require.Empty(t, deposits)
	return nil
}
