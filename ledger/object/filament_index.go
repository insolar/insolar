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
	"fmt"
	"sort"
	"sync"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/pkg/errors"
	"go.opencensus.io/stats"
)

// FilamentCacheStorage is a in-memory storage, that stores a collection of IndexBuckets
type FilamentCacheStorage struct {
	recordStorage RecordStorage
	pcs           insolar.PlatformCryptographyScheme

	bucketsLock sync.RWMutex
	buckets     map[insolar.PulseNumber]map[insolar.ID]*pendingMeta
}

// pendingMeta contains info for calculating pending requests states
// The structure contains full chain of pendings (they are grouped by pulse).
// Groups are sorted by pulse, from a lowest to a highest
// There are a few maps inside, that help not to search through full fillament every SetRequest/SetResult
type pendingMeta struct {
	stateCalculationBarrier chan struct{}
	isStateCalculated       bool

	fullFilament []chainLink

	notClosedRequestsIds      []insolar.ID
	notClosedRequestsIdsIndex map[insolar.PulseNumber]map[insolar.ID]struct{}

	resultsForOutOfLimitRequests map[insolar.ID]struct{}
}

type chainLink struct {
	PN             insolar.PulseNumber
	MetaRecordsIDs []insolar.ID
}

// NewInMemoryIndex creates a new FilamentCacheStorage
func NewInMemoryIndex(recordStorage RecordStorage, pcs insolar.PlatformCryptographyScheme) *FilamentCacheStorage {
	return &FilamentCacheStorage{
		pcs:           pcs,
		recordStorage: recordStorage,
		buckets:       map[insolar.PulseNumber]map[insolar.ID]*filamentCache{},
	}
}

func (i *FilamentCacheStorage) createBucket(ctx context.Context, pn insolar.PulseNumber, objID insolar.ID) *filamentCache {
	i.bucketsLock.Lock()
	defer i.bucketsLock.Unlock()

	bucket := &filamentCache{
		objectMeta: FilamentIndex{
			ObjID:          objID,
			PendingRecords: []insolar.ID{},
		},
		pendingMeta: pendingMeta{
			fullFilament:              []chainLink{},
			notClosedRequestsIds:      []insolar.ID{},
			notClosedRequestsIdsIndex: map[insolar.PulseNumber]map[insolar.ID]struct{}{},
			stateCalculationBarrier:   make(chan struct{}),
			isStateCalculated:         false,
		},
	}

	objsByPn, ok := i.buckets[pn]
	if !ok {
		objsByPn = map[insolar.ID]*filamentCache{}
		i.buckets[pn] = objsByPn
	}

	_, ok = objsByPn[objID]
	if !ok {
		objsByPn[objID] = bucket
	}

	inslogger.FromContext(ctx).Debugf("[createBucket] create bucket for obj - %v was created successfully", objID.DebugString())
	return bucket
}

func (i *FilamentCacheStorage) bucket(pn insolar.PulseNumber, objID insolar.ID) *filamentCache {
	i.bucketsLock.RLock()
	defer i.bucketsLock.RUnlock()

	objsByPn, ok := i.buckets[pn]
	if !ok {
		return nil
	}

	return objsByPn[objID]
}

// SetBucket adds a bucket with provided pulseNumber and ID
func (i *FilamentCacheStorage) SetBucket(ctx context.Context, pn insolar.PulseNumber, bucket FilamentIndex) error {
	i.bucketsLock.Lock()
	defer i.bucketsLock.Unlock()

	bucks, ok := i.buckets[pn]
	if !ok {
		bucks = map[insolar.ID]*filamentCache{}
		i.buckets[pn] = bucks
	}

	bucks[bucket.ObjID] = &filamentCache{
		objectMeta: bucket,
		pendingMeta: pendingMeta{
			notClosedRequestsIds:      []insolar.ID{},
			fullFilament:              []chainLink{},
			isStateCalculated:         false,
			stateCalculationBarrier:   make(chan struct{}),
			notClosedRequestsIdsIndex: map[insolar.PulseNumber]map[insolar.ID]struct{}{},
		},
	}

	stats.Record(ctx,
		statBucketAddedCount.M(1),
	)

	return nil
}

