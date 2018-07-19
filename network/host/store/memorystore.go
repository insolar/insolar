/*
 *    Copyright 2018 INS Ecosystem
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

package store

import (
	"sync"
	"time"
)

// memoryStore is a simple in-memory key/value store used for unit testing, and
// the CLI example.
type memoryStore struct {
	mutex        *sync.RWMutex
	data         map[string][]byte
	replicateMap map[string]time.Time
	expireMap    map[string]time.Time
}

// NewMemoryStore creates new memory store.
func NewMemoryStore() Store {
	return newMemoryStore()
}

func newMemoryStore() *memoryStore {
	return &memoryStore{
		mutex:        &sync.RWMutex{},
		data:         make(map[string][]byte),
		replicateMap: make(map[string]time.Time),
		expireMap:    make(map[string]time.Time),
	}
}

// Store will store a key/value pair for the local node with the given
// replication and expiration times.
func (ms *memoryStore) Store(key Key, data []byte, replication time.Time, expiration time.Time, publisher bool) error {
	ms.mutex.Lock()
	defer ms.mutex.Unlock()

	keyStr := key.String()

	ms.replicateMap[keyStr] = replication
	ms.expireMap[keyStr] = expiration
	ms.data[keyStr] = data
	return nil
}

// Retrieve will return the local key/value if it exists.
func (ms *memoryStore) Retrieve(key Key) ([]byte, bool) {
	ms.mutex.RLock()
	defer ms.mutex.RUnlock()

	data, found := ms.data[key.String()]
	return data, found
}

// Delete deletes a key/value pair from the memoryStore.
func (ms *memoryStore) Delete(key Key) {
	ms.mutex.Lock()
	defer ms.mutex.Unlock()

	keyStr := key.String()

	delete(ms.replicateMap, keyStr)
	delete(ms.expireMap, keyStr)
	delete(ms.data, keyStr)
}

// GetKeysReadyToReplicate should return the keys of all data to be
// replicated across the insolar. Typically all data should be
// replicated every tReplicate seconds.
func (ms *memoryStore) GetKeysReadyToReplicate() []Key {
	ms.mutex.RLock()
	defer ms.mutex.RUnlock()

	var keys []Key
	for k := range ms.data {
		if time.Now().After(ms.replicateMap[k]) {
			keys = append(keys, []byte(k))
		}
	}
	return keys
}

// ExpireKeys should expire all key/values due for expiration.
func (ms *memoryStore) ExpireKeys() {
	ms.mutex.Lock()
	defer ms.mutex.Unlock()

	for k, v := range ms.expireMap {
		if time.Now().After(v) {
			delete(ms.replicateMap, k)
			delete(ms.expireMap, k)
			delete(ms.data, k)
		}
	}
}
