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

// NewGetObjectRedirectReply return new GetObjectRedirectReply
func NewGetObjectRedirectReply(
	factory core.DelegationTokenFactory, parcel core.Parcel, receiver *core.RecordRef, state *core.RecordID,
) (*GetObjectRedirectReply, error) {
	var err error
	rep := GetObjectRedirectReply{
		Receiver: receiver,
		StateID:  state,
	}
	redirectedMessage := rep.Redirected(parcel.Message())
	sender := parcel.GetSender()
	rep.Token, err = factory.IssueGetObjectRedirect(&sender, redirectedMessage)
	if err != nil {
		return nil, err
	}
	return &rep, nil
}

// GetReceiver returns node reference to send message to.
func (r *GetObjectRedirectReply) GetReceiver() *core.RecordRef {
	return r.Receiver
}

// GetToken returns delegation token.
func (r *GetObjectRedirectReply) GetToken() core.DelegationToken {
	return r.Token
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

// GetChildrenRedirect is a redirect reply for get children.
type GetChildrenRedirect struct {
	Receiver *core.RecordRef
	Token    core.DelegationToken
}

// NewGetChildrenRedirect creates a new instance of GetChildrenRedirect.
func NewGetChildrenRedirect(
	factory core.DelegationTokenFactory, parcel core.Parcel, receiver *core.RecordRef,
) (*GetChildrenRedirect, error) {
	var err error
	rep := GetChildrenRedirect{
		Receiver: receiver,
	}
	redirectedMessage := rep.Redirected(parcel.Message())
	sender := parcel.GetSender()
	rep.Token, err = factory.IssueGetChildrenRedirect(&sender, redirectedMessage)
	if err != nil {
		return nil, err
	}
	return &rep, nil
}

// GetReceiver returns node reference to send message to.
func (r *GetChildrenRedirect) GetReceiver() *core.RecordRef {
	return r.Receiver
}

// GetToken returns delegation token.
func (r *GetChildrenRedirect) GetToken() core.DelegationToken {
	return r.Token
}

// Type returns type of the reply
func (r *GetChildrenRedirect) Type() core.ReplyType {
	return TypeGetChildrenRedirect
}

// Redirected creates redirected message from redirect data.
func (r *GetChildrenRedirect) Redirected(genericMsg core.Message) core.Message {
	return genericMsg
}
