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
	"log"
	"path/filepath"
	"sync"

	"github.com/dgraph-io/badger"
	"github.com/insolar/insolar/core"
	"github.com/pkg/errors"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/ledger/hash"
	"github.com/insolar/insolar/ledger/index"
	"github.com/insolar/insolar/ledger/jetdrop"
	"github.com/insolar/insolar/ledger/record"
)

const (
	scopeIDLifeline byte = 1
	scopeIDRecord   byte = 2
	scopeIDJetDrop  byte = 3
	scopeIDEntropy  byte = 4

	rootKey = "0"
)

// DB represents BadgerDB storage implementation.
type DB struct {
	db           *badger.DB
	currentPulse core.PulseNumber
	rootRef      *record.Reference

	// dropWG guards inflight updates before jet drop calculated.
	dropWG sync.WaitGroup

	// for BadgerDB it is normal to have transaction conflicts
	// and these conflicts we should resolve by ourself
	// so txretiries is our knob to tune up retry logic.
	txretiries int
}

// SetTxRetiries sets number of retries on conflict in Update
func (db *DB) SetTxRetiries(n int) {
	db.txretiries = n
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
func NewDB(conf configuration.Ledger, opts *badger.Options) (*DB, error) {
	opts = setOptions(opts)
	dir, err := filepath.Abs(conf.DataDirectory)
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
		db:         bdb,
		txretiries: conf.TxRetriesOnConflict,
	}
	return db, nil
}

// Bootstrap creates initial records in storage.
func (db *DB) Bootstrap() error {
	getRootRef := func() (*record.Reference, error) {
		rootRefBuff, err := db.Get([]byte(rootKey))
		if err != nil {
			return nil, err
		}
		var coreRootRef core.RecordRef
		copy(coreRootRef[:], rootRefBuff)
		rootRef := record.Core2Reference(coreRootRef)
		return &rootRef, nil
	}

	createRootRecord := func() (*record.Reference, error) {
		rootRef, err := db.SetRecord(&record.ObjectActivateRecord{
			ActivationRecord: record.ActivationRecord{
				StatefulResult: record.StatefulResult{
					ResultRecord: record.ResultRecord{
						RequestRecord: record.Core2Reference(core.RecordRef{}),
						DomainRecord:  record.Core2Reference(core.RecordRef{}),
					},
				},
			},
		})
		if err != nil {
			return nil, err
		}
		err = db.SetObjectIndex(rootRef, &index.ObjectLifeline{LatestStateRef: *rootRef})
		if err != nil {
			return nil, err
		}

		// TODO: temporary fake entropy
		err = db.SetEntropy(db.GetCurrentPulse(), core.Entropy{})
		if err != nil {
			return nil, err
		}

		return rootRef, db.Set([]byte(rootKey), rootRef.CoreRef()[:])
	}

	var err error
	db.rootRef, err = getRootRef()
	if err == ErrNotFound {
		db.rootRef, err = createRootRecord()
	}
	if err != nil {
		return errors.Wrap(err, "bootstrap failed")
	}

	return nil
}

// RootRef returns the root record reference.
//
// Root record is the parent for all top-level records.
func (db *DB) RootRef() *record.Reference {
	return db.rootRef
}

// Close wraps BadgerDB Close method.
//
// From https://godoc.org/github.com/dgraph-io/badger#DB.Close:
// «It's crucial to call it to ensure all the pending updates make their way to disk.
// Calling DB.Close() multiple times is not safe and wouldcause panic.»
func (db *DB) Close() error {
	// TODO: add close flag and mutex guard on Close method
	return db.db.Close()
}

// Get wraps matching transaction manager method.
func (db *DB) Get(key []byte) ([]byte, error) {
	tx := db.BeginTransaction(false)
	defer tx.Discard()
	return tx.Get(key)
}

// Set wraps matching transaction manager method.
func (db *DB) Set(key, value []byte) error {
	return db.Update(func(tx *TransactionManager) error {
		return tx.Set(key, value)
	})
}

// GetRecord wraps matching transaction manager method.
func (db *DB) GetRecord(ref *record.Reference) (record.Record, error) {
	tx := db.BeginTransaction(false)
	defer tx.Discard()
	rec, err := tx.GetRecord(ref)
	if err != nil {
		return nil, err
	}
	return rec, nil
}

// SetRecord wraps matching transaction manager method.
func (db *DB) SetRecord(rec record.Record) (ref *record.Reference, err error) {
	err = db.Update(func(tx *TransactionManager) error {
		ref, err = tx.SetRecord(rec)
		return err
	})
	if err != nil {
		ref = nil
	}
	return ref, err
}

