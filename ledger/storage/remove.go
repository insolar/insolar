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
	"github.com/insolar/insolar/core"
)

var rmScanFromPulse = core.PulseNumber(core.FirstPulseNumber + 1).Bytes()

// RemoveJetIndexesUntil removes for provided JetID all lifelines older than provided pulse number.
func (db *DB) RemoveJetIndexesUntil(ctx context.Context, jetID core.RecordID, pn core.PulseNumber) (int, error) {
	count := 0
	jetprefix := prefixkey(scopeIDLifeline, jetID[:])
	startprefix := prefixkey(scopeIDLifeline, jetID[:], rmScanFromPulse)

	return count, db.db.Update(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()
		for it.Seek(startprefix); it.ValidForPrefix(jetprefix); it.Next() {
			key := it.Item().Key()
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
