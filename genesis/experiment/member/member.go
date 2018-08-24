package member

import (
	"github.com/insolar/insolar/toolkit/go/foundation"
)

//var TypeReference = foundation.Reference("member")

type Member struct {
	foundation.BaseContract
	Name      string
	PublicKey []byte
}

func (m *Member) GetName() string {
	return m.Name
}
func (m *Member) GetPublicKey() []byte {
	return m.PublicKey
}

func NewMember(name string) (*Member, foundation.Reference) {
	member := &Member{
		Name: name,
	}
	//fmt.Printf("%x\n", &member)
	//fmt.Printf("%x\n", &(member.BaseContract))
	//fmt.Println(member.MyReference())
	reference := foundation.SaveToLedger(member)
	return member, reference
}
