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

package object

import (
	"context"
	"crypto/sha256"
	"io/ioutil"
	"math/rand"
	"os"
	"testing"
	"time"

	fuzz "github.com/google/gofuzz"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/insolar/store"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRecordKey(t *testing.T) {
	t.Parallel()

	expectedKey := recordKey(gen.ID())

	rawID := expectedKey.ID()

	actualKey := newRecordKey(rawID)
	require.Equal(t, expectedKey, actualKey)
}

func TestRecordPositionKey(t *testing.T) {
	t.Parallel()

	expectedKey := recordPositionKey{pn: insolar.GenesisPulse.PulseNumber, number: 42}

	rawID := expectedKey.ID()

	actualKey := newRecordPositionKey(insolar.GenesisPulse.PulseNumber, 42)

	actualKeyFromBytes := newRecordPositionKeyFromBytes(rawID)
	require.Equal(t, expectedKey, actualKeyFromBytes)
	require.Equal(t, expectedKey, actualKey)
}

func TestLastKnownRecordPositionKey(t *testing.T) {
	t.Parallel()

	expectedKey := lastKnownRecordPositionKey{pn: insolar.GenesisPulse.PulseNumber}

	rawID := expectedKey.ID()

	actualKey := newLastKnownRecordPositionKey(rawID)
	require.Equal(t, expectedKey, actualKey)
}

func TestRecordStorage_TruncateHead(t *testing.T) {
	t.Parallel()

	ctx := inslogger.TestContext(t)
	tmpdir, err := ioutil.TempDir("", "bdb-test-")
	defer os.RemoveAll(tmpdir)
	assert.NoError(t, err)

	ops := BadgerDefaultOptions(tmpdir)
	dbMock, err := store.NewBadgerDB(ops)
	defer dbMock.Stop(ctx)
	require.NoError(t, err)

	recordStore := NewRecordDB(dbMock)

	numElements := 10

	// it's used for writing pulses in random order to db
	indexes := make([]int, numElements)
	for i := 0; i < numElements; i++ {
		indexes[i] = i
	}
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(indexes), func(i, j int) { indexes[i], indexes[j] = indexes[j], indexes[i] })

	startPulseNumber := insolar.GenesisPulse.PulseNumber
	ids := make([]insolar.ID, numElements)
	for _, idx := range indexes {
		pulse := startPulseNumber + insolar.PulseNumber(idx)
		ids[idx] = *insolar.NewID(pulse, []byte(testutils.RandomString()))

		err := recordStore.Set(ctx, record.Material{JetID: *insolar.NewJetID(uint8(idx), ids[idx].Bytes()), ID: ids[idx]})
		require.NoError(t, err)
	}

	for i := 0; i < numElements; i++ {
		_, err := recordStore.ForID(ctx, ids[i])
		require.NoError(t, err)
	}

	numLeftElements := numElements / 2
	err = recordStore.TruncateHead(ctx, startPulseNumber+insolar.PulseNumber(numLeftElements))
	require.NoError(t, err)

	for i := 0; i < numLeftElements; i++ {
		_, err := recordStore.ForID(ctx, ids[i])
		require.NoError(t, err)

		_, err = recordStore.AtPosition(ids[i].Pulse(), 1)
		require.NoError(t, err)

		_, err = recordStore.LastKnownPosition(ids[i].Pulse())
		require.NoError(t, err)
	}

	for i := numElements - 1; i >= numLeftElements; i-- {
		_, err := recordStore.ForID(ctx, ids[i])
		require.EqualError(t, err, ErrNotFound.Error())

		_, err = recordStore.AtPosition(ids[i].Pulse(), 1)
		require.EqualError(t, err, ErrNotFound.Error())

		_, err = recordStore.LastKnownPosition(ids[i].Pulse())
		require.EqualError(t, err, ErrNotFound.Error())
	}
}

func TestRecordStorage_NewStorageMemory(t *testing.T) {
	t.Parallel()

	recordStorage := NewRecordMemory()
	assert.Equal(t, 0, len(recordStorage.recsStor))
}

