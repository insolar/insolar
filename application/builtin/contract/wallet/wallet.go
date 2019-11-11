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

	"github.com/insolar/insolar/application/builtin/proxy/account"
	"github.com/insolar/insolar/application/builtin/proxy/deposit"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/logicrunner/builtin/foundation"
)

const XNS = "XNS"

// Wallet - basic wallet contract.
type Wallet struct {
	foundation.BaseContract
	Accounts foundation.StableMap
	Deposits foundation.StableMap
}

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

// GetAccount returns account ref
// ins:immutable
func (w *Wallet) GetAccount(assetName string) (*insolar.Reference, error) {
	accountReference, ok := w.Accounts[assetName]
	if !ok {
		return nil, fmt.Errorf("asset not found: %s", assetName)
	}
	return insolar.NewObjectReferenceFromString(accountReference)
}

// Transfer transfers money to given wallet.
// ins:immutable
func (w *Wallet) Transfer(
	assetName string, amountStr string, toMember *insolar.Reference,
	fromMember insolar.Reference, request insolar.Reference,
) (interface{}, error) {
	accRef, err := w.GetAccount(assetName)
	if err != nil {
		return nil, fmt.Errorf("failed to get account by asset name: %s", err.Error())
	}
	acc := account.GetObject(*accRef)
	return acc.Transfer(amountStr, toMember, fromMember, request)
}

// GetBalance gets balance by asset name.
// ins:immutable
func (w *Wallet) GetBalance(assetName string) (string, error) {
	accRef, err := w.GetAccount(assetName)
	if err != nil {
		return "", fmt.Errorf("failed to get account by asset: %s", err.Error())
	}
	acc := account.GetObject(*accRef)
	return acc.GetBalance()
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
func (w *Wallet) GetDeposits() ([]interface{}, error) {
	result := make([]interface{}, 0)
	for _, dRef := range w.Deposits {

		reference, err := insolar.NewObjectReferenceFromString(dRef)
		if err != nil {
			return nil, err
		}
		d := deposit.GetObject(*reference)

		depositInfo, err := d.Itself()
		if err != nil {
			return nil, fmt.Errorf("failed to get deposit itself: %s", err.Error())
		}

		result = append(result, depositInfo)
	}
	return result, nil
}

// FindDeposit finds deposit for this wallet with this transaction hash.
// ins:immutable
func (w *Wallet) FindDeposit(transactionHash string) (bool, *insolar.Reference, error) {
	if depositReferenceStr, ok := w.Deposits[transactionHash]; ok {
		depositReference, _ := insolar.NewObjectReferenceFromString(depositReferenceStr)
		return true, depositReference, nil
	}
	return false, nil, nil
}
