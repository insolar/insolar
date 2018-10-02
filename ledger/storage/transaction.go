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
	"github.com/dgraph-io/badger"
	"github.com/insolar/insolar/log"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/ledger/index"
	"github.com/insolar/insolar/ledger/record"
)

// TransactionManager is used to ensure persistent writes to disk.
type TransactionManager struct {
	db     *DB
	txn    *badger.Txn
	update bool
}

func prefixkey(prefix byte, key []byte) []byte {
	k := make([]byte, core.RecordRefSize+1)
	k[0] = prefix
	_ = copy(k[1:], key)
	return k
}

// Commit tries to write transaction on disk. Returns error on fail.
func (m *TransactionManager) Commit() error {
	return m.txn.Commit(nil)
}

// Discard terminates transaction without disk writes.
func (m *TransactionManager) Discard() {
	if m.update {
		m.db.dropWG.Done()
	}
	m.txn.Discard()
}

// GetRecord returns record from BadgerDB by *record.Reference.
//
// It returns ErrNotFound if the DB does not contain the key.
func (m *TransactionManager) GetRecord(id *record.ID) (record.Record, error) {
	k := prefixkey(scopeIDRecord, record.ID2Bytes(*id))
	log.Debugf("Getting record %s", id)
	item, err := m.txn.Get(k)
	if err != nil {
		if err == badger.ErrKeyNotFound {
			return nil, ErrNotFound
		}
		return nil, err
	}
	buf, err := item.Value()
	if err != nil {
		return nil, err
	}
	raw, err := record.DecodeToRaw(buf)
	if err != nil {
		return nil, err
	}
	return raw.ToRecord(), nil
}

// SetRecord stores record in BadgerDB and returns *record.Reference of new record.
//
// If record exists returns ErrOverride error.
func (m *TransactionManager) SetRecord(rec record.Record) (*record.ID, error) {
	raw, err := record.EncodeToRaw(rec)
	if err != nil {
		return nil, err
	}
	id := record.ID{
		Pulse: m.db.GetCurrentPulse(),
		Hash:  raw.Hash(),
	}
	k := prefixkey(scopeIDRecord, record.ID2Bytes(id))
	_, geterr := m.txn.Get(k)
	if geterr == nil {
		return nil, ErrOverride
	}
	if geterr != badger.ErrKeyNotFound {
		return nil, geterr
	}

	err = m.txn.Set(k, record.MustEncodeRaw(raw))
	if err != nil {
		return nil, err
	}
	return &id, nil
}

// GetClassIndex fetches class lifeline's index.
func (m *TransactionManager) GetClassIndex(id *record.ID) (*index.ClassLifeline, error) {
	k := prefixkey(scopeIDLifeline, record.ID2Bytes(*id))
	item, err := m.txn.Get(k)
	if err != nil {
		if err == badger.ErrKeyNotFound {
			return nil, ErrNotFound
		}
		return nil, err
	}
	buf, err := item.Value()
	if err != nil {
		return nil, err
	}
	return index.DecodeClassLifeline(buf)
}

// SetClassIndex stores class lifeline index.
func (m *TransactionManager) SetClassIndex(id *record.ID, idx *index.ClassLifeline) error {
	k := prefixkey(scopeIDLifeline, record.ID2Bytes(*id))
	encoded, err := index.EncodeClassLifeline(idx)
	if err != nil {
		return err
	}
	return m.txn.Set(k, encoded)
}

// GetObjectIndex fetches object lifeline index.
func (m *TransactionManager) GetObjectIndex(id *record.ID) (*index.ObjectLifeline, error) {
	k := prefixkey(scopeIDLifeline, record.ID2Bytes(*id))
	item, err := m.txn.Get(k)
	if err != nil {
		if err == badger.ErrKeyNotFound {
			return nil, ErrNotFound
		}
		return nil, err
	}
	buf, err := item.Value()
	if err != nil {
		return nil, err
	}
	return index.DecodeObjectLifeline(buf)
}

// SetObjectIndex stores object lifeline index.
func (m *TransactionManager) SetObjectIndex(id *record.ID, idx *index.ObjectLifeline) error {
	k := prefixkey(scopeIDLifeline, record.ID2Bytes(*id))
	if idx.Delegates == nil {
		idx.Delegates = map[core.RecordRef]record.Reference{}
	}
	encoded, err := index.EncodeObjectLifeline(idx)
	if err != nil {
		return err
	}
	return m.txn.Set(k, encoded)
}

// GetEntropy returns entropy from storage for given pulse.
//
// GeneratedEntropy is used for calculating node roles.
func (m *TransactionManager) GetEntropy(pulse core.PulseNumber) (*core.Entropy, error) {
	k := prefixkey(scopeIDEntropy, pulse.Bytes())
	item, err := m.txn.Get(k)
	if err != nil {
		if err == badger.ErrKeyNotFound {
			return nil, ErrNotFound
		}
		return nil, err
	}
	buf, err := item.Value()
	if err != nil {
		return nil, err
	}
	var entropy core.Entropy
	copy(entropy[:], buf)
	return &entropy, nil
}

// SetEntropy stores given entropy for given pulse in storage.
//
// GeneratedEntropy is used for calculating node roles.
func (m *TransactionManager) SetEntropy(pulse core.PulseNumber, entropy core.Entropy) error {
	k := prefixkey(scopeIDEntropy, pulse.Bytes())
	return m.txn.Set(k, entropy[:])
}

// Get returns value by key.
func (m *TransactionManager) Get(key []byte) ([]byte, error) {
	// var buf []byte
	item, err := m.txn.Get(key)
	if err != nil {
		if err == badger.ErrKeyNotFound {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return item.ValueCopy(nil)
}

// Set stores value by key.
func (m *TransactionManager) Set(key, value []byte) error {
	return m.txn.Set(key, value)
}
