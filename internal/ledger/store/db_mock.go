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

package store

import (
	"bytes"
	"sync"

	"github.com/google/btree"
	"github.com/pkg/errors"
)

// MockDB is a mock DB implementation. It can be used as a stub for other implementations in component tests.
type MockDB struct {
	lock    sync.RWMutex
	backend map[string][]byte
	keys    *btree.BTree
}

// NewMemoryMockDB creates new mock DB instance.
func NewMemoryMockDB() *MockDB {
	db := &MockDB{
		backend: map[string][]byte{},
		keys:    btree.New(2),
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
	b.keys.ReplaceOrInsert(KeyItem(fullKey))
	return nil
}

// NewIterator returns new Iterator over the memory storage.
func (b *MockDB) NewIterator(scope Scope) Iterator {
	mi := memoryIterator{scope: scope, fullPrefix: scope.Bytes(), db: b}
	b.lock.RLock()
	return &mi
}

type KeyItem []byte

func (k KeyItem) Less(than btree.Item) bool {
	return bytes.Compare(k, than.(KeyItem)) == -1
}

type memoryIterator struct {
	scope      Scope
	fullPrefix []byte
	db         *MockDB
	once       sync.Once
	items      []KeyItem
	current    int
}

func (mi *memoryIterator) Close() {
	mi.db.lock.RUnlock()
}

func (mi *memoryIterator) Seek(prefix []byte) {
	mi.fullPrefix = append(mi.scope.Bytes(), prefix...)
	mi.searchKeys()
}

func (mi *memoryIterator) Next() bool {
	firstTime := false
	mi.once.Do(func() {
		mi.searchKeys()

		firstTime = true
	})
	if firstTime {
		return mi.current >= 0 && mi.current < len(mi.items)
	}

	mi.current++
	return mi.current >= 0 && mi.current < len(mi.items)
}

func (mi *memoryIterator) Key() []byte {
	if mi.current < 0 || mi.current >= len(mi.items) {
		return nil
	}
	val := mi.items[mi.current][len(mi.scope.Bytes()):]
	return val
}

func (mi *memoryIterator) Value() ([]byte, error) {
	if mi.current < 0 || mi.current >= len(mi.items) {
		return nil, errors.New("invalid iterator")
	}
	key := mi.items[mi.current]
	value, ok := mi.db.backend[string(key)]
	if !ok {
		return nil, ErrNotFound
	}
	return append([]byte{}, value...), nil
}

func (mi *memoryIterator) searchKeys() {
	mi.items = []KeyItem{}
	mi.db.keys.AscendGreaterOrEqual(KeyItem(mi.fullPrefix), func(i btree.Item) bool {
		if !bytes.HasPrefix(i.(KeyItem), mi.fullPrefix) {
			return false
		}
		mi.items = append(mi.items, i.(KeyItem))
		return true
	})
	if len(mi.items) == 0 {
		mi.current = -2
		return
	}
	mi.current = 0
}
