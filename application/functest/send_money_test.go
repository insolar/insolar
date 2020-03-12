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

	"github.com/insolar/insolar/applicationbase/testutils/launchnet"
	"github.com/insolar/insolar/applicationbase/testutils/testrequest"
	"github.com/insolar/insolar/insolar/gen"
)

const times = 5
const feeSize = "100000000"

func checkBalanceFewTimes(t *testing.T, caller *AppUser, ref string, expected *big.Int) {
	for i := 0; i < times; i++ {
		balance := getBalanceNoErr(t, caller, ref)
		if balance.String() == expected.String() {
			return
		}
		time.Sleep(time.Second)
	}
	t.Error("Received balance is not equal expected")
}

// TODO: uncomment after undoing of all transaction in failed request will be supported
func TestTransferMoney(t *testing.T) {
	firstMember := createMember(t)
	secondMember := createMember(t)

	// init money on members
	_, err := testrequest.SignedRequest(t, launchnet.TestRPCUrlPublic, &Root, "member.transfer",
		map[string]interface{}{"amount": "2000000000", "toMemberReference": firstMember.Ref})
	require.NoError(t, err)
	_, err = testrequest.SignedRequest(t, launchnet.TestRPCUrlPublic, &Root, "member.transfer",
		map[string]interface{}{"amount": "2000000000", "toMemberReference": secondMember.Ref})
	require.NoError(t, err)

	oldFirstBalance := getBalanceNoErr(t, firstMember, firstMember.Ref)
	oldSecondBalance := getBalanceNoErr(t, secondMember, secondMember.Ref)

	feeRes, err := testrequest.SignedRequest(t, launchnet.TestRPCUrlPublic, &FeeMember, "member.get", nil)
	require.Nil(t, err)
	feeMemberRef, ok := feeRes.(map[string]interface{})["reference"].(string)
	require.True(t, ok)
	FeeMember.Ref = feeMemberRef
	feeBalance := getBalanceNoErr(t, &FeeMember, feeMemberRef)

	amountStr := "10"
	amount, _ := new(big.Int).SetString(amountStr, 10)
	fee, _ := new(big.Int).SetString(feeSize, 10)
	expectedFirstBalance := new(big.Int).Sub(oldFirstBalance, amount)
	expectedFirstBalance.Sub(expectedFirstBalance, fee)
	expectedSecondBalance := new(big.Int).Add(oldSecondBalance, amount)
	expectedFeeBalance := new(big.Int).Add(feeBalance, fee)

	_, err = testrequest.SignedRequest(t, launchnet.TestRPCUrlPublic, firstMember, "member.transfer",
		map[string]interface{}{"amount": amountStr, "toMemberReference": secondMember.Ref})
	require.NoError(t, err)

	checkBalanceFewTimes(t, secondMember, secondMember.Ref, expectedSecondBalance)
	newFirstBalance := getBalanceNoErr(t, firstMember, firstMember.Ref)
	require.Equal(t, expectedFirstBalance, newFirstBalance)

	checkBalanceFewTimes(t, &FeeMember, feeMemberRef, expectedFeeBalance)
	newFeeBalance := getBalanceNoErr(t, &FeeMember, feeMemberRef)
	require.Equal(t, expectedFeeBalance, newFeeBalance)
}

func TestTransferMoneyToNotObjectRef(t *testing.T) {
	firstMember := createMember(t)
	secondMember := createMember(t)
	oldFirstBalance := getBalanceNoErr(t, firstMember, firstMember.Ref)
	oldSecondBalance := getBalanceNoErr(t, secondMember, secondMember.Ref)

	feeRes, err := testrequest.SignedRequest(t, launchnet.TestRPCUrlPublic, &FeeMember, "member.get", nil)
	require.Nil(t, err)
	feeMemberRef, ok := feeRes.(map[string]interface{})["reference"].(string)
	require.True(t, ok)
	FeeMember.Ref = feeMemberRef
	feeBalance := getBalanceNoErr(t, &FeeMember, feeMemberRef)

	amountStr := "10"

	_, _, err = testrequest.MakeSignedRequest(launchnet.TestRPCUrlPublic, firstMember, "member.transfer",
		map[string]interface{}{"amount": amountStr, "toMemberReference": secondMember.Ref + ".record"})
	require.Error(t, err)

	data := checkConvertRequesterError(t, err).Data
	require.Contains(t, strings.Join(data.Trace, ": "), "OpenAPI schema validation")

	newFirstBalance := getBalanceNoErr(t, firstMember, firstMember.Ref)
	newSecondBalance := getBalanceNoErr(t, secondMember, secondMember.Ref)
	newFeeBalance := getBalanceNoErr(t, &FeeMember, feeMemberRef)
	require.Equal(t, oldFirstBalance, newFirstBalance)
	require.Equal(t, oldSecondBalance, newSecondBalance)
	require.Equal(t, feeBalance, newFeeBalance)
}

