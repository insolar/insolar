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

package deposit

import (
	"fmt"
	"math/big"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/logicrunner/builtin/foundation/safemath"
	"github.com/insolar/insolar/logicrunner/builtin/proxy/account"
	"github.com/insolar/insolar/logicrunner/builtin/proxy/member"
	"github.com/insolar/insolar/logicrunner/builtin/proxy/migrationadmin"
	"github.com/insolar/insolar/logicrunner/builtin/proxy/wallet"

	"github.com/insolar/insolar/logicrunner/builtin/foundation"
)

type status string

const (
	// TODO: https://insolar.atlassian.net/browse/WLT-768
	// day   = 24 * 60 * 60
	day   = 10
	month = 30 * day

	vestingPeriodInDays = 360

	// TODO: https://insolar.atlassian.net/browse/WLT-768
	// offsetDepositPulse insolar.PulseNumber = 6 * month
	offsetDepositPulse insolar.PulseNumber = 10

	XNS = "XNS"
)

// Deposit is like wallet. It holds migrated money.
type Deposit struct {
	foundation.BaseContract
	Balance                 string               `json:"balance"`
	PulseDepositUnHold      insolar.PulseNumber  `json:"holdReleaseDate"`
	MigrationDaemonConfirms foundation.StableMap `json:"confirmerReferences"`
	Amount                  string               `json:"amount"`
	Bonus                   string               `json:"bonus"`
	TxHash                  string               `json:"ethTxHash"`
}

// GetTxHash gets transaction hash.
// ins:immutable
func (d Deposit) GetTxHash() (string, error) {
	return d.TxHash, nil
}

// GetAmount gets amount.
// ins:immutable
func (d Deposit) GetAmount() (string, error) {
	return d.Amount, nil
}

// Return pulse of unhold deposit.
// ins:immutable
func (d *Deposit) GetPulseUnHold() (insolar.PulseNumber, error) {
	return d.PulseDepositUnHold, nil
}

// New creates new deposit.
func New(migrationDaemonRef insolar.Reference, txHash string, amount string) (*Deposit, error) {

	migrationDaemonConfirms := make(foundation.StableMap)
	migrationDaemonConfirms[migrationDaemonRef.String()] = amount

	return &Deposit{
		Balance:                 "0",
		MigrationDaemonConfirms: migrationDaemonConfirms,
		Amount:                  "0",
		TxHash:                  txHash,
	}, nil
}

func calculateUnHoldPulse(currentPulse insolar.PulseNumber) insolar.PulseNumber {
	return currentPulse + offsetDepositPulse
}

// Itself gets deposit information.
// ins:immutable
func (d Deposit) Itself() (interface{}, error) {
	return d, nil
}

// Confirm adds confirm for deposit by migration daemon.
func (d *Deposit) Confirm(migrationDaemonRef string, txHash string, amountStr string) error {
	if txHash != d.TxHash {
		return fmt.Errorf("transaction hash is incorrect")
	}
	if _, ok := d.MigrationDaemonConfirms[migrationDaemonRef]; ok {
		return fmt.Errorf("confirm from this migration daemon already exists: '%s' ", migrationDaemonRef)
	}
	d.MigrationDaemonConfirms[migrationDaemonRef] = amountStr

	if len(d.MigrationDaemonConfirms) > 2 {
		migrationAdminContract := migrationadmin.GetObject(foundation.GetMigrationAdmin())
		activeDaemons, err := migrationAdminContract.GetActiveDaemons()
		if err != nil {
			return fmt.Errorf("failed to get list active daemons: %s", err.Error())
		}
		err = d.checkAmount(activeDaemons)
		if err != nil {
			return fmt.Errorf("failed to check amount in confirmation from migration daemon: '%s'", err.Error())
		}
		currentPulse, err := foundation.GetPulseNumber()
		if err != nil {
			return fmt.Errorf("failed to get current pulse: %s", err.Error())
		}
		d.Amount = amountStr
		d.PulseDepositUnHold = calculateUnHoldPulse(currentPulse)

		ma := member.GetObject(foundation.GetMigrationAdminMember())
		accountRef, err := ma.GetAccount(XNS)
		if err != nil {
			return fmt.Errorf("get account ref failed: %s", err.Error())
		}
		a := account.GetObject(*accountRef)
		err = a.TransferToDeposit(amountStr, d.GetReference())
		if err != nil {
			return fmt.Errorf("failed to transfer from migration wallet to deposit: %s", err.Error())
		}
	}
	return nil
}

