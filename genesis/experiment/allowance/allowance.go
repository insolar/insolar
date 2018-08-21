package allowance

import(
	mfm "ilya/v2/mockMagic"
)

type Allowance struct {
	mfm.MockMagic
	to *mfm.Reference
	amount uint
	expTime uint
}

func ProxyGetChildrenOf(reference *mfm.Reference) []*Allowance {
	return make([]*Allowance, 3)
}

func ( a *Allowance ) IsExpired() bool{
	var currTime uint = 0 // TODO: Get real time
	return currTime > a.expTime
}

func (a *Allowance) TakeAmount() uint {
	if a.MockGetCaller() == a.to && !a.IsExpired() {
		a.MockSelfDestructRequest()
		return a.amount
	}
	return 0
}

func (a *Allowance) GetBalanceForOwner() (uint) {
	if !a.IsExpired() {
		return a.amount
	}
	return 0
}

func (a *Allowance) DeleteExpiredAllowance() (uint) {
	if a.MockGetCaller() == a.MockGetMyOwner() && !a.IsExpired() {
		a.MockSelfDestructRequest()
		return a.amount
	}
	return 0
}

