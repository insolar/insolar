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

package wallet

import (
	"fmt"

	"github.com/insolar/insolar/application/contract/wallet/safemath"
	"github.com/insolar/insolar/application/proxy/wallet"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/logicrunner/goplugin/foundation"
)

// Wallet - basic wallet contract
type Wallet struct {
	foundation.BaseContract
	Balance uint
}

// New creates new wallet
func New(balance uint) (*Wallet, error) {
	return &Wallet{
		Balance: balance,
	}, nil
}

// Transfer transfers money to given wallet
func (w *Wallet) Transfer(amount uint, to *insolar.Reference) error {

	toWallet, err := wallet.GetImplementationFrom(*to)
	if err != nil {
		return fmt.Errorf("[ Transfer ] Can't get implementation: %s", err.Error())
	}

	newBalance, err := safemath.Sub(w.Balance, amount)
	if err != nil {
		return fmt.Errorf("[ Transfer ] Not enough balance for transfer: %s", err.Error())
	}
	w.Balance = newBalance

	acceptErr := toWallet.Accept(amount)
	if acceptErr != nil {
		newBalance, err := safemath.Add(w.Balance, amount)
		if err != nil {
			return fmt.Errorf("[ Transfer ] Couldn't add amount back to balance: %s", err.Error())
		}
		w.Balance = newBalance

		return fmt.Errorf("[ Transfer ] Cant accept balance to wallet: %s", acceptErr.Error())
	} else {
		return nil
	}
}

// Accept transfer to balance
func (w *Wallet) Accept(amount uint) (err error) {
	w.Balance, err = safemath.Add(w.Balance, amount)
	if err != nil {
		return fmt.Errorf("[ Accept ] Couldn't add amount to balance: %s", err.Error())
	}
	return nil
}

// GetBalance gets total balance
func (w *Wallet) GetBalance() (uint, error) {
	return w.Balance, nil
}
