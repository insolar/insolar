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

package object_test

import (
	"context"
	"io/ioutil"
	"math/rand"
	"os"
	"testing"

	fuzz "github.com/google/gofuzz"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/internal/ledger/store"
	"github.com/insolar/insolar/ledger/object"
)

func TestRecord_Components(t *testing.T) {
	ctx := inslogger.TestContext(t)

	type tempRecord struct {
		id  insolar.ID
		rec record.Material
	}

	var records []tempRecord

	f := fuzz.New().Funcs(func(t *tempRecord, c fuzz.Continue) {
		t.id = gen.ID()
		t.rec = getMaterialRecord()
	})
	f.NilChance(0)
	f.NumElements(10, 20)
	f.Fuzz(&records)

	t.Run("saves correct record", func(t *testing.T) {
		t.Parallel()

		memStorage := object.NewRecordMemory()
		tmpdir, err := ioutil.TempDir("", "bdb-test-")
		defer os.RemoveAll(tmpdir)
		require.NoError(t, err)

		db, err := store.NewTestBadgerDB(tmpdir)
		require.NoError(t, err)
		defer db.Stop(context.Background())
		dbStorage := object.NewRecordDB(db)

		for _, r := range records {
			memErr := memStorage.Set(ctx, r.id, r.rec)
			dbErr := dbStorage.Set(ctx, r.id, r.rec)
			require.NoError(t, memErr)
			require.NoError(t, dbErr)
		}

		for _, r := range records {
			memRecord, memErr := memStorage.ForID(ctx, r.id)
			dbRecord, dbErr := dbStorage.ForID(ctx, r.id)
			require.NoError(t, memErr)
			require.NoError(t, dbErr)

			assert.Equal(t, r.rec, memRecord)
			assert.Equal(t, r.rec, dbRecord)
		}
	})

	t.Run("returns error when no record for id", func(t *testing.T) {
		t.Parallel()
		memStorage := object.NewRecordMemory()
		tmpdir, err := ioutil.TempDir("", "bdb-test-")
		defer os.RemoveAll(tmpdir)
		require.NoError(t, err)

		db, err := store.NewTestBadgerDB(tmpdir)
		require.NoError(t, err)
		defer db.Stop(context.Background())
		dbStorage := object.NewRecordDB(db)

		for i := int32(0); i < rand.Int31n(10); i++ {
			_, memErr := memStorage.ForID(ctx, gen.ID())
			_, dbErr := dbStorage.ForID(ctx, gen.ID())
			require.Error(t, memErr)
			require.Error(t, dbErr)
			assert.Equal(t, object.ErrNotFound, memErr)
			assert.Equal(t, object.ErrNotFound, dbErr)
		}
	})

	t.Run("returns override error when saving with the same id", func(t *testing.T) {
		t.Parallel()

		memStorage := object.NewRecordMemory()
		tmpdir1, err := ioutil.TempDir("", "bdb-test-")
		defer os.RemoveAll(tmpdir1)
		require.NoError(t, err)

		db1, err := store.NewTestBadgerDB(tmpdir1)
		defer db1.Stop(context.Background())
		dbStorage := object.NewRecordDB(db1)

		for _, r := range records {
			memErr := memStorage.Set(ctx, r.id, r.rec)
			dbErr := dbStorage.Set(ctx, r.id, r.rec)
			require.NoError(t, memErr)
			require.NoError(t, dbErr)
		}

		for _, r := range records {
			memErr := memStorage.Set(ctx, r.id, r.rec)
			dbErr := dbStorage.Set(ctx, r.id, r.rec)
			require.Error(t, memErr)
			require.Error(t, dbErr)
			assert.Equal(t, object.ErrOverride, memErr)
			assert.Equal(t, object.ErrOverride, dbErr)
		}
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
		Virtual: &virtRec,
		JetID:   gen.JetID(),
	}

	return materialRecord
}
