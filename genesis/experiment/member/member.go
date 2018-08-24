package member

import (
	"github.com/insolar/insolar/toolkit/go/foundation"
)

var TypeReference = foundation.Reference("member")

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
	reference := foundation.SaveToLedger(member)
	return member, reference
}
