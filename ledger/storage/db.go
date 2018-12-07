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
	"context"
	"path/filepath"
	"sync"

	"github.com/dgraph-io/badger"
	"github.com/pkg/errors"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/core/message"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/index"
	"github.com/insolar/insolar/ledger/jetdrop"
	"github.com/insolar/insolar/ledger/record"
)

const (
	scopeIDLifeline byte = 1
	scopeIDRecord   byte = 2
	scopeIDJetDrop  byte = 3
	scopeIDPulse    byte = 4
	scopeIDSystem   byte = 5
	scopeIDMessage  byte = 6
	scopeIDBlob     byte = 7
	scopeIDLocal    byte = 8

	sysGenesis                byte = 1
	sysLatestPulse            byte = 2
	sysReplicatedPulse        byte = 3
	sysLastSyncedPulseOnHeavy byte = 4
)

// DB represents BadgerDB storage implementation.
type DB struct {
	PlatformCryptographyScheme core.PlatformCryptographyScheme `inject:""`

	db         *badger.DB
	genesisRef *core.RecordRef

	// dropWG guards inflight updates before jet drop calculated.
	dropWG sync.WaitGroup

	// for BadgerDB it is normal to have transaction conflicts
	// and these conflicts we should resolve by ourself
	// so txretiries is our knob to tune up retry logic.
	txretiries int

	idlocker *IDLocker

	// NodeHistory is an in-memory active node storage for each pulse. It's required to calculate node roles
	// for past pulses to locate data.
	// It should only contain previous N pulses. It should be stored on disk.
	nodeHistory     map[core.PulseNumber][]core.Node
	nodeHistoryLock sync.Mutex
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
		db:          bdb,
		txretiries:  conf.Storage.TxRetriesOnConflict,
		idlocker:    NewIDLocker(),
		nodeHistory: map[core.PulseNumber][]core.Node{},
	}
	return db, nil
}

