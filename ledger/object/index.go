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
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/internal/ledger/store"
	"go.opencensus.io/stats"
)

//go:generate minimock -i github.com/insolar/insolar/ledger/object.LifelineIndex -o ./ -s _mock.go

// LifelineIndex is a base storage for lifelines.
type LifelineIndex interface {
	// LifelineAccessor provides methods for fetching lifelines.
	LifelineAccessor
	// LifelineModifier provides methods for modifying lifelines.
	LifelineModifier
}

//go:generate minimock -i github.com/insolar/insolar/ledger/object.LifelineAccessor -o ./ -s _mock.go

// LifelineAccessor provides methods for fetching lifelines.
type LifelineAccessor interface {
	// ForID returns a lifeline from a bucket with provided PN and ObjID
	ForID(ctx context.Context, pn insolar.PulseNumber, objID insolar.ID) (Lifeline, error)
}

//go:generate minimock -i github.com/insolar/insolar/ledger/object.LifelineModifier -o ./ -s _mock.go

// LifelineModifier provides methods for modifying lifelines.
type LifelineModifier interface {
	// Set set a lifeline to a bucket with provided pulseNumber and ID
	Set(ctx context.Context, pn insolar.PulseNumber, objID insolar.ID, lifeline Lifeline) error
}

//go:generate minimock -i github.com/insolar/insolar/ledger/object.IndexPendingModifier -o ./ -s _mock.go

// IndexPendingModifier provides methods for modifying pending requests.
type IndexPendingModifier interface {
	// SetRequest adds a request to a bucket with provided pulseNumber and ID
	SetRequest(ctx context.Context, pn insolar.PulseNumber, objID insolar.ID, reqID insolar.ID) error
	// SetResultRecord adds a result record to a bucket with provided pulseNumber and ID
	SetResultRecord(ctx context.Context, pn insolar.PulseNumber, objID insolar.ID, resID insolar.ID) error
}

//go:generate minimock -i github.com/insolar/insolar/ledger/object.IndexBucketModifier -o ./ -s _mock.go

// IndexBucketModifier provides methods for modifying buckets of index.
// Index contains buckets with pn->objID->Bucket hierarchy.
// With using of IndexBucketModifier there is a possibility to set buckets from outside of an index.
type IndexBucketModifier interface {
	// SetBucket adds a bucket with provided pulseNumber and ID
	SetBucket(ctx context.Context, pn insolar.PulseNumber, bucket IndexBucket) error
}

//go:generate minimock -i github.com/insolar/insolar/ledger/object.LifelineStateModifier -o ./ -s _mock.go

// LifelineStateModifier provides an interface for changing a state of lifeline.
type LifelineStateModifier interface {
	// SetLifelineUsage updates a last usage fields of a bucket for a provided pulseNumber and an object id
	SetLifelineUsage(ctx context.Context, pn insolar.PulseNumber, objID insolar.ID) error
}

//go:generate minimock -i github.com/insolar/insolar/ledger/object.IndexCleaner -o ./ -s _mock.go

// IndexCleaner provides an interface for removing backets from a storage.
type IndexCleaner interface {
	// DeleteForPN method removes indexes from a storage for a provided
	DeleteForPN(ctx context.Context, pn insolar.PulseNumber)
}

//go:generate minimock -i github.com/insolar/insolar/ledger/object.IndexBucketAccessor -o ./ -s _mock.go

// IndexBucketAccessor provides an interface for fetching buckets from an index.
type IndexBucketAccessor interface {
	// ForPNAndJet returns a collection of buckets for a provided pn and jetID
	ForPNAndJet(ctx context.Context, pn insolar.PulseNumber, jetID insolar.JetID) []IndexBucket
}

// lockedIndexBucket is a thread-safe wrapper around IndexBucket struct.
// Due to IndexBucket is a protobuf-generated struct,
// lockedIndexBucket was created for creating an opportunity for using of IndexBucket struct  in a thread-safe way.
type lockedIndexBucket struct {
	sync.RWMutex

	bucket IndexBucket
}

func (i *lockedIndexBucket) lifeline() (Lifeline, error) {
	i.RLock()
	defer i.RUnlock()

	return CloneIndex(i.bucket.Lifeline), nil
}

func (i *lockedIndexBucket) setLifeline(lifeline Lifeline, pn insolar.PulseNumber) {
	i.Lock()
	defer i.Unlock()

	i.bucket.Lifeline = lifeline
	i.bucket.LifelineLastUsed = pn
}

