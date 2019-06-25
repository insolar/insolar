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

// func TestInMemoryIndex_SetRequest(t *testing.T) {
// 	t.Run("err when no lifeline", func(t *testing.T) {
// 		ctx := inslogger.TestContext(t)
// 		pn := gen.PulseNumber()
// 		objID := gen.ID()
//
// 		idx := NewFilamentCacheStorage(NewIndexStorageMemory(), nil, nil, nil, nil, nil, nil, nil, nil)
//
// 		err := idx.SetRequest(ctx, pn, objID, insolar.JetID{}, insolar.ID{})
//
// 		require.Error(t, err, ErrLifelineNotFound)
// 	})
//
// 	t.Run("set first request on the object", func(t *testing.T) {
// 		ctx := inslogger.TestContext(t)
// 		pn := insolar.PulseNumber(10)
// 		objID := *insolar.NewID(pn, nil)
//
// 		reqID := insolar.NewID(100, []byte{1})
// 		prevPending := *insolar.NewID(pn, []byte{1})
//
// 		pf := record.PendingFilament{
// 			RecordID:       *reqID,
// 			PreviousRecord: &prevPending,
// 		}
//
// 		pfv := record.Wrap(pf)
// 		hash := record.HashVirtual(platformpolicy.NewPlatformCryptographyScheme().ReferenceHasher(), pfv)
// 		metaReqID := *insolar.NewID(pn, hash)
//
// 		rsm := NewRecordStorageMock(t)
// 		rsm.SetFunc = func(p context.Context, p1 insolar.ID, p2 record.Material) (r error) {
// 			require.Equal(t, metaReqID, p1)
// 			cast, ok := record.Unwrap(p2.Virtual).(*record.PendingFilament)
// 			require.Equal(t, true, ok)
// 			require.Equal(t, *reqID, cast.RecordID)
// 			require.Equal(t, prevPending, *cast.PreviousRecord)
// 			return nil
// 		}
//
// 		idxStor := NewIndexStorageMemory()
// 		filamnetCache := NewFilamentCacheStorage(idxStor, idxStor, NewIndexLocker(), rsm, nil, platformpolicy.NewPlatformCryptographyScheme(), nil, nil, nil)
// 		idxStor.CreateIndex(ctx, pn, objID)
// 		idxStor.buckets[pn][objID].Lifeline.PendingPointer = &prevPending
//
// 		err := filamnetCache.SetRequest(ctx, pn, objID, insolar.JetID{}, *reqID)
// 		require.NoError(t, err)
//
// 		buck := idxStor.buckets[pn][objID]
// 		filBuck := filamnetCache.buckets[pn][objID]
//
// 		require.NotEqual(t, prevPending, *buck.Lifeline.PendingPointer)
//
// 		require.Equal(t, 1, len(buck.PendingRecords))
// 		require.Equal(t, 1, len(filBuck.fullFilament))
// 		require.Equal(t, 1, len(filBuck.fullFilament[0].MetaRecordsIDs))
//
// 		require.Equal(t, metaReqID, buck.PendingRecords[0])
// 		require.Equal(t, pn, filBuck.fullFilament[0].PN)
// 		require.Equal(t, metaReqID, filBuck.fullFilament[0].MetaRecordsIDs[0])
//
// 		require.Equal(t, 1, len(filBuck.notClosedRequestsIdsIndex))
// 		require.Equal(t, 1, len(filBuck.notClosedRequestsIds))
//
// 		_, ok := filBuck.notClosedRequestsIdsIndex[pn][*reqID]
// 		require.Equal(t, true, ok)
// 		require.Equal(t, *reqID, filBuck.notClosedRequestsIds[0])
// 	})
//
// 	t.Run("set two request on the object", func(t *testing.T) {
// 		ctx := inslogger.TestContext(t)
// 		pn := insolar.PulseNumber(10)
// 		objID := *insolar.NewID(pn, nil)
//
// 		reqID := insolar.NewID(100, []byte{1})
// 		firstPending := *insolar.NewID(pn, []byte{1})
//
// 		pf := record.PendingFilament{
// 			RecordID:       *reqID,
// 			PreviousRecord: &firstPending,
// 		}
//
// 		pfv := record.Wrap(pf)
// 		hash := record.HashVirtual(platformpolicy.NewPlatformCryptographyScheme().ReferenceHasher(), pfv)
// 		firstMeta := *insolar.NewID(pn, hash)
//
// 		reqSID := insolar.NewID(200, []byte{2})
// 		pfs := record.PendingFilament{
// 			RecordID:       *reqSID,
// 			PreviousRecord: &firstMeta,
// 		}
// 		pfsv := record.Wrap(pfs)
// 		hash = record.HashVirtual(platformpolicy.NewPlatformCryptographyScheme().ReferenceHasher(), pfsv)
// 		secondMeta := *insolar.NewID(pn, hash)
//
// 		rsm := NewRecordStorageMock(t)
// 		rsm.SetFunc = func(p context.Context, p1 insolar.ID, p2 record.Material) (r error) {
// 			switch p1 {
// 			case firstMeta:
// 				concrete, ok := record.Unwrap(p2.Virtual).(*record.PendingFilament)
// 				require.Equal(t, true, ok)
// 				require.Equal(t, *reqID, concrete.RecordID)
// 				require.Equal(t, &firstPending, concrete.PreviousRecord)
// 			case secondMeta:
// 				concrete, ok := record.Unwrap(p2.Virtual).(*record.PendingFilament)
// 				require.Equal(t, true, ok)
// 				require.Equal(t, *reqSID, concrete.RecordID)
// 				require.Equal(t, &firstMeta, concrete.PreviousRecord)
// 			default:
// 				t.Fatal("test is totally broken")
// 			}
// 			return nil
// 		}
//
// 		idxStor := NewIndexStorageMemory()
// 		filamnetCache := NewFilamentCacheStorage(idxStor, idxStor, NewIndexLocker(), rsm, nil, platformpolicy.NewPlatformCryptographyScheme(), nil, nil, nil)
// 		idxStor.CreateIndex(ctx, pn, objID)
// 		buck := idxStor.buckets[pn][objID]
// 		buck.Lifeline.PendingPointer = &firstPending
//
// 		err := filamnetCache.SetRequest(ctx, pn, objID, insolar.JetID{}, *reqID)
// 		require.NoError(t, err)
//
// 		err = filamnetCache.SetRequest(ctx, pn, objID, insolar.JetID{}, *reqSID)
// 		require.NoError(t, err)
//
// 		filBuck := filamnetCache.buckets[pn][objID]
//
// 		buck = idxStor.buckets[pn][objID]
// 		require.Equal(t, 2, len(buck.PendingRecords))
// 		require.Equal(t, 1, len(filBuck.fullFilament))
// 		require.Equal(t, 2, len(filBuck.fullFilament[0].MetaRecordsIDs))
//
// 		require.Equal(t, firstMeta, buck.PendingRecords[0])
// 		require.Equal(t, secondMeta, buck.PendingRecords[1])
// 		require.Equal(t, pn, filBuck.fullFilament[0].PN)
// 		require.Equal(t, firstMeta, filBuck.fullFilament[0].MetaRecordsIDs[0])
// 		require.Equal(t, secondMeta, filBuck.fullFilament[0].MetaRecordsIDs[1])
//
// 		require.Equal(t, 1, len(filBuck.notClosedRequestsIdsIndex))
// 		require.Equal(t, 2, len(filBuck.notClosedRequestsIds))
//
// 		_, ok := filBuck.notClosedRequestsIdsIndex[pn][*reqID]
// 		require.Equal(t, true, ok)
// 		require.Equal(t, *reqID, filBuck.notClosedRequestsIds[0])
//
// 		_, ok = filBuck.notClosedRequestsIdsIndex[pn][*reqSID]
// 		require.Equal(t, true, ok)
// 		require.Equal(t, *reqSID, filBuck.notClosedRequestsIds[1])
// 	})
//
// 	t.Run("failed with older pulse", func(t *testing.T) {
// 		ctx := inslogger.TestContext(t)
// 		pn := insolar.PulseNumber(111)
// 		objID := *insolar.NewID(pn, nil)
//
// 		reqID := insolar.NewID(55, []byte{1})
//
// 		rsm := NewRecordStorageMock(t)
// 		rsm.SetMock.Return(nil)
//
// 		idxStor := NewIndexStorageMemory()
// 		filamnetCache := NewFilamentCacheStorage(idxStor, idxStor, NewIndexLocker(), rsm, nil, platformpolicy.NewPlatformCryptographyScheme(), nil, nil, nil)
// 		idxStor.CreateIndex(ctx, pn, objID)
// 		idxStor.buckets[pn][objID].Lifeline.PendingPointer = &objID
//
// 		err := filamnetCache.SetRequest(ctx, pn, objID, insolar.JetID{}, *reqID)
// 		require.Error(t, err)
// 	})
// }
//
// func TestInMemoryIndex_SetFilament(t *testing.T) {
// 	t.Run("works fine", func(t *testing.T) {
// 		ctx := inslogger.TestContext(t)
// 		pn := insolar.PulseNumber(123)
// 		objID := gen.ID()
//
// 		rsm := NewRecordStorageMock(t)
//
// 		reqID := *insolar.NewID(222, nil)
// 		reqRec := record.Request{Object: insolar.NewReference(reqID)}
// 		reqRecV := record.Wrap(reqRec)
// 		pf := record.PendingFilament{
// 			RecordID: reqID,
// 		}
// 		pfv := record.Wrap(pf)
// 		hash := record.HashVirtual(platformpolicy.NewPlatformCryptographyScheme().ReferenceHasher(), pfv)
// 		firstMeta := *insolar.NewID(pn, hash)
//
// 		rsm.SetFunc = func(p context.Context, p1 insolar.ID, p2 record.Material) (r error) {
// 			switch p1 {
// 			case firstMeta:
// 				concrete, ok := record.Unwrap(p2.Virtual).(*record.PendingFilament)
// 				require.Equal(t, true, ok)
// 				require.Equal(t, pf, *concrete)
// 			case reqID:
// 				concrete, ok := record.Unwrap(p2.Virtual).(*record.Request)
// 				require.Equal(t, true, ok)
// 				require.Equal(t, reqRec, *concrete)
// 			default:
// 				t.Fatal("test is totally broken")
// 			}
//
// 			return nil
// 		}
//
// 		idxStor := NewIndexStorageMemory()
// 		filamnetCache := NewFilamentCacheStorage(idxStor, idxStor, NewIndexLocker(), rsm, nil, platformpolicy.NewPlatformCryptographyScheme(), nil, nil, nil)
// 		idxStor.CreateIndex(ctx, pn, objID)
// 		filamnetCache.createPendingBucket(ctx, pn, objID)
//
// 		fBuck := filamnetCache.buckets[pn][objID]
// 		fBuck.fullFilament = append(fBuck.fullFilament, chainLink{PN: pn + 1, MetaRecordsIDs: []insolar.ID{}})
// 		fBuck.fullFilament = append(fBuck.fullFilament, chainLink{PN: pn - 10, MetaRecordsIDs: []insolar.ID{}})
//
// 		fill := []record.CompositeFilamentRecord{{MetaID: firstMeta, Meta: record.Material{Virtual: &pfv}, RecordID: reqID, Record: record.Material{Virtual: &reqRecV}}}
// 		err := filamnetCache.setFilament(ctx, fBuck, pn, fill)
// 		require.NoError(t, err)
//
// 		require.Equal(t, 3, len(fBuck.fullFilament))
// 		require.Equal(t, pn-10, fBuck.fullFilament[0].PN)
// 		require.Equal(t, pn, fBuck.fullFilament[1].PN)
// 		require.Equal(t, pn+1, fBuck.fullFilament[2].PN)
//
// 		require.Equal(t, firstMeta, fBuck.fullFilament[1].MetaRecordsIDs[0])
// 		rsm.MinimockFinish()
// 	})
// }
//
// func TestInMemoryIndex_Records(t *testing.T) {
// 	t.Run("err when no lifeline", func(t *testing.T) {
// 		ctx := inslogger.TestContext(t)
// 		pn := gen.PulseNumber()
// 		objID := gen.ID()
// 		idxStor := NewIndexStorageMemory()
// 		filCache := NewFilamentCacheStorage(idxStor, idxStor, nil, nil, nil, nil, nil, nil, nil)
//
// 		_, err := filCache.Records(ctx, pn, objID)
//
// 		require.Error(t, err, ErrLifelineNotFound)
// 	})
//
// 	t.Run("works fine", func(t *testing.T) {
// 		ctx := inslogger.TestContext(t)
// 		pn := gen.PulseNumber()
// 		objID := gen.ID()
//
// 		objRef := insolar.NewReference(*insolar.NewID(123, nil))
// 		req := record.Request{Object: objRef}
// 		reqV := record.Wrap(req)
// 		reqID := insolar.NewID(444, nil)
// 		metaReq := record.PendingFilament{RecordID: *reqID}
// 		metaReqV := record.Wrap(metaReq)
//
// 		rsm := NewRecordStorageMock(t)
// 		var savedReqID insolar.ID
// 		rsm.SetFunc = func(p context.Context, p1 insolar.ID, p2 record.Material) (r error) {
// 			savedReqID = p1
// 			return nil
// 		}
// 		rsm.ForIDFunc = func(p context.Context, p1 insolar.ID) (r record.Material, r1 error) {
// 			switch p1 {
// 			case savedReqID:
// 				return record.Material{
// 					Virtual: &metaReqV,
// 				}, nil
// 			case *reqID:
// 				return record.Material{
// 					Virtual: &reqV,
// 				}, nil
// 			default:
// 				panic("everything is broken")
// 			}
// 		}
//
// 		idxStor := NewIndexStorageMemory()
// 		filCache := NewFilamentCacheStorage(idxStor, idxStor, NewIndexLocker(), rsm, nil, platformpolicy.NewPlatformCryptographyScheme(), nil, nil, nil)
// 		idxStor.CreateIndex(ctx, pn, objID)
//
// 		_ = filCache.SetRequest(ctx, pn, objID, insolar.JetID{}, *reqID)
//
// 		data, err := filCache.Records(ctx, pn, objID)
//
// 		require.NoError(t, err)
// 		require.Equal(t, 1, len(data))
// 		require.Equal(t, record.Wrap(req), *data[0].Record.Virtual)
// 		require.Equal(t, *reqID, data[0].RecordID)
// 		require.Equal(t, record.Material{Virtual: &metaReqV}, data[0].Meta)
// 		require.Equal(t, savedReqID, data[0].MetaID)
// 	})
// }
//
// func TestInMemoryIndex_OpenRequestsForObjID(t *testing.T) {
// 	t.Run("err when no lifeline", func(t *testing.T) {
// 		ctx := inslogger.TestContext(t)
// 		pn := gen.PulseNumber()
// 		objID := gen.ID()
// 		filCache := NewFilamentCacheStorage(nil, nil, nil, nil, nil, platformpolicy.NewPlatformCryptographyScheme(), nil, nil, nil)
//
// 		_, err := filCache.OpenRequestsForObjID(ctx, pn, objID, 0)
//
// 		require.Error(t, err, ErrLifelineNotFound)
// 	})
//
// 	t.Run("works fine", func(t *testing.T) {
// 		ctx := inslogger.TestContext(t)
// 		pn := insolar.PulseNumber(10)
// 		objID := gen.ID()
//
// 		rms := NewRecordMemory()
//
// 		objRef := insolar.NewReference(*insolar.NewID(123, nil))
// 		req := record.Request{Object: objRef}
// 		reqV := record.Wrap(req)
// 		reqID := insolar.NewID(333, nil)
//
// 		_ = rms.Set(ctx, *reqID, record.Material{Virtual: &reqV})
//
// 		objRefS := insolar.NewReference(*insolar.NewID(234, nil))
// 		reqS := record.Request{Object: objRefS}
// 		reqSV := record.Wrap(reqS)
// 		reqSID := insolar.NewID(666, nil)
//
// 		_ = rms.Set(ctx, *reqSID, record.Material{Virtual: &reqSV})
//
// 		idxStor := NewIndexStorageMemory()
// 		filCache := NewFilamentCacheStorage(idxStor, idxStor, NewIndexLocker(), rms, nil, platformpolicy.NewPlatformCryptographyScheme(), nil, nil, nil)
// 		idxStor.CreateIndex(ctx, pn, objID)
//
// 		err := filCache.SetRequest(ctx, pn, objID, insolar.JetID{}, *reqID)
// 		require.NoError(t, err)
// 		err = filCache.SetRequest(ctx, pn, objID, insolar.JetID{}, *reqSID)
// 		require.NoError(t, err)
//
// 		pBuck := filCache.buckets[pn][objID]
// 		pBuck.isStateCalculated = true
//
// 		t.Run("query all", func(t *testing.T) {
// 			reqs, err := filCache.OpenRequestsForObjID(ctx, pn, objID, 10)
// 			require.NoError(t, err)
// 			require.Equal(t, 2, len(reqs))
// 			require.Equal(t, req, reqs[0])
// 			require.Equal(t, reqS, reqs[1])
// 		})
//
// 		t.Run("query one", func(t *testing.T) {
// 			reqs, err := filCache.OpenRequestsForObjID(ctx, pn, objID, 1)
// 			require.NoError(t, err)
// 			require.Equal(t, 1, len(reqs))
// 			require.Equal(t, req, reqs[0])
// 		})
// 	})
// }
//
// func TestInMemoryIndex_SetResult(t *testing.T) {
// 	t.Run("err when no lifeline", func(t *testing.T) {
// 		ctx := inslogger.TestContext(t)
// 		pn := gen.PulseNumber()
// 		objID := gen.ID()
// 		idxStor := NewIndexStorageMemory()
// 		filCache := NewFilamentCacheStorage(idxStor, idxStor, NewIndexLocker(), nil, nil, platformpolicy.NewPlatformCryptographyScheme(), nil, nil, nil)
//
// 		err := filCache.SetResult(ctx, pn, objID, insolar.JetID{}, insolar.ID{}, record.Result{})
//
// 		require.Error(t, err, ErrLifelineNotFound)
// 	})
//
// 	t.Run("set result, when no requests", func(t *testing.T) {
// 		t.Skip("until https://insolar.atlassian.net/browse/INS-2705")
// 		ctx := inslogger.TestContext(t)
// 		pn := gen.PulseNumber()
// 		objID := gen.ID()
//
// 		objRef := insolar.NewReference(*insolar.NewID(123, nil))
// 		res := record.Result{Request: *objRef}
// 		resID := insolar.NewID(999, nil)
//
// 		rsm := NewRecordStorageMock(t)
// 		idxStor := NewIndexStorageMemory()
// 		filCache := NewFilamentCacheStorage(idxStor, idxStor, NewIndexLocker(), rsm, nil, platformpolicy.NewPlatformCryptographyScheme(), nil, nil, nil)
//
// 		idxStor.CreateIndex(ctx, pn, objID)
//
// 		err := filCache.SetResult(ctx, pn, objID, insolar.JetID{}, *resID, res)
//
// 		require.Error(t, err, ErrResultWithoutRequest)
// 	})
//
// 	t.Run("close requests work fine", func(t *testing.T) {
// 		ctx := inslogger.TestContext(t)
// 		pn := gen.PulseNumber()
// 		objID := gen.ID()
//
// 		objRef := insolar.NewReference(*insolar.NewID(pn, []byte{1}))
// 		req := record.Request{Object: objRef}
// 		reqV := record.Wrap(req)
//
// 		objRefS := insolar.NewReference(*insolar.NewID(pn, []byte{2}))
// 		reqS := record.Request{Object: objRefS}
// 		reqSV := record.Wrap(reqS)
//
// 		res := record.Result{Request: *objRef}
// 		resV := record.Wrap(res)
// 		resID := insolar.NewID(3, nil)
//
// 		rms := NewRecordMemory()
// 		_ = rms.Set(ctx, *objRef.Record(), record.Material{Virtual: &reqV})
// 		_ = rms.Set(ctx, *objRefS.Record(), record.Material{Virtual: &reqSV})
// 		_ = rms.Set(ctx, *resID, record.Material{Virtual: &resV})
//
// 		idxStor := NewIndexStorageMemory()
// 		filCache := NewFilamentCacheStorage(idxStor, idxStor, NewIndexLocker(), rms, nil, platformpolicy.NewPlatformCryptographyScheme(), nil, nil, nil)
// 		idxStor.CreateIndex(ctx, pn, objID)
//
// 		_ = filCache.SetRequest(ctx, pn, objID, insolar.JetID{}, *objRef.Record())
// 		_ = filCache.SetRequest(ctx, pn, objID, insolar.JetID{}, *objRefS.Record())
//
// 		err := filCache.SetResult(ctx, pn, objID, insolar.JetID{}, *resID, res)
// 		require.NoError(t, err)
//
// 		pBuck := filCache.buckets[pn][objID]
// 		pBuck.isStateCalculated = true
//
// 		open, err := filCache.OpenRequestsForObjID(ctx, pn, objID, 10)
// 		require.NoError(t, err)
//
// 		require.Equal(t, 1, len(open))
// 		require.Equal(t, reqS, open[0])
//
// 		fbuck := filCache.buckets[pn][objID]
// 		require.Equal(t, 1, len(fbuck.fullFilament))
// 		require.Equal(t, pn, fbuck.fullFilament[0].PN)
// 		require.Equal(t, 3, len(fbuck.fullFilament[0].MetaRecordsIDs))
//
// 		require.Equal(t, 1, len(fbuck.notClosedRequestsIdsIndex[pn]))
// 		_, ok := fbuck.notClosedRequestsIdsIndex[pn][*reqS.Object.Record()]
// 		require.Equal(t, true, ok)
//
// 		require.Equal(t, 1, len(fbuck.notClosedRequestsIds))
// 		require.Equal(t, *objRefS.Record(), fbuck.notClosedRequestsIds[0])
// 	})
//
// 	t.Run("close requests work fine. change of earliestState", func(t *testing.T) {
// 		ctx := inslogger.TestContext(t)
// 		pn := gen.PulseNumber()
// 		objID := gen.ID()
//
// 		objRef := insolar.NewReference(*insolar.NewID(pn, []byte{1}))
// 		req := record.Request{Object: objRef}
// 		reqV := record.Wrap(req)
//
// 		res := record.Result{Request: *objRef}
// 		resV := record.Wrap(res)
// 		resID := insolar.NewID(3, nil)
//
// 		rms := NewRecordMemory()
// 		_ = rms.Set(ctx, *objRef.Record(), record.Material{Virtual: &reqV})
// 		_ = rms.Set(ctx, *resID, record.Material{Virtual: &resV})
//
// 		idxStor := NewIndexStorageMemory()
// 		filCache := NewFilamentCacheStorage(idxStor, idxStor, NewIndexLocker(), rms, nil, platformpolicy.NewPlatformCryptographyScheme(), nil, nil, nil)
// 		idxStor.CreateIndex(ctx, pn, objID)
// 		fBuck := filCache.createPendingBucket(ctx, pn, objID)
// 		fBuck.isStateCalculated = true
//
// 		buck := idxStor.buckets[pn][objID]
// 		buck.Lifeline.EarliestOpenRequest = &insolar.GenesisPulse.PulseNumber
// 		_ = filCache.SetRequest(ctx, pn, objID, insolar.JetID{}, *objRef.Record())
//
// 		err := filCache.SetResult(ctx, pn, objID, insolar.JetID{}, *resID, res)
// 		require.NoError(t, err)
//
// 		open, err := filCache.OpenRequestsForObjID(ctx, pn, objID, 10)
// 		require.NoError(t, err)
//
// 		require.Equal(t, 0, len(open))
//
// 		buck = idxStor.buckets[pn][objID]
// 		require.Nil(t, buck.Lifeline.EarliestOpenRequest)
// 	})
// }
//
// func TestInMemoryIndex_RefreshState(t *testing.T) {
// 	t.Run("works fine. req and res", func(t *testing.T) {
// 		ctx := inslogger.TestContext(t)
// 		pn := gen.PulseNumber()
// 		objID := gen.ID()
//
// 		reqID := insolar.NewID(pn+1, nil)
// 		pf := record.PendingFilament{
// 			RecordID: *reqID,
// 		}
// 		pfv := record.Wrap(pf)
// 		hash := record.HashVirtual(platformpolicy.NewPlatformCryptographyScheme().ReferenceHasher(), pfv)
// 		metaReqID := *insolar.NewID(pn+1, hash)
//
// 		resID := insolar.NewID(222, nil)
// 		pfr := record.PendingFilament{RecordID: *resID, PreviousRecord: reqID}
// 		pfrv := record.Wrap(pfr)
// 		hash = record.HashVirtual(platformpolicy.NewPlatformCryptographyScheme().ReferenceHasher(), pfrv)
// 		metaResID := *insolar.NewID(pn, hash)
//
// 		rsm := NewRecordStorageMock(t)
// 		rsm.ForIDFunc = func(p context.Context, p1 insolar.ID) (r record.Material, r1 error) {
// 			switch p1 {
// 			case metaReqID:
// 				return record.Material{Virtual: &pfv}, nil
// 			case metaResID:
// 				return record.Material{Virtual: &pfrv}, nil
// 			case *reqID:
// 				reqV := record.Wrap(record.Request{Object: insolar.NewReference(*reqID)})
// 				return record.Material{Virtual: &reqV}, nil
// 			case *resID:
// 				resV := record.Wrap(record.Result{Request: *insolar.NewReference(*reqID)})
// 				return record.Material{Virtual: &resV}, nil
// 			default:
// 				panic("test is totaly broken")
// 			}
// 		}
//
// 		idxStor := NewIndexStorageMemory()
// 		filCache := NewFilamentCacheStorage(idxStor, idxStor, NewIndexLocker(), rsm, nil, platformpolicy.NewPlatformCryptographyScheme(), nil, nil, nil)
// 		filCache.createPendingBucket(ctx, pn, objID)
// 		fBuck := filCache.buckets[pn][objID]
// 		idxStor.CreateIndex(ctx, pn, objID)
// 		buck := idxStor.buckets[pn][objID]
//
// 		fBuck.notClosedRequestsIdsIndex = map[insolar.PulseNumber]map[insolar.ID]struct{}{
// 			pn + 1: {},
// 			pn:     {},
// 		}
//
// 		fBuck.fullFilament = append(fBuck.fullFilament, chainLink{PN: pn + 1, MetaRecordsIDs: []insolar.ID{metaReqID}})
// 		fBuck.fullFilament = append(fBuck.fullFilament, chainLink{PN: pn, MetaRecordsIDs: []insolar.ID{metaResID}})
//
// 		err := filCache.refresh(ctx, buck, fBuck)
// 		require.NoError(t, err)
//
// 		require.Equal(t, true, fBuck.isStateCalculated)
// 		require.Equal(t, 2, len(fBuck.fullFilament))
// 		require.Equal(t, 0, len(fBuck.notClosedRequestsIds))
// 		require.Equal(t, 0, len(fBuck.notClosedRequestsIdsIndex[pn]))
// 		require.Equal(t, 0, len(fBuck.notClosedRequestsIdsIndex[pn+1]))
// 	})
//
// 	t.Run("works fine. open pending", func(t *testing.T) {
// 		ctx := inslogger.TestContext(t)
// 		pn := gen.PulseNumber()
// 		objID := gen.ID()
//
// 		reqID := insolar.NewID(111, nil)
// 		pf := record.PendingFilament{RecordID: *reqID}
// 		pfID := gen.ID()
//
// 		rsm := NewRecordStorageMock(t)
// 		rsm.ForIDFunc = func(p context.Context, p1 insolar.ID) (r record.Material, r1 error) {
// 			switch p1 {
// 			case pfID:
// 				metaV := record.Wrap(pf)
// 				return record.Material{Virtual: &metaV}, nil
// 			case *reqID:
// 				reqV := record.Wrap(record.Request{Object: insolar.NewReference(*reqID)})
// 				return record.Material{Virtual: &reqV}, nil
// 			default:
// 				panic("test is totally broken")
// 			}
// 		}
//
// 		idxStor := NewIndexStorageMemory()
// 		filCache := NewFilamentCacheStorage(idxStor, idxStor, NewIndexLocker(), rsm, nil, platformpolicy.NewPlatformCryptographyScheme(), nil, nil, nil)
// 		filCache.createPendingBucket(ctx, pn, objID)
// 		fBuck := filCache.buckets[pn][objID]
// 		fBuck.notClosedRequestsIdsIndex = map[insolar.PulseNumber]map[insolar.ID]struct{}{
// 			pn: {},
// 		}
// 		idxStor.CreateIndex(ctx, pn, objID)
// 		buck := idxStor.buckets[pn][objID]
//
// 		fBuck.fullFilament = append(fBuck.fullFilament, chainLink{PN: pn, MetaRecordsIDs: []insolar.ID{pfID}})
//
// 		err := filCache.refresh(ctx, buck, fBuck)
// 		require.NoError(t, err)
//
// 		require.Equal(t, true, fBuck.isStateCalculated)
// 		require.Equal(t, 1, len(fBuck.fullFilament))
// 		require.Equal(t, 1, len(fBuck.notClosedRequestsIds))
// 		require.Equal(t, *reqID, fBuck.notClosedRequestsIds[0])
// 		require.Equal(t, 1, len(fBuck.notClosedRequestsIdsIndex[pn]))
// 	})
//
// 	t.Run("calculates EarliestOpenRequest properly", func(t *testing.T) {
// 		ctx := inslogger.TestContext(t)
// 		pn := gen.PulseNumber()
// 		objID := gen.ID()
//
// 		objRef := insolar.NewReference(*insolar.NewID(pn, nil))
// 		req := record.Request{Object: objRef}
// 		reqV := record.Wrap(req)
// 		pf := record.PendingFilament{
// 			RecordID: *req.Object.Record(),
// 		}
// 		pfv := record.Wrap(pf)
// 		hash := record.HashVirtual(platformpolicy.NewPlatformCryptographyScheme().ReferenceHasher(), pfv)
// 		metaReqID := *insolar.NewID(pn, hash)
//
// 		res := record.Result{Request: *objRef}
// 		resV := record.Wrap(res)
// 		resID := *insolar.NewID(321, nil)
// 		pfRes := record.PendingFilament{
// 			RecordID:       resID,
// 			PreviousRecord: &metaReqID,
// 		}
// 		pfResV := record.Wrap(pfRes)
// 		hash = record.HashVirtual(platformpolicy.NewPlatformCryptographyScheme().ReferenceHasher(), pfResV)
// 		metaResID := *insolar.NewID(pn, hash)
//
// 		objRefS := insolar.NewReference(*insolar.NewID(pn+1, nil))
// 		reqS := record.Request{Object: objRefS}
// 		reqSV := record.Wrap(reqS)
//
// 		pfReqS := record.PendingFilament{
// 			RecordID:       *reqS.Object.Record(),
// 			PreviousRecord: &metaResID,
// 		}
// 		pfReqSV := record.Wrap(pfReqS)
// 		hash = record.HashVirtual(platformpolicy.NewPlatformCryptographyScheme().ReferenceHasher(), pfReqSV)
// 		metaReqSID := *insolar.NewID(pn+1, hash)
//
// 		objRefT := insolar.NewReference(*insolar.NewID(pn+2, nil))
// 		reqT := record.Request{Object: objRefT}
// 		pfReqT := record.PendingFilament{
// 			RecordID:       *reqT.Object.Record(),
// 			PreviousRecord: &metaReqSID,
// 		}
// 		pfReqTV := record.Wrap(pfReqT)
// 		hash = record.HashVirtual(platformpolicy.NewPlatformCryptographyScheme().ReferenceHasher(), pfReqTV)
// 		metaReqTID := *insolar.NewID(pn, hash)
//
// 		resT := record.Result{Request: *objRefT}
// 		reqTV := record.Wrap(reqT)
// 		resTV := record.Wrap(resT)
// 		resTID := *insolar.NewID(999, nil)
// 		pfResT := record.PendingFilament{
// 			RecordID:       resTID,
// 			PreviousRecord: &metaReqTID,
// 		}
// 		pfResTV := record.Wrap(pfResT)
// 		hash = record.HashVirtual(platformpolicy.NewPlatformCryptographyScheme().ReferenceHasher(), pfResTV)
// 		metaResTID := *insolar.NewID(pn, hash)
//
// 		rsm := NewRecordStorageMock(t)
// 		rsm.ForIDFunc = func(p context.Context, p1 insolar.ID) (r record.Material, r1 error) {
// 			switch p1 {
// 			case metaReqID:
// 				return record.Material{Virtual: &pfv}, nil
// 			case metaResID:
// 				return record.Material{Virtual: &pfResV}, nil
// 			case metaReqSID:
// 				return record.Material{Virtual: &pfReqSV}, nil
// 			case metaReqTID:
// 				return record.Material{Virtual: &pfReqTV}, nil
// 			case metaResTID:
// 				return record.Material{Virtual: &pfResTV}, nil
// 			case *req.Object.Record():
// 				return record.Material{Virtual: &reqV}, nil
// 			case resID:
// 				return record.Material{Virtual: &resV}, nil
// 			case *reqS.Object.Record():
// 				return record.Material{Virtual: &reqSV}, nil
// 			case *reqT.Object.Record():
// 				return record.Material{Virtual: &reqTV}, nil
// 			case resTID:
// 				return record.Material{Virtual: &resTV}, nil
// 			default:
// 				panic("test is totally broken")
// 			}
// 		}
//
// 		idxStor := NewIndexStorageMemory()
// 		filCache := NewFilamentCacheStorage(idxStor, idxStor, NewIndexLocker(), rsm, nil, platformpolicy.NewPlatformCryptographyScheme(), nil, nil, nil)
// 		filCache.createPendingBucket(ctx, pn, objID)
// 		fBuck := filCache.buckets[pn][objID]
// 		fBuck.notClosedRequestsIdsIndex = map[insolar.PulseNumber]map[insolar.ID]struct{}{
// 			pn: {},
// 		}
// 		idxStor.CreateIndex(ctx, pn, objID)
// 		buck := idxStor.buckets[pn][objID]
//
// 		fBuck.fullFilament = append(fBuck.fullFilament, chainLink{PN: pn, MetaRecordsIDs: []insolar.ID{metaReqID, metaResID}})
// 		fBuck.fullFilament = append(fBuck.fullFilament, chainLink{PN: pn + 1, MetaRecordsIDs: []insolar.ID{metaReqSID}})
//
// 		fBuck.fullFilament = append(fBuck.fullFilament, chainLink{PN: pn + 2, MetaRecordsIDs: []insolar.ID{metaReqTID, metaResTID}})
//
// 		fBuck.notClosedRequestsIdsIndex = map[insolar.PulseNumber]map[insolar.ID]struct{}{
// 			pn:     {},
// 			pn + 1: {},
// 			pn + 2: {},
// 		}
//
// 		err := filCache.refresh(ctx, buck, fBuck)
// 		require.NoError(t, err)
//
// 		require.Equal(t, pn+1, *buck.Lifeline.EarliestOpenRequest)
// 	})
// }
