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
	"encoding/gob"
	"io"

	"github.com/dgraph-io/badger"
	"github.com/insolar/insolar/core"
)

// SetHeavySyncedPulse saves last successfuly synced pulse number on heavy node.
func (db *DB) SetHeavySyncedPulse(ctx context.Context, jetID core.RecordID, pulsenum core.PulseNumber) error {
	return db.Update(ctx, func(tx *TransactionManager) error {
		return tx.set(ctx, prefixkey(scopeIDSystem, jetID[:], []byte{sysLastSyncedPulseOnHeavy}), pulsenum.Bytes())
	})
}

// GetHeavySyncedPulse returns last successfuly synced pulse number on heavy node.
func (db *DB) GetHeavySyncedPulse(ctx context.Context, jetID core.RecordID) (pn core.PulseNumber, err error) {
	var buf []byte
	buf, err = db.get(ctx, prefixkey(scopeIDSystem, jetID[:], []byte{sysLastSyncedPulseOnHeavy}))
	if err == nil {
		pn = core.NewPulseNumber(buf)
	} else if err == ErrNotFound {
		err = nil
	}
	return
}

var sysHeavyClientStatePrefix = prefixkey(scopeIDSystem, []byte{sysHeavyClientState})

func sysHeavyClientStateKeyForJet(jetID []byte) []byte {
	return bytes.Join([][]byte{sysHeavyClientStatePrefix, jetID[:]}, nil)
}

// GetSyncClientJetPulses returns all jet's pulses not synced to heavy.
func (db *DB) GetSyncClientJetPulses(ctx context.Context, jetID core.RecordID) ([]core.PulseNumber, error) {
	db.jetHeavyClientLocker.Lock(&jetID)
	defer db.jetHeavyClientLocker.Unlock(&jetID)
	return db.getSyncClientJetPulses(ctx, jetID)
}

func (db *DB) getSyncClientJetPulses(ctx context.Context, jetID core.RecordID) ([]core.PulseNumber, error) {
	k := sysHeavyClientStateKeyForJet(jetID[:])
	buf, err := db.get(ctx, k)
	if err != nil {
		if err == ErrNotFound {
			err = nil
		}
		return nil, err
	}
	return decodePulsesList(bytes.NewReader(buf))
}

func decodePulsesList(r io.Reader) (pns []core.PulseNumber, err error) {
	enc := gob.NewDecoder(r)
	err = enc.Decode(&pns)
	return
}

// SetSyncClientJetPulses saves all jet's pulses not synced to heavy.
func (db *DB) SetSyncClientJetPulses(ctx context.Context, jetID core.RecordID, pns []core.PulseNumber) error {
	db.jetHeavyClientLocker.Lock(&jetID)
	defer db.jetHeavyClientLocker.Unlock(&jetID)
	return db.setSyncClientJetPulses(ctx, jetID, pns)
}

func (db *DB) setSyncClientJetPulses(ctx context.Context, jetID core.RecordID, pns []core.PulseNumber) error {
	k := sysHeavyClientStateKeyForJet(jetID[:])
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(pns)
	if err != nil {
		return err
	}
	return db.set(ctx, k, buf.Bytes())
}

// GetAllSyncClientJets returns map of all jet's if they have non empty list pulses to sync.
func (db *DB) GetAllSyncClientJets(ctx context.Context) (map[core.RecordID][]core.PulseNumber, error) {
	jets := map[core.RecordID][]core.PulseNumber{}
	err := db.db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()

		for it.Seek(sysHeavyClientStatePrefix); it.ValidForPrefix(sysHeavyClientStatePrefix); it.Next() {
			item := it.Item()
			if item == nil {
				break
			}
			key := item.Key()
			value, err := it.Item().Value()
			if err != nil {
				return err
			}
			syncPulses, err := decodePulsesList(bytes.NewReader(value))
			if err != nil {
				return err
			}
			if len(syncPulses) > 0 {
				var jetID core.RecordID
				offset := len(sysHeavyClientStatePrefix)
				copy(jetID[:], key[offset:offset+len(jetID)])
				jets[jetID] = syncPulses
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return jets, nil
}
