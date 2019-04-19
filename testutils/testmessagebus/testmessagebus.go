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

package testmessagebus

import (
	"context"
	"encoding/gob"
	"fmt"
	"io"
	"reflect"
	"testing"

	"github.com/insolar/insolar/ledger/storage/pulse"
	"github.com/pkg/errors"

	"github.com/insolar/insolar/component"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/delegationtoken"
	"github.com/insolar/insolar/insolar/message"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/messagebus"
	"github.com/insolar/insolar/platformpolicy"
	"github.com/insolar/insolar/testutils"
)

type TapeRecord struct {
	Message insolar.Message
	Reply   insolar.Reply
	Error   error
}

type TestMessageBus struct {
	handlers      map[insolar.MessageType]insolar.MessageHandler
	pf            message.ParcelFactory
	PulseAccessor pulse.Accessor
	ReadingTape   []TapeRecord
	WritingTape   []TapeRecord
}

func (mb *TestMessageBus) NewPlayer(ctx context.Context, reader io.Reader) (insolar.MessageBus, error) {
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

func (mb *TestMessageBus) NewRecorder(ctx context.Context, currentPulse insolar.Pulse) (insolar.MessageBus, error) {
	tape := make([]TapeRecord, 0)
	res := *mb
	res.WritingTape = tape
	return &res, nil
}

func NewTestMessageBus(t *testing.T) *TestMessageBus {
	cryptoServiceMock := testutils.NewCryptographyServiceMock(t)
	cryptoServiceMock.SignFunc = func(p []byte) (r *insolar.Signature, r1 error) {
		signature := insolar.SignatureFromBytes(nil)
		return &signature, nil
	}

	delegationTokenFactory := delegationtoken.NewDelegationTokenFactory()

	parcelFactory := messagebus.NewParcelFactory()

	cm := &component.Manager{}
	cm.Register(platformpolicy.NewPlatformCryptographyScheme())
	cm.Inject(delegationTokenFactory, parcelFactory, cryptoServiceMock)

	return &TestMessageBus{handlers: map[insolar.MessageType]insolar.MessageHandler{}, pf: parcelFactory}
}

func (mb *TestMessageBus) Register(p insolar.MessageType, handler insolar.MessageHandler) error {
	_, ok := mb.handlers[p]
	if ok {
		return errors.New("handler for this type already exists")
	}

	mb.handlers[p] = handler
	return nil
}

func (mb *TestMessageBus) ReRegister(p insolar.MessageType, handler insolar.MessageHandler) {
	mb.handlers[p] = handler
}

func (mb *TestMessageBus) MustRegister(p insolar.MessageType, handler insolar.MessageHandler) {
	err := mb.Register(p, handler)
	if err != nil {
		panic(err)
	}
}

func (mb *TestMessageBus) Send(ctx context.Context, m insolar.Message, _ *insolar.MessageSendOptions) (insolar.Reply, error) {
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

	currentPulse, err := mb.PulseAccessor.Latest(ctx)
	if err != nil {
		return nil, err
	}

	parcel, err := mb.pf.Create(ctx, m, testutils.RandomRef(), nil, insolar.Pulse{PulseNumber: currentPulse.PulseNumber, Entropy: insolar.Entropy{}})
	if err != nil {
		return nil, err
	}
	t := parcel.Message().Type()
	handler, ok := mb.handlers[t]
	if !ok {
		return nil, errors.New(fmt.Sprint("[ TestMessageBus ] no handler for message type:", t.String()))
	}

	ctx = parcel.Context(context.Background())

	reply, err := handler(ctx, parcel)
	if mb.WritingTape != nil {
		// WARNING! The following commented line of code is cursed.
		// It makes some test (e.g. TestNilResults) hang under the debugger, and we have no idea why.
		// Don't uncomment unless you solved this mystery.
		// inslogger.FromContext(ctx).Debugf("Writing message %+v on the tape", m)
		mb.WritingTape = append(mb.WritingTape, TapeRecord{Message: m, Reply: reply, Error: err})
	}

	return reply, err
}

func (mb *TestMessageBus) SendViaWatermill(ctx context.Context, m insolar.Message, _ *insolar.MessageSendOptions) (insolar.Reply, error) {
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

	currentPulse, err := mb.PulseAccessor.Latest(ctx)
	if err != nil {
		return nil, err
	}

	parcel, err := mb.pf.Create(ctx, m, testutils.RandomRef(), nil, insolar.Pulse{PulseNumber: currentPulse.PulseNumber, Entropy: insolar.Entropy{}})
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
		// WARNING! The following commented line of code is cursed.
		// It makes some test (e.g. TestNilResults) hang under the debugger, and we have no idea why.
		// Don't uncomment unless you solved this mystery.
		// inslogger.FromContext(ctx).Debugf("Writing message %+v on the tape", m)
		mb.WritingTape = append(mb.WritingTape, TapeRecord{Message: m, Reply: reply, Error: err})
	}

	return reply, err
}

func (mb *TestMessageBus) OnPulse(context.Context, insolar.Pulse) error {
	return nil
}

func (mb *TestMessageBus) Lock(context.Context) {
}

func (mb *TestMessageBus) Unlock(context.Context) {
}
