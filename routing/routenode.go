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
	"github.com/insolar/network/node"
)

// RouteNode represents a node in the network locally
// a separate struct due to the fact that we may want to add some metadata
// here later such as RTT, or LastSeen time.
type RouteNode struct {
	*node.Node
}

// NewRouteNode creates new RouteNode.
func NewRouteNode(node *node.Node) *RouteNode {
	return &RouteNode{
		Node: node,
	}
}

// RouteNodesFrom creates list of RouteNodes from a list of Nodes.
func RouteNodesFrom(nodes []*node.Node) []*RouteNode {
	routeNodes := make([]*RouteNode, len(nodes))

	for i, n := range nodes {
		routeNodes[i] = NewRouteNode(n)
	}

	return routeNodes
}
