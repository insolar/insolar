/*
 *    Copyright 2018 INS Ecosystem
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

package experiment

import (
	"testing"

	"github.com/insolar/insolar/genesis/experiment/member"
	"github.com/insolar/insolar/genesis/experiment/wallet"
	"github.com/insolar/insolar/toolkit/go/foundation"
	"github.com/stretchr/testify/assert"
)

func TestTransferViaReceive(t *testing.T) {
	// Create member which balance will increase
	toMember, toMemberRef := member.NewMember("Vasya")
	toMember.SetContext(&foundation.CallContext{
		Me: toMemberRef,
	})

	// Create member which balance will decrease
	fromMember, fromMemberRef := member.NewMember("Petya")
	fromMember.SetContext(&foundation.CallContext{
		Me: fromMemberRef,
	})

	// Create wallet for toMember
	toWallet, toWalletRef := wallet.NewWallet(1000)
	toWallet.SetContext(&foundation.CallContext{
		Me: toWalletRef,
	})

	// Create wallet for fromMember
	fromWallet, fromWalletRef := wallet.NewWallet(2000)
	fromWallet.SetContext(&foundation.CallContext{
		Me: fromWalletRef,
	})

	// Make fromWallet delegate of fromMember
	fromMember.TakeDelegate(fromWallet, &wallet.TypeReference)
	// Make toWallet delegate of toMember
	toMember.TakeDelegate(toWallet, &wallet.TypeReference)

	// Get fromMember as wallet instance
	fromMemberAsWallet, ok := foundation.GetImplementationFor(fromMemberRef, &wallet.TypeReference).(*wallet.Wallet)
	assert.True(t, ok)
	assert.NotNil(t, fromMemberAsWallet)

	// Get toMember as wallet instance
	toMemberAsWallet, ok := foundation.GetImplementationFor(toMemberRef, &wallet.TypeReference).(*wallet.Wallet)
	assert.True(t, ok)
	assert.NotNil(t, toMemberAsWallet)

	// Inject fake context of Caller for test
	foundation.InjectFakeContext(3, &foundation.CallContext{Caller: toWalletRef})

	// Call to get money from one member to another
	toMemberAsWallet.Receive(500, fromMemberRef)

	// Check balance
	assert.Equal(t, 1500, int(fromWallet.GetTotalBalance()))
	assert.Equal(t, 1500, int(toWallet.GetTotalBalance()))
}

func TestTransferViaTransfer(t *testing.T) {
	// Create member which balance will increase
	toMember, toMemberRef := member.NewMember("Vasya")
	toMember.SetContext(&foundation.CallContext{
		Me: toMemberRef,
	})

	// Create member which balance will decrease
	fromMember, fromMemberRef := member.NewMember("Petya")
	fromMember.SetContext(&foundation.CallContext{
		Me: fromMemberRef,
	})

	// Create wallet for toMember
	toWallet, toWalletRef := wallet.NewWallet(1000)
	toWallet.SetContext(&foundation.CallContext{
		Me: toWalletRef,
	})

	// Create wallet for fromMember
	fromWallet, fromWalletRef := wallet.NewWallet(2000)
	fromWallet.SetContext(&foundation.CallContext{
		Me: fromWalletRef,
	})

	// Make fromWallet delegate of fromMember
	fromMember.TakeDelegate(fromWallet, &wallet.TypeReference)
	// Make toWallet delegate of toMember
	toMember.TakeDelegate(toWallet, &wallet.TypeReference)

	// Get fromMember as wallet instance
	fromMemberAsWallet, ok := foundation.GetImplementationFor(fromMemberRef, &wallet.TypeReference).(*wallet.Wallet)
	assert.True(t, ok)
	assert.NotNil(t, fromMemberAsWallet)

	// Get toMember as wallet instance
	toMemberAsWallet, ok := foundation.GetImplementationFor(toMemberRef, &wallet.TypeReference).(*wallet.Wallet)
	assert.True(t, ok)
	assert.NotNil(t, toMemberAsWallet)

	// Inject fake context of Caller for test
	foundation.InjectFakeContext(2, &foundation.CallContext{Caller: toMemberRef})

	// Call to get money from one member to another
	fromMemberAsWallet.Transfer(500, toMemberRef)

	// Check balance
	assert.Equal(t, 1500, int(fromWallet.GetTotalBalance()))
	assert.Equal(t, 1500, int(toWallet.GetTotalBalance()))
}
