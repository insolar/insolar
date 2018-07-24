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
	"github.com/insolar/insolar/network/host/node"
)

// RouteSet is used in order to sort a list of arbitrary nodes against a
// comparator. These nodes are sorted by xor distance.
type RouteSet struct {
	// nodes are a list of nodes to be compared.
	nodes []*node.Node

	// comparator is the requestID to compare to.
	comparator []byte
}

// NewRouteSet creates new RouteSet.
func NewRouteSet() *RouteSet {
	return &RouteSet{}
}

// Nodes returns list of RouteSet nodes.
func (rs *RouteSet) Nodes() []*node.Node {
	nodes := make([]*node.Node, len(rs.nodes))
	copy(nodes, rs.nodes)
	return nodes
}

// FirstNode returns first Node from RouteSet.
func (rs *RouteSet) FirstNode() *node.Node {
	return rs.nodes[0]
}

// Contains checks if RouteSet contains given RouteNode.
func (rs *RouteSet) Contains(node *RouteNode) bool {
	exists := false
	for _, n := range rs.nodes {
		if node.ID.Equal(n.ID) {
			exists = true
		}
	}
	return exists
}

// Append adds single RouteNode to RouteSet.
func (rs *RouteSet) Append(node *RouteNode) {
	if !rs.Contains(node) {
		rs.nodes = append(rs.nodes, node.Node)
	}
}

// Remove removes node from RouteSet.
func (rs *RouteSet) Remove(node *RouteNode) {
	for i, n := range rs.nodes {
		if n.ID.Equal(node.ID) {
			rs.nodes = append(rs.nodes[:i], rs.nodes[i+1:]...)
			return
		}
	}
}

// RemoveMany removes list of RoutNodes from RouteSet
func (rs *RouteSet) RemoveMany(nodes []*RouteNode) {
	for _, n := range nodes {
		rs.Remove(n)
	}
}

// AppendMany adds a list of RouteNodes to RouteSet.
func (rs *RouteSet) AppendMany(nodes []*RouteNode) {
	for _, n := range nodes {
		rs.Append(n)
	}
}

// Len returns number of nodes in RouteSet.
func (rs *RouteSet) Len() int {
	return len(rs.nodes)
}

// Swap swaps two nodes in RouteSet.
func (rs *RouteSet) Swap(i, j int) {
	rs.nodes[i], rs.nodes[j] = rs.nodes[j], rs.nodes[i]
}

// Less is a sorting function for RouteSet.
func (rs *RouteSet) Less(i, j int) bool {
	iDist := getDistance(rs.nodes[i].ID, rs.comparator)
	jDist := getDistance(rs.nodes[j].ID, rs.comparator)

	return iDist.Cmp(jDist) == -1
}
