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

package storage_test

import (
	"bytes"
	"context"
	"fmt"
	"sort"
	"testing"

	"github.com/dgraph-io/badger"
	"github.com/insolar/insolar/component"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/storage"
	"github.com/insolar/insolar/ledger/storage/index"
	"github.com/insolar/insolar/ledger/storage/jet"
	"github.com/insolar/insolar/ledger/storage/record"
	"github.com/insolar/insolar/ledger/storage/storagetest"
	"github.com/insolar/insolar/platformpolicy"
	"github.com/insolar/insolar/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type replicaIterSuite struct {
	suite.Suite

	cm      *component.Manager
	ctx     context.Context
	cleaner func()
	db      storage.DBContext

	objectStorage storage.ObjectStorage
	dropStorage   storage.DropStorage
}

func NewReplicaIterSuite() *replicaIterSuite {
	return &replicaIterSuite{
		Suite: suite.Suite{},
	}
}

// Init and run suite
func TestReplicaIter(t *testing.T) {
	suite.Run(t, NewReplicaIterSuite())
}

func (s *replicaIterSuite) BeforeTest(suiteName, testName string) {
	s.cm = &component.Manager{}
	s.ctx = inslogger.TestContext(s.T())

	db, cleaner := storagetest.TmpDB(s.ctx, s.T())
	s.db = db
	s.cleaner = cleaner

	s.objectStorage = storage.NewObjectStorage()
	s.dropStorage = storage.NewDropStorage(10)

	s.cm.Inject(
		platformpolicy.NewPlatformCryptographyScheme(),
		s.db,
		s.objectStorage,
		s.dropStorage,
	)

	err := s.cm.Init(s.ctx)
	if err != nil {
		s.T().Error("ComponentManager init failed", err)
	}
	err = s.cm.Start(s.ctx)
	if err != nil {
		s.T().Error("ComponentManager start failed", err)
	}
}

func (s *replicaIterSuite) AfterTest(suiteName, testName string) {
	err := s.cm.Stop(s.ctx)
	if err != nil {
		s.T().Error("ComponentManager stop failed", err)
	}
	s.cleaner()
}

func pulseDelta(n int) core.PulseNumber { return core.PulseNumber(core.FirstPulseNumber + n) }

