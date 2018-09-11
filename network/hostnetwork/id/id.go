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
	"fmt"

	"github.com/jbenet/go-base58"
)

// ID is host id.
type ID struct {
	key []byte
}

// MarshalBinary is binary marshaler.
func (id ID) MarshalBinary() ([]byte, error) {
	var res bytes.Buffer
	key := base58.Encode(id.key)
	fmt.Fprintln(&res, key)
	return res.Bytes(), nil
}

// UnmarshalBinary is binary unmarshaler.
func (id *ID) UnmarshalBinary(data []byte) error {
	res := bytes.NewBuffer(data)
	var key string
	_, err := fmt.Fscanln(res, &key)
	id.key = base58.Decode(key)
	return err
}

// NewID returns random host id.
func NewID() (ID, error) {
	key := make([]byte, 20) // TODO: choose hash func
	_, err := random.Read(key)
	id := ID{key: key}
	return id, err
}

// KeyEqual checks if id is equal ot another.
func (id ID) KeyEqual(other []byte) bool {
	return bytes.Equal(id.key, other)
}

// KeyString is a base58-encoded string representation of host public key.
func (id ID) KeyString() string {
	return base58.Encode(id.key)
}

// GetKey returns a raw key.
func (id ID) GetKey() []byte {
	return id.key
}
