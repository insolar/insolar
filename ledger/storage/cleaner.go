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

package storage

import (
	"context"

	"github.com/dgraph-io/badger"
	multierror "github.com/hashicorp/go-multierror"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/instrumentation/insmetrics"
	"github.com/insolar/insolar/ledger/recentstorage"
	"github.com/pkg/errors"
	"go.opencensus.io/stats"
)

// Cleaner cleans lights after sync to heavy
//go:generate minimock -i github.com/insolar/insolar/ledger/storage.Cleaner -o ./ -s _mock.go
type Cleaner interface {
	CleanJetRecordsUntilPulse(
		ctx context.Context,
		jetID insolar.ID,
		pn insolar.PulseNumber,
	) (map[string]RmStat, error)

	CleanJetIndexes(
		ctx context.Context,
		jetID insolar.ID,
		recent recentstorage.RecentIndexStorage,
		candidates []insolar.ID,
	) (RmStat, error)
}

type cleaner struct {
	DB DBContext `inject:""`
}

// NewCleaner is a constructor for Cleaner.
func NewCleaner() Cleaner {
	return new(cleaner)
}

var rmScanFromPulse = insolar.PulseNumber(insolar.FirstPulseNumber + 1).Bytes()

// RmStat holds removal statistics
type RmStat struct {
	Scanned int64
	Removed int64
	Errors  int64
}

func recordCleanupMetrics(ctx context.Context, stat map[string]RmStat) {
	for name, value := range stat {
		mctx := insmetrics.InsertTag(ctx, recordType, name)
		stats.Record(mctx,
			statCleanScanned.M(value.Scanned),
			statCleanRemoved.M(value.Removed),
		)
	}
}

// CleanJetRecordsUntilPulse removes all records synced on heavy, except indexes until pn pulse number for jetID.
//
// Returns removal statistics and cummulative error of sub cleanup methods.
func (c *cleaner) CleanJetRecordsUntilPulse(
	ctx context.Context,
	jetID insolar.ID,
	pn insolar.PulseNumber,
) (map[string]RmStat, error) {
	allstat := map[string]RmStat{}
	var result error

	var err error
	var stat RmStat
	if stat, err = c.RemoveJetBlobsUntil(ctx, jetID, pn); err != nil {
		result = multierror.Append(result, errors.Wrap(err, "RemoveJetBlobsUntil"))
		stat.Errors = stat.Scanned
		stat.Removed = 0
	}
	allstat["blobs"] = stat

	if stat, err = c.RemoveJetRecordsUntil(ctx, jetID, pn); err != nil {
		result = multierror.Append(result, errors.Wrap(err, "RemoveJetRecordsUntil"))
		stat.Errors = stat.Scanned
		stat.Removed = 0
	}
	allstat["records"] = stat

	recordCleanupMetrics(ctx, allstat)

	return allstat, result
}

// RemoveJetBlobsUntil removes for provided JetID all blobs older than provided pulse number.
func (c *cleaner) RemoveJetBlobsUntil(ctx context.Context, jetID insolar.ID, pn insolar.PulseNumber) (RmStat, error) {
	return c.removeJetRecordsUntil(ctx, scopeIDBlob, jetID, pn)
}

// RemoveJetRecordsUntil removes for provided JetID all records older than provided pulse number.
// In recods pending requests live, so we need recent storage here
func (c *cleaner) RemoveJetRecordsUntil(ctx context.Context, jetID insolar.ID, pn insolar.PulseNumber) (RmStat, error) {
	return c.removeJetRecordsUntil(ctx, scopeIDRecord, jetID, pn)
}

func (c *cleaner) removeJetRecordsUntil(
	ctx context.Context,
	namespace byte,
	jetID insolar.ID,
	pn insolar.PulseNumber,
) (RmStat, error) {
	var stat RmStat
	prefix := insolar.JetID(jetID).Prefix()
	jetprefix := prefixkey(namespace, prefix)
	startprefix := prefixkey(namespace, prefix, rmScanFromPulse)

	return stat, c.DB.GetBadgerDB().Update(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 0
		opts.PrefetchValues = false
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Seek(startprefix); it.ValidForPrefix(jetprefix); it.Next() {
			key := it.Item().KeyCopy(nil)
			if pulseFromKey(key) >= pn {
				break
			}
			stat.Scanned++

			if err := txn.Delete(key); err != nil {
				return err
			}
			stat.Removed++
		}
		return nil
	})
}

// CleanJetIndexes removes indexes from candidates list,
// call recent storage is list still valid.
//
// It locks recent jet store while works.
func (c *cleaner) CleanJetIndexes(
	ctx context.Context,
	jetID insolar.ID,
	recent recentstorage.RecentIndexStorage,
	candidates []insolar.ID,
) (RmStat, error) {
	var stat RmStat
	prefix := insolar.JetID(jetID).Prefix()

	recent.FilterNotExistWithLock(ctx, candidates, func(fordelete []insolar.ID) {
		for _, recID := range fordelete {
			stat.Scanned++
			key := prefixkey(scopeIDLifeline, prefix, recID[:])
			err := c.DB.GetBadgerDB().Update(func(txn *badger.Txn) error {
				return txn.Delete(key)
			})
			if err != nil {
				stat.Errors++
			} else {
				stat.Removed++
			}
		}
	})

	mctx := insmetrics.InsertTag(ctx, recordType, "indexes")
	stats.Record(mctx,
		statCleanScanned.M(stat.Scanned),
		statCleanRemoved.M(stat.Removed),
		statCleanFailed.M(stat.Errors),
	)
	return stat, nil
}
