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

	"github.com/insolar/insolar/network/hostnetwork/host"
	"github.com/insolar/insolar/network/hostnetwork/id"

	"github.com/stretchr/testify/assert"
)

func createRouteHost(addrString string) *RouteHost {
	addr, _ := host.NewAddress(addrString)
	newHost := host.NewHost(addr)
	newHost.ID, _ = id.NewID(id.GetRandomKey())
	return NewRouteHost(newHost)
}

func TestNewRouteSet(t *testing.T) {
	rs := NewRouteSet()

	assert.Equal(t, &RouteSet{}, rs)
	assert.Implements(t, (*sort.Interface)(nil), rs)
}

func TestRouteSet_Hosts(t *testing.T) {
	rs := NewRouteSet()

	assert.Empty(t, rs.Hosts())

	host1 := createRouteHost("127.0.0.1:31337")
	host2 := createRouteHost("10.10.11.11:12345")

	rs.Append(host1)
	rs.Append(host2)

	assert.Equal(t, []*host.Host{host1.Host, host2.Host}, rs.Hosts())
	assert.Equal(t, rs.hosts, rs.Hosts())
}

func TestRouteSet_Hosts_ReturnsCopy(t *testing.T) {
	rs := NewRouteSet()
	rs.Append(createRouteHost("127.0.0.1:31337"))
	rs.Append(createRouteHost("10.10.11.11:12345"))

	hostsCopy := rs.Hosts()
	hostsCopy[0] = nil

	assert.NotNil(t, rs.hosts[0])
}

func TestRouteSet_FirstHost(t *testing.T) {
	rs := NewRouteSet()
	host1 := createRouteHost("127.0.0.1:31337")
	host2 := createRouteHost("10.10.11.11:12345")
	rs.Append(host1)
	rs.Append(host2)

	assert.Equal(t, host1.Host, rs.FirstHost())
}

func TestRouteSet_Contains(t *testing.T) {
	rs := NewRouteSet()
	host1 := createRouteHost("127.0.0.1:31337")
	host2 := createRouteHost("10.10.11.11:12345")
	host3 := createRouteHost("192.168.1.1:13666")
	rs.Append(host1)
	rs.Append(host2)

	assert.True(t, rs.Contains(host1))
	assert.True(t, rs.Contains(host2))
	assert.False(t, rs.Contains(host3))
}

func TestRouteSet_Append(t *testing.T) {
	rs := NewRouteSet()

	assert.Empty(t, rs.hosts)

	host1 := createRouteHost("127.0.0.1:31337")
	host2 := createRouteHost("10.10.11.11:12345")

	rs.Append(host1)
	rs.Append(host2)

	assert.Equal(t, []*host.Host{host1.Host, host2.Host}, rs.hosts)
}

func TestRouteSet_Remove(t *testing.T) {
	rs := NewRouteSet()
	host1 := createRouteHost("127.0.0.1:31337")
	host2 := createRouteHost("10.10.11.11:12345")
	host3 := createRouteHost("192.168.1.1:13666")
	rs.Append(host1)
	rs.Append(host2)
	rs.Append(host3)

	assert.True(t, rs.Contains(host1))
	assert.True(t, rs.Contains(host2))
	assert.True(t, rs.Contains(host3))

	rs.Remove(host2)

	assert.True(t, rs.Contains(host1))
	assert.False(t, rs.Contains(host2))
	assert.True(t, rs.Contains(host3))
}

func TestRouteSet_RemoveMany(t *testing.T) {
	rs := NewRouteSet()
	var hosts []*RouteHost
	hosts = append(hosts, createRouteHost("127.0.0.1:31337"))
	hosts = append(hosts, createRouteHost("10.10.11.11:12345"))
	hosts = append(hosts, createRouteHost("192.168.1.1:13666"))
	for _, n := range hosts {
		rs.Append(n)
	}

	for _, n := range hosts {
		assert.True(t, rs.Contains(n))
	}

	rs.RemoveMany(hosts)

	assert.Empty(t, rs.Hosts())
}

func TestRouteSet_AppendMany(t *testing.T) {
	rs := NewRouteSet()

	assert.Empty(t, rs.hosts)

	host1 := createRouteHost("127.0.0.1:31337")
	host2 := createRouteHost("10.10.11.11:12345")

	rs.AppendMany([]*RouteHost{host1, host2})

	assert.Equal(t, []*host.Host{host1.Host, host2.Host}, rs.hosts)
}

func TestRouteSet_Len(t *testing.T) {
	rs := NewRouteSet()

	assert.Equal(t, 0, rs.Len())

	host1 := createRouteHost("127.0.0.1:31337")
	host2 := createRouteHost("10.10.11.11:12345")

	rs.Append(host1)
	rs.Append(host2)

	assert.Equal(t, 2, rs.Len())
}

func TestRouteSet_Swap(t *testing.T) {
	rs := NewRouteSet()

	assert.Empty(t, rs.hosts)

	host1 := createRouteHost("127.0.0.1:31337")
	host2 := createRouteHost("10.10.11.11:12345")

	rs.Append(host1)
	rs.Append(host2)

	assert.Equal(t, []*host.Host{host1.Host, host2.Host}, rs.hosts)

	rs.Swap(0, 1)

	assert.Equal(t, []*host.Host{host2.Host, host1.Host}, rs.hosts)
}

func TestRouteSet_Less(t *testing.T) {
	addr, _ := host.NewAddress("127.0.0.1:31337")
	host1 := host.NewHost(addr)
	host1.ID = getIDWithValues(5)
	host2 := host.NewHost(addr)
	host2.ID = getIDWithValues(7)
	rs := NewRouteSet()
	rs.Append(NewRouteHost(host1))
	rs.Append(NewRouteHost(host2))

	assert.True(t, rs.Less(0, 1))
	assert.False(t, rs.Less(1, 0))
}
