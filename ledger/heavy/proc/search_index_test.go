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

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/gojuno/minimock"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/heavy/proc"
	"github.com/insolar/insolar/ledger/object"
	"github.com/stretchr/testify/require"
)

func TestSearchIndex_Proceed(t *testing.T) {
	t.Parallel()
	mc := minimock.NewController(t)
	ctx := inslogger.TestContext(t)

	var (
		index  *object.IndexAccessorMock
		pulses *pulse.CalculatorMock
		sender *bus.SenderMock
	)

	resetComponents := func() {
		index = object.NewIndexAccessorMock(t)
		pulses = pulse.NewCalculatorMock(t)
		sender = bus.NewSenderMock(t)
	}

	newProc := func(msg payload.Meta) *proc.SearchIndex {
		p := proc.NewSearchIndex(msg)
		p.Dep(index, pulses, sender)
		return p
	}

	resetComponents()
	t.Run("search through 3 pulsed and find the index", func(t *testing.T) {
		objID := *insolar.NewID(100, []byte{1, 2, 3, 4})
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

		index.ForIDMock.When(ctx, insolar.PulseNumber(100), objID).Then(record.Index{}, object.ErrIndexNotFound)
		index.ForIDMock.When(ctx, insolar.PulseNumber(99), objID).Then(record.Index{}, object.ErrIndexNotFound)
		index.ForIDMock.When(ctx, insolar.PulseNumber(98), objID).Then(expectedIdx, nil)

		pulses.BackwardsMock.When(ctx, insolar.PulseNumber(100), 1).Then(insolar.Pulse{
			PulseNumber: 99,
		}, nil)

		pulses.BackwardsMock.When(ctx, insolar.PulseNumber(99), 1).Then(insolar.Pulse{
			PulseNumber: 98,
		}, nil)

		msg := payload.SearchIndex{
			ObjectID: objID,
			Until:    insolar.PulseNumber(90),
		}
		buf, err := msg.Marshal()
		require.NoError(t, err)
		receivedMeta := payload.Meta{Payload: buf}
		p := newProc(receivedMeta)

		sender.ReplyMock.Set(func(_ context.Context, origin payload.Meta, rep *message.Message) {
			require.Equal(t, receivedMeta, origin)

			resp, err := payload.Unmarshal(rep.Payload)
			require.NoError(t, err)

			res, ok := resp.(*payload.Index)
			require.True(t, ok)

			lfl := record.Lifeline{}
			err = lfl.Unmarshal(res.Index)
			require.NoError(t, err)

			require.Equal(t, lflParent, lfl.Parent)
		})

		err = p.Proceed(ctx)
		require.NoError(t, err)

		mc.Finish()
	})

	resetComponents()
	t.Run("fails if no lifeline and reaches the limit", func(t *testing.T) {
		objID := *insolar.NewID(100, []byte{1, 2, 3, 4})

		index.ForIDMock.When(ctx, insolar.PulseNumber(100), objID).Then(record.Index{}, object.ErrIndexNotFound)
		index.ForIDMock.When(ctx, insolar.PulseNumber(99), objID).Then(record.Index{}, object.ErrIndexNotFound)
		index.ForIDMock.When(ctx, insolar.PulseNumber(98), objID).Then(record.Index{}, object.ErrIndexNotFound)

		pulses.BackwardsMock.When(ctx, insolar.PulseNumber(100), 1).Then(insolar.Pulse{
			PulseNumber: 99,
		}, nil)

		pulses.BackwardsMock.When(ctx, insolar.PulseNumber(99), 1).Then(insolar.Pulse{
			PulseNumber: 98,
		}, nil)
		pulses.BackwardsMock.When(ctx, insolar.PulseNumber(98), 1).Then(insolar.Pulse{
			PulseNumber: 97,
		}, nil)

		msg := payload.SearchIndex{
			ObjectID: objID,
			Until:    insolar.PulseNumber(98),
		}
		buf, err := msg.Marshal()
		require.NoError(t, err)
		receivedMeta := payload.Meta{Payload: buf}
		p := newProc(receivedMeta)

		err = p.Proceed(ctx)
		require.Error(t, err)

		codedError := err.(*payload.CodedError)
		require.Equal(t, uint32(payload.CodeNotFound), codedError.Code)

		mc.Finish()
	})
}
