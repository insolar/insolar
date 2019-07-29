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
	"crypto/sha256"
	"io/ioutil"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/internal/ledger/store"
	"github.com/insolar/insolar/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRecordKey(t *testing.T) {
	t.Parallel()

	expectedKey := recordKey(testutils.RandomID())

	rawID := expectedKey.ID()

	actualKey := newRecordKey(rawID)
	require.Equal(t, expectedKey, actualKey)
}

func TestRecordStorage_TruncateHead(t *testing.T) {
	t.Parallel()

	ctx := inslogger.TestContext(t)
	tmpdir, err := ioutil.TempDir("", "bdb-test-")
	defer os.RemoveAll(tmpdir)
	assert.NoError(t, err)

	dbMock, err := store.NewBadgerDB(tmpdir)
	defer dbMock.Stop(ctx)
	require.NoError(t, err)

	recordStore := NewRecordDB(dbMock)

	numElements := 100

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

		recordStore.Set(ctx, ids[idx], record.Material{JetID: *insolar.NewJetID(uint8(idx), nil)})

		for i := 0; i < 5; i++ {
			recordStore.Set(ctx, ids[idx], record.Material{JetID: *insolar.NewJetID(uint8(i), nil)})
		}

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
	}

	for i := numElements - 1; i >= numLeftElements; i-- {
		_, err := recordStore.ForID(ctx, ids[i])
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

	id := gen.ID()
	rec := getMaterialRecord()

	t.Run("saves correct record-value", func(t *testing.T) {
		t.Parallel()

		recordStorage := NewRecordMemory()

		err := recordStorage.Set(ctx, id, rec)
		require.NoError(t, err)
		assert.Equal(t, 1, len(recordStorage.recsStor))
		assert.Equal(t, rec, recordStorage.recsStor[id])
	})

	t.Run("returns override error when saving with the same id", func(t *testing.T) {
		t.Parallel()

		recordStorage := NewRecordMemory()

		err := recordStorage.Set(ctx, id, rec)
		require.NoError(t, err)

		err = recordStorage.Set(ctx, id, rec)
		require.Error(t, err)
		assert.Equal(t, ErrOverride, err)
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
			err := recordStorage.Set(ctx, *id, record.Material{})
			require.NoError(t, err)
		}

		for i := int32(0); i < countSecondPulse; i++ {
			randID := gen.ID()
			id := insolar.NewID(secondPulse, randID.Hash())
			err := recordStorage.Set(ctx, *id, record.Material{})
			require.NoError(t, err)
		}
		assert.Equal(t, countFirstPulse+countSecondPulse, int32(len(recordStorage.recsStor)))

		recordStorage.DeleteForPN(ctx, firstPulse)
		assert.Equal(t, countSecondPulse, int32(len(recordStorage.recsStor)))
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

		id := insolar.NewID(searchPN, hash)

		searchRecs[*id] = struct{}{}
		err := recordMemory.Set(ctx, *id, rec)
		require.NoError(t, err)
	}

	for i := int32(0); i < rand.Int31n(512); i++ {
		rec := getMaterialRecord()

		randID := gen.ID()
		rID := insolar.NewID(gen.PulseNumber(), randID.Hash())
		err := recordMemory.Set(ctx, *rID, rec)
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
