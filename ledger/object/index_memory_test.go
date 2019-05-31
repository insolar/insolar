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
		require.NotNil(t, meta.ObjectIndex)

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
	buck := ObjectIndex{
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

		err := index.SetObjectIndex(ctx, pn, buck)
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

		err := index.SetObjectIndex(ctx, pn, buck)
		require.NoError(t, err)

		sLlflID := gen.ID()
		sJetID := gen.JetID()
		sBuck := ObjectIndex{
			ObjID: objID,
			Lifeline: Lifeline{
				LatestState: &sLlflID,
				JetID:       sJetID,
				Delegates:   []LifelineDelegate{},
			},
		}

		err = index.SetObjectIndex(ctx, pn, sBuck)
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

	index.buckets[fPn] = map[insolar.ID]*extendedObjectIndex{}
	index.buckets[sPn] = map[insolar.ID]*extendedObjectIndex{}
	index.buckets[tPn] = map[insolar.ID]*extendedObjectIndex{}

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

		require.Equal(t, insolar.PulseNumber(0), buck.PreviousPendingFilament)

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

func TestInMemoryIndex_Records(t *testing.T) {
	t.Run("err when no lifeline", func(t *testing.T) {
		ctx := inslogger.TestContext(t)
		pn := gen.PulseNumber()
		objID := gen.ID()
		idx := NewInMemoryIndex()

		_, err := idx.Records(ctx, pn, objID)

		require.Error(t, err, ErrLifelineNotFound)
	})

	t.Run("works fine", func(t *testing.T) {
		ctx := inslogger.TestContext(t)
		pn := gen.PulseNumber()
		objID := gen.ID()
		idx := NewInMemoryIndex()
		idx.createBucket(ctx, pn, objID)

		objRef := insolar.NewReference(insolar.ID{}, *insolar.NewID(123, nil))
		req := record.Request{Object: objRef}

		_ = idx.SetRequest(ctx, pn, objID, req)

		data, err := idx.Records(ctx, pn, objID)

		require.NoError(t, err)
		require.Equal(t, 1, len(data))
		require.Equal(t, record.Wrap(req), data[0])
	})
}

func TestInMemoryIndex_OpenRequestsForObjID(t *testing.T) {
	t.Run("err when no lifeline", func(t *testing.T) {
		ctx := inslogger.TestContext(t)
		pn := gen.PulseNumber()
		objID := gen.ID()
		idx := NewInMemoryIndex()

		_, err := idx.OpenRequestsForObjID(ctx, pn, objID, 0)

		require.Error(t, err, ErrLifelineNotFound)
	})

	t.Run("works fine", func(t *testing.T) {
		ctx := inslogger.TestContext(t)
		pn := gen.PulseNumber()
		objID := gen.ID()
		idx := NewInMemoryIndex()

		idx.createBucket(ctx, pn, objID)

		objRef := insolar.NewReference(insolar.ID{}, *insolar.NewID(123, nil))
		req := record.Request{Object: objRef}

		objRefS := insolar.NewReference(insolar.ID{}, *insolar.NewID(234, nil))
		reqS := record.Request{Object: objRefS}

		err := idx.SetRequest(ctx, pn, objID, req)
		require.NoError(t, err)
		err = idx.SetRequest(ctx, pn, objID, reqS)
		require.NoError(t, err)

		t.Run("query all", func(t *testing.T) {
			reqs, err := idx.OpenRequestsForObjID(ctx, pn, objID, 10)
			require.NoError(t, err)
			require.Equal(t, 2, len(reqs))
			require.Equal(t, req, reqs[0])
			require.Equal(t, reqS, reqs[1])
		})

		t.Run("query one", func(t *testing.T) {
			reqs, err := idx.OpenRequestsForObjID(ctx, pn, objID, 1)
			require.NoError(t, err)
			require.Equal(t, 1, len(reqs))
			require.Equal(t, req, reqs[0])
		})
	})
}

func TestInMemoryIndex_MetaForObjID(t *testing.T) {
	t.Run("err when no lifeline", func(t *testing.T) {
		ctx := inslogger.TestContext(t)
		pn := gen.PulseNumber()
		objID := gen.ID()
		idx := NewInMemoryIndex()

		_, err := idx.MetaForObjID(ctx, pn, objID)

		require.Error(t, err, ErrLifelineNotFound)
	})

	t.Run("works fine. first in a row", func(t *testing.T) {
		ctx := inslogger.TestContext(t)
		pn := gen.PulseNumber()
		objID := gen.ID()
		idx := NewInMemoryIndex()
		idx.createBucket(ctx, pn, objID)

		meta, err := idx.MetaForObjID(ctx, pn, objID)
		require.NoError(t, err)

		require.Nil(t, meta.PreviousPN)
		require.Nil(t, meta.ReadUntil)
		require.Equal(t, false, meta.IsStateCalculated)
	})

	t.Run("works fine. second in a row", func(t *testing.T) {
		ctx := inslogger.TestContext(t)
		pn := gen.PulseNumber()
		objID := gen.ID()
		idx := NewInMemoryIndex()
		idx.createBucket(ctx, pn, objID)

		ru := insolar.PulseNumber(888)
		idx.buckets[pn][objID].PreviousPendingFilament = insolar.PulseNumber(888)
		idx.buckets[pn][objID].readPendingUntil = &ru

		meta, err := idx.MetaForObjID(ctx, pn, objID)
		require.NoError(t, err)

		require.NotNil(t, meta.PreviousPN)
		require.NotNil(t, meta.ReadUntil)

		require.Equal(t, false, meta.IsStateCalculated)
		require.Equal(t, insolar.PulseNumber(888), *meta.PreviousPN)
		require.Equal(t, insolar.PulseNumber(888), *meta.ReadUntil)
	})
}

func TestInMemoryIndex_SetReadUntil(t *testing.T) {
	t.Run("err when no lifeline", func(t *testing.T) {
		ctx := inslogger.TestContext(t)
		pn := gen.PulseNumber()
		objID := gen.ID()
		idx := NewInMemoryIndex()

		err := idx.SetReadUntil(ctx, pn, objID, nil)

		require.Error(t, err, ErrLifelineNotFound)
	})

	t.Run("works fine", func(t *testing.T) {
		ctx := inslogger.TestContext(t)
		pn := gen.PulseNumber()
		objID := gen.ID()
		idx := NewInMemoryIndex()
		idx.createBucket(ctx, pn, objID)

		rupn := insolar.PulseNumber(10)
		err := idx.SetReadUntil(ctx, pn, objID, &rupn)

		require.NoError(t, err)
		require.Equal(t, insolar.PulseNumber(10), *idx.buckets[pn][objID].readPendingUntil)
	})
}

func TestInMemoryIndex_SetResult(t *testing.T) {
	t.Run("err when no lifeline", func(t *testing.T) {
		ctx := inslogger.TestContext(t)
		pn := gen.PulseNumber()
		objID := gen.ID()
		idx := NewInMemoryIndex()

		err := idx.SetResult(ctx, pn, objID, record.Result{})

		require.Error(t, err, ErrLifelineNotFound)
	})

	t.Run("set result, when no requests", func(t *testing.T) {
		ctx := inslogger.TestContext(t)
		pn := gen.PulseNumber()
		objID := gen.ID()
		idx := NewInMemoryIndex()
		idx.createBucket(ctx, pn, objID)

		objRef := insolar.NewReference(insolar.ID{}, *insolar.NewID(123, nil))
		res := record.Result{Request: *objRef}

		err := idx.SetResult(ctx, pn, objID, res)

		require.NoError(t, err)
		buck := idx.buckets[pn][objID]

		require.Equal(t, 1, len(buck.fullFilament))
		require.Equal(t, record.Wrap(res), buck.fullFilament[0].Records[0])
		require.Equal(t, pn, buck.fullFilament[0].PN)
		require.Equal(t, 1, len(buck.PendingRecords))
		require.Equal(t, record.Wrap(res), buck.PendingRecords[0])
	})

	t.Run("set 2 results, when no requests", func(t *testing.T) {
		ctx := inslogger.TestContext(t)
		pn := gen.PulseNumber()
		objID := gen.ID()
		idx := NewInMemoryIndex()
		idx.createBucket(ctx, pn, objID)

		objRef := insolar.NewReference(insolar.ID{}, *insolar.NewID(123, nil))
		res := record.Result{Request: *objRef}
		resS := record.Result{Request: *objRef, Payload: []byte{1, 2, 3, 4, 5, 6}}

		err := idx.SetResult(ctx, pn, objID, res)
		require.NoError(t, err)
		err = idx.SetResult(ctx, pn, objID, resS)
		require.NoError(t, err)

		buck := idx.buckets[pn][objID]

		require.Equal(t, 1, len(buck.fullFilament))

		require.Equal(t, 2, len(buck.fullFilament[0].Records))
		require.Equal(t, pn, buck.fullFilament[0].PN)
		require.Equal(t, record.Wrap(res), buck.fullFilament[0].Records[0])
		require.Equal(t, record.Wrap(resS), buck.fullFilament[0].Records[1])

		require.Equal(t, 2, len(buck.PendingRecords))
		require.Equal(t, record.Wrap(res), buck.PendingRecords[0])
		require.Equal(t, record.Wrap(resS), buck.PendingRecords[1])
	})

	t.Run("close requests work fine", func(t *testing.T) {
		ctx := inslogger.TestContext(t)
		pn := gen.PulseNumber()
		objID := gen.ID()
		idx := NewInMemoryIndex()
		idx.createBucket(ctx, pn, objID)

		objRef := insolar.NewReference(insolar.ID{}, *insolar.NewID(123, nil))
		req := record.Request{Object: objRef}
		_ = idx.SetRequest(ctx, pn, objID, req)

		objRefS := insolar.NewReference(insolar.ID{}, *insolar.NewID(321, nil))
		reqS := record.Request{Object: objRefS}
		_ = idx.SetRequest(ctx, pn, objID, reqS)

		res := record.Result{Request: *objRef}

		err := idx.SetResult(ctx, pn, objID, res)
		require.NoError(t, err)

		open, err := idx.OpenRequestsForObjID(ctx, pn, objID, 10)
		require.NoError(t, err)

		require.Equal(t, 1, len(open))
		require.Equal(t, reqS, open[0])

		buck := idx.buckets[pn][objID]
		require.Equal(t, 1, len(buck.fullFilament))
		require.Equal(t, pn, buck.fullFilament[0].PN)
		require.Equal(t, 3, len(buck.fullFilament[0].Records))

		require.Equal(t, 1, len(buck.notClosedRequestsIndex[pn]))
		_, ok := buck.notClosedRequestsIndex[pn][*reqS.Object.Record()]
		require.Equal(t, true, ok)

		require.Equal(t, 1, len(buck.notClosedRequests))
		require.Equal(t, reqS, buck.notClosedRequests[0])

		require.Equal(t, 2, len(buck.requestPNIndex))
	})

	t.Run("set result, other there are other fillaments", func(t *testing.T) {
		ctx := inslogger.TestContext(t)
		pn := gen.PulseNumber()
		objID := gen.ID()
		idx := NewInMemoryIndex()
		idx.createBucket(ctx, pn, objID)
		buck := idx.buckets[pn][objID]
		buck.fullFilament = append(buck.fullFilament, chainLink{PN: pn + 1, Records: []record.Virtual{}})

		objRef := insolar.NewReference(insolar.ID{}, *insolar.NewID(123, nil))
		res := record.Result{Request: *objRef}

		err := idx.SetResult(ctx, pn, objID, res)

		require.NoError(t, err)

		require.Equal(t, 2, len(buck.fullFilament))
		require.Equal(t, pn, buck.fullFilament[0].PN)
		require.Equal(t, pn+1, buck.fullFilament[1].PN)
	})
}

func TestInMemoryIndex_RefreshState(t *testing.T) {
	t.Run("err when no lifeline", func(t *testing.T) {
		ctx := inslogger.TestContext(t)
		pn := gen.PulseNumber()
		objID := gen.ID()
		idx := NewInMemoryIndex()

		err := idx.RefreshState(ctx, pn, objID)

		require.Error(t, err, ErrLifelineNotFound)
	})

	t.Run("works fine. req and res", func(t *testing.T) {
		ctx := inslogger.TestContext(t)
		pn := gen.PulseNumber()
		objID := gen.ID()
		idx := NewInMemoryIndex()
		idx.createBucket(ctx, pn, objID)
		buck := idx.buckets[pn][objID]
		buck.notClosedRequestsIndex = map[insolar.PulseNumber]map[insolar.ID]*record.Request{
			pn + 1: {},
			pn:     {},
		}
		buck.requestPNIndex = map[insolar.ID]insolar.PulseNumber{}

		objRef := insolar.NewReference(insolar.ID{}, *insolar.NewID(123, nil))
		req := record.Request{Object: objRef}
		buck.fullFilament = append(buck.fullFilament, chainLink{PN: pn + 1, Records: []record.Virtual{record.Wrap(req)}})

		res := record.Result{Request: *objRef}
		buck.fullFilament = append(buck.fullFilament, chainLink{PN: pn, Records: []record.Virtual{record.Wrap(res)}})

		err := idx.RefreshState(ctx, pn, objID)
		require.NoError(t, err)

		require.Equal(t, true, buck.isStateCalculated)
		require.Equal(t, 2, len(buck.fullFilament))
		require.Equal(t, 0, len(buck.notClosedRequests))
		require.Equal(t, 0, len(buck.notClosedRequestsIndex[pn]))
		require.Equal(t, 0, len(buck.notClosedRequestsIndex[pn+1]))
	})

	t.Run("works fine. req and res", func(t *testing.T) {
		ctx := inslogger.TestContext(t)
		pn := gen.PulseNumber()
		objID := gen.ID()
		idx := NewInMemoryIndex()
		idx.createBucket(ctx, pn, objID)
		buck := idx.buckets[pn][objID]
		buck.notClosedRequestsIndex = map[insolar.PulseNumber]map[insolar.ID]*record.Request{
			pn + 1: {},
			pn:     {},
		}
		buck.requestPNIndex = map[insolar.ID]insolar.PulseNumber{}

		objRef := insolar.NewReference(insolar.ID{}, *insolar.NewID(123, nil))
		req := record.Request{Object: objRef}
		buck.fullFilament = append(buck.fullFilament, chainLink{PN: pn + 1, Records: []record.Virtual{record.Wrap(req)}})

		res := record.Result{Request: *objRef}
		buck.fullFilament = append(buck.fullFilament, chainLink{PN: pn, Records: []record.Virtual{record.Wrap(res)}})

		err := idx.RefreshState(ctx, pn, objID)
		require.NoError(t, err)

		require.Equal(t, true, buck.isStateCalculated)
		require.Equal(t, 2, len(buck.fullFilament))
		require.Equal(t, 0, len(buck.notClosedRequests))
		require.Equal(t, 0, len(buck.notClosedRequestsIndex[pn]))
		require.Equal(t, 0, len(buck.notClosedRequestsIndex[pn+1]))
	})

	t.Run("works fine. open pending", func(t *testing.T) {
		ctx := inslogger.TestContext(t)
		pn := gen.PulseNumber()
		objID := gen.ID()
		idx := NewInMemoryIndex()
		idx.createBucket(ctx, pn, objID)
		buck := idx.buckets[pn][objID]
		buck.notClosedRequestsIndex = map[insolar.PulseNumber]map[insolar.ID]*record.Request{
			pn: {},
		}
		buck.requestPNIndex = map[insolar.ID]insolar.PulseNumber{}

		objRef := insolar.NewReference(insolar.ID{}, *insolar.NewID(123, nil))
		req := record.Request{Object: objRef}
		buck.fullFilament = append(buck.fullFilament, chainLink{PN: pn, Records: []record.Virtual{record.Wrap(req)}})

		err := idx.RefreshState(ctx, pn, objID)
		require.NoError(t, err)

		require.Equal(t, true, buck.isStateCalculated)
		require.Equal(t, 1, len(buck.fullFilament))
		require.Equal(t, 1, len(buck.notClosedRequests))
		require.Equal(t, req, buck.notClosedRequests[0])
		require.Equal(t, 1, len(buck.notClosedRequestsIndex[pn]))
	})

	t.Run("calculates readPendingUntil properly", func(t *testing.T) {
		ctx := inslogger.TestContext(t)
		pn := gen.PulseNumber()
		objID := gen.ID()
		idx := NewInMemoryIndex()
		idx.createBucket(ctx, pn, objID)
		buck := idx.buckets[pn][objID]

		objRef := insolar.NewReference(insolar.ID{}, *insolar.NewID(123, nil))
		req := record.Request{Object: objRef}
		res := record.Result{Request: *objRef}
		buck.fullFilament = append(buck.fullFilament, chainLink{PN: pn, Records: []record.Virtual{record.Wrap(req), record.Wrap(res)}})

		objRefS := insolar.NewReference(insolar.ID{}, *insolar.NewID(567, nil))
		reqS := record.Request{Object: objRefS}
		buck.fullFilament = append(buck.fullFilament, chainLink{PN: pn + 1, Records: []record.Virtual{record.Wrap(reqS)}})

		objRefT := insolar.NewReference(insolar.ID{}, *insolar.NewID(888, nil))
		reqT := record.Request{Object: objRefT}
		resT := record.Result{Request: *objRefT}
		buck.fullFilament = append(buck.fullFilament, chainLink{PN: pn + 2, Records: []record.Virtual{record.Wrap(reqT), record.Wrap(resT)}})

		buck.notClosedRequestsIndex = map[insolar.PulseNumber]map[insolar.ID]*record.Request{
			pn:     {},
			pn + 1: {},
			pn + 2: {},
		}
		buck.requestPNIndex = map[insolar.ID]insolar.PulseNumber{}

		err := idx.RefreshState(ctx, pn, objID)
		require.NoError(t, err)

		require.Equal(t, pn, *buck.readPendingUntil)
	})
}
