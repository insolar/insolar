/*
 *    Copyright 2019 Insolar
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

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

// GetChildrenRedirectReply is a redirect reply for get children.
type GetChildrenRedirectReply struct {
	Receiver *core.RecordRef
	Token    core.DelegationToken

	FromChild core.RecordID
}

// NewGetChildrenRedirect creates a new instance of GetChildrenRedirectReply.
func NewGetChildrenRedirect(
	factory core.DelegationTokenFactory, parcel core.Parcel, receiver *core.RecordRef, fromChild core.RecordID,
) (*GetChildrenRedirectReply, error) {
	var err error
	rep := GetChildrenRedirectReply{
		Receiver:  receiver,
		FromChild: fromChild,
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
func (r *GetChildrenRedirectReply) GetReceiver() *core.RecordRef {
	return r.Receiver
}

// GetToken returns delegation token.
func (r *GetChildrenRedirectReply) GetToken() core.DelegationToken {
	return r.Token
}

// Type returns type of the reply
func (r *GetChildrenRedirectReply) Type() core.ReplyType {
	return TypeGetChildrenRedirect
}

// Redirected creates redirected message from redirect data.
func (r *GetChildrenRedirectReply) Redirected(genericMsg core.Message) core.Message {
	msg := genericMsg.(*message.GetChildren)
	return &message.GetChildren{
		Parent:    msg.Parent,
		FromChild: &r.FromChild,
		FromPulse: msg.FromPulse,
		Amount:    msg.Amount,
	}
}

// GetCodeRedirectReply is a redirect reply for get children.
type GetCodeRedirectReply struct {
	Receiver *core.RecordRef
	Token    core.DelegationToken
}

// NewGetCodeRedirect creates a new instance of GetChildrenRedirectReply.
func NewGetCodeRedirect(
	factory core.DelegationTokenFactory, parcel core.Parcel, receiver *core.RecordRef,
) (*GetCodeRedirectReply, error) {
	var err error
	rep := GetCodeRedirectReply{
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
func (r *GetCodeRedirectReply) GetReceiver() *core.RecordRef {
	return r.Receiver
}

// GetToken returns delegation token.
func (r *GetCodeRedirectReply) GetToken() core.DelegationToken {
	return r.Token
}

// Type returns type of the reply
func (r *GetCodeRedirectReply) Type() core.ReplyType {
	return TypeGetCodeRedirect
}

// Redirected creates redirected message from redirect data.
func (r *GetCodeRedirectReply) Redirected(genericMsg core.Message) core.Message {
	return genericMsg
}