// ForPNAndJet returns a collection of buckets for a provided pn and jetID
func (i *FilamentCacheStorage) ForPNAndJet(ctx context.Context, pn insolar.PulseNumber, jetID insolar.JetID) []FilamentIndex {
	i.bucketsLock.Lock()
	defer i.bucketsLock.Unlock()

	bucks, ok := i.buckets[pn]
	if !ok {
		return nil
	}

	res := []FilamentIndex{}

	for _, b := range bucks {
		if b.objectMeta.Lifeline.JetID != jetID {
			continue
		}

		clonedLfl := CloneIndex(b.objectMeta.Lifeline)
		var clonedRecords []insolar.ID

		clonedRecords = append(clonedRecords, b.objectMeta.PendingRecords...)

		res = append(res, FilamentIndex{
			ObjID:            b.objectMeta.ObjID,
			Lifeline:         clonedLfl,
			LifelineLastUsed: b.objectMeta.LifelineLastUsed,
			PendingRecords:   clonedRecords,
		})
	}

	return res
}

// DeleteForPN deletes all buckets for a provided pulse number
func (i *FilamentCacheStorage) DeleteForPN(ctx context.Context, pn insolar.PulseNumber) {
	i.bucketsLock.Lock()
	defer i.bucketsLock.Unlock()

	bucks, ok := i.buckets[pn]
	if !ok {
		return
	}

	delete(i.buckets, pn)

	stats.Record(ctx,
		statBucketRemovedCount.M(int64(len(bucks))),
	)

	for _, buck := range bucks {
		stats.Record(ctx,
			statObjectPendingRecordsInMemoryRemovedCount.M(int64(len(buck.objectMeta.PendingRecords))),
		)
	}
}

// SetRequest sets a request for a specific object
func (i *FilamentCacheStorage) SetRequest(ctx context.Context, pn insolar.PulseNumber, objID insolar.ID, jetID insolar.JetID, reqID insolar.ID) error {
	b := i.bucket(pn, objID)
	if b == nil {
		return ErrLifelineNotFound
	}

	b.Lock()
	defer b.Unlock()

	pf := record.PendingFilament{
		RecordID:       reqID,
		PreviousRecord: b.objectMeta.Lifeline.PendingPointer,
	}

	if b.objectMeta.Lifeline.EarliestOpenRequest == nil {
		b.objectMeta.Lifeline.EarliestOpenRequest = &pn
	}

	pfv := record.Wrap(pf)
	hash := record.HashVirtual(i.pcs.ReferenceHasher(), pfv)
	metaID := *insolar.NewID(pn, hash)

	err := i.recordStorage.Set(ctx, metaID, record.Material{Virtual: &pfv, JetID: jetID})
	if err != nil {
		return errors.Wrap(err, "failed to create a meta-record about pending request")
	}

	b.addMetaIDToFilament(pn, metaID)

	_, ok := b.pendingMeta.notClosedRequestsIdsIndex[pn]
	if !ok {
		b.pendingMeta.notClosedRequestsIdsIndex[pn] = map[insolar.ID]struct{}{}
	}
	b.pendingMeta.notClosedRequestsIdsIndex[pn][reqID] = struct{}{}
	b.pendingMeta.notClosedRequestsIds = append(b.pendingMeta.notClosedRequestsIds, reqID)

	stats.Record(ctx,
		statObjectPendingRequestsInMemoryAddedCount.M(int64(1)),
	)

	inslogger.FromContext(ctx).Debugf("open requests - %v for - %v", len(b.pendingMeta.notClosedRequestsIds), objID.DebugString())

	return nil

}

