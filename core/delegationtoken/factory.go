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
	"github.com/pkg/errors"
)

type delegationTokenFactory struct {
	Cryptography core.CryptographyService `inject:""`
}

func NewDelegationTokenFactory() core.DelegationTokenFactory {
	return &delegationTokenFactory{}
}

func (f *delegationTokenFactory) IssuePendingExecution(
	msg core.Message, pulse core.PulseNumber,
) (core.DelegationToken, error) {
	var buff bytes.Buffer
	enc := gob.NewEncoder(&buff)
	err := enc.Encode(msg)
	if err != nil {
		return []byte{}, err
	}

	sign, err := f.Cryptography.Sign(buff.Bytes())
	if err != nil {
		return []byte{}, err
	}

	return append([]byte{byte(core.DTTypePendingExecution)}, sign.Bytes()...), nil
}

func (f *delegationTokenFactory) Verify(token core.DelegationToken, msg core.Message) (bool, error) {
	token, err := f.newFromBytes(data)
	if err != nil {
		return false, err
	}
	if token == nil {
		return false, nil
	}

	return token.Verify(msg)
}

func (f *delegationTokenFactory) newFromBytes(data []byte) (core.DelegationToken, error) {
	if len(data) == 0 {
		return nil, nil
	}

	res, err := empty(core.DelegationTokenType(data[0]))
	if err != nil {
		return nil, err
	}

	err = gob.NewDecoder(bytes.NewReader(data[1:])).Decode(res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func empty(t core.DelegationTokenType) (core.DelegationToken, error) {
	switch t {

	case core.DTTypePendingExecution:
		return &PendingExecution{}, nil
	default:
		return nil, errors.Errorf("unimplemented delegation token type %d", t)
	}
}
