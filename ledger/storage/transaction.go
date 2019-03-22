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

	"github.com/dgraph-io/badger"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/ledger/storage/object"
)

type keyval struct {
	k []byte
	v []byte
}

// TransactionManager is used to ensure persistent writes to disk.
type TransactionManager struct {
	db        *DB
	update    bool
	locks     []*insolar.ID
	txupdates map[string]keyval
}

func (m *TransactionManager) lockOnID(id *insolar.ID) {
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
func (m *TransactionManager) GetRequest(ctx context.Context, jetID insolar.ID, id *insolar.ID) (object.Request, error) {
	rec, err := m.GetRecord(ctx, jetID, id)
	if err != nil {
		return nil, err
	}
	// TODO: return error if record is not a request.
	req := rec.(object.Request)
	return req, nil
}

// GetBlob returns binary value stored by record ID.
func (m *TransactionManager) GetBlob(ctx context.Context, jetID insolar.ID, id *insolar.ID) ([]byte, error) {
	jetPrefix := insolar.JetID(jetID).Prefix()
	k := prefixkey(scopeIDBlob, jetPrefix, id[:])
	return m.get(ctx, k)
}

// SetBlob saves binary value for provided pulse.
func (m *TransactionManager) SetBlob(ctx context.Context, jetID insolar.ID, pulseNumber insolar.PulseNumber, blob []byte) (*insolar.ID, error) {
	id := object.CalculateIDForBlob(m.db.PlatformCryptographyScheme, pulseNumber, blob)
	jetPrefix := insolar.JetID(jetID).Prefix()
	k := prefixkey(scopeIDBlob, jetPrefix, id[:])

	// TODO: @andreyromancev. 16.01.19. Blob override is ok.
	// geterr := m.db.db.View(func(tx *badger.Txn) error {
	// 	_, err := tx.Get(k)
	// 	return err
	// })
	// if geterr == nil {
	// 	return id, ErrOverride
	// }
	// if geterr != badger.ErrKeyNotFound {
	// 	return nil, ErrNotFound
	// }

	err := m.set(ctx, k, blob)
	if err != nil {
		return nil, err
	}
	return id, nil
}

// GetRecord returns record from BadgerDB by *record.Reference.
//
// It returns ErrNotFound if the DB does not contain the key.
func (m *TransactionManager) GetRecord(ctx context.Context, jetID insolar.ID, id *insolar.ID) (object.VirtualRecord, error) {
	jetPrefix := insolar.JetID(jetID).Prefix()
	k := prefixkey(scopeIDRecord, jetPrefix, id[:])
	buf, err := m.get(ctx, k)
	if err != nil {
		return nil, err
	}
	return object.DeserializeRecord(buf), nil
}

// SetRecord stores record in BadgerDB and returns *record.ID of new record.
//
// If record exists returns both *record.ID and ErrOverride error.
// If record not found returns nil and ErrNotFound error
func (m *TransactionManager) SetRecord(ctx context.Context, jetID insolar.ID, pulseNumber insolar.PulseNumber, rec object.VirtualRecord) (*insolar.ID, error) {
	id := object.NewRecordIDFromRecord(m.db.PlatformCryptographyScheme, pulseNumber, rec)
	prefix := insolar.JetID(jetID).Prefix()
	k := prefixkey(scopeIDRecord, prefix, id[:])
	geterr := m.db.db.View(func(tx *badger.Txn) error {
		_, err := tx.Get(k)
		return err
	})
	if geterr == nil {
		return id, ErrOverride
	}
	if geterr != badger.ErrKeyNotFound {
		return nil, geterr
	}

	err := m.set(ctx, k, object.SerializeRecord(rec))
	if err != nil {
		return nil, err
	}
	return id, nil
}

// GetObjectIndex fetches object lifeline index.
func (m *TransactionManager) GetObjectIndex(
	ctx context.Context,
	jetID insolar.ID,
	id *insolar.ID,
	forupdate bool,
) (*object.Lifeline, error) {
	if forupdate {
		m.lockOnID(id)
	}
	prefix := insolar.JetID(jetID).Prefix()
	k := prefixkey(scopeIDLifeline, prefix, id[:])
	buf, err := m.get(ctx, k)
	if err != nil {
		return nil, err
	}
	res := object.DecodeIndex(buf)
	return &res, nil
}

// SetObjectIndex stores object lifeline index.
func (m *TransactionManager) SetObjectIndex(
	ctx context.Context,
	jetID insolar.ID,
	id *insolar.ID,
	idx *object.Lifeline,
) error {
	prefix := insolar.JetID(jetID).Prefix()
	k := prefixkey(scopeIDLifeline, prefix, id[:])
	if idx.Delegates == nil {
		idx.Delegates = map[insolar.Reference]insolar.Reference{}
	}
	encoded := object.EncodeIndex(*idx)
	return m.set(ctx, k, encoded)
}

// RemoveObjectIndex removes an index of an object
func (m *TransactionManager) RemoveObjectIndex(
	ctx context.Context,
	jetID insolar.ID,
	ref *insolar.ID,
) error {
	m.lockOnID(ref)
	prefix := insolar.JetID(jetID).Prefix()
	k := prefixkey(scopeIDLifeline, prefix, ref[:])
	return m.remove(ctx, k)
}

// set stores value by key.
func (m *TransactionManager) set(ctx context.Context, key, value []byte) error {
	m.txupdates[string(key)] = keyval{k: key, v: value}
	return nil
}

// get returns value by key.
func (m *TransactionManager) get(ctx context.Context, key []byte) ([]byte, error) {
	if kv, ok := m.txupdates[string(key)]; ok {
		return kv.v, nil
	}

	txn := m.db.db.NewTransaction(false)
	defer txn.Discard()
	item, err := txn.Get(key)
	if err != nil {
		if err == badger.ErrKeyNotFound {
			return nil, insolar.ErrNotFound
		}
		return nil, err
	}
	return item.ValueCopy(nil)
}

// removes value by key
func (m *TransactionManager) remove(ctx context.Context, key []byte) error {
	debugf(ctx, "get key %v", bytes2hex(key))

	txn := m.db.db.NewTransaction(true)
	defer txn.Discard()

	err := txn.Delete(key)
	if err != nil {
		return err
	}

	return txn.Commit(nil)
}
