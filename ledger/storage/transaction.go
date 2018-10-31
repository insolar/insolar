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
	"encoding/binary"

	"github.com/dgraph-io/badger"
	"github.com/insolar/insolar/cryptohelpers/hash"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/ledger/index"
	"github.com/insolar/insolar/ledger/record"
	"github.com/insolar/insolar/log"
)

type keyval struct {
	k []byte
	v []byte
}

// TransactionManager is used to ensure persistent writes to disk.
type TransactionManager struct {
	db        *DB
	update    bool
	locks     []*core.RecordID
	txupdates map[string]keyval
}

func prefixkey(prefix byte, key []byte) []byte {
	k := make([]byte, core.RecordIDSize+1)
	k[0] = prefix
	_ = copy(k[1:], key)
	return k
}

func (m *TransactionManager) lockOnID(id *core.RecordID) {
	m.db.idlocker.Lock(id)
	m.locks = append(m.locks, id)
}

func (m *TransactionManager) releaseLocks() {
	for _, id := range m.locks {
		m.db.idlocker.Unlock(id)
	}
}

// Commit tries to write transaction on disk. Returns error on fail.
func (m *TransactionManager) Commit() error {
	if len(m.txupdates) == 0 {
		return nil
	}
	var err error
	tx := m.db.db.NewTransaction(m.update)
	defer tx.Discard()
	for _, rec := range m.txupdates {
		err = tx.Set(rec.k, rec.v)
		if err != nil {
			break
		}
	}
	m.txupdates = nil
	if err != nil {
		return err
	}
	return tx.Commit(nil)
}

// Discard terminates transaction without disk writes.
func (m *TransactionManager) Discard() {
	m.releaseLocks()
	if m.update {
		m.db.dropWG.Done()
	}
}

// GetRequest returns request record from BadgerDB by *record.Reference.
//
// It returns ErrNotFound if the DB does not contain the key.
func (m *TransactionManager) GetRequest(id *core.RecordID) (record.Request, error) {
	rec, err := m.GetRecord(id)
	if err != nil {
		return nil, err
	}
	// TODO: return error if record is not a request.
	req := rec.(record.Request)
	return req, nil
}

func (m *TransactionManager) set(key, val []byte) {
	m.txupdates[string(key)] = keyval{k: key, v: val}
}

func (m *TransactionManager) GetBlob(id *core.RecordID) ([]byte, error) {
	k := prefixkey(scopeIDBlob, id[:])
	log.Debugf("GetRecord by id %+v (key=%x)", id, k)
	return m.Get(k)
}

func (m *TransactionManager) SetBlob(pulseNumber core.PulseNumber, blob []byte) (*core.RecordID, error) {
	id := record.CalculateIDForBlob(pulseNumber, blob)
	k := prefixkey(scopeIDBlob, id[:])
	geterr := m.db.db.View(func(tx *badger.Txn) error {
		_, err := tx.Get(k)
		return err
	})
	if geterr == nil {
		return id, ErrOverride
	}
	if geterr != badger.ErrKeyNotFound {
		return nil, ErrNotFound
	}

	m.set(k, blob)
	return id, nil
}

// GetRecord returns record from BadgerDB by *record.Reference.
//
// It returns ErrNotFound if the DB does not contain the key.
func (m *TransactionManager) GetRecord(id *core.RecordID) (record.Record, error) {
	k := prefixkey(scopeIDRecord, id[:])
	log.Debugf("GetRecord by id %+v (key=%x)", id, k)
	buf, err := m.Get(k)
	if err != nil {
		return nil, err
	}
	return record.DeserializeRecord(buf), nil
}

// SetRecord stores record in BadgerDB and returns *record.ID of new record.
//
// If record exists returns both *record.ID and ErrOverride error.
// If record not found returns nil and ErrNotFound error
func (m *TransactionManager) SetRecord(pulseNumber core.PulseNumber, rec record.Record) (*core.RecordID, error) {
	recHash := hash.NewIDHash()
	_, err := rec.WriteHashData(recHash)
	if err != nil {
		return nil, err
	}
	id := core.NewRecordID(pulseNumber, recHash.Sum(nil))
	k := prefixkey(scopeIDRecord, id[:])
	geterr := m.db.db.View(func(tx *badger.Txn) error {
		_, err := tx.Get(k)
		return err
	})
	if geterr == nil {
		return id, ErrOverride
	}
	if geterr != badger.ErrKeyNotFound {
		return nil, ErrNotFound
	}

	m.set(k, record.SerializeRecord(rec))
	return id, nil
}

// GetClassIndex fetches class lifeline's index.
func (m *TransactionManager) GetClassIndex(id *core.RecordID, forupdate bool) (*index.ClassLifeline, error) {
	if forupdate {
		m.lockOnID(id)
	}
	k := prefixkey(scopeIDLifeline, id[:])
	buf, err := m.Get(k)
	if err != nil {
		return nil, err
	}
	return index.DecodeClassLifeline(buf)
}

// SetClassIndex stores class lifeline index.
func (m *TransactionManager) SetClassIndex(id *core.RecordID, idx *index.ClassLifeline) error {
	k := prefixkey(scopeIDLifeline, id[:])
	encoded, err := index.EncodeClassLifeline(idx)
	if err != nil {
		return err
	}
	m.set(k, encoded)
	return nil
}

// GetObjectIndex fetches object lifeline index.
func (m *TransactionManager) GetObjectIndex(id *core.RecordID, forupdate bool) (*index.ObjectLifeline, error) {
	if forupdate {
		m.lockOnID(id)
	}
	k := prefixkey(scopeIDLifeline, id[:])
	buf, err := m.Get(k)
	if err != nil {
		return nil, err
	}
	return index.DecodeObjectLifeline(buf)
}

// SetObjectIndex stores object lifeline index.
func (m *TransactionManager) SetObjectIndex(id *core.RecordID, idx *index.ObjectLifeline) error {
	k := prefixkey(scopeIDLifeline, id[:])
	if idx.Delegates == nil {
		idx.Delegates = map[core.RecordRef]core.RecordRef{}
	}
	encoded, err := index.EncodeObjectLifeline(idx)
	if err != nil {
		return err
	}
	m.set(k, encoded)
	return nil
}

// Get returns value by key.
func (m *TransactionManager) Get(key []byte) ([]byte, error) {
	if kv, ok := m.txupdates[string(key)]; ok {
		return kv.v, nil
	}

	txn := m.db.db.NewTransaction(false)
	defer txn.Discard()
	item, err := txn.Get(key)
	if err != nil {
		if err == badger.ErrKeyNotFound {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return item.ValueCopy(nil)
}

// GetLatestPulseNumber returns current pulse number.
func (m *TransactionManager) GetLatestPulseNumber() (core.PulseNumber, error) {
	buf, err := m.Get(prefixkey(scopeIDSystem, []byte{sysLatestPulse}))
	if err != nil {
		return 0, err
	}
	return core.PulseNumber(binary.BigEndian.Uint32(buf)), nil
}

// Set stores value by key.
func (m *TransactionManager) Set(key, value []byte) error {
	m.set(key, value)
	return nil
}
