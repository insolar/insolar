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

package routing

import (
	"github.com/insolar/insolar/network/transport/host"
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
func NewRouteSet(comparator []byte) *RouteSet {
	return &RouteSet{comparator: comparator}
}

// Hosts returns list of RouteSet hosts.
func (rs *RouteSet) Hosts() []*host.Host {
	hosts := make([]*host.Host, len(rs.hosts))
	copy(hosts, rs.hosts)
	return hosts
}

// FirstHost returns first Host from RouteSet.
func (rs *RouteSet) FirstHost() *host.Host {
	return rs.hosts[0]
}

// Contains checks if RouteSet contains given RouteHost.
func (rs *RouteSet) Contains(host *RouteHost) bool {
	exists := false
	for _, n := range rs.hosts {
		if host.ID.Equal(n.ID.Bytes()) {
			exists = true
		}
	}
	return exists
}

// Append adds single RouteHost to RouteSet.
func (rs *RouteSet) Append(host *RouteHost) {
	if !rs.Contains(host) {
		rs.hosts = append(rs.hosts, host.Host)
	}
}

// Remove removes host from RouteSet.
func (rs *RouteSet) Remove(host *RouteHost) {
	for i, n := range rs.hosts {
		if n.ID.Equal(host.ID.Bytes()) {
			rs.hosts = append(rs.hosts[:i], rs.hosts[i+1:]...)
			return
		}
	}
}

// RemoveMany removes list of RoutHosts from RouteSet
func (rs *RouteSet) RemoveMany(hosts []*RouteHost) {
	for _, n := range hosts {
		rs.Remove(n)
	}
}

// AppendMany adds a list of RouteHosts to RouteSet.
func (rs *RouteSet) AppendMany(hosts []*RouteHost) {
	for _, n := range hosts {
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
	iDist := getDistance(rs.hosts[i].ID.Bytes(), rs.comparator)
	jDist := getDistance(rs.hosts[j].ID.Bytes(), rs.comparator)

	return iDist.Cmp(jDist) == -1
}
