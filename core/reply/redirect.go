package reply

import (
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/core/message"
)

type GetObjectRedirectReply struct {
	core.Reply

	To *core.RecordRef
	StateID *core.RecordID

	Sign *core.Signature
}


func NewGetObjectRedirectReply(to *core.RecordRef, state *core.RecordID) *GetObjectRedirectReply {
	return &GetObjectRedirectReply{
		To: to,
		StateID: state,
	}
}

func (r *GetObjectRedirectReply) RecreateMessage(genericMessage core.Message) core.Message {
	getObjectRequest := genericMessage.(*message.GetObject)
	getObjectRequest.State = r.StateID
	return getObjectRequest
}

func (r *GetObjectRedirectReply) CreateToken(genericMessage core.Parcel) []byte {
	newMessage := r.RecreateMessage(genericMessage.Message())
	dataForSign := append(genericMessage.GetSender().Bytes(), message.ToBytes(newMessage)...)
	return dataForSign
}

