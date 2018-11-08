/*
 *    Copyright 2018 Insolar
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package functest

import (
	"testing"

	"github.com/insolar/insolar/testutils"
	"github.com/stretchr/testify/assert"
)

func TestTransferMoney(t *testing.T) {
	firstMember := createMember(t, "Member1")
	secondMember := createMember(t, "Member2")
	oldFirstBalance := getBalanceNoErr(t, firstMember, firstMember.ref)
	oldSecondBalance := getBalanceNoErr(t, secondMember, secondMember.ref)

	amount := 111

	_, err := signedRequest(firstMember, "Transfer", amount, secondMember.ref)
	assert.NoError(t, err)

	newFirstBalance := getBalanceNoErr(t, firstMember, firstMember.ref)
	newSecondBalance := getBalanceNoErr(t, secondMember, secondMember.ref)
	assert.Equal(t, oldFirstBalance-amount, newFirstBalance)
	assert.Equal(t, oldSecondBalance+amount, newSecondBalance)
}

func TestTransferMoneyFromNotExist(t *testing.T) {
	firstMember := createMember(t, "Member1")
	firstMember.ref = testutils.RandomRef().String()

	secondMember := createMember(t, "Member2")
	oldSecondBalance := getBalanceNoErr(t, secondMember, secondMember.ref)

	amount := 111

	_, err := signedRequest(firstMember, "Transfer", amount, secondMember.ref)
	assert.EqualError(t, err, "Can't get public key: couldn't get object message: couldn't get object: failed to fetch object index: storage object not found")

	newSecondBalance := getBalanceNoErr(t, secondMember, secondMember.ref)
	assert.Equal(t, oldSecondBalance, newSecondBalance)
}

func TestTransferMoneyToNotExist(t *testing.T) {
	firstMember := createMember(t, "Member1")
	oldFirstBalance := getBalanceNoErr(t, firstMember, firstMember.ref)

	amount := 111

	_, err := signedRequest(firstMember, "Transfer", amount, testutils.RandomRef())
	assert.EqualError(t, err, "[ Transfer ] Can't get implementation: on calling main API: failed to fetch object index: storage object not found")

	newFirstBalance := getBalanceNoErr(t, firstMember, firstMember.ref)
	assert.Equal(t, oldFirstBalance, newFirstBalance)
}

func TestTransferNegativeAmount(t *testing.T) {
	firstMember := createMember(t, "Member1")
	secondMember := createMember(t, "Member2")
	oldFirstBalance := getBalanceNoErr(t, firstMember, firstMember.ref)
	oldSecondBalance := getBalanceNoErr(t, secondMember, secondMember.ref)

	amount := -111

	_, err := signedRequest(firstMember, "Transfer", amount, secondMember.ref)
	assert.EqualError(t, err, "[ transferCall ] Amount must be positive")

	newFirstBalance := getBalanceNoErr(t, firstMember, firstMember.ref)
	newSecondBalance := getBalanceNoErr(t, secondMember, secondMember.ref)
	assert.Equal(t, oldFirstBalance, newFirstBalance)
	assert.Equal(t, oldSecondBalance, newSecondBalance)
}

func TestTransferAllAmount(t *testing.T) {
	firstMember := createMember(t, "Member1")
	secondMember := createMember(t, "Member2")
	oldFirstBalance := getBalanceNoErr(t, firstMember, firstMember.ref)
	oldSecondBalance := getBalanceNoErr(t, secondMember, secondMember.ref)

	amount := oldFirstBalance

	_, err := signedRequest(firstMember, "Transfer", amount, secondMember.ref)
	assert.NoError(t, err)

	newFirstBalance := getBalanceNoErr(t, firstMember, firstMember.ref)
	newSecondBalance := getBalanceNoErr(t, secondMember, secondMember.ref)
	assert.Equal(t, 0, newFirstBalance)
	assert.Equal(t, oldSecondBalance+oldFirstBalance, newSecondBalance)
}

func TestTransferMoreThanAvailableAmount(t *testing.T) {
	firstMember := createMember(t, "Member1")
	secondMember := createMember(t, "Member2")
	oldFirstBalance := getBalanceNoErr(t, firstMember, firstMember.ref)
	oldSecondBalance := getBalanceNoErr(t, secondMember, secondMember.ref)

	amount := oldFirstBalance + 100

	_, err := signedRequest(firstMember, "Transfer", amount, secondMember.ref)
	assert.EqualError(t, err, "[ Transfer ] Not enough balance for transfer: subtrahend must be smaller than minuend")

	newFirstBalance := getBalanceNoErr(t, firstMember, firstMember.ref)
	newSecondBalance := getBalanceNoErr(t, secondMember, secondMember.ref)
	assert.Equal(t, oldFirstBalance, newFirstBalance)
	assert.Equal(t, oldSecondBalance, newSecondBalance)
}

func TestTransferToMyself(t *testing.T) {
	member := createMember(t, "Member1")
	oldMemberBalance := getBalanceNoErr(t, member, member.ref)

	amount := 100

	_, err := signedRequest(member, "Transfer", amount, member.ref)
	assert.EqualError(t, err, "[ transferCall ] Recipient must be different from the sender")

	newMemberBalance := getBalanceNoErr(t, member, member.ref)
	assert.Equal(t, oldMemberBalance, newMemberBalance)
}

// TODO: test to check overflow of balance
// TODO: check transfer zero amount

func TestTransferTwoTimes(t *testing.T) {
	firstMember := createMember(t, "Member1")
	secondMember := createMember(t, "Member2")
	oldFirstBalance := getBalanceNoErr(t, firstMember, firstMember.ref)
	oldSecondBalance := getBalanceNoErr(t, secondMember, secondMember.ref)

	amount := 100

	_, err := signedRequest(firstMember, "Transfer", amount, secondMember.ref)
	assert.NoError(t, err)
	_, err = signedRequest(firstMember, "Transfer", amount, secondMember.ref)
	assert.NoError(t, err)

	newFirstBalance := getBalanceNoErr(t, firstMember, firstMember.ref)
	newSecondBalance := getBalanceNoErr(t, secondMember, secondMember.ref)
	assert.Equal(t, oldFirstBalance-2*amount, newFirstBalance)
	assert.Equal(t, oldSecondBalance+2*amount, newSecondBalance)
}
