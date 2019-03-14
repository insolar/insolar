/*
 *    Copyright 2019 Insolar
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

package jet

import (
	"context"
	"sync"

	"github.com/insolar/insolar/core"
)

type Store struct {
	sync.RWMutex
	trees map[core.PulseNumber]*Tree
}

var (
	_ Accessor = &Store{}
	_ Modifier = &Store{}
)

func NewStore() *Store {
	return &Store{
		trees: map[core.PulseNumber]*Tree{},
	}
}

// TODO: add test if empty tree has at least one jetID
func (s *Store) All(ctx context.Context, pulse core.PulseNumber) []core.JetID {
	s.RLock()
	defer s.RUnlock()
	return s.treeForPulse(ctx, pulse).LeafIDs()
}

// ForID finds jet for specified pulse and object.
func (s *Store) ForID(ctx context.Context, pulse core.PulseNumber, recordID core.RecordID) (core.JetID, bool) {
	var t *Tree
	s.RLock()
	t, _ = s.trees[pulse]
	s.RUnlock()
	if t == nil {
		t = s.TreeForPulse(ctx, pulse)
	}
	return t.Find(recordID)
}

// Update updates jet tree for specified pulse.
func (s *Store) Update(ctx context.Context, pulse core.PulseNumber, setActual bool, ids ...core.JetID) {
	s.Lock()
	defer s.Unlock()

	tree := s.treeForPulse(ctx, pulse)
	for _, id := range ids {
		tree.Update(id, setActual)
	}
}

// Split performs jet split and returns resulting jet ids.
func (s *Store) Split(
	ctx context.Context, pulse core.PulseNumber, id core.JetID,
) (core.JetID, core.JetID, error) {
	s.Lock()
	defer s.Unlock()

	tree := s.treeForPulse(ctx, pulse)
	left, right, err := tree.Split(id)
	if err != nil {
		return core.ZeroJetID, core.ZeroJetID, err
	}
	return left, right, nil
}

// TODO: rename?
// Clone copies tree from one pulse to another. Use it to copy past tree into new pulse.
func (s *Store) Clone(
	ctx context.Context, from, to core.PulseNumber,
) {
	s.Lock()
	defer s.Unlock()
	s.trees[to] = s.treeForPulse(ctx, from).Clone(false)
}

// Delete concurrent safe
func (s *Store) Delete(
	ctx context.Context, pulse core.PulseNumber,
) {
	s.Lock()
	defer s.Unlock()
	delete(s.trees, pulse)
}

// TreeForPulse concurrent safe ...
func (s *Store) TreeForPulse(ctx context.Context, pulse core.PulseNumber) *Tree {
	s.Lock()
	defer s.Unlock()
	return s.treeForPulse(ctx, pulse)
}

func (s *Store) treeForPulse(ctx context.Context, pulse core.PulseNumber) *Tree {
	if t, ok := s.trees[pulse]; ok {
		return t
	}

	actualDefault := pulse == core.GenesisPulse.PulseNumber
	tree := NewTree(actualDefault)
	s.trees[pulse] = tree
	return tree
}
