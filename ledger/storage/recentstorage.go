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
	"context"
	"sync"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/instrumentation/insmetrics"
	"github.com/insolar/insolar/ledger/recentstorage"
	"go.opencensus.io/stats"
)

// RecentStorageProvider provides a recent storage for jet
type RecentStorageProvider struct {
	// TODO: @andreyromancev. 15.01.19. Use byte array for key.
	storage    map[core.RecordID]*RecentStorage
	lock       sync.Mutex
	DefaultTTL int
}

// NewRecentStorageProvider creates new provider
func NewRecentStorageProvider(defaultTTL int) *RecentStorageProvider {
	return &RecentStorageProvider{DefaultTTL: defaultTTL, storage: map[core.RecordID]*RecentStorage{}}
}

// GetStorage returns a recent storage for jet
func (p *RecentStorageProvider) GetStorage(ctx context.Context, jetID core.RecordID) recentstorage.RecentStorage {
	p.lock.Lock()
	defer p.lock.Unlock()

	storage, ok := p.storage[jetID]
	if !ok {
		storage = NewRecentStorage(jetID, p.DefaultTTL)
		p.storage[jetID] = storage
	}
	return storage
}

// CloneStorage clones a recent storage from one jet to another
func (p *RecentStorageProvider) CloneStorage(ctx context.Context, fromJetID, toJetID core.RecordID) {
	p.lock.Lock()
	defer p.lock.Unlock()

	fromStorage, ok := p.storage[fromJetID]
	if !ok {
		return
	}
	toStorage := &RecentStorage{
		jetID:           toJetID,
		recentObjects:   make(map[core.RecordID]recentObjectMeta, len(fromStorage.recentObjects)),
		pendingRequests: make(map[core.RecordID]map[core.RecordID]struct{}, len(fromStorage.pendingRequests)),
		DefaultTTL:      p.DefaultTTL,
		objectLock:      sync.Mutex{},
	}
	for k, v := range fromStorage.recentObjects {
		clone := v
		toStorage.recentObjects[k] = clone
	}
	for objID, objRequests := range fromStorage.pendingRequests {
		clone := make(map[core.RecordID]struct{}, len(objRequests))
		for reqID, v := range objRequests {
			clone[reqID] = v
		}
		toStorage.pendingRequests[objID] = clone
	}
	p.storage[toJetID] = toStorage
}

// RemoveStorage removes storage from provider
func (p *RecentStorageProvider) RemoveStorage(ctx context.Context, id core.RecordID) {
	p.lock.Lock()
	defer p.lock.Unlock()

	if storage, ok := p.storage[id]; ok {
		storage.objectLock.Lock()
		defer storage.objectLock.Unlock()

		ctx = insmetrics.InsertTag(ctx, tagJet, storage.jetID.DebugString())
		stats.Record(ctx,
			statRecentStorageObjectsRemoved.M(int64(len(storage.recentObjects))),
			statRecentStoragePendingsRemoved.M(int64(len(storage.pendingRequests))),
		)

		delete(p.storage, id)
	}
}

// RecentStorage is a base structure
type RecentStorage struct {
	jetID           core.RecordID
	recentObjects   map[core.RecordID]recentObjectMeta
	objectLock      sync.Mutex
	pendingRequests map[core.RecordID]map[core.RecordID]struct{}
	requestLock     sync.RWMutex
	DefaultTTL      int
}

type recentObjectMeta struct {
	ttl int
}

// NewRecentStorage creates default RecentStorage object
func NewRecentStorage(jetID core.RecordID, defaultTTL int) *RecentStorage {
	return &RecentStorage{
		jetID:           jetID,
		recentObjects:   map[core.RecordID]recentObjectMeta{},
		pendingRequests: map[core.RecordID]map[core.RecordID]struct{}{},
		DefaultTTL:      defaultTTL,
		objectLock:      sync.Mutex{},
	}
}

// AddObject adds object to cache
func (r *RecentStorage) AddObject(ctx context.Context, id core.RecordID) {
	r.AddObjectWithTLL(ctx, id, r.DefaultTTL)
}