func TestRecordStorage_ForID(t *testing.T) {
	t.Parallel()

	ctx := inslogger.TestContext(t)

	id := gen.ID()
	rec := getMaterialRecord()

	t.Run("returns correct record-value", func(t *testing.T) {
		t.Parallel()

		recordStorage := NewRecordMemory()
		recordStorage.recsStor[id] = rec

		resultRec, err := recordStorage.ForID(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, rec, resultRec)
	})

	t.Run("returns error when no record-value for id", func(t *testing.T) {
		t.Parallel()

		recordStorage := NewRecordMemory()
		recordStorage.recsStor[id] = rec

		_, err := recordStorage.ForID(ctx, gen.ID())
		require.Error(t, err)
		assert.Equal(t, ErrNotFound, err)
	})
}

func TestRecordStorage_Set(t *testing.T) {
	t.Parallel()

	ctx := inslogger.TestContext(t)

	t.Run("saves correct record-value", func(t *testing.T) {
		t.Parallel()

		recordStorage := NewRecordMemory()
		rec := getMaterialRecord()
		rec.ID = gen.ID()

		err := recordStorage.SetAtomic(ctx, rec)
		require.NoError(t, err)
		assert.Equal(t, 1, len(recordStorage.recsStor))
		assert.Equal(t, rec, recordStorage.recsStor[rec.ID])
	})

	t.Run("returns override error when saving with the same id", func(t *testing.T) {
		t.Parallel()

		recordStorage := NewRecordMemory()
		rec := getMaterialRecord()
		rec.ID = gen.ID()

		err := recordStorage.SetAtomic(ctx, rec)
		require.NoError(t, err)

		err = recordStorage.SetAtomic(ctx, rec)
		require.Error(t, err)
		assert.Equal(t, ErrOverride, err)
	})

	t.Run("saves multiple records", func(t *testing.T) {
		t.Parallel()

		recordStorage := NewRecordMemory()
		var recs []record.Material
		fuzz.New().NumElements(10, 20).NilChance(0).Funcs(func(r *record.Material, c fuzz.Continue) {
			r.ID = gen.ID()
		}).Fuzz(&recs)
		err := recordStorage.SetAtomic(ctx, recs...)
		require.NoError(t, err)

		for _, r := range recs {
			rec, err := recordStorage.ForID(ctx, r.ID)
			require.NoError(t, err)
			require.Equal(t, rec, r)
		}
	})

	t.Run("override on single record saves none", func(t *testing.T) {
		t.Parallel()

		recordStorage := NewRecordMemory()
		var recs []record.Material
		fuzz.New().NumElements(10, 20).NilChance(0).Funcs(func(r *record.Material, c fuzz.Continue) {
			r.ID = gen.ID()
		}).Fuzz(&recs)

		err := recordStorage.SetAtomic(ctx, recs[0])
		require.NoError(t, err)

		err = recordStorage.SetAtomic(ctx, recs...)
		require.Equal(t, ErrOverride, err)

		for _, r := range recs[1:] {
			_, err := recordStorage.ForID(ctx, r.ID)
			require.Equal(t, ErrNotFound, err)
		}
	})
}

