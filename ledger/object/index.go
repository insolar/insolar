package object

import (
	"context"
	"sync"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/internal/ledger/store"
)

type Index interface {
	IndexAccessor
	IndexModifier
}

type IndexAccessor interface {
	LifelineForID(ctx context.Context, pn insolar.PulseNumber, objID insolar.ID) (Lifeline, error)
}

type IndexModifier interface {
	SetLifeline(ctx context.Context, pn insolar.PulseNumber, objID insolar.ID, lifeline Lifeline) error
	SetRequest(ctx context.Context, pn insolar.PulseNumber, objID insolar.ID, reqID insolar.ID) error
	SetResultRecord(ctx context.Context, pn insolar.PulseNumber, objID insolar.ID, resID insolar.ID) error
}

type IndexStateModifier interface {
	SetLifelineUsage(ctx context.Context, pn insolar.PulseNumber, objID insolar.ID) error
}

type IndexBucket struct {
	lifelineLock sync.RWMutex
	lifeline     *Lifeline

	requestLock sync.RWMutex
	requests    []insolar.ID

	resultLock sync.RWMutex
	results    []insolar.ID
}

func (i *IndexBucket) getLifeline() (*Lifeline, error) {
	i.lifelineLock.RLock()
	defer i.lifelineLock.RUnlock()
	if i.lifeline == nil {
		return nil, ErrLifelineNotFound
	}

	return i.lifeline, nil
}

func (i *IndexBucket) setLifeline(lifeline *Lifeline) {
	i.lifelineLock.Lock()
	defer i.lifelineLock.Unlock()

	i.lifeline = lifeline
}

func (i *IndexBucket) setRequest(reqID insolar.ID) {
	i.requestLock.Lock()
	defer i.requestLock.Unlock()

	i.requests = append(i.requests, reqID)
}

func (i *IndexBucket) setResult(resID insolar.ID) {
	i.resultLock.Lock()
	defer i.resultLock.Unlock()

	i.results = append(i.results, resID)
}

type InMemoryIndex struct {
	bucketsLock sync.RWMutex
	buckets     map[insolar.PulseNumber]map[insolar.ID]*IndexBucket
}

func NewInMemoryIndex() *InMemoryIndex {
	return &InMemoryIndex{buckets: map[insolar.PulseNumber]map[insolar.ID]*IndexBucket{}}
}

func (i *InMemoryIndex) getBucket(ctx context.Context, pn insolar.PulseNumber, objID insolar.ID) *IndexBucket {
	i.bucketsLock.Lock()
	defer i.bucketsLock.Unlock()

	var objsByPn map[insolar.ID]*IndexBucket
	objsByPn, ok := i.buckets[pn]
	if !ok {
		objsByPn = map[insolar.ID]*IndexBucket{}
		i.buckets[pn] = objsByPn
	}

	bucket := objsByPn[objID]
	if bucket == nil {
		bucket = &IndexBucket{
			requests: []insolar.ID{},
			results:  []insolar.ID{},
		}
		objsByPn[objID] = bucket
	}

	return bucket
}

func (i *InMemoryIndex) SetLifeline(ctx context.Context, pn insolar.PulseNumber, objID insolar.ID, lifeline Lifeline) error {
	b := i.getBucket(ctx, pn, objID)
	b.setLifeline(&lifeline)

	return nil
}

func (i *InMemoryIndex) SetRequest(ctx context.Context, pn insolar.PulseNumber, objID insolar.ID, reqID insolar.ID) error {
	b := i.getBucket(ctx, pn, objID)
	b.setRequest(reqID)

	return nil
}

func (i *InMemoryIndex) SetResultRecord(ctx context.Context, pn insolar.PulseNumber, objID insolar.ID, resID insolar.ID) error {
	b := i.getBucket(ctx, pn, objID)
	b.setResult(resID)

	return nil
}

func (i *InMemoryIndex) LifelineForID(ctx context.Context, pn insolar.PulseNumber, objID insolar.ID) (Lifeline, error) {
	b := i.getBucket(ctx, pn, objID)
	lfl, err := b.getLifeline()
	if err != nil {
		return Lifeline{}, err
	}

	return *lfl, nil
}

