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
	"fmt"
)

type TransportUnique int32

type TransportBoundID struct {
	id     ID
	unique TransportUnique
}

func NewTransportBoundID(unique TransportUnique) (*TransportBoundID, error) {
	id, err := NewID()
	if err != nil {
		return nil, err
	}
	return &TransportBoundID{
		id:     id,
		unique: unique,
	}, nil
}

// Equal checks if id is equal to another.
func (id TransportBoundID) Equal(other TransportBoundID) bool {
	return id.unique == other.unique && id.id.Equal(other.id)
}

// String is a base58-encoded string representation of node id.
func (id TransportBoundID) String() string {
	return fmt.Sprintf("%s (%d)", id.id.String(), id.unique)
}
