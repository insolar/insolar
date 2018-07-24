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

package routing

import (
	"sort"
	"testing"

	"github.com/insolar/insolar/network/host/node"

	"github.com/stretchr/testify/assert"
)

func createRouteNode(addrString string) *RouteNode {
	addr, _ := node.NewAddress(addrString)
	newNode := node.NewNode(addr)
	newNode.ID, _ = node.NewID()
	return NewRouteNode(newNode)
}

func TestNewRouteSet(t *testing.T) {
	rs := NewRouteSet()

	assert.Equal(t, &RouteSet{}, rs)
	assert.Implements(t, (*sort.Interface)(nil), rs)
}

func TestRouteSet_Nodes(t *testing.T) {
	rs := NewRouteSet()

	assert.Empty(t, rs.Nodes())

	node1 := createRouteNode("127.0.0.1:31337")
	node2 := createRouteNode("10.10.11.11:12345")

	rs.Append(node1)
	rs.Append(node2)

	assert.Equal(t, []*node.Node{node1.Node, node2.Node}, rs.Nodes())
	assert.Equal(t, rs.nodes, rs.Nodes())
}

func TestRouteSet_Nodes_ReturnsCopy(t *testing.T) {
	rs := NewRouteSet()
	rs.Append(createRouteNode("127.0.0.1:31337"))
	rs.Append(createRouteNode("10.10.11.11:12345"))

	nodesCopy := rs.Nodes()
	nodesCopy[0] = nil

	assert.NotNil(t, rs.nodes[0])
}

func TestRouteSet_FirstNode(t *testing.T) {
	rs := NewRouteSet()
	node1 := createRouteNode("127.0.0.1:31337")
	node2 := createRouteNode("10.10.11.11:12345")
	rs.Append(node1)
	rs.Append(node2)

	assert.Equal(t, node1.Node, rs.FirstNode())
}

func TestRouteSet_Contains(t *testing.T) {
	rs := NewRouteSet()
	node1 := createRouteNode("127.0.0.1:31337")
	node2 := createRouteNode("10.10.11.11:12345")
	node3 := createRouteNode("192.168.1.1:13666")
	rs.Append(node1)
	rs.Append(node2)

	assert.True(t, rs.Contains(node1))
	assert.True(t, rs.Contains(node2))
	assert.False(t, rs.Contains(node3))
}

func TestRouteSet_Append(t *testing.T) {
	rs := NewRouteSet()

	assert.Empty(t, rs.nodes)

	node1 := createRouteNode("127.0.0.1:31337")
	node2 := createRouteNode("10.10.11.11:12345")

	rs.Append(node1)
	rs.Append(node2)

	assert.Equal(t, []*node.Node{node1.Node, node2.Node}, rs.nodes)
}

func TestRouteSet_Remove(t *testing.T) {
	rs := NewRouteSet()
	node1 := createRouteNode("127.0.0.1:31337")
	node2 := createRouteNode("10.10.11.11:12345")
	node3 := createRouteNode("192.168.1.1:13666")
	rs.Append(node1)
	rs.Append(node2)
	rs.Append(node3)

	assert.True(t, rs.Contains(node1))
	assert.True(t, rs.Contains(node2))
	assert.True(t, rs.Contains(node3))

	rs.Remove(node2)

	assert.True(t, rs.Contains(node1))
	assert.False(t, rs.Contains(node2))
	assert.True(t, rs.Contains(node3))
}

func TestRouteSet_RemoveMany(t *testing.T) {
	rs := NewRouteSet()
	var nodes []*RouteNode
	nodes = append(nodes, createRouteNode("127.0.0.1:31337"))
	nodes = append(nodes, createRouteNode("10.10.11.11:12345"))
	nodes = append(nodes, createRouteNode("192.168.1.1:13666"))
	for _, n := range nodes {
		rs.Append(n)
	}

	for _, n := range nodes {
		assert.True(t, rs.Contains(n))
	}

	rs.RemoveMany(nodes)

	assert.Empty(t, rs.Nodes())
}

func TestRouteSet_AppendMany(t *testing.T) {
	rs := NewRouteSet()

	assert.Empty(t, rs.nodes)

	node1 := createRouteNode("127.0.0.1:31337")
	node2 := createRouteNode("10.10.11.11:12345")

	rs.AppendMany([]*RouteNode{node1, node2})

	assert.Equal(t, []*node.Node{node1.Node, node2.Node}, rs.nodes)
}

func TestRouteSet_Len(t *testing.T) {
	rs := NewRouteSet()

	assert.Equal(t, 0, rs.Len())

	node1 := createRouteNode("127.0.0.1:31337")
	node2 := createRouteNode("10.10.11.11:12345")

	rs.Append(node1)
	rs.Append(node2)

	assert.Equal(t, 2, rs.Len())
}

func TestRouteSet_Swap(t *testing.T) {
	rs := NewRouteSet()

	assert.Empty(t, rs.nodes)

	node1 := createRouteNode("127.0.0.1:31337")
	node2 := createRouteNode("10.10.11.11:12345")

	rs.Append(node1)
	rs.Append(node2)

	assert.Equal(t, []*node.Node{node1.Node, node2.Node}, rs.nodes)

	rs.Swap(0, 1)

	assert.Equal(t, []*node.Node{node2.Node, node1.Node}, rs.nodes)
}

func TestRouteSet_Less(t *testing.T) {
	addr, _ := node.NewAddress("127.0.0.1:31337")
	node1 := node.NewNode(addr)
	node1.ID = getIDWithValues(5)
	node2 := node.NewNode(addr)
	node2.ID = getIDWithValues(7)
	rs := NewRouteSet()
	rs.Append(NewRouteNode(node1))
	rs.Append(NewRouteNode(node2))

	assert.True(t, rs.Less(0, 1))
	assert.False(t, rs.Less(1, 0))
}
