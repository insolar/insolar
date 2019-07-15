///
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
///

package executor

import (
	"context"
	"errors"
	"fmt"

	"github.com/insolar/insolar/insolar"
)

type FinalizationKeeper struct {
	jetKeeper      JetKeeper
	networkStopper insolar.TerminationHandler
	limit          int
}

func NewFinalizationKeeper(jetKeeper JetKeeper, networkStopper insolar.TerminationHandler, limit int) *FinalizationKeeper {
	return &FinalizationKeeper{
		jetKeeper:      jetKeeper,
		networkStopper: networkStopper,
		limit:          limit,
	}
}

func (f *FinalizationKeeper) OnPulse(ctx context.Context, current insolar.PulseNumber) error {
	lastConfirmedPulse := f.jetKeeper.TopSyncPulse()
	if current < lastConfirmedPulse {
		return errors.New(fmt.Sprintf("Current pulse ( %d ) is less than last confirmed ( %d )", current, lastConfirmedPulse))
	}
	if int(current-lastConfirmedPulse) > f.limit {
		f.networkStopper.Leave(ctx, 0)
		return errors.New(fmt.Sprintf("Last finalized pulse falls behind too much. Stop node. Current pulse: %d, last confirmed: %d", current, lastConfirmedPulse))
	}

	return nil
}
