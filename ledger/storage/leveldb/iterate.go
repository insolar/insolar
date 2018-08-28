/*
 *    Copyright 2018 INS Ecosystem
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

package leveldb

import (
	"bytes"

	"github.com/syndtr/goleveldb/leveldb/iterator"
	"github.com/syndtr/goleveldb/leveldb/util"

	"github.com/insolar/insolar/ledger/record"
)

// HashIterator iterates overs a DB's keys.
// An iterator provides methods for record hash extraction.
//
// HashIterator supposed to be used only in functions like ProcessSlotRecords.
// Any release of iterator's resources like in native leveldb iterators
// (see github.com/syndtr/goleveldb/leveldb/iterator) not needed.
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

// ProcessSlotRecords executes a iteration function and provides HashIterator
// to iterate all records with the same record.PulseNum.
//
// Error returned by the ProcessSlotRecords is based on iteration function
// result or leveldb iterator error if any.
func (ll *LevelLedger) ProcessSlotRecords(n record.PulseNum, ifn func(it HashIterator) error) error {
	prefix := make([]byte, record.PulseNumSize+1)
	prefix[0] = scopeIDRecord
	buf := bytes.NewBuffer(prefix[1:1])
	n.MustWrite(buf)

	ldbIter := ll.ldb.NewIterator(util.BytesPrefix(prefix), nil)
	defer ldbIter.Release()

	err := ifn(&iter{i: ldbIter})
	if err != nil {
		return err
	}
	return ldbIter.Error()
}

// GetSlotHashes returns array of all record's hashes in provided PulseNum.
func (ll *LevelLedger) GetSlotHashes(n record.PulseNum) ([][]byte, error) {
	var hashes [][]byte
	err := ll.ProcessSlotRecords(n, func(it HashIterator) error {
		for i := 1; it.Next(); i++ {
			hashes = append(hashes, it.Hash())
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return hashes, nil
}

// levelDB iterator wrapper code
type iter struct {
	i iterator.Iterator
}

func (it *iter) Next() bool {
	return it.i.Next()
}

func (it *iter) Hash() []byte {
	key := it.i.Key()
	hash := make([]byte, len(key)-1)
	_ = copy(hash, key[1:])
	return hash
}

func (it *iter) ShallowHash() []byte {
	return it.i.Key()[1:]
}
