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

package object

import (
	"testing"

	"github.com/insolar/insolar"
	"github.com/insolar/insolar/gen"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/storage/db"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRecordStorage_NewStorageMemory(t *testing.T) {
	t.Parallel()

	recordStorage := NewRecordMemory()
	assert.Equal(t, 0, len(recordStorage.memory))
}

func TestRecordStorage_ForID(t *testing.T) {
	t.Parallel()

	ctx := inslogger.TestContext(t)

	jetID := gen.JetID()
	id := gen.ID()
	rec := MaterialRecord{
		Record: &ResultRecord{},
		JetID:  jetID,
	}

	t.Run("returns correct record-value", func(t *testing.T) {
		t.Parallel()

		recordStorage := &RecordMemory{
			memory: map[insolar.ID]MaterialRecord{},
		}
		recordStorage.memory[id] = rec

		resultRec, err := recordStorage.ForID(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, rec, resultRec)
		assert.Equal(t, jetID, resultRec.JetID)
	})

	t.Run("returns error when no record-value for id", func(t *testing.T) {
		t.Parallel()

		recordStorage := &RecordMemory{
			memory: map[insolar.ID]MaterialRecord{},
		}
		recordStorage.memory[id] = rec

		_, err := recordStorage.ForID(ctx, gen.ID())
		require.Error(t, err)
		assert.Equal(t, RecNotFound, err)
	})
}

func TestRecordStorage_Set(t *testing.T) {
	t.Parallel()

	ctx := inslogger.TestContext(t)

	jetID := gen.JetID()
	id := gen.ID()
	rec := MaterialRecord{
		Record: &ResultRecord{},
		JetID:  jetID,
	}

	jetIndex := db.NewJetIndexModifierMock(t)
	jetIndex.AddMock.Expect(id, jetID)

	t.Run("saves correct record-value", func(t *testing.T) {
		t.Parallel()

		recordStorage := &RecordMemory{
			memory:   map[insolar.ID]MaterialRecord{},
			jetIndex: jetIndex,
		}
		err := recordStorage.Set(ctx, id, rec)
		require.NoError(t, err)
		assert.Equal(t, 1, len(recordStorage.memory))
		assert.Equal(t, rec, recordStorage.memory[id])
		assert.Equal(t, jetID, recordStorage.memory[id].JetID)
	})

	t.Run("returns override error when saving with the same id", func(t *testing.T) {
		t.Parallel()

		recordStorage := &RecordMemory{
			memory:   map[insolar.ID]MaterialRecord{},
			jetIndex: jetIndex,
		}
		err := recordStorage.Set(ctx, id, rec)
		require.NoError(t, err)

		err = recordStorage.Set(ctx, id, rec)
		require.Error(t, err)
		assert.Equal(t, ErrOverride, err)
	})
}
