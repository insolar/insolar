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

package heavy

import (
	"context"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/heavy/replica"
	"github.com/pkg/errors"
)

type dropTruncater interface {
	TruncateHead(ctx context.Context, lastPulse insolar.PulseNumber) error
}

type dbRollback struct {
	drops     dropTruncater
	jetKeeper replica.JetKeeper
}

func (d *dbRollback) Start(ctx context.Context) error {
	logger := inslogger.FromContext(ctx)
	pn := d.jetKeeper.TopSyncPulse()

	logger.Debug("[ dbRollback.Start ] last finalized pulse number: ", pn)
	if pn == insolar.GenesisPulse.PulseNumber {
		logger.Debug("[ dbRollback.Start ] No finalized data. Nothing done")
		return nil
	}

	err := d.drops.TruncateHead(ctx, pn)
	if err != nil {
		return errors.Wrapf(err, "can't truncate db to pulse: %d", pn)
	}

	return nil
}
