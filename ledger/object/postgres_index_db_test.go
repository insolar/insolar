// +build slowtest

package object

import (
	"bytes"
	"context"
	"math/rand"
	"sort"
	"testing"
	"time"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func truncateIndexAndRecordTables() {
	_, err := getPool().Exec(context.Background(),
		"TRUNCATE last_known_pulse_for_indexes, indexes, records, records_last_position")
	if err != nil {
		panic(err)
	}
}

const indexCount = 5

func TestPostgresIndexDB_DontLooseIndexAfterTruncate(t *testing.T) {
	defer truncateIndexAndRecordTables()

	ctx := inslogger.TestContext(t)

	indexStore := NewPostgresIndexDB(getPool(), nil)

	testPulse := insolar.GenesisPulse.PulseNumber
	nextPulse := testPulse + 1
	bucket := record.Index{}
	bucket.ObjID = gen.ID()

	err := indexStore.SetIndex(ctx, testPulse, bucket)
	require.NoError(t, err)
	_, err = indexStore.ForID(ctx, testPulse, bucket.ObjID)
	require.NoError(t, err)

	err = indexStore.SetIndex(ctx, nextPulse, bucket)
	require.NoError(t, err)

	_, err = indexStore.ForID(ctx, nextPulse, bucket.ObjID)
	require.NoError(t, err)

	err = indexStore.TruncateHead(ctx, nextPulse)
	require.NoError(t, err)

	_, err = indexStore.ForID(ctx, nextPulse, bucket.ObjID)
	require.EqualError(t, err, ErrIndexNotFound.Error())

	// no update such object in that pulse -> try to get last known pulse but it refers to nextPulse
	// , but we Truncate index with that pulse -> couldn't find that object
	_, err = indexStore.ForID(ctx, nextPulse+1, bucket.ObjID)
	require.EqualError(t, err, ErrIndexNotFound.Error())

	err = indexStore.UpdateLastKnownPulse(ctx, testPulse)
	require.NoError(t, err)
	_, err = indexStore.ForID(ctx, testPulse+2, bucket.ObjID)
	require.NoError(t, err)
}

func TestPostgresIndexDB_TruncateHead(t *testing.T) {
	defer truncateIndexAndRecordTables()

	ctx := inslogger.TestContext(t)

	indexStore := NewPostgresIndexDB(getPool(), NewPostgresRecordDB(getPool()))

	numElements := 10

	// it's used for writing pulses in random order to db
	indexes := make([]int, numElements)
	for i := 0; i < numElements; i++ {
		indexes[i] = i
	}
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(indexes), func(i, j int) { indexes[i], indexes[j] = indexes[j], indexes[i] })

	startPulseNumber := insolar.GenesisPulse.PulseNumber
	objects := make([]insolar.ID, numElements)
	for _, idx := range indexes {
		pulse := startPulseNumber + insolar.PulseNumber(idx)
		objects[idx] = gen.ID()

		bucket := record.Index{}

		bucket.ObjID = objects[idx]
		err := indexStore.SetIndex(ctx, pulse, bucket)
		require.NoError(t, err)

		for i := 0; i < indexCount; i++ {
			bucket := record.Index{}

			bucket.ObjID = gen.ID()
			err := indexStore.SetIndex(ctx, pulse, bucket)
			require.NoError(t, err)
		}

	}

	for i := 0; i < numElements; i++ {
		_, err := indexStore.ForID(ctx, startPulseNumber+insolar.PulseNumber(i), objects[i])
		require.NoError(t, err)
	}

	numLeftElements := numElements / 2
	err := indexStore.TruncateHead(ctx, startPulseNumber+insolar.PulseNumber(numLeftElements))
	require.NoError(t, err)

	for i := 0; i < numLeftElements; i++ {
		_, err := indexStore.ForID(ctx, startPulseNumber+insolar.PulseNumber(i), objects[i])
		require.NoError(t, err)
	}

	for i := numElements - 1; i >= numLeftElements; i-- {
		_, err := indexStore.ForID(ctx, startPulseNumber+insolar.PulseNumber(i), objects[i])
		require.EqualError(t, err, ErrIndexNotFound.Error())
	}
}

