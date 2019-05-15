package object

import (
	"context"
	"sync"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/internal/ledger/store"
)

//go:generate minimock -i github.com/insolar/insolar/ledger/object.LifelineIndex -o ./ -s _mock.go

type LifelineIndex interface {
	IndexLifelineAccessor
	IndexLifelineModifier
}

//go:generate minimock -i github.com/insolar/insolar/ledger/object.IndexLifelineAccessor -o ./ -s _mock.go

type IndexLifelineAccessor interface {
	LifelineForID(ctx context.Context, pn insolar.PulseNumber, objID insolar.ID) (Lifeline, error)
}

//go:generate minimock -i github.com/insolar/insolar/ledger/object.IndexLifelineModifier -o ./ -s _mock.go

type IndexLifelineModifier interface {
	SetLifeline(ctx context.Context, pn insolar.PulseNumber, objID insolar.ID, lifeline Lifeline) error
}

//go:generate minimock -i github.com/insolar/insolar/ledger/object.IndexPendingModifier -o ./ -s _mock.go

type IndexPendingModifier interface {
	SetRequest(ctx context.Context, pn insolar.PulseNumber, objID insolar.ID, reqID insolar.ID) error
	SetResultRecord(ctx context.Context, pn insolar.PulseNumber, objID insolar.ID, resID insolar.ID) error
}

//go:generate minimock -i github.com/insolar/insolar/ledger/object.IndexBucketModifier -o ./ -s _mock.go

type IndexBucketModifier interface {
	SetBucket(ctx context.Context, pn insolar.PulseNumber, bucket IndexBucket) error
}

//go:generate minimock -i github.com/insolar/insolar/ledger/object.IndexLifelineStateModifier -o ./ -s _mock.go

type IndexLifelineStateModifier interface {
	SetLifelineUsage(ctx context.Context, pn insolar.PulseNumber, objID insolar.ID) error
}

//go:generate minimock -i github.com/insolar/insolar/ledger/object.IndexCleaner -o ./ -s _mock.go

// IndexCleaner provides an interface for removing backets from a storage.
type IndexCleaner interface {
	// DeleteForPN method removes indexes from a storage for a provided
	DeleteForPN(ctx context.Context, pn insolar.PulseNumber)
}

//go:generate minimock -i github.com/insolar/insolar/ledger/object.IndexBucketAccessor -o ./ -s _mock.go

type IndexBucketAccessor interface {
	ForPNAndJet(ctx context.Context, pn insolar.PulseNumber, jetID insolar.JetID) []IndexBucket
}

// type IndexBucketModifier interface {
// 	Clone(ctx context.Context, fromPN insolar.PulseNumber, toPN insolar.PulseNumber, from insolar.JetID, to insolar.JetID)
// }

type LockedIndexBucket struct {
	lifelineLock         sync.RWMutex
	lifelineLastUsedLock sync.RWMutex
	requestLock          sync.RWMutex
	resultLock           sync.RWMutex

	bucket IndexBucket
}

func (i *LockedIndexBucket) lifeline() (Lifeline, error) {
	i.lifelineLock.RLock()
	defer i.lifelineLock.RUnlock()
	if i.bucket.Lifeline == nil {
		return Lifeline{}, ErrLifelineNotFound
	}

	return CloneIndex(*i.bucket.Lifeline), nil
}

func (i *LockedIndexBucket) setLifeline(lifeline *Lifeline, pn insolar.PulseNumber) {
	i.lifelineLock.Lock()
	defer i.lifelineLock.Unlock()

	i.bucket.Lifeline = lifeline
	i.bucket.LifelineLastUsed = pn
}

func (i *LockedIndexBucket) setLifelineLastUsed(pn insolar.PulseNumber) {
	i.lifelineLastUsedLock.Lock()
	defer i.lifelineLastUsedLock.Unlock()

	i.bucket.LifelineLastUsed = pn
}

func (i *LockedIndexBucket) setRequest(reqID insolar.ID) {
	i.requestLock.Lock()
	defer i.requestLock.Unlock()

	i.bucket.Requests = append(i.bucket.Requests, reqID)
}

func (i *LockedIndexBucket) setResult(resID insolar.ID) {
	i.resultLock.Lock()
	defer i.resultLock.Unlock()

	i.bucket.Results = append(i.bucket.Results, resID)
}

type InMemoryIndex struct {
	bucketsLock sync.RWMutex
	buckets     map[insolar.PulseNumber]map[insolar.ID]*LockedIndexBucket
}

func NewInMemoryIndex() *InMemoryIndex {
	return &InMemoryIndex{
		buckets: map[insolar.PulseNumber]map[insolar.ID]*LockedIndexBucket{},
	}
}

func (i *InMemoryIndex) bucket(ctx context.Context, pn insolar.PulseNumber, objID insolar.ID) *LockedIndexBucket {
	i.bucketsLock.Lock()
	defer i.bucketsLock.Unlock()

	var objsByPn map[insolar.ID]*LockedIndexBucket
	objsByPn, ok := i.buckets[pn]
	if !ok {
		objsByPn = map[insolar.ID]*LockedIndexBucket{}
		i.buckets[pn] = objsByPn
	}

	bucket := objsByPn[objID]
	if bucket == nil {
		bucket = &LockedIndexBucket{
			bucket: IndexBucket{
				ObjID:    objID,
				Results:  []insolar.ID{},
				Requests: []insolar.ID{},
			},
		}
		objsByPn[objID] = bucket
	}

	return bucket
}

func (i *InMemoryIndex) SetLifeline(ctx context.Context, pn insolar.PulseNumber, objID insolar.ID, lifeline Lifeline) error {
	b := i.bucket(ctx, pn, objID)
	b.setLifeline(&lifeline, pn)

	return nil
}

