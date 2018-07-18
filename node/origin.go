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

// Origin is “self” variant of Node.
// Unlike ordinary node it can have multiple IDs.
type Origin struct {
	IDs     []ID
	Address *Address
}

// NewOrigin creates origin node from list of ids and network address.
func NewOrigin(ids []ID, address *Address) (*Origin, error) {
	var err error

	if len(ids) == 0 {
		ids, err = NewIDs(1)
	}

	if err != nil {
		return nil, err
	}

	return &Origin{
		IDs:     ids,
		Address: address,
	}, nil
}

func (s *Origin) containsID(id ID) bool {
	for _, myID := range s.IDs {
		if id.Equal(myID) {
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
