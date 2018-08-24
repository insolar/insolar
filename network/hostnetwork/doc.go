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

/*
Package hostnetwork is an implementation of Kademlia DHT. It is mostly based on original specification but has multiple backward-incompatible changes.

Usage:

	package main

	import (
		"github.com/insolar/insolar/network/hostnetwork"
		"github.com/insolar/insolar/network/hostnetwork/connection"
		"github.com/insolar/insolar/network/hostnetwork/host"
		"github.com/insolar/insolar/network/hostnetwork/relay"
		"github.com/insolar/insolar/network/hostnetwork/resolver"
		"github.com/insolar/insolar/network/hostnetwork/rpc"
		"github.com/insolar/insolar/network/hostnetwork/store"
		"github.com/insolar/insolar/network/hostnetwork/transport"
	)

	func main() {
		configuration := hostnetwork.NewNetworkConfiguration(
			resolver.NewStunResolver(""),
			connection.NewConnectionFactory(),
			transport.NewUTPTransportFactory(),
			store.NewMemoryStoreFactory(),
			rpc.NewRPCFactory(map[string]rpc.RemoteProcedure{}),
			relay.NewProxy())

		dhtNetwork, err := configuration.NewHostNetwork("0.0.0.0:31337", &hostnetwork.Options{})
		if err != nil {
			panic(err)
		}
		defer configuration.CloseNetwork()

		dhtNetwork.Listen()
	}
*/
package hostnetwork
