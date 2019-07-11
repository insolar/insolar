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

package sequence

import (
	"bytes"
	"context"
	"sync"

	"github.com/pkg/errors"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/internal/ledger/store"
)

// Item is a structure to hold unified db record.
type Item struct {
	Key   []byte
	Value []byte
}

//go:generate minimock -i github.com/insolar/insolar/ledger/heavy/sequence.Sequencer -o ./ -s _mock.go

// Sequencer is an interface to work with db entity sequence.
type Sequencer interface {
	// Len returns count of records in db scope by the pulse.
	Len(scope store.Scope, pulse insolar.PulseNumber) int
	// First returns first item in the db scope.
	First(scope store.Scope) *Item
	// Last returns last item in the db scope.
	Last(scope store.Scope) *Item
	// Slice returns slice of records from provided position with corresponding limit.
	Slice(scope store.Scope, from insolar.PulseNumber, skip uint32, to insolar.PulseNumber, limit uint32) []Item
	// Upsert updates or inserts sequence of db records.
	Upsert(scope store.Scope, sequence []Item)
}

type sequencer struct {
	sync.RWMutex
	db store.DB
}

func NewSequencer(db store.DB) Sequencer {
	return &sequencer{db: db}
}

func (s *sequencer) Len(scope store.Scope, pulse insolar.PulseNumber) int {
	s.RLock()
	defer s.RUnlock()

	result := 0
	it := s.db.NewIterator(polyKey{[]byte{}, scope}, false)
	defer it.Close()
	for it.Next() && bytes.HasPrefix(it.Key(), pulse.Bytes()) {
		result++
	}
	return result
}

func (s *sequencer) First(scope store.Scope) *Item {
	it := s.db.NewIterator(polyKey{id: []byte{}, scope: scope}, false)
	defer it.Close()

	if !it.Next() {
		return nil
	}
	return &Item{Key: it.Key(), Value: it.Value()}
}

func (s *sequencer) Last(scope store.Scope) *Item {
	it := s.db.NewIterator(polyKey{id: []byte{0xFF, 0xFF, 0xFF, 0xFF}, scope: scope}, true)
	defer it.Close()

	if !it.Next() {
		return nil
	}
	return &Item{Key: it.Key(), Value: it.Value()}
}

func (s *sequencer) Slice(scope store.Scope, from insolar.PulseNumber, skip uint32, to insolar.PulseNumber, limit uint32) []Item {
	s.RLock()
	defer s.RUnlock()

	var result []Item
	it := s.db.NewIterator(polyKey{id: from.Bytes(), scope: scope}, false)
	defer it.Close()

	skipped := 0
	for it.Next() && insolar.NewPulseNumber(it.Key()[:4]) < to && len(result) < int(limit) {
		if skipped < int(skip) {
			skipped++
			continue
		}
		result = append(result, Item{
			Key:   it.Key(),
			Value: it.Value(),
		})
	}
	return result
}

func (s *sequencer) Upsert(scope store.Scope, sequence []Item) {
	s.Lock()
	defer s.Unlock()

	for _, item := range sequence {
		err := s.db.Set(polyKey{item.Key, scope}, item.Value)
		if err != nil {
			inslogger.FromContext(context.Background()).Error(errors.Wrapf(err, "failed to save item of sequence"))
		}
	}
}

type polyKey struct {
	id    []byte
	scope store.Scope
}

func (k polyKey) ID() []byte {
	return k.id
}

func (k polyKey) Scope() store.Scope {
	return k.scope
}
