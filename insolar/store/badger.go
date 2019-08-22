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
	"io"
	"sync"
	"time"

	"github.com/dgraph-io/badger"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/pkg/errors"
)

// BadgerDB is a badger DB implementation.
type BadgerDB struct {
	backend *badger.DB

	stopGC  chan struct{}
	forceGC chan chan struct{}
	doneGC  chan struct{}
}

// NewBadgerDB creates new BadgerDB instance.
// Creates new badger.DB instance with provided working dir and use it as backend for BadgerDB.
func NewBadgerDB(ops badger.Options) (*BadgerDB, error) {
	bdb, err := badger.Open(ops)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open badger")
	}

	b := &BadgerDB{backend: bdb}
	b.runGC(context.Background())
	return b, nil
}

func (b *BadgerDB) Backend() *badger.DB {
	return b.backend
}

type state struct {
	mu    sync.RWMutex
	state bool
}

func (s *state) set(val bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.state = val
}

func (s *state) check() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.state
}

func (b *BadgerDB) runGC(ctx context.Context) {
	db := b.backend
	logger := inslogger.FromContext(ctx)
	ticker := time.NewTicker(5 * time.Minute)

	b.forceGC = make(chan chan struct{})
	b.stopGC = make(chan struct{})
	b.doneGC = make(chan struct{})

	do := func() {
		logger.Info("BadgerDB: values GC start")
		defer logger.Info("BadgerDB: values GC end")

		err := db.RunValueLogGC(0.7)
		if err != nil && err != badger.ErrNoRewrite {
			logger.Errorf("BadgerDB: GC failed with error: %v", err.Error())
		}
	}

	var gcWait sync.WaitGroup
	gcAsync := &state{}

	go func() {
		for {
			select {
			case done := <-b.forceGC:
				func() {
					defer close(done)
					if gcAsync.check() {
						// blocks ForceValueGC call (on done channel) until end of GC
						gcWait.Wait()
						return
					}
					do()
				}()
			case <-ticker.C:
				func() {
					if gcAsync.check() {
						return
					}
					gcAsync.set(true)
					gcWait.Add(1)
					go func() {
						do()
						gcAsync.set(false)
						gcWait.Done()
					}()
				}()
			case <-b.stopGC:
				gcWait.Wait()
				close(b.doneGC)
				return
			}
		}
	}()
}

// ForceValueGC forces badger values garbage collection.
func (b *BadgerDB) ForceValueGC(ctx context.Context) {
	fin := make(chan struct{})
	b.forceGC <- fin
	<-fin
}

// Stop gracefully stops all disk writes. After calling this, it's safe to kill the process without losing data.
func (b *BadgerDB) Stop(ctx context.Context) error {
	logger := inslogger.FromContext(ctx)
	defer logger.Info("BadgerDB: database closed")

	logger.Info("BadgerDB: wait values GC")
	close(b.stopGC)
	<-b.doneGC

	logger.Info("BadgerDB: closing database...")

	return b.backend.Close()
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
		return err
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
		return txn.Set(fullKey, value)
	})

	return err
}

// Delete deletes value for a key.
func (b *BadgerDB) Delete(key Key) error {
	fullKey := append(key.Scope().Bytes(), key.ID()...)
	err := b.backend.Update(func(txn *badger.Txn) error {
		return txn.Delete(fullKey)
	})

	return err
}

// Backup creates backup.
func (b *BadgerDB) Backup(w io.Writer, since uint64) (uint64, error) {
	return b.backend.Backup(w, since)
}

// NewReadIterator returns new Iterator over the store.
func NewReadIterator(db *badger.DB, pivot Key, reverse bool) Iterator {
	bi := badgerIterator{pivot: pivot, reverse: reverse}
	bi.txn = db.NewTransaction(false)
	opts := badger.DefaultIteratorOptions
	opts.Reverse = reverse
	bi.it = bi.txn.NewIterator(opts)
	return &bi
}

// NewIterator returns new Iterator over the store.
func (b *BadgerDB) NewIterator(pivot Key, reverse bool) Iterator {
	bi := badgerIterator{pivot: pivot, reverse: reverse}
	bi.txn = b.backend.NewTransaction(false)
	opts := badger.DefaultIteratorOptions
	opts.Reverse = reverse
	bi.it = bi.txn.NewIterator(opts)
	return &bi
}

type badgerIterator struct {
	once      sync.Once
	pivot     Key
	reverse   bool
	txn       *badger.Txn
	it        *badger.Iterator
	prevKey   []byte
	prevValue []byte
}

func (bi *badgerIterator) Close() {
	bi.it.Close()
	bi.txn.Discard()
}

func (bi *badgerIterator) Next() bool {
	scope := bi.pivot.Scope().Bytes()
	bi.once.Do(func() {
		bi.it.Seek(append(bi.pivot.Scope().Bytes(), bi.pivot.ID()...))
	})
	if !bi.it.ValidForPrefix(scope) {
		return false
	}

	prev := bi.it.Item().KeyCopy(nil)
	bi.prevKey = prev[len(scope):]
	prev, err := bi.it.Item().ValueCopy(nil)
	if err != nil {
		return false
	}
	bi.prevValue = prev

	bi.it.Next()
	return true
}

func (bi *badgerIterator) Key() []byte {
	return bi.prevKey
}

func (bi *badgerIterator) Value() ([]byte, error) {
	return bi.prevValue, nil
}
