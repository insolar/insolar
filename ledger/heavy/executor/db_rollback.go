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

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/object"
	"github.com/pkg/errors"
)

//go:generate minimock -i github.com/insolar/insolar/ledger/heavy/executor.headTruncater -o ./ -s _gen_mock.go -g
type headTruncater interface {
	TruncateHead(ctx context.Context, from insolar.PulseNumber) error
}

// DBRollback is used for rollback all data which is not finalized
// It removes all data which was added after pulse which we consider as finalized
type DBRollback struct {
	dbs             []headTruncater
	jetKeeper       JetKeeper
	pulseCalculator pulse.Calculator
}

func NewDBRollback(jetKeeper JetKeeper, pulseCalculator pulse.Calculator, dbs ...headTruncater) *DBRollback {

	return &DBRollback{
		jetKeeper:       jetKeeper,
		dbs:             dbs,
		pulseCalculator: pulseCalculator,
	}
}

func (d *DBRollback) Start(ctx context.Context) error {
	logger := inslogger.FromContext(ctx)
	lastSyncPulseNumber := d.jetKeeper.TopSyncPulse()

	logger.Debug("[ DBRollback.Start ] last finalized pulse number: ", lastSyncPulseNumber)
	if lastSyncPulseNumber == insolar.GenesisPulse.PulseNumber {
		logger.Debug("[ DBRollback.Start ] No finalized data. Nothing done")
		return nil
	}

	pn, err := d.pulseCalculator.Forwards(ctx, lastSyncPulseNumber, 1)
	if err != nil {
		if err == pulse.ErrNotFound {
			inslogger.FromContext(ctx).Debug("No pulse after: ", lastSyncPulseNumber, ". Nothing done.")
			return nil
		}
		return errors.Wrap(err, "pulseCalculator.Forwards returns error")
	}

	for idx, db := range d.dbs {
		if indexDB, ok := db.(object.IndexModifier); ok {
			if err := indexDB.UpdateLastKnownPulse(ctx, lastSyncPulseNumber); err != nil {
				return errors.Wrap(err, "can't update last sync pulse")
			}
		}

		err := db.TruncateHead(ctx, pn.PulseNumber)
		if err != nil {
			return errors.Wrapf(err, "can't truncate %d db to pulse: %d", idx, pn.PulseNumber)
		}
	}

	return nil
}
