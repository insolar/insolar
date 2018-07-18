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

package relay

import (
	"strconv"
	"testing"

	"github.com/insolar/network/node"
	"github.com/stretchr/testify/assert"
)

func TestNewRelay(t *testing.T) {
	relay := NewRelay()

	check := true

	if relay == nil {
		check = false
	}

	assert.Equal(t, true, check)
}

func makeAddresses(count int, t *testing.T) []*node.Address {
	ip := "127.0.0.1:"
	addresses := make([]*node.Address, 0)

	for i := 0; i < count; i++ {
		address, err := node.NewAddress(ip + strconv.Itoa(i+20000))
		if err != nil {
			assert.Errorf(t, nil, "error: %s", err.Error())
			continue
		}
		addresses = append(addresses, address)
	}

	return addresses
}

func makeNodes(count int, t *testing.T) []*node.Node {
	result := make([]*node.Node, 0)
	addresses := makeAddresses(count, t)

	for i := 0; i < count; i++ {
		id, err := node.NewID()

		if err != nil {
			assert.Errorf(t, nil, "error: %s", err.Error())
			continue
		}

		result = append(result, &node.Node{ID: id, Address: addresses[i]})
	}

	return result
}

func TestRelay_AddClient(t *testing.T) {
	relay := NewRelay()
	count := 20

	nodes := makeNodes(count, t)

	for i := range nodes {
		relay.AddClient(nodes[i])
	}

	assert.Equal(t, count, relay.ClientsCount())
}

func TestRelay_RemoveClient(t *testing.T) {
	relay := NewRelay()
	count := 20

	nodes := makeNodes(count, t)

	for i := range nodes {
		err := relay.AddClient(nodes[i])
		if err != nil {
			assert.Errorf(t, nil, "error: %s", err.Error())
		}
	}
	assert.Equal(t, count, relay.ClientsCount())

	for i := range nodes {
		relay.RemoveClient(nodes[i])
	}

	assert.Equal(t, 0, relay.ClientsCount())
}

func TestRelay_NeedToRelay(t *testing.T) {
	relay := NewRelay()
	count := 20
	ip := "127.0.0.2:"

	nodes := makeNodes(count, t)

	for i := range nodes {
		relay.AddClient(nodes[i])
	}

	assert.Equal(t, count, relay.ClientsCount())

	for i := range nodes {
		res := relay.NeedToRelay(nodes[i].Address.String())
		assert.Equal(t, true, res)
	}

	for i := 0; i < count; i++ {
		address, err := node.NewAddress(ip + strconv.Itoa(i+20000))

		if err != nil {
			assert.Errorf(t, nil, "error: %s", err.Error())
			continue
		}
		res := relay.NeedToRelay(address.String())
		assert.Equal(t, false, res)
	}
}

func TestRelay_Count(t *testing.T) {
	relay := NewRelay()
	count := 20

	nodes := makeNodes(count, t)

	for i := range nodes {
		relay.AddClient(nodes[i])
	}

	assert.Equal(t, count, relay.ClientsCount())
}

func TestProxy_AddProxyNode(t *testing.T) {
	count := 20
	addresses := makeAddresses(count, t)
	proxy := CreateProxy()

	for i := range addresses {
		proxy.AddProxyNode(addresses[i].String())
	}

	assert.Equal(t, count, proxy.ProxyNodesCount())
}

func TestProxy_RemoveProxyNode(t *testing.T) {
	count := 20
	addresses := makeAddresses(count, t)
	proxy := CreateProxy()

	for i := range addresses {
		proxy.AddProxyNode(addresses[i].String())
	}

	assert.Equal(t, count, proxy.ProxyNodesCount())

	for i := range addresses {
		proxy.RemoveProxyNode(addresses[i].String())
	}

	assert.Equal(t, 0, proxy.ProxyNodesCount())
}

func TestProxy_GetNextProxyAddress(t *testing.T) {
	count := 20
	addresses := makeAddresses(count, t)
	proxy := CreateProxy()
	idx := make(map[int]string, count)

	for i := range addresses {
		proxy.AddProxyNode(addresses[i].String())
		idx[i] = addresses[i].String()
	}

	assert.Equal(t, count, proxy.ProxyNodesCount())
	assert.Equal(t, count, len(idx))

	for i := 0; i < proxy.ProxyNodesCount(); i++ {
		assert.Equal(t, idx[i], proxy.GetNextProxyAddress())
	}
}

func TestProxy_Count(t *testing.T) {
	count := 20
	addresses := makeAddresses(count, t)
	proxy := CreateProxy()

	for i := range addresses {
		proxy.AddProxyNode(addresses[i].String())
	}

	assert.Equal(t, count, proxy.ProxyNodesCount())
}

func TestCreateProxy(t *testing.T) {
	proxy := CreateProxy()

	check := true

	if proxy == nil {
		check = false
	}

	assert.Equal(t, true, check)
}
