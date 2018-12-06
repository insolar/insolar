/*
 *    Copyright 2018 Insolar
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

package delegationtoken

import (
	"bytes"
	"encoding/gob"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/core/message"
)

type delegationTokenFactory struct {
	Cryptography core.CryptographyService `inject:""`
}

// NewDelegationTokenFactory creates new token factory instance.
func NewDelegationTokenFactory() core.DelegationTokenFactory {
	return &delegationTokenFactory{}
}

// IssuePendingExecution creates new token for provided message.
func (f *delegationTokenFactory) IssuePendingExecution(
	msg core.Message, pulse core.PulseNumber,
) (core.DelegationToken, error) {
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
	token := &PendingExecution{}
	token.Signature = sign.Bytes()

	return token, nil
}

// IssueGetObjectRedirect creates new token for provided message.
func (f *delegationTokenFactory) IssueGetObjectRedirect(
	sender *core.RecordRef, redirectedMessage core.Message,
) (core.DelegationToken, error) {
	parsedMessage := redirectedMessage.(*message.GetObject)
	dataForSign := append(sender.Bytes(), message.ToBytes(parsedMessage)...)
	sign, err := f.Cryptography.Sign(dataForSign)
	if err != nil {
		return nil, err
	}
	return &GetObjectRedirect{Signature: sign.Bytes()}, nil
}

// IssueGetChildrenRedirect creates new token for provided message.
func (f *delegationTokenFactory) IssueGetChildrenRedirect(
	sender *core.RecordRef, redirectedMessage core.Message,
) (core.DelegationToken, error) {
	parsedMessage := redirectedMessage.(*message.GetChildren)
	dataForSign := append(sender.Bytes(), message.ToBytes(parsedMessage)...)
	sign, err := f.Cryptography.Sign(dataForSign)
	if err != nil {
		return nil, err
	}
	return &GetChildrenRedirect{Signature: sign.Bytes()}, nil
}

// IssueGetCodeRedirect creates new token for provided message.
func (f *delegationTokenFactory) IssueGetCodeRedirect(
	sender *core.RecordRef, redirectedMessage core.Message,
) (core.DelegationToken, error) {
	parsedMessage := redirectedMessage.(*message.GetCode)
	dataForSign := append(sender.Bytes(), message.ToBytes(parsedMessage)...)
	sign, err := f.Cryptography.Sign(dataForSign)
	if err != nil {
		return nil, err
	}
	return &GetCodeRedirect{Signature: sign.Bytes()}, nil
}

// Verify performs token validation.
func (f *delegationTokenFactory) Verify(parcel core.Parcel) (bool, error) {
	if parcel.DelegationToken() == nil {
		return false, nil
	}

	return parcel.DelegationToken().Verify(parcel)
}
