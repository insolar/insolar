/*
 *    Copyright 2018 Insolar
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

package phases

import (
	"context"
	"time"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/log"
)

type PhaseManager struct {
	FirstPhase           *FirstPhase           `inject:""`
	SecondPhase          *SecondPhase          `inject:""`
}

// NewPhaseManager creates and returns a new phase manager.
func NewPhaseManager() *PhaseManager {
	return &PhaseManager{}
}

// Start starts calculate args on phases.
func (pm *PhaseManager) OnPulse(ctx context.Context, pulse *core.Pulse) error {
	pulseDuration := pm.getPulseDuration()

	checkError(runPhase(contextTimeout(ctx, pulseDuration, 0.2), func() error {
		return pm.FirstPhase.Execute(pulse)
	}))

	firstPhaseState := pm.FirstPhase.State
	pm.FirstPhase.State = &FirstPhaseState{}

	return nil
}

func (pm *PhaseManager) getPulseDuration() time.Duration {
	// TODO: calculate
	return 10 * time.Second
}

func runPhase(ctx context.Context, phase func() error) error {
	done := make(chan error, 1)
	go func() {
		done <- phase()
	}()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-done:
		return err
	}
}

func contextTimeout(ctx context.Context, duration time.Duration, k float64) context.Context {
	timeout := time.Duration(k * float64(duration))
	timedCtx, _ := context.WithTimeout(ctx, timeout)
	return timedCtx
}

func checkError(err error) {
	if err != nil {
		log.Error(err)
	}
}
