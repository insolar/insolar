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
	"log"

	"github.com/insolar/insolar/network/hostnetwork"
	"github.com/insolar/insolar/network/hostnetwork/id"
)

// Node is an essence which provides communication between network level and MessageRouter.
type Node struct {
	id       id.ID
	role     string
	dht      *hostnetwork.DHT
	domainID string
}

// NewNode creates a node with given args.
func NewNode(newID id.ID, newDHT *hostnetwork.DHT, newDomainID string, newRole string) (*Node, error) {
	node := &Node{
		id:       newID,
		dht:      newDHT,
		domainID: newDomainID,
		role:     newRole,
	}
	err := node.dht.StartCheckNodesRole(node.createContext(), node.domainID)
	return node, err
}

// SetNodeRole sets a new Node role.
func (node *Node) SetNodeRole(newRole string) {
	node.role = newRole
}

// SetNodeID sets a new Node id.
func (node *Node) SetNodeID(newID id.ID) {
	node.id = newID
}

// SetDHT sets a new Node host.
func (node *Node) SetDHT(newDHT *hostnetwork.DHT) {
	node.dht = newDHT
}

// SetDomainID sets a new domain ID.
func (node *Node) SetDomainID(newDomainID string) {
	node.domainID = newDomainID
}

// GetNodeRole returns a Node role.
func (node Node) GetNodeRole() string {
	return node.role
}

// GetNodeID returns a Node ID.
func (node Node) GetNodeID() id.ID {
	return node.id
}

// GetDomainIDs returns a Node domain ID.
func (node Node) GetDomainID() string {
	return node.domainID
}

// SendPacket sends packet from service to target.
func (node Node) SendPacket(targetID, method string, args [][]byte) error {
	_, err := node.dht.RemoteProcedureCall(node.createContext(), targetID, method, args)
	return err
}

// CreateContext returns a dht context.
func (node Node) createContext() hostnetwork.Context {
	ctx, err := hostnetwork.NewContextBuilder(node.dht).SetDefaultHost().Build()
	if err != nil {
		log.Fatalln("Failed to create context:", err.Error())
	}
	return ctx
}
