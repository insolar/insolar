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

package relay

import (
	"strconv"
	"testing"

	"github.com/insolar/insolar/network/transport/host"
	"github.com/insolar/insolar/testutils"
	"github.com/stretchr/testify/assert"
)

func TestRelay_ClientsCount(t *testing.T) {
	relay := NewRelay()
	count := 20

	hosts := makeHosts(count, t)

	for i := range hosts {
		relay.AddClient(hosts[i])
	}

	assert.Equal(t, count, relay.ClientsCount())
}

func TestNewProxy(t *testing.T) {
	proxy := NewProxy()
	assert.NotNil(t, proxy)
}

func TestNewRelay(t *testing.T) {
	relay := NewRelay()
	assert.NotNil(t, relay)
}

func makeAddresses(count int, t *testing.T) []*host.Address {
	ip := "127.0.0.1:"
	addresses := make([]*host.Address, 0)

	for i := 0; i < count; i++ {
		address, err := host.NewAddress(ip + strconv.Itoa(i+20000))
		if err != nil {
			assert.Errorf(t, nil, "error: %s", err.Error())
			continue
		}
		addresses = append(addresses, address)
	}

	return addresses
}

func makeHosts(count int, t *testing.T) []*host.Host {
	result := make([]*host.Host, 0)
	addresses := makeAddresses(count, t)

	for i := 0; i < count; i++ {
		id := testutils.RandomRef()
		result = append(result, &host.Host{NodeID: id, Address: addresses[i]})
	}

	return result
}

func TestRelay_AddClient(t *testing.T) {
	relay := NewRelay()
	count := 20

	hosts := makeHosts(count, t)

	for i := range hosts {
		err := relay.AddClient(hosts[i])
		assert.NoError(t, err)
		err = relay.AddClient(hosts[i]) // adding existing host
		assert.EqualError(t, err, "client exists already")
	}

	assert.Equal(t, count, relay.ClientsCount())
}

func TestRelay_RemoveClient(t *testing.T) {
	relay := NewRelay()
	count := 20

	hosts := makeHosts(count, t)

	for i := range hosts {
		err := relay.AddClient(hosts[i])
		if err != nil {
			assert.Errorf(t, nil, "error: %s", err.Error())
		}
	}
	assert.Equal(t, count, relay.ClientsCount())

	for i := range hosts {
		err := relay.RemoveClient(hosts[i])
		assert.NoError(t, err)
		err = relay.RemoveClient(hosts[i])
		assert.EqualError(t, err, "client not found")
	}

	assert.Equal(t, 0, relay.ClientsCount())
}

func TestRelay_NeedToRelay(t *testing.T) {
	relay := NewRelay()
	count := 20
	ip := "127.0.0.2:"

	hosts := makeHosts(count, t)

	for i := range hosts {
		relay.AddClient(hosts[i])
	}

	assert.Equal(t, count, relay.ClientsCount())

	for i := range hosts {
		res := relay.NeedToRelay(hosts[i].Address.String())
		assert.Equal(t, true, res)
	}

	for i := 0; i < count; i++ {
		address, err := host.NewAddress(ip + strconv.Itoa(i+20000))

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

	hosts := makeHosts(count, t)

	for i := range hosts {
		relay.AddClient(hosts[i])
	}

	assert.Equal(t, count, relay.ClientsCount())
}

func TestProxy_AddProxyHost(t *testing.T) {
	count := 20
	addresses := makeAddresses(count, t)
	proxy := NewProxy()

	for i := range addresses {
		proxy.AddProxyHost(addresses[i].String())
		proxy.AddProxyHost(addresses[i].String()) // adding existed host
	}

	assert.Equal(t, count, proxy.ProxyHostsCount())
}

func TestProxy_RemoveProxyHost(t *testing.T) {
	count := 20
	addresses := makeAddresses(count, t)
	proxy := NewProxy()

	for i := range addresses {
		proxy.AddProxyHost(addresses[i].String())
	}

	assert.Equal(t, count, proxy.ProxyHostsCount())

	for i := range addresses {
		proxy.RemoveProxyHost(addresses[i].String())
		proxy.RemoveProxyHost(addresses[i].String()) // remove removed host
	}

	assert.Equal(t, 0, proxy.ProxyHostsCount())
}

func TestProxy_GetNextProxyAddress(t *testing.T) {
	count := 20
	addresses := makeAddresses(count, t)
	proxy := NewProxy()
	idx := make(map[int]string, count)

	assert.Equal(t, "", proxy.GetNextProxyAddress())

	for i := range addresses {
		proxy.AddProxyHost(addresses[i].String())
		idx[i] = addresses[i].String()
	}

	assert.Equal(t, count, proxy.ProxyHostsCount())
	assert.Equal(t, count, len(idx))

	for i := 0; i < proxy.ProxyHostsCount(); i++ {
		assert.Equal(t, idx[i], proxy.GetNextProxyAddress())
	}
	for i := 0; i < proxy.ProxyHostsCount(); i++ {
		assert.Equal(t, idx[i], proxy.GetNextProxyAddress())
	}
}

func TestProxy_Count(t *testing.T) {
	count := 20
	addresses := makeAddresses(count, t)
	proxy := NewProxy()

	for i := range addresses {
		proxy.AddProxyHost(addresses[i].String())
	}

	assert.Equal(t, count, proxy.ProxyHostsCount())
}

func TestCreateProxy(t *testing.T) {
	proxy := NewProxy()

	check := true

	if proxy == nil {
		check = false
	}

	assert.Equal(t, true, check)
}
