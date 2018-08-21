package memberProxy

import (
	mfm "ilya/v2/mockMagic"
	"ilya/v2/member"
	"ilya/v2/wallet"
)

type MemberProxy struct {
	member.Member
}

func (mp *MemberProxy) ProxyGetImplementation(ref *mfm.Reference) interface{} {
	// TODO magic
	return wallet.Wallet{}
}

func ProxyGetObject(addressRef *mfm.Reference) *MemberProxy {
	// TODO get object from ledger
	return &MemberProxy{}
}

func (m *MemberProxy) GetName() string {
	return m.Member.GetName()
}

func (m *MemberProxy) GetPublicKey() []byte {
	return m.Member.GetPublicKey()
}