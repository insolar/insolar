// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

// +build functest

package functest

import (
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/api"
	"github.com/insolar/insolar/application/testutils/launchnet"
	"github.com/insolar/insolar/logicrunner/builtin/foundation"
	"github.com/insolar/insolar/testutils"
)

// TODO: https://insolar.atlassian.net/browse/WLT-768
func TestDepositTransferToken(t *testing.T) {
	ethHash := testutils.RandomEthHash()

	member := fullMigration(t, ethHash)

	firstBalance := getBalanceNoErr(t, member, member.Ref)
	secondBalance := new(big.Int).Add(firstBalance, big.NewInt(1000))

	anon := func() api.CallMethodReply {
		_, _, err := makeSignedRequest(launchnet.TestRPCUrlPublic, member,
			"deposit.transfer", map[string]interface{}{"amount": "1000", "ethTxHash": ethHash})

		data := checkConvertRequesterError(t, err).Data
		for _, v := range data.Trace {
			if !strings.Contains(v, "hold period didn't end") {
				return api.CallMethodReply{}
			}
		}
		return api.CallMethodReply{
			Error: &foundation.Error{S: err.Error()},
		}
	}

	_, err := waitUntilRequestProcessed(anon, time.Second*30, time.Second, 30)
	require.NoError(t, err)
	anon = func() api.CallMethodReply {
		_, _, err := makeSignedRequest(launchnet.TestRPCUrlPublic, member,
			"deposit.transfer", map[string]interface{}{"amount": "1000", "ethTxHash": ethHash})
		if err == nil {
			return api.CallMethodReply{}
		}
		return api.CallMethodReply{
			Error: &foundation.Error{S: err.Error()},
		}
	}
	_, err = waitUntilRequestProcessed(anon, time.Second*30, time.Second, 30)
	require.NoError(t, err)
	checkBalanceFewTimes(t, member, member.Ref, secondBalance)
}

func TestDepositTransferBeforeUnhold(t *testing.T) {
	ethHash := testutils.RandomEthHash()

	member := fullMigration(t, ethHash)

	_, err := signedRequestWithEmptyRequestRef(t, launchnet.TestRPCUrlPublic, member,
		"deposit.transfer", map[string]interface{}{"amount": "100", "ethTxHash": ethHash})
	data := checkConvertRequesterError(t, err).Data
	require.Contains(t, data.Trace, "hold period didn't end", "check lockup_pulse_period param at bootstrap.yaml: it maybe too low")
}

func TestDepositTransferBiggerAmount(t *testing.T) {
	ethHash := testutils.RandomEthHash()

	member := fullMigration(t, ethHash)

	_, err := signedRequestWithEmptyRequestRef(t, launchnet.TestRPCUrlPublic, member,
		"deposit.transfer", map[string]interface{}{"amount": "10000000000000", "ethTxHash": ethHash})
	data := checkConvertRequesterError(t, err).Data
	require.Contains(t, data.Trace, "not enough balance for transfer")
}

func TestDepositTransferAnotherTx(t *testing.T) {
	ethHash := testutils.RandomEthHash()

	member := fullMigration(t, ethHash)

	_, err := signedRequestWithEmptyRequestRef(t, launchnet.TestRPCUrlPublic, member,
		"deposit.transfer", map[string]interface{}{"amount": "100", "ethTxHash": testutils.RandomEthHash()})
	data := checkConvertRequesterError(t, err).Data
	require.Contains(t, data.Trace, "can't find deposit")
}

func TestDepositTransferWrongValueAmount(t *testing.T) {
	ethHash := testutils.RandomEthHash()

	member := fullMigration(t, ethHash)

	_, err := signedRequestWithEmptyRequestRef(t, launchnet.TestRPCUrlPublic, member,
		"deposit.transfer", map[string]interface{}{"amount": "foo", "ethTxHash": ethHash})
	require.Error(t, err)
	data := checkConvertRequesterError(t, err).Data
	expectedError(t, data.Trace, `Error at "/params/callParams/amount":JSON string doesn't match the regular expression '^[1-9][0-9]*$`)
}

func TestDepositTransferNotEnoughConfirms(t *testing.T) {
	ethHash := testutils.RandomEthHash()

	activateDaemons(t, countTwoActiveDaemon)
	member := createMigrationMemberForMA(t)
	_ = migrate(t, member.Ref, "1000", ethHash, member.MigrationAddress, 2)

	_, err := signedRequestWithEmptyRequestRef(t, launchnet.TestRPCUrlPublic, member,
		"deposit.transfer", map[string]interface{}{"amount": "100", "ethTxHash": ethHash})
	data := checkConvertRequesterError(t, err).Data
	require.Contains(t, data.Trace, "not enough balance for transfer")
}
