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

package executor

import (
	"context"
	"fmt"
	"time"

	"github.com/ThreeDotsLabs/watermill/message"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/backoff"
	"github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/node"
	"github.com/insolar/insolar/insolar/payload"
	insolarPulse "github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/drop"
	"github.com/insolar/insolar/ledger/object"
	"github.com/insolar/insolar/pulse"

	"github.com/pkg/errors"
)

//go:generate minimock -i github.com/insolar/insolar/ledger/light/executor.StateIniter -o ./ -s _mock.go -g

type StateIniter interface {
	// PrepareState prepares actual data to get the light started.
	// Fetch necessary jets and drops from heavy.
	PrepareState(ctx context.Context, pulse insolar.PulseNumber) (justJoined bool, jets []insolar.JetID, err error)
}

const timeout = 10 * time.Second

// NewStateIniter creates StateIniterDefault with all required components.
func NewStateIniter(
	jetModifier jet.Modifier,
	jetReleaser JetReleaser,
	drops drop.Modifier,
	nodes node.Accessor,
	sender bus.Sender,
	pulseAppender insolarPulse.Appender,
	pulseAccessor insolarPulse.Accessor,
	calc JetCalculator,
	indexes object.MemoryIndexModifier,
) *StateIniterDefault {
	return &StateIniterDefault{
		jetModifier:   jetModifier,
		jetReleaser:   jetReleaser,
		drops:         drops,
		nodes:         nodes,
		sender:        sender,
		pulseAppender: pulseAppender,
		pulseAccessor: pulseAccessor,
		jetCalculator: calc,
		indexes:       indexes,
		backoff: backoff.Backoff{
			Factor: 2,
			Jitter: true,
			Min:    50 * time.Millisecond,
			Max:    time.Second,
		},
	}
}

// StateIniterDefault implements StateIniter.
type StateIniterDefault struct {
	jetModifier   jet.Modifier
	jetReleaser   JetReleaser
	drops         drop.Modifier
	nodes         node.Accessor
	sender        bus.Sender
	pulseAppender insolarPulse.Appender
	pulseAccessor insolarPulse.Accessor
	jetCalculator JetCalculator
	backoff       backoff.Backoff
	indexes       object.MemoryIndexModifier
}

func (s *StateIniterDefault) PrepareState(
	ctx context.Context,
	forPulse insolar.PulseNumber,
) (bool, []insolar.JetID, error) {
	if forPulse < pulse.MinTimePulse {
		return false, nil, errors.Errorf("invalid pulse %s for light state initialization ", forPulse)
	}

	// If we have any pulse, it means we already working. No need to fetch any initial data.
	latestPulse, err := s.pulseAccessor.Latest(ctx)
	if err == nil {
		myJets, err := s.jetCalculator.MineForPulse(ctx, latestPulse.PulseNumber)
		if err != nil {
			return false, nil, errors.Wrap(err, "failed to calculate my jets")
		}
		return false, myJets, nil
	}
	if err != insolarPulse.ErrNotFound {
		return false, nil, errors.Wrap(err, "failed to fetch latest pulse")
	}

	heavy, err := s.heavy(forPulse)
	if err != nil {
		return false, nil, err
	}
	msg, err := payload.NewMessage(&payload.GetLightInitialState{
		Pulse: forPulse,
	})
	if err != nil {
		return false, nil, errors.Wrap(err, "failed to create GetInitialState message")
	}

	ctx, cancel := context.WithCancel(ctx)
	go func() {
		select {
		case <-ctx.Done():
		case <-time.After(timeout):
			cancel()
		}
	}()

	jets, err := s.loadStateRetry(ctx, msg, heavy, forPulse)
	if err != nil {
		return false, nil, err
	}
	cancel()

	return true, jets, nil
}

func (s *StateIniterDefault) heavy(pn insolar.PulseNumber) (insolar.Reference, error) {
	candidates, err := s.nodes.InRole(pn, insolar.StaticRoleHeavyMaterial)
	if err != nil {
		return *insolar.NewEmptyReference(), errors.Wrap(err, "failed to calculate heavy node for pulse")
	}
	if len(candidates) == 0 {
		return *insolar.NewEmptyReference(), errors.New("failed to calculate heavy node for pulse")
	}
	return candidates[0].ID, nil
}

func (s *StateIniterDefault) loadStateRetry(
	ctx context.Context,
	msg *message.Message,
	heavy insolar.Reference,
	pn insolar.PulseNumber,
) ([]insolar.JetID, error) {
	reps, done := s.sender.SendTarget(ctx, msg, heavy)
	defer done()

	res, ok := <-reps
	if !ok {
		return nil, errors.New("no reply for light state initialization")
	}

	pl, err := payload.UnmarshalFromMeta(res.Payload)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal reply")
	}

	if errPayload, ok := pl.(*payload.Error); ok {
		if errPayload.Code != payload.CodeNoStartPulse {
			return nil, errors.Wrap(errors.New(errPayload.Text), "failed to fetch state from heavy")
		}
		select {
		case <-ctx.Done():
			return nil, errors.New("retry timeout")
		case <-time.After(s.backoff.Duration()):
			return s.loadStateRetry(ctx, msg, heavy, pn)
		}
	}

	state, ok := pl.(*payload.LightInitialState)
	if !ok {
		return nil, fmt.Errorf("unexpected reply %T", pl)
	}

	prevPulse := insolarPulse.FromProto(&state.Pulse)
	err = s.pulseAppender.Append(ctx, *prevPulse)
	if err != nil {
		return nil, errors.Wrap(err, "failed to append pulse")
	}

	inslogger.FromContext(ctx).WithFields(map[string]interface{}{
		"jets":          insolar.JetIDCollection(state.JetIDs).DebugString(),
		"prev_pulse":    prevPulse.PulseNumber,
		"network_start": state.NetworkStart,
	}).Debug("received initial state from heavy")

	if len(state.JetIDs) < len(state.Drops) {
		return nil, errors.New("Jets count must be greater or equal than drops count")
	}

	// If not network start, we should wait for other lights to give us data.
	if !state.NetworkStart {
		inslogger.FromContext(ctx).Info("Not network start. Wait for other light")
		return nil, nil
	}

	err = s.jetModifier.Update(ctx, pn, true, state.JetIDs...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to update jets")
	}

	for _, jetID := range state.JetIDs {
		err = s.jetReleaser.Unlock(ctx, pn, jetID)
		if err != nil {
			return nil, errors.Wrap(err, "failed to unlock jet")
		}
	}

	for _, d := range state.Drops {
		if d.Pulse != prevPulse.PulseNumber {
			return nil, errors.New("received drop with wrong pulse")
		}
		err = s.drops.Set(ctx, d)
		if err != nil {
			return nil, errors.Wrap(err, "failed to set drop")
		}
	}

	for _, idx := range state.Indexes {
		s.indexes.Set(ctx, pn, idx)
	}

	return state.JetIDs, nil
}
