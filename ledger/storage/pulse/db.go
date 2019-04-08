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
	"bytes"
	"context"
	"sync"

	"github.com/ugorji/go/codec"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/internal/ledger/store"
)

// DB is a DB storage implementation. It saves pulses to disk and does not allow removal.
type DB struct {
	db   store.DB
	lock sync.RWMutex
}

type pulseKey insolar.PulseNumber

func (k pulseKey) Scope() store.Scope {
	return store.ScopePulse
}

func (k pulseKey) ID() []byte {
	return append([]byte{prefixPulse}, insolar.PulseNumber(k).Bytes()...)
}

type metaKey byte

func (k metaKey) Scope() store.Scope {
	return store.ScopePulse
}

func (k metaKey) ID() []byte {
	return []byte{prefixMeta, byte(k)}
}

type dbNode struct {
	Pulse      insolar.Pulse
	Prev, Next *insolar.PulseNumber
}

var (
	prefixPulse byte = 1
	prefixMeta  byte = 2
)

var (
	keyHead metaKey = 1
)

// NewDB creates new DB storage instance.
func NewDB(db store.DB) *DB {
	return &DB{db: db}
}

// ForPulseNumber returns pulse for provided a pulse number. If not found, ErrNotFound will be returned.
func (s *DB) ForPulseNumber(ctx context.Context, pn insolar.PulseNumber) (pulse insolar.Pulse, err error) {
	nd, err := s.get(pn)
	if err != nil {
		return
	}
	return nd.Pulse, nil
}

// Latest returns a latest pulse saved in DB. If not found, ErrNotFound will be returned.
func (s *DB) Latest(ctx context.Context) (pulse insolar.Pulse, err error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	head, err := s.head()
	if err != nil {
		return
	}
	nd, err := s.get(head)
	if err != nil {
		return
	}
	return nd.Pulse, nil
}

// Append appends provided pulse to current storage. Pulse number should be greater than currently saved for preserving
// pulse consistency. If a provided pulse does not meet the requirements, ErrBadPulse will be returned.
func (s *DB) Append(ctx context.Context, pulse insolar.Pulse) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	var insertWithHead = func(head insolar.PulseNumber) error {
		oldHead, err := s.get(head)
		if err != nil {
			return err
		}
		oldHead.Next = &pulse.PulseNumber

		// Set new pulse.
		err = s.set(pulse.PulseNumber, dbNode{
			Prev:  &oldHead.Pulse.PulseNumber,
			Pulse: pulse,
		})
		if err != nil {
			return err
		}
		// Set old updated tail.
		err = s.set(oldHead.Pulse.PulseNumber, oldHead)
		if err != nil {
			return err
		}
		// Set head meta record.
		return s.setHead(pulse.PulseNumber)
	}
	var insertWithoutHead = func() error {
		// Set new pulse.
		err := s.set(pulse.PulseNumber, dbNode{
			Pulse: pulse,
		})
		if err != nil {
			return err
		}
		// Set head meta record.
		return s.setHead(pulse.PulseNumber)
	}

	head, err := s.head()
	if err == ErrNotFound {
		return insertWithoutHead()
	}

	if pulse.PulseNumber <= head {
		return ErrBadPulse
	}
	return insertWithHead(head)
}

// Forwards calculates steps pulses forwards from provided pulse. If calculated pulse does not exist, ErrNotFound will
// be returned.
func (s *DB) Forwards(ctx context.Context, pn insolar.PulseNumber, steps int) (pulse insolar.Pulse, err error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	node, err := s.get(pn)
	if err != nil {
		return
	}

	iterator := node
	for i := 0; i < steps; i++ {
		if iterator.Next == nil {
			err = ErrNotFound
			return
		}
		iterator, err = s.get(*iterator.Next)
		if err != nil {
			return
		}
	}

	return iterator.Pulse, nil
}

// Backwards calculates steps pulses backwards from provided pulse. If calculated pulse does not exist, ErrNotFound will
// be returned.
func (s *DB) Backwards(ctx context.Context, pn insolar.PulseNumber, steps int) (pulse insolar.Pulse, err error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	node, err := s.get(pn)
	if err != nil {
		return
	}

	iterator := node
	for i := 0; i < steps; i++ {
		if iterator.Prev == nil {
			err = ErrNotFound
			return
		}
		iterator, err = s.get(*iterator.Prev)
		if err != nil {
			return
		}
	}

	return iterator.Pulse, nil
}

func (s *DB) get(pn insolar.PulseNumber) (nd dbNode, err error) {
	buf, err := s.db.Get(pulseKey(pn))
	if err == store.ErrNotFound {
		err = ErrNotFound
		return
	}
	if err != nil {
		return
	}
	nd = deserialize(buf)
	return
}

func (s *DB) set(pn insolar.PulseNumber, nd dbNode) error {
	return s.db.Set(pulseKey(pn), serialize(nd))
}

func (s *DB) head() (pn insolar.PulseNumber, err error) {
	buf, err := s.db.Get(keyHead)
	if err == store.ErrNotFound {
		err = ErrNotFound
		return
	}
	if err != nil {
		return
	}
	pn = insolar.NewPulseNumber(buf)
	return
}

func (s *DB) setHead(pn insolar.PulseNumber) error {
	return s.db.Set(keyHead, pn.Bytes())
}

func serialize(nd dbNode) []byte {
	buff := bytes.NewBuffer(nil)
	enc := codec.NewEncoder(buff, &codec.CborHandle{})
	enc.MustEncode(nd)
	return buff.Bytes()
}

func deserialize(buf []byte) (nd dbNode) {
	dec := codec.NewDecoderBytes(buf, &codec.CborHandle{})
	dec.MustDecode(&nd)
	return nd
}
