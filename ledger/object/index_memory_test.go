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
	"context"
	"testing"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/platformpolicy"
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

		storage := NewInMemoryIndex(nil, nil)
		pn := gen.PulseNumber()

		err := storage.Set(ctx, pn, id, idx)

		require.NoError(t, err)
		require.Equal(t, 1, len(storage.buckets))

		buck, buckOK := storage.buckets[pn]
		require.Equal(t, true, buckOK)
		require.Equal(t, 1, len(buck))

		meta, metaOK := buck[id]
		require.Equal(t, true, metaOK)
		require.NotNil(t, meta)
		require.NotNil(t, meta.objectMeta)

		require.Equal(t, meta.objectMeta.Lifeline, idx)
		require.Equal(t, meta.objectMeta.LifelineLastUsed, pn)
		require.Equal(t, meta.objectMeta.ObjID, id)
	})

	t.Run("save multiple values", func(t *testing.T) {
		fID := insolar.NewID(1, nil)
		sID := insolar.NewID(2, nil)
		tID := insolar.NewID(3, nil)
		fthID := insolar.NewID(4, nil)

		storage := NewInMemoryIndex(nil, nil)
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
		require.Equal(t, *fID, storage.buckets[1][*fID].objectMeta.ObjID)
		require.Equal(t, *sID, storage.buckets[1][*sID].objectMeta.ObjID)
		require.Equal(t, *tID, storage.buckets[2][*tID].objectMeta.ObjID)
		require.Equal(t, *fthID, storage.buckets[2][*fthID].objectMeta.ObjID)
	})

	t.Run("override indices is ok", func(t *testing.T) {
		t.Parallel()

		storage := NewInMemoryIndex(nil, nil)
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

		storage := NewInMemoryIndex(nil, nil)
		pn := gen.PulseNumber()

		err := storage.Set(ctx, pn, id, idx)
		require.NoError(t, err)

		res, err := storage.ForID(ctx, pn, id)

		require.NoError(t, err)
		require.Equal(t, idx, res)
	})

	t.Run("returns error when no index-value for id", func(t *testing.T) {
		t.Parallel()

		storage := NewInMemoryIndex(nil, nil)
		pn := gen.PulseNumber()

		_, err := storage.ForID(ctx, pn, id)

		require.Equal(t, ErrLifelineNotFound, err)
	})
}

