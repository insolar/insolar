// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package jet

import (
	"context"
	"sync"

	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/pkg/errors"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/store"
)

type BadgerDBStore struct {
	sync.RWMutex
	db store.DB
}

func NewBadgerDBStore(db store.DB) *BadgerDBStore {
	return &BadgerDBStore{db: db}
}

func (s *BadgerDBStore) All(ctx context.Context, pulse insolar.PulseNumber) []insolar.JetID {
	s.RLock()
	defer s.RUnlock()

	tree := s.get(pulse)
	return tree.LeafIDs()
}

func (s *BadgerDBStore) ForID(ctx context.Context, pulse insolar.PulseNumber, recordID insolar.ID) (insolar.JetID, bool) {
	s.RLock()
	defer s.RUnlock()

	tree := s.get(pulse)
	return tree.Find(recordID)
}

// TruncateHead remove all records starting with 'from'
func (s *BadgerDBStore) TruncateHead(ctx context.Context, from insolar.PulseNumber) error {
	s.Lock()
	defer s.Unlock()

	it := s.db.NewIterator(pulseKey(from), false)
	defer it.Close()

	var hasKeys bool
	for it.Next() {
		hasKeys = true
		key := newPulseKey(it.Key())
		err := s.db.Delete(&key)
		if err != nil {
			return errors.Wrapf(err, "can't delete key: %+v", key)
		}

		inslogger.FromContext(ctx).Debugf("Erased key with pulse number: %s", insolar.PulseNumber(key))
	}

	if !hasKeys {
		inslogger.FromContext(ctx).Debugf("No records. Nothing done. Pulse number: %s", from.String())
	}

	return nil
}

func (s *BadgerDBStore) Update(ctx context.Context, pulse insolar.PulseNumber, actual bool, ids ...insolar.JetID) error {
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

func (s *BadgerDBStore) Split(ctx context.Context, pulse insolar.PulseNumber, id insolar.JetID) (insolar.JetID, insolar.JetID, error) {
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
func (s *BadgerDBStore) Clone(ctx context.Context, from, to insolar.PulseNumber, keepActual bool) error {
	s.Lock()
	defer s.Unlock()

	tree := s.get(from)
	newTree := tree.Clone(keepActual)
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

func (s *BadgerDBStore) get(pn insolar.PulseNumber) *Tree {
	serializedTree, err := s.db.Get(pulseKey(pn))
	if err != nil {
		return NewTree(pn == insolar.GenesisPulse.PulseNumber)
	}

	recovered := &Tree{}
	err = recovered.Unmarshal(serializedTree)
	if err != nil {
		return nil
	}
	return recovered
}

func (s *BadgerDBStore) set(pn insolar.PulseNumber, jt *Tree) error {
	key := pulseKey(pn)

	serialized, err := jt.Marshal()
	if err != nil {
		return errors.Wrap(err, "failed to serialize jet.Tree")
	}

	return s.db.Set(key, serialized)
}
