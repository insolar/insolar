/*
 *    Copyright 2019 Insolar Technologies
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
	indexStorages   map[core.RecordID]*RecentIndexStorage
	pendingStorages map[core.RecordID]*PendingStorage

	indexLock   sync.Mutex
	pendingLock sync.Mutex

	DefaultTTL int
}

// NewRecentStorageProvider creates new provider
func NewRecentStorageProvider(defaultTTL int) *RecentStorageProvider {
	return &RecentStorageProvider{
		DefaultTTL:      defaultTTL,
		indexStorages:   map[core.RecordID]*RecentIndexStorage{},
		pendingStorages: map[core.RecordID]*PendingStorage{},
	}
}

// GetIndexStorage returns a recent indexes for a specific jet
func (p *RecentStorageProvider) GetIndexStorage(ctx context.Context, jetID core.RecordID) recentstorage.RecentIndexStorage {
	p.indexLock.Lock()
	defer p.indexLock.Unlock()

	storage, ok := p.indexStorages[jetID]
	if !ok {
		storage = NewRecentIndexStorage(jetID, p.DefaultTTL)
		p.indexStorages[jetID] = storage
	}
	return storage
}

// GetPendingStorage returns pendings for a specific jet
func (p *RecentStorageProvider) GetPendingStorage(ctx context.Context, jetID core.RecordID) recentstorage.PendingStorage {
	p.pendingLock.Lock()
	defer p.pendingLock.Unlock()

	storage, ok := p.pendingStorages[jetID]
	if !ok {
		storage = NewPendingStorage(jetID)
		p.pendingStorages[jetID] = storage
	}
	return storage
}

// CloneIndexStorage clones indexes from one jet to another one
func (p *RecentStorageProvider) CloneIndexStorage(ctx context.Context, fromJetID, toJetID core.RecordID) {
	p.indexLock.Lock()
	defer p.indexLock.Unlock()

	fromStorage, ok := p.indexStorages[fromJetID]
	if !ok {
		return
	}
	toStorage := &RecentIndexStorage{
		jetID:      toJetID,
		indexes:    map[core.RecordID]recentObjectMeta{},
		DefaultTTL: p.DefaultTTL,
	}
	for k, v := range fromStorage.indexes {
		clone := v
		toStorage.indexes[k] = clone
	}
	p.indexStorages[toJetID] = toStorage
}

// ClonePendingStorage clones pending requests from one jet to another one
func (p *RecentStorageProvider) ClonePendingStorage(ctx context.Context, fromJetID, toJetID core.RecordID) {
	p.pendingLock.Lock()
	defer p.pendingLock.Unlock()

	fromStorage, ok := p.pendingStorages[fromJetID]
	if !ok {
		return
	}
	toStorage := &PendingStorage{
		jetID:    toJetID,
		requests: map[core.RecordID]map[core.RecordID]struct{}{},
	}
	for objID, objRequests := range fromStorage.requests {
		clone := make(map[core.RecordID]struct{}, len(objRequests))
		for reqID, v := range objRequests {
			clone[reqID] = v
		}
		toStorage.requests[objID] = clone
	}
	p.pendingStorages[toJetID] = toStorage
}

// DecreaseIndexesTTL decrease ttl of all indexes in all storages
// If storage contains indexes with zero ttl, they are removed from storage and returned to a caller
// If there are no indexes with ttl more then 0, storage is removed
func (p *RecentStorageProvider) DecreaseIndexesTTL(ctx context.Context) map[core.RecordID][]core.RecordID {
	p.indexLock.Lock()
	defer p.indexLock.Unlock()

	resMapLock := sync.Mutex{}
	resMap := map[core.RecordID][]core.RecordID{}
	wg := sync.WaitGroup{}
	wg.Add(len(p.indexStorages))

	for jetID, storage := range p.indexStorages {
		go func(jetID core.RecordID, s *RecentIndexStorage) {
			res := s.DecreaseIndexTTL(ctx)

			if len(res) > 0 {
				resMapLock.Lock()
				resMap[jetID] = res
				resMapLock.Unlock()
			}

			if len(s.indexes) == 0 {
				delete(p.indexStorages, jetID)
			}

			wg.Done()
		}(jetID, storage)
	}

	wg.Wait()

	return resMap
}

// RemoveIndexStorage removes indexes for a specific jet from provider
// If there is a reference to RecentIndexStorage somewhere, it won't be affected
func (p *RecentStorageProvider) RemoveIndexStorage(ctx context.Context, id core.RecordID) {
	p.indexLock.Lock()
	defer p.indexLock.Unlock()

	if storage, ok := p.indexStorages[id]; ok {
		storage.lock.Lock()
		defer storage.lock.Unlock()

		ctx = insmetrics.InsertTag(ctx, tagJet, storage.jetID.DebugString())
		stats.Record(ctx,
			statRecentStorageObjectsRemoved.M(int64(len(storage.indexes))),
		)

		delete(p.indexStorages, id)
	}
}

// RemovePendingStorage removes pending requests for a specific jet from provider
// If there is a reference to RecentIndexStorage somewhere, it won't be affected
func (p *RecentStorageProvider) RemovePendingStorage(ctx context.Context, id core.RecordID) {
	p.pendingLock.Lock()
	defer p.pendingLock.Unlock()

	if storage, ok := p.pendingStorages[id]; ok {
		storage.lock.Lock()
		defer storage.lock.Unlock()

		ctx = insmetrics.InsertTag(ctx, tagJet, storage.jetID.DebugString())
		stats.Record(ctx,
			statRecentStoragePendingsRemoved.M(int64(len(storage.requests))),
		)

		delete(p.pendingStorages, id)
	}
}

// RecentIndexStorage is an implementation of RecentIndexStorage interface
// This is a in-memory cache for indexes` ids
type RecentIndexStorage struct {
	jetID      core.RecordID
	indexes    map[core.RecordID]recentObjectMeta
	lock       sync.Mutex
	DefaultTTL int
}

type recentObjectMeta struct {
	ttl int
}

// NewRecentIndexStorage creates new *RecentIndexStorage
func NewRecentIndexStorage(jetID core.RecordID, defaultTTL int) *RecentIndexStorage {
	return &RecentIndexStorage{jetID: jetID, DefaultTTL: defaultTTL}
}

// AddObject adds index's id to an in-memory cache and sets DefaultTTL for it
func (r *RecentIndexStorage) AddObject(ctx context.Context, id core.RecordID) {
	r.AddObjectWithTLL(ctx, id, r.DefaultTTL)
}

// AddObjectWithTLL adds index's id to an in-memory cache with provided ttl
func (r *RecentIndexStorage) AddObjectWithTLL(ctx context.Context, id core.RecordID, ttl int) {
	r.lock.Lock()
	defer r.lock.Unlock()

	r.indexes[id] = recentObjectMeta{ttl: ttl}

	ctx = insmetrics.InsertTag(ctx, tagJet, r.jetID.DebugString())
	stats.Record(ctx, statRecentStorageObjectsAdded.M(1))
}

// GetObjects returns deep copy of indexes' ids
func (r *RecentIndexStorage) GetObjects() map[core.RecordID]int {
	r.lock.Lock()
	defer r.lock.Unlock()

	targetMap := make(map[core.RecordID]int, len(r.indexes))
	for key, value := range r.indexes {
		targetMap[key] = value.ttl
	}

	return targetMap
}

func (r *RecentStorage) FilterNotExistWithLock(
	ctx context.Context,
	candidates []core.RecordID,
	lockedFn func(filtered []core.RecordID),
) {
	r.objectLock.Lock()
	markedCandidates := r.markForDelete(candidates)
	lockedFn(markedCandidates)
	r.objectLock.Unlock()
}

// DecreaseIndexTTL decreases ttls and remove indexes with 0 ttl
// Removed indexes will be returned as a functon's result
func (r *RecentIndexStorage) DecreaseIndexTTL(ctx context.Context) []core.RecordID {
	r.lock.Lock()
	defer r.lock.Unlock()

	var clearedObjects []core.RecordID
	for key, value := range r.indexes {
		value.ttl--
		if value.ttl == 0 {
			clearedObjects = append(clearedObjects, key)
			delete(r.indexes, key)
			continue
		}
		r.indexes[key] = value
	}
	return clearedObjects
}

// PendingStorage contains indexes of unclosed requests (pendings) for a specific object id
type PendingStorage struct {
	jetID    core.RecordID
	requests map[core.RecordID]map[core.RecordID]struct{}
	lock     sync.RWMutex
}

// NewPendingStorage creates *PendingStorage
func NewPendingStorage(jetID core.RecordID) *PendingStorage {
	return &PendingStorage{jetID: jetID}
}

// AddPendingRequest adds an id of pending request to memory
// The id stores in a collection ids of a specific object
func (r *PendingStorage) AddPendingRequest(ctx context.Context, obj, req core.RecordID) {
	r.lock.Lock()
	defer r.lock.Unlock()

	if _, ok := r.requests[obj]; !ok {
		r.requests[obj] = map[core.RecordID]struct{}{}
	}
	r.requests[obj][req] = struct{}{}

	ctx = insmetrics.InsertTag(ctx, tagJet, r.jetID.DebugString())
	stats.Record(ctx, statRecentStoragePendingsAdded.M(1))
}

// GetRequests returns a deep-copy of requests collections
func (r *PendingStorage) GetRequests() map[core.RecordID]map[core.RecordID]struct{} {
	r.lock.RLock()
	defer r.lock.RUnlock()

	requestsClone := make(map[core.RecordID]map[core.RecordID]struct{})
	for objID, objRequests := range r.requests {
		objRequestsClone := make(map[core.RecordID]struct{}, len(objRequests))
		for reqID, v := range objRequests {
			objRequestsClone[reqID] = v
		}
		requestsClone[objID] = objRequestsClone
	}

	return requestsClone
}

// GetRequestsForObject returns a deep-copy of requests collection for a specific object
func (r *PendingStorage) GetRequestsForObject(obj core.RecordID) []core.RecordID {
	r.lock.RLock()
	defer r.lock.RUnlock()

	forObject, ok := r.requests[obj]
	if !ok {
		return nil
	}
	results := make([]core.RecordID, 0, len(forObject))
	for reqID := range forObject {
		results = append(results, reqID)
	}

	return results
}

// RemovePendingRequest removes an id from cache
func (r *PendingStorage) RemovePendingRequest(ctx context.Context, obj, req core.RecordID) {
	r.lock.Lock()
	defer r.lock.Unlock()

	if _, ok := r.requests[obj]; !ok {
		return
	}
	delete(r.requests[obj], req)
	if len(r.requests[obj]) == 0 {
		delete(r.requests, obj)
	}

	ctx = insmetrics.InsertTag(ctx, tagJet, r.jetID.DebugString())
	stats.Record(ctx, statRecentStoragePendingsRemoved.M(1))
}

func (r *RecentStorage) markForDelete(candidates []core.RecordID) []core.RecordID {
	result := make([]core.RecordID, 0, len(candidates))

	for _, c := range candidates {
		_, exists := r.recentObjects[c]
		if !exists {
			result = append(result, c)
		}
	}
	return result
}
