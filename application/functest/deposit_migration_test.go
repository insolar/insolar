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
	"fmt"
	"math/big"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/application/testutils/launchnet"
)

func TestMigrationToken(t *testing.T) {
	activeDaemons := activateDaemons(t, countTwoActiveDaemon)
	member := createMigrationMemberForMA(t)

	deposit := migrate(t, member.Ref, "1000", "Test_TxHash", member.MigrationAddress, 0)
	firstMemberBalance := deposit["balance"].(string)

	require.Equal(t, "0", firstMemberBalance)
	firstMABalance, err := getAdminDepositBalance(t, &launchnet.MigrationAdmin, launchnet.MigrationAdmin.Ref)
	require.NoError(t, err)

	for i := 1; i < len(activeDaemons); i++ {
		deposit = migrate(t, member.Ref, "1000", "Test_TxHash", member.MigrationAddress, i)
	}

	confirmations := deposit["confirmerReferences"].(map[string]interface{})

	for _, daemons := range activeDaemons {
		require.Equal(t, confirmations[daemons.Ref], "10000")
	}

	require.Equal(t, deposit["ethTxHash"], "Test_TxHash")
	require.Equal(t, deposit["amount"], "10000")

	secondMemberBalance := deposit["balance"].(string)
	require.Equal(t, "10000", secondMemberBalance)
	secondMABalance, err := getAdminDepositBalance(t, &launchnet.MigrationAdmin, launchnet.MigrationAdmin.Ref)
	require.NoError(t, err)

	dif := new(big.Int).Sub(firstMABalance, secondMABalance)
	require.Equal(t, "10000", dif.String())
}

func TestMigrationTokenThreeActiveDaemons(t *testing.T) {
	activeDaemons := activateDaemons(t, countThreeActiveDaemon)
	member := createMigrationMemberForMA(t)
	for i := 0; i < len(activeDaemons)-1; i++ {
		_ = migrate(t, member.Ref, "1000", "Test_TxHash", member.MigrationAddress, i)
	}

	_, err := signedRequest(t,
		launchnet.TestRPCUrl,
		launchnet.MigrationDaemons[countThreeActiveDaemon-1],
		"deposit.migration",
		map[string]interface{}{"amount": "1000", "ethTxHash": "Test_TxHash", "migrationAddress": member.MigrationAddress})
	require.NoError(t, err)
}

func TestMigrationTokenOnDifferentDeposits(t *testing.T) {
	activateDaemons(t, countTwoActiveDaemon)
	member := createMigrationMemberForMA(t)

	_ = migrate(t, member.Ref, "1000", "Test_TxHash", member.MigrationAddress, 0)
	deposit := migrate(t, member.Ref, "1000", "Test_TxHash", member.MigrationAddress, 1)

	confirmations := deposit["confirmerReferences"].(map[string]interface{})
	require.Equal(t, confirmations[launchnet.MigrationDaemons[0].Ref], "10000")
	require.Equal(t, confirmations[launchnet.MigrationDaemons[1].Ref], "10000")
}

func TestMigrationTokenNotInTheList(t *testing.T) {
	migrationAddress, err := generateMigrationAddress()
	require.NoError(t, err)
	_, err = signedRequestWithEmptyRequestRef(t, launchnet.TestRPCUrl,
		&launchnet.MigrationAdmin,
		"deposit.migration",
		map[string]interface{}{"amount": "1000", "ethTxHash": "TxHash", "migrationAddress": migrationAddress})
	data := checkConvertRequesterError(t, err).Data
	require.Contains(t, data.Trace, "the member is not migration daemon")
}

func TestMigrationTokenZeroAmount(t *testing.T) {
	member := createMigrationMemberForMA(t)

	result, err := signedRequestWithEmptyRequestRef(t,
		launchnet.TestRPCUrl,
		launchnet.MigrationDaemons[0],
		"deposit.migration",
		map[string]interface{}{"amount": "0", "ethTxHash": "TxHash", "migrationAddress": member.MigrationAddress})

	data := checkConvertRequesterError(t, err).Data
	require.Contains(t, data.Trace, "amount must be greater than zero")
	require.Nil(t, result)
}

