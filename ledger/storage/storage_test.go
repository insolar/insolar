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
	"sort"
	"testing"

	"github.com/insolar/insolar/core"
	base58 "github.com/jbenet/go-base58"
	"github.com/stretchr/testify/assert"

	"github.com/insolar/insolar/ledger/index"
	"github.com/insolar/insolar/ledger/jetdrop"
	"github.com/insolar/insolar/ledger/record"
	"github.com/insolar/insolar/ledger/storage"

	"github.com/insolar/insolar/ledger/storage/storagetest"
)

func TestDB_GetRecordNotFound(t *testing.T) {
	t.Parallel()
	db, cleaner := storagetest.TmpDB(t, "")
	defer cleaner()

	rec, err := db.GetRecord(&record.ID{})
	assert.Equal(t, err, storage.ErrNotFound)
	assert.Nil(t, rec)
}

func TestDB_SetRecord(t *testing.T) {
	t.Parallel()
	db, cleaner := storagetest.TmpDB(t, "")
	defer cleaner()

	rec := &record.CallRequest{}
	gotRef, err := db.SetRecord(rec)
	assert.Nil(t, err)

	gotRec, err := db.GetRecord(gotRef)
	assert.Nil(t, err)
	assert.Equal(t, rec, gotRec)

	_, err = db.SetRecord(rec)
	assert.Equalf(t, err, storage.ErrOverride, "records override should be forbidden")
}

func TestDB_GetClassIndex_ReturnsNotFoundIfNoIndex(t *testing.T) {
	t.Parallel()
	db, cleaner := storagetest.TmpDB(t, "")
	defer cleaner()

	idx, err := db.GetClassIndex(&record.ID{Pulse: 1}, false)
	assert.Equal(t, err, storage.ErrNotFound)
	assert.Nil(t, idx)
}

func TestDB_SetClassIndex_StoresCorrectDataInStorage(t *testing.T) {
	t.Parallel()
	db, cleaner := storagetest.TmpDB(t, "")
	defer cleaner()

	idgen := func() record.ID {
		return record.ID{Hash: randhash()}
	}
	latestRef := idgen()
	idx := index.ClassLifeline{
		LatestState: latestRef,
		AmendRefs:   []record.ID{idgen(), idgen(), idgen()},
	}
	zeroID := record.ID{
		Hash: hexhash("122444"),
	}
	err := db.SetClassIndex(&zeroID, &idx)
	assert.Nil(t, err)

	storedIndex, err := db.GetClassIndex(&zeroID, false)
	assert.NoError(t, err)
	assert.Equal(t, *storedIndex, idx)
}

func TestDB_SetObjectIndex_ReturnsNotFoundIfNoIndex(t *testing.T) {
	t.Parallel()
	db, cleaner := storagetest.TmpDB(t, "")
	defer cleaner()

	idx, err := db.GetObjectIndex(&record.ID{Hash: hexhash("5000")}, false)
	assert.Equal(t, storage.ErrNotFound, err)
	assert.Nil(t, idx)
}

func TestDB_SetObjectIndex_StoresCorrectDataInStorage(t *testing.T) {
	t.Parallel()
	db, cleaner := storagetest.TmpDB(t, "")
	defer cleaner()

	idx := index.ObjectLifeline{
		ClassRef:    referenceWithHashes("50", "60"),
		LatestState: record.ID{Hash: hexhash("20")},
	}
	zeroid := record.ID{Hash: hexhash("")}
	err := db.SetObjectIndex(&zeroid, &idx)
	assert.Nil(t, err)

	storedIndex, err := db.GetObjectIndex(&zeroid, false)
	assert.NoError(t, err)
	assert.Equal(t, *storedIndex, idx)
}

func TestDB_GetDrop_ReturnsNotFoundIfNoDrop(t *testing.T) {
	t.Parallel()
	db, cleaner := storagetest.TmpDB(t, "")
	defer cleaner()

	drop, err := db.GetDrop(1)
	assert.Equal(t, err, storage.ErrNotFound)
	assert.Nil(t, drop)
}