func TestPostgresDBIndexStorage_ForID(t *testing.T) {
	defer truncateIndexAndRecordTables()

	ctx := inslogger.TestContext(t)

	id := gen.ID()

	t.Run("returns error when no index-value for id", func(t *testing.T) {
		storage := NewPostgresIndexDB(getPool(), NewPostgresRecordDB(getPool()))
		pn := gen.PulseNumber()

		_, err := storage.ForID(ctx, pn, id)

		assert.Equal(t, ErrIndexNotFound, err)
	})
}

func TestPostgresDBIndexStorage_ForPulse(t *testing.T) {
	defer truncateIndexAndRecordTables()

	ctx := inslogger.TestContext(t)
	prevPn := gen.PulseNumber()
	pn := prevPn + 10
	nextPn := pn + 20

	// Sort indexes for proper compare
	// For now badger iterator already sorted by key but this behavior can change
	sortIndexes := func(slice []record.Index) {
		cmp := func(i, j int) bool {
			cmp := bytes.Compare(slice[i].ObjID.Bytes(), slice[j].ObjID.Bytes())
			return cmp < 0
		}
		sort.Slice(slice, cmp)
	}

	t.Run("empty index storage", func(t *testing.T) {
		defer truncateIndexAndRecordTables()
		storage := NewPostgresIndexDB(getPool(), nil)

		indexes, err := storage.ForPulse(ctx, pn)
		require.Error(t, err)
		require.Equal(t, err, ErrIndexNotFound)
		require.Nil(t, indexes)
	})

	t.Run("index storage with couple values", func(t *testing.T) {
		defer truncateIndexAndRecordTables()
		storage := NewPostgresIndexDB(getPool(), nil)

		var indexes []record.Index
		for i := 0; i < indexCount; i++ {
			indexes = append(indexes, record.Index{ObjID: gen.ID()})
			err := storage.SetIndex(ctx, pn, indexes[i])
			require.NoError(t, err)
		}

		realIndexes, err := storage.ForPulse(ctx, pn)
		require.NoError(t, err)
		require.Equal(t, len(indexes), len(realIndexes))

		sortIndexes(realIndexes)
		sortIndexes(indexes)
		for i := 0; i < indexCount; i++ {
			require.Equal(t, indexes[i], realIndexes[i])
		}
	})

	t.Run("index storage with couple values in different pulses", func(t *testing.T) {
		defer truncateIndexAndRecordTables()
		storage := NewPostgresIndexDB(getPool(), nil)

		var indexes []record.Index
		for i := 0; i < indexCount; i++ {
			indexes = append(indexes, record.Index{ObjID: gen.ID()})
			err := storage.SetIndex(ctx, pn, indexes[i])
			require.NoError(t, err)
		}

		// add some values in prev pulse
		for i := 0; i < indexCount; i++ {
			err := storage.SetIndex(ctx, prevPn, record.Index{ObjID: gen.ID()})
			require.NoError(t, err)
		}

		// add some values in next pulse
		for i := 0; i < indexCount; i++ {
			err := storage.SetIndex(ctx, nextPn, record.Index{ObjID: gen.ID()})
			require.NoError(t, err)
		}

		realIndexes, err := storage.ForPulse(ctx, pn)
		require.NoError(t, err)
		require.Equal(t, len(indexes), len(realIndexes))

		sortIndexes(realIndexes)
		sortIndexes(indexes)
		for i := 0; i < indexCount; i++ {
			require.Equal(t, indexes[i], realIndexes[i])
		}
	})
}

func TestPostgresDBIndex_SetBucket(t *testing.T) {
	ctx := inslogger.TestContext(t)
	objID := gen.ID()
	lflID := gen.ID()
	buck := record.Index{
		ObjID: objID,
		Lifeline: record.Lifeline{
			LatestState: &lflID,
		},
	}

	t.Run("saves correct bucket", func(t *testing.T) {
		defer truncateIndexAndRecordTables()

		pn := gen.PulseNumber()
		index := NewPostgresIndexDB(getPool(), NewPostgresRecordDB(getPool()))

		err := index.SetIndex(ctx, pn, buck)
		require.NoError(t, err)

		res, err := index.ForID(ctx, pn, objID)
		require.NoError(t, err)

		idxBuf, _ := buck.Marshal()
		resBuf, _ := res.Marshal()

		assert.Equal(t, idxBuf, resBuf)
	})

	t.Run("re-save works fine", func(t *testing.T) {
		defer truncateIndexAndRecordTables()

		pn := gen.PulseNumber()
		index := NewPostgresIndexDB(getPool(), NewPostgresRecordDB(getPool()))

		err := index.SetIndex(ctx, pn, buck)
		require.NoError(t, err)

		sLlflID := gen.ID()
		sBuck := record.Index{
			ObjID: objID,
			Lifeline: record.Lifeline{
				LatestState: &sLlflID,
			},
		}

		err = index.SetIndex(ctx, pn, sBuck)
		require.NoError(t, err)

		res, err := index.ForID(ctx, pn, objID)
		require.NoError(t, err)

		idxBuf, _ := sBuck.Marshal()
		resBuf, _ := res.Marshal()

		assert.Equal(t, idxBuf, resBuf)
	})
}

