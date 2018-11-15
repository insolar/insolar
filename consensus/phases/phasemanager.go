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
	"fmt"
	"time"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/log"
	"github.com/pkg/errors"
)

type PhaseManager struct {
	FirstPhase           *FirstPhase           `inject:""`
	SecondPhase          *SecondPhase          `inject:""`
	ThirdPhasePulse      *ThirdPhasePulse      `inject:""`
	ThirdPhaseReferendum *ThirdPhaseReferendum `inject:""`

	PulseManager core.PulseManager `inject:""`
}

// NewPhaseManager creates and returns a new phase manager.
func NewPhaseManager() *PhaseManager {
	return &PhaseManager{}
}

// Start starts calculate args on phases.
func (pm *PhaseManager) OnPulse(ctx context.Context, pulse *core.Pulse) error {
	var err error

	currentPulse, err := pm.PulseManager.Current(ctx)
	if err != nil {
		return errors.Wrap(err, " [ OnPulse ] Failed to get current pulse")
	}

	if pulse.PulseNumber <= currentPulse.PulseNumber {
		return errors.Wrap(err, " [ OnPulse ] New pulse number equal or less current")
	}

	pulseDuration, err := getPulseDuration(currentPulse)
	if err != nil {
		return errors.Wrap(err, "[ OnPulse ] Failed to get pulse duration")
	}

	var tctx context.Context
	var cancel context.CancelFunc

	tctx, cancel = contextTimeout(ctx, *pulseDuration, 0.2)
	defer cancel()

	firstPhaseState, err := pm.FirstPhase.Execute(tctx, pulse)
	checkError(err)

	tctx, cancel = contextTimeout(ctx, *pulseDuration, 0.2)
	defer cancel()

	secondPhaseState, err := pm.SecondPhase.Execute(ctx, firstPhaseState)
	checkError(err)

	fmt.Println(secondPhaseState) // TODO: remove after use

	// checkError(runPhase(ctx, func() error {
	// 	return pm.ThirdPhasePulse.Execute(secondPhaseState)
	// }))
	// checkError(runPhase(ctx, func() error {
	// 	return pm.ThirdPhaseReferendum.Execute(secondPhaseState)
	// }))

	return nil
}

func getPulseDuration(pulse *core.Pulse) (*time.Duration, error) {
	duration := time.Duration(pulse.PulseNumber-pulse.PrevPulseNumber) * time.Second
	return &duration, nil
}

func contextTimeout(ctx context.Context, duration time.Duration, k float64) (context.Context, context.CancelFunc) {
	timeout := time.Duration(k * float64(duration))
	timedCtx, cancelFund := context.WithTimeout(ctx, timeout)
	return timedCtx, cancelFund
}

func checkError(err error) {
	if err != nil {
		log.Error(err)
	}
}
