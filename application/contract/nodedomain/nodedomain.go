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
	"errors"

	"github.com/insolar/insolar/application/proxy/noderecord"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/cryptohelpers/ecdsa"
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

// RegisterNode registers node in system
func (nd *NodeDomain) RegisterNode(pk string, role string) (core.RecordRef, error) {
	// TODO: what should be done when record already exists?
	newRecord := noderecord.NewNodeRecord(pk, role)
	record := newRecord.AsChild(nd.GetReference())
	return record.GetReference(), nil
}

func (nd *NodeDomain) getNodeRecord(ref core.RecordRef) *noderecord.NodeRecord {
	return noderecord.GetObject(ref)
}

// RemoveNode deletes node from registry
func (nd *NodeDomain) RemoveNode(nodeRef core.RecordRef) error {
	node := nd.getNodeRecord(nodeRef)
	return node.Destroy()
}

// IsAuthorized checks is signature correct
func (nd *NodeDomain) IsAuthorized(nodeRef core.RecordRef, seed []byte, signatureRaw []byte) (bool, error) {
	pubKey, err := nd.getNodeRecord(nodeRef).GetPublicKey()
	if err != nil {
		return false, err
	}
	ok, err := ecdsa.Verify(seed, signatureRaw, pubKey)
	if err != nil {
		return false, err
	}
	return ok, nil
}

// Authorize checks node and returns node info
func (nd *NodeDomain) Authorize(nodeRef core.RecordRef, seed []byte, signatureRaw []byte) (pubKey string, role core.NodeRole, err error) {
	// TODO: this should be removed when proxies stop panic
	defer func() {
		if r := recover(); r != nil {
			pubKey = ""
			role = core.RoleUnknown
			err, ok := r.(error)
			errTxt := ""
			if ok {
				errTxt = err.Error()
			}
			err = errors.New("[ Authorize ] Recover after panic: " + errTxt)
		}
	}()
	nodeR := nd.getNodeRecord(nodeRef)
	role, pubKey, err = nodeR.GetRoleAndPublicKey()
	if err != nil {
		return "", core.RoleUnknown, errors.New("[ Authorize ] Problem with getting role and key: " + err.Error())
	}
	ok, err := ecdsa.Verify(seed, signatureRaw, pubKey)
	if err != nil {
		return "", core.RoleUnknown, errors.New("[ Authorize ] Problem with verifying of signature: " + err.Error())
	}
	if !ok {
		return "", core.RoleUnknown, errors.New("[ Authorize ] Can't verify signature: " + err.Error())
	}

	return pubKey, role, nil
}
