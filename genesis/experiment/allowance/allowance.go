/*
 *    Copyright 2018 Insolar
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package allowance

import (
	"time"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/logicrunner/goplugin/foundation"
)

type Allowance struct {
	foundation.BaseContract
	To         core.RecordRef
	Amount     uint
	ExpireTime int64
}

func (a *Allowance) IsExpired() bool {
	return a.GetContext().Time.After(time.Unix(a.ExpireTime, 0))
}

func (a *Allowance) TakeAmount() uint {
	caller := a.GetContext().Caller
	if caller == a.To && !a.IsExpired() {
		a.SelfDestructRequest()
		r := a.Amount
		a.Amount = 0
		return r
	}
	return 0
}

func (a *Allowance) GetBalanceForOwner() uint {
	if !a.IsExpired() {
		return a.Amount
	}
	return 0
}

func (a *Allowance) DeleteExpiredAllowance() uint {
	if a.GetContext().Caller == a.GetContext().Parent && !a.IsExpired() {
		a.SelfDestructRequest()
		return a.Amount
	}
	return 0
}

func New(to core.RecordRef, amount uint, expire int64) *Allowance {
	return &Allowance{To: to, Amount: amount, ExpireTime: expire}
}
