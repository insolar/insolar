package object

import (
	"context"
	"sync"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/record"
	"github.com/pkg/errors"
)

type filamentCache struct {
	// lock  sync.RWMutex
	records RecordAccessor
	cache   map[insolar.PulseNumber][]record.CompositeFilamentRecord
}

func newFilamentCache(r RecordAccessor) *filamentCache {
	return &filamentCache{
		cache:   map[insolar.PulseNumber][]record.CompositeFilamentRecord{},
		records: r,
	}
}

func (c *filamentCache) Append(rec record.CompositeFilamentRecord) error {
	// c.lock.RLock()
	// defer c.lock.RUnlock()
	// if rec.MetaID.Pulse() < c.last {
	// 	return errors.New("wrong pulse")
	// }

	recs, ok := c.cache[rec.MetaID.Pulse()]
	if !ok {
		c.cache[rec.MetaID.Pulse()] = []record.CompositeFilamentRecord{rec}
		return nil
	}

	c.cache[rec.MetaID.Pulse()] = append(recs, rec)
	return nil
}

func (c *filamentCache) Set(pn insolar.PulseNumber, recs []record.CompositeFilamentRecord) {
	c.cache[pn] = recs
}

func (c *filamentCache) Get(ctx context.Context, from insolar.ID) ([]record.CompositeFilamentRecord, error) {
	// c.lock.RLock()
	// defer c.lock.RUnlock()

	forPulse := from.Pulse()
	recs, ok := c.cache[forPulse]
	if ok {
		return recs, nil
	}

	// Not found in cache. Trying to fetch from primary storage.
	iter := &from
	for iter != nil && iter.Pulse() == forPulse {
		var composite record.CompositeFilamentRecord
		// Fetching filament record.
		filamentRecord, err := c.records.ForID(ctx, *iter)
		if err != nil {
			return nil, err
		}
		composite.RecordID = *iter
		composite.Record = filamentRecord

		// Fetching primary.
		virtual := record.Unwrap(filamentRecord.Virtual)
		filament, ok := virtual.(*record.PendingFilament)
		if !ok {
			return nil, errors.New("failed to convert filament record")
		}
		rec, err := c.records.ForID(ctx, filament.RecordID)
		if err != nil {
			return nil, err
		}
		composite.RecordID = filament.RecordID
		composite.Record = rec

		// Adding to cache.
		recs = append(recs, composite)
		iter = filament.PreviousRecord
	}

	c.cache[forPulse] = recs

	return recs, nil
}

func (c *filamentCache) Delete(pn insolar.PulseNumber) {
	// c.lock.Lock()
	// defer c.lock.Unlock()

	delete(c.cache, pn)
}

type objectCache struct {
	Incoming *filamentCache
}

func newObjectCache(recs RecordAccessor) *objectCache {
	return &objectCache{
		Incoming: newFilamentCache(recs),
	}
}

type cacheStore struct {
	records RecordAccessor

	lock   sync.Mutex
	caches map[insolar.ID]*objectCache
}

func newCacheStore(recs RecordAccessor) *cacheStore {
	return &cacheStore{
		caches:  map[insolar.ID]*objectCache{},
		records: recs,
	}
}

func (c *cacheStore) Get(id insolar.ID) *objectCache {
	c.lock.Lock()
	defer c.lock.Unlock()

	obj, ok := c.caches[id]
	if !ok {
		obj = newObjectCache(c.records)
		c.caches[id] = obj
	}

	return obj
}

func (c *cacheStore) Delete(id insolar.ID) {
	c.lock.Lock()
	defer c.lock.Unlock()
	delete(c.caches, id)
}
