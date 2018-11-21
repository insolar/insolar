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

package id

import (
	"bytes"

	"github.com/jbenet/go-base58"
)

// ID is host id.
type ID []byte

// NewID returns random host id.
func NewID() (ID, error) {
	key := make([]byte, 20)
	_, err := random.Read(key)
	id := ID(key)
	return id, err
}

// FromBase58 returns decoded host id.
func FromBase58(encoded string) ID {
	return ID(base58.Decode(encoded))
}

// Equal checks if id is equal ot another.
func (id ID) Equal(other []byte) bool {
	return bytes.Equal(id, other)
}

// String is a base58-encoded string representation of host public key.
func (id ID) String() string {
	return base58.Encode(id)
}

// Bytes returns a raw key.
func (id ID) Bytes() []byte {
	return id
}
