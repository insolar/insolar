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
	"fmt"

	"github.com/jbenet/go-base58"
)

// ID is node id.
type ID struct {
	key  []byte
	hash []byte
}

// GetRandomKey generates and returns a random key for ID.
func GetRandomKey() []byte {
	key := make([]byte, 20)
	_, _ = random.Read(key)
	return key
}

// GetHash returns hash of key.
func (id ID) GetHash() []byte {
	return id.hash
}

// SetHash sets new hash.
func (id *ID) SetHash(newHash []byte) {
	id.hash = newHash
}

// MarshalBinary is binary marshaler.
func (id ID) MarshalBinary() ([]byte, error) {
	var res bytes.Buffer
	key := base58.Encode(id.key)
	hash := base58.Encode(id.hash)
	fmt.Fprintln(&res, key, hash)
	return res.Bytes(), nil
}

// UnmarshalBinary is binary unmarshaler.
func (id *ID) UnmarshalBinary(data []byte) error {
	res := bytes.NewBuffer(data)
	var key string
	var hash string
	_, err := fmt.Fscanln(res, &key, &hash)
	id.key = base58.Decode(key)
	id.hash = base58.Decode(hash)
	return err
}

// NewID returns random node id.
func NewID(key []byte) (ID, error) {
	hash := make([]byte, 20) // TODO: choose hash func
	_, err := random.Read(hash)
	id := ID{hash: hash, key: key}
	return id, err
}

// HashEqual checks if hash is equal to another.
func (id ID) HashEqual(other []byte) bool {
	return bytes.Equal(id.hash, other)
}

// KeyEqual checks if id is equal ot another.
func (id ID) KeyEqual(other []byte) bool {
	return bytes.Equal(id.key, other)
}

// KeyString is a base58-encoded string representation of node public key.
func (id ID) KeyString() string {
	return base58.Encode(id.key)
}

// HashString is a base58-encoded string representation hash of public key.
func (id ID) HashString() string {
	return base58.Encode(id.hash)
}