// Init creates initial records in storage.
func (db *DB) Init(ctx context.Context) error {
	inslog := inslogger.FromContext(ctx)
	inslog.Debug("start storage bootstrap")
	getGenesisRef := func() (*core.RecordRef, error) {
		buff, err := db.get(ctx, prefixkey(scopeIDSystem, []byte{sysGenesis}))
		if err != nil {
			return nil, err
		}
		var genesisRef core.RecordRef
		copy(genesisRef[:], buff)
		return &genesisRef, nil
	}

	createGenesisRecord := func() (*core.RecordRef, error) {
		err := db.AddPulse(
			ctx,
			core.Pulse{
				PulseNumber: core.GenesisPulse.PulseNumber,
				Entropy:     core.GenesisPulse.Entropy,
			},
		)
		if err != nil {
			return nil, err
		}
		// It should be 0. Becase pulse after 65537 will try to use a hash of drop between 0 - 65537
		err = db.SetDrop(ctx, &jetdrop.JetDrop{})
		if err != nil {
			return nil, err
		}

		lastPulse, err := db.GetLatestPulse(ctx)
		if err != nil {
			return nil, err
		}
		genesisID, err := db.SetRecord(ctx, lastPulse.Pulse.PulseNumber, &record.GenesisRecord{})
		if err != nil {
			return nil, err
		}
		err = db.SetObjectIndex(
			ctx,
			genesisID,
			&index.ObjectLifeline{LatestState: genesisID, LatestStateApproved: genesisID},
		)
		if err != nil {
			return nil, err
		}

		genesisRef := core.NewRecordRef(*genesisID, *genesisID)
		return genesisRef, db.set(ctx, prefixkey(scopeIDSystem, []byte{sysGenesis}), genesisRef[:])
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

// Stop stops DB component.
func (db *DB) Stop(ctx context.Context) error {
	return db.Close()
}

// GetBlob returns binary value stored by record ID.
func (db *DB) GetBlob(ctx context.Context, id *core.RecordID) ([]byte, error) {
	var (
		blob []byte
		err  error
	)

	err = db.View(ctx, func(tx *TransactionManager) error {
		blob, err = tx.GetBlob(ctx, id)
		return err
	})
	if err != nil {
		return nil, err
	}
	return blob, nil
}

// SetBlob saves binary value for provided pulse.
func (db *DB) SetBlob(ctx context.Context, pulseNumber core.PulseNumber, blob []byte) (*core.RecordID, error) {
	var (
		id  *core.RecordID
		err error
	)
	err = db.Update(ctx, func(tx *TransactionManager) error {
		id, err = tx.SetBlob(ctx, pulseNumber, blob)
		return err
	})
	if err != nil {
		return nil, err
	}
	return id, nil
}

// GetRecord wraps matching transaction manager method.
func (db *DB) GetRecord(ctx context.Context, id *core.RecordID) (record.Record, error) {
	var (
		fetchedRecord record.Record
		err           error
	)

	err = db.View(ctx, func(tx *TransactionManager) error {
		fetchedRecord, err = tx.GetRecord(ctx, id)
		return err
	})
	if err != nil {
		return nil, err
	}
	return fetchedRecord, nil
}

// SetRecord wraps matching transaction manager method.
func (db *DB) SetRecord(ctx context.Context, pulseNumber core.PulseNumber, rec record.Record) (*core.RecordID, error) {
	var (
		id  *core.RecordID
		err error
	)
	err = db.Update(ctx, func(tx *TransactionManager) error {
		id, err = tx.SetRecord(ctx, pulseNumber, rec)
		return err
	})
	if err != nil {
		return nil, err
	}
	return id, nil
}

// GetObjectIndex wraps matching transaction manager method.
func (db *DB) GetObjectIndex(
	ctx context.Context,
	id *core.RecordID,
	forupdate bool,
) (*index.ObjectLifeline, error) {
	tx := db.BeginTransaction(false)
	defer tx.Discard()

	idx, err := tx.GetObjectIndex(ctx, id, forupdate)
	if err != nil {
		return nil, err
	}
	return idx, nil
}

// SetObjectIndex wraps matching transaction manager method.
func (db *DB) SetObjectIndex(
	ctx context.Context,
	id *core.RecordID,
	idx *index.ObjectLifeline,
) error {
	return db.Update(ctx, func(tx *TransactionManager) error {
		return tx.SetObjectIndex(ctx, id, idx)
	})
}

// RemoveObjectIndex removes an index of an object
func (db *DB) RemoveObjectIndex(ctx context.Context, ref *core.RecordID) error {
	return db.Update(ctx, func(tx *TransactionManager) error {
		return tx.RemoveObjectIndex(ctx, ref)
	})
}

// GetDrop returns jet drop for a given pulse number.
func (db *DB) GetDrop(ctx context.Context, pulse core.PulseNumber) (*jetdrop.JetDrop, error) {
	k := prefixkey(scopeIDJetDrop, pulse.Bytes())
	buf, err := db.get(ctx, k)
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
func (db *DB) CreateDrop(ctx context.Context, pulse core.PulseNumber, prevHash []byte) (
	*jetdrop.JetDrop,
	[][]byte,
	error,
) {
	var err error
	db.waitinflight()

	hw := db.PlatformCryptographyScheme.ReferenceHasher()
	_, err = hw.Write(prevHash)
	if err != nil {
		return nil, nil, err
	}

	wg := sync.WaitGroup{}
	wg.Add(2)

	var messages [][]byte
	var messagesError error

	go func() {
		messagesPrefix := make([]byte, core.PulseNumberSize+1)
		messagesPrefix[0] = scopeIDMessage
		copy(messagesPrefix[1:], pulse.Bytes())

		messagesError = db.db.View(func(txn *badger.Txn) error {
			it := txn.NewIterator(badger.DefaultIteratorOptions)
			defer it.Close()

			for it.Seek(messagesPrefix); it.ValidForPrefix(messagesPrefix); it.Next() {
				val, err := it.Item().ValueCopy(nil)
				if err != nil {
					return err
				}
				messages = append(messages, val)
			}
			return nil
		})

		wg.Done()
	}()

	var jetDropHashError error

	go func() {
		recordPrefix := make([]byte, core.PulseNumberSize+1)
		recordPrefix[0] = scopeIDRecord
		copy(recordPrefix[1:], pulse.Bytes())

		jetDropHashError = db.db.View(func(txn *badger.Txn) error {
			it := txn.NewIterator(badger.DefaultIteratorOptions)
			defer it.Close()

			for it.Seek(recordPrefix); it.ValidForPrefix(recordPrefix); it.Next() {
				val, err := it.Item().ValueCopy(nil)
				if err != nil {
					return err
				}
				hw.Sum(val)
			}
			return nil
		})

		wg.Done()
	}()

	wg.Wait()

	if messagesError != nil {
		return nil, nil, messagesError
	}
	if jetDropHashError != nil {
		return nil, nil, jetDropHashError
	}

	drop := jetdrop.JetDrop{
		Pulse:    pulse,
		PrevHash: prevHash,
		Hash:     hw.Sum(nil),
	}
	return &drop, messages, nil
}

// SetDrop saves provided JetDrop in db.
func (db *DB) SetDrop(ctx context.Context, drop *jetdrop.JetDrop) error {
	k := prefixkey(scopeIDJetDrop, drop.Pulse.Bytes())
	_, err := db.get(ctx, k)
	if err == nil {
		return ErrOverride
	}

	encoded, err := jetdrop.Encode(drop)
	if err != nil {
		return err
	}
	return db.set(ctx, k, encoded)
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
func (db *DB) View(ctx context.Context, fn func(*TransactionManager) error) error {
	tx := db.BeginTransaction(false)
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
				err = ErrConflict
			}
			break
		}
		tries--
		tx.Discard()
	}
	tx.Discard()

	if err != nil {
		inslogger.FromContext(ctx).Errorln("DB Update error:", err)
	}
	return err
}

// GetBadgerDB return badger.DB instance (for internal usage, like tests)
func (db *DB) GetBadgerDB() *badger.DB {
	return db.db
}

// SetMessage persists message to the database
func (db *DB) SetMessage(ctx context.Context, pulseNumber core.PulseNumber, genericMessage core.Message) error {
	messageBytes := message.ToBytes(genericMessage)
	hw := db.PlatformCryptographyScheme.ReferenceHasher()
	_, err := hw.Write(messageBytes)
	if err != nil {
		return err
	}
	hw.Sum(nil)

	return db.set(
		ctx,
		prefixkey(scopeIDMessage, bytes.Join([][]byte{pulseNumber.Bytes(), hw.Sum(nil)}, nil)),
		messageBytes,
	)
}

// SetLocalData saves provided data to storage.
func (db *DB) SetLocalData(ctx context.Context, pulse core.PulseNumber, key []byte, data []byte) error {
	return db.set(
		ctx,
		bytes.Join([][]byte{{scopeIDLocal}, pulse.Bytes(), key}, nil),
		data,
	)
}

// GetLocalData retrieves data from storage.
func (db *DB) GetLocalData(ctx context.Context, pulse core.PulseNumber, key []byte) ([]byte, error) {
	return db.get(
		ctx,
		bytes.Join([][]byte{{scopeIDLocal}, pulse.Bytes(), key}, nil),
	)
}

// IterateLocalData iterates over all record with specified prefix and calls handler with key and value of that record.
//
// The key will be returned without prefix (e.g. the remaining slice) and value will be returned as it was saved.
func (db *DB) IterateLocalData(
	ctx context.Context,
	pulse core.PulseNumber,
	prefix []byte,
	handler func(k, v []byte) error,
) error {
	fullPrefix := bytes.Join([][]byte{{scopeIDLocal}, pulse.Bytes(), prefix}, nil)
	return db.iterate(ctx, fullPrefix, handler)
}

// IterateRecords iterates over records.
func (db *DB) IterateRecords(
	ctx context.Context,
	pulse core.PulseNumber,
	handler func(id core.RecordID, rec record.Record) error,
) error {
	prefix := bytes.Join([][]byte{{scopeIDRecord}, pulse.Bytes()}, nil)

	return db.iterate(ctx, prefix, func(k, v []byte) error {
		id := core.NewRecordID(pulse, k)
		rec := record.DeserializeRecord(v)
		err := handler(*id, rec)
		if err != nil {
			return err
		}
		return nil
	})
}

// SetActiveNodes saves active nodes for pulse in memory.
func (db *DB) SetActiveNodes(pulse core.PulseNumber, nodes []core.Node) error {
	db.nodeHistoryLock.Lock()
	defer db.nodeHistoryLock.Unlock()

	if _, ok := db.nodeHistory[pulse]; ok {
		return errors.New("node history override is forbidden")
	}

	db.nodeHistory[pulse] = nodes

	return nil
}

// GetActiveNodes return active nodes for specified pulse.
func (db *DB) GetActiveNodes(pulse core.PulseNumber) ([]core.Node, error) {
	nodes, ok := db.nodeHistory[pulse]
	if !ok {
		return nil, errors.New("no nodes for this pulse")
	}

	return nodes, nil
}

// StoreKeyValues stores provided key/value pairs.
func (db *DB) StoreKeyValues(ctx context.Context, kvs []core.KV) error {
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

// get wraps matching transaction manager method.
func (db *DB) get(ctx context.Context, key []byte) ([]byte, error) {
	tx := db.BeginTransaction(false)
	defer tx.Discard()
	return tx.get(ctx, key)
}

// set wraps matching transaction manager method.
func (db *DB) set(ctx context.Context, key, value []byte) error {
	return db.Update(ctx, func(tx *TransactionManager) error {
		return tx.set(ctx, key, value)
	})
}

func (db *DB) iterate(
	ctx context.Context,
	prefix []byte,
	handler func(k, v []byte) error,
) error {
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
