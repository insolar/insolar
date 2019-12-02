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
	"github.com/gojuno/minimock/v3"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/heavy/proc"
)

func TestSendPulse_Proceed(t *testing.T) {
	mc := minimock.NewController(t)
	ctx := inslogger.TestContext(t)

	var (
		sender *bus.SenderMock
		pulses *pulse.AccessorMock
	)

	setup := func() {
		sender = bus.NewSenderMock(mc)
		pulses = pulse.NewAccessorMock(mc)
	}

	newProc := func(msg payload.Meta) *proc.SendPulse {
		p := proc.NewSendPulse(msg)
		p.Dep(pulses, sender)
		return p
	}

	t.Run("pulse not found", func(t *testing.T) {
		setup()
		defer mc.Finish()

		p := newProc(payload.Meta{})
		pulses.ForPulseNumberMock.Return(insolar.Pulse{}, pulse.ErrNotFound)

		err := p.Proceed(ctx)
		require.Error(t, err)
		assert.Equal(t, pulse.ErrNotFound, errors.Cause(err))
	})

	t.Run("happy basic", func(t *testing.T) {
		setup()
		defer mc.Finish()

		pulseNumber := insolar.PulseNumber(42)
		msg := payload.GetPulse{
			PulseNumber: pulseNumber,
		}
		buf, err := msg.Marshal()
		require.NoError(t, err)
		receivedMeta := payload.Meta{Payload: buf}
		p := newProc(receivedMeta)

		expectedPulse := insolar.Pulse{PulseNumber: pulseNumber}
		pulses.ForPulseNumberMock.Return(expectedPulse, nil)
		sender.ReplyMock.Set(func(_ context.Context, origin payload.Meta, rep *message.Message) {
			require.Equal(t, receivedMeta, origin)

			resp, err := payload.Unmarshal(rep.Payload)
			require.NoError(t, err)

			res, ok := resp.(*payload.Pulse)
			require.True(t, ok)
			require.Equal(t, expectedPulse.PulseNumber, pulse.FromProto(&res.Pulse).PulseNumber)
		})

		err = p.Proceed(ctx)
		require.NoError(t, err)
	})
}
