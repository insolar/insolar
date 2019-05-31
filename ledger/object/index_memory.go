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
	"sort"
	"sync"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"go.opencensus.io/stats"
)

// extendedIndexBucket is a thread-safe wrapper around IndexBucket struct.
// Due to IndexBucket is a protobuf-generated struct,
// extendedIndexBucket was created for creating an opportunity for using of IndexBucket struct  in a thread-safe way.
// Also it stores some meta-info, that is required for the work process
type extendedIndexBucket struct {
	sync.RWMutex
	IndexBucket

	pendingMeta pendingMeta
}

type pendingMeta struct {
	isStateCalculated bool
	fullFilament      []chainLink

	notClosedRequests      []record.Request
	notClosedRequestsIndex map[insolar.PulseNumber]map[insolar.ID]*record.Request
	requestPNIndex         map[insolar.ID]insolar.PulseNumber

	readPendingUntil *insolar.PulseNumber
}

type chainLink struct {
	PN      insolar.PulseNumber
	Records []record.Virtual
}

func (i *extendedIndexBucket) lifeline() (Lifeline, error) {
	i.RLock()
	defer i.RUnlock()

	return CloneIndex(i.Lifeline), nil
}

func (i *extendedIndexBucket) setLifeline(lifeline Lifeline, pn insolar.PulseNumber) {
	i.Lock()
	defer i.Unlock()

	i.Lifeline = lifeline
	i.LifelineLastUsed = pn
}

func (i *extendedIndexBucket) setLifelineLastUsed(pn insolar.PulseNumber) {
	i.Lock()
	defer i.Unlock()

	i.LifelineLastUsed = pn
}

// InMemoryIndex is a in-memory storage, that stores a collection of IndexBuckets
type InMemoryIndex struct {
	bucketsLock sync.RWMutex
	buckets     map[insolar.PulseNumber]map[insolar.ID]*extendedIndexBucket
}

// NewInMemoryIndex creates a new InMemoryIndex
func NewInMemoryIndex() *InMemoryIndex {
	return &InMemoryIndex{
		buckets: map[insolar.PulseNumber]map[insolar.ID]*extendedIndexBucket{},
	}
}

func (i *InMemoryIndex) createBucket(ctx context.Context, pn insolar.PulseNumber, objID insolar.ID) *extendedIndexBucket {
	i.bucketsLock.Lock()
	defer i.bucketsLock.Unlock()

	bucket := &extendedIndexBucket{
		IndexBucket: IndexBucket{
			ObjID:          objID,
			PendingRecords: []record.Virtual{},
		},
		pendingMeta: pendingMeta{
			fullFilament:           []chainLink{},
			notClosedRequests:      []record.Request{},
			notClosedRequestsIndex: map[insolar.PulseNumber]map[insolar.ID]*record.Request{},
			requestPNIndex:         map[insolar.ID]insolar.PulseNumber{},
			readPendingUntil:       nil,
			isStateCalculated:      false,
		},
	}

	objsByPn, ok := i.buckets[pn]
	if !ok {
		objsByPn = map[insolar.ID]*extendedIndexBucket{}
		i.buckets[pn] = objsByPn
	}
	objsByPn[objID] = bucket

	inslogger.FromContext(ctx).Debugf("[createBucket] create bucket for obj - %v was created successfully", objID.DebugString())
	return bucket
}

func (i *InMemoryIndex) bucket(pn insolar.PulseNumber, objID insolar.ID) *extendedIndexBucket {
	i.bucketsLock.RLock()
	defer i.bucketsLock.RUnlock()

	objsByPn, ok := i.buckets[pn]
	if !ok {
		return nil
	}

	return objsByPn[objID]
}

// Set sets a lifeline to a bucket with provided pulseNumber and ID
func (i *InMemoryIndex) Set(ctx context.Context, pn insolar.PulseNumber, objID insolar.ID, lifeline Lifeline) error {
	b := i.bucket(pn, objID)
	if b == nil {
		b = i.createBucket(ctx, pn, objID)
	}
	b.setLifeline(lifeline, pn)

	stats.Record(ctx,
		statBucketAddedCount.M(1),
	)

	inslogger.FromContext(ctx).Debugf("[Set] lifeline for obj - %v was set successfully", objID.DebugString())
	return nil
}

