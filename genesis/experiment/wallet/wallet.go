package wallet

import (
	"github.com/insolar/insolar/genesis/experiment/allowance"
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
	a := allowance.Allowance{To: to, Amount: amount, ExpireTime: w.GetContext().Time.Unix() + 10}
	w.AddChild(&a, &allowance.TypeReference)
	return &a
}

func (w *Wallet) Receive(amount uint, from *foundation.Reference) {
	//intr := foundation.GetObject(from)
	fromWallet := foundation.GetImplementationFor(from, &TypeReference).(*Wallet)

	a := fromWallet.Allocate(amount, w.GetContext().Me)
	w.balance += a.TakeAmount()
}

func (w *Wallet) Transfer(amount uint, to *foundation.Reference) {
	w.balance -= amount
	a := allowance.Allowance{To: to, Amount: amount, ExpireTime: w.GetContext().Time.Unix() + 10}
	w.AddChild(&a, &allowance.TypeReference)

	toWallet := foundation.GetImplementationFor(to, &TypeReference).(*Wallet)
	toWallet.Accept(&a)
}

func (w *Wallet) Accept(a *allowance.Allowance) {
	w.balance += a.TakeAmount()
}

func (w *Wallet) GetTotalBalance() uint {
	var totalAllowanced uint = 0
	for _, c := range w.GetChildrenTyped(&allowance.TypeReference) {
		Allowance := c.(*allowance.Allowance)
		totalAllowanced += Allowance.GetBalanceForOwner()
	}
	return w.balance + totalAllowanced
}

func (w *Wallet) ReturnAndDeleteExpiriedAllowances() {
	for _, c := range w.GetChildrenTyped(&allowance.TypeReference) {
		Allowance := c.(*allowance.Allowance)
		w.balance += Allowance.DeleteExpiredAllowance()
	}
}

func NewWallet(balance uint) (*Wallet, *foundation.Reference) {
	wallet := &Wallet{
		balance: balance,
	}
	reference := foundation.SaveToLedger(wallet)
	return wallet, reference
}
