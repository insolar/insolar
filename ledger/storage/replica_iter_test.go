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

func pulseDelta(n int) core.PulseNumber { return core.PulseNumber(core.FirstPulseNumber + n) }

func Test_StoreKeyValues(t *testing.T) {
	t.Parallel()
	ctx := inslogger.TestContext(t)

	var (
		expectedrecs []key
		expectedidxs []key
	)
	var allKVs []core.KV
	pulsescount := 3

	func() {
		db, cleaner := storagetest.TmpDB(ctx, t, storagetest.DisableBootstrap())
		defer cleaner()
		for n := 0; n < pulsescount; n++ {
			lastPulse := core.PulseNumber(pulseDelta(n))
			addRecords(ctx, t, db, lastPulse)
			setDrop(ctx, t, db, lastPulse)
		}

		for n := 0; n < pulsescount; n++ {
			start, end := pulseDelta(n), pulseDelta(n+1)
			replicator := storage.NewReplicaIter(ctx, db, start, end, 99)

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
		gotrecs []key
		gotidxs []key
	)
	func() {
		db, cleaner := storagetest.TmpDB(ctx, t, storagetest.DisableBootstrap())
		defer cleaner()
		err := db.StoreKeyValues(ctx, allKVs)
		require.NoError(t, err)
		gotrecs, gotidxs = getallkeys(db.GetBadgerDB())
	}()

	require.Equal(t, expectedrecs, gotrecs, "records are the same after restore")
	require.Equal(t, expectedidxs, gotidxs, "indexes are the same after restore")
}

func Test_ReplicaIter_FirstPulse(t *testing.T) {
	t.Parallel()
	ctx := inslogger.TestContext(t)
	db, cleaner := storagetest.TmpDB(ctx, t)
	defer cleaner()

	addRecords(ctx, t, db, core.FirstPulseNumber)

	replicator := storage.NewReplicaIter(ctx, db, core.FirstPulseNumber, core.FirstPulseNumber+1, 100500)
	var got []key
	for i := 0; ; i++ {
		if i > 50 {
			t.Fatal("too many loops")
		}

		recs, err := replicator.NextRecords()
		if err == storage.ErrReplicatorDone {
			break
		}
		if err != nil {
			panic(err)
		}

		for _, rec := range recs {
			got = append(got, rec.K)
		}
	}

	got = sortkeys(got)
	all, idxs := getallkeys(db.GetBadgerDB())
	all = append(all, idxs...)
	all = sortkeys(all)

	// fmt.Println("All:")
	// printkeys(all, "  ")
	// fmt.Println("Got:")
	// printkeys(got, "  ")

	require.Equal(t, all, got, "get expected records for first pulse")
}

func Test_ReplicaIter_Base(t *testing.T) {
	t.Parallel()
	ctx := inslogger.TestContext(t)
	db, cleaner := storagetest.TmpDB(ctx, t, storagetest.DisableBootstrap())
	defer cleaner()

	var lastPulse core.PulseNumber
	pulsescount := 2

	recsBefore, idxBefore := getallkeys(db.GetBadgerDB())
	require.Nil(t, recsBefore)
	require.Nil(t, idxBefore)

	ttPerPulse := make(map[int][]key)
	ttRange := make(map[int][]key)

	recsPerPulse := make(map[int][]key)
	for i := 0; i < pulsescount; i++ {
		lastPulse = pulseDelta(i)

		addRecords(ctx, t, db, lastPulse)
		setDrop(ctx, t, db, lastPulse)

		recs, _ := getallkeys(db.GetBadgerDB())
		recKeys := getdelta(recsBefore, recs)
		recsBefore = recs

		_, idxAll := getallkeys(db.GetBadgerDB())

		recsPerPulse[i] = recKeys
		ttPerPulse[i] = append(ttPerPulse[i], recKeys...)
		ttPerPulse[i] = append(ttPerPulse[i], idxAll...)
	}
	_, idxsAfter := getallkeys(db.GetBadgerDB())

	for i := 0; i < pulsescount; i++ {
		// in range should be all record from the next pulses
		for j := i; j < pulsescount; j++ {
			ttRange[i] = append(ttRange[i], recsPerPulse[j]...)
		}
		// and all current indexes
		ttRange[i] = append(ttRange[i], idxsAfter...)
	}

	// BEWARE: test expects limit 100is enough to have at least `atLeastIterations` iterations
	maxsize := 100
	atLeastIterations := 2

	for n := 0; n < pulsescount; n++ {
		p := pulseDelta(n)
		replicator := storage.NewReplicaIter(ctx, db, p, p+1, maxsize)
		var got []key

		iterations := 1
		for ; ; iterations++ {
			if iterations > 500 {
				t.Fatal("too many loops")
			}

			recs, err := replicator.NextRecords()
			if err == storage.ErrReplicatorDone {
				break
			}
			if err != nil {
				panic(err)
			}

			for _, rec := range recs {
				got = append(got, rec.K)
			}
		}

		assert.Truef(t, iterations >= atLeastIterations,
			"expect at least %v iterations", atLeastIterations)

		ttPerPulse[n] = sortkeys(ttPerPulse[n])
		got = sortkeys(got)
		require.Equalf(t, ttPerPulse[n], got, "get expected records at pulse %v", p)
	}

	lastPulse = lastPulse + 1
	// addRecords here is for purpose:
	// new records on +1 pulse should not affect iterator result on previous pulse range
	addRecords(ctx, t, db, lastPulse)
	for n := 0; n < pulsescount; n++ {
		p := pulseDelta(n)

		replicator := storage.NewReplicaIter(ctx, db, p, lastPulse, maxsize)
		var got []key
		for {
			recs, err := replicator.NextRecords()
			if err == storage.ErrReplicatorDone {
				break
			}
			if err != nil {
				panic(err)
			}
			for _, rec := range recs {
				got = append(got, rec.K)
			}
		}

		got = sortkeys(got)
		ttRange[n] = sortkeys(ttRange[n])

		require.Equalf(t, ttRange[n], got,
			"get expected records in pulse range [%v:%v]", p, lastPulse)
	}
}

func setDrop(
	ctx context.Context,
	t *testing.T,
	db *storage.DB,
	pulsenum core.PulseNumber,
) {
	prevDrop, err := db.GetDrop(ctx, pulsenum-1)
	var prevhash []byte
	if err == nil {
		prevhash = prevDrop.Hash
	} else if err != storage.ErrNotFound {
		require.NoError(t, err)
	}
	drop, _, err := db.CreateDrop(ctx, pulsenum, prevhash)
	if err != nil {
		require.NoError(t, err)
	}
	err = db.SetDrop(ctx, drop)
	require.NoError(t, err)
}

func addRecords(
	ctx context.Context,
	t *testing.T,
	db *storage.DB,
	pulsenum core.PulseNumber,
) {
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
	require.NoError(t, err)

	// set blob
	_, err = db.SetBlob(ctx, pulsenum, []byte("100500"))
	require.NoError(t, err)

	// set index of record
	err = db.SetObjectIndex(ctx, parentID, &index.ObjectLifeline{
		LatestState: parentID,
	})
	require.NoError(t, err)

	return
}

func getdelta(before []key, after []key) (delta []key) {
CHECKIFCONTAINS:
	for _, k1 := range after {
		for _, k2 := range before {
			if bytes.Compare(k1, k2) == 0 {
				continue CHECKIFCONTAINS
			}
		}
		// not found
		delta = append(delta, k1)
	}
	return
}

var (
	scopeIDLifeline = byte(1)
	scopeIDRecord   = byte(2)
	scopeIDJetDrop  = byte(3)
	scopeIDBlob     = byte(7)
)

func getallkeys(db *badger.DB) (records []key, indexes []key) {
	txn := db.NewTransaction(true)
	defer txn.Discard()

	it := txn.NewIterator(badger.DefaultIteratorOptions)
	defer it.Close()
	for it.Rewind(); it.Valid(); it.Next() {
		item := it.Item()
		k := item.KeyCopy(nil)
		if key(k).pulse() == 0 {
			continue
		}

		switch k[0] {
		case scopeIDRecord:
			records = append(records, k)
		case scopeIDBlob:
			records = append(records, k)
		case scopeIDJetDrop:
			records = append(records, k)
		case scopeIDLifeline:
			indexes = append(indexes, k)
		}
	}
	return
}

type key []byte

func (b key) pulse() core.PulseNumber {
	return core.NewPulseNumber(b[1 : 1+core.PulseNumberSize])
}

func (b key) String() string {
	return hex.EncodeToString(b)
}

func sortkeys(keys []key) []key {
	sort.Slice(keys, func(i, j int) bool {
		return bytes.Compare(keys[i], keys[j]) < 0
	})
	return keys
}

func printkeys(keys []key, prefix string) {
	for _, k := range keys {
		fmt.Printf("%v%v (%v)\n", prefix, k, k.pulse())
	}
}
