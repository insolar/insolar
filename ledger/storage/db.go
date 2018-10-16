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
	"bytes"
	"path/filepath"
	"sync"

	"github.com/dgraph-io/badger"
	"github.com/pkg/errors"
	"github.com/ugorji/go/codec"

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
	scopeIDPulse    byte = 4
	scopeIDSystem   byte = 5

	sysGenesis     byte = 1
	sysLatestPulse byte = 2
)

// DB represents BadgerDB storage implementation.
type DB struct {
	db         *badger.DB
	genesisRef *record.Reference

	// dropWG guards inflight updates before jet drop calculated.
	dropWG sync.WaitGroup

	// for BadgerDB it is normal to have transaction conflicts
	// and these conflicts we should resolve by ourself
	// so txretiries is our knob to tune up retry logic.
	txretiries int

	idlocker *IDLocker
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
		idlocker:   NewIDLocker(),
	}
	return db, nil
}

// Bootstrap creates initial records in storage.
func (db *DB) Bootstrap() error {
	getGenesisRef := func() (*record.Reference, error) {
		rootRefBuff, err := db.Get(prefixkey(scopeIDSystem, []byte{sysGenesis}))
		if err != nil {
			return nil, err
		}
		var coreRootRef core.RecordRef
		copy(coreRootRef[:], rootRefBuff)
		rootRef := record.Core2Reference(coreRootRef)
		return &rootRef, nil
	}

	createGenesisRecord := func() (*record.Reference, error) {
		err := db.AddPulse(core.Pulse{
			PulseNumber: core.GenesisPulse.PulseNumber,
			Entropy:     core.GenesisPulse.Entropy,
		})
		if err != nil {
			return nil, err
		}
		err = db.SetDrop(&jetdrop.JetDrop{})
		if err != nil {
			return nil, err
		}

		genesisID, err := db.SetRecord(&record.GenesisRecord{})
		if err != nil {
			return nil, err
		}
		err = db.SetObjectIndex(genesisID, &index.ObjectLifeline{LatestState: *genesisID})
		if err != nil {
			return nil, err
		}

		genesisRef := record.Reference{Domain: *genesisID, Record: *genesisID}
		return &genesisRef, db.Set(prefixkey(scopeIDSystem, []byte{sysGenesis}), genesisRef.CoreRef()[:])
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

// SetRecordBinary saves binary record for specified key.
//
// This method is used for data replication.
func (db *DB) SetRecordBinary(key, rec []byte) error {
	return db.Set(prefixkey(scopeIDRecord, key), rec)
}

// GetClassIndex wraps matching transaction manager method.
func (db *DB) GetClassIndex(id *record.ID, forupdate bool) (*index.ClassLifeline, error) {
	tx := db.BeginTransaction(forupdate)
	defer tx.Discard()

	idx, err := tx.GetClassIndex(id, false)
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
func (db *DB) GetObjectIndex(id *record.ID, forupdate bool) (*index.ObjectLifeline, error) {
	tx := db.BeginTransaction(false)
	defer tx.Discard()

	idx, err := tx.GetObjectIndex(id, forupdate)
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

// CreateDrop creates and stores jet drop for given pulse number.
//
// Previous JetDrop hash should be provided. On success returns saved drop and slot records.
func (db *DB) CreateDrop(pulse core.PulseNumber, prevHash []byte) (*jetdrop.JetDrop, [][2][]byte, error) {
	db.waitinflight()

	prefix := make([]byte, core.PulseNumberSize+1)
	prefix[0] = scopeIDRecord
	copy(prefix[1:], pulse.Bytes())

	// We need to look for the closest key that is bigger because we need to reverse iterate from the last record.
	seekFor := make([]byte, len(prefix))
	copy(seekFor, prefix)
	seekFor[len(prefix)-1]++

	hw := hash.NewIDHash()
	_, err := hw.Write(prevHash)
	if err != nil {
		return nil, nil, err
	}

	var records [][2][]byte
	err = db.db.View(func(txn *badger.Txn) error {
		ops := badger.DefaultIteratorOptions
		ops.Reverse = true
		it := txn.NewIterator(ops)
		defer it.Close()
		it.Seek(seekFor)

		for {
			if !it.Valid() {
				break
			}

			item := it.Item()
			key := item.Key()
			if !bytes.Equal(key[:core.PulseNumberSize+1], prefix) {
				break
			}

			_, err := hw.Write(key[1:])
			if err != nil {
				return err
			}
			value, err := item.ValueCopy(nil)
			if err != nil {
				return err
			}
			records = append(records, [2][]byte{key[1:], value})

			it.Next()
		}

		return nil
	})
	if err != nil {
		return nil, nil, err
	}

	drop := jetdrop.JetDrop{
		Pulse:    pulse,
		PrevHash: prevHash,
		Hash:     hw.Sum(nil),
	}
	return &drop, records, nil
}

// SetDrop saves provided JetDrop in db.
func (db *DB) SetDrop(drop *jetdrop.JetDrop) error {
	k := prefixkey(scopeIDJetDrop, drop.Pulse.Bytes())
	_, err := db.Get(k)
	if err == nil {
		return ErrOverride
	}

	encoded, err := jetdrop.Encode(drop)
	if err != nil {
		return err
	}

	return db.Set(k, encoded)
}

// AddPulse saves new pulse data and updates index.
func (db *DB) AddPulse(pulse core.Pulse) error {
	return db.Update(func(tx *TransactionManager) error {
		var latest core.PulseNumber
		latest, err := tx.GetLatestPulseNumber()
		if err != nil && err != ErrNotFound {
			return err
		}
		pulseRec := record.PulseRecord{
			PrevPulse:          latest,
			Entropy:            pulse.Entropy,
			PredictedNextPulse: pulse.NextPulseNumber,
		}
		var buf bytes.Buffer
		enc := codec.NewEncoder(&buf, &codec.CborHandle{})
		err = enc.Encode(pulseRec)
		if err != nil {
			return err
		}
		err = tx.Set(prefixkey(scopeIDPulse, pulse.PulseNumber.Bytes()), buf.Bytes())
		if err != nil {
			return err
		}
		return tx.Set(prefixkey(scopeIDSystem, []byte{sysLatestPulse}), pulse.PulseNumber.Bytes())
	})
}

// GetPulse returns pulse for provided pulse number.
func (db *DB) GetPulse(num core.PulseNumber) (*record.PulseRecord, error) {
	buf, err := db.Get(prefixkey(scopeIDPulse, num.Bytes()))
	if err != nil {
		return nil, err
	}

	dec := codec.NewDecoder(bytes.NewReader(buf), &codec.CborHandle{})
	var rec record.PulseRecord
	err = dec.Decode(&rec)
	if err != nil {
		return nil, err
	}
	return &rec, nil
}

// GetLatestPulseNumber returns current pulse number.
func (db *DB) GetLatestPulseNumber() (core.PulseNumber, error) {
	tx := db.BeginTransaction(false)
	defer tx.Discard()

	return tx.GetLatestPulseNumber()
}

// BeginTransaction opens a new transaction.
// All methods called on returned transaction manager will persist changes
// only after success on "Commit" call.
func (db *DB) BeginTransaction(update bool) *TransactionManager {
	if update {
		db.dropWG.Add(1)
	}
	return &TransactionManager{
		db:        db,
		update:    update,
		txupdates: make(map[string]keyval),
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
