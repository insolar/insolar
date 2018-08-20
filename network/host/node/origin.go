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
	"github.com/insolar/insolar/network/host/id"
)

// Origin is “self” variant of Node.
// Unlike ordinary node it can have multiple IDs.
type Origin struct {
	IDs     []id.ID
	Address *Address
}

// NewOrigin creates origin node from list of ids and network address.
func NewOrigin(ids []id.ID, address *Address) (*Origin, error) {
	if len(ids) == 0 {
		id1, err := id.NewID([]byte{49, 50, 51, 52, 53, 54, 55, 56, 57, 48, 97, 98, 99, 100, 101, 102, 103, 104, 105, 106})
		ids = append(ids, id1)

		if err != nil {
			return nil, err
		}
	}

	return &Origin{
		IDs:     ids,
		Address: address,
	}, nil
}

func (s *Origin) containsID(id id.ID) bool {
	for _, myID := range s.IDs {
		if id.HashEqual(myID.Hash) {
			return true
		}
	}
	return false
}

// Contains checks if origin node “contains” network node.
// It checks if node's and origin's addresses match and node's id is in origin's ids list.
func (s *Origin) Contains(node *Node) bool {
	return node.Address.Equal(*s.Address) && s.containsID(node.ID)
}