func TestInMemoryIndex_ForPulse(t *testing.T) {
	t.Skip() // waiting filaments by @egorikas (INS-1939)

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

	index := NewInMemoryIndex(nil, nil)

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
	buck := FilamentIndex{
		ObjID: objID,
		Lifeline: Lifeline{
			LatestState: &lflID,
			JetID:       jetID,
			Delegates:   []LifelineDelegate{},
		},
	}

	t.Run("saves correct bucket", func(t *testing.T) {
		pn := gen.PulseNumber()
		index := NewInMemoryIndex(nil, nil)

		err := index.SetBucket(ctx, pn, buck)
		require.NoError(t, err)

		savedBuck := index.buckets[pn][objID].objectMeta
		require.NotNil(t, savedBuck)

		buckBuf, _ := buck.Marshal()
		savedBuckBuf, _ := savedBuck.Marshal()

		require.Equal(t, buckBuf, savedBuckBuf)
	})

	t.Run("re-save works fine", func(t *testing.T) {
		pn := gen.PulseNumber()
		index := NewInMemoryIndex(nil, nil)

		err := index.SetBucket(ctx, pn, buck)
		require.NoError(t, err)

		sLlflID := gen.ID()
		sJetID := gen.JetID()
		sBuck := FilamentIndex{
			ObjID: objID,
			Lifeline: Lifeline{
				LatestState: &sLlflID,
				JetID:       sJetID,
				Delegates:   []LifelineDelegate{},
			},
		}

		err = index.SetBucket(ctx, pn, sBuck)
		require.NoError(t, err)

		savedBuck := index.buckets[pn][objID].objectMeta
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

		index := NewInMemoryIndex(nil, nil)

		_ = index.Set(ctx, pn, id, idx)

		require.Equal(t, pn, index.buckets[pn][id].objectMeta.LifelineLastUsed)

		index.buckets[newPN] = index.buckets[pn]

		err := index.SetLifelineUsage(ctx, newPN, id)

		require.NoError(t, err)
		require.Equal(t, newPN, index.buckets[newPN][id].objectMeta.LifelineLastUsed)
	})

	t.Run("returns ErrLifelineNotFound if no bucket", func(t *testing.T) {
		t.Parallel()

		index := NewInMemoryIndex(nil, nil)
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

	index := NewInMemoryIndex(nil, nil)

	index.buckets[fPn] = map[insolar.ID]*filamentCache{}
	index.buckets[sPn] = map[insolar.ID]*filamentCache{}
	index.buckets[tPn] = map[insolar.ID]*filamentCache{}

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
		idx := NewInMemoryIndex(nil, platformpolicy.NewPlatformCryptographyScheme())

		err := idx.SetRequest(ctx, pn, objID, insolar.ID{})

		require.Error(t, err, ErrLifelineNotFound)
	})

	t.Run("set first request on the object", func(t *testing.T) {
		ctx := inslogger.TestContext(t)
		pn := gen.PulseNumber()
		objID := gen.ID()

		reqID := insolar.NewID(1, []byte{1})
		prevPending := gen.ID()

		pf := record.PendingFilament{
			RecordID:       *reqID,
			PreviousRecord: &prevPending,
		}

		pfv := record.Wrap(pf)
		hash := record.HashVirtual(platformpolicy.NewPlatformCryptographyScheme().ReferenceHasher(), pfv)
		metaReqID := *insolar.NewID(pn, hash)

		rsm := NewRecordStorageMock(t)
		rsm.SetFunc = func(p context.Context, p1 insolar.ID, p2 record.Material) (r error) {
			require.Equal(t, metaReqID, p1)
			cast, ok := record.Unwrap(p2.Virtual).(*record.PendingFilament)
			require.Equal(t, true, ok)
			require.Equal(t, *reqID, cast.RecordID)
			require.Equal(t, prevPending, *cast.PreviousRecord)
			return nil
		}

		idx := NewInMemoryIndex(rsm, platformpolicy.NewPlatformCryptographyScheme())
		idx.createBucket(ctx, pn, objID)
		idx.buckets[pn][objID].objectMeta.Lifeline.PendingPointer = &prevPending

		err := idx.SetRequest(ctx, pn, objID, *reqID)
		require.NoError(t, err)

		buck := idx.buckets[pn][objID]

		require.NotEqual(t, prevPending, *buck.objectMeta.Lifeline.PendingPointer)

		require.Equal(t, 1, len(buck.objectMeta.PendingRecords))
		require.Equal(t, 1, len(buck.pendingMeta.fullFilament))
		require.Equal(t, 1, len(buck.pendingMeta.fullFilament[0].MetaRecordsIDs))

		require.Equal(t, metaReqID, buck.objectMeta.PendingRecords[0])
		require.Equal(t, pn, buck.pendingMeta.fullFilament[0].PN)
		require.Equal(t, metaReqID, buck.pendingMeta.fullFilament[0].MetaRecordsIDs[0])

		require.Equal(t, 1, len(buck.pendingMeta.notClosedRequestsIdsIndex))
		require.Equal(t, 1, len(buck.pendingMeta.notClosedRequestsIds))

		_, ok := buck.pendingMeta.notClosedRequestsIdsIndex[pn][*reqID]
		require.Equal(t, true, ok)
		require.Equal(t, *reqID, buck.pendingMeta.notClosedRequestsIds[0])
	})

	t.Run("set two request on the object", func(t *testing.T) {
		ctx := inslogger.TestContext(t)
		pn := gen.PulseNumber()
		objID := gen.ID()

		reqID := insolar.NewID(1, []byte{1})
		firstPending := gen.ID()

		pf := record.PendingFilament{
			RecordID:       *reqID,
			PreviousRecord: &firstPending,
		}

		pfv := record.Wrap(pf)
		hash := record.HashVirtual(platformpolicy.NewPlatformCryptographyScheme().ReferenceHasher(), pfv)
		firstMeta := *insolar.NewID(pn, hash)

		reqSID := insolar.NewID(2, []byte{2})
		pfs := record.PendingFilament{
			RecordID:       *reqSID,
			PreviousRecord: &firstMeta,
		}
		pfsv := record.Wrap(pfs)
		hash = record.HashVirtual(platformpolicy.NewPlatformCryptographyScheme().ReferenceHasher(), pfsv)
		secondMeta := *insolar.NewID(pn, hash)

		rsm := NewRecordStorageMock(t)
		rsm.SetFunc = func(p context.Context, p1 insolar.ID, p2 record.Material) (r error) {
			switch p1 {
			case firstMeta:
				concrete, ok := record.Unwrap(p2.Virtual).(*record.PendingFilament)
				require.Equal(t, true, ok)
				require.Equal(t, *reqID, concrete.RecordID)
				require.Equal(t, &firstPending, concrete.PreviousRecord)
			case secondMeta:
				concrete, ok := record.Unwrap(p2.Virtual).(*record.PendingFilament)
				require.Equal(t, true, ok)
				require.Equal(t, *reqSID, concrete.RecordID)
				require.Equal(t, &firstMeta, concrete.PreviousRecord)
			default:
				t.Fatal("test is totally broken")
			}
			return nil
		}

		idx := NewInMemoryIndex(rsm, platformpolicy.NewPlatformCryptographyScheme())
		idx.createBucket(ctx, pn, objID)
		buck := idx.buckets[pn][objID]
		buck.objectMeta.Lifeline.PendingPointer = &firstPending

		err := idx.SetRequest(ctx, pn, objID, *reqID)
		require.NoError(t, err)

		err = idx.SetRequest(ctx, pn, objID, *reqSID)
		require.NoError(t, err)

		require.Equal(t, 2, len(buck.objectMeta.PendingRecords))
		require.Equal(t, 1, len(buck.pendingMeta.fullFilament))
		require.Equal(t, 2, len(buck.pendingMeta.fullFilament[0].MetaRecordsIDs))

		require.Equal(t, firstMeta, buck.objectMeta.PendingRecords[0])
		require.Equal(t, secondMeta, buck.objectMeta.PendingRecords[1])
		require.Equal(t, pn, buck.pendingMeta.fullFilament[0].PN)
		require.Equal(t, firstMeta, buck.pendingMeta.fullFilament[0].MetaRecordsIDs[0])
		require.Equal(t, secondMeta, buck.pendingMeta.fullFilament[0].MetaRecordsIDs[1])

		require.Equal(t, 1, len(buck.pendingMeta.notClosedRequestsIdsIndex))
		require.Equal(t, 2, len(buck.pendingMeta.notClosedRequestsIds))

		_, ok := buck.pendingMeta.notClosedRequestsIdsIndex[pn][*reqID]
		require.Equal(t, true, ok)
		require.Equal(t, *reqID, buck.pendingMeta.notClosedRequestsIds[0])

		_, ok = buck.pendingMeta.notClosedRequestsIdsIndex[pn][*reqSID]
		require.Equal(t, true, ok)
		require.Equal(t, *reqSID, buck.pendingMeta.notClosedRequestsIds[1])
	})

	t.Run("test rebalanced fillaments buckets list", func(t *testing.T) {
		ctx := inslogger.TestContext(t)
		pn := insolar.PulseNumber(123)
		objID := gen.ID()

		reqID := insolar.NewID(1, []byte{1})

		rsm := NewRecordStorageMock(t)
		rsm.SetMock.Return(nil)

		idx := NewInMemoryIndex(rsm, platformpolicy.NewPlatformCryptographyScheme())
		idx.createBucket(ctx, pn, objID)

		buck := idx.buckets[pn][objID]
		buck.pendingMeta.fullFilament = append(buck.pendingMeta.fullFilament, chainLink{PN: pn + 1, MetaRecordsIDs: []insolar.ID{}})

		err := idx.SetRequest(ctx, pn, objID, *reqID)
		require.NoError(t, err)

		require.Equal(t, 2, len(buck.pendingMeta.fullFilament))
		require.Equal(t, pn, buck.pendingMeta.fullFilament[0].PN)
		require.Equal(t, pn+1, buck.pendingMeta.fullFilament[1].PN)
	})
}

