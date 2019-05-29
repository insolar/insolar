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
	"github.com/insolar/insolar/logicrunner/goplugin/foundation"
	"math/big"
	"time"
)

type DepositStatus string

const (
	DEPOSIT_CONFIRMS uint = 3

	OPEN    DepositStatus = "Open"
	HOLDING DepositStatus = "Holding"
	CLOSE   DepositStatus = "Close"
)

type Deposit struct {
	foundation.BaseContract
	Status         DepositStatus
	OracleConfirms map[string]bool
	Confirms       uint
	TxHash         string
	UnHoldDate     time.Time
	Amount         big.Int
}

func (d *Deposit) GetTxHash() (string, error) {
	return d.TxHash, nil
}

func (d *Deposit) GetAmount() (big.Int, error) {
	return d.Amount, nil
}

func New(oracleConfirms map[string]bool, txHash string, amount big.Int, unHoldDate time.Time) (*Deposit, error) {
	return &Deposit{
		Status:         OPEN,
		OracleConfirms: oracleConfirms,
		Confirms:       0,
		TxHash:         txHash,
		UnHoldDate:     unHoldDate,
		Amount:         amount,
	}, nil
}

func (d *Deposit) Confirm(oracleName string, txHash string, amount big.Int) (uint, error) {
	if txHash != d.TxHash {
		return 0, fmt.Errorf("[ Confirm ] Transaction hash is incorrect")
	}

	if (&amount).Cmp(&d.Amount) != 0 {
		return 0, fmt.Errorf("[ Confirm ] Amount is incorrect")
	}

	if confirm, ok := d.OracleConfirms[oracleName]; ok {
		if confirm {
			return 0, fmt.Errorf("[ Confirm ] Confirm from the oracle " + oracleName + " already exists")
		} else {
			d.OracleConfirms[oracleName] = true
			d.Confirms++
			if d.Confirms == DEPOSIT_CONFIRMS {
				d.Status = HOLDING
			}
			return d.Confirms, nil
		}
	} else {
		return 0, fmt.Errorf("[ Confirm ] Oracle name is incorrect")
	}
}