func (i *InMemoryIndex) SetRequest(ctx context.Context, pn insolar.PulseNumber, objID insolar.ID, req record.Request) error {
	b := i.bucket(pn, objID)
	if b == nil {
		return ErrLifelineNotFound
	}

	b.Lock()
	defer b.Unlock()

	// if b.PreviousPendingFilament == 0 {
	// 	b.PreviousPendingFilament = pn
	// }

	b.PendingRecords = append(b.PendingRecords, record.Wrap(req))

	isInserted := false
	for i, chainPart := range b.pendingMeta.fullFilament {
		if chainPart.PN == pn {
			b.pendingMeta.fullFilament[i].Records = append(b.pendingMeta.fullFilament[i].Records, record.Wrap(req))
			isInserted = true
		}
	}

	if !isInserted {
		b.pendingMeta.fullFilament = append(b.pendingMeta.fullFilament, chainLink{Records: []record.Virtual{record.Wrap(req)}, PN: pn})
		sort.Slice(b.pendingMeta.fullFilament, func(i, j int) bool {
			return b.pendingMeta.fullFilament[i].PN < b.pendingMeta.fullFilament[j].PN
		})
	}

	b.pendingMeta.requestPNIndex[*req.Object.Record()] = pn

	_, ok := b.pendingMeta.notClosedRequestsIndex[pn]
	if !ok {
		b.pendingMeta.notClosedRequestsIndex[pn] = map[insolar.ID]*record.Request{}
	}
	b.pendingMeta.notClosedRequestsIndex[pn][*req.Object.Record()] = &req
	b.pendingMeta.notClosedRequests = append(b.pendingMeta.notClosedRequests, req)

	stats.Record(ctx,
		statObjectPendingRequestsInMemoryAddedCount.M(int64(1)),
	)

	return nil

}
func (i *InMemoryIndex) SetResult(ctx context.Context, pn insolar.PulseNumber, objID insolar.ID, res record.Result) error {
	b := i.bucket(pn, objID)
	if b == nil {
		return ErrLifelineNotFound
	}

	b.Lock()
	defer b.Unlock()

	b.PendingRecords = append(b.PendingRecords, record.Wrap(res))

	isInserted := false
	for i, chainPart := range b.pendingMeta.fullFilament {
		if chainPart.PN == pn {
			b.pendingMeta.fullFilament[i].Records = append(b.pendingMeta.fullFilament[i].Records, record.Wrap(res))
			isInserted = true
		}
	}

	if !isInserted {
		b.pendingMeta.fullFilament = append(b.pendingMeta.fullFilament, chainLink{Records: []record.Virtual{record.Wrap(res)}, PN: pn})
		sort.Slice(b.pendingMeta.fullFilament, func(i, j int) bool {
			return b.pendingMeta.fullFilament[i].PN < b.pendingMeta.fullFilament[j].PN
		})
	}

	reqPN, ok := b.pendingMeta.requestPNIndex[*res.Request.Record()]
	if ok {
		delete(b.pendingMeta.notClosedRequestsIndex[reqPN], *res.Request.Record())
		for i := 0; i < len(b.pendingMeta.notClosedRequests); i++ {
			if *b.pendingMeta.notClosedRequests[i].Object.Record() == *res.Request.Record() {
				b.pendingMeta.notClosedRequests = append(b.pendingMeta.notClosedRequests[:i], b.pendingMeta.notClosedRequests[i+1:]...)
				break
			}
		}
	}

	stats.Record(ctx,
		statObjectPendingRequestsInMemoryAddedCount.M(int64(1)),
	)

	return nil
}

func (i *InMemoryIndex) SetFilament(ctx context.Context, pn insolar.PulseNumber, objID insolar.ID, filPN insolar.PulseNumber, recs []record.Virtual) error {
	b := i.bucket(pn, objID)
	if b == nil {
		return ErrLifelineNotFound
	}

	b.Lock()
	defer b.Unlock()

	b.pendingMeta.fullFilament = append(b.pendingMeta.fullFilament, chainLink{Records: recs, PN: filPN})
	sort.Slice(b.pendingMeta.fullFilament, func(i, j int) bool {
		return b.pendingMeta.fullFilament[i].PN < b.pendingMeta.fullFilament[j].PN
	})

	return nil
}

func (i *InMemoryIndex) RefreshState(ctx context.Context, pn insolar.PulseNumber, objID insolar.ID) error {
	b := i.bucket(pn, objID)
	if b == nil {
		return ErrLifelineNotFound
	}

	b.Lock()
	defer b.Unlock()

	for _, chainLink := range b.pendingMeta.fullFilament {
		for _, chainPart := range chainLink.Records {
			tChainPart := chainPart
			switch r := record.Unwrap(&tChainPart).(type) {
			case *record.Request:
				b.pendingMeta.notClosedRequestsIndex[chainLink.PN][*r.Object.Record()] = r
				b.pendingMeta.requestPNIndex[*r.Object.Record()] = chainLink.PN
			case *record.Result:
				reqPN, ok := b.pendingMeta.requestPNIndex[*r.Request.Record()]
				if ok {
					delete(b.pendingMeta.notClosedRequestsIndex[reqPN], *r.Request.Record())
				}
			}
		}
	}

	isEarliestFound := false

	for i, chainLink := range b.pendingMeta.fullFilament {
		if len(b.pendingMeta.notClosedRequestsIndex[chainLink.PN]) == 0 {
			if isEarliestFound {
				continue
			}

			b.pendingMeta.readPendingUntil = &b.pendingMeta.fullFilament[i].PN
			isEarliestFound = true

		} else {
			for _, ncr := range b.pendingMeta.notClosedRequestsIndex[chainLink.PN] {
				b.pendingMeta.notClosedRequests = append(b.pendingMeta.notClosedRequests, *ncr)
			}
		}
	}

	b.pendingMeta.isStateCalculated = true

	return nil
}

