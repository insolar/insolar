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

package wallet

import (
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/genesis/experiment/allowance"
	"github.com/insolar/insolar/toolkit/go/foundation"
)

var TypeReference = core.String2Ref("w")

type Wallet struct {
	foundation.BaseContract
	Balance uint
}

func (w *Wallet) Allocate(amount uint, to core.RecordRef) *allowance.Allowance {
	// TODO check balance is enough
	w.Balance -= amount
	a := allowance.NewAllowance(to, amount, w.GetContext().Time.Unix()+10)
	w.AddChild(a, allowance.TypeReference)
	return a
}

func (w *Wallet) Receive(amount uint, from core.RecordRef) {
	fromWallet := foundation.GetImplementationFor(from, TypeReference).(*Wallet)

	a := fromWallet.Allocate(amount, w.GetContext().Me)
	w.Balance += a.TakeAmount()
}

func (w *Wallet) Transfer(amount uint, to core.RecordRef) {
	w.Balance -= amount

	toWallet := foundation.GetImplementationFor(to, TypeReference).(*Wallet)
	toWalletRef := toWallet.GetReference()

	a := allowance.NewAllowance(toWalletRef, amount, w.GetContext().Time.Unix()+10)
	w.AddChild(a, allowance.TypeReference)

	toWallet.Accept(a)
}

func (w *Wallet) Accept(a *allowance.Allowance) {
	w.Balance += a.TakeAmount()
}

func (w *Wallet) GetTotalBalance() uint {
	var totalAllowanced uint
	for _, c := range w.GetChildrenTyped(allowance.TypeReference) {
		Allowance := c.(*allowance.Allowance)
		totalAllowanced += Allowance.GetBalanceForOwner()
	}
	return w.Balance + totalAllowanced
}

func (w *Wallet) ReturnAndDeleteExpiriedAllowances() {
	for _, c := range w.GetChildrenTyped(allowance.TypeReference) {
		Allowance := c.(*allowance.Allowance)
		w.Balance += Allowance.DeleteExpiredAllowance()
	}
}

func NewWallet(balance uint) *Wallet {
	wallet := &Wallet{
		Balance: balance,
	}
	return wallet
}
