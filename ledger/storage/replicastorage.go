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
	"bytes"
	"context"
	"encoding/gob"
	"io"

	"github.com/dgraph-io/badger"
	"github.com/insolar/insolar/insolar"
	"github.com/pkg/errors"
)

//go:generate minimock -i github.com/insolar/insolar/ledger/storage.ReplicaStorage -o ./ -s _mock.go

// ReplicaStorage is a heavy-based storage
type ReplicaStorage interface {
	SetHeavySyncedPulse(ctx context.Context, jetID insolar.ID, pulsenum insolar.PulseNumber) error
	GetHeavySyncedPulse(ctx context.Context, jetID insolar.ID) (pn insolar.PulseNumber, err error)
	GetSyncClientJetPulses(ctx context.Context, jetID insolar.ID) ([]insolar.PulseNumber, error)
	SetSyncClientJetPulses(ctx context.Context, jetID insolar.ID, pns []insolar.PulseNumber) error
	GetAllSyncClientJets(ctx context.Context) (map[insolar.ID][]insolar.PulseNumber, error)
	GetAllNonEmptySyncClientJets(ctx context.Context) (map[insolar.ID][]insolar.PulseNumber, error)
}

type replicaStorage struct {
	DB DBContext `inject:""`
}

func NewReplicaStorage() ReplicaStorage {
	return new(replicaStorage)
}

// SetHeavySyncedPulse saves last successfuly synced pulse number on heavy node.
func (rs *replicaStorage) SetHeavySyncedPulse(ctx context.Context, jetID insolar.ID, pulsenum insolar.PulseNumber) error {
	return rs.DB.Update(ctx, func(tx *TransactionManager) error {
		return tx.set(ctx, prefixkey(scopeIDSystem, jetID[:], []byte{sysLastSyncedPulseOnHeavy}), pulsenum.Bytes())
	})
}

// GetHeavySyncedPulse returns last successfuly synced pulse number on heavy node.
func (rs *replicaStorage) GetHeavySyncedPulse(ctx context.Context, jetID insolar.ID) (pn insolar.PulseNumber, err error) {
	var buf []byte
	buf, err = rs.DB.Get(ctx, prefixkey(scopeIDSystem, jetID[:], []byte{sysLastSyncedPulseOnHeavy}))
	if err == nil {
		pn = insolar.NewPulseNumber(buf)
	} else if err == insolar.ErrNotFound {
		err = nil
	}
	return
}

var sysHeavyClientStatePrefix = prefixkey(scopeIDSystem, []byte{sysHeavyClientState})

func sysHeavyClientStateKeyForJet(jetID []byte) []byte {
	return bytes.Join([][]byte{sysHeavyClientStatePrefix, jetID[:]}, nil)
}

// GetSyncClientJetPulses returns all jet's pulses not synced to heavy.
func (rs *replicaStorage) GetSyncClientJetPulses(ctx context.Context, jetID insolar.ID) ([]insolar.PulseNumber, error) {
	k := sysHeavyClientStateKeyForJet(jetID[:])
	buf, err := rs.DB.Get(ctx, k)
	if err == insolar.ErrNotFound {
		return nil, nil
	} else if err != nil {
		return nil, errors.Wrap(err, "GetSyncClientJetPulses failed")
	}
	return decodePulsesList(bytes.NewReader(buf))
}

func decodePulsesList(r io.Reader) (pns []insolar.PulseNumber, err error) {
	enc := gob.NewDecoder(r)
	err = enc.Decode(&pns)
	return
}

// SetSyncClientJetPulses saves all jet's pulses not synced to heavy.
func (rs *replicaStorage) SetSyncClientJetPulses(ctx context.Context, jetID insolar.ID, pns []insolar.PulseNumber) error {
	k := sysHeavyClientStateKeyForJet(jetID[:])
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(pns)
	if err != nil {
		return err
	}
	return rs.DB.Set(ctx, k, buf.Bytes())
}

// GetAllSyncClientJets returns map of all jet's processed by node.
func (rs *replicaStorage) GetAllSyncClientJets(ctx context.Context) (map[insolar.ID][]insolar.PulseNumber, error) {
	jets := map[insolar.ID][]insolar.PulseNumber{}
	err := rs.DB.GetBadgerDB().View(func(txn *badger.Txn) error {
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

			var jetID insolar.ID
			offset := len(sysHeavyClientStatePrefix)
			copy(jetID[:], key[offset:offset+len(jetID)])
			jets[jetID] = syncPulses
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return jets, nil
}

// GetAllNonEmptySyncClientJets returns map of all jet's if they have non empty list pulses to sync.
func (rs *replicaStorage) GetAllNonEmptySyncClientJets(ctx context.Context) (map[insolar.ID][]insolar.PulseNumber, error) {
	states, err := rs.GetAllSyncClientJets(ctx)
	if err != nil {
		return nil, err
	}
	for jetID, syncPulses := range states {
		if len(syncPulses) == 0 {
			delete(states, jetID)
		}
	}
	return states, nil
}
