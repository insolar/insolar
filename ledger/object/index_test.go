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

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/internal/ledger/store"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInMemoryIndex_SetLifeline(t *testing.T) {
	t.Parallel()

	ctx := inslogger.TestContext(t)

	jetID := gen.JetID()
	id := gen.ID()
	idx := Lifeline{
		LatestState: &id,
		JetID:       jetID,
		Delegates:   []LifelineDelegate{},
	}

	t.Run("saves correct index-value", func(t *testing.T) {
		t.Parallel()

		storage := NewInMemoryIndex()
		pn := gen.PulseNumber()

		err := storage.SetLifeline(ctx, pn, id, idx)

		require.NoError(t, err)
		assert.Equal(t, 1, len(storage.buckets))

		buck, buckOK := storage.buckets[pn]
		require.Equal(t, true, buckOK)
		require.Equal(t, 1, len(buck))

		meta, metaOK := buck[id]
		require.Equal(t, true, metaOK)
		require.NotNil(t, meta)
		require.NotNil(t, meta.bucket)

		require.Equal(t, meta.bucket.Lifeline, idx)
		require.Equal(t, meta.bucket.LifelineLastUsed, pn)
		require.Equal(t, meta.bucket.ObjID, id)
	})

	t.Run("save multiple values", func(t *testing.T) {
		fID := insolar.NewID(1, nil)
		sID := insolar.NewID(2, nil)
		tID := insolar.NewID(3, nil)
		fthID := insolar.NewID(4, nil)

		storage := NewInMemoryIndex()
		err := storage.SetLifeline(ctx, 1, *fID, idx)
		require.NoError(t, err)
		err = storage.SetLifeline(ctx, 1, *sID, idx)
		require.NoError(t, err)
		err = storage.SetLifeline(ctx, 2, *tID, idx)
		require.NoError(t, err)
		err = storage.SetLifeline(ctx, 2, *fthID, idx)
		require.NoError(t, err)

		require.Equal(t, 2, len(storage.buckets))
		require.Equal(t, 2, len(storage.buckets[1]))
		require.Equal(t, 2, len(storage.buckets[2]))
		require.Equal(t, *fID, storage.buckets[1][*fID].bucket.ObjID)
		require.Equal(t, *sID, storage.buckets[1][*sID].bucket.ObjID)
		require.Equal(t, *tID, storage.buckets[2][*tID].bucket.ObjID)
		require.Equal(t, *fthID, storage.buckets[2][*fthID].bucket.ObjID)
	})

	t.Run("override indices is ok", func(t *testing.T) {
		t.Parallel()

		storage := NewInMemoryIndex()
		pn := gen.PulseNumber()

		err := storage.SetLifeline(ctx, pn, id, idx)
		require.NoError(t, err)

		err = storage.SetLifeline(ctx, pn, id, idx)
		require.NoError(t, err)
	})
}

func TestIndexStorage_ForID(t *testing.T) {
	t.Parallel()

	ctx := inslogger.TestContext(t)

	jetID := gen.JetID()
	id := gen.ID()
	idx := Lifeline{
		LatestState: &id,
		JetID:       jetID,
		Delegates:   []LifelineDelegate{},
	}

	t.Run("returns correct index-value", func(t *testing.T) {
		t.Parallel()

		storage := NewInMemoryIndex()
		pn := gen.PulseNumber()

		err := storage.SetLifeline(ctx, pn, id, idx)
		require.NoError(t, err)

		res, err := storage.LifelineForID(ctx, pn, id)

		require.NoError(t, err)
		assert.Equal(t, idx, res)
	})

	t.Run("returns error when no index-value for id", func(t *testing.T) {
		t.Parallel()

		storage := NewInMemoryIndex()
		pn := gen.PulseNumber()

		_, err := storage.LifelineForID(ctx, pn, id)

		assert.Equal(t, ErrLifelineNotFound, err)
	})
}

