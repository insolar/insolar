// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package proc_test

import (
	"context"
	"errors"
	"testing"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/light/executor"
	"github.com/insolar/insolar/ledger/light/proc"
	"github.com/insolar/insolar/pulse"
	"github.com/insolar/insolar/testutils"
)

func TestFetchJet_Proceed(t *testing.T) {
	ctx := inslogger.TestContext(t)
	mc := minimock.NewController(t)

	var (
		jetAccessor *jet.AccessorMock
		jetFetcher  *executor.JetFetcherMock
		coordinator *jet.CoordinatorMock
		sender      *bus.SenderMock
	)
	setup := func() {
		jetAccessor = jet.NewAccessorMock(mc)
		jetFetcher = executor.NewJetFetcherMock(mc)
		coordinator = jet.NewCoordinatorMock(mc)
		sender = bus.NewSenderMock(mc)

	}

	t.Run("error wrong jetID", func(t *testing.T) {
		setup()
		defer mc.Finish()

		jetFetcher.FetchMock.Return(&insolar.ID{}, errors.New("test"))

		p := proc.NewFetchJet(gen.ID(), pulse.MinTimePulse, payload.Meta{}, false)
		p.Dep(jetAccessor, jetFetcher, coordinator, sender)
		err := p.Proceed(ctx)
		assert.Error(t, err, "expected error 'failed to fetch jet'")
	})

	t.Run("virtual is executor", func(t *testing.T) {
		setup()
		defer mc.Finish()

		id := gen.ID()
		jetFetcher.FetchMock.Return(&id, nil)
		me := insolar.NewReference(gen.ID())
		coordinator.LightExecutorForJetMock.Return(me, nil)
		coordinator.MeMock.Return(*me)

		p := proc.NewFetchJet(gen.ID(), pulse.MinTimePulse, payload.Meta{}, false)
		p.Dep(jetAccessor, jetFetcher, coordinator, sender)
		err := p.Proceed(ctx)
		assert.NoError(t, err)
		assert.Equal(t, p.Result.Jet, insolar.JetID(id))
	})

	t.Run("virtual is not executor", func(t *testing.T) {
		setup()
		defer mc.Finish()

		id := gen.ID()
		jetFetcher.FetchMock.Return(&id, nil)
		me := insolar.NewReference(gen.ID())
		notMe := insolar.NewReference(gen.ID())
		coordinator.LightExecutorForJetMock.Return(me, nil)
		coordinator.MeMock.Return(*notMe)

		p := proc.NewFetchJet(gen.ID(), pulse.MinTimePulse, payload.Meta{}, false)
		p.Dep(jetAccessor, jetFetcher, coordinator, sender)
		err := p.Proceed(ctx)
		assert.Error(t, err)
	})

	t.Run("virtual passing", func(t *testing.T) {
		setup()
		defer mc.Finish()

		id := gen.ID()
		jetFetcher.FetchMock.Return(&id, nil)
		msgSender := insolar.NewReference(gen.ID())
		me := insolar.NewReference(gen.ID())
		notMe := insolar.NewReference(gen.ID())
		coordinator.LightExecutorForJetMock.Return(notMe, nil)
		coordinator.MeMock.Return(*me)

		reps := make(chan *message.Message, 1)

		sender.SendTargetMock.Inspect(func(ctx context.Context, msg *message.Message, target insolar.Reference) {
			assert.True(t, target == *notMe || target == *msgSender, "expected messages to msgSender and to right executor")
		}).Return(reps, func() {})

		p := proc.NewFetchJet(gen.ID(), pulse.MinTimePulse, payload.Meta{Sender: *msgSender}, true)
		p.Dep(jetAccessor, jetFetcher, coordinator, sender)
		err := p.Proceed(ctx)
		assert.Error(t, err)
		assert.Equal(t, err, proc.ErrNotExecutor)
	})

}

func TestWaitHot_Proceed(t *testing.T) {
	ctx := inslogger.TestContext(t)
	mc := minimock.NewController(t)

	var (
		jetWaiter *executor.JetWaiterMock
	)
	setup := func() {
		jetWaiter = executor.NewJetWaiterMock(mc)
	}

	t.Run("happy", func(t *testing.T) {
		setup()
		defer mc.Finish()

		j := gen.JetID()
		pn := gen.PulseNumber()

		jetWaiter.WaitMock.Inspect(func(ctx context.Context, jetID insolar.JetID, pulse insolar.PulseNumber) {
			assert.Equal(t, j, jetID)
			assert.Equal(t, pn, pulse)
		}).Return(nil)

		p := proc.NewWaitHot(j, pn, payload.Meta{})
		p.Dep(jetWaiter)
		err := p.Proceed(ctx)
		assert.NoError(t, err)
	})
}

func TestCalculateID_Proceed(t *testing.T) {
	ctx := inslogger.TestContext(t)
	mc := minimock.NewController(t)

	var (
		pcs insolar.PlatformCryptographyScheme
	)
	setup := func() {
		pcs = testutils.NewPlatformCryptographyScheme()
	}

	t.Run("happy", func(t *testing.T) {
		setup()
		defer mc.Finish()

		pl, _ := (&payload.Meta{
			Polymorph:  0,
			Payload:    nil,
			Sender:     insolar.Reference{},
			Receiver:   insolar.Reference{},
			Pulse:      0,
			ID:         nil,
			OriginHash: payload.MessageHash{},
		}).Marshal()

		pn := gen.PulseNumber()

		p := proc.NewCalculateID(pl, pn)
		p.Dep(pcs)
		err := p.Proceed(ctx)
		assert.NoError(t, err)
	})
}
