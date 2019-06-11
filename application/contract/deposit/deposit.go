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
	"strconv"
	"time"

	"github.com/insolar/insolar/logicrunner/goplugin/foundation"
)

type DepositStatus string

const (
	DepositConfirms uint = 3

	Open    DepositStatus = "Open"
	Holding DepositStatus = "Holding"
	Close   DepositStatus = "Close"
)

type Deposit struct {
	foundation.BaseContract
	Status         DepositStatus
	OracleConfirms map[string]bool
	Confirms       uint
	TxHash         string
	CreationDate   time.Time
	UnHoldDate     time.Time
	Amount         string
	Bonus          string
}

func (d *Deposit) GetTxHash() (string, error) {
	return d.TxHash, nil
}

func (d *Deposit) GetAmount() (string, error) {
	return d.Amount, nil
}

func New(oracleConfirms map[string]bool, txHash string, amount string, unHoldDate time.Time) (*Deposit, error) {
	return &Deposit{
		Status:         Open,
		OracleConfirms: oracleConfirms,
		Confirms:       0,
		TxHash:         txHash,
		UnHoldDate:     unHoldDate,
		Amount:         amount,
	}, nil
}

func (d *Deposit) MapMarshal() (map[string]string, error) {
	return map[string]string{
		"Status":   string(d.Status),
		"Confirms": strconv.Itoa(int(d.Confirms)),
		"TxHash":   d.TxHash,
		"Amount":   d.Amount,
	}, nil
}

func (d *Deposit) Confirm(oracleName string, txHash string, amountStr string) (uint, error) {
	if txHash != d.TxHash {
		return 0, fmt.Errorf("[ Confirm ] Transaction hash is incorrect")
	}

	inputAmount := new(big.Int)
	inputAmount, ok := inputAmount.SetString(amountStr, 10)
	if !ok {
		return 0, fmt.Errorf("[ Confirm ] can't parse input amount")
	}
	depositAmount := new(big.Int)
	depositAmount, ok = depositAmount.SetString(d.Amount, 10)
	if !ok {
		return 0, fmt.Errorf("[ Confirm ] can't parse Deposit amount")
	}

	if (inputAmount).Cmp(depositAmount) != 0 {
		return 0, fmt.Errorf("[ Confirm ] Amount is incorrect")
	}

	if confirm, ok := d.OracleConfirms[oracleName]; ok {
		if confirm {
			return 0, fmt.Errorf("[ Confirm ] Confirm from the oracle " + oracleName + " already exists")
		} else {
			d.OracleConfirms[oracleName] = true
			d.Confirms++
			if d.Confirms == DepositConfirms {
				d.Status = Holding
			}
			return d.Confirms, nil
		}
	} else {
		return 0, fmt.Errorf("[ Confirm ] Oracle name is incorrect")
	}
}
