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

	"github.com/insolar/insolar/application/proxy/allowance"
	"github.com/insolar/insolar/application/proxy/wallet"
)

// Wallet - basic wallet contract
type Wallet struct {
	foundation.BaseContract
	Balance uint
}

// Allocate - returns reference to a new allowance
func (w *Wallet) Allocate(amount uint, to *core.RecordRef) (core.RecordRef, error) {
	// TODO check balance is enough
	w.Balance -= amount
	ah := allowance.New(to, amount, w.GetContext().Time.Unix()+10)
	a := ah.AsChild(w.GetReference())
	return a.GetReference(), nil
}

func (w *Wallet) Receive(amount uint, from *core.RecordRef) error {
	fromWallet := wallet.GetImplementationFrom(*from)

	v := w.GetReference()
	aRef, err := fromWallet.Allocate(amount, &v)
	if err != nil {
		return err
	}

	b, err := allowance.GetObject(aRef).TakeAmount()
	if err != nil {
		return err
	}

	w.Balance += b

	return nil
}

func (w *Wallet) Transfer(amount uint, to *core.RecordRef) error {
	w.Balance -= amount

	toWallet := wallet.GetImplementationFrom(*to)
	toWalletRef := toWallet.GetReference()

	ah := allowance.New(&toWalletRef, amount, w.GetContext().Time.Unix()+10)
	a := ah.AsChild(w.GetReference())

	r := a.GetReference()
	toWallet.Accept(&r)
	return nil
}

func (w *Wallet) Accept(aRef *core.RecordRef) error {
	b, err := allowance.GetObject(*aRef).TakeAmount()
	if err != nil {
		return err
	}
	w.Balance += b
	return nil
}

func (w *Wallet) GetTotalBalance() (uint, error) {
	var totalAllowanced uint
	crefs, err := w.GetChildrenTyped(allowance.GetClass())
	if err != nil {
		return 0, err
	}
	for _, cref := range crefs {
		a := allowance.GetObject(cref)
		balance, err := a.GetBalanceForOwner()
		if err != nil {
			return 0, err
		}

		totalAllowanced += balance
	}
	return w.Balance + totalAllowanced, nil
}

func (w *Wallet) ReturnAndDeleteExpiredAllowances() error {
	crefs, err := w.GetChildrenTyped(allowance.GetClass())
	if err != nil {
		return err
	}
	for _, cref := range crefs {
		Allowance := allowance.GetObject(cref)
		balance, err := Allowance.DeleteExpiredAllowance()
		if err != nil {
			return err
		}
		w.Balance += balance
	}
	return err
}

func New(balance uint) (*Wallet, error) {
	return &Wallet{
		Balance: balance,
	}, nil
}
