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

package nodedomain

import (
	"fmt"

	"github.com/insolar/insolar/application/proxy/noderecord"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/logicrunner/goplugin/foundation"
)

// NodeDomain holds noderecords
type NodeDomain struct {
	foundation.BaseContract
}

// NewNodeDomain create new NodeDomain
func NewNodeDomain() (*NodeDomain, error) {
	return &NodeDomain{}, nil
}

func (nd *NodeDomain) getNodeRecord(ref core.RecordRef) *noderecord.NodeRecord {
	return noderecord.GetObject(ref)
}

var INSATTR_RegisterNode_API = true

// RegisterNode registers node in system
func (nd *NodeDomain) RegisterNode(publicKey string, role string) (string, error) {

	newNode := noderecord.NewNodeRecord(publicKey, role)
	node, err := newNode.AsChild(nd.GetReference())
	if err != nil {
		return "", fmt.Errorf("[ RegisterNode ] Can't save as child: %s", err.Error())
	}

	return node.Reference.String(), err
}

// RemoveNode deletes node from registry
func (nd *NodeDomain) RemoveNode(nodeRef core.RecordRef) error {
	node := nd.getNodeRecord(nodeRef)
	return node.Destroy()
}

// IsAuthorized checks is signature correct
/*func (nd *NodeDomain) IsAuthorized(nodeRef core.RecordRef, seed []byte, signatureRaw []byte) (bool, error) {
	pubKey, err := nd.getNodeRecord(nodeRef).GetPublicKey()
	if err != nil {
		return false, fmt.Errorf("[ IsAuthorized ] Can't get nodes: %s", err.Error())
	}

	publicKey, err := foundation.ImportPublicKey(pubKey)
	if err != nil {
		return false, fmt.Errorf("[ verifySig ] Invalid public key")
	}

	ok := foundation.Verify(seed, signatureRaw, publicKey)
	return ok, nil
}*/

//var INSATTR_Authorize_API = true

// Authorize checks node and returns node info
/*func (nd *NodeDomain) Authorize(nodeRef core.RecordRef, seed []byte, signatureRaw []byte) (string, core.NodeRole, error) {
	nodeR := nd.getNodeRecord(nodeRef)
	nodeInfo, err := nodeR.GetNodeInfo()
	if err != nil {
		return "", core.RoleUnknown, fmt.Errorf("[ Authorize ] Problem with Getting info: %s", err.Error())
	}

	pubKey := nodeInfo.PublicKey
	role := nodeInfo.Role

	publicKey, err := foundation.ImportPublicKey(pubKey)
	if err != nil {
		return "", core.RoleUnknown, fmt.Errorf("[ verifySig ] Invalid public key")
	}

	ok := foundation.Verify(seed, signatureRaw, publicKey)
	if !ok {
		return "", core.RoleUnknown, fmt.Errorf("[ Authorize ] Can't verify signature")
	}

	return pubKey, role, nil
}*/