func TestTransferMoneyToNotSelfScopedRef(t *testing.T) {
	firstMember := createMember(t)
	secondMember := createMember(t)
	oldFirstBalance := getBalanceNoErr(t, firstMember, firstMember.Ref)
	oldSecondBalance := getBalanceNoErr(t, secondMember, secondMember.Ref)

	feeRes, err := testrequest.SignedRequest(t, launchnet.TestRPCUrlPublic, &FeeMember, "member.get", nil)
	require.Nil(t, err)
	feeMemberRef, ok := feeRes.(map[string]interface{})["reference"].(string)
	require.True(t, ok)
	FeeMember.Ref = feeMemberRef
	feeBalance := getBalanceNoErr(t, &FeeMember, feeMemberRef)

	amountStr := "10"

	sepIdx := strings.Index(firstMember.Ref, ":")
	_, _, err = testrequest.MakeSignedRequest(
		launchnet.TestRPCUrlPublic,
		firstMember,
		"member.transfer",
		map[string]interface{}{
			"amount":            amountStr,
			"toMemberReference": secondMember.Ref + "." + firstMember.Ref[sepIdx+1:],
		},
	)
	require.Error(t, err)

	data := checkConvertRequesterError(t, err).Data
	require.Contains(t, strings.Join(data.Trace, ": "), "OpenAPI schema validation")

	newFirstBalance := getBalanceNoErr(t, firstMember, firstMember.Ref)
	newSecondBalance := getBalanceNoErr(t, secondMember, secondMember.Ref)
	newFeeBalance := getBalanceNoErr(t, &FeeMember, feeMemberRef)
	require.Equal(t, oldFirstBalance, newFirstBalance)
	require.Equal(t, oldSecondBalance, newSecondBalance)
	require.Equal(t, feeBalance, newFeeBalance)
}

func TestTransferMoneyFromNotExist(t *testing.T) {
	firstMember := createMember(t)
	firstMember.Ref = gen.Reference().String()

	secondMember := createMember(t)
	oldSecondBalance := getBalanceNoErr(t, secondMember, secondMember.Ref)

	amount := "10"

	_, err := testrequest.SignedRequestWithEmptyRequestRef(t, launchnet.TestRPCUrlPublic, firstMember, "member.transfer",
		map[string]interface{}{"amount": amount, "toMemberReference": secondMember.Ref})
	data := checkConvertRequesterError(t, err).Data
	require.Contains(t, data.Trace, "failed to fetch index from heavy")
	newSecondBalance := getBalanceNoErr(t, secondMember, secondMember.Ref)
	require.Equal(t, oldSecondBalance, newSecondBalance)
}

func TestTransferMoneyToNotExist(t *testing.T) {
	firstMember := createMember(t)
	oldFirstBalance := getBalanceNoErr(t, firstMember, firstMember.Ref)

	amount := "10"

	_, err := testrequest.SignedRequestWithEmptyRequestRef(t, launchnet.TestRPCUrlPublic, firstMember, "member.transfer",
		map[string]interface{}{"amount": amount, "toMemberReference": gen.Reference().String()})
	data := checkConvertRequesterError(t, err).Data
	require.Contains(t, data.Trace, "recipient member does not exist")

	newFirstBalance := getBalanceNoErr(t, firstMember, firstMember.Ref)
	require.Equal(t, oldFirstBalance, newFirstBalance)
}

func TestTransferNegativeAmount(t *testing.T) {
	firstMember := createMember(t)
	secondMember := createMember(t)
	oldFirstBalance := getBalanceNoErr(t, firstMember, firstMember.Ref)
	oldSecondBalance := getBalanceNoErr(t, secondMember, secondMember.Ref)

	amount := "-111"

	_, err := testrequest.SignedRequestWithEmptyRequestRef(t, launchnet.TestRPCUrlPublic, firstMember, "member.transfer",
		map[string]interface{}{"amount": amount, "toMemberReference": secondMember.Ref})
	require.Error(t, err)

	newFirstBalance := getBalanceNoErr(t, firstMember, firstMember.Ref)
	newSecondBalance := getBalanceNoErr(t, secondMember, secondMember.Ref)
	require.Equal(t, oldFirstBalance, newFirstBalance)
	require.Equal(t, oldSecondBalance, newSecondBalance)
}

