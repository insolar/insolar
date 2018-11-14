package storage

import (
	"sync"

	"github.com/insolar/insolar/core"
)

type RecentObjectsIndex struct {
	RecentObjects map[string]*RecentObjectsIndexMeta
	lock          sync.Mutex
	DefaultTtl    int
}

type RecentObjectsIndexMeta struct {
	Ttl int
}

func NewRecentObjectsIndex(defaultTtl int) *RecentObjectsIndex {
	return &RecentObjectsIndex{
		RecentObjects: map[string]*RecentObjectsIndexMeta{},
		DefaultTtl:    defaultTtl,
		lock:          sync.Mutex{},
	}
}

func (r *RecentObjectsIndex) AddId(id *core.RecordID) {
	r.lock.Lock()
	defer r.lock.Unlock()

	value, ok := r.RecentObjects[string(id.Bytes())]

	if !ok {
		r.RecentObjects[string(id.Bytes())] = &RecentObjectsIndexMeta{
			Ttl: r.DefaultTtl,
		}
		return
	}

	value.Ttl = r.DefaultTtl
}

func (r *RecentObjectsIndex) RemoveWithTtlMoreThen(ttl int) {
	r.lock.Lock()
	defer r.lock.Unlock()

	for key, value := range r.RecentObjects {
		if value.Ttl == 0 {
			delete(r.RecentObjects, key)
		}
	}
}

func (r *RecentObjectsIndex) Clear() {
	r.lock.Lock()
	defer r.lock.Unlock()

	r.RecentObjects = map[string]*RecentObjectsIndexMeta{}
}