func TestInMemoryIndex_SetFilament(t *testing.T) {
	t.Run("err when no lifeline", func(t *testing.T) {
		ctx := inslogger.TestContext(t)
		pn := gen.PulseNumber()
		objID := gen.ID()
		idx := NewInMemoryIndex(nil, platformpolicy.NewPlatformCryptographyScheme())

		err := idx.SetFilament(ctx, pn, objID, gen.PulseNumber(), nil)

		require.Error(t, err, ErrLifelineNotFound)
	})

	t.Run("works fine", func(t *testing.T) {
		ctx := inslogger.TestContext(t)
		pn := insolar.PulseNumber(123)
		objID := gen.ID()

		rsm := NewRecordStorageMock(t)

		reqID := *insolar.NewID(222, nil)
		reqRec := record.Request{Object: insolar.NewReference(reqID)}
		reqRecV := record.Wrap(reqRec)
		pf := record.PendingFilament{
			RecordID: reqID,
		}
		pfv := record.Wrap(pf)
		hash := record.HashVirtual(platformpolicy.NewPlatformCryptographyScheme().ReferenceHasher(), pfv)
		firstMeta := *insolar.NewID(pn, hash)

		rsm.SetFunc = func(p context.Context, p1 insolar.ID, p2 record.Material) (r error) {
			switch p1 {
			case firstMeta:
				concrete, ok := record.Unwrap(p2.Virtual).(*record.PendingFilament)
				require.Equal(t, true, ok)
				require.Equal(t, pf, *concrete)
			case reqID:
				concrete, ok := record.Unwrap(p2.Virtual).(*record.Request)
				require.Equal(t, true, ok)
				require.Equal(t, reqRec, *concrete)
			default:
				t.Fatal("test is totally broken")
			}

			return nil
		}

		idx := NewInMemoryIndex(rsm, platformpolicy.NewPlatformCryptographyScheme())
		idx.createBucket(ctx, pn, objID)

		buck := idx.buckets[pn][objID]
		buck.pendingMeta.fullFilament = append(buck.pendingMeta.fullFilament, chainLink{PN: pn + 1, MetaRecordsIDs: []insolar.ID{}})
		buck.pendingMeta.fullFilament = append(buck.pendingMeta.fullFilament, chainLink{PN: pn - 10, MetaRecordsIDs: []insolar.ID{}})

		fill := []record.CompositeFilamentRecord{{MetaID: firstMeta, Meta: pf, RecordID: reqID, Record: record.Material{Virtual: &reqRecV}}}
		err := idx.SetFilament(ctx, pn, objID, pn, fill)
		require.NoError(t, err)

		require.Equal(t, 3, len(buck.pendingMeta.fullFilament))
		require.Equal(t, pn-10, buck.pendingMeta.fullFilament[0].PN)
		require.Equal(t, pn, buck.pendingMeta.fullFilament[1].PN)
		require.Equal(t, pn+1, buck.pendingMeta.fullFilament[2].PN)

		require.Equal(t, firstMeta, buck.pendingMeta.fullFilament[1].MetaRecordsIDs[0])
		rsm.MinimockFinish()
	})
}

