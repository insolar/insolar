//
// Copyright 2019 Insolar Technologies GmbH
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package reply

import (
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/message"
)

// GetChildrenRedirectReply is a redirect reply for get children.
type GetChildrenRedirectReply struct {
	Receiver *insolar.Reference
	Token    insolar.DelegationToken

	FromChild insolar.ID
}

// NewGetChildrenRedirect creates a new instance of GetChildrenRedirectReply.
func NewGetChildrenRedirect(
	factory insolar.DelegationTokenFactory, parcel insolar.Parcel, receiver *insolar.Reference, fromChild insolar.ID,
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
func (r *GetChildrenRedirectReply) GetReceiver() *insolar.Reference {
	return r.Receiver
}

// GetToken returns delegation token.
func (r *GetChildrenRedirectReply) GetToken() insolar.DelegationToken {
	return r.Token
}

// Type returns type of the reply
func (r *GetChildrenRedirectReply) Type() insolar.ReplyType {
	return TypeGetChildrenRedirect
}

// Redirected creates redirected message from redirect data.
func (r *GetChildrenRedirectReply) Redirected(genericMsg insolar.Message) insolar.Message {
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
	Receiver *insolar.Reference
	Token    insolar.DelegationToken
}

// GetReceiver returns node reference to send message to.
func (r *GetCodeRedirectReply) GetReceiver() *insolar.Reference {
	return r.Receiver
}

// GetToken returns delegation token.
func (r *GetCodeRedirectReply) GetToken() insolar.DelegationToken {
	return r.Token
}

// Type returns type of the reply
func (r *GetCodeRedirectReply) Type() insolar.ReplyType {
	return TypeGetCodeRedirect
}

// Redirected creates redirected message from redirect data.
func (r *GetCodeRedirectReply) Redirected(genericMsg insolar.Message) insolar.Message {
	return genericMsg
}
