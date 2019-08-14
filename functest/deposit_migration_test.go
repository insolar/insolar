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
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/testutils"
)

func TestMigrationToken(t *testing.T) {
	migrationAddress := testutils.RandomString()
	member := createMigrationMemberForMA(t, migrationAddress)

	deposit := migrate(t, member.ref, "1000", "Test_TxHash", migrationAddress, 0)

	firstMemberBalance := deposit["balance"].(string)
	require.Equal(t, firstMemberBalance, "0")
	firstMABalance := getBalanceNoErr(t, &migrationAdmin, migrationAdmin.ref)
	require.Equal(t, firstMABalance.String(), "100000000000000000000")

	confirmerReferences, ok := deposit["confirmerReferences"].([]interface{})
	require.True(t, ok, fmt.Sprintf("failed to cast result: expected []string, got %T", deposit["confirmerReferences"]))
	require.Equal(t, confirmerReferences[0], migrationDaemons[0].ref)

	deposit = migrate(t, member.ref, "1000", "Test_TxHash", migrationAddress, 2)

	confirmerReferences, ok = deposit["confirmerReferences"].([]interface{})
	require.True(t, ok)
	require.Equal(t, confirmerReferences[0], migrationDaemons[0].ref)
	require.Equal(t, confirmerReferences[2], migrationDaemons[2].ref)

	deposit = migrate(t, member.ref, "1000", "Test_TxHash", migrationAddress, 1)

	confirmerReferences, ok = deposit["confirmerReferences"].([]interface{})
	require.True(t, ok)
	require.Equal(t, confirmerReferences[0], migrationDaemons[0].ref)
	require.Equal(t, confirmerReferences[1], migrationDaemons[1].ref)
	require.Equal(t, confirmerReferences[2], migrationDaemons[2].ref)

	secondMemberBalance := deposit["balance"].(string)
	require.Equal(t, secondMemberBalance, "1000")
	secondMABalance1 := getBalanceNoErr(t, &migrationAdmin, migrationAdmin.ref)
	require.Equal(t, secondMABalance1.String(), "99999999999999999000")
}

func TestMigrationTokenOnDifferentDeposits(t *testing.T) {
	migrationAddress := testutils.RandomString()
	member := createMigrationMemberForMA(t, migrationAddress)

	deposit := migrate(t, member.ref, "1000", "Test_TxHash1", migrationAddress, 1)

	confirmerReferences, ok := deposit["confirmerReferences"].([]interface{})
	require.True(t, ok, fmt.Sprintf("failed to cast result: expected []string, got %T", deposit["confirmerReferences"]))
	require.Equal(t, confirmerReferences[1], migrationDaemons[1].ref)

	deposit = migrate(t, member.ref, "1000", "Test_TxHash2", migrationAddress, 1)

	confirmerReferences, ok = deposit["confirmerReferences"].([]interface{})
	require.True(t, ok)
	require.Equal(t, confirmerReferences[1], migrationDaemons[1].ref)
}

func TestMigrationTokenNotInTheList(t *testing.T) {
	migrationAddress := generateMigrationAddress()
	_, err := signedRequestWithEmptyRequestRef(t, &migrationAdmin,
		"deposit.migration",
		map[string]interface{}{"amount": "1000", "ethTxHash": "TxHash", "migrationAddress": migrationAddress})
	require.Error(t, err)
	require.Contains(t, err.Error(), "this migration daemon is not in the list")
}

func TestMigrationTokenZeroAmount(t *testing.T) {
	migrationAddress := generateMigrationAddress()
	_ = createMigrationMemberForMA(t, migrationAddress)

	result, err := signedRequestWithEmptyRequestRef(t,
		&migrationDaemons[0],
		"deposit.migration",
		map[string]interface{}{"amount": "0", "ethTxHash": "TxHash", "migrationAddress": migrationAddress})

	require.Error(t, err)
	require.Contains(t, err.Error(), "amount must be greater than zero")
	require.Nil(t, result)

}

func TestMigrationTokenMistakeField(t *testing.T) {
	migrationAddress := generateMigrationAddress()
	_ = createMigrationMemberForMA(t, migrationAddress)

	result, err := signedRequestWithEmptyRequestRef(t,
		&migrationDaemons[0],
		"deposit.migration",
		map[string]interface{}{"amount1": "0", "ethTxHash": "TxHash", "migrationAddress": migrationAddress})
	require.Error(t, err)
	require.Contains(t, err.Error(), " incorect input: failed to get 'amount' param")
	require.Nil(t, result)
}

func TestMigrationTokenNilValue(t *testing.T) {
	migrationAddress := generateMigrationAddress()
	_ = createMigrationMemberForMA(t, migrationAddress)

	result, err := signedRequestWithEmptyRequestRef(t, &migrationDaemons[0], "deposit.migration", map[string]interface{}{"amount": "20", "ethTxHash": nil, "migrationAddress": migrationAddress})
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to get 'ethTxHash' param")
	require.Nil(t, result)

}

func TestMigrationTokenMaxAmount(t *testing.T) {
	migrationAddress := generateMigrationAddress()
	member := createMigrationMemberForMA(t, migrationAddress)

	result, err := signedRequest(t,
		&migrationDaemons[0],
		"deposit.migration",
		map[string]interface{}{"amount": "500000000000000000", "ethTxHash": "ethTxHash", "migrationAddress": migrationAddress})
	require.NoError(t, err)
	require.Equal(t, result.(map[string]interface{})["memberReference"].(string), member.ref)
}

func TestMigrationDoubleMigrationFromSameDaemon(t *testing.T) {
	migrationAddress := generateMigrationAddress()
	member := createMigrationMemberForMA(t, migrationAddress)

	resultMigr1, err := signedRequest(t,
		&migrationDaemons[0], "deposit.migration", map[string]interface{}{"amount": "20", "ethTxHash": "ethTxHash", "migrationAddress": migrationAddress})
	require.NoError(t, err)
	require.Equal(t, resultMigr1.(map[string]interface{})["memberReference"].(string), member.ref)

	_, err = signedRequestWithEmptyRequestRef(t,
		&migrationDaemons[0],
		"deposit.migration",
		map[string]interface{}{"amount": "20", "ethTxHash": "ethTxHash", "migrationAddress": migrationAddress})
	require.Error(t, err)
	require.Contains(t, err.Error(), "confirmed failed: confirm from the")
}

func TestMigrationAnotherAmountSameTx(t *testing.T) {
	migrationAddress := generateMigrationAddress()
	member := createMigrationMemberForMA(t, migrationAddress)

	resultMigr1, err := signedRequest(t,
		&migrationDaemons[0], "deposit.migration", map[string]interface{}{"amount": "20", "ethTxHash": "ethTxHash", "migrationAddress": migrationAddress})
	require.NoError(t, err)
	require.Equal(t, resultMigr1.(map[string]interface{})["memberReference"].(string), member.ref)

	_, err = signedRequestWithEmptyRequestRef(t,
		&migrationDaemons[1],
		"deposit.migration",
		map[string]interface{}{"amount": "30", "ethTxHash": "ethTxHash", "migrationAddress": migrationAddress})
	require.Error(t, err)
	require.Contains(t, err.Error(), "deposit with this transaction hash has different amount")
}