func TestRecordStorage_DB_Set(t *testing.T) {
	t.Parallel()

	ctx := inslogger.TestContext(t)

	t.Run("saves correct record-value", func(t *testing.T) {
		t.Parallel()

		id := gen.ID()
		rec := getMaterialRecord()

		tmpdir, err := ioutil.TempDir("", "bdb-test-")
		defer os.RemoveAll(tmpdir)
		require.NoError(t, err)

		db, err := store.NewBadgerDB(BadgerDefaultOptions(tmpdir))
		require.NoError(t, err)
		defer db.Stop(context.Background())

		recordStorage := NewRecordDB(db)

		rec.ID = id
		err = recordStorage.Set(ctx, rec)
		require.NoError(t, err)
	})

	t.Run("saves correct record-value. batch version", func(t *testing.T) {
		t.Parallel()

		id := gen.ID()
		rec := getMaterialRecord()

		tmpdir, err := ioutil.TempDir("", "bdb-test-")
		defer os.RemoveAll(tmpdir)
		require.NoError(t, err)

		db, err := store.NewBadgerDB(BadgerDefaultOptions(tmpdir))
		require.NoError(t, err)
		defer db.Stop(context.Background())

		recordStorage := NewRecordDB(db)

		rec.ID = id
		err = recordStorage.BatchSet(ctx, []record.Material{rec})
		require.NoError(t, err)
	})

	t.Run("saves correct record-value. batch version. cross-pulses", func(t *testing.T) {
		t.Parallel()

		id := gen.ID()
		sID := *insolar.NewID(id.Pulse()+1, []byte{2})
		tID := *insolar.NewID(id.Pulse()+2, []byte{3})

		rec := getMaterialRecord()
		rec.ID = id
		sRec := getMaterialRecord()
		sRec.ID = sID

		tRec := getMaterialRecord()
		tRec.ID = tID

		tmpdir, err := ioutil.TempDir("", "bdb-test-")
		defer os.RemoveAll(tmpdir)
		require.NoError(t, err)

		db, err := store.NewBadgerDB(BadgerDefaultOptions(tmpdir))
		require.NoError(t, err)
		defer db.Stop(context.Background())

		recordStorage := NewRecordDB(db)

		rec.ID = id
		err = recordStorage.BatchSet(ctx, []record.Material{rec, sRec, tRec})
		require.NoError(t, err)

		position, err := recordStorage.LastKnownPosition(id.Pulse())
		require.NoError(t, err)
		require.Equal(t, uint32(1), position)
		savedID, err := recordStorage.AtPosition(id.Pulse(), position)
		require.NoError(t, err)
		require.Equal(t, id, savedID)

		position, err = recordStorage.LastKnownPosition(sID.Pulse())
		require.NoError(t, err)
		require.Equal(t, uint32(1), position)
		savedID, err = recordStorage.AtPosition(sID.Pulse(), position)
		require.NoError(t, err)
		require.Equal(t, sID, savedID)

		position, err = recordStorage.LastKnownPosition(tID.Pulse())
		require.NoError(t, err)
		require.Equal(t, uint32(1), position)
		savedID, err = recordStorage.AtPosition(tID.Pulse(), position)
		require.NoError(t, err)
		require.Equal(t, tID, savedID)
	})

	t.Run("returns override error when saving with the same id", func(t *testing.T) {
		t.Parallel()

		id := gen.ID()
		rec := getMaterialRecord()

		tmpdir, err := ioutil.TempDir("", "bdb-test-")
		defer os.RemoveAll(tmpdir)
		require.NoError(t, err)

		db, err := store.NewBadgerDB(BadgerDefaultOptions(tmpdir))
		require.NoError(t, err)
		defer db.Stop(context.Background())

		recordStorage := NewRecordDB(db)

		rec.ID = id
		err = recordStorage.Set(ctx, rec)
		require.NoError(t, err)

		err = recordStorage.Set(ctx, rec)
		require.Error(t, err)
		assert.Equal(t, ErrOverride, err)
	})

	t.Run("returns override error when saving with the same id. batch version", func(t *testing.T) {
		t.Parallel()

		id := gen.ID()
		rec := getMaterialRecord()

		tmpdir, err := ioutil.TempDir("", "bdb-test-")
		defer os.RemoveAll(tmpdir)
		require.NoError(t, err)

		db, err := store.NewBadgerDB(BadgerDefaultOptions(tmpdir))
		require.NoError(t, err)
		defer db.Stop(context.Background())

		recordStorage := NewRecordDB(db)

		rec.ID = id

		err = recordStorage.BatchSet(ctx, []record.Material{rec, rec})

		require.Equal(t, ErrOverride, err)
	})
}

func TestRecordStorage_Delete(t *testing.T) {
	t.Parallel()

	ctx := inslogger.TestContext(t)

	firstPulse := gen.PulseNumber()
	secondPulse := firstPulse + 1

	t.Run("delete all records for selected pulse", func(t *testing.T) {
		t.Parallel()

		recordStorage := NewRecordMemory()

		countFirstPulse := rand.Int31n(256)
		countSecondPulse := rand.Int31n(256)

		for i := int32(0); i < countFirstPulse; i++ {
			randID := gen.ID()
			id := insolar.NewID(firstPulse, randID.Hash())
			err := recordStorage.SetAtomic(ctx, record.Material{ID: *id})
			require.NoError(t, err)
		}

		for i := int32(0); i < countSecondPulse; i++ {
			randID := gen.ID()
			id := insolar.NewID(secondPulse, randID.Hash())
			err := recordStorage.SetAtomic(ctx, record.Material{ID: *id})
			require.NoError(t, err)
		}
		require.Equal(t, countFirstPulse+countSecondPulse, int32(len(recordStorage.recsStor)))

		recordStorage.DeleteForPN(ctx, firstPulse)
		require.Equal(t, countSecondPulse, int32(len(recordStorage.recsStor)))
	})
}