func TestInMemoryIndex_Records(t *testing.T) {
	t.Run("err when no lifeline", func(t *testing.T) {
		ctx := inslogger.TestContext(t)
		pn := gen.PulseNumber()
		objID := gen.ID()
		idx := NewInMemoryIndex(nil, platformpolicy.NewPlatformCryptographyScheme())

		_, err := idx.Records(ctx, pn, objID)

		require.Error(t, err, ErrLifelineNotFound)
	})

	t.Run("works fine", func(t *testing.T) {
		ctx := inslogger.TestContext(t)
		pn := gen.PulseNumber()
		objID := gen.ID()

		objRef := insolar.NewReference(*insolar.NewID(123, nil))
		req := record.Request{Object: objRef}
		reqV := record.Wrap(req)
		reqID := insolar.NewID(444, nil)
		metaReq := record.PendingFilament{RecordID: *reqID}
		metaReqV := record.Wrap(metaReq)

		rsm := NewRecordStorageMock(t)
		var savedReqID insolar.ID
		rsm.SetFunc = func(p context.Context, p1 insolar.ID, p2 record.Material) (r error) {
			savedReqID = p1
			return nil
		}
		rsm.ForIDFunc = func(p context.Context, p1 insolar.ID) (r record.Material, r1 error) {
			switch p1 {
			case savedReqID:
				return record.Material{
					Virtual: &metaReqV,
				}, nil
			case *reqID:
				return record.Material{
					Virtual: &reqV,
				}, nil
			default:
				panic("everything is broken")
			}
		}

		idx := NewInMemoryIndex(rsm, platformpolicy.NewPlatformCryptographyScheme())
		idx.createBucket(ctx, pn, objID)

		_ = idx.SetRequest(ctx, pn, objID, *reqID)

		data, err := idx.Records(ctx, pn, objID)

		require.NoError(t, err)
		require.Equal(t, 1, len(data))
		require.Equal(t, record.Wrap(req), *data[0].Record.Virtual)
		require.Equal(t, *reqID, data[0].RecordID)
		require.Equal(t, metaReq, data[0].Meta)
		require.Equal(t, savedReqID, data[0].MetaID)
	})
}

