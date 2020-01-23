// Copyright 2020 Insolar Network Ltd.
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

// +build functest

package functest

import (
	"fmt"
	"math/big"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/application/testutils/launchnet"
	"github.com/insolar/insolar/testutils"
)

func TestMigrationToken(t *testing.T) {
	activeDaemons := activateDaemons(t, countTwoActiveDaemon)
	member := createMigrationMemberForMA(t)

	ethHash := testutils.RandomEthHash()

	deposit := migrate(t, member.Ref, "1000", ethHash, member.MigrationAddress, 0)
	firstMemberBalance := deposit["balance"].(string)

	require.Equal(t, "0", firstMemberBalance)
	firstMABalance, err := getAdminDepositBalance(t, &launchnet.MigrationAdmin, launchnet.MigrationAdmin.Ref)
	require.NoError(t, err)

	for i := 1; i < len(activeDaemons); i++ {
		deposit = migrate(t, member.Ref, "1000", ethHash, member.MigrationAddress, i)
	}

	confirmations := deposit["confirmerReferences"].(map[string]interface{})

	for _, daemons := range activeDaemons {
		require.Equal(t, "1000", confirmations[daemons.Ref])
	}

	require.Equal(t, ethHash, deposit["ethTxHash"])
	require.Equal(t, "1000", deposit["amount"])

	secondMemberBalance := deposit["balance"].(string)
	require.Equal(t, "1000", secondMemberBalance)
	secondMABalance, err := getAdminDepositBalance(t, &launchnet.MigrationAdmin, launchnet.MigrationAdmin.Ref)
	require.NoError(t, err)

	dif := new(big.Int).Sub(firstMABalance, secondMABalance)
	require.Equal(t, "1000", dif.String())
}

func TestMigrationTokenOneActiveDaemon(t *testing.T) {
	// one daemon confirmation can't change balance
	activateDaemons(t, countOneActiveDaemon)
	daemonIndex := 0
	member := createMigrationMemberForMA(t)

	ethHash := testutils.RandomEthHash()

	deposit := migrate(t, member.Ref, "1000", ethHash, member.MigrationAddress, daemonIndex)
	balance := deposit["balance"].(string)
	require.Equal(t, "0", balance)

	confirmations := deposit["confirmerReferences"].(map[string]interface{})
	require.Equal(t, "1000", confirmations[launchnet.MigrationDaemons[daemonIndex].Ref])
}

func TestMigrationTokenThreeActiveDaemons(t *testing.T) {
	activeDaemons := activateDaemons(t, countThreeActiveDaemon)
	member := createMigrationMemberForMA(t)

	ethHash := testutils.RandomEthHash()

	for i := 0; i < len(activeDaemons)-1; i++ {
		_ = migrate(t, member.Ref, "1000", ethHash, member.MigrationAddress, i)
	}

	_, err := signedRequest(t,
		launchnet.TestRPCUrl,
		launchnet.MigrationDaemons[countThreeActiveDaemon-1],
		"deposit.migration",
		map[string]interface{}{"amount": "1000", "ethTxHash": ethHash, "migrationAddress": member.MigrationAddress})
	require.NoError(t, err)
}

func TestMigrationTokenOnDifferentDeposits(t *testing.T) {
	activateDaemons(t, countTwoActiveDaemon)
	member := createMigrationMemberForMA(t)

	ethHash := testutils.RandomEthHash()

	_ = migrate(t, member.Ref, "1000", ethHash, member.MigrationAddress, 0)
	deposit := migrate(t, member.Ref, "1000", ethHash, member.MigrationAddress, 1)

	confirmations := deposit["confirmerReferences"].(map[string]interface{})
	require.Equal(t, "1000", confirmations[launchnet.MigrationDaemons[0].Ref])
	require.Equal(t, "1000", confirmations[launchnet.MigrationDaemons[1].Ref])
}

func TestMigrationTokenNotInTheList(t *testing.T) {
	migrationAddress := testutils.RandomEthMigrationAddress()
	_, err := signedRequestWithEmptyRequestRef(t, launchnet.TestRPCUrl,
		&launchnet.MigrationAdmin,
		"deposit.migration",
		map[string]interface{}{"amount": "1000", "ethTxHash": testutils.RandomEthHash(), "migrationAddress": migrationAddress})
	data := checkConvertRequesterError(t, err).Data
	require.Contains(t, data.Trace, "the member is not migration daemon")
}

