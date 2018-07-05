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

package node

import (
	"bytes"

	"github.com/jbenet/go-base58"
)

// ID is node id
type ID []byte

// NewID returns random node id
// TODO: Should test errors produced here
func NewID() (ID, error) {
	result := make([]byte, 20)
	_, err := random.Read(result)
	return result, err
}

// NewIDs returns given number of random node ids
func NewIDs(num int) ([]ID, error) {
	result := make([]ID, num)

	for i := range result {
		id, err := NewID()

		if err != nil {
			return nil, err
		}

		result[i] = id
	}

	return result, nil
}

// Equal checks if id is equal to another
func (id ID) Equal(other ID) bool {
	return bytes.Equal(id, other)
}

// String is a base58-encoded string representation of node id
func (id ID) String() string {
	return base58.Encode(id)
}
