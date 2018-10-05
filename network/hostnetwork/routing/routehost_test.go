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
	"testing"

	"github.com/insolar/insolar/network/hostnetwork/host"

	"github.com/stretchr/testify/assert"
)

func TestNewRouteHost(t *testing.T) {
	testAddr, _ := host.NewAddress("127.0.0.1:31337")
	testHost := host.NewHost(testAddr)

	expectedRouteHost := &RouteHost{testHost}
	actualRouteHost := NewRouteHost(testHost)

	assert.Equal(t, expectedRouteHost, actualRouteHost)
}

func TestRouteHostsFrom(t *testing.T) {
	testAddr1, _ := host.NewAddress("127.0.0.1:31337")
	testAddr2, _ := host.NewAddress("10.10.11.11:31338")
	hosts := []*host.Host{host.NewHost(testAddr1), host.NewHost(testAddr2)}

	routeHosts := RouteHostsFrom(hosts)

	assert.Len(t, routeHosts, 2)
	assert.Equal(t, hosts[0].String(), routeHosts[0].String())
	assert.Equal(t, hosts[1].String(), routeHosts[1].String())
}
