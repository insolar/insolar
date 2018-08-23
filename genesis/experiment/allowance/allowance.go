package allowance

import (
	"time"

	"github.com/insolar/insolar/toolkit/go/foundation"
)

var TypeReference = foundation.Reference("allowance")

type Allowance struct {
	foundation.BaseContract
	To         *foundation.Reference
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
