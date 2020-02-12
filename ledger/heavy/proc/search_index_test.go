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
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/heavy/proc"
	"github.com/insolar/insolar/ledger/object"
	pulse_core "github.com/insolar/insolar/pulse"
	"github.com/stretchr/testify/require"
)

func TestSearchIndex_Proceed(t *testing.T) {
	t.Parallel()
	mc := minimock.NewController(t)
	ctx := inslogger.TestContext(t)

	var (
		index           *object.IndexAccessorMock
		pulseCalculator *pulse.CalculatorMock
		pulseAccessor   *pulse.AccessorMock
		sender          *bus.SenderMock
		recordAccessor  *object.RecordAccessorMock
	)

	resetComponents := func() {
		index = object.NewIndexAccessorMock(t)
		pulseCalculator = pulse.NewCalculatorMock(t)
		sender = bus.NewSenderMock(t)
		pulseAccessor = pulse.NewAccessorMock(t)
		recordAccessor = object.NewRecordAccessorMock(t)
	}

	newProc := func(msg payload.Meta) *proc.SearchIndex {
		p := proc.NewSearchIndex(msg)
		p.Dep(index, pulseCalculator, pulseAccessor, recordAccessor, sender)
		return p
	}

	resetComponents()
	t.Run("fails if until is less than MinPulseTime", func(t *testing.T) {
		msg := payload.SearchIndex{
			Until: pulse_core.MinTimePulse - 1,
		}
		buf, err := msg.Marshal()
		require.NoError(t, err)
		receivedMeta := payload.Meta{Payload: buf}
		p := newProc(receivedMeta)

		err = p.Proceed(ctx)
		require.Error(t, err)

		mc.Finish()
	})

	resetComponents()
	t.Run("not fails if until is bigger than current pulse", func(t *testing.T) {
		msg := payload.SearchIndex{
			Until: pulse_core.MinTimePulse + 1,
		}
		buf, err := msg.Marshal()
		require.NoError(t, err)
		receivedMeta := payload.Meta{Payload: buf}
		p := newProc(receivedMeta)
		sender.ReplyMock.Set(func(_ context.Context, origin payload.Meta, rep *message.Message) {
			require.Equal(t, receivedMeta, origin)

			resp, err := payload.Unmarshal(rep.Payload)
			require.NoError(t, err)

			res, ok := resp.(*payload.SearchIndexInfo)
			require.True(t, ok)

			require.Nil(t, res.Index)
		})

		pulseAccessor.LatestMock.Return(*insolar.GenesisPulse, nil)

		err = p.Proceed(ctx)
		require.NoError(t, err)

		mc.Finish()
	})

	resetComponents()
	t.Run("search through 3 pulsed and find the index", func(t *testing.T) {
		objID := *insolar.NewID(pulse_core.MinTimePulse+100, []byte{1, 2, 3, 4})
		lflParent := gen.RecordReference()
		expectedIdx := record.Index{
			Polymorph: 0,
			ObjID:     objID,
			Lifeline: record.Lifeline{
				Parent: lflParent,
			},
			LifelineLastUsed: 0,
			PendingRecords:   nil,
		}

		pulseAccessor.LatestMock.Return(insolar.Pulse{PulseNumber: pulse_core.MinTimePulse + 100}, nil)

		recordAccessor.ForIDMock.When(ctx, objID).Then(record.Material{}, object.ErrNotFound)
		recordAccessor.ForIDMock.When(ctx, *insolar.NewID(pulse_core.MinTimePulse+99, []byte{1, 2, 3, 4})).Then(record.Material{}, object.ErrNotFound)
		recordAccessor.ForIDMock.When(ctx, *insolar.NewID(pulse_core.MinTimePulse+98, []byte{1, 2, 3, 4})).Then(record.Material{}, nil)

		index.LastKnownForIDMock.When(ctx, *insolar.NewID(pulse_core.MinTimePulse+98, []byte{1, 2, 3, 4})).Then(expectedIdx, nil)

		pulseCalculator.BackwardsMock.When(ctx, insolar.PulseNumber(pulse_core.MinTimePulse+100), 1).Then(insolar.Pulse{
			PulseNumber: pulse_core.MinTimePulse + 99,
		}, nil)

		pulseCalculator.BackwardsMock.When(ctx, insolar.PulseNumber(pulse_core.MinTimePulse+99), 1).Then(insolar.Pulse{
			PulseNumber: pulse_core.MinTimePulse + 98,
		}, nil)

		msg := payload.SearchIndex{
			ObjectID: objID,
			Until:    insolar.PulseNumber(pulse_core.MinTimePulse + 90),
		}
		buf, err := msg.Marshal()
		require.NoError(t, err)
		receivedMeta := payload.Meta{Payload: buf}
		p := newProc(receivedMeta)

		sender.ReplyMock.Set(func(_ context.Context, origin payload.Meta, rep *message.Message) {
			require.Equal(t, receivedMeta, origin)

			resp, err := payload.Unmarshal(rep.Payload)
			require.NoError(t, err)

			res, ok := resp.(*payload.SearchIndexInfo)
			require.True(t, ok)

			require.Equal(t, lflParent, res.Index.Lifeline.Parent)
		})

		err = p.Proceed(ctx)
		require.NoError(t, err)

		mc.Finish()
	})

	resetComponents()
	t.Run("search through 3 pulsed with objID bigger than Latest", func(t *testing.T) {
		objID := *insolar.NewID(pulse_core.MinTimePulse+400, []byte{1, 2, 3, 4})
		lflParent := gen.RecordReference()
		expectedIdx := record.Index{
			Polymorph: 0,
			ObjID:     objID,
			Lifeline: record.Lifeline{
				Parent: lflParent,
			},
			LifelineLastUsed: 0,
			PendingRecords:   nil,
		}

		pulseAccessor.LatestMock.Return(insolar.Pulse{PulseNumber: pulse_core.MinTimePulse + 100}, nil)

		recordAccessor.ForIDMock.When(ctx, *insolar.NewID(pulse_core.MinTimePulse+100, []byte{1, 2, 3, 4})).Then(record.Material{}, object.ErrNotFound)
		recordAccessor.ForIDMock.When(ctx, *insolar.NewID(pulse_core.MinTimePulse+99, []byte{1, 2, 3, 4})).Then(record.Material{}, object.ErrNotFound)
		recordAccessor.ForIDMock.When(ctx, *insolar.NewID(pulse_core.MinTimePulse+98, []byte{1, 2, 3, 4})).Then(record.Material{}, nil)

		index.LastKnownForIDMock.When(ctx, *insolar.NewID(pulse_core.MinTimePulse+98, []byte{1, 2, 3, 4})).Then(expectedIdx, nil)

		pulseCalculator.BackwardsMock.When(ctx, insolar.PulseNumber(pulse_core.MinTimePulse+100), 1).Then(insolar.Pulse{
			PulseNumber: pulse_core.MinTimePulse + 99,
		}, nil)

		pulseCalculator.BackwardsMock.When(ctx, insolar.PulseNumber(pulse_core.MinTimePulse+99), 1).Then(insolar.Pulse{
			PulseNumber: pulse_core.MinTimePulse + 98,
		}, nil)

		msg := payload.SearchIndex{
			ObjectID: objID,
			Until:    insolar.PulseNumber(pulse_core.MinTimePulse + 90),
		}
		buf, err := msg.Marshal()
		require.NoError(t, err)
		receivedMeta := payload.Meta{Payload: buf}
		p := newProc(receivedMeta)

		sender.ReplyMock.Set(func(_ context.Context, origin payload.Meta, rep *message.Message) {
			require.Equal(t, receivedMeta, origin)

			resp, err := payload.Unmarshal(rep.Payload)
			require.NoError(t, err)

			res, ok := resp.(*payload.SearchIndexInfo)
			require.True(t, ok)

			require.Equal(t, lflParent, res.Index.Lifeline.Parent)
		})

		err = p.Proceed(ctx)
		require.NoError(t, err)

		mc.Finish()
	})

	resetComponents()
	t.Run("fails if no lifeline and reaches the limit", func(t *testing.T) {
		objID := *insolar.NewID(pulse_core.MinTimePulse+100, []byte{1, 2, 3, 4})

		pulseAccessor.LatestMock.Return(insolar.Pulse{PulseNumber: pulse_core.MinTimePulse + 100}, nil)

		recordAccessor.ForIDMock.When(ctx, objID).Then(record.Material{}, object.ErrNotFound)
		recordAccessor.ForIDMock.When(ctx, *insolar.NewID(pulse_core.MinTimePulse+99, []byte{1, 2, 3, 4})).Then(record.Material{}, object.ErrNotFound)
		recordAccessor.ForIDMock.When(ctx, *insolar.NewID(pulse_core.MinTimePulse+98, []byte{1, 2, 3, 4})).Then(record.Material{}, object.ErrNotFound)

		pulseCalculator.BackwardsMock.When(ctx, insolar.PulseNumber(pulse_core.MinTimePulse+100), 1).Then(insolar.Pulse{
			PulseNumber: pulse_core.MinTimePulse + 99,
		}, nil)

		pulseCalculator.BackwardsMock.When(ctx, insolar.PulseNumber(pulse_core.MinTimePulse+99), 1).Then(insolar.Pulse{
			PulseNumber: pulse_core.MinTimePulse + 98,
		}, nil)
		pulseCalculator.BackwardsMock.When(ctx, insolar.PulseNumber(pulse_core.MinTimePulse+98), 1).Then(insolar.Pulse{
			PulseNumber: pulse_core.MinTimePulse + 97,
		}, nil)

		msg := payload.SearchIndex{
			ObjectID: objID,
			Until:    insolar.PulseNumber(pulse_core.MinTimePulse + 98),
		}
		buf, err := msg.Marshal()
		require.NoError(t, err)
		receivedMeta := payload.Meta{Payload: buf}
		sender.ReplyMock.Set(func(_ context.Context, origin payload.Meta, rep *message.Message) {
			require.Equal(t, receivedMeta, origin)

			resp, err := payload.Unmarshal(rep.Payload)
			require.NoError(t, err)

			res, ok := resp.(*payload.SearchIndexInfo)
			require.True(t, ok)

			require.Nil(t, res.Index)
		})

		p := newProc(receivedMeta)

		err = p.Proceed(ctx)
		require.NoError(t, err)

		mc.Finish()
	})

}
