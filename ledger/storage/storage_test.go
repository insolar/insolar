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
	"testing"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/core/message"
	"github.com/jbenet/go-base58"
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

	rec, err := db.GetRecord(&core.RecordID{})
	assert.Equal(t, err, storage.ErrNotFound)
	assert.Nil(t, rec)
}

func TestDB_SetRecord(t *testing.T) {
	t.Parallel()
	db, cleaner := storagetest.TmpDB(t, "")
	defer cleaner()

	rec := &record.CallRequest{}
	gotRef, err := db.SetRecord(core.GenesisPulse.PulseNumber, rec)
	assert.Nil(t, err)

	gotRec, err := db.GetRecord(gotRef)
	assert.Nil(t, err)
	assert.Equal(t, rec, gotRec)

	_, err = db.SetRecord(core.GenesisPulse.PulseNumber, rec)
	assert.Equalf(t, err, storage.ErrOverride, "records override should be forbidden")
}

func TestDB_GetClassIndex_ReturnsNotFoundIfNoIndex(t *testing.T) {
	t.Parallel()
	db, cleaner := storagetest.TmpDB(t, "")
	defer cleaner()

	idx, err := db.GetClassIndex(core.NewRecordID(1, nil), false)
	assert.Equal(t, err, storage.ErrNotFound)
	assert.Nil(t, idx)
}

func TestDB_SetClassIndex_StoresCorrectDataInStorage(t *testing.T) {
	t.Parallel()
	db, cleaner := storagetest.TmpDB(t, "")
	defer cleaner()

	idgen := func() core.RecordID {
		return *core.NewRecordID(0, randhash())
	}
	latestRef := idgen()
	idx := index.ClassLifeline{
		LatestState: &latestRef,
	}
	zeroID := core.NewRecordID(0, hexhash("122444"))
	err := db.SetClassIndex(zeroID, &idx)
	assert.Nil(t, err)

	storedIndex, err := db.GetClassIndex(zeroID, false)
	assert.NoError(t, err)
	assert.Equal(t, *storedIndex, idx)
}

func TestDB_SetObjectIndex_ReturnsNotFoundIfNoIndex(t *testing.T) {
	t.Parallel()
	db, cleaner := storagetest.TmpDB(t, "")
	defer cleaner()

	idx, err := db.GetObjectIndex(core.NewRecordID(0, hexhash("5000")), false)
	assert.Equal(t, storage.ErrNotFound, err)
	assert.Nil(t, idx)
}

func TestDB_SetObjectIndex_StoresCorrectDataInStorage(t *testing.T) {
	t.Parallel()
	db, cleaner := storagetest.TmpDB(t, "")
	defer cleaner()

	idx := index.ObjectLifeline{
		ClassRef:    referenceWithHashes("50", "60"),
		LatestState: core.NewRecordID(0, hexhash("20")),
	}
	zeroid := core.NewRecordID(0, hexhash(""))
	err := db.SetObjectIndex(zeroid, &idx)
	assert.Nil(t, err)

	storedIndex, err := db.GetObjectIndex(zeroid, false)
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
	for i := 1; i < 4; i++ {
		setRecordMessage := message.SetRecord{
			Record: record.SerializeRecord(&record.CodeRecord{
				Code: []byte{byte(i)},
			}),
		}
		db.SetMessage(pulse, &setRecordMessage)
	}

	drop, messages, err := db.CreateDrop(pulse, []byte{4, 5, 6})
	assert.NoError(t, err)
	assert.Equal(t, 3, len(messages))
	assert.Equal(t, pulse, drop.Pulse)
	assert.Equal(t, "2aCdao6DhZSWQNTrtrxJW7QQZRb6UJ1ssRi9cg", base58.Encode(drop.Hash))

	for _, rawMessage := range messages {
		formatedMessage, err := message.Deserialize(bytes.NewBuffer(rawMessage))
		assert.NoError(t, err)
		assert.Equal(t, core.TypeSetRecord, formatedMessage.Message().Type())
	}
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
