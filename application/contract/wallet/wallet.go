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
	"fmt"

	"github.com/insolar/insolar/application/contract/wallet/safemath"
	"github.com/insolar/insolar/application/proxy/allowance"
	"github.com/insolar/insolar/application/proxy/wallet"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/logicrunner/goplugin/foundation"
)

// Wallet - basic wallet contract
type Wallet struct {
	foundation.BaseContract
	Balance uint
}

// Transfer transfers money to given wallet
func (w *Wallet) Transfer(amount uint, to *core.RecordRef) error {

	toWallet, err := wallet.GetImplementationFrom(*to)
	if err != nil {
		return fmt.Errorf("[ Transfer ] Can't get implementation: %s", err.Error())
	}

	toWalletRef := toWallet.GetReference()

	newBalance, err := safemath.Sub(w.Balance, amount)
	if err != nil {
		return fmt.Errorf("[ Transfer ] Not enough balance for transfer: %s", err.Error())
	}

	ah := allowance.New(&toWalletRef, amount, w.GetContext().Time.Unix()+10)
	a, err := ah.AsChild(w.GetReference())
	if err != nil {
		return fmt.Errorf("[ Transfer ] Can't save as child: %s", err.Error())
	}

	// Changing balance only after allowance was successfully create
	w.Balance = newBalance

	r := a.GetReference()
	err = toWallet.AcceptNoWait(&r)
	return err
}

// Accept transforms allowance to balance
func (w *Wallet) Accept(aRef *core.RecordRef) error {
	b, err := allowance.GetObject(*aRef).TakeAmount()
	if err != nil {
		return fmt.Errorf("[ Accept ] Can't take amount: %s", err.Error())
	}
	w.Balance, err = safemath.Add(w.Balance, b)
	if err != nil {
		return fmt.Errorf("[ Accept ] Couldn't add amount to balance: %s", err.Error())
	}
	return nil
}

// GetBalance gets total balance
func (w *Wallet) GetBalance() (uint, error) {
	iterator, err := w.NewChildrenTypedIterator(allowance.GetPrototype())
	if err != nil {
		return 0, fmt.Errorf("[ GetBalance ] Can't get children: %s", err.Error())
	}

	for iterator.HasNext() {
		cref, err := iterator.Next()
		if err != nil {
			return 0, fmt.Errorf("[ GetBalance ] Can't get next child: %s", err.Error())
		}

		if !cref.IsEmpty() {
			a := allowance.GetObject(cref)
			balance, err := a.GetExpiredBalance()

			if err != nil {
				balance = 0
				//return 0, fmt.Errorf("[ GetBalance ] Can't get balance for owner: %s", err.Error())
			}

			w.Balance, err = safemath.Add(w.Balance, balance)
			if err != nil {
				return 0, fmt.Errorf("[ GetTotalBalance ] Couldn't add expired allowance to balance: %s", err.Error())
			}
		}
	}
	return w.Balance, nil
}

// New creates new allowance
func New(balance uint) (*Wallet, error) {
	return &Wallet{
		Balance: balance,
	}, nil
}