func TestDB_CreateDrop(t *testing.T) {
	t.Parallel()
	db, cleaner := storagetest.TmpDB(t, "")
	defer cleaner()

	pulse := core.PulseNumber(core.FirstPulseNumber + 10)
	err := db.AddPulse(core.Pulse{
		PulseNumber: pulse,
		Entropy:     core.Entropy{1, 2, 3},
	})
	records := []record.ObjectActivateRecord{
		{Memory: []byte{1}},
		{Memory: []byte{2}},
		{Memory: []byte{3}},
	}

	var (
		expectedRecData [][2][]byte
		expectedIdxData [][2][]byte
	)
	for _, rec := range records {
		recid, err := db.SetRecord(&rec)
		assert.NoError(t, err)

		idx := &index.ObjectLifeline{LatestState: *recid}
		err = db.SetObjectIndex(recid, idx)
		assert.NoError(t, err)

		idxbytes, _ := index.EncodeObjectLifeline(idx)
		expectedIdxData = append(expectedIdxData, [2][]byte{
			record.ID2Bytes(*recid), idxbytes,
		})
		expectedRecData = append(expectedRecData, [2][]byte{
			record.ID2Bytes(*recid),
			record.MustEncodeRaw(record.MustEncodeToRaw(&rec)),
		})
	}
	sortrecords(true, expectedRecData)
	sortrecords(true, expectedIdxData)

	drop, recData, idxData, err := db.CreateDrop(pulse, []byte{4, 5, 6})
	assert.NoError(t, err)
	assert.Equal(t, pulse, drop.Pulse)
	assert.Equal(t, "23qdSCk6Q1L8m1neGWRFYcZnyRdrWfhhopP3pEu", base58.Encode(drop.Hash))
	assert.Equal(t, expectedRecData, recData)

	genesisPulseBytes := core.PulseNumber(core.FirstPulseNumber).Bytes()
	// filter idxData: remove prefix and skip genesis indexes
	idxData = func() [][2][]byte {
		var out [][2][]byte
		for _, pair := range idxData {
			pair[0] = pair[0][1:]
			if bytes.HasPrefix(pair[0], genesisPulseBytes) {
				continue
			}
			out = append(out, pair)
		}
		return out
	}()

	assert.Equal(t, expectedIdxData, idxData)
}

func sortrecords(ascendant bool, in [][2][]byte) [][2][]byte {
	sort.Slice(in, func(i, j int) bool {
		res := bytes.Compare(in[i][0], in[j][0])
		if ascendant {
			return res < 0
		}
		return res > 0
	})
	return in
}

func TestDB_SetDrop(t *testing.T) {
	t.Parallel()
	db, cleaner := storagetest.TmpDB(t, "")
	defer cleaner()

	drop42 := jetdrop.JetDrop{
		Pulse: 42,
		Hash:  []byte{0xFF},
	}
	err := db.SetDrop(&drop42)
	assert.NoError(t, err)

	got, err := db.GetDrop(42)
	assert.NoError(t, err)
	assert.Equal(t, *got, drop42)
}

func TestDB_AddPulse(t *testing.T) {
	t.Parallel()
	db, cleaner := storagetest.TmpDB(t, "")
	defer cleaner()

	err := db.AddPulse(core.Pulse{PulseNumber: 42, Entropy: core.Entropy{1, 2, 3}})
	assert.NoError(t, err)
	latestPulse, err := db.GetLatestPulseNumber()
	assert.Equal(t, core.PulseNumber(42), latestPulse)
	pulse, err := db.GetPulse(latestPulse)
	assert.NoError(t, err)
	assert.Equal(t, record.PulseRecord{PrevPulse: core.FirstPulseNumber, Entropy: core.Entropy{1, 2, 3}}, *pulse)
}