func TestInMemoryIndex_ForPNAndJet(t *testing.T) {
	t.Parallel()

	ctx := inslogger.TestContext(t)

	fJetId := insolar.NewJetID(1, []byte{1})
	sJetId := insolar.NewJetID(1, []byte{2})
	tJetId := insolar.NewJetID(1, []byte{3})

	fId := insolar.NewID(123, []byte{})
	sId := insolar.NewID(124, []byte{})
	tId := insolar.NewID(125, []byte{})

	fPn := insolar.PulseNumber(1)
	sPn := insolar.PulseNumber(2)

	fIdx := Lifeline{
		LatestState: insolar.NewID(123, []byte{}),
		JetID:       *fJetId,
		Delegates:   []LifelineDelegate{},
	}
	sIdx := Lifeline{
		LatestState: insolar.NewID(124, []byte{}),
		JetID:       *sJetId,
		Delegates:   []LifelineDelegate{},
	}
	tIdx := Lifeline{
		LatestState: insolar.NewID(125, []byte{}),
		JetID:       *tJetId,
		Delegates:   []LifelineDelegate{},
	}

	index := NewInMemoryIndex()

	_ = index.SetLifeline(ctx, fPn, *fId, fIdx)
	_ = index.SetLifeline(ctx, fPn, *sId, sIdx)
	_ = index.SetLifeline(ctx, sPn, *tId, tIdx)

	res := index.ForPNAndJet(ctx, fPn, *fJetId)
	require.Equal(t, 1, len(res))
	require.NotNil(t, res[0].Lifeline)
	require.Equal(t, *fId, res[0].ObjID)
	require.Equal(t, fIdx.LatestState, res[0].Lifeline.LatestState)

	res = index.ForPNAndJet(ctx, fPn, *sJetId)
	require.Equal(t, 1, len(res))
	require.NotNil(t, res[0].Lifeline)
	require.Equal(t, *sId, res[0].ObjID)
	require.Equal(t, sIdx.LatestState, res[0].Lifeline.LatestState)

}

func TestInMemoryIndex_SetBucket(t *testing.T) {
	t.Parallel()

	ctx := inslogger.TestContext(t)
	objID := gen.ID()
	lflID := gen.ID()
	jetID := gen.JetID()
	buck := IndexBucket{
		ObjID: objID,
		Lifeline: Lifeline{
			LatestState: &lflID,
			JetID:       jetID,
			Delegates:   []LifelineDelegate{},
		},
	}

	t.Run("saves correct bucket", func(t *testing.T) {
		pn := gen.PulseNumber()
		index := NewInMemoryIndex()

		err := index.SetBucket(ctx, pn, buck)
		require.NoError(t, err)

		savedBuck := index.buckets[pn][objID]
		require.NotNil(t, savedBuck)

		buckBuf, _ := buck.Marshal()
		savedBuckBuf, _ := savedBuck.bucket.Marshal()

		require.Equal(t, buckBuf, savedBuckBuf)
	})

	t.Run("re-save works fine", func(t *testing.T) {
		pn := gen.PulseNumber()
		index := NewInMemoryIndex()

		err := index.SetBucket(ctx, pn, buck)
		require.NoError(t, err)

		sLlflID := gen.ID()
		sJetID := gen.JetID()
		sBuck := IndexBucket{
			ObjID: objID,
			Lifeline: Lifeline{
				LatestState: &sLlflID,
				JetID:       sJetID,
				Delegates:   []LifelineDelegate{},
			},
		}

		err = index.SetBucket(ctx, pn, sBuck)
		require.NoError(t, err)

		savedBuck := index.buckets[pn][objID]
		require.NotNil(t, savedBuck)

		sBuckBuf, _ := sBuck.Marshal()
		savedBuckBuf, _ := savedBuck.bucket.Marshal()

		require.Equal(t, sBuckBuf, savedBuckBuf)
	})
}

func TestInMemoryIndex_SetLifelineUsage(t *testing.T) {
	t.Parallel()

	ctx := inslogger.TestContext(t)

	jetID := gen.JetID()
	id := gen.ID()
	idx := Lifeline{
		LatestState: &id,
		JetID:       jetID,
		Delegates:   []LifelineDelegate{},
	}

	pn := gen.PulseNumber()
	newPN := pn + 1

	t.Run("works fine", func(t *testing.T) {
		t.Parallel()

		index := NewInMemoryIndex()

		_ = index.SetLifeline(ctx, pn, id, idx)

		require.Equal(t, pn, index.buckets[pn][id].bucket.LifelineLastUsed)

		index.buckets[newPN] = index.buckets[pn]

		err := index.SetLifelineUsage(ctx, newPN, id)

		require.NoError(t, err)
		require.Equal(t, newPN, index.buckets[newPN][id].bucket.LifelineLastUsed)
	})

	t.Run("returns ErrLifelineNotFound if no bucket", func(t *testing.T) {
		t.Parallel()

		index := NewInMemoryIndex()
		err := index.SetLifelineUsage(ctx, pn, id)
		require.Error(t, ErrLifelineNotFound, err)
	})
}

