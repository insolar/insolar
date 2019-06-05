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
	"math/big"

	"github.com/insolar/insolar/application/contract/wallet/safemath"
	"github.com/insolar/insolar/application/proxy/wallet"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/logicrunner/goplugin/foundation"
)

// Wallet - basic wallet contract
type Wallet struct {
	foundation.BaseContract
	Balance string
}

// New creates new wallet
func New(balance string) (*Wallet, error) {
	return &Wallet{
		Balance: balance,
	}, nil
}

// Transfer transfers money to given wallet
func (w *Wallet) Transfer(amountStr string, toMember *insolar.Reference) error {

	amount := new(big.Int)
	amount, ok := amount.SetString(amountStr, 10)
	if !ok {
		return fmt.Errorf("[ Transfer ] can't parse returned balance")
	}
	balance := new(big.Int)
	balance, ok = balance.SetString(w.Balance, 10)
	if !ok {
		return fmt.Errorf("[ Transfer ] can't parse returned balance")
	}

	toWallet, err := wallet.GetImplementationFrom(*toMember)
	if err != nil {
		return fmt.Errorf("[ Transfer ] Can't get implementation: %s", err.Error())
	}

	newBalance, err := safemath.Sub(balance, amount)
	if err != nil {
		return fmt.Errorf("[ Transfer ] Not enough balance for transfer: %s", err.Error())
	}
	w.Balance = newBalance.String()

	acceptErr := toWallet.Accept(amount.String())
	if acceptErr != nil {
		newBalance, err := safemath.Add(balance, amount)
		if err != nil {
			return fmt.Errorf("[ Transfer ] Couldn't add amount back to balance: %s", err.Error())
		}
		w.Balance = newBalance.String()

		return fmt.Errorf("[ Transfer ] Cant accept balance to wallet: %s", acceptErr.Error())
	} else {
		return nil
	}
}

// Accept transfer to balance
func (w *Wallet) Accept(amountStr string) (err error) {

	amount := new(big.Int)
	amount, ok := amount.SetString(amountStr, 10)
	if !ok {
		return fmt.Errorf("[ Accept ] can't parse returned balance")
	}

	balance := new(big.Int)
	balance, ok = balance.SetString(w.Balance, 10)
	if !ok {
		return fmt.Errorf("[ Accept ] can't parse returned balance")
	}

	b, err := safemath.Add(balance, amount)
	if err != nil {
		return fmt.Errorf("[ Accept ] Couldn't add amount to balance: %s", err.Error())
	}
	w.Balance = b.String()

	return nil
}

// GetBalance gets total balance
func (w *Wallet) GetBalance() (string, error) {
	return w.Balance, nil
}
