// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package executor_test

import (
	"testing"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/node"
	"github.com/insolar/insolar/insolar/payload"
	insolarPulse "github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/drop"
	"github.com/insolar/insolar/ledger/light/executor"
	"github.com/insolar/insolar/ledger/object"
	"github.com/insolar/insolar/pulse"
)

func TestStateIniterDefault_PrepareState(t *testing.T) {
	ctx := inslogger.TestContext(t)
	mc := minimock.NewController(t)

	var (
		jetModifier   *jet.ModifierMock
		jetReleaser   *executor.JetReleaserMock
		drops         *drop.ModifierMock
		nodes         *node.AccessorMock
		sender        *bus.SenderMock
		pulseAppender *insolarPulse.AppenderMock
		pulseAccessor *insolarPulse.AccessorMock
		jetCalculator *executor.JetCalculatorMock
		indexes       *object.MemoryIndexModifierMock
	)

	setup := func() {
		jetModifier = jet.NewModifierMock(mc)
		jetReleaser = executor.NewJetReleaserMock(mc)
		drops = drop.NewModifierMock(mc)
		nodes = node.NewAccessorMock(mc)
		sender = bus.NewSenderMock(mc)
		pulseAppender = insolarPulse.NewAppenderMock(mc)
		pulseAccessor = insolarPulse.NewAccessorMock(mc)
		jetCalculator = executor.NewJetCalculatorMock(mc)
		indexes = object.NewMemoryIndexModifierMock(mc)
	}

	t.Run("wrong pulse", func(t *testing.T) {
		setup()
		defer mc.Finish()

		s := executor.NewStateIniter(
			0,
			jetModifier,
			jetReleaser,
			drops,
			nodes,
			sender,
			pulseAppender,
			pulseAccessor,
			jetCalculator,
			indexes,
		)

		_, _, err := s.PrepareState(ctx, pulse.MinTimePulse/2)
		assert.Error(t, err, "must return error 'invalid pulse'")
	})

	t.Run("wrong heavy", func(t *testing.T) {
		setup()
		defer mc.Finish()

		var heavy []insolar.Node
		s := executor.NewStateIniter(
			0,
			jetModifier,
			jetReleaser,
			drops,
			nodes.InRoleMock.Return(heavy, nil),
			sender,
			pulseAppender,
			pulseAccessor.LatestMock.Return(insolar.Pulse{}, insolarPulse.ErrNotFound),
			jetCalculator,
			indexes,
		)

		justAdded, jetsReturned, err := s.PrepareState(ctx, pulse.MinTimePulse)
		assert.Error(t, err, "must return error 'failed to calculate heavy node for pulse'")
		assert.Nil(t, jetsReturned)
		assert.False(t, justAdded)
	})

	t.Run("no need to fetch init data", func(t *testing.T) {
		setup()
		defer mc.Finish()

		jets := []insolar.JetID{gen.JetID(), gen.JetID(), gen.JetID()}
		s := executor.NewStateIniter(
			0,
			jetModifier,
			jetReleaser,
			drops,
			nodes,
			sender,
			pulseAppender,
			pulseAccessor.LatestMock.Return(insolar.Pulse{PulseNumber: pulse.MinTimePulse + 10}, nil),
			jetCalculator.MineForPulseMock.Return(jets, nil),
			indexes,
		)

		justAdded, jetsReturned, err := s.PrepareState(ctx, pulse.MinTimePulse)
		assert.NoError(t, err, "must be nil")
		assert.Equal(t, jets, jetsReturned)
		assert.False(t, justAdded)
	})

	t.Run("fetching init data failing on heavy", func(t *testing.T) {
		setup()
		defer mc.Finish()

		reps := make(chan *message.Message, 1)
		reps <- payload.MustNewMessage(&payload.Meta{
			Payload: payload.MustMarshal(&payload.Error{
				Code: payload.CodeUnknown,
			}),
		})
		sender.SendTargetMock.Return(reps, func() {})

		heavy := []insolar.Node{{ID: *insolar.NewReference(gen.ID()), Role: insolar.StaticRoleHeavyMaterial}}
		s := executor.NewStateIniter(
			0,
			jetModifier,
			jetReleaser,
			drops,
			nodes.InRoleMock.Return(heavy, nil),
			sender,
			pulseAppender,
			pulseAccessor.LatestMock.Return(insolar.Pulse{}, insolarPulse.ErrNotFound),
			jetCalculator,
			indexes,
		)

		justAdded, jetsReturned, err := s.PrepareState(ctx, pulse.MinTimePulse)
		assert.Error(t, err, "must be error 'failed to fetch state from heavy'")
		assert.Nil(t, jetsReturned)
		assert.False(t, justAdded)
	})

	t.Run("panic because of configuration mismatch", func(t *testing.T) {
		setup()
		defer mc.Finish()

		heavy := []insolar.Node{{ID: *insolar.NewReference(gen.ID()), Role: insolar.StaticRoleHeavyMaterial}}

		reps := make(chan *message.Message, 1)
		reps <- payload.MustNewMessage(&payload.Meta{
			Payload: payload.MustMarshal(&payload.LightInitialState{
				LightChainLimit: 10,
			}),
		})

		s := executor.NewStateIniter(
			5,
			jetModifier,
			jetReleaser,
			drops,
			nodes.InRoleMock.Return(heavy, nil),
			sender.SendTargetMock.Return(reps, func() {}),
			pulseAppender,
			pulseAccessor.LatestMock.Return(insolar.Pulse{}, insolarPulse.ErrNotFound),
			jetCalculator,
			indexes,
		)

		// configuration mismatch: LightChainLimit: from heavy 10, from light 5
		require.Panics(t, func() {
			_, _, _ = s.PrepareState(ctx, pulse.MinTimePulse+10)
		})
	})

	t.Run("fetching init data", func(t *testing.T) {
		setup()
		defer mc.Finish()

		j1 := gen.JetID()
		j2 := gen.JetID()

		jets := []insolar.JetID{j1, j2}
		heavy := []insolar.Node{{ID: *insolar.NewReference(gen.ID()), Role: insolar.StaticRoleHeavyMaterial}}

		reps := make(chan *message.Message, 1)
		reps <- payload.MustNewMessage(&payload.Meta{
			Payload: payload.MustMarshal(&payload.LightInitialState{
				NetworkStart: true,
				JetIDs:       jets,
				Pulse: insolarPulse.PulseProto{
					PulseNumber: pulse.MinTimePulse,
				},
				Drops: []drop.Drop{
					{JetID: j1, Pulse: pulse.MinTimePulse},
					{JetID: j2, Pulse: pulse.MinTimePulse},
				},
			}),
		})

		s := executor.NewStateIniter(
			0,
			jetModifier.UpdateMock.Return(nil),
			jetReleaser.UnlockMock.Return(nil),
			drops.SetMock.Return(nil),
			nodes.InRoleMock.Return(heavy, nil),
			sender.SendTargetMock.Return(reps, func() {}),
			pulseAppender.AppendMock.Return(nil),
			pulseAccessor.LatestMock.Return(insolar.Pulse{}, insolarPulse.ErrNotFound),
			jetCalculator,
			indexes,
		)

		justAdded, jetsReturned, err := s.PrepareState(ctx, pulse.MinTimePulse+10)
		assert.NoError(t, err, "must be nil")
		assert.Equal(t, jets, jetsReturned)
		assert.True(t, justAdded)
	})
}
