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
	"github.com/insolar/insolar/api"
	"github.com/insolar/insolar/logicrunner/builtin/foundation"
	"github.com/insolar/insolar/testutils"
	"github.com/stretchr/testify/require"
	"math/big"
	"strings"
	"testing"
	"time"
)

// TODO: https://insolar.atlassian.net/browse/WLT-768
func TestDepositTransferToken(t *testing.T) {
	member := fullMigration(t, "Eth_TxHash_test")

	firstBalance := getBalanceNoErr(t, member, member.Ref)
	secondBalance := new(big.Int).Add(firstBalance, big.NewInt(1000))

	anon := func() api.CallMethodReply {
		_, _, err := makeSignedRequest(member, "deposit.transfer", map[string]interface{}{"amount": "1000", "ethTxHash": "Eth_TxHash_test"})
		require.Error(t, err)
		fmt.Println("Lockup STEP: --- RETURNED   " + err.Error())
		if !strings.Contains(err.Error(), "hold period didn't end") {
			return api.CallMethodReply{}
		}
		return api.CallMethodReply{
			Error: &foundation.Error{err.Error()},
		}
	}
	_, err := waitUntilRequestProcessed(anon, time.Second*20, time.Second, 20)
	require.NoError(t, err)

	anon = func() api.CallMethodReply {
		_, _, err := makeSignedRequest(member, "deposit.transfer", map[string]interface{}{"amount": "1000", "ethTxHash": "Eth_TxHash_test"})
		if err == nil {
			fmt.Println("Vesting STEP: === SUCCESS")
			return api.CallMethodReply{}
		}
		fmt.Println("Vesting STEP: --- RETURNED   " + err.Error())
		return api.CallMethodReply{
			Error: &foundation.Error{err.Error()},
		}
	}
	_, err = waitUntilRequestProcessed(anon, time.Second*20, time.Second, 20)
	require.NoError(t, err)
	checkBalanceFewTimes(t, member, member.Ref, secondBalance)
}

func TestDepositTransferBeforeUnhold(t *testing.T) {
	member := fullMigration(t, "Eth_TxHash_test")

	_, err := signedRequestWithEmptyRequestRef(t, member, "deposit.transfer", map[string]interface{}{"amount": "100", "ethTxHash": "Eth_TxHash_test"})
	require.Error(t, err)
	require.Contains(t, err.Error(), "hold period didn't end")
}

func TestDepositTransferBiggerAmount(t *testing.T) {
	member := fullMigration(t, "Eth_TxHash_test")

	_, err := signedRequestWithEmptyRequestRef(t, member, "deposit.transfer", map[string]interface{}{"amount": "10000000000000", "ethTxHash": "Eth_TxHash_test"})
	require.Error(t, err)
	require.Contains(t, err.Error(), "not enough balance for transfer")
}

func TestDepositTransferAnotherTx(t *testing.T) {
	member := fullMigration(t, "Eth_TxHash_test")

	_, err := signedRequestWithEmptyRequestRef(t, member, "deposit.transfer", map[string]interface{}{"amount": "100", "ethTxHash": "Eth_TxHash_testNovalid"})
	require.Error(t, err)
	require.Contains(t, err.Error(), "can't find deposit")
}

func TestDepositTransferWrongValueAmount(t *testing.T) {
	member := fullMigration(t, "Eth_TxHash_test")

	_, err := signedRequestWithEmptyRequestRef(t, member, "deposit.transfer", map[string]interface{}{"amount": "foo", "ethTxHash": "Eth_TxHash_test"})
	require.Error(t, err)
	require.Contains(t, err.Error(), "can't parse input amount")
}

func TestDepositTransferNotEnoughConfirms(t *testing.T) {
	activateDaemons(t)
	migrationAddress := testutils.RandomString()
	member := createMigrationMemberForMA(t, migrationAddress)
	_ = migrate(t, member.Ref, "1000", "Eth_TxHash_test", migrationAddress, 2)

	_ = migrate(t, member.Ref, "1000", "Eth_TxHash_test", migrationAddress, 0)

	_, err := signedRequestWithEmptyRequestRef(t, member, "deposit.transfer", map[string]interface{}{"amount": "100", "ethTxHash": "Eth_TxHash_test"})
	require.Error(t, err)
	require.Contains(t, err.Error(), "not enough balance for transfer")
}
