package reply

import (
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/core/message"
)

// GetObjectRedirectReply is a redirect-reply for get object
type GetObjectRedirectReply struct {
	Receiver *core.RecordRef
	Token    core.DelegationToken

	StateID *core.RecordID
}

// GetReceiver returns node reference to send message to.
func (r *GetObjectRedirectReply) GetReceiver() *core.RecordRef {
	return r.Receiver
}

// GetToken returns delegation token.
func (r *GetObjectRedirectReply) GetToken() core.DelegationToken {
	return r.Token
}

// NewGetObjectRedirectReply return new GetObjectRedirectReply
func NewGetObjectRedirectReply(to *core.RecordRef, state *core.RecordID) *GetObjectRedirectReply {
	return &GetObjectRedirectReply{
		Receiver: to,
		StateID:  state,
	}
}

// Type returns type of the reply
func (r *GetObjectRedirectReply) Type() core.ReplyType {
	return TypeGetObjectRedirect
}

// Redirected creates redirected message from redirect data.
func (r *GetObjectRedirectReply) Redirected(genericMsg core.Message) core.Message {
	msg := genericMsg.(*message.GetObject)
	return &message.GetObject{
		State:    r.StateID,
		Head:     msg.Head,
		Approved: msg.Approved,
	}
}