// Check amount field in confirmation from migration daemons.
func (d *Deposit) checkAmount(activeDaemons []string) error {
	if len(activeDaemons) > 0 {
		amount := d.MigrationDaemonConfirms[activeDaemons[0]]
		for i := 0; i < insolar.GenesisAmountActiveMigrationDaemonMembers; i++ {
			if amount != d.MigrationDaemonConfirms[activeDaemons[i]] {
				return fmt.Errorf(" several migration daemons send different amount  ")
			}
		}
		return nil
	}
	return fmt.Errorf(" list with migration daemons member is empty ")
}
func (d *Deposit) canTransfer(transferAmount *big.Int) error {
	c := 0
	for _, r := range d.MigrationDaemonConfirms {
		if r != "" {
			c++
		}
	}
	if c < 3 {
		return fmt.Errorf("number of confirms is less then 3")
	}

	currentPulse, err := foundation.GetPulseNumber()
	if err != nil {
		return fmt.Errorf("failed to get pulse number: %s", err.Error())
	}
	if d.PulseDepositUnHold > currentPulse {
		return fmt.Errorf("hold period didn't end")
	}

	spentPeriodInDays := big.NewInt(int64((currentPulse - d.PulseDepositUnHold) / day))
	amount, ok := new(big.Int).SetString(d.Amount, 10)
	if !ok {
		return fmt.Errorf("can't parse derposit amount")
	}
	balance, ok := new(big.Int).SetString(d.Balance, 10)
	if !ok {
		return fmt.Errorf("can't parse derposit balance")
	}

	// How much can we transfer for this time
	availableForNow := new(big.Int).Div(
		new(big.Int).Mul(amount, spentPeriodInDays),
		big.NewInt(vestingPeriodInDays),
	)

	if new(big.Int).Sub(amount, availableForNow).Cmp(
		new(big.Int).Sub(balance, transferAmount),
	) == 1 {
		return fmt.Errorf("not enough unholded balance for transfer")
	}

	return nil
}

// Transfer transfers money from deposit to wallet. It can be called only after deposit hold period.
func (d *Deposit) Transfer(amountStr string, wallerRef insolar.Reference) (interface{}, error) {

	amount, ok := new(big.Int).SetString(amountStr, 10)
	if !ok {
		return nil, fmt.Errorf("can't parse input amount")
	}
	zero, _ := new(big.Int).SetString("0", 10)
	if amount.Cmp(zero) == -1 {
		return nil, fmt.Errorf("amount must be larger then zero")
	}

	balance, ok := new(big.Int).SetString(d.Amount, 10)
	if !ok {
		return nil, fmt.Errorf("can't parse deposit balance")
	}
	newBalance, err := safemath.Sub(balance, amount)
	if err != nil {
		return nil, fmt.Errorf("not enough balance for transfer: %s", err.Error())
	}

	err = d.canTransfer(amount)
	if err != nil {
		return nil, fmt.Errorf("can't start transfer: %s", err.Error())
	}

	d.Amount = newBalance.String()

	w := wallet.GetObject(wallerRef)

	acceptWalletErr := w.Accept(amountStr, XNS)
	if acceptWalletErr == nil {
		return nil, nil
	}

	newBalance, err = safemath.Add(balance, amount)
	if err != nil {
		return nil, fmt.Errorf("failed to add amount back to balance: %s", err.Error())
	}
	d.Amount = newBalance.String()
	return nil, fmt.Errorf("failed to transfer amount: %s", acceptWalletErr.Error())
}

// Accept accepts transfer to balance.
// ins:saga(INS_FLAG_NO_ROLLBACK_METHOD)
func (d *Deposit) Accept(amountStr string) error {

	amount := new(big.Int)
	amount, ok := amount.SetString(amountStr, 10)
	if !ok {
		return fmt.Errorf("can't parse input amount")
	}

	balance := new(big.Int)
	balance, ok = balance.SetString(d.Balance, 10)
	if !ok {
		return fmt.Errorf("can't parse deposit balance")
	}

	b, err := safemath.Add(balance, amount)
	if err != nil {
		return fmt.Errorf("failed to add amount to balance: %s", err.Error())
	}
	d.Balance = b.String()

	return nil
}
