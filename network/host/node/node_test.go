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
	"testing"

	"github.com/insolar/insolar/network/host/id"
	"github.com/stretchr/testify/assert"
)

func TestNewNode(t *testing.T) {
	addr, _ := NewAddress("127.0.0.1:31337")
	actualNode := NewNode(addr)
	id1, _ := id.NewID(nil)
	id1.SetHash(nil)
	expectedNode := &Node{
		ID:      id1,
		Address: addr,
	}

	assert.Equal(t, expectedNode, actualNode)
}

func TestNode_String(t *testing.T) {
	addr, _ := NewAddress("127.0.0.1:31337")
	nd := NewNode(addr)
	id1, _ := id.NewID(id.GetRandomKey())
	id1.SetHash([]byte{49, 50, 51, 52, 53, 54, 55, 56, 57, 48, 97, 98, 99, 100, 101, 102, 103, 104, 105, 106})
	nd.ID = id1

	assert.Equal(t, "gkdhQDvLi23xxjXjhpMWaTt5byb (127.0.0.1:31337)", nd.String())
}

func TestNode_Equal(t *testing.T) {
	id1, _ := id.NewID(id.GetRandomKey())
	id2, _ := id.NewID(id.GetRandomKey())
	idNil, _ := id.NewID(id.GetRandomKey())
	idNil.SetHash(nil)
	addr1, _ := NewAddress("127.0.0.1:31337")
	addr2, _ := NewAddress("10.10.11.11:12345")

	tests := []struct {
		id1   id.ID
		addr1 *Address
		id2   id.ID
		addr2 *Address
		equal bool
		name  string
	}{
		{id1, addr1, id1, addr1, true, "same id and address"},
		{id1, addr1, id1, addr2, false, "different addresses"},
		{id1, addr1, id2, addr1, false, "different ids"},
		{id1, addr1, id2, addr2, false, "different id and address"},
		{id1, addr1, id2, addr2, false, "different id and address"},
		{id1, nil, id1, nil, false, "nil addresses"},
		{idNil, addr1, idNil, addr1, true, "nil ids"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.equal, Node{test.id1, test.addr1}.Equal(Node{test.id2, test.addr2}))
		})
	}
}