// SetResult sets a result for a specific object. Also, if there is a not closed request for a provided result,
// the request will be closed
func (i *FilamentCacheStorage) SetResult(ctx context.Context, pn insolar.PulseNumber, objID insolar.ID, jetID insolar.JetID, resID insolar.ID, res record.Result) error {
	logger := inslogger.FromContext(ctx)
	b := i.bucket(pn, objID)
	if b == nil {
		logger.Error("% for id - %v", ErrLifelineNotFound, objID.DebugString())
		return ErrLifelineNotFound
	}

	b.Lock()
	defer b.Unlock()

	pf := record.PendingFilament{
		RecordID:       resID,
		PreviousRecord: b.objectMeta.Lifeline.PendingPointer,
	}

	pfv := record.Wrap(pf)
	hash := record.HashVirtual(i.pcs.ReferenceHasher(), pfv)
	metaID := *insolar.NewID(pn, hash)

	err := i.recordStorage.Set(ctx, metaID, record.Material{Virtual: &pfv, JetID: jetID})
	if err != nil {
		panic(errors.Wrapf(err, "obj id - %v", metaID.DebugString()))
		return errors.Wrap(err, "failed to create a meta-record about pending request")
	}

	b.addMetaIDToFilament(pn, metaID)

	reqsIDs, ok := b.pendingMeta.notClosedRequestsIdsIndex[res.Request.Record().Pulse()]
	if ok {
		delete(reqsIDs, *res.Request.Record())
		for i := 0; i < len(b.pendingMeta.notClosedRequestsIds); i++ {
			if b.pendingMeta.notClosedRequestsIds[i] == *res.Request.Record() {
				b.pendingMeta.notClosedRequestsIds = append(b.pendingMeta.notClosedRequestsIds[:i], b.pendingMeta.notClosedRequestsIds[i+1:]...)
				break
			}
		}
	}

	if len(b.pendingMeta.notClosedRequestsIds) == 0 {
		logger.Debugf("no open requests for - %v", objID.DebugString())
		logger.Debugf("RefreshPendingFilament set EarliestOpenRequest - %v, val - %v", objID.DebugString(), b.objectMeta.Lifeline.EarliestOpenRequest)
		b.objectMeta.Lifeline.EarliestOpenRequest = nil
	}

	stats.Record(ctx,
		statObjectPendingResultsInMemoryAddedCount.M(int64(1)),
	)

	return nil
}

// SetFilament adds a slice of records to an object with provided id and pulse. It's assumed, that the method is
// called for setting records from another light, during the process of filling full chaing of pendings
func (i *FilamentCacheStorage) SetFilament(ctx context.Context, pn insolar.PulseNumber, objID insolar.ID, filPN insolar.PulseNumber, recs []record.CompositeFilamentRecord) error {
	b := i.bucket(pn, objID)
	if b == nil {
		return ErrLifelineNotFound
	}

	b.Lock()
	defer b.Unlock()

	recsIds := make([]insolar.ID, len(recs))
	for idx, rec := range recs {
		recsIds[idx] = rec.MetaID

		err := i.recordStorage.Set(ctx, rec.MetaID, rec.Meta)
		if err != nil && err != ErrOverride {
			panic(errors.Wrapf(err, "obj id - %v", rec.MetaID.DebugString()))
			return errors.Wrap(err, "filament update failed")
		}
		err = i.recordStorage.Set(ctx, rec.RecordID, rec.Record)
		if err != nil && err != ErrOverride {
			panic(errors.Wrapf(err, "obj id - %v", rec.MetaID.DebugString()))
			return errors.Wrap(err, "filament update failed")
		}
	}

	b.pendingMeta.fullFilament = append(b.pendingMeta.fullFilament, chainLink{MetaRecordsIDs: recsIds, PN: filPN})
	sort.Slice(b.pendingMeta.fullFilament, func(i, j int) bool {
		return b.pendingMeta.fullFilament[i].PN < b.pendingMeta.fullFilament[j].PN
	})

	return nil
}

func (i *FilamentCacheStorage) WaitForRefresh(ctx context.Context, pn insolar.PulseNumber, objID insolar.ID) (<-chan struct{}, error) {
	b := i.bucket(pn, objID)
	if b == nil {
		return nil, ErrLifelineNotFound
	}

	b.RLock()
	defer b.RUnlock()

	return b.pendingMeta.stateCalculationBarrier, nil
}