func TestMigrationTokenZeroAmount(t *testing.T) {
	member := createMigrationMemberForMA(t)

	result, err := signedRequestWithEmptyRequestRef(t,
		launchnet.TestRPCUrl,
		launchnet.MigrationDaemons[0],
		"deposit.migration",
		map[string]interface{}{"amount": "0", "ethTxHash": testutils.RandomEthHash(), "migrationAddress": member.MigrationAddress})

	data := checkConvertRequesterError(t, err).Data
	expectedError(t, data.Trace, `Error at "/params/callParams/amount":JSON string doesn't match the regular expression '^[1-9][0-9]*$`)
	require.Nil(t, result)
}

func TestMigrationTokenMistakeField(t *testing.T) {
	member := createMigrationMemberForMA(t)

	result, err := signedRequestWithEmptyRequestRef(t,
		launchnet.TestRPCUrl,
		launchnet.MigrationDaemons[0],
		"deposit.migration",
		map[string]interface{}{"amount1": "0", "ethTxHash": testutils.RandomEthHash(), "migrationAddress": member.MigrationAddress})
	data := checkConvertRequesterError(t, err).Data
	expectedError(t, data.Trace, "Property 'amount' is missing")
	require.Nil(t, result)
}

func TestMigrationTokenNilValue(t *testing.T) {
	member := createMigrationMemberForMA(t)

	result, err := signedRequestWithEmptyRequestRef(t, launchnet.TestRPCUrl, launchnet.MigrationDaemons[0],
		"deposit.migration", map[string]interface{}{"amount": "20", "ethTxHash": nil, "migrationAddress": member.MigrationAddress})
	data := checkConvertRequesterError(t, err).Data
	expectedError(t, data.Trace, `Error at "/params/callParams/ethTxHash":Value is not nullable`)
	require.Nil(t, result)

}

func TestMigrationTokenMaxAmount(t *testing.T) {
	activateDaemons(t, countTwoActiveDaemon)
	member := createMigrationMemberForMA(t)

	result, err := signedRequest(t,
		launchnet.TestRPCUrl,
		launchnet.MigrationDaemons[0],
		"deposit.migration",
		map[string]interface{}{"amount": "500000000000000000", "ethTxHash": testutils.RandomEthHash(), "migrationAddress": member.MigrationAddress})
	require.NoError(t, err)
	require.Equal(t, result.(map[string]interface{})["memberReference"].(string), member.Ref)
}

func TestMigrationDoubleMigrationFromSameDaemon(t *testing.T) {
	activateDaemons(t, countTwoActiveDaemon)
	member := createMigrationMemberForMA(t)

	ethHash := testutils.RandomEthHash()

	resultMigr1, err := signedRequest(t,
		launchnet.TestRPCUrl,
		launchnet.MigrationDaemons[0], "deposit.migration",
		map[string]interface{}{"amount": "20", "ethTxHash": ethHash, "migrationAddress": member.MigrationAddress})
	require.NoError(t, err)
	require.Equal(t, member.Ref, resultMigr1.(map[string]interface{})["memberReference"].(string))

	resultMigr2, err := signedRequest(t,
		launchnet.TestRPCUrl,
		launchnet.MigrationDaemons[0], "deposit.migration",
		map[string]interface{}{"amount": "20", "ethTxHash": ethHash, "migrationAddress": member.MigrationAddress})
	require.NoError(t, err)
	require.Equal(t, member.Ref, resultMigr2.(map[string]interface{})["memberReference"].(string))
}

func TestMigrationDoubleMigrationFromSameDaemon_WithDifferentAmount(t *testing.T) {
	activateDaemons(t, countTwoActiveDaemon)
	member := createMigrationMemberForMA(t)

	ethHash := testutils.RandomEthHash()

	resultMigr1, err := signedRequest(t,
		launchnet.TestRPCUrl,
		launchnet.MigrationDaemons[0], "deposit.migration",
		map[string]interface{}{"amount": "20", "ethTxHash": ethHash, "migrationAddress": member.MigrationAddress})
	require.NoError(t, err)
	require.Equal(t, member.Ref, resultMigr1.(map[string]interface{})["memberReference"].(string))

	_, err = signedRequestWithEmptyRequestRef(t,
		launchnet.TestRPCUrl,
		launchnet.MigrationDaemons[0],
		"deposit.migration",
		map[string]interface{}{"amount": "30", "ethTxHash": ethHash, "migrationAddress": member.MigrationAddress})
	data := checkConvertRequesterError(t, err).Data
	require.Contains(
		t,
		data.Trace,
		fmt.Sprintf("confirm from this migration daemon %s already exists with different amount", launchnet.MigrationDaemons[0].Ref),
	)
}

