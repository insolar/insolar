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
	"testing"
	"time"

	"github.com/insolar/insolar/application/genesisrefs"
	"github.com/insolar/insolar/application/testutils/launchnet"
	"github.com/stretchr/testify/require"
)

func TestFoundationMemberCreate(t *testing.T) {
	for _, m := range launchnet.Foundation {
		err := verifyFundsMembersAndDeposits(t, m)
		if err != nil {
			require.NoError(t, err)
		}
	}
}

func TestEnterpriseMemberCreate(t *testing.T) {
	for _, m := range launchnet.Enterprise {
		err := verifyFundsMembersExist(t, m)
		if err != nil {
			require.NoError(t, err)
		}
	}
}

func TestNetworkIncentivesMemberCreate(t *testing.T) {
	for _, m := range launchnet.NetworkIncentives {
		err := verifyFundsMembersAndDeposits(t, m)
		if err != nil {
			require.NoError(t, err)
		}
	}
}

func TestApplicationIncentivesMemberCreate(t *testing.T) {
	for _, m := range launchnet.ApplicationIncentives {
		err := verifyFundsMembersAndDeposits(t, m)
		if err != nil {
			require.NoError(t, err)
		}
	}
}

func TestFundsMemberCreate(t *testing.T) {
	for _, m := range launchnet.Funds {
		err := verifyFundsMembersExist(t, m)
		if err != nil {
			require.NoError(t, err)
		}
	}
}

func checkBalanceAndDepositFewTimes(t *testing.T, m *launchnet.User, expectedBalance string, expectedDeposit string) {
	var balance *big.Int
	var depositStr string
	for i := 0; i < times; i++ {
		balance, deposits := getBalanceAndDepositsNoErr(t, m, m.Ref)
		depositStr = deposits[genesisrefs.FundsDepositName].(map[string]interface{})["balance"].(string)
		if balance.String() == expectedBalance && depositStr == expectedDeposit {
			return
		}
		time.Sleep(time.Second)
	}
	t.Errorf("Received balance or deposite is not equal expected: current balance %s, expected %s;"+
		" current deposite %s, expected %s",
		balance, expectedBalance,
		depositStr, expectedDeposit)
}

func TestFundsTransferDeposit(t *testing.T) {
	for _, m := range launchnet.Funds {
		res2, err := signedRequest(t, launchnet.TestRPCUrlPublic, m, "member.get", nil)
		require.NoError(t, err)
		decodedRes2, ok := res2.(map[string]interface{})
		m.Ref = decodedRes2["reference"].(string)
		require.True(t, ok, fmt.Sprintf("failed to decode: expected map[string]interface{}, got %T", res2))

		oldBalance, deposits := getBalanceAndDepositsNoErr(t, m, m.Ref)
		oldDepositStr := deposits[genesisrefs.FundsDepositName].(map[string]interface{})["balance"].(string)
		amount, _ := new(big.Int).SetString("100", 10)
		expectedBalance := new(big.Int).Add(oldBalance, amount)
		oldDeposit, _ := new(big.Int).SetString(oldDepositStr, 10)
		expectedDeposit := new(big.Int).Sub(oldDeposit, amount)

		_, err = signedRequest(t, launchnet.TestRPCUrlPublic, m,
			"deposit.transfer", map[string]interface{}{"amount": "100", "ethTxHash": genesisrefs.FundsDepositName},
		)
		require.NoError(t, err)

		checkBalanceAndDepositFewTimes(t, m, expectedBalance.String(), expectedDeposit.String())
	}
}

func TestEnterpriseTransferDeposit(t *testing.T) {
	for _, m := range launchnet.Enterprise {
		res2, err := signedRequest(t, launchnet.TestRPCUrlPublic, m, "member.get", nil)
		require.NoError(t, err)
		decodedRes2, ok := res2.(map[string]interface{})
		m.Ref = decodedRes2["reference"].(string)
		require.True(t, ok, fmt.Sprintf("failed to decode: expected map[string]interface{}, got %T", res2))

		oldBalance, deposits := getBalanceAndDepositsNoErr(t, m, m.Ref)
		oldDepositStr := deposits[genesisrefs.FundsDepositName].(map[string]interface{})["balance"].(string)
		amount, _ := new(big.Int).SetString("100", 10)
		expectedBalance := new(big.Int).Add(oldBalance, amount)
		oldDeposit, _ := new(big.Int).SetString(oldDepositStr, 10)
		expectedDeposit := new(big.Int).Sub(oldDeposit, amount)

		_, err = signedRequest(t, launchnet.TestRPCUrlPublic, m,
			"deposit.transfer", map[string]interface{}{"amount": "100", "ethTxHash": genesisrefs.FundsDepositName},
		)
		require.NoError(t, err)

		checkBalanceAndDepositFewTimes(t, m, expectedBalance.String(), expectedDeposit.String())
	}
}

