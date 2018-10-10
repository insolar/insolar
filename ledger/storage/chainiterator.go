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
	"errors"

	"github.com/insolar/insolar/ledger/record"
)

// ChainRecord is an interface for iterable records.
type ChainRecord interface {
	Next() *record.ID
}

// ChainIterator iterates over objects children.
type ChainIterator struct {
	db      *DB
	current *record.ID
}

// NewChainIterator creates new record iterator.
func NewChainIterator(db *DB, from *record.ID) *ChainIterator {
	return &ChainIterator{
		db:      db,
		current: from,
	}
}

// HasNext checks if any elements left in iterator.
func (i *ChainIterator) HasNext() bool {
	return i.current != nil
}

// Next returns element and fetches ref for the next one.
func (i *ChainIterator) Next() (*record.ID, ChainRecord, error) {
	id := i.current
	rec, err := i.db.GetRecord(id)
	if err != nil {
		return nil, nil, err
	}
	iterable, ok := rec.(ChainRecord)
	if !ok {
		return nil, nil, errors.New("wrong record type")
	}

	i.current = iterable.Next()
	return id, iterable, nil
}
