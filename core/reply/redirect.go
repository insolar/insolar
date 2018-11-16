package reply

import (
	"github.com/insolar/insolar/core"
)

type Redirect interface {
	core.Reply
	GetTo() *core.RecordRef
	GetSign() *core.Signature
}

type GenericRedirect struct {
	To   *core.RecordRef
	Sign *core.Signature
}

func (r *GenericRedirect) GetTo() *core.RecordRef {
	return r.To
}

func (r *GenericRedirect) GetSign() *core.Signature {
	return r.Sign
}

func (r *GenericRedirect) Type() core.ReplyType {
	return TypeRedirect
}

type ObjectRedirect struct {
	GenericRedirect
	StateID core.RecordID
}

func (r *ObjectRedirect) Type() core.ReplyType {
	return TypeDefinedStateRedirect
}