func (i *lockedIndexBucket) setLifelineLastUsed(pn insolar.PulseNumber) {
	i.Lock()
	defer i.Unlock()

	i.bucket.LifelineLastUsed = pn
}

// func (i *lockedIndexBucket) setRequest(reqID insolar.ID) {
// 	i.requestLock.Lock()
// 	defer i.requestLock.Unlock()
//
// 	i.bucket.Requests = append(i.bucket.Requests, reqID)
// }
//
// func (i *lockedIndexBucket) setResult(resID insolar.ID) {
// 	i.resultLock.Lock()
// 	defer i.resultLock.Unlock()
//
// 	i.bucket.Results = append(i.bucket.Results, resID)
// }

// InMemoryIndex is a in-memory storage, that stores a collection of IndexBuckets
type InMemoryIndex struct {
	bucketsLock sync.RWMutex
	buckets     map[insolar.PulseNumber]map[insolar.ID]*lockedIndexBucket
}

// NewInMemoryIndex creates a new InMemoryIndex
func NewInMemoryIndex() *InMemoryIndex {
	return &InMemoryIndex{
		buckets: map[insolar.PulseNumber]map[insolar.ID]*lockedIndexBucket{},
	}
}

func (i *InMemoryIndex) createBucket(ctx context.Context, pn insolar.PulseNumber, objID insolar.ID) *lockedIndexBucket {
	i.bucketsLock.Lock()
	defer i.bucketsLock.Unlock()

	bucket := &lockedIndexBucket{
		bucket: IndexBucket{
			ObjID:    objID,
			Results:  []insolar.ID{},
			Requests: []insolar.ID{},
		},
	}

	objsByPn, ok := i.buckets[pn]
	if !ok {
		objsByPn = map[insolar.ID]*lockedIndexBucket{}
		i.buckets[pn] = objsByPn
	}
	objsByPn[objID] = bucket

	inslogger.FromContext(ctx).Debugf("[createBucket] create bucket for obj - %v was created successfully", objID.DebugString())
	return bucket
}

func (i *InMemoryIndex) bucket(pn insolar.PulseNumber, objID insolar.ID) *lockedIndexBucket {
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
		statIndexInMemoryAddedCount.M(1),
	)

	inslogger.FromContext(ctx).Debugf("[Set] lifeline for obj - %v was set successfully", objID.DebugString())
	return nil
}

// SetRequest adds a request to a bucket with provided pulseNumber and ID
// func (i *InMemoryIndex) SetRequest(ctx context.Context, pn insolar.PulseNumber, objID insolar.ID, reqID insolar.ID) error {
// 	b := i.bucket(pn, objID)
// 	if b == nil {
// 		b = i.createBucket(ctx, pn, objID)
// 	}
// 	b.setRequest(reqID)
//
// 	return nil
// }

// SetResultRecord adds a result record to a bucket with provided pulseNumber and ID
// func (i *InMemoryIndex) SetResultRecord(ctx context.Context, pn insolar.PulseNumber, objID insolar.ID, resID insolar.ID) error {
// 	b := i.bucket(pn, objID)
// 	if b == nil {
// 		b = i.createBucket(ctx, pn, objID)
// 	}
// 	b.setResult(resID)
//
// 	return nil
// }