func TestMigrationAnotherAmountSameTx(t *testing.T) {
	activateDaemons(t, countThreeActiveDaemon)

	member := createMigrationMemberForMA(t)

	ethHash := testutils.RandomEthHash()

	_, err := signedRequest(t,
		launchnet.TestRPCUrl,
		launchnet.MigrationDaemons[0], "deposit.migration",
		map[string]interface{}{"amount": "20", "ethTxHash": ethHash, "migrationAddress": member.MigrationAddress})
	require.NoError(t, err)

	_, _, err = makeSignedRequest(
		launchnet.TestRPCUrl,
		launchnet.MigrationDaemons[2],
		"deposit.migration",
		map[string]interface{}{"amount": "30", "ethTxHash": ethHash, "migrationAddress": member.MigrationAddress})
	data := checkConvertRequesterError(t, err).Data
	require.Contains(t, data.Trace, "failed to check amount in confirmation from migration daemon")
	require.Contains(t, data.Trace, fmt.Sprintf("%s send amount 20", launchnet.MigrationDaemons[0].Ref))
	require.Contains(t, data.Trace, fmt.Sprintf("%s send amount 30", launchnet.MigrationDaemons[2].Ref))
}

func TestMigration_WrongSecondAMount(t *testing.T) {
	activateDaemons(t, countThreeActiveDaemon)

	member := createMigrationMemberForMA(t)

	ethHash := testutils.RandomEthHash()

	_, err := signedRequest(t,
		launchnet.TestRPCUrl,
		launchnet.MigrationDaemons[0], "deposit.migration",
		map[string]interface{}{"amount": "100", "ethTxHash": ethHash, "migrationAddress": member.MigrationAddress})
	require.NoError(t, err)

	_, _, err = makeSignedRequest(
		launchnet.TestRPCUrl,
		launchnet.MigrationDaemons[1],
		"deposit.migration",
		map[string]interface{}{"amount": "200", "ethTxHash": ethHash, "migrationAddress": member.MigrationAddress})
	data := checkConvertRequesterError(t, err).Data
	require.Contains(t, data.Trace, "several migration daemons send different amount")
	require.Contains(t, data.Trace, fmt.Sprintf("%s send amount 100", launchnet.MigrationDaemons[0].Ref))
	require.Contains(t, data.Trace, fmt.Sprintf("%s send amount 200", launchnet.MigrationDaemons[1].Ref))

	_, _, err = makeSignedRequest(
		launchnet.TestRPCUrl,
		launchnet.MigrationDaemons[2],
		"deposit.migration",
		map[string]interface{}{"amount": "100", "ethTxHash": ethHash, "migrationAddress": member.MigrationAddress})
	data = checkConvertRequesterError(t, err).Data
	require.Contains(t, data.Trace, "several migration daemons send different amount")
	require.Contains(t, data.Trace, fmt.Sprintf("%s send amount 200", launchnet.MigrationDaemons[1].Ref))

	_, deposits := getBalanceAndDepositsNoErr(t, member, member.Ref)
	deposit, ok := deposits[ethHash].(map[string]interface{})
	require.True(t, ok)
	require.Equal(t, ethHash, deposit["ethTxHash"])
	require.Equal(t, "100", deposit["amount"])
	memberBalance := deposit["balance"].(string)
	require.Equal(t, "100", memberBalance)
	confirmations := deposit["confirmerReferences"].(map[string]interface{})
	require.Equal(t, "100", confirmations[launchnet.MigrationDaemons[0].Ref])
	require.Equal(t, "200", confirmations[launchnet.MigrationDaemons[1].Ref])
	require.Equal(t, "100", confirmations[launchnet.MigrationDaemons[2].Ref])
}

