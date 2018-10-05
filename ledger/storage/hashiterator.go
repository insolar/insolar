/*
 *    Copyright 2018 Insolar
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

package storage

import (
	"github.com/dgraph-io/badger"
	"github.com/insolar/insolar/core"
)

// HashIterator iterates over a database record's hashes.
// An iterator provides methods for record's hash access.
//
// HashIterator supposed to be used only in functions like ProcessSlotHashes.
// Any release of iterator resources not needed.
type HashIterator interface {
	// Next moves the iterator to the next key/value pair.
	// It returns false then iterator is exhausted.
	Next() bool

	// Hash returns record's hash copy. That allows use returned value
	// on any iteration step or outside of iteration function.
	Hash() []byte

	// ShallowHash returns unsafe record's hash, that could be used only
	// in current iteration step. It could be useful for processing hashes
	// on the fly to avoid unnecessary copy and memory allocations.
	ShallowHash() []byte
}

// Iterator iterates over a database records.
// An iterator provides methods for record access hashes and values
type Iterator interface {
	HashIterator

	Value() ([]byte, error)
}

func pulseNumRecordPrefix(pulse core.PulseNumber) []byte {
	prefix := make([]byte, core.PulseNumberSize+1)
	prefix[0] = scopeIDRecord
	copy(prefix[1:], pulse.Bytes())
	return prefix
}

// ProcessSlotHashes executes a iteration function ifn and provides HashIterator
// inside it to iterate over all records hashes with the same record.PulseNum.
//
// Error returned by the ProcessSlotRecords is based on iteration function
// result or BadgerDB iterator error if any.
func (db *DB) ProcessSlotHashes(n core.PulseNumber, ifn func(it HashIterator) error) error {
	iopts := badger.DefaultIteratorOptions
	iopts.PrefetchValues = false
	return db.processPulse(n, iopts, func(it Iterator) error { return ifn(it.(HashIterator)) })
}

// ProcessSlot executes a iteration function ifn and provides Iterator
// inside it to iterate over all records with the same record.PulseNum.
//
// Error returned is based on iteration function result or BadgerDB
// iterator error if any.
func (db *DB) ProcessSlot(n core.PulseNumber, ifn func(it Iterator) error) error {
	iopts := badger.DefaultIteratorOptions
	iopts.PrefetchValues = true
	return db.processPulse(n, iopts, ifn)
}

func (db *DB) processPulse(n core.PulseNumber, iopts badger.IteratorOptions, ifn func(it Iterator) error) error {
	prefix := pulseNumRecordPrefix(n)

	// TODO: add transaction conflict processing
	return db.db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(iopts)
		it.Seek(prefix)
		defer it.Close()
		return ifn(&iter{i: it, prefix: prefix})
	})
}

// GetSlotHashes returns array of all record's hashes in provided PulseNum.
func (db *DB) GetSlotHashes(n core.PulseNumber) ([][]byte, error) {
	var hashes [][]byte
	err := db.ProcessSlotHashes(n, func(it HashIterator) error {
		for it.Next() {
			hashes = append(hashes, it.Hash())
		}
		return nil
	})
	if err != nil {
		hashes = nil
	}
	return hashes, err
}

// iter is a BadgerDB's iterator wrapper code.
type iter struct {
	i       *badger.Iterator
	started bool
	prefix  []byte
}

func (it *iter) valid() bool {
	return it.i.Valid() && it.i.ValidForPrefix(it.prefix)
}

func (it *iter) Next() bool {
	if it.started {
		it.i.Next()
	}
	it.started = true
	return it.valid()
}

func (it *iter) Hash() []byte {
	item := it.i.Item()
	key := item.Key()
	hash := make([]byte, len(key)-1)
	_ = copy(hash, key[1:])
	return hash
}

func (it *iter) ShallowHash() []byte {
	item := it.i.Item()
	return item.Key()[1:]
}

func (it *iter) Value() ([]byte, error) {
	item := it.i.Item()
	return item.ValueCopy(nil)
}
