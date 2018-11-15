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

	var (
		expectedrecs []string
		expectedidxs []string
	)
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
		expectedrecs, expectedidxs = getallkeys(db.GetBadgerDB())
	}()

	var (
		gotrecs []string
		gotidxs []string
	)
	func() {
		db, cleaner := storagetest.TmpDB(ctx, t, storagetest.DisableBootstrap())
		defer cleaner()
		err := db.StoreKeyValues(ctx, allKVs)
		require.NoError(t, err)
		gotrecs, gotidxs = getallkeys(db.GetBadgerDB())
	}()

	require.Equal(t, expectedrecs, gotrecs)
	require.Equal(t, expectedidxs, gotidxs)
}

func Test_ReplicaIter(t *testing.T) {
	t.Parallel()
	ctx := inslogger.TestContext(t)
	db, cleaner := storagetest.TmpDB(ctx, t, storagetest.DisableBootstrap())
	defer cleaner()

	recsBefore, idxBefore := getallkeys(db.GetBadgerDB())
	require.Nil(t, recsBefore)
	require.Nil(t, idxBefore)

	// TODO: remove assertpulse struct
	// tt represents test case PulseNumber -> expected record keys
	tt := make(map[int][]string)

	recordsPerPulse := make(map[int][]string)
	for i := 1; i < 3; i++ {
		addRecords(ctx, t, db, core.PulseNumber(i))
		recs, _ := getallkeys(db.GetBadgerDB())
		recKeys := getdelta(recsBefore, recs)
		recsBefore = recs

		recordsPerPulse[i] = recKeys
	}
	_, idxsAfter := getallkeys(db.GetBadgerDB())
	for i := 1; i < 3; i++ {
		for j := i; j < 3; j++ {
			tt[i] = append(tt[i], recordsPerPulse[j]...)
		}
		tt[i] = append(tt[i], idxsAfter...)
	}

	// BEWARE: test expects limit 512 is enougth to have at least `atLeastIterations` iterations
	// it could be fragile, probably I should figure out how to write this test in more stable way.
	// (now there is no so much time for that)
	maxsize := 512
	atLeastIterations := 2

	for n := 1; n < 3; n++ {
		fmt.Println("=================== Pulse:", n, " ====================")
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

		assert.Truef(t, iterations >= atLeastIterations,
			"expect at least %v iterations", atLeastIterations)

		sort.Strings(tt[n])
		sort.Strings(got)
		require.Equal(t, tt[n], got)
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

func getdelta(before []string, after []string) (delta []string) {
CHECKIFCONTAINS:
	for _, k1 := range after {
		for _, k2 := range before {
			if k1 == k2 {
				continue CHECKIFCONTAINS
			}
		}
		// not found
		delta = append(delta, k1)
	}
	return
}

// strip namesapce
func getallkeys(db *badger.DB) (records []string, indexes []string) {
	txn := db.NewTransaction(true)
	defer txn.Discard()

	it := txn.NewIterator(badger.DefaultIteratorOptions)
	defer it.Close()
	for it.Rewind(); it.Valid(); it.Next() {
		item := it.Item()
		k := item.KeyCopy(nil)
		kstr := hex.EncodeToString(k)
		switch k[0] {
		case scopeIDRecord:
			records = append(records, kstr)
		case scopeIDLifeline:
			indexes = append(indexes, kstr)
		}
	}
	return
}
