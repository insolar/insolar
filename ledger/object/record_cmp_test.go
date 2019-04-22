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
	"math/rand"
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
	memStorage := object.NewRecordMemory()
	dbStorage := object.NewRecordDB(store.NewMemoryMockDB())

	type tempRecord struct {
		id  insolar.ID
		rec record.MaterialRecord
	}

	var records []tempRecord

	f := fuzz.New().Funcs(func(t *tempRecord, c fuzz.Continue) {
		t.id = gen.ID()
		t.rec = record.MaterialRecord{
			Record: &object.ResultRecord{},
			JetID:  gen.JetID(),
		}
	})
	f.NilChance(0)
	f.NumElements(10, 20)
	f.Fuzz(&records)

	t.Run("saves correct record", func(t *testing.T) {
		t.Parallel()

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
		dbStorage := object.NewRecordDB(store.NewMemoryMockDB())

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