func Test_StoreKeyValues(t *testing.T) {
	t.Parallel()
	ctx := inslogger.TestContext(t)
	jetID := testutils.RandomJet()
	// fmt.Printf("random jetID: %v\n", jetID.DebugString())

	var (
		expectedrecs []key
		expectedidxs []key
	)
	var allKVs []core.KV
	pulsescount := 3

	func() {
		db, cleaner := storagetest.TmpDB(ctx, t)
		defer cleaner()

		os := storage.NewObjectStorage()
		ds := storage.NewDropStorage(10)

		cm := &component.Manager{}
		cm.Inject(
			platformpolicy.NewPlatformCryptographyScheme(),
			db,
			os,
			ds,
		)
		err := cm.Init(ctx)
		if err != nil {
			t.Error("ComponentManager init failed", err)
		}
		err = cm.Start(ctx)
		if err != nil {
			t.Error("ComponentManager start failed", err)
		}
		defer cm.Stop(ctx)

		for n := 0; n < pulsescount; n++ {
			lastPulse := core.PulseNumber(pulseDelta(n))
			addRecords(ctx, t, os, jetID, lastPulse)
			setDrop(ctx, t, ds, jetID, lastPulse)
		}

		for n := 0; n < pulsescount; n++ {
			start, end := pulseDelta(n), pulseDelta(n+1)
			replicator := storage.NewReplicaIter(ctx, db, jetID, start, end, 99)

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
		nullifyJetInKeys(expectedrecs)
		nullifyJetInKeys(expectedidxs)
		sortkeys(expectedrecs)
		sortkeys(expectedidxs)
	}()

	var (
		gotrecs []key
		gotidxs []key
	)
	func() {
		db, cleaner := storagetest.TmpDB(ctx, t)
		defer cleaner()
		err := db.StoreKeyValues(ctx, allKVs)
		require.NoError(t, err)
		gotrecs, gotidxs = getallkeys(db.GetBadgerDB())
	}()

	assert.Equal(t, len(expectedrecs), len(gotrecs), "records counts are the same after restore")
	assert.Equal(t, len(expectedidxs), len(gotidxs), "indexes count are the same after restore")

	require.Equal(t, expectedrecs, gotrecs, "records are the same after restore")
	require.Equal(t, expectedidxs, gotidxs, "indexes are the same after restore")
}

func (s *replicaIterSuite) Test_ReplicaIter_FirstPulse() {
	// it's easy to test simple case with zero Jet
	jetID := *jet.NewID(0, nil)

	addRecords(s.ctx, s.T(), s.objectStorage, jetID, core.FirstPulseNumber)
	replicator := storage.NewReplicaIter(s.ctx, s.db, jetID, core.FirstPulseNumber, core.FirstPulseNumber+1, 100500)
	var got []key
	for i := 0; ; i++ {
		if i > 50 {
			s.T().Fatal("too many loops")
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
	all, idxs := getallkeys(s.db.GetBadgerDB())
	all = append(all, idxs...)
	all = sortkeys(all)

	require.Equal(s.T(), all, got, "get expected records for first pulse")
}

func Test_ReplicaIter_Base(t *testing.T) {
	ctx := inslogger.TestContext(t)
	db, cleaner := storagetest.TmpDB(ctx, t, storagetest.DisableBootstrap())
	defer cleaner()

	os := storage.NewObjectStorage()
	ds := storage.NewDropStorage(10)

	cm := &component.Manager{}
	cm.Inject(
		platformpolicy.NewPlatformCryptographyScheme(),
		db,
		os,
		ds,
	)
	err := cm.Init(ctx)
	if err != nil {
		t.Error("ComponentManager init failed", err)
	}
	err = cm.Start(ctx)
	if err != nil {
		t.Error("ComponentManager start failed", err)
	}
	defer cm.Stop(ctx)

	var lastPulse core.PulseNumber
	pulsescount := 2
	// it's easy to test simple case with zero Jet
	jetID := *jet.NewID(0, nil)

	recsBefore, idxBefore := getallkeys(db.GetBadgerDB())
	require.Nil(t, recsBefore)
	require.Nil(t, idxBefore)

	ttPerPulse := make(map[int][]key)
	ttRange := make(map[int][]key)

	recsPerPulse := make(map[int][]key)
	for i := 0; i < pulsescount; i++ {
		lastPulse = pulseDelta(i)

		addRecords(ctx, t, os, jetID, lastPulse)
		setDrop(ctx, t, ds, jetID, lastPulse)

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
		replicator := storage.NewReplicaIter(ctx, db, jetID, p, p+1, maxsize)
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
	addRecords(ctx, t, os, jetID, lastPulse)
	for n := 0; n < pulsescount; n++ {
		p := pulseDelta(n)

		replicator := storage.NewReplicaIter(ctx, db, jetID, p, lastPulse, maxsize)
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
	dropStorage storage.DropStorage,
	jetID core.RecordID,
	pulsenum core.PulseNumber,
) {
	prevDrop, err := dropStorage.GetDrop(ctx, jetID, pulsenum-1)
	var prevhash []byte
	if err == nil {
		prevhash = prevDrop.Hash
	} else if err != core.ErrNotFound {
		require.NoError(t, err)
	}
	drop, _, dropSize, err := dropStorage.CreateDrop(ctx, jetID, pulsenum, prevhash)
	if err != nil {
		require.NoError(t, err)
	}
	require.NotEqual(t, 0, dropSize)
	err = dropStorage.SetDrop(ctx, jetID, drop)
	require.NoError(t, err)
}

func addRecords(
	ctx context.Context,
	t *testing.T,
	objectStorage storage.ObjectStorage,
	jetID core.RecordID,
	pulsenum core.PulseNumber,
) {
	// set record
	parentID, err := objectStorage.SetRecord(
		ctx,
		jetID,
		pulsenum,
		&record.ObjectActivateRecord{
			SideEffectRecord: record.SideEffectRecord{
				Domain: testutils.RandomRef(),
			},
		},
	)
	require.NoError(t, err)

	// set blob
	_, err = objectStorage.SetBlob(ctx, jetID, pulsenum, []byte("100500"))
	require.NoError(t, err)

	// set index of record
	err = objectStorage.SetObjectIndex(ctx, jetID, parentID, &index.ObjectLifeline{
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
		pn := storage.Key(k).PulseNumber()
		if pn == 0 {
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

type key storage.Key

func (k key) String() string {
	return storage.Key(k).String()
}

func sortkeys(keys []key) []key {
	sort.Slice(keys, func(i, j int) bool {
		return bytes.Compare(keys[i], keys[j]) < 0
	})
	return keys
}

func printkeys(keys []key, prefix string) {
	for _, k := range keys {
		fmt.Printf("%v%v (%v)\n", prefix, k, storage.Key(k).PulseNumber())
	}
}

func nullifyJetInKeys(keys []key) {
	for _, k := range keys {
		storage.NullifyJetInKey(k)
	}
}
