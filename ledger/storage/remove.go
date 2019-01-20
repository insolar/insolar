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

package storage

import (
	"context"

	"github.com/dgraph-io/badger"
	"github.com/hashicorp/go-multierror"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/ledger/storage/jet"
)

var rmScanFromPulse = core.PulseNumber(core.FirstPulseNumber + 1).Bytes()

// RemoveAllForJetUntilPulse removes all syncing on heavy records until pulse number for provided jetID
// returns removal stat and cummulative error
func (db *DB) RemoveAllForJetUntilPulse(ctx context.Context, jetID core.RecordID, pn core.PulseNumber) (map[string]int, error) {
	stat := map[string]int{}
	var result error

	var err error
	var removed int
	if removed, err = db.RemoveJetIndexesUntil(ctx, jetID, pn); err != nil {
		result = multierror.Append(result, err)
	}
	stat["indexes"] = removed
	if removed, err = db.RemoveJetBlobsUntil(ctx, jetID, pn); err != nil {
		result = multierror.Append(result, err)
	}
	stat["blobs"] = removed
	if removed, err = db.RemoveJetRecordsUntil(ctx, jetID, pn); err != nil {
		result = multierror.Append(result, err)
	}
	stat["records"] = removed
	if removed, err = db.RemoveJetDropsUntil(ctx, jetID, pn); err != nil {
		result = multierror.Append(result, err)
	}
	stat["drops"] = removed

	return stat, result
}

// RemoveJetIndexesUntil removes for provided JetID all lifelines older than provided pulse number.
func (db *DB) RemoveJetIndexesUntil(ctx context.Context, jetID core.RecordID, pn core.PulseNumber) (int, error) {
	return db.removeJetRecordsUntil(ctx, scopeIDLifeline, jetID, pn)
}

// RemoveJetBlobsUntil removes for provided JetID all blobs older than provided pulse number.
func (db *DB) RemoveJetBlobsUntil(ctx context.Context, jetID core.RecordID, pn core.PulseNumber) (int, error) {
	return db.removeJetRecordsUntil(ctx, scopeIDBlob, jetID, pn)
}

// RemoveJetRecordsUntil removes for provided JetID all records older than provided pulse number.
func (db *DB) RemoveJetRecordsUntil(ctx context.Context, jetID core.RecordID, pn core.PulseNumber) (int, error) {
	return db.removeJetRecordsUntil(ctx, scopeIDRecord, jetID, pn)
}

// RemoveJetDropsUntil removes for provided JetID all jet drops older than provided pulse number.
func (db *DB) RemoveJetDropsUntil(ctx context.Context, jetID core.RecordID, pn core.PulseNumber) (int, error) {
	return db.removeJetRecordsUntil(ctx, scopeIDJetDrop, jetID, pn)
}

func (db *DB) removeJetRecordsUntil(ctx context.Context, namespace byte, jetID core.RecordID, pn core.PulseNumber) (int, error) {
	_, prefix := jet.Jet(jetID)
	jetprefix := prefixkey(namespace, prefix)
	startprefix := prefixkey(namespace, prefix, rmScanFromPulse)

	count := 0
	return count, db.db.Update(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()
		for it.Seek(startprefix); it.ValidForPrefix(jetprefix); it.Next() {
			key := it.Item().KeyCopy(nil)
			if pulseFromKey(key) >= pn {
				break
			}
			if err := txn.Delete(key); err != nil {
				return err
			}
			count++
		}
		return nil
	})
}
