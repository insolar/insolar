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

//
// func TestLifelineStorage_Set(t *testing.T) {
// 	t.Parallel()
//
// 	ctx := inslogger.TestContext(t)
//
// 	jetID := gen.JetID()
// 	id := gen.ID()
// 	idx := Lifeline{
// 		LatestState: &id,
// 		JetID:       jetID,
// 		Delegates:   []LifelineDelegate{},
// 	}
//
// 	t.Run("saves correct index-value", func(t *testing.T) {
// 		t.Parallel()
//
// 		idxStor := NewIndexStorageMemory()
// 		storage := NewLifelineStorage(idxStor, idxStor)
// 		pn := gen.PulseNumber()
//
// 		err := storage.Set(ctx, pn, id, idx)
//
// 		require.NoError(t, err)
// 		require.Equal(t, 1, len(idxStor.buckets))
//
// 		buck, buckOK := idxStor.buckets[pn]
// 		require.Equal(t, true, buckOK)
// 		require.Equal(t, 1, len(buck))
//
// 		meta, metaOK := buck[id]
// 		require.Equal(t, true, metaOK)
// 		require.NotNil(t, meta)
// 		require.NotNil(t, meta)
//
// 		require.Equal(t, meta.Lifeline, idx)
// 		require.Equal(t, meta.LifelineLastUsed, pn)
// 		require.Equal(t, meta.ObjID, id)
// 	})
//
// 	t.Run("save multiple values", func(t *testing.T) {
// 		fID := insolar.NewID(1, nil)
// 		sID := insolar.NewID(2, nil)
// 		tID := insolar.NewID(3, nil)
// 		fthID := insolar.NewID(4, nil)
//
// 		idxStor := NewIndexStorageMemory()
// 		storage := NewLifelineStorage(idxStor, idxStor)
// 		err := storage.Set(ctx, 1, *fID, idx)
// 		require.NoError(t, err)
// 		err = storage.Set(ctx, 1, *sID, idx)
// 		require.NoError(t, err)
// 		err = storage.Set(ctx, 2, *tID, idx)
// 		require.NoError(t, err)
// 		err = storage.Set(ctx, 2, *fthID, idx)
// 		require.NoError(t, err)
//
// 		require.Equal(t, 2, len(idxStor.buckets))
// 		require.Equal(t, 2, len(idxStor.buckets[1]))
// 		require.Equal(t, 2, len(idxStor.buckets[2]))
// 		require.Equal(t, *fID, idxStor.buckets[1][*fID].ObjID)
// 		require.Equal(t, *sID, idxStor.buckets[1][*sID].ObjID)
// 		require.Equal(t, *tID, idxStor.buckets[2][*tID].ObjID)
// 		require.Equal(t, *fthID, idxStor.buckets[2][*fthID].ObjID)
// 	})
//
// 	t.Run("override indices is ok", func(t *testing.T) {
// 		t.Parallel()
//
// 		idxStor := NewIndexStorageMemory()
// 		storage := NewLifelineStorage(idxStor, idxStor)
// 		pn := gen.PulseNumber()
//
// 		err := storage.Set(ctx, pn, id, idx)
// 		require.NoError(t, err)
//
// 		err = storage.Set(ctx, pn, id, idx)
// 		require.NoError(t, err)
// 	})
// }
//
// func TestIndexStorage_ForID(t *testing.T) {
// 	t.Parallel()
//
// 	ctx := inslogger.TestContext(t)
//
// 	jetID := gen.JetID()
// 	id := gen.ID()
// 	idx := Lifeline{
// 		LatestState: &id,
// 		JetID:       jetID,
// 		Delegates:   []LifelineDelegate{},
// 	}
//
// 	t.Run("returns correct index-value", func(t *testing.T) {
// 		t.Parallel()
//
// 		idxStor := NewIndexStorageMemory()
// 		storage := NewLifelineStorage(idxStor, idxStor)
// 		pn := gen.PulseNumber()
//
// 		err := storage.Set(ctx, pn, id, idx)
// 		require.NoError(t, err)
//
// 		res, err := storage.ForID(ctx, pn, id)
//
// 		require.NoError(t, err)
// 		require.Equal(t, idx, res)
// 	})
//
// 	t.Run("returns error when no index-value for id", func(t *testing.T) {
// 		t.Parallel()
//
// 		idxStor := NewIndexStorageMemory()
// 		storage := NewLifelineStorage(idxStor, idxStor)
// 		pn := gen.PulseNumber()
//
// 		_, err := storage.ForID(ctx, pn, id)
//
// 		require.Equal(t, ErrLifelineNotFound, err)
// 	})
// }
//
// func TestInMemoryIndex_ForPNAndJet(t *testing.T) {
// 	t.Parallel()
//
// 	ctx := inslogger.TestContext(t)
//
// 	fJetId := insolar.NewJetID(1, []byte{1})
// 	sJetId := insolar.NewJetID(1, []byte{2})
// 	tJetId := insolar.NewJetID(1, []byte{3})
//
// 	fId := insolar.NewID(123, []byte{})
// 	sId := insolar.NewID(124, []byte{})
// 	tId := insolar.NewID(125, []byte{})
//
// 	fPn := insolar.PulseNumber(1)
// 	sPn := insolar.PulseNumber(2)
//
// 	fIdx := Lifeline{
// 		LatestState: insolar.NewID(123, []byte{}),
// 		JetID:       *fJetId,
// 		Delegates:   []LifelineDelegate{},
// 	}
// 	sIdx := Lifeline{
// 		LatestState: insolar.NewID(124, []byte{}),
// 		JetID:       *sJetId,
// 		Delegates:   []LifelineDelegate{},
// 	}
// 	tIdx := Lifeline{
// 		LatestState: insolar.NewID(125, []byte{}),
// 		JetID:       *tJetId,
// 		Delegates:   []LifelineDelegate{},
// 	}
//
// 	idxStor := NewIndexStorageMemory()
// 	storage := NewLifelineStorage(idxStor, idxStor)
//
// 	_ = storage.Set(ctx, fPn, *fId, fIdx)
// 	_ = storage.Set(ctx, fPn, *sId, sIdx)
// 	_ = storage.Set(ctx, sPn, *tId, tIdx)
//
// 	res := idxStor.ForPNAndJet(ctx, fPn, *fJetId)
// 	require.Equal(t, 1, len(res))
// 	require.NotNil(t, res[0].Lifeline)
// 	require.Equal(t, *fId, res[0].ObjID)
// 	require.Equal(t, fIdx.LatestState, res[0].Lifeline.LatestState)
//
// 	res = idxStor.ForPNAndJet(ctx, fPn, *sJetId)
// 	require.Equal(t, 1, len(res))
// 	require.NotNil(t, res[0].Lifeline)
// 	require.Equal(t, *sId, res[0].ObjID)
// 	require.Equal(t, sIdx.LatestState, res[0].Lifeline.LatestState)
// }
//
// func TestInMemoryIndex_SetLifelineUsage(t *testing.T) {
// 	t.Parallel()
//
// 	ctx := inslogger.TestContext(t)
//
// 	jetID := gen.JetID()
// 	id := gen.ID()
// 	idx := Lifeline{
// 		LatestState: &id,
// 		JetID:       jetID,
// 		Delegates:   []LifelineDelegate{},
// 	}
//
// 	pn := gen.PulseNumber()
// 	newPN := pn + 1
//
// 	t.Run("works fine", func(t *testing.T) {
// 		t.Parallel()
//
// 		idxStor := NewIndexStorageMemory()
// 		storage := NewLifelineStorage(idxStor, idxStor)
//
// 		_ = storage.Set(ctx, pn, id, idx)
//
// 		require.Equal(t, pn, idxStor.buckets[pn][id].LifelineLastUsed)
//
// 		idxStor.buckets[newPN] = idxStor.buckets[pn]
//
// 		err := storage.SetLifelineUsage(ctx, newPN, id)
//
// 		require.NoError(t, err)
// 		require.Equal(t, newPN, idxStor.buckets[newPN][id].LifelineLastUsed)
// 	})
//
// 	t.Run("returns ErrLifelineNotFound if no bucket", func(t *testing.T) {
// 		t.Parallel()
//
// 		idxStor := NewIndexStorageMemory()
// 		storage := NewLifelineStorage(idxStor, idxStor)
//
// 		err := storage.SetLifelineUsage(ctx, pn, id)
// 		require.Error(t, ErrLifelineNotFound, err)
// 	})
// }