func TestInMemoryIndex_OpenRequestsForObjID(t *testing.T) {
	t.Run("err when no lifeline", func(t *testing.T) {
		ctx := inslogger.TestContext(t)
		pn := gen.PulseNumber()
		objID := gen.ID()
		idx := NewInMemoryIndex(nil, platformpolicy.NewPlatformCryptographyScheme())

		_, err := idx.OpenRequestsForObjID(ctx, pn, objID, 0)

		require.Error(t, err, ErrLifelineNotFound)
	})

	t.Run("works fine", func(t *testing.T) {
		ctx := inslogger.TestContext(t)
		pn := gen.PulseNumber()
		objID := gen.ID()

		rms := NewRecordMemory()

		objRef := insolar.NewReference(*insolar.NewID(123, nil))
		req := record.Request{Object: objRef}
		reqV := record.Wrap(req)
		reqID := insolar.NewID(333, nil)

		_ = rms.Set(ctx, *reqID, record.Material{Virtual: &reqV})

		objRefS := insolar.NewReference(*insolar.NewID(234, nil))
		reqS := record.Request{Object: objRefS}
		reqSV := record.Wrap(reqS)
		reqSID := insolar.NewID(666, nil)

		_ = rms.Set(ctx, *reqSID, record.Material{Virtual: &reqSV})

		idx := NewInMemoryIndex(rms, platformpolicy.NewPlatformCryptographyScheme())
		idx.createBucket(ctx, pn, objID)

		err := idx.SetRequest(ctx, pn, objID, *reqID)
		require.NoError(t, err)
		err = idx.SetRequest(ctx, pn, objID, *reqSID)
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

func TestInMemoryIndex_SetResult(t *testing.T) {
	t.Run("err when no lifeline", func(t *testing.T) {
		ctx := inslogger.TestContext(t)
		pn := gen.PulseNumber()
		objID := gen.ID()
		idx := NewInMemoryIndex(nil, platformpolicy.NewPlatformCryptographyScheme())

		err := idx.SetResult(ctx, pn, objID, insolar.ID{}, record.Result{})

		require.Error(t, err, ErrLifelineNotFound)
	})

	t.Run("set result, when no requests", func(t *testing.T) {
		ctx := inslogger.TestContext(t)
		pn := gen.PulseNumber()
		objID := gen.ID()

		objRef := insolar.NewReference(*insolar.NewID(123, nil))
		res := record.Result{Request: *objRef}
		resID := insolar.NewID(999, nil)

		previousPen := gen.ID()

		pf := record.PendingFilament{
			RecordID:       *resID,
			PreviousRecord: &previousPen,
		}
		pfv := record.Wrap(pf)
		hash := record.HashVirtual(platformpolicy.NewPlatformCryptographyScheme().ReferenceHasher(), pfv)
		metaID := *insolar.NewID(pn, hash)

		rsm := NewRecordStorageMock(t)
		rsm.SetFunc = func(p context.Context, p1 insolar.ID, p2 record.Material) (r error) {
			require.Equal(t, metaID, p1)
			concrete, ok := record.Unwrap(p2.Virtual).(*record.PendingFilament)
			require.Equal(t, true, ok)
			require.Equal(t, *concrete, pf)
			return nil
		}

		idx := NewInMemoryIndex(rsm, platformpolicy.NewPlatformCryptographyScheme())
		idx.createBucket(ctx, pn, objID)
		buck := idx.buckets[pn][objID]
		buck.objectMeta.Lifeline.PendingPointer = &previousPen

		err := idx.SetResult(ctx, pn, objID, *resID, res)

		require.NoError(t, err)

		require.Equal(t, 1, len(buck.pendingMeta.fullFilament))
		require.Equal(t, metaID, buck.pendingMeta.fullFilament[0].MetaRecordsIDs[0])
		require.Equal(t, pn, buck.pendingMeta.fullFilament[0].PN)
		require.Equal(t, 1, len(buck.objectMeta.PendingRecords))
		require.Equal(t, metaID, buck.objectMeta.PendingRecords[0])
	})

	t.Run("set 2 results, when no requests", func(t *testing.T) {
		ctx := inslogger.TestContext(t)
		pn := gen.PulseNumber()
		objID := gen.ID()
		idx := NewInMemoryIndex(NewRecordMemory(), platformpolicy.NewPlatformCryptographyScheme())
		idx.createBucket(ctx, pn, objID)

		objRef := insolar.NewReference(*insolar.NewID(123, nil))
		res := record.Result{Request: *objRef}
		resID := insolar.NewID(333, nil)
		resS := record.Result{Request: *objRef, Payload: []byte{1, 2, 3, 4, 5, 6}}
		resSID := insolar.NewID(222, nil)

		pf := record.PendingFilament{
			RecordID: *resID,
		}
		pfv := record.Wrap(pf)
		hash := record.HashVirtual(platformpolicy.NewPlatformCryptographyScheme().ReferenceHasher(), pfv)
		firstMetaID := *insolar.NewID(pn, hash)

		pfs := record.PendingFilament{
			RecordID:       *resSID,
			PreviousRecord: &firstMetaID,
		}
		pfsv := record.Wrap(pfs)
		hash = record.HashVirtual(platformpolicy.NewPlatformCryptographyScheme().ReferenceHasher(), pfsv)
		secondMetaID := *insolar.NewID(pn, hash)

		err := idx.SetResult(ctx, pn, objID, *resID, res)
		require.NoError(t, err)
		err = idx.SetResult(ctx, pn, objID, *resSID, resS)
		require.NoError(t, err)

		buck := idx.buckets[pn][objID]

		require.Equal(t, 1, len(buck.pendingMeta.fullFilament))

		require.Equal(t, 2, len(buck.pendingMeta.fullFilament[0].MetaRecordsIDs))
		require.Equal(t, pn, buck.pendingMeta.fullFilament[0].PN)
		require.Equal(t, firstMetaID, buck.pendingMeta.fullFilament[0].MetaRecordsIDs[0])
		require.Equal(t, secondMetaID, buck.pendingMeta.fullFilament[0].MetaRecordsIDs[1])

		require.Equal(t, 2, len(buck.objectMeta.PendingRecords))
		require.Equal(t, firstMetaID, buck.objectMeta.PendingRecords[0])
		require.Equal(t, secondMetaID, buck.objectMeta.PendingRecords[1])
	})

	t.Run("close requests work fine", func(t *testing.T) {
		ctx := inslogger.TestContext(t)
		pn := gen.PulseNumber()
		objID := gen.ID()

		objRef := insolar.NewReference(*insolar.NewID(pn, []byte{1}))
		req := record.Request{Object: objRef}
		reqV := record.Wrap(req)

		objRefS := insolar.NewReference(*insolar.NewID(pn, []byte{2}))
		reqS := record.Request{Object: objRefS}
		reqSV := record.Wrap(reqS)

		res := record.Result{Request: *objRef}
		resV := record.Wrap(res)
		resID := insolar.NewID(3, nil)

		rms := NewRecordMemory()
		_ = rms.Set(ctx, *objRef.Record(), record.Material{Virtual: &reqV})
		_ = rms.Set(ctx, *objRefS.Record(), record.Material{Virtual: &reqSV})
		_ = rms.Set(ctx, *resID, record.Material{Virtual: &resV})

		idx := NewInMemoryIndex(rms, platformpolicy.NewPlatformCryptographyScheme())
		idx.createBucket(ctx, pn, objID)
		_ = idx.SetRequest(ctx, pn, objID, *objRef.Record())
		_ = idx.SetRequest(ctx, pn, objID, *objRefS.Record())

		err := idx.SetResult(ctx, pn, objID, *resID, res)
		require.NoError(t, err)

		open, err := idx.OpenRequestsForObjID(ctx, pn, objID, 10)
		require.NoError(t, err)

		require.Equal(t, 1, len(open))
		require.Equal(t, reqS, open[0])

		buck := idx.buckets[pn][objID]
		require.Equal(t, 1, len(buck.pendingMeta.fullFilament))
		require.Equal(t, pn, buck.pendingMeta.fullFilament[0].PN)
		require.Equal(t, 3, len(buck.pendingMeta.fullFilament[0].MetaRecordsIDs))

		require.Equal(t, 1, len(buck.pendingMeta.notClosedRequestsIdsIndex[pn]))
		_, ok := buck.pendingMeta.notClosedRequestsIdsIndex[pn][*reqS.Object.Record()]
		require.Equal(t, true, ok)

		require.Equal(t, 1, len(buck.pendingMeta.notClosedRequestsIds))
		require.Equal(t, *objRefS.Record(), buck.pendingMeta.notClosedRequestsIds[0])
	})

	t.Run("set result, there are other fillaments", func(t *testing.T) {
		ctx := inslogger.TestContext(t)
		pn := gen.PulseNumber()
		objID := gen.ID()
		idx := NewInMemoryIndex(NewRecordMemory(), platformpolicy.NewPlatformCryptographyScheme())
		idx.createBucket(ctx, pn, objID)
		buck := idx.buckets[pn][objID]
		buck.pendingMeta.fullFilament = append(buck.pendingMeta.fullFilament, chainLink{PN: pn + 1, MetaRecordsIDs: []insolar.ID{}})

		objRef := insolar.NewReference(*insolar.NewID(123, nil))
		res := record.Result{Request: *objRef}
		resID := insolar.NewID(1, nil)

		err := idx.SetResult(ctx, pn, objID, *resID, res)

		require.NoError(t, err)

		require.Equal(t, 2, len(buck.pendingMeta.fullFilament))
		require.Equal(t, pn, buck.pendingMeta.fullFilament[0].PN)
		require.Equal(t, pn+1, buck.pendingMeta.fullFilament[1].PN)
	})
}

func TestInMemoryIndex_RefreshState(t *testing.T) {
	t.Run("err when no lifeline", func(t *testing.T) {
		ctx := inslogger.TestContext(t)
		pn := gen.PulseNumber()
		objID := gen.ID()
		idx := NewInMemoryIndex(nil, platformpolicy.NewPlatformCryptographyScheme())

		err := idx.RefreshState(ctx, pn, objID)

		require.Error(t, err, ErrLifelineNotFound)
	})

	t.Run("works fine. req and res", func(t *testing.T) {
		ctx := inslogger.TestContext(t)
		pn := gen.PulseNumber()
		objID := gen.ID()

		reqID := insolar.NewID(pn+1, nil)
		pf := record.PendingFilament{
			RecordID: *reqID,
		}
		pfv := record.Wrap(pf)
		hash := record.HashVirtual(platformpolicy.NewPlatformCryptographyScheme().ReferenceHasher(), pfv)
		metaReqID := *insolar.NewID(pn+1, hash)

		resID := insolar.NewID(222, nil)
		pfr := record.PendingFilament{RecordID: *resID, PreviousRecord: reqID}
		pfrv := record.Wrap(pfr)
		hash = record.HashVirtual(platformpolicy.NewPlatformCryptographyScheme().ReferenceHasher(), pfrv)
		metaResID := *insolar.NewID(pn, hash)

		rsm := NewRecordStorageMock(t)
		rsm.ForIDFunc = func(p context.Context, p1 insolar.ID) (r record.Material, r1 error) {
			switch p1 {
			case metaReqID:
				return record.Material{Virtual: &pfv}, nil
			case metaResID:
				return record.Material{Virtual: &pfrv}, nil
			case *reqID:
				reqV := record.Wrap(record.Request{Object: insolar.NewReference(*reqID)})
				return record.Material{Virtual: &reqV}, nil
			case *resID:
				resV := record.Wrap(record.Result{Request: *insolar.NewReference(*reqID)})
				return record.Material{Virtual: &resV}, nil
			default:
				panic("test is totaly broken")
			}
		}

		idx := NewInMemoryIndex(rsm, platformpolicy.NewPlatformCryptographyScheme())
		idx.createBucket(ctx, pn, objID)
		buck := idx.buckets[pn][objID]

		buck.pendingMeta.notClosedRequestsIdsIndex = map[insolar.PulseNumber]map[insolar.ID]struct{}{
			pn + 1: {},
			pn:     {},
		}

		buck.pendingMeta.fullFilament = append(buck.pendingMeta.fullFilament, chainLink{PN: pn + 1, MetaRecordsIDs: []insolar.ID{metaReqID}})
		buck.pendingMeta.fullFilament = append(buck.pendingMeta.fullFilament, chainLink{PN: pn, MetaRecordsIDs: []insolar.ID{metaResID}})

		err := idx.RefreshState(ctx, pn, objID)
		require.NoError(t, err)

		require.Equal(t, true, buck.pendingMeta.isStateCalculated)
		require.Equal(t, 2, len(buck.pendingMeta.fullFilament))
		require.Equal(t, 0, len(buck.pendingMeta.notClosedRequestsIds))
		require.Equal(t, 0, len(buck.pendingMeta.notClosedRequestsIdsIndex[pn]))
		require.Equal(t, 0, len(buck.pendingMeta.notClosedRequestsIdsIndex[pn+1]))
	})

	t.Run("works fine. open pending", func(t *testing.T) {
		ctx := inslogger.TestContext(t)
		pn := gen.PulseNumber()
		objID := gen.ID()

		reqID := insolar.NewID(111, nil)
		pf := record.PendingFilament{RecordID: *reqID}
		pfID := gen.ID()

		rsm := NewRecordStorageMock(t)
		rsm.ForIDFunc = func(p context.Context, p1 insolar.ID) (r record.Material, r1 error) {
			switch p1 {
			case pfID:
				metaV := record.Wrap(pf)
				return record.Material{Virtual: &metaV}, nil
			case *reqID:
				reqV := record.Wrap(record.Request{Object: insolar.NewReference(*reqID)})
				return record.Material{Virtual: &reqV}, nil
			default:
				panic("test is totally broken")
			}
		}

		idx := NewInMemoryIndex(rsm, platformpolicy.NewPlatformCryptographyScheme())
		idx.createBucket(ctx, pn, objID)
		buck := idx.buckets[pn][objID]
		buck.pendingMeta.notClosedRequestsIdsIndex = map[insolar.PulseNumber]map[insolar.ID]struct{}{
			pn: {},
		}

		buck.pendingMeta.fullFilament = append(buck.pendingMeta.fullFilament, chainLink{PN: pn, MetaRecordsIDs: []insolar.ID{pfID}})

		err := idx.RefreshState(ctx, pn, objID)
		require.NoError(t, err)

		require.Equal(t, true, buck.pendingMeta.isStateCalculated)
		require.Equal(t, 1, len(buck.pendingMeta.fullFilament))
		require.Equal(t, 1, len(buck.pendingMeta.notClosedRequestsIds))
		require.Equal(t, *reqID, buck.pendingMeta.notClosedRequestsIds[0])
		require.Equal(t, 1, len(buck.pendingMeta.notClosedRequestsIdsIndex[pn]))
	})

	t.Run("calculates EarliestOpenRequest properly", func(t *testing.T) {
		ctx := inslogger.TestContext(t)
		pn := gen.PulseNumber()
		objID := gen.ID()

		objRef := insolar.NewReference(*insolar.NewID(pn, nil))
		req := record.Request{Object: objRef}
		reqV := record.Wrap(req)
		pf := record.PendingFilament{
			RecordID: *req.Object.Record(),
		}
		pfv := record.Wrap(pf)
		hash := record.HashVirtual(platformpolicy.NewPlatformCryptographyScheme().ReferenceHasher(), pfv)
		metaReqID := *insolar.NewID(pn, hash)

		res := record.Result{Request: *objRef}
		resV := record.Wrap(res)
		resID := *insolar.NewID(321, nil)
		pfRes := record.PendingFilament{
			RecordID:       resID,
			PreviousRecord: &metaReqID,
		}
		pfResV := record.Wrap(pfRes)
		hash = record.HashVirtual(platformpolicy.NewPlatformCryptographyScheme().ReferenceHasher(), pfResV)
		metaResID := *insolar.NewID(pn, hash)

		objRefS := insolar.NewReference(*insolar.NewID(pn+1, nil))
		reqS := record.Request{Object: objRefS}
		reqSV := record.Wrap(reqS)

		pfReqS := record.PendingFilament{
			RecordID:       *reqS.Object.Record(),
			PreviousRecord: &metaResID,
		}
		pfReqSV := record.Wrap(pfReqS)
		hash = record.HashVirtual(platformpolicy.NewPlatformCryptographyScheme().ReferenceHasher(), pfReqSV)
		metaReqSID := *insolar.NewID(pn+1, hash)

		objRefT := insolar.NewReference(*insolar.NewID(pn+2, nil))
		reqT := record.Request{Object: objRefT}
		pfReqT := record.PendingFilament{
			RecordID:       *reqT.Object.Record(),
			PreviousRecord: &metaReqSID,
		}
		pfReqTV := record.Wrap(pfReqT)
		hash = record.HashVirtual(platformpolicy.NewPlatformCryptographyScheme().ReferenceHasher(), pfReqTV)
		metaReqTID := *insolar.NewID(pn, hash)

		resT := record.Result{Request: *objRefT}
		reqTV := record.Wrap(reqT)
		resTV := record.Wrap(resT)
		resTID := *insolar.NewID(999, nil)
		pfResT := record.PendingFilament{
			RecordID:       resTID,
			PreviousRecord: &metaReqTID,
		}
		pfResTV := record.Wrap(pfResT)
		hash = record.HashVirtual(platformpolicy.NewPlatformCryptographyScheme().ReferenceHasher(), pfResTV)
		metaResTID := *insolar.NewID(pn, hash)

		rsm := NewRecordStorageMock(t)
		rsm.ForIDFunc = func(p context.Context, p1 insolar.ID) (r record.Material, r1 error) {
			switch p1 {
			case metaReqID:
				return record.Material{Virtual: &pfv}, nil
			case metaResID:
				return record.Material{Virtual: &pfResV}, nil
			case metaReqSID:
				return record.Material{Virtual: &pfReqSV}, nil
			case metaReqTID:
				return record.Material{Virtual: &pfReqTV}, nil
			case metaResTID:
				return record.Material{Virtual: &pfResTV}, nil
			case *req.Object.Record():
				return record.Material{Virtual: &reqV}, nil
			case resID:
				return record.Material{Virtual: &resV}, nil
			case *reqS.Object.Record():
				return record.Material{Virtual: &reqSV}, nil
			case *reqT.Object.Record():
				return record.Material{Virtual: &reqTV}, nil
			case resTID:
				return record.Material{Virtual: &resTV}, nil
			default:
				panic("test is totally broken")
			}
		}

		idx := NewInMemoryIndex(rsm, platformpolicy.NewPlatformCryptographyScheme())
		idx.createBucket(ctx, pn, objID)
		buck := idx.buckets[pn][objID]

		buck.pendingMeta.fullFilament = append(buck.pendingMeta.fullFilament, chainLink{PN: pn, MetaRecordsIDs: []insolar.ID{metaReqID, metaResID}})
		buck.pendingMeta.fullFilament = append(buck.pendingMeta.fullFilament, chainLink{PN: pn + 1, MetaRecordsIDs: []insolar.ID{metaReqSID}})

		buck.pendingMeta.fullFilament = append(buck.pendingMeta.fullFilament, chainLink{PN: pn + 2, MetaRecordsIDs: []insolar.ID{metaReqTID, metaResTID}})

		buck.pendingMeta.notClosedRequestsIdsIndex = map[insolar.PulseNumber]map[insolar.ID]struct{}{
			pn:     {},
			pn + 1: {},
			pn + 2: {},
		}

		err := idx.RefreshState(ctx, pn, objID)
		require.NoError(t, err)

		require.Equal(t, pn+1, buck.objectMeta.Lifeline.EarliestOpenRequest)
	})
}
