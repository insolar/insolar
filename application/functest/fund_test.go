// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

// +build functest

package functest

import (
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/application"
	"github.com/insolar/insolar/application/genesisrefs"
	"github.com/insolar/insolar/applicationbase/testutils/launchnet"
	"github.com/insolar/insolar/applicationbase/testutils/testrequest"
)

func TestFoundationMemberCreate(t *testing.T) {
	for _, m := range Foundation {
		err := verifyFundsMembersAndDeposits(t, m, application.FoundationDistributionAmount)
		if err != nil {
			require.NoError(t, err)
		}
	}
}

func TestEnterpriseMemberCreate(t *testing.T) {
	for _, m := range Enterprise {
		err := verifyFundsMembersExist(t, m, application.EnterpriseDistributionAmount)
		if err != nil {
			require.NoError(t, err)
		}
	}
}

func TestNetworkIncentivesMemberCreate(t *testing.T) {
	// for speed up test check only last member
	m := NetworkIncentives[application.GenesisAmountNetworkIncentivesMembers-1]

	err := verifyFundsMembersAndDeposits(t, m, application.NetworkIncentivesDistributionAmount)
	if err != nil {
		require.NoError(t, err)
	}
}

func TestApplicationIncentivesMemberCreate(t *testing.T) {
	for _, m := range ApplicationIncentives {
		err := verifyFundsMembersAndDeposits(t, m, application.AppIncentivesDistributionAmount)
		if err != nil {
			require.NoError(t, err)
		}
	}
}

func checkBalanceAndDepositFewTimes(t *testing.T, m *AppUser, expectedBalance string, expectedDeposit string) {
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

func TestNetworkIncentivesTransferDeposit(t *testing.T) {
	// for speed up test check only last member
	m := NetworkIncentives[application.GenesisAmountNetworkIncentivesMembers-1]

	res2, err := testrequest.SignedRequest(t, launchnet.TestRPCUrlPublic, m, "member.get", nil)
	require.NoError(t, err)
	decodedRes2, ok := res2.(map[string]interface{})
	m.Ref = decodedRes2["reference"].(string)
	require.True(t, ok, fmt.Sprintf("failed to decode: expected map[string]interface{}, got %T", res2))

	_, err = testrequest.SignedRequestWithEmptyRequestRef(t, launchnet.TestRPCUrlPublic, m,
		"deposit.transfer", map[string]interface{}{"amount": "100", "ethTxHash": genesisrefs.FundsDepositName},
	)
	data := checkConvertRequesterError(t, err).Data
	require.Contains(t, data.Trace, "hold period didn't end")

	checkBalanceAndDepositFewTimes(t, m, "0", application.NetworkIncentivesDistributionAmount)
}

func TestApplicationIncentivesTransferDeposit(t *testing.T) {
	for _, m := range ApplicationIncentives {
		res2, err := testrequest.SignedRequest(t, launchnet.TestRPCUrlPublic, m, "member.get", nil)
		require.NoError(t, err)
		decodedRes2, ok := res2.(map[string]interface{})
		m.Ref = decodedRes2["reference"].(string)
		require.True(t, ok, fmt.Sprintf("failed to decode: expected map[string]interface{}, got %T", res2))

		_, err = testrequest.SignedRequestWithEmptyRequestRef(t, launchnet.TestRPCUrlPublic, m,
			"deposit.transfer", map[string]interface{}{"amount": "100", "ethTxHash": genesisrefs.FundsDepositName},
		)
		data := checkConvertRequesterError(t, err).Data
		require.Contains(t, data.Trace, "hold period didn't end")

		checkBalanceAndDepositFewTimes(t, m, "0", application.AppIncentivesDistributionAmount)
	}
}

func TestFoundationTransferDeposit(t *testing.T) {
	for _, m := range Foundation {
		res2, err := testrequest.SignedRequest(t, launchnet.TestRPCUrlPublic, m, "member.get", nil)
		require.NoError(t, err)
		decodedRes2, ok := res2.(map[string]interface{})
		m.Ref = decodedRes2["reference"].(string)
		require.True(t, ok, fmt.Sprintf("failed to decode: expected map[string]interface{}, got %T", res2))

		_, err = testrequest.SignedRequestWithEmptyRequestRef(t, launchnet.TestRPCUrlPublic, m,
			"deposit.transfer", map[string]interface{}{"amount": "100", "ethTxHash": genesisrefs.FundsDepositName},
		)
		data := checkConvertRequesterError(t, err).Data
		require.Contains(t, data.Trace, "hold period didn't end")

		checkBalanceAndDepositFewTimes(t, m, "0", application.FoundationDistributionAmount)
	}
}

func TestMigrationDaemonTransferDeposit(t *testing.T) {
	m := &MigrationAdmin

	res2, err := testrequest.SignedRequest(t, launchnet.TestRPCUrlPublic, m, "member.get", nil)
	require.NoError(t, err)
	decodedRes2, ok := res2.(map[string]interface{})
	require.True(t, ok, fmt.Sprintf("failed to decode: expected map[string]interface{}, got %T", res2))
	m.Ref = decodedRes2["reference"].(string)

	oldBalance, deposits := getBalanceAndDepositsNoErr(t, m, m.Ref)
	oldDepositStr := deposits[genesisrefs.FundsDepositName].(map[string]interface{})["balance"].(string)

	_, err = testrequest.SignedRequestWithEmptyRequestRef(t, launchnet.TestRPCUrlPublic, m,
		"deposit.transfer", map[string]interface{}{"amount": "100", "ethTxHash": genesisrefs.FundsDepositName},
	)
	data := checkConvertRequesterError(t, err).Data
	require.Contains(t, data.Trace, "hold period didn't end")

	newBalance, newDeposits := getBalanceAndDepositsNoErr(t, m, m.Ref)
	newDepositStr := newDeposits[genesisrefs.FundsDepositName].(map[string]interface{})["balance"].(string)
	require.Equal(t, oldBalance.String(), newBalance.String())
	require.Equal(t, oldDepositStr, newDepositStr)
}
