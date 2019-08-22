package proc_test

import (
	"context"
	"testing"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/gojuno/minimock"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/light/executor"
	"github.com/insolar/insolar/ledger/light/proc"
	"github.com/insolar/insolar/ledger/object"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEnsureIndex_Proceed(t *testing.T) {
	ctx := inslogger.TestContext(t)
	mc := minimock.NewController(t)

	var (
		locker        *object.IndexLockerMock
		indexes       *object.MemoryIndexStorageMock
		cord          *jet.CoordinatorMock
		sender        *bus.SenderMock
		writeAccessor *executor.WriteAccessorMock
	)
	setup := func() {
		locker = object.NewIndexLockerMock(mc)
		indexes = object.NewMemoryIndexStorageMock(mc)
		cord = jet.NewCoordinatorMock(mc)
		sender = bus.NewSenderMock(mc)
		writeAccessor = executor.NewWriteAccessorMock(mc)

		locker.LockMock.Return()
		locker.UnlockMock.Return()
	}

	t.Run("returns CodeNotFound if no index", func(t *testing.T) {
		setup()
		defer mc.Finish()

		indexes.ForIDMock.Return(record.Index{}, object.ErrIndexNotFound)
		cord.HeavyMock.Return(&insolar.Reference{}, nil)
		reps := make(chan *message.Message, 1)
		reps <- payload.MustNewMessage(&payload.Meta{
			Payload: payload.MustMarshal(&payload.Error{
				Code: payload.CodeNotFound,
			}),
		})
		sender.SendTargetMock.Return(reps, func() {})

		p := proc.NewEnsureIndex(gen.ID(), gen.JetID(), payload.Meta{}, insolar.FirstPulseNumber)
		p.Dep(locker, indexes, cord, sender, writeAccessor)
		err := p.Proceed(ctx)
		assert.Error(t, err)
		coded, ok := err.(*payload.CodedError)
		require.True(t, ok, "wrong error type")
		assert.Equal(t, uint32(payload.CodeNotFound), coded.Code, "wrong error code")
	})

	t.Run("fetches from heavy if not found", func(t *testing.T) {
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

		p := proc.NewEnsureIndex(objectID, gen.JetID(), payload.Meta{}, insolar.FirstPulseNumber)
		p.Dep(locker, indexes, cord, sender, writeAccessor)
		err := p.Proceed(ctx)
		assert.Error(t, err)
		coded, ok := err.(*payload.CodedError)
		require.True(t, ok, "wrong error type")
		assert.Equal(t, uint32(payload.CodeNotFound), coded.Code, "wrong error code")
	})
}
