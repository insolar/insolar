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

package proc_test

import (
	"context"
	"testing"
	"time"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/gojuno/minimock"
	"github.com/stretchr/testify/assert"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/light/executor"
	"github.com/insolar/insolar/ledger/light/proc"
	"github.com/insolar/insolar/ledger/object"
)

func TestSendObject_Proceed(t *testing.T) {
	ctx := flow.TestContextWithPulse(inslogger.TestContext(t), gen.PulseNumber())
	mc := minimock.NewController(t)

	var (
		indexes     *object.IndexAccessorMock
		jets        *jet.StorageMock
		records     *object.RecordAccessorMock
		sender      *bus.SenderMock
		coordinator *jet.CoordinatorMock
		fetcher     *executor.JetFetcherMock
	)

	setup := func() {
		records = object.NewRecordAccessorMock(mc)
		sender = bus.NewSenderMock(mc)
		coordinator = jet.NewCoordinatorMock(mc)
		fetcher = executor.NewJetFetcherMock(mc)
		indexes = object.NewIndexAccessorMock(mc)
		jets = jet.NewStorageMock(mc)
	}

	t.Run("Error deactivated object", func(t *testing.T) {
		setup()
		defer mc.Finish()

		objectID := gen.ID()
		latestState := gen.ID()
		indexes.ForIDMock.Return(record.Index{
			ObjID: objectID,
			Lifeline: record.Lifeline{
				LatestState: &latestState,
				StateID:     record.StateDeactivation,
			},
		}, nil)

		msg := payload.Meta{}

		expectedError, _ := payload.NewMessage(&payload.Error{
			Text: "object is deactivated",
			Code: payload.CodeDeactivated,
		})
		sender.ReplyMock.Inspect(func(ctx context.Context, origin payload.Meta, reply *message.Message) {
			assert.Equal(t, expectedError.Payload, reply.Payload)
			assert.Equal(t, msg, origin)
		}).Return()

		p := proc.NewSendObject(msg, objectID)
		p.Dep(coordinator, jets, fetcher, records, indexes, sender)

		err := p.Proceed(ctx)
		assert.NoError(t, err)
	})

	t.Run("Simple success", func(t *testing.T) {
		setup()
		defer mc.Finish()

		objectID := gen.ID()
		latestState := gen.ID()
		index := record.Index{
			ObjID: objectID,
			Lifeline: record.Lifeline{
				LatestState: &latestState,
				StateID:     record.StateActivation,
			},
		}

		indexes.ForIDMock.Return(index, nil)
		buf, err := index.Lifeline.Marshal()
		expectedIndex, _ := payload.NewMessage(&payload.Index{
			Index: buf,
		})

		rec := record.Material{
			Virtual:  record.Wrap(&record.Activate{}),
			ID:       gen.ID(),
			ObjectID: objectID,
		}
		records.ForIDMock.Return(rec, nil)

		buf, _ = rec.Marshal()
		expectedMsg, _ := payload.NewMessage(&payload.State{
			Record: buf,
		})
		msg := payload.Meta{}

		sender.ReplyMock.Inspect(func(ctx context.Context, origin payload.Meta, reply *message.Message) {
			// First message, index
			if sender.ReplyAfterCounter() == 0 {
				assert.Equal(t, expectedIndex.Payload, reply.Payload)
			}

			// Second message, record
			if sender.ReplyAfterCounter() == 1 {
				assert.Equal(t, expectedMsg.Payload, reply.Payload)
			}
			assert.Equal(t, msg, origin)
		}).Return()

		p := proc.NewSendObject(msg, objectID)
		p.Dep(coordinator, jets, fetcher, records, indexes, sender)

		err = p.Proceed(ctx)
		assert.NoError(t, err)

	})

	t.Run("Error reply, Deactivated from State", func(t *testing.T) {
		setup()
		defer mc.Finish()

		objectID := gen.ID()
		latestState := gen.ID()
		index := record.Index{
			ObjID: objectID,
			Lifeline: record.Lifeline{
				LatestState: &latestState,
				StateID:     record.StateActivation,
			},
		}

		indexes.ForIDMock.Return(index, nil)
		buf, err := index.Lifeline.Marshal()
		expectedIndex, _ := payload.NewMessage(&payload.Index{
			Index: buf,
		})

		rec := record.Material{
			Virtual:  record.Wrap(&record.Deactivate{}),
			ID:       gen.ID(),
			ObjectID: objectID,
		}
		records.ForIDMock.Return(rec, nil)

		buf, _ = rec.Marshal()
		expectedError, _ := payload.NewMessage(&payload.Error{
			Text: "object is deactivated",
			Code: payload.CodeDeactivated,
		})

		msg := payload.Meta{}

		sender.ReplyMock.Inspect(func(ctx context.Context, origin payload.Meta, reply *message.Message) {
			// First message, index
			if sender.ReplyAfterCounter() == 0 {
				assert.Equal(t, expectedIndex.Payload, reply.Payload)
			}

			// Second message, record
			if sender.ReplyAfterCounter() == 1 {
				assert.Equal(t, expectedError.Payload, reply.Payload)
			}
			assert.Equal(t, msg, origin)
		}).Return()

		p := proc.NewSendObject(msg, objectID)
		p.Dep(coordinator, jets, fetcher, records, indexes, sender)

		err = p.Proceed(ctx)
		assert.NoError(t, err)
	})

	t.Run("Send PassState on heavy", func(t *testing.T) {
		setup()
		defer mc.Finish()
		defer mc.Wait(10 * time.Second)

		objectID := gen.ID()
		latestState := gen.ID()
		index := record.Index{
			ObjID: objectID,
			Lifeline: record.Lifeline{
				LatestState: &latestState,
				StateID:     record.StateActivation,
			},
		}

		indexes.ForIDMock.Return(index, nil)
		buf, err := index.Lifeline.Marshal()
		expectedIndex, _ := payload.NewMessage(&payload.Index{
			Index: buf,
		})

		records.ForIDMock.Return(record.Material{}, object.ErrNotFound)

		msg := payload.Meta{}
		buf, _ = msg.Marshal()
		expectedMsg, _ := payload.NewMessage(&payload.PassState{
			Origin:  buf,
			StateID: latestState,
		})

		expectedTarget := insolar.NewReference(gen.ID())
		coordinator.IsBeyondLimitMock.Return(true, nil)
		coordinator.HeavyMock.Return(expectedTarget, nil)

		sender.ReplyMock.Inspect(func(ctx context.Context, origin payload.Meta, reply *message.Message) {
			assert.Equal(t, expectedIndex.Payload, reply.Payload)
			assert.Equal(t, msg, origin)
		}).Return()

		sender.SendTargetMock.Inspect(func(ctx context.Context, msg *message.Message, target insolar.Reference) {
			assert.Equal(t, expectedMsg.Payload, msg.Payload)
			assert.Equal(t, expectedTarget, &target)
		}).Return(make(chan *message.Message), func() {})

		p := proc.NewSendObject(msg, objectID)
		p.Dep(coordinator, jets, fetcher, records, indexes, sender)

		err = p.Proceed(ctx)
		assert.NoError(t, err)
	})
}
