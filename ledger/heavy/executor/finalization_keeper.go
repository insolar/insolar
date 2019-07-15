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

//go:generate minimock -i github.com/insolar/insolar/ledger/heavy/executor.FinalizationKeeper -o ./ -s _gen_mock.go

// FinalizationKeeper check how far last finalized pulse and current one from each other
// and if distance is more than limit it stops network
type FinalizationKeeper interface {
	OnPulse(ctx context.Context, current insolar.PulseNumber) error
}

type finalizationKeeper struct {
	jetKeeper       JetKeeper
	networkStopper  insolar.TerminationHandler
	limit           int
	pulseCalculator pulse.Calculator
}

func NewFinalizationKeeper(jk JetKeeper, ns insolar.TerminationHandler, pc pulse.Calculator, limit int) *finalizationKeeper {
	return &finalizationKeeper{
		jetKeeper:       jk,
		networkStopper:  ns,
		limit:           limit - 1,
		pulseCalculator: pc,
	}
}

func (f *finalizationKeeper) OnPulse(ctx context.Context, current insolar.PulseNumber) error {
	logger := inslogger.FromContext(ctx)
	bottomLevel, err := f.pulseCalculator.Backwards(ctx, current, f.limit)
	if err != nil {
		if err == pulse.ErrNotFound {
			logger.Debug("finalizationKeeper: possibly we started not so long ago. Do nothing. Current pulse: ", current)
			return nil
		}
		return errors.Wrap(err, "Can't get old pulse: ")

	}

	lastConfirmedPulse := f.jetKeeper.TopSyncPulse()
	if current < lastConfirmedPulse {
		return errors.New(fmt.Sprintf("Current pulse ( %d ) is less than last confirmed ( %d )", current, lastConfirmedPulse))
	}

	if lastConfirmedPulse <= bottomLevel.PulseNumber {
		f.networkStopper.Leave(ctx, 0)
		return errors.New(fmt.Sprintf("last finalized pulse falls behind too much. Stop node. bottomLevel.PulseNumber: %d, last confirmed: %d", bottomLevel.PulseNumber, lastConfirmedPulse))
	}

	logger.Debugf("FinalizationKeeper: everything is ok. Current pulse: %d, last confirmed: %d, limit: %d", current, lastConfirmedPulse, f.limit)

	return nil
}
