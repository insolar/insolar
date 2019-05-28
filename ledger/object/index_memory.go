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

	requestCache map[insolar.ID]*record.Result

	hasFullChain     bool
	pulsePendingMap  map[insolar.PulseNumber]int
	fullPendingChain []record.Virtual
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
		requestCache:     map[insolar.ID]*record.Result{},
		fullPendingChain: []record.Virtual{},
		hasFullChain:     false,
		pulsePendingMap:  map[insolar.PulseNumber]int{},
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

	b.PendingRecords = append(b.PendingRecords, record.Wrap(req))

	b.requestCache[*req.Object.Record()] = nil

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

	cachedRes, ok := b.requestCache[*res.Request.Record()]
	if !ok {
		return ErrPendingRequestNotFound
	}
	if cachedRes != nil {
		return ErrPendingResultAlreadySet
	}

	b.PendingRecords = append(b.PendingRecords, record.Wrap(res))
	b.requestCache[*res.Request.Record()] = &res

	stats.Record(ctx,
		statObjectPendingRequestsInMemoryAddedCount.M(int64(1)),
	)

	return nil
}

func (i *InMemoryIndex) SetFilament(ctx context.Context, pn insolar.PulseNumber, objID insolar.ID, recs []record.Virtual) error {
	b := i.bucket(pn, objID)
	if b == nil {
		return ErrLifelineNotFound
	}

	b.Lock()
	defer b.Unlock()

	newChain := make([]record.Virtual, len(b.PendingRecords)+len(recs))
	for _, rec := range recs {
		switch r := record.Unwrap(&rec).(type) {
		case *record.Request:
			b.requestCache[*r.Object.Record()] = nil
			newChain = append(newChain, rec)
		case *record.Result:
			res, ok := b.requestCache[*r.Request.Record()]
			if !ok || res != nil {
				panic("inconsistent chain state")
			}
			b.requestCache[*r.Request.Record()] = r
			newChain = append(newChain, rec)
		default:
			panic("unknown record")
		}
	}
	newChain = append(newChain, b.PendingRecords...)
	b.PendingRecords = newChain

	return nil
}

func (i *InMemoryIndex) Meta(ctx context.Context, pn insolar.PulseNumber, objID insolar.ID) (bool, *insolar.PulseNumber, error) {
	b := i.bucket(pn, objID)
	if b == nil {
		return false, nil, ErrLifelineNotFound
	}

	b.Lock()
	defer b.Unlock()

	for _, rec := range b.PendingRecords {
		concrete := record.Unwrap(&rec)
		switch r := concrete.(type) {
		case *record.Request:
			cachedRes, ok := b.requestCache[*r.Object.Record()]
			if !ok || cachedRes == nil {
				return true, &pn, nil
			}
		case *record.Result:
			panic("cases for closing result for previous slots")
		}
	}

	if b.HasOpenRequestsBehind {
		return true, b.LastKnownPendingPN, nil
	}

	// special case for the first pulse ever for an object
	// if there are some requests, we know them
	var lastKnowPN *insolar.PulseNumber = nil
	if len(b.PendingRecords) > 0 {
		lastKnowPN = &pn
	}

	return false, lastKnowPN, nil
}

func (i *InMemoryIndex) HasPendingBehind(ctx context.Context, currentPN insolar.PulseNumber, objID insolar.ID) (bool, error) {
	b := i.bucket(currentPN, objID)
	if b == nil {
		return false, ErrLifelineNotFound
	}

	b.RLock()
	defer b.RUnlock()

	return b.HasOpenRequestsBehind, nil
}

func (i *InMemoryIndex) LastKnownPN(ctx context.Context, currentPN insolar.PulseNumber, objID insolar.ID) (*insolar.PulseNumber, error) {
	b := i.bucket(currentPN, objID)
	if b == nil {
		return nil, ErrLifelineNotFound
	}

	b.RLock()
	defer b.RUnlock()

	return b.LastKnownPendingPN, nil
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
		IndexBucket:  bucket,
		requestCache: map[insolar.ID]*record.Result{},
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
