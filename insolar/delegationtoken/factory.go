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

package delegationtoken

import (
	"bytes"
	"encoding/gob"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/message"
)

type delegationTokenFactory struct {
	Cryptography insolar.CryptographyService `inject:""`
}

// NewDelegationTokenFactory creates new token factory instance.
func NewDelegationTokenFactory() insolar.DelegationTokenFactory {
	return &delegationTokenFactory{}
}

// IssuePendingExecution creates new token for provided message.
func (f *delegationTokenFactory) IssuePendingExecution(
	msg insolar.Message, pulse insolar.PulseNumber,
) (insolar.DelegationToken, error) {
	var buff bytes.Buffer
	enc := gob.NewEncoder(&buff)
	err := enc.Encode(msg)
	if err != nil {
		return nil, err
	}

	sign, err := f.Cryptography.Sign(buff.Bytes())
	if err != nil {
		return nil, err
	}
	token := &PendingExecutionToken{}
	token.Signature = sign.Bytes()

	return token, nil
}

// IssueGetChildrenRedirect creates new token for provided message.
func (f *delegationTokenFactory) IssueGetChildrenRedirect(
	sender *insolar.Reference, redirectedMessage insolar.Message,
) (insolar.DelegationToken, error) {
	parsedMessage := redirectedMessage.(*message.GetChildren)
	dataForSign := append(sender.Bytes(), message.ToBytes(parsedMessage)...)
	sign, err := f.Cryptography.Sign(dataForSign)
	if err != nil {
		return nil, err
	}
	return &GetChildrenRedirectToken{Signature: sign.Bytes()}, nil
}

// Verify performs token validation.
func (f *delegationTokenFactory) Verify(parcel insolar.Parcel) (bool, error) {
	if parcel.DelegationToken() == nil {
		return false, nil
	}

	return parcel.DelegationToken().Verify(parcel)
}
