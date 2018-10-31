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
	"github.com/insolar/insolar/core/message"
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
	scopeIDMessage  byte = 6
	scopeIDBlob     byte = 7

	sysGenesis     byte = 1
	sysLatestPulse byte = 2
)

// DB represents BadgerDB storage implementation.
type DB struct {
	db         *badger.DB
	genesisRef *core.RecordRef

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
	getGenesisRef := func() (*core.RecordRef, error) {
		buff, err := db.Get(prefixkey(scopeIDSystem, []byte{sysGenesis}))
		if err != nil {
			return nil, err
		}
		var genesisRef core.RecordRef
		copy(genesisRef[:], buff)
		return &genesisRef, nil
	}

	createGenesisRecord := func() (*core.RecordRef, error) {
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

		lastPulse, err := db.GetLatestPulseNumber()
		if err != nil {
			return nil, err
		}
		genesisID, err := db.SetRecord(lastPulse, &record.GenesisRecord{})
		if err != nil {
			return nil, err
		}
		err = db.SetObjectIndex(
			genesisID,
			&index.ObjectLifeline{LatestState: genesisID, LatestStateApproved: genesisID},
		)
		if err != nil {
			return nil, err
		}

		genesisRef := core.NewRecordRef(*genesisID, *genesisID)
		return genesisRef, db.Set(prefixkey(scopeIDSystem, []byte{sysGenesis}), genesisRef[:])
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
func (db *DB) GenesisRef() *core.RecordRef {
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

func (db *DB) GetBlob(id *core.RecordID) ([]byte, error) {
	var (
		blob []byte
		err  error
	)

	err = db.View(func(tx *TransactionManager) error {
		blob, err = tx.GetBlob(id)
		return err
	})
	if err != nil {
		return nil, err
	}
	return blob, nil
}

func (db *DB) SetBlob(pulseNumber core.PulseNumber, blob []byte) (*core.RecordID, error) {
	var (
		id  *core.RecordID
		err error
	)
	err = db.Update(func(tx *TransactionManager) error {
		id, err = tx.SetBlob(pulseNumber, blob)
		return err
	})
	if err != nil {
		return nil, err
	}
	return id, nil
}

// GetRequest wraps matching transaction manager method.
func (db *DB) GetRequest(id *core.RecordID) (record.Request, error) {
	tx := db.BeginTransaction(false)
	defer tx.Discard()
	return tx.GetRequest(id)
}

// GetRecord wraps matching transaction manager method.
func (db *DB) GetRecord(id *core.RecordID) (record.Record, error) {
	var (
		fetchedRecord record.Record
		err           error
	)

	err = db.View(func(tx *TransactionManager) error {
		fetchedRecord, err = tx.GetRecord(id)
		return err
	})
	if err != nil {
		return nil, err
	}
	return fetchedRecord, nil
}

// SetRecord wraps matching transaction manager method.
func (db *DB) SetRecord(pulseNumber core.PulseNumber, rec record.Record) (*core.RecordID, error) {
	var (
		id  *core.RecordID
		err error
	)
	err = db.Update(func(tx *TransactionManager) error {
		id, err = tx.SetRecord(pulseNumber, rec)
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
func (db *DB) GetClassIndex(id *core.RecordID, forupdate bool) (*index.ClassLifeline, error) {
	tx := db.BeginTransaction(forupdate)
	defer tx.Discard()

	idx, err := tx.GetClassIndex(id, false)
	if err != nil {
		return nil, err
	}
	return idx, nil
}

// SetClassIndex wraps matching transaction manager method.
func (db *DB) SetClassIndex(id *core.RecordID, idx *index.ClassLifeline) error {
	return db.Update(func(tx *TransactionManager) error {
		return tx.SetClassIndex(id, idx)
	})
}

// GetObjectIndex wraps matching transaction manager method.
func (db *DB) GetObjectIndex(id *core.RecordID, forupdate bool) (*index.ObjectLifeline, error) {
	tx := db.BeginTransaction(false)
	defer tx.Discard()

	idx, err := tx.GetObjectIndex(id, forupdate)
	if err != nil {
		return nil, err
	}
	return idx, nil
}

// SetObjectIndex wraps matching transaction manager method.
func (db *DB) SetObjectIndex(id *core.RecordID, idx *index.ObjectLifeline) error {
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
func (db *DB) CreateDrop(pulse core.PulseNumber, prevHash []byte) (
	*jetdrop.JetDrop,
	[][]byte,
	error,
) {
	var err error
	db.waitinflight()

	hw := hash.NewIDHash()
	_, err = hw.Write(prevHash)
	if err != nil {
		return nil, nil, err
	}

	prefix := make([]byte, core.PulseNumberSize+1)
	prefix[0] = scopeIDMessage
	copy(prefix[1:], pulse.Bytes())

	var messages [][]byte
	err = db.db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()

		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			val, err := it.Item().ValueCopy(nil)
			if err != nil {
				return err
			}
			messages = append(messages, val)
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
	return &drop, messages, nil
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

// GetBadgerDB return badger.DB instance (for internal usage, like tests)
func (db *DB) GetBadgerDB() *badger.DB {
	return db.db
}

// SetMessage persists message to the database
func (db *DB) SetMessage(pulseNumber core.PulseNumber, genericMessage core.Message) error {
	messageBytes, err := message.ToBytes(genericMessage)
	if err != nil {
		return err
	}

	hw := hash.NewIDHash()
	_, err = hw.Write(messageBytes)
	if err != nil {
		return err
	}
	hw.Sum(nil)

	return db.Set(
		prefixkey(scopeIDMessage, bytes.Join([][]byte{pulseNumber.Bytes(), hw.Sum(nil)}, nil)),
		messageBytes,
	)
}
