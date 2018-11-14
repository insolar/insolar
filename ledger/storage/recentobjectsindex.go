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
func NewRecentObjectsIndex(defaultTTL int) *RecentObjectsIndex {
	return &RecentObjectsIndex{
		recentObjects: map[string]*RecentObjectsIndexMeta{},
		DefaultTTL:    defaultTTL,
		lock:          sync.Mutex{},
	}
}

// AddID adds object to cache
func (r *RecentObjectsIndex) AddID(id *core.RecordID) {
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
