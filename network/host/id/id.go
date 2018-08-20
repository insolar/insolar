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

package id

import (
	"bytes"

	"github.com/jbenet/go-base58"
)

// ID is node id.
type ID struct {
	Key  []byte
	Hash []byte
}

// NewID returns random node id.
func NewID(key []byte) (ID, error) {
	hash := make([]byte, 20) // TODO: choose hash func
	_, err := random.Read(hash)
	id := ID{Hash: hash, Key: key}
	return id, err
}

// HashEqual checks if hash is equal to another.
func (id ID) HashEqual(other []byte) bool {
	return bytes.Equal(id.Hash, other)
}

// KeyEqual checks if id is equal ot another.
func (id ID) KeyEqual(other []byte) bool {
	return bytes.Equal(id.Key, other)
}

// KeyString is a base58-encoded string representation of node public key.
func (id ID) KeyString() string {
	return base58.Encode(id.Key)
}

// HashString is a base58-encoded string representation hash of public key.
func (id ID) HashString() string {
	return base58.Encode(id.Hash)
}