func (i *InMemoryIndex) SetLifelineUsage(ctx context.Context, pn insolar.PulseNumber, objID insolar.ID) error {
	b := i.getBucket(ctx, pn, objID)
	lfl, err := b.getLifeline()
	if err != nil {
		return err
	}

	lfl.LastUsed = pn
	b.setLifeline(lfl)

	return nil
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

func NewIndexDB(db store.DB) *IndexDB {
	return &IndexDB{db: db}
}

func (i *IndexDB) SetLifeline(ctx context.Context, pn insolar.PulseNumber, objID insolar.ID, lifeline Lifeline) error {
	i.lock.Lock()
	defer i.lock.Unlock()

	if lifeline.Delegates == nil {
		lifeline.Delegates = map[insolar.Reference]insolar.Reference{}
	}

	buc, err := i.get(pn, objID)
	if err == store.ErrNotFound {
		buc = &IndexBucket{}
	} else if err != nil {
		return err
	}

	buc.lifeline = &lifeline
	return i.set(pn, objID, buc)
}

func (i *IndexDB) SetRequest(ctx context.Context, pn insolar.PulseNumber, objID insolar.ID, reqID insolar.ID) error {
	i.lock.Lock()
	defer i.lock.Unlock()

	buc, err := i.get(pn, objID)
	if err == store.ErrNotFound {
		buc = &IndexBucket{}
	} else if err != nil {
		return err
	}

	buc.requests = append(buc.requests, reqID)
	return i.set(pn, objID, buc)
}

func (i *IndexDB) SetResultRecord(ctx context.Context, pn insolar.PulseNumber, objID insolar.ID, resID insolar.ID) error {
	i.lock.Lock()
	defer i.lock.Unlock()

	buc, err := i.get(pn, objID)
	if err == store.ErrNotFound {
		buc = &IndexBucket{}
	} else if err != nil {
		return err
	}

	buc.results = append(buc.requests, resID)
	return i.set(pn, objID, buc)
}

func (i *IndexDB) set(pn insolar.PulseNumber, objID insolar.ID, bucket *IndexBucket) error {
	key := indexKey{pn: pn, objID: objID}

	return i.db.Set(key, MustEncodeBucket(bucket))
}

func (i *IndexDB) get(pn insolar.PulseNumber, objID insolar.ID) (*IndexBucket, error) {
	buff, err := i.db.Get(indexKey{pn: pn, objID: objID})
	if err == store.ErrNotFound {
		return nil, ErrIndexBucketNotFound

	}
	if err != nil {
		return nil, err
	}
	bucket := MustDecodeBucket(buff)
	return bucket, nil
}

func MustEncodeBucket(buck *IndexBucket) []byte {
	panic("bux")
}

func MustDecodeBucket(buff []byte) *IndexBucket {
	panic("bux")
}

//
// type LifelineDB struct {
// 	lock sync.RWMutex
// 	db   store.DB
// }
//
// type lifelineKey insolar.ID
//
// func (k lifelineKey) Scope() store.Scope {
// 	return store.ScopeLifeline
// }
//
// func (k lifelineKey) ID() []byte {
// 	res := insolar.ID(k)
// 	return (&res).Bytes()
// }
//
// // NewLifelineDB creates new DB storage instance.
// func NewLifelineDB(db store.DB) *LifelineDB {
// 	return &LifelineDB{db: db}
// }
//
// // Set saves new index-value in storage.
// func (i *LifelineDB) Set(ctx context.Context, id insolar.ID, index Lifeline) error {
// 	i.lock.Lock()
// 	defer i.lock.Unlock()
//
// 	if index.Delegates == nil {
// 		index.Delegates = map[insolar.Reference]insolar.Reference{}
// 	}
//
// 	return i.set(id, index)
// }
//
// // ForID returns index for provided id.
// func (i *LifelineDB) ForID(ctx context.Context, id insolar.ID) (index Lifeline, err error) {
// 	i.lock.RLock()
// 	defer i.lock.RUnlock()
//
// 	return i.get(id)
// }
//
// func (i *LifelineDB) set(id insolar.ID, index Lifeline) error {
// 	key := lifelineKey(id)
//
// 	return i.db.Set(key, EncodeIndex(index))
// }
//
// func (i *LifelineDB) get(id insolar.ID) (index Lifeline, err error) {
// 	buff, err := i.db.Get(lifelineKey(id))
// 	if err == store.ErrNotFound {
// 		err = ErrLifelineNotFound
// 		return
// 	}
// 	if err != nil {
// 		return
// 	}
// 	index = MustDecodeIndex(buff)
// 	return
// }
