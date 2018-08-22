package wallet

import (
	"github.com/insolar/insolar/logicrunner/goplugin/experiment/allowance"
	"github.com/insolar/insolar/logicrunner/goplugin/experiment/foundation"
	"github.com/insolar/insolar/logicrunner/goplugin/experiment/member"
)

// todo make this investigation through reflection
var TypeReference = foundation.Reference("wallet")

type Wallet struct {
	foundation.BaseContract
	balance uint
}

func (w *Wallet) Allocate(amount uint, to *foundation.Reference) *allowance.Allowance {
	// TODO check balance is enough
	w.balance -= amount
	return &allowance.Allowance{To: to, Amount: amount, ExpireTime: w.GetContext().Time.Unix() + 10}
}

func (w *Wallet) Receive(amount uint, from *foundation.Reference) {
	fromMember := member.GetObject(from)
	fromWallet := fromMember.GetImplementationFor(&TypeReference).(Wallet)
	Allowance := fromWallet.Allocate(amount, w.GetContext().Me)
	w.balance += Allowance.TakeAmount()
}

func (w *Wallet) GetTotalBalance() uint {
	var totalAllowanced uint = 0
	for _, c := range w.GetChildrenTyped(&allowance.TypeReference) {
		totalAllowanced += c.(allowance.Allowance).GetBalanceForOwner()
	}
	return w.balance + totalAllowanced
}

func (w *Wallet) ReturnAndDeleteExpiriedAllowances() {
	for _, c := range w.GetChildrenTyped(&allowance.TypeReference) {
		w.balance += c.(allowance.Allowance).DeleteExpiredAllowance()
	}
}
