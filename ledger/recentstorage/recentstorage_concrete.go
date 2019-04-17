//
// Copyright 2019 Insolar Technologies GmbH
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package recentstorage

import (
	"context"
	"sync"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/instrumentation/insmetrics"
	"go.opencensus.io/stats"
)

// RecentStorageProvider provides a recent storage for jet
type RecentStorageProvider struct { //nolint: golint
	pendingStorages map[insolar.ID]*PendingStorageConcrete

	indexLock   sync.Mutex
	pendingLock sync.Mutex
}

// NewRecentStorageProvider creates new provider
func NewRecentStorageProvider() *RecentStorageProvider {
	return &RecentStorageProvider{
		pendingStorages: map[insolar.ID]*PendingStorageConcrete{},
	}
}

// GetPendingStorage returns pendings for a specific jet
func (p *RecentStorageProvider) GetPendingStorage(ctx context.Context, jetID insolar.ID) PendingStorage {
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

// ClonePendingStorage clones pending requests from one jet to another one
func (p *RecentStorageProvider) ClonePendingStorage(ctx context.Context, fromJetID, toJetID insolar.ID) {
	p.pendingLock.Lock()
	fromStorage, ok := p.pendingStorages[fromJetID]
	p.pendingLock.Unlock()

	if !ok {
		return
	}

	fromStorage.lock.RLock()
	toStorage := &PendingStorageConcrete{
		jetID:    toJetID,
		requests: map[insolar.ID]*lockedPendingObjectContext{},
	}
	for objID, pendingContext := range fromStorage.requests {
		if len(pendingContext.Context.Requests) == 0 {
			continue
		}

		pendingContext.lock.Lock()

		clone := PendingObjectContext{
			Active:   pendingContext.Context.Active,
			Requests: []insolar.ID{},
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

// RemovePendingStorage removes pending requests for a specific jet from provider
// If there is a reference to RecentIndexStorage somewhere, it won't be affected
func (p *RecentStorageProvider) RemovePendingStorage(ctx context.Context, id insolar.ID) {
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

// PendingStorageConcrete contains indexes of unclosed requests (pendings) for a specific object id
type PendingStorageConcrete struct {
	lock sync.RWMutex

	jetID insolar.ID

	requests map[insolar.ID]*lockedPendingObjectContext
}

// PendingObjectContext contains a list of requests for an object
// Also it contains a boolean-flag for determination object's status
// If It's false, current LME, when it gets a hot-data, needs to send
// notifications about forgotten requests
type PendingObjectContext struct {
	Active   bool
	Requests []insolar.ID
}

type lockedPendingObjectContext struct {
	Context *PendingObjectContext
	lock    sync.RWMutex
}

// NewPendingStorage creates *PendingStorage
func NewPendingStorage(jetID insolar.ID) *PendingStorageConcrete {
	return &PendingStorageConcrete{
		jetID:    jetID,
		requests: map[insolar.ID]*lockedPendingObjectContext{},
	}
}

// AddPendingRequest adds an id of pending request to memory
// The id stores in a collection ids of a specific object
func (r *PendingStorageConcrete) AddPendingRequest(ctx context.Context, obj, req insolar.ID) {
	r.lock.Lock()
	defer r.lock.Unlock()

	var objectContext *lockedPendingObjectContext
	var ok bool
	if objectContext, ok = r.requests[obj]; !ok {
		objectContext = &lockedPendingObjectContext{
			lock: sync.RWMutex{},
			Context: &PendingObjectContext{
				Active:   true,
				Requests: []insolar.ID{},
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
func (r *PendingStorageConcrete) SetContextToObject(ctx context.Context, obj insolar.ID, objContext PendingObjectContext) {
	r.lock.Lock()
	defer r.lock.Unlock()

	r.requests[obj] = &lockedPendingObjectContext{
		Context: &objContext,
	}
	ctx = insmetrics.InsertTag(ctx, tagJet, r.jetID.DebugString())
	stats.Record(ctx, statRecentStoragePendingsAdded.M(int64(len(objContext.Requests))))
}

// GetRequests returns a deep-copy of requests collections
func (r *PendingStorageConcrete) GetRequests() map[insolar.ID]PendingObjectContext {
	r.lock.RLock()
	defer r.lock.RUnlock()

	requestsClone := map[insolar.ID]PendingObjectContext{}
	for objID, objContext := range r.requests {
		objContext.lock.RLock()

		objectClone := PendingObjectContext{
			Active:   objContext.Context.Active,
			Requests: []insolar.ID{},
		}
		objectClone.Requests = append(objectClone.Requests, objContext.Context.Requests...)
		requestsClone[objID] = objectClone

		objContext.lock.RUnlock()
	}

	return requestsClone
}

// GetRequestsForObject returns a deep-copy of requests collection for a specific object
func (r *PendingStorageConcrete) GetRequestsForObject(obj insolar.ID) []insolar.ID {
	r.lock.RLock()
	defer r.lock.RUnlock()

	forObject, ok := r.requests[obj]
	if !ok {
		return nil
	}

	forObject.lock.RLock()
	defer forObject.lock.RUnlock()

	results := make([]insolar.ID, 0, len(forObject.Context.Requests))
	results = append(results, forObject.Context.Requests...)

	return results
}

// RemovePendingRequest removes a request on object from cache
func (r *PendingStorageConcrete) RemovePendingRequest(ctx context.Context, obj, req insolar.ID) {
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
		objContext.Context.Requests = []insolar.ID{}
		return
	}

	objContext.Context.Requests = append(objContext.Context.Requests[:index], objContext.Context.Requests[index+1:]...)

	ctx = insmetrics.InsertTag(ctx, tagJet, r.jetID.DebugString())
	stats.Record(ctx, statRecentStoragePendingsRemoved.M(1))
}