func TestMigration_WrongFirstAmount(t *testing.T) {
	activateDaemons(t, countThreeActiveDaemon)

	member := createMigrationMemberForMA(t)

	ethHash := testutils.RandomEthHash()

	_, err := signedRequest(t,
		launchnet.TestRPCUrl,
		launchnet.MigrationDaemons[0], "deposit.migration",
		map[string]interface{}{"amount": "200", "ethTxHash": ethHash, "migrationAddress": member.MigrationAddress})
	require.NoError(t, err)

	_, _, err = makeSignedRequest(
		launchnet.TestRPCUrl,
		launchnet.MigrationDaemons[1],
		"deposit.migration",
		map[string]interface{}{"amount": "100", "ethTxHash": ethHash, "migrationAddress": member.MigrationAddress})
	data := checkConvertRequesterError(t, err).Data
	require.Contains(t, data.Trace, "several migration daemons send different amount")
	require.Contains(t, data.Trace, fmt.Sprintf("%s send amount 200", launchnet.MigrationDaemons[0].Ref))
	require.Contains(t, data.Trace, fmt.Sprintf("%s send amount 100", launchnet.MigrationDaemons[1].Ref))

	_, _, err = makeSignedRequest(
		launchnet.TestRPCUrl,
		launchnet.MigrationDaemons[2],
		"deposit.migration",
		map[string]interface{}{"amount": "100", "ethTxHash": ethHash, "migrationAddress": member.MigrationAddress})
	data = checkConvertRequesterError(t, err).Data
	require.Contains(t, data.Trace, "several migration daemons send different amount")
	require.Contains(t, data.Trace, fmt.Sprintf("%s send amount 200", launchnet.MigrationDaemons[0].Ref))

	_, deposits := getBalanceAndDepositsNoErr(t, member, member.Ref)
	deposit, ok := deposits[ethHash].(map[string]interface{})
	require.True(t, ok)
	require.Equal(t, ethHash, deposit["ethTxHash"])
	require.Equal(t, "100", deposit["amount"])
	memberBalance := deposit["balance"].(string)
	require.Equal(t, "100", memberBalance)
	confirmations := deposit["confirmerReferences"].(map[string]interface{})
	require.Equal(t, "200", confirmations[launchnet.MigrationDaemons[0].Ref])
	require.Equal(t, "100", confirmations[launchnet.MigrationDaemons[1].Ref])
	require.Equal(t, "100", confirmations[launchnet.MigrationDaemons[2].Ref])
}

func TestMigration_WrongThirdAmount(t *testing.T) {
	activateDaemons(t, countThreeActiveDaemon)

	member := createMigrationMemberForMA(t)

	ethHash := testutils.RandomEthHash()

	_, err := signedRequest(t,
		launchnet.TestRPCUrl,
		launchnet.MigrationDaemons[0], "deposit.migration",
		map[string]interface{}{"amount": "100", "ethTxHash": ethHash, "migrationAddress": member.MigrationAddress})
	require.NoError(t, err)

	_, err = signedRequest(t,
		launchnet.TestRPCUrl,
		launchnet.MigrationDaemons[1], "deposit.migration",
		map[string]interface{}{"amount": "100", "ethTxHash": ethHash, "migrationAddress": member.MigrationAddress})
	require.NoError(t, err)

	_, _, err = makeSignedRequest(
		launchnet.TestRPCUrl,
		launchnet.MigrationDaemons[2],
		"deposit.migration",
		map[string]interface{}{"amount": "200", "ethTxHash": ethHash, "migrationAddress": member.MigrationAddress})
	data := checkConvertRequesterError(t, err).Data
	require.Contains(t, data.Trace, fmt.Sprintf("migration is done for this deposit %s, but with different amount", ethHash))

	_, deposits := getBalanceAndDepositsNoErr(t, member, member.Ref)
	deposit, ok := deposits[ethHash].(map[string]interface{})
	require.True(t, ok)
	require.Equal(t, ethHash, deposit["ethTxHash"])
	require.Equal(t, "100", deposit["amount"])
	memberBalance := deposit["balance"].(string)
	require.Equal(t, "100", memberBalance)
	confirmations := deposit["confirmerReferences"].(map[string]interface{})
	require.Equal(t, "100", confirmations[launchnet.MigrationDaemons[0].Ref])
	require.Equal(t, "100", confirmations[launchnet.MigrationDaemons[1].Ref])
	require.Equal(t, "200", confirmations[launchnet.MigrationDaemons[2].Ref])
}

