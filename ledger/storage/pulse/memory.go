package pulse

import (
	"context"
	"sync"

	"github.com/insolar/insolar/core"
	"github.com/pkg/errors"
)

type StorageMem struct {
	lock    sync.RWMutex
	storage map[core.PulseNumber]*memNode
	head    *memNode
	tail    *memNode
}

type memNode struct {
	pulse      core.Pulse
	prev, next *memNode
}

func NewStorageMem() *StorageMem {
	return &StorageMem{
		storage: make(map[core.PulseNumber]*memNode),
	}
}

func (s *StorageMem) ForPulseNumber(ctx context.Context, pn core.PulseNumber) (pulse core.Pulse, err error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	node, ok := s.storage[pn]
	if !ok {
		err = ErrNotFound
		return
	}

	return node.pulse, nil
}

func (s *StorageMem) Latest(ctx context.Context) (pulse core.Pulse, err error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	if s.head == nil {
		err = ErrNotFound
		return
	}

	return s.head.pulse, nil
}

func (s *StorageMem) Append(ctx context.Context, pulse core.Pulse) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	var insertWithHead = func() {
		oldHead := s.head
		newHead := &memNode{
			prev:  oldHead,
			pulse: pulse,
		}
		oldHead.next = newHead
		newHead.prev = oldHead
		s.storage[newHead.pulse.PulseNumber] = newHead
		s.head = newHead
	}
	var insertWithoutHead = func() {
		s.head = &memNode{
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

func (s *StorageMem) Shift(ctx context.Context) (pulse core.Pulse, err error) {
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

func (s *StorageMem) Forwards(ctx context.Context, pn core.PulseNumber, steps int) (pulse core.Pulse, err error) {
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

func (s *StorageMem) Backwards(ctx context.Context, pn core.PulseNumber, steps int) (pulse core.Pulse, err error) {
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