// SetBucket adds a bucket with provided pulseNumber and ID
func (i *InMemoryIndex) SetBucket(ctx context.Context, pn insolar.PulseNumber, bucket IndexBucket) error {
	i.bucketsLock.Lock()
	defer i.bucketsLock.Unlock()

	bucks, ok := i.buckets[pn]
	if !ok {
		bucks = map[insolar.ID]*lockedIndexBucket{}
		i.buckets[pn] = bucks
	}

	bucks[bucket.ObjID] = &lockedIndexBucket{
		bucket: bucket,
	}

	stats.Record(ctx,
		statIndexInMemoryAddedCount.M(1),
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
		if b.bucket.Lifeline.JetID != jetID {
			continue
		}

		clonedLfl := CloneIndex(b.bucket.Lifeline)
		var clonedResults []insolar.ID
		var clonedRequests []insolar.ID

		clonedRequests = append(clonedRequests, b.bucket.Requests...)
		clonedResults = append(clonedResults, b.bucket.Results...)

		res = append(res, IndexBucket{
			ObjID:            b.bucket.ObjID,
			Lifeline:         clonedLfl,
			LifelineLastUsed: b.bucket.LifelineLastUsed,
			Results:          clonedResults,
			Requests:         clonedRequests,
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

	buck, ok := i.buckets[pn]
	if ok {
		stats.Record(ctx,
			statIndexInMemoryRemovedCount.M(int64(len(buck))),
		)

		delete(i.buckets, pn)
	}
}

// IndexDB is a db-based storage, that stores a collection of IndexBuckets
type IndexDB struct {
	lock sync.RWMutex
	db   store.DB
}

type indexKey struct {
	pn    insolar.PulseNumber
	objID insolar.ID
}

func (k indexKey) Scope() store.Scope {
	return store.ScopeIndex
}

func (k indexKey) ID() []byte {
	return append(k.pn.Bytes(), k.objID.Bytes()...)
}

type lastKnownIndexPNKey struct {
	objID insolar.ID
}

func (k lastKnownIndexPNKey) Scope() store.Scope {
	return store.ScopeLastKnownIndexPN
}

func (k lastKnownIndexPNKey) ID() []byte {
	return k.objID.Bytes()
}

// NewIndexDB creates a new instance of IndexDB
func NewIndexDB(db store.DB) *IndexDB {
	return &IndexDB{db: db}
}

// Set sets a lifeline to a bucket with provided pulseNumber and ID
func (i *IndexDB) Set(ctx context.Context, pn insolar.PulseNumber, objID insolar.ID, lifeline Lifeline) error {
	i.lock.Lock()
	defer i.lock.Unlock()

	if lifeline.Delegates == nil {
		lifeline.Delegates = []LifelineDelegate{}
	}

	buc, err := i.getBucket(pn, objID)
	if err == ErrIndexBucketNotFound {
		buc = &IndexBucket{}
	} else if err != nil {
		return err
	}

	buc.Lifeline = lifeline
	err = i.setBucket(pn, objID, buc)
	if err != nil {
		return err
	}

	err = i.setLastKnownPN(pn, objID)
	if err != nil {
		return err
	}

	stats.Record(ctx,
		statIndexDBAddedCount.M(1),
	)

	inslogger.FromContext(ctx).Debugf("[Set] lifeline for obj - %v was set successfully", objID.DebugString())

	return nil
}

// SetBucket adds a bucket with provided pulseNumber and ID
func (i *IndexDB) SetBucket(ctx context.Context, pn insolar.PulseNumber, bucket IndexBucket) error {
	i.lock.Lock()
	defer i.lock.Unlock()

	err := i.setBucket(pn, bucket.ObjID, &bucket)
	if err != nil {
		return err
	}

	stats.Record(ctx,
		statIndexDBAddedCount.M(1),
	)

	inslogger.FromContext(ctx).Debugf("[SetBucket] bucket for obj - %v was set successfully", bucket.ObjID.DebugString())
	return i.setLastKnownPN(pn, bucket.ObjID)
}

// ForID returns a lifeline from a bucket with provided PN and ObjID
func (i *IndexDB) ForID(ctx context.Context, pn insolar.PulseNumber, objID insolar.ID) (Lifeline, error) {
	var buck *IndexBucket
	buck, err := i.getBucket(pn, objID)
	if err == ErrIndexBucketNotFound {
		lastPN, err := i.getLastKnownPN(objID)
		if err != nil {
			return Lifeline{}, ErrLifelineNotFound
		}

		buck, err = i.getBucket(lastPN, objID)
		if err != nil {
			return Lifeline{}, err
		}
	} else if err != nil {
		return Lifeline{}, err
	}

	return buck.Lifeline, nil
}

func (i *IndexDB) setBucket(pn insolar.PulseNumber, objID insolar.ID, bucket *IndexBucket) error {
	key := indexKey{pn: pn, objID: objID}

	buff, err := bucket.Marshal()
	if err != nil {
		return err
	}

	return i.db.Set(key, buff)
}

func (i *IndexDB) getBucket(pn insolar.PulseNumber, objID insolar.ID) (*IndexBucket, error) {
	buff, err := i.db.Get(indexKey{pn: pn, objID: objID})
	if err == store.ErrNotFound {
		return nil, ErrIndexBucketNotFound

	}
	if err != nil {
		return nil, err
	}
	bucket := IndexBucket{}
	err = bucket.Unmarshal(buff)
	return &bucket, err
}

func (i *IndexDB) setLastKnownPN(pn insolar.PulseNumber, objID insolar.ID) error {
	key := lastKnownIndexPNKey{objID: objID}
	return i.db.Set(key, pn.Bytes())
}

func (i *IndexDB) getLastKnownPN(objID insolar.ID) (insolar.PulseNumber, error) {
	buff, err := i.db.Get(lastKnownIndexPNKey{objID: objID})
	if err != nil {
		return insolar.FirstPulseNumber, err
	}
	return insolar.NewPulseNumber(buff), err
}
