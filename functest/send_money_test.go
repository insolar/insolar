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

func _TestTransferNegativeAmount(t *testing.T) {
	firstMember := createMember(t, "Member1")
	secondMember := createMember(t, "Member2")

	amount := -111

	_, err := signedRequest(firstMember, "Transfer", amount, secondMember.ref)
	assert.Error(t, err)
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

/*func _TestTransferMoreThanAvailableAmount(t *testing.T) {
	firstMemberRef := createMember(t)
	secondMemberRef := createMember(t)
	oldFirstBalance := getBalance(t, firstMemberRef)
	oldSecondBalance := getBalance(t, secondMemberRef)

	body := getResponseBody(t, postParams{
		"query_type": "send_money",
		"from":       secondMemberRef,
		"to":         firstMemberRef,
		"amount":     10000000000,
	})

	transferResponse := &sendMoneyResponse{}
	unmarshalResponse(t, body, transferResponse)

	// Add checking than contract gives specific error

	newFirstBalance := getBalance(t, firstMemberRef)
	newSecondBalance := getBalance(t, secondMemberRef)

	assert.Equal(t, oldFirstBalance, newFirstBalance)
	assert.Equal(t, oldSecondBalance, newSecondBalance)
}

func _TestTransferToMyself(t *testing.T) {
	memberRef := createMember(t)
	oldBalance := getBalance(t, memberRef)

	body := getResponseBody(t, postParams{
		"query_type": "send_money",
		"from":       memberRef,
		"to":         memberRef,
		"amount":     oldBalance - 1,
	})

	transferResponse := &sendMoneyResponse{}
	unmarshalResponse(t, body, transferResponse)

	assert.True(t, transferResponse.Success)

	newBalance := getBalance(t, memberRef)

	assert.Equal(t, oldBalance, newBalance)
}*/

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
