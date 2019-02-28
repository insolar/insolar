/*
 *    Copyright 2019 Insolar
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

package recentstorage

import (
	"context"
	"sync"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/instrumentation/insmetrics"
	"go.opencensus.io/stats"
)

// RecentStorageProvider provides a recent storage for jet
type RecentStorageProvider struct { //nolint: golint
	indexStorages   map[core.RecordID]*RecentIndexStorageConcrete
	pendingStorages map[core.RecordID]*PendingStorageConcrete

	indexLock   sync.Mutex
	pendingLock sync.Mutex

	DefaultTTL int
}

// NewRecentStorageProvider creates new provider
func NewRecentStorageProvider(defaultTTL int) *RecentStorageProvider {
	return &RecentStorageProvider{
		DefaultTTL:      defaultTTL,
		indexStorages:   map[core.RecordID]*RecentIndexStorageConcrete{},
		pendingStorages: map[core.RecordID]*PendingStorageConcrete{},
	}
}

// GetIndexStorage returns a recent indexes for a specific jet
func (p *RecentStorageProvider) GetIndexStorage(ctx context.Context, jetID core.RecordID) RecentIndexStorage {
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
func (p *RecentStorageProvider) GetPendingStorage(ctx context.Context, jetID core.RecordID) PendingStorage {
	p.pendingLock.Lock()
	defer p.pendingLock.Unlock()

	storage, ok := p.pendingStorages[jetID]
	if !ok {
		storage = NewPendingStorage(jetID)
		p.pendingStorages[jetID] = storage
	}
	return storage
}

// Count returns count of pendings in all storages
func (p *RecentStorageProvider) Count() int {
	p.pendingLock.Lock()
	defer p.pendingLock.Unlock()

	count := 0
	for _, storage := range p.pendingStorages {
		count += len(storage.requests)
	}

	return count
}

// CloneIndexStorage clones indexes from one jet to another one
func (p *RecentStorageProvider) CloneIndexStorage(ctx context.Context, fromJetID, toJetID core.RecordID) {
	p.indexLock.Lock()
	defer p.indexLock.Unlock()

	fromStorage, ok := p.indexStorages[fromJetID]
	if !ok {
		return
	}
	toStorage := &RecentIndexStorageConcrete{
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
	fromStorage, ok := p.pendingStorages[fromJetID]
	p.pendingLock.Unlock()

	if !ok {
		return
	}

	fromStorage.lock.RLock()
	toStorage := &PendingStorageConcrete{
		jetID:    toJetID,
		requests: map[core.RecordID]*lockedPendingObjectContext{},
	}
	for objID, pendingContext := range fromStorage.requests {
		if len(pendingContext.Context.Requests) == 0 {
			continue
		}

		pendingContext.lock.Lock()

		clone := PendingObjectContext{
			Active:   pendingContext.Context.Active,
			Requests: []core.RecordID{},
		}

		clone.Requests = append(clone.Requests, pendingContext.Context.Requests...)
		toStorage.requests[objID] = &lockedPendingObjectContext{Context: &clone}

		pendingContext.lock.Unlock()
	}
	fromStorage.lock.RUnlock()

	p.pendingLock.Lock()
	p.pendingStorages[toJetID] = toStorage
	p.pendingLock.Unlock()
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
		go func(jetID core.RecordID, s *RecentIndexStorageConcrete) {
			res := s.DecreaseIndexTTL(ctx)

			if len(res) > 0 {
				resMapLock.Lock()
				resMap[jetID] = res
				resMapLock.Unlock()
			}

			wg.Done()
		}(jetID, storage)
	}

	wg.Wait()

	for jetID, storage := range p.indexStorages {
		if len(storage.indexes) == 0 {
			delete(p.indexStorages, jetID)
		}
	}

	return resMap
}

// RemovePendingStorage removes pending requests for a specific jet from provider
// If there is a reference to RecentIndexStorage somewhere, it won't be affected
func (p *RecentStorageProvider) RemovePendingStorage(ctx context.Context, id core.RecordID) {
	p.pendingLock.Lock()
	defer p.pendingLock.Unlock()

	if storage, ok := p.pendingStorages[id]; ok {

		ctx = insmetrics.InsertTag(ctx, tagJet, storage.jetID.DebugString())
		stats.Record(ctx,
			statRecentStoragePendingsRemoved.M(int64(len(storage.requests))),
		)

		delete(p.pendingStorages, id)
	}
}

// RecentIndexStorageConcrete is an implementation of RecentIndexStorage interface
// This is a in-memory cache for indexes` ids
type RecentIndexStorageConcrete struct {
	jetID      core.RecordID
	indexes    map[core.RecordID]recentObjectMeta
	lock       sync.Mutex
	DefaultTTL int
}

type recentObjectMeta struct {
	ttl int
}

// NewRecentIndexStorage creates new *RecentIndexStorage
func NewRecentIndexStorage(jetID core.RecordID, defaultTTL int) *RecentIndexStorageConcrete {
	return &RecentIndexStorageConcrete{jetID: jetID, DefaultTTL: defaultTTL, indexes: map[core.RecordID]recentObjectMeta{}}
}

// AddObject adds index's id to an in-memory cache and sets DefaultTTL for it
func (r *RecentIndexStorageConcrete) AddObject(ctx context.Context, id core.RecordID) {
	r.AddObjectWithTLL(ctx, id, r.DefaultTTL)
}

// AddObjectWithTLL adds index's id to an in-memory cache with provided ttl
func (r *RecentIndexStorageConcrete) AddObjectWithTLL(ctx context.Context, id core.RecordID, ttl int) {
	r.lock.Lock()
	defer r.lock.Unlock()

	if ttl < 0 {
		inslogger.FromContext(ctx).Error("below zero ttl happened")
		panic("below zero ttl happened")
	}

	r.indexes[id] = recentObjectMeta{ttl: ttl}

	ctx = insmetrics.InsertTag(ctx, tagJet, r.jetID.DebugString())
	stats.Record(ctx, statRecentStorageObjectsAdded.M(1))
}

// GetObjects returns deep copy of indexes' ids
func (r *RecentIndexStorageConcrete) GetObjects() map[core.RecordID]int {
	r.lock.Lock()
	defer r.lock.Unlock()

	targetMap := make(map[core.RecordID]int, len(r.indexes))
	for key, value := range r.indexes {
		targetMap[key] = value.ttl
	}

	return targetMap
}

// FilterNotExistWithLock filters candidates from params
// found indexes, which aren't removed from cache
// and the passes it to lockedFn
func (r *RecentIndexStorageConcrete) FilterNotExistWithLock(
	ctx context.Context,
	candidates []core.RecordID,
	lockedFn func(filtered []core.RecordID),
) {
	r.lock.Lock()
	markedCandidates := r.markForDelete(candidates)
	lockedFn(markedCandidates)
	r.lock.Unlock()
}

func (r *RecentIndexStorageConcrete) markForDelete(candidates []core.RecordID) []core.RecordID {
	result := make([]core.RecordID, 0, len(candidates))

	for _, c := range candidates {
		_, exists := r.indexes[c]
		if !exists {
			result = append(result, c)
		}
	}
	return result
}

// DecreaseIndexTTL decreases ttls and remove indexes with 0 ttl
// Removed indexes will be returned as a functon's result
func (r *RecentIndexStorageConcrete) DecreaseIndexTTL(ctx context.Context) []core.RecordID {
	r.lock.Lock()
	defer r.lock.Unlock()

	var clearedObjects []core.RecordID
	for key, value := range r.indexes {
		value.ttl--
		if value.ttl <= 0 {
			clearedObjects = append(clearedObjects, key)
			delete(r.indexes, key)
			continue
		}
		r.indexes[key] = value
	}
	return clearedObjects
}

// PendingStorageConcrete contains indexes of unclosed requests (pendings) for a specific object id
type PendingStorageConcrete struct {
	lock sync.RWMutex

	jetID core.RecordID

	requests map[core.RecordID]*lockedPendingObjectContext
}

// PendingObjectContext contains a list of requests for an object
// Also it contains a boolean-flag for determination object's status
// If It's false, current LME, when it gets a hot-data, needs to send
// notifications about forgotten requests
type PendingObjectContext struct {
	Active   bool
	Requests []core.RecordID
}

type lockedPendingObjectContext struct {
	Context *PendingObjectContext
	lock    sync.RWMutex
}

// NewPendingStorage creates *PendingStorage
func NewPendingStorage(jetID core.RecordID) *PendingStorageConcrete {
	return &PendingStorageConcrete{
		jetID:    jetID,
		requests: map[core.RecordID]*lockedPendingObjectContext{},
	}
}

// AddPendingRequest adds an id of pending request to memory
// The id stores in a collection ids of a specific object
func (r *PendingStorageConcrete) AddPendingRequest(ctx context.Context, obj, req core.RecordID) {
	r.lock.Lock()
	defer r.lock.Unlock()

	var objectContext *lockedPendingObjectContext
	var ok bool
	if objectContext, ok = r.requests[obj]; !ok {
		objectContext = &lockedPendingObjectContext{
			lock: sync.RWMutex{},
			Context: &PendingObjectContext{
				Active:   true,
				Requests: []core.RecordID{},
			},
		}
		r.requests[obj] = objectContext
	}

	objectContext.lock.Lock()
	defer objectContext.lock.Unlock()

	objectContext.Context.Requests = append(objectContext.Context.Requests, req)

	ctx = insmetrics.InsertTag(ctx, tagJet, r.jetID.DebugString())
	stats.Record(ctx, statRecentStoragePendingsAdded.M(1))
}

// SetContextToObject add a context to a provided object
func (r *PendingStorageConcrete) SetContextToObject(ctx context.Context, obj core.RecordID, objContext PendingObjectContext) {
	r.lock.Lock()
	defer r.lock.Unlock()

	r.requests[obj] = &lockedPendingObjectContext{
		Context: &objContext,
	}
	ctx = insmetrics.InsertTag(ctx, tagJet, r.jetID.DebugString())
	stats.Record(ctx, statRecentStoragePendingsAdded.M(int64(len(objContext.Requests))))
}

// GetRequests returns a deep-copy of requests collections
func (r *PendingStorageConcrete) GetRequests() map[core.RecordID]PendingObjectContext {
	r.lock.RLock()
	defer r.lock.RUnlock()

	requestsClone := map[core.RecordID]PendingObjectContext{}
	for objID, objContext := range r.requests {
		objContext.lock.RLock()

		objectClone := PendingObjectContext{
			Active:   objContext.Context.Active,
			Requests: []core.RecordID{},
		}
		objectClone.Requests = append(objectClone.Requests, objContext.Context.Requests...)
		requestsClone[objID] = objectClone

		objContext.lock.RUnlock()
	}

	return requestsClone
}

// GetRequestsForObject returns a deep-copy of requests collection for a specific object
func (r *PendingStorageConcrete) GetRequestsForObject(obj core.RecordID) []core.RecordID {
	r.lock.RLock()
	defer r.lock.RUnlock()

	forObject, ok := r.requests[obj]
	if !ok {
		return nil
	}

	forObject.lock.RLock()
	defer forObject.lock.RUnlock()

	results := make([]core.RecordID, 0, len(forObject.Context.Requests))
	results = append(results, forObject.Context.Requests...)

	return results
}

// RemovePendingRequest removes a request on object from cache
func (r *PendingStorageConcrete) RemovePendingRequest(ctx context.Context, obj, req core.RecordID) {
	r.lock.RLock()
	defer r.lock.RUnlock()

	objContext, ok := r.requests[obj]
	if !ok {
		return
	}

	objContext.lock.Lock()
	defer objContext.lock.Unlock()

	if len(objContext.Context.Requests) == 0 {
		return
	}

	firstRequest := objContext.Context.Requests[0]
	if firstRequest.Pulse() == req.Pulse() {
		objContext.Context.Active = true
	}

	index := -1
	for internalIndex, objReq := range objContext.Context.Requests {
		if objReq == req {
			index = internalIndex
			break
		}
	}

	if index == -1 {
		return
	}

	if len(objContext.Context.Requests) == 1 {
		objContext.Context.Requests = []core.RecordID{}
		return
	}

	objContext.Context.Requests = append(objContext.Context.Requests[:index], objContext.Context.Requests[index+1:]...)

	ctx = insmetrics.InsertTag(ctx, tagJet, r.jetID.DebugString())
	stats.Record(ctx, statRecentStoragePendingsRemoved.M(1))
}
