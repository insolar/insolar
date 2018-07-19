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

package store

import (
	"time"
)

// Store is the interface for implementing the storage mechanism for the
// DHT.
type Store interface {
	// Store should store a key/value pair for the local node with the
	// given replication and expiration times.
	Store(key Key, data []byte, replication time.Time, expiration time.Time, publisher bool) error

	// Retrieve should return the local key/value if it exists.
	Retrieve(key Key) (data []byte, found bool)

	// Delete should delete a key/value pair from the Store.
	Delete(key Key)

	// GetKeysReadyToReplicate should return the keys of all data to be
	// replicated across the insolar. Typically all data should be
	// replicated every tReplicate seconds.
	GetKeysReadyToReplicate() []Key

	// ExpireKeys should expire all key/values due for expiration.
	ExpireKeys()
}

// NewStore creates new memory store.
func NewStore() Store {
	return NewMemoryStore()
}
