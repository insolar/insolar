// Copyright 2020 Insolar Network Ltd.
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

package proc_test

import (
	"context"
	"testing"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"

	"github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/light/proc"
	"github.com/insolar/insolar/ledger/object"
)

func TestPassState_Proceed(t *testing.T) {
	ctx := inslogger.TestContext(t)
	mc := minimock.NewController(t)

	var (
		records *object.RecordAccessorMock
		sender  *bus.SenderMock
	)

	setup := func() {
		records = object.NewRecordAccessorMock(mc)
		sender = bus.NewSenderMock(mc)
	}

	t.Run("Simple success", func(t *testing.T) {
		setup()
		defer mc.Finish()

		stateID := gen.ID()
		origMsg := payload.Meta{}
		origMsgBuf, _ := (&origMsg).Marshal()
		passed, _ := (&payload.PassState{
			Origin:  origMsgBuf,
			StateID: stateID,
		}).Marshal()

		msg := payload.Meta{
			Payload: passed,
		}

		rec := record.Material{
			Virtual:  record.Wrap(&record.Activate{}),
			ID:       gen.ID(),
			ObjectID: gen.ID(),
		}
		records.ForIDMock.Return(rec, nil)

		buf, err := rec.Marshal()
		expectedMsg, _ := payload.NewMessage(&payload.State{
			Record: buf,
		})

		sender.ReplyMock.Inspect(func(ctx context.Context, origin payload.Meta, reply *message.Message) {
			assert.Equal(t, expectedMsg.Payload, reply.Payload)
			assert.Equal(t, origMsg, origin)
		}).Return()

		p := proc.NewPassState(msg, stateID, origMsg)
		p.Dep(records, sender)

		err = p.Proceed(ctx)
		assert.NoError(t, err)
	})

	t.Run("Object not found sends error to origin and last sender", func(t *testing.T) {
		setup()
		defer mc.Finish()

		stateID := gen.ID()
		origMsg := payload.Meta{
			Receiver: gen.Reference(),
		}
		origMsgBuf, _ := (&origMsg).Marshal()
		passed, _ := (&payload.PassState{
			Origin:  origMsgBuf,
			StateID: gen.ID(),
		}).Marshal()

		msg := payload.Meta{
			Payload: passed,
		}

		records.ForIDMock.Return(record.Material{}, object.ErrNotFound)

		expectedError, _ := payload.NewMessage(&payload.Error{
			Text: "state not found",
			Code: payload.CodeNotFound,
		})
		sender.ReplyMock.Inspect(func(ctx context.Context, origin payload.Meta, reply *message.Message) {
			assert.Equal(t, expectedError.Payload, reply.Payload)
			assert.Equal(t, origMsg, origin)
		}).Return()

		p := proc.NewPassState(msg, stateID, origMsg)
		p.Dep(records, sender)

		err := p.Proceed(ctx)
		assert.Error(t, err)
	})

	t.Run("Deactivated object sends error to origin and last sender", func(t *testing.T) {
		setup()
		defer mc.Finish()

		stateID := gen.ID()
		origMsg := payload.Meta{
			Receiver: gen.Reference(),
		}
		origMsgBuf, _ := (&origMsg).Marshal()
		passed, _ := (&payload.PassState{
			Origin:  origMsgBuf,
			StateID: gen.ID(),
		}).Marshal()

		msg := payload.Meta{
			Payload: passed,
		}

		rec := record.Material{
			Virtual:  record.Wrap(&record.Deactivate{}),
			ID:       gen.ID(),
			ObjectID: gen.ID(),
		}

		records.ForIDMock.Return(rec, nil)

		expectedError, _ := payload.NewMessage(&payload.Error{
			Text: "object is deactivated",
			Code: payload.CodeDeactivated,
		})
		sender.ReplyMock.Inspect(func(ctx context.Context, origin payload.Meta, reply *message.Message) {
			assert.Equal(t, expectedError.Payload, reply.Payload)
			assert.Equal(t, origMsg, origin)
		}).Return()

		p := proc.NewPassState(msg, stateID, origMsg)
		p.Dep(records, sender)

		err := p.Proceed(ctx)
		assert.Error(t, err)
	})
}