func TestMigrationTokenMistakeField(t *testing.T) {
	member := createMigrationMemberForMA(t)

	result, err := signedRequestWithEmptyRequestRef(t,
		launchnet.TestRPCUrl,
		launchnet.MigrationDaemons[0],
		"deposit.migration",
		map[string]interface{}{"amount1": "0", "ethTxHash": "TxHash", "migrationAddress": member.MigrationAddress})
	data := checkConvertRequesterError(t, err).Data
	require.Contains(t, data.Trace, "failed to get 'amount' param")
	require.Nil(t, result)
}

func TestMigrationTokenNilValue(t *testing.T) {
	member := createMigrationMemberForMA(t)

	result, err := signedRequestWithEmptyRequestRef(t, launchnet.TestRPCUrl, launchnet.MigrationDaemons[0],
		"deposit.migration", map[string]interface{}{"amount": "20", "ethTxHash": nil, "migrationAddress": member.MigrationAddress})
	data := checkConvertRequesterError(t, err).Data
	require.Contains(t, data.Trace, "failed to get 'ethTxHash' param")
	require.Nil(t, result)

}

func TestMigrationTokenMaxAmount(t *testing.T) {
	activateDaemons(t, countTwoActiveDaemon)
	member := createMigrationMemberForMA(t)

	result, err := signedRequest(t,
		launchnet.TestRPCUrl,
		launchnet.MigrationDaemons[0],
		"deposit.migration",
		map[string]interface{}{"amount": "500000000000000000", "ethTxHash": "ethTxHash", "migrationAddress": member.MigrationAddress})
	require.NoError(t, err)
	require.Equal(t, result.(map[string]interface{})["memberReference"].(string), member.Ref)
}

func TestMigrationDoubleMigrationFromSameDaemon(t *testing.T) {
	activateDaemons(t, countTwoActiveDaemon)
	member := createMigrationMemberForMA(t)

	resultMigr1, err := signedRequest(t,
		launchnet.TestRPCUrl,
		launchnet.MigrationDaemons[0], "deposit.migration",
		map[string]interface{}{"amount": "20", "ethTxHash": "ethTxHash", "migrationAddress": member.MigrationAddress})
	require.NoError(t, err)
	require.Equal(t, resultMigr1.(map[string]interface{})["memberReference"].(string), member.Ref)

	resultMigr2, err := signedRequest(t,
		launchnet.TestRPCUrl,
		launchnet.MigrationDaemons[0], "deposit.migration",
		map[string]interface{}{"amount": "20", "ethTxHash": "ethTxHash", "migrationAddress": member.MigrationAddress})
	require.NoError(t, err)
	require.Equal(t, resultMigr2.(map[string]interface{})["memberReference"].(string), member.Ref)
}

func TestMigrationDoubleMigrationFromSameDaemon_WithDifferentAmount(t *testing.T) {
	activateDaemons(t, countTwoActiveDaemon)
	member := createMigrationMemberForMA(t)

	resultMigr1, err := signedRequest(t,
		launchnet.TestRPCUrl,
		launchnet.MigrationDaemons[0], "deposit.migration",
		map[string]interface{}{"amount": "20", "ethTxHash": "ethTxHash", "migrationAddress": member.MigrationAddress})
	require.NoError(t, err)
	require.Equal(t, resultMigr1.(map[string]interface{})["memberReference"].(string), member.Ref)

	_, err = signedRequestWithEmptyRequestRef(t,
		launchnet.TestRPCUrl,
		launchnet.MigrationDaemons[0],
		"deposit.migration",
		map[string]interface{}{"amount": "30", "ethTxHash": "ethTxHash", "migrationAddress": member.MigrationAddress})
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

	_, err := signedRequest(t,
		launchnet.TestRPCUrl,
		launchnet.MigrationDaemons[0], "deposit.migration",
		map[string]interface{}{"amount": "20", "ethTxHash": "ethTxHash", "migrationAddress": member.MigrationAddress})
	require.NoError(t, err)

	_, _, err = makeSignedRequest(
		launchnet.TestRPCUrl,
		launchnet.MigrationDaemons[2],
		"deposit.migration",
		map[string]interface{}{"amount": "30", "ethTxHash": "ethTxHash", "migrationAddress": member.MigrationAddress})
	data := checkConvertRequesterError(t, err).Data
	require.Contains(t, data.Trace, "failed to check amount in confirmation from migration daemon")
	require.Contains(t, data.Trace, fmt.Sprintf("%s send amount 200", launchnet.MigrationDaemons[0].Ref))
	require.Contains(t, data.Trace, fmt.Sprintf("%s send amount 300", launchnet.MigrationDaemons[2].Ref))
}

