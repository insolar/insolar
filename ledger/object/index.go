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
	SetRequest(ctx context.Context, pn insolar.PulseNumber, objID insolar.ID, reqID insolar.ID)
	SetResultRecord(ctx context.Context, pn insolar.PulseNumber, objID insolar.ID, resID insolar.ID)
}

type IndexStateModifier interface {
	SetLifelineUsage(ctx context.Context, pn insolar.PulseNumber, objID insolar.ID) error
}

type indexBucket struct {
	lifelineLock sync.RWMutex
	lifeline     *Lifeline

	requestLock sync.RWMutex
	requests    []insolar.ID

	resultLock sync.RWMutex
	results    []insolar.ID
}

func (i *indexBucket) getLifeline() (*Lifeline, error) {
	i.lifelineLock.RLock()
	defer i.lifelineLock.RUnlock()
	if i.lifeline == nil {
		return nil, ErrLifelineNotFound
	}

	return i.lifeline, nil
}

func (i *indexBucket) setLifeline(lifeline *Lifeline) {
	i.lifelineLock.Lock()
	defer i.lifelineLock.Unlock()

	i.lifeline = lifeline
}

func (i *indexBucket) setRequest(reqID insolar.ID) {
	i.requestLock.Lock()
	defer i.requestLock.Unlock()

	i.requests = append(i.requests, reqID)
}

func (i *indexBucket) setResult(resID insolar.ID) {
	i.resultLock.Lock()
	defer i.resultLock.Unlock()

	i.results = append(i.results, resID)
}

type InMemoryIndex struct {
	bucketsLock sync.RWMutex
	buckets     map[insolar.PulseNumber]map[insolar.ID]*indexBucket
}

func NewInMemoryIndex() *InMemoryIndex {
	return &InMemoryIndex{buckets: map[insolar.PulseNumber]map[insolar.ID]*indexBucket{}}
}

func (i *InMemoryIndex) getBucket(ctx context.Context, pn insolar.PulseNumber, objID insolar.ID) *indexBucket {
	i.bucketsLock.Lock()
	defer i.bucketsLock.Unlock()

	var objsByPn map[insolar.ID]*indexBucket
	objsByPn, ok := i.buckets[pn]
	if !ok {
		objsByPn = map[insolar.ID]*indexBucket{}
		i.buckets[pn] = objsByPn
	}

	bucket := objsByPn[objID]
	if bucket == nil {
		bucket = &indexBucket{
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

func (i *InMemoryIndex) SetRequest(ctx context.Context, pn insolar.PulseNumber, objID insolar.ID, reqID insolar.ID) {
	b := i.getBucket(ctx, pn, objID)
	b.setRequest(reqID)
}

func (i *InMemoryIndex) SetResultRecord(ctx context.Context, pn insolar.PulseNumber, objID insolar.ID, resID insolar.ID) {
	b := i.getBucket(ctx, pn, objID)
	b.setResult(resID)
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

type InDBIndex struct {
	lock sync.RWMutex
	db   store.DB
}

func NewInDBIndex(db store.DB) *LifelineDB {
	return &LifelineDB{db: db}
}

func (*InDBIndex) SetLifeline(ctx context.Context, pn insolar.PulseNumber, objID insolar.ID, lifeline Lifeline) error {
	panic("implement me")
}

func (*InDBIndex) SetRequest(ctx context.Context, pn insolar.PulseNumber, objID insolar.ID, reqID insolar.ID) {
	panic("implement me")
}

func (*InDBIndex) SetResultRecord(ctx context.Context, pn insolar.PulseNumber, objID insolar.ID, resID insolar.ID) {
	panic("implement me")
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
// // NewIndexDB creates new DB storage instance.
// func NewIndexDB(db store.DB) *LifelineDB {
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
