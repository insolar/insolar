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

package nodenetwork

import (
	"github.com/insolar/insolar/network/hostnetwork/id"
)

// Node is an essence which provides communication between network level and MessageRouter.
type Node struct {
	id       []byte
	role     string
	hostID   string
	domainID string
}

// NewNode creates a node with given args.
func NewNode(hostID, domainID string) *Node {
	return &Node{
		id:       id.GetRandomKey(),
		hostID:   hostID,
		domainID: domainID,
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
func (node Node) GetNodeID() []byte {
	return node.id
}

// GetDomainIDs returns a Node domain ID.
func (node Node) GetDomainID() string {
	return node.domainID
}
