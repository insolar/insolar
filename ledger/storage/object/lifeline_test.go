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
	"testing"

	"github.com/google/gofuzz"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/internal/ledger/store"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIndexStorage_NewStorageMemory(t *testing.T) {
	t.Parallel()

	indexStorage := NewIndexMemory()
	assert.Equal(t, 0, len(indexStorage.memory))
}

func TestIndexStorage_ForID(t *testing.T) {
	t.Parallel()

	ctx := inslogger.TestContext(t)

	jetID := gen.JetID()
	id := gen.ID()
	idx := Lifeline{
		LatestState: &id,
		JetID:       jetID,
	}

	t.Run("returns correct index-value", func(t *testing.T) {
		t.Parallel()

		indexStorage := &IndexMemory{
			memory: map[insolar.ID]Lifeline{},
		}
		indexStorage.memory[id] = idx

		resultIdx, err := indexStorage.ForID(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, idx, resultIdx)
		assert.Equal(t, jetID, resultIdx.JetID)
	})

	t.Run("returns error when no index-value for id", func(t *testing.T) {
		t.Parallel()

		indexStorage := &IndexMemory{
			memory: map[insolar.ID]Lifeline{},
		}
		indexStorage.memory[id] = idx

		_, err := indexStorage.ForID(ctx, gen.ID())
		require.Error(t, err)
		assert.Equal(t, ErrIndexNotFound, err)
	})
}

func TestIndexStorage_Set(t *testing.T) {
	t.Parallel()

	ctx := inslogger.TestContext(t)

	jetID := gen.JetID()
	id := gen.ID()
	idx := Lifeline{
		LatestState: &id,
		JetID:       jetID,
	}

	jetIndex := store.NewJetIndexModifierMock(t)
	jetIndex.AddMock.Expect(id, jetID)

	t.Run("saves correct index-value", func(t *testing.T) {
		t.Parallel()

		indexStorage := &IndexMemory{
			memory:   map[insolar.ID]Lifeline{},
			jetIndex: jetIndex,
		}
		err := indexStorage.Set(ctx, id, idx)
		require.NoError(t, err)
		assert.Equal(t, 1, len(indexStorage.memory))
		assert.Equal(t, idx, indexStorage.memory[id])
		assert.Equal(t, jetID, indexStorage.memory[id].JetID)
	})

	t.Run("override indices is ok", func(t *testing.T) {
		t.Parallel()

		indexStorage := &IndexMemory{
			memory:   map[insolar.ID]Lifeline{},
			jetIndex: jetIndex,
		}
		err := indexStorage.Set(ctx, id, idx)
		require.NoError(t, err)

		err = indexStorage.Set(ctx, id, idx)
		assert.NoError(t, err)
	})
}

func TestIndexStorage_Set_SaveLastUpdate(t *testing.T) {
	t.Parallel()

	ctx := inslogger.TestContext(t)

	jetID := gen.JetID()
	id := gen.ID()
	pn := gen.PulseNumber()
	idx := Lifeline{
		LatestState:  &id,
		LatestUpdate: pn,
		JetID:        jetID,
	}

	jetIndex := store.NewJetIndexModifierMock(t)
	jetIndex.AddMock.Expect(id, jetID)

	t.Run("saves correct LastUpdate field in index", func(t *testing.T) {
		t.Parallel()

		indexStorage := &IndexMemory{
			memory:   map[insolar.ID]Lifeline{},
			jetIndex: jetIndex,
		}
		err := indexStorage.Set(ctx, id, idx)
		require.NoError(t, err)
		assert.Equal(t, pn, indexStorage.memory[id].LatestUpdate)
	})
}

func TestCloneObjectLifeline(t *testing.T) {
	t.Parallel()

	currentIdx := lifeline()

	clonedIdx := CloneIndex(currentIdx)

	assert.Equal(t, currentIdx, clonedIdx)
	assert.False(t, &currentIdx == &clonedIdx)
}

func id() (id *insolar.ID) {
	fuzz.New().NilChance(0.5).Fuzz(&id)
	return
}

func delegates() (result map[insolar.Reference]insolar.Reference) {
	fuzz.New().NilChance(0.5).NumElements(1, 10).Fuzz(&result)
	return
}

func state() (state StateID) {
	fuzz.New().NilChance(0).Fuzz(&state)
	return
}

func lifeline() Lifeline {
	var index Lifeline
	fuzz.New().NilChance(0).Funcs(
		func(idx *Lifeline, c fuzz.Continue) {
			idx.LatestState = id()
			idx.LatestStateApproved = id()
			idx.ChildPointer = id()
			idx.Delegates = delegates()
			idx.State = state()
			idx.Parent = gen.Reference()
			idx.LatestUpdate = gen.PulseNumber()
			idx.JetID = gen.JetID()
		},
	).Fuzz(&index)

	return index
}
