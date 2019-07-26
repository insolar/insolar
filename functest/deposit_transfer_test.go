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
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestDepositTransferToken(t *testing.T) {
	member := fullMigration(t, "Eth_TxHash_test")

	firstBalance := getBalanceNoErr(t, member, member.ref)
	secondBalance := new(big.Int).Add(firstBalance, big.NewInt(100))

	for i := 0; i <= 10; i++ {
		time.Sleep(time.Second)
		_, err := signedRequest(member, "deposit.transfer", map[string]interface{}{"amount": "100", "ethTxHash": "Eth_TxHash_test"})
		if err != nil {
			require.Contains(t, err.Error(), "hold period didn't end")
		} else {
			break
		}
	}

	finalBalance := getBalanceNoErr(t, member, member.ref)

	require.Equal(t, secondBalance, finalBalance)
}

func TestDepositTransferBeforeUnhold(t *testing.T) {
	member := fullMigration(t, "Eth_TxHash_test")

	_, err := signedRequest(member, "deposit.transfer", map[string]interface{}{"amount": "100", "ethTxHash": "Eth_TxHash_test"})
	require.Contains(t, err.Error(), "hold period didn't end")
}

func TestDepositTransferBiggerAmount(t *testing.T) {
	member := fullMigration(t, "Eth_TxHash_test")

	_, err := signedRequest(member, "deposit.transfer", map[string]interface{}{"amount": "1001", "ethTxHash": "Eth_TxHash_test"})
	require.Contains(t, err.Error(), "not enough balance for transfer")
}

func TestDepositTransferAnotherTx(t *testing.T) {
	member := fullMigration(t, "Eth_TxHash_test")

	_, err := signedRequest(member, "deposit.transfer", map[string]interface{}{"amount": "100", "ethTxHash": "Eth_TxHash_foo"})
	require.Contains(t, err.Error(), "can't find deposit")
}

func TestDepositTransferWrongValueAmount(t *testing.T) {
	member := fullMigration(t, "Eth_TxHash_test")

	_, err := signedRequest(member, "deposit.transfer", map[string]interface{}{"amount": "foo", "ethTxHash": "Eth_TxHash_test"})
	require.Contains(t, err.Error(), "can't parse input amount")
}

func TestDepositTransferNotEnoughConfirms(t *testing.T) {
	member, err := newUserWithKeys()
	require.NoError(t, err)
	migrationAddress := generateMigrationAddress()
	_, err = signedRequest(&migrationAdmin, "migration.addBurnAddresses", map[string]interface{}{"burnAddresses": []string{migrationAddress}})
	require.NoError(t, err)
	_, err = retryableMemberMigrationCreate(member, true)
	require.NoError(t, err)

	migrate(t, member.ref, "1000", "Eth_TxHash_test", migrationAddress, 2)
	migrate(t, member.ref, "1000", "Eth_TxHash_test", migrationAddress, 0)

	_, err = signedRequest(member, "deposit.transfer", map[string]interface{}{"amount": "100", "ethTxHash": "Eth_TxHash_test"})
	require.Contains(t, err.Error(), "number of confirms is less then 3")
}
