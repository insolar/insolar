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

	"github.com/insolar/insolar/application/appfoundation"
	"github.com/insolar/insolar/application/builtin/proxy/noderecord"
	"github.com/insolar/insolar/logicrunner/builtin/foundation"
)

// NodeDomain holds node records.
type NodeDomain struct {
	foundation.BaseContract

	NodeIndexPublicKey foundation.StableMap
}

// NewNodeDomain create new NodeDomain.
func NewNodeDomain() (*NodeDomain, error) {
	return &NodeDomain{
		NodeIndexPublicKey: make(foundation.StableMap),
	}, nil
}

// RegisterNode registers node in system.
func (nd *NodeDomain) RegisterNode(publicKey string, role string) (string, error) {
	root := appfoundation.GetRootMember()
	if *nd.GetContext().Caller != root {
		return "", fmt.Errorf("only root member can register node")
	}

	newNode := noderecord.NewNodeRecord(publicKey, role)
	node, err := newNode.AsChild(nd.GetReference())
	if err != nil {
		return "", fmt.Errorf("failed to save as child: %s", err.Error())
	}

	newNodeRef := node.GetReference().String()
	nd.NodeIndexPublicKey[publicKey] = newNodeRef

	return newNodeRef, err
}

// GetNodeRefByPublicKey returns node reference.
// ins:immutable
func (nd *NodeDomain) GetNodeRefByPublicKey(publicKey string) (string, error) {
	nodeRef, ok := nd.NodeIndexPublicKey[publicKey]
	if !ok {
		return "", fmt.Errorf("network node not found by public key: %s", publicKey)
	}
	return nodeRef, nil
}
