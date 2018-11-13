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
	"context"
	"encoding/binary"
	"encoding/hex"
	"sync"

	"github.com/dgraph-io/badger"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/cryptohelpers/hash"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/index"
	"github.com/insolar/insolar/ledger/record"
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

	recentObjectsIndex *recentObjectsIndex
}

type recentObjectsIndex struct{
	fetchedObjects map[string]*core.RecordID
	updatedObjects map[string]*core.RecordID

	fetchedObjectsLock sync.Mutex
	updatedObjectsLock sync.Mutex
}

func (r *recentObjectsIndex) addToFetched(id *core.RecordID){
	r.fetchedObjectsLock.Lock()
	defer r.fetchedObjectsLock.Unlock()

	r.fetchedObjects[id.String()] = id
}

func (r *recentObjectsIndex) addToUpdated(id *core.RecordID){
	r.updatedObjectsLock.Lock()
	defer r.updatedObjectsLock.Unlock()

	r.updatedObjects[id.String()] = id
}




type byte2hex byte

func (b byte2hex) String() string {
	return hex.EncodeToString([]byte{byte(b)})
}

type bytes2hex []byte

func (h bytes2hex) String() string {
	return hex.EncodeToString(h)
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
	if err != nil {
		return err
	}
	return tx.Commit(nil)
}

// Discard terminates transaction without disk writes.
func (m *TransactionManager) Discard() {
	m.txupdates = nil
	m.releaseLocks()
	if m.update {
		m.db.dropWG.Done()
	}
}

// GetRequest returns request record from BadgerDB by *record.Reference.
//
// It returns ErrNotFound if the DB does not contain the key.
func (m *TransactionManager) GetRequest(ctx context.Context, id *core.RecordID) (record.Request, error) {
	rec, err := m.GetRecord(ctx, id)
	if err != nil {
		return nil, err
	}
	// TODO: return error if record is not a request.
	req := rec.(record.Request)
	return req, nil
}

// GetBlob returns binary value stored by record ID.
func (m *TransactionManager) GetBlob(ctx context.Context, id *core.RecordID) ([]byte, error) {
	k := prefixkey(scopeIDBlob, id[:])
	inslogger.FromContext(ctx).Debugf(
		"GetRecord by id %v (prefix=%v)", id, byte2hex(scopeIDBlob))
	return m.get(ctx, k)
}

// SetBlob saves binary value for provided pulse.
func (m *TransactionManager) SetBlob(ctx context.Context, pulseNumber core.PulseNumber, blob []byte) (*core.RecordID, error) {
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

	err := m.set(ctx, k, blob)
	if err != nil {
		return nil, err
	}
	return id, nil
}

// GetRecord returns record from BadgerDB by *record.Reference.
//
// It returns ErrNotFound if the DB does not contain the key.
func (m *TransactionManager) GetRecord(ctx context.Context, id *core.RecordID) (record.Record, error) {
	k := prefixkey(scopeIDRecord, id[:])
	inslogger.FromContext(ctx).Debugf(
		"GetRecord by id %v (prefix=%v)", id, byte2hex(scopeIDRecord))
	buf, err := m.get(ctx, k)
	if err != nil {
		return nil, err
	}
	return record.DeserializeRecord(buf), nil
}

// SetRecord stores record in BadgerDB and returns *record.ID of new record.
//
// If record exists returns both *record.ID and ErrOverride error.
// If record not found returns nil and ErrNotFound error
func (m *TransactionManager) SetRecord(ctx context.Context, pulseNumber core.PulseNumber, rec record.Record) (*core.RecordID, error) {
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

	err = m.set(ctx, k, record.SerializeRecord(rec))
	if err != nil {
		return nil, err
	}
	return id, nil
}

// GetObjectIndex fetches object lifeline index.
func (m *TransactionManager) GetObjectIndex(
	ctx context.Context,
	id *core.RecordID,
	forupdate bool,
) (*index.ObjectLifeline, error) {
	if forupdate {
		m.lockOnID(id)
	}
	k := prefixkey(scopeIDLifeline, id[:])
	buf, err := m.get(ctx, k)
	if err != nil {
		return nil, err
	}
	index, err := index.DecodeObjectLifeline(buf)
	if err == nil{
		m.recentObjectsIndex.addToFetched(id)
	}
	return index,err
}

// SetObjectIndex stores object lifeline index.
func (m *TransactionManager) SetObjectIndex(
	ctx context.Context,
	id *core.RecordID,
	idx *index.ObjectLifeline,
) error {
	k := prefixkey(scopeIDLifeline, id[:])
	if idx.Delegates == nil {
		idx.Delegates = map[core.RecordRef]core.RecordRef{}
	}
	id.String()
	encoded, err := index.EncodeObjectLifeline(idx)
	if err != nil {
		return err
	}
	err = m.set(ctx, k, encoded)
	if err == nil{
		m.recentObjectsIndex.addToUpdated(id)
	}
	return err
}

// GetLatestPulseNumber returns current pulse number.
func (m *TransactionManager) GetLatestPulseNumber(ctx context.Context) (core.PulseNumber, error) {
	buf, err := m.get(ctx, prefixkey(scopeIDSystem, []byte{sysLatestPulse}))
	if err != nil {
		return 0, err
	}
	return core.PulseNumber(binary.BigEndian.Uint32(buf)), nil
}

func (m *TransactionManager) GetLatestObjects(ctx context.Context) ([]*index.ObjectLifeline) {
	indexCount := len(m.recentObjectsIndex.fetchedObjects) + len(m.recentObjectsIndex.updatedObjects)
	result :=  make([]*index.ObjectLifeline, 0, indexCount)
	resultLock := sync.Mutex{}

	wg := sync.WaitGroup{}
	wg.Add(indexCount)

	fetchIndexFunc := func(id *core.RecordID) {
		defer wg.Done()
		k := prefixkey(scopeIDLifeline, id[:])
		buf, err := m.get(ctx, k)
		if err != nil {
			inslogger.FromContext(ctx).Errorf("problems with fetching index - %v", err)
		}
		index, err := index.DecodeObjectLifeline(buf)
		if err != nil {
			inslogger.FromContext(ctx).Errorf("problems with decoding index - %v", err)
		}

		resultLock.Lock()
		result = append(result, index)
		resultLock.Unlock()
	}

	for _, value := range m.recentObjectsIndex.fetchedObjects{
		go fetchIndexFunc(value)
	}
	for _, value := range m.recentObjectsIndex.updatedObjects{
		go fetchIndexFunc(value)
	}

	wg.Wait()

	return result
}

// set stores value by key.
func (m *TransactionManager) set(ctx context.Context, key, value []byte) error {
	inslogger.FromContext(ctx).Debugf("set key %v", bytes2hex(key))

	m.txupdates[string(key)] = keyval{k: key, v: value}
	return nil
}

// get returns value by key.
func (m *TransactionManager) get(ctx context.Context, key []byte) ([]byte, error) {
	inslogger.FromContext(ctx).Debugf("get key %v", bytes2hex(key))

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
