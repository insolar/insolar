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

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/network/hostnetwork/id"
)

// Nodenetwork is nodes manager.
type Nodenetwork struct {
	nodes map[string]*Node // key - reference ID, value - node ID.
}

// NewNodeNetwork creates a new nodenetwork.
func NewNodeNetwork(nodeCfg configuration.NodeNetwork) *Nodenetwork {
	nodes := make(map[string]*Node)
	for _, cfg := range nodeCfg.Nodes {
		node := NewNode(cfg.HostID, cfg.ReferenceID)
		nodes[cfg.ReferenceID] = node
	}
	network := &Nodenetwork{
		nodes: nodes,
	}
	return network
}

// AddNode adds a new node to nodes map.
func (network *Nodenetwork) AddNode(hostAddress, domainID string) error {
	node, err := network.createNode(hostAddress, domainID)
	if err != nil {
		return err
	}
	return network.addNode(node)
}

// GetReferenceHostID returns a host found by reference.
func (network *Nodenetwork) GetReferenceHostID(ref []byte) (string, error) {
	target := getReferenceString(ref)
	if _, ok := network.nodes[target]; !ok {
		return "", errors.New("reference doesn't exist")
	}
	if _, ok := network.nodes[target]; !ok {
		return "", errors.New("reference ID doesn't exist")
	}
	return network.nodes[target].hostID, nil
}

func (network *Nodenetwork) createNode(hostAddress, domainID string) (*Node, error) {
	key := id.GetRandomKey()
	stringID := string(key[:bytes.IndexByte(key, 0)])
	node := NewNode(stringID, domainID)
	return node, nil
}

func (network *Nodenetwork) addNode(node *Node) error {
	if node == nil {
		return errors.New("node is nil")
	}
	network.nodes[node.GetDomainID()] = node
	return nil
}

func getReferenceString(ref []byte) string {
	if ref == nil {
		return ""
	}
	return string(ref)
}