func (i *InMemoryIndex) SetRequest(ctx context.Context, pn insolar.PulseNumber, objID insolar.ID, reqID insolar.ID) error {
	b := i.bucket(ctx, pn, objID)
	b.setRequest(reqID)

	return nil
}

func (i *InMemoryIndex) SetResultRecord(ctx context.Context, pn insolar.PulseNumber, objID insolar.ID, resID insolar.ID) error {
	b := i.bucket(ctx, pn, objID)
	b.setResult(resID)

	return nil
}

func (i *InMemoryIndex) SetBucket(ctx context.Context, pn insolar.PulseNumber, bucket IndexBucket) error {
	i.bucketsLock.Lock()
	defer i.bucketsLock.Unlock()

	bucks, ok := i.buckets[pn]
	if !ok {
		bucks = map[insolar.ID]*LockedIndexBucket{}
		i.buckets[pn] = bucks
	}

	bucks[bucket.ObjID] = &LockedIndexBucket{
		bucket: bucket,
	}

	return nil
}

func (i *InMemoryIndex) LifelineForID(ctx context.Context, pn insolar.PulseNumber, objID insolar.ID) (Lifeline, error) {
	b := i.bucket(ctx, pn, objID)
	return b.lifeline()
}

func (i *InMemoryIndex) ForPNAndJet(ctx context.Context, pn insolar.PulseNumber, jetID insolar.JetID) []IndexBucket {
	i.bucketsLock.Lock()
	defer i.bucketsLock.Unlock()

	bucks, ok := i.buckets[pn]
	if !ok {
		return nil
	}

	var res []IndexBucket

	for _, b := range bucks {
		if b.bucket.Lifeline == nil {
			panic("empty lifeline")
		}
		if b.bucket.Lifeline.JetID != jetID {
			continue
		}

		clonedLfl := CloneIndex(*b.bucket.Lifeline)
		var clonedResults []insolar.ID
		var clonedRequests []insolar.ID

		for _, r := range b.bucket.Requests {
			clonedRequests = append(clonedRequests, r)
		}
		for _, r := range b.bucket.Results {
			clonedResults = append(clonedResults, r)
		}

		res = append(res, IndexBucket{
			ObjID:            b.bucket.ObjID,
			Lifeline:         &clonedLfl,
			LifelineLastUsed: b.bucket.LifelineLastUsed,
			Results:          clonedResults,
			Requests:         clonedRequests,
		})

	}

	return res
}

func (i *InMemoryIndex) SetLifelineUsage(ctx context.Context, pn insolar.PulseNumber, objID insolar.ID) error {
	b := i.bucket(ctx, pn, objID)
	_, err := b.lifeline()
	if err != nil {
		return err
	}

	b.setLifelineLastUsed(pn)

	return nil
}

func (i *InMemoryIndex) DeleteForPN(ctx context.Context, pn insolar.PulseNumber) {
	i.bucketsLock.Lock()
	defer i.bucketsLock.Unlock()

	delete(i.buckets, pn)
}

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

func NewIndexDB(db store.DB) *IndexDB {
	return &IndexDB{db: db}
}

func (i *IndexDB) SetLifeline(ctx context.Context, pn insolar.PulseNumber, objID insolar.ID, lifeline Lifeline) error {
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

	buc.Lifeline = &lifeline
	err = i.setBucket(pn, objID, buc)
	if err != nil {
		return err
	}
	return i.setLastKnownPN(ctx, pn, objID)
}

// func (i *IndexDB) SetRequest(ctx context.Context, pn insolar.PulseNumber, objID insolar.ID, reqID insolar.ID) error {
// 	i.lock.Lock()
// 	defer i.lock.Unlock()
//
// 	buc, err := i.get(pn, objID)
// 	if err == store.ErrNotFound {
// 		buc = &IndexBucket{}
// 	} else if err != nil {
// 		return err
// 	}
//
// 	buc.Requests = append(buc.Requests, reqID)
// 	return i.set(pn, objID, buc)
// }
//
// func (i *IndexDB) SetResultRecord(ctx context.Context, pn insolar.PulseNumber, objID insolar.ID, resID insolar.ID) error {
// 	i.lock.Lock()
// 	defer i.lock.Unlock()
//
// 	buc, err := i.get(pn, objID)
// 	if err == store.ErrNotFound {
// 		buc = &IndexBucket{}
// 	} else if err != nil {
// 		return err
// 	}
//
// 	buc.Results = append(buc.Results, resID)
// 	return i.set(pn, objID, buc)
// }

func (i *IndexDB) SetBucket(ctx context.Context, pn insolar.PulseNumber, bucket IndexBucket) error {
	i.lock.Lock()
	defer i.lock.Unlock()

	return i.setBucket(pn, bucket.ObjID, &bucket)
}

func (i *IndexDB) LifelineForID(ctx context.Context, pn insolar.PulseNumber, objID insolar.ID) (Lifeline, error) {
	var buck *IndexBucket
	var err error
	buck, err = i.getBucket(pn, objID)
	if err == ErrIndexBucketNotFound {
		var lastPN insolar.PulseNumber
		lastPN, err = i.getLastKnownPN(objID)
		if err != nil {
			return Lifeline{}, ErrLifelineNotFound
		}

		buck, err = i.getBucket(lastPN, objID)
	}
	if err != nil {
		panic(err)
		return Lifeline{}, err
	}
	if buck.Lifeline == nil {
		return Lifeline{}, ErrLifelineNotFound
	}

	return *buck.Lifeline, nil
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

func (i *IndexDB) setLastKnownPN(ctx context.Context, pn insolar.PulseNumber, objID insolar.ID) error {
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
