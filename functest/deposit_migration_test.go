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
	member, err := newUserWithKeys()
	require.NoError(t, err)
	ba := testutils.RandomString()
	_, err = signedRequest(&migrationAdmin, "migration.addBurnAddresses", map[string]interface{}{"burnAddresses": []string{ba}})
	require.NoError(t, err)
	_, err = retryableMemberMigrationCreate(member, true)
	require.NoError(t, err)

	_, err = signedRequest(
		&migrationDaemons[0],
		"deposit.migration",
		map[string]interface{}{"amount": "1000", "ethTxHash": "Test_TxHash", "migrationAddress": ba})
	require.NoError(t, err)

	anotherMember := createMember(t)
	res, err := signedRequest(anotherMember, "wallet.getBalance", map[string]interface{}{"reference": member.ref})
	require.NoError(t, err)
	deposits, ok := res.(map[string]interface{})["deposits"].([]interface{})
	require.True(t, ok)
	deposit, ok := deposits[0].(map[string]interface{})
	require.True(t, ok)
	require.Equal(t, deposit["amount"], "1000")
	require.Equal(t, deposit["ethTxHash"], "Test_TxHash")
	require.Equal(t, deposit["status"], "Open")
	require.Equal(t, deposit["confirms"], float64(1))

	migrationDaemonConfirms, ok := deposit["migrationDaemonConfirms"].([]interface{})
	require.True(t, ok, fmt.Sprintf("failed to cast result: expected []string, got %T", deposit["migrationDaemonConfirms"]))
	require.Equal(t, migrationDaemonConfirms[0], migrationDaemons[0].ref)

	_, err = signedRequest(
		&migrationDaemons[1],
		"deposit.migration",
		map[string]interface{}{"amount": "1000", "ethTxHash": "Test_TxHash", "migrationAddress": ba})
	require.NoError(t, err)
	res, err = signedRequest(anotherMember, "wallet.getBalance", map[string]interface{}{"reference": member.ref})
	require.NoError(t, err)
	deposits, ok = res.(map[string]interface{})["deposits"].([]interface{})
	require.True(t, ok)
	deposit, ok = deposits[0].(map[string]interface{})
	require.True(t, ok)
	require.Equal(t, deposit["amount"], "1000")
	require.Equal(t, deposit["ethTxHash"], "Test_TxHash")
	require.Equal(t, deposit["status"], "Open")
	require.Equal(t, deposit["confirms"], float64(2))

	migrationDaemonConfirms, ok = deposit["migrationDaemonConfirms"].([]interface{})
	require.True(t, ok)
	require.Equal(t, migrationDaemonConfirms[0], migrationDaemons[0].ref)
	require.Equal(t, migrationDaemonConfirms[1], migrationDaemons[1].ref)

	_, err = signedRequest(
		&migrationDaemons[2],
		"deposit.migration",
		map[string]interface{}{"amount": "1000", "ethTxHash": "Test_TxHash", "migrationAddress": ba})
	require.NoError(t, err)
	res, err = signedRequest(anotherMember, "wallet.getBalance", map[string]interface{}{"reference": member.ref})
	require.NoError(t, err)
	deposits, ok = res.(map[string]interface{})["deposits"].([]interface{})
	require.True(t, ok)
	deposit, ok = deposits[0].(map[string]interface{})
	require.True(t, ok)
	require.Equal(t, deposit["amount"], "1000")
	require.Equal(t, deposit["ethTxHash"], "Test_TxHash")
	require.Equal(t, deposit["status"], "Holding")
	require.Equal(t, deposit["confirms"], float64(3))

	migrationDaemonConfirms, ok = deposit["migrationDaemonConfirms"].([]interface{})
	require.True(t, ok)
	require.Equal(t, migrationDaemonConfirms[0], migrationDaemons[0].ref)
	require.Equal(t, migrationDaemonConfirms[1], migrationDaemons[1].ref)
	require.Equal(t, migrationDaemonConfirms[2], migrationDaemons[2].ref)
}

