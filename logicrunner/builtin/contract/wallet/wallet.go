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
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/logicrunner/builtin/foundation"
	"github.com/insolar/insolar/logicrunner/builtin/foundation/safemath"
	"github.com/insolar/insolar/logicrunner/builtin/proxy/account"
	"github.com/insolar/insolar/logicrunner/builtin/proxy/costcenter"
	proxyMember "github.com/insolar/insolar/logicrunner/builtin/proxy/member"
	"github.com/insolar/insolar/logicrunner/builtin/proxy/rootdomain"
	"math/big"
)

// Wallet - basic wallet contract.
type Wallet struct {
	foundation.BaseContract
	Accounts foundation.StableMap
}

// New creates new wallet.
func New(rootDomain insolar.Reference) (*Wallet, error) {
	accs := make(foundation.StableMap)
	newAccount, _ := account.New("0").AsChild(rootDomain)

	accs["XNS"] = newAccount.GetReference().String()

	return &Wallet{
		Accounts: accs,
	}, nil
}

func (w *Wallet) GetAccount(assetName string) (*insolar.Reference, error) {
	return insolar.NewReferenceFromBase58(w.Accounts[assetName])
}

// Transfer transfers money to given wallet.
func (w *Wallet) Transfer(rootDomainRef insolar.Reference, assetName string, amountStr string, toMember *insolar.Reference) (interface{}, error) {
	amount, ok := new(big.Int).SetString(amountStr, 10)
	if !ok {
		return nil, fmt.Errorf("can't parse input amount")
	}
	zero, _ := new(big.Int).SetString("0", 10)
	if amount.Cmp(zero) == -1 {
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

	amount, _ = new(big.Int).SetString(amountStr, 10)
	fee, _ := new(big.Int).SetString(feeStr, 10)
	totalSum, err := safemath.Add(fee, amount)
	currentBalanceStr, err := w.GetBalance("XNS")
	currentBalance, _ := new(big.Int).SetString(currentBalanceStr, 10)
	if totalSum.Cmp(currentBalance) > 0  {
		return nil, fmt.Errorf("balance is too low: %s", currentBalanceStr)
	}

	toWallet, err := proxyMember.GetObject(*toMember).GetWallet()
	if err != nil {
		return nil, fmt.Errorf("failed to get destination member wallet: %s", err.Error())
	}
	accRef, _ := w.GetAccount(assetName)
	acc := account.GetObject(*accRef)
	err = acc.Transfer(amountStr, &toWallet)


	feeWalletRef,_ := cc.GetFeeWalletRef()
	err = acc.Transfer(feeStr, &feeWalletRef)

	return nil, fmt.Errorf("failed to accept balance to wallet: %s", err)
}

func (w *Wallet) GetBalance(assetName string) (string, error) {
	accRef, _ := w.GetAccount(assetName)
	acc := account.GetObject(*accRef)
	return acc.GetBalance()
}

func(w *Wallet) Accept(amountStr string, assetName string) error {
	accRef, _ := w.GetAccount(assetName)
	acc := account.GetObject(*accRef)
	return acc.Accept(amountStr)
}
