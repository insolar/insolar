/*
 *    Copyright 2019 Insolar Technologies
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

package storage

import (
	"context"

	"github.com/dgraph-io/badger"
	multierror "github.com/hashicorp/go-multierror"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/instrumentation/insmetrics"
	"github.com/insolar/insolar/ledger/recentstorage"
	"github.com/insolar/insolar/ledger/storage/jet"
	"github.com/pkg/errors"
	"go.opencensus.io/stats"
)

// Cleaner cleans lights after sync to heavy
//go:generate minimock -i github.com/insolar/insolar/ledger/storage.Cleaner -o ./ -s _mock.go
type Cleaner interface {
	RemoveAllForJetUntilPulse(
		ctx context.Context,
		jetID core.RecordID,
		pn core.PulseNumber,
		recent recentstorage.RecentStorage,
	) (map[string]RmStat, error)

	RemoveJetIndexesUntil(ctx context.Context, jetID core.RecordID, pn core.PulseNumber, recent recentstorage.RecentStorage) (RmStat, error)
	RemoveJetBlobsUntil(ctx context.Context, jetID core.RecordID, pn core.PulseNumber) (RmStat, error)
	RemoveJetRecordsUntil(ctx context.Context, jetID core.RecordID, pn core.PulseNumber) (RmStat, error)
	RemoveJetDropsUntil(ctx context.Context, jetID core.RecordID, pn core.PulseNumber) (RmStat, error)
}

var rmScanFromPulse = core.PulseNumber(core.FirstPulseNumber + 1).Bytes()

type RmStat struct {
	Scanned int64
	Removed int64
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

// RemoveAllForJetUntilPulse removes all syncing on heavy records until pulse number for provided jetID
// returns removal stat and cummulative error
func (db *DB) RemoveAllForJetUntilPulse(
	ctx context.Context,
	jetID core.RecordID,
	pn core.PulseNumber,
	recent recentstorage.RecentStorage,
) (map[string]RmStat, error) {
	allstat := map[string]RmStat{}
	var result error

	var err error
	var stat RmStat
	if stat, err = db.RemoveJetIndexesUntil(ctx, jetID, pn, recent); err != nil {
		result = multierror.Append(result, errors.Wrap(err, "RemoveJetIndexesUntil"))
	}
	allstat["indexes"] = stat
	if stat, err = db.RemoveJetBlobsUntil(ctx, jetID, pn); err != nil {
		result = multierror.Append(result, errors.Wrap(err, "RemoveJetBlobsUntil"))
	}
	allstat["blobs"] = stat
	if stat, err = db.RemoveJetRecordsUntil(ctx, jetID, pn); err != nil {
		result = multierror.Append(result, errors.Wrap(err, "RemoveJetRecordsUntil"))
	}
	allstat["records"] = stat
	if stat, err = db.RemoveJetDropsUntil(ctx, jetID, pn); err != nil {
		result = multierror.Append(result, errors.Wrap(err, "RemoveJetDropsUntil"))
	}
	allstat["drops"] = stat
	recordCleanupMetrics(ctx, allstat)

	return allstat, result
}

// RemoveJetIndexesUntil removes for provided JetID all lifelines older than provided pulse number.
// Indexes caches by recent storage, we should avoid them deletion.
func (db *DB) RemoveJetIndexesUntil(ctx context.Context, jetID core.RecordID, pn core.PulseNumber, recent recentstorage.RecentStorage) (RmStat, error) {
	return db.removeJetRecordsUntil(ctx, scopeIDLifeline, jetID, pn, recent)
}

// RemoveJetBlobsUntil removes for provided JetID all blobs older than provided pulse number.
func (db *DB) RemoveJetBlobsUntil(ctx context.Context, jetID core.RecordID, pn core.PulseNumber) (RmStat, error) {
	return db.removeJetRecordsUntil(ctx, scopeIDBlob, jetID, pn, nil)
}

// RemoveJetRecordsUntil removes for provided JetID all records older than provided pulse number.
// In recods pending requests live, so we need recent storage here
func (db *DB) RemoveJetRecordsUntil(ctx context.Context, jetID core.RecordID, pn core.PulseNumber) (RmStat, error) {
	return db.removeJetRecordsUntil(ctx, scopeIDRecord, jetID, pn, nil)
}

// RemoveJetDropsUntil removes for provided JetID all jet drops older than provided pulse number.
func (db *DB) RemoveJetDropsUntil(ctx context.Context, jetID core.RecordID, pn core.PulseNumber) (RmStat, error) {
	return db.removeJetRecordsUntil(ctx, scopeIDJetDrop, jetID, pn, nil)
}

func (db *DB) removeJetRecordsUntil(
	ctx context.Context,
	namespace byte,
	jetID core.RecordID,
	pn core.PulseNumber,
	recent recentstorage.RecentStorage,
) (RmStat, error) {
	var stat RmStat
	_, prefix := jet.Jet(jetID)
	jetprefix := prefixkey(namespace, prefix)
	startprefix := prefixkey(namespace, prefix, rmScanFromPulse)

	return stat, db.db.Update(func(txn *badger.Txn) error {
		var id core.RecordID
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 0
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Seek(startprefix); it.ValidForPrefix(jetprefix); it.Next() {
			key := it.Item().KeyCopy(nil)
			if pulseFromKey(key) >= pn {
				break
			}
			stat.Scanned++

			if recent != nil {
				copy(id[:], key[len(jetprefix):])
				if recent.IsRecordIDCached(id) {
					continue
				}
			}

			if err := txn.Delete(key); err != nil {
				return err
			}
			stat.Removed++
		}
		return nil
	})
}
