package storage

import (
	"sync"

	"github.com/insolar/insolar/core"
)

type RecentObjectsIndex struct {
	recentObjects map[string]struct{}
	lock          sync.Mutex
}

func NewRecentObjectsIndex() *RecentObjectsIndex {
	return &RecentObjectsIndex{
		recentObjects: map[string]struct{}{},
	}
}

func (r *RecentObjectsIndex) addId(id *core.RecordID) {
	r.lock.Lock()
	defer r.lock.Unlock()

	r.recentObjects[string(id.Bytes())] = struct{}{}
}

func (r *RecentObjectsIndex) clear() {
	r.lock.Lock()
	defer r.lock.Unlock()

	r.recentObjects = map[string]struct{}{}
}
