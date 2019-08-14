/*
 *    Copyright 2019 Insolar Technologies
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package proc

import (
	"testing"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/stretchr/testify/require"
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/ledger/heavy/executor"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/ledger/drop"
	"github.com/insolar/insolar/insolar/bus"
	"github.com/ThreeDotsLabs/watermill/message"
	"context"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/insolar/store"
)

func TestNewSendInitialState(t *testing.T) {
	meta := payload.Meta{}
	expected := &SendInitialState{
		meta: meta,
	}
	is := NewSendInitialState(meta)
	require.Equal(t, expected, is)
}

func TestSendInitialState_Dep(t *testing.T) {
	startPulse := pulse.NewStartPulse()
	jetKeeper := executor.NewJetKeeperMock(t)
	jetTree := jet.NewStorageMock(t)
	jetCoordinator := jet.NewCoordinatorMock(t)
	dropDB := &drop.DB{}
	pulseAccessor := pulse.NewAccessorMock(t)
	sender := bus.NewSenderMock(t)

	is := NewSendInitialState(payload.Meta{})
	is.Dep(startPulse, jetKeeper, jetTree, jetCoordinator, dropDB, pulseAccessor, sender)
	require.Equal(t, startPulse, is.dep.startPulse)
	require.Equal(t, jetKeeper, is.dep.jetKeeper)
	require.Equal(t, jetTree, is.dep.jetTree)
	require.Equal(t, jetCoordinator, is.dep.jetCoordinator)
	require.Equal(t, dropDB, is.dep.dropDB)
	require.Equal(t, pulseAccessor, is.dep.pulseAccessor)
	require.Equal(t, sender, is.dep.sender)
}

func TestSendInitialState_ProceedNoPulse(t *testing.T) {
	sender := bus.NewSenderMock(t)
	sender.ReplyMock.Set(func(ctx context.Context, origin payload.Meta, reply *message.Message) {
		result, err := payload.Unmarshal(reply.Payload)
		require.NoError(t, err)
		resErr, ok := result.(*payload.Error)
		require.True(t, ok)
		require.Equal(t, uint32(payload.CodeNoStartPulse), resErr.Code)
		require.Equal(t, "Couldn't get start pulse", resErr.Text)
	})
	is := NewSendInitialState(payload.Meta{})
	is.dep.startPulse = pulse.NewStartPulse()
	is.dep.sender = sender
	err  := is.Proceed(context.Background())
	require.NoError(t, err)
}

func TestSendInitialState_ProceedUnknownRequest(t *testing.T) {
	p, err := payload.Marshal(&payload.Request{})
	require.NoError(t, err)
	is := NewSendInitialState(payload.Meta{ Payload: p })
	is.dep.startPulse = pulse.NewStartPulse()
	is.dep.startPulse.SetStartPulse(context.Background(), insolar.Pulse{PulseNumber: 1000})
	err = is.Proceed(context.Background())
	require.EqualError(t, err, "unexpected payload type *payload.Request")
}

func TestSendInitialState_ProceedForNetworkStart(t *testing.T) {
	p, err := payload.Marshal(&payload.GetLightInitialState{ Pulse: 1000 })
	require.NoError(t, err)
	startPulse := pulse.NewStartPulse()
	startPulse.SetStartPulse(context.Background(), insolar.Pulse{ PulseNumber: 1000 })
	topSyncPulse := insolar.Pulse{ PulseNumber: 999}
	jetKeeper := executor.NewJetKeeperMock(t)
	jetKeeper.TopSyncPulseMock.Return(topSyncPulse.PulseNumber)
	pulseAccessor := pulse.NewAccessorMock(t)
	pulseAccessor.ForPulseNumberMock.Return(topSyncPulse, nil)
	jetID := gen.JetID()
	jetTree := jet.NewStorageMock(t)
	jetTree.AllMock.Return([]insolar.JetID{ jetID })
	light := gen.Reference()
	jetCoordinator := jet.NewCoordinatorMock(t)
	jetCoordinator.LightExecutorForJetMock.Return(&light, nil)
	dropItem := drop.MustEncode(&drop.Drop{
		Pulse: 999,
		JetID: jetID,
	})
	db := store.NewDBMock(t)
	db.GetMock.Return(dropItem, nil)
	dropDB := drop.NewDB(db)
	sender := bus.NewSenderMock(t)
	sender.ReplyMock.Set(func(ctx context.Context, origin payload.Meta, reply *message.Message) {
		result, err := payload.Unmarshal(reply.Payload)
		require.NoError(t, err)
		state, ok := result.(*payload.LightInitialState)
		require.True(t, ok)
		require.Equal(t, topSyncPulse.PulseNumber, state.Pulse.PulseNumber)
		require.Equal(t, 1, len(state.JetIDs))
		require.Equal(t, jetID, state.JetIDs[0])
		require.Equal(t, 1, len(state.Drops))
		require.Equal(t, dropItem, state.Drops[0])
		require.True(t, state.NetworkStart)
	})
	is := NewSendInitialState(payload.Meta{
		Payload: p,
		Sender: light,
		Pulse: 1000,
	})
	is.dep.startPulse = startPulse
	is.dep.jetKeeper = jetKeeper
	is.dep.pulseAccessor = pulseAccessor
	is.dep.jetTree = jetTree
	is.dep.jetCoordinator = jetCoordinator
	is.dep.dropDB = dropDB
	is.dep.sender = sender
	err = is.Proceed(context.Background())
	require.NoError(t, err)
}

func TestSendInitialState_ProceedForNetworkStartWithSplit(t *testing.T) {
	p, err := payload.Marshal(&payload.GetLightInitialState{ Pulse: 1000 })
	require.NoError(t, err)
	startPulse := pulse.NewStartPulse()
	startPulse.SetStartPulse(context.Background(), insolar.Pulse{ PulseNumber: 1000 })
	topSyncPulse := insolar.Pulse{ PulseNumber: 999}
	jetKeeper := executor.NewJetKeeperMock(t)
	jetKeeper.TopSyncPulseMock.Return(topSyncPulse.PulseNumber)
	pulseAccessor := pulse.NewAccessorMock(t)
	pulseAccessor.ForPulseNumberMock.Return(topSyncPulse, nil)
	jetID := gen.JetID()
	jetTree := jet.NewStorageMock(t)
	jetTree.AllMock.Return([]insolar.JetID{ jetID })
	light := gen.Reference()
	jetCoordinator := jet.NewCoordinatorMock(t)
	jetCoordinator.LightExecutorForJetMock.Return(&light, nil)
	dropItem := drop.MustEncode(&drop.Drop{
		Pulse: 999,
		JetID: jetID,
		Split: true,
	})
	db := store.NewDBMock(t)
	db.GetMock.Return(dropItem, nil)
	dropDB := drop.NewDB(db)
	sender := bus.NewSenderMock(t)
	sender.ReplyMock.Set(func(ctx context.Context, origin payload.Meta, reply *message.Message) {
		result, err := payload.Unmarshal(reply.Payload)
		require.NoError(t, err)
		state, ok := result.(*payload.LightInitialState)
		require.True(t, ok)
		require.Equal(t, topSyncPulse.PulseNumber, state.Pulse.PulseNumber)
		require.Equal(t, 2, len(state.JetIDs))
		left, right := jet.Siblings(jetID)
		require.Equal(t, left, state.JetIDs[0])
		require.Equal(t, right, state.JetIDs[1])
		require.Equal(t, 1, len(state.Drops))
		require.Equal(t, dropItem, state.Drops[0])
		require.True(t, state.NetworkStart)
	})
	is := NewSendInitialState(payload.Meta{
		Payload: p,
		Sender: light,
		Pulse: 1000,
	})
	is.dep.startPulse = startPulse
	is.dep.jetKeeper = jetKeeper
	is.dep.pulseAccessor = pulseAccessor
	is.dep.jetTree = jetTree
	is.dep.jetCoordinator = jetCoordinator
	is.dep.dropDB = dropDB
	is.dep.sender = sender
	err = is.Proceed(context.Background())
	require.NoError(t, err)
}

func TestSendInitialState_ProceedForJoiner(t *testing.T) {
	p, err := payload.Marshal(&payload.GetLightInitialState{ Pulse: 1001 })
	require.NoError(t, err)
	startPulse := pulse.NewStartPulse()
	startPulse.SetStartPulse(context.Background(), insolar.Pulse{ PulseNumber: 1000 })
	topSyncPulse := insolar.Pulse{ PulseNumber: 999}
	jetKeeper := executor.NewJetKeeperMock(t)
	jetKeeper.TopSyncPulseMock.Return(topSyncPulse.PulseNumber)
	pulseAccessor := pulse.NewAccessorMock(t)
	pulseAccessor.ForPulseNumberMock.Return(topSyncPulse, nil)
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
		Sender: light,
		Pulse: 1001,
	})
	is.dep.startPulse = startPulse
	is.dep.jetKeeper = jetKeeper
	is.dep.pulseAccessor = pulseAccessor
	is.dep.sender = sender
	err = is.Proceed(context.Background())
	require.NoError(t, err)
}