func TestMigration_WrongSecondAMount(t *testing.T) {
	activateDaemons(t, countThreeActiveDaemon)

	member := createMigrationMemberForMA(t)

	_, err := signedRequest(t,
		launchnet.TestRPCUrl,
		launchnet.MigrationDaemons[0], "deposit.migration",
		map[string]interface{}{"amount": "100", "ethTxHash": "ethTxHash", "migrationAddress": member.MigrationAddress})
	require.NoError(t, err)

	_, _, err = makeSignedRequest(
		launchnet.TestRPCUrl,
		launchnet.MigrationDaemons[1],
		"deposit.migration",
		map[string]interface{}{"amount": "200", "ethTxHash": "ethTxHash", "migrationAddress": member.MigrationAddress})
	data := checkConvertRequesterError(t, err).Data
	require.Contains(t, data.Trace, "several migration daemons send different amount")
	require.Contains(t, data.Trace, fmt.Sprintf("%s send amount 1000", launchnet.MigrationDaemons[0].Ref))
	require.Contains(t, data.Trace, fmt.Sprintf("%s send amount 2000", launchnet.MigrationDaemons[1].Ref))

	_, _, err = makeSignedRequest(
		launchnet.TestRPCUrl,
		launchnet.MigrationDaemons[2],
		"deposit.migration",
		map[string]interface{}{"amount": "100", "ethTxHash": "ethTxHash", "migrationAddress": member.MigrationAddress})
	data = checkConvertRequesterError(t, err).Data
	require.Contains(t, data.Trace, "several migration daemons send different amount")
	require.Contains(t, data.Trace, fmt.Sprintf("%s send amount 2000", launchnet.MigrationDaemons[1].Ref))

	_, deposits := getBalanceAndDepositsNoErr(t, member, member.Ref)
	deposit, ok := deposits["ethTxHash"].(map[string]interface{})
	require.True(t, ok)
	require.Equal(t, deposit["ethTxHash"], "ethTxHash")
	require.Equal(t, deposit["amount"], "1000")
	memberBalance := deposit["balance"].(string)
	require.Equal(t, "1000", memberBalance)
	confirmations := deposit["confirmerReferences"].(map[string]interface{})
	require.Equal(t, confirmations[launchnet.MigrationDaemons[0].Ref], "1000")
	require.Equal(t, confirmations[launchnet.MigrationDaemons[1].Ref], "2000")
	require.Equal(t, confirmations[launchnet.MigrationDaemons[2].Ref], "1000")
}

