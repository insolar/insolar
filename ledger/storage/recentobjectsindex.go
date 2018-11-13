package storage

import (
	"sync"

	"github.com/insolar/insolar/core"
)

type RecentObjectsIndex struct {
	fetchedObjects map[string]*core.RecordID
	updatedObjects map[string]*core.RecordID

	fetchedObjectsLock sync.Mutex
	updatedObjectsLock sync.Mutex
}

func NewRecentObjectsIndex() *RecentObjectsIndex {
	return &RecentObjectsIndex{
		fetchedObjects: map[string]*core.RecordID{},
		updatedObjects: map[string]*core.RecordID{},
	}
}

func (r *RecentObjectsIndex) addToFetched(id *core.RecordID) {
	r.fetchedObjectsLock.Lock()
	defer r.fetchedObjectsLock.Unlock()

	r.fetchedObjects[id.String()] = id
}

func (r *RecentObjectsIndex) addToUpdated(id *core.RecordID) {
	r.updatedObjectsLock.Lock()
	defer r.updatedObjectsLock.Unlock()

	r.updatedObjects[id.String()] = id
}

func (r *RecentObjectsIndex) clear() {
	r.updatedObjectsLock.Lock()
	defer r.updatedObjectsLock.Unlock()
	r.fetchedObjectsLock.Lock()
	defer r.fetchedObjectsLock.Unlock()

	r.updatedObjects = map[string]*core.RecordID{}
	r.fetchedObjects = map[string]*core.RecordID{}
}
