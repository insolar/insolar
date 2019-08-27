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
	"github.com/insolar/insolar/logicrunner/builtin/proxy/deposit"
	"math/big"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/logicrunner/builtin/foundation"
	"github.com/insolar/insolar/logicrunner/builtin/foundation/safemath"
	"github.com/insolar/insolar/logicrunner/builtin/proxy/account"
	"github.com/insolar/insolar/logicrunner/builtin/proxy/costcenter"
	"github.com/insolar/insolar/logicrunner/builtin/proxy/member"
)

// Wallet - basic wallet contract.
type Wallet struct {
	foundation.BaseContract
	Accounts foundation.StableMap
	Deposits foundation.StableMap
}

const XNS = "XNS"

// New creates new wallet.
func New(accountReference insolar.Reference) (*Wallet, error) {
	if accountReference.IsEmpty() {
		return nil, fmt.Errorf("reference is empty")
	}
	accounts := make(foundation.StableMap)
	// TODO: Think about creating of new types of assets and initial balance
	accounts[XNS] = accountReference.String()

	return &Wallet{
		Accounts: accounts,
		Deposits: make(foundation.StableMap),
	}, nil
}

func (w *Wallet) GetAccount(assetName string) (*insolar.Reference, error) {
	accountReference, ok := w.Accounts[assetName]
	if !ok {
		return nil, fmt.Errorf("asset not found: %s", assetName)
	}
	return insolar.NewReferenceFromBase58(accountReference)
}

// Transfer transfers money to given wallet.
func (w *Wallet) Transfer(rootDomainRef insolar.Reference, assetName string, amountStr string, toMember *insolar.Reference) (interface{}, error) {
	amount, ok := new(big.Int).SetString(amountStr, 10)
	if !ok {
		return nil, fmt.Errorf("can't parse input amount")
	}
	zero, _ := new(big.Int).SetString("0", 10)
	if amount.Cmp(zero) < 1 {
		return nil, fmt.Errorf("amount must be larger then zero")
	}

	ccRef := foundation.GetCostCenter()

	cc := costcenter.GetObject(ccRef)
	feeStr, err := cc.CalcFee(amountStr)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate fee for amount: %s", err.Error())
	}

	amount, ok = new(big.Int).SetString(amountStr, 10)
	if !ok {
		return nil, fmt.Errorf("can't parse input amountStr")
	}
	fee, ok := new(big.Int).SetString(feeStr, 10)
	if !ok {
		return nil, fmt.Errorf("can't parse input feeStr")
	}
	totalSum, err := safemath.Add(fee, amount)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate totalSum for amount: %s", err.Error())
	}
	currentBalanceStr, err := w.GetBalance(assetName)
	if err != nil {
		return nil, fmt.Errorf("failed to get balance for asset: %s", err.Error())
	}
	currentBalance, _ := new(big.Int).SetString(currentBalanceStr, 10)
	if totalSum.Cmp(currentBalance) > 0 {
		return nil, fmt.Errorf("balance is too low: %s", currentBalanceStr)
	}

	toAccount, err := member.GetObject(*toMember).GetAccount(assetName)
	if err != nil {
		return nil, fmt.Errorf("failed to get account by asset name: %s", err.Error())
	}

	accRef, err := w.GetAccount(assetName)
	if err != nil {
		return nil, fmt.Errorf("failed to get account by asset name: %s", err.Error())
	}
	acc := account.GetObject(*accRef)
	err = acc.TransferToAccount(amountStr, *toAccount)
	if err != nil {
		return nil, fmt.Errorf("failed to transfer: %s", err.Error())
	}

	toFeeAccount, err := cc.GetFeeAccount()
	if err != nil {
		return nil, fmt.Errorf("failed to get fee account: %s", err.Error())
	}

	err = acc.TransferToAccount(feeStr, toFeeAccount)
	if err != nil {
		return nil, fmt.Errorf("failed to transfer: %s", err.Error())
	}

	return member.TransferResponse{Fee: feeStr}, nil
}

// GetBalance gets balance by asset name.
func (w *Wallet) GetBalance(assetName string) (string, error) {
	accRef, err := w.GetAccount(assetName)
	if err != nil {
		return "", fmt.Errorf("failed to get account by asset: %s", err.Error())
	}
	acc := account.GetObject(*accRef)
	return acc.GetBalance()
}

// Accept accepts transfer.
func (w *Wallet) Accept(amountStr string, assetName string) error {
	accRef, err := w.GetAccount(assetName)
	if err != nil {
		return fmt.Errorf("failed to get account by asset: %s", err.Error())
	}
	acc := account.GetObject(*accRef)
	return acc.Accept(amountStr)
}

// AddDeposit method stores deposit reference in member it belongs to
func (w *Wallet) AddDeposit(txId string, deposit insolar.Reference) error {
	if _, ok := w.Deposits[txId]; ok {
		return fmt.Errorf("deposit for this transaction already exist")
	}
	w.Deposits[txId] = deposit.String()
	return nil
}

// GetDeposits get all deposits for this wallet
// ins:immutable
func (w *Wallet) GetDeposits() (map[string]interface{}, error) {
	result := map[string]interface{}{}
	for tx, dRef := range w.Deposits {

		reference, err := insolar.NewReferenceFromBase58(dRef)
		if err != nil {
			return nil, err
		}
		d := deposit.GetObject(*reference)

		depositInfo, err := d.Itself()
		if err != nil {
			return nil, fmt.Errorf("failed to get deposit itself: %s", err.Error())
		}

		result[tx] = depositInfo
	}
	return result, nil
}

// FindDeposit finds deposit for this wallet with this transaction hash.
// ins:immutable
func (w *Wallet) FindDeposit(transactionHash string) (bool, *insolar.Reference, error) {
	if depositReferenceStr, ok := w.Deposits[transactionHash]; ok {
		depositReference, _ := insolar.NewReferenceFromBase58(depositReferenceStr)
		return true, depositReference, nil
	}

	return false, nil, nil
}
