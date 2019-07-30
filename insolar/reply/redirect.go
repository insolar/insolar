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
)

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
