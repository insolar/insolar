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

	"github.com/insolar/insolar/logicrunner/goplugin/foundation"
)

type status string

const month = 30 * 24 * 60 * 60
const (
	confirms           uint                = 3
	offSetDepositPulse insolar.PulseNumber = 6 * month
)

const (
	statusOpen    status = "Open"
	statusHolding status = "Holding"
	statusClose   status = "Close"
)

// Deposit is like wallet. It holds migrated money.
type Deposit struct {
	foundation.BaseContract
	PulseDepositCreate      insolar.PulseNumber
	PulseUnHoldDeposit      insolar.PulseNumber
	MigrationDaemonConfirms map[insolar.Reference]bool
	Confirms                uint
	Amount                  string
	Bonus                   string
	TxHash                  string
	Status                  status
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
func New(migrationDaemonConfirms map[insolar.Reference]bool, txHash string, amount string, currentPulse insolar.PulseNumber) (*Deposit, error) {
	return &Deposit{
		PulseDepositCreate:      currentPulse,
		PulseUnHoldDeposit:      calculateUnHoldPulse(currentPulse),
		MigrationDaemonConfirms: migrationDaemonConfirms,
		Confirms:                0,
		Amount:                  amount,
		TxHash:                  txHash,
		Status:                  statusOpen,
	}, nil
}

func calculateUnHoldPulse(currentPulse insolar.PulseNumber) insolar.PulseNumber {
	return currentPulse + offSetDepositPulse
}

// MapMarshal gets deposit information.
func (d *Deposit) MapMarshal() (map[string]string, error) {
	return map[string]string{
		"pulseDepositCreate": d.PulseDepositCreate.String(),
		"pulseUnHoldDeposit": d.PulseUnHoldDeposit.String(),
		"amount":             d.Amount,
		"bonus":              d.Bonus,
		"txId":               d.TxHash,
	}, nil
}

// Confirm adds confirm for deposit by migration daemon.
func (d *Deposit) Confirm(migrationDaemon insolar.Reference, txHash string, amountStr string) error {
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
		return fmt.Errorf("amount is incorrect")
	}
	if confirm, ok := d.MigrationDaemonConfirms[migrationDaemon]; ok {
		if confirm {
			return fmt.Errorf("confirm from the migration daemon '%s' already exists", migrationDaemon.String())
		} else {
			d.MigrationDaemonConfirms[migrationDaemon] = true
			d.Confirms++
			if d.Confirms == confirms {
				d.Status = statusHolding
			}
			return nil
		}
	} else {
		return fmt.Errorf("migration daemon name is incorrect")
	}
}
