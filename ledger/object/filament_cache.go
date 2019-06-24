package object

import (
	"context"
	"sync"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/record"
	"github.com/pkg/errors"
)

type cacheStore struct {
	records RecordAccessor

	lock   sync.Mutex
	caches map[insolar.ID]*filamentCache
}

func newCacheStore(r RecordAccessor) *cacheStore {
	return &cacheStore{
		caches:  map[insolar.ID]*filamentCache{},
		records: r,
	}
}

func (c *cacheStore) Get(id insolar.ID) *filamentCache {
	c.lock.Lock()
	defer c.lock.Unlock()

	obj, ok := c.caches[id]
	if !ok {
		obj = newFilamentCache(c.records)
		c.caches[id] = obj
	}

	return obj
}

func (c *cacheStore) Delete(id insolar.ID) {
	c.lock.Lock()
	defer c.lock.Unlock()
	delete(c.caches, id)
}

type filamentCache struct {
	cache map[insolar.ID]record.CompositeFilamentRecord

	records     RecordAccessor
	jetFetcher  jet.Fetcher
	coordinator jet.Coordinator
	sender      bus.Sender
}

func newFilamentCache(r RecordAccessor) *filamentCache {
	return &filamentCache{
		cache:   map[insolar.ID]record.CompositeFilamentRecord{},
		records: r,
	}
}

func (c *filamentCache) Update(recs []record.CompositeFilamentRecord) {
	for _, rec := range recs {
		c.cache[rec.MetaID] = rec
	}
}

func (c *filamentCache) NewIterator(ctx context.Context, from insolar.ID) filamentIterator {
	return filamentIterator{
		currentID: &from,
		cache:     c,
	}
}

func (c *filamentCache) Clear() {
	c.cache = map[insolar.ID]record.CompositeFilamentRecord{}
}

type filamentIterator struct {
	currentID *insolar.ID
	cache     *filamentCache
}

func (i *filamentIterator) PrevID() *insolar.ID {
	return i.currentID
}

func (i *filamentIterator) HasPrev() bool {
	return i.currentID != nil
}

func (i *filamentIterator) Prev(ctx context.Context) (record.CompositeFilamentRecord, error) {
	if i.currentID == nil {
		return record.CompositeFilamentRecord{}, ErrNotFound
	}

	composite, ok := i.cache.cache[*i.currentID]
	if ok {
		virtual := record.Unwrap(composite.Meta.Virtual)
		filament, ok := virtual.(*record.PendingFilament)
		if !ok {
			return record.CompositeFilamentRecord{}, errors.New("failed to convert filament record")
		}
		i.currentID = filament.PreviousRecord
		return composite, nil
	}

	// Fetching filament record.
	filamentRecord, err := i.cache.records.ForID(ctx, *i.currentID)
	if err != nil {
		return record.CompositeFilamentRecord{}, err
	}
	composite.RecordID = *i.currentID
	composite.Record = filamentRecord

	// Fetching primary record.
	virtual := record.Unwrap(filamentRecord.Virtual)
	filament, ok := virtual.(*record.PendingFilament)
	if !ok {
		return record.CompositeFilamentRecord{}, errors.New("failed to convert filament record")
	}
	rec, err := i.cache.records.ForID(ctx, filament.RecordID)
	if err != nil {
		return record.CompositeFilamentRecord{}, err
	}
	composite.RecordID = filament.RecordID
	composite.Record = rec

	// Adding to cache.
	i.cache.cache[*i.currentID] = composite
	i.currentID = filament.PreviousRecord

	return composite, nil
}