func TestNewInMemoryIndex_DeleteForPN(t *testing.T) {
	t.Parallel()

	ctx := inslogger.TestContext(t)

	fPn := gen.PulseNumber()
	sPn := fPn + 1
	tPn := sPn + 1

	index := NewInMemoryIndex()

	index.buckets[fPn] = map[insolar.ID]*LockedIndexBucket{}
	index.buckets[sPn] = map[insolar.ID]*LockedIndexBucket{}
	index.buckets[tPn] = map[insolar.ID]*LockedIndexBucket{}

	index.DeleteForPN(ctx, sPn)

	_, ok := index.buckets[fPn]
	require.Equal(t, true, ok)
	_, ok = index.buckets[sPn]
	require.Equal(t, false, ok)
	_, ok = index.buckets[tPn]
	require.Equal(t, true, ok)
}

func TestDBIndex_SetLifeline(t *testing.T) {
	t.Parallel()

	ctx := inslogger.TestContext(t)

	jetID := gen.JetID()
	id := gen.ID()
	idx := Lifeline{
		LatestState: &id,
		JetID:       jetID,
		Delegates:   []LifelineDelegate{},
	}

	t.Run("saves correct index-value", func(t *testing.T) {
		t.Parallel()

		storage := NewIndexDB(store.NewMemoryMockDB())
		pn := gen.PulseNumber()

		err := storage.SetLifeline(ctx, pn, id, idx)

		require.NoError(t, err)
	})

	t.Run("override indices is ok", func(t *testing.T) {
		t.Parallel()

		storage := NewInMemoryIndex()
		pn := gen.PulseNumber()

		err := storage.SetLifeline(ctx, pn, id, idx)
		require.NoError(t, err)

		err = storage.SetLifeline(ctx, pn, id, idx)
		require.NoError(t, err)
	})
}

func TestDBIndexStorage_ForID(t *testing.T) {
	t.Parallel()

	ctx := inslogger.TestContext(t)

	jetID := gen.JetID()
	id := gen.ID()
	idx := Lifeline{
		LatestState: &id,
		JetID:       jetID,
		Delegates:   []LifelineDelegate{},
	}

	t.Run("returns correct index-value", func(t *testing.T) {
		t.Parallel()

		storage := NewIndexDB(store.NewMemoryMockDB())
		pn := gen.PulseNumber()

		err := storage.SetLifeline(ctx, pn, id, idx)
		require.NoError(t, err)

		res, err := storage.LifelineForID(ctx, pn, id)
		require.NoError(t, err)

		idxBuf, _ := idx.Marshal()
		resBuf, _ := res.Marshal()

		assert.Equal(t, idxBuf, resBuf)
	})

	t.Run("returns error when no index-value for id", func(t *testing.T) {
		t.Parallel()

		storage := NewIndexDB(store.NewMemoryMockDB())
		pn := gen.PulseNumber()

		_, err := storage.LifelineForID(ctx, pn, id)

		assert.Equal(t, ErrLifelineNotFound, err)
	})
}

func TestDBIndex_SetBucket(t *testing.T) {
	t.Parallel()

	ctx := inslogger.TestContext(t)
	objID := gen.ID()
	lflID := gen.ID()
	jetID := gen.JetID()
	buck := IndexBucket{
		ObjID: objID,
		Lifeline: Lifeline{
			LatestState: &lflID,
			JetID:       jetID,
			Delegates:   []LifelineDelegate{},
		},
	}

	t.Run("saves correct bucket", func(t *testing.T) {
		pn := gen.PulseNumber()
		index := NewIndexDB(store.NewMemoryMockDB())

		err := index.SetBucket(ctx, pn, buck)
		require.NoError(t, err)

		res, err := index.LifelineForID(ctx, pn, objID)
		require.NoError(t, err)

		idxBuf, _ := buck.Lifeline.Marshal()
		resBuf, _ := res.Marshal()

		assert.Equal(t, idxBuf, resBuf)
	})

	t.Run("re-save works fine", func(t *testing.T) {
		pn := gen.PulseNumber()
		index := NewIndexDB(store.NewMemoryMockDB())

		err := index.SetBucket(ctx, pn, buck)
		require.NoError(t, err)

		sLlflID := gen.ID()
		sJetID := gen.JetID()
		sBuck := IndexBucket{
			ObjID: objID,
			Lifeline: Lifeline{
				LatestState: &sLlflID,
				JetID:       sJetID,
				Delegates:   []LifelineDelegate{},
			},
		}

		err = index.SetBucket(ctx, pn, sBuck)
		require.NoError(t, err)

		res, err := index.LifelineForID(ctx, pn, objID)
		require.NoError(t, err)

		idxBuf, _ := sBuck.Lifeline.Marshal()
		resBuf, _ := res.Marshal()

		assert.Equal(t, idxBuf, resBuf)
	})
}
