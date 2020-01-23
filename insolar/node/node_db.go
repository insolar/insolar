//
// Copyright 2020 Insolar Technologies GmbH
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

package node

import (
	"context"

	"github.com/dgraph-io/badger"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/store"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/pkg/errors"
)

// AALEKSEEV TODO use PostgreSQL + see node_db_test

// StorageDB is a badger-based impl of a node storage.
type StorageDB struct {
	db *store.BadgerDB
}

// NewStorageDB create new instance of StorageDB.
func NewStorageDB(db *store.BadgerDB) *StorageDB {
	return &StorageDB{db: db}
}

type nodeHistoryKey insolar.PulseNumber

func (k nodeHistoryKey) Scope() store.Scope {
	return store.ScopeNodeHistory
}

func (k nodeHistoryKey) DebugString() string {
	pn := insolar.PulseNumber(k)
	return "nodeHistoryKey. " + pn.String()
}

func (k nodeHistoryKey) ID() []byte {
	pn := insolar.PulseNumber(k)
	return pn.Bytes()
}

// Set saves active nodes for pulse in memory.
func (s *StorageDB) Set(pulse insolar.PulseNumber, nodes []insolar.Node) error {
	nodesList := &insolar.NodeList{}
	if len(nodes) != 0 {
		nodesList.Nodes = nodes
	}
	rawNodes, err := nodesList.Marshal()
	if err != nil {
		return err
	}
	return s.db.Backend().Update(func(txn *badger.Txn) error {
		key := nodeHistoryKey(pulse)
		fullKey := append(key.Scope().Bytes(), key.ID()...)
		_, err = txn.Get(fullKey)
		if err != nil && err != badger.ErrKeyNotFound {
			return err
		}
		if err == nil {
			return ErrOverride
		}

		return txn.Set(fullKey, rawNodes)
	})
}

// All return active nodes for specified pulse.
func (s *StorageDB) All(pulse insolar.PulseNumber) ([]insolar.Node, error) {
	var res []insolar.Node
	err := s.db.Backend().View(func(txn *badger.Txn) error {
		key := nodeHistoryKey(pulse)
		fullKey := append(key.Scope().Bytes(), key.ID()...)
		item, err := txn.Get(fullKey)
		if err != nil {
			if err == badger.ErrKeyNotFound {
				return ErrNoNodes
			}
			return err
		}

		buff, err := item.ValueCopy(nil)
		if err != nil {
			return err
		}
		list := &insolar.NodeList{}
		err = list.Unmarshal(buff)
		if err != nil {
			return err
		}
		res = list.Nodes

		return nil
	})
	return res, err
}

// InRole return active nodes for specified pulse and role.
func (s *StorageDB) InRole(pulse insolar.PulseNumber, role insolar.StaticRole) ([]insolar.Node, error) {
	nodes, err := s.All(pulse)
	if err != nil {
		return nil, err
	}
	var inRole []insolar.Node
	for _, node := range nodes {
		if node.Role == role {
			inRole = append(inRole, node)
		}
	}

	return inRole, nil
}

// DeleteForPN erases nodes for specified pulse.
func (s *StorageDB) DeleteForPN(pulse insolar.PulseNumber) {
	panic("implement me")
}

// TruncateHead remove all records after lastPulse
func (s *StorageDB) TruncateHead(ctx context.Context, from insolar.PulseNumber) error {
	it := s.db.NewIterator(nodeHistoryKey(from), false)
	defer it.Close()

	var hasKeys bool
	for it.Next() {
		hasKeys = true
		key := nodeHistoryKey(insolar.NewPulseNumber(it.Key()))
		err := s.db.Delete(&key)
		if err != nil {
			return errors.Wrapf(err, "can't delete key: %+v", key)
		}

		inslogger.FromContext(ctx).Debugf("Node db. Erased key. Pulse number: %s", key.DebugString())
	}
	if !hasKeys {
		inslogger.FromContext(ctx).Debug("Node db. No records. Nothing done. Pulse number: " + from.String())
	}

	return nil
}
