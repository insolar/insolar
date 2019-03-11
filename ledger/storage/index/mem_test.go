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

package index

import (
	"testing"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/gen"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/storage/db"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIndexStorage_NewStorageMem(t *testing.T) {
	t.Parallel()

	indexStorage := NewStorageMem()
	assert.Equal(t, 0, len(indexStorage.memory))
}

func TestIndexStorage_ForID(t *testing.T) {
	t.Parallel()

	ctx := inslogger.TestContext(t)

	jetID := gen.JetID()
	id := gen.ID()
	idx := ObjectLifeline{
		LatestState: &id,
		JetID:       jetID,
	}

	t.Run("returns correct index-value", func(t *testing.T) {
		t.Parallel()

		indexStorage := &StorageMem{
			memory: map[core.RecordID]ObjectLifeline{},
		}
		indexStorage.memory[id] = idx

		resultIdx, err := indexStorage.ForID(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, idx, resultIdx)
		assert.Equal(t, jetID, resultIdx.JetID)
	})

	t.Run("returns error when no index-value for id", func(t *testing.T) {
		t.Parallel()

		indexStorage := &StorageMem{
			memory: map[core.RecordID]ObjectLifeline{},
		}
		indexStorage.memory[id] = idx

		_, err := indexStorage.ForID(ctx, gen.ID())
		require.Error(t, err)
		assert.Equal(t, ErrNotFound, err)
	})
}

func TestIndexStorage_Set(t *testing.T) {
	t.Parallel()

	ctx := inslogger.TestContext(t)

	jetID := gen.JetID()
	id := gen.ID()
	idx := ObjectLifeline{
		LatestState: &id,
		JetID:       jetID,
	}

	jetIndex := db.NewJetIndexModifierMock(t)
	jetIndex.AddMock.Expect(id, jetID)

	t.Run("saves correct index-value", func(t *testing.T) {
		t.Parallel()

		indexStorage := &StorageMem{
			memory:   map[core.RecordID]ObjectLifeline{},
			jetIndex: jetIndex,
		}
		err := indexStorage.Set(ctx, id, idx)
		require.NoError(t, err)
		assert.Equal(t, 1, len(indexStorage.memory))
		assert.Equal(t, idx, indexStorage.memory[id])
		assert.Equal(t, jetID, indexStorage.memory[id].JetID)
	})

	t.Run("returns override error when saving with the same id", func(t *testing.T) {
		t.Parallel()

		indexStorage := &StorageMem{
			memory:   map[core.RecordID]ObjectLifeline{},
			jetIndex: jetIndex,
		}
		err := indexStorage.Set(ctx, id, idx)
		require.NoError(t, err)

		err = indexStorage.Set(ctx, id, idx)
		require.Error(t, err)
		assert.Equal(t, ErrOverride, err)
	})
}
