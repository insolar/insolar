//
// Modified BSD 3-Clause Clear License
//
// Copyright (c) 2019 Insolar Technologies GmbH
//
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without modification,
// are permitted (subject to the limitations in the disclaimer below) provided that
// the following conditions are met:
//  * Redistributions of source code must retain the above copyright notice, this list
//    of conditions and the following disclaimer.
//  * Redistributions in binary form must reproduce the above copyright notice, this list
//    of conditions and the following disclaimer in the documentation and/or other materials
//    provided with the distribution.
//  * Neither the name of Insolar Technologies GmbH nor the names of its contributors
//    may be used to endorse or promote products derived from this software without
//    specific prior written permission.
//
// NO EXPRESS OR IMPLIED LICENSES TO ANY PARTY'S PATENT RIGHTS ARE GRANTED
// BY THIS LICENSE. THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS
// AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES,
// INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY
// AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL
// THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT,
// INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING,
// BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS
// OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
// ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
//
// Notwithstanding any other provisions of this license, it is prohibited to:
//    (a) use this software,
//
//    (b) prepare modifications and derivative works of this software,
//
//    (c) distribute this software (including without limitation in source code, binary or
//        object code form), and
//
//    (d) reproduce copies of this software
//
//    for any commercial purposes, and/or
//
//    for the purposes of making available this software to third parties as a service,
//    including, without limitation, any software-as-a-service, platform-as-a-service,
//    infrastructure-as-a-service or other similar online service, irrespective of
//    whether it competes with the products or services of Insolar Technologies GmbH.
//

package relay

import (
	"strconv"
	"testing"

	"github.com/insolar/insolar/network/transport/host"
	"github.com/insolar/insolar/testutils"
	"github.com/stretchr/testify/require"
)

func TestRelay_ClientsCount(t *testing.T) {
	relay := NewRelay()
	count := 20

	hosts := makeHosts(count, t)

	for i := range hosts {
		relay.AddClient(hosts[i])
	}

	require.Equal(t, count, relay.ClientsCount())
}

func TestNewProxy(t *testing.T) {
	proxy := NewProxy()
	require.NotNil(t, proxy)
}

func TestNewRelay(t *testing.T) {
	relay := NewRelay()
	require.NotNil(t, relay)
}

func makeAddresses(count int, t *testing.T) []*host.Address {
	ip := "127.0.0.1:"
	addresses := make([]*host.Address, 0)

	for i := 0; i < count; i++ {
		address, err := host.NewAddress(ip + strconv.Itoa(i+20000))
		if err != nil {
			require.Errorf(t, nil, "error: %s", err.Error())
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
		require.NoError(t, err)
		err = relay.AddClient(hosts[i]) // adding existing host
		require.EqualError(t, err, "client exists already")
	}

	require.Equal(t, count, relay.ClientsCount())
}

func TestRelay_RemoveClient(t *testing.T) {
	relay := NewRelay()
	count := 20

	hosts := makeHosts(count, t)

	for i := range hosts {
		err := relay.AddClient(hosts[i])
		if err != nil {
			require.Errorf(t, nil, "error: %s", err.Error())
		}
	}
	require.Equal(t, count, relay.ClientsCount())

	for i := range hosts {
		err := relay.RemoveClient(hosts[i])
		require.NoError(t, err)
		err = relay.RemoveClient(hosts[i])
		require.EqualError(t, err, "client not found")
	}

	require.Equal(t, 0, relay.ClientsCount())
}

func TestRelay_NeedToRelay(t *testing.T) {
	relay := NewRelay()
	count := 20
	ip := "127.0.0.2:"

	hosts := makeHosts(count, t)

	for i := range hosts {
		relay.AddClient(hosts[i])
	}

	require.Equal(t, count, relay.ClientsCount())

	for i := range hosts {
		res := relay.NeedToRelay(hosts[i].Address.String())
		require.Equal(t, true, res)
	}

	for i := 0; i < count; i++ {
		address, err := host.NewAddress(ip + strconv.Itoa(i+20000))

		if err != nil {
			require.Errorf(t, nil, "error: %s", err.Error())
			continue
		}
		res := relay.NeedToRelay(address.String())
		require.Equal(t, false, res)
	}
}

func TestRelay_Count(t *testing.T) {
	relay := NewRelay()
	count := 20

	hosts := makeHosts(count, t)

	for i := range hosts {
		relay.AddClient(hosts[i])
	}

	require.Equal(t, count, relay.ClientsCount())
}

func TestProxy_AddProxyHost(t *testing.T) {
	count := 20
	addresses := makeAddresses(count, t)
	proxy := NewProxy()

	for i := range addresses {
		proxy.AddProxyHost(addresses[i].String())
		proxy.AddProxyHost(addresses[i].String()) // adding existed host
	}

	require.Equal(t, count, proxy.ProxyHostsCount())
}

func TestProxy_RemoveProxyHost(t *testing.T) {
	count := 20
	addresses := makeAddresses(count, t)
	proxy := NewProxy()

	for i := range addresses {
		proxy.AddProxyHost(addresses[i].String())
	}

	require.Equal(t, count, proxy.ProxyHostsCount())

	for i := range addresses {
		proxy.RemoveProxyHost(addresses[i].String())
		proxy.RemoveProxyHost(addresses[i].String()) // remove removed host
	}

	require.Equal(t, 0, proxy.ProxyHostsCount())
}

func TestProxy_GetNextProxyAddress(t *testing.T) {
	count := 20
	addresses := makeAddresses(count, t)
	proxy := NewProxy()
	idx := make(map[int]string, count)

	require.Equal(t, "", proxy.GetNextProxyAddress())

	for i := range addresses {
		proxy.AddProxyHost(addresses[i].String())
		idx[i] = addresses[i].String()
	}

	require.Equal(t, count, proxy.ProxyHostsCount())
	require.Equal(t, count, len(idx))

	for i := 0; i < proxy.ProxyHostsCount(); i++ {
		require.Equal(t, idx[i], proxy.GetNextProxyAddress())
	}
	for i := 0; i < proxy.ProxyHostsCount(); i++ {
		require.Equal(t, idx[i], proxy.GetNextProxyAddress())
	}
}

func TestProxy_Count(t *testing.T) {
	count := 20
	addresses := makeAddresses(count, t)
	proxy := NewProxy()

	for i := range addresses {
		proxy.AddProxyHost(addresses[i].String())
	}

	require.Equal(t, count, proxy.ProxyHostsCount())
}

func TestCreateProxy(t *testing.T) {
	proxy := NewProxy()

	check := true

	if proxy == nil {
		check = false
	}

	require.Equal(t, true, check)
}