func TestNetworkIncentivesTransferDeposit(t *testing.T) {
	for _, m := range launchnet.NetworkIncentives {
		res2, err := signedRequest(t, launchnet.TestRPCUrlPublic, m, "member.get", nil)
		require.NoError(t, err)
		decodedRes2, ok := res2.(map[string]interface{})
		m.Ref = decodedRes2["reference"].(string)
		require.True(t, ok, fmt.Sprintf("failed to decode: expected map[string]interface{}, got %T", res2))

		_, err = signedRequestWithEmptyRequestRef(t, launchnet.TestRPCUrlPublic, m,
			"deposit.transfer", map[string]interface{}{"amount": "100", "ethTxHash": genesisrefs.FundsDepositName},
		)
		data := checkConvertRequesterError(t, err).Data
		require.Contains(t, data.Trace, "hold period didn't end")

		checkBalanceAndDepositFewTimes(t, m, "0", TestDepositAmount)
	}
}

func TestApplicationIncentivesTransferDeposit(t *testing.T) {
	for _, m := range launchnet.ApplicationIncentives {
		res2, err := signedRequest(t, launchnet.TestRPCUrlPublic, m, "member.get", nil)
		require.NoError(t, err)
		decodedRes2, ok := res2.(map[string]interface{})
		m.Ref = decodedRes2["reference"].(string)
		require.True(t, ok, fmt.Sprintf("failed to decode: expected map[string]interface{}, got %T", res2))

		_, err = signedRequestWithEmptyRequestRef(t, launchnet.TestRPCUrlPublic, m,
			"deposit.transfer", map[string]interface{}{"amount": "100", "ethTxHash": genesisrefs.FundsDepositName},
		)
		data := checkConvertRequesterError(t, err).Data
		require.Contains(t, data.Trace, "hold period didn't end")

		checkBalanceAndDepositFewTimes(t, m, "0", TestDepositAmount)
	}
}

func TestFoundationTransferDeposit(t *testing.T) {
	for _, m := range launchnet.Foundation {
		res2, err := signedRequest(t, launchnet.TestRPCUrlPublic, m, "member.get", nil)
		require.NoError(t, err)
		decodedRes2, ok := res2.(map[string]interface{})
		m.Ref = decodedRes2["reference"].(string)
		require.True(t, ok, fmt.Sprintf("failed to decode: expected map[string]interface{}, got %T", res2))

		_, err = signedRequestWithEmptyRequestRef(t, launchnet.TestRPCUrlPublic, m,
			"deposit.transfer", map[string]interface{}{"amount": "100", "ethTxHash": genesisrefs.FundsDepositName},
		)
		data := checkConvertRequesterError(t, err).Data
		require.Contains(t, data.Trace, "hold period didn't end")

		checkBalanceAndDepositFewTimes(t, m, "0", TestDepositAmount)
	}
}

func TestMigrationDaemonTransferDeposit(t *testing.T) {
	m := &launchnet.MigrationAdmin

	res2, err := signedRequest(t, launchnet.TestRPCUrlPublic, m, "member.get", nil)
	require.NoError(t, err)
	decodedRes2, ok := res2.(map[string]interface{})
	require.True(t, ok, fmt.Sprintf("failed to decode: expected map[string]interface{}, got %T", res2))
	m.Ref = decodedRes2["reference"].(string)

	oldBalance, deposits := getBalanceAndDepositsNoErr(t, m, m.Ref)
	oldDepositStr := deposits[genesisrefs.FundsDepositName].(map[string]interface{})["balance"].(string)

	_, err = signedRequestWithEmptyRequestRef(t, launchnet.TestRPCUrlPublic, m,
		"deposit.transfer", map[string]interface{}{"amount": "100", "ethTxHash": genesisrefs.FundsDepositName},
	)
	data := checkConvertRequesterError(t, err).Data
	require.Contains(t, data.Trace, "hold period didn't end")

	newBalance, newDeposits := getBalanceAndDepositsNoErr(t, m, m.Ref)
	newDepositStr := newDeposits[genesisrefs.FundsDepositName].(map[string]interface{})["balance"].(string)
	require.Equal(t, oldBalance.String(), newBalance.String())
	require.Equal(t, oldDepositStr, newDepositStr)
}
