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

package nodenetwork

import (
	"crypto/sha1"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/core"
	"golang.org/x/crypto/sha3"
)

// NodeNetwork is node manager.
type NodeNetwork struct {
	node *Node
}

// NewNodeNetwork creates a new node network.
func NewNodeNetwork(nodeCfg configuration.NodeNetwork) *NodeNetwork {
	node := NewNode(core.String2Ref(nodeCfg.Node.ID))
	network := &NodeNetwork{
		node: node,
	}
	return network
}

// ResolveHostID returns a host found by reference.
func (network *NodeNetwork) ResolveHostID(ref core.RecordRef) string {
	sha3digest := sha3.Sum512(ref[:])
	sha1digest := sha1.Sum(sha3digest[:])
	return string(sha1digest[:])
}

// GetID returns current node id
func (network *NodeNetwork) GetID() core.RecordRef {
	return network.node.GetID()
}
