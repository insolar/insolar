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

	"github.com/insolar/insolar/application/appfoundation"
	"github.com/insolar/insolar/application/builtin/proxy/costcenter"
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

// GetBalance gets total balance.
// ins:immutable
func (a *Account) GetBalance() (string, error) {
	return a.Balance, nil
}

// Transfer transfers money to given member.
func (a *Account) Transfer(
	amountStr string, toMember *insolar.Reference,
	fromMember insolar.Reference, request insolar.Reference,
) (interface{}, error) {

	amount, ok := new(big.Int).SetString(amountStr, 10)
	if !ok {
		return nil, fmt.Errorf("can't parse input amount")
	}
	if amount.Sign() <= 0 {
		return nil, fmt.Errorf("amount must be larger then zero")
	}

	ccRef := appfoundation.GetCostCenter()

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

	err = a.transferToMember(amountStr, *toMember, fromMember, request)
	if err != nil {
		return nil, fmt.Errorf("failed to transfer amount: %s", err.Error())
	}

	if feeStr != "0" {
		err = a.transferToMember(feeStr, *toFeeMember, fromMember, request)
		if err != nil {
			return nil, fmt.Errorf("failed to transfer fee: %s", err.Error())
		}
	}

	return member.TransferResponse{Fee: feeStr}, nil
}

// transferToMember transfers funds to member.
func (a *Account) transferToMember(
	amountStr string, toMember insolar.Reference, fromMember insolar.Reference, request insolar.Reference,
) error {
	to := member.GetObject(toMember)
	return to.Accept(appfoundation.SagaAcceptInfo{
		Amount:     amountStr,
		FromMember: fromMember,
		Request:    request,
	})
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
