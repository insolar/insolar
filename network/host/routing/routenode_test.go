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
	"testing"

	"github.com/insolar/insolar/network/host/node"

	"github.com/stretchr/testify/assert"
)

func TestNewRouteNode(t *testing.T) {
	testAddr, _ := node.NewAddress("127.0.0.1:31337")
	testNode := node.NewNode(testAddr)

	expectedRouteNode := &RouteNode{testNode}
	actualRouteNode := NewRouteNode(testNode)

	assert.Equal(t, expectedRouteNode, actualRouteNode)
}

func TestRouteNodesFrom(t *testing.T) {
	testAddr1, _ := node.NewAddress("127.0.0.1:31337")
	testAddr2, _ := node.NewAddress("10.10.11.11:31338")
	nodes := []*node.Node{node.NewNode(testAddr1), node.NewNode(testAddr2)}

	routeNodes := RouteNodesFrom(nodes)

	assert.Len(t, routeNodes, 2)
	assert.Equal(t, nodes[0].String(), routeNodes[0].String())
	assert.Equal(t, nodes[1].String(), routeNodes[1].String())
}