// RefreshState recalculates state of the chain, marks requests as closed and opened.
func (i *FilamentCacheStorage) RefreshState(ctx context.Context, pn insolar.PulseNumber, objID insolar.ID) error {
	println("RefreshState")
	logger := inslogger.FromContext(ctx)
	logger.Debugf("RefreshState for objID: %v pn: %v", objID.DebugString(), pn)
	b := i.bucket(pn, objID)
	if b == nil {
		return ErrLifelineNotFound
	}

	b.Lock()
	defer b.Unlock()

	if b.pendingMeta.isStateCalculated {
		return nil
	}

	for _, chainLink := range b.pendingMeta.fullFilament {
		for _, metaID := range chainLink.MetaRecordsIDs {
			metaRec, err := i.recordStorage.ForID(ctx, metaID)
			if err != nil {
				panic(errors.Wrapf(err, "obj id - %v", metaID.DebugString()))
				return errors.Wrap(err, "failed to refresh an index state")
			}

			concreteMeta := record.Unwrap(metaRec.Virtual).(*record.PendingFilament)
			rec, err := i.recordStorage.ForID(ctx, concreteMeta.RecordID)
			if err != nil {
				panic(errors.Wrapf(err, "obj id - %v", concreteMeta.RecordID.DebugString()))
				return errors.Wrap(err, "failed to refresh an index state")
			}

			switch r := record.Unwrap(rec.Virtual).(type) {
			case *record.Request:
				_, ok := b.pendingMeta.notClosedRequestsIdsIndex[chainLink.PN]
				if !ok {
					b.pendingMeta.notClosedRequestsIdsIndex[chainLink.PN] = map[insolar.ID]struct{}{}
				}
				b.pendingMeta.notClosedRequestsIdsIndex[chainLink.PN][*r.Object.Record()] = struct{}{}
			case *record.Result:
				println(r.Request.Record().Pulse())
				openReqs, ok := b.pendingMeta.notClosedRequestsIdsIndex[r.Request.Record().Pulse()]
				if ok {
					delete(openReqs, *r.Request.Record())
					//
					// _, ok = openReqs[*r.Request.Record()]
					// if ok {
					// 	delete(openReqs, *r.Request.Record())
					// } else {
					// 	panic("uxu, it's for a oudated req")
					// 	b.pendingMeta.resultsForOutOfLimitRequests[*r.Request.Record()] = struct{}{}
					// }
				}
			default:
				panic(fmt.Sprintf("unknow type - %v", r))
			}
		}
	}

	isEarliestFound := false

	for i, chainLink := range b.pendingMeta.fullFilament {
		if len(b.pendingMeta.notClosedRequestsIdsIndex[chainLink.PN]) != 0 {
			if !isEarliestFound {
				b.objectMeta.Lifeline.EarliestOpenRequest = &b.pendingMeta.fullFilament[i].PN
				logger.Debugf("RefreshPendingFilament set EarliestOpenRequest - %v, val - %v", objID.DebugString(), b.objectMeta.Lifeline.EarliestOpenRequest)
				isEarliestFound = true
			}

			for openReqID := range b.pendingMeta.notClosedRequestsIdsIndex[chainLink.PN] {
				b.pendingMeta.notClosedRequestsIds = append(b.pendingMeta.notClosedRequestsIds, openReqID)
			}
		}
	}

	logger.Debugf("RefreshState. Close channel for objID: %v pn: %v", objID.DebugString(), pn)
	close(b.pendingMeta.stateCalculationBarrier)
	b.pendingMeta.isStateCalculated = true

	if len(b.pendingMeta.notClosedRequestsIds) == 0 {
		b.objectMeta.Lifeline.EarliestOpenRequest = nil
		logger.Debugf("RefreshPendingFilament set EarliestOpenRequest - %v, val - %v", objID.DebugString(), b.objectMeta.Lifeline.EarliestOpenRequest)
		logger.Debugf("no open requests for - %v", objID.DebugString())
	}

	return nil
}

func (i *FilamentCacheStorage) ExpireRequests(ctx context.Context, pn insolar.PulseNumber, objID insolar.ID, jetID insolar.JetID, reqs []insolar.ID) error {
	b := i.bucket(pn, objID)
	if b == nil {
		return ErrLifelineNotFound
	}

	b.Lock()
	defer b.Unlock()

	for _, req := range reqs {
		_, ok := b.pendingMeta.resultsForOutOfLimitRequests[req]
		if ok {
			delete(b.pendingMeta.resultsForOutOfLimitRequests, req)
			continue
		}

		expRes := record.Wrap(
			record.Result{
				Status:  record.Expired,
				Request: *insolar.NewReference(req),
				Object:  objID,
			},
		)
		hash := record.HashVirtual(i.pcs.ReferenceHasher(), expRes)
		resID := insolar.NewID(pn, hash)
		err := i.recordStorage.Set(ctx, *resID, record.Material{Virtual: &expRes, JetID: jetID})
		if err != nil {
			return err
		}

		pf := record.PendingFilament{
			RecordID:       *resID,
			PreviousRecord: b.objectMeta.Lifeline.PendingPointer,
		}
		pf.PreviousRecord = b.objectMeta.Lifeline.PendingPointer

		pfv := record.Wrap(pf)
		hash = record.HashVirtual(i.pcs.ReferenceHasher(), pfv)
		metaID := *insolar.NewID(pn, hash)

		err = i.recordStorage.Set(ctx, metaID, record.Material{Virtual: &pfv, JetID: jetID})
		if err != nil {
			panic(errors.Wrapf(err, "obj id - %v", metaID.DebugString()))
			return errors.Wrap(err, "failed to create a meta-record about pending request")
		}

		b.addMetaIDToFilament(pn, metaID)

		delete(b.pendingMeta.resultsForOutOfLimitRequests, req)

		stats.Record(ctx,
			statObjectPendingRecordsInMemoryRemovedCount.M(int64(1)),
		)
	}

	if len(b.pendingMeta.resultsForOutOfLimitRequests) != 0 {
		panic("we have no idea about request. it's impossible situation")
	}

	return nil
}

