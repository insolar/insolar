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
	"github.com/insolar/insolar/core"
)

// Node is an essence which provides communication between network level and MessageRouter.
type Node struct {
	id        string
	role      string
	hostID    string
	reference core.RecordRef
}

// NewNode creates a node with given args.
func NewNode(nodeID, hostID string, domainID core.RecordRef) *Node {
	return &Node{
		id:        nodeID,
		hostID:    hostID,
		reference: domainID,
	}
}

// nolint
func (node *Node) setRole(newRole string) {
	node.role = newRole
}

// GetNodeRole returns a Node role.
func (node Node) GetNodeRole() string {
	return node.role
}

// GetNodeID returns a Node ID.
func (node Node) GetNodeID() string {
	return node.id
}

// GetReference returns a Node domain ID.
func (node Node) GetReference() core.RecordRef {
	return node.reference
}
