package storage

import (
	"sync"

	"github.com/insolar/insolar/core"
)

// RecentObjectsIndex is a base structure
type RecentObjectsIndex struct {
	recentObjects map[string]*RecentObjectsIndexMeta
	lock          sync.Mutex
	DefaultTTL    int
}

// RecentObjectsIndexMeta contains meta about indexes
type RecentObjectsIndexMeta struct {
	TTL int
}

// NewRecentObjectsIndex creates default RecentObjectsIndex object
func NewRecentObjectsIndex(defaultTTl int) *RecentObjectsIndex {
	return &RecentObjectsIndex{
		recentObjects: map[string]*RecentObjectsIndexMeta{},
		DefaultTTL:    defaultTTl,
		lock:          sync.Mutex{},
	}
}

// AddId adds object to cache
func (r *RecentObjectsIndex) AddId(id *core.RecordID) {
	r.lock.Lock()
	defer r.lock.Unlock()

	value, ok := r.recentObjects[string(id.Bytes())]

	if !ok {
		r.recentObjects[string(id.Bytes())] = &RecentObjectsIndexMeta{
			TTL: r.DefaultTTL,
		}
		return
	}

	value.TTL = r.DefaultTTL
}

// GetObjects returns hot-indexes
func (r *RecentObjectsIndex) GetObjects() map[string]*RecentObjectsIndexMeta {
	r.lock.Lock()
	defer r.lock.Unlock()

	targetMap := map[string]*RecentObjectsIndexMeta{}
	for key, value := range r.recentObjects {
		targetMap[key] = value
	}

	return targetMap
}

// ClearZeroTTL clears objects with zero TTL
func (r *RecentObjectsIndex) ClearZeroTTL() {
	r.lock.Lock()
	defer r.lock.Unlock()

	for key, value := range r.recentObjects {
		if value.TTL == 0 {
			delete(r.recentObjects, key)
		}
	}
}

// Clear clears the whole cache
func (r *RecentObjectsIndex) Clear() {
	r.lock.Lock()
	defer r.lock.Unlock()

	r.recentObjects = map[string]*RecentObjectsIndexMeta{}
}
