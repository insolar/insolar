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
	"github.com/insolar/insolar/pulse"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEnsureIndex_Proceed(t *testing.T) {
	ctx := inslogger.TestContext(t)
	mc := minimock.NewController(t)

	var (
		indexes       *object.MemoryIndexStorageMock
		cord          *jet.CoordinatorMock
		sender        *bus.SenderMock
		writeAccessor *executor.WriteAccessorMock
	)
	setup := func() {
		indexes = object.NewMemoryIndexStorageMock(mc)
		cord = jet.NewCoordinatorMock(mc)
		sender = bus.NewSenderMock(mc)
		writeAccessor = executor.NewWriteAccessorMock(mc)
	}

	t.Run("Simple success", func(t *testing.T) {
		setup()
		defer mc.Finish()

		pulse := gen.PulseNumber()
		idx := record.Index{
			ObjID:            insolar.ID{},
			Lifeline:         record.Lifeline{},
			LifelineLastUsed: pulse,
			PendingRecords:   nil,
		}
		indexes.ForIDMock.Return(idx, nil)

		p := proc.NewEnsureIndex(gen.ID(), gen.JetID(), payload.Meta{}, pulse)
		p.Dep(indexes, cord, sender, writeAccessor)
		err := p.Proceed(ctx)
		assert.NoError(t, err)
	})

	t.Run("fetches from heavy if index not found, returns flow cancelled error", func(t *testing.T) {
		setup()
		defer mc.Finish()

		indexes.ForIDMock.Return(record.Index{}, object.ErrIndexNotFound)
		cord.HeavyMock.Return(&insolar.Reference{}, nil)
		idx, err := (&record.Lifeline{}).Marshal()
		reps := make(chan *message.Message, 1)
		reps <- payload.MustNewMessage(&payload.Meta{
			Payload: payload.MustMarshal(&payload.Index{
				Polymorph: 0,
				Index:     idx,
			}),
		})
		sender.SendTargetMock.Return(reps, func() {})
		writeAccessor.BeginMock.Return(func() {}, executor.ErrWriteClosed)

		p := proc.NewEnsureIndex(gen.ID(), gen.JetID(), payload.Meta{}, pulse.MinTimePulse)
		p.Dep(indexes, cord, sender, writeAccessor)
		err = p.Proceed(ctx)
		assert.Error(t, err)
		assert.Equal(t, err, flow.ErrCancelled)
	})

	t.Run("success, fetches from heavy if index not found", func(t *testing.T) {
		setup()
		defer mc.Finish()

		indexes.ForIDMock.Return(record.Index{}, object.ErrIndexNotFound)
		cord.HeavyMock.Return(&insolar.Reference{}, nil)
		idx, err := (&record.Lifeline{}).Marshal()
		reps := make(chan *message.Message, 1)
		reps <- payload.MustNewMessage(&payload.Meta{
			Payload: payload.MustMarshal(&payload.Index{
				Polymorph: 0,
				Index:     idx,
			}),
		})
		sender.SendTargetMock.Return(reps, func() {})
		writeAccessor.BeginMock.Return(func() {}, nil)
		indexes.SetIfNoneMock.Return()

		p := proc.NewEnsureIndex(gen.ID(), gen.JetID(), payload.Meta{}, pulse.MinTimePulse)
		p.Dep(indexes, cord, sender, writeAccessor)
		err = p.Proceed(ctx)
		assert.NoError(t, err)
	})

	t.Run("fetches from heavy if not found, returns CodeNotFound", func(t *testing.T) {
		setup()
		defer mc.Finish()

		objectID := gen.ID()
		indexes.ForIDMock.Set(func(ctx context.Context, pn insolar.PulseNumber, objID insolar.ID) (record.Index, error) {
			assert.Equal(t, insolar.GenesisPulse.PulseNumber, pn)
			assert.Equal(t, objectID, objID)

			return record.Index{}, object.ErrIndexNotFound
		})
		cord.HeavyMock.Return(&insolar.Reference{}, nil)
		reps := make(chan *message.Message, 1)
		reps <- payload.MustNewMessage(&payload.Meta{
			Payload: payload.MustMarshal(&payload.Error{
				Code: payload.CodeNotFound,
			}),
		})
		sender.SendTargetMock.Return(reps, func() {})

		p := proc.NewEnsureIndex(objectID, gen.JetID(), payload.Meta{}, pulse.MinTimePulse)
		p.Dep(indexes, cord, sender, writeAccessor)
		err := p.Proceed(ctx)
		assert.Error(t, err)
		coded, ok := err.(*payload.CodedError)
		require.True(t, ok, "wrong error type")
		assert.Equal(t, payload.CodeNotFound, coded.Code, "wrong error code")
	})
}