func TestPostgresIndexDB_FetchFilament(t *testing.T) {
	defer truncateIndexAndRecordTables()

	ctx := inslogger.TestContext(t)
	recordStorage := NewPostgresRecordDB(getPool())
	index := NewPostgresIndexDB(getPool(), NewPostgresRecordDB(getPool()))

	first := insolar.NewID(1, nil)
	second := insolar.NewID(2, nil)

	firstMeta := *insolar.NewID(11, nil)
	secondMeta := *insolar.NewID(22, nil)

	firstFil := record.PendingFilament{
		RecordID: *first,
	}
	firstFilV := record.Wrap(&firstFil)
	secondFil := record.PendingFilament{
		RecordID:       *second,
		PreviousRecord: first,
	}
	secondFilV := record.Wrap(&secondFil)

	_ = recordStorage.Set(ctx, record.Material{ID: *first})
	_ = recordStorage.Set(ctx, record.Material{ID: *second})
	_ = recordStorage.Set(ctx, record.Material{Virtual: firstFilV, ID: firstMeta})
	_ = recordStorage.Set(ctx, record.Material{Virtual: secondFilV, ID: secondMeta})

	fi := &record.Index{
		PendingRecords: []insolar.ID{firstMeta, secondMeta},
	}

	res, err := index.filament(fi)

	require.NoError(t, err)
	require.Equal(t, 2, len(res))

	require.Equal(t, *first, res[0].RecordID)
	require.Equal(t, firstMeta, res[0].MetaID)

	require.Equal(t, *second, res[1].RecordID)
	require.Equal(t, secondMeta, res[1].MetaID)
}

func TestPostgresIndexDB_NextFilament(t *testing.T) {
	defer truncateIndexAndRecordTables()
	ctx := inslogger.TestContext(t)

	first := insolar.NewID(1, nil)
	firstMeta := *insolar.NewID(11, nil)

	t.Run("previous exists", func(t *testing.T) {
		defer truncateIndexAndRecordTables()
		recordStorage := NewPostgresRecordDB(getPool())
		index := NewPostgresIndexDB(getPool(), NewPostgresRecordDB(getPool()))

		firstFil := record.PendingFilament{
			PreviousRecord: first,
		}
		firstFilV := record.Wrap(&firstFil)

		_ = recordStorage.Set(ctx, record.Material{Virtual: firstFilV, ID: firstMeta})

		fi := &record.Index{
			PendingRecords: []insolar.ID{firstMeta},
		}

		cc, npn, err := index.nextFilament(fi)

		require.NoError(t, err)
		require.Equal(t, true, cc)

		require.Equal(t, insolar.PulseNumber(1), npn)
	})

	t.Run("previous doesn't exist", func(t *testing.T) {
		defer truncateIndexAndRecordTables()
		recordStorage := NewPostgresRecordDB(getPool())
		index := NewPostgresIndexDB(getPool(), NewPostgresRecordDB(getPool()))

		firstFil := record.PendingFilament{}
		firstFilV := record.Wrap(&firstFil)

		_ = recordStorage.Set(ctx, record.Material{Virtual: firstFilV, ID: firstMeta})

		fi := &record.Index{
			PendingRecords: []insolar.ID{firstMeta},
		}

		cc, _, err := index.nextFilament(fi)

		require.NoError(t, err)
		require.Equal(t, false, cc)
	})

	t.Run("doesn't exist", func(t *testing.T) {
		defer truncateIndexAndRecordTables()
		index := NewPostgresIndexDB(getPool(), NewPostgresRecordDB(getPool()))

		fi := &record.Index{
			PendingRecords: []insolar.ID{firstMeta},
		}

		cc, _, err := index.nextFilament(fi)

		require.Error(t, err, ErrNotFound)
		require.Equal(t, false, cc)
	})
}

