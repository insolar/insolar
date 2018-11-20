package reply

import (
	"github.com/insolar/insolar/core"
)

type GetObjectRedirectReply struct {
	core.Reply

	To *core.RecordRef
	StateID *core.RecordID

	Token *core.DelegationToken
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

