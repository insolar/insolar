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
	"context"
	"path/filepath"

	"github.com/dgraph-io/badger"
	"github.com/pkg/errors"
)

// BadgerDB is a badger DB implementation.
type BadgerDB struct {
	backend *badger.DB
}

// NewBadgerDB creates new BadgerDB instance.
// Creates new badger.DB instance with provided working dir and use it as backend for BadgerDB.
func NewBadgerDB(dir string) (*BadgerDB, error) {
	dir, err := filepath.Abs(dir)
	if err != nil {
		return nil, err
	}

	ops := badger.DefaultOptions
	ops.ValueDir = dir
	ops.Dir = dir
	bdb, err := badger.Open(ops)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open badger")
	}

	return &BadgerDB{backend: bdb}, nil
}

// Get returns value for specified key or an error. A copy of a value will be returned (i.e. getting large value can be
// long).
func (b *BadgerDB) Get(key Key) (value []byte, err error) {
	fullKey := append(key.Scope().Bytes(), key.ID()...)

	err = b.backend.View(func(txn *badger.Txn) error {
		item, err := txn.Get(fullKey)
		if err != nil {
			return err
		}
		value, err = item.ValueCopy(nil)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		if err == badger.ErrKeyNotFound {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return
}

// Set stores value for a key.
func (b *BadgerDB) Set(key Key, value []byte) error {
	fullKey := append(key.Scope().Bytes(), key.ID()...)

	err := b.backend.Update(func(txn *badger.Txn) error {
		err := txn.Set(fullKey, value)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

// NewIterator returns new Iterator over the store.
func (b *BadgerDB) NewIterator(scope Scope) Iterator {
	bi := badgerIterator{scope: scope, fullPrefix: scope.Bytes()}
	bi.txn = b.backend.NewTransaction(false)
	opts := badger.DefaultIteratorOptions
	bi.it = bi.txn.NewIterator(opts)
	return &bi
}

// Stop gracefully stops all disk writes. After calling this, it's safe to kill the process without losing data.
func (b *BadgerDB) Stop(ctx context.Context) error {
	return b.backend.Close()
}

type badgerIterator struct {
	scope      Scope
	fullPrefix []byte
	txn        *badger.Txn
	it         *badger.Iterator
}

func (bi *badgerIterator) Close() {
	bi.it.Close()
	bi.txn.Discard()
}

func (bi *badgerIterator) Seek(prefix []byte) {
	bi.fullPrefix = append(bi.scope.Bytes(), prefix...)
	bi.it.Seek(bi.fullPrefix)
}

func (bi *badgerIterator) Next() bool {
	bi.it.Next()
	return bi.it.ValidForPrefix(bi.fullPrefix)
}

func (bi *badgerIterator) Key() []byte {
	return bi.it.Item().Key()[len(bi.scope.Bytes()):]
}

func (bi *badgerIterator) Value() []byte {
	value, _ := bi.it.Item().ValueCopy(nil)
	return value
}
