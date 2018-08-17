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

Package node is a fundamental part of networking system. Each node has:

 - one real network address (IP or any other transport protocol address)
 - multiple abstract network IDs (either node's own or ones belonging to relayed nodes)

Contains structures to describe network entities in code.

Usage:

 	originAddress, err := node.NewAddress(address)
	if err != nil {
		...
	}

	origin, err := node.NewOrigin(nil, originAddress)
	if err != nil {
		...
	}

*/
package node
