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
)

type Deposit struct {
	foundation.BaseContract
	OracleConfirms map[string]bool
	Confirms       uint
	TxHash         string
	UnHoldDate     string
	Amount         uint
}

func (d *Deposit) GetTxHash() (string, error) {
	return d.TxHash, nil
}

func (d *Deposit) GetAmount() (uint, error) {
	return d.Amount, nil
}

func New(oracleConfirms map[string]bool, txHash string, amount uint) (*Deposit, error) {
	return &Deposit{
		OracleConfirms: oracleConfirms,
		Confirms:       0,
		TxHash:         txHash,
		Amount:         amount,
	}, nil
}

func (d *Deposit) Confirm(oracleName string, txHash string, amount uint) (bool, error) {
	if txHash != d.TxHash {
		return false, fmt.Errorf("[ Confirm ] Transaction hash is incorrect")
	}

	if amount != d.Amount {
		return false, fmt.Errorf("[ Confirm ] Amount is incorrect")
	}

	if confirm, ok := d.OracleConfirms[oracleName]; ok {
		if confirm {
			return false, fmt.Errorf("[ Confirm ] Confirm from the oracle " + oracleName + " already exists")
		} else {
			d.OracleConfirms[oracleName] = true
			d.Confirms++
			if d.Confirms == 1 {
				return true, nil
			} else {
				return false, nil
			}
		}
	} else {
		return false, fmt.Errorf("[ Confirm ] Oracle name is incorrect")
	}
}