func TestRecordStorage_ForPulse(t *testing.T) {
	t.Parallel()

	ctx := inslogger.TestContext(t)
	recordMemory := NewRecordMemory()

	searchJetID := gen.JetID()
	searchPN := gen.PulseNumber()

	searchRecs := map[insolar.ID]struct{}{}
	for i := int32(0); i < rand.Int31n(256); i++ {
		rec := getMaterialRecord()
		rec.JetID = searchJetID

		h := sha256.New()
		hash := record.HashVirtual(h, rec.Virtual)

		rec.ID = *insolar.NewID(searchPN, hash)

		searchRecs[rec.ID] = struct{}{}
		err := recordMemory.SetAtomic(ctx, rec)
		require.NoError(t, err)
	}

	for i := int32(0); i < rand.Int31n(512); i++ {
		rec := getMaterialRecord()

		randID := gen.ID()
		rec.ID = *insolar.NewID(gen.PulseNumber(), randID.Hash())
		err := recordMemory.SetAtomic(ctx, rec)
		require.NoError(t, err)
	}

	res := recordMemory.ForPulse(ctx, searchJetID, searchPN)
	require.Equal(t, len(searchRecs), len(res))

	for _, r := range res {
		h := sha256.New()
		hash := record.HashVirtual(h, r.Virtual)

		rID := insolar.NewID(searchPN, hash)
		_, ok := searchRecs[*rID]
		require.Equal(t, true, ok)
	}
}

