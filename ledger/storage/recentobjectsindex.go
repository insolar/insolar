package storage

import (
	"sync"

	"github.com/insolar/insolar/core"
)

type RecentObjectsIndex struct {
	RecentObjects map[string]*RecentObjectsIndexMeta
	lock          sync.Mutex
	DefaultTTL    int
}

type RecentObjectsIndexMeta struct {
	TTL int
}

func NewRecentObjectsIndex(defaultTtl int) *RecentObjectsIndex {
	return &RecentObjectsIndex{
		RecentObjects: map[string]*RecentObjectsIndexMeta{},
		DefaultTTL:    defaultTtl,
		lock:          sync.Mutex{},
	}
}

func (r *RecentObjectsIndex) AddId(id *core.RecordID) {
	r.lock.Lock()
	defer r.lock.Unlock()

	value, ok := r.RecentObjects[string(id.Bytes())]

	if !ok {
		r.RecentObjects[string(id.Bytes())] = &RecentObjectsIndexMeta{
			TTL: r.DefaultTTL,
		}
		return
	}

	value.TTL = r.DefaultTTL
}

func (r *RecentObjectsIndex) RemoveWithTtlMoreThen(ttl int) {
	r.lock.Lock()
	defer r.lock.Unlock()

	for key, value := range r.RecentObjects {
		if value.TTL == 0 {
			delete(r.RecentObjects, key)
		}
	}
}

func (r *RecentObjectsIndex) Clear() {
	r.lock.Lock()
	defer r.lock.Unlock()

	r.RecentObjects = map[string]*RecentObjectsIndexMeta{}
}
