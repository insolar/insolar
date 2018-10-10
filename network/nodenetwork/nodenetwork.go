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
	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/core"
	"github.com/jbenet/go-base58"
	"golang.org/x/crypto/sha3"
)

// NodeNetwork is node manager.
type NodeNetwork struct {
	node *Node
}

// NewNodeNetwork creates a new node network.
func NewNodeNetwork(nodeCfg configuration.NodeNetwork) *NodeNetwork {
	node := NewNode(core.NewRefFromBase58(nodeCfg.Node.ID))
	network := &NodeNetwork{
		node: node,
	}
	return network
}

// ResolveHostID returns a host found by reference.
func ResolveHostID(ref core.RecordRef) string {
	hash := make([]byte, 20)
	sha3.ShakeSum128(hash, ref[:])
	return base58.Encode(hash)
}

// GetID returns current node id
func (network *NodeNetwork) GetID() core.RecordRef {
	return network.node.GetID()
}
