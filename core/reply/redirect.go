package reply

import (
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/core/message"
)

type ObjectRedirect struct {
	core.Reply

	To *core.RecordRef
	StateID *core.RecordID
}

func (r *ObjectRedirect) RecreateMessage(genericMessage core.Message) core.Message {
	getObjectRequest := genericMessage.(*message.GetObject)
	getObjectRequest.State = r.StateID
	return getObjectRequest
}

func NewObjectRedirect(to *core.RecordRef, state *core.RecordID) *ObjectRedirect {
	return &ObjectRedirect{
		To: to,
		StateID: state,
	}
}
