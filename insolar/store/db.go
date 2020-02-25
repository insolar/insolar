// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package store

import (
	"io"
)

//go:generate minimock -i github.com/insolar/insolar/insolar/store.DB -o ./ -s _gen_mock.go -g

// DB provides a simple key-value store interface for persisting data.
// But it is internally ordered ( lexicographically by key bytes )
// so if you want you can iterate over store using Iterator interface.
type DB interface {
	// Backend returns the underlying badger.DB object. Use with care.
	Get(key Key) (value []byte, err error)
	Set(key Key, value []byte) error
	Delete(key Key) error
	NewIterator(pivot Key, reverse bool) Iterator
}

// Backuper provides interface for making backups
type Backuper interface {
	// Backup does incremental backup starting from 'since' timestamp and write result to 'to' parameter.
	// It returns a timestamp indicating when the entries were dumped which can be passed into a
	// later invocation to generate an incremental dump.
	Backup(to io.Writer, since uint64) (uint64, error)
}

//go:generate minimock -i github.com/insolar/insolar/insolar/store.Iterator -o ./ -s _gen_mock.go -g

// Iterator provides an interface for walking through the storage record sequence (where records are sorted lexicographically).
type Iterator interface {
	// Next moves the iterator to the next key-value pair.
	Next() bool
	// Close frees resources within the iterator and invalidates it.
	Close()
	// Key returns only the second part of the composite key - (ID) without scope id.
	// Warning: Key is only valid as long as item is valid (until iterator.Next() called), or transaction is valid.
	// If you need to use it outside its validity, please copy the key.
	Key() []byte
	// Value returns value itself (ex: record, drop, blob, etc).
	// Warning: Value is only valid as long as item is valid (until iterator.Next() called), or transaction is valid.
	// If you need to use it outside its validity, please copy the value.
	Value() ([]byte, error)
}

// Key represents a key for the key-value store. Scope is required to separate different DB clients and should be
// unique.
type Key interface {
	// Scope returns a first part for constructing a composite key for storing record in db
	Scope() Scope
	// ID returns a second part for constructing a composite key for storing record in db
	ID() []byte
}

//go:generate stringer -type=Scope

// Scope separates DB clients.
type Scope byte

// Bytes returns binary scope representation.
func (s Scope) Bytes() []byte {
	return []byte{byte(s)}
}

const (
	// ScopePulse is the scope for pulse storage.
	ScopePulse Scope = 1
	// ScopeRecord is the scope for record storage.
	ScopeRecord Scope = 2
	// ScopeJetDrop is the scope for a jet drop storage.
	ScopeJetDrop Scope = 3
	// ScopeIndex is the scope for an index records.
	ScopeIndex Scope = 4
	// ScopeLastKnownIndexPN is the scope for a last known pulse number of the index bucket
	ScopeLastKnownIndexPN Scope = 5
	// ScopeGenesis is the scope for a genesis records.
	ScopeGenesis Scope = 6
	// ScopeJetTree is the scope for a jet tree storage.
	ScopeJetTree Scope = 7
	// ScopeJetKeeper is the scope for a jet id storage.
	ScopeJetKeeper Scope = 8
	// ScopeJetKeeperSyncPulse is the scope for a top sync pulse storage.
	ScopeJetKeeperSyncPulse Scope = 9
	// ScopeRecordPosition is the scope for records' positions.
	ScopeRecordPosition Scope = 10
	// ScopeBackupStart is the scope for backup starts.
	ScopeBackupStart Scope = 11
	// ScopeDBInit is scope for one key which means db is initialized.
	ScopeDBInit Scope = 12
	// ScopeNodeHistory is scope for list of nodes for every pulse
	ScopeNodeHistory Scope = 13
)
