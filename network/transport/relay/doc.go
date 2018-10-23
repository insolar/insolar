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

/*
Package relay is an implementation of relay mechanism. Proxy contains info about hosts which can relaying packets.
Relay contains info about hosts of which packets current host have to relay.

Usage:
	package relay

	relay := NewRelay()
	relay.AddClient(host)

	if relay.NeedToRelay(host.Address()) {
		//relay packet
	}

	relay.RemoveClient(host)

	//-----------------------------------

	proxy := NewProxy()
	proxy.AddProxyHost(host.Address())

	if proxy.ProxyHostsCount > 0 {
		address := proxy.GetNextProxyAddress
		//send packet to next proxy
	}

	proxy.RemoveProxyHost(host.Address)

*/
package relay
