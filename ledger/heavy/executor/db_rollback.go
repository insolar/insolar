// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package executor

import (
	"context"

	"github.com/insolar/insolar/insolar"
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
	dbs       []headTruncater
	jetKeeper JetKeeper
}

func NewDBRollback(jetKeeper JetKeeper, dbs ...headTruncater) *DBRollback {
	return &DBRollback{
		jetKeeper: jetKeeper,
		dbs:       dbs,
	}
}

func (d *DBRollback) Start(ctx context.Context) error {
	lastSyncPulseNumber := d.jetKeeper.TopSyncPulse()

	inslogger.FromContext(ctx).Info("db rollback starts. topSyncPulse: ", lastSyncPulseNumber)

	nextPulse := lastSyncPulseNumber + 1

	for idx, db := range d.dbs {
		err := db.TruncateHead(ctx, nextPulse)
		if err != nil {
			return errors.Wrapf(err, "can't truncate %d db since pulse: %d", idx, nextPulse)
		}

		if indexDB, ok := db.(object.IndexModifier); ok {
			if err := indexDB.UpdateLastKnownPulse(ctx, lastSyncPulseNumber); err != nil {
				return errors.Wrap(err, "can't update last sync pulse")
			}
		}
	}

	return nil
}
