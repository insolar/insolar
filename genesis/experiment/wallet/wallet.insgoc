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
	"github.com/insolar/insolar/logicrunner/goplugin/foundation"

	"contract-proxy/allowance"
	"contract-proxy/wallet"
)

// Wallet - basic wallet contract
type Wallet struct {
	foundation.BaseContract
	Balance uint
}

// Allocate - returns reference to a new allowance
func (w *Wallet) Allocate(amount uint, to core.RecordRef) core.RecordRef {
	// TODO check balance is enough
	w.Balance -= amount
	ah := allowance.New(to, amount, w.GetContext().Time.Unix()+10)
	a := ah.AsChild(w.GetReference())
	return a.GetReference()
}

func (w *Wallet) Receive(amount uint, from core.RecordRef) {
	fromWallet := foundation.GetImplementationFor(from, w.GetClass()).(*wallet.Wallet)

	aRef := fromWallet.Allocate(amount, w.GetReference())
	w.Balance += allowance.GetObject(aRef).TakeAmount()
}

func (w *Wallet) Transfer(amount uint, to core.RecordRef) {
	w.Balance -= amount

	toWallet := foundation.GetImplementationFor(to, w.GetClass()).(*wallet.Wallet)
	toWalletRef := toWallet.GetReference()

	ah := allowance.New(toWalletRef, amount, w.GetContext().Time.Unix()+10)
	a := ah.AsChild(w.GetReference())

	toWallet.Accept(a.GetReference())
}

func (w *Wallet) Accept(aRef core.RecordRef) {
	w.Balance += allowance.GetObject(aRef).TakeAmount()
}

func (w *Wallet) GetTotalBalance() uint {
	var totalAllowanced uint
	for _, c := range w.GetChildrenTyped(allowance.GetClass()) {
		a := c.(*allowance.Allowance)
		totalAllowanced += a.GetBalanceForOwner()
	}
	return w.Balance + totalAllowanced
}

func (w *Wallet) ReturnAndDeleteExpiriedAllowances() {
	for _, c := range w.GetChildrenTyped(allowance.GetClass()) {
		Allowance := c.(*allowance.Allowance)
		w.Balance += Allowance.DeleteExpiredAllowance()
	}
}

func New(balance uint) *Wallet {
	return &Wallet{
		Balance: balance,
	}
}
