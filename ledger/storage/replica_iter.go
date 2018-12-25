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
	"bytes"
	"context"
	"errors"

	"github.com/dgraph-io/badger"

	"github.com/insolar/insolar/core"
)

// iterstate stores iterator state
type iterstate struct {
	prefix []byte
	start  []byte
	end    []byte
}

// ReplicaIter provides partial iterator over BadgerDB key/value pairs
// required for replication to Heavy Material node in provided pulses range.
//
// "Required KV pairs" are all keys with namespace 'scopeIDRecord' (TODO: 'add scopeIDBlob')
// in provided pulses range and all indexes from zero pulse to the end of provided range.
//
// "Partial" means it fetches data in chunks of the specified size.
// After a chunk has been fetched, an iterator saves current position.
//
// NOTE: This is not an "honest" alogrithm, because the last record size can exceed the limit.
// Better implementation is for the future work.
type ReplicaIter struct {
	ctx        context.Context
	db         *DB
	limitBytes int
	istates    []*iterstate
	lastpulse  core.PulseNumber
}

// NewReplicaIter creates ReplicaIter what iterates over records on jet,
// required for heavy material replication.
//
// Params 'start' and 'end' defines pulses from which scan should happen,
// and on which it should be stopped, but indexes scan are always started
// from core.FirstPulseNumber.
//
// Param 'limit' sets per message limit.
func NewReplicaIter(
	ctx context.Context,
	db *DB,
	jetID core.RecordID,
	start core.PulseNumber,
	end core.PulseNumber,
	limit int,
) *ReplicaIter {
	newit := func(prefixbyte byte, jet *core.RecordID, start, end core.PulseNumber) *iterstate {
		prefix := []byte{prefixbyte}
		iter := &iterstate{prefix: prefix}
		if jet == nil {
			iter.start = bytes.Join([][]byte{prefix, start.Bytes()}, nil)
			iter.end = bytes.Join([][]byte{prefix, end.Bytes()}, nil)
		} else {
			iter.start = bytes.Join([][]byte{prefix, jet[:], start.Bytes()}, nil)
			iter.end = bytes.Join([][]byte{prefix, jet[:], end.Bytes()}, nil)
		}
		return iter
	}

	return &ReplicaIter{
		ctx:        ctx,
		db:         db,
		limitBytes: limit,
		// record iterators (order matters for heavy node consistency)
		istates: []*iterstate{
			newit(scopeIDRecord, &jetID, start, end),
			newit(scopeIDBlob, &jetID, start, end),
			newit(scopeIDLifeline, &jetID, core.FirstPulseNumber, end),
			newit(scopeIDJetDrop, &jetID, start, end),
		},
	}
}

// NextRecords fetches next part of key value pairs.
func (r *ReplicaIter) NextRecords() ([]core.KV, error) {
	if r.isDone() {
		return nil, ErrReplicatorDone
	}
	fc := &fetchchunk{
		db:    r.db.db,
		limit: r.limitBytes,
	}
	for _, is := range r.istates {
		if is.start == nil {
			continue
		}
		var fetcherr error
		var lastpulse core.PulseNumber
		is.start, lastpulse, fetcherr = fc.fetch(r.ctx, is.prefix, is.start, is.end)
		if fetcherr != nil {
			return nil, fetcherr
		}
		if lastpulse > r.lastpulse {
			r.lastpulse = lastpulse
		}
	}
	return fc.records, nil
}

// LastPulse returns maximum pulse number of returned keys after each fetch.
func (r *ReplicaIter) LastSeenPulse() core.PulseNumber {
	return r.lastpulse
}

// ErrReplicatorDone is returned by an Replicator NextRecords method when the iteration is complete.
var ErrReplicatorDone = errors.New("no more items in iterator")

func (r *ReplicaIter) isDone() bool {
	for _, is := range r.istates {
		if is.start != nil {
			return false
		}
	}
	return true
}

type fetchchunk struct {
	db      *badger.DB
	records []core.KV
	size    int
	limit   int
}

func (fc *fetchchunk) fetch(
	ctx context.Context,
	prefix []byte,
	start []byte,
	end []byte,
) ([]byte, core.PulseNumber, error) {
	if fc.size > fc.limit {
		return start, 0, nil
	}

	var nextstart []byte
	var lastpulse core.PulseNumber
	err := fc.db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()

		for it.Seek(start); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			if item == nil {
				break
			}
			// key prefix < end
			if bytes.Compare(item.Key()[:len(end)], end) != -1 {
				break
			}

			key := item.KeyCopy(nil)
			if fc.size > fc.limit {
				nextstart = key
				// inslogger.FromContext(ctx).Warnf("size > r.limit: %v > %v (nextstart=%v)",
				// 	fc.size, fc.limit, hex.EncodeToString(key))
				return nil
			}

			lastpulse = pulseFromKey(key)
			// fmt.Printf("key: %v (pulse=%v)\n", hex.EncodeToString(key), lastpulse)

			value, err := it.Item().ValueCopy(nil)
			if err != nil {
				return err
			}
			fc.records = append(fc.records, core.KV{K: key, V: value})
			fc.size += len(key) + len(value)
		}
		nextstart = nil
		return nil
	})
	return nextstart, lastpulse, err
}