// TODO: unskip test after undoing of all transaction in failed request will be supported
func TestTransferAllAmount(t *testing.T) {
	t.Skip()
	firstMember := createMember(t)
	secondMember := createMember(t)
	oldFirstBalance := getBalanceNoErr(t, firstMember, firstMember.Ref)
	oldSecondBalance := getBalanceNoErr(t, secondMember, secondMember.Ref)

	amount := oldFirstBalance

	summ := new(big.Int)
	summ.Add(oldSecondBalance, oldFirstBalance)
	_, err := testrequest.SignedRequest(t, launchnet.TestRPCUrlPublic, firstMember, "member.transfer",
		map[string]interface{}{"amount": amount, "toMemberReference": secondMember.Ref})
	require.NoError(t, err)

	checkBalanceFewTimes(t, secondMember, secondMember.Ref, summ)
	newFirstBalance := getBalanceNoErr(t, firstMember, firstMember.Ref)
	require.Equal(t, 0, newFirstBalance)
}

func TestTransferMoreThanAvailableAmount(t *testing.T) {
	firstMember := createMember(t)
	secondMember := createMember(t)
	oldFirstBalance := getBalanceNoErr(t, firstMember, firstMember.Ref)
	oldSecondBalance := getBalanceNoErr(t, secondMember, secondMember.Ref)

	amount := new(big.Int)
	amount.Add(oldFirstBalance, big.NewInt(10))

	_, err := testrequest.SignedRequestWithEmptyRequestRef(t, launchnet.TestRPCUrlPublic, firstMember, "member.transfer",
		map[string]interface{}{"amount": amount.String(), "toMemberReference": secondMember.Ref})
	data := checkConvertRequesterError(t, err).Data
	require.Contains(t, data.Trace, "balance is too low")
	newFirstBalance := getBalanceNoErr(t, firstMember, firstMember.Ref)
	newSecondBalance := getBalanceNoErr(t, secondMember, secondMember.Ref)
	require.Equal(t, oldFirstBalance, newFirstBalance)
	require.Equal(t, oldSecondBalance, newSecondBalance)
}

func TestTransferToMyself(t *testing.T) {
	member := createMember(t)
	oldMemberBalance := getBalanceNoErr(t, member, member.Ref)

	amount := "20"

	_, err := testrequest.SignedRequestWithEmptyRequestRef(t, launchnet.TestRPCUrlPublic, member, "member.transfer",
		map[string]interface{}{"amount": amount, "toMemberReference": member.Ref})
	data := checkConvertRequesterError(t, err).Data
	require.Contains(t, data.Trace, "recipient must be different from the sender")
	newMemberBalance := getBalanceNoErr(t, member, member.Ref)
	require.Equal(t, oldMemberBalance, newMemberBalance)
}

// TODO: test to check overflow of balance
// TODO: check transfer zero amount

// TODO: uncomment after undoing of all transaction in failed request will be supported
func TestTransferTwoTimes(t *testing.T) {
	firstMember := createMember(t)
	secondMember := createMember(t)
	// Skip validation of balance before/after transfer
	// oldFirstBalance := getBalanceNoErr(t, firstMember, firstMember.Ref)
	// oldSecondBalance := getBalanceNoErr(t, secondMember, secondMember.Ref)

	// init money on members
	_, err := testrequest.SignedRequest(t, launchnet.TestRPCUrlPublic, &Root, "member.transfer",
		map[string]interface{}{"amount": "5000000000", "toMemberReference": firstMember.Ref})
	require.NoError(t, err)
	_, err = testrequest.SignedRequest(t, launchnet.TestRPCUrlPublic, &Root, "member.transfer",
		map[string]interface{}{"amount": "5000000000", "toMemberReference": secondMember.Ref})
	require.NoError(t, err)

	amount := "10"

	_, err = testrequest.SignedRequest(t, launchnet.TestRPCUrlPublic, firstMember, "member.transfer",
		map[string]interface{}{"amount": amount, "toMemberReference": secondMember.Ref})
	require.NoError(t, err)
	_, err = testrequest.SignedRequest(t, launchnet.TestRPCUrlPublic, firstMember, "member.transfer",
		map[string]interface{}{"amount": amount, "toMemberReference": secondMember.Ref})
	require.NoError(t, err)

	// Skip validation of balance before/after transfer
	// checkBalanceFewTimes(t, secondMember, secondMember.Ref, oldSecondBalance+2*amount)
	// newFirstBalance := getBalanceNoErr(t, firstMember, firstMember.Ref)
	// require.Equal(t, oldFirstBalance-2*amount, newFirstBalance)
}
