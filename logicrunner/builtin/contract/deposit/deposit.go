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

	"github.com/insolar/insolar/logicrunner/builtin/foundation"
)

type status string

const (
	month = 30 * 24 * 60 * 60

	confirms           uint                = 3
	offsetDepositPulse insolar.PulseNumber = 6 * month

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
func NewDeposit(migrationDaemonConfirms [3]string, txHash string, amount string) (*Deposit, error) {
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
