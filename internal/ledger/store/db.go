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

package store

//go:generate minimock -i github.com/insolar/insolar/internal/ledger/store.DB -o ./ -s _gen_mock.go

// DB provides a simple key-value store interface for persisting data.
// But it is internally ordered ( lexicographically by key bytes )
// so if you want you can iterate over store using Iterator interface.
type DB interface {
	Get(key Key) (value []byte, err error)
	Set(key Key, value []byte) error
	NewIterator(scope Scope) Iterator
}

// Iterator provides an interface for walking through the storage record sequence (where records are sorted lexicographically).
type Iterator interface {
	// Seek moves the iterator to storage record that starts with prefix bytes.
	Seek(prefix []byte)
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

	// ScopeBlob is the scope for a blobs records.
	ScopeBlob Scope = 7

	// ScopeGenesis is the scope for a genesis records.
	ScopeGenesis Scope = 8
)
