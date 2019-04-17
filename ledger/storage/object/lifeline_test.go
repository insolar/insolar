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

	fuzz "github.com/google/gofuzz"
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
	assert.Equal(t, 0, len(indexStorage.indexStorage))
}

func TestIndexStorage_ForID(t *testing.T) {
	t.Parallel()

	ctx := inslogger.TestContext(t)

	jetID := gen.JetID()
	id := gen.ID()
	idx := Lifeline{
		LatestState: &id,
		JetID:       jetID,
		Delegates:   map[insolar.Reference]insolar.Reference{},
	}

	t.Run("returns correct index-value", func(t *testing.T) {
		t.Parallel()

		indexStorage := &IndexMemory{
			indexStorage: map[insolar.ID]Lifeline{},
		}
		indexStorage.indexStorage[id] = idx

		resultIdx, err := indexStorage.ForID(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, idx, resultIdx)
		assert.Equal(t, jetID, resultIdx.JetID)
	})

	t.Run("returns error when no index-value for id", func(t *testing.T) {
		t.Parallel()

		indexStorage := &IndexMemory{
			indexStorage: map[insolar.ID]Lifeline{},
		}
		indexStorage.indexStorage[id] = idx

		_, err := indexStorage.ForID(ctx, gen.ID())
		require.Error(t, err)
		assert.Equal(t, ErrIndexNotFound, err)
	})
}

func TestIndexDB_Set(t *testing.T) {
	t.Parallel()

	ctx := inslogger.TestContext(t)

	jetID := gen.JetID()
	id := gen.ID()
	idx := Lifeline{
		LatestState: &id,
		JetID:       jetID,
		Delegates:   map[insolar.Reference]insolar.Reference{},
	}

	jetIndex := store.NewJetIndexModifierMock(t)
	jetIndex.AddMock.Expect(id, jetID)

	t.Run("saves correct index-value", func(t *testing.T) {
		t.Parallel()

		indexStorage := NewIndexDB(store.NewMemoryMockDB())
		err := indexStorage.Set(ctx, id, idx)
		require.NoError(t, err)
		savedIdx, err := indexStorage.ForID(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, idx, savedIdx)
		assert.Equal(t, jetID, savedIdx.JetID)
	})

	t.Run("override indices is ok", func(t *testing.T) {
		t.Parallel()

		indexStorage := &IndexMemory{
			indexStorage:     map[insolar.ID]Lifeline{},
			jetIndexModifier: jetIndex,
		}
		err := indexStorage.Set(ctx, id, idx)
		require.NoError(t, err)

		err = indexStorage.Set(ctx, id, idx)
		assert.NoError(t, err)
	})

	t.Run("init delegates, when nil", func(t *testing.T) {
		t.Parallel()

		indexStorage := NewIndexDB(store.NewMemoryMockDB())
		err := indexStorage.Set(ctx, id, Lifeline{Delegates: nil})
		require.NoError(t, err)
		savedIdx, err := indexStorage.ForID(ctx, id)
		require.NoError(t, err)
		assert.NotNil(t, idx, savedIdx.Delegates)
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
		Delegates:   map[insolar.Reference]insolar.Reference{},
	}

	jetIndex := store.NewJetIndexModifierMock(t)
	jetIndex.AddMock.Expect(id, jetID)

	t.Run("saves correct index-value", func(t *testing.T) {
		t.Parallel()

		indexStorage := &IndexMemory{
			indexStorage:     map[insolar.ID]Lifeline{},
			jetIndexModifier: jetIndex,
		}
		err := indexStorage.Set(ctx, id, idx)
		require.NoError(t, err)
		assert.Equal(t, 1, len(indexStorage.indexStorage))
		assert.Equal(t, idx, indexStorage.indexStorage[id])
		assert.Equal(t, jetID, indexStorage.indexStorage[id].JetID)
	})

	t.Run("override indices is ok", func(t *testing.T) {
		t.Parallel()

		indexStorage := &IndexMemory{
			indexStorage:     map[insolar.ID]Lifeline{},
			jetIndexModifier: jetIndex,
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
			indexStorage:     map[insolar.ID]Lifeline{},
			jetIndexModifier: jetIndex,
		}
		err := indexStorage.Set(ctx, id, idx)
		require.NoError(t, err)
		assert.Equal(t, pn, indexStorage.indexStorage[id].LatestUpdate)
	})
}

func TestCloneObjectLifeline(t *testing.T) {
	t.Parallel()

	currentIdx := lifeline()

	clonedIdx := CloneIndex(currentIdx)

	assert.Equal(t, currentIdx, clonedIdx)
	assert.False(t, &currentIdx == &clonedIdx)
}

func TestCloneObjectLifeline_AlwaysFillInDelegates(t *testing.T) {
	t.Parallel()

	idx := Lifeline{}

	clonedIdx := CloneIndex(idx)

	assert.NotNil(t, clonedIdx.Delegates)
}

func TestCloneObjectLifeline_InsureDelegatesMapNotNil(t *testing.T) {
	t.Parallel()

	idx := Lifeline{}

	cloneIdx := CloneIndex(idx)

	require.NotNil(t, cloneIdx.Delegates)
}

func TestIndexMemory_ForPulseAndJet(t *testing.T) {
	t.Parallel()
	memStor := NewIndexMemory()
	ctx := inslogger.TestContext(t)

	jetID := gen.JetID()
	fPulse := gen.PulseNumber()
	sPulse := gen.PulseNumber()
	tPulse := gen.PulseNumber()

	_ = memStor.Set(ctx, *insolar.NewID(fPulse, []byte{1}), Lifeline{JetID: jetID, LatestUpdate: gen.PulseNumber()})
	memStor.SetUsageForPulse(ctx, *insolar.NewID(fPulse, []byte{1}), fPulse)
	_ = memStor.Set(ctx, *insolar.NewID(fPulse, []byte{2}), Lifeline{JetID: jetID, LatestUpdate: gen.PulseNumber()})
	memStor.SetUsageForPulse(ctx, *insolar.NewID(fPulse, []byte{2}), fPulse)
	_ = memStor.Set(ctx, *insolar.NewID(sPulse, nil), Lifeline{JetID: jetID})
	_ = memStor.Set(ctx, *insolar.NewID(tPulse, nil), Lifeline{JetID: jetID})

	res := memStor.ForPulseAndJet(ctx, fPulse, jetID)

	require.Equal(t, 2, len(res))
	_, ok := memStor.indexStorage[*insolar.NewID(fPulse, []byte{1})]
	require.Equal(t, true, ok)
	_, ok = memStor.indexStorage[*insolar.NewID(fPulse, []byte{2})]
	require.Equal(t, true, ok)
}

func id() (id *insolar.ID) {
	fuzz.New().NilChance(0.5).Fuzz(&id)
	return
}

func delegates() (result map[insolar.Reference]insolar.Reference) {
	fuzz.New().NilChance(0).NumElements(1, 10).Fuzz(&result)
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