func TestRecordPositionDB(t *testing.T) {
	t.Parallel()

	t.Run("Las returns error, when no info", func(t *testing.T) {
		tmpdir, err := ioutil.TempDir("", "bdb-test-")
		defer os.RemoveAll(tmpdir)
		require.NoError(t, err)

		db, err := store.NewBadgerDB(BadgerDefaultOptions(tmpdir))
		require.NoError(t, err)
		defer db.Stop(context.Background())

		recordStorage := NewRecordDB(db)
		pn := gen.PulseNumber()

		position, err := recordStorage.LastKnownPosition(pn)

		require.Error(t, err)
		require.Equal(t, err.Error(), ErrNotFound.Error())
		require.Equal(t, uint32(0), position)
	})

	t.Run("LastKnownPosition works fine", func(t *testing.T) {
		tmpdir, err := ioutil.TempDir("", "bdb-test-")
		defer os.RemoveAll(tmpdir)
		require.NoError(t, err)

		db, err := store.NewBadgerDB(BadgerDefaultOptions(tmpdir))
		require.NoError(t, err)
		defer db.Stop(context.Background())

		recordStorage := NewRecordDB(db)
		pn := gen.PulseNumber()

		id := gen.IDWithPulse(pn)

		err = recordStorage.Set(context.TODO(), record.Material{ID: id})
		require.NoError(t, err)

		next, err := recordStorage.LastKnownPosition(pn)

		require.NoError(t, err)
		require.Equal(t, uint32(1), next)
	})

	t.Run("IncrementPosition works fine", func(t *testing.T) {
		tmpdir, err := ioutil.TempDir("", "bdb-test-")
		defer os.RemoveAll(tmpdir)
		require.NoError(t, err)

		db, err := store.NewBadgerDB(BadgerDefaultOptions(tmpdir))
		require.NoError(t, err)
		defer db.Stop(context.Background())

		recordStorage := NewRecordDB(db)
		pn := gen.PulseNumber()

		id := gen.IDWithPulse(pn)
		sID := gen.IDWithPulse(pn)
		tID := gen.IDWithPulse(pn)

		err = recordStorage.Set(context.TODO(), record.Material{ID: id})
		require.NoError(t, err)
		err = recordStorage.Set(context.TODO(), record.Material{ID: sID})
		require.NoError(t, err)
		err = recordStorage.Set(context.TODO(), record.Material{ID: tID})
		require.NoError(t, err)

		next, err := recordStorage.LastKnownPosition(pn)

		require.NoError(t, err)
		require.Equal(t, uint32(3), next)
	})

	t.Run("IncrementPosition works fine. batch version", func(t *testing.T) {
		tmpdir, err := ioutil.TempDir("", "bdb-test-")
		defer os.RemoveAll(tmpdir)
		require.NoError(t, err)

		db, err := store.NewBadgerDB(BadgerDefaultOptions(tmpdir))
		require.NoError(t, err)
		defer db.Stop(context.Background())

		recordStorage := NewRecordDB(db)
		pn := gen.PulseNumber()

		id := *insolar.NewID(pn, []byte{1})
		sID := *insolar.NewID(pn, []byte{2})
		tID := *insolar.NewID(pn, []byte{3})

		err = recordStorage.BatchSet(
			context.TODO(), []record.Material{
				{ID: id},
				{ID: sID},
				{ID: tID},
			})
		require.NoError(t, err)

		next, err := recordStorage.LastKnownPosition(pn)

		require.NoError(t, err)
		require.Equal(t, uint32(3), next)
	})

	t.Run("AtPosition works fine", func(t *testing.T) {
		tmpdir, err := ioutil.TempDir("", "bdb-test-")
		defer os.RemoveAll(tmpdir)
		require.NoError(t, err)

		db, err := store.NewBadgerDB(BadgerDefaultOptions(tmpdir))
		require.NoError(t, err)
		defer db.Stop(context.Background())

		recordStorage := NewRecordDB(db)
		pn := gen.PulseNumber()

		id := gen.IDWithPulse(pn)
		sID := gen.IDWithPulse(pn)
		tID := gen.IDWithPulse(pn)

		err = recordStorage.Set(context.TODO(), record.Material{ID: id})
		require.NoError(t, err)
		savedID, err := recordStorage.AtPosition(pn, 1)
		require.NoError(t, err)
		require.Equal(t, id, savedID)

		err = recordStorage.Set(context.TODO(), record.Material{ID: sID})
		require.NoError(t, err)
		savedID, err = recordStorage.AtPosition(pn, 2)
		require.NoError(t, err)
		require.Equal(t, sID, savedID)

		err = recordStorage.Set(context.TODO(), record.Material{ID: tID})
		require.NoError(t, err)
		savedID, err = recordStorage.AtPosition(pn, 3)
		require.NoError(t, err)
		require.Equal(t, tID, savedID)
	})

	t.Run("AtPosition returns error, when the passed position is biggest then saved", func(t *testing.T) {
		tmpdir, err := ioutil.TempDir("", "bdb-test-")
		defer os.RemoveAll(tmpdir)
		require.NoError(t, err)

		db, err := store.NewBadgerDB(BadgerDefaultOptions(tmpdir))
		require.NoError(t, err)
		defer db.Stop(context.Background())

		recordStorage := NewRecordDB(db)
		pn := gen.PulseNumber()

		_, err = recordStorage.AtPosition(pn, 1)
		require.Error(t, err)
		require.Equal(t, err, ErrNotFound)

		id := gen.IDWithPulse(pn)

		err = recordStorage.Set(context.TODO(), record.Material{ID: id})
		require.NoError(t, err)
		savedID, err := recordStorage.AtPosition(pn, 1)
		require.NoError(t, err)
		require.Equal(t, id, savedID)
	})
}

// getVirtualRecord generates random Virtual record
func getVirtualRecord() record.Virtual {
	var requestRecord record.IncomingRequest

	obj := gen.Reference()
	requestRecord.Object = &obj

	virtualRecord := record.Virtual{
		Union: &record.Virtual_IncomingRequest{
			IncomingRequest: &requestRecord,
		},
	}

	return virtualRecord
}

// getMaterialRecord generates random Material record
func getMaterialRecord() record.Material {
	virtRec := getVirtualRecord()

	materialRecord := record.Material{
		Virtual: virtRec,
		JetID:   gen.JetID(),
	}

	return materialRecord
}