func (b *filamentCache) addMetaIDToFilament(pn insolar.PulseNumber, metaID insolar.ID) {
	b.objectMeta.PendingRecords = append(b.objectMeta.PendingRecords, metaID)
	b.objectMeta.Lifeline.PendingPointer = &metaID

	isInserted := false
	for i, chainPart := range b.pendingMeta.fullFilament {
		if chainPart.PN == pn {
			b.pendingMeta.fullFilament[i].MetaRecordsIDs = append(b.pendingMeta.fullFilament[i].MetaRecordsIDs, metaID)
			isInserted = true
		}
	}

	if !isInserted {
		b.pendingMeta.fullFilament = append(b.pendingMeta.fullFilament, chainLink{MetaRecordsIDs: []insolar.ID{metaID}, PN: pn})
		sort.Slice(b.pendingMeta.fullFilament, func(i, j int) bool {
			return b.pendingMeta.fullFilament[i].PN < b.pendingMeta.fullFilament[j].PN
		})
	}
}

func (i *FilamentCacheStorage) FirstPending(ctx context.Context, currentPN insolar.PulseNumber, objID insolar.ID) (*record.PendingFilament, error) {
	b := i.bucket(currentPN, objID)
	if b == nil {
		return nil, ErrLifelineNotFound
	}

	b.RLock()
	defer b.RUnlock()

	if len(b.pendingMeta.fullFilament) == 0 {
		return nil, nil
	}
	if len(b.pendingMeta.fullFilament[0].MetaRecordsIDs) == 0 {
		return nil, nil
	}

	metaID := b.pendingMeta.fullFilament[0].MetaRecordsIDs[0]
	rec, err := i.recordStorage.ForID(ctx, metaID)
	if err != nil {
		return nil, err
	}

	return record.Unwrap(rec.Virtual).(*record.PendingFilament), nil

}

// OpenRequestsForObjID returns open requests for a specific object
func (i *FilamentCacheStorage) OpenRequestsForObjID(ctx context.Context, currentPN insolar.PulseNumber, objID insolar.ID, count int) ([]record.Request, error) {
	b := i.bucket(currentPN, objID)
	if b == nil {
		return nil, ErrLifelineNotFound
	}

	b.RLock()
	defer b.RUnlock()

	if len(b.pendingMeta.notClosedRequestsIds) < count {
		count = len(b.pendingMeta.notClosedRequestsIds)
	}

	res := make([]record.Request, count)

	for idx := 0; idx < count; idx++ {
		rec, err := i.recordStorage.ForID(ctx, b.pendingMeta.notClosedRequestsIds[idx])
		if err != nil {
			return nil, err
		}

		switch r := record.Unwrap(rec.Virtual).(type) {
		case *record.Request:
			res[idx] = *r
		default:
			panic("filament is totally broken")
		}
	}

	return res, nil
}

// Records returns all the records for a provided object
func (i *FilamentCacheStorage) Records(ctx context.Context, currentPN insolar.PulseNumber, objID insolar.ID) ([]record.CompositeFilamentRecord, error) {
	b := i.bucket(currentPN, objID)
	if b == nil {
		return nil, ErrLifelineNotFound
	}

	b.RLock()
	defer b.RLock()

	res := make([]record.CompositeFilamentRecord, len(b.objectMeta.PendingRecords))
	for idx, id := range b.objectMeta.PendingRecords {
		metaRec, err := i.recordStorage.ForID(ctx, id)
		if err != nil {
			return nil, err
		}

		concreteMeta := record.Unwrap(metaRec.Virtual).(*record.PendingFilament)
		rec, err := i.recordStorage.ForID(ctx, concreteMeta.RecordID)
		if err != nil {
			return nil, err
		}

		res[idx] = record.CompositeFilamentRecord{
			Record:   rec,
			RecordID: concreteMeta.RecordID,
			Meta:     metaRec,
			MetaID:   id,
		}
	}

	return res, nil
}
