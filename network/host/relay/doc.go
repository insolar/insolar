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
Package relay is an implementation of relay mechanism. Proxy contains info about nodes which can relaying messages.
Relay contains info about nodes of which messages current node have to relay.

Usage:
	package relay

	relay := NewRelay()
	relay.AddClient(node)

	if relay.NeedToRelay(node.Address()) {
		//relay packet
	}

	relay.RemoveClient(node)

	//-----------------------------------

	proxy := NewProxy()
	proxy.AddProxyNode(node.Address())

	if proxy.ProxyNodesCount > 0 {
		address := proxy.GetNextProxyAddress
		//send packet to next proxy
	}

	proxy.RemoveProxyNode(node.Address)

*/
package relay
