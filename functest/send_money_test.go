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

	"github.com/insolar/insolar/testutils"
	"github.com/stretchr/testify/require"
)

const times = 5

func checkBalanceFewTimes(t *testing.T, caller *user, ref string, expected big.Int) {
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
	// Skip validation of balance before/after transfer
	// oldFirstBalance := getBalanceNoErr(t, firstMember, firstMember.ref)
	// oldSecondBalance := getBalanceNoErr(t, secondMember, secondMember.ref)

	amount := "10"

	_, err := signedRequest(firstMember, "wallet.transfer", map[string]interface{}{"amount": amount, "toMemberReference": secondMember.ref})
	require.NoError(t, err)

	// Skip validation of balance before/after transfer
	// checkBalanceFewTimes(t, secondMember, secondMember.ref, oldSecondBalance+amount)
	// newFirstBalance := getBalanceNoErr(t, firstMember, firstMember.ref)
	// require.Equal(t, oldFirstBalance-amount, newFirstBalance)
}

func TestTransferMoneyFromNotExist(t *testing.T) {
	firstMember := createMember(t)
	firstMember.ref = testutils.RandomRef().String()

	secondMember := createMember(t)
	oldSecondBalance := getBalanceNoErr(t, secondMember, secondMember.ref)

	amount := "10"

	_, err := signedRequest(firstMember, "wallet.transfer", map[string]interface{}{"amount": amount, "toMemberReference": secondMember.ref})
	require.NotNil(t, err)
	require.Contains(t, err.Error(), "index not found")

	newSecondBalance := getBalanceNoErr(t, secondMember, secondMember.ref)
	require.Equal(t, oldSecondBalance, newSecondBalance)
}

func TestTransferMoneyToNotExist(t *testing.T) {
	firstMember := createMember(t)
	oldFirstBalance := getBalanceNoErr(t, firstMember, firstMember.ref)

	amount := "10"

	_, err := signedRequest(firstMember, "wallet.transfer", map[string]interface{}{"amount": amount, "toMemberReference": testutils.RandomRef().String()})
	require.NotNil(t, err)
	require.Contains(t, err.Error(), "index not found")

	newFirstBalance := getBalanceNoErr(t, firstMember, firstMember.ref)
	require.Equal(t, oldFirstBalance, newFirstBalance)
}

func TestTransferNegativeAmount(t *testing.T) {
	firstMember := createMember(t)
	secondMember := createMember(t)
	oldFirstBalance := getBalanceNoErr(t, firstMember, firstMember.ref)
	oldSecondBalance := getBalanceNoErr(t, secondMember, secondMember.ref)

	amount := "-111"

	_, err := signedRequest(firstMember, "wallet.transfer", map[string]interface{}{"amount": amount, "toMemberReference": secondMember.ref})
	require.Error(t, err)

	newFirstBalance := getBalanceNoErr(t, firstMember, firstMember.ref)
	newSecondBalance := getBalanceNoErr(t, secondMember, secondMember.ref)
	require.Equal(t, oldFirstBalance, newFirstBalance)
	require.Equal(t, oldSecondBalance, newSecondBalance)
}

// TODO: unskip test after undoing of all transaction in failed request will be supported
func TestTransferAllAmount(t *testing.T) {
	t.Skip()
	firstMember := createMember(t)
	secondMember := createMember(t)
	oldFirstBalance := getBalanceNoErr(t, firstMember, firstMember.ref)
	oldSecondBalance := getBalanceNoErr(t, secondMember, secondMember.ref)

	amount := oldFirstBalance

	summ := new(big.Int)
	summ.Add(oldSecondBalance, oldFirstBalance)

	_, err := signedRequest(firstMember, "wallet.transfer", map[string]interface{}{"amount": amount, "toMemberReference": secondMember.ref})
	require.NoError(t, err)

	checkBalanceFewTimes(t, secondMember, secondMember.ref, *summ)
	newFirstBalance := getBalanceNoErr(t, firstMember, firstMember.ref)
	require.Equal(t, 0, newFirstBalance)
}

func TestTransferMoreThanAvailableAmount(t *testing.T) {
	firstMember := createMember(t)
	secondMember := createMember(t)
	oldFirstBalance := getBalanceNoErr(t, firstMember, firstMember.ref)
	oldSecondBalance := getBalanceNoErr(t, secondMember, secondMember.ref)

	amount := new(big.Int)
	amount.Add(oldFirstBalance, big.NewInt(10))

	_, err := signedRequest(firstMember, "wallet.transfer", map[string]interface{}{"amount": amount.String(), "toMemberReference": secondMember.ref})
	require.NotNil(t, err)
	require.Contains(t, err.Error(), "subtrahend must be smaller than minuend")
	newFirstBalance := getBalanceNoErr(t, firstMember, firstMember.ref)
	newSecondBalance := getBalanceNoErr(t, secondMember, secondMember.ref)
	require.Equal(t, oldFirstBalance, newFirstBalance)
	require.Equal(t, oldSecondBalance, newSecondBalance)
}

func TestTransferToMyself(t *testing.T) {
	member := createMember(t)
	oldMemberBalance := getBalanceNoErr(t, member, member.ref)

	amount := "20"

	_, err := signedRequest(member, "wallet.transfer", map[string]interface{}{"amount": amount, "toMemberReference": member.ref})
	require.NotNil(t, err)
	require.Contains(t, err.Error(), "recipient must be different from the sender")
	newMemberBalance := getBalanceNoErr(t, member, member.ref)
	require.Equal(t, oldMemberBalance, newMemberBalance)
}

// TODO: test to check overflow of balance
// TODO: check transfer zero amount

// TODO: uncomment after undoing of all transaction in failed request will be supported
func TestTransferTwoTimes(t *testing.T) {
	firstMember := createMember(t)
	secondMember := createMember(t)
	// Skip validation of balance before/after transfer
	// oldFirstBalance := getBalanceNoErr(t, firstMember, firstMember.ref)
	// oldSecondBalance := getBalanceNoErr(t, secondMember, secondMember.ref)

	amount := "10"

	_, err := signedRequest(firstMember, "wallet.transfer", map[string]interface{}{"amount": amount, "toMemberReference": secondMember.ref})
	require.NoError(t, err)
	_, err = signedRequest(firstMember, "wallet.transfer", map[string]interface{}{"amount": amount, "toMemberReference": secondMember.ref})
	require.NoError(t, err)

	// Skip validation of balance before/after transfer
	// checkBalanceFewTimes(t, secondMember, secondMember.ref, oldSecondBalance+2*amount)
	// newFirstBalance := getBalanceNoErr(t, firstMember, firstMember.ref)
	// require.Equal(t, oldFirstBalance-2*amount, newFirstBalance)
}