func TestMigration_WrongFirstAmount(t *testing.T) {
	activateDaemons(t, countThreeActiveDaemon)

	member := createMigrationMemberForMA(t)

	_, err := signedRequest(t,
		launchnet.TestRPCUrl,
		launchnet.MigrationDaemons[0], "deposit.migration",
		map[string]interface{}{"amount": "200", "ethTxHash": "ethTxHash", "migrationAddress": member.MigrationAddress})
	require.NoError(t, err)

	_, _, err = makeSignedRequest(
		launchnet.TestRPCUrl,
		launchnet.MigrationDaemons[1],
		"deposit.migration",
		map[string]interface{}{"amount": "100", "ethTxHash": "ethTxHash", "migrationAddress": member.MigrationAddress})
	data := checkConvertRequesterError(t, err).Data
	require.Contains(t, data.Trace, "several migration daemons send different amount")
	require.Contains(t, data.Trace, fmt.Sprintf("%s send amount 2000", launchnet.MigrationDaemons[0].Ref))
	require.Contains(t, data.Trace, fmt.Sprintf("%s send amount 1000", launchnet.MigrationDaemons[1].Ref))

	_, _, err = makeSignedRequest(
		launchnet.TestRPCUrl,
		launchnet.MigrationDaemons[2],
		"deposit.migration",
		map[string]interface{}{"amount": "100", "ethTxHash": "ethTxHash", "migrationAddress": member.MigrationAddress})
	data = checkConvertRequesterError(t, err).Data
	require.Contains(t, data.Trace, "several migration daemons send different amount")
	require.Contains(t, data.Trace, fmt.Sprintf("%s send amount 2000", launchnet.MigrationDaemons[0].Ref))

	_, deposits := getBalanceAndDepositsNoErr(t, member, member.Ref)
	deposit, ok := deposits["ethTxHash"].(map[string]interface{})
	require.True(t, ok)
	require.Equal(t, deposit["ethTxHash"], "ethTxHash")
	require.Equal(t, deposit["amount"], "1000")
	memberBalance := deposit["balance"].(string)
	require.Equal(t, "1000", memberBalance)
	confirmations := deposit["confirmerReferences"].(map[string]interface{})
	require.Equal(t, confirmations[launchnet.MigrationDaemons[0].Ref], "2000")
	require.Equal(t, confirmations[launchnet.MigrationDaemons[1].Ref], "1000")
	require.Equal(t, confirmations[launchnet.MigrationDaemons[2].Ref], "1000")
}

func TestMigration_WrongThirdAmount(t *testing.T) {
	activateDaemons(t, countThreeActiveDaemon)

	member := createMigrationMemberForMA(t)

	_, err := signedRequest(t,
		launchnet.TestRPCUrl,
		launchnet.MigrationDaemons[0], "deposit.migration",
		map[string]interface{}{"amount": "100", "ethTxHash": "ethTxHash", "migrationAddress": member.MigrationAddress})
	require.NoError(t, err)

	_, err = signedRequest(t,
		launchnet.TestRPCUrl,
		launchnet.MigrationDaemons[1], "deposit.migration",
		map[string]interface{}{"amount": "100", "ethTxHash": "ethTxHash", "migrationAddress": member.MigrationAddress})
	require.NoError(t, err)

	_, _, err = makeSignedRequest(
		launchnet.TestRPCUrl,
		launchnet.MigrationDaemons[2],
		"deposit.migration",
		map[string]interface{}{"amount": "200", "ethTxHash": "ethTxHash", "migrationAddress": member.MigrationAddress})
	data := checkConvertRequesterError(t, err).Data
	require.Contains(t, data.Trace, "migration is done for this deposit ethTxHash, but with different amount")

	_, deposits := getBalanceAndDepositsNoErr(t, member, member.Ref)
	deposit, ok := deposits["ethTxHash"].(map[string]interface{})
	require.True(t, ok)
	require.Equal(t, deposit["ethTxHash"], "ethTxHash")
	require.Equal(t, deposit["amount"], "1000")
	memberBalance := deposit["balance"].(string)
	require.Equal(t, "1000", memberBalance)
	confirmations := deposit["confirmerReferences"].(map[string]interface{})
	require.Equal(t, confirmations[launchnet.MigrationDaemons[0].Ref], "1000")
	require.Equal(t, confirmations[launchnet.MigrationDaemons[1].Ref], "1000")
	require.Equal(t, confirmations[launchnet.MigrationDaemons[2].Ref], "2000")
}

