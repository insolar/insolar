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

	"github.com/dgraph-io/badger"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/store"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/pkg/errors"
)

// DB is a DB storage implementation. It saves pulses to disk and does not allow removal.
type DB struct {
	lock sync.RWMutex
	db   store.DB

	latest *insolar.Pulse
}

type pulseKey insolar.PulseNumber

func (k pulseKey) Scope() store.Scope {
	return store.ScopePulse
}

func (k pulseKey) ID() []byte {
	return insolar.PulseNumber(k).Bytes()
}

func newPulseKey(raw []byte) pulseKey {
	key := pulseKey(insolar.NewPulseNumber(raw))
	return key
}

type dbNode struct {
	Pulse      insolar.Pulse
	Prev, Next *insolar.PulseNumber
}

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
	/*for {
		err = s.db.Backend().View(func(txn *badger.Txn) error {
			opts := badger.DefaultIteratorOptions
			opts.Reverse = true
			pivot := pulseKey(insolar.PulseNumber(0xFFFFFFFF))
			scope := pivot.Scope().Bytes()
			prefix := append(pivot.Scope().Bytes(), pivot.ID()...)
			it := txn.NewIterator(opts)
			it.Seek(prefix)
		}

		if err == nil {
			inslogger.FromContext(ctx).Debugf("DB.Latest -  s.db.Backend().View returned an error, retrying: %s", err.Error())
			break
		}
	} */

	s.lock.RLock()
	defer s.lock.RUnlock()

	if s.latest != nil {
		return *s.latest, nil
	}

	// get head
	head, err := func() (insolar.PulseNumber, error) {
		rit := s.db.NewIterator(pulseKey(insolar.PulseNumber(0xFFFFFFFF)), true)
		defer rit.Close()

		if !rit.Next() {
			return insolar.GenesisPulse.PulseNumber, ErrNotFound
		}
		return insolar.NewPulseNumber(rit.Key()), nil
	}()

	if err != nil {
		return
	}

	// get(head)
	nd, err := func() (nd dbNode, err error) {
		key := pulseKey(head)
		// get()
		buf, err := func() (value []byte, err error) {
			fullKey := append(key.Scope().Bytes(), key.ID()...)

			err = s.db.Backend().View(func(txn *badger.Txn) error {
				item, err := txn.Get(fullKey)
				if err != nil {
					return err
				}
				value, err = item.ValueCopy(nil)
				return err
			})

			if err != nil {
				if err == badger.ErrKeyNotFound {
					return nil, ErrNotFound
				}
				return nil, err
			}

			return
		}()

		if err == store.ErrNotFound {
			err = ErrNotFound
			return
		}
		if err != nil {
			return
		}
		nd = deserialize(buf)
		return
	}()

	if err != nil {
		return
	}
	return nd.Pulse, nil
}

// TruncateHead remove all records after lastPulse
func (s *DB) TruncateHead(ctx context.Context, from insolar.PulseNumber) error {
	it := s.db.NewIterator(pulseKey(from), false)
	defer it.Close()

	s.lock.Lock()
	s.latest = nil
	s.lock.Unlock()

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
		inslogger.FromContext(ctx).Debug("No records. Nothing done. Pulse number: " + from.String())
	}

	return nil
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
		return s.set(oldHead.Pulse.PulseNumber, oldHead)
	}
	var insertWithoutHead = func() error {
		// Set new pulse.
		return s.set(pulse.PulseNumber, dbNode{
			Pulse: pulse,
		})
	}

	head, err := s.head()
	if err == ErrNotFound {
		err := insertWithoutHead()
		if err == nil {
			s.latest = &pulse
		}
		return err
	}

	if pulse.PulseNumber <= head {
		return ErrBadPulse
	}
	err = insertWithHead(head)
	if err == nil {
		s.latest = &pulse
	}
	return err
}

// Forwards calculates steps pulses forwards from provided pulse. If calculated pulse does not exist, ErrNotFound will
// be returned.
func (s *DB) Forwards(ctx context.Context, pn insolar.PulseNumber, steps int) (insolar.Pulse, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	_, err := s.db.Get(pulseKey(pn))
	if err != nil {
		return *insolar.GenesisPulse, err
	}

	it := s.db.NewIterator(pulseKey(pn), false)
	defer it.Close()
	for i := 0; it.Next(); i++ {
		if i == steps {
			buf, err := it.Value()
			if err != nil {
				return *insolar.GenesisPulse, err
			}
			nd := deserialize(buf)
			return nd.Pulse, nil
		}
	}
	return *insolar.GenesisPulse, ErrNotFound
}

// Backwards calculates steps pulses backwards from provided pulse. If calculated pulse does not exist, ErrNotFound will
// be returned.
func (s *DB) Backwards(ctx context.Context, pn insolar.PulseNumber, steps int) (insolar.Pulse, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	_, err := s.db.Get(pulseKey(pn))
	if err != nil {
		return *insolar.GenesisPulse, err
	}

	rit := s.db.NewIterator(pulseKey(pn), true)
	defer rit.Close()
	for i := 0; rit.Next(); i++ {
		if i == steps {
			buf, err := rit.Value()
			if err != nil {
				return *insolar.GenesisPulse, err
			}
			nd := deserialize(buf)
			return nd.Pulse, nil
		}
	}
	return *insolar.GenesisPulse, ErrNotFound
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

	rit := s.db.NewIterator(pulseKey(insolar.PulseNumber(0xFFFFFFFF)), true)
	defer rit.Close()

	if !rit.Next() {
		return insolar.GenesisPulse.PulseNumber, ErrNotFound
	}
	return insolar.NewPulseNumber(rit.Key()), nil
}

func serialize(nd dbNode) []byte {
	return insolar.MustSerialize(nd)
}

func deserialize(buf []byte) (nd dbNode) {
	insolar.MustDeserialize(buf, &nd)
	return nd
}
