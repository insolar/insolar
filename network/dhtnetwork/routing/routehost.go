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

// RouteHost represents a host in the network locally
// a separate struct due to the fact that we may want to add some metadata
// here later such as RTT, or LastSeen time.
type RouteHost struct {
	*host.Host
}

// NewRouteHost creates new RouteHost.
func NewRouteHost(host *host.Host) *RouteHost {
	return &RouteHost{
		Host: host,
	}
}

// RouteHostsFrom creates list of RouteHosts from a list of Hosts.
func RouteHostsFrom(hosts []*host.Host) []*RouteHost {
	routeHosts := make([]*RouteHost, len(hosts))

	for i, n := range hosts {
		routeHosts[i] = NewRouteHost(n)
	}

	return routeHosts
}
