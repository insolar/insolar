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

package executor_test

import (
	"testing"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/gojuno/minimock"
	"github.com/stretchr/testify/assert"

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

		heavy := []insolar.Node{{*insolar.NewReference(gen.ID()), insolar.StaticRoleHeavyMaterial}}
		s := executor.NewStateIniter(
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

	t.Run("fetching init data", func(t *testing.T) {
		setup()
		defer mc.Finish()

		j1 := gen.JetID()
		j2 := gen.JetID()

		jets := []insolar.JetID{j1, j2}
		heavy := []insolar.Node{{*insolar.NewReference(gen.ID()), insolar.StaticRoleHeavyMaterial}}

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
