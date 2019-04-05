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

package pulse

import (
	"context"
	"sync"

	"github.com/insolar/insolar/insolar"
	"github.com/pkg/errors"
)

// StorageMem is a memory storage implementation. It saves pulses to memory and allows removal.
type StorageMem struct {
	lock    sync.RWMutex
	storage map[insolar.PulseNumber]*memNode
	head    *memNode
	tail    *memNode
}

type memNode struct {
	pulse      insolar.Pulse
	prev, next *memNode
}

// NewStorageMem creates new memory storage instance.
func NewStorageMem() *StorageMem {
	return &StorageMem{
		storage: make(map[insolar.PulseNumber]*memNode),
	}
}

// ForPulseNumber returns pulse for provided Pulse number. If not found, ErrNotFound will be returned.
func (s *StorageMem) ForPulseNumber(ctx context.Context, pn insolar.PulseNumber) (pulse insolar.Pulse, err error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	node, ok := s.storage[pn]
	if !ok {
		err = ErrNotFound
		return
	}

	return node.pulse, nil
}

// Latest returns a latest pulse saved in memory. If not found, ErrNotFound will be returned.
func (s *StorageMem) Latest(ctx context.Context) (pulse insolar.Pulse, err error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	if s.tail == nil {
		err = ErrNotFound
		return
	}

	return s.tail.pulse, nil
}

// Append appends provided a pulse to current storage. Pulse number should be greater than currently saved for preserving
// pulse consistency. If provided Pulse does not meet the requirements, ErrBadPulse will be returned.
func (s *StorageMem) Append(ctx context.Context, pulse insolar.Pulse) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	var appendTail = func() {
		oldTail := s.tail
		newTail := &memNode{
			prev:  oldTail,
			pulse: pulse,
		}
		oldTail.next = newTail
		newTail.prev = oldTail
		s.storage[newTail.pulse.PulseNumber] = newTail
		s.tail = newTail
	}
	var appendHead = func() {
		s.tail = &memNode{
			pulse: pulse,
		}
		s.storage[pulse.PulseNumber] = s.tail
		s.head = s.tail
	}

	if s.head == nil {
		appendHead()
		return nil
	}

	if pulse.PulseNumber <= s.tail.pulse.PulseNumber {
		return ErrBadPulse
	}
	appendTail()

	return nil
}

// Shift removes youngest pulse from storage. If the storage is empty, an error will be returned.
func (s *StorageMem) Shift(ctx context.Context, pn insolar.PulseNumber) (err error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	if s.head == nil {
		err = errors.New("nothing to shift")
		return
	}

	h := s.head
	for h != nil && h.pulse.PulseNumber <= pn {
		delete(s.storage, h.pulse.PulseNumber)
		h = h.next
	}

	s.head = h
	if s.head == nil {
		s.tail = nil
	} else {
		s.head.prev = nil
	}

	return nil
}

// Forwards calculates steps pulses forwards from provided Pulse. If calculated pulse does not exist, ErrNotFound will
// be returned.
func (s *StorageMem) Forwards(ctx context.Context, pn insolar.PulseNumber, steps int) (pulse insolar.Pulse, err error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	node, ok := s.storage[pn]
	if !ok {
		err = ErrNotFound
		return
	}

	iterator := node
	for i := 0; i < steps; i++ {
		if iterator.next == nil {
			err = ErrNotFound
			return
		}
		iterator = iterator.next
	}

	return iterator.pulse, nil
}

// Backwards calculates steps pulses backwards from provided pulse. If calculated pulse does not exist, ErrNotFound will
// be returned.
func (s *StorageMem) Backwards(ctx context.Context, pn insolar.PulseNumber, steps int) (pulse insolar.Pulse, err error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	node, ok := s.storage[pn]
	if !ok {
		err = ErrNotFound
		return
	}

	iterator := node
	for i := 0; i < steps; i++ {
		if iterator.prev == nil {
			err = ErrNotFound
			return
		}
		iterator = iterator.prev
	}

	return iterator.pulse, nil
}
