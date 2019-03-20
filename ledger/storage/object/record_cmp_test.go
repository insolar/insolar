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

package object_test

import (
	"math/rand"
	"testing"

	fuzz "github.com/google/gofuzz"
	"github.com/insolar/insolar"
	"github.com/insolar/insolar/gen"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/storage/object"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInMemoryRecord(t *testing.T) {
	t.Parallel()

	ctx := inslogger.TestContext(t)

	recordStorage := object.NewRecordMemory()

	type tempRecord struct {
		id  insolar.ID
		rec object.MaterialRecord
	}

	var records []tempRecord

	f := fuzz.New().Funcs(func(t *tempRecord, c fuzz.Continue) {
		t.id = gen.ID()
		t.rec = object.MaterialRecord{
			Record: &object.ResultRecord{},
			JetID:  gen.JetID(),
		}
	})
	f.NumElements(5, 10).NilChance(0).Fuzz(&records)

	t.Run("saves correct record", func(t *testing.T) {
		for _, r := range records {
			err := recordStorage.Set(ctx, r.id, r.rec)
			require.NoError(t, err)
		}

		for _, r := range records {
			resRecord, err := recordStorage.ForID(ctx, r.id)
			require.NoError(t, err)

			assert.Equal(t, r.rec, resRecord)
			assert.Equal(t, r.rec.Record, resRecord.Record)
			assert.Equal(t, r.rec.JetID, resRecord.JetID)
		}
	})

	t.Run("returns error when no record for id", func(t *testing.T) {
		t.Parallel()

		for i := int32(0); i < rand.Int31n(10); i++ {
			_, err := recordStorage.ForID(ctx, gen.ID())
			require.Error(t, err)
			assert.Equal(t, object.RecNotFound, err)
		}
	})

	t.Run("returns override error when saving with the same id", func(t *testing.T) {
		t.Parallel()

		recordStorage := object.NewRecordMemory()
		for _, r := range records {
			err := recordStorage.Set(ctx, r.id, r.rec)
			require.NoError(t, err)
		}

		for _, r := range records {
			err := recordStorage.Set(ctx, r.id, r.rec)
			require.Error(t, err)
			assert.Equal(t, object.ErrOverride, err)
		}
	})
}
