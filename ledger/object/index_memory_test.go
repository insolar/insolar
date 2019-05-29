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
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/instrumentation/inslogger"
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

		err := storage.Set(ctx, pn, id, idx)

		require.NoError(t, err)
		assert.Equal(t, 1, len(storage.buckets))

		buck, buckOK := storage.buckets[pn]
		require.Equal(t, true, buckOK)
		require.Equal(t, 1, len(buck))

		meta, metaOK := buck[id]
		require.Equal(t, true, metaOK)
		require.NotNil(t, meta)
		require.NotNil(t, meta.IndexBucket)

		require.Equal(t, meta.Lifeline, idx)
		require.Equal(t, meta.LifelineLastUsed, pn)
		require.Equal(t, meta.ObjID, id)
	})

	t.Run("save multiple values", func(t *testing.T) {
		fID := insolar.NewID(1, nil)
		sID := insolar.NewID(2, nil)
		tID := insolar.NewID(3, nil)
		fthID := insolar.NewID(4, nil)

		storage := NewInMemoryIndex()
		err := storage.Set(ctx, 1, *fID, idx)
		require.NoError(t, err)
		err = storage.Set(ctx, 1, *sID, idx)
		require.NoError(t, err)
		err = storage.Set(ctx, 2, *tID, idx)
		require.NoError(t, err)
		err = storage.Set(ctx, 2, *fthID, idx)
		require.NoError(t, err)

		require.Equal(t, 2, len(storage.buckets))
		require.Equal(t, 2, len(storage.buckets[1]))
		require.Equal(t, 2, len(storage.buckets[2]))
		require.Equal(t, *fID, storage.buckets[1][*fID].ObjID)
		require.Equal(t, *sID, storage.buckets[1][*sID].ObjID)
		require.Equal(t, *tID, storage.buckets[2][*tID].ObjID)
		require.Equal(t, *fthID, storage.buckets[2][*fthID].ObjID)
	})

	t.Run("override indices is ok", func(t *testing.T) {
		t.Parallel()

		storage := NewInMemoryIndex()
		pn := gen.PulseNumber()

		err := storage.Set(ctx, pn, id, idx)
		require.NoError(t, err)

		err = storage.Set(ctx, pn, id, idx)
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

		err := storage.Set(ctx, pn, id, idx)
		require.NoError(t, err)

		res, err := storage.ForID(ctx, pn, id)

		require.NoError(t, err)
		assert.Equal(t, idx, res)
	})

	t.Run("returns error when no index-value for id", func(t *testing.T) {
		t.Parallel()

		storage := NewInMemoryIndex()
		pn := gen.PulseNumber()

		_, err := storage.ForID(ctx, pn, id)

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

	_ = index.Set(ctx, fPn, *fId, fIdx)
	_ = index.Set(ctx, fPn, *sId, sIdx)
	_ = index.Set(ctx, sPn, *tId, tIdx)

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
		savedBuckBuf, _ := savedBuck.Marshal()

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
		savedBuckBuf, _ := savedBuck.Marshal()

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

		_ = index.Set(ctx, pn, id, idx)

		require.Equal(t, pn, index.buckets[pn][id].LifelineLastUsed)

		index.buckets[newPN] = index.buckets[pn]

		err := index.SetLifelineUsage(ctx, newPN, id)

		require.NoError(t, err)
		require.Equal(t, newPN, index.buckets[newPN][id].LifelineLastUsed)
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

	index.buckets[fPn] = map[insolar.ID]*extendedIndexBucket{}
	index.buckets[sPn] = map[insolar.ID]*extendedIndexBucket{}
	index.buckets[tPn] = map[insolar.ID]*extendedIndexBucket{}

	index.DeleteForPN(ctx, sPn)

	_, ok := index.buckets[fPn]
	require.Equal(t, true, ok)
	_, ok = index.buckets[sPn]
	require.Equal(t, false, ok)
	_, ok = index.buckets[tPn]
	require.Equal(t, true, ok)
}

func TestInMemoryIndex_SetRequest(t *testing.T) {
	t.Run("err when no lifeline", func(t *testing.T) {
		ctx := inslogger.TestContext(t)
		pn := gen.PulseNumber()
		objID := gen.ID()
		idx := NewInMemoryIndex()

		err := idx.SetRequest(ctx, pn, objID, record.Request{})

		require.Error(t, err, ErrLifelineNotFound)
	})

	t.Run("set first request on the object", func(t *testing.T) {
		ctx := inslogger.TestContext(t)
		pn := gen.PulseNumber()
		objID := gen.ID()
		idx := NewInMemoryIndex()
		idx.createBucket(ctx, pn, objID)

		objRef := gen.Reference()
		req := record.Request{Object: &objRef}

		err := idx.SetRequest(ctx, pn, objID, req)
		require.NoError(t, err)

		buck := idx.buckets[pn][objID]

		require.Equal(t, pn, buck.PreviousPendingFilament)

		require.Equal(t, 1, len(buck.PendingRecords))
		require.Equal(t, 1, len(buck.fullFilament))
		require.Equal(t, 1, len(buck.fullFilament[0].Records))

		require.Equal(t, record.Wrap(req), buck.PendingRecords[0])
		require.Equal(t, pn, buck.fullFilament[0].PN)
		require.Equal(t, record.Wrap(req), buck.fullFilament[0].Records[0])

		require.Equal(t, 1, len(buck.requestPNIndex))
		require.Equal(t, 1, len(buck.notClosedRequestsIndex))
		require.Equal(t, 1, len(buck.notClosedRequests))

		require.Equal(t, pn, buck.requestPNIndex[*req.Object.Record()])
		require.Equal(t, req, *buck.notClosedRequestsIndex[pn][*req.Object.Record()])
		require.Equal(t, req, buck.notClosedRequests[0])
	})

	t.Run("set two request on the object", func(t *testing.T) {
		ctx := inslogger.TestContext(t)
		pn := gen.PulseNumber()
		objID := gen.ID()
		idx := NewInMemoryIndex()
		idx.createBucket(ctx, pn, objID)

		objRef := insolar.NewReference(insolar.ID{}, *insolar.NewID(123, nil))
		req := record.Request{Object: objRef}

		err := idx.SetRequest(ctx, pn, objID, req)
		require.NoError(t, err)

		objRefS := insolar.NewReference(insolar.ID{}, *insolar.NewID(321, nil))
		reqS := record.Request{Object: objRefS}

		err = idx.SetRequest(ctx, pn, objID, reqS)
		require.NoError(t, err)

		buck := idx.buckets[pn][objID]
		require.Equal(t, 2, len(buck.PendingRecords))
		require.Equal(t, 1, len(buck.fullFilament))
		require.Equal(t, 2, len(buck.fullFilament[0].Records))

		require.Equal(t, record.Wrap(req), buck.PendingRecords[0])
		require.Equal(t, record.Wrap(reqS), buck.PendingRecords[1])
		require.Equal(t, pn, buck.fullFilament[0].PN)
		require.Equal(t, record.Wrap(req), buck.fullFilament[0].Records[0])
		require.Equal(t, record.Wrap(reqS), buck.fullFilament[0].Records[1])

		require.Equal(t, 2, len(buck.requestPNIndex))
		require.Equal(t, 1, len(buck.notClosedRequestsIndex))
		require.Equal(t, 2, len(buck.notClosedRequests))

		require.Equal(t, pn, buck.requestPNIndex[*req.Object.Record()])
		require.Equal(t, req, *buck.notClosedRequestsIndex[pn][*req.Object.Record()])
		require.Equal(t, req, buck.notClosedRequests[0])

		require.Equal(t, pn, buck.requestPNIndex[*reqS.Object.Record()])
		require.Equal(t, reqS, *buck.notClosedRequestsIndex[pn][*reqS.Object.Record()])
		require.Equal(t, reqS, buck.notClosedRequests[1])
	})

	t.Run("test rebalanced fillaments buckets list", func(t *testing.T) {
		ctx := inslogger.TestContext(t)
		pn := insolar.PulseNumber(123)
		objID := gen.ID()
		idx := NewInMemoryIndex()
		idx.createBucket(ctx, pn, objID)

		buck := idx.buckets[pn][objID]
		buck.fullFilament = append(buck.fullFilament, chainLink{PN: pn + 1, Records: []record.Virtual{}})

		objRef := gen.Reference()
		req := record.Request{Object: &objRef}

		err := idx.SetRequest(ctx, pn, objID, req)
		require.NoError(t, err)

		require.Equal(t, 2, len(buck.fullFilament))
		require.Equal(t, pn, buck.fullFilament[0].PN)
		require.Equal(t, pn+1, buck.fullFilament[1].PN)
	})

}

func TestInMemoryIndex_SetFilament(t *testing.T) {
	t.Run("err when no lifeline", func(t *testing.T) {
		ctx := inslogger.TestContext(t)
		pn := gen.PulseNumber()
		objID := gen.ID()
		idx := NewInMemoryIndex()

		err := idx.SetFilament(ctx, pn, objID, gen.PulseNumber(), nil)

		require.Error(t, err, ErrLifelineNotFound)
	})

	t.Run("works fine", func(t *testing.T) {
		ctx := inslogger.TestContext(t)
		pn := insolar.PulseNumber(123)
		objID := gen.ID()
		idx := NewInMemoryIndex()
		idx.createBucket(ctx, pn, objID)

		buck := idx.buckets[pn][objID]
		buck.fullFilament = append(buck.fullFilament, chainLink{PN: pn + 1, Records: []record.Virtual{}})
		buck.fullFilament = append(buck.fullFilament, chainLink{PN: pn - 10, Records: []record.Virtual{}})

		objRef := gen.Reference()
		req := record.Request{Object: &objRef}

		err := idx.SetFilament(ctx, pn, objID, pn, []record.Virtual{record.Wrap(req)})
		require.NoError(t, err)

		require.Equal(t, 3, len(buck.fullFilament))
		require.Equal(t, pn-10, buck.fullFilament[0].PN)
		require.Equal(t, pn, buck.fullFilament[1].PN)
		require.Equal(t, pn+1, buck.fullFilament[2].PN)

		require.Equal(t, record.Wrap(req), buck.fullFilament[1].Records[0])
	})
}
