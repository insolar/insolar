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

	"github.com/insolar/insolar/application/builtin/proxy/costcenter"
	"github.com/insolar/insolar/application/builtin/proxy/deposit"
	"github.com/insolar/insolar/application/builtin/proxy/member"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/logicrunner/builtin/foundation"
	"github.com/insolar/insolar/logicrunner/builtin/foundation/safemath"
)

type Account struct {
	foundation.BaseContract
	Balance string
}

func New(balance string) (*Account, error) {
	return &Account{Balance: balance}, nil
}

type destination interface {
	Accept(string) error
}

// Transfer transfers funds to giver reference.
func (a *Account) transfer(amountStr string, destinationObject destination) error {
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
	return destinationObject.Accept(amountStr)
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

// TransferToDeposit transfers funds to deposit.
func (a *Account) TransferToDeposit(amountStr string, toDeposit insolar.Reference) error {
	to := deposit.GetObject(toDeposit)
	return a.transfer(amountStr, to)
}

// TransferToMember transfers funds to member.
func (a *Account) TransferToMember(amountStr string, toMember insolar.Reference) error {
	to := member.GetObject(toMember)
	return to.Accept(amountStr)
}

// GetBalance gets total balance.
// ins:immutable
func (a *Account) GetBalance() (string, error) {
	return a.Balance, nil
}

// Transfer transfers money to given member.
func (a *Account) Transfer(rootDomainRef insolar.Reference, amountStr string, toMember *insolar.Reference) (interface{}, error) {

	amount, ok := new(big.Int).SetString(amountStr, 10)
	if !ok {
		return nil, fmt.Errorf("can't parse input amount")
	}
	if amount.Sign() <= 0 {
		return nil, fmt.Errorf("amount must be larger then zero")
	}

	ccRef := foundation.GetCostCenter()

	cc := costcenter.GetObject(ccRef)
	feeStr, err := cc.CalcFee(amountStr)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate fee for amount: %s", err.Error())
	}

	var toFeeMember *insolar.Reference
	if feeStr != "0" {
		fee, ok := new(big.Int).SetString(feeStr, 10)
		if !ok {
			return nil, fmt.Errorf("can't parse input feeStr")
		}

		toFeeMember, err = cc.GetFeeMember()
		if err != nil {
			return nil, fmt.Errorf("failed to get fee member: %s", err.Error())
		}

		amount, err = safemath.Add(fee, amount)
		if err != nil {
			return nil, fmt.Errorf("failed to calculate totalSum for amount: %s", err.Error())
		}
	}

	currentBalanceStr, err := a.GetBalance()
	if err != nil {
		return nil, fmt.Errorf("failed to get balance for asset: %s", err.Error())
	}
	currentBalance, ok := new(big.Int).SetString(currentBalanceStr, 10)
	if !ok {
		return nil, fmt.Errorf("can't parse account balance")
	}
	if amount.Cmp(currentBalance) > 0 {
		return nil, fmt.Errorf("balance is too low: %s", currentBalanceStr)
	}

	newBalance, err := safemath.Sub(currentBalance, amount)
	if err != nil {
		return nil, fmt.Errorf("not enough balance for transfer: %s", err.Error())
	}
	a.Balance = newBalance.String()

	err = a.TransferToMember(amountStr, *toMember)
	if err != nil {
		return nil, fmt.Errorf("failed to transfer amount: %s", err.Error())
	}

	if feeStr != "0" {
		err = a.TransferToMember(feeStr, *toFeeMember)
		if err != nil {
			return nil, fmt.Errorf("failed to transfer fee: %s", err.Error())
		}
	}

	return member.TransferResponse{Fee: feeStr}, nil
}

// IncreaseBalance increases the current balance by the amount.
func (a *Account) IncreaseBalance(amountStr string) error {
	amount, ok := new(big.Int).SetString(amountStr, 10)
	if !ok {
		return fmt.Errorf("can't parse input amount")
	}
	if amount.Sign() <= 0 {
		return fmt.Errorf("amount should be greater then zero")
	}
	balance, ok := new(big.Int).SetString(a.Balance, 10)
	if !ok {
		return fmt.Errorf("can't parse account balance")
	}
	newBalance, err := safemath.Add(balance, amount)
	if err != nil {
		return fmt.Errorf("failed to add amount to balance: %s", err.Error())
	}
	a.Balance = newBalance.String()
	return nil
}