func TestPostgresIndexDB_Records(t *testing.T) {
	defer truncateIndexAndRecordTables()
	ctx := inslogger.TestContext(t)

	t.Run("returns err, if readUntil > readFrom", func(t *testing.T) {
		defer truncateIndexAndRecordTables()
		index := NewPostgresIndexDB(getPool(), NewPostgresRecordDB(getPool()))

		res, err := index.Records(ctx, 1, 10, insolar.ID{})

		require.Error(t, err)
		require.Nil(t, res)
	})

	t.Run("works fine", func(t *testing.T) {
		defer truncateIndexAndRecordTables()
		index := NewPostgresIndexDB(getPool(), NewPostgresRecordDB(getPool()))
		rms := NewPostgresRecordDB(getPool())

		pn := insolar.PulseNumber(3)
		pnS := insolar.PulseNumber(2)
		pnT := insolar.PulseNumber(1)

		// Records
		idT := insolar.NewID(pnT, nil)
		rT := record.IncomingRequest{Object: insolar.NewReference(gen.ID())}
		rTV := record.Wrap(&rT)
		_ = rms.Set(ctx, record.Material{Virtual: rTV, ID: *idT})

		idS := insolar.NewID(pnS, nil)
		rS := record.IncomingRequest{Object: insolar.NewReference(gen.ID())}
		rSV := record.Wrap(&rS)
		_ = rms.Set(ctx, record.Material{Virtual: rSV, ID: *idS})

		id := insolar.NewID(pn, nil)
		r := record.IncomingRequest{Object: insolar.NewReference(gen.ID())}
		rv := record.Wrap(&r)
		_ = rms.Set(ctx, record.Material{Virtual: rv, ID: *id})

		// Pending filaments
		midT := insolar.NewID(pnT, []byte{1})
		mT := record.PendingFilament{RecordID: *idT}
		mTV := record.Wrap(&mT)
		_ = rms.Set(ctx, record.Material{Virtual: mTV, ID: *midT})

		midS := insolar.NewID(pnS, []byte{1})
		mS := record.PendingFilament{RecordID: *idS, PreviousRecord: midT}
		mSV := record.Wrap(&mS)
		_ = rms.Set(ctx, record.Material{Virtual: mSV, ID: *midS})

		mid := insolar.NewID(pn, []byte{1})
		m := record.PendingFilament{RecordID: *id, PreviousRecord: midS}
		mV := record.Wrap(&m)
		_ = rms.Set(ctx, record.Material{Virtual: mV, ID: *mid})

		objID := gen.ID()

		third := record.Index{ObjID: objID, PendingRecords: []insolar.ID{*midT}}
		second := record.Index{ObjID: objID, PendingRecords: []insolar.ID{*midS}}
		first := record.Index{ObjID: objID, PendingRecords: []insolar.ID{*mid}}

		err := index.SetIndex(ctx, pn, first)
		require.NoError(t, err)
		err = index.SetIndex(ctx, pnS, second)
		require.NoError(t, err)
		err = index.SetIndex(ctx, pnT, third)
		require.NoError(t, err)

		res, err := index.Records(ctx, insolar.PulseNumber(3), insolar.PulseNumber(2), objID)

		require.NoError(t, err)
		require.Equal(t, 2, len(res))

		require.Equal(t, *idS, res[0].RecordID)
		require.Equal(t, *id, res[1].RecordID)

		require.Equal(t, *midS, res[0].MetaID)
		require.Equal(t, *mid, res[1].MetaID)
	})

}

func TestPostgresIndexDB_UpdateLastKnownPulse(t *testing.T) {
	defer truncateIndexAndRecordTables()

	ctx := inslogger.TestContext(t)
	objectID := gen.ID()

	t.Run("insert once and then only updates", func(t *testing.T) {
		defer truncateIndexAndRecordTables()
		storage := NewPostgresIndexDB(getPool(), nil)

		pn := gen.PulseNumber()
		for i := 0; i < indexCount; i++ {
			err := storage.SetIndex(ctx, pn, record.Index{
				ObjID:            objectID,
				LifelineLastUsed: pn,
			})
			require.NoError(t, err)

			err = storage.UpdateLastKnownPulse(ctx, pn)
			require.NoError(t, err)

			index, err := storage.LastKnownForID(ctx, objectID)
			require.NoError(t, err)

			require.Equal(t, pn, index.LifelineLastUsed)

			pn = pn + 10
		}
	})
}
