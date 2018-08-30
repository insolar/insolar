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
	"bytes"
	"errors"
	"log"

	"github.com/insolar/insolar/network/hostnetwork"
	"github.com/insolar/insolar/network/hostnetwork/host"
	"github.com/insolar/insolar/network/hostnetwork/id"
)

// Nodenetwork is nodes manager.
type Nodenetwork struct {
	dht   *hostnetwork.DHT
	nodes map[string]*Node // key - reference ID, value - node ID.
}

// NewNodeNetwork creates a new nodenetwork.
func NewNodeNetwork(DHT *hostnetwork.DHT) *Nodenetwork {
	network := &Nodenetwork{
		dht:   DHT,
		nodes: make(map[string]*Node, 1),
	}
	return network
}

// AddNode adds a new node to nodes map.
func (network *Nodenetwork) AddNode(hostAddress, hostID, domainID string) error {
	node, err := network.createNode(hostAddress, hostID, domainID)
	if err != nil {
		return err
	}
	return network.addNode(node)
}

// SendPacket sends packet from service to target.
func (network Nodenetwork) SendPacket(reference, method string, args [][]byte) ([]byte, error) {
	host := network.nodes[reference].host
	if host == nil {
		return nil, errors.New("host doesn't exist")
	}
	res, err := network.dht.RemoteProcedureCall(network.createContext(), host.ID.HashString(), method, args)
	return res, err
}

func (network *Nodenetwork) createNode(hostAddress, hostID, domainID string) (*Node, error) {
	address, err := host.NewAddress(hostAddress)
	if err != nil {
		return nil, err
	}
	newHost := host.NewHost(address)
	key := id.GetRandomKey()
	stringID := string(key[:bytes.IndexByte(key, 0)])
	node := NewNode(stringID, newHost, domainID)
	return node, nil
}

func (network *Nodenetwork) addNode(node *Node) error {
	if node == nil {
		return errors.New("node is nil")
	}
	network.nodes[node.GetDomainID()] = node
	return nil
}

func (network Nodenetwork) createContext() hostnetwork.Context {

	ctx, err := hostnetwork.NewContextBuilder(network.dht).SetDefaultHost().Build()
	if err != nil {
		log.Fatalln("Failed to create context:", err.Error())
	}
	return ctx
}
