package reply

import (
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/core/message"
)

type Redirect interface {
	core.Reply
	GetTo() *core.RecordRef
	GetSign() *core.Signature
	SetSign(sign *core.Signature)
	RecreateMessage(genericMessage core.Message) core.Message
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

type ObjectRedirect struct {
	GenericRedirect
	StateID core.RecordID
}

func (r *ObjectRedirect) SetSign(sign *core.Signature) {
	r.Sign = sign
}

func (r *ObjectRedirect) RecreateMessage(genericMessage core.Message) core.Message {
	getObjectRequest := genericMessage.(*message.GetObject)
	getObjectRequest.State = &r.StateID
	return getObjectRequest
}

func (r *ObjectRedirect) Type() core.ReplyType {
	return TypeGetObjectRedirect
}

func NewObjectRedirect(to *core.RecordRef, state *core.RecordID) *ObjectRedirect {
	return &ObjectRedirect{
		GenericRedirect: GenericRedirect{
			To: to,
		},
		StateID: *state,
	}
}
