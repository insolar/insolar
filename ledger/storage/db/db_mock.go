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

package db

import (
	"sync"
)

// MockDB is a mock DB implementation. It can be used as a stub for other implementations in component tests.
type MockDB struct {
	lock    sync.RWMutex
	backend map[string][]byte
}

// NewMemoryMockDB creates new mock DB instance.
func NewMemoryMockDB() *MockDB {
	db := &MockDB{
		backend: map[string][]byte{},
	}
	return db
}

// Get returns a copy of the value for specified key from memory.
func (b *MockDB) Get(key Key) (value []byte, err error) {
	fullKey := append(key.Scope().Bytes(), key.ID()...)

	b.lock.RLock()
	defer b.lock.RUnlock()
	value, ok := b.backend[string(fullKey)]
	if !ok {
		return nil, ErrNotFound
	}
	return append([]byte{}, value...), nil
}

// Set stores value for a key in memory storage.
func (b *MockDB) Set(key Key, value []byte) error {
	fullKey := append(key.Scope().Bytes(), key.ID()...)
	b.lock.Lock()
	defer b.lock.Unlock()
	b.backend[string(fullKey)] = append([]byte{}, value...)
	return nil
}
