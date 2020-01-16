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

	"github.com/dgraph-io/badger"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/store"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/pkg/errors"
)

// AALEKSEEV TODO use PostgreSQL + see db_test & pulse_cmp_test

// DB is a DB storage implementation. It saves pulses to disk and does not allow removal.
type DB struct {
	db *badger.DB
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
func NewDB(db *store.BadgerDB) *DB {
	return &DB{db: db.Backend()}
}

// ForPulseNumber returns pulse for provided a pulse number. If not found, ErrNotFound will be returned.
func (s *DB) ForPulseNumber(ctx context.Context, pn insolar.PulseNumber) (retPulse insolar.Pulse, retErr error) {
	for {
		err := s.db.View(func(txn *badger.Txn) error {
			node, err := get(txn, pulseKey(pn))
			if err != nil {
				retErr = err
				return nil
			}

			retPulse = node.Pulse
			return nil
		})

		if err == nil {
			break
		}

		inslogger.FromContext(ctx).Debugf("DB.ForPulseNumber -  s.db.Backend().View returned an error, retrying: %s", err.Error())
	}
	return
}

// Latest returns a latest pulse saved in DB. If not found, ErrNotFound will be returned.
func (s *DB) Latest(ctx context.Context) (retPulse insolar.Pulse, retErr error) {
	for {
		err := s.db.View(func(txn *badger.Txn) error {
			head, err := head(txn)
			if err != nil {
				retErr = err
				return nil
			}

			node, err := get(txn, pulseKey(head))
			if err != nil {
				retErr = err
				return nil
			}

			retPulse = node.Pulse
			return nil
		})

		if err == nil {
			break
		}

		inslogger.FromContext(ctx).Debugf("DB.Latest -  s.db.Backend().View returned an error, retrying: %s", err.Error())
	}
	return
}

// TruncateHead remove all records after lastPulse
func (s *DB) TruncateHead(ctx context.Context, from insolar.PulseNumber) error {
	var hasKeys bool
	for {
		hasKeys = false
		err := s.db.Update(func(txn *badger.Txn) error {
			it := txn.NewIterator(badger.DefaultIteratorOptions)
			defer it.Close()

			pivot := pulseKey(from)
			prefix := append(pivot.Scope().Bytes(), pivot.ID()...)
			scope := pivot.Scope().Bytes()
			it.Seek(prefix)
			for {
				if !it.ValidForPrefix(scope) {
					break
				}

				hasKeys = true
				k := it.Item().KeyCopy(nil)
				loggedKey := newPulseKey(k[len(scope):])
				it.Next()

				err := txn.Delete(k)
				if err != nil {
					txn.Discard()
					return errors.Wrapf(err, "can't delete key: %+v", loggedKey)
				}

				// It's not very good to write logs from inside of the transaction, but since
				// TruncateHead() is not called often it's OK in this case.
				inslogger.FromContext(ctx).Debugf("DB.TruncateHead - Erased key with pulse number: %s", insolar.PulseNumber(loggedKey))
			}

			return nil
		})

		if err == nil {
			break
		}

		inslogger.FromContext(ctx).Debugf("DB.TruncateHead - s.db.Backend().Update returned an error, retrying: %s", err.Error())
	}

	if !hasKeys {
		inslogger.FromContext(ctx).Debug("DB.TruncateHead - No records to delete from pulse number: " + from.String())
	}

	return nil
}

// Append appends provided pulse to current storage. Pulse number should be greater than currently saved for preserving
// pulse consistency. If a provided pulse does not meet the requirements, ErrBadPulse will be returned.
func (s *DB) Append(ctx context.Context, pulse insolar.Pulse) error {
	var retErr error

	for k, s := range pulse.Signs {
		if s.PulseNumber != pulse.PulseNumber {
			return errors.New("Signatures check failed for pulse: pulse numbers mismatch")
		}

		if k != s.ChosenPublicKey {
			return errors.New("Signatures check failed for pulse: public keys mismatch")
		}
	}

	for {
		err := s.db.Update(func(txn *badger.Txn) error {
			var insertWithHead = func(head insolar.PulseNumber) error {
				oldHead, err := get(txn, pulseKey(head))
				if err != nil {
					return err
				}
				oldHead.Next = &pulse.PulseNumber

				// Set new pulse.
				err = set(txn, pulse.PulseNumber, dbNode{
					Prev:  &oldHead.Pulse.PulseNumber,
					Pulse: pulse,
				})
				if err != nil {
					return err
				}
				// Set old updated tail.
				return set(txn, oldHead.Pulse.PulseNumber, oldHead)
			}
			var insertWithoutHead = func() error {
				// Set new pulse.
				return set(txn, pulse.PulseNumber, dbNode{
					Pulse: pulse,
				})
			}

			head, err := head(txn)
			if err == ErrNotFound {
				err = insertWithoutHead()
				if err != nil {
					txn.Discard()
				}
				return err
			}

			if pulse.PulseNumber <= head {
				retErr = ErrBadPulse
				return nil
			}

			err = insertWithHead(head)
			if err != nil {
				txn.Discard()
			}
			return err
		})

		if err == nil {
			break
		}

		inslogger.FromContext(ctx).Debugf("DB.Append -  s.db.Backend().Update returned an error, retrying: %s", err.Error())
	}
	return retErr
}

// Forwards calculates steps pulses forwards from provided pulse. If calculated pulse does not exist, ErrNotFound will
// be returned.
func (s *DB) Forwards(ctx context.Context, pn insolar.PulseNumber, steps int) (insolar.Pulse, error) {
	return s.traverse(ctx, pn, steps, false)
}

// Backwards calculates steps pulses backwards from provided pulse. If calculated pulse does not exist, ErrNotFound will
// be returned.
func (s *DB) Backwards(ctx context.Context, pn insolar.PulseNumber, steps int) (insolar.Pulse, error) {
	return s.traverse(ctx, pn, steps, true)
}

func (s *DB) traverse(ctx context.Context, pn insolar.PulseNumber, steps int, reverse bool) (insolar.Pulse, error) {
	if steps < 0 {
		return *insolar.GenesisPulse, errors.New("DB.traverse - `steps` argument should be not negative")
	}

	var (
		retPulse insolar.Pulse
		retErr   error
	)
	for {
		err := s.db.View(func(txn *badger.Txn) error {
			opts := badger.DefaultIteratorOptions
			opts.Reverse = reverse
			opts.PrefetchSize = steps + 1
			it := txn.NewIterator(opts)
			defer it.Close()

			pivot := pulseKey(pn)
			prefix := append(pivot.Scope().Bytes(), pivot.ID()...)
			scope := pivot.Scope().Bytes()
			it.Seek(prefix)
			i := 0
			for {
				if !it.ValidForPrefix(scope) {
					break
				}

				if i == steps {
					buf, err := it.Item().ValueCopy(nil)
					if err != nil {
						retPulse = *insolar.GenesisPulse
						retErr = err
						return nil
					}
					node := deserialize(buf)
					retPulse = node.Pulse
					retErr = nil
					return nil
				}

				it.Next()
				i++
			}

			// not found
			retPulse = *insolar.GenesisPulse
			retErr = ErrNotFound
			return nil
		})

		if err == nil {
			break
		}

		inslogger.FromContext(ctx).Debugf("DB.traverse - s.db.Backend().View returned an error, retrying: %s", err.Error())
	}

	return retPulse, retErr
}

func head(txn *badger.Txn) (insolar.PulseNumber, error) {
	opts := badger.DefaultIteratorOptions
	opts.Reverse = true
	// we need only one last key
	opts.PrefetchSize = 1
	it := txn.NewIterator(opts)
	defer it.Close()

	pivot := pulseKey(insolar.PulseNumber(0xFFFFFFFF))
	scope := pivot.Scope().Bytes()
	prefix := append(pivot.Scope().Bytes(), pivot.ID()...)
	it.Seek(prefix)
	if !it.ValidForPrefix(scope) {
		return insolar.GenesisPulse.PulseNumber, ErrNotFound
	}

	k := it.Item().KeyCopy(nil)
	return insolar.NewPulseNumber(k[len(scope):]), nil
}

func get(txn *badger.Txn, key pulseKey) (retNode dbNode, retErr error) {
	fullKey := append(key.Scope().Bytes(), key.ID()...)
	item, err := txn.Get(fullKey)
	if err != nil {
		if err == badger.ErrKeyNotFound {
			err = ErrNotFound
		}
		retErr = err
		return
	}
	buf, err := item.ValueCopy(nil)
	if err != nil {
		retErr = err
		return
	}

	retNode = deserialize(buf)
	return
}

func set(txn *badger.Txn, pn insolar.PulseNumber, node dbNode) error {
	key := pulseKey(pn)
	fullKey := append(key.Scope().Bytes(), key.ID()...)
	return txn.Set(fullKey, serialize(node))
}

func serialize(nd dbNode) []byte {
	return insolar.MustSerialize(nd)
}

func deserialize(buf []byte) (nd dbNode) {
	insolar.MustDeserialize(buf, &nd)
	return nd
}
