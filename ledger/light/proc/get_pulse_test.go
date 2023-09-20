package proc_test

import (
	"testing"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/payload"
	insPulse "github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/light/proc"
	"github.com/insolar/insolar/pulse"
)

func TestGetPulse_Proceed(t *testing.T) {
	ctx := inslogger.TestContext(t)
	mc := minimock.NewController(t)

	var (
		sender      *bus.SenderMock
		coordinator *jet.CoordinatorMock
	)

	setup := func() {
		sender = bus.NewSenderMock(mc)
		coordinator = jet.NewCoordinatorMock(mc)
	}

	t.Run("Simple success", func(t *testing.T) {
		setup()
		defer mc.Finish()

		coordinator.HeavyMock.Return(&insolar.Reference{}, nil)
		reps := make(chan *message.Message, 1)
		reps <- payload.MustNewMessage(&payload.Meta{
			Polymorph: uint32(payload.TypeMeta),
			Payload: payload.MustMarshal(&payload.Pulse{
				Polymorph: uint32(payload.TypePulse),
				Pulse:     *insPulse.ToProto(&insolar.Pulse{PulseNumber: pulse.MinTimePulse}),
			}),
		})
		sender.
			SendTargetMock.Return(reps, func() {}).
			ReplyMock.Return()

		p := proc.NewGetPulse(payload.Meta{}, pulse.MinTimePulse)
		p.Dep(coordinator, sender)
		err := p.Proceed(ctx)
		assert.NoError(t, err)
	})

	t.Run("fetches from heavy if not found, returns CodeNotFound", func(t *testing.T) {
		setup()
		defer mc.Finish()

		coordinator.HeavyMock.Return(&insolar.Reference{}, nil)
		reps := make(chan *message.Message, 1)
		reps <- payload.MustNewMessage(&payload.Meta{
			Polymorph: uint32(payload.TypeMeta),
			Payload: payload.MustMarshal(&payload.Error{
				Polymorph: uint32(payload.TypeError),
				Code:      payload.CodeNotFound,
			}),
		})
		sender.SendTargetMock.Return(reps, func() {})

		p := proc.NewGetPulse(payload.Meta{}, pulse.MinTimePulse)
		p.Dep(coordinator, sender)
		err := p.Proceed(ctx)
		assert.Error(t, err)
		coded, ok := err.(*payload.CodedError)
		require.True(t, ok, "wrong error type")
		assert.Equal(t, payload.CodeNotFound, coded.Code, "wrong error code")
	})
}
