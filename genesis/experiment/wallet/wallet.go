package wallet

import (
	"github.com/insolar/insolar/genesis/experiment/allowance"
	"github.com/insolar/insolar/genesis/experiment/member"
	"github.com/insolar/insolar/toolkit/go/foundation"
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
	fromWallet := fromMember.GetImplementationFor(&TypeReference).(*Wallet)
	Allowance := fromWallet.Allocate(amount, w.GetContext().Me)
	Allowance.SetContext(&foundation.CallContext{ // todo this is hack for testing
		Caller: w.GetContext().Me,
	})
	w.balance += Allowance.TakeAmount()
}

func (w *Wallet) GetTotalBalance() uint {
	var totalAllowanced uint = 0
	for _, c := range w.GetChildrenTyped(&allowance.TypeReference) {
		allowance := c.(allowance.Allowance)
		totalAllowanced += allowance.GetBalanceForOwner()
	}
	return w.balance + totalAllowanced
}

func (w *Wallet) ReturnAndDeleteExpiriedAllowances() {
	for _, c := range w.GetChildrenTyped(&allowance.TypeReference) {
		allowance := c.(allowance.Allowance)
		w.balance += allowance.DeleteExpiredAllowance()
	}
}

func NewWallet(balance uint) (*Wallet, *foundation.Reference) {
	wallet := &Wallet{
		balance: balance,
	}
	reference := foundation.SaveToLedger(wallet)
	return wallet, reference
}
