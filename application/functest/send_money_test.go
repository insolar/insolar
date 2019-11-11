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
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/application/testutils/launchnet"
	"github.com/insolar/insolar/insolar/gen"
)

const times = 5
const feeSize = "1000000000"

func checkBalanceFewTimes(t *testing.T, caller *launchnet.User, ref string, expected *big.Int) {
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
	oldFirstBalance := getBalanceNoErr(t, firstMember, firstMember.Ref)
	oldSecondBalance := getBalanceNoErr(t, secondMember, secondMember.Ref)

	feeRes, err := signedRequest(t, launchnet.TestRPCUrlPublic, &launchnet.FeeMember, "member.get", nil)
	require.Nil(t, err)
	feeMemberRef, ok := feeRes.(map[string]interface{})["reference"].(string)
	require.True(t, ok)
	launchnet.FeeMember.Ref = feeMemberRef
	feeBalance := getBalanceNoErr(t, &launchnet.FeeMember, feeMemberRef)

	amountStr := "10"
	amount, _ := new(big.Int).SetString(amountStr, 10)
	fee, _ := new(big.Int).SetString(feeSize, 10)
	expectedFirstBalance := new(big.Int).Sub(oldFirstBalance, amount)
	expectedFirstBalance.Sub(expectedFirstBalance, fee)
	expectedSecondBalance := new(big.Int).Add(oldSecondBalance, amount)
	expectedFeeBalance := new(big.Int).Add(feeBalance, fee)

	_, err = signedRequest(t, launchnet.TestRPCUrlPublic, firstMember, "member.transfer",
		map[string]interface{}{"amount": amountStr, "toMemberReference": secondMember.Ref})
	require.NoError(t, err)

	checkBalanceFewTimes(t, secondMember, secondMember.Ref, expectedSecondBalance)
	newFirstBalance := getBalanceNoErr(t, firstMember, firstMember.Ref)
	require.Equal(t, expectedFirstBalance, newFirstBalance)

	checkBalanceFewTimes(t, &launchnet.FeeMember, feeMemberRef, expectedFeeBalance)
	newFeeBalance := getBalanceNoErr(t, &launchnet.FeeMember, feeMemberRef)
	require.Equal(t, expectedFeeBalance, newFeeBalance)
}

func TestTransferMoneyToNotObjectRef(t *testing.T) {
	firstMember := createMember(t)
	secondMember := createMember(t)
	oldFirstBalance := getBalanceNoErr(t, firstMember, firstMember.Ref)
	oldSecondBalance := getBalanceNoErr(t, secondMember, secondMember.Ref)

	feeRes, err := signedRequest(t, launchnet.TestRPCUrlPublic, &launchnet.FeeMember, "member.get", nil)
	require.Nil(t, err)
	feeMemberRef, ok := feeRes.(map[string]interface{})["reference"].(string)
	require.True(t, ok)
	launchnet.FeeMember.Ref = feeMemberRef
	feeBalance := getBalanceNoErr(t, &launchnet.FeeMember, feeMemberRef)

	amountStr := "10"

	_, _, err = makeSignedRequest(launchnet.TestRPCUrlPublic, firstMember, "member.transfer",
		map[string]interface{}{"amount": amountStr, "toMemberReference": secondMember.Ref + ".record"})
	require.Error(t, err)

	data := checkConvertRequesterError(t, err).Data
	require.Contains(t, strings.Join(data.Trace, ": "), "provided reference is not object")

	newFirstBalance := getBalanceNoErr(t, firstMember, firstMember.Ref)
	newSecondBalance := getBalanceNoErr(t, secondMember, secondMember.Ref)
	newFeeBalance := getBalanceNoErr(t, &launchnet.FeeMember, feeMemberRef)
	require.Equal(t, oldFirstBalance, newFirstBalance)
	require.Equal(t, oldSecondBalance, newSecondBalance)
	require.Equal(t, feeBalance, newFeeBalance)
}

func TestTransferMoneyToNotSelfScopedRef(t *testing.T) {
	firstMember := createMember(t)
	secondMember := createMember(t)
	oldFirstBalance := getBalanceNoErr(t, firstMember, firstMember.Ref)
	oldSecondBalance := getBalanceNoErr(t, secondMember, secondMember.Ref)

	feeRes, err := signedRequest(t, launchnet.TestRPCUrlPublic, &launchnet.FeeMember, "member.get", nil)
	require.Nil(t, err)
	feeMemberRef, ok := feeRes.(map[string]interface{})["reference"].(string)
	require.True(t, ok)
	launchnet.FeeMember.Ref = feeMemberRef
	feeBalance := getBalanceNoErr(t, &launchnet.FeeMember, feeMemberRef)

	amountStr := "10"

	sepIdx := strings.Index(firstMember.Ref, ":")
	_, _, err = makeSignedRequest(
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
	require.Contains(t, strings.Join(data.Trace, ": "), "provided reference is not self-scoped")

	newFirstBalance := getBalanceNoErr(t, firstMember, firstMember.Ref)
	newSecondBalance := getBalanceNoErr(t, secondMember, secondMember.Ref)
	newFeeBalance := getBalanceNoErr(t, &launchnet.FeeMember, feeMemberRef)
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

	_, err := signedRequestWithEmptyRequestRef(t, launchnet.TestRPCUrlPublic, firstMember, "member.transfer",
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

	_, err := signedRequestWithEmptyRequestRef(t, launchnet.TestRPCUrlPublic, firstMember, "member.transfer",
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

	_, err := signedRequestWithEmptyRequestRef(t, launchnet.TestRPCUrlPublic, firstMember, "member.transfer",
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
	_, err := signedRequest(t, launchnet.TestRPCUrlPublic, firstMember, "member.transfer",
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

	_, err := signedRequestWithEmptyRequestRef(t, launchnet.TestRPCUrlPublic, firstMember, "member.transfer",
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

	_, err := signedRequestWithEmptyRequestRef(t, launchnet.TestRPCUrlPublic, member, "member.transfer",
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

	amount := "10"

	_, err := signedRequest(t, launchnet.TestRPCUrlPublic, firstMember, "member.transfer",
		map[string]interface{}{"amount": amount, "toMemberReference": secondMember.Ref})
	require.NoError(t, err)
	_, err = signedRequest(t, launchnet.TestRPCUrlPublic, firstMember, "member.transfer",
		map[string]interface{}{"amount": amount, "toMemberReference": secondMember.Ref})
	require.NoError(t, err)

	// Skip validation of balance before/after transfer
	// checkBalanceFewTimes(t, secondMember, secondMember.Ref, oldSecondBalance+2*amount)
	// newFirstBalance := getBalanceNoErr(t, firstMember, firstMember.Ref)
	// require.Equal(t, oldFirstBalance-2*amount, newFirstBalance)
}
