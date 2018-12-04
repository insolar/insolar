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
	"fmt"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/logicrunner/goplugin/foundation"
)

// RecordInfo holds record info
type RecordInfo struct {
	PublicKey string
	Role      core.StaticRole
}

// NodeRecord contains info about node
type NodeRecord struct {
	foundation.BaseContract

	Record RecordInfo
}

// NewNodeRecord creates new NodeRecord
func NewNodeRecord(publicKey string, roleStr string) (*NodeRecord, error) {
	if len(publicKey) == 0 {
		return nil, fmt.Errorf("[ NewNodeRecord ] public key is required")
	}
	if len(roleStr) == 0 {
		return nil, fmt.Errorf("[ NewNodeRecord ] role is required")
	}

	role := core.GetStaticRoleFromString(roleStr)
	if role == core.StaticRoleUnknown {
		return nil, fmt.Errorf("Role is not supported: %s", roleStr)
	}

	return &NodeRecord{
		Record: RecordInfo{
			PublicKey: publicKey,
			Role:      role,
		},
	}, nil
}

var INSATTR_GetNodeInfo_API = true

// GetNodeInfo returns RecordInfo
func (nr *NodeRecord) GetNodeInfo() (RecordInfo, error) {
	return nr.Record, nil
}

var INSATTR_GetPublicKey_API = true

// GetPublicKey returns public key
func (nr *NodeRecord) GetPublicKey() (string, error) {
	return nr.Record.PublicKey, nil
}

// GetRole returns role
func (nr *NodeRecord) GetRole() (core.StaticRole, error) {
	return nr.Record.Role, nil
}

// Destroy makes request to destroy current node record
func (nr *NodeRecord) Destroy() error {
	nr.SelfDestruct()
	return nil
}
