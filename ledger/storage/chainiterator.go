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

// Next returns next element.
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
