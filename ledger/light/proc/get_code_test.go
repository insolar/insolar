package proc_test

import (
	"context"
	"testing"
	"time"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/gojuno/minimock/v3"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"

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
)

func TestGetCode_Proceed(t *testing.T) {
	ctx := inslogger.TestContext(t)
	mc := minimock.NewController(t)

	var (
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
	}

	t.Run("Simple success", func(t *testing.T) {
		setup()
		defer mc.Finish()

		rec := record.Material{
			Virtual:  record.Wrap(&record.Code{}),
			ID:       gen.ID(),
			ObjectID: gen.ID(),
		}
		records.ForIDMock.Return(rec, nil)

		msg := payload.Meta{}
		buf, err := rec.Marshal()
		expectedMsg, _ := payload.NewMessage(&payload.Code{
			Record: buf,
		})

		sender.ReplyMock.Inspect(func(ctx context.Context, origin payload.Meta, reply *message.Message) {
			assert.Equal(t, expectedMsg.Payload, reply.Payload)
			assert.Equal(t, msg, origin)
		}).Return()

		p := proc.NewGetCode(msg, gen.ID(), true)
		p.Dep(records, coordinator, fetcher, sender)

		err = p.Proceed(ctx)
		assert.NoError(t, err)
	})

	t.Run("Passing request on heavy", func(t *testing.T) {
		setup()
		defer mc.Finish()
		defer mc.Wait(10 * time.Minute)

		records.ForIDMock.Return(record.Material{}, object.ErrNotFound)

		expectedTarget := insolar.NewReference(gen.ID())

		coordinator.IsBeyondLimitMock.Return(true, nil)
		coordinator.HeavyMock.Return(expectedTarget, nil)

		meta := payload.Meta{}
		buf, _ := meta.Marshal()
		expectedPass, _ := payload.NewMessage(&payload.Pass{
			Origin: buf,
		})

		sender.SendTargetMock.Inspect(func(ctx context.Context, msg *message.Message, target insolar.Reference) {
			assert.Equal(t, expectedPass.Payload, msg.Payload)
			assert.Equal(t, expectedTarget, &target)
		}).Return(
			make(chan *message.Message), func() {})

		p := proc.NewGetCode(meta, gen.ID(), true)
		p.Dep(records, coordinator, fetcher, sender)

		err := p.Proceed(ctx)
		assert.NoError(t, err)
	})

	t.Run("Passing request on light", func(t *testing.T) {
		setup()
		defer mc.Finish()
		defer mc.Wait(10 * time.Minute)

		records.ForIDMock.Return(record.Material{}, object.ErrNotFound)

		expectedTarget := insolar.NewReference(gen.ID())

		coordinator.IsBeyondLimitMock.Return(false, nil)
		expectedJetID := gen.ID()
		fetcher.FetchMock.Return(&expectedJetID, nil)
		coordinator.LightExecutorForJetMock.Inspect(func(ctx context.Context, jetID insolar.ID, pulse insolar.PulseNumber) {
			assert.Equal(t, jetID, expectedJetID)
		}).Return(expectedTarget, nil)

		meta := payload.Meta{}
		buf, _ := meta.Marshal()
		expectedPass, _ := payload.NewMessage(&payload.Pass{
			Origin: buf,
		})

		sender.SendTargetMock.Inspect(func(ctx context.Context, msg *message.Message, target insolar.Reference) {
			assert.Equal(t, expectedPass.Payload, msg.Payload)
			assert.Equal(t, expectedTarget, &target)
		}).Return(
			make(chan *message.Message), func() {})

		p := proc.NewGetCode(meta, gen.ID(), true)
		p.Dep(records, coordinator, fetcher, sender)

		err := p.Proceed(ctx)
		assert.NoError(t, err)
	})

	t.Run("Not passing, returns error", func(t *testing.T) {
		setup()
		defer mc.Finish()

		msg := payload.Meta{}
		records.ForIDMock.Return(record.Material{}, object.ErrNotFound)

		p := proc.NewGetCode(msg, gen.ID(), false)
		p.Dep(records, coordinator, fetcher, sender)

		err := p.Proceed(ctx)
		assert.Error(t, err)
		insError, ok := errors.Cause(err).(*payload.CodedError)
		assert.True(t, ok)
		assert.Equal(t, payload.CodeNotFound, insError.GetCode())
	})
}
