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
	"github.com/insolar/insolar/logicrunner/goplugin/experiment/foundation"
	"github.com/stretchr/testify/assert"
)

func TestFirst(t *testing.T) {
	toMember, toMemberRef := member.NewMember("Vasya")
	toMember.SetContext(&foundation.CallContext{
		Me: toMemberRef,
	})

	fromMember, fromMemberRef := member.NewMember("Petya")
	fromMember.SetContext(&foundation.CallContext{
		Me: fromMemberRef,
	})

	toWallet, toWalletRef := wallet.NewWallet(1000)
	toWallet.SetContext(&foundation.CallContext{
		Me: toWalletRef,
	})

	fromWallet, fromWalletRef := wallet.NewWallet(2000)
	fromWallet.SetContext(&foundation.CallContext{
		Me: fromWalletRef,
	})

	foundation.SetDelegate(fromMemberRef, &wallet.TypeReference, fromWallet)

	_, ok := fromMember.GetImplementationFor(&wallet.TypeReference).(*wallet.Wallet)
	assert.True(t, ok)

	toWallet.Receive(500, fromMemberRef)

	assert.Equal(t, uint(1500), fromWallet.GetTotalBalance())
	assert.Equal(t, uint(1500), toWallet.GetTotalBalance())
}
