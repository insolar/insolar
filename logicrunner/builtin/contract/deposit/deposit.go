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
	"time"

	"github.com/insolar/insolar/insolar"

	"github.com/insolar/insolar/logicrunner/builtin/foundation"
)

type status string

const (
	confirms uint = 3
)

const (
	statusOpen    status = "Open"
	statusHolding status = "Holding"
	statusClose   status = "Close"
)

// Deposit is like wallet. It holds migrated money.
type Deposit struct {
	foundation.BaseContract
	Timestamp               time.Time
	HoldReleaseDate         time.Time
	MigrationDaemonConfirms foundation.StableMap
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
func New(migrationDaemonConfirms foundation.StableMap, txHash string, amount string, holdReleaseDate time.Time) (*Deposit, error) {
	return &Deposit{

		MigrationDaemonConfirms: migrationDaemonConfirms,
		Confirms:                0,
		TxHash:                  txHash,
		HoldReleaseDate:         holdReleaseDate,
		Amount:                  amount,
		Status:                  statusOpen,
	}, nil
}

// MapMarshal gets deposit information.
func (d *Deposit) MapMarshal() ([][2]string, error) {
	return [][2]string{
		{"timestamp", d.Timestamp.String()},
		{"holdReleaseDate", d.HoldReleaseDate.String()},
		{"amount", d.Amount},
		{"bonus", d.Bonus},
		{"txId", d.TxHash},
	}, nil
}

// Confirm adds confirm for deposit by migration daemon.
func (d *Deposit) Confirm(migrationDaemon insolar.Reference, txHash string, amountStr string) (uint, error) {
	if txHash != d.TxHash {
		return 0, fmt.Errorf("transaction hash is incorrect")
	}

	inputAmount := new(big.Int)
	inputAmount, ok := inputAmount.SetString(amountStr, 10)
	if !ok {
		return 0, fmt.Errorf("failed to parse input amount")
	}
	depositAmount := new(big.Int)
	depositAmount, ok = depositAmount.SetString(d.Amount, 10)
	if !ok {
		return 0, fmt.Errorf("failed to parse deposit amount")
	}

	if (inputAmount).Cmp(depositAmount) != 0 {
		return 0, fmt.Errorf("amount is incorrect")
	}

	if confirm, ok := d.MigrationDaemonConfirms.Get(migrationDaemon); ok {
		if confirm, ok := confirm.(bool); ok && confirm {
			return 0, fmt.Errorf("confirm from the migration daemon '%s' already exists", migrationDaemon.String())
		} else {
			d.MigrationDaemonConfirms.Set(migrationDaemon, true)
			d.Confirms++
			if d.Confirms == confirms {
				d.Status = statusHolding
			}
			return d.Confirms, nil
		}
	} else {
		return 0, fmt.Errorf("migration daemon name is incorrect")
	}
}
