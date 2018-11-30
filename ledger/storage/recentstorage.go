/*
 *    Copyright 2018 Insolar
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package storage

import (
	"sync"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/core/message"
)

// RecentStorage is a base structure
type RecentStorage struct {
	recentObjects   map[core.RecordID]*message.RecentObjectsIndexMeta
	objectLock      sync.Mutex
	pendingRequests map[core.RecordID]struct{}
	requestLock     sync.Mutex
	DefaultTTL      int
}

// NewRecentStorage creates default RecentStorage object
func NewRecentStorage(defaultTTL int) *RecentStorage {
	return &RecentStorage{
		recentObjects:   map[core.RecordID]*message.RecentObjectsIndexMeta{},
		pendingRequests: map[core.RecordID]struct{}{},
		DefaultTTL:      defaultTTL,
		objectLock:      sync.Mutex{},
	}
}

// AddObject adds object to cache
func (r *RecentStorage) AddObject(id core.RecordID) {
	r.objectLock.Lock()
	defer r.objectLock.Unlock()

	value, ok := r.recentObjects[id]

	if !ok {
		r.recentObjects[id] = &message.RecentObjectsIndexMeta{
			TTL: r.DefaultTTL,
		}
		return
	}

	value.TTL = r.DefaultTTL
}

// AddObjectWithMeta adds object with meta to cache
func (r *RecentStorage) AddObjectWithMeta(id core.RecordID, meta *message.RecentObjectsIndexMeta) {
	r.objectLock.Lock()
	defer r.objectLock.Unlock()
	if meta == nil {
		meta = &message.RecentObjectsIndexMeta{
			TTL: r.DefaultTTL,
		}
	}

	r.recentObjects[id] = meta
}

// AddPendingRequest adds request to cache.
func (r *RecentStorage) AddPendingRequest(id core.RecordID) {
	r.requestLock.Lock()
	defer r.requestLock.Unlock()

	if _, ok := r.pendingRequests[id]; !ok {
		r.pendingRequests[id] = struct{}{}
		return
	}
}

// RemovePendingRequest removes request from cache.
func (r *RecentStorage) RemovePendingRequest(id core.RecordID) {
	r.requestLock.Lock()
	defer r.requestLock.Unlock()

	delete(r.pendingRequests, id)
}

// GetObjects returns object hot-indexes.
func (r *RecentStorage) GetObjects() map[core.RecordID]*message.RecentObjectsIndexMeta {
	r.objectLock.Lock()
	defer r.objectLock.Unlock()

	targetMap := make(map[core.RecordID]*message.RecentObjectsIndexMeta, len(r.recentObjects))
	for key, value := range r.recentObjects {
		targetMap[key] = value
	}

	return targetMap
}

// GetRequests returns request hot-indexes.
func (r *RecentStorage) GetRequests() []core.RecordID {
	r.requestLock.Lock()
	defer r.requestLock.Unlock()

	requests := make([]core.RecordID, 0, len(r.pendingRequests))
	for id := range r.pendingRequests {
		requests = append(requests, id)
	}

	return requests
}

// ClearZeroTTLObjects clears objects with zero TTL
func (r *RecentStorage) ClearZeroTTLObjects() {
	r.objectLock.Lock()
	defer r.objectLock.Unlock()

	for key, value := range r.recentObjects {
		if value.TTL == 0 {
			delete(r.recentObjects, key)
		}
	}
}

// ClearObjects clears the whole cache
func (r *RecentStorage) ClearObjects() {
	r.objectLock.Lock()
	defer r.objectLock.Unlock()

	r.recentObjects = map[core.RecordID]*message.RecentObjectsIndexMeta{}
	r.pendingRequests = map[core.RecordID]struct{}{}
}
