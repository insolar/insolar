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

package account

import (
	"fmt"
	"math/big"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/logicrunner/builtin/foundation"
	"github.com/insolar/insolar/logicrunner/builtin/foundation/safemath"
	"github.com/insolar/insolar/logicrunner/builtin/proxy/account"
)

type Account struct {
	foundation.BaseContract
	Balance string
}

func New(balance string) (*Account, error) {
	return &Account{Balance: balance}, nil
}

// Transfer transfers funds to giver reference.
func (a *Account) Transfer(amountStr string, toAccount *insolar.Reference) error {
	amount, ok := new(big.Int).SetString(amountStr, 10)
	if !ok {
		return fmt.Errorf("can't parse amountStr")
	}
	balance, ok := new(big.Int).SetString(a.Balance, 10)
	if !ok {
		return fmt.Errorf("can't parse account balance")
	}

	newBalance, err := safemath.Sub(balance, amount)
	if err != nil {
		return fmt.Errorf("not enough balance for transfer: %s", err.Error())
	}
	a.Balance = newBalance.String()
	destWallet := account.GetObject(*toAccount)
	return destWallet.Accept(amountStr)
}

// Accept accepts transfer to balance.
//ins:saga(INS_FLAG_NO_ROLLBACK_METHOD)
func (a *Account) Accept(amountStr string) error {

	amount := new(big.Int)
	amount, ok := amount.SetString(amountStr, 10)
	if !ok {
		return fmt.Errorf("can't parse input amount")
	}

	balance := new(big.Int)
	balance, ok = balance.SetString(a.Balance, 10)
	if !ok {
		return fmt.Errorf("can't parse account balance")
	}

	b, err := safemath.Add(balance, amount)
	if err != nil {
		return fmt.Errorf("failed to add amount to balance: %s", err.Error())
	}
	a.Balance = b.String()

	return nil
}

// RollBack rolls back transfer to balance.
func (a *Account) RollBack(amountStr string) error {

	amount := new(big.Int)
	amount, ok := amount.SetString(amountStr, 10)
	if !ok {
		return fmt.Errorf("can't parse input amount")
	}

	balance := new(big.Int)
	balance, ok = balance.SetString(a.Balance, 10)
	if !ok {
		return fmt.Errorf("can't parse account balance")
	}

	b, err := safemath.Sub(balance, amount)
	if err != nil {
		return fmt.Errorf("failed to sub amount from balance: %s", err.Error())
	}
	a.Balance = b.String()

	return nil
}

// GetBalance gets total balance.
// ins:immutable
func (a *Account) GetBalance() (string, error) {
	return a.Balance, nil
}
