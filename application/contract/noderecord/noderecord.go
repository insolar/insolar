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

package noderecord

import (
	"errors"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/logicrunner/goplugin/foundation"
)

// NodeRecord contains info about node
type NodeRecord struct {
	foundation.BaseContract

	PublicKey string
	Role      core.NodeRole
}

// New creates new NodeRecord
func NewNodeRecord(pk string, roleS string) (*NodeRecord, error) {

	role := core.GetRoleFromString(roleS)
	if role == core.RoleUnknown {
		return nil, errors.New("Can't unsupported role")
	}

	return &NodeRecord{
		PublicKey: pk,
		Role:      role,
	}, nil
}

// GetPublicKey returns public key
func (nr *NodeRecord) GetPublicKey() (string, error) {
	return nr.PublicKey, nil
}

// GetRole returns role
func (nr *NodeRecord) GetRole() (core.NodeRole, error) {
	return nr.Role, nil
}

// GetRoleAndPublicKey returns role-pubKey pair
func (nr *NodeRecord) GetRoleAndPublicKey() (core.NodeRole, string, error) {
	return nr.Role, nr.PublicKey, nil
}

// SelfDestroy makes request to destroy current node record
func (nr *NodeRecord) Destroy() error {
	nr.SelfDestruct()
	return nil
}
