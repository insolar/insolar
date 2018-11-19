package reply

import (
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/core/message"
)

type ObjectRedirectReply struct {
	core.Reply

	To *core.RecordRef
	StateID *core.RecordID
}

func (r *ObjectRedirectReply) RecreateMessage(genericMessage core.Message) core.Message {
	getObjectRequest := genericMessage.(*message.GetObject)
	getObjectRequest.State = r.StateID
	return getObjectRequest
}

func NewObjectRedirect(to *core.RecordRef, state *core.RecordID) *ObjectRedirectReply {
	return &ObjectRedirectReply{
		To: to,
		StateID: state,
	}
}