func TestMigration_WrongAllAmount(t *testing.T) {
	activateDaemons(t, countThreeActiveDaemon)

	member := createMigrationMemberForMA(t)

	ethHash := testutils.RandomEthHash()

	_, err := signedRequest(t,
		launchnet.TestRPCUrl,
		launchnet.MigrationDaemons[0], "deposit.migration",
		map[string]interface{}{"amount": "100", "ethTxHash": ethHash, "migrationAddress": member.MigrationAddress})
	require.NoError(t, err)

	_, _, err = makeSignedRequest(
		launchnet.TestRPCUrl,
		launchnet.MigrationDaemons[1],
		"deposit.migration",
		map[string]interface{}{"amount": "200", "ethTxHash": ethHash, "migrationAddress": member.MigrationAddress})
	data := checkConvertRequesterError(t, err).Data
	require.Contains(t, data.Trace, "several migration daemons send different amount")
	require.Contains(t, data.Trace, fmt.Sprintf("%s send amount 100", launchnet.MigrationDaemons[0].Ref))
	require.Contains(t, data.Trace, fmt.Sprintf("%s send amount 200", launchnet.MigrationDaemons[1].Ref))

	_, _, err = makeSignedRequest(
		launchnet.TestRPCUrl,
		launchnet.MigrationDaemons[2],
		"deposit.migration",
		map[string]interface{}{"amount": "300", "ethTxHash": ethHash, "migrationAddress": member.MigrationAddress})
	data = checkConvertRequesterError(t, err).Data
	require.Contains(t, data.Trace, "several migration daemons send different amount")
	require.Contains(t, data.Trace, fmt.Sprintf("%s send amount 100", launchnet.MigrationDaemons[0].Ref))
	require.Contains(t, data.Trace, fmt.Sprintf("%s send amount 200", launchnet.MigrationDaemons[1].Ref))
	require.Contains(t, data.Trace, fmt.Sprintf("%s send amount 300", launchnet.MigrationDaemons[2].Ref))

	_, deposits := getBalanceAndDepositsNoErr(t, member, member.Ref)
	deposit, ok := deposits[ethHash].(map[string]interface{})
	require.True(t, ok)
	require.Equal(t, ethHash, deposit["ethTxHash"])
	require.Equal(t, "0", deposit["amount"])
	memberBalance := deposit["balance"].(string)
	require.Equal(t, "0", memberBalance)
	confirmations := deposit["confirmerReferences"].(map[string]interface{})
	require.Equal(t, "100", confirmations[launchnet.MigrationDaemons[0].Ref])
	require.Equal(t, "200", confirmations[launchnet.MigrationDaemons[1].Ref])
	require.Equal(t, "300", confirmations[launchnet.MigrationDaemons[2].Ref])
}

func TestMigrationTokenDoubleSpend(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(2)
	_ = activateDaemons(t, countThreeActiveDaemon)
	member := createMigrationMemberForMA(t)
	anotherMember := createMember(t)

	ethHash := testutils.RandomEthHash()

	deposit := migrate(t, member.Ref, "1000", ethHash, member.MigrationAddress, 0)
	firstMemberBalance := deposit["balance"].(string)

	require.Equal(t, "0", firstMemberBalance)
	firstMABalance, err := getAdminDepositBalance(t, &launchnet.MigrationAdmin, launchnet.MigrationAdmin.Ref)
	require.NoError(t, err)

	for i := 1; i < countThreeActiveDaemon; i++ {
		go func(i int) {
			res, _, err := makeSignedRequest(
				launchnet.TestRPCUrl,
				launchnet.MigrationDaemons[i],
				"deposit.migration",
				map[string]interface{}{"amount": "1000", "ethTxHash": ethHash, "migrationAddress": member.MigrationAddress})
			if err != nil {
				requestErrorData := checkConvertRequesterError(t, err).Data
				t.Log(requestErrorData)
			} else {
				t.Log(res)
			}
			wg.Done()
		}(i)
	}
	wg.Wait()

	_, deposits := getBalanceAndDepositsNoErr(t, anotherMember, member.Ref)
	deposit, ok := deposits[ethHash].(map[string]interface{})
	require.True(t, ok)

	require.Equal(t, ethHash, deposit["ethTxHash"])
	require.Equal(t, "1000", deposit["amount"])
	secondMemberBalance := deposit["balance"].(string)
	require.Equal(t, "1000", secondMemberBalance)
	secondMABalance, err := getAdminDepositBalance(t, &launchnet.MigrationAdmin, launchnet.MigrationAdmin.Ref)
	require.NoError(t, err)
	dif := new(big.Int).Sub(firstMABalance, secondMABalance)
	require.Equal(t, "1000", dif.String())
}
