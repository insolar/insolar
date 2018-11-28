package reply

import (
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/core/message"
)

// GetObjectRedirect is a redirect-reply for get object
type GetObjectRedirect struct {
	Receiver *core.RecordRef
	Token    core.DelegationToken

	StateID *core.RecordID
}

// NewGetObjectRedirectReply return new GetObjectRedirect
func NewGetObjectRedirectReply(
	factory core.DelegationTokenFactory, parcel core.Parcel, receiver *core.RecordRef, state *core.RecordID,
) (*GetObjectRedirect, error) {
	var err error
	rep := GetObjectRedirect{
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
func (r *GetObjectRedirect) GetReceiver() *core.RecordRef {
	return r.Receiver
}

// GetToken returns delegation token.
func (r *GetObjectRedirect) GetToken() core.DelegationToken {
	return r.Token
}

// Type returns type of the reply
func (r *GetObjectRedirect) Type() core.ReplyType {
	return TypeGetObjectRedirect
}

// Redirected creates redirected message from redirect data.
func (r *GetObjectRedirect) Redirected(genericMsg core.Message) core.Message {
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

// GetCodeRedirect is a redirect reply for get children.
type GetCodeRedirect struct {
	Receiver *core.RecordRef
	Token    core.DelegationToken
}

// NewGetCodeRedirect creates a new instance of GetChildrenRedirect.
func NewGetCodeRedirect(
	factory core.DelegationTokenFactory, parcel core.Parcel, receiver *core.RecordRef,
) (*GetCodeRedirect, error) {
	var err error
	rep := GetCodeRedirect{
		Receiver: receiver,
	}
	redirectedMessage := rep.Redirected(parcel.Message())
	sender := parcel.GetSender()
	rep.Token, err = factory.IssueGetCodeRedirect(&sender, redirectedMessage)
	if err != nil {
		return nil, err
	}
	return &rep, nil
}

// GetReceiver returns node reference to send message to.
func (r *GetCodeRedirect) GetReceiver() *core.RecordRef {
	return r.Receiver
}

// GetToken returns delegation token.
func (r *GetCodeRedirect) GetToken() core.DelegationToken {
	return r.Token
}

// Type returns type of the reply
func (r *GetCodeRedirect) Type() core.ReplyType {
	return TypeGetCodeRedirect
}

// Redirected creates redirected message from redirect data.
func (r *GetCodeRedirect) Redirected(genericMsg core.Message) core.Message {
	return genericMsg
}
