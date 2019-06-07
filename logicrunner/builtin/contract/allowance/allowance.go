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

package allowance

import (
	"fmt"
	"time"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/logicrunner/builtin/proxy/wallet"
	"github.com/insolar/insolar/logicrunner/goplugin/foundation"
)

type Allowance struct {
	foundation.BaseContract
	To         insolar.Reference
	Amount     uint
	ExpireTime int64
}

func (a *Allowance) isExpired() bool {
	return a.GetContext().Time.After(time.Unix(a.ExpireTime, 0))
}

// TakeAmount allows take amount and delete allowance
func (a *Allowance) TakeAmount() (uint, error) {
	if *(a.GetContext().Caller) != a.To {
		return 0, fmt.Errorf("[ TakeAmount ] Only recepient can take amount")
	}
	if a.isExpired() {
		return 0, fmt.Errorf("[ TakeAmount ] Allowance expiried")
	}
	if err := a.SelfDestruct(); err != nil {
		return 0, err
	}
	return a.Amount, nil
}

// GetBalanceForOwner returns balance
func (a *Allowance) GetBalanceForOwner() (uint, error) {
	return a.Amount, nil
}

// GetExpiredBalance gets balance from expired allowance and delete allowance
func (a *Allowance) GetExpiredBalance() (uint, error) {
	if *(a.GetContext().Caller) != *(a.GetContext().Parent) {
		return 0, fmt.Errorf("[ DeleteExpiredAllowance ] Only owner can delete expiried Allowance")
	}
	if a.isExpired() {
		if err := a.SelfDestruct(); err != nil {
			return 0, err
		}
		return a.Amount, nil
	}
	return 0, nil
}

// New check is caller wallet and makes new allowance
func New(to *insolar.Reference, amount uint, expire int64) (*Allowance, error) {
	if !wallet.PrototypeReference.Equal(*foundation.GetContext().CallerPrototype) {
		return nil, fmt.Errorf("[ New Allowance ] : Can't create allowance from not wallet contract")
	}
	return &Allowance{To: *to, Amount: amount, ExpireTime: expire}, nil
}
