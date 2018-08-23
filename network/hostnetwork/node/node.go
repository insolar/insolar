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

package node

import (
	"fmt"

	"github.com/insolar/insolar/network/hostnetwork/id"
)

// Node is the over-the-wire representation of a node.
type Node struct {
	// ID is a 20 byte unique identifier.
	ID id.ID

	// Address is IP and port.
	Address *Address
}

// NewNode creates a new Node for bootstrapping.
func NewNode(address *Address) *Node {
	return &Node{
		Address: address,
	}
}

// String representation of Node.
func (node Node) String() string {
	return fmt.Sprintf("%s (%s)", node.ID.HashString(), node.Address.String())
}

// Equal checks if node equals to other node (e.g. nodes' IDs and network addresses match).
func (node Node) Equal(other Node) bool {
	return node.ID.HashEqual(other.ID.GetHash()) && other.Address != nil && node.Address.Equal(*other.Address)
}
