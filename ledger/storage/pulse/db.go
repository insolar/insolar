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

package pulse

import (
	"context"
	"sync"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/ledger/storage/db"
	"github.com/pkg/errors"
)

type dbStorage struct {
	lock    sync.RWMutex
	storage map[core.PulseNumber]*memoryNode
	head    *memoryNode
	tail    *memoryNode
}

type dbKey core.PulseNumber

func (k dbKey) Scope() db.Scope {
	return db.ScopePulse
}

func (k dbKey) Key() []byte {
	return core.PulseNumber(k).Bytes()
}

type dbNode struct {
	pulse      core.Pulse
	prev, next dbKey
}

func NewDBStorage() *dbStorage {
	return &dbStorage{
		storage: make(map[core.PulseNumber]*memoryNode),
	}
}

func (s *dbStorage) ForPulseNumber(ctx context.Context, pn core.PulseNumber) (pulse core.Pulse, err error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	node, ok := s.storage[pn]
	if !ok {
		err = core.ErrNotFound
		return
	}

	return node.pulse, nil
}

func (s *dbStorage) Latest(ctx context.Context) (pulse core.Pulse, err error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	if s.head == nil {
		err = core.ErrNotFound
		return
	}

	return s.head.pulse, nil
}

func (s *dbStorage) Append(ctx context.Context, pulse core.Pulse) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	var insertWithHead = func() {
		oldHead := s.head
		newHead := &memoryNode{
			prev:  oldHead,
			pulse: pulse,
		}
		oldHead.next = newHead
		newHead.prev = oldHead
		s.storage[newHead.pulse.PulseNumber] = newHead
		s.head = newHead
	}
	var insertWithoutHead = func() {
		s.head = &memoryNode{
			pulse: pulse,
		}
		s.storage[pulse.PulseNumber] = s.head
		s.tail = s.head
	}

	if s.head == nil {
		insertWithoutHead()
		return nil
	}

	if pulse.PulseNumber <= s.head.pulse.PulseNumber {
		return errors.New("pulse should be greater than the latest")
	}
	insertWithHead()

	return nil
}

func (s *dbStorage) Shift(ctx context.Context) (pulse core.Pulse, err error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	if s.tail == nil {
		err = errors.New("nothing to shift")
		return
	}

	delete(s.storage, s.tail.pulse.PulseNumber)
	if s.tail == s.head {
		tail := s.tail
		s.tail, s.head = nil, nil
		return tail.pulse, nil
	}

	tail := s.tail
	tail.next.prev = nil
	s.tail = tail.next
	return tail.pulse, nil
}

func (s *dbStorage) Forwards(ctx context.Context, pn core.PulseNumber, steps int) (pulse core.Pulse, err error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	node, ok := s.storage[pn]
	if !ok {
		err = core.ErrNotFound
		return
	}

	iterator := node
	for i := 0; i < steps; i++ {
		if iterator.next == nil {
			err = core.ErrNotFound
			return
		}
		iterator = iterator.next
	}

	return iterator.pulse, nil
}

func (s *dbStorage) Backwards(ctx context.Context, pn core.PulseNumber, steps int) (pulse core.Pulse, err error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	node, ok := s.storage[pn]
	if !ok {
		err = core.ErrNotFound
		return
	}

	iterator := node
	for i := 0; i < steps; i++ {
		if iterator.prev == nil {
			err = core.ErrNotFound
			return
		}
		iterator = iterator.prev
	}

	return iterator.pulse, nil
}
