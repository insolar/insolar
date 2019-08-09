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
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/pkg/errors"
)

// FinalizationKeeper check how far from each other last finalized pulse and current one
// and if distance is more than limit it stops network
type FinalizationKeeper interface {
	OnPulse(ctx context.Context, current insolar.PulseNumber) error
}

type FinalizationKeeperDefault struct {
	jetKeeper       JetKeeper
	networkStopper  insolar.TerminationHandler
	limit           int
	pulseCalculator pulse.Calculator
}

func NewFinalizationKeeperDefault(jk JetKeeper, ns insolar.TerminationHandler, pc pulse.Calculator, limit int) *FinalizationKeeperDefault {
	return &FinalizationKeeperDefault{
		jetKeeper:       jk,
		networkStopper:  ns,
		limit:           limit,
		pulseCalculator: pc,
	}
}

func (f *FinalizationKeeperDefault) OnPulse(ctx context.Context, current insolar.PulseNumber) error {
	logger := inslogger.FromContext(ctx)
	bottomLevel, err := f.pulseCalculator.Backwards(ctx, current, f.limit)
	if err != nil {
		if err == pulse.ErrNotFound {
			logger.Debug("finalizationKeeper: possibly we started not so long ago. Do nothing. Current pulse: ", current)
			return nil
		}
		return errors.Wrap(err, "Can't get old pulse")
	}

	lastConfirmedPulse := f.jetKeeper.TopSyncPulse()
	if current < lastConfirmedPulse {
		return errors.New(fmt.Sprintf("Current pulse ( %d ) is less than last confirmed ( %d )", current, lastConfirmedPulse))
	}

	if lastConfirmedPulse <= bottomLevel.PulseNumber {
		// f.networkStopper.Leave(ctx, 0)
		panic(fmt.Sprintf("it shouldn't be called. lastConfirmedPulse:%v, bottomLevel.PulseNumber:%v", lastConfirmedPulse, bottomLevel.PulseNumber))
		return errors.New(fmt.Sprintf("last finalized pulse falls behind too much. Stop node. bottomLevel.PulseNumber: %d, last confirmed: %d", bottomLevel.PulseNumber, lastConfirmedPulse))
	}

	logger.Debugf("FinalizationKeeper: everything is ok. Current pulse: %d, last confirmed: %d, limit: %d", current, lastConfirmedPulse, f.limit)

	return nil
}
