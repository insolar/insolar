//
// Copyright 2019 Insolar Technologies GmbH
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package nodedomain

import (
	"fmt"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/logicrunner/builtin/proxy/noderecord"
	"github.com/insolar/insolar/logicrunner/builtin/proxy/rootdomain"
	"github.com/insolar/insolar/logicrunner/goplugin/foundation"
)

// NodeDomain holds noderecords
type NodeDomain struct {
	foundation.BaseContract

	NodeIndexPK map[string]string
}

// NewNodeDomain create new NodeDomain
func NewNodeDomain() (*NodeDomain, error) {
	return &NodeDomain{
		NodeIndexPK: make(map[string]string),
	}, nil
}

func (nd *NodeDomain) getNodeRecord(ref insolar.Reference) *noderecord.NodeRecord {
	return noderecord.GetObject(ref)
}

// RegisterNode registers node in system
func (nd *NodeDomain) RegisterNode(publicKey string, role string) (string, error) {

	root, err := rootdomain.GetObject(*nd.GetContext().Parent).GetRootMemberRef()
	if err != nil {
		return "", fmt.Errorf("[ RegisterNode ] Couldn't get root member reference: %s", err.Error())
	}
	if *nd.GetContext().Caller != *root {
		return "", fmt.Errorf("[ RegisterNode ] Only Root member can register node")
	}

	newNode := noderecord.NewNodeRecord(publicKey, role)
	node, err := newNode.AsChild(nd.GetReference())
	if err != nil {
		return "", fmt.Errorf("[ RegisterNode ] Can't save as child: %s", err.Error())
	}

	newNodeRef := node.GetReference().String()
	nd.NodeIndexPK[publicKey] = newNodeRef

	return newNodeRef, err
}

// GetNodeRefByPK returns node ref
func (nd *NodeDomain) GetNodeRefByPK(publicKey string) (string, error) {
	nodeRef, ok := nd.NodeIndexPK[publicKey]
	if !ok {
		return nodeRef, fmt.Errorf("[ GetNodeRefByPK ] NetworkNode not found by PK: %s", publicKey)
	}
	return nodeRef, nil
}

// RemoveNode deletes node from registry
func (nd *NodeDomain) RemoveNode(nodeRef insolar.Reference) error {
	node := nd.getNodeRecord(nodeRef)
	nodePK, err := node.GetPublicKey()
	if err != nil {
		return fmt.Errorf("[ RemoveNode ] NetworkNode not found by PK: %s", nodePK)
	}

	delete(nd.NodeIndexPK, nodePK)
	return node.Destroy()
}