// AddObjectWithTLL adds object with specified TTL to the cache
func (r *RecentStorage) AddObjectWithTLL(ctx context.Context, id core.RecordID, ttl int) {
	r.objectLock.Lock()
	defer r.objectLock.Unlock()

	r.recentObjects[id] = recentObjectMeta{ttl: ttl}

	ctx = insmetrics.InsertTag(ctx, tagJet, r.jetID.DebugString())
	stats.Record(ctx, statRecentStorageObjectsAdded.M(1))
}

// AddPendingRequest adds request to cache.
func (r *RecentStorage) AddPendingRequest(ctx context.Context, obj, req core.RecordID) {
	r.requestLock.Lock()
	defer r.requestLock.Unlock()

	if _, ok := r.pendingRequests[obj]; !ok {
		r.pendingRequests[obj] = map[core.RecordID]struct{}{}
	}
	r.pendingRequests[obj][req] = struct{}{}

	ctx = insmetrics.InsertTag(ctx, tagJet, r.jetID.DebugString())
	stats.Record(ctx, statRecentStoragePendingsAdded.M(1))
}

// RemovePendingRequest removes request from cache.
func (r *RecentStorage) RemovePendingRequest(ctx context.Context, obj, req core.RecordID) {
	r.requestLock.Lock()
	defer r.requestLock.Unlock()

	if _, ok := r.pendingRequests[obj]; !ok {
		return
	}
	delete(r.pendingRequests[obj], req)
	if len(r.pendingRequests[obj]) == 0 {
		delete(r.pendingRequests, obj)
	}

	ctx = insmetrics.InsertTag(ctx, tagJet, r.jetID.DebugString())
	stats.Record(ctx, statRecentStoragePendingsRemoved.M(1))
}

// GetObjects returns object hot-indexes.
func (r *RecentStorage) GetObjects() map[core.RecordID]int {
	r.objectLock.Lock()
	defer r.objectLock.Unlock()

	targetMap := make(map[core.RecordID]int, len(r.recentObjects))
	for key, value := range r.recentObjects {
		targetMap[key] = value.ttl
	}

	return targetMap
}

// GetRequests returns request hot-indexes.
func (r *RecentStorage) GetRequests() map[core.RecordID]map[core.RecordID]struct{} {
	r.requestLock.RLock()
	defer r.requestLock.RUnlock()

	requestsClone := make(map[core.RecordID]map[core.RecordID]struct{})
	for objID, objRequests := range r.pendingRequests {
		objRequestsClone := make(map[core.RecordID]struct{}, len(objRequests))
		for reqID, v := range objRequests {
			objRequestsClone[reqID] = v
		}
		requestsClone[objID] = objRequestsClone
	}

	return requestsClone
}

// GetRequestsForObject returns request hot-indexes for object.
func (r *RecentStorage) GetRequestsForObject(obj core.RecordID) []core.RecordID {
	r.requestLock.RLock()
	defer r.requestLock.RUnlock()

	forObject, ok := r.pendingRequests[obj]
	if !ok {
		return nil
	}
	results := make([]core.RecordID, 0, len(forObject))
	for reqID := range forObject {
		results = append(results, reqID)
	}

	return results
}

// IsRecordIDCached checks recordID inside caches
func (r *RecentStorage) IsRecordIDCached(obj core.RecordID) bool {
	r.objectLock.Lock()
	_, ok := r.recentObjects[obj]
	if ok {
		r.objectLock.Unlock()
		return ok
	}
	r.objectLock.Unlock()

	r.requestLock.RLock()
	_, ok = r.pendingRequests[obj]
	r.requestLock.RUnlock()
	return ok
}

// DecreaseTTL decreases ttl and clears objects if their ttl is zero
func (r *RecentStorage) DecreaseTTL(ctx context.Context) {
	r.objectLock.Lock()
	defer r.objectLock.Unlock()

	for key, value := range r.recentObjects {
		value.ttl--
		if value.ttl == 0 {
			delete(r.recentObjects, key)
			continue
		}
		r.recentObjects[key] = value
	}
}
