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
	"github.com/insolar/insolar/logicrunner/builtin/contract/deposit/safemath"
	"github.com/insolar/insolar/logicrunner/builtin/proxy/wallet"

	"github.com/insolar/insolar/logicrunner/builtin/foundation"
)

type status string

const (
	month = 30 * 24 * 60 * 60

	confirms uint = 3
	// offsetDepositPulse insolar.PulseNumber = 6 * month
	offsetDepositPulse insolar.PulseNumber = 10

	statusOpen    status = "Open"
	statusHolding status = "Holding"
	statusClose   status = "Close"
)

// Deposit is like wallet. It holds migrated money.
type Deposit struct {
	foundation.BaseContract
	PulseDepositCreate      insolar.PulseNumber `json:"timestamp"`
	PulseDepositHold        insolar.PulseNumber `json:"holdStartDate"`
	PulseDepositUnHold      insolar.PulseNumber `json:"holdReleaseDate"`
	MigrationDaemonConfirms [3]string           `json:"confirmerReferences"`
	Amount                  string              `json:"amount"`
	Bonus                   string              `json:"bonus"`
	TxHash                  string              `json:"ethTxHash"`
}

// GetTxHash gets transaction hash.
func (d *Deposit) GetTxHash() (string, error) {
	return d.TxHash, nil
}

// GetAmount gets amount.
func (d *Deposit) GetAmount() (string, error) {
	return d.Amount, nil
}

// New creates new deposit.
func New(migrationDaemonConfirms [3]string, txHash string, amount string) (*Deposit, error) {
	currentPulse, err := foundation.GetPulseNumber()
	if err != nil {
		return nil, fmt.Errorf("failed to get current pulse: %s", err.Error())
	}
	return &Deposit{
		PulseDepositCreate:      currentPulse,
		MigrationDaemonConfirms: migrationDaemonConfirms,
		Amount:                  amount,
		TxHash:                  txHash,
	}, nil
}

func calculateUnHoldPulse(currentPulse insolar.PulseNumber) insolar.PulseNumber {
	return currentPulse + offsetDepositPulse
}

// Itself gets deposit information.
func (d *Deposit) Itself() (interface{}, error) {
	return *d, nil
}

// Confirm adds confirm for deposit by migration daemon.
func (d *Deposit) Confirm(migrationDaemonIndex int, migrationDaemonRef string, txHash string, amountStr string) error {
	if txHash != d.TxHash {
		return fmt.Errorf("transaction hash is incorrect")
	}

	inputAmount := new(big.Int)
	inputAmount, ok := inputAmount.SetString(amountStr, 10)
	if !ok {
		return fmt.Errorf("failed to parse input amount")
	}
	depositAmount := new(big.Int)
	depositAmount, ok = depositAmount.SetString(d.Amount, 10)
	if !ok {
		return fmt.Errorf("failed to parse deposit amount")
	}

	if (inputAmount).Cmp(depositAmount) != 0 {
		return fmt.Errorf("deposit with this transaction hash has different amount")
	}

	if d.MigrationDaemonConfirms[migrationDaemonIndex] != "" {
		return fmt.Errorf("confirm from the '%v' migration daemon already exists; member '%s' already confirmed it", migrationDaemonIndex, migrationDaemonRef)
	} else {
		d.MigrationDaemonConfirms[migrationDaemonIndex] = migrationDaemonRef

		n := 0
		for _, c := range d.MigrationDaemonConfirms {
			if c != "" {
				n++
			}
		}
		if uint(n) >= confirms {
			currentPulse, err := foundation.GetPulseNumber()
			if err != nil {
				return fmt.Errorf("failed to get current pulse: %s", err.Error())
			}
			d.PulseDepositHold = currentPulse
			d.PulseDepositUnHold = calculateUnHoldPulse(currentPulse)
		}
		return nil
	}
}

func (d *Deposit) canTransfer() error {
	c := 0
	for _, r := range d.MigrationDaemonConfirms {
		if r != "" {
			c++
		}
	}
	if c < 3 {
		return fmt.Errorf("number of confirms is less then 3")
	}

	p, err := foundation.GetPulseNumber()
	if err != nil {
		return fmt.Errorf("failed to get pulse number: %s", err.Error())
	}
	if d.PulseDepositUnHold > p {
		return fmt.Errorf("hold period didn't end")
	}

	return nil
}

// Transfer transfers money from deposit to wallet.It can be called only after deposit hold period.
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

	err = d.canTransfer()
	if err != nil {
		return nil, fmt.Errorf("can't start transfer: %s", err.Error())
	}

	d.Amount = newBalance.String()

	w := wallet.GetObject(wallerRef)

	acceptWalletErr := w.Accept(amountStr)
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
