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
	"path/filepath"
	"sync"

	"github.com/dgraph-io/badger"
	"github.com/pkg/errors"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/cryptohelpers/hash"
	"github.com/insolar/insolar/ledger/index"
	"github.com/insolar/insolar/ledger/jetdrop"
	"github.com/insolar/insolar/ledger/record"
	"github.com/insolar/insolar/log"
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
	genesisRef   *record.Reference

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
		db:         bdb,
		txretiries: conf.Storage.TxRetriesOnConflict,
	}
	return db, nil
}

// Bootstrap creates initial records in storage.
func (db *DB) Bootstrap() error {
	getGenesisRef := func() (*record.Reference, error) {
		rootRefBuff, err := db.Get([]byte(rootKey))
		if err != nil {
			return nil, err
		}
		var coreRootRef core.RecordRef
		copy(coreRootRef[:], rootRefBuff)
		rootRef := record.Core2Reference(coreRootRef)
		return &rootRef, nil
	}

	createGenesisRecord := func() (*record.Reference, error) {
		genesisID, err := db.SetRecord(&record.GenesisRecord{})
		if err != nil {
			return nil, err
		}
		err = db.SetObjectIndex(genesisID, &index.ObjectLifeline{LatestState: *genesisID})
		if err != nil {
			return nil, err
		}

		db.SetCurrentPulse(core.GenesisPulse.PulseNumber)
		err = db.SetEntropy(core.GenesisPulse.PulseNumber, core.GenesisPulse.Entropy)
		if err != nil {
			return nil, err
		}
		_, err = db.SetDrop(core.GenesisPulse.PulseNumber, &jetdrop.JetDrop{})
		if err != nil {
			return nil, err
		}

		genesisRef := record.Reference{Domain: *genesisID, Record: *genesisID}
		return &genesisRef, db.Set([]byte(rootKey), genesisRef.CoreRef()[:])
	}

	var err error
	db.genesisRef, err = getGenesisRef()
	if err == ErrNotFound {
		db.genesisRef, err = createGenesisRecord()
	}
	if err != nil {
		return errors.Wrap(err, "bootstrap failed")
	}

	return nil
}

// GenesisRef returns the genesis record reference.
//
// Genesis record is the parent for all top-level records.
func (db *DB) GenesisRef() *record.Reference {
	return db.genesisRef
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

// GetRequest wraps matching transaction manager method.
func (db *DB) GetRequest(id *record.ID) (record.Request, error) {
	tx := db.BeginTransaction(false)
	defer tx.Discard()
	return tx.GetRequest(id)
}

// SetRequest wraps matching transaction manager method.
func (db *DB) SetRequest(req record.Request) (*record.ID, error) {
	var (
		id  *record.ID
		err error
	)
	txerr := db.Update(func(tx *TransactionManager) error {
		id, err = tx.SetRequest(req)
		return err
	})
	if txerr != nil {
		return nil, txerr
	}
	return id, nil
}

// GetRecord wraps matching transaction manager method.
func (db *DB) GetRecord(id *record.ID) (record.Record, error) {
	tx := db.BeginTransaction(false)
	defer tx.Discard()
	rec, err := tx.GetRecord(id)
	if err != nil {
		return nil, err
	}
	return rec, nil
}

// SetRecord wraps matching transaction manager method.
func (db *DB) SetRecord(rec record.Record) (*record.ID, error) {
	var (
		id  *record.ID
		err error
	)
	err = db.Update(func(tx *TransactionManager) error {
		id, err = tx.SetRecord(rec)
		return err
	})
	if err != nil {
		return nil, err
	}
	return id, nil
}

// GetClassIndex wraps matching transaction manager method.
func (db *DB) GetClassIndex(id *record.ID) (*index.ClassLifeline, error) {
	tx := db.BeginTransaction(false)
	defer tx.Discard()

	idx, err := tx.GetClassIndex(id)
	if err != nil {
		return nil, err
	}
	return idx, nil
}

// SetClassIndex wraps matching transaction manager method.
func (db *DB) SetClassIndex(id *record.ID, idx *index.ClassLifeline) error {
	return db.Update(func(tx *TransactionManager) error {
		return tx.SetClassIndex(id, idx)
	})
}

// GetObjectIndex wraps matching transaction manager method.
func (db *DB) GetObjectIndex(id *record.ID) (*index.ObjectLifeline, error) {
	tx := db.BeginTransaction(false)
	defer tx.Discard()

	idx, err := tx.GetObjectIndex(id)
	if err != nil {
		return nil, err
	}
	return idx, nil
}

// SetObjectIndex wraps matching transaction manager method.
func (db *DB) SetObjectIndex(id *record.ID, idx *index.ObjectLifeline) error {
	return db.Update(func(tx *TransactionManager) error {
		return tx.SetObjectIndex(id, idx)
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

	hw := hash.NewIDHash()
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
				log.Info("local storage transaction conflict")
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
