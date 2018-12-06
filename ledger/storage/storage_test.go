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

	"github.com/insolar/insolar/platformpolicy"
	"github.com/jbenet/go-base58"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/core/message"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/index"
	"github.com/insolar/insolar/ledger/jetdrop"
	"github.com/insolar/insolar/ledger/record"
	"github.com/insolar/insolar/ledger/storage"
	"github.com/insolar/insolar/ledger/storage/storagetest"
)

func TestDB_GetRecordNotFound(t *testing.T) {
	t.Parallel()
	ctx := inslogger.TestContext(t)
	db, cleaner := storagetest.TmpDB(ctx, t)
	defer cleaner()

	rec, err := db.GetRecord(ctx, &core.RecordID{})
	assert.Equal(t, err, storage.ErrNotFound)
	assert.Nil(t, rec)
}

func TestDB_SetRecord(t *testing.T) {
	t.Parallel()
	ctx := inslogger.TestContext(t)
	db, cleaner := storagetest.TmpDB(ctx, t)
	defer cleaner()

	rec := &record.CallRequest{}
	gotRef, err := db.SetRecord(ctx, core.GenesisPulse.PulseNumber, rec)
	assert.Nil(t, err)

	gotRec, err := db.GetRecord(ctx, gotRef)
	assert.Nil(t, err)
	assert.Equal(t, rec, gotRec)

	_, err = db.SetRecord(ctx, core.GenesisPulse.PulseNumber, rec)
	assert.Equalf(t, err, storage.ErrOverride, "records override should be forbidden")
}

func TestDB_SetObjectIndex_ReturnsNotFoundIfNoIndex(t *testing.T) {
	t.Parallel()
	ctx := inslogger.TestContext(t)
	db, cleaner := storagetest.TmpDB(ctx, t)
	defer cleaner()

	idx, err := db.GetObjectIndex(ctx, core.NewRecordID(0, hexhash("5000")), false)
	assert.Equal(t, storage.ErrNotFound, err)
	assert.Nil(t, idx)
}

func TestDB_SetObjectIndex_StoresCorrectDataInStorage(t *testing.T) {
	t.Parallel()
	ctx := inslogger.TestContext(t)
	db, cleaner := storagetest.TmpDB(ctx, t)
	defer cleaner()

	idx := index.ObjectLifeline{
		LatestState: core.NewRecordID(0, hexhash("20")),
	}
	zeroid := core.NewRecordID(0, hexhash(""))
	err := db.SetObjectIndex(ctx, zeroid, &idx)
	assert.Nil(t, err)

	storedIndex, err := db.GetObjectIndex(ctx, zeroid, false)
	assert.NoError(t, err)
	assert.Equal(t, *storedIndex, idx)
}

func TestDB_GetDrop_ReturnsNotFoundIfNoDrop(t *testing.T) {
	t.Parallel()
	ctx := inslogger.TestContext(t)
	db, cleaner := storagetest.TmpDB(ctx, t)
	defer cleaner()

	drop, err := db.GetDrop(ctx, 1)
	assert.Equal(t, err, storage.ErrNotFound)
	assert.Nil(t, drop)
}

func TestDB_CreateDrop(t *testing.T) {
	t.Parallel()
	ctx := inslogger.TestContext(t)
	db, cleaner := storagetest.TmpDB(ctx, t)
	defer cleaner()

	pulse := core.PulseNumber(core.FirstPulseNumber + 10)
	err := db.AddPulse(
		ctx,
		core.Pulse{
			PulseNumber: pulse,
			Entropy:     core.Entropy{1, 2, 3},
		},
	)
	cs := platformpolicy.NewPlatformCryptographyScheme()

	for i := 1; i < 4; i++ {
		setRecordMessage := message.SetRecord{
			Record: record.SerializeRecord(&record.CodeRecord{
				Code: record.CalculateIDForBlob(cs, pulse, []byte{byte(i)}),
			}),
		}
		db.SetMessage(ctx, pulse, &setRecordMessage)
		db.SetBlob(ctx, pulse, []byte{byte(i)})
	}

	drop, messages, err := db.CreateDrop(ctx, pulse, []byte{4, 5, 6})
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
	ctx := inslogger.TestContext(t)
	db, cleaner := storagetest.TmpDB(ctx, t)
	defer cleaner()

	drop42 := jetdrop.JetDrop{
		Pulse: 42,
		Hash:  []byte{0xFF},
	}
	err := db.SetDrop(ctx, &drop42)
	assert.NoError(t, err)

	got, err := db.GetDrop(ctx, 42)
	assert.NoError(t, err)
	assert.Equal(t, *got, drop42)
}

func TestDB_AddPulse(t *testing.T) {
	t.Parallel()
	ctx := inslogger.TestContext(t)
	db, cleaner := storagetest.TmpDB(ctx, t)
	defer cleaner()

	err := db.AddPulse(
		ctx,
		core.Pulse{PulseNumber: 42, Entropy: core.Entropy{1, 2, 3}},
	)
	assert.NoError(t, err)
	latestPulse, err := db.GetLatestPulse(ctx)
	assert.Equal(t, core.PulseNumber(42), latestPulse.Pulse.PulseNumber)
	pulse, err := db.GetPulse(ctx, latestPulse.Pulse.PulseNumber)
	assert.NoError(t, err)
	prev := core.PulseNumber(core.FirstPulseNumber)
	assert.Equal(t, storage.Pulse{Prev: &prev, Pulse: core.Pulse{Entropy: core.Entropy{1, 2, 3}, PulseNumber: 42}}, *pulse)
}

func TestDB_SetLocalData(t *testing.T) {
	t.Parallel()
	ctx := inslogger.TestContext(t)
	db, cleaner := storagetest.TmpDB(ctx, t)
	defer cleaner()

	err := db.SetLocalData(ctx, 0, []byte{1}, []byte{2})
	require.NoError(t, err)

	data, err := db.GetLocalData(ctx, 0, []byte{1})
	require.NoError(t, err)
	assert.Equal(t, []byte{2}, data)

	_, err = db.GetLocalData(ctx, 1, []byte{1})
	assert.Equal(t, storage.ErrNotFound, err)
}

func TestDB_IterateLocalData(t *testing.T) {
	t.Parallel()
	ctx := inslogger.TestContext(t)
	db, cleaner := storagetest.TmpDB(ctx, t)
	defer cleaner()

	err := db.SetLocalData(ctx, 1, []byte{1, 1}, []byte{1})
	require.NoError(t, err)
	err = db.SetLocalData(ctx, 1, []byte{1, 2}, []byte{2})
	require.NoError(t, err)
	err = db.SetLocalData(ctx, 1, []byte{2, 1}, []byte{3})
	require.NoError(t, err)
	err = db.SetLocalData(ctx, 2, []byte{1, 1}, []byte{4})
	require.NoError(t, err)

	type tuple struct {
		k []byte
		v []byte
	}
	var results []tuple
	err = db.IterateLocalData(ctx, 1, []byte{1}, func(k, v []byte) error {
		results = append(results, tuple{k: k, v: v})
		return nil
	})
	require.NoError(t, err)
	assert.Equal(t, []tuple{
		{k: []byte{1}, v: []byte{1}},
		{k: []byte{2}, v: []byte{2}},
	}, results)
}