func TestMigrationTokenNotInTheList(t *testing.T) {
	ba := testutils.RandomString()
	_, err := signedRequest(&migrationAdmin,
		"deposit.migration",
		map[string]interface{}{"amount": "1000", "ethTxHash": "TxHash", "migrationAddress": ba})
	require.Contains(t, err.Error(), "this migration daemon is not in the list")
}

func TestMigrationTokenZeroAmount(t *testing.T) {
	ba := testutils.RandomString()
	err := createMemberWithMigrationAddress(ba)
	require.NoError(t, err)

	result, err := signedRequest(
		&migrationDaemons[0],
		"deposit.migration",
		map[string]interface{}{"amount": "0", "ethTxHash": "TxHash", "migrationAddress": ba})

	require.Contains(t, err.Error(), "amount must be greater than zero")
	require.Nil(t, result)

}

func TestMigrationTokenMistakeField(t *testing.T) {
	ba := testutils.RandomString()
	err := createMemberWithMigrationAddress(ba)
	require.NoError(t, err)

	result, err := signedRequest(
		&migrationDaemons[0],
		"deposit.migration",
		map[string]interface{}{"amount1": "0", "ethTxHash": "TxHash", "migrationAddress": ba})
	require.Contains(t, err.Error(), " incorect input: failed to get 'amount' param")
	require.Nil(t, result)
}

func TestMigrationTokenNilValue(t *testing.T) {
	ba := testutils.RandomString()
	err := createMemberWithMigrationAddress(ba)
	require.NoError(t, err)

	result, err := signedRequest(&migrationDaemons[0], "deposit.migration", map[string]interface{}{"amount": "20", "ethTxHash": nil, "migrationAddress": ba})
	require.Contains(t, err.Error(), "failed to get 'ethTxHash' param")
	require.Nil(t, result)

}

func TestMigrationTokenMaxAmount(t *testing.T) {
	ba := testutils.RandomString()
	err := createMemberWithMigrationAddress(ba)
	require.NoError(t, err)

	result, err := signedRequest(
		&migrationDaemons[0],
		"deposit.migration",
		map[string]interface{}{"amount": "500000000000000000", "ethTxHash": "ethTxHash", "migrationAddress": ba})
	require.NoError(t, err)
	require.Nil(t, result)
}

func TestMigrationDoubleMigration(t *testing.T) {
	ba := testutils.RandomString()
	err := createMemberWithMigrationAddress(ba)
	require.NoError(t, err)

	resultMigr1, err := signedRequest(
		&migrationDaemons[0], "deposit.migration", map[string]interface{}{"amount": "20", "ethTxHash": "ethTxHash", "migrationAddress": ba})
	require.NoError(t, err)
	require.Nil(t, resultMigr1)

	_, err = signedRequest(
		&migrationDaemons[0],
		"deposit.migration",
		map[string]interface{}{"amount": "20", "ethTxHash": "ethTxHash", "migrationAddress": ba})
	require.Contains(t, err.Error(), "confirmed failed: confirm from the")
}

func TestMigrationAnotherAmountSameTx(t *testing.T) {
	ba := testutils.RandomString()
	err := createMemberWithMigrationAddress(ba)
	require.NoError(t, err)

	resultMigr1, err := signedRequest(
		&migrationDaemons[0], "deposit.migration", map[string]interface{}{"amount": "20", "ethTxHash": "ethTxHash", "migrationAddress": ba})
	require.NoError(t, err)
	require.Nil(t, resultMigr1)

	_, err = signedRequest(
		&migrationDaemons[0],
		"deposit.migration",
		map[string]interface{}{"amount": "30", "ethTxHash": "ethTxHash", "migrationAddress": ba})
	require.Contains(t, err.Error(), "deposit with this transaction hash has different amount")
}

func createMemberWithMigrationAddress(migrationAddress string) error {
	member, err := newUserWithKeys()
	if err != nil {
		return err
	}

	_, err = signedRequest(&migrationAdmin, "migration.addBurnAddresses", map[string]interface{}{"burnAddresses": []string{migrationAddress}})
	if err != nil {
		return err
	}
	_, err = retryableMemberMigrationCreate(member, true)
	if err != nil {
		return err
	}
	return nil
}
