/*
 *
 *  *    Copyright 2018 Insolar
 *  *
 *  *    Licensed under the Apache License, Version 2.0 (the "License");
 *  *    you may not use this file except in compliance with the License.
 *  *    You may obtain a copy of the License at
 *  *
 *  *        http://www.apache.org/licenses/LICENSE-2.0
 *  *
 *  *    Unless required by applicable law or agreed to in writing, software
 *  *    distributed under the License is distributed on an "AS IS" BASIS,
 *  *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  *    See the License for the specific language governing permissions and
 *  *    limitations under the License.
 *
 */

package testmessagebus

import (
	"context"
	"encoding/gob"
	"fmt"
	"io"
	"reflect"
	"testing"

	"github.com/pkg/errors"

	"github.com/insolar/insolar/component"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/core/delegationtoken"
	"github.com/insolar/insolar/core/message"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/messagebus"
	"github.com/insolar/insolar/platformpolicy"
	"github.com/insolar/insolar/testutils"
)

type TapeRecord struct {
	Message core.Message
	Reply   core.Reply
	Error   error
}

type TestMessageBus struct {
	handlers    map[core.MessageType]core.MessageHandler
	pf          message.ParcelFactory
	PulseNumber core.PulseNumber
	ReadingTape []TapeRecord
	WritingTape []TapeRecord
}

func (mb *TestMessageBus) NewPlayer(ctx context.Context, reader io.Reader) (core.MessageBus, error) {
	tape := make([]TapeRecord, 0)
	enc := gob.NewDecoder(reader)
	err := enc.Decode(&tape)
	if err != nil {
		return nil, err
	}
	res := *mb
	res.ReadingTape = tape
	return &res, nil
}

func (mb *TestMessageBus) WriteTape(ctx context.Context, writer io.Writer) error {
	if mb.WritingTape == nil {
		return errors.New("Not writing message bus")
	}
	enc := gob.NewEncoder(writer)
	err := enc.Encode(mb.WritingTape)
	if err != nil {
		return err
	}

	return nil
}

func (mb *TestMessageBus) NewRecorder(ctx context.Context, currentPulse core.Pulse) (core.MessageBus, error) {
	tape := make([]TapeRecord, 0)
	res := *mb
	res.WritingTape = tape
	return &res, nil
}

func NewTestMessageBus(t *testing.T) *TestMessageBus {
	mock := testutils.NewCryptographyServiceMock(t)
	mock.SignFunc = func(p []byte) (r *core.Signature, r1 error) {
		signature := core.SignatureFromBytes(nil)
		return &signature, nil
	}
	delegationTokenFactory := delegationtoken.NewDelegationTokenFactory()
	parcelFactory := messagebus.NewParcelFactory()
	cm := &component.Manager{}
	cm.Register(platformpolicy.NewPlatformCryptographyScheme())
	cm.Inject(delegationTokenFactory, parcelFactory, mock)

	return &TestMessageBus{handlers: map[core.MessageType]core.MessageHandler{}, pf: parcelFactory}
}

func (mb *TestMessageBus) Register(p core.MessageType, handler core.MessageHandler) error {
	_, ok := mb.handlers[p]
	if ok {
		return errors.New("handler for this type already exists")
	}

	mb.handlers[p] = handler
	return nil
}

func (mb *TestMessageBus) ReRegister(p core.MessageType, handler core.MessageHandler) {
	mb.handlers[p] = handler
}

func (mb *TestMessageBus) MustRegister(p core.MessageType, handler core.MessageHandler) {
	err := mb.Register(p, handler)
	if err != nil {
		panic(err)
	}
}

func (mb *TestMessageBus) Send(
	ctx context.Context, m core.Message, currentPulse core.Pulse, ops *core.MessageSendOptions,
) (core.Reply, error) {
	if mb.ReadingTape != nil {
		if len(mb.ReadingTape) == 0 {
			return nil, errors.Errorf("No expected messages, got %+v", m)
		}
		head, tail := mb.ReadingTape[0], mb.ReadingTape[1:]
		mb.ReadingTape = tail

		inslogger.FromContext(ctx).Debugf("Reading message %+v off the tape", head.Message)

		if !reflect.DeepEqual(head.Message, m) {
			return nil, errors.Errorf("Message in the tape and sended arn't equal; got: %+v, expected: %+v", m, head.Message)
		}
		return head.Reply, head.Error
	}
	parcel, err := mb.pf.Create(ctx, m, testutils.RandomRef(), nil, currentPulse)
	if err != nil {
		return nil, err
	}
	t := parcel.Message().Type()
	handler, ok := mb.handlers[t]
	if !ok {
		return nil, errors.New(fmt.Sprint("no handler for message type:", t.String()))
	}

	ctx = parcel.Context(context.Background())

	reply, err := handler(ctx, parcel)
	if mb.WritingTape != nil {
		inslogger.FromContext(ctx).Debugf("Writing message %+v on the tape", m)
		mb.WritingTape = append(mb.WritingTape, TapeRecord{Message: m, Reply: reply, Error: err})
	}

	return reply, err
}