func (i *InMemoryIndex) SetReadUntil(ctx context.Context, pn insolar.PulseNumber, objID insolar.ID, readUntil *insolar.PulseNumber) error {
	b := i.bucket(pn, objID)
	if b == nil {
		return ErrLifelineNotFound
	}

	b.Lock()
	defer b.Unlock()

	b.pendingMeta.readPendingUntil = readUntil

	return nil
}

func (i *InMemoryIndex) MetaForObjID(ctx context.Context, currentPN insolar.PulseNumber, objID insolar.ID) (PendingMeta, error) {
	b := i.bucket(currentPN, objID)
	if b == nil {
		return PendingMeta{}, ErrLifelineNotFound
	}

	b.RLock()
	defer b.RUnlock()

	var ppf *insolar.PulseNumber
	if b.PreviousPendingFilament != 0 {
		ppf = &b.PreviousPendingFilament
	}

	return PendingMeta{
		IsStateCalculated: b.pendingMeta.isStateCalculated,
		ReadUntil:         b.pendingMeta.readPendingUntil,
		PreviousPN:        ppf,
	}, nil
}

func (i *InMemoryIndex) OpenRequestsForObjID(ctx context.Context, currentPN insolar.PulseNumber, objID insolar.ID, count int) ([]record.Request, error) {
	b := i.bucket(currentPN, objID)
	if b == nil {
		return nil, ErrLifelineNotFound
	}

	b.RLock()
	defer b.RUnlock()

	if len(b.pendingMeta.notClosedRequests) > count {
		return append([]record.Request{}, b.pendingMeta.notClosedRequests[:count]...), nil
	}

	return append([]record.Request{}, b.pendingMeta.notClosedRequests...), nil
}

func (i *InMemoryIndex) Records(ctx context.Context, currentPN insolar.PulseNumber, objID insolar.ID) ([]record.Virtual, error) {
	b := i.bucket(currentPN, objID)
	if b == nil {
		return nil, ErrLifelineNotFound
	}

	b.RLock()
	defer b.RLock()

	return b.PendingRecords, nil
}

// SetBucket adds a bucket with provided pulseNumber and ID
func (i *InMemoryIndex) SetBucket(ctx context.Context, pn insolar.PulseNumber, bucket IndexBucket) error {
	i.bucketsLock.Lock()
	defer i.bucketsLock.Unlock()

	bucks, ok := i.buckets[pn]
	if !ok {
		bucks = map[insolar.ID]*extendedIndexBucket{}
		i.buckets[pn] = bucks
	}

	bucks[bucket.ObjID] = &extendedIndexBucket{
		IndexBucket: bucket,
		pendingMeta: pendingMeta{
			notClosedRequests:      []record.Request{},
			fullFilament:           []chainLink{},
			isStateCalculated:      false,
			requestPNIndex:         map[insolar.ID]insolar.PulseNumber{},
			notClosedRequestsIndex: map[insolar.PulseNumber]map[insolar.ID]*record.Request{},
		},
	}

	stats.Record(ctx,
		statBucketAddedCount.M(1),
	)

	return nil
}

// ForID returns a lifeline from a bucket with provided PN and ObjID
func (i *InMemoryIndex) ForID(ctx context.Context, pn insolar.PulseNumber, objID insolar.ID) (Lifeline, error) {
	b := i.bucket(pn, objID)
	if b == nil {
		return Lifeline{}, ErrLifelineNotFound
	}
	return b.lifeline()
}

// ForPNAndJet returns a collection of buckets for a provided pn and jetID
func (i *InMemoryIndex) ForPNAndJet(ctx context.Context, pn insolar.PulseNumber, jetID insolar.JetID) []IndexBucket {
	i.bucketsLock.Lock()
	defer i.bucketsLock.Unlock()

	bucks, ok := i.buckets[pn]
	if !ok {
		return nil
	}

	res := []IndexBucket{}

	for _, b := range bucks {
		if b.Lifeline.JetID != jetID {
			continue
		}

		clonedLfl := CloneIndex(b.Lifeline)
		var clonedRecords []record.Virtual

		clonedRecords = append(clonedRecords, b.PendingRecords...)

		res = append(res, IndexBucket{
			ObjID:            b.ObjID,
			Lifeline:         clonedLfl,
			LifelineLastUsed: b.LifelineLastUsed,
			PendingRecords:   clonedRecords,
		})
	}

	return res
}

// SetLifelineUsage updates a last usage fields of a bucket for a provided pulseNumber and an object id
func (i *InMemoryIndex) SetLifelineUsage(ctx context.Context, pn insolar.PulseNumber, objID insolar.ID) error {
	b := i.bucket(pn, objID)
	if b == nil {
		return ErrLifelineNotFound
	}

	b.setLifelineLastUsed(pn)

	return nil
}

// DeleteForPN deletes all buckets for a provided pulse number
func (i *InMemoryIndex) DeleteForPN(ctx context.Context, pn insolar.PulseNumber) {
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
			statObjectPendingRequestsInMemoryRemovedCount.M(int64(len(buck.PendingRecords))),
		)
	}
}
