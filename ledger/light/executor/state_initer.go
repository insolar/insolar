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

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/ledger/drop"
	"github.com/insolar/insolar/ledger/light/hot"
	"github.com/pkg/errors"
)

//go:generate minimock -i github.com/insolar/insolar/ledger/light/executor.StateIniter -o ./ -s _mock.go -g

type StateIniter interface {
	// PrepareState prepares actual data to get the light started.
	// Fetch necessary jets and drops from heavy.
	PrepareState(ctx context.Context, pulse insolar.PulseNumber) error
}

// NewStateIniter creates StateIniterDefault with all required components.
func NewStateIniter(
	jetModifier jet.Modifier,
	jetReleaser hot.JetReleaser,
	drops drop.Modifier,
	coordinator jet.Coordinator,
	sender bus.Sender,
) *StateIniterDefault {
	return &StateIniterDefault{
		jetModifier: jetModifier,
		jetReleaser: jetReleaser,
		drops:       drops,
		coordinator: coordinator,
		sender:      sender,
	}
}

// StateIniterDefault implements StateIniter.
type StateIniterDefault struct {
	jetModifier jet.Modifier
	jetReleaser hot.JetReleaser
	drops       drop.Modifier
	coordinator jet.Coordinator
	sender      bus.Sender
}

func (s *StateIniterDefault) PrepareState(ctx context.Context, pulse insolar.PulseNumber) error {
	if pulse < insolar.FirstPulseNumber {
		return errors.Errorf("invalid pulse %s for light state initialization ", pulse)
	}

	heavyRef, err := s.coordinator.Heavy(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to get heavy ref from coordinator")
	}

	msg, err := payload.NewMessage(&payload.GetLightInitialState{
		Pulse: pulse,
	})
	if err != nil {
		return errors.Wrap(err, "failed to create GetInitialState message")
	}

	reps, done := s.sender.SendTarget(ctx, msg, *heavyRef)
	defer done()

	res, ok := <-reps
	if !ok {
		return errors.New("no reply for light state initialization")
	}

	pl, err := payload.UnmarshalFromMeta(res.Payload)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal reply")
	}
	initialState, ok := pl.(*payload.LightInitialState)
	if !ok {
		return fmt.Errorf("unexpected reply %T", pl)
	}

	jets := initialState.JetIDs
	err = s.jetModifier.Update(ctx, pulse, true, jets...)
	if err != nil {
		return errors.Wrap(err, "failed to update jets")
	}

	for _, jetID := range jets {
		err = s.jetReleaser.Unlock(ctx, insolar.ID(jetID))
		if err != nil {
			return errors.Wrap(err, "failed to unlock jet")
		}
	}

	for _, buf := range initialState.Drops {
		d, err := drop.Decode(buf)
		if err != nil {
			return errors.Wrap(err, "failed to decode drop")
		}
		err = s.drops.Set(ctx, *d)
		if err != nil {
			return errors.Wrap(err, "failed to set drop")
		}
	}

	return nil
}
