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

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/logicrunner/builtin/contract/member"
	"github.com/insolar/insolar/logicrunner/builtin/foundation"
	"github.com/insolar/insolar/logicrunner/builtin/foundation/safemath"
	"github.com/insolar/insolar/logicrunner/builtin/proxy/costcenter"
	proxyMember "github.com/insolar/insolar/logicrunner/builtin/proxy/member"
	"github.com/insolar/insolar/logicrunner/builtin/proxy/rootdomain"
	"github.com/insolar/insolar/logicrunner/builtin/proxy/wallet"
)

// Wallet - basic wallet contract.
type Wallet struct {
	foundation.BaseContract
	Balance string
}

// New creates new wallet.
func New(balance string) (*Wallet, error) {
	return &Wallet{
		Balance: balance,
	}, nil
}

// Transfer transfers money to given wallet.
func (w *Wallet) Transfer(rootDomainRef insolar.Reference, amountStr string, toMember *insolar.Reference) (interface{}, error) {

	amount, ok := new(big.Int).SetString(amountStr, 10)
	if !ok {
		return nil, fmt.Errorf("can't parse input amount")
	}
	zero, _ := new(big.Int).SetString("0", 10)
	if amount.Cmp(zero) < 1 {
		return nil, fmt.Errorf("amount must be larger then zero")
	}

	rd := rootdomain.GetObject(rootDomainRef)
	ccRef, err := rd.GetCostCenter()
	if err != nil {
		return nil, fmt.Errorf("failed to get cost center reference: %s", err.Error())
	}

	cc := costcenter.GetObject(ccRef)
	feeStr, err := cc.CalcFee(amountStr)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate fee for amount: %s", err.Error())
	}
	fee, _ := new(big.Int).SetString(feeStr, 10)

	amountWithFee := new(big.Int).Add(fee, amount)

	balance, ok := new(big.Int).SetString(w.Balance, 10)
	if !ok {
		return nil, fmt.Errorf("can't parse wallet balance")
	}

	memberWallet, err := proxyMember.GetObject(*toMember).GetWallet()
	if err != nil {
		return nil, fmt.Errorf("failed to get member wallet: %s", err.Error())
	}

	newBalance, err := safemath.Sub(balance, amountWithFee)
	if err != nil {
		return nil, fmt.Errorf("not enough balance for transfer: %s", err.Error())
	}
	w.Balance = newBalance.String()

	fwRef, err := cc.GetFeeWalletRef()
	if err != nil {
		return nil, fmt.Errorf("failed to get fee wallet reference: %s", err.Error())
	}

	feeWallet := wallet.GetObject(fwRef)

	acceptFeeErr := feeWallet.Accept(feeStr)
	if acceptFeeErr != nil {
		newBalance, err = safemath.Add(balance, amountWithFee)
		if err != nil {
			return nil, fmt.Errorf("failed to add amount back to balance: %s", err.Error())
		}
		w.Balance = newBalance.String()
		return nil, fmt.Errorf("failed to transfer fee: %s", acceptFeeErr.Error())
	}

	toWallet := wallet.GetObject(memberWallet)
	acceptErr := toWallet.Accept(amount.String())
	if acceptErr == nil {
		return member.TransferResponse{Fee: feeStr}, nil
	}

	newBalance, err = safemath.Add(balance, amountWithFee)
	if err != nil {
		return nil, fmt.Errorf("failed to add amount back to balance: %s", err.Error())
	}
	w.Balance = newBalance.String()

	err = feeWallet.RollBack(feeStr)
	if err != nil {
		return nil, fmt.Errorf("failed to roll back fee: %s", err.Error())
	}

	return nil, fmt.Errorf("failed to accept balance to wallet: %s", acceptErr.Error())
}

// Accept accepts transfer to balance.
func (w *Wallet) Accept(amountStr string) (err error) {

	amount := new(big.Int)
	amount, ok := amount.SetString(amountStr, 10)
	if !ok {
		return fmt.Errorf("can't parse input amount")
	}

	balance := new(big.Int)
	balance, ok = balance.SetString(w.Balance, 10)
	if !ok {
		return fmt.Errorf("can't parse wallet balance")
	}

	b, err := safemath.Add(balance, amount)
	if err != nil {
		return fmt.Errorf("failed to add amount to balance: %s", err.Error())
	}
	w.Balance = b.String()

	return nil
}

// RollBack rolls back transfer to balance.
func (w *Wallet) RollBack(amountStr string) (err error) {

	amount := new(big.Int)
	amount, ok := amount.SetString(amountStr, 10)
	if !ok {
		return fmt.Errorf("can't parse input amount")
	}

	balance := new(big.Int)
	balance, ok = balance.SetString(w.Balance, 10)
	if !ok {
		return fmt.Errorf("can't parse wallet balance")
	}

	b, err := safemath.Sub(balance, amount)
	if err != nil {
		return fmt.Errorf("failed to sub amount from balance: %s", err.Error())
	}
	w.Balance = b.String()

	return nil
}

// GetBalance gets total balance.
func (w *Wallet) GetBalance() (string, error) {
	return w.Balance, nil
}
