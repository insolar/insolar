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

package jet

import (
	"context"
	"sync"

	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/pkg/errors"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/internal/ledger/store"
)

type DBStore struct {
	sync.RWMutex
	db store.DB
}

func NewDBStore(db store.DB) *DBStore {
	return &DBStore{db: db}
}

func (s *DBStore) All(ctx context.Context, pulse insolar.PulseNumber) []insolar.JetID {
	s.RLock()
	defer s.RUnlock()

	tree := s.get(pulse)
	return tree.LeafIDs()
}

func (s *DBStore) ForID(ctx context.Context, pulse insolar.PulseNumber, recordID insolar.ID) (insolar.JetID, bool) {
	s.RLock()
	defer s.RUnlock()

	tree := s.get(pulse)
	return tree.Find(recordID)
}

// TruncateHead remove all records after lastPulse
func (s *DBStore) TruncateHead(ctx context.Context, lastPulse insolar.PulseNumber) error {
	s.Lock()
	defer s.Unlock()

	it := s.db.NewIterator(pulseKey(lastPulse), false)
	defer it.Close()

	if !it.Next() {
		inslogger.FromContext(ctx).Infof("[ DBStore.TruncateHead ] No records. Nothing done. Pulse number: %s", lastPulse.String())
		return nil
	}

	for it.Next() {
		key := newPulseKey(it.Key())
		err := s.db.Delete(&key)
		if err != nil {
			return errors.Wrapf(err, "[ DBStore.TruncateHead ] Can't Delete key: %+v", key)
		}

		inslogger.FromContext(ctx).Infof("[ DBStore.TruncateHead ] erased key. Pulse number: %s", insolar.PulseNumber(key))
	}
	return nil
}

func (s *DBStore) Update(ctx context.Context, pulse insolar.PulseNumber, actual bool, ids ...insolar.JetID) error {
	s.Lock()
	defer s.Unlock()

	tree := s.get(pulse)

	for _, id := range ids {
		tree.Update(id, actual)
	}
	err := s.set(pulse, tree)
	if err != nil {
		return errors.Wrapf(err, "failed to update jets")
	}
	return nil
}

func (s *DBStore) Split(ctx context.Context, pulse insolar.PulseNumber, id insolar.JetID) (insolar.JetID, insolar.JetID, error) {
	s.Lock()
	defer s.Unlock()

	tree := s.get(pulse)
	left, right, err := tree.Split(id)
	if err != nil {
		return insolar.ZeroJetID, insolar.ZeroJetID, err
	}
	err = s.set(pulse, tree)
	if err != nil {
		return insolar.ZeroJetID, insolar.ZeroJetID, err
	}
	return left, right, nil
}
func (s *DBStore) Clone(ctx context.Context, from, to insolar.PulseNumber) error {
	s.Lock()
	defer s.Unlock()

	tree := s.get(from)
	newTree := tree.Clone(false)
	err := s.set(to, newTree)
	if err != nil {
		return errors.Wrapf(err, "failed to clone jet.Tree")
	}
	return nil
}

type pulseKey insolar.PulseNumber

func (k pulseKey) Scope() store.Scope {
	return store.ScopeJetTree
}

func (k pulseKey) ID() []byte {
	return insolar.PulseNumber(k).Bytes()
}

func newPulseKey(raw []byte) pulseKey {
	key := pulseKey(insolar.NewPulseNumber(raw))
	return key
}

func (s *DBStore) get(pn insolar.PulseNumber) *Tree {
	serializedTree, err := s.db.Get(pulseKey(pn))
	if err != nil {
		return NewTree(pn == insolar.GenesisPulse.PulseNumber)
	}

	recovered := &Tree{}
	err = insolar.Deserialize(serializedTree, recovered)
	if err != nil {
		return nil
	}
	return recovered
}

func (s *DBStore) set(pn insolar.PulseNumber, jt *Tree) error {
	key := pulseKey(pn)

	serialized, err := insolar.Serialize(jt)
	if err != nil {
		return errors.Wrap(err, "failed to serialize jet.Tree")
	}

	return s.db.Set(key, serialized)
}