// GetClassIndex wraps matching transaction manager method.
func (db *DB) GetClassIndex(ref *record.Reference) (*index.ClassLifeline, error) {
	tx := db.BeginTransaction(false)
	defer tx.Discard()

	idx, err := tx.GetClassIndex(ref)
	if err != nil {
		return nil, err
	}
	return idx, nil
}

// SetClassIndex wraps matching transaction manager method.
func (db *DB) SetClassIndex(ref *record.Reference, idx *index.ClassLifeline) error {
	return db.Update(func(tx *TransactionManager) error {
		return tx.SetClassIndex(ref, idx)
	})
}

// GetObjectIndex wraps matching transaction manager method.
func (db *DB) GetObjectIndex(ref *record.Reference) (*index.ObjectLifeline, error) {
	tx := db.BeginTransaction(false)
	defer tx.Discard()

	idx, err := tx.GetObjectIndex(ref)
	if err != nil {
		return nil, err
	}
	return idx, nil
}

// SetObjectIndex wraps matching transaction manager method.
func (db *DB) SetObjectIndex(ref *record.Reference, idx *index.ObjectLifeline) error {
	return db.Update(func(tx *TransactionManager) error {
		return tx.SetObjectIndex(ref, idx)
	})
}

// GetDrop returns jet drop for a given pulse number.
func (db *DB) GetDrop(pulse core.PulseNumber) (*jetdrop.JetDrop, error) {
	k := prefixkey(scopeIDJetDrop, pulse.Bytes())
	buf, err := db.Get(k)
	if err != nil {
		return nil, err
	}
	drop, err := jetdrop.Decode(buf)
	if err != nil {
		return nil, err
	}
	return drop, nil
}

func (db *DB) waitinflight() {
	db.dropWG.Wait()
}

// SetDrop stores jet drop for given pulse number.
//
// Previous JetDrop should be provided. On success returns saved drop hash.
func (db *DB) SetDrop(pulse core.PulseNumber, prevdrop *jetdrop.JetDrop) (*jetdrop.JetDrop, error) {
	db.waitinflight()

	hw := hash.NewSHA3()
	err := db.ProcessSlotHashes(pulse, func(it HashIterator) error {
		for i := 1; it.Next(); i++ {
			b := it.ShallowHash()
			_, err := hw.Write(b)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	drophash := hw.Sum(nil)

	drop := &jetdrop.JetDrop{
		Pulse:    pulse,
		PrevHash: prevdrop.Hash,
		Hash:     drophash,
	}
	encoded, err := jetdrop.Encode(drop)
	if err != nil {
		return nil, err
	}

	k := prefixkey(scopeIDJetDrop, pulse.Bytes())
	err = db.Set(k, encoded)
	if err != nil {
		drop = nil
	}
	return drop, err
}

// GetEntropy wraps matching transaction manager method.
func (db *DB) GetEntropy(pulse core.PulseNumber) (*core.Entropy, error) {
	tx := db.BeginTransaction(false)
	defer tx.Discard()

	idx, err := tx.GetEntropy(pulse)
	if err != nil {
		return nil, err
	}
	return idx, nil
}

// SetEntropy wraps matching transaction manager method.
func (db *DB) SetEntropy(pulse core.PulseNumber, entropy core.Entropy) error {
	return db.Update(func(tx *TransactionManager) error {
		return tx.SetEntropy(pulse, entropy)
	})
}

// SetCurrentPulse sets current pulse number.
func (db *DB) SetCurrentPulse(pulse core.PulseNumber) {
	db.currentPulse = pulse
}

// GetCurrentPulse returns current pulse number.
func (db *DB) GetCurrentPulse() core.PulseNumber {
	return db.currentPulse
}

// BeginTransaction opens a new transaction.
// All methods called on returned transaction manager will persist changes
// only after success on "Commit" call.
func (db *DB) BeginTransaction(update bool) *TransactionManager {
	if update {
		db.dropWG.Add(1)
	}
	return &TransactionManager{
		db:     db,
		txn:    db.db.NewTransaction(update),
		update: update,
	}
}

// View accepts transaction function. All calls to received transaction manager will be consistent.
func (db *DB) View(fn func(*TransactionManager) error) error {
	tx := db.BeginTransaction(false)
	defer tx.Discard()
	return fn(tx)
}

// Update accepts transaction function and commits changes. All calls to received transaction manager will be
// consistent and written tp disk or an error will be returned.
func (db *DB) Update(fn func(*TransactionManager) error) error {
	tries := db.txretiries
	var tx *TransactionManager
	var err error
	for {
		tx = db.BeginTransaction(true)
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
				log.Println(">>> ErrConflict:", ErrConflict)
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
