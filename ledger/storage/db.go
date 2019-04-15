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

package storage

import (
	"context"
	"encoding/hex"
	"path/filepath"
	"sync"

	"github.com/dgraph-io/badger"
	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/insolar"
	"github.com/pkg/errors"
)

const (
	scopeIDSystem byte = 5

	sysGenesis                byte = 1
	sysHeavyClientState       byte = 3
	sysLastSyncedPulseOnHeavy byte = 4
)

//go:generate minimock -i github.com/insolar/insolar/ledger/storage.DBContext -o ./ -s _mock.go

// DBContext provides base db methods
type DBContext interface {
	BeginTransaction(update bool) (*TransactionManager, error)
	View(ctx context.Context, fn func(*TransactionManager) error) error
	Update(ctx context.Context, fn func(*TransactionManager) error) error

	StoreKeyValues(ctx context.Context, kvs []insolar.KV) error

	GetBadgerDB() *badger.DB

	Close() error

	Set(ctx context.Context, key, value []byte) error
	Get(ctx context.Context, key []byte) ([]byte, error)

	WaitingFlight()

	iterate(ctx context.Context,
		prefix []byte,
		handler func(k, v []byte) error,
	) error
}

// DB represents BadgerDB storage implementation.
type DB struct {
	PlatformCryptographyScheme insolar.PlatformCryptographyScheme `inject:""`

	db *badger.DB

	// dropLock protects dropWG from concurrent calls to Add and Wait
	dropLock sync.Mutex
	// dropWG guards inflight updates before jet drop calculated.
	dropWG sync.WaitGroup

	// for BadgerDB it is normal to have transaction conflicts
	// and these conflicts we should resolve by ourself
	// so txretiries is our knob to tune up retry logic.
	txretiries int

	jetHeavyClientLocker IDLocker

	closeLock sync.RWMutex
	isClosed  bool
}

func setOptions(o *badger.Options) *badger.Options {
	newo := &badger.Options{}
	if o != nil {
		*newo = *o
	} else {
		*newo = badger.DefaultOptions
	}
	return newo
}

// NewDB returns storage.DB with BadgerDB instance initialized by opts.
// Creates database in provided dir or in current directory if dir parameter is empty.
func NewDB(conf configuration.Ledger, opts *badger.Options) (DBContext, error) {
	opts = setOptions(opts)
	dir, err := filepath.Abs(conf.Storage.DataDirectory)
	if err != nil {
		return nil, err
	}

	opts.Dir = dir
	opts.ValueDir = dir

	bdb, err := badger.Open(*opts)
	if err != nil {
		return nil, errors.Wrap(err, "local database open failed")
	}

	db := &DB{
		db:                   bdb,
		txretiries:           conf.Storage.TxRetriesOnConflict,
		jetHeavyClientLocker: NewIDLocker(),
	}
	return db, nil
}

// Close wraps BadgerDB Close method.
//
// From https://godoc.org/github.com/dgraph-io/badger#DB.Close:
// «It's crucial to call it to ensure all the pending updates make their way to disk.
// Calling DB.Close() multiple times is not safe and wouldcause panic.»
func (db *DB) Close() error {
	db.closeLock.Lock()
	defer db.closeLock.Unlock()
	if db.isClosed {
		return ErrClosed
	}
	db.isClosed = true

	return db.db.Close()
}

// Stop stops DB component.
func (db *DB) Stop(ctx context.Context) error {
	return db.Close()
}

// BeginTransaction opens a new transaction.
// All methods called on returned transaction manager will persist changes
// only after success on "Commit" call.
func (db *DB) BeginTransaction(update bool) (*TransactionManager, error) {
	db.closeLock.RLock()
	defer db.closeLock.RUnlock()
	if db.isClosed {
		return nil, ErrClosed
	}

	if update {
		db.dropLock.Lock()
		db.dropWG.Add(1)
		db.dropLock.Unlock()
	}
	return &TransactionManager{
		db:        db,
		update:    update,
		txupdates: make(map[string]keyval),
	}, nil
}

// View accepts transaction function. All calls to received transaction manager will be consistent.
func (db *DB) View(ctx context.Context, fn func(*TransactionManager) error) error {
	tx, err := db.BeginTransaction(false)
	if err != nil {
		return err
	}
	defer tx.Discard()
	return fn(tx)
}

// Update accepts transaction function and commits changes. All calls to received transaction manager will be
// consistent and written tp disk or an error will be returned.
func (db *DB) Update(ctx context.Context, fn func(*TransactionManager) error) error {
	tries := db.txretiries
	var tx *TransactionManager
	var err error
	for {
		tx, err = db.BeginTransaction(true)
		if err != nil {
			return err
		}
		err = fn(tx)
		if err != nil {
			break
		}
		err = tx.Commit()
		if err == nil {
			break
		}
		if err != badger.ErrConflict {
			break
		}
		if tries < 1 {
			if db.txretiries > 0 {
				err = ErrConflictRetriesOver
			} else {
				err = ErrConflict
			}
			break
		}
		tries--
		tx.Discard()
	}
	tx.Discard()

	return err
}

// GetBadgerDB return badger.DB instance (for internal usage, like tests)
func (db *DB) GetBadgerDB() *badger.DB {
	return db.db
}

// StoreKeyValues stores provided key/value pairs.
func (db *DB) StoreKeyValues(ctx context.Context, kvs []insolar.KV) error {
	return db.Update(ctx, func(tx *TransactionManager) error {
		for _, rec := range kvs {
			err := tx.set(ctx, rec.K, rec.V)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (db *DB) GetPlatformCryptographyScheme() insolar.PlatformCryptographyScheme {
	return db.PlatformCryptographyScheme
}

// get wraps matching transaction manager method.
func (db *DB) Get(ctx context.Context, key []byte) ([]byte, error) {
	tx, err := db.BeginTransaction(false)
	if err != nil {
		return nil, err
	}
	defer tx.Discard()
	return tx.get(ctx, key)
}

// set wraps matching transaction manager method.
func (db *DB) Set(ctx context.Context, key, value []byte) error {
	return db.Update(ctx, func(tx *TransactionManager) error {
		return tx.set(ctx, key, value)
	})
}

func (db *DB) WaitingFlight() {
	db.dropLock.Lock()
	db.dropWG.Wait()
	db.dropLock.Unlock()
}

func (db *DB) iterate(
	ctx context.Context,
	prefix []byte,
	handler func(k, v []byte) error,
) error {
	db.closeLock.RLock()
	defer db.closeLock.RUnlock()
	if db.isClosed {
		return ErrClosed
	}

	return db.db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()

		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			key := it.Item().KeyCopy(nil)[len(prefix):]
			value, err := it.Item().ValueCopy(nil)
			if err != nil {
				return err
			}
			err = handler(key, value)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

// Key type for wrapping storage binary key.
type Key []byte

// PulseNumber returns pulse number for provided storage binary key.
func (b Key) PulseNumber() insolar.PulseNumber {
	// by default expect jetID after:
	// offset in this case: is 1 + RecordHashSize (jet length) - 1 minus jet prefix
	from := insolar.RecordHashSize
	switch b[0] {
	case scopeIDSystem:
		// for specific system records is different rules
		// pulse number could exist or not
		return 0
	}
	return pulseNumFromKey(from, b)
}

// String string hex representation
func (b Key) String() string {
	return hex.EncodeToString(b)
}
