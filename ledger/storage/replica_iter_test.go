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

package storage_test

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"sort"
	"testing"

	"github.com/dgraph-io/badger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/index"
	"github.com/insolar/insolar/ledger/record"
	"github.com/insolar/insolar/ledger/storage"
	"github.com/insolar/insolar/ledger/storage/storagetest"
	"github.com/insolar/insolar/testutils"
)

var (
	scopeIDLifeline = byte(1)
	scopeIDRecord   = byte(2)
)

func Test_StoreKeyValues(t *testing.T) {
	t.Parallel()
	ctx := inslogger.TestContext(t)

	var expected []keySize
	var allKVs []core.KV
	pulsescount := 3

	func() {
		db, cleaner := storagetest.TmpDB(ctx, t, storagetest.DisableBootstrap())
		defer cleaner()
		for i := 1; i <= pulsescount; i++ {
			addRecords(ctx, t, db, core.PulseNumber(i))
		}

		for n := 1; n <= pulsescount; n++ {
			replicator := storage.NewReplicaIter(ctx, db, core.PulseNumber(n), 99)

			for i := 0; ; i++ {
				recs, err := replicator.NextRecords()
				if err == storage.ErrReplicatorDone {
					break
				}
				if err != nil {
					panic(err)
				}
				allKVs = append(allKVs, recs...)
			}
		}
		expected = getallkeys(db.GetBadgerDB())
	}()

	var got []keySize
	func() {
		db, cleaner := storagetest.TmpDB(ctx, t, storagetest.DisableBootstrap())
		defer cleaner()
		err := db.StoreKeyValues(ctx, allKVs)
		require.NoError(t, err)
		got = getallkeys(db.GetBadgerDB())
	}()

	require.Equal(t, expected, got)
	// fmt.Println("expect:", outputKeySizes(expected))
	// fmt.Println("got:", outputKeySizes(got))
	// fmt.Printf("allKVs: %#v\n", allKVs)
}

func Test_ReplicaIter(t *testing.T) {
	t.Parallel()
	ctx := inslogger.TestContext(t)
	db, cleaner := storagetest.TmpDB(ctx, t, storagetest.DisableBootstrap())
	defer cleaner()

	keysBefore := getallkeys(db.GetBadgerDB())
	require.Nil(t, keysBefore)
	// tt is test cases: PulseNumber -> expected records count
	type assertpulse struct {
		count int
		keys  []string
	}
	tt := make(map[int]*assertpulse)

	var allIndexKeys []string
	createdCount := 0
	indexesCount := 0
	for i := 1; i < 3; i++ {
		recN, idxN := addRecords(ctx, t, db, core.PulseNumber(i))
		indexesCount += idxN
		createdCount += recN + idxN

		currentstate := getallkeys(db.GetBadgerDB())
		recKeys, idxKeys := getdelta(keysBefore, currentstate)
		// _, _ = recKeySizes, idxKeySizes
		keysBefore = currentstate

		allIndexKeys = append(allIndexKeys, idxKeys...)
		tt[i] = &assertpulse{
			count: recN,
			keys:  recKeys,
		}
	}
	for i := 1; i < 3; i++ {
		tt[i].count += indexesCount
		tt[i].keys = append(tt[i].keys, allIndexKeys...)
	}
	keysAfter := getallkeys(db.GetBadgerDB())
	// fmt.Println("keysAfter", outputKeySizes(keysAfter))
	require.Equal(t, createdCount, len(keysAfter))

	// BEWARE: test expects limit 512 is enougth to have at least `atLeastIterations` iterations
	// it could be fragile, probably I should figure out how to write this test in more stable way.
	// (now there is no so much time for that)
	maxsize := 512
	atLeastIterations := 2

	for n := 1; n < 3; n++ {
		// fmt.Println("=================== Pulse:", n, " ====================")
		replicator := storage.NewReplicaIter(ctx, db, core.PulseNumber(n), maxsize)
		var got []string

		iterations := 0
		for i := 0; ; i++ {
			recs, err := replicator.NextRecords()
			if err == storage.ErrReplicatorDone {
				break
			}
			if err != nil {
				panic(err)
			}

			iterations = i + 1
			if i > 5 {
				fmt.Println("~~~~~~~~~~~ BREAK LOOP ~~~~~~~~~~~~~~")
				break
			}

			for _, rec := range recs {
				got = append(got, hex.EncodeToString(rec.K))
			}
		}
		// fmt.Println("pulse:", n, "iterations:", iterations)

		assert.Equal(t, tt[n].count, len(got))
		assert.Truef(t, iterations >= atLeastIterations,
			"expect at least %v iterations", atLeastIterations)

		sort.Strings(tt[n].keys)
		sort.Strings(got)
		require.Equal(t, tt[n].keys, got)
	}
}

func addRecords(
	ctx context.Context,
	t *testing.T,
	db *storage.DB,
	pulsenum core.PulseNumber,
) (records, indexes int) {
	// set record
	parentID, err := db.SetRecord(
		ctx,
		pulsenum,
		&record.ObjectActivateRecord{
			SideEffectRecord: record.SideEffectRecord{
				Domain: testutils.RandomRef(),
			},
		},
	)
	records++

	// set index of record
	err = db.SetObjectIndex(ctx, parentID, &index.ObjectLifeline{
		LatestState: parentID,
	})
	require.NoError(t, err)
	indexes++

	return
}

type keySize struct {
	key  hexbytes
	size int
}

type hexbytes []byte

// String implements Stringer on bytes slice.
func (b hexbytes) String() string {
	return hex.EncodeToString(b)
}

func outputKeySizes(ks []keySize) (s string) {
	s += fmt.Sprintf("Found %v keys:\n", len(ks))
	for _, k := range ks {
		s += fmt.Sprintf("  key=%s (size=%v)\n", k.key, k.size)
	}
	return
}

// func getallindexes(db *badger.DB) (out []keySize) {
// 	keys := getallkeys(db)
// 	for _, k := range keys {
// 		if k.key[0] == scopeIDLifeline {
// 			out = append(out, k)
// 		}
// 	}
// 	return
// }

func getdelta(before []keySize, after []keySize) (
	recs []string,
	idx []string,
) {
CHECKIFCONTAINS:
	for _, k1 := range after {
		for _, k2 := range before {
			if bytes.Equal(k1.key, k2.key) {
				continue CHECKIFCONTAINS
			}
		}
		// not found
		key := k1.key
		if key[0] == scopeIDRecord {
			recs = append(recs, hex.EncodeToString(key))
		}
		if key[0] == scopeIDLifeline {
			idx = append(idx, hex.EncodeToString(key))
		}
	}
	return
}

// strip namesapce
func getallkeys(db *badger.DB) (keys []keySize) {
	txn := db.NewTransaction(true)
	defer txn.Discard()

	it := txn.NewIterator(badger.DefaultIteratorOptions)
	defer it.Close()
	for it.Rewind(); it.Valid(); it.Next() {
		item := it.Item()
		k := item.KeyCopy(nil)
		val, err := item.Value()
		if err != nil {
			panic(err)
		}
		keys = append(keys, keySize{
			key:  k,
			size: len(k) + len(val),
		})
	}
	return
}
