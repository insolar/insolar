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

	"github.com/pkg/errors"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/utils"
	"github.com/insolar/insolar/instrumentation/inslogger"
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

func (s *DBStore) Update(ctx context.Context, pulse insolar.PulseNumber, actual bool, ids ...insolar.JetID) {
	s.Lock()
	defer s.Unlock()

	logger := inslogger.FromContext(ctx).WithFields(map[string]interface{}{
		"pulse":  pulse,
		"actual": actual,
		"ids":    ids,
	})

	tree := s.get(pulse)

	for _, id := range ids {
		tree.Update(id, actual)
	}
	err := s.set(pulse, tree)
	if err != nil {
		logger.Error(errors.Wrapf(err, "failed to update jets"))
	}
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
func (s *DBStore) Clone(ctx context.Context, from, to insolar.PulseNumber) {
	s.Lock()
	defer s.Unlock()

	logger := inslogger.FromContext(ctx).WithFields(map[string]interface{}{
		"from_pulse": from,
		"to_pulse":   to,
	})

	tree := s.get(from)
	newTree := tree.Clone(false)
	err := s.set(to, newTree)
	if err != nil {
		logger.Error(errors.Wrapf(err, "failed to clone jet.Tree"))
	}
}

func (s *DBStore) DeleteForPN(ctx context.Context, pulse insolar.PulseNumber) {
	logger := inslogger.FromContext(ctx).WithFields(map[string]interface{}{
		"pulse": pulse,
	})

	logger.Errorf("deleting records from db jet store is not provided")
}

type pulseKey insolar.PulseNumber

func (k pulseKey) Scope() store.Scope {
	return store.ScopeJetTree
}

func (k pulseKey) ID() []byte {
	return utils.UInt32ToBytes(uint32(k))
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
