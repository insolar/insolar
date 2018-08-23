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
	"github.com/insolar/insolar/network/hostnetwork/host"
)

// RouteSet is used in order to sort a list of arbitrary hosts against a
// comparator. These hosts are sorted by xor distance.
type RouteSet struct {
	// hosts are a list of hosts to be compared.
	hosts []*host.Host

	// comparator is the requestID to compare to.
	comparator []byte
}

// NewRouteSet creates new RouteSet.
func NewRouteSet() *RouteSet {
	return &RouteSet{}
}

// Nodes returns list of RouteSet hosts.
func (rs *RouteSet) Nodes() []*host.Host {
	nodes := make([]*host.Host, len(rs.hosts))
	copy(nodes, rs.hosts)
	return nodes
}

// FirstHost returns first Host from RouteSet.
func (rs *RouteSet) FirstHost() *host.Host {
	return rs.hosts[0]
}

// Contains checks if RouteSet contains given RouteNode.
func (rs *RouteSet) Contains(node *RouteNode) bool {
	exists := false
	for _, n := range rs.hosts {
		if node.ID.HashEqual(n.ID.GetHash()) {
			exists = true
		}
	}
	return exists
}

// Append adds single RouteNode to RouteSet.
func (rs *RouteSet) Append(node *RouteNode) {
	if !rs.Contains(node) {
		rs.hosts = append(rs.hosts, node.Host)
	}
}

// Remove removes host from RouteSet.
func (rs *RouteSet) Remove(node *RouteNode) {
	for i, n := range rs.hosts {
		if n.ID.HashEqual(node.ID.GetHash()) {
			rs.hosts = append(rs.hosts[:i], rs.hosts[i+1:]...)
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

// Len returns number of hosts in RouteSet.
func (rs *RouteSet) Len() int {
	return len(rs.hosts)
}

// Swap swaps two hosts in RouteSet.
func (rs *RouteSet) Swap(i, j int) {
	rs.hosts[i], rs.hosts[j] = rs.hosts[j], rs.hosts[i]
}

// Less is a sorting function for RouteSet.
func (rs *RouteSet) Less(i, j int) bool {
	iDist := getDistance(rs.hosts[i].ID.GetHash(), rs.comparator)
	jDist := getDistance(rs.hosts[j].ID.GetHash(), rs.comparator)

	return iDist.Cmp(jDist) == -1
}
