package wallet


import (
	"ilya/v2/allowance"
	"ilya/v2/memberProxy"
	mfm "ilya/v2/mockMagic"
)

type Wallet struct {
	mfm.MockMagic
	balance uint
}

func (w *Wallet) Allocate(amount uint, to *mfm.Reference) *allowance.Allowance {
	// TODO check balance is enough
	w.balance -= amount
	return &allowance.Allowance{} //to: to, amount: amount, expTime: 0} // TODO Set real exp time
}

func (w *Wallet) Receive(amount uint, from *mfm.Reference) {
	memberSender := memberProxy.ProxyGetObject(from)
	interfaceSender := memberSender.ProxyGetImplementation(&mfm.Reference{}) // TODO set reference to wallet class
	walletSender := interfaceSender.(Wallet)
	walletReceiver := w.MockGetSelfReference()
	allowance := walletSender.Allocate(amount, walletReceiver)
	w.balance += allowance.TakeAmount()
}

func (w *Wallet) GetTotalBalance() uint {
	children := allowance.ProxyGetChildrenOf(&mfm.Reference{})
	var totalAllowanced uint = 0
	for _, child := range children {
		totalAllowanced += child.GetBalanceForOwner()
	}
	return w.balance + totalAllowanced
}

func (w *Wallet) ReturnAndDeleteExpiriedAllowances() {
	children := allowance.ProxyGetChildrenOf(&mfm.Reference{})
	for _, child := range children {
		w.balance += child.DeleteExpiredAllowance()
	}
}