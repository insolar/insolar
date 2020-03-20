// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package proc

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
	"github.com/insolar/insolar/ledger/drop"
	"github.com/insolar/insolar/ledger/heavy/executor"
	"github.com/stretchr/testify/require"
)

func TestNewSendInitialState(t *testing.T) {
	meta := payload.Meta{}
	expected := &SendInitialState{
		meta: meta,
	}
	is := NewSendInitialState(meta, 0)
	require.Equal(t, expected, is)
}

func TestSendInitialState_Dep(t *testing.T) {
	startPulse := pulse.NewStartPulse()
	jetKeeper := executor.NewJetKeeperMock(t)
	stateAccessor := executor.NewInitialStateAccessorMock(t)
	pulseAccessor := pulse.NewAccessorMock(t)
	sender := bus.NewSenderMock(t)

	is := NewSendInitialState(payload.Meta{}, 0)
	is.Dep(startPulse, jetKeeper, stateAccessor, pulseAccessor, sender)
	require.Equal(t, startPulse, is.dep.startPulse)
	require.Equal(t, jetKeeper, is.dep.jetKeeper)
	require.Equal(t, pulseAccessor, is.dep.pulseAccessor)
	require.Equal(t, sender, is.dep.sender)
}

func TestSendInitialState_ProceedNoPulse(t *testing.T) {
	ctx := inslogger.TestContext(t)
	sender := bus.NewSenderMock(t)
	sender.ReplyMock.Set(func(ctx context.Context, origin payload.Meta, reply *message.Message) {
		result, err := payload.Unmarshal(reply.Payload)
		require.NoError(t, err)
		resErr, ok := result.(*payload.Error)
		require.True(t, ok)
		require.Equal(t, payload.CodeNoStartPulse, resErr.Code)
		require.Equal(t, "Couldn't get start pulse", resErr.Text)
	})
	is := NewSendInitialState(payload.Meta{}, 0)
	is.dep.startPulse = pulse.NewStartPulse()
	is.dep.sender = sender
	err := is.Proceed(ctx)
	require.NoError(t, err)
}

func TestSendInitialState_ProceedUnknownRequest(t *testing.T) {
	ctx := inslogger.TestContext(t)
	p, err := payload.Marshal(&payload.Request{})
	require.NoError(t, err)
	is := NewSendInitialState(payload.Meta{Payload: p}, 0)
	is.dep.startPulse = pulse.NewStartPulse()
	is.dep.startPulse.SetStartPulse(ctx, insolar.Pulse{PulseNumber: 1000})

	err = is.Proceed(ctx)
	require.EqualError(t, err, "unexpected payload type *payload.Request")
}

func TestSendInitialState_ProceedForNetworkStart(t *testing.T) {
	mc := minimock.NewController(t)
	defer mc.Finish()
	ctx := inslogger.TestContext(t)

	p, err := payload.Marshal(&payload.GetLightInitialState{Pulse: 1000})
	require.NoError(t, err)

	sp := insolar.Pulse{PulseNumber: 1000}
	startPulse := pulse.NewStartPulse()
	startPulse.SetStartPulse(ctx, sp)

	topSyncPulse := insolar.Pulse{PulseNumber: 999}
	jetKeeper := executor.NewJetKeeperMock(mc)
	jetKeeper.TopSyncPulseMock.Return(topSyncPulse.PulseNumber)

	pulseAccessor := pulse.NewAccessorMock(mc)
	pulseAccessor.ForPulseNumberMock.Return(topSyncPulse, nil)

	light := gen.Reference()

	JetIDs := make([]insolar.JetID, 0)
	Drops := make([]drop.Drop, 0)
	Indexes := make([]record.Index, 0)
	initialStateAccessor := executor.NewInitialStateAccessorMock(mc)
	initialStateAccessor.GetMock.Expect(ctx, light, sp.PulseNumber).Return(&executor.InitialState{
		JetIDs:  JetIDs,
		Drops:   Drops,
		Indexes: Indexes,
	})

	sender := bus.NewSenderMock(t)
	sender.ReplyMock.Set(func(ctx context.Context, origin payload.Meta, reply *message.Message) {
		result, err := payload.Unmarshal(reply.Payload)
		require.NoError(t, err)
		state, ok := result.(*payload.LightInitialState)
		require.True(t, ok)
		require.Equal(t, topSyncPulse.PulseNumber, state.Pulse.PulseNumber)
		require.Equal(t, 0, len(state.JetIDs))
		require.Equal(t, 0, len(state.Drops))
		require.Equal(t, 0, len(state.Indexes))
		require.True(t, state.NetworkStart)
	})

	is := NewSendInitialState(payload.Meta{
		Payload: p,
		Sender:  light,
		Pulse:   1000,
	}, 0)
	is.dep.startPulse = startPulse
	is.dep.jetKeeper = jetKeeper
	is.dep.pulseAccessor = pulseAccessor
	is.dep.sender = sender
	is.dep.initialState = initialStateAccessor

	err = is.Proceed(ctx)
	require.NoError(t, err)
}

func TestSendInitialState_ProceedForJoiner(t *testing.T) {
	mc := minimock.NewController(t)
	defer mc.Finish()
	ctx := inslogger.TestContext(t)

	p, err := payload.Marshal(&payload.GetLightInitialState{Pulse: 1001})
	require.NoError(t, err)

	startPulse := pulse.NewStartPulse()
	startPulse.SetStartPulse(ctx, insolar.Pulse{PulseNumber: 1000})

	topSyncPulse := insolar.Pulse{PulseNumber: 999}
	jetKeeper := executor.NewJetKeeperMock(mc)
	jetKeeper.TopSyncPulseMock.Return(topSyncPulse.PulseNumber)

	pulseAccessor := pulse.NewAccessorMock(mc)
	pulseAccessor.ForPulseNumberMock.Expect(ctx, topSyncPulse.PulseNumber).Return(topSyncPulse, nil)

	light := gen.Reference()
	sender := bus.NewSenderMock(t)
	sender.ReplyMock.Set(func(ctx context.Context, origin payload.Meta, reply *message.Message) {
		result, err := payload.Unmarshal(reply.Payload)
		require.NoError(t, err)
		state, ok := result.(*payload.LightInitialState)
		require.True(t, ok)
		require.Equal(t, topSyncPulse.PulseNumber, state.Pulse.PulseNumber)
		require.False(t, state.NetworkStart)
	})

	is := NewSendInitialState(payload.Meta{
		Payload: p,
		Sender:  light,
		Pulse:   1001,
	}, 0)

	is.dep.startPulse = startPulse
	is.dep.jetKeeper = jetKeeper
	is.dep.pulseAccessor = pulseAccessor
	is.dep.sender = sender

	err = is.Proceed(ctx)
	require.NoError(t, err)
}
