package reply

import (
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/core/message"
)

// GetObjectRedirectReply is a redirect-reply for get object
type GetObjectRedirectReply struct {
	core.Reply

	To      *core.RecordRef
	StateID *core.RecordID

	Token core.DelegationToken
}

// NewGetObjectRedirectReply return new GetObjectRedirectReply
func NewGetObjectRedirectReply(to *core.RecordRef, state *core.RecordID) *GetObjectRedirectReply {
	return &GetObjectRedirectReply{
		To:      to,
		StateID: state,
	}
}

// Type returns type of the reply
func (r *GetObjectRedirectReply) Type() core.ReplyType {
	return TypeGetObjectRedirect
}

// RecreateMessage recreates the message on the base of token
func (r *GetObjectRedirectReply) RecreateMessage(msg *message.GetObject) *message.GetObject {
	return &message.GetObject{
		State:    r.StateID,
		Head:     msg.Head,
		Approved: msg.Approved,
	}
}