func TestMigration_WrongAllAmount(t *testing.T) {
	activateDaemons(t, countThreeActiveDaemon)

	member := createMigrationMemberForMA(t)

	_, err := signedRequest(t,
		launchnet.TestRPCUrl,
		launchnet.MigrationDaemons[0], "deposit.migration",
		map[string]interface{}{"amount": "100", "ethTxHash": "ethTxHash", "migrationAddress": member.MigrationAddress})
	require.NoError(t, err)

	_, _, err = makeSignedRequest(
		launchnet.TestRPCUrl,
		launchnet.MigrationDaemons[1],
		"deposit.migration",
		map[string]interface{}{"amount": "200", "ethTxHash": "ethTxHash", "migrationAddress": member.MigrationAddress})
	data := checkConvertRequesterError(t, err).Data
	require.Contains(t, data.Trace, "several migration daemons send different amount")
	require.Contains(t, data.Trace, fmt.Sprintf("%s send amount 1000", launchnet.MigrationDaemons[0].Ref))
	require.Contains(t, data.Trace, fmt.Sprintf("%s send amount 2000", launchnet.MigrationDaemons[1].Ref))

	_, _, err = makeSignedRequest(
		launchnet.TestRPCUrl,
		launchnet.MigrationDaemons[2],
		"deposit.migration",
		map[string]interface{}{"amount": "300", "ethTxHash": "ethTxHash", "migrationAddress": member.MigrationAddress})
	data = checkConvertRequesterError(t, err).Data
	require.Contains(t, data.Trace, "several migration daemons send different amount")
	require.Contains(t, data.Trace, fmt.Sprintf("%s send amount 1000", launchnet.MigrationDaemons[0].Ref))
	require.Contains(t, data.Trace, fmt.Sprintf("%s send amount 2000", launchnet.MigrationDaemons[1].Ref))
	require.Contains(t, data.Trace, fmt.Sprintf("%s send amount 3000", launchnet.MigrationDaemons[2].Ref))

	_, deposits := getBalanceAndDepositsNoErr(t, member, member.Ref)
	fmt.Println(deposits)
	deposit, ok := deposits["ethTxHash"].(map[string]interface{})
	require.True(t, ok)
	require.Equal(t, deposit["ethTxHash"], "ethTxHash")
	require.Equal(t, deposit["amount"], "0")
	memberBalance := deposit["balance"].(string)
	require.Equal(t, "0", memberBalance)
	confirmations := deposit["confirmerReferences"].(map[string]interface{})
	require.Equal(t, confirmations[launchnet.MigrationDaemons[0].Ref], "1000")
	require.Equal(t, confirmations[launchnet.MigrationDaemons[1].Ref], "2000")
	require.Equal(t, confirmations[launchnet.MigrationDaemons[2].Ref], "3000")
}

func TestMigrationTokenDoubleSpend(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(2)
	_ = activateDaemons(t, countThreeActiveDaemon)
	member := createMigrationMemberForMA(t)
	anotherMember := createMember(t)

	deposit := migrate(t, member.Ref, "1000", "Test_TxHash", member.MigrationAddress, 0)
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
				map[string]interface{}{"amount": "1000", "ethTxHash": "Test_TxHash", "migrationAddress": member.MigrationAddress})
			if err != nil {
				requestErrorData := checkConvertRequesterError(t, err).Data
				fmt.Println(requestErrorData)
			} else {
				fmt.Println(res)
			}
			wg.Done()
		}(i)
	}
	wg.Wait()

	_, deposits := getBalanceAndDepositsNoErr(t, anotherMember, member.Ref)
	deposit, ok := deposits["Test_TxHash"].(map[string]interface{})
	require.True(t, ok)

	require.Equal(t, deposit["ethTxHash"], "Test_TxHash")
	require.Equal(t, deposit["amount"], "10000")
	secondMemberBalance := deposit["balance"].(string)
	require.Equal(t, "10000", secondMemberBalance)
	secondMABalance, err := getAdminDepositBalance(t, &launchnet.MigrationAdmin, launchnet.MigrationAdmin.Ref)
	require.NoError(t, err)
	dif := new(big.Int).Sub(firstMABalance, secondMABalance)
	require.Equal(t, "10000", dif.String())
}
