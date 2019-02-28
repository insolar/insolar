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

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/core"
)

type memoryDB struct {
	lock    sync.RWMutex
	backend map[string][]byte
}

func NewMemoryDB(conf configuration.Ledger) (*memoryDB, error) {
	db := &memoryDB{
		backend: map[string][]byte{},
	}
	return db, nil
}

func (b *memoryDB) Get(key Key) (value []byte, err error) {
	fullKey := append(key.Scope().Bytes(), key.Key()...)

	b.lock.RLock()
	defer b.lock.RUnlock()
	value, ok := b.backend[string(fullKey)]
	if !ok {
		return nil, core.ErrNotFound
	}
	return
}

func (b *memoryDB) Set(key Key, value []byte) error {
	fullKey := append(key.Scope().Bytes(), key.Key()...)
	b.lock.Lock()
	defer b.lock.Unlock()
	b.backend[string(fullKey)] = append([]byte{}, value...)
	return nil
}
