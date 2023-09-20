package proc_test

import (
	"context"
	"testing"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/gojuno/minimock/v3"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"

	"github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/light/executor"
	"github.com/insolar/insolar/ledger/light/proc"
)

func TestSendFilament_Proceed(t *testing.T) {
	ctx := inslogger.TestContext(t)
	mc := minimock.NewController(t)

	var (
		sender    *bus.SenderMock
		filaments *executor.FilamentCalculatorMock
	)

	setup := func() {
		sender = bus.NewSenderMock(mc)
		filaments = executor.NewFilamentCalculatorMock(mc)
	}

	t.Run("simple success", func(t *testing.T) {
		setup()
		defer mc.Finish()

		obj := gen.ID()
		pl, _ := (&payload.GetFilament{}).Marshal()

		msg := payload.Meta{
			Payload: pl,
		}

		storageRecs := make([]record.CompositeFilamentRecord, 5)
		filaments.RequestsMock.Return(storageRecs, nil)
		expectedMsg, _ := payload.NewMessage(&payload.FilamentSegment{
			ObjectID: obj,
			Records:  storageRecs,
		})

		sender.ReplyMock.Inspect(func(ctx context.Context, origin payload.Meta, reply *message.Message) {
			assert.Equal(t, expectedMsg.Payload, reply.Payload)
			assert.Equal(t, msg, origin)
		}).Return()

		p := proc.NewSendFilament(msg, obj, gen.ID(), gen.PulseNumber())
		p.Dep(sender, filaments)

		err := p.Proceed(ctx)
		assert.NoError(t, err)
	})

	t.Run("requests not found sends error", func(t *testing.T) {
		setup()
		defer mc.Finish()

		obj := gen.ID()
		pl, _ := (&payload.GetFilament{}).Marshal()

		msg := payload.Meta{
			Payload: pl,
		}

		var storageRecs []record.CompositeFilamentRecord
		filaments.RequestsMock.Return(storageRecs, nil)

		p := proc.NewSendFilament(msg, obj, gen.ID(), gen.PulseNumber())
		p.Dep(sender, filaments)

		err := p.Proceed(ctx)
		assert.Error(t, err)
		insError, ok := errors.Cause(err).(*payload.CodedError)
		assert.True(t, ok)
		assert.Equal(t, payload.CodeNotFound, insError.GetCode())

	})
}
