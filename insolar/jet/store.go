// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package jet

import (
	"context"
	"sync"

	"github.com/insolar/insolar/insolar"
)

var (
	_ Accessor = &Store{}
	_ Modifier = &Store{}
)

type lockedTree struct {
	sync.RWMutex
	t *Tree
}

func (lt *lockedTree) find(recordID insolar.ID) (insolar.JetID, bool) {
	lt.RLock()
	defer lt.RUnlock()
	return lt.t.Find(recordID)
}

func (lt *lockedTree) update(id insolar.JetID, setActual bool) {
	lt.Lock()
	defer lt.Unlock()
	lt.t.Update(id, setActual)
}

func (lt *lockedTree) leafIDs() []insolar.JetID {
	lt.RLock()
	defer lt.RUnlock()
	return lt.t.LeafIDs()
}

func (lt *lockedTree) clone(keep bool) *Tree {
	lt.RLock()
	defer lt.RUnlock()
	return lt.t.Clone(keep)
}

func (lt *lockedTree) split(id insolar.JetID) (insolar.JetID, insolar.JetID, error) {
	lt.RLock()
	defer lt.RUnlock()
	return lt.t.Split(id)
}

// Store stores jet trees per pulse.
// It provides methods for querying and modification this trees.
type Store struct {
	sync.RWMutex
	trees map[insolar.PulseNumber]*lockedTree
}

// NewStore creates new Store instance.
func NewStore() *Store {
	return &Store{
		trees: map[insolar.PulseNumber]*lockedTree{},
	}
}

// All returns all jet from jet tree for provided pulse.
func (s *Store) All(ctx context.Context, pulse insolar.PulseNumber) []insolar.JetID {
	return s.ltreeForPulse(pulse).leafIDs()
}

// ForID finds jet in jet tree for provided pulse and object.
// Always returns jet id and activity flag for this jet.
func (s *Store) ForID(ctx context.Context, pulse insolar.PulseNumber, recordID insolar.ID) (insolar.JetID, bool) {
	return s.ltreeForPulse(pulse).find(recordID)
}

// Update updates jet tree for specified pulse.
func (s *Store) Update(ctx context.Context, pulse insolar.PulseNumber, setActual bool, ids ...insolar.JetID) error {
	s.Lock()
	defer s.Unlock()

	ltree := s.ltreeForPulseUnsafe(pulse)
	for _, id := range ids {
		ltree.update(id, setActual)
	}
	// required because TreeForPulse could return new tree.
	s.trees[pulse] = ltree
	return nil
}

func (s *Store) Split(
	ctx context.Context, pulse insolar.PulseNumber, id insolar.JetID,
) (insolar.JetID, insolar.JetID, error) {
	ltree := s.ltreeForPulse(pulse)
	left, right, err := ltree.split(id)
	if err != nil {
		return insolar.ZeroJetID, insolar.ZeroJetID, err
	}
	return left, right, nil
}

// Clone copies tree from one pulse to another. Use it to copy the past tree into new pulse.
func (s *Store) Clone(
	ctx context.Context, from, to insolar.PulseNumber, keepActual bool,
) error {
	newTree := s.ltreeForPulse(from).clone(keepActual)

	s.Lock()
	defer s.Unlock()

	s.trees[to] = &lockedTree{
		t: newTree,
	}
	return nil
}

// Delete jets for pulse (concurrent safe).
func (s *Store) DeleteForPN(
	ctx context.Context, pulse insolar.PulseNumber,
) {
	s.Lock()
	defer s.Unlock()
	delete(s.trees, pulse)
}

// ltreeForPulse returns jet tree with lock for pulse, it's concurrent safe.
func (s *Store) ltreeForPulse(pulse insolar.PulseNumber) *lockedTree {
	s.Lock()
	defer s.Unlock()
	return s.ltreeForPulseUnsafe(pulse)
}

// ltreeForPulseUnsafe returns jet tree with lock for pulse, it's concurrent unsafe and requires write lock.
func (s *Store) ltreeForPulseUnsafe(pulse insolar.PulseNumber) *lockedTree {
	if ltree, ok := s.trees[pulse]; ok {
		return ltree
	}

	ltree := &lockedTree{
		t: NewTree(pulse == insolar.GenesisPulse.PulseNumber),
	}
	s.trees[pulse] = ltree
	return ltree
}
