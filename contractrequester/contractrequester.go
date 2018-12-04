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

package contractrequester

import (
	"context"
	"crypto/rand"
	"encoding/binary"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/core/message"
	"github.com/pkg/errors"
)

// ContractRequester helps to call contracts
type ContractRequester struct {
	MessageBus core.MessageBus `inject:""`
}

// New creates new ContractRequester
func New() (*ContractRequester, error) {
	return &ContractRequester{}, nil
}

func randomUint64() uint64 {
	buf := make([]byte, 8)
	_, err := rand.Read(buf)
	if err != nil {
		panic(err)
	}

	return binary.LittleEndian.Uint64(buf)
}

func (cr *ContractRequester) routeCall(ctx context.Context, ref core.RecordRef, method string, args core.Arguments) (core.Reply, error) {
	if cr.MessageBus == nil {
		return nil, errors.New("[ ContractRequester::routeCall ] message bus was not set during initialization")
	}

	e := &message.CallMethod{
		BaseLogicMessage: message.BaseLogicMessage{Nonce: randomUint64()},
		ObjectRef:        ref,
		Method:           method,
		Arguments:        args,
	}

	res, err := cr.MessageBus.Send(ctx, e, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "[ ContractRequester::routeCall ] couldn't send message: %s", ref.String())
	}

	return res, nil
}

// SendRequest makes call to method of contract by its ref
func (cr *ContractRequester) SendRequest(ctx context.Context, ref *core.RecordRef, method string, argsIn []interface{}) (core.Reply, error) {
	args, err := core.MarshalArgs(argsIn...)
	if err != nil {
		return nil, errors.Wrap(err, "[ ContractRequester::SendRequest ] Can't marshal")
	}

	routResult, err := cr.routeCall(ctx, *ref, method, args)
	if err != nil {
		return nil, errors.Wrap(err, "[ ContractRequester::SendRequest ] Can't route call")
	}

	return routResult, nil
}
